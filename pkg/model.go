package pkg

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type CoreEntity struct {
	ctx    context.Context
	cancel context.CancelFunc

	cli   *clientv3.Client
	lease clientv3.LeaseID

	ttl              int64
	maxRetry         uint
	countRetry       uint
	leaseRetryBefore func()
	leaseRetryAfter  func()

	logger *zap.Logger
}

type TLSEntity struct {
	CaCert        string `json:"ca_cert" bson:"ca_cert" yaml:"ca_cert" mapstructure:"ca_cert"`
	ClientCert    string `json:"client_cert" bson:"client_cert" yaml:"client_cert" mapstructure:"client_cert"`
	ClientCertKey string `json:"client_cert_key" bson:"client_cert_key" yaml:"client_cert_key" mapstructure:"client_cert_key"`
}

type ConfigEntity struct {
	Tls TLSEntity `json:"tls" bson:"tls" yaml:"tls" mapstructure:"tls"`

	Account  string `json:"account" bson:"account" yaml:"account" mapstructure:"account"`
	Password string `json:"password" bson:"password" yaml:"password" mapstructure:"password"`
	Address  string `json:"address" yaml:"address" mapstructure:"address"`

	TTL      int64 `json:"ttl" yaml:"ttl" mapstructure:"ttl"`
	MaxRetry uint  `json:"max_retry" yaml:"max_retry" mapstructure:"max_retry"`
	Mode     bool  `json:"mode" yaml:"mode" mapstructure:"mode"` // Mode is true cluster
}

type RawEntity struct {
	Key   string           `json:"key"`
	Value any              `json:"value"`
	Lease clientv3.LeaseID `json:"lease"`
}

type ServiceEntity struct {
	Name      string       `json:"name"`
	Namespace string       `json:"namespace"`
	Endpoint  string       `json:"endpoint"`
	Srv       *interface{} `json:"srv"`
}
