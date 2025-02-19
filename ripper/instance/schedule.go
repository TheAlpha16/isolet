package instance

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/TheAlpha16/isolet/ripper/config"
	"github.com/TheAlpha16/isolet/ripper/models"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var JobChan = make(chan models.Job, 100)

func StartWatch() {
	k8sclient, err := NewK8sClient()
	if err != nil {
		log.Fatalf("[ERROR] unable to get k8s client: %s", err.Error())
	}

	InitWorkers(k8sclient)

	for {
		deployments, err := k8sclient.ClientSet.AppsV1().Deployments(config.INSTANCE_NAMESPACE).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Printf("[ERROR] could not list deployments: %s\n", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		for _, deployment := range deployments.Items {
			ProcessDeadline(deployment, k8sclient)
		}

		time.Sleep(5 * time.Second)
	}
}

func ProcessDeadline(deployment appsv1.Deployment, k8sclient *K8sClient) {
	deadline_string := deployment.Annotations["deadline"]
	if deadline_string == "" {
		return
	}

	deadline, err := strconv.ParseInt(deadline_string, 10, 64)
	if err != nil {
		return
	}

	endTime := time.UnixMilli(deadline)
	if !time.Now().Before(endTime) {
		teamid_int, _ := strconv.Atoi(deployment.Labels["teamid"])
		teamid := int64(teamid_int)
		chall_id, _ := strconv.Atoi(deployment.Labels["chall_id"])

		log.Printf("[LOG] deadline reached -  (teamid: %d, chall_id: %d)\n", teamid, chall_id)

		JobChan <- models.Job{
			TeamID:       teamid,
			ChallID:      chall_id,
			InstanceName: deployment.Name,
			Ctx:          context.Background(),
		}
	}
}

func InitWorkers(k8sclient *K8sClient) {
	workers, err := strconv.Atoi(config.WORKER_COUNT)
	if err != nil {
		workers = 5
	}

	for i := 0; i < workers; i++ {
		go Worker(i, k8sclient)
	}
}

func Worker(id int, k8sclient *K8sClient) {
	log.Printf("[LOG] Worker %d started\n", id)
	for job := range JobChan {
		EvictInstance(&job, k8sclient)
	}

	log.Printf("[LOG] Worker %d stopped\n", id)
}
