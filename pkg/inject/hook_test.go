package inject

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/gomega"
	"github.com/wandera/scccmd/internal/testcerts"
	"github.com/wandera/scccmd/internal/testutil"
	"gopkg.in/yaml.v2"
	"k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	rotatedKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA3Tr24CaBegyfkdDGWckqMHEWvpJBThjXlMz/FKcg1bgq57OD
oNHXN4dcyPCHWWEY3Eo3YG1es4pqTkvzK0+1JoY6/K88Lu1ePj5PeSFuWfPWi1BW
9oyWJW+AAzqqGkZmSo4z26N+E7N8ht5bTBMNVD3jqz9+MaqCTVmQ6dAgdFKH07wd
XWh6kKoh2g9bgBKB+qWrezmVRb31i93sM1pJos35cUmIbWgiSQYuSXEInitejcGZ
cBjqRy61SiZB7nbmMoew0G0aGXe20Wx+QiyJjbt9XNUm0IvjAJ1SSiPqfFQ4F+tx
K4q3xAwp1smyiMv57RNC2ny8YMntZYgQDDkhBQIDAQABAoIBAQDZHK396yw0WEEd
vFN8+CRUaBfXLPe0KkMgAFLxtNdPhy9sNsueP3HESC7x8MQUHmtkfd1837kJ4HRV
pMnfnpj8Vs17AIrCzycnVMVv7jQ7SUcrb8v4qJ4N3TA3exJHOQHYd1hDXF81/Hbg
cUYOEcCKBTby8BvrqBe6y4ShQiUnoaeeM5j9x32+QB652/9PMuZJ9xfwyoEBjoVA
cccp+u3oBX864ztaG9Gn0wbgRVeafsPfuAOUmShykohV1mVJddiA0wayxGi0TmoK
dwrltdToI7BmpmmTLc59O44JFGwO67aJQHsrHBjEnpWlxFDwbfZuf93FgdFUFFjr
tVx2dPF9AoGBAPkIaUYxMSW78MD9862eJlGS6F/SBfPLTve3l9M+UFnEsWapNz+V
anlupaBtlfRJxLDzjheMVZTv/TaFIrrMdN/xc/NoYj6lw5RV8PEfKPB0FjSAqSEl
iVOA5q4kuI1xEeV7xLE4uJUF3wdoHz9cSmjrXDVZXq/KsaInABlAl1xjAoGBAONr
bgFUlQ+fbsGNjeefrIXBSU93cDo5tkd4+lefrmUVTwCo5V/AvMMQi88gV/sz6qCJ
gR0hXOqdQyn++cuqfMTDj4lediDGFEPho+n4ToIvkA020NQ05cKxcmn/6Ei+P9pk
v+zoT9RiVnkBje2n/KU2d/PEL9Nl4gvvAgPLt8V3AoGAZ6JZdQ15n3Nj0FyecKz0
01Oogl+7fGYqGap8cztmYsUY8lkPFdXPNnOWV3njQoMEaIMiqagL4Wwx2uNyvXvi
U2N+1lelMt720h8loqJN/irBJt44BARD7s0gsm2zo6DfSrnD8+Bf6BxGYSWyg0Kb
8KepesYTQmK+o3VJdDjOBHMCgYAIxbwYkQqu75d2H9+5b49YGXyadCEAHfnKCACg
IKi5fXjurZUrfGPLonfCJZ0/M2F5j9RLK15KLobIt+0qzgjCDkkbI2mrGfjuJWYN
QGbG3s7Ps62age/a8r1XGWf8ZlpQMlK08MEjkCeFw2mWIUS9mrxFyuuNXAC8NRv+
yXztQQKBgQDWTFFQdeYfuiKHrNmgOmLVuk1WhAaDgsdK8RrnNZgJX9bd7n7bm7No
GheN946AYsFge4DX7o0UXXJ3h5hTFn/hSWASI54cO6WyWNEiaP5HRlZqK7Jfej7L
mz+dlU3j/BY19RLmYeg4jFV4W66CnkDqpneOJs5WdmFFoWnHn7gRBw==
-----END RSA PRIVATE KEY-----`)

	// ServerCert is a test cert for dynamic admission controller.
	rotatedCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDATCCAemgAwIBAgIJAJwGb32Zn8sDMA0GCSqGSIb3DQEBCwUAMA4xDDAKBgNV
BAMMA19jYTAgFw0xODAzMTYxNzI0NDJaGA8yMjkxMTIzMDE3MjQ0MlowEjEQMA4G
A1UEAwwHX3NlcnZlcjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN06
9uAmgXoMn5HQxlnJKjBxFr6SQU4Y15TM/xSnINW4Kuezg6DR1zeHXMjwh1lhGNxK
N2BtXrOKak5L8ytPtSaGOvyvPC7tXj4+T3khblnz1otQVvaMliVvgAM6qhpGZkqO
M9ujfhOzfIbeW0wTDVQ946s/fjGqgk1ZkOnQIHRSh9O8HV1oepCqIdoPW4ASgfql
q3s5lUW99Yvd7DNaSaLN+XFJiG1oIkkGLklxCJ4rXo3BmXAY6kcutUomQe525jKH
sNBtGhl3ttFsfkIsiY27fVzVJtCL4wCdUkoj6nxUOBfrcSuKt8QMKdbJsojL+e0T
Qtp8vGDJ7WWIEAw5IQUCAwEAAaNcMFowCQYDVR0TBAIwADALBgNVHQ8EBAMCBeAw
HQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMBMCEGA1UdEQQaMBiHBH8AAAGH
EAAAAAAAAAAAAAAAAAAAAAEwDQYJKoZIhvcNAQELBQADggEBACbBlWo/pY/OIJaW
RwkfSRVzEIWpHt5OF6p93xfyy4/zVwwhH1AQB7Euji8vOaVNOpMfGYLNH3KIRReC
CIvGEH4yZDbpiH2cOshqMCuV1CMRUTdl4mq6M0PtGm6b8OG3uIFTLIR973LBWOl5
wCR1yrefT1NHuIScGaBXUGAV4JAx37pfg84hDD73T2j1TDD3Lrmsb9WCP+L26TG6
ICN61cIhgz8wChQpF8/fFAI5Fjbjrz5C1Xw/EUHLf/TTn/7Yfp2BHsGm126Et+k+
+MLBzBfrHKwPaGqDvNHUDrI6c3GI0Qp7jW93FbL5ul8JQ+AowoMF2dIEbN9qQEVP
ZOQ5UvU=
-----END CERTIFICATE-----`)
)

