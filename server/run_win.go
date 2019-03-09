package server

import (
	"github.com/Saburo90/statistical_report/conf"
	"github.com/Saburo90/statistical_report/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func run() error {
	if !conf.StatisC.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	var (
		handler = gin.New()
		addr    = conf.StatisC.ServerC.Addr
		isTLS   = conf.StatisC.ServerC.CertFile != "" && conf.StatisC.ServerC.KeyFile != ""
	)

	handler.Use(router.Logger(), gin.Recovery())

	router.SetupRouter(handler)

	zap.L().Info("HTTP LISTENING", zap.String("Address", addr), zap.Bool("IsTLS", isTLS))
	if isTLS {
		return http.ListenAndServeTLS(addr, conf.StatisC.ServerC.CertFile, conf.StatisC.ServerC.KeyFile, handler)
	} else {
		return http.ListenAndServe(addr, handler)
	}
}
