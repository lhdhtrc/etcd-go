package etcd

import (
	"github.com/lhdhtrc/func-go/array"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
)

// ServiceDiscoverAdapter 服务发现适配器
func ServiceDiscoverAdapter(service map[string]*[]string) func(e *clientv3.Event) {
	return func(e *clientv3.Event) {
		var (
			key string
			val string
		)

		if e.PrevKv != nil {
			key = string(e.PrevKv.Key)
			val = string(e.PrevKv.Value)
		} else {
			key = string(e.Kv.Key)
			val = string(e.Kv.Value)
		}

		kt := strings.Split(key, "/")
		kt = kt[:len(kt)-1]
		key = strings.Join(kt, "/")

		switch e.Type {
		// PUT，新增或替换
		case clientv3.EventTypePut:
			*service[key] = append(*service[key], val)
			*service[key] = array.Unique[string](*service[key], func(index int, item string) string {
				return item
			})
		// DELETE
		case clientv3.EventTypeDelete:
			*service[key] = array.Filter(*service[key], func(index int, item string) bool {
				return item != val
			})
		}
	}
}
