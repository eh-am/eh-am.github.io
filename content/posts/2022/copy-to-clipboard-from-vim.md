
+++
date = "2022-10-25"
title = "Copy to clipboard from vim using visual selection"
tags = ["vim"]
categories = ["technical"]
+++

I don't remember the first time I found this tip, but the other day I was trying to find the exact syntax and couldn't, hence this post.

# Option 1
To copy to clipboard from vim using visual selection, select the thing you want (`v` or `shift + v`), then press `:`, it will open the command mode box (forgot what's called) with `:'<,'> ` prepopulated. Then just append `:w !pbcopy`

ie.

```vim
:'<,'>:w !pbcopy
```


The `:w` bit writes to a file.

The `!pbcopy` bit (not quite sure) means that instead of writing to a file `pbcopy`, to evaluate it as a command.


Since I've been using a mac, `pbcopy` is already available. For linux you can alias it to `xclip -selection clipboard -in`.


The downside is that it copies the entire line, which is kinda annoying.

# Option 2 (better)
On Mac and Windows we can write to the `*` register which is the system one.

So select as normal, then press `"*y`.

I think for linux it's the `+` register.
