package handler

import (
	"giftServer/internal/model"
	"log"
	"testing"
	"time"
)

var Code = ""

func init() {
	filepath := "../../config/app.ini"
	model.GetAppIni(filepath)
	model.Init()
}

// 测试创建礼品码
func TestCreateGiftHandler(t *testing.T) {
	formData := model.CreateGiftModels{}
	formData.CodeType = model.CodeTypeThree
	formData.DrawCount = 20
	formData.Des = "des 这是金币的礼品码"
	formData.ValidTime = time.Now().Unix() + 7*24*60
	formData.Content = `{"1":1000,"2":10000}`
	formData.CreateUser = "qq"
	formData.UserId = 123456
	code, err := CreateGiftHandler(formData)
	if err != nil {
		log.Println("err 测试创建礼品码")
	}
	if len(code) == 8 {
		log.Println("ok 测试创建礼品码")
	}
	Code = code
	return
}

// 测试获取礼品码
func TestGetGiftHandler(t *testing.T) {
	resData, err := GetGiftHandler(Code)
	if err != nil {
		log.Println("err 测试获取礼品码")
	}
	if len(resData) < 1 {
		log.Println("err 测试获取礼品码")
	}
	log.Println("ok 测试获取礼品码", resData)
	return
}

// 测试验证礼品码
func TestGetGiftRewardHandler(t *testing.T) {
	uid := "0010"
	content, msg, err := GetGiftRewardHandler(uid, Code)
	if err != nil {
		log.Println("err 测试获取礼品码", msg)
	}
	if len(content) < 1 {
		log.Println("err 获取礼品码内容", msg)
		return
	}
	log.Println("ok 测试验证礼品码获取礼品码", content)
	return
}
