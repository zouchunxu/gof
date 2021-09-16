package etcd

import "google.golang.org/grpc/resolver"

type builder struct {
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	return &gofResolver{}, nil
}

type gofResolver struct {
}

func (g gofResolver) ResolveNow(options resolver.ResolveNowOptions) {
	panic("implement me")
}

func (g gofResolver) Close() {
	panic("implement me")
}
