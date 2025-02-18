package instance

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteInstance(c *fiber.Ctx, chall_id int, teamid int64) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if _, err := database.IsRunning(ctx, chall_id, teamid); err != nil {
		return err
	}

	instance_name := utils.GetInstanceName(chall_id, teamid)
	k8sclient, err := NewK8sClient()
	if err != nil {
		log.Println(err)
		return errors.New("error in stopping the instance, contact admin")
	}

	if err := deleteDeployment(instance_name, c.Context(), k8sclient); err != nil {
		log.Println("error in stop (deployment) - chall_id: ", chall_id, " teamid: ", teamid)
		return errors.New("error in stopping the instance, contact admin")
	}

	if err := deleteService(fmt.Sprintf("%s-svc", instance_name), c.Context(), k8sclient); err != nil {
		log.Println("error in stop (service) - chall_id: ", chall_id, " teamid: ", teamid)
		return errors.New("error in stopping the instance, contact admin")
	}

	if err := deleteIngress(fmt.Sprintf("%s-ingress", instance_name), "IngressRoute", c.Context(), k8sclient); err != nil {
		log.Println("error in stop (ingress) - chall_id: ", chall_id, " teamid: ", teamid)
		return errors.New("error in stopping the instance, contact admin")
	}

	if err := database.DeleteFlag(c, chall_id, teamid); err != nil {
		return err
	}

	if err := database.DeleteRunning(c, chall_id, teamid); err != nil {
		return err
	}

	return nil
}

func deleteDeployment(instance_name string, ctx context.Context, k8sclient *K8sClient) error {
	deployment, err := NewK8sDeployment(k8sclient, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	err = deployment.Delete(ctx, instance_name, config.INSTANCE_NAMESPACE)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			log.Println(err)
			return err
		}
	}

	for {
		_, err := deployment.Get(ctx, instance_name, config.INSTANCE_NAMESPACE)
		if err != nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}

func deleteService(instance_name string, ctx context.Context, k8sclient *K8sClient) error {
	service, err := NewK8sService(k8sclient, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	err = service.Delete(ctx, instance_name, config.INSTANCE_NAMESPACE)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			log.Println(err)
			return err
		}
	}

	for {
		_, err := service.Get(ctx, instance_name, config.INSTANCE_NAMESPACE)
		if err != nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}

func deleteIngress(instance_name string, kind string, ctx context.Context, k8sclient *K8sClient) error {
	ingress := NewK8sResource(k8sclient)

	err := ingress.Delete(ctx, instance_name, config.INSTANCE_NAMESPACE, "traefik.io", "v1alpha1", kind)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			log.Println(err)
			return err
		}
	}

	for {
		_, err := ingress.Get(ctx, instance_name, config.INSTANCE_NAMESPACE, "traefik.io", "v1alpha1", kind)
		if err != nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}
