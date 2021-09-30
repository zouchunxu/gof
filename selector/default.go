package selector

import (
	"context"
	"fmt"
	"sync"
)

// Default is composite selector.
type Default struct {
	NodeBuilder WeightedNodeBuilder
	Balancer    Balancer
	Filters     []Filter

	lk            sync.RWMutex
	weightedNodes []Node
}

// Select select one node.
func (d *Default) Select(ctx context.Context, opts ...SelectOption) (selected Node, done DoneFunc, err error) {
	d.lk.RLock()
	weightedNodes := d.weightedNodes
	d.lk.RUnlock()
	fmt.Println("weightedNodes: ", weightedNodes[0].Address())
	for _, f := range d.Filters {
		weightedNodes = f(ctx, weightedNodes)
	}
	var options SelectOptions
	for _, o := range opts {
		o(&options)
	}
	for _, f := range options.Filters {
		weightedNodes = f(ctx, weightedNodes)
	}
	candidates := make([]WeightedNode, 0, len(weightedNodes))
	for _, n := range weightedNodes {
		candidates = append(candidates, n.(WeightedNode))
	}
	if len(candidates) == 0 {
		return nil, nil, ErrNoAvailable
	}
	return d.Balancer.Pick(ctx, candidates)
}

// Apply update nodes info.
func (d *Default) Apply(nodes []Node) {
	weightedNodes := make([]Node, 0, len(nodes))
	for _, n := range nodes {
		weightedNodes = append(weightedNodes, d.NodeBuilder.Build(n))
	}
	d.lk.Lock()
	//  不能删除未改变的node
	d.weightedNodes = weightedNodes
	d.lk.Unlock()
}
