package k8s

import (
	"fmt"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type K8SClient struct {
	Clientset *kubernetes.Clientset
	Flags     *genericclioptions.ConfigFlags
}

func InitClient() (*K8SClient, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	flags := genericclioptions.NewConfigFlags(true)
	flags.KubeConfig = &kubeconfig

	config, err := flags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("Download error kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("client creation error: %w", err)
	}

	return &K8SClient{Clientset: clientset, Flags: flags}, nil

}
