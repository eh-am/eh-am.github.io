+++ 
date = 2023-01-19T11:32:25Z
title = "DeepCopy in Javascript"
categories = ["technical"]
+++

Sometimes you need to perform a deep copy in Javascript.
Spread obviously doesn't work, since it does a shallow copy.
I've used `JSON.parse(JSON.stringify)`, but that doesn't work with `Date` or other
data structures, like `Map` and `Set` and others.

But recently [I've came across](structuredClone) the `structuredClone` which does what's expected.
And the [support seems okay](https://caniuse.com/?search=structuredClone), although
`core-js` has a [polyfill](https://github.com/zloirock/core-js#structuredclone).

For the full list of supported types, see the [MDN docs](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Structured_clone_algorithm#supported_types).
