package model

import (
	"errors"
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"log"
	"time"
)

/**
redis 初始化
*/

var (
	RedisPool *redigo.Pool
)

// redis pool

func PoolInitRedis(server string, password string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     2, //空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   3, //最大数
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Init() {
	var addr = "127.0.0.1:6379"
	//var addr = "redis:6379"
	var password = "123456"
	//var password = ""
	RedisPool = PoolInitRedis(addr, password)
	// 测试
	conn := RedisPool.Get()
	res, err := redigo.String(conn.Do("set", "qq", "start"))
	log.Println("--res", res, err)
	res, err = redigo.String(conn.Do("get", "qq"))
	log.Println("-2-res", res, err)
}

// 加锁

func Lock() (ok bool, err error) {
	c := RedisPool.Get()
	defer c.Close()
	//设置锁key-value和过期时间
	//_, err = redigo.String(c.Do("SET", "lock_key", "lock_value", "EX", 10*time.Second, "NX"))
	_, err = redigo.String(c.Do("SET", "lock_key", "lock_value", "EX", 10, "NX"))
	if err != nil {
		if err == redigo.ErrNil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 解锁

func Unlock(value string) (err error) {
	c := RedisPool.Get()
	defer c.Close()
	//获取锁value
	setValue, err := redigo.String(c.Do("GET", "lock_key"))
	if err != nil {
		return
	}
	//判断锁是否属于该释放锁的线程
	if setValue != value {
		err = errors.New("非法用户，无法释放该锁")
		return
	}
	//属于该用户，直接删除该key
	_, err = c.Do("DEL", "lock_key")
	return
}

func GetAction() {
	count := 1
	for i := 1; i <= 10; i++ {
		log.Println("--test")
	RETRY:
		count += 1
		lock, err := Lock()
		log.Println("-err", lock, err)
		if !lock {
			// 取消设置
			if count > 10 {
				return
			}
			//return
			// 重试
			goto RETRY
		}
		log.Println("--ll", i)
		log.Println("-2-ll", i)
	}
	fmt.Printf("end")
}
