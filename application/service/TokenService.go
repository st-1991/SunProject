package service

import (
	"SunProject/config"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/robbert229/jwt"
	"strings"
	"time"
)

const secret = "keepHappySecret888"

type Jwt struct {
	Algorithm jwt.Algorithm
	Claims *jwt.Claims
}

type Token string

type UserData struct {
	ID int `json:"id"`
	Phone string `json:"phone"`
}

func (j *Jwt) initJwt() {
	j.Algorithm = jwt.HmacSha256(secret)
	j.Claims = jwt.NewClaim()
}


func (j Jwt) CreateToken(key string, value string, seconds int64) (Token, error) {
	j.initJwt()
	j.Claims.Set(key, value)
	j.Claims.SetTime("exp", time.Now().Add(time.Second * time.Duration(seconds)))
	jwtStr, err := j.Algorithm.Encode(j.Claims)
	if err != nil {
		return "", err
	}
	token := strings.Split(jwtStr, ".")[2]

	ok := config.RedisKey("tokens:" + token).PrefixKey().Set(jwtStr).Expire(seconds)
	if !ok {
		return "", fmt.Errorf("redis set error")
	}
	return Token(token), nil
}


func (t Token) Validate() bool {
	redisKey := config.RedisKey("tokens:" + t)
	jwtStr, err := redis.String(redisKey.PrefixKey().Get())
	if err != nil {
		config.Logger().Error("redis get error:", err)
		return false
	}
	j := Jwt{}
	j.initJwt()
	if j.Algorithm.Validate(jwtStr) != nil {
		return false
	}
	return true
}


func (t Token) GetUserInfo(key string) (UserData, error) {
	redisKey := config.RedisKey("tokens:" + t)
	jwtStr, err := redis.String(redisKey.PrefixKey().Get())
	if err != nil {
		return UserData{}, err
	}
	j := Jwt{}
	j.initJwt()
	loadedClaims, err := j.Algorithm.Decode(jwtStr)
	if err != nil {
		return UserData{}, err
	}

	userStr, err := loadedClaims.Get(key)
	if err != nil {
		return UserData{}, err
	}
	userData := UserData{}
	err = json.Unmarshal([]byte(userStr.(string)), &userData)
	if err != nil {
		return UserData{}, err
	}
	return userData, nil
}