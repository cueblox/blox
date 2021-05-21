# blox build

Validate & Build dataset

## Synopsis

The build command will ensure that your dataset is correct by
	validating it against your schemata. Once validated, it will render all
	your content into a single JSON file, which can be consumed by your tooling
	of choice.
	
	Referential Integrity can be enforced with -i. This ensures that any fields
	ending with _id are valid references to identifiers within the other content type.
	
	The build process will create an 'image' record for images in your 'static_dir' if you use the -g or --images flag.
	
	Images will be pushed to blob storage if you use -s/--sync and set the appropriate environment variables. 
	Currently only Azure blob storate is implemented. See https://gocloud.dev/howto/blob/#services for required 
	environment variables and setup information.
	

```
blox build [flags]
```

## Options

```
  -h, --help                    help for build
  -g, --images                  Create 'image' records for images in static directory
  -i, --referential-integrity   Verify referential integrity
  -s, --sync                    Sync images to blob storage
```

## Options inherited from parent commands

```
      --debug   enable debug logging, overrides 'quiet' flag
      --quiet   disable logging
```

## See also

* [blox](/cmd/blox)	 - CueBlox is a suite of slightly opinionated tools for managing and sharing content repositories of YAML and Markdown documents.

