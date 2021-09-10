# blox render

Render templates with compiled data

## Synopsis

Render templates with compiled data. 
Use the 'with' parameter to restrict the data set to a single content type.
Use the 'each' parameter to execute the template once for each item.

```
blox render [flags]
```

## Options

```
  -e, --each              render template once per item
  -h, --help              help for render
  -t, --template string   template to render
  -w, --with string       dataset to use
```

## Options inherited from parent commands

```
      --debug   enable debug logging, overrides 'quiet' flag
      --quiet   disable logging
```

## See also

* [blox](/cmd/blox)	 - CueBlox is a suite of slightly opinionated tools for managing and sharing content repositories of YAML and Markdown documents.

