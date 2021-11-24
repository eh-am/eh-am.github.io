+++
date = "2021-11-24"
title = "Continuously check code coverage while developing tests in go"
categories = ["blog", "go", "technical"]
+++

I have this workflow where while developing tests I also check for code coverage. The idea is to see if I am really touching that piece of code I am adding tests for. Similar to TDD, there's a certain primal endorphine rush from seeing the coverage go from red to green.


In a nutshell, it works like this:

```shell
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

The first bit will create a `coverage.out` file, which we will then process using `go tool cover` and generate a HTML.

It also calls, on Linux, [`xdg-open`](https://github.com/golang/go/blob/2ebe77a2fda1ee9ff6fd9a3e08933ad1ebaea039/src/cmd/internal/browser/browser.go#L29) which opens up the file in your (supposedly) preferred browser.


However, by default it creates a temporary HTML file (on Linux, under `/tmp`). For example, `file:///tmp/cover617771716/coverage.html#file0`. So every time you run the command again, it will generate another file, which will open in another tab, making the experience quite annoying.

A simple solution is to, instead of letting `cover` generate a random file, to output to a fixed location file, like `coverage.html`, which you can then refresh manually in your browser.

The whole thing then becomes:
```shell
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

# you only have to run this once
xdg-open coverage.html
```


Theoretically one can add a server that watches that file and refreshes the html.

I may do that sometime.
