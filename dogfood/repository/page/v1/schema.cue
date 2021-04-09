{
	_schema: {
		name:      "Page"
		namespace: "schemas.cueblox.com"
	}

	#Page: {
		_model: {
			plural: "pages"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

	title: string
	excerpt:  string
    draft: bool | *false
    publish_date: string
    image?: string
	body?:      string
	tags?: [...string]
	section_id?: string
	}

}
