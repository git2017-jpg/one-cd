package deployer

import (
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Ingress 获取ingress
func (d *Deployer) Ingress(cluster string, namespace string, name string) (ingress *v1beta1.Ingress, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	if ingress, err = client.ExtensionsV1beta1().Ingresses(namespace).Get(d.ctx, name, metav1.GetOptions{}); err != nil {
		return
	}
	return
}
