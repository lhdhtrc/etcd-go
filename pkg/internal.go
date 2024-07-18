package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

func (core *CoreEntity) install(config *ConfigEntity) *clientv3.Client {
	logPrefix := "install etcd"
	core.logger.Info(fmt.Sprintf("%s %s", logPrefix, "start ->"))

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
			core.logger.Error(fmt.Sprintf("error: %s", err.Error()))
			return nil
		}
		clientOptions.TLS = tlsConfig
	}

	cli, err := clientv3.New(clientOptions)
	if err != nil {
		core.logger.Error(fmt.Sprintf("error: %s", err.Error()))
		return nil
	}

	core.logger.Info(fmt.Sprintf("%s %s", logPrefix, "success ->"))

	return cli
}

func (core *CoreEntity) createLease() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grant, err := core.cli.Grant(ctx, core.ttl)
	if err != nil {
		return err
	}
	core.lease = grant.ID

	return nil
}

func (core *CoreEntity) sustainLease() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kac, ke := core.cli.KeepAlive(ctx, core.lease)
	if ke != nil {
		core.retryLease()
		core.logger.Error(ke.Error())
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-kac:
			if !ok {
				core.retryLease()
				return
			}
			if core.countRetry != 0 {
				core.countRetry = 0
			}
		}
	}
}

func (core *CoreEntity) retryLease() {
	if core.countRetry < core.maxRetry {
		if core.leaseRetryBefore != nil {
			core.leaseRetryBefore()
		}
		time.Sleep(5 * time.Second)

		core.countRetry++
		core.logger.Info(fmt.Sprintf("etcd retry lease: %d/%d", core.countRetry, core.maxRetry))

		if err := core.createLease(); err != nil {
			core.logger.Error(err.Error())
			core.retryLease()
			return
		}
		go core.sustainLease()

		if core.leaseRetryAfter != nil {
			core.leaseRetryAfter()
		}
	}
}
