package config

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"strings"
	"time"
)

var Redis *redis.Pool

const (
	Prefix = "KeepHappy:"
	Db = 1
)

type RedisConn struct {
	Redis redis.Conn
}

type RedisKey string

func init() {
	var (
		host string
		auth string
		db   int
	)
	host = "127.0.0.1:6379"
	auth = ""
	db = 0
	Redis = &redis.Pool{
		MaxIdle:     100,
		MaxActive:   4000,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host, redis.DialPassword(auth), redis.DialDatabase(db))
			if nil != err {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	//var err error
	//Redis, err = redis.Dial("tcp", "127.0.0.1:6379", redis.DialDatabase(Db))
	//if err != nil {
	//	panic("conn redis failed" + err.Error())
	//}
}

func IsConnError(err error) bool {
	var needNewConn bool

	if err == nil {
		return false
	}

	if err == io.EOF {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connect: connection refused") {
		needNewConn = true
	}
	return needNewConn
}

func Redo(command string, opt ...interface{}) (interface{}, error) {
	rd := Redis.Get()
	defer rd.Close()

	var conn redis.Conn
	var err error
	var maxretry = 3
	var needNewConn bool

	resp, err := rd.Do(command, opt...)
	needNewConn = IsConnError(err)
	if needNewConn == false {
		return resp, err
	} else {
		conn, err = Redis.Dial()
	}

	for index := 0; index < maxretry; index++ {
		if conn == nil && index+1 > maxretry {
			return resp, err
		}
		if conn == nil {
			conn, err = Redis.Dial()
		}
		if err != nil {
			continue
		}

		resp, err := conn.Do(command, opt...)
		needNewConn = IsConnError(err)
		if needNewConn == false {
			return resp, err
		} else {
			conn, err = Redis.Dial()
		}
	}

	conn.Close()
	return "", errors.New("redis error")
}


func (k RedisKey) PrefixKey() RedisKey {
	return Prefix + k
}

func (k RedisKey) Set(val interface{}) RedisKey {
	_, err := Redo("set", k, val)
	if err != nil {
		Logger().Error(fmt.Sprintf("redis错误：%s", err))
		return k
	}
	return k
}

func (k RedisKey) Get() (reply interface{}, err error) {
	return Redo("get", k)
}

func (k RedisKey) Expire(seconds int64) bool {
	_, err := Redo("expire", k, seconds)
	if err != nil {
		return false
	}
	return true
}

func (k RedisKey) Del() bool {
	_, err := Redo("del", k)
	if err != nil {
		return false
	}
	return true
}


