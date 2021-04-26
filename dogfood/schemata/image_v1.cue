{
	_schema: {
		name:      "Image"
		namespace: "schemas.cueblox.com"
	}

	#Image: {
		_dataset: {
			plural: "images"
			supportedExtensions: ["yaml", "yml"]
		}

		file_name:        	string 
		width:      		int
		height:        		int
		alt_text?:			string
		caption?:			string
		attribution?:		string
		attribution_link?:	string
	}

}
