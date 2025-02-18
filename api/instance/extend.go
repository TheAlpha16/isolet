package instance

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/models"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func AddTime(c *fiber.Ctx, chall_id int, teamid int64, deadline *models.ExtendDeadline) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if _, err := database.IsRunning(ctx, chall_id, teamid); err != nil {
		return err
	}

	instance_name := utils.GetInstanceName(chall_id, teamid)
	k8sclient, err := NewK8sClient()
	if err != nil {
		log.Println(err)
		return errors.New("error in extending the deadline, please contact admin")
	}

	newdeadline, err := database.AddTime(c, chall_id, teamid)
	if err != nil {
		return err
	}

	err = updateDeadline(instance_name, newdeadline, c.Context(), k8sclient)
	if err != nil {
		log.Println(err)
		return errors.New("error in extending the deadline, please contact admin")
	}

	deadline.Deadline = newdeadline
	return nil
}

func updateDeadline(instance_name string, deadline int64, ctx context.Context, k8sclient *K8sClient) error {
	deployment, err := NewK8sDeployment(k8sclient, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	deploy, err := deployment.Get(ctx, instance_name, config.INSTANCE_NAMESPACE)
	if err != nil {
		log.Println(err)
		return err
	}

	deploy.ObjectMeta.Annotations["deadline"] =
		fmt.Sprintf("%d", deadline)

	deployment.Deployment = deploy

	err = deployment.Update(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
