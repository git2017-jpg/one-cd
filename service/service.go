package service

import (
	"one-cd/conf"
	"one-cd/deployer"
)

// Service ...
type Service struct {
	Deployer *deployer.Deployer
}

// New ...
func New() (s *Service) {
	s = &Service{
		Deployer: deployer.New(conf.Conf.KubeConfigPath),
	}
	return
}

// Init ...
func (s *Service) Init() (err error) {
	return
}

// Ping ...
func (s *Service) Ping() (err error) {
	return
}
