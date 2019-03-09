package conf

type (
	grpcCli struct {
		Addr string `toml:"addr"`
		Port int    `toml:"port"`
	}

	grpcService struct {
		Addr string `toml:"addr"`
	}
)
