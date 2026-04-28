package schema

type Config struct {
	RPCBindAddr   string `json:"rpc_bind_addr" toml:"rpc_bind_addr"`
	RPCPublicAddr string `json:"rpc_public_addr" toml:"rpc_public_addr"`
	RPCSecret     string `json:"rpc_secret" toml:"rpc_secret"`
	Admin         Admin  `json:"admin" toml:"admin"`
	S3API         S3API  `json:"s3_api" toml:"s3_api"`
	S3Web         S3Web  `json:"s3_web" toml:"s3_web"`
}

type Admin struct {
	AdminToken   string `json:"admin_token" toml:"admin_token"`
	APIBindAddr  string `json:"api_bind_addr" toml:"api_bind_addr"`
	MetricsToken string `json:"metrics_token" toml:"metrics_token"`
}

type S3API struct {
	APIBindAddr string `json:"api_bind_addr" toml:"api_bind_addr"`
	RootDomain  string `json:"root_domain" toml:"root_domain"`
	S3Region    string `json:"s3_region" toml:"s3_region"`
}

type S3Web struct {
	BindAddr   string `json:"bind_addr" toml:"bind_addr"`
	Index      string `json:"index" toml:"index"`
	RootDomain string `json:"root_domain" toml:"root_domain"`
}

type PublicConfig struct {
	RPCBindAddr   string      `json:"rpc_bind_addr"`
	RPCPublicAddr string      `json:"rpc_public_addr"`
	Admin         PublicAdmin `json:"admin"`
	S3API         S3API       `json:"s3_api"`
	S3Web         S3Web       `json:"s3_web"`
}

type PublicAdmin struct {
	APIBindAddr string `json:"api_bind_addr"`
}

func (c Config) Public() PublicConfig {
	return PublicConfig{
		RPCBindAddr:   c.RPCBindAddr,
		RPCPublicAddr: c.RPCPublicAddr,
		Admin: PublicAdmin{
			APIBindAddr: c.Admin.APIBindAddr,
		},
		S3API: c.S3API,
		S3Web: c.S3Web,
	}
}
