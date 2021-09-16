package etcd

import (
	"google.golang.org/grpc/resolver"
)

type builder struct {
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	r := &gofResolver{}
	return r, nil
}

// Scheme return scheme of discovery
func (*builder) Scheme() string {
	return "grpc"
}

type gofResolver struct {
}

func (g *gofResolver) watch() {
	for {

	}
}

//func (r *discoveryResolver) watch() {
//	for {
//		select {
//		case <-r.ctx.Done():
//			return
//		default:
//		}
//		ins, err := r.w.Next()
//		if err != nil {
//			if errors.Is(err, context.Canceled) {
//				return
//			}
//			r.log.Errorf("[resolver] Failed to watch discovery endpoint: %v", err)
//			time.Sleep(time.Second)
//			continue
//		}
//		r.update(ins)
//	}
//}

func (g *gofResolver) ResolveNow(options resolver.ResolveNowOptions) {
	panic("implement me")
}

func (g *gofResolver) Close() {
	panic("implement me")
}
