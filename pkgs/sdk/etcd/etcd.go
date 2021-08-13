package etcd

import (
	"context"
	naming2 "github.com/shipengqi/qim/library/naming"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
)

type etcd struct {
	sync.RWMutex

	cli     *clientv3.Client
	lease   clientv3.Lease
	leaseID clientv3.LeaseID
}

func New(url string) (naming2.Interface, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            []string{url},
		DialTimeout:          15 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	cli.KV = namespace.NewKV(cli.KV, "my-prefix/")
	cli.Watcher = namespace.NewWatcher(cli.Watcher, "my-prefix/")
	cli.Lease = namespace.NewLease(cli.Lease, "my-prefix/")

	return &etcd{
		cli: cli,
	}, nil
}

func (e *etcd) Find(name string, tags ...string) ([]naming2.Registration, error) {
	_, err := e.cli.Get(context.TODO(), name, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (e *etcd) Register(s naming2.Registration) error {
	// _, err := e.cli.Put(context.TODO(), key, val)
	return nil
}

func (e *etcd) Deregister(serviceID string) error {
	// _, err := e.lease.Revoke(context.TODO(), e.leaseResp.ID)
	return nil
}

func (e *etcd) Subscribe(serviceName string, callback func([]naming2.Registration)) error {
	return nil
}

func (e *etcd) Unsubscribe(serviceName string) error {
	return nil
}
