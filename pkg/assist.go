package pkg

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (core *CoreEntity) Cli() *clientv3.Client {
	return core.cli
}

func (core *CoreEntity) Lease(key string) clientv3.LeaseID {
	lease, ok := core.lease[key]
	if !ok {
		return 0
	}
	return lease
}
