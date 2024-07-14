package etcd

import (
	"fmt"
	"reflect"
)

// Register service register
func (core *CoreEntity) Register(service *ServiceEntity) {
	ref := reflect.TypeOf(service.Srv).Elem()
	length := ref.NumMethod()

	for i := 0; i < length; i++ {
		core.Pub(&RawEntity{
			Key:   fmt.Sprintf("%s/%s/%s/%d", service.Namespace, service.Name, ref.Method(i).Name, core.lease),
			Value: service.Endpoint,
			Lease: core.lease,
		})
	}
}
