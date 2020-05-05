package deployer

import (
	"context"
	"sync"
)

// Deployer ...
type Deployer struct {
	ctx     context.Context
	clients sync.Map
}

// New ...
func New() *Deployer {
	return &Deployer{
		ctx: context.TODO(),
	}
}
