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
func (core *CoreEntity) PubRaw(lk string, raw []*RawEntity) {
	lease := core.Lease(lk)
	for _, row := range raw {
		core.Pub(lease, row)
	}
}
