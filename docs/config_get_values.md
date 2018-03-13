## config get values

Get the config values in specified format from the given config server

### Synopsis


Get the config values in specified format from the given config server

```
config get values [flags]
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
  -p, --profile string       configuration profile (default "default")
  -s, --source string        address of the config server
  -v, --verbose              verbose output
```

### SEE ALSO
* [config get](config_get.md)	 - Get the config from the given config server

