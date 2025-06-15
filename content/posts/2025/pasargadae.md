+++ 
date = 2025-06-15T17:13:10+01:00
title = "I’m leaving for Pasargada"
description = "Ugly looking design, weird user flows, bad implementations, nitpicks and more"
categories = ["technical"]
+++

Starting a new analysis series of things that I see and bother me.

Disclaimer: there's no shade to any of the people involved in the design and implementation of these subjects. The world is complicated and we can't do our best all the time.

# Analysis
Today we will analyze Medium's comment section.

The first image illustrates the start of the page. All good so far.
![](./1.jpg)

Scrolling down we found the comment section:
![](./2.jpg)

It says there are 22 responses. It allows me to write my own response.

Then there's a list of responses, what are they sorted by? I see they are not by date. Neither by number of "likes"? So is it defined by their own "recommended" metric? I have no idea.

I want to see more responses. Hence I click the button.
![](./3.jpg)


Now it opened a sidebar with more comments. A bit weird to show on the side, but ok. At least it says it's ordered by **Most Relevant**.

Also there are 2 scrollbars, one from the page itself and one from the sidebar. Not a big fan.
![](./4.jpg)

Some of the comments are clearly truncated, with a "more" button.
![](./5.png)

Some others are not truncated, but are inside a box:
![](./6.png)

Some are truncated AND inside a box.
![](./7.png)

Some are just truncated, but wait a minute, these are user written ellipsis!
![](./8.png)

A bit confusing. But you can still hover over a response, and if the cursor changes, it means it's a truncated reply.

So let's see more of that.
![](./9.jpg)

Now we've been taken to a completely different page!

Ok, cool, let's hit the back button and continue seeing the other posts.
![](./1.jpg)

Nope. We are back to an "initial state", the comment section is collapsed and we lost our position.

Anyway, let's see a different post now. Let's click "more" on this one:
![](./5.png)

![](./10.jpg)

It just expanded in place, alright.

Let's click another comment
![](./11.png)

This one just scrolls to the highlighted section from the original article. Not a very clear distinction except that the box contains a green highlight.

Valid to say that the sidebar comment thing uses infinite scroll to load more comments, it's well done, but I think it's something in the way to keep context between navigation. Of course it can be done, but it's more work.

Also replies can be shown in the same sidebar:
![](./12.jpg)

The usual padding way of creating hierarchy. I wonder how many layers they use until migrate to a different solution (probably linking to the "Individual Response" page).

And last but not least, if you click the number of "Claps" in the original post, it opens a modal, which in this case has enough content to require a scrollbar:
![](./13.jpg)

But if you click the **browser** scrollbar, it closes the modal!

So you have to use the mouse wheel/equivalent, which at the end has a "Show more claps" button.
![](./14.jpg)

No infinite scroll here huh?

Also, differently from the "Responses" button, which both icon and number correspond to the same action, the "Claps" have 2 different actions depending on where you click: if it's the icon, it adds a "Clap", if it's the number, it shows the modal previously discussed.

![](./15.png)

# Conclusion
I've briefly went over the "Response" section of a Medium article. Pointing what I found was OK, what was weird and or/inconsistent.

My personal design preferences are to simplify and keep it consistent. But it's known the organizations are big and complex and end up lacking a unified view.
