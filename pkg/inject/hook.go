package inject

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/howeyc/fsnotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

var (
	runtimeScheme     = runtime.NewScheme()
	codecs            = serializer.NewCodecFactory(runtimeScheme)
	deserializer      = codecs.UniversalDeserializer()
	ignoredNamespaces = []string{
		metav1.NamespaceSystem,
		metav1.NamespacePublic,
	}
)

const (
	watchDebounceDelay = 100 * time.Millisecond
)

// WebhookConfigDefaults configures default init container values.
type WebhookConfigDefaults struct {
	ContainerName string `yaml:"container-name"`
	Label         string `yaml:"label"`
	Profile       string `yaml:"profile"`
	VolumeName    string `yaml:"volume-name"`
	VolumeMount   string `yaml:"volume-mount"`
	Source        string `yaml:"source"`
}

// InitContainerResourcesList resources for init container
type InitContainerResourcesList struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// InitContainerResources resources for init container
type InitContainerResources struct {
	Requests InitContainerResourcesList `yaml:"requests"`
	Limits   InitContainerResourcesList `yaml:"limits"`
}

// WebhookConfig struct representing webhook configuration values.
type WebhookConfig struct {
	AnnotationPrefix string                 `yaml:"annotation-prefix"`
	Policy           InjectionPolicy        `yaml:"policy"`
	ContainerImage   string                 `yaml:"container-image"`
	Default          WebhookConfigDefaults  `yaml:"default"`
	Resources        InitContainerResources `yaml:"resources"`
}

// Webhook implements a mutating webhook for automatic config injection.
type Webhook struct {
	mu     sync.RWMutex
	config *WebhookConfig

	healthCheckInterval time.Duration
	healthCheckFile     string

	server     *http.Server
	configFile string
	watcher    *fsnotify.Watcher
	certFile   string
	keyFile    string
	cert       *tls.Certificate
}

// WebhookParameters configures parameters for the config injection webhook.
type WebhookParameters struct {
	// ConfigFile is the path to the injection configuration file.
	ConfigFile string

	// CertFile is the path to the x509 certificate for https.
	CertFile string

	// KeyFile is the path to the x509 private key matching `CertFile`.
	KeyFile string

	// Port is the webhook port, e.g. typically 443 for https.
	Port int
}

