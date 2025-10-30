package k8s

import (
	"net/http"
	"path/filepath"

	"github.com/TheAlpha16/isolet/api/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getRestConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	kubeconfig := utils.GetConfig().K8s.KubeConfigFilePath
	if kubeconfig == "" {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func NewK8sClient(httpClient *http.Client) (client.Client, error) {
	config, err := getRestConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := client.New(config, client.Options{
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
