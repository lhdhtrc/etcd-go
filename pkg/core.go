package pkg

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

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

func (core *CoreEntity) Cli(key string) (*clientv3.Client, error) {
	cli, ok := core.cli[key]
	if !ok {
		return nil, fmt.Errorf("etcd cli key not found")
	}
	return cli, nil
}

func (core *CoreEntity) Lease(key string) clientv3.LeaseID {
	lease, ok := core.lease[key]
	if !ok {
		return 0
	}
	return lease
}

func (core *CoreEntity) Pub(info *PubEntity) {
	if cli, ce := core.Cli(info.CK); ce == nil {
		lease := core.Lease(info.LK)
		for _, raw := range info.Raw {
			val, _ := json.Marshal(raw.Value)
			if lease != 0 {
				if _, err := cli.Put(context.Background(), raw.Key, string(val), clientv3.WithLease(lease)); err != nil {
					core.logger.Error(err.Error())
				}
			} else {
				if _, err := cli.Put(context.Background(), raw.Key, string(val)); err != nil {
					core.logger.Error(err.Error())
				}
			}
		}
	} else {
		core.logger.Error(ce.Error())
	}
}
