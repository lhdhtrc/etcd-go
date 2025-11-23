package etcd

import (
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func New(conf *Conf) (*clientv3.Client, error) {
	config := clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   conf.Endpoint,
	}

	if conf.Username != "" && conf.Password != "" {
		config.Username = conf.Username
		config.Password = conf.Password
	}
	if conf.Tls.CaCert != "" && conf.Tls.ClientCert != "" && conf.Tls.ClientCertKey != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      conf.Tls.ClientCert,
			KeyFile:       conf.Tls.ClientCertKey,
			TrustedCAFile: conf.Tls.CaCert,
		}

		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
		config.TLS = tlsConfig
	}

	return clientv3.New(config)
}
