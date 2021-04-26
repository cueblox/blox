---
title: Validation
excerpt: Data Validation
publish_date: '2021-03-19'
section_id: getting-started
weight: 1
---

# Validation

Validation is the first thing that CueBlox can bring to your YAML and Markdown content. In a CueBlox root, you can run `blox build` to validate your content against the defined schema. If you don't have a CueBlox project yet, you can create one with `blox init` or use the `cueblox/blox` repository itself.

Build will validate all of the content and transform it into a JSON file to be consumed by other tools.

```shell
# From the cueblox/blox Git clone
blox build
 INFO  processing images in dogfood/static
 SUCCESS  Validation Complete
 SUCCESS  Data blox written to '_build/data.json'
```
