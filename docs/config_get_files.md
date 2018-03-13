## config get files

Get the config files from the given config server

### Synopsis


Get the config files from the given config server

```
config get files [flags]
```

### Options

```
  -f, --files FileMappings   files to get in form of source:destination pairs, example '--files application.yaml:config.yaml'
  -h, --help                 help for files
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
