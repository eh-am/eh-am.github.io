+++
date = "2022-07-17"
title = "from i3 to yabai / setting up a new mac"
categories = ["blog", "i3", "yabai", "linux", "osx", "mac"]
+++

I've been using linux (fedora/arch) + i3 for a few years by now.

Then I got a mac for work.

To be fair, I >asked< for a mac because I wanted to experiment the whole arm hype. Spoiler: it's quite good indeed. Worth it disrupting an existing flow? Don't know yet.

But anyway, I had a whole setup in linux that I had to adapt to this new machine. Some of it it's just [second-system syndrome](https://en.wikipedia.org/wiki/Second-system_effect) (like keeping track of installed tools), some of it are indeed things I had to adapt.

As always, primarily intend is my own documentation, therefore I may overlook certain aspects.

# tl;dr
{{< video src="./capture.mp4" type="video/mp4" preload="auto" >}}



# Brew
Bit a of a no-brainer, everybody uses it and I had no problems using it when I had a mac.

One thing I am doing now is using a [Brewfile](https://github.com/Homebrew/homebrew-bundle) to keep track of my installed tools. Preferably with a comment reminding me why the heck I installed that. 

For comparison, in my Fedora desktop I haven
```bash
> sudo dnf list installed | wc -l
5818
```

Although this is not a fair comparison, since it also show dependencies (ie packages I have not >actively< installed). Those are:


```
sudo dnf history userinstalled | wc -l
764
```

Source: [DNF Command Reference — dnf latest documentation](https://dnf.readthedocs.io/en/latest/command_ref.html#history-store-command-label)

I may go back and use brew with linux as well.

# yadm
For managing the dotfiles, there're [plentiful of tools](https://dotfiles.github.io/utilities/). I decided to use [yadm](https://yadm.io/) for a few reasons:
* written in Bash. I mean this is not a good reason, it's just that I have a preference for it over tools written in python.
* git-like interface. In the end I am still modifying raw text files. One thing that particular annoyed me with chezmoi is the use of [chezmoi edit](https://www.chezmoi.io/reference/commands/edit/) to edit files. I want to just `vim MYFILE.txt`!
* support for [alternative files](https://yadm.io/docs/alternates), which is just like "if linux use this file, if mac use this other file"
* support for [templates](https://yadm.io/docs/templates#) in case you want to use more or less the same file with slight differences. I started with alternative files and once I see the patterns emerge I may merge the files back into a template.


## yadm bootstrap
running `yadm bootstrap` allows you to execute a script/binary when {first cloning/pulling} the dotfiles repo.

It's a good idea to make it idempotent.

In my case, I made it run [`brew bundle`](https://yadm.io/docs/bootstrap#install-homebrew-and-a-bundle-of-recipes)
Also made it run our next tool, `ansible`


# ansible
Ideally everything would be setup using `brew` + config files. But it's not always possible.
So I added `ansible` to run some stuff like cloning the [tmux plugin manager repo](https://github.com/tmux-plugins/tpm),
or hiding the dock.


# yabai/skhd
To get the most similar experience to i3, we need both [koekeishiya/yabai: A tiling window manager for macOS based on binary space partitioning](https://github.com/koekeishiya/yabai) and [koekeishiya/skhd: Simple hotkey daemon for macOS](https://github.com/koekeishiya/skhd). yabai is the tiling manager, and skhd is a hotkey daemon.


One slightly inconvenience is having to [disable System Integrity Protection](https://github.com/koekeishiya/yabai/wiki/Disabling-System-Integrity-Protection).

I won't go over my hotkeys, feel free to check the dotfiles repo.

# statusbar
Installed [cmacrae/spacebar](https://github.com/cmacrae/spacebar) to have an indicator for the workspace.

Looking at my notes I found that was the first moment I was required to install rosetta.
```
sudo softwareupdate --install-rosetta --agree-to-license
```

# Keyboard setup
Mapped certain keys:

| original | to|
|-|-|
| Caps Lock | escape |
| Option | Command |
| Command | Option |


The Caps lock is a classic one for vim users.
Option / Command it's just that I am already used to use $mod key (normally command), since I run arch (btw) in a old Macbook air.
However if I mapped skhd to CMD, lots of things would break since OSX uses it quite a bit (as opposed to alt).

Also set up key repeat / delay until repeat to the fastest option available.
I tried setting this up with ansible but it didn't work.

# Wrapping up
Of course there are plenty more things to be configured or things I've omitted.
Feel free to explore the dotfiles repo here: [eh-am/dotfiles](https://github.com/eh-am/dotfiles)

