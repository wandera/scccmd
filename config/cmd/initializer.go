package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/api/apps/v1"
	"k8s.io/client-go/rest"
	"os"
	"os/signal"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type config struct {
	AnnotationPrefix     string `yaml:"annotation-prefix"`
	ContainerImage       string `yaml:"container-image"`
	DefaultContainerName string `yaml:"default-container-name"`
	DefaultLabel         string `yaml:"default-label"`
	DefaultProfile       string `yaml:"default-profile"`
	DefaultVolumeName    string `yaml:"default-volume-name"`
	DefaultVolumeMount   string `yaml:"default-volume-mount"`
	DefaultSource        string `yaml:"default-source"`
}

type dynamicConfig struct {
	containerName string
	volumeName    string
	volumeMount   string
	imageArgs     []string
}

var ip = struct {
	configmap       string
	initializerName string
	namespace       string
	kubeconfig      bool
}{}

var initializerCmd = &cobra.Command{
	Use:   "initializer",
	Short: "Runs K8s initializer for injecting config from Cloud Config Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeInitializer(args)
	},
}

func executeInitializer(args []string) error {
	fmt.Println("Starting the Config initializer...")

	if Verbose {
		fmt.Printf("Initializer name set to: %s \n", ip.initializerName)
	}

	clusterConfig, err := restConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return err
	}

	// Load the Envoy Initializer configuration from a Kubernetes ConfigMap.
	cm, err := clientset.CoreV1().ConfigMaps(ip.namespace).Get(ip.configmap, metav1.GetOptions{})
	if err != nil {
		return err
	}

	c, err := configmapToConfig(cm)
	if err != nil {
		return err
	}

	// Watch uninitialized Deployments in all namespaces.
	restClient := clientset.AppsV1().RESTClient()
	watchlist := cache.NewListWatchFromClient(restClient, "deployments", corev1.NamespaceAll, fields.Everything())

	// Wrap the returned watchlist to workaround the inability to include
	// the `IncludeUninitialized` list option when setting up watch clients.
	includeUninitializedWatchlist := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.IncludeUninitialized = true
			return watchlist.List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.IncludeUninitialized = true
			return watchlist.Watch(options)
		},
	}

	syncPeriod := 30 * time.Second

	_, controller := cache.NewInformer(includeUninitializedWatchlist, &v1.Deployment{}, syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				err := initializeDeployment(obj.(*v1.Deployment), c, clientset)
				if err != nil {
					fmt.Println(err)
				}
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Config initializer started")

	<-signalChan

	fmt.Println("Shutdown signal received, exiting...")
	close(stop)

	return nil
}

func initializeDeployment(deployment *v1.Deployment, c *config, clientset *kubernetes.Clientset) error {
	if deployment.ObjectMeta.GetInitializers() != nil {
		pendingInitializers := deployment.ObjectMeta.GetInitializers().Pending

		if ip.initializerName == pendingInitializers[0].Name {
			fmt.Printf("Checking deployment '%s' \n", deployment.Name)

			initializedDeployment := deployment.DeepCopy()

			// Remove self from the list of pending Initializers while preserving ordering.
			if len(pendingInitializers) == 1 {
				initializedDeployment.ObjectMeta.Initializers = nil
			} else {
				initializedDeployment.ObjectMeta.Initializers.Pending = append(pendingInitializers[:0], pendingInitializers[1:]...)
			}

			a := deployment.ObjectMeta.GetAnnotations()
			ea := c.AnnotationPrefix + "enabled"
			_, ok := a[ea]
			if !ok {
				if Verbose {
					fmt.Printf("Required '%s' annotation missing skipping injection \n", ea)
				}
				_, err := clientset.AppsV1().Deployments(deployment.Namespace).Update(initializedDeployment)
				if err != nil {
					return err
				}
				return nil
			}

			if Verbose {
				fmt.Printf("Patching deployment '%s'\n", deployment.Name)
			}

			d, err := calculateDynamicConfig(c, a, deployment)
			if err != nil {
				return err
			}

			volumeMount := corev1.VolumeMount{
				Name:      d.volumeName,
				MountPath: d.volumeMount,
			}

			// Add InitContainer
			initializedDeployment.Spec.Template.Spec.InitContainers = append(deployment.Spec.Template.Spec.InitContainers, corev1.Container{
				Name:         d.containerName,
				Image:        c.ContainerImage,
				Args:         d.imageArgs,
				VolumeMounts: []corev1.VolumeMount{volumeMount},
			})

			// Add volume mount to the first regular container
			initializedDeployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(initializedDeployment.Spec.Template.Spec.Containers[0].VolumeMounts, volumeMount)

			// Add volume definition
			initializedDeployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
				Name:         volumeMount.Name,
				VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
			})

			oldData, err := json.Marshal(deployment)
			if err != nil {
				return err
			}

			newData, err := json.Marshal(initializedDeployment)
			if err != nil {
				return err
			}

			patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.Deployment{})
			if err != nil {
				return err
			}

			_, err = clientset.AppsV1().Deployments(deployment.Namespace).Patch(deployment.Name, types.StrategicMergePatchType, patchBytes)
			if err != nil {
				return err
			}

			if Verbose {
				fmt.Printf("Deployment '%s' successfully patched \n", deployment.Name)
			}
		}
	}

	return nil
}

