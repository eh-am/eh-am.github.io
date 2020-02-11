+++
date = "2020-01-20"
title = "Unnoficial Bash Strict Mode"
slug = "unnoficial-bash-strict-mode"
categories = ["blog", "bash", "technical"]
+++

Bash (and shell scripting in general) is NOT straightforward. It's easy to mess up if you don't know what you are doing. If you come from a traditional programming background and just want to plumb a few lines of code, there are a few behaviours that will confuse the hell out of you.

To help with that, the [unnoficial bash strict mode](http://redsymbol.net/articles/unofficial-bash-strict-mode/) was created. In this post, we will go over misleading behaviors and how the strict mode can be helpful in each case (quirks included).

## errexit

### errexit: The problem

In the following bash script:
{{< file "/posts/bash-strict-mode/errexit.sh" "bash" >}}

Given the file does not exist, should it run all the way through or should it fail?

Turns out the script continues running just fine!

To tackle the problem, I tend to resort to the "Fail Fast" approach. From the [Fail Fast - C2 wiki](https://wiki.c2.com/?FailFast):

> This is done upon encountering an error so serious that it is possible that the process state is corrupt or inconsistent, and immediate exit is the best way to ensure that no (more) damage is done.

Sounds reasonable. So why is "Fail silently" the normal behaviour in a shell script? Well, think that in the context of a shell you DO NOT want to exit when there's an error (imagine crashing your shell when you `cat` a file that does not exist). Looks like the behaviour was simply carried out to the non-interactive shell.

### errexit: The solution

So how can we improve this behaviour? By setting the flag `errexit`.

```sh
set -o errexit
```

or the shorthand version (more commonly used):

```sh
set -e
```

What does this do? As per the docs:

> Exit immediately if a pipeline (...) returns a non-zero status.
>
> -- https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html#The-Set-Builtin

Going back to our example, we would do it instead:
{{< file "/posts/bash-strict-mode/errexit2.sh" "bash" >}}

Which would then fail. Since the file does not exist, `cat` returns a non-zero exit code. This behaviour is described in the following [bats](https://github.com/sstephenson/bats) unit test:

{{< file "/posts/bash-strict-mode/errexit.bats" "bash" >}}

That works and will definitely [IMO](http://mywiki.wooledge.org/BashFAQ/105) help you, but be aware:

### Errexit: Quirks

#### Quirk #1: Programs that return non-zero status

Not all commands return 0 on successfull runs. The most proeminent example is `grep`. From the docs:

> Normally the exit status is
>
> 1. 0 if a line is selected,
> 2. 1 if no lines were selected,
> 3. and 2 if an error occurred.
>
> However, if the -q or --quiet or --silent is used and a line is selected, the exit status is 0 even if an error occurred.

So in the example below, `echo` will never be run.
{{< file "/posts/bash-strict-mode/grep_fail.sh" "bash" >}}

What can we do in that situation? Thankfully there's a bit in the bash manual on `errexit` section that can help us (reformatted for clarity):

> The shell does not exit if the command that fails is
>
> 1. part of the command list immediately following a while or until keyword,
> 2. part of the test in an if statement,
> 3. part of any command executed in a && or || list except the command following the final && or ||
> 4. any command in a pipeline but the last, or if the command’s return status is being inverted with !.

In our case, we can simply rewrite to comply with 2:

{{< file "/posts/bash-strict-mode/grep_correct.sh" "bash" >}}

This behaviour can be verified by the following bats test:

{{< file "/posts/bash-strict-mode/grep.bats" "bash" >}}

#### Quirk 2: What if you are ok with a command failing/returning non-zero?

In that case, simply OR with `true`:

```sh
rm *.log || true
```

Since we do not want to fail if there are no log files.

Let's think a little bit why this works. From the docs (which we already read in a previous point):

> The shell does not exit if the command that fails is
> (...)
>
> 3. part of any command executed in a && or || list except the command following the final && or ||

As the command following the final `||` is `true`, there's no way for the whole line to fail.

Another option would be to turn it off momentaneously:

```sh
set +e
command_allowed_to_fail
set -e
```

The `+` syntax means "remove" and `-` means "to add" (go figure). Therefore, we are simply disabling that feature while our `command_allowed_to_fail` is called!

#### Bonus point: How do I know which command failed?

Not really a quirk, and not specific to `errexit`, but often you need to know where it failed.

{{< file "/posts/bash-strict-mode/which_command_failed.sh" "bash" >}}

How can you tell which command failed? (Apart from looking at the very obvious mistake)

- 1. `echo` everything you are doing \
     **Pros**: straightforward \
     **Cons**: quite boring to do so

- 2. `set -x`, which will print every instruction. \
     **Pros**: simple to add \
     **Cons**: you may end up exposing more than you want (imagine printing a variable with secrets, now imagine that running in a CI environment)

- 3. put a `trap` to [print the line number when a command fails](https://intoli.com/blog/exit-on-errors-in-bash-scripts/) \
     **Pros**: can be added globally \
     **Cons**: bit verbose

#### Is there anything more?

Yup. Once you get the gist of it, read the entry on [BashFAQ](http://mywiki.wooledge.org/BashFAQ/105) and its linked resources.

## Pipefail

### Pipefail: The problem

{{< file "/posts/bash-strict-mode/pipefail_first.sh" "bash" >}}
This would run just fine!

Unfortunately `errexit` is not enough to save us here. From the docs, again:

> The shell does not exit if the command that fails is:
> (...)
>
> 4. any command in a pipeline but the last, or if the command’s return status is being inverted with !.

> The exit status of a pipeline is the exit status of the last command in the pipeline

### Pipefail: The solution

Let's set `pipefail`:

> If pipefail is enabled, the pipeline’s return status is the value of the last (rightmost) command to exit with a non-zero status, or zero if all commands exit successfully.

In other words, it will only return 0 if all parts of the pipeline return 0.

As opposed to `errexit`, `pipefail` can only be set by its full form:

```
set -o pipefail
```

Let's fix the example shown before:
{{< file "/posts/bash-strict-mode/pipefail_first_correct.sh" "bash" >}}

Both behaviours are verified by the following Bats test:
{{< file "/posts/bash-strict-mode/pipefail_first.bats" "bash" >}}

### Pipefail: Quirks

#### Quirk 1:

> the pipeline’s return status is the value of the last (rightmost) command to exit with a non-zero status

{{< file "/posts/bash-strict-mode/pipefail_quirk_1.sh" "bash" >}}

`cat`'s exit code is `1` for when the file does not exist. And `xargs`'s exit code is `123` "if any invocation of the command exited with status 1-12".
Obviously both calls are broken, but what exit code do we get here?

The answer is `123`, which is not ideal.

My recommendation for this case is to simply break it down into different instructions:

{{< file "/posts/bash-strict-mode/pipefail_quirk_1_correct.sh" "bash" >}}

This behaviour can be confirmed by the following bats test:

{{< file "/posts/bash-strict-mode/pipefail_quirk_1.bats" "bash" >}}

#### Quirk 2:

Be careful with what you pipe:

{{< file "/posts/bash-strict-mode/pipefail_quirk_2.sh" "bash" >}}

In this example, we are loading a whitelist file, feeding it to another command (here implemented as a function) that passes it to yet another service (e.g. a CLI tool). Even though the file does not exist, the pipeline does not fail. This ends up passing an empty string to `remove_hosts`, which could have catastrophic effects! (Deleting more than you expect).

Ideally, you want to fail as soon as possible. The best way to do so is to break it down into more instructions and just be more careful ¯\_(ツ)\_/¯

{{< file "/posts/bash-strict-mode/pipefail_quirk_2_correct.sh" "bash" >}}

As always, this behaviour is described by the following bats file:
{{< file "/posts/bash-strict-mode/pipefail_quirk_2.bats" "bash" >}}

For more examples, check [Examples of why pipefail is really important to use](https://gist.github.com/yoramvandevelde/fb4d8aa6fa6d8b1eab6da81b62373d85).

## nounset

Last but not least, this one is very straightforward.

### nounset: The problem

{{< file "/posts/bash-strict-mode/nounset.sh" "bash" >}}

### nounset: The solution

> Treat unset variables and parameters other than the special parameters ‘@’ or ‘\*’ as an error when performing parameter expansion. An error message will be written to the standard error, and a non-interactive shell will exit.

{{< file "/posts/bash-strict-mode/nounset.bats" "bash" >}}

### nounset: Quirks

#### Quirk 1:

As mentioned in the docs, `@` and `*` are treated differently:

{{< file "/posts/bash-strict-mode/nounset_quirk_1.sh" "bash" >}}

So always verify the arguments you are getting are actually correct:
{{< file "/posts/bash-strict-mode/nounset_quirk_1_correct.sh" "bash" >}}

{{< file "/posts/bash-strict-mode/nounset_quirk_1.bats" "bash" >}}

## Conclusion

I hope this is enough to:

- illustrate how the expectations we often have are not true;
- how 'unnoficial strict mode' can help;
- and how the strict mode is not a panacea!

## References/Recommended Readings

- [The Let It Crash Philosophy Outside Erlang](http://stratus3d.com/blog/2020/01/20/applying-the-let-it-crash-philosophy-outside-erlang/)
- [Bash Strict Mode](http://redsymbol.net/articles/unofficial-bash-strict-mode/)
- [Dacav's Home: More shell inconsistencies](https://dacav.roundhousecode.com/blog/2019-10/02-more-shell-inconsistencies.html)
- [How to Exit When Errors Occur in Bash Scripts](https://intoli.com/blog/exit-on-errors-in-bash-scripts/)
- [Bash Errexit Inconsistency - Stratus3D](http://stratus3d.com/blog/2019/11/29/bash-errexit-inconsistency/)
- [Examples of why pipefail is really important to use](https://gist.github.com/yoramvandevelde/fb4d8aa6fa6d8b1eab6da81b62373d85)