// UnmarshalYAML implements Unmarshaler interface for WebhookConfig
func (w *WebhookConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawWebhookConfig WebhookConfig
	raw := rawWebhookConfig{
		Policy:           InjectionPolicyEnabled,
		ContainerImage:   "wanderadock/scccmd",
		AnnotationPrefix: "config.scccmd.github.com/",
		Default: WebhookConfigDefaults{
			ContainerName: "config-init",
			VolumeMount:   "/config",
			VolumeName:    "config-volume",
			Label:         "master",
			Profile:       "default",
			Source:        "http://config-service.default.svc:8080",
		},
		Resources: InitContainerResources{
			Requests: InitContainerResourcesList{
				CPU:    resource.NewScaledQuantity(100, resource.Milli).String(),
				Memory: resource.NewScaledQuantity(10, resource.Mega).String(),
			},
			Limits: InitContainerResourcesList{
				CPU:    resource.NewScaledQuantity(100, resource.Milli).String(),
				Memory: resource.NewScaledQuantity(50, resource.Mega).String(),
			},
		},
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*w = WebhookConfig(raw)
	return nil
}

// NewWebhook creates a new instance of a mutating webhook for automatic sidecar injection.
func NewWebhook(p WebhookParameters) (*Webhook, error) {
	config, err := loadConfig(p.ConfigFile)
	if err != nil {
		return nil, err
	}
	pair, err := tls.LoadX509KeyPair(p.CertFile, p.KeyFile)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	// watch the parent directory of the target files so we can catch
	// symlink updates of k8s ConfigMaps volumes.
	for _, file := range []string{p.ConfigFile, p.CertFile, p.KeyFile} {
		watchDir, _ := filepath.Split(file)
		if err := watcher.Watch(watchDir); err != nil {
			return nil, fmt.Errorf("could not watch %v: %v", file, err)
		}
	}

	wh := &Webhook{
		server: &http.Server{
			Addr: fmt.Sprintf(":%v", p.Port),
		},
		config:     config,
		configFile: p.ConfigFile,
		watcher:    watcher,
		certFile:   p.CertFile,
		keyFile:    p.KeyFile,
		cert:       &pair,
	}
	// mtls disabled because apiserver webhook cert usage is still TBD.
	wh.server.TLSConfig = &tls.Config{GetCertificate: wh.getCert}
	h := http.NewServeMux()
	h.HandleFunc("/inject", wh.serveInject)
	wh.server.Handler = h

	return wh, nil
}

// Run starts the webhook control loop
func (wh *Webhook) Run(stop <-chan struct{}) {
	go func() {
		if err := wh.server.ListenAndServeTLS("", ""); err != nil {
			log.Printf("ListenAndServeTLS for admission webhook returned error: %v", err)
		}
	}()
	defer wh.watcher.Close() // nolint: errcheck
	defer wh.server.Close()  // nolint: errcheck

	var healthC <-chan time.Time
	if wh.healthCheckInterval != 0 && wh.healthCheckFile != "" {
		t := time.NewTicker(wh.healthCheckInterval)
		healthC = t.C
		defer t.Stop()
	}
	var timerC <-chan time.Time

	for {
		select {
		case <-timerC:
			config, err := loadConfig(wh.configFile)
			if err != nil {
				log.Printf("update error: %v", err)
				break
			}

			pair, err := tls.LoadX509KeyPair(wh.certFile, wh.keyFile)
			if err != nil {
				log.Printf("reload cert error: %v", err)
				break
			}
			wh.mu.Lock()
			wh.config = config
			wh.cert = &pair
			wh.mu.Unlock()
		case event := <-wh.watcher.Event:
			// use a timer to debounce configuration updates
			if event.IsModify() || event.IsCreate() {
				timerC = time.After(watchDebounceDelay)
			}
		case err := <-wh.watcher.Error:
			log.Printf("Watcher error: %v", err)
		case <-healthC:
			content := []byte(`ok`)
			if err := ioutil.WriteFile(wh.healthCheckFile, content, 0644); err != nil {
				log.Printf("Health check update of %q failed: %v", wh.healthCheckFile, err)
			}
		case <-stop:
			return
		}
	}
}

func (wh *Webhook) serveInject(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		log.Println("no body found")
		http.Error(w, "no body found", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Printf("contentType=%s, expect application/json \n", contentType)
		http.Error(w, "invalid Content-Type, want `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var reviewResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		log.Printf("Could not decode body: %v", err)
		reviewResponse = toAdmissionResponse(err)
	} else {
		reviewResponse = wh.inject(&ar)
	}

	response := v1beta1.AdmissionReview{}
	if reviewResponse != nil {
		response.Response = reviewResponse
		if ar.Request != nil {
			response.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(response)
	if err != nil {
		log.Printf("Could not encode response: %v", err)
		http.Error(w, fmt.Sprintf("could encode response: %v", err), http.StatusInternalServerError)
	}
	if _, err := w.Write(resp); err != nil {
		log.Printf("Could not write response: %v", err)
		http.Error(w, fmt.Sprintf("could write response: %v", err), http.StatusInternalServerError)
	}
}

func (wh *Webhook) inject(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	statusKey := wh.config.AnnotationPrefix + "status"
	injectKey := wh.config.AnnotationPrefix + "inject"

	req := ar.Request
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		log.Printf("Could not unmarshal raw object: %v", err)
		return toAdmissionResponse(err)
	}

	log.Printf("AdmissionReview for Kind=%v Namespace=%v Name=%v (%v) UID=%v Rfc6902PatchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, pod.Name, req.UID, req.Operation, req.UserInfo)
	log.Printf("Object: %v", string(req.Object.Raw))
	log.Printf("OldObject: %v", string(req.OldObject.Raw))

	if !injectRequired(ignoredNamespaces, wh.config.Policy, &pod.ObjectMeta, injectKey, statusKey) {
		log.Printf("Skipping %s/%s due to policy check", pod.Namespace, pod.Name)
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	spec, status, err := injectionData(&pod.Spec, &pod.ObjectMeta, wh.config) // nolint: lll
	if err != nil {
		return toAdmissionResponse(err)
	}

	annotations := map[string]string{statusKey: status}

	patchBytes, err := createPatch(&pod, injectionStatus(&pod, statusKey), annotations, spec)
	if err != nil {
		return toAdmissionResponse(err)
	}

	log.Printf("AdmissionResponse: patch=%v\n", string(patchBytes))

	reviewResponse := v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
	return &reviewResponse
}

func (wh *Webhook) getCert(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	wh.mu.Lock()
	defer wh.mu.Unlock()
	return wh.cert, nil
}

func loadConfig(injectFile string) (*WebhookConfig, error) {
	data, err := ioutil.ReadFile(injectFile)
	if err != nil {
		return nil, err
	}
	var c WebhookConfig
	if err := yaml.Unmarshal(data, &c); err != nil { // nolint: vetshadow
		return nil, err
	}

	log.Printf("New configuration: sha256sum %x", sha256.Sum256(data))

	return &c, nil
}

func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{Result: &metav1.Status{Message: err.Error()}}
}

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	//_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
	//
	//// The `v1` package from k8s.io/kubernetes/pkgp/apis/core/v1 has
	//// the object defaulting functions which are not included in
	//// k8s.io/api/corev1. The default functions are required by
	//// runtime.ObjectDefaulter to workaround lack of server-side
	//// defaulting with webhooks (see
	//// https://github.com/kubernetes/kubernetes/issues/57982).
	//_ = v1.AddToScheme(runtimeScheme)
}
