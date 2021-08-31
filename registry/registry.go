package registry

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

func Register(client *clientv3.Client, key, addr string) error {
	em, err := endpoints.NewManager(client, key)
	if err != nil {
		return err
	}
	return em.AddEndpoint(context.TODO(), key+"/"+addr, endpoints.Endpoint{Addr: addr})
}
