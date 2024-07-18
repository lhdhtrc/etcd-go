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
	core.lease = grant.ID

	kac, ke := core.cli.KeepAlive(ctx, grant.ID)
	if ke != nil {
		core.retryLease()
		core.logger.Error(fmt.Sprintf("%s %s\n", logPrefix, ke.Error()))
		return
	}
	core.logger.Info(fmt.Sprintf("%s %s", logPrefix, "success ->"))

	for {
		select {
		case <-ctx.Done():
			return
		case r, ok := <-kac:
			fmt.Println(r, ok)
			if !ok {
				core.retryLease()
				return
			}
			if core.countRetry != 0 {
				core.countRetry = 0
			}
		}
	}
}

func (core *CoreEntity) Uninstall() {
	if _, err := core.cli.Revoke(context.Background(), core.lease); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("uninstall etcd success")
}

func (core *CoreEntity) Pub(ctx context.Context, raw *RawEntity) {
	var val string
	if str, ok := raw.Value.(string); ok {
		val = str
	} else {
		t, _ := json.Marshal(raw.Value)
		val = string(t)
	}

	if raw.Lease != 0 {
		if _, err := core.cli.Put(ctx, raw.Key, val, clientv3.WithLease(raw.Lease)); err != nil {
			core.logger.Error(err.Error())
		}
	} else {
		if _, err := core.cli.Put(ctx, raw.Key, val); err != nil {
			core.logger.Error(err.Error())
		}
	}
}

func (core *CoreEntity) Sub(prefix string, adapter func(e *clientv3.Event)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wc := core.cli.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wc {
		for _, e := range v.Events {
			adapter(e)
		}
	}
}
