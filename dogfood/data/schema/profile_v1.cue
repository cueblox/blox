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

		// First name, please
		first_name: string @template("Forename")
		// WHat's your last name?
		last_name:  string @template("Surname")
		// Yes or No, huh?
		yesno: bool @template(false)
		// How old are you even? You better be 21
		age?:       int    @template(21)
		company?:   string @template("CueBlox")
		title?:     string @template("Cue Slinger")
		body?:      string @template("☕️ Required")
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
