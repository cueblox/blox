---
title: Quick Start
excerpt: Get up and running quickly
publish_date: "2021-03-19"
section: getting-started
weight: 1
---

# Quick Start

You can install `blox` the CueBlox cli by downloading a release from our [releases](https://github.com/cueblox/blox/releases) page on GitHub. Binaries are available for Mac and Linux in both ARM and AMD64 variants, and for Windows in AMD64 only. After downloading a bundle, extract the binary and put it somewhere in your path.

## Using a Package Manager

`blox` is also available in [Homebrew](https://brew.sh) as a cask. Add the CueBlox tap:

```bash
brew tap cueblox/tap
```

Then install the `brew` formula:

```bash
brew install blox
```

Support for other package managers is in our roadmap, if you'd like to contribute, see [issue 55](https://github.com/cueblox/blox/issues/55).

## Using the `blox` Tool

The `blox` tool has several commands to help you manage your content. To get started you'll want to create a new content repository. Create a new directory for your content and navigate into it in your terminal:

```bash
mkdir democontent
cd democontent
```

Now use the `blox init` command to create a scaffolded content repository.

```
blox init
```

This creates a set of folders and a configuration file.

```
drwxr-xr-x   2 bjk  staff    64 Apr 14 07:02 _build
-rw-r--r--   1 bjk  staff    71 Apr 14 07:02 blox.cue
drwxr-xr-x   2 bjk  staff    64 Apr 14 07:02 schemata
drwxr-xr-x   2 bjk  staff    64 Apr 14 07:02 content
drwxr-xr-x   2 bjk  staff    64 Apr 14 07:02 static
```

You can control the names and locations of these folders by using flags on the `init` command:

```
blox init --help
Create a group of folders to store your content. A directory for your data,
schemata, and build output will be created.

Usage:
  blox init [flags]

Flags:
  -d, --destination string   where post-processed content will be stored (output json) (default "_build")
  -e, --extension string     default file extension for new content (default ".md")
  -h, --help                 help for init
  -a, --schemata string      where schema definitions will be stored (default "schemata")
  -c, --skip                 skip creation of a configuration file
  -s, --source string        where pre-processed content will be stored (source markdown) (default "content")
  -k, --static string        where static files will be stored (default "static")
```

## Directory Layout

The `blox init` command created four folders for you. If you used the defaults, they are `_build`, `schema`, `source` and `static`.

The `_build` folder is where the `blox` command will assemble your content into a JSON file. Since it is generated, and not source content, you may choose to ignore this folder, or the `data.json` file in this folder and depend on your build system to regenerate this file. That is the recommended approach to ensure that your content is always validated and up to date.

The `source` folder is the top level of your content store. Folders under `source` will represent unique content types.

The `schemata` folder is where the definitions of your content types are stored. Each file in the `schemata` folder defines one or more content types that can be created in your content store.

The `static` folder is provided as a convenience for you to store static assets that belong with the content. You can use this static directory to store images and other assets, but it is your responsibility to serve these assets in a way that they can be reached. One possibility is to serve the static directory from the same application you use to serve GraphQL and REST endpoints for your content.

## Your First Schema

Schema definitions are written in [Cue](https://cuelang.org), which provides a powerful syntax for validations and defaults in your schemata. Let's create your first schema by defining a `Page` content type:

```cue
{
	_schema: {
		name:      "Page"
		namespace: "schemas.cueblox.com"
	}

	#Page: {
		_dataset: {
			plural: "pages"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

		title:        string @template("My New Page")
		excerpt:      string @template("Small description about my page")
		draft:        bool | *false
		publish_date: string @template("2020-01-01")
		image?:       string
		body?:        string
		tags?: [...string]
		weight?:     int
	}

}
```

There's a bit to unwrap here if you're new to Cue definitions. For an introduction, read the [Cue documentation](https://cuelang.org/docs/concepts/). We'll only be using a small subset of the features of Cue for this schema, so let's walk through the definition starting at the top.

Cue is a superset of JSON, so Cue definitions should look familiar to you if you've used JSON before. Our `Page` schema starts with some metadata required by CueBlox:

```
	_schema: {
		name:      "Page"
		namespace: "schemas.cueblox.com"
	}
```

This block of data defines the namespace of your schema, which allows CueBlox to help prevent naming collisions in your schemata.

The second block of data defines a model called `#Page`:

```
    #Page: {
		_dataset: {
			plural: "pages"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

		title:        string @template("My New Page")
		excerpt:      string @template("Small description about my page")
		draft:        bool | *false
		publish_date: string @template("2020-01-01")
		image?:       string
		body?:        string
		tags?: [...string]
		weight?:     int
	}
```

Inside the model definition is another set of required metadata, the `_dataset` block. `_dataset` has two required fields: `plural`, and `supportedExtensions`.

`plural` is the name used for resolving your content on the filesystem. A `Page` model has a `plural` definition of `pages`, so CueBlox will search for `Page` content in the `pages` subdirectory of `source`, which is defined in your `blox.cue` configuration file.

`supportedExtensions` tells CueBlox which files should be included in processing. You can include any file extension that you want to be treated as YAML or Markdown. This allows you to use `mdx` or `svx` or any future tool that looks like markdown but includes embedded components.

After the `_dataset` metadata is a list of field definitions for our `#Page` model. Fields are defined with a name and a type, optionally followed by a default value and a `@template` definition which will be used when you create new content.

Our model defines four required fields: `title`, `excerpt`, `draft`, `publish_date`. Optional fields have a question mark `?` at the end of the field name, so fields without the `?` are required.

Default values are optionally specified by a `disjunction` which is written with the pipe operator `|` and a value specified with an asterisk `*`. Read the definition for the `draft` field as `a boolean value, which also has a default value of false`.

The `tags` field specifies an optional array of `string` values, and the `weight` field specifies an optional integer value that we can use to sort our `Page` models in a list.

Copy the `#Page` definition below and save it in your `schemata` directory with a descriptive name like `page.cue`:

```
{
	_schema: {
		name:      "Page"
		namespace: "schemas.cueblox.com"
	}

	#Page: {
		_dataset: {
			plural: "pages"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

		title:        string @template("My New Page")
		excerpt:      string @template("Small description about my page")
		draft:        bool | *false
		publish_date: string @template("2020-01-01")
		image?:       string
		body?:        string
		tags?: [...string]
		weight?:     int
	}

}
```

## Creating Content

With your first schema in place you can begin creating content. Since the `plural` definition of the `#Page` type is `pages`, CueBlox expects your Page content to be in a `pages` subdirectory. You can create a new Markdown file in `content/pages`, or let CueBlox do it for you:

```
blox new --dataset page hello
```

The `blox new` command takes an argument to specify which content type you're creating, and an argument defining the `ID` or slug of your new file. The above command will create `content/pages/hello.yaml` in a default configuration.

The contents of the created `hello` file will include the templated values from the schema definition.

```
title: My New Page
excerpt: Small description about my page
publish_date: "2020-01-01"
```

Now you can edit these fields and try out your first `blox build`.

```
blox build
 INFO  Validating ...
 SUCCESS  Validations complete
```

If everything worked, you'll see no errors and a new data file will be written in your output directory which is `_build` by default:

```
tree .
.
├── _build
│   └── data.json
├── blox.cue
├── content
│   └── pages
│       └── hello.yaml
├── schemata
│   └── page.cue
└── static
```

Let's take a closer look at the data that was output:

```
cat _build/data.json

{
  "pages": [
    {
      "title": "My New Page",
      "excerpt": "Small description about my page",
      "draft": false,
      "publish_date": "2020-01-01",
      "id": "hello"
    }
  ]
}
```

We've formatted the output here to make it easier to read, but the data in the file doesn't have extra whitespace. Notice that we have a JSON object representing a map with a key for each of our content models. The value of that key is an array of content models.

Also notice that CueBlox automatically inserted default values for fields that weren't specified like our `draft` field, and it inserted an `id` field with the filename of the file on disk.

At this point you've generated an output dataset that represents your content which has been validated against a schema you've written with defaults automatically applied. _Congratulations_, you're now ready to dive deeper into the more advanced features of CueBlox!
