package deployer

import (
	"errors"
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
)

type statusPrinter = func(cluster, namespace, deploymentName string, info string)

// WaitForPodContainersRunning ...
func (d *Deployer) WaitForPodContainersRunning(cluster, namespace, deploymentName string, threshold, checkInterval time.Duration,
	printer statusPrinter) (err error) {
	var (
		ready      bool
		status     string
		deployment *v1.Deployment
	)
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
		status, ready = d.deploymentStatus(deployment)
		printer(cluster, namespace, deploymentName, status)
		if ready {
			return
		}
		if time.Now().After(end) {
			status = "error: timed out waiting for the condition"
			err = errors.New(status)
			printer(cluster, namespace, deploymentName, status)
			return
		}
	}
}

func (d *Deployer) deploymentStatus(deployment *v1.Deployment) (string, bool) {
	if deployment.Generation <= deployment.Status.ObservedGeneration {
		if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
			return fmt.Sprintf("Waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated...",
				deployment.Name, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas), false
		}
		if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
			return fmt.Sprintf("Waiting for deployment %q rollout to finish: %d old replicas are pending termination...",
				deployment.Name, deployment.Status.Replicas-deployment.Status.UpdatedReplicas), false
		}
		if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
			return fmt.Sprintf("Waiting for deployment %q rollout to finish: %d of %d updated replicas are available...",
				deployment.Name, deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas), false
		}
		return fmt.Sprintf("deployment %q successfully rolled out\n", deployment.Name), true
	}
	return "Waiting for deployment spec update to be observed...", false
}
