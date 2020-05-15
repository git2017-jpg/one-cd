package deployer

import (
	"errors"
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
)

// WaitForPodContainersRunning ...
func (d *Deployer) WaitForPodContainersRunning(cluster string, namespace string, deploymentName string, threshold, checkInterval time.Duration) (err error) {
	var (
		running bool
		info    string
	)
	end := time.Now().Add(threshold)
	for true {
		<-time.NewTimer(checkInterval).C
		running, info, err = d.podContainersRunning(cluster, namespace, deploymentName)
		fmt.Println(info)
		if err != nil {
			break
		}
		if running {
			info = fmt.Sprintf("deployment %q successfully rolled out\n", deploymentName)
			fmt.Println(info)
			return
		}
		if time.Now().After(end) {
			info = "error: timed out waiting for the condition\n"
			err = errors.New(info)
			fmt.Println(info)
			return
		}
	}
	return
}

func (d *Deployer) podContainersRunning(cluster string, namespace string, deploymentName string) (running bool, info string, err error) {
	var (
		deployment *v1.Deployment
		pods       []*coreV1.Pod
	)
	if deployment, err = d.Deployment(cluster, namespace, deploymentName); err != nil {
		info = fmt.Sprintf("get deployment %s failed", deploymentName)
		return
	}
	if pods, err = d.PodList(cluster, namespace, deploymentName); err != nil {
		info = fmt.Sprintf("get PodList %s failed", deploymentName)
		return
	}
	if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
		info = fmt.Sprintf("Waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated...", deployment.Name, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas)
	}
	if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
		info = fmt.Sprintf("Waiting for deployment %q rollout to finish: %d old replicas are pending termination...", deployment.Name, deployment.Status.Replicas-deployment.Status.UpdatedReplicas)
	}
	if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
		info = fmt.Sprintf("Waiting for deployment %q rollout to finish: %d of %d updated replicas are available...", deployment.Name, deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas)
	}
	for _, item := range pods {
		for _, status := range item.Status.ContainerStatuses {
			if !status.Ready {
				running = false
				return
			}
		}
	}
	running = true
	return
}
