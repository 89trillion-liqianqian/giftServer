package handler

import (
	"encoding/json"
	"giftServer/internal/model"
	"giftServer/internal/service"
	"strconv"
)

// 创建礼品码
func CreateGiftHandler(formData model.CreateGiftModels) (code string, err error) {
	//codeType:="1" 	// 1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
	codeType := formData.CodeType
	// 校验可领取参数
	if codeType == model.CodeTypeOne {
		if formData.DrawCount != 1 {
			formData.DrawCount = 1
		}
	}
	//code="SFDSHFUISD33"
	code = service.GetGiftCode()
	formData.Code = code
	err = model.SaveGiftRedis(formData)
	if codeType == model.CodeTypeTwo {
		// 保存
		model.SaveGiftRedisType(code)
	}
	return
}

// 查询礼品码
func GetGiftHandler(code string) (resData map[string]string, err error) {
	resData, err = model.GetGiftRedis(code)
	if len(resData) > 0 {
		// 获取历史记录
		historyData, _ := model.GetGiftRedis(code + "_history")
		b, _ := json.Marshal(historyData)
		resData["historyData"] = string(b)
	}
	return
}

// 领取礼品
func GetGiftRewardHandler(uid, code string) (content, msg string, err error) {
	resData, err := model.GetGiftRedis(code)
	if len(resData) < 1 {
		return
	}
	//  判有效性
	if _, ok := resData["validTime"]; ok {
		validTime := resData["validTime"]
		if !service.CheckTime(validTime) {
			// 过期，领取失败，
			msg = "过期，领取失败"
			return
		}
	}
	//uid:=0
	codeType := "1" // 1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
	if _, ok := resData["CodeType"]; ok {
		codeType = resData["CodeType"]
	}
	costCount := resData["CostCount"]
	costCountInt, _ := strconv.Atoi(costCount)
	if codeType == "1" {
		uidStr := resData["UserId"]
		if uid == uidStr && costCountInt < 1 {
			// 返礼品
			content = resData["Content"]
			//  消耗、保存
			err = model.SaveGiftCostCountRedis(code)
			err = model.SaveGiftCostHistoryRedis(code, uid)
			return
		} else {
			// 领取失败，
			msg = "失败，已领取"
			return
		}
	} else if codeType == "2" {
		//  限制次数
		drawCount := resData["DrawCount"]
		drawCountInt, _ := strconv.Atoi(drawCount)
		if costCountInt >= drawCountInt {
			msg = "礼品码，领取次数已完"
			return
		}
		// 领取礼品
		isOk, _ := model.GetGiftTypeTwo(code, drawCountInt)
		if isOk {
			// 返礼品
			content = resData["Content"]
			// 保存历史
			err = model.SaveGiftCostHistoryRedis(code, uid)
			return
		}
		msg = "领取失败,次数已使用完"
		return
	} else if codeType == "3" {
		// 返礼品
		content = resData["Content"]
		//  消耗、更新次数，保存历史
		err = model.SaveGiftCostCountRedis(code)
		err = model.SaveGiftCostHistoryRedis(code, uid)
		if err != nil {
			msg = "领取失败，请重试"
		}
		return
	}
	return
}
