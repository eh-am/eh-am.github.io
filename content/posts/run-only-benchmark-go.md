+++
date = "2021-07-22"
title = "Run only benchmark in go(lang)"
slug = "run-only-benchmark-go"
categories = ["blog", "go", "golang", "technical"]
+++

a quick go(lang) tip

let's say you have tests and benchmarks in your codebase
and you want to **run only benchmarks**, or to put in other words, to **filter out all tests** from being executed

unfortunately there's no `--bench-only` nor anything similar available,
but there's a trick popular in the community that does the job

the idea is to basically filter out all tests with an unlikely regex

```
go test ./... -bench=. \
  -run=dwioajd2qjdowajdq2jodfwioaaaadlzndwadjwionmilafmaeopdwpaod \
  pakdwioadjwioajdwaiojdwaiowdjaiodwajwdaiojdaiowdjaiodwajiodwajiodwajwd \
  aiwdoajdwioajwdioadwa
```
![cat bashing the keyboard](https://media.giphy.com/media/unQ3IJU2RG7DO/giphy.gif)

but of course we can do better with a clever regex

```
go test ./... -bench=. -run=^$
```

as you may be aware, `^` matches the beginning of a string and `$` the end
so it tries to match a string that's empty

naturally there's no test with empty name, otherwise the parser would fail, i think

therefore it ends up filtering out all tests, so only benchmarks run


