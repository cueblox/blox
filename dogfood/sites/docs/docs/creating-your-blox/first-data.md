---
title: Creating Data
publish_date: '2021-03-19'
section_id: getting-started
weight: 1
---

# Creating Data

Continuing from the previous section, we now want to add some data to our blox that conforms to the required schema. Lets remind ourselves of the schema, trimmed of metadata for brevity:

```cue
# schemata/artist.cue
{
    _schema: { ... }

	#Artist: {
		_dataset: {
            plural: "artists"
            supportedExtensions: ["yaml", "yml", "md"]
        }

        name: string
        url?: string
        hungry: bool | *false
    }
}
```

## Manually Creating Data

The `plural` field from our DataSet metadata tells us which directory our data, artists, needs to be in to be validating correctly against our schema.

Blox supports both YAML and Markdown (with frontmatter) content, and we'd welcome other loaders for more. We can constrain the filetypes we want to load for this particular DataSet through `supportedExtensions`.

This means it's really trivial to add new data to be validated. In this instance, we create some YAML or Markdown files in the `artists` directory.

```shell
# linkin-park.yml
name: Linkin Park
url: https://linkinpark.com

# fleetwood-mac.md
---
name: Fleetwood Mac
url: https://fleetwoodmac.com
---

Fleetwood Mac are a British-American rock band, formed in London in 1967.
```

### The "Body" Field

When using markdown files with YAML frontmatter and a body, CueBlox does some "magic" translation. Lets use the Fleetwood Mac example above. This actually translate, in YAML, to:

```yaml
name: Fleetwood Mac
url: https://fleetwoodmac.com
body: |
    Fleetwood Mac are a British-American rock band, formed in London in 1967.
```

This will actually **fail** validation with our current schema, as `body` isn't a valid property on `#Artist`.  So remember to include a `body` field when using this data format.

## Using `blox` to Create Data

CueBlox's CLI, `blox`, also ships with a helper to scaffold new data.

```shell
❯ blox new
Error: requires at least 1 arg(s), only received 0
Usage:
  blox new [flags]

Flags:
      --dataset string   Which DataSet to create content for?
  -h, --help             help for new
```

Using `blox new --dataset artist soilwork`, we can have CueBlox create a new data file for us.

Note: Currently, this only supports writing YAML. PR's welcome.

The `blox` command works by leveraging `@template` annotations within the schema. Let's update our schema to take advantage of this awesome-sauce.

```cue
# schemata/artist.cue
{
    _schema: { ... }

	#Artist: {
		_dataset: { ... }

        name: string           @template(Chevelle)
        url?: string           @template(https://google.com)
        hungry: bool | *false
    }
}
```

Now, using `blox new --dataset artist soilwork`, we'll get a data file at `./artists/soilwork.yaml` that looks like:

```
# soilwork.yaml
name: Chevelle
url: https://google.com
hungry: false
```

This is early stage for templates, but we're going to be adding more awesome soon; including an interactive mode to accept user input for fields.
