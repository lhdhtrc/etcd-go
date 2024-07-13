package pkg

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (core *CoreEntity) Cli(key string) (*clientv3.Client, error) {
	cli, ok := core.cli[key]
	if !ok {
		return nil, fmt.Errorf("etcd cli key not found")
	}
	return cli, nil
}

func (core *CoreEntity) Lease(key string) clientv3.LeaseID {
	lease, ok := core.lease[key]
	if !ok {
		return 0
	}
	return lease
}
