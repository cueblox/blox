{
	_schema: {
		name:      "Category"
		namespace: "schemas.cueblox.com"
	}

	#Category: {
		_dataset: {
			plural: "categories"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

		name:        string @template("Name")
		description: string @template("Description")
		image?: 	 string
		body?:       string @template("This is my category for ...")
	}

}
