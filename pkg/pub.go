package pkg

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (core *CoreEntity) Pub(ck, lk string, raw *RawEntity) {
	if cli, ce := core.Cli(ck); ce == nil {
		var val string
		if str, ok := raw.Value.(string); ok {
			val = str
		} else {
			t, _ := json.Marshal(raw.Value)
			val = string(t)
		}

		lease := core.Lease(lk)
		if lease != 0 {
			if _, err := cli.Put(context.Background(), raw.Key, val, clientv3.WithLease(lease)); err != nil {
				core.logger.Error(err.Error())
			}
		} else {
			if _, err := cli.Put(context.Background(), raw.Key, val); err != nil {
				core.logger.Error(err.Error())
			}
		}
	}
}

func (core *CoreEntity) PubRaw(info *PubEntity) {
	if cli, ce := core.Cli(info.CK); ce == nil {
		lease := core.Lease(info.LK)
		for _, raw := range info.Raw {
			var val string
			if str, ok := raw.Value.(string); ok {
				val = str
			} else {
				t, _ := json.Marshal(raw.Value)
				val = string(t)
			}
			if lease != 0 {
				if _, err := cli.Put(context.Background(), raw.Key, val, clientv3.WithLease(lease)); err != nil {
					core.logger.Error(err.Error())
				}
			} else {
				if _, err := cli.Put(context.Background(), raw.Key, val); err != nil {
					core.logger.Error(err.Error())
				}
			}
		}
	} else {
		core.logger.Error(ce.Error())
	}
}
