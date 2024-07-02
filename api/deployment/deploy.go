package deployment

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/gofiber/fiber/v2"

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

func DeployInstance(c *fiber.Ctx, userid int, level int) (int64, string, int32, string, error) {
	instance_name := utils.GetInstanceName(userid, level)
	password := database.GenerateRandom()[0:32]
	flag := config.CTF_NAME + "{" + database.GenerateRandom()[0:32] + "}"
	var deadline int64 = 1893456000000

	// Hostname to be known when using subdomains for connections
	hostname := utils.GetHostName(userid, level)

	kubeclient, err := GetKubeClient()
	if err != nil {
		log.Println(err)
		return deadline, "", -1, "", err
	}

	pod := getPodObject(instance_name, level, userid, password, flag)
	pod, err = kubeclient.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Println(err)
			return deadline, "", -1, "", err
		}
	}

	for {
		createdPod, err := kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Get(context.TODO(), instance_name, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
			return deadline, "", -1, "", err
		}

		if len(createdPod.Status.ContainerStatuses) > 0 {
			if createdPod.Status.ContainerStatuses[0].State.Waiting != nil && (createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "CrashLoopBackOff" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "ErrImagePull" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "ImagePullBackOff" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "CreateContainerConfigError" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "InvalidImageName" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "CreateContainerError") {
				kubeclient.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
				log.Printf("Error in launch: level%d reason: %s", level, createdPod.Status.ContainerStatuses[0].State.Waiting.Reason)
				return deadline, "", -1, "", fmt.Errorf("runtime error in image - level%d not found in registry", level)
			}
		}

		if createdPod.Status.Phase == "Running" && createdPod.Status.StartTime != nil {
			deadline = createdPod.Status.StartTime.Add(time.Minute * time.Duration(config.INSTANCE_TIME)).UnixMilli()
			break
		}
	}

	err = UpdateDeadline(kubeclient, instance_name, deadline)
	if err != nil {
		_ = DeleteInstance(c, userid, level)
		log.Println(err)
		return deadline, "", -1, "", err
	}

	svcConfig := getServiceObject(instance_name, level, userid)
	_, err = kubeclient.CoreV1().Services(svcConfig.Namespace).Create(context.TODO(), svcConfig, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		if !strings.Contains(err.Error(), "already exists") {
			return deadline, "", -1, "", err
		}
	}

	createdService, err := kubeclient.CoreV1().Services(svcConfig.Namespace).Get(context.TODO(), svcConfig.Name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return deadline, "", -1, "", err
	}
	port := createdService.Spec.Ports[0].NodePort

	nodes, err := kubeclient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return deadline, "", -1, "", err
	}
	nodeip := nodes.Items[0].Status.Addresses
	for i := 0; i < len(nodeip); i++ {
		if nodeip[i].Type == "ExternalIP" {
			hostname = nodeip[i].Address
			break
		}
	}

	if err := database.NewFlag(c, userid, level, password, flag, port, hostname, deadline); err != nil {
		log.Println(err)
		return deadline, "", -1, "", err
	}

	return deadline, password, port, hostname, nil
}

