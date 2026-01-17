package k8s

import (
	"path/filepath"

	"github.com/TheAlpha16/isolet/herald/utils"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetRestConfig() (*rest.Config, error) {
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
