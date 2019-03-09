package server

import (
	"github.com/Saburo90/statistical_report/components"
	"github.com/Saburo90/statistical_report/components/grpccli"
	"github.com/Saburo90/statistical_report/conf"
	"github.com/Saburo90/statistical_report/grpc"
	"github.com/Saburo90/statistical_report/timed"
	"go.uber.org/zap"
	"log"
	"os"
)

func Start() {
	if len(os.Args) == 1 {
		log.Fatal("Empty Config File")
	}

	if err := conf.LoadConfigFromToml(os.Args[1]); err != nil {
		log.Fatal("Load Config Failure, Error:%v", err)
	}

	if err := components.SetupLogger(conf.StatisC.Debug); err != nil {
		log.Fatal("Setup Zap Log Failure, Error:%v", err)
	}

	zap.L().Info("Start Connect Database")
	if err := components.SetupDatabase(&conf.StatisC.DbC); err != nil {
		zap.L().Fatal("Connect To Database Failure", zap.Error(err))
	}

	zap.L().Info("Start Connect Redis")
	if err := components.SetupRedis(&conf.StatisC.RedisC); err != nil {
		zap.L().Fatal("Connect To Redis Failure", zap.Error(err))
	}

	zap.L().Info("Start Grpc Service")
	if err := grpc.SetupGrpc(); err != nil {
		zap.L().Fatal("Start Grpc Service Failure", zap.Error(err))
	}

	zap.L().Info("Start Grpc Client")
	if err := grpccli.SetupGrpcClis(); err != nil {
		zap.L().Fatal("Start Grpc Client Failure", zap.Error(err))
	}

	zap.L().Info("Start Timed Task")
	if err := timed.SetupTimed(); err != nil {
		zap.L().Fatal("Start Timed Task Failure")
	}

	zap.L().Info("terminated", zap.Error(run()))
}
