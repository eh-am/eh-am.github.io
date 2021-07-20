+++
date = "2021-06-20"
title = "search under cursor in vim"
slug = "search-under-cursor-vim"
categories = ["blog", "vim", "tech"]
+++


here's a vim tip i've learned recently

let's say you want to search other occurences of a keyword
in that case, just press `*`

and then `n` to move forward and `N` (that's capital N) to move backwards


keep in mind [vim's definition of keyword](http://vimdoc.sourceforge.net/htmldoc/options.html#'isfname'), ie **spider-man** is considered 2 words

in that (and other more advanced) cases, you may want to use visual selection

unfortunately it's not as straightforward as we would want, but here's how it is done:

* press `v` to enter visual mode,
* select the desired text
* `y` to copy to register `"`
* `/` to enter search
* `Ctrl + r` + `"` to paste from register `"`

keep in mind the query search is considered a string

so let's say you have a file with
```
spider-.*
spider-pig
spider-ham
spider-man
```
searching for `spider-.*` will naturally highlight all lines

![Spider-Ham giving thumbs up](./spider-ham.png)

# source
* `:help *`
* `:help iskeyword`
* [How to search for selected text in Vim? - Super User](https://superuser.com/questions/41378/how-to-search-for-selected-text-in-vim)

