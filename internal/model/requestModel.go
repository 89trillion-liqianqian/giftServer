package model

import (
	redigo "github.com/gomodule/redigo/redis"
	"log"
	"time"
)

const GiftType = "_type"

const (
	//1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
	CodeTypeOne   = 1 //指定用户一次性消耗
	CodeTypeTwo   = 2 //不指定用户限制兑换次数
	CodeTypeThree = 3 //不限用户不限次数兑换

	CodeTypeOneStr   = "1" //指定用户一次性消耗
	CodeTypeTwoStr   = "2" //不指定用户限制兑换次数
	CodeTypeThreeStr = "3" //不限用户不限次数兑换
)

// 礼品码数据结构
type CreateGiftModels struct {
	Code       string `form:"code" binding:""`               // 礼品码
	CodeType   int    `form:"codeType" binding:"required"`   //礼品码类型，1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
	DrawCount  int    `form:"drawCount" binding:"required"`  //可领取次数
	Des        string `form:"des" binding:"required"`        //描述
	ValidTime  int64  `form:"validTime" binding:"required"`  //过期时间戳
	Content    string `form:"content" binding:"required"`    //奖品内容
	CreateUser string `form:"createUser" binding:"required"` //创建着
	CostCount  int    `form:"costCount" binding:""`          //已领取次数
	UserId     int    `form:"userId" binding:""`             //知道用户
}

// 领取限制次数的礼品码
func GetGiftTypeTwo(code string, drawCount int) (isOk bool, err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	codeType := code + "_type"
	count := 1
RETRY:
	count += 1
	lock, err := Lock()
	if !lock {
		// 取消设置
		if count > 100 {
			return
		}
		// 重试
		goto RETRY
	}
	// 获取领取次数
	costCount, err := redigo.Int(conn.Do("GET", codeType))
	if err != nil {
		Unlock("lock_value")
		log.Println("--getGiftTypeTwo", costCount, err)
		return
	}
	if costCount >= drawCount {
		Unlock("lock_value")
		log.Println("--getGiftTypeTwo,领取失败，次数已使用完", costCount, err)
		return
	}
	res, err := redigo.Int64(conn.Do("INCR", codeType))
	if err != nil {
		Unlock("lock_value")
		log.Println("--getGiftTypeTwo -add ", codeType, res, err)
		return
	}
	res, err = redigo.Int64(conn.Do("HINCRBY", code, "CostCount", 1))
	if err != nil {
		Unlock("lock_value")
		log.Println("--saveGiftCostCountRedis", res, err)
	}
	isOk = true
	Unlock("lock_value")
	return
}

// 保存礼品码的领取用户历史
func SaveGiftCostHistoryRedis(code, uid string) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	code += "_history"
	nowTime := time.Now().Unix()
	res, err := redigo.Int64(conn.Do("HSET", code, uid, nowTime))
	if err != nil {
		log.Println("--saveGiftCostHistoryRedis", res, err)
	}
	return
}

// 保存礼品码的领取次数数据
func SaveGiftCostCountRedis(code string) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()
	log.Println("----test", code, 1)
	res, err := redigo.Int64(conn.Do("HINCRBY", code, "CostCount", 1))
	if err != nil {
		log.Println("--saveGiftCostCountRedis", res, err)
	}
	return
}

// 获取礼品数据
func GetGiftRedis(code string) (resMap map[string]string, err error) {
	conn := RedisPool.Get()
	defer conn.Close()

	resMap, err = redigo.StringMap(conn.Do("HGETAll", code))
	if err != nil {
		log.Println("--getGiftRedis", resMap, err)
	}
	return
}

// 保存礼品数据
func SaveGiftRedis(formData CreateGiftModels) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()

	res, err := redigo.String(conn.Do("HMSET", redigo.Args{formData.Code}.AddFlat(formData)...))
	if err != nil {
		log.Println("--saveGiftRedis", res, err)
	}
	return
}

// 保存礼品数据
func SaveGiftRedisType(code string) (err error) {
	conn := RedisPool.Get()
	defer conn.Close()

	// 限制兑换次数
	code += GiftType
	res, err := redigo.String(conn.Do("set", code, 0))
	if err != nil {
		log.Println("--saveGiftRedisType", res, err)
	}
	return
}
