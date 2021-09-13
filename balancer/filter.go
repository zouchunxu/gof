package balancer

import "context"

var _ Selector = &FilterSelector{}

// Filter is nodes filter
type Filter interface {
	Filter(ctx context.Context, nodes []Node) (filtered []Node, err error)
}

// FilterSelector is a selector contains filters
type FilterSelector struct {
	selector Selector
	filters  []Filter
}

// BuildFilterSelector create FilterSelector
func BuildFilterSelector(selector Selector, filters ...Filter) *FilterSelector {
	return &FilterSelector{
		selector: selector,
		filters:  filters,
	}
}

// Select one node via node filters
func (fs *FilterSelector) Select(ctx context.Context, nodes []Node) (Node, Done, error) {
	for _, filter := range fs.filters {
		var err error
		nodes, err = filter.Filter(ctx, nodes)
		if err != nil {
			return nil, nil, err
		}
	}
	if len(nodes) == 0 {
		return nil, nil, ErrNoAvailable
	}
	return fs.selector.Select(ctx, nodes)
}
