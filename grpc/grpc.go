package grpc

import (
	"context"
	"fmt"
	"github.com/Saburo90/statistical_report/conf"
	"github.com/Saburo90/statistical_report/grpc/server"
	"github.com/Saburo90/statistical_report/protos/statistical"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
)

func SetupGrpc() error {
	grpcListenner, err := net.Listen("tcp", conf.StatisC.GRpcService.Addr)

	if err != nil {
		return err
	}
	originServer := grpc.NewServer(grpc.UnaryInterceptor(unaryServerInterceptor))
	grpcStatistical := server.NewGRPCStatisticalServer()
	daily_statistical.RegisterStatisticalServer(originServer, grpcStatistical)
	go func() {
		zap.L().Info("GRPC LISTENING", zap.String("transport", "gRPC"), zap.String("ADDRESS", conf.StatisC.GRpcService.Addr))

		if err = originServer.Serve(grpcListenner); err != nil {
			zap.L().Error("Start GRPC FAILURE", zap.Error(err))
			os.Exit(-1)
		}
	}()

	return nil
}

// interceptor log printed
func unaryServerInterceptor(cotxt context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	zap.L().Info(info.FullMethod, zap.String("request", fmt.Sprintf("%+v", req)))

	resp, err := handler(cotxt, req)

	if err != nil {
		zap.L().Error(info.FullMethod, zap.Error(err))
		return nil, err
	}

	return resp, err
}
