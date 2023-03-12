+++
date = "2022-06-26"
title = "'Too many authentication failures' when using SSH with password"
tags = ["ssh"]
categories = ["technical"]
+++

Disclaimer: This happened some time ago and the details are fuzzy. This post only purposes is for my own documentation.

I have a raspberry pi running here and didn't bother to set up pubkey authentication.

After adding new SSH keys to my personal machine, SSHing to the raspberry pi failed with

```
Too many authentication failures
```

Which did not make sense at all.


I remember looking at the SSH Server logs and realizing it was still trying to login using my SSH keys,
but since none of them were set up in this server, it was hitting the MaxAuthTries limit (I think by default it's 6).

To not use Pubkey authentication you can use this handy flag `-o PubkeyAuthentication=no`.
