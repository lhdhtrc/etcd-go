package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

func New(logger *zap.Logger, config *ConfigEntity) *CoreEntity {
	core := &CoreEntity{
		logger:   logger,
		ttl:      config.TTL,
		maxRetry: config.MaxRetry,
	}
	core.cli = core.install(config)

	return core
}

func (core *CoreEntity) InitLease() {
	logPrefix := "etcd init lease"
	core.logger.Info(fmt.Sprintf("%s %s", logPrefix, "start ->"))

	if core.cli == nil {
		core.logger.Error(fmt.Sprintf("%s %s\n", logPrefix, "etcd client not found"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grant, ge := core.cli.Grant(ctx, core.ttl)
	if ge != nil {
		core.retryLease()
		core.logger.Error(fmt.Sprintf("%s %s\n", logPrefix, ge.Error()))
		return
	}

	kac, ke := core.cli.KeepAlive(ctx, grant.ID)
	if ke != nil {
		core.retryLease()
		core.logger.Error(fmt.Sprintf("%s %s\n", logPrefix, ke.Error()))
		return
	}
	core.lease = grant.ID
	core.countRetry = 0

	go func() {
		for range kac {
		}
		core.retryLease()
	}()

	core.logger.Info(fmt.Sprintf("%s %s", logPrefix, "success ->"))
}

func (core *CoreEntity) Uninstall() {
	if _, err := core.cli.Revoke(context.Background(), core.lease); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("uninstall etcd success")
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

func (core *CoreEntity) Sub(prefix string, adapter func(e *clientv3.Event)) {
	wc := core.cli.Watch(core.ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	go func() {
		for v := range wc {
			for _, e := range v.Events {
				adapter(e)
			}
		}
	}()
}
