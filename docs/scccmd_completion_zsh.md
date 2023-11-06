## scccmd completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(scccmd completion zsh)

To load completions for every new session, execute once:

#### Linux:

	scccmd completion zsh > "${fpath[1]}/_scccmd"

#### macOS:

	scccmd completion zsh > $(brew --prefix)/share/zsh/site-functions/_scccmd

You will need to start a new shell for this setup to take effect.


```
scccmd completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --log-level string   command log level (options: [panic fatal error warning info debug trace]) (default "info")
```

### SEE ALSO

* [scccmd completion](scccmd_completion.md)	 - Generate the autocompletion script for the specified shell

