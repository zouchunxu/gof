package random

import (
	"context"
	"github.com/zouchunxu/gof/selector"
	"github.com/zouchunxu/gof/selector/node/direct"
	"math/rand"
)

const (
	// Name is random balancer name
	Name = "random"
)

var _ selector.Balancer = &Balancer{} // Name is balancer name

// WithFilter 过滤filter
func WithFilter(filters ...selector.Filter) Option {
	return func(o *options) {
		o.filters = filters
	}
}

// Option is random builder option.
type Option func(o *options)

// options is random builder options
type options struct {
	filters []selector.Filter
}

// Balancer is a random balancer.
type Balancer struct{}

// New 随机select
func New(opts ...Option) selector.Selector {
	var option options
	for _, opt := range opts {
		opt(&option)
	}

	return &selector.Default{
		Balancer:    &Balancer{},
		NodeBuilder: &direct.Builder{},
		Filters:     option.filters,
	}
}

// Pick pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	return selected, d, nil
}
