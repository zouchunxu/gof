package grpc

import (
	"fmt"
	"github.com/zouchunxu/gof/selector"
	"github.com/zouchunxu/gof/selector/node"
	gBalancer "google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
	"math/rand"
)

func init() {
	b := base.NewBalancerBuilder(
		"gof",
		&GofBuilder{},
		base.Config{HealthCheck: true},
	)
	gBalancer.Register(b)
}

type GofPicker struct {
	subConns map[string]gBalancer.SubConn
	nodes    []selector.Node
}

func (g *GofPicker) Pick(info gBalancer.PickInfo) (gBalancer.PickResult, error) {
	cur := rand.Intn(len(g.nodes))
	selected := g.nodes[cur]
	sub := g.subConns[selected.Address()]
	fmt.Println("sub: ", selected.Address())

	return gBalancer.PickResult{
		SubConn: sub,
		Done: func(di gBalancer.DoneInfo) {
			i := selector.DoneInfo{
				Err:           di.Err,
				BytesSent:     di.BytesSent,
				BytesReceived: di.BytesReceived,
				ReplyMeta:     GofTrailer(di.Trailer),
			}
			fmt.Printf("done i %+v", i)
		},
	}, nil
}

type GofBuilder struct {
	selector selector.Selector
}

func (g *GofBuilder) Build(info base.PickerBuildInfo) gBalancer.Picker {
	nodes := make([]selector.Node, 0)
	subConns := make(map[string]gBalancer.SubConn)
	for conn, info := range info.ReadySCs {
		if _, ok := subConns[info.Address.Addr]; ok {
			continue
		}
		subConns[info.Address.Addr] = conn

		var w int64 = 100
		n := &node.Node{
			Addr:   info.Address.Addr,
			Name:   info.Address.ServerName,
			Weight: &w,
		}
		nodes = append(nodes, n)
	}

	p := &GofPicker{
		//selector: g.selector,
		subConns: subConns,
	}
	p.nodes = nodes
	//p.selector.Apply(nodes)
	return p
}

type GofTrailer metadata.MD

func (g GofTrailer) Get(key string) string {
	k := metadata.MD(g).Get(key)
	if len(k) == 0 {
		return ""
	}
	return k[len(k)-1]
}
