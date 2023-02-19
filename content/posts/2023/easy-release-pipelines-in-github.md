+++ 
date = 2023-02-19T11:11:52Z
title = "Easy release pipelines in Github with Release Please"
categories = ["technical"]
+++

There's this handy approach I use when releasing open source projects in Github that I would like to share.

The idea boils down to:
* Using [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) to write commits in a semantic manner which can be used by tools
* Using [Release Please](https://github.com/googleapis/release-please) to automatically create Pull Requests, which upon merge creates a Release with all the Changelog jazz.

To enforce people will use conventional commits, there are couple things you need to do:


Create a Github Actions workflow that validates that PR is indeed using `conventional commits`:

```yaml
name: conventional-pr
on:
  pull_request:
    branches:
      - main
    types:
      - opened
      - edited
      - synchronize
jobs:
  lint-pr:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: CondeNast/conventional-pull-request-action@v0.2.0
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          commitTitleMatch: "false"
          ignoreCommits: "true"
```

The catch here is to in the Repository configuration, to set
```
Allow squash merging

Combine all commits from the head branch into a single commit in the base branch.
```

And

```
Default to pull request title and description
```

That way, the commits themselves don't have to follow `conventional commits`, only the PR title (which can then be edited via the `UI`). Also, upon changing it the checks run again.

![The Setup you need to do in the GitHub repository](../publish-please.png)


Then, of course, the `Publish Please` workflow:
```yaml
on:
  push:
    branches:
      - main
name: release-please
jobs:
  release-please:
    runs-on: ubuntu-latest
    steps:
      - uses: google-github-actions/release-please-action@v3
        with:
          release-type: go
          package-name: release-please-action
          token: ${{ secrets.RELEASE_PLEASE_TOKEN }}
```

Lastly, you need to do build the application when it's released. Here I used `wangyoucao577/go-release-action` since it's a `go` project and I needed to build for different OSes/Architectures but realistically just change to your needs:

```yaml
on:
  release:
    # https://stackoverflow.com/a/61066906
    types: [published]

jobs:
  publishmatrix:
    name: Release Go Binary
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [darwin, linux]
        goarch: [amd64, arm64]
    steps:
    - uses: actions/checkout@v3
    - uses: wangyoucao577/go-release-action@v1.35
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        goos: ${{ matrix.goos }}
        goarch: ${{ matrix.goarch }}
        binary_name: "my-binary"
        extra_files: LICENSE README.md
        # The default naming pattern uses the release name
        # Which makes difficult for scripts to download the correct version
        asset_name: my-binary-ci-${{ matrix.goos }}-${{ matrix.goarch }}
        pre_command: export CGO_ENABLED=0
```

One important point is that you need to react to the `published` event, since when publishing via automated tools (like `Release Please`), the `created` step is bypassed.

Here's an example of it being in use:
![](../publish-please-2.png)
![](../publish-please-3.png)
