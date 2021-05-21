# blox completion

Prints shell autocompletion scripts for blox

## Synopsis

Allows you to setup your shell to completions blox commands and flags.

### Bash

	$ source <(blox completion bash)

To load completions for each session, execute once:

#### Linux

	$ blox completion bash > /etc/bash_completion.d/blox

#### MacOS

	$ blox completion bash > /usr/local/etc/bash_completion.d/blox

### ZSH

If shell completion is not already enabled in your environment you will need to enable it.
You can execute the following once:

	$ echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions for each session, execute once:

	$ blox completion zsh > "${fpath[1]}/_blox"

You will need to start a new shell for this setup to take effect.

### Fish

	$ blox completion fish | source

To load completions for each session, execute once:

	$ blox completion fish > ~/.config/fish/completions/blox.fish

**NOTE**: If you are using an official blox package, it should setup completions for you out of the box.


```
blox completion [bash|zsh|fish]
```

## Options

```
  -h, --help   help for completion
```

## Options inherited from parent commands

```
      --debug   enable debug logging, overrides 'quiet' flag
      --quiet   disable logging
```

## See also

* [blox](/cmd/blox)	 - CueBlox is a suite of slightly opinionated tools for managing and sharing content repositories of YAML and Markdown documents.

