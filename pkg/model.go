package etcd

type TLS struct {
	CaCert        string `json:"ca_cert" bson:"ca_cert" yaml:"ca_cert" mapstructure:"ca_cert"`
	ClientCert    string `json:"client_cert" bson:"client_cert" yaml:"client_cert" mapstructure:"client_cert"`
	ClientCertKey string `json:"client_cert_key" bson:"client_cert_key" yaml:"client_cert_key" mapstructure:"client_cert_key"`
}

type Config struct {
	Username string   `json:"account" bson:"account" yaml:"account" mapstructure:"account"`
	Password string   `json:"password" bson:"password" yaml:"password" mapstructure:"password"`
	Tls      TLS      `json:"tls" bson:"tls" yaml:"tls" mapstructure:"tls"`
	Endpoint []string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
}
