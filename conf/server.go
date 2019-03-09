package conf

type ServerConfig struct {
	Network  string
	Addr     string
	CertFile string
	KeyFile  string
	IP       string `toml:"ip"`
}
