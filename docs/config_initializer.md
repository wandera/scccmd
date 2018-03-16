## config initializer

Runs K8s initializer for injecting config from Cloud Config Server

### Synopsis


Runs K8s initializer for injecting config from Cloud Config Server

```
config initializer [flags]
```

### Options

```
  -c, --configmap string          The config initializer configuration configmap (default "initializer-config")
  -h, --help                      help for initializer
  -i, --initializer-name string   The initializer name (default "config.initializer.kubernetes.io")
      --kubeconfig                If kubeconfig should b e used for connecting to the cluster, mainly for debugging purposes, when false command autodiscover configuration from within the cluster
  -n, --namespace string          The configuration namespace (default "default")
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO
* [config](config.md)	 - Spring Cloud Config management tool

