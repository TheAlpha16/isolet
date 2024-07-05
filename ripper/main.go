package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/ripper/config"
	"github.com/TheAlpha16/isolet/ripper/database"
	"github.com/TheAlpha16/isolet/ripper/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

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

func DeleteInstance(userid int, level int) error {
	instance_name := utils.GetInstanceName(userid, level)
	kubeclient, err := GetKubeClient()
	if err != nil {
		log.Printf("ERROR: userid=%d, level=%d %s\n", userid, level, err.Error())
		return err
	}

	err = kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Delete(context.TODO(), instance_name, metav1.DeleteOptions{})
	if err != nil {

		if !strings.Contains(err.Error(), "not found") {
			log.Printf("ERROR: userid=%d, level=%d %s\n", userid, level, err.Error())
			return err
		}

		err = kubeclient.CoreV1().Services(config.INSTANCE_NAMESPACE).Delete(context.TODO(), fmt.Sprintf("svc-%s", instance_name), metav1.DeleteOptions{})
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				log.Printf("ERROR: userid=%d, level=%d %s\n", userid, level, err.Error())
				return err
			}
		}

		if err := database.DeleteFlag(userid, level); err != nil {
			log.Printf("ERROR: userid=%d, level=%d %s\n", userid, level, err.Error())
			return err
		}

		if err := database.DeleteRunning(userid, level); err != nil {
			log.Printf("ERROR: userid=%d, level=%d %s\n", userid, level, err.Error())
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
			log.Printf("ERROR: userid=%d, level=%d %s\n", userid, level, err.Error())
			return err
		}
	}

	if err := database.DeleteFlag(userid, level); err != nil {
		log.Printf("ERROR: userid=%d, level=%d %s\n", userid, level, err.Error())
		return err
	}

	if err := database.DeleteRunning(userid, level); err != nil {
		log.Printf("ERROR: userid=%d, level=%d %s\n", userid, level, err.Error())
		return err
	}

	return nil
}

func EvictPods() error {
	kubeclient, err := GetKubeClient()
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return err
	}

	pods, err := kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return err
	}

	for _, pod := range pods.Items {
		userid, _ := strconv.Atoi(pod.Labels["userid"])
		level, _ := strconv.Atoi(pod.Labels["level"])

		if pod.Status.Phase == core.PodPending {
			continue
		}

		ann := pod.Annotations["deadline"]
		if ann == "" {
			continue
		}

		deadline, err := strconv.ParseInt(ann, 10, 64)
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
			return err
		}

		endTime := time.UnixMilli(deadline)

		if !time.Now().Before(endTime) {
			log.Printf("DELETE: userid=%d, level=%d", userid, level)
			_ = DeleteInstance(userid, level)
		}
	}

	return nil
}

func main() {
	log.Printf("LOG: Connecting to DB...\n")

	for {
		if err := database.Connect(); err != nil {
			log.Println(err)
			log.Println("sleep for 1 minute")
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	log.Printf("LOG: DB connection established\n")

	for {
		err := EvictPods()
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
		}
	}
}
