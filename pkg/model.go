package pkg

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type CoreEntity struct {
	ctx    context.Context
	cancel context.CancelFunc
	cli    map[string]*clientv3.Client
	lease  map[string]clientv3.LeaseID
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

	Mode         bool `json:"mode" yaml:"mode" mapstructure:"mode"` // Mode is true cluster
	LoggerEnable bool `json:"logger_enable" bson:"logger_enable" yaml:"logger_enable" mapstructure:"logger_enable"`
}

type RawEntity struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type PubEntity struct {
	CK  string       `json:"ck"` // CK Cli Key
	LK  string       `json:"lk"` // LK Lease Key
	Raw []*RawEntity `json:"raw"`
}

type ServiceEntity struct {
	CK        string       `json:"ck"` // CK Cli Key
	LK        string       `json:"lk"` // LK Cli
	Name      string       `json:"name"`
	Namespace string       `json:"namespace"`
	Endpoint  string       `json:"endpoint"`
	Srv       *interface{} `json:"srv"`
}
