package deployer

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// WaitForPodContainersRunning ...
func (d *Deployer) WaitForPodContainersRunning(cluster string, namespace string, deploymentName string, threshold, checkInterval time.Duration) error {
	end := time.Now().Add(threshold)
	for true {
		<-time.NewTimer(checkInterval).C
		running, err := d.podContainersRunning(cluster, namespace, deploymentName)
		if running {
			return nil
		}
		if err != nil {
			println(fmt.Sprintf("Encountered an error checking for running pods: %s", err))
		}
		if time.Now().After(end) {
			return fmt.Errorf("Failed to get all running containers")
		}
	}
	return nil
}

func (d *Deployer) podContainersRunning(cluster string, namespace string, deploymentName string) (running bool, err error) {
	var (
		client *kubernetes.Clientset
	)
	if client, err = d.Client(cluster); err != nil {
		return
	}
	pods, err := client.CoreV1().Pods(namespace).List(d.ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deploymentName),
	})
	if err != nil {
		return false, err
	}

	for _, item := range pods.Items {
		fmt.Println(item.Status.Message)
		for _, status := range item.Status.ContainerStatuses {
			if !status.Ready {
				return false, nil
			}
		}
	}
	return true, nil
}
