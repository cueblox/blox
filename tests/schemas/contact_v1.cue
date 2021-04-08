{
	_schema: {
		name:      "Contact"
		namespace: "schemas.cueblox.com"
	}

	#Contact: {
		_model: {
			plural: "contacts"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]

		}

		name:    #Name
		address: #Address
		phone:   string
		email:   string
	}

	#Name: {
		forename: string
		surname:  string
	}

	#Address: {
		number:   string
		street:   string
		city:     string
		country:  string
		postcode: string
	}
}
