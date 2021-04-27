## blox init

Create folders and configuration to maintain content with the blox toolset

### Synopsis

Create a group of folders to store your content. A directory for your data,
	schemata, and build output will be created.

```
blox init [flags]
```

### Options

```
  -b, --build string      where post-processed content will be stored (output json) (default "_build")
  -d, --data string       where pre-processed content will be stored (source markdown or yaml) (default "data")
  -h, --help              help for init
  -s, --schemata string   where the schemata will be stored (default "schemata")
  -c, --skip              don't write a configuration file
  -t, --starter string    use a pre-defined starter in the CURRENT directory
  -a, --static string     where the static originals will be found (default "static")
```

### Options inherited from parent commands

```
      --debug   enable debug logging, overrides 'quiet' flag
      --quiet   disable logging
```

### SEE ALSO

* [blox](/cmd/blox)	 - CueBlox is a suite of slightly opinionated tools for managing and sharing content repositories of YAML and Markdown documents.

