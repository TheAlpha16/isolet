package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CyberLabs-Infosec/isolet/ripper/config"
	"github.com/CyberLabs-Infosec/isolet/ripper/database"
	"github.com/CyberLabs-Infosec/isolet/ripper/utils"

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
		log.Printf("%s - ERROR: userid=%d, level=%d %s\n", time.Now().Format(time.UnixDate), userid, level, err.Error())
		return err
	}

	err = kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).Delete(context.TODO(), instance_name, metav1.DeleteOptions{})
	if err != nil {

		if !strings.Contains(err.Error(), "not found") {
			log.Printf("%s - ERROR: userid=%d, level=%d %s\n", time.Now().Format(time.UnixDate), userid, level, err.Error())
			return err
		}

		err = kubeclient.CoreV1().Services(config.INSTANCE_NAMESPACE).Delete(context.TODO(), fmt.Sprintf("svc-%s", instance_name), metav1.DeleteOptions{})
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				log.Printf("%s - ERROR: userid=%d, level=%d %s\n", time.Now().Format(time.UnixDate), userid, level, err.Error())
				return err
			}
		}

		if err := database.DeleteFlag(userid, level); err != nil {
			log.Printf("%s - ERROR: userid=%d, level=%d %s\n", time.Now().Format(time.UnixDate), userid, level, err.Error())
			return err
		}

		if err := database.DeleteRunning(userid, level); err != nil {
			log.Printf("%s - ERROR: userid=%d, level=%d %s\n", time.Now().Format(time.UnixDate), userid, level, err.Error())
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
			log.Printf("%s - ERROR: userid=%d, level=%d %s\n", time.Now().Format(time.UnixDate), userid, level, err.Error())
			return err
		}
	}

	if err := database.DeleteFlag(userid, level); err != nil {
		log.Printf("%s - ERROR: userid=%d, level=%d %s\n", time.Now().Format(time.UnixDate), userid, level, err.Error())
		return err
	}

	if err := database.DeleteRunning(userid, level); err != nil {
		log.Printf("%s - ERROR: userid=%d, level=%d %s\n", time.Now().Format(time.UnixDate), userid, level, err.Error())
		return err
	}

	return nil
}

func EvictPods() error {
	kubeclient, err := GetKubeClient()
	if err != nil {
		log.Printf("%s - ERROR: %s\n", time.Now().Format(time.UnixDate), err.Error())
		return err
	}

	pods, err := kubeclient.CoreV1().Pods(config.INSTANCE_NAMESPACE).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("%s - ERROR: %s\n", time.Now().Format(time.UnixDate), err.Error())
		return err
	}

	for i := 0; i < len(pods.Items); i++ {
		pod := pods.Items[i]
		userid, _ := strconv.Atoi(pod.Labels["userid"])
		level, _ := strconv.Atoi(pod.Labels["level"])

		// remove failed and unknown pods
		// if pod.Status.Phase == core.PodFailed {
		// 	log.Printf("%s - DELETE: userid=%d, level=%d started=%s\n", time.Now().Format(time.UnixDate), userid, level, startTime.Format(time.UnixDate))
		// 	_ = DeleteInstance(userid, level)
		// 	continue
		// }

		// if pod.Status.Phase == core.PodUnknown {
		// 	log.Printf("%s - DELETE: userid=%d, level=%d started=%s\n", time.Now().Format(time.UnixDate), userid, level, startTime.Format(time.UnixDate))
		// 	_ = DeleteInstance(userid, level)
		// 	continue
		// }

		// if pod.Status.Phase == core.PodSucceeded {
		// 	log.Printf("%s - DELETE: userid=%d, level=%d started=%s\n", time.Now().Format(time.UnixDate), userid, level, startTime.Format(time.UnixDate))
		// 	_ = DeleteInstance(userid, level)
		// 	continue
		// }

		if pod.Status.Phase != core.PodRunning {
			log.Println("continuing..")
			continue
		}

		startTime := pod.Status.StartTime.Time
		if startTime.Add(time.Minute * time.Duration(config.INSTANCE_TIME)).After(metav1.Now().Time) {
			log.Printf("%s - DELETE: userid=%d, level=%d started=%s\n", time.Now().Format(time.UnixDate), userid, level, startTime.Format(time.UnixDate))
			_ = DeleteInstance(userid, level)
		}
	}

	return nil
}

func main() {
	log.Printf("%s - LOG: Starting ripper...\n", time.Now().Format(time.UnixDate))
	log.Printf("%s - LOG: Connecting to DB...\n", time.Now().Format(time.UnixDate))

	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s - LOG: DB connection established\n", time.Now().Format(time.UnixDate))

	for {
		err := EvictPods()
		if err != nil {
			log.Printf("%s - ERROR: %s\n", time.Now().Format(time.UnixDate), err.Error())
		}
		// time.Sleep(time.Minute * time.Duration(1))
	}
}
