package etcd

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

// Register service register
func (core *CoreEntity) Register(service *ServiceEntity) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ref := reflect.TypeOf(service.Srv).Elem()
	length := ref.NumMethod()

	for i := 0; i < length; i++ {
		core.Pub(ctx, &RawEntity{
			Key:   fmt.Sprintf("%s/%s/%s/%d", service.Namespace, service.Name, ref.Method(i).Name, core.lease),
			Value: service.Endpoint,
			Lease: core.lease,
		})
	}
}
