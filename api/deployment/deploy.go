package deployment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/models"
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

func DeployInstance(c *fiber.Ctx, chall_id int, teamid int64) (*models.Flag, error) {
	challenge := new(models.Challenge)
	image := new(models.Image)

	if err := database.ValidOnDemandChallenge(c, chall_id, teamid, challenge, image); err != nil {
		return nil, err
	}

	if err := database.CanStartInstance(c, chall_id, teamid); err != nil {
		return nil, err
	}

	instance_name := utils.GetInstanceName(chall_id, teamid)

	flagObject := models.Flag{
		TeamID:   teamid,
		ChallID:  chall_id,
		Flag:     "",
		Password: database.GenerateRandom()[0:32],
		Port:     image.Port,
		Hostname: utils.GetHostName(chall_id, teamid),
		Deadline: 1893456000000,
	}

	if challenge.Flag != "" {
		flagObject.Flag = strings.TrimSuffix(challenge.Flag, "}")
		flagObject.Flag = flagObject.Flag + "_" + database.GenerateRandom()[0:16] + "}"
	} else {
		flagObject.Flag = config.CTF_NAME + "{" + database.GenerateRandom()[0:32] + "}"
	}

	kubeclient, err := GetKubeClient()
	if err != nil {
		log.Println(err)
		_ = database.DeleteRunning(c, chall_id, teamid)
		return nil, errors.New("error in deployment, please contact admin")
	}

	pod := getPodObject(instance_name, flagObject, image)
	pod, err = kubeclient.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			_ = database.DeleteRunning(c, chall_id, teamid)
			log.Println(err)
			return nil, errors.New("error in deployment, please contact admin")
		}
	}

	for {
		createdPod, err := kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Get(context.TODO(), instance_name, metav1.GetOptions{})
		if err != nil {
			_ = database.DeleteRunning(c, chall_id, teamid)
			log.Println(err)
			return nil, errors.New("error in deployment, please contact admin")
		}

		if len(createdPod.Status.ContainerStatuses) > 0 {
			if createdPod.Status.ContainerStatuses[0].State.Waiting != nil && (createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "CrashLoopBackOff" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "ErrImagePull" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "ImagePullBackOff" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "CreateContainerConfigError" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "InvalidImageName" || createdPod.Status.ContainerStatuses[0].State.Waiting.Reason == "CreateContainerError") {
				kubeclient.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
				log.Printf("Error in launch: chall_id-%d reason: %s", chall_id, createdPod.Status.ContainerStatuses[0].State.Waiting.Reason)
				_ = database.DeleteRunning(c, chall_id, teamid)
				return nil, errors.New("error in deployment, please contact admin")
			}
		}

		if createdPod.Status.Phase == "Running" && createdPod.Status.StartTime != nil {
			flagObject.Deadline = createdPod.Status.StartTime.Add(time.Minute * time.Duration(config.INSTANCE_TIME)).UnixMilli()
			break
		}
	}

	err = UpdateDeadline(kubeclient, instance_name, flagObject.Deadline)
	if err != nil {
		_ = DeleteInstance(c, chall_id, teamid)
		log.Println(err)
		return nil, errors.New("error in deployment, please contact admin")
	}

	svcConfig := getServiceObject(instance_name, flagObject)
	_, err = kubeclient.CoreV1().Services(svcConfig.Namespace).Create(context.TODO(), svcConfig, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		if !strings.Contains(err.Error(), "already exists") {
			return nil, errors.New("error in deployment, please contact admin")
		}
	}

	createdService, err := kubeclient.CoreV1().Services(svcConfig.Namespace).Get(context.TODO(), svcConfig.Name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, errors.New("error in deployment, please contact admin")
	}

	port := createdService.Spec.Ports[0].NodePort
	hostname := config.INSTANCE_HOSTNAME

	nodes, err := kubeclient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, errors.New("error in deployment, please contact admin")
	}

	nodeip := nodes.Items[0].Status.Addresses
	for i := 0; i < len(nodeip); i++ {
		if nodeip[i].Type == "ExternalIP" {
			hostname = nodeip[i].Address
			break
		}
	}

	flagObject.Port = int(port)
	flagObject.Hostname = hostname
	flagObject.Deployment = image.Deployment

	if err := database.NewFlag(c, &flagObject); err != nil {
		return nil, err
	}

	return &flagObject, nil
}

func DeleteInstance(c *fiber.Ctx, chall_id int, teamid int64) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if _, err := database.IsRunning(ctx, chall_id, teamid); err != nil {
		return err
	}

	instance_name := utils.GetInstanceName(chall_id, teamid)
	kubeclient, err := GetKubeClient()
	if err != nil {
		log.Println(err)
		return errors.New("error in deletion, please contact admin")
	}

	err = kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Delete(context.TODO(), instance_name, metav1.DeleteOptions{})
	if err != nil {

		if !strings.Contains(err.Error(), "not found") {
			log.Println(err)
			return errors.New("error in deletion, please contact admin")
		}

		err = kubeclient.CoreV1().Services(config.INSTANCE_NAMESPACE).Delete(context.TODO(), fmt.Sprintf("svc-%s", instance_name), metav1.DeleteOptions{})
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				log.Println(err)
				return errors.New("error in deletion, please contact admin")
			}
		}

		if err := database.DeleteFlag(c, chall_id, teamid); err != nil {
			return err
		}

		if err := database.DeleteRunning(c, chall_id, teamid); err != nil {
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
			return errors.New("error in deletion, please contact admin")
		}
	}

	if err := database.DeleteFlag(c, chall_id, teamid); err != nil {
		return err
	}

	if err := database.DeleteRunning(c, chall_id, teamid); err != nil {
		return err
	}

	return nil
}

func AddTime(c *fiber.Ctx, chall_id int, teamid int64, deadline *models.ExtendDeadline) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if _, err := database.IsRunning(ctx, chall_id, teamid); err != nil {
		return err
	}

	instance_name := utils.GetInstanceName(chall_id, teamid)
	kubeclient, err := GetKubeClient()
	if err != nil {
		log.Println(err)
		return errors.New("error in extension, please contact admin")
	}

	_, err = kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Get(c.Context(), instance_name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			log.Println(err)
			return errors.New("error in extension, please contact admin")
		}

		if err := database.DeleteFlag(c, chall_id, teamid); err != nil {
			return errors.New("instance not running")
		}

		if err := database.DeleteRunning(c, chall_id, teamid); err != nil {
			return errors.New("instance not running")
		}

		return errors.New("instance not running")
	}

	newdeadline, err := database.AddTime(c, chall_id, teamid)
	if err != nil {
		return err
	}

	err = UpdateDeadline(kubeclient, instance_name, newdeadline)
	if err != nil {
		log.Println(err)
		return errors.New("error in extension, please contact admin")
	}

	deadline.Deadline = newdeadline
	return nil
}

func getPodObject(instance_name string, flagObject models.Flag, image *models.Image) *core.Pod {
	var imagePath string
	var cpu string
	var memory string

	imagePath = config.IMAGE_REGISTRY
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
				"chall_id": fmt.Sprintf("%d", flagObject.ChallID),
				"teamid":   fmt.Sprintf("%d", flagObject.TeamID),
				"app":      "instance",
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
							Name:  "USERNAME",
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
				"chall_id": fmt.Sprintf("%d", flagObject.ChallID),
				"teamid":   fmt.Sprintf("%d", flagObject.TeamID),
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
				"chall_id": fmt.Sprintf("%d", flagObject.ChallID),
				"teamid":   fmt.Sprintf("%d", flagObject.TeamID),
				"app":      "instance",
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
