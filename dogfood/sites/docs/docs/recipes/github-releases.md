# Host Your Dataset on GitHub Using Releases

GitHub Releases allow you to host pre-built assets, source code, etc. for download. This recipe takes advantage of GitHub Releases using a fixed release tag name to build your dataset and make it available at a fixed download URL.

## This recipe is the foundation of nearly all the other things we do with cueblox

## Prerequisites

* content repository hosted on GitHub
* blox.cue configuration is complete

## Presumptions

Since your configuration can vary drastically based on your needs, we'll be working with the following assumptions:

Directory structure:

```bash
$ tree -L 1
.
├── README.md
├── data
├── unrelated_directory
└── other_unrelated_directory
```

Our CueBlox managed data lives in the `data` directory, and the configuration file is in that directory as well.

blox.cue contents:
```json
{
  data_dir: "."
  schemata_dir: "schemata"
  build_dir: ".build" <-- Take note of this
  template_dir: "tpl"
  static_dir: "static"
}
```

Important to note: the `data_dir` is set to `.` -- the current directory. This isn't required for this recipe to work, it simply allows our directory structure to be a little more flat. The important piece is `build_dir` which is set to `.build`. These directories are relative to the location of the `blox.cue` file, so in our example, the output from `blox build` will be at `$REPO_ROOT/data/.build/data.json`.

## Releasing with GitHub Actions

Create a new GitHub Action by placing a file in `$REPO_ROOT/.github/workflows`. You can call it anything, if you don't specify a name in the Action's YAML definition, the file name will be used. We will use `data.yaml` in this recipe, to give us a workflow called `data`.

data.yaml:

```yaml
on:
  push:
    paths:
      - .github/**
      - data/** <-- run when files change here

jobs:
  build:
    runs-on: ubuntu-latest
    name: Publish CueBlox
    steps:
      - name: Checkout
        uses: actions/checkout@v2

      - name: Build & Validate Blox Data
        id: build
        uses: cueblox/github-action@v0.0.8
        with:
          directory: data  <-- location of blox.cue

      - uses: marvinpinto/action-automatic-releases@latest
        with:
          repo_token: "${{ secrets.GITHUB_TOKEN }}"
          automatic_release_tag: "blox"
          prerelease: true
          title: "CueBlox"
          files: |
            data/.build/data.json <-- build_dir + data.json
```

This workflow uses the `cueblox/github-action` action to compile and validate your dataset, passing in the working directory `data` to tell `blox` where to look for your `blox.cue` file.

The last step uses `marvinpinto/action-automatic-releases` to create a release of your dataset. Because it specifies `prerelease: true` and `automatic_release_tag: "blox"`, the release will always have the tag `blox`, which means it will always be available at the same URL.

If you follow this pattern your releases will be available at a URL that looks like this:

```
https://github.com/you/reponame/releases/download/blox/data.json
```

Now you have a fixed location to download your dataset which will be updated automatically every time you push new files to your content repository. 

