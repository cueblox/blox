---
title: Creating Your Blox
publish_date: '2021-03-19'
section_id: getting-started
weight: 1
---

# Creating Your Blox

CueBlox needs a configuration file and some directories to start validating and transforming your data. You can have CueBlox create this structure for you with `blox init`.

```shell
blox init
 INFO  Initialized folder structures.
```

This creates a `blox.cue` file with the default configuration.

```shell
cat blox.cue
{
  build_dir:    "_build"
  data_dir:     "data"
  schemata_dir: "schemata"
  static_dir:   "static"
}
```
