package service

import (
  "gopkg.in/redis.v5"
  "github.com/wangyibin/alilog"
  "os"
)

var AppLog = alilog.New("cailian-cron", "api_log")
var redisLog = AppLog.With("file", "service.redis.go")

var redisCache  *redis.Client
func  OpenRedis() *redis.Client{
  addr := os.Getenv("REDIS_ADDR")
  pwd := os.Getenv("REDIS_PWD")
  redisCache = redis.NewClient(&redis.Options{Addr:addr, Password:pwd, DB:0})
  sc := redisCache.Ping()
  if sc.Err() == nil {
    redisLog.With("method", "OpenRedis").Infof(redisCache.String() + "\n")
    redisLog.With("method", "OpenRedis").Infof("Redis conection set up successfully \n")
  } else {
    redisLog.With("method", "OpenRedis").Error(sc.Err())
  }
  return redisCache
}

func CloseRedis() {
  redisCache.Close()
}

func GetRedisCache() *redis.Client {
  return redisCache
}
