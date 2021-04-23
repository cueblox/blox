---
title: Foundations
excerpt: Fundamentals of CueBlox
publish_date: "2021-03-19"
section_id: foundations
weight: 2
---

# Definitions

Let's start by getting on the same page about what the words in these documents mean. Naming is hard!

## Dataset

A `dataset` is a group of related documents that share the same properties. 
The `dataset` that powers this documentation website is called `pages`. Every document in 
the `pages` `dataset` has the same metadata.

## Metadata

Every document stores two different types of information:

* Content
* Information about the content

The content of the document is the data, the words, the information that you create 
to be consumed by an application like a website. You're reading the content of a document called
`foundations.mdx`. 

>Because much of the data that you create will be consumed by a web 
>server, we'll often use the term `body` to refer to the content as well. Consider the 
>two terms interchangeable.

In addition to the content stored in this document, there is additional data about the content.
Data about the data is called `metadata`.  

CueBlox can parse documents stored in two different formats:  [yaml](https://yaml.org/) and [markdown](https://daringfireball.net/projects/markdown/).


`foundations.mdx` is stored in the `mdx` format
which is an enhanced version of markdown that allows the author to include React components.

### YAML
YAML documents are parsed by reading key/value pairs. A simple YAML document might look like this:

```yaml
url: https://www.cueblox.com
```

### Markdown

Markdown documents provide a simplified syntax to format a document. 
You can add metadata to a Markdown document by adding a YAML to the top of the Markdown
file.

```markdown
---
title: Foundations
excerpt: Fundamentals of CueBlox
publish_date: "2021-03-19"
section_id: foundations
weight: 2
---

# Definitions

Let's start by getting on the same page about what the words in these documents mean. Naming is hard!
```
The YAML added to the markdown document that defines the page you are reading is separated
from the body by enclosing it in three dash characters: ```---```. When a markdown document 
includes metadata in this form, it's called `FrontMatter`.

Whether your document is stored in YAML format, or markdown with YAML frontmatter, we refer to 
the information about the document as `metadata`.

## Schema

In order to enforce consistency in a dataset, CueBlox uses a `schema`, which is a set of 
rules, defaults, and metadata about a dataset. We use `schemata` as the plural of `schema`.

The `schema` is defined using [cue](https://cuelang.org), but you don't need to become an 
expert in Cue to start making your own `schema`.

### Example Schema
Here's the schema that defines the `page` dataset that powers the documention you're reading.

```json
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
		section_id?: string
		weight?:     int
	}

}
```

It starts with some metadata about the schema in a field called (appropriately) `_schema`.
The `_schema` field defines the name and namespace of the `schema`. This allows CueBlox to 
find the right `schema` to validate and process your documents.

The next section in the `schema` defines a `Page`. It contains more metadata in the `_dataset` field
which we use to find your content. 

The most important part of the `Page` definition is the list of fields that are defined in a `Page`
document. Each field is defined with a name, a data type, and other optional metadata about 
the field, such as defaults.

Fields that are optional have a `?` at the end of the field name.

Fields definitions can provide a default value by specifying it after the data type:
```json
        draft:        bool | *false
```

The most common data types you will use are `string`, `int`, and `bool`. You can find all
the supported data types in the documentation on [cuelang.org](https://cuelang.org)
