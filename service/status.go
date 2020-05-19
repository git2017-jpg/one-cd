package service

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"
)

func (s *Service) statusPrinter(cluster, namespace, deploymentName, info string) {
	fmt.Println(info)
}

// WaitForDeployment ...
func (s *Service) WaitForDeployment(cluster, namespace, deploymentName string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, string(debug.Stack()))
		}
	}()
	s.Deployer.WaitForPodContainersRunning(cluster, namespace, deploymentName,
		time.Second*3600, time.Second*3, s.statusPrinter)
}
