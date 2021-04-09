{
	_schema: {
		name:      "Section"
		namespace: "schemas.cueblox.com"
	}

	#Section: {
		_model: {
			plural: "sections"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

	name:  string
	description: string
    body?: string
	}


}
