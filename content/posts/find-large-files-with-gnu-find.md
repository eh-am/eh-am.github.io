+++
date = "2021-11-21"
title = "Find large files with GNU find"
slug = "find-large-files-with-gnu-find"
categories = ["blog", "shell", "technical"]
+++

I edit my [podcast](www.triceratops.show) using Audacity, which creates big (>1GB) project files (.aup3).
In an ideal world, I would keep these project files, for a nice remaster when I get better at editing.
In the real world, however, these files occupy space. I can technically backup them. But contrary to what people say, for the end user, storage is still expensive. I don't have the calculations in my hand right now, but it's either store in the cloud ($$$ expensive) or manage your backups at time (time expensive).

So I decided to take the L and just delete these files from time to time.

In practice, I already know the files I want to delete (the ones named `.aup3`). But more than that, I want to first identify what are them and their sizes.

# Find files bigger than X
First step is to find what are the files bigger than X. Find has a `-size` parameter that can help here:

```shell
$ find . -size +1G

>
./episodio-8/decupagem.aup3
./episodio-8/episodio-8.aup3
./episodio-10/ep-10.aup3
./episodio-15/ep15.aup3
./episodio-15/decupagem halloween.aup3
./episodio-7/episodio-7.aup3
./episodio-12/decupagem.aup3
./episodio-12/ep12.aup3
./episodio-13/decupagem.aup3
./episodio-13/ep13.aup3
./episodio-18/episodio-18.aup3
./episodio-5/Episodio-5.aup3
./episodio-6/episodio-6-de-verdade.aup3
./episodio-6/episodio-6.aup3
./episodio-9/epi-9.aup3
./episodio-9/decupagem.aup3
./episodio-17/ep17.aup3
./episodio-17/ep 17du1.aup3
./episodio-11/2021-09-23_21-15-20 audio track 1.aup3
./episodio-11/epi-11.aup3
./bah.aup3
./episodio-16/episodio-16.aup3
./episodio-14/decupagem.aup3
```

Cool, but what are their sizes?

