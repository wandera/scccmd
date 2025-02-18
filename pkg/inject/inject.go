package inject

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InjectionPolicy determines the policy for injecting the
// config init container into the watched namespace(s).
type InjectionPolicy string

// SidecarInjectionStatus contains basic information about the
// injected sidecar. This includes the names of added containers and
// volumes.
type SidecarInjectionStatus struct {
	InitContainers []string `json:"initContainers"`
	VolumeMounts   []string `json:"volumeMounts"`
	Volumes        []string `json:"volumes"`
}

// SidecarInjectionSpec collects all container types and volumes for
// sidecar mesh injection.
type SidecarInjectionSpec struct {
	InitContainers []v1.Container   `yaml:"initContainers"`
	VolumeMounts   []v1.VolumeMount `yaml:"volumeMounts"`
	Volumes        []v1.Volume      `yaml:"volumes"`
}

type dynamicConfig struct {
	containerName string
	volumeName    string
	volumeMount   string
	imageArgs     []string
}

const (
	// InjectionPolicyDisabled specifies that the sidecar injector
	// will not inject the sidecar into resources by default for the
	// namespace(s) being watched. Resources can enable injection
	// using the "<annotation prefix>/inject" annotation with value of
	// true.
	InjectionPolicyDisabled InjectionPolicy = "disabled"

	// InjectionPolicyEnabled specifies that the sidecar injector will
	// inject the sidecar into resources by default for the
	// namespace(s) being watched. Resources can disable injection
	// using the "<annotation prefix>/inject" annotation with value of
	// false.
	InjectionPolicyEnabled InjectionPolicy = "enabled"
)

// InjectionStatus extracts the injection status from the pod.
func injectionStatus(pod *corev1.Pod, annotationStatusKey string) *SidecarInjectionStatus {
	var statusBytes []byte
	if pod.ObjectMeta.Annotations != nil {
		if value, ok := pod.ObjectMeta.Annotations[annotationStatusKey]; ok {
			statusBytes = []byte(value)
		}
	}

	// default case when injected pod has explicit status
	var status SidecarInjectionStatus
	if err := json.Unmarshal(statusBytes, &status); err == nil {
		// heuristic assumes status is valid if any of the resource
		// lists is non-empty.
		if len(status.InitContainers) != 0 ||
			len(status.VolumeMounts) != 0 ||
			len(status.Volumes) != 0 {
			return &status
		}
	}
	return &SidecarInjectionStatus{}
}

func injectionData(spec *v1.PodSpec, metadata *metav1.ObjectMeta, config *WebhookConfig) (*SidecarInjectionSpec, string, error) { // nolint: lll
	d, err := calculateDynamicConfig(config, metadata.GetAnnotations(), spec)
	if err != nil {
		return nil, "", err
	}

	volumeMount := corev1.VolumeMount{
		Name:      d.volumeName,
		MountPath: d.volumeMount,
	}

	sic := SidecarInjectionSpec{
		InitContainers: []corev1.Container{
			{
				Name:         d.containerName,
				Image:        config.ContainerImage,
				Args:         d.imageArgs,
				VolumeMounts: []corev1.VolumeMount{volumeMount},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						"cpu":    resource.MustParse(config.Resources.Requests.CPU),
						"memory": resource.MustParse(config.Resources.Requests.Memory),
					},
					Limits: corev1.ResourceList{
						"cpu":    resource.MustParse(config.Resources.Limits.CPU),
						"memory": resource.MustParse(config.Resources.Limits.Memory),
					},
				},
				SecurityContext: &corev1.SecurityContext{
					AllowPrivilegeEscalation: &config.SecurityContext.AllowPrivilegeEscalation,
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{volumeMount},
		Volumes: []corev1.Volume{
			{
				Name:         volumeMount.Name,
				VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
			},
		},
	}

	status := &SidecarInjectionStatus{}
	for _, c := range sic.InitContainers {
		status.InitContainers = append(status.InitContainers, c.Name)
	}
	for _, c := range sic.VolumeMounts {
		status.VolumeMounts = append(status.VolumeMounts, c.Name)
	}
	for _, c := range sic.Volumes {
		status.Volumes = append(status.Volumes, c.Name)
	}
	statusAnnotationValue, err := json.Marshal(status)
	if err != nil {
		return nil, "", fmt.Errorf("error encoded injection status: %v", err)
	}
	return &sic, string(statusAnnotationValue), nil
}

func injectRequired(ignored []string, namespacePolicy InjectionPolicy, metadata *metav1.ObjectMeta, annotationInjectKey, annotationStatusKey string) bool { // nolint: lll
	// skip special kubernetes system namespaces
	for _, namespace := range ignored {
		if metadata.Namespace == namespace {
			return false
		}
	}

	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	var useDefault bool
	var inject bool
	switch strings.ToLower(annotations[annotationInjectKey]) {
	// http://yaml.org/type/bool.html
	case "y", "yes", "true", "on":
		inject = true
	case "":
		useDefault = true
	}

	var required bool
	switch namespacePolicy {
	default: // InjectionPolicyOff
		required = false
	case InjectionPolicyDisabled:
		if useDefault {
			required = false
		} else {
			required = inject
		}
	case InjectionPolicyEnabled:
		if useDefault {
			required = true
		} else {
			required = inject
		}
	}

	status := annotations[annotationStatusKey]

	log.Infof("Sidecar injection policy for %v/%v: namespacePolicy:%v useDefault:%v inject:%v status:%q required:%v",
		metadata.Namespace, metadata.Name, namespacePolicy, useDefault, inject, status, required)

	return required
}

func calculateImageArgs(c *WebhookConfig, a map[string]string, podSpec *corev1.PodSpec) ([]string, error) {
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
		source = c.Default.Source
	}

	if profile, ok = a[c.AnnotationPrefix+"profile"]; !ok {
		profile = c.Default.Profile
	}

	if label, ok = a[c.AnnotationPrefix+"label"]; !ok {
		label = c.Default.Label
	}

	if application, ok = a[c.AnnotationPrefix+"application"]; !ok {
		if len(podSpec.Containers) > 0 {
			application = podSpec.Containers[0].Name
			log.Debugf("defaulting application name to %s\n", application)
		}
	}

	return append([]string{"get", mode, "--source", source, "--application", application, "--profile", profile, "--label", label}, extra...), nil
}

func calculateDynamicConfig(c *WebhookConfig, a map[string]string, podSpec *corev1.PodSpec) (*dynamicConfig, error) {
	d := dynamicConfig{}
	var ok bool

	if d.containerName, ok = a[c.AnnotationPrefix+"container-name"]; !ok {
		d.containerName = c.Default.ContainerName
	}

	if d.volumeName, ok = a[c.AnnotationPrefix+"volume-name"]; !ok {
		d.volumeName = c.Default.VolumeName
	}

	if d.volumeMount, ok = a[c.AnnotationPrefix+"volume-mount"]; !ok {
		d.volumeMount = c.Default.VolumeMount
	}

	imageArgs, err := calculateImageArgs(c, a, podSpec)
	d.imageArgs = imageArgs
	return &d, err
}
