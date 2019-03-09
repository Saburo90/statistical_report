package grpccli

import (
	"fmt"
	"github.com/Saburo90/statistical_report/conf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"strconv"
)

var (
	gRpcConns = map[string]*grpc.ClientConn{}
)

func SetupGrpcClis() error {
	for name, v := range conf.StatisC.GRpcCli {
		conn, err := grpc.Dial(v.Addr+":"+strconv.Itoa(v.Port), grpc.WithInsecure())

		if err != nil {
			return err
		}

		zap.L().Info("gRPC Connect", zap.String("name", name), zap.String("Addr", v.Addr), zap.Int("Port", v.Port))

		gRpcConns[name] = conn
	}

	return nil
}

func GetgRPCCli(name string) (*grpc.ClientConn, error) {
	cli, ok := gRpcConns[name]

	if !ok {
		return nil, fmt.Errorf("gRPC Conns Not Found [%s]", name)
	}

	return cli, nil
}
