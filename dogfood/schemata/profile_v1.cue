{
	_schema: {
		name:      "Profile"
		namespace: "schemas.cueblox.com"
	}

	#Profile: {
		_dataset: {
			plural: "profiles"
			supportedExtensions: ["yaml", "yml", "md", "mdx"]
		}

		first_name: string @template("Forename")
		last_name:  string @template("Surname")
		age?:       int    @template(21)
		company?:   string @template("CueBlox")
		title?:     string @template("Cue Slinger")
		body?:      string @template("☕️ Required")
		image?:  string
		social_accounts?: [...#TwitterAccount | #GitHubAccount | #MiscellaneousAccount]
	}

	#TwitterAccount: {
		network:  "twitter"
		username: string @template("twitter-handle")
		url:      *"https://twitter.com/\(username)" | string
	}

	#GitHubAccount: {
		network:  "github"
		username: string @template("github-handle")
		url:      *"https://github.com/\(username)" | string
	}

	#MiscellaneousAccount: {
		network: string @template("some_network")
		url:     string @template("https://some_url")
	}
}
