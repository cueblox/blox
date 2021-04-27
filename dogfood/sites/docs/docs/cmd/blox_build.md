## blox build

Validate & Build dataset

### Synopsis

The build command will ensure that your dataset is correct by
	validating it against your schemata. Once validated, it will render all
	your content into a single JSON file, which can be consumed by your tooling
	of choice.
	
	Referential Integrity can be enforced with -i. This ensures that any fields
	ending with _id are valid references to identifiers within the other content type.

```
blox build [flags]
```

### Options

```
  -h, --help                    help for build
  -i, --referential-integrity   Verify referential integrity
```

### Options inherited from parent commands

```
      --debug   enable debug logging, overrides 'quiet' flag
      --quiet   disable logging
```

### SEE ALSO

* [blox](/cmd/blox)	 - CueBlox is a suite of slightly opinionated tools for managing and sharing content repositories of YAML and Markdown documents.

