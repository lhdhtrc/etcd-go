package pkg

import (
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (core *CoreEntity) Pub(lease clientv3.LeaseID, raw *RawEntity) {
	var val string
	if str, ok := raw.Value.(string); ok {
		val = str
	} else {
		t, _ := json.Marshal(raw.Value)
		val = string(t)
	}

	if lease != 0 {
		if _, err := core.cli.Put(core.ctx, raw.Key, val, clientv3.WithLease(lease)); err != nil {
			core.logger.Error(err.Error())
		}
	} else {
		if _, err := core.cli.Put(core.ctx, raw.Key, val); err != nil {
			core.logger.Error(err.Error())
		}
	}
}

// PubRaw batch send kv
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
				if _, err := cli.Put(core.ctx, raw.Key, val, clientv3.WithLease(lease)); err != nil {
					core.logger.Error(err.Error())
				}
			} else {
				if _, err := cli.Put(core.ctx, raw.Key, val); err != nil {
					core.logger.Error(err.Error())
				}
			}
		}
	} else {
		core.logger.Error(ce.Error())
	}
}
