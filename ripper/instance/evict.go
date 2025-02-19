package instance

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/ripper/config"
	"github.com/TheAlpha16/isolet/ripper/database"
	"github.com/TheAlpha16/isolet/ripper/models"
)

func EvictInstance(job *models.Job, k8sclient *K8sClient) error {
	if err := deleteIngress(fmt.Sprintf("%s-ingress", job.InstanceName), "IngressRoute", job.Ctx, k8sclient); err != nil {
		log.Printf("[ERROR] (teamid: %d, chall_id: %d) (ingress) %s\n", job.TeamID, job.ChallID, err.Error())
	}

	if err := deleteService(fmt.Sprintf("%s-svc", job.InstanceName), job.Ctx, k8sclient); err != nil {
		log.Printf("[ERROR] (teamid: %d, chall_id: %d) (service) %s\n", job.TeamID, job.ChallID, err.Error())
	}

	if err := deleteDeployment(job.InstanceName, job.Ctx, k8sclient); err != nil {
		log.Printf("[ERROR] (teamid: %d, chall_id: %d) (deployment) %s\n", job.TeamID, job.ChallID, err.Error())
	}

	if err := database.DeleteFlag(job.TeamID, job.ChallID); err != nil {
		log.Printf("[ERROR] (teamid: %d, chall_id: %d) (flag) %s\n", job.TeamID, job.ChallID, err.Error())
	}

	if err := database.DeleteRunning(job.TeamID, job.ChallID); err != nil {
		log.Printf("[ERROR] (teamid: %d, chall_id: %d) (running) %s\n", job.TeamID, job.ChallID, err.Error())
	}

	return nil
}

func deleteDeployment(instance_name string, ctx context.Context, k8sclient *K8sClient) error {
	deployment, err := NewK8sDeployment(k8sclient, nil)
	if err != nil {
		return err
	}

	err = deployment.Delete(ctx, instance_name, config.INSTANCE_NAMESPACE)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
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
		return err
	}

	err = service.Delete(ctx, instance_name, config.INSTANCE_NAMESPACE)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
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
