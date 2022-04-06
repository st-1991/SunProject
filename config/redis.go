package config

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var Redis redis.Conn

const (
	Prefix = "KeepHappy:"
	Db = 1
)

type RedisConn struct {
	Redis redis.Conn
}

type RedisKey string

func init() {
	var err error

	Redis, err = redis.Dial("tcp", "127.0.0.1:6379", redis.DialDatabase(Db))
	if err != nil {
		fmt.Println("conn redis failed", err)
	}
}


func (k RedisKey) PrefixKey() RedisKey {
	return Prefix + k
}

func (k RedisKey) Set(val interface{}) RedisKey {
		_, err := Redis.Do("set", k, val)
		if err != nil {
			Logger().Error(fmt.Sprintf("redis错误：%s", err))
			return k
		}
		return k
}

func (k RedisKey) Get(key string) (reply interface{}, err error) {
	return Redis.Do("get", key)
}

func (k RedisKey) Expire(seconds int) bool {
		_, err := Redis.Do("expire", k, seconds)
		if err != nil {
			return false
		}
		return true
}