const (
	annotationPrefix    = "config.scccmd.github.com/"
	annotationInjectKey = "config.scccmd.github.com/inject"
)

func TestInjectRequired(t *testing.T) {
	podSpec := &corev1.PodSpec{}
	podSpecHostNetwork := &corev1.PodSpec{
		HostNetwork: true,
	}

	cases := []struct {
		policy  InjectionPolicy
		podSpec *corev1.PodSpec
		meta    *metav1.ObjectMeta
		want    bool
	}{
		{
			policy:  InjectionPolicyEnabled,
			podSpec: podSpec,
			meta: &metav1.ObjectMeta{
				Name:        "no-policy",
				Namespace:   "test-namespace",
				Annotations: map[string]string{},
			},
			want: true,
		},
		{
			policy:  InjectionPolicyEnabled,
			podSpec: podSpec,
			meta: &metav1.ObjectMeta{
				Name:      "default-policy",
				Namespace: "test-namespace",
			},
			want: true,
		},
		{
			policy:  InjectionPolicyEnabled,
			podSpec: podSpec,
			meta: &metav1.ObjectMeta{
				Name:        "force-on-policy",
				Namespace:   "test-namespace",
				Annotations: map[string]string{annotationInjectKey: "true"},
			},
			want: true,
		},
		{
			policy:  InjectionPolicyEnabled,
			podSpec: podSpec,
			meta: &metav1.ObjectMeta{
				Name:        "force-off-policy",
				Namespace:   "test-namespace",
				Annotations: map[string]string{annotationInjectKey: "false"},
			},
			want: false,
		},
		{
			policy:  InjectionPolicyDisabled,
			podSpec: podSpec,
			meta: &metav1.ObjectMeta{
				Name:        "no-policy",
				Namespace:   "test-namespace",
				Annotations: map[string]string{},
			},
			want: false,
		},
		{
			policy:  InjectionPolicyDisabled,
			podSpec: podSpec,
			meta: &metav1.ObjectMeta{
				Name:      "default-policy",
				Namespace: "test-namespace",
			},
			want: false,
		},
		{
			policy:  InjectionPolicyDisabled,
			podSpec: podSpec,
			meta: &metav1.ObjectMeta{
				Name:        "force-on-policy",
				Namespace:   "test-namespace",
				Annotations: map[string]string{annotationInjectKey: "true"},
			},
			want: true,
		},
		{
			policy:  InjectionPolicyDisabled,
			podSpec: podSpec,
			meta: &metav1.ObjectMeta{
				Name:        "force-off-policy",
				Namespace:   "test-namespace",
				Annotations: map[string]string{},
			},
			want: false,
		},
		{
			policy:  InjectionPolicyEnabled,
			podSpec: podSpecHostNetwork,
			meta: &metav1.ObjectMeta{
				Name:        "force-off-policy",
				Namespace:   "test-namespace",
				Annotations: map[string]string{annotationInjectKey: "false"},
			},
			want: false,
		},
	}

	for _, c := range cases {
		if got := injectRequired(ignoredNamespaces, c.policy, c.meta, annotationInjectKey, "config.scccmd.github.com/status"); got != c.want {
			t.Errorf("injectRequired(%v, %v) got %v want %v", c.policy, c.meta, got, c.want)
		}
	}
}

