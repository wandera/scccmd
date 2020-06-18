## scccmd diff files

Diff the config files from the given config server

### Synopsis

Diff the config files from the given config server

```
scccmd diff files [flags]
```

### Options

```
  -f, --files string   files to get in form of file1,file2, example '--files application.yaml,config.yaml'
  -h, --help           help for files
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