func calculateDynamicConfig(c *config, a map[string]string, deployment *v1.Deployment) (*dynamicConfig, error) {
	var d = dynamicConfig{}
	var ok bool

	if d.containerName, ok = a[c.AnnotationPrefix+"container-name"]; !ok {
		d.containerName = c.DefaultContainerName
	}

	if d.volumeName, ok = a[c.AnnotationPrefix+"volume-name"]; !ok {
		d.volumeName = c.DefaultVolumeName
	}

	if d.volumeMount, ok = a[c.AnnotationPrefix+"volume-mount"]; !ok {
		d.volumeMount = c.DefaultVolumeMount
	}

	imageArgs, err := calculateImageArgs(c, a, deployment)
	d.imageArgs = imageArgs
	return &d, err
}

func calculateImageArgs(c *config, a map[string]string, deployment *v1.Deployment) ([]string, error) {
	var ok bool
	var mode string
	var application string
	var profile string
	var label string
	var source string
	var extra []string

	if mapping, ok := a[c.AnnotationPrefix+"mapping"]; ok {
		mode = "files"
		extra = append(extra, "--files", mapping)
	} else if destination, ok := a[c.AnnotationPrefix+"destination"]; ok {
		mode = "values"
		extra = append(extra, "--destination", destination)
	} else {
		return nil, fmt.Errorf("one of '%s' or '%s' annotations should be specified", c.AnnotationPrefix+"mapping", c.AnnotationPrefix+"destination")
	}

	if source, ok = a[c.AnnotationPrefix+"source"]; !ok {
		source = c.DefaultSource
	}

	if profile, ok = a[c.AnnotationPrefix+"profile"]; !ok {
		profile = c.DefaultProfile
	}

	if label, ok = a[c.AnnotationPrefix+"label"]; !ok {
		label = c.DefaultLabel
	}

	if application, ok = a[c.AnnotationPrefix+"application"]; !ok {
		application = deployment.Name
	}

	if Verbose {
		extra = append(extra, "-v")
	}

	return append([]string{"get", mode, "--source", source, "--application", application, "--profile", profile, "--label", label}, extra...), nil
}

func configmapToConfig(configmap *corev1.ConfigMap) (*config, error) {
	var c config
	confAsYaml, err := yaml.Marshal(configmap.Data)
	err = yaml.Unmarshal(confAsYaml, &c)
	if err != nil {
		return nil, err
	}

	defaultIfEmpty(&c.AnnotationPrefix, "config.initializer.kubernetes.io/")
	defaultIfEmpty(&c.DefaultContainerName, "config-init")
	defaultIfEmpty(&c.DefaultLabel, "master")
	defaultIfEmpty(&c.DefaultProfile, "default")
	defaultIfEmpty(&c.DefaultSource, "http://config-manager-controller.config.svc:8080")
	defaultIfEmpty(&c.DefaultVolumeName, "config-volume")
	defaultIfEmpty(&c.DefaultVolumeMount, "/config")

	return &c, nil
}

func restConfig() (*rest.Config, error) {
	if ip.kubeconfig {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		if Verbose {
			fmt.Println("Loading kubeconfig cluster configuration")
		}
		return kubeConfig.ClientConfig()
	}
	if Verbose {
		fmt.Println("Loading incluster cluster configuration")
	}
	return rest.InClusterConfig()
}

func defaultIfEmpty(val *string, def string) {
	if *val == "" {
		*val = def
	}
}

func init() {
	initializerCmd.Flags().StringVarP(&ip.namespace, "namespace", "n", "default", "The configuration namespace")
	initializerCmd.Flags().StringVarP(&ip.configmap, "configmap", "c", "initializer-config", "The config initializer configuration configmap")
	initializerCmd.Flags().StringVarP(&ip.initializerName, "initializer-name", "i", "config.initializer.kubernetes.io", "The initializer name")
	initializerCmd.Flags().BoolVar(&ip.kubeconfig, "kubeconfig", false, "If kubeconfig should b e used for connecting to the cluster, mainly for debugging purposes, when false command autodiscover configuration from within the cluster")
}
