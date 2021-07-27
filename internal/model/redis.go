package model

import (
	"errors"
	redigo "github.com/gomodule/redigo/redis"
	"log"
	"time"
)

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

// redis 初始化
func Init() {
	//var addr = "127.0.0.1:6379"
	//var addr = "redis:6379"
	//var password = "123456"
	addr, err := GetKey("main-redis", "Host")
	if err != nil {
		addr = "127.0.0.1:6379"
	}
	password, err := GetKey("main-redis", "Password")
	if err != nil {
		password = ""
	}
	log.Println("---redis", addr, password)
	//var password = ""
	RedisPool = PoolInitRedis(addr, password)
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
