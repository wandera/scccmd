## config init

Initialize the config from the given config server

### Synopsis


Initialize the config from the given config server

```
config init [flags]
```

### Options

```
  -a, --application string   name of the application to get the config for
  -f, --files fileMappings   files to get in form of source:destination pairs, example '--files application.yaml:config.yaml'
  -h, --help                 help for init
  -l, --label string         configuration label (default "master")
  -p, --profile string       configuration profile (default "default")
  -s, --source string        address of the config server
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO
* [config](config.md)	 - Spring Cloud Config management tool

