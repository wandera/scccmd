## scccmd webhook

Runs K8s webhook for injecting config from Cloud Config Server

```
scccmd webhook [flags]
```

### Options

```
  -c, --cert-file string     location of public part of SSL certificate (default "keys/publickey.cer")
  -f, --config-file string   the configuration file (default "config/config.yaml")
  -h, --help                 help for webhook
  -k, --key-file string      location of private key of SSL certificate (default "keys/private.key")
  -p, --port int             webhook port (default 443)
```

### Options inherited from parent commands

```
      --log-level string   command log level (options: [panic fatal error warning info debug trace]) (default "info")
```

### SEE ALSO

* [scccmd](scccmd.md)	 - Spring Cloud Config management tool

