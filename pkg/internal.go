package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

func (core *CoreEntity) install(config *ConfigEntity) (*clientv3.Client, error) {
	logPrefix := "setup etcd"
	fmt.Printf("%s %s\n", logPrefix, "start ->")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientOptions := clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   strings.Split(config.Address, ","),
		Context:     ctx,
	}

	if config.Account != "" && config.Password != "" {
		clientOptions.Username = config.Account
		clientOptions.Password = config.Password
	}
	if config.Tls.CaCert != "" && config.Tls.ClientCert != "" && config.Tls.ClientCertKey != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      config.Tls.ClientCert,
			KeyFile:       config.Tls.ClientCertKey,
			TrustedCAFile: config.Tls.CaCert,
		}

		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
		clientOptions.TLS = tlsConfig
	}

	cli, err := clientv3.New(clientOptions)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s %s", logPrefix, "success ->")

	return cli, nil
}

func (core *CoreEntity) initLease() {
	logPrefix := "init lease"
	fmt.Printf("%s %s\n", logPrefix, "start ->")

	if core.cli == nil {
		fmt.Printf("%s %s\n", logPrefix, "etcd client not found")
		return
	}

	grant, ge := core.cli.Grant(core.ctx, core.ttl)
	if ge != nil {
		core.retryLease()
		fmt.Printf("%s %s\n", logPrefix, ge.Error())
		return
	}

	kac, ke := core.cli.KeepAlive(core.ctx, grant.ID)
	if ke != nil {
		core.retryLease()
		fmt.Printf("%s %s\n", logPrefix, ke.Error())
		return
	}
	core.lease = grant.ID

	go func() {
		for range kac {
		}
		core.retryLease()
		fmt.Println("lease stop success")
	}()
	fmt.Printf("%s %s\n", logPrefix, "success ->")
}

func (core *CoreEntity) retryLease() {
	if core.countRetry < core.maxRetry {
		if core.leaseRetryBefore != nil {
			core.leaseRetryBefore()
		}
		time.Sleep(5 * time.Second)

		core.countRetry++
		fmt.Printf("retry lease: %d/%d\n", core.countRetry, core.maxRetry)
		core.initLease()

		if core.leaseRetryAfter != nil {
			core.leaseRetryAfter()
		}
	}
}
