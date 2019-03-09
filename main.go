package main

import (
	_ "github.com/Saburo90/statistical_report/docs"
	"github.com/Saburo90/statistical_report/server"
)

// @title 统计报表API
// @version 1.0
// @description
// @host 39.106.39.7:8092
// @BasePath /statistical
func main() {
	server.Start()
}
