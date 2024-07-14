## Etcd Go
Provides easy to use API to operate Etcd.

### How to use it?
`go get github.com/lhdhtrc/etcd-go`
```go
package main

import (
	"fmt"
	etcd "github.com/lhdhtrc/etcd-go/pkg"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	instance, err := etcd.New(logger, &etcd.ConfigEntity{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// If you need a lease?
	instance.InitLease()

	// If you need to register microservices? (If you need service registration, please initialize the lease first)
	instance.Register(&etcd.ServiceEntity{})

	// If you need service discovery? (If you need another adapter, refer to the etcd.ServiceDiscoverAdapter)
	service := make(map[string]*[]string)
	instance.Sub("/microservice/lhdht", etcd.ServiceDiscoverAdapter(service))

	// If you need to add KV? 
	instance.Pub(&etcd.RawEntity{})
	
	// Note that at the end of the process, please reclaim your lease!
	instance.Uninstall()
}
```

### Finally
- If you feel good, click on star.
- If you have a good suggestion, please ask the issue.