package instance

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/models"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func DeployInstance(
	c *fiber.Ctx,
	chall_id int,
	teamid int64,
) (*models.Flag, error) {
	challenge := new(models.Challenge)

	if err := database.ValidOnDemandChallenge(c, chall_id, teamid, challenge); err != nil {
		return nil, err
	}

	if err := database.CanStartInstance(c, chall_id, teamid); err != nil {
		return nil, err
	}

	instance_name := utils.GetInstanceName(chall_id, teamid)
	challenge_name := utils.GetChallengeSubdomain(challenge.Name)

	flagObject := models.Flag{
		TeamID:   teamid,
		ChallID:  chall_id,
		Flag:     challenge.Flag,
		Password: "",
		Port:     challenge.Port,
		Hostname: utils.GetHostName([]string{instance_name}),
		Deadline: 1893456000000,
		Deployment: challenge.Deployment,
	}

	// if challenge.Flag != "" {
	// 	flagObject.Flag = strings.TrimSuffix(challenge.Flag, "}")
	// 	flagObject.Flag = fmt.Sprintf("%s_%s}", flagObject.Flag, database.GenerateRandom()[0:16])
	// } else {
	// 	flagObject.Flag = fmt.Sprintf("%s{%s}", flagObject.Flag, database.GenerateRandom()[0:16])
	// }

	k8sclient, err := NewK8sClient()
	if err != nil {
		log.Println(err)
		_ = database.DeleteRunning(c, chall_id, teamid)
		return nil, errors.New("error in deployment, please contact admin")
	}

	configMap, err := getConfigMap(challenge_name, c.Context(), k8sclient)
	if err != nil {
		log.Println(err)
		_ = database.DeleteRunning(c, chall_id, teamid)
		return nil, errors.New("error in deployment, please contact admin")
	}

	yamlData := configMap.Data["deployment.yaml"]
	err = createHandler(yamlData, instance_name, &flagObject, c, k8sclient)
	if err != nil {
		_ = database.DeleteRunning(c, chall_id, teamid)
		return nil, errors.New("error in deployment, please contact admin")
	}

	if err := database.NewFlag(c, &flagObject); err != nil {
		_ = DeleteInstance(c, flagObject.ChallID, flagObject.TeamID)
		return nil, err
	}

	return &flagObject, nil
}

func getConfigMap(challenge_name string, ctx context.Context, k8sclient *K8sClient) (*core.ConfigMap, error) {
	return k8sclient.ClientSet.CoreV1().
		ConfigMaps("store").
		Get(ctx, fmt.Sprintf("%s-cm", challenge_name), metav1.GetOptions{})
}

func createHandler(yamlData string, instance_name string, flagObject *models.Flag, c *fiber.Ctx, k8sclient *K8sClient) error {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(yamlData)), len(yamlData))

	for {
		var obj unstructured.Unstructured
		err := decoder.Decode(&obj)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return err
			}
			break
		}

		switch obj.GetKind() {
		case "Deployment":
			err = createDeployment(&obj, instance_name, flagObject, c.Context(), k8sclient)
			if err != nil {
				log.Println("error in launch - deployment: ", obj.GetName())
				return err
			}

			err = updateDeadline(instance_name, flagObject.Deadline, c.Context(), k8sclient)
			if err != nil {
				_ = DeleteInstance(c, flagObject.ChallID, flagObject.TeamID)
				log.Println("error in launch - deadline: ", obj.GetName())
				return err
			}
		case "Service":
			err = createService(&obj, instance_name, flagObject, c.Context(), k8sclient)
			if err != nil {
				_ = DeleteInstance(c, flagObject.ChallID, flagObject.TeamID)
				log.Println("error in launch - service: ", obj.GetName())
				return err
			}
		default:
			err = createDefault(&obj, instance_name, flagObject, c.Context(), k8sclient)
			if err != nil {
				_ = DeleteInstance(c, flagObject.ChallID, flagObject.TeamID)
				log.Println("error in launch - default: ", obj.GetName())
				return err
			}
		}
	}

	return nil
}

