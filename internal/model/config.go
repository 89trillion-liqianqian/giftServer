package model

import (
	"github.com/Unknwon/goconfig"
	"log"
)

/**
app配置文件解析
*/
var cfg *goconfig.ConfigFile

func GetAppIni(filepath string) (err error) {
	// 解析app配置
	config, err := goconfig.LoadConfigFile(filepath)
	if err != nil {
		log.Println("配置文件读取错误,找不到配置文件", err)
		return err
	}
	cfg = config
	return nil
}

// 获取端口号
func GetAppPort() (HttpPort string, err error) {
	// 获取app端口
	if HttpPort, err = cfg.GetValue("server", "HttpPort"); err != nil {
		log.Println("配置文件中不存在types", err)
		return HttpPort, nil
	}

	return HttpPort, nil
}

// 获取配置
func GetKey(sec, keyStr string) (value string, err error) {
	if value, err = cfg.GetValue(sec, keyStr); err != nil {
		log.Println("配置文件中不存在types", err)
		return value, nil
	}
	return
}
