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

		title:             string @template("My New Article")
		excerpt:           string @template("Small Description")
		featured:          bool | *false
		draft:             bool | *false
		publish_date:      string @template("2020-01-01")
		image?:         string @relationship(Image)
		last_edit_date?:   string
		edit_description?: string
		body?:             string @template("My Awesome Article")
		tags?: 			   [...string]
		category?: 	   string @relationship(Category)
		profile?:  	   string @relationship(Profile)

	}

}
