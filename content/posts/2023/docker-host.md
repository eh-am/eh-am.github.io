+++ 
date = 2023-02-19T11:04:12Z
title = "Acessing the host from a docker container"
tags = ["docker"]
categories = ["technical"]
+++

This is one of the things I wish I knew earlier.

When developing locally, sometimes you need a database, in those cases you can simply
run `docker`/`docker-compose`, expose the port and hit that port from the development server.

When running tests, I like to use `testcontainers` to spin up a test db.

So `host` -> `container` communication is straightforward.


However, there are cases when I need to access the `host` from a container:
* Running `prometheus` in a container, and pulling `/metrics` from the `host`;
* Running some benchmark tool against a local development server;

I knew that when running `Docker for Mac` (and I think `Docker for Windows`), the `host.docker.internal` DNS name
is resolved to the host. But that didn't work for Linux, which requires passing `--add-host=host.docker.internal:host-gateway` for it to work.

The point is that there would be 2 different setups depending on the OS, which doesn't scale
very well in a team with non-homogeneous OSes. Boo!

But recently I found about `qoomon/docker-host`, which cleverly makes the setup transparent.

[It works by checking if the `host.docker.internal` resolves (Mac), otherwise falls back to the
default gateway (for Linux)](https://github.com/qoomon/docker-host/blob/5b9cfbc9d2410bf65accea42f25625ec9eb95ff8/entrypoint.sh#L38-L55).

Then it uses [iptables rules](https://github.com/qoomon/docker-host/blob/5b9cfbc9d2410bf65accea42f25625ec9eb95ff8/entrypoint.sh#L78-L85) to forward any traffic received on that container to the same port, in the host.


For illustration, run an HTTP server on port 8000 using `python`:

```bash
python3 -m http.server
```

Then, in a `docker-compose.yaml` file, we `curl` that server via `http://docker-host:8000`:

```yaml
version: '3.9'
services:
  docker-host:
    image: "qoomon/docker-host"
    cap_add:
      - "NET_ADMIN"
      - "NET_RAW"
  curl:
    image: alpine/curl
    command:
    - sh
    - -c
    - 'while true; do curl -s http://docker-host:8000; sleep 10; done'
```

Which returns
```
curl_1         | <!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
curl_1         | <html>
curl_1         | <head>
curl_1         | <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
curl_1         | <title>Directory listing for /</title>
curl_1         | </head>
curl_1         | <body>
curl_1         | <h1>Directory listing for /</h1>
curl_1         | <hr>
curl_1         | <ul>
curl_1         | <li><a href="hello-world">hello-world</a></li>
curl_1         | </ul>
curl_1         | <hr>
curl_1         | </body>
curl_1         | </html>
curl_1         | <!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
curl_1         | <html>
curl_1         | <head>
curl_1         | <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
curl_1         | <title>Directory listing for /</title>
curl_1         | </head>
curl_1         | <body>
curl_1         | <h1>Directory listing for /</h1>
curl_1         | <hr>
curl_1         | <ul>
curl_1         | <li><a href="hello-world">hello-world</a></li>
curl_1         | </ul>
curl_1         | <hr>
curl_1         | </body>
curl_1         | </html>
```

Showing that it is able to hit the host!

# Quirks
1. If you are running rootless containers, [you need to manually set the `DOCKER_HOST` env var with the host's machine `HOSTNAME`](https://github.com/qoomon/docker-host/issues/49#issuecomment-1353064725)


2. Of course, the `NET_ADMIN` and `NET_RAW` capabilities need to be added, since you are messing up with the network.
