package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
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
		core.logger.Error(fmt.Sprintf("%s %s", logPrefix, "etcd client not found"))
		return
	}

	if err := core.createLease(); err != nil {
		core.logger.Error(fmt.Sprintf("%s %s", logPrefix, err.Error()))
		core.retryLease()
		return
	}
	go core.sustainLease()

	core.logger.Info(fmt.Sprintf("%s %s", logPrefix, "success ->"))
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

func (core *CoreEntity) Sub(prefix string, init func(count int64, kvs []*mvccpb.KeyValue), adapter func(e *clientv3.Event)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core.Find(prefix, init)
	wc := core.cli.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wc {
		for _, e := range v.Events {
			adapter(e)
		}
	}
}

func (core *CoreEntity) Find(prefix string, handle func(count int64, kvs []*mvccpb.KeyValue)) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := core.cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		core.logger.Error(err.Error())
		return
	}
	handle(res.Count, res.Kvs)
}
