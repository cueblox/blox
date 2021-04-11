{
	_schema: {
		name:      "Profile"
		namespace: "schemas.cueblox.com"
	}

	#Website: {
		_dataset: {
			plural: "websites"
			supportedExtensions: ["yaml", "yml"]
		}

		url:         string
		profile_id?: string
		body?:       string
	}
}
