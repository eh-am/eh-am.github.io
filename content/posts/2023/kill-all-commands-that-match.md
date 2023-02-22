+++ 
date = 2023-02-22T16:41:46Z
title = "Kill all commands that match a query"
categories = ["technical"]
+++

I often kill a bunch of processes with some command line incantation, like `ps aux | grep $MY_QUERY | awk '{ print $2 }' | xargs kill`,
which as you can see, it's pretty verbose.

I am ashamed of myself, but only recently I learned `pkill` has a `-f` flag that kills all commands 
(ie. the whole thing, not only the process name) that match a certain query.

We can even pair it with `-i` to make it case insensitive, and `-l` to print the PID of the affected process.

```bash
pkill -fi $MY_QUERY
```
