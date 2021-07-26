package router

import (
	"giftServer/internal/ctrl"
	"github.com/gin-gonic/gin"
)

// 路由管理
func Router(r *gin.Engine) {
	r.GET("/ping", ctrl.PingFunc)
	r.GET("/getGift", ctrl.GetGift)
	r.POST("/createGift", ctrl.CreateGift)
	r.POST("/checkCode", ctrl.CheckCode)
}
