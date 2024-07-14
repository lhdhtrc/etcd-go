package pkg

import (
	"fmt"
	"reflect"
)

// Register service register
func (core *CoreEntity) Register(service *ServiceEntity) {
	ref := reflect.TypeOf(service.Srv).Elem()
	length := ref.NumMethod()

	lease := core.Lease(service.LK)
	for i := 0; i < length; i++ {
		core.Pub(lease, &RawEntity{
			Key:   fmt.Sprintf("%s/%s/%s/%d", service.Namespace, service.Name, ref.Method(i).Name, lease),
			Value: service.Endpoint,
		})
	}
}
