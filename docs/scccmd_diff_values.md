## scccmd diff values

Diff the config values in specified format from the given config server

### Synopsis

Diff the config values in specified format from the given config server

```
scccmd diff values [flags]
```

### Options

```
  -f, --format string   output format might be one of 'json|yaml|properties' (default "yaml")
  -h, --help            help for values
```

### Options inherited from parent commands

```
  -a, --application string      name of the application to get the config for
      --label string            configuration label (default "master")
      --log-level string        command log level (options: [panic fatal error warning info debug trace]) (default "info")
      --profile string          configuration profile (default "default")
  -s, --source string           address of the config server
      --target-label string     second label to diff with
      --target-profile string   second profile to diff with, --profile value will be used, if not defined
```

### SEE ALSO

* [scccmd diff](scccmd_diff.md)	 - Diff the config from the given config server

