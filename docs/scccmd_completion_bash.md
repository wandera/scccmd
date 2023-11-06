## scccmd completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(scccmd completion bash)

To load completions for every new session, execute once:

#### Linux:

	scccmd completion bash > /etc/bash_completion.d/scccmd

#### macOS:

	scccmd completion bash > $(brew --prefix)/etc/bash_completion.d/scccmd

You will need to start a new shell for this setup to take effect.


```
scccmd completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --log-level string   command log level (options: [panic fatal error warning info debug trace]) (default "info")
```

### SEE ALSO

* [scccmd completion](scccmd_completion.md)	 - Generate the autocompletion script for the specified shell

