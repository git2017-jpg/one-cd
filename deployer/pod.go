package deployer

import (
	"bytes"
	"fmt"
	"io"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodEvents 获取pod事件
func (d *Deployer) PodEvents(cluster string, namespace string, podName string) (list *coreV1.EventList, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	fieldSelector := fmt.Sprintf("involvedObject.kind=Pod,involvedObject.name=%s", podName)
	if list, err = client.CoreV1().Events(namespace).List(d.ctx, metav1.ListOptions{FieldSelector: fieldSelector}); err != nil {
		return
	}
	return
}

// PodList ...
func (d *Deployer) PodList(cluster string, namespace string, deploymentName string) (podList *coreV1.PodList, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	labelSelector := fmt.Sprintf("app=%s", deploymentName)
	if podList, err = client.CoreV1().Pods(namespace).List(d.ctx, metav1.ListOptions{LabelSelector: labelSelector}); err != nil {
		return
	}
	return
}

// PodLog 获取pod日志
func (d *Deployer) PodLog(cluster string, namespace string, podName, container string, sinceSeconds int64, previous bool) (log string, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	req := client.CoreV1().Pods(namespace).GetLogs(podName, &coreV1.PodLogOptions{
		Container:    container,
		SinceSeconds: &sinceSeconds,
		Previous:     previous})
	podLogs, err := req.Stream(d.ctx)
	if err != nil {
		return
	}
	defer podLogs.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, podLogs); err != nil {
		return
	}
	log = buf.String()
	return
}

// PodDelete 删除pod
func (d *Deployer) PodDelete(cluster string, namespace string, podName string) (err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	if err = client.CoreV1().Pods(namespace).Delete(d.ctx, podName, metav1.DeleteOptions{}); err != nil {
		return
	}
	return
}
