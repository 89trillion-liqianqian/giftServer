package ctrl

import (
	"giftServer/internal/handler"
	"giftServer/internal/model"
	"giftServer/internal/myerr"
	"github.com/gin-gonic/gin"
	"log"
)

// 测试ping
func PingFunc(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
	return
}

// 创建礼品码
func CreateGift(c *gin.Context) {
	msg := "管理后台-创建礼品码 ok"
	var formData model.CreateGiftModels
	if err := c.ShouldBind(&formData); err != nil {
		log.Println("--CreateGift err", err)
		msg = "参数错误，管理后台-创建礼品码"
		myerr.ResponseErr(c, msg)
		return
	}
	// 参数校验
	if !(formData.CodeType == model.CodeTypeOne || formData.CodeType == model.CodeTypeTwo || formData.CodeType == model.CodeTypeThree) {
		msg = "参数错误,礼品码类型：1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换"
		log.Println("--CreateGift err", msg)
		myerr.ResponseErr(c, msg)
		return
	}
	// 创建礼品码
	code, err := handler.CreateGiftHandler(formData)
	if err != nil {
		msg = "创建礼品码失败，redis失败，请确认"
		myerr.ResponseErr(c, msg)
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  msg,
		"data": code,
	})
	return
}

// 查询礼品码
func GetGift(c *gin.Context) {
	code := c.Query("code")
	msg := "查询礼品码成功"
	// 校验参数
	if len(code) != 8 {
		msg = "err 礼品码是八位字符串"
		myerr.ResponseErr(c, msg)
		return
	}
	// 读取数据
	giftData, err := handler.GetGiftHandler(code)
	if err != nil {
		msg = "查询礼品码失败，redis失败，请确认 "
		myerr.ResponseErr(c, msg)
		return
	}
	if len(giftData) < 1 {
		msg = "查询礼品码失败，礼品码不存在，请确认"
		myerr.ResponseErr(c, msg)
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  msg,
		"data": giftData,
	})
	return
}

// 验证礼品码
func CheckCode(c *gin.Context) {
	code := c.PostForm("code")
	uid := c.PostForm("uid")
	msg := "验证礼品码成功"
	if code == "" || uid == "" {
		log.Println("--CheckCode-code", code, uid)
		msg = "验证礼品码 。参数错误"
		myerr.ResponseErr(c, msg)
		return
	}
	// 读取数据
	content, msg, err := handler.GetGiftRewardHandler(uid, code)
	if err != nil {
		myerr.ResponseErr(c, msg)
		return
	}
	if len(content) < 1 {
		msg += "获取失败"
		myerr.ResponseErr(c, msg)
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  msg,
		"data": content,
	})
	return
}
