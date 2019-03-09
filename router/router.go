package router

import (
	"github.com/Saburo90/statistical_report/handler/mini_statis"
	"github.com/Saburo90/statistical_report/handler/recycle_statis"
	"github.com/Saburo90/statistical_report/handler/sale_statis"
	"github.com/Saburo90/statistical_report/handler/user_statis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func SetupRouter(e *gin.Engine) {
	e.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Context-Type", "STATIS-TOKEN"},
	}))

	if gin.Mode() == gin.DebugMode {
		e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	e.NoRoute(func(c *gin.Context) { c.String(http.StatusNotFound, "") })

	statistical := e.Group("/statistical")

	// user statis module
	userStatisticsAPI := statistical.Group("/user")
	{
		userStatisticsAPI.POST("/getOverview", user_statis.GetOverviewHandler)

		userStatisticsAPI.POST("/getWxAPIData", user_statis.GetWxAPIDataHandler)
	}

	// sales statis module
	salesStatisticsAPI := statistical.Group("/sales")
	{
		salesStatisticsAPI.POST("/getSalesData", sale_statis.GetSalesDataHandler)
	}

	// miniprogram statis module
	miniStatisticsAPI := statistical.Group("/miniprogram")
	{
		miniStatisticsAPI.POST("/getMiniData", mini_statis.GetRoamMiniWxAPIDataHandler)
	}

	rcyStatisticsAPI := statistical.Group("/recycle")
	{
		rcyStatisticsAPI.POST("/getRecycleData", recycle_statis.GetRecyleOrderDataHandler)
	}
}