func makeTestData(t testing.TB, skip bool) []byte {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Annotations: map[string]string{
				annotationPrefix + "destination": "config.yaml",
			},
		},
		Spec: corev1.PodSpec{
			Volumes:        []corev1.Volume{{Name: "v0"}},
			InitContainers: []corev1.Container{{Name: "c0"}},
			Containers:     []corev1.Container{{Name: "c1"}, {Name: "c2"}},
		},
	}

	if skip {
		pod.ObjectMeta.Annotations[annotationInjectKey] = "false"
	}

	raw, err := json.Marshal(&pod)
	if err != nil {
		t.Fatalf("Could not create test pod: %v", err)
	}

	review := v1.AdmissionReview{
		Request: &v1.AdmissionRequest{
			Kind: metav1.GroupVersionKind{},
			Object: runtime.RawExtension{
				Raw: raw,
			},
			Operation: v1.Create,
		},
	}
	reviewJSON, err := json.Marshal(review)
	if err != nil {
		t.Fatalf("Failed to create AdmissionReview: %v", err)
	}
	return reviewJSON
}

func createWebhook(t testing.TB) (*Webhook, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "webhook_test")
	if err != nil {
		t.Fatalf("TempDir() failed: %v", err)
	}
	cleanup := func() {
		os.RemoveAll(dir) // nolint: errcheck
	}

	config := &WebhookConfig{
		Policy: InjectionPolicyEnabled,
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
				CPU:    resource.NewScaledQuantity(10, resource.Milli).String(),
				Memory: resource.NewScaledQuantity(10, resource.Mega).String(),
			},
			Limits: InitContainerResourcesList{
				CPU:    resource.NewScaledQuantity(50, resource.Milli).String(),
				Memory: resource.NewScaledQuantity(50, resource.Mega).String(),
			},
		},
	}

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		cleanup()
		t.Fatalf("Could not marshal test injection config: %v", err)
	}

	var (
		configFile = filepath.Join(dir, "config-file.yaml")
		certFile   = filepath.Join(dir, "cert-file.yaml")
		keyFile    = filepath.Join(dir, "key-file.yaml")
		port       = 0
	)

	if err := os.WriteFile(configFile, configBytes, 0o644); err != nil { // nolint: vetshadow
		cleanup()
		t.Fatalf("WriteFile(%v) failed: %v", configFile, err)
	}

	// cert
	if err := os.WriteFile(certFile, testcerts.ServerCert, 0o644); err != nil { // nolint: vetshadow
		cleanup()
		t.Fatalf("WriteFile(%v) failed: %v", certFile, err)
	}
	// key
	if err := os.WriteFile(keyFile, testcerts.ServerKey, 0o644); err != nil { // nolint: vetshadow
		cleanup()
		t.Fatalf("WriteFile(%v) failed: %v", keyFile, err)
	}

	wh, err := NewWebhook(WebhookParameters{
		ConfigFile: configFile, CertFile: certFile, KeyFile: keyFile, Port: port,
	})
	if err != nil {
		cleanup()
		t.Fatalf("NewWebhook() failed: %v", err)
	}
	return wh, cleanup
}

