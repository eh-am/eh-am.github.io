+++
date = "2022-12-22"
title = "open tmux in a certain directory"
categories = ["technical"]
tags = ["tmux"]
+++

When I start working on a project, I start tmux and then `cd` into its directory.
There's a simpler way to do it, though:

```shell
tmux new -s $SESSION_NAME -c $DIRECTORY
```

Which then can be aliased, so I can do something like `tmux-blog` to start in the correct directory.

We can go ever further, and start the section with a specific layout running the commands you want:

```shell
tmux new -s blog -c $BLOG_DIR \; \
  split-window 'make serve' \; \
  detach-client && \
  tmux attach -t blog
```

For horizontal split use `split-window -h`.
