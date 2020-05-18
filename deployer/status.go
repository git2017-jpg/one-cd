package deployer

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	v1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
)

type statusPrinter = func(cluster, namespace, deploymentName string, info string)

// WaitForPodContainersRunning ...
func (d *Deployer) WaitForPodContainersRunning(cluster, namespace, deploymentName string, threshold, checkInterval time.Duration,
	printer statusPrinter) (err error) {
	var (
		running    bool
		status     string
		deployment *v1.Deployment
		pods       []*coreV1.Pod
	)
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, string(debug.Stack()))
		}
	}()
	if printer == nil {
		err = errors.New("parameter printer is not allow nil")
		return
	}
	end := time.Now().Add(threshold)
	for {
		<-time.NewTimer(checkInterval).C
		deployment, err = d.Deployment(cluster, namespace, deploymentName)
		if err != nil {
			status = fmt.Sprintf("deployments.apps %s not found", deploymentName)
			printer(cluster, namespace, deploymentName, status)
			return
		}
		status = d.deploymentStatus(deployment)
		printer(cluster, namespace, deploymentName, status)
		if pods, err = d.PodList(cluster, namespace, deploymentName); err != nil {
			status = fmt.Sprintf("get PodList %s failed", deploymentName)
			printer(cluster, namespace, deploymentName, status)
			return
		}
		if running, err = d.podContainersRunning(pods); err != nil {
			break
		}
		if running {
			status = fmt.Sprintf("deployment %q successfully rolled out", deploymentName)
			printer(cluster, namespace, deploymentName, status)
			return
		}
		if time.Now().After(end) {
			status = "error: timed out waiting for the condition"
			err = errors.New(status)
			printer(cluster, namespace, deploymentName, status)
			return
		}
	}
	return
}

func (d *Deployer) deploymentStatus(deployment *v1.Deployment) (status string) {
	if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
		status = fmt.Sprintf("Waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated...",
			deployment.Name, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas)
		return
	}
	if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
		status = fmt.Sprintf("Waiting for deployment %q rollout to finish: %d old replicas are pending termination...",
			deployment.Name, deployment.Status.Replicas-deployment.Status.UpdatedReplicas)
		return
	}
	if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
		status = fmt.Sprintf("Waiting for deployment %q rollout to finish: %d of %d updated replicas are available...",
			deployment.Name, deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas)
		return
	}
	status = "Waiting for deployment spec update to be observed..."
	return
}

func (d *Deployer) podContainersRunning(pods []*coreV1.Pod) (running bool, err error) {
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
