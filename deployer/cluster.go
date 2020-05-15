package deployer

import (
	"fmt"
	"log"

	corev1informer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var stopChan = make(chan struct{})

// ClientInfo ...
type ClientInfo struct {
	Client   *kubernetes.Clientset
	Informer cache.SharedIndexInformer
}

// AddCluster 增加集群
func (d *Deployer) AddCluster(cluster, kubeConfig string) (clientInfo *ClientInfo, err error) {
	var conf *rest.Config
	if conf, err = clientcmd.BuildConfigFromFlags("", kubeConfig); err != nil {
		log.Println("BuildConfigFromFlags", err)
		return
	}
	client, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Println("NewForConfig", err)
		return
	}
	informer := corev1informer.NewPodInformer(client, "", 0, cache.Indexers{})
	go informer.Run(stopChan)
	if !cache.WaitForCacheSync(stopChan, informer.HasSynced) {
		log.Println(fmt.Errorf("sync cache err %s", cluster))
		return
	}
	clientInfo = &ClientInfo{
		Client:   client,
		Informer: informer,
	}
	d.clients.Store(cluster, clientInfo)
	return
}

// Client ...
func (d *Deployer) Client(cluster string) (client *kubernetes.Clientset, err error) {
	v, ok := d.clients.Load(cluster)
	if !ok {
		err = fmt.Errorf("无效的或未启用的集群 %s", cluster)
		return
	}
	client = v.(*ClientInfo).Client
	return
}

// Informer ...
func (d *Deployer) Informer(cluster string) (informer cache.SharedIndexInformer, err error) {
	v, ok := d.clients.Load(cluster)
	if !ok {
		err = fmt.Errorf("无效的或未启用的集群 %s", cluster)
		return
	}
	informer = v.(*ClientInfo).Informer
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

// Finalize ...
func (d *Deployer) Finalize() {
	close(stopChan)
}
