package main

import (
	"fmt"
	"io/ioutil"
	"one-cd/deployer"
	"time"
)

func main() {
	d := deployer.New()
	d.AddCluster("default", "/Users/rongchang/.kube/config")

	yml, err := ioutil.ReadFile("/Users/rongchang/codes/projects/k8s-demo/k8s.yaml")
	if err != nil {
		fmt.Print(err)
		return
	}
	data, err := d.Deploy(string(yml))
	if err != nil {
		fmt.Println("错误信息：", err)
		return
	}
	fmt.Println(data)

	d.WaitForPodContainersRunning("default", "default", "k8s-demo", time.Second*10, time.Second*2)
}
