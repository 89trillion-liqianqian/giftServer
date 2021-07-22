package utils

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// 随机字符串

func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

//获取八位礼品码

func GetGiftCode() (code string) {
	code = string(Krand(8, 3))
	return
}

// 校验是否过期

func CheckTime(validTime string) (isValid bool) {
	validTimeInt, _ := strconv.Atoi(validTime)
	validTimeInt64 := int64(validTimeInt)
	nowTime := time.Now().Unix()
	if validTimeInt64 > nowTime {
		isValid = true
		return
	}
	return
}
