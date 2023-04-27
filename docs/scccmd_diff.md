## scccmd diff

Diff the config from the given config server

### Synopsis

Diff the config from the given config server

### Options

```
  -a, --application string      name of the application to get the config for
  -h, --help                    help for diff
      --label string            configuration label (default "master")
      --profile string          configuration profile (default "default")
  -s, --source string           address of the config server
      --target-label string     second label to diff with
      --target-profile string   second profile to diff with, --profile value will be used, if not defined
```

### Options inherited from parent commands

```
      --log-level string   command log level (options: [panic fatal error warning info debug]) (default "info")
```

### SEE ALSO

* [scccmd](scccmd.md)	 - Spring Cloud Config management tool
* [scccmd diff files](scccmd_diff_files.md)	 - Diff the config files from the given config server
* [scccmd diff values](scccmd_diff_values.md)	 - Diff the config values in specified format from the given config server

