package deployer

import (
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetScale ...
func (d *Deployer) GetScale(cluster string, namespace string, deploymentName string) (scale *autoscalingv1.Scale, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	if scale, err = client.AppsV1().Deployments(namespace).GetScale(d.ctx, deploymentName, metav1.GetOptions{}); err != nil {
		return
	}
	return
}

// UpdateScale 扩缩容
func (d *Deployer) UpdateScale(cluster string, namespace string, deploymentName string, replicas int32) (scale *autoscalingv1.Scale, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	scale = &autoscalingv1.Scale{
		Spec: autoscalingv1.ScaleSpec{
			Replicas: replicas},
	}
	scale.Name = deploymentName
	scale.Namespace = namespace
	if scale, err = client.AppsV1().Deployments(namespace).UpdateScale(d.ctx, deploymentName, scale, metav1.UpdateOptions{}); err != nil {
		return
	}
	return
}
