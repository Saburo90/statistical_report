package server

import (
	"crypto/tls"
	"gitee.com/NotOnlyBooks/statistical_report/conf"
	"gitee.com/NotOnlyBooks/statistical_report/router"
	"github.com/facebookgo/grace/gracehttp"
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
		isTLS   = conf.StatisC.ServerC.KeyFile != "" && conf.StatisC.ServerC.CertFile != ""
	)

	handler.Use(router.Logger(), gin.Recovery())

	router.SetupRouter(handler)

	zap.L().Info("HTTP SERVER LISTENING(linux)", zap.String("Address", addr), zap.Bool("isTLS", isTLS))

	if false && conf.StatisC.Debug {
		if isTLS {
			return http.ListenAndServeTLS(addr, conf.StatisC.ServerC.CertFile, conf.StatisC.ServerC.KeyFile, handler)
		} else {
			return http.ListenAndServe(addr, handler)
		}
	} else {
		server := &http.Server{Addr: addr, Handler: handler}

		if isTLS {
			certificates, err := tls.LoadX509KeyPair(conf.StatisC.ServerC.CertFile, conf.StatisC.ServerC.KeyFile)

			if err != nil {
				return err
			}

			config := tls.Config{}
			config.Certificates = make([]tls.Certificate, 1)
			config.Certificates[0] = certificates
			server.TLSConfig = &config
		}
		return gracehttp.Serve(server)
	}
}
