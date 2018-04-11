## scccmd webhook

Runs K8s webhook for injecting config from Cloud Config Server

### Synopsis

Runs K8s webhook for injecting config from Cloud Config Server

```
scccmd webhook [flags]
```

### Options

```
  -c, --cert-file string     Location of public part of SSL certificate (default "keys/publickey.cer")
  -f, --config-file string   The configuration namespace (default "config/config.yaml")
  -h, --help                 help for webhook
  -k, --key-file string      Location of private key of SSL certificate (default "keys/private.key")
  -p, --port int             Webhook port (default 443)
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO

* [scccmd](scccmd.md)	 - Spring Cloud Config management tool

