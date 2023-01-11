package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/gredis"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/routers"
)

// 初始化
func init() {

	// 初始化配置
	setting.Setup()

	// models 初始化
	models.Setup()

	// 初始化日志
	logging.Setup()

	// 初始化 redis
	gredis.Setup()

	// 初始化工具
	util.Setup()
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/EDDYCJY/go-gin-example
// @license.name MIT
// @license.url https://github.com/EDDYCJY/go-gin-example/blob/master/LICENSE
func main() {

	// 设置 gin 模式
	// 生产环境下使用 gin.ReleaseMode
	// 开发环境下使用 gin.DebugMode
	// 默认 gin.DebugMode
	gin.SetMode(setting.ServerSetting.RunMode)

	// 初始化路由
	routersInit := routers.InitRouter()

	// 设置 读取超时时间
	readTimeout := setting.ServerSetting.ReadTimeout

	// 设置 写入超时时间
	writeTimeout := setting.ServerSetting.WriteTimeout

	// 设置 端口
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	// 设置头信息 大小 1 << 20 = 1048576 = 1M
	maxHeaderBytes := 1 << 20

	// 设置服务器参数
	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	// 启动服务
	server.ListenAndServe()

	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	//endless.DefaultReadTimeOut = readTimeout
	//endless.DefaultWriteTimeOut = writeTimeout
	//endless.DefaultMaxHeaderBytes = maxHeaderBytes
	//server := endless.NewServer(endPoint, routersInit)
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
}
