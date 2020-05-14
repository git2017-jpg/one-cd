package deployer

import (
	"encoding/json"
	"errors"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

// Deployment 获取Deployment信息
func (d *Deployer) Deployment(cluster string, namespace string, deploymentName string) (deployment *v1.Deployment, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	if deployment, err = client.AppsV1().Deployments(namespace).Get(d.ctx, deploymentName, metav1.GetOptions{}); err != nil {
		return
	}
	return
}

// DeploymentDelete 删除Deployment
func (d *Deployer) DeploymentDelete(cluster string, namespace string, deploymentName string) (err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	if err = client.AppsV1().Deployments(namespace).Delete(d.ctx, deploymentName, metav1.DeleteOptions{}); err != nil {
		return
	}
	return
}

// DeploymentEvents 获取Deployment事件
func (d *Deployer) DeploymentEvents(cluster string, namespace string, deploymentName string) (list *coreV1.EventList, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	fieldSelector := fmt.Sprintf("involvedObject.kind=Deployment,involvedObject.name=%s", deploymentName)
	if list, err = client.CoreV1().Events(namespace).List(d.ctx, metav1.ListOptions{FieldSelector: fieldSelector}); err != nil {
		return
	}

	return
}

// Deploy ...
func (d *Deployer) Deploy(yaml string) (deployment *v1.Deployment, err error) {
	var (
		data   []byte
		client *kubernetes.Clientset
	)
	if data, err = yaml2.ToJSON([]byte(yaml)); err != nil {
		return
	}
	deployment = &v1.Deployment{}
	if err = json.Unmarshal(data, deployment); err != nil {
		return
	}
	cluster := deployment.ObjectMeta.ClusterName
	namespace := deployment.ObjectMeta.Namespace
	if cluster == "" {
		cluster = "default"
	}
	if namespace == "" {
		namespace = "default"
	}
	deploymentName := deployment.ObjectMeta.Name
	if client, err = d.Client(cluster); err != nil {
		return
	}
	if _, err = d.Deployment(cluster, namespace, deploymentName); err != nil {
		if deployment, err = client.AppsV1().Deployments(namespace).Create(d.ctx, deployment, metav1.CreateOptions{}); err != nil {
			return
		}
	} else {
		if deployment, err = client.AppsV1().Deployments(namespace).Update(d.ctx, deployment, metav1.UpdateOptions{}); err != nil {
			return
		}
	}
	return
}

// Update 更新镜像版本
func (d *Deployer) Update(cluster string, namespace string, deploymentName, image string) (deployment *v1.Deployment, err error) {
	var client *kubernetes.Clientset
	if deployment, err = d.Deployment(cluster, namespace, deploymentName); err != nil {
		return
	}
	if client, err = d.Client(cluster); err != nil {
		return
	}
	deployment.Spec.Template.Spec.Containers[0].Image = image
	if deployment, err = client.AppsV1().Deployments(namespace).Update(d.ctx, deployment, metav1.UpdateOptions{}); err != nil {
		return
	}
	return
}

// RollBack 回滚需要指定版本
func (d *Deployer) RollBack(cluster string, namespace string, deploymentName, rs string) (deployment *v1.Deployment, err error) {
	var (
		client         *kubernetes.Clientset
		replicaSetList *v1.ReplicaSetList
	)
	if replicaSetList, err = d.ReplicaSetList(cluster, namespace, deploymentName); err != nil {
		return
	}
	if len(replicaSetList.Items) <= 1 {
		err = errors.New("回滚未执行，没有可回滚的版本")
		return
	}
	if deployment, err = d.Deployment(cluster, namespace, deploymentName); err != nil {
		return
	}
	if client, err = d.Client(cluster); err != nil {
		return
	}
	for _, v := range replicaSetList.Items {
		if v.ObjectMeta.Name == rs {
			deployment.Spec.Template = v.Spec.Template
			if deployment, err = client.AppsV1().Deployments(namespace).Update(d.ctx, deployment, metav1.UpdateOptions{}); err != nil {
				return
			}
			return
		}
	}
	err = errors.New("回滚未执行，没有找到指定的版本")
	return
}

// ReplicaSetList 获取rs列表
func (d *Deployer) ReplicaSetList(cluster string, namespace string, deploymentName string) (replicaSetList *v1.ReplicaSetList, err error) {
	var (
		client *kubernetes.Clientset
	)
	if client, err = d.Client(cluster); err != nil {
		return
	}
	labelSelector := fmt.Sprintf("app=%s", deploymentName)
	if replicaSetList, err = client.AppsV1().ReplicaSets(namespace).List(d.ctx, metav1.ListOptions{LabelSelector: labelSelector}); err != nil {
		return
	}
	return
}
