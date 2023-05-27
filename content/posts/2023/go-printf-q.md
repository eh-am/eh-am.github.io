+++ 
date = 2023-05-27T17:21:53+01:00
title = "Formatting strings with %q in go"
tags = ["go"]
categories = ["technical"]
+++

TIL you can use `%d` in `Printf` statements (and its friends, like `Fatalf/Errorf`) when printing strings/runes.

What it does is:
* Escapes any special characters
* Add quotes around it (single if it's a single rune)

I really like how it adds quotes around the string. I've had innumerous times
when a go program was supposed to printed something (mostly likely in errors),
but since the string was empty, it didn't print anything.

For example,

```
Error. Invalid URL:
```

Which always takes me a while to notice. By having an error like
```
Error. Invalid URL: ""
```

It makes it clear how somehow an empty string is being passed.

Beware that his only works for strings (and runes) though.

Example can be found in the [go playground](https://go.dev/play/p/YgxdWAS4-yh).