func DeleteInstance(c *fiber.Ctx, userid int, level int) error {
	instance_name := utils.GetInstanceName(userid, level)
	kubeclient, err := GetKubeClient()
	if err != nil {
		log.Println(err)
		return err
	}

	err = kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Delete(context.TODO(), instance_name, metav1.DeleteOptions{})
	if err != nil {

		if !strings.Contains(err.Error(), "not found") {
			log.Println(err)
			return err
		}

		err = kubeclient.CoreV1().Services(config.INSTANCE_NAMESPACE).Delete(context.TODO(), fmt.Sprintf("svc-%s", instance_name), metav1.DeleteOptions{})
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				log.Println(err)
				return err
			}
		}

		if err := database.DeleteFlag(c, userid, level); err != nil {
			log.Println(err)
			return err
		}

		if err := database.DeleteRunning(c, userid, level); err != nil {
			log.Println(err)
			return err
		}

		return nil
	}

	for {
		_, err := kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Get(context.TODO(), instance_name, metav1.GetOptions{})
		if err != nil {
			break
		}
	}

	err = kubeclient.CoreV1().Services(config.INSTANCE_NAMESPACE).Delete(context.TODO(), fmt.Sprintf("svc-%s", instance_name), metav1.DeleteOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			log.Println(err)
			return err
		}
	}

	if err := database.DeleteFlag(c, userid, level); err != nil {
		log.Println(err)
		return err
	}

	if err := database.DeleteRunning(c, userid, level); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func AddTime(c *fiber.Ctx, userid int, level int) (bool, string, int64) {
	instance_name := utils.GetInstanceName(userid, level)
	kubeclient, err := GetKubeClient()
	if err != nil {
		log.Println(err)
		return false, "error in extension, please contact admin", 1
	}

	_, err = kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Get(c.Context(), instance_name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			log.Println(err)
			return false, "error in extension, please contact admin", 1
		}

		if err := database.DeleteFlag(c, userid, level); err != nil {
			log.Println(err)
			return false, "instance not running", 1
		}

		if err := database.DeleteRunning(c, userid, level); err != nil {
			log.Println(err)
			return false, "instance not running", 1
		}

		return false, "instance not running", 1
	}

	isOK, message, newdeadline := database.AddTime(c, userid, level)
	if !isOK {
		return isOK, message, 1
	}

	err = UpdateDeadline(kubeclient, instance_name, newdeadline)
	if err != nil {
		log.Println(err)
		return false, "error in extension, please contact admin", 1
	}

	return true, "updated deadline successfully", newdeadline
}

func getPodObject(instance_name string, level int, userid int, password string, flag string) *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance_name,
			Namespace: config.INSTANCE_NAMESPACE,
			Labels: map[string]string{
				"level":  fmt.Sprintf("%d", level),
				"userid": fmt.Sprintf("%d", userid),
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
					Image: fmt.Sprintf("%slevel%d:latest", config.IMAGE_REGISTRY_PREFIX, level),
					Ports: []core.ContainerPort{
						{
							ContainerPort: 22,
						},
					},
					Resources: core.ResourceRequirements{
						Limits: core.ResourceList{
							core.ResourceName(core.ResourceCPU):              resource.MustParse(config.CPU_LIMIT),
							core.ResourceName(core.ResourceMemory):           resource.MustParse(config.MEMORY_LIMIT),
							core.ResourceName(core.ResourceEphemeralStorage): resource.MustParse(config.DISK_LIMIT),
						},
						Requests: core.ResourceList{
							core.ResourceName(core.ResourceCPU):              resource.MustParse(config.CPU_REQUEST),
							core.ResourceName(core.ResourceMemory):           resource.MustParse(config.MEMORY_REQUEST),
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
							Name:  "USER_PASSWORD",
							Value: password,
						},
						{
							Name:  "FLAG",
							Value: flag,
						},
					},
				},
			},
		},
	}
}

func getServiceObject(instance_name string, level int, userid int) *core.Service {
	return &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("svc-%s", instance_name),
			Namespace: config.INSTANCE_NAMESPACE,
			Labels: map[string]string{
				"level":  fmt.Sprintf("%d", level),
				"userid": fmt.Sprintf("%d", userid),
			},
		},
		Spec: core.ServiceSpec{
			Type: "NodePort",
			Ports: []core.ServicePort{
				{
					Port: 22,
				},
			},
			Selector: map[string]string{
				"level":  fmt.Sprintf("%d", level),
				"userid": fmt.Sprintf("%d", userid),
				"app":    "instance",
			},
		},
	}
}

func UpdateDeadline(kubeclient *kubernetes.Clientset, instance_name string, deadline int64) error {
	pod, err := kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Get(context.TODO(), instance_name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	newPod := pod.DeepCopy()
	ann := newPod.ObjectMeta.Annotations
	if ann == nil {
		ann = make(map[string]string)
	}
	ann["deadline"] = strconv.Itoa(int(deadline))
	newPod.ObjectMeta.Annotations = ann

	_, err = kubeclient.CoreV1().Pods(newPod.ObjectMeta.Namespace).Update(context.TODO(), newPod, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
