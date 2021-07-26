package myerr

import "github.com/gin-gonic/gin"

// 错误返回
func ResponseErr(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  msg,
		"data": "",
	})
	return
}
