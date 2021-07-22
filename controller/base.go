package controller

import (
	"giftServer/model"
	"github.com/gin-gonic/gin"
	"log"
)

/**
http get post 处理
*/

func Routers(r *gin.Engine) {
	r.GET("/ping", pingFunc)
	r.GET("/getGift", GetGift)
	r.POST("/createGift", CreateGift)
	r.POST("/checkCode", CheckCode)
	return
}

// 测试
func pingFunc(c *gin.Context) {
	model.GetAction()
	c.JSON(200, gin.H{
		"message": "pong22",
	})
	return
}

// 创建礼品码

func CreateGift(c *gin.Context) {
	var formData model.CreateGiftModels
	if err := c.ShouldBind(&formData); err != nil {
		log.Println("--CreateGift err", err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "参数错误，管理后台-创建礼品码",
			"data": "",
		})
		return
	}
	// 保存redis
	code, err := model.CreateGiftModel(formData)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "创建礼品码失败，redis失败，请确认 ",
			"data": "",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "管理后台-创建礼品码 ok",
		"data": code,
	})
	return
}

// 查询礼品码

func GetGift(c *gin.Context) {
	code := c.Query("code")
	log.Println("--GetGift-code", code)
	// 从redis 读取数据 redis
	giftData, err := model.GetGiftModel(code)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "查询礼品码失败，redis失败，请确认 ",
			"data": "",
		})
		return
	}
	if len(giftData) < 1 {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "查询礼品码失败，礼品码不存在，请确认 ",
			"data": "",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "查询礼品码成功",
		"data": giftData,
	})
	return
}

// 验证礼品码

func CheckCode(c *gin.Context) {
	code := c.PostForm("code")
	uid := c.PostForm("uid")
	if code == "" || uid == "" {
		log.Println("--CheckCode-code", code, uid)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "验证礼品码 。参数错误",
			"data": "",
		})
		return
	}
	log.Println("--CheckCode-code", code, uid)
	// 从redis 读取数据 redis
	content, msg, err := model.GetGiftReward(uid, code)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
			"data": "",
		})
		return
	}
	if len(content) < 1 {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
			"data": "",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "验证礼品码成功",
		"data": content,
	})
	return
}
