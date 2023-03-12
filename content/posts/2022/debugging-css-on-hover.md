+++
date = "2022-09-06"
title = "debugging a trigged on hover tooltip style"
categories = ["technical"]
tags = ["css", "hover", "tooltip"]
+++

At $WORK we have a tooltip component triggered on hover via javascript:

![Showing tooltip behaviour](./tooltip.gif)

Normal stuff.

But I needed to debug its styling. I know how to do that when it's a css state, just [use dev tools and set up an element state to `hover`](https://www.smashingmagazine.com/2021/02/useful-chrome-firefox-devtools-tips-shortcuts/#switching-between-dock-states-chrome-edge-firefox).

Here, however, state is triggered via javascript, so it won't work.

The solution, simply enough, is to go to the `Debugger` section, do the action (here hovering), then press F8, which is shortcut for pausing the execution.

Then it's just a matter of using the inspector to play with the styles:

![State when debugger is turned on](./debugger.png)
