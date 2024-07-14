package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func New(logger *zap.Logger, config *ConfigEntity) (*CoreEntity, error) {
	ctx, cancel := context.WithCancel(context.Background())

	core := &CoreEntity{
		ctx:      ctx,
		cancel:   cancel,
		logger:   logger,
		ttl:      config.TTL,
		maxRetry: config.MaxRetry,
	}
	core.cli = core.install(config)

	return core, nil
}

func (core *CoreEntity) InitLease() {
	logPrefix := "etcd init lease"
	core.logger.Info(fmt.Sprintf("%s %s", logPrefix, "start ->"))

	if core.cli == nil {
		core.logger.Error(fmt.Sprintf("%s %s\n", logPrefix, "etcd client not found"))
		return
	}

	grant, ge := core.cli.Grant(core.ctx, core.ttl)
	if ge != nil {
		core.retryLease()
		core.logger.Error(fmt.Sprintf("%s %s\n", logPrefix, ge.Error()))
		return
	}

	kac, ke := core.cli.KeepAlive(core.ctx, grant.ID)
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
	if _, err := core.cli.Revoke(core.ctx, core.lease); err != nil {
		fmt.Println(err.Error())
		return
	}

	core.cancel()
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
