package service

import (
	"SunProject/application/models"
	"SunProject/config"

	"github.com/garyburd/redigo/redis"

	"strconv"
	"time"
)

type UserDynamic struct {
	Id, UserId int
}

func DynamicThumbUp(ud UserDynamic) {
	// 判断是否已经点赞
	if ud.IsThumbUp() {
		return
	}
	ud.setThumbUpCache()
	db := config.DB.Begin()
	// 增加动态点赞数
	d := models.Dynamic{Id: ud.Id}
	if dOk := d.IncrColumn(db, "thumb_up"); !dOk {
		ud.delThumbUpCache()
		db.Rollback()
		return
	}
	// 增加用户获赞数
	u := models.User{Id: ud.UserId}
	if uOk := u.SetThumbUp(db); !uOk {
		ud.delThumbUpCache()
		db.Rollback()
		return
	}
	db.Commit()
	config.Logger().Info("点赞成功")
}

func (ud *UserDynamic) getCacheKey() config.RedisKey {
	redisKey := config.RedisKey("thumb_up:" + strconv.Itoa(ud.Id))
	return redisKey.PrefixKey()
}

// SetCache 设置点赞缓存
func (ud *UserDynamic) setThumbUpCache() bool {
	if _, err := config.Redis.Do("HMSET", ud.getCacheKey(), ud.UserId, time.Now().Format("2006-01-02 15:04:05")); err != nil {
		config.Logger().Error("点赞失败 redis HMSET error:", err)
		return false
	}
	return true
}

func (ud UserDynamic) delThumbUpCache() bool {
	// 删除点赞缓存
	if _, err := config.Redis.Do("HDEL", ud.getCacheKey(), ud.UserId); err != nil {
		config.Logger().Error("点赞失败 redis HDEL error:", err)
		return false
	}
	return true
}

func (ud UserDynamic) IsThumbUp() bool {
	// 判断是否已经点赞
	isTU, err := redis.Bool(config.Redis.Do("HEXISTS", ud.getCacheKey(), ud.UserId))
	if err != nil {
		config.Logger().Error("点赞失败 redis HEXISTS error:", err)
		return false
	}
	return isTU
}

func AddCommentNum(id, commentId,level int) {
	if level == 1 {
		// 增加评论数
		d := models.Dynamic{Id: id}
		if dOk := d.IncrColumn(config.DB, "comment_num"); !dOk {
			config.Logger().Error("增加评论数失败")
		}
	} else {
		// 增加回复数
		c := models.Comments{Id: commentId}
		if cOk := c.IncrColumn(config.DB, "comment_num"); !cOk {
			config.Logger().Error("增加回复数失败")
		}
	}
}