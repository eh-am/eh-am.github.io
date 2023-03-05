+++ 
date = 2023-03-05T18:20:55Z
title = "Optimizing a deploy to go from ~8 minutes to ~14 seconds"
tags = ["hugo", "profiling"]
categories = ["technical"]
+++

As you may know, I have an alternative rock podcast called [triceratops show](www.triceratops.show).
Its stack is pretty simple: a static website using hugo. Posts are created via Netlify CMS, which triggers a pipeline that builds the app and deploys to s3.

I noticed deploys were getting pretty slow. 5, 6, 7 minutes and sometimes even 8. It didn't make sense, since building locally is pretty fast, and pushing to s3 should be fast, since most files won't change.

So I decided to investigate.

First step was to run `hugo deploy` with `--dryRun` and a bunch of debug/log flags (why so many I wonder?):

```bash
hugo deploy --log --verbose --verboseLog --debug --maxDeletes 0 --dryRun
```

As you can imagine, logs were pretty verbose, but I managed to notice that
certain XML files didn't need to be uploaded, but their HTML versions did

```
DEBUG 2023/03/04 15:02:56 artistas/melvins/index.html needs to be uploaded: size differs
DEBUG 2023/03/04 15:02:56 artistas/melvins/index.xml exists at target and does not need to be uploaded
```

However, the file itself is pretty small (~2kB), so that shouldn't account for all the slowness
```
[DRY RUN] Would upload: artistas/melvins/index.html (2.0 kB, Content-Encoding: "gzip", Content-Type: "text/html"): size differs
```

I decided to profile hugo to see if the CPU was doing anything unexpected.
So I cloned hugo, installed the pyroscope agent and told to upload to Pyroscope Cloud. I just put to run anywhere in the `commands/deploy.go` file:

```go
pyroscope.Start(pyroscope.Config{
  ApplicationName: "hugo.debug",
  ServerAddress:   "https://ingest.pyroscope.cloud",
  AuthToken:       "MY_AUTH_KEY",
})
```

Pyroscope CPU profile showed me that most of the work was being done in the `deploy.walkRemote` function, mostly reading (`s3blob.Read`) and calculating the hash (`cripto/md5.Write`). What for?

<iframe frameBorder="0" width="100%" height="400" src="https://flamegraph.com/share/dcb094dc-baae-11ed-9a53-164a641be00c/iframe?onlyDisplay=flamegraph&colorMode=light"></iframe>

I took that as a clue and started looking into hugo's sourcecode.

I saw that it first checked the file's MD5 hash, and if doesn't exist, it reads the file and compute the hash.

https://github.com/gohugoio/hugo/blob/32ea40aa82736d33aebac246c2552d1342988d9d/deploy/deploy.go#L573-L603

The hash may come from the blob storage itself, or via a `metaMd5Hash` metadata field, set on upload.

https://github.com/gohugoio/hugo/blob/32ea40aa82736d33aebac246c2552d1342988d9d/deploy/deploy.go#L320


Assuming the worse, I sprinkled some log messages and figured that mp3 files were the ones being downloaded and their hash computed.

Why does hugo need to know about these mp3 files? Well, it checks the files that are local and in the remote, so that it can figure out if any files need to be deleted. If it only taken the local files into account, a file that was previously deployed but doesn't exist anymore would have been kept!

Anyway, the thing is, I don't need that to happen at all, since mp3 files are not in git, so they are uploaded via a completely different workflow. I updated my hugo config file to simply ignore mp3 (and mp4) files:

```yaml
deployment:
  targets:
    (...)

    # exclude mp3 since they are uploaded manually
    exclude: "**.{mp3,mp4}"
```

The last thing was figuring out why html files didn't match:
```
[DRY RUN] Would upload: episodios/17/index.html (14 kB, Content-Type: "text/html"): size differs
```

I first thought it was because my files were being built without minified, but uploaded a minified version. Then I thought it was due to gzip (maybe different versions) computing different sizes.

Upon inspecting the html from both the deployed version and locally, I realized I was baking the `GIT_SHA` in each file. Of course it would change with every commit!

**So I removed that, pushed, and now my deploys only take ~14 seconds!**

The final profile still looks pretty similar, except for the time spent, so I won't bore you.

There are still one question, like why doesn't S3 return the MD5 for my mp3 files? In their docs they have
>  The ETag might or might not be an MD5 digest of the object data. Whether or not it is depends on how the object was created and how it is encrypted, as follows:
> * (...)
> * Objects created by either the Multipart Upload or Upload Part Copy operation have ETags that are not MD5 digests, regardless of the method of encryption.

Source:
https://docs.aws.amazon.com/AmazonS3/latest/API/RESTCommonResponseHeaders.html

So since I upload via the AWS Console, maybe it uses multi-part upload?

In any case it doesn't change my solution, none of this work even needs to be done, since mp3 files are part of a completely different workflow.

