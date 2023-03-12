+++
date = "2022-07-21"
title = "yarn and 'Extracting tar content of undefined failed, the file appears to be corrupt'"
categories = ["technical"]
tags = ["yarn", "npm", "node"]
+++

Here's a funny one.

One day, at $WORK, CI started to fail when running `yarn install`:

 ```
 error https://registry.yarnpkg.com/ol/-/ol-6.12.0.tgz: 
 Extracting tar content of undefined failed, the file appears to be corrupt: 
 "ENOENT: no such file or directory, 
 chmod '/home/runner/.cache/yarn/v6/npm-ol-6.12.0-0de51abad0aaeb0eca41cba3c6d26ee485a3e92b-integrity/node_modules/ol/format/KML.js'"
 ```

After some digging we found that's due to a race condition that made yarn pull the same dependency in parallel.

The workaround? `yarn --network-concurrency=1`

And that was fine for a while.

Although it was slow. So I decided to investigate it a bit further.

[Somebody found out that it happens when you import a git repository, and that repository contains a `prepare` script](https://github.com/yarnpkg/yarn/issues/6312#issuecomment-422806004).


Luckily (or by design), there was only a repo dependency, that we control.

Removing the `prepare` upstream solved the issue.
