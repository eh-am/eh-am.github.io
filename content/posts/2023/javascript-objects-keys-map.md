+++ 
date = 2023-03-23T10:18:14Z
title = "Debug session: Javascript, toString, object, maps e outras porcarias"
description = ""
slug = "" 
tags = []
categories = []
externalLink = ""
series = []
+++

Fun debugging session that happened at $WORK.

I found that a component, published to npm was crashingwhen giving a certain input.

Since it was minified, it was hard to see where the error was. I added the sourcemap. It happened in a loop, a certain property which was expected to be an array, ended up being a `Number`:

```javascript
existingNode.total[0] += currentNode.total[0];
```


Firefox error message was kinda confusing:
```
Uncaught TypeError: can't assign to property 0 on 288: not an object
```

Chrome's was better:
```
TypeError: Cannot create property '0' on number '144'
```

That's the thing with TypeScript: **types are not real**.
![An allusion to the "is x in the room with us right now" meme](./types.jpg)

Anyway, after some debugging I found that it happened when the string `toString` is used. I handcrafted a [minimum reproducible example](https://en.wikipedia.org/wiki/Minimal_reproducible_example), which helped me debug. Debuggers are cool, but Conditional Breakpoints often let me down.

To give some context, the code iterates a tree and merges nodes' values with the same name. It's a feature for [sandwich viewing](https://jamie-wong.com/post/speedscope/) if you are interested.

What I found curious is that when I logged the node, it told me it was a function with properties!

```js
function toString()
length: 0
name: "toString"
self: 8
total: 72
type: "single"
```

This is something I kinda forgot about, but why wouldn't be possible?

```js
function myFn() {}
myFn.foo = 'bar';

> myFn().foo
'bar'
```

Back in the day we used to write constructors using functions:
```js
function Person() { this.name = 'John'; }

> new Person().name
'John'
```

And of course, ES6 classes have some syntax sugar over this implementation ([but not only that!](https://stackoverflow.com/a/54861781)).


Anyway, what was happening is that we created a map-like object. And then to not initialize twice, we did a check for the presence of an existing item:

```js
childrenMap[node.children[i].name] ||= node.children[i];
```

However, since `toString` always exist in an object, it wasn't being initialized correctly!

Then it would spiral in some crazy madness I won't go deep about.

The solution was to just replace the map-like object for a real [Map](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map):

> A Map does not contain any keys by default. It only contains what is explicitly put into it.

> An Object has a prototype, so it contains default keys that could collide with your own keys if you're not careful.



I didn't explain where the properties of `toString` were coming form, to figure out I used a handy Proxy:
```js
Object.prototype.toString = new Proxy(Object.prototype.toString, {
  set: () => {
    debugger;
  },
});
```

Whenever `toString` is `set` (ie overridden), the `debugger` is called.

With that I managed to find the offending code:
```js
hash[name] = hash[name] || { name: name || '<empty>' };

// latter on hash[name] becomes c
c.type = 'single';
c.self = zero(c.self) + ff.getBarSelf(level, j);
c.total = zero(c.total) + ff.getBarTotal(level, j);
```
So if `name` is for example, `toString`, `hash[name]` returns the function from the prototype chain, ie from `Object.prototype.toString`, which attributes are added, leaking to every single object.
