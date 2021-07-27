package myerr

import "github.com/gin-gonic/gin"

// 错误返回
func ResponseErr(c *gin.Context, msg string, code int) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": "",
	})
	return
}
