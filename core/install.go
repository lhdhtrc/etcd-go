package etcd

import (
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

func New(config *Config) (*clientv3.Client, error) {
	conf := clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   strings.Split(config.Address, ","),
	}

	if config.Account != "" && config.Password != "" {
		conf.Username = config.Account
		conf.Password = config.Password
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
		conf.TLS = tlsConfig
	}

	return clientv3.New(conf)
}
