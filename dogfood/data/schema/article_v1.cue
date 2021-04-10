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

		title:             string
		excerpt:           string
		featured:          bool | *false
		draft:             bool | *false
		publish_date:      string
		image?:            string
		last_edit_date?:   string
		edit_description?: string
		body?:             string
		tags?: [...string]
		category_id?: string
		profile_id?:  string
	}

}
