---
title: Defining Your First Schema
publish_date: '2021-03-19'
section_id: getting-started
weight: 1
---

# Defining Your First Schema

Now that we have our configuration and a directory structure, we need to provide a schema to validate any data we add to the data directory.

Schemata is written in Cue, which comes with a little learning curve; but once you're comfortable with it, it's extremely powerful.

## Schema Metadata

Lets add our first schema to the `schemata` directory by creating a new file called `artist.cue`.

```cue
# schemata/artist.cue
{
}
```

In order for the schema to be loaded correctly by CueBlox, we need to provide a little boilerplate that helps identify how to work with the schema.

```cue
# schemata/artist.cue
{
	_schema: {
		name:      "Artist"
		namespace: "schemas.cueblox.com"
	}
```

This metadata provides enough information for CueBlox to begin scanning your schema file for DataSets.

## DataSet Metadata

DataSets also need some metadata.

```cue
# schemata/artist.cue
{
    # Same as above, omitted for brevity
    _schema: { ... }

	#Artist: {
		_dataset: {
			plural: "artists"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}
    }
}
```

The `#Artist` key is called a "Definition" in Cue. Definitions are structures that define constraints for the data associated with them. You'll recognise definitions in Cue by the `#` prefix.

The `plural` value is rather important, as this is the name of a directory within your data directory that contains the data to be loaded and validated against this DataSet.

The `supportedExtensions` key allows you to provide a list of file extensions to attempt and load the data from. YAML and Markdown are supported at this time.

## DataSet Constraints

That's all the boilerplate. Now we can define what our `Artist` structure should look like. This definition will be used to validate all data within the `data/artists` directory.

```cue
# schemata/artist.cue
{
    # Same as above, omitted for brevity
    _schema: { ... }

	#Artist: {
        # Same as above, omitted for brevity
		_dataset: { ... }

        name: string
        url?: string

        hungry: bool | *false
    }
}
```

We've now added our first three fields / property to the Artist definition. Firstly, we've added a mandatory field called `name` and an optional field called `url`. Fields with `?` are optional and can be missing from our data files. The last field, `hungry`, is another CUE construct, this time for default fields. This allows you to have mandatory fields with a sensible default when they're missing. In this instance, we're setting `hungry` to be of type `bool` with a default of `false`. `*` is used to indicate a default value.

### Further Reading

We won't go into more details of CUE itself, but we will encourage you to read the following to understand the full potential of using CUE for validation:

- [Official Docs](https://cuelang.org/docs/usecases/validation/)
- [CUEtorials](https://cuetorials.com/)
