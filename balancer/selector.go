package balancer

import (
	"context"

	"errors"
)

// ErrNoAvailable is no available node
var ErrNoAvailable = errors.New("no_available_node")

// Selector is node pick balancer
type Selector interface {
	// Select nodes
	// if err == nil, selected and done must not be empty.
	Select(ctx context.Context, nodes []Node) (selected Node, done Done, err error)
}
