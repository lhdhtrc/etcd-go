package pkg

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strings"
	"time"
)

func New(logger *zap.Logger, config *ConfigEntity) *CoreEntity {
	ctx, cancel := context.WithCancel(context.Background())

	core := &CoreEntity{
		ctx:    ctx,
		cancel: cancel,
		lease:  make(map[string]clientv3.LeaseID),
		logger: logger,
	}

	if cli, err := core.Setup(config); err == nil {
		core.cli = cli
	} else {
		logger.Error(err.Error())
	}

	core.initLease()

	return core, nil
}

func (core *CoreEntity) Setup(config *ConfigEntity) (*clientv3.Client, error) {
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
