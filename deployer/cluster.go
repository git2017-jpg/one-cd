package deployer

import (
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// AddCluster 增加集群
func (d *Deployer) AddCluster(cluster, kubeConfig string) (client *kubernetes.Clientset, err error) {
	var conf *rest.Config
	if conf, err = clientcmd.BuildConfigFromFlags("", kubeConfig); err != nil {
		log.Println("BuildConfigFromFlags", err)
		return
	}
	if client, err = kubernetes.NewForConfig(conf); err != nil {
		log.Println("NewForConfig", err)
		return
	}
	d.clients.Store(cluster, client)
	return
}

// Client ...
func (d *Deployer) Client(cluster string) (client *kubernetes.Clientset, err error) {
	v, ok := d.clients.Load(cluster)
	if !ok {
		err = fmt.Errorf("无效的或未启用的集群 %s", cluster)
		return
	}
	client = v.(*kubernetes.Clientset)
	return
}

// Clusters ...
func (d *Deployer) Clusters() (list []string) {
	d.clients.Range(func(k, v interface{}) bool {
		list = append(list, k.(string))
		return true
	})
	return
}
