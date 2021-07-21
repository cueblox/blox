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
		section?: string
		weight?:     int
	}

}
