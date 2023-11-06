## scccmd completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	scccmd completion fish | source

To load completions for every new session, execute once:

	scccmd completion fish > ~/.config/fish/completions/scccmd.fish

You will need to start a new shell for this setup to take effect.


```
scccmd completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --log-level string   command log level (options: [panic fatal error warning info debug trace]) (default "info")
```

### SEE ALSO

* [scccmd completion](scccmd_completion.md)	 - Generate the autocompletion script for the specified shell