func createDeployment(obj *unstructured.Unstructured, instance_name string, flagObject *models.Flag, ctx context.Context, k8sclient *K8sClient) error {
	deployment, err := NewK8sDeployment(k8sclient, obj)
	if err != nil {
		log.Println(err)
		return err
	}

	// change deployment name to instance name
	deployment.Deployment.Name = instance_name

	// add labels to identify the instance
	deployment.Deployment.Labels["chall_id"] = fmt.Sprintf("%d", flagObject.ChallID)
	deployment.Deployment.Labels["teamid"] = fmt.Sprintf("%d", flagObject.TeamID)
	deployment.Deployment.Spec.Template.Labels["chall_id"] = fmt.Sprintf("%d", flagObject.ChallID)
	deployment.Deployment.Spec.Template.Labels["teamid"] = fmt.Sprintf("%d", flagObject.TeamID)
	deployment.Deployment.Spec.Selector.MatchLabels["chall_id"] = fmt.Sprintf("%d", flagObject.ChallID)
	deployment.Deployment.Spec.Selector.MatchLabels["teamid"] = fmt.Sprintf("%d", flagObject.TeamID)

	// update environment variables
	// for i, container := range deployment.Deployment.Spec.Template.Spec.Containers {
	// 	envVars := container.Env
	// 	envVars = append(envVars, core.EnvVar{
	// 		Name:  "FLAG",
	// 		Value: flagObject.Flag,
	// 	})
	// 	envVars = append(envVars, core.EnvVar{
	// 		Name:  "USERNAME",
	// 		Value: config.DEFAULT_USERNAME,
	// 	})
	// 	envVars = append(envVars, core.EnvVar{
	// 		Name:  "PASSWORD",
	// 		Value: flagObject.Password,
	// 	})
	// 	envVars = append(envVars, core.EnvVar{
	// 		Name:  "CTF_NAME",
	// 		Value: config.CTF_NAME,
	// 	})

	// 	deployment.Deployment.Spec.Template.Spec.Containers[i].Env = envVars
	// }

	if deployment.Deployment.Spec.Replicas == nil || *deployment.Deployment.Spec.Replicas == 0 {
		deployment.Deployment.Spec.Replicas = utils.Int32Addr("1")
	}

	err = deployment.Create(ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Println(err)
			return err
		}
	}

	for {
		deploy, err := deployment.Get(ctx, deployment.Deployment.Name, deployment.Deployment.Namespace)
		if err != nil {
			log.Println(err)
			return err
		}

		if deploy.Status.ReadyReplicas == deploy.Status.Replicas && deploy.Status.Replicas > 0 {
			pods, err := k8sclient.ClientSet.CoreV1().
				Pods(deploy.Namespace).
				List(ctx, metav1.ListOptions{
					LabelSelector: fmt.Sprintf("teamid=%d,chall_id=%d", flagObject.TeamID, flagObject.ChallID),
				})
			if err != nil {
				log.Println(err)
				return err
			}

			allReady := true
			for _, pod := range pods.Items {
				for _, cs := range pod.Status.ContainerStatuses {
					if cs.State.Waiting != nil {
						reason := cs.State.Waiting.Reason
						if reason == "CrashLoopBackOff" ||
							reason == "ErrImagePull" ||
							reason == "ImagePullBackOff" ||
							reason == "CreateContainerConfigError" ||
							reason == "InvalidImageName" ||
							reason == "CreateContainerError" {
							_ = deployment.Delete(ctx, deploy.Name, deploy.Namespace)
							log.Printf("Error in launch: challenege: %s, reason: %s", deploy.Name, reason)
							return errors.New("unable to create deployment")
						}
					}
				}

				if pod.Status.Phase != core.PodRunning {
					allReady = false
					break
				}
			}

			if allReady {
				if len(pods.Items) == 0 {
					_ = deployment.Delete(ctx, deploy.Name, deploy.Namespace)
					log.Println("No pods found")
					return errors.New("no pods found")
				}

				flagObject.Deadline = pods.Items[0].Status.StartTime.Add(time.Minute * time.Duration(config.INSTANCE_TIME)).UnixMilli()
				break
			}
		}
	}

	return nil
}

func createService(obj *unstructured.Unstructured, instance_name string, flagObject *models.Flag, ctx context.Context, k8sclient *K8sClient) error {
	service, err := NewK8sService(k8sclient, obj)
	if err != nil {
		log.Println(err)
		return err
	}

	// update service name
	service.Service.Name = fmt.Sprintf("%s-svc", instance_name)

	// update service labels
	service.Service.Labels["chall_id"] = fmt.Sprintf("%d", flagObject.ChallID)
	service.Service.Labels["teamid"] = fmt.Sprintf("%d", flagObject.TeamID)

	// update service selector
	service.Service.Spec.Selector["chall_id"] = fmt.Sprintf("%d", flagObject.ChallID)
	service.Service.Spec.Selector["teamid"] = fmt.Sprintf("%d", flagObject.TeamID)

	err = service.Create(ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Println(err)
			return err
		}
	}

	return nil
}

func createDefault(obj *unstructured.Unstructured, instance_name string, flagObject *models.Flag, ctx context.Context, k8sclient *K8sClient) error {
	resource := NewK8sResource(k8sclient)

	// update object name
	obj.SetName(fmt.Sprintf("%s-ingress", instance_name))

	// update object labels
	currentLabels := obj.GetLabels()
	if currentLabels == nil {
		currentLabels = make(map[string]string)
	}
	currentLabels["chall_id"] = fmt.Sprintf("%d", flagObject.ChallID)
	currentLabels["teamid"] = fmt.Sprintf("%d", flagObject.TeamID)

	obj.SetLabels(currentLabels)

	// update entry point matches
	if obj.GetKind() == "IngressRoute" {
		specObj, ok := obj.Object["spec"].(map[string]interface{})
		if !ok {
			return errors.New("invalid spec in ingressroute")
		}

		routes, ok := specObj["routes"].([]interface{})
		if !ok {
			return errors.New("invalid routes in ingressroute")
		}

		if len(routes) == 0 {
			return errors.New("no routes set in ingressroute")
		}

		// update host match
		route_0, ok := routes[0].(map[string]interface{})
		if !ok {
			return errors.New("invalid route in ingressroute")
		}
		route_0["match"] = fmt.Sprintf("Host(`%s`)", flagObject.Hostname)

		// update service name
		services, ok := route_0["services"].([]interface{})
		if !ok {
			return errors.New("services not set in ingressroute")
		}

		if len(services) == 0 {
			return errors.New("no services set in ingressroute")
		}

		service_0, ok := services[0].(map[string]interface{})
		if !ok {
			return errors.New("invalid service in ingressroute")
		}
		service_0["name"] = fmt.Sprintf("%s-svc", instance_name)

		// cleanup
		services[0] = service_0
		route_0["services"] = services
		routes[0] = route_0
		specObj["routes"] = routes
		obj.Object["spec"] = specObj
	}

	err := resource.Create(ctx, obj)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Println(err)
			return err
		}
	}

	return nil
}
