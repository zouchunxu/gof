package node

// Node is slector node
type Node struct {
	Addr   string
	Weight *int64
	Ver    string
	Name   string
	Met    map[string]string
}

// Address is node address
func (n *Node) Address() string {
	return n.Addr
}

// ServiceName is node serviceName
func (n *Node) ServiceName() string {
	return n.Name
}

// InitialWeight is node initialWeight
func (n *Node) InitialWeight() *int64 {
	return n.Weight
}

// Version is node version
func (n *Node) Version() string {
	return n.Ver
}

// Metadata is node metadata
func (n *Node) Metadata() map[string]string {
	return n.Met
}
