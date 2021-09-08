package random

import (
	"context"
	"github.com/zouchunxu/gof/balancer"
	"math/rand"
)

var (
	_ balancer.Selector = &Selector{}

	// Name is balancer name
	Name = "random"
)

type Selector struct{}

func New() *Selector {
	return &Selector{}
}

func (p *Selector) Select(_ context.Context, nodes []balancer.Node) (balancer.Node, balancer.Done, error) {
	if len(nodes) == 0 {
		err := balancer.ErrNoAvailable
		return nil, nil, err
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	done := func(ctx context.Context, info balancer.DoneInfo) {
		d(ctx, info)
	}
	return selected, done, nil
}
