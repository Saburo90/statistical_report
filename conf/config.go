package conf

import (
	"errors"
	"github.com/BurntSushi/toml"
)

var (
	StatisC GlobalC
)

type GlobalC struct {
	Debug       bool               `toml:"debug"`
	DbC         DBConfig           `toml:"db"`
	RedisC      RedisConfig        `toml:"redis"`
	ServerC     ServerConfig       `toml:"server"`
	WXApi       WxAPI              `toml:"wxapi"`
	WXSecret    WxSecret           `toml:"wxsecret"`
	PHPApi      PHPapi             `toml:"phpapi"`
	GRpcCli     map[string]grpcCli `toml:"grpc_cli"`
	GRpcService grpcService        `toml:"grpc_service"`
}

func LoadConfigFromToml(fileName string) error {
	if _, err := toml.DecodeFile(fileName, &StatisC); err != nil {
		return err
	}
	return nil
}

func (gC *GlobalC) VerifyConfig() error {
	if gC.ServerC.Addr == "" {
		return errors.New("Server Addr Must be Configured")
	}

	if gC.ServerC.IP == "" {
		return errors.New("Server IP Must be Configured")
	}

	return nil
}
