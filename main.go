package main

import (
	_ "gitee.com/NotOnlyBooks/statistical_report/docs"
	"gitee.com/NotOnlyBooks/statistical_report/server"
)

// @title 统计报表API
// @version 1.0
// @description
// @host 39.106.39.7:8092
// @BasePath /statistical
func main() {
	server.Start()
}
