package main

import (
	"giftServer/app/http"
	"giftServer/internal/model"
)

// 第三题，入口
func main() {
	// 加载配置
	filepath := "../config/app.ini"
	model.GetAppIni(filepath)
	//初始化redis
	model.Init()
	// 启动http 服务
	http.HttpServer()
}
