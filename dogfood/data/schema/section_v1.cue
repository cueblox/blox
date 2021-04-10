{
	_schema: {
		name:      "Section"
		namespace: "schemas.cueblox.com"
	}

	#Section: {
		_dataset: {
			plural: "sections"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

		name:        string
		description: string
		body?:       string
		weight?:     int
	}

}
