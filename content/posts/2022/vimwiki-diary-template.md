+++
date = "2022-01-22"
title = "Vimwiki/Taskwiki diary template"
categories = ["blog", "technical", "vimwiki"]
+++

I've been using a `vim + vimwiki + taskwiki + taskwarrior + timewarrior` setup for a while. It works really well!

(One day I will try org-mode. Not today though)

One of the bits that got a good impression is the diary template.
Basically it contains:
* A header
* The list of tasks that haven't been completed
* The list of tasks completed on that day


```
#!/usr/bin/env bash
set -euo pipefail

project=$1
date=$(date +'%A, %B %d, %Y')
dateUTC=$(date -u +"%Y-%m-%d")
tomorrowUTC=$(date --date=tomorrow +"%Y-%m-%d")

cat <<EOF
# $date
# Tasks
## Outstanding tasks | status:pending and project:$project and entry.before:$tomorrowUTC
## Finished today | status:completed and project:$project and end.after:$dateUTC and end.before:$tomorrowUTC
EOF
```

Some explanations/disclosures:

* It's invoked with `./diary-template.sh myproject` (see below). I reuse this same template for 2 taskwarrior projects:
**work** and **life** (everything else).
* `entry.before` upon saving back old entries, don't want to update with new data
* pretty sure that `date` syntax only works with GNU date


Then in the vim config:

```vim
au BufNewFile ~/projects/vimwiki/life/diary/*.md :silent 0r !~/projects/vimwiki/templates/task.sh 'life'
```

When a new buffer is created in `projects/vimwiki/life/diary/*.md`, we run that script file and use its contents as the template.

Which gives something like this (example from work):


```markdown
# Friday, October 01, 2021
# Tasks
## Outstanding tasks | status:pending and project:work and entry.before:2021-10-02

## Finished today | status:completed and project:work and end.after:2021-10-01 and end.before:2021-10-02
```

Which `taskwiki` will dynamically render into

```markdown
# Friday, October 01, 2021
# Tasks
## Outstanding tasks | status:pending and project:work and entry.before:2021-10-02
* [ ] add docs for feature XXX #26bf1a75
* [ ] setup pipeline for project #ff6e9878

## Finished today | status:completed and project:work and end.after:2021-10-01 and end.before:2021-10-02
 * [X] loosen linting rules for golang  #3719f8e2
```


If you are interested in creating vimwiki pages for each task, check out my poorly made fork [github.com/eh-am/taskwiki](github.com/eh-am/taskwiki)
