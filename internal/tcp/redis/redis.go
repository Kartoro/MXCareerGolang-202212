package redis

import (
	"MXCareerGolang-202212/config"
	"MXCareerGolang-202212/internal/model"
	"MXCareerGolang-202212/internal/util"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

var cache *redis.Client

func init() {
	cache = redis.NewClient(&redis.Options{
		Addr:         config.RedisAddr,
		MinIdleConns: 5000,
		ReadTimeout:  -1,
		WriteTimeout: -1,
		PoolTimeout:  120 * time.Second,
		//PoolSize:     config.RedisPoolSize,
	})
	_, err := cache.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis connect failed")
		return
	}
	log.Println("redis connected")

	WarmUpCache()

	return
}

func set(key string, value interface{}, expiration int64) error {
	err := cache.Set(context.Background(), key, value, time.Duration(expiration*1e9)).Err()
	if err != nil {
		log.Printf("Redis Set Failed! key: %v value: %v err:%v\n", key, value, err.Error())
		return err
	}
	return nil
}

func get(key string) (string, error) {
	value, err := cache.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("Redis Get Failed! key: %v value:%v err: %v\n", key, value, err)
		return "", err
	}
	return value, err
}

func LoginAuth(userName string, password string) bool {
	var pwdKey = userName + "_pwd"
	if isExits := cache.Exists(context.Background(), pwdKey).Val(); isExits == 0 { // 0 为不存在
		return false
	}
	pwd, err := get(pwdKey)
	if err != nil {
		return false
	}
	if pwd == util.SHA256(password) {
		return true
	}
	return false
}

func SetPassword(userName string, password string) {
	_ = set(userName+"_pwd", util.SHA256(password), int64(config.MaxExTime))
}

func GetProfile(userName string) (model.User, bool) {
	if isExits := cache.Exists(context.Background(), userName).Val(); isExits == 0 {
		return model.User{}, false
	}

	values, err := cache.HGetAll(context.Background(), userName).Result()
	if err != nil {
		return model.User{}, false
	}
	return model.User{NickName: values["nick_name"], Avatar: values["pic_name"]}, true
}

func SetProfile(userName string, nickName string, picName string) error {
	fields := map[string]interface{}{
		"nick_name": nickName,
		"pic_name":  picName,
	}
	err := cache.HMSet(context.Background(), userName, fields).Err()
	if err != nil {
		return err
	}
	return nil
}

// DelCache - del before write in mysql to ensure consistency
func DelCache(userName string) error {
	err := cache.Del(context.Background(), userName).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetToken - save token and expiration to redis
func SetToken(userName string, token string, expiration int64) error {
	return set("auth_"+userName, token, expiration)
}

func CheckToken(userName string, token string) (bool, error) {
	val, err := get("auth_" + userName)
	if err != nil {
		return false, err
	}
	return token == val, nil
}

func WarmUpCache() {
	log.Println("redis warm up start")
	for i := 1; i <= 200; i++ {
		s := strconv.Itoa(i)
		SetPassword(s, s)
	}
	log.Println("redis warm up done")
}
