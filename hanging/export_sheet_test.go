package hanging

import (
	"gitee.com/NotOnlyBooks/statistical_report/components"
	"gitee.com/NotOnlyBooks/statistical_report/conf"
	"testing"
)

func TestCompleteBnameAndBpriceAndDScore(t *testing.T) {
	if err := conf.LoadConfigFromToml("/home/gopath/src/statistical_report/config.toml"); err != nil {
		t.Error(err.Error())
		return
	}

	components.SetupLogger(true)

	if err := components.SetupDatabase(&conf.StatisC.DbC); err != nil {
		t.Error(err.Error())
		return
	}

	if err := components.SetupRedis(&conf.StatisC.RedisC); err != nil {
		t.Error(err.Error())
		return
	}
	CompleteBnameAndBpriceAndDScore()
}
