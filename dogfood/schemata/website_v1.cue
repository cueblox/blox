{
	_schema: {
		name:      "Website"
		namespace: "schemas.cueblox.com"
	}

	#Website: {
		_dataset: {
			plural: "websites"
			supportedExtensions: ["yaml", "yml"]
		}

		url:       string @template("https://google.com")
		profile?:  string @relationship(Profile)
		body?:     string
	}
}
