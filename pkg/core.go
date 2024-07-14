package pkg

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strings"
	"time"
)

func New(logger *zap.Logger, config *ConfigEntity) *CoreEntity {
	ctx, cancel := context.WithCancel(context.Background())

	core := &CoreEntity{
		ctx:    ctx,
		cancel: cancel,
		lease:  make(map[string]clientv3.LeaseID),
		logger: logger,
	}

	if cli, err := core.Setup(config); err == nil {
		core.cli = cli
	} else {
		logger.Error(err.Error())
	}

	core.initLease()

	return core, nil
}

func (core *CoreEntity) InitLease() {
	logPrefix := "etcd init lease"
	fmt.Printf("%s %s\n", logPrefix, "start ->")

	if core.cli == nil {
		fmt.Printf("%s %s\n", logPrefix, "etcd client not found")
		return
	}

	grant, ge := core.cli.Grant(core.ctx, core.ttl)
	if ge != nil {
		core.retryLease()
		fmt.Printf("%s %s\n", logPrefix, ge.Error())
		return
	}

	kac, ke := core.cli.KeepAlive(core.ctx, grant.ID)
	if ke != nil {
		core.retryLease()
		fmt.Printf("%s %s\n", logPrefix, ke.Error())
		return
	}
	core.lease = grant.ID
	core.countRetry = 0

	go func() {
		for range kac {
		}
		core.retryLease()
		fmt.Println("stop etcd lease success")
	}()
	fmt.Printf("%s %s\n", logPrefix, "success ->")
}

func (core *CoreEntity) Pub(raw *RawEntity) {
	var val string
	if str, ok := raw.Value.(string); ok {
		val = str
	} else {
		t, _ := json.Marshal(raw.Value)
		val = string(t)
	}

	if raw.Lease != 0 {
		if _, err := core.cli.Put(core.ctx, raw.Key, val, clientv3.WithLease(raw.Lease)); err != nil {
			core.logger.Error(err.Error())
		}
	} else {
		if _, err := core.cli.Put(core.ctx, raw.Key, val); err != nil {
			core.logger.Error(err.Error())
		}
	}
}
