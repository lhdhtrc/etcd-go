package pkg

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (core *CoreEntity) Cli() *clientv3.Client {
	return core.cli
}

func (core *CoreEntity) Lease() clientv3.LeaseID {
	return core.lease
}
