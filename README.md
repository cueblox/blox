# Blox

Blox is the CLI for working with [CueBlox](https://cueblox.com).

## What is Blox?

See our rapidly-evolving documentation [here](https://github.com/cueblox/blox/blob/main/dogfood/data/source/pages/index.md)

## Vocabulary

### Blox

A Blox is a collection of DataSets, grouped into a Schema, and distributed as a repository.

These Blox can be consumed using the `blox` CLI to provide data validation and generation for your content repositories, ensuring type safety across your content.

### DataSet

A DataSet is a type with a strongly defined schema, using [Cue](https://cuelang.org).

See [examples](./dogfood/schemata/profile_v1.cue)

### Schema

A Schema is a Cue file definition of one or more DataSets, with some metadata to help connect some dots for the `blox` CLI.

See [examples](./dogfood/schemata)

### Repository

Collection of schemas, distributed via HTTP with a `manifest.json`. Can be downloaded by the `blox` CLI.

## Shoulders of Giants

Blox would not be possible if not for all of the great projects it depends on. Please see [SHOULDERS.md](SHOULDERS.md) to see a list of them.
