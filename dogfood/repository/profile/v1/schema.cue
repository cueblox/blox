{
	_schema: {
		name:      "Profile"
		namespace: "schemas.cueblox.com"
	}

	#Profile: {
		_model: {
			plural: "profiles"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

		first_name: string
		last_name:  string
		age?:       int
		company?:   string
		title?:     string
		body?:      string
		social_accounts?: [...#TwitterAccount | #GitHubAccount | #MiscellaneousAccount]
	}

	#TwitterAccount: {
		network:  "twitter"
		username: string
		url:      *"https://twitter.com/\(username)" | string
	}

	#GitHubAccount: {
		network:  "github"
		username: string
		url:      *"https://github.com/\(username)" | string
	}

	#MiscellaneousAccount: {
		network: string
		url:     string
	}
}
