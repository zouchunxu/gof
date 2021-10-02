package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/apimachinery/pkg/util/json"
	"math/rand"
	"time"
)

type Registry struct {
	kv     clientv3.KV
	client *clientv3.Client
	lease  clientv3.Lease
}

var namespace = "gof"

//New creates etcd registry
func New(client *clientv3.Client) *Registry {
	return &Registry{
		kv:     clientv3.NewKV(client),
		client: client,
	}
}

type val struct {
	Ips []string               `json:"ips"`
	Md  map[string]interface{} `json:"md"`
}

func (r *Registry) Registry(ctx context.Context, key string, values []string, md map[string]interface{}) error {
	key = namespace + "/" + key + "/" + "service"
	v := &val{
		Ips: values,
		Md:  md,
	}
	tmp, err := json.Marshal(v)
	if err != nil {
		return err
	}
	value := string(tmp)
	if r.lease != nil {
		_ = r.lease.Close()
	}
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		return err
	}
	go r.heartBeat(ctx, leaseID, key, value)
	return nil
}

func (r *Registry) heartBeat(ctx context.Context, leaseID clientv3.LeaseID, key string, value string) {
	curLeaseID := leaseID
	kac, err := r.client.KeepAlive(ctx, leaseID)
	if err != nil {
		curLeaseID = 0
	}
	rand.Seed(time.Now().Unix())
	for {
		if curLeaseID == 0 {
			var retreat []int
			for retryCnt := 0; retryCnt < 5; retryCnt++ {
				if ctx.Err() != nil {
					return
				}
				idChan := make(chan clientv3.LeaseID)
				errChan := make(chan error)
				cancelCtx, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()
					id, err := r.registerWithKV(cancelCtx, key, value)
					if err != nil {
						errChan <- err
					} else {
						idChan <- id
					}
				}()
				select {
				case <-time.After(3 * time.Second):
					cancel()
					continue
				case <-errChan:
					continue
				case curLeaseID = <-idChan:
				}

				kac, err = r.client.KeepAlive(ctx, curLeaseID)
				if err == nil {
					break
				}
				retreat = append(retreat, 1<<retryCnt) // 2的retryCnt次方
				time.Sleep(time.Duration(retreat[rand.Intn(len(retreat))]) * time.Second)
			}
			//err
			if _, ok := <-kac; !ok {
				return
			}
		}
		select {
		case _, ok := <-kac:
			if !ok {
				fmt.Println("not ok")
				if ctx.Err() != nil {
					return
				}
				curLeaseID = 0
				continue
			}
		case <-ctx.Done():
		}

	}
}

//registerWithKV create a new lease, return current leaseID
func (r *Registry) registerWithKV(ctx context.Context, key, value string) (clientv3.LeaseID, error) {
	grant, err := r.lease.Grant(ctx, int64(time.Second.Seconds()))
	if err != nil {
		return 0, err
	}
	_, err = r.client.Put(ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return 0, err
	}
	return grant.ID, nil
}

func (r *Registry) GetService(ctx context.Context, name string) ([]val, error) {
	key := namespace + "/" + name
	resp, err := r.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var v []val
	for _, kv := range resp.Kvs {
		var list val
		_ = json.Unmarshal(kv.Value, &list)
		v = append(v, list)
	}
	return v, nil
}

func (r *Registry) Watch(ctx context.Context, name string) (*watcher, error) {
	key := fmt.Sprintf("%s/%s", namespace, name)
	return newWatcher(ctx, key, r.client), nil
}
