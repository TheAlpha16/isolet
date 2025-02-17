package instance

import (
	"path/filepath"

	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/config"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClient struct {
	ClientSet     *kubernetes.Clientset
	DynamicClient dynamic.Interface
	RESTMapper    *restmapper.DeferredDiscoveryRESTMapper
}

func getRestConfig() (*rest.Config, error) {
	var configRest *rest.Config
	configRest, err := rest.InClusterConfig()

	if err == nil {
		return configRest, nil
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = utils.StringAddr(filepath.Join(home, ".kube", "config"))
	} else {
		kubeconfig = utils.StringAddr(config.KUBECONFIG_FILE_PATH)
	}

	configRest, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return configRest, nil
}

func NewK8sClient() (*K8sClient, error) {
	config, err := getRestConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))

	return &K8sClient{
		ClientSet:     clientset,
		DynamicClient: dynamicClient,
		RESTMapper:    restMapper,
	}, nil
}
