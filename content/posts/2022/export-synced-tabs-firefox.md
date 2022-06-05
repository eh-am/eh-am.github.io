+++
date = "2022-06-05"
title = "Export Firefox Synced Tabs"
slug = "export-firefox-synced-tabs"
categories = ["blog", "technical", "firefox"]
+++

I tend to have many tabs open. Like way too many.
I first noticed this behaviour when the [tabs counter in my browser became a smile](https://www.dailymail.co.uk/sciencetech/article-4356434/Google-Chrome-secret-message-incognito-users.html).

But that is fine, everybody knows someone like that, hell, you may even be that person.

What I am here to talk about is a healthy routine I've developed: from time to time (ideally once a month), close all tabs. But hey, firefox has a feature that allos closing old tabs automatically, somebody may say. Yeah, but when possible I prefer to have control over what happens in my life. The world moves too fast! Bring back the 00s! Get off my lawn! \*sips monster energy drink\*


There are also couple other useful things: 

a) Do the export thing via a computer. I don't want to connect my smartphone to my computer.
b) Be able to export these tabs, in the vain hope that I would check them again.

Ok so how do we do it? Easy.

1. Run **firefox**, enable [Firefox Sync](https://www.mozilla.org/en-US/firefox/sync/) and install the [About Sync extension](https://stackoverflow.com/questions/50829318/export-synced-tabs-in-firefox)
2. In your computer, access `about:sync`
3. In the `tabs` section, click `Record Editor (server)` and Select the desired device
4. Copy the json and save somewhere else (preferrably in a directory `$YEAR/$MONTH/$DEVICE.json`)


Just be careful that sometimes devices are not synced, most likely devices that you haven't used in a while. Or you forgot to turn wifi on turning on airplane mode or something.


At some point I may write a script that exports all automatically.

Also it would be nice to be able to close tabs remotely (via Firefox Sync), but I don't see that possible just yet.


# References
[Export Synced Tabs in Firefox - Stack Overflow](https://stackoverflow.com/questions/50829318/export-synced-tabs-in-firefox)
