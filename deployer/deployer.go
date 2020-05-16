package deployer

import (
	"context"
	"sync"
)

// Deployer ...
type Deployer struct {
	ctx        context.Context
	clients    sync.Map
	configPath string
}

// New ...
func New(configPath string) *Deployer {
	return &Deployer{
		ctx:        context.TODO(),
		configPath: configPath,
	}
}
