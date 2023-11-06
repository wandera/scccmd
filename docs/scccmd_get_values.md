## scccmd get values

Get the config values in specified format from the given config server

```
scccmd get values [flags]
```

### Options

```
  -d, --destination string   destination file name
  -f, --format string        output format might be one of 'json|yaml|properties' (default "yaml")
  -h, --help                 help for values
```

### Options inherited from parent commands

```
  -a, --application string   name of the application to get the config for
  -l, --label string         configuration label (default "master")
      --log-level string     command log level (options: [panic fatal error warning info debug trace]) (default "info")
  -p, --profile string       configuration profile (default "default")
  -s, --source string        address of the config server
```

### SEE ALSO

* [scccmd get](scccmd_get.md)	 - Get the config from the given config server

