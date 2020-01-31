+++
date = "2020-01-20"
title = "Unnoficial Bash Strict Mode"
slug = "unnoficial-bash-strict-mode"
categories = ["blog"]
+++

Bash (and shell scripting in general) is a mess. It's easy to mess up if you don't know what you are doing. If you come from traditional programming background and just want to plumb a few lines of code (specially common in this container world), there are a few behaviours that will confuse the hell out of you.

To help wih that the [unnoficial bash strict mode](http://redsymbol.net/articles/unofficial-bash-strict-mode/) was created. In this post we will go over each problem and how does the strict mode help (and its quirks)

## Errexit

### The problem

Given the following bash script (assuming the file referenced does not exist)
{{< file "/posts/bash-strict-mode/errexit.sh" "bash" >}}

Should it run all the way through, or should it fail due to the file being missed?

Turns out even though the file does not exit, the script continues running just fine! In this trivial case it does not do any harm, but it can easily bite you:

{{< file "/posts/bash-strict-mode/errexit_fairly_complicated.sh" "bash" >}}

Assume `generate_whitelist` is not installed for some reason. You got the idea.

To tackle the problem I tend to resort to the "Fail Fast" approach. From the [Fail Fast - C2 wiki](https://wiki.c2.com/?FailFast):

> This is done upon encountering an error so serious that it is possible that the process state is corrupt or inconsistent, and immediate exit is the best way to ensure that no (more) damage is done.

Sounds reasonable. So why is that "Fail silentyl" the normal behaviour in a shell script? Well, think that in a context of a shell you DO not want to exit when there's an error (imagine crashing your shell when you cat a file that does not exist). Looks like the behaviour was carried out to shellscript.

### The solution

So how do you fix this? By setting the flag `errexit`.

```sh
#!/usr/bin/env bash
set -e
```

or if you like ~~verbosity~~ explictness:

```sh
#!/usr/bin/env bash
set -o errexit
```

What does this do? As per the docs

> Exit immediately if a pipeline (...) returns a non-zero status.
>
> -- https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html#The-Set-Builtin

So going back to the first example, we would do it like this
{{< file "/posts/bash-strict-mode/errexit2.sh" "bash" >}}

Which would then fail, since the file does not exist. The behaviour can be described on the following [bats](https://github.com/sstephenson/bats) unit test.

{{< file "/posts/bash-strict-mode/errexit.bats" "bash" >}}

That works wonderfully and will definitely (IMO) help you, but be aware:

### Quirks

#### Quirk #1: Programs that return non-zero status

Not all commands return 0 when it runs correctly. The most proeminent example is `grep`. From the docs:

> Normally the exit status is 0 if a line is selected, 1 if no lines were selected, and 2 if an error occurred.
> However, if the -q or --quiet or --silent is used and a line is selected, the exit status is 0 even if an error occurred.

So given example below
{{% file "/posts/bash-strict-mode/grep_fail.sh" "bash" %}}
the `echo` will never be run.

So what can we do in that situation? Thankfully there's a bit in the bash manual that can help us (reformatted for clarity):

> The shell does not exit if the command that fails is
>
> 1. part of the command list immediately following a while or until keyword,
> 2. part of the test in an if statement,
> 3. part of any command executed in a && or || list except the command following the final && or ||
> 4. any command in a pipeline but the last, or if the command’s return status is being inverted with !.

In our case, we can simply rewrite to comply with 2:

{{< file "/posts/bash-strict-mode/grep_correct.sh" "bash" >}}

The behaviour can be verified by the following bats test

{{< file "/posts/bash-strict-mode/grep.bats" "bash" >}}

<!--
#### Quirk #2: It does not propagate to subshells (sometimes)

Not exclusive to `errexit`, but it's good to know that in:
{{% file "/posts/bash-strict-mode/subshell.sh" "bash" %}}
The echo from the inner script will run. Which is not what we want.

##### A little detour

If you try to run the example above it may not work. Since `errexit` inheritance depends on `inherit_errexit` to be enabled (which is only available from bash 4.4 onwards). From the docs:

> inherit_errexit
> If set, command substitution inherits the value of the errexit option, instead of unsetting it in the subshell environment. This option is enabled when POSIX mode is enabled.

To turn it off (to reproduce the original example behaviour, add to the shellscript. Here `-u` removes and `-s` adds)

```
shopt -u inherit_errexit
```

Ok, now assuming we don't have to deal with any `inherit_errexit` bullcrap. To properly export to subshells, add

```bash
export SHELLOPTS
```

Which will basically reexport "SHELLOPTS" var to the subshell.

-->

#### Quirk 3: What if you DO want to allow a command to fail?

In that case, OR with true:

```sh
kubectl create namespace my-app || true
```

(just FYI I would do this declaratively with `kubectl create ns my-app -o yaml --dry-run | kubectl apply -f -`)

Let's think a little bit why this works, from the docs (which we read in previous point)

> 3. part of any command executed in a && or || list except the command following the final && or ||

That has some implications:
{{< file "/posts/bash-strict-mode/errexit_or.sh" "bash" >}}
{{< file "/posts/bash-strict-mode/errexit_or.bats" "bash" >}}

Another option would be to turn it off momentaneously

```sh
set +e
command_allowed_to_fail
set -e
```

The `+` syntax mean "remove" and `-` means "to add" (go figure). Therefore we are simply disabling that feature between our `command_allowed_to_fail` calls!

#### Quirk 4: How do I know which command failed?

Not really a quirk, but often you need to know where it failed.

{{< file "/posts/bash-strict-mode/which_command_failed.sh" "bash" >}}

How can you tell which command failed? (Apart from looking at the very obvious mistake)

- 1. echo everything you are doing
     Pros: straightforward
     Cons: quite boring to do so

- 2. `set -x`, which will print every instruction.
     Pros: simple to add
     Cons: you may end up exposing more than you want (imagine printing a variable with secrets, now imagine that being run in a ci)

- 3. put a `trap` to [print the line number when a command fails](https://intoli.com/blog/exit-on-errors-in-bash-scripts/)
     Pros: can be added globally
     Cons: I've found any yet

#### It only works in "pipelines"

```
#!/usr/bin/env bash

set -e


divisor="0"
local division=$((10/$divisor))

echo "hello world"
```

#### Is there anything more?

Yup. Read the entry on [BashFAQ](http://mywiki.wooledge.org/BashFAQ/105) and its subsequent links.

# References

- http://redsymbol.net/articles/unofficial-bash-strict-mode/
- https://dacav.roundhousecode.com/blog/2019-10/02-more-shell-inconsistencies.html
- https://intoli.com/blog/exit-on-errors-in-bash-scripts/
- http://stratus3d.com/blog/2019/11/29/bash-errexit-inconsistency/
- https://gist.github.com/yoramvandevelde/fb4d8aa6fa6d8b1eab6da81b62373d85

```

```
