---
title: Introduction
excerpt: What is CueBlox
publish_date: "2021-03-19"
section: introduction
weight: 1
---

# Introduction

CueBlox is a suite of slightly opinionated tools for managing and sharing content repositories of YAML and Markdown documents.

## Motivation

We built CueBlox out of a desire to share data in a group or team setting without worrying about incompatible metadata. Markdown (with YAML frontmatter) and just plain YAML are great tools for content authors, but we were missing a reliable way of ensuring the content created by different people -- possibly in different repositories -- has the same metadata. Agreeing to use a common template for your frontmatter is a good start, but even that falls short when inconsistencies in metadata values creep in.

With data consistency at the core, our next concern was making the data easily accessible. We explored methods of serving content repositories over GraphQL and REST by trying to build a complicated server that syncs remote git repositories to server temp storage, parses and validates them, then serves the combined data over GraphQL. While technically possible, there were many moving parts and it felt like too much manual work for a task that should be relatively simple.

Finally, we wanted to make it trivial to make a content repository available for consumption publicly or privately over industry-standard protocols. You shouldn't have to install custom tools to consume data. It should be available in standard encoding formats over standard transfer protocols.

## The CueBlox Solution

CueBlox provides a scaffold system that allows you to create a content repository with well-defined schemata. Schema definitions are stored with the content so that any consumer of the content has the information required to validate and parse the content.

Schema definitions are written in [Cue](https://cuelang.org), providing a powerful method of defining field-level requirements and providing default values for fields.

CueBlox also supports a convention-based relational data system by automatically linking and validating foreign keys which follow basic naming conventions. A document can reference a document of another type by adding a foreign key field to the metadata. Much like the relational database concepts on which this is based, CueBlox will validate that the record referenced in your foreign key field exists. This allows documents of different schemas to reference relational data in a familiar format without having to specify relationships in configuration files.

A concrete example of this is the typical blog post. A blog post written in Markdown might have YAML metadata in the frontmatter that describes the post and provides additional information which is used to format and display the post. Here's a the frontmatter for the page you're reading now:

```yaml
title: Introduction
excerpt: What is CueBlox
publish_date: "2021-03-19"
section_id: introduction
weight: 1
```

This frontmatter lives in a file called `introduction.md` in the `pages` folder of our content repository. The `section_id` field is a foreign key reference to a document in the `sections` folder which is named `introduction.md`. In relational database terms, the Introduction page belongs to the Introduction section. `Sections` have many `Pages`, `Pages` belong to `Sections`. When validating the content, CueBlox will ensure that the `Section` referenced actually exists, and throw an error if it doesn't.

Field-level validations and relational validations ensure that your data exists in the shape you intended, without common errors. Most publishing platforms will happily allow you to have frontmatter fields that are not known to the system consuming them. If you mis-spell `excerpt` in your frontmatter, you're not likely to know it until the publishing system gets the data and there's no summary in the listing page. CueBlox prevents this by giving you the ability to define required and optional fields, and allowing you to "close" a definition, preventing extra fields from being added.

## Beyond Validation

Providing schema and relational validation is the foundation of CueBlox, but it only solves the consistency issue. The next layer of CueBlox takes a content repository and assembles it into a JSON object that can be consumed anywhere. You can assemble your content into JSON locally next the applications that consume it, or you can use other tools to serve that data over the Internet.

CueBlox provides pre-built examples that allow you to build automated content deployment pipelines hosted on your favorite cloud. We've built easy to integrate GitHub actions to automate validation and deployment, as well. In a matter of minutes you can publish a content repository over GraphQL and REST APIs with automatic deployments on all the major hosting providers. CueBlox was built to be deployed inexpensively, and it's likely that your favorite hosting provider has a free tier of hosting that will be more than sufficient for CueBlox's minimal hosting requirements.

## Teamwork Makes the Dream Work

A fundamental problem we set out to solve was collaboration across teams of any size. CueBlox enables this by allowing teams to create and publish their own Schemata, including built-in support for versioning. With shared schemata, your team can be confident that everyone will create data that can be consumed and aggregated without manual intervention.

## Beyond The Blog

While the toolset is well suited for managing markdown content that you intend to publish, that's far from the only use case we envisioned. One of our primary goals was to allow a git-based workflow for other types of data. As a concrete example, the primary authors of CueBlox are in Developer Relations. We intend to create shared schemata for our teams to allow us to share information about events, speaking engagements, conferences and other types of data that might be aggregated to create a team calendar, or provide metrics for management. Because CueBlox is markdown or YAML stored in a git repository, we can even enable git workflows for approval. Imagine creating a travel request by submitting a PR to a shared team content repository. Approval and merging automatically create the data about the event and travel costs. Your workflows are only limited by your imagination, but CueBlox enables everyone to create and consume the data.