func TestRunAndServe(t *testing.T) {
	wh, cleanup := createWebhook(t)
	defer cleanup()
	stop := make(chan struct{})
	defer func() { close(stop) }()
	go wh.Run(stop)

	validReview := makeTestData(t, false)
	skipReview := makeTestData(t, true)

	// nolint: lll
	validPatch := []byte(`[
		{
			"op": "add",
			"path": "/spec/initContainers/0/volumeMounts",
			"value": [
			  {
				"name": "config-volume",
				"mountPath": "/config"
			  }
			]
		},
		{
			"op":"add",
			"path":"/spec/containers/0/volumeMounts",
			"value":[
				{
					"name":"config-volume",
					"mountPath":"/config"
				}
			]
		},
		{
			"op":"add",
			"path":"/spec/containers/1/volumeMounts",
			"value":[
				{
					"name":"config-volume",
					"mountPath":"/config"
				}
			]
		},
		{
			"op":"add",
			"path":"/spec/volumes/-",
			"value":{
				"name":"config-volume",
				"emptyDir":{}
			}
		},
		{
			"op":"add",
			"path":"/spec/initContainers/0",
			"value":{
				"name":"config-init",
				"image":"wanderadock/scccmd",
				"args":["get","values","--source","http://config-service.default.svc:8080","--application","c1","--profile","default","--label","master","--destination","config.yaml"],
				"resources":{"limits":{"cpu":"50m","memory":"50M"},"requests":{"cpu":"10m","memory":"10M"}},
				"volumeMounts":[{"name":"config-volume","mountPath":"/config"}]
			}
		},
		{
			"op":"add",
			"path":"/metadata/annotations/config.scccmd.github.com~1status",
			"value":"{\"initContainers\":[\"config-init\"],\"volumeMounts\":[\"config-volume\"],\"volumes\":[\"config-volume\"]}"
		}
	]`)

	cases := []struct {
		name           string
		body           []byte
		contentType    string
		wantAllowed    bool
		wantStatusCode int
		wantPatch      []byte
	}{
		{
			name:           "valid",
			body:           validReview,
			contentType:    "application/json",
			wantAllowed:    true,
			wantStatusCode: http.StatusOK,
			wantPatch:      validPatch,
		},
		{
			name:           "skipped",
			body:           skipReview,
			contentType:    "application/json",
			wantAllowed:    true,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "wrong content-type",
			body:           validReview,
			contentType:    "application/yaml",
			wantAllowed:    false,
			wantStatusCode: http.StatusUnsupportedMediaType,
		},
		{
			name:           "bad content",
			body:           []byte{0, 1, 2, 3, 4, 5}, // random data
			contentType:    "application/json",
			wantAllowed:    false,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "missing body",
			contentType:    "application/json",
			wantAllowed:    false,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("[%d] %s", i, c.name), func(t *testing.T) {
			req := httptest.NewRequest("POST", "http://sidecar-injector/inject", bytes.NewReader(c.body))
			req.Header.Add("Content-Type", c.contentType)

			w := httptest.NewRecorder()
			wh.serveInject(w, req)
			res := w.Result()

			if res.StatusCode != c.wantStatusCode {
				t.Fatalf("wrong status code: \ngot %v \nwant %v", res.StatusCode, c.wantStatusCode)
			}

			if res.StatusCode != http.StatusOK {
				return
			}

			gotBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read body: %v", err)
			}
			var gotReview v1.AdmissionReview
			if err := json.Unmarshal(gotBody, &gotReview); err != nil {
				t.Fatalf("could not decode response body: %v", err)
			}
			if gotReview.Response.Allowed != c.wantAllowed {
				t.Fatalf("AdmissionReview.Response.Allowed is wrong : got %v want %v",
					gotReview.Response.Allowed, c.wantAllowed)
			}

			var gotPatch bytes.Buffer
			if len(gotReview.Response.Patch) > 0 {
				if err := json.Compact(&gotPatch, gotReview.Response.Patch); err != nil {
					t.Fatal(err.Error())
				}
			}
			var wantPatch bytes.Buffer
			if len(c.wantPatch) > 0 {
				if err := json.Compact(&wantPatch, c.wantPatch); err != nil {
					t.Fatal(err.Error())
				}
			}
			testutil.AssertString(t, "got bad patch", wantPatch.String(), wantPatch.String())
		})
	}
}

func TestReloadCert(t *testing.T) {
	wh, cleanup := createWebhook(t)
	defer cleanup()
	stop := make(chan struct{})
	defer func() { close(stop) }()
	go wh.Run(stop)
	checkCert(t, wh, testcerts.ServerCert, testcerts.ServerKey)
	// Update cert/key files.
	if err := os.WriteFile(wh.certFile, rotatedCert, 0o644); err != nil { // nolint: vetshadow
		cleanup()
		t.Fatalf("WriteFile(%v) failed: %v", wh.certFile, err)
	}
	if err := os.WriteFile(wh.keyFile, rotatedKey, 0o644); err != nil { // nolint: vetshadow
		cleanup()
		t.Fatalf("WriteFile(%v) failed: %v", wh.keyFile, err)
	}
	g := gomega.NewGomegaWithT(t)
	g.Eventually(func() bool {
		return checkCert(t, wh, rotatedCert, rotatedKey)
	}, "10s", "100ms").Should(gomega.BeTrue())
}

func checkCert(t *testing.T, wh *Webhook, cert, key []byte) bool {
	t.Helper()
	actual, err := wh.getCert(nil)
	if err != nil {
		t.Fatalf("fail to get certificate from webhook: %s", err)
	}
	expected, err := tls.X509KeyPair(cert, key)
	if err != nil {
		t.Fatalf("fail to load test certs.")
	}
	return bytes.Equal(actual.Certificate[0], expected.Certificate[0])
}
