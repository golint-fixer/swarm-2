package strategy

import (
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/scheduler/node"
)

// SpecificPlacementStrategy places the container into the cluster with specific nodes.
type SpecificPlacementStrategy struct {
}

// Initialize a SpecificPlacementStrategy.
func (p *SpecificPlacementStrategy) Initialize() error {
	return nil
}

// Name returns the name of the strategy.
func (p *SpecificPlacementStrategy) Name() string {
	return "specific"
}

// RankAndSort randomly sorts the list of nodes.
func (p *SpecificPlacementStrategy) RankAndSort(config *cluster.ContainerConfig, nodes []*node.Node) ([]*node.Node, error) {
	return nodes, nil
}
