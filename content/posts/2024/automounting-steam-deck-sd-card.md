+++ 
date = 2024-08-25T12:43:04+01:00
title = "Automounting Steam Deck SD Card"
categories = ["technical"]
+++

Normally (I guess), you need to first enter Desktop Mode, which will automount the SD card.
Then you can go back to Game Mode and play your games from the SD card (mostly emulation :) ).

A good thing of the steam deck is that it just runs Linux, so we can use same old tools we are used too.
What we want here is to add the sd card to `fstab`

Assuming you have SSH access to the deck.

I recommend first automounting it via Desktop mode. It should auto mount to `/run/media/deck/{UUID}`,
you can double check with `df -T`:


```
/dev/mmcblk0p1 btrfs    999868416 269506996 729257116  27% /run/media/deck/7408fbbe-276b-4f85-8494-d9549daa4ba0
```

So mine is using `btrfs`, pay attention to this, if yours is different the fstab config may change.
Not gonna go into filesystems and whatnot, if you want to be able to access the SD card in a non Linux machine
I recommend exFAT (which works in Mac too) and NTFS (Windows only, although I think you can use a third party tool to have write access in a Mac).

Then, since it's already mounted, you can peek into `/etc/mtab` and literally just copy

```
cat /etc/mtab 

/dev/mmcblk0p1 /run/media/deck/7408fbbe-276b-4f85-8494-d9549daa4ba0 btrfs rw,nosuid,nodev,relatime,ssd,space_cache=v2,subvolid=5,subvol=/ 0 0
```

Although we need to modify slightly.
Find that partition's UUID using `lsblk -oNAME,UUID`

Then edit `/etc/fstab`, we need to do 2 changes from what we got from `/etc/mtab`:

1. Prefix with `UUID={YOUR_UUID}`. This is necessary so that in case the drive name changes (right now it's `/dev/mmcblk0p1`), it can be looked up via the UUID.
2.  Add `nofail` to the options, this is super important so that it won't prevent your steam deck to boot in case the SD card can't be read.

So in my case I added all this to `fstab`

```
UUID=7408fbbe-276b-4f85-8494-d9549daa4ba0 /run/media/deck/7408fbbe-276b-4f85-8494-d9549daa4ba0 btrfs nofail,rw,nosuid,nodev,relatime,ssd,space_cache=v2,subvolid=5,subvol=/ 0 0
```

Unmount the existing mount point `umount /run/media/deck/7408fbbe-276b-4f85-8494-d9549daa4ba0` in my case. You may get a `Resource busy`, in that case you can technically track all the places that are using that disk using `lsof` and kill them. Honestly if that happens, I recommend instead to just comment out that line in `fstab`, then rebooting and immediately going to desktop mode, then you can try to umount again.

Run `mount -a` everything that's in `fstab`.

If it works, reboot and it should be automounting.

Now to make it available in Game Mode (to be able to install Steam games), go to Desktop Mode, open Steam -> Settings -> Storage -> Add Drive, then add your SD card.

It should work now.

# bonus: breaking down my fstab config
As I just copied whatever was in `/etc/mtab`, I didn't fully understand everything. So let's look at the docs:

```
/dev/mmcblk0p1 /run/media/deck/7408fbbe-276b-4f85-8494-d9549daa4ba0 btrfs nofail,rw,nosuid,nodev,relatime,ssd,space_cache=v2,subvolid=5,subvol=/ 0 0
```

`nofail` -> already explained, won't prevent the boot to fail if the partition can't be found
`rw` -> obvious enough, to allow reading and writing
`nosuid` -> don't allow having the suid and sgid bits, this is a little complicated, but long story short is a bit that a file can have which will allow anyone to run that as if they were the user that owns the file, the most common case I know is to run as root, like for the `sudo` binary, since sudo itself allows you to run as root (the owner of the file)
`nodev` -> not quite sure, I believe it doesn't allow the filesystem to contain special files, which usually allows talking directly to the device hardware
`relatime` -> [some low level stuff about access time](https://en.wikipedia.org/wiki/Stat_(system_call)#Criticism_of_atime), don't care really much
`ssd` -> for what I understand it enables auto detection of ssd and btrfs add some optimizations (https://man.archlinux.org/man/core/btrfs-progs/btrfs.5.en)
`space_cache=v2` -> defaults to v2 which is better, not sure why it's set (https://man.archlinux.org/man/core/btrfs-progs/btrfs.5.en)
`subvol=/` ->  mounts from a path, not sure it's useful here

Now looking at all this, it looks like we can simplify it a bit if we use `defaults`. I found that the kernel has the following defaults ` rw, suid, dev, exec, auto, nouser, and async.` (https://man7.org/linux/man-pages/man5/fstab.5.html), but in reality it may depend on the filesystem, apparently btrfs [don't set anything](https://www.jwillikers.com/btrfs-mount-options)? But then [someone says there are defaults](https://old.reddit.com/r/btrfs/comments/14t7mjr/where_to_find_default_mount_options/jr21uj7/) which are different from the kernel's one.

You know what? I will leave as it is.

