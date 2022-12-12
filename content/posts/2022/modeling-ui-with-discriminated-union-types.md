
+++
date = "2022-12-12"
title = "Modeling UI with Discriminated Union Types"
categories = ["blog", "technical", "typescript"]
+++


I like to model my app's state using discriminated union types:

```ts
type MyState<T> = {
  type: 'pristine',
} | {
  type: 'loading',
} | {
  type: 'loaded',
  data: T,
} | {
  type: 'failed',
}
```

Reason is that it forces me to think about all states in the UI.

`pristine` means it's never been touched. In most places this is the initial state which will change in a split second, normally upon first mounting the component which then starts a data load. But sometimes it's different, for example, a Dropdown that lazy loads data only when clicked.


`loading` means data is being loaded (duh). Don't assume data will be loaded instantly, and instead, code the UI such that it has meaningful visuals. Loading spinner, a skeleton screen etc.

`loaded` means we now have data. This is the "good path".

`failed` means that data fetching failed. Honestly I don't always add it, since sometimes it's dealt with at a different layer. But you could show a message in the component body, a retry button, or just an indication that the component is incomplete.


Sometimes I also add a `reloading` state:
```ts
{ type: 'reloading', data: T }
```

The difference from `loading` is that this time, we have data. Useful for stuff like graphs and dashboards, where you don't want to lose previous data while you are `reloading`.


Then I like to use `switch case` with exhaustiveness checks:
```tsx
switch (s.type) {
  case 'pristine': {
    return <></>
  }
  case 'loading': {
    return <div>Loading</div>
  }
  case 'failed': {
    return <div>Huhh something went wrong</div>
  }
  case 'loaded': {
    const data = s.data;
    return <MyHappyPathComponent data={data} />;
  }

  default: {
    const exhaustiveCheck: never = s;
    return exhaustiveCheck;
  }
}
```

The trick is the use of `never`. It forces us to deal with any additional `type` when they are added (eg a new `reloading` type). This is a pretty useful technique when dealing with unions (eg `A | B` then you add `A | B | C`).

Source: https://basarat.gitbook.io/typescript/type-system/discriminated-unions#exhaustive-checks


We can put that into a component
```ts
function MyComponent(s: MyState) {
  switch (s.type) {
    case 'pristine': {
      return <></>
    }
    case 'loading': {
      return <div>Loading</div>
    }
    case 'failed': {
      return <div>Huhh something went wrong</div>
    }
    case 'loaded': {
      const data = s.data;
      return <MyHappyPathComponent data={data} />;
    }

    default: {
      const exhaustiveCheck: never = s;
      return exhaustiveCheck;
    }
  }
}
```

Then for `MyHappyPathComponent`, how can one type its props?
We can't simply access like `MyState["data"]`, since that field only exists for the `type: "loaded"` variant.

The solution here is to use Typescript's `Extract`

```ts
type MyHappyComponentProps = {
  data: Extract<MyState, { type: "loaded" }>['data']
}
```

The way it works is as follows: `Extract` only the types that are assignable to `{type: "loaded"}`. In this case, only loaded is, which we then get only its `data`.

For an example (with required adaptations) see it in the [typescript playground](https://www.typescriptlang.org/play?#code/C4TwDgpgBAsiDKwCGwIB4AqA+KBeKA3gFBRSiQBcUA5GAE4CWAzsAwHYTUA0RAvlAB9CJMuAhVqAGwD2SACbsA5tz6DhpcuJoz5EOStJyUSKhh78hxDWIkAzJA0l6V-IgHo3UAKpsAxtIBbAIg2YDJpKCYIaGAAC2gIAA9YpABXFgYAN2hfeN8Aayh7RyV3T0sy0mtKGjoIHQU2ZR4PKqgjZFNK3iIifzYWSORUKjhEFHQCKDYkYKoWRibeHHwrKs0JemZWDmo+AG5e21S-Vmk2WBAAYUCwc5DgAAomKk1pWyGJgEp1KqYAdwYwFyUGeADpND81m1SL4kFEaFsMrsqNCYW06sBUnQLtRqId0aQeoSoHCEVJZI1lKiRCTSJjsbiGkowaz8bSYcTCWToNRik59DS6VUGTiaAAxBwC9kkrnonnaSnOIXCqCii5wAASSDAYBAAAUULEbgE7hxQuCOkgvgTCT0OW05BB7KlJMAVcL+oMkil0qxslc8vkqBxsnQ8JFbcL1VAfWkMgGg1HORyevbjqcGOdLtrdQajSazQ9nlQAKKJYB0JC+YBoN4fFgTLiEUQ1Cm6fTLADa1Ct1AAulCOTHqCk81AwEbqFAANSRMEzYKHXhAA).

# Downsides
As with everything in life, there's also downsides:

* There's quite a bit of boilerplate, although it's only because we are dealing with cases we would otherwise deal with implicitly or not at all.
* It gets coupled with the state. So if you use redux, with an approach like this you end up not using selectors very often, only to get the state itself. All the further manipulations are done in the component. For example, `selectMyData(s: State) => Data` isn't something you would write anymore.

