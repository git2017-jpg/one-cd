package deployer

import (
	"fmt"
	"log"
	"path"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var stopper = make(chan struct{})

// ClientInfo ...
type ClientInfo struct {
	Client      *kubernetes.Clientset
	PodInformer cache.SharedIndexInformer
}

func (d *Deployer) loadCluster(cluster string) (clientInfo *ClientInfo, err error) {
	clientInfo, err = d.AddCluster(cluster, path.Join(d.configPath, cluster))
	return
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
	factory := informers.NewSharedInformerFactory(client, 0)
	podInformer := factory.Core().V1().Pods().Informer()
	go factory.Start(stopper)
	if !cache.WaitForCacheSync(stopper, podInformer.HasSynced) {
		log.Println(fmt.Errorf("Timed out waiting for caches to sync, cluster: %s", cluster))
		return
	}
	clientInfo = &ClientInfo{
		Client:      client,
		PodInformer: podInformer,
	}
	d.clients.Store(cluster, clientInfo)
	return
}

// Client ...
func (d *Deployer) Client(cluster string) (client *kubernetes.Clientset, err error) {
	v, ok := d.clients.Load(cluster)
	if !ok {
		if v, err = d.loadCluster(cluster); err != nil {
			err = fmt.Errorf("未找到集群 %s 配置文件", cluster)
			return
		}
	}
	client = v.(*ClientInfo).Client
	return
}

// Informer ...
func (d *Deployer) PodInformer(cluster string) (informer cache.SharedIndexInformer, err error) {
	v, ok := d.clients.Load(cluster)
	if !ok {
		if v, err = d.loadCluster(cluster); err != nil {
			err = fmt.Errorf("未找到集群 %s 配置文件", cluster)
			return
		}
	}
	informer = v.(*ClientInfo).PodInformer
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
	close(stopper)
}
