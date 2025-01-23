package deployment

import (
	"fmt"
	"strings"
	"path/filepath"


	"github.com/TheAlpha16/isolet/admin/config"
	"github.com/TheAlpha16/isolet/admin/models"
	"github.com/TheAlpha16/isolet/admin/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"k8s.io/apimachinery/pkg/api/resource"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetKubeClient() (*kubernetes.Clientset, error) {
	var configRest *rest.Config
	configRest, err := rest.InClusterConfig()

	if err == nil {
		clientset, err := kubernetes.NewForConfig(configRest)
		if err != nil {
			return nil, err
		}
		return clientset, nil
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

	clientset, err := kubernetes.NewForConfig(configRest)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getPodObject(instance_name string, flagObject models.Flag, image *models.Image) *core.Pod {
	var imagePath string
	var cpu string
	var memory string

	if image.Registry != "" {
		imagePath = image.Registry
	} else {
		imagePath = config.IMAGE_REGISTRY
	}

	imagePath = strings.TrimSuffix(imagePath, "/")
	imagePath = fmt.Sprintf("%s/%s", imagePath, image.Image)

	if image.CPU == 0 {
		cpu = config.CPU_REQUEST
	} else {
		cpu = fmt.Sprintf("%dm", image.CPU)
	}

	if image.Memory == 0 {
		memory = config.MEMORY_REQUEST
	} else {
		memory = fmt.Sprintf("%dMi", image.Memory)
	}

	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance_name,
			Namespace: config.INSTANCE_NAMESPACE,
			Labels: map[string]string{
				"chall_id":  fmt.Sprintf("%d", flagObject.ChallID),
				"teamid": fmt.Sprintf("%d", flagObject.TeamID),
				"app":    "instance",
			},
		},
		Spec: core.PodSpec{
			AutomountServiceAccountToken:  utils.BoolAddr(false),
			EnableServiceLinks:            utils.BoolAddr(false),
			TerminationGracePeriodSeconds: utils.Int64Addr(config.TERMINATION_PERIOD),
			Containers: []core.Container{
				{
					Name:  instance_name,
					Image: imagePath,
					Ports: []core.ContainerPort{
						{
							ContainerPort: int32(image.Port),
						},
					},
					Resources: core.ResourceRequirements{
						Limits: core.ResourceList{
							core.ResourceName(core.ResourceCPU):              resource.MustParse(config.CPU_LIMIT),
							core.ResourceName(core.ResourceMemory):           resource.MustParse(config.MEMORY_LIMIT),
							core.ResourceName(core.ResourceEphemeralStorage): resource.MustParse(config.DISK_LIMIT),
						},
						Requests: core.ResourceList{
							core.ResourceName(core.ResourceCPU):              resource.MustParse(cpu),
							core.ResourceName(core.ResourceMemory):           resource.MustParse(memory),
							core.ResourceName(core.ResourceEphemeralStorage): resource.MustParse(config.DISK_REQUEST),
						},
					},
					ImagePullPolicy: core.PullAlways,
					Env: []core.EnvVar{
						{
							Name:  "CTF_NAME",
							Value: config.CTF_NAME,
						},
						{
							Name: "USERNAME",
							Value: config.DEFAULT_USERNAME,
						},
						{
							Name:  "USER_PASSWORD",
							Value: flagObject.Password,
						},
						{
							Name:  "FLAG",
							Value: flagObject.Flag,
						},
					},
				},
			},
		},
	}
}

func getServiceObject(instance_name string, flagObject models.Flag) *core.Service {
	return &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("svc-%s", instance_name),
			Namespace: config.INSTANCE_NAMESPACE,
			Labels: map[string]string{
				"chall_id":  fmt.Sprintf("%d", flagObject.ChallID),
				"teamid": fmt.Sprintf("%d", flagObject.TeamID),
			},
		},
		Spec: core.ServiceSpec{
			Type: "NodePort",
			Ports: []core.ServicePort{
				{
					Port: int32(flagObject.Port),
				},
			},
			Selector: map[string]string{
				"chall_id":  fmt.Sprintf("%d", flagObject.ChallID),
				"teamid": fmt.Sprintf("%d", flagObject.TeamID),
				"app":    "instance",
			},
		},
	}
}

