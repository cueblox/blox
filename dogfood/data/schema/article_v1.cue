{
	_schema: {
		name:      "Article"
		namespace: "schemas.cueblox.com"
	}

	#Article: {
		_dataset: {
			plural: "articles"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

		// Usually used for the articles h1 tags
		title:             string @template("My New Article")
		// The Except should be a small description
		excerpt:           string @template("Small Description")
		// Should this article be featured?
		featured:          bool | *false
		// Drafts won't be published on the website
		draft:             bool | *false
		// ISO8601, please
		publish_date:      string @template("2020-01-01")
		image?:            string
		last_edit_date?:   string
		edit_description?: string
		body?:             string @template("My Awesome Article")
		tags?: [...string]
		category_id?: string
		profile_id?:  string
	}

}