# Find files bigger than X and printing their sizes
`find` has a parameter `-printf`, [which allows us to do many things](https://man7.org/linux/man-pages/man1/find.1.html):

```
 %s     File's size in bytes.
```

```shell
$ find . -size +1G -printf "%s"
50197299202410676224237626982414056161281581056000369596825633111408642186543104446188748827419607041681915904167877017623684055041110310912281247744038853017601447886848151873126452540866563044802560529668505615465840642232942592⏎
```

Huh oh, I need to format that. And add the file names back.

```shell
$ find . -size +1G -printf "%s\t%p\n"

>
5019729920	./episodio-8/decupagem.aup3
2410676224	./episodio-8/episodio-8.aup3
2376269824	./episodio-10/ep-10.aup3
1405616128	./episodio-15/ep15.aup3
1581056000	./episodio-15/decupagem halloween.aup3
3695968256	./episodio-7/episodio-7.aup3
3311140864	./episodio-12/decupagem.aup3
2186543104	./episodio-12/ep12.aup3
4461887488	./episodio-13/decupagem.aup3
2741960704	./episodio-13/ep13.aup3
1681915904	./episodio-18/episodio-18.aup3
1678770176	./episodio-5/Episodio-5.aup3
2368405504	./episodio-6/episodio-6-de-verdade.aup3
1110310912	./episodio-6/episodio-6.aup3
2812477440	./episodio-9/epi-9.aup3
3885301760	./episodio-9/decupagem.aup3
1447886848	./episodio-17/ep17.aup3
1518731264	./episodio-17/ep 17du1.aup3
5254086656	./episodio-11/2021-09-23_21-15-20 audio track 1.aup3
3044802560	./episodio-11/epi-11.aup3
5296685056	./bah.aup3
1546584064	./episodio-16/episodio-16.aup3
2232942592	./episodio-14/decupagem.aup3
```

Nice! But would be nice to know their sizes in human readable format.

# Find files bigger than X and printing their sizes in human readable format
Here we will use the `numfmt` tool, which allows us to **format numbers**. God I love boring, clear names.

To be fair, as it description says, **Convert numbers from/to human-readable strings**. So the human part is not that obvious.

Anyway, since I am not familiar with that tool, let's see an example on [tl;dr](https://tldr.sh/)

```

Convert 5th field (1-indexed) to IEC Units without converting header:

ls -l | numfmt --header={{1}} --field={{5}} --to={{iec}}
```

Cool this is more or less what I want to do. I want to convert the first field to IEC.


```shell
$ find . -size +1G -printf "%s\t%p\n" | numfmt --field={{1}} --to={{iec}}

numfmt: invalid field value ‘{{1}}’
Try 'numfmt --help' for more information.
```

Right. Let me try without the curly brackets.

```shell
$ find . -size +1G -printf "%s\t%p\n" | numfmt --field=1 --to=iec

>
4.7G ./episodio-8/decupagem.aup3
2.3G ./episodio-8/episodio-8.aup3
2.3G ./episodio-10/ep-10.aup3
1.4G ./episodio-15/ep15.aup3
1.5G ./episodio-15/decupagem halloween.aup3
3.5G ./episodio-7/episodio-7.aup3
3.1G ./episodio-12/decupagem.aup3
2.1G ./episodio-12/ep12.aup3
4.2G ./episodio-13/decupagem.aup3
2.6G ./episodio-13/ep13.aup3
1.6G ./episodio-18/episodio-18.aup3
1.6G ./episodio-5/Episodio-5.aup3
2.3G ./episodio-6/episodio-6-de-verdade.aup3
1.1G ./episodio-6/episodio-6.aup3
2.7G ./episodio-9/epi-9.aup3
3.7G ./episodio-9/decupagem.aup3
1.4G ./episodio-17/ep17.aup3
1.5G ./episodio-17/ep 17du1.aup3
4.9G ./episodio-11/2021-09-23_21-15-20 audio track 1.aup3
2.9G ./episodio-11/epi-11.aup3
5.0G ./bah.aup3
1.5G ./episodio-16/episodio-16.aup3
2.1G ./episodio-14/decupagem.aup3
```

Awesome! But it's missing something!

# Find files bigger than X and printing their sizes in human readable format sorted by size

In this case we can simply use `sort` tool, in `--reverse`.

```shell
$ find . -size +1G -printf "%s\t%p\n" | numfmt --field=1 --to=iec | sort -h --reverse

>
5.0G ./bah.aup3
4.9G ./episodio-11/2021-09-23_21-15-20 audio track 1.aup3
4.7G ./episodio-8/decupagem.aup3
4.2G ./episodio-13/decupagem.aup3
3.7G ./episodio-9/decupagem.aup3
3.5G ./episodio-7/episodio-7.aup3
3.1G ./episodio-12/decupagem.aup3
2.9G ./episodio-11/epi-11.aup3
2.7G ./episodio-9/epi-9.aup3
2.6G ./episodio-13/ep13.aup3
2.3G ./episodio-8/episodio-8.aup3
2.3G ./episodio-6/episodio-6-de-verdade.aup3
2.3G ./episodio-10/ep-10.aup3
2.1G ./episodio-14/decupagem.aup3
2.1G ./episodio-12/ep12.aup3
1.6G ./episodio-5/Episodio-5.aup3
1.6G ./episodio-18/episodio-18.aup3
1.5G ./episodio-17/ep 17du1.aup3
1.5G ./episodio-16/episodio-16.aup3
1.5G ./episodio-15/decupagem halloween.aup3
1.4G ./episodio-17/ep17.aup3
1.4G ./episodio-15/ep15.aup3
1.1G ./episodio-6/episodio-6.aup3
```


To be fair. In my first try I didn't assume sort could work with human readable values. So as an alternative, I a) printed another column with the sizes in non human format and used that to sort or b) sorted before `numfmt`:


```shell
$ find . -size +1G -printf "%s\t%p\n" | sort --reverse | numfmt --field=1 --to=iec

5.0G ./bah.aup3
4.9G ./episodio-11/2021-09-23_21-15-20 audio track 1.aup3
4.7G ./episodio-8/decupagem.aup3
4.2G ./episodio-13/decupagem.aup3
3.7G ./episodio-9/decupagem.aup3
3.5G ./episodio-7/episodio-7.aup3
3.1G ./episodio-12/decupagem.aup3
2.9G ./episodio-11/epi-11.aup3
2.7G ./episodio-9/epi-9.aup3
2.6G ./episodio-13/ep13.aup3
2.3G ./episodio-8/episodio-8.aup3
2.3G ./episodio-10/ep-10.aup3
2.3G ./episodio-6/episodio-6-de-verdade.aup3
2.1G ./episodio-14/decupagem.aup3
2.1G ./episodio-12/ep12.aup3
1.6G ./episodio-18/episodio-18.aup3
1.6G ./episodio-5/Episodio-5.aup3
1.5G ./episodio-15/decupagem halloween.aup3
1.5G ./episodio-16/episodio-16.aup3
1.5G ./episodio-17/ep 17du1.aup3
1.4G ./episodio-17/ep17.aup3
1.4G ./episodio-15/ep15.aup3
1.1G ./episodio-6/episodio-6.aup3
```

# Bonus: finding the total size
Here we throw my beloved `awk` into the mix. Basically we print all sizes, and iterate over each line summing into a variable `s`, which we then print in the `END`. Then we just format that to `iec` (see how you don't actually need to pass `--field` if it's the first column?):

```shell
$ find . -size +1G -printf "%s\n" | awk '{ s+=$1 } END {print s}' | numfmt --to=iec
> 59G
```

# Conclusion
Then from that point on I just deleted the files manually, since I needed to manually think "do I want to delete this?".

But anyway, this post shows one thing I really love about these unix tools: composability. I can go iteratively adding tools on top of each other and it just...works (most of the time).
