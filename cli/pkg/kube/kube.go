package kube

import (
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type (
	Config struct {
		Kubeconfig string
	}

	Client struct {
		*kubernetes.Clientset
	}
)

// New creates a kubernetes client
func New(cfg *Config) (*Client, error) {
	var kubeconfig string
	var config *rest.Config
	var err error

	if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	if cfg.Kubeconfig != "" {
		kubeconfig = cfg.Kubeconfig
	}
	if kubeconfig != "" && isExist(kubeconfig) {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{clientSet}, nil
}

// GetNodes get nodes
func (c *Client) GetNodes() (*corev1.NodeList, error) {
	return c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

// GetSecrets get secrets
func (c *Client) GetSecrets(namespace string) (*corev1.SecretList, error) {
	return c.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetSecret(namespace, name string) (*corev1.Secret, error) {
	return c.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *Client) GetSecretsWithLabel(namespace, label string) (*corev1.SecretList, error) {
	return c.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})
}

// CreateSecret Create secret
func (c *Client) CreateSecret(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	return c.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
}

// UpdateSecret Update secret
func (c *Client) UpdateSecret(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	return c.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
}

func (c *Client) PatchSecret(namespace, name string, data []byte) (*corev1.Secret, error) {
	return c.CoreV1().Secrets(namespace).Patch(
		context.TODO(),
		name,
		types.StrategicMergePatchType,
		data,
		metav1.PatchOptions{},
	)
}

func (c *Client) ApplySecret(namespace, name string, data map[string]string) (*corev1.Secret, error) {
	olds, err := c.GetSecret(namespace, name)
	if err != nil || olds == nil {
		var news corev1.Secret
		news.SetName(name)
		news.SetNamespace(namespace)
		news.StringData = data
		return c.CreateSecret(namespace, &news)
	}
	olds.StringData = data
	return c.UpdateSecret(namespace, olds)
}

func (c *Client) ApplySecretBytes(namespace, name string, data map[string][]byte) (*corev1.Secret, error) {
	olds, err := c.GetSecret(namespace, name)
	if err != nil || olds == nil {
		var news corev1.Secret
		news.SetName(name)
		news.SetNamespace(namespace)
		news.Data = data
		return c.CreateSecret(namespace, &news)
	}
	olds.Data = data
	return c.UpdateSecret(namespace, olds)
}

// DeleteSecret Delete secret
func (c *Client) DeleteSecret(namespace, name string) error {
	return c.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// GetServices get services
func (c *Client) GetServices(namespace string) (*corev1.ServiceList, error) {
	return c.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetServicesWithLabel(namespace, label string) (*corev1.ServiceList, error) {
	return c.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})
}

// GetServices get services
func (c *Client) GetService(namespace, name string) (*corev1.Service, error) {
	return c.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// GetConfigMap get ConfigMap
func (c *Client) GetConfigMap(namespace, name string) (*corev1.ConfigMap, error) {
	return c.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *Client) CreateConfigMap(namespace string, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return c.CoreV1().ConfigMaps(namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
}

func (c *Client) UpdateConfigMap(namespace string, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return c.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
}

func (c *Client) PatchConfigMap(namespace, name string, data []byte) (*corev1.ConfigMap, error) {
	return c.CoreV1().ConfigMaps(namespace).Patch(
		context.TODO(),
		name,
		types.StrategicMergePatchType,
		data,
		metav1.PatchOptions{},
	)
}

func (c *Client) ApplyConfigMap(namespace, name string, data map[string]string) (*corev1.ConfigMap, error) {
	olds, err := c.GetConfigMap(namespace, name)
	if err != nil || olds == nil {
		var news corev1.ConfigMap
		news.SetName(name)
		news.SetNamespace(namespace)
		news.Data = data
		return c.CreateConfigMap(namespace, &news)
	}
	olds.Data = data
	return c.UpdateConfigMap(namespace, olds)
}

// ----------------------------------------------------------------------------
// Helpers...

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func isExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
