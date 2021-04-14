---
title: Home
weight: 0
excerpt: CueBlox documentation
publish_date: "2021-03-19"
---

# CueBlox

CueBlox is a set of tools that allow you to create and consume datasets from YAML or Markdown files.

At the core is a tool that aggregates similar files into collections of data. If you've ever built a website with a static site generator like Hugo, this will be familiar.

Where CueBlox really shines though is in the additional functionality it enables. CueBlox has several features that enable some interesting and novel integrations for your data:

- Validate your datasets against a schema.

  Ensure your data is always valid by providing defaults and validation rules using the [Cue](https://cuelang.org) language to define schemata for your data.

  The schema that validates this page looks like this:

  ```
  	#Page: {
  	_dataset: {
  		plural: "pages"
  		supportedExtensions: ["yaml", "yml", "md", "mdx"]
  	}

  	title:        string
  	excerpt:      string
  	draft:        bool | *false
  	publish_date: string
  	image?:       string
  	body?:        string
  	tags?: [...string]
  	section_id?: string
  	weight?:     int
  }
  ```

  In this schema we've defined required fields and optional fields by adding or omitting the `?` in the field definition. A `Page` must have a `title`, but `weight` is optional.

- Aggregate your datasets and export them into a JSON file

  Data that is locked in your git repository is only useful in that repository. CueBlox allows you to validate, aggregate, and export your data in the integration-friendly JSON format for consumption elsewhere. The data that drives this website is processed through CueBlox and automatically published as a GitHub release. You can see it [here](https://github.com/cueblox/blox/releases/tag/blox)

- Leverage third party tools to make your data easily accessible

  CueBlox includes a `hosting` command that generates a website suitable for deployment on Azure Static Web Apps, Vercel, or Netlify. The site includes a serverless function powered by [json-graphql-server](https://github.com/marmelab/json-graphql-server). The serverless function pulls the data from a URL (the GitHub release mentioned above) and serves it over GraphQL as a read-only dataset. By following a [few conventions](https://github.com/marmelab/json-graphql-server#generated-types-and-queries), we even get real referential integrity between our markdown files. You can play with the GraphQL server that powers this site [here](https://api.cueblox.com/api/graphql). Click on the "Docs" link at the top and see all the wonderful GraphQL data that is generated from our directories of Markdown files.

- Create and consume standardized schemata for you and your team

  The schemata you create are available to you locally in your content directory. But you can also create a set of schemata that is published for others to consume. The schema we use to build all of the CueBlox website are published on [this website](https://schemas.cueblox.com/) All the information about the schemata available, including version information is available in the `manifest.json` file linked in the HTML.

  The `blox` cli tool allows you to add these remote schemata to your project:

  ```
  ❯ blox remote list schemas.cueblox.com
  Namespace           | Schema   | Version
  schemas.cueblox.com | article  | v1
  schemas.cueblox.com | category | v1
  schemas.cueblox.com | page     | v1
  schemas.cueblox.com | profile  | v1
  schemas.cueblox.com | section  | v1
  schemas.cueblox.com | website  | v1

  ❯ blox remote get schemas.cueblox.com page v1
  ```

  Using published schemata this way allows you to create and consume standardized data across several projects, teams, and people. In fact, it is the core driver for our development of CueBlox -- enabling teams to publish information individually, but have it aggregated and consumed at a higher level without data validation worries. When everyone uses the same schema, all the data is always consistent.

- Leverage all of this functionality using GitOps principles

  CueBlox was built for lazy people. If it can be generated, we generate it. If it can be automated, we automate it. The end result is that the only thing you need to do to publish your content is check it into your git repository. GitHub actions take care of the rest!
