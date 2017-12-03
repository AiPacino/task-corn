package service

import (
  "cailian/cron/model"
  "fmt"
  "encoding/json"
  "gopkg.in/redis.v5"
  "strconv"
  "time"
)

var KEY_ROLL_LIST_ALL = "cls_roll_list_all"
var KEY_ROLL_LIST_ = "cls_roll_list_"
var KEY_ARTICLE_LIST_ = "cls_article_list_"
var KEY_SPECIAL_TOPIC = "cls_special_topic_list"
var KEY_COLUMN_LIST = "cls_column_list"
var KEY_CONTENTION_LIST = "cls_contention_list"
var KEY_THEME_LIST = "cls_theme_list"
var KEY_ROLL_SHOW = "cls_roll_show"


func CacheArticle() {
  var handleRollList4s = time.NewTicker(time.Second * 4)    // 每4秒钟获取电报其他几种类型列表
  var handleRollList10s = time.NewTicker(time.Second * 10)  // 每10秒钟调度任务

  // 每4秒钟调度任务
  go func() {
    defer func() {
      if err := recover(); err != nil {
        fmt.Println("每4秒钟调度任务 recover error", err)
      }
    }()
    for range handleRollList4s.C {
      fmt.Println("每4秒钟调度任务")

      // 电报页面滚动的专题列表
      cacheSpecialTopicList()

      // 缓存电报列表
      cacheRollListByRedis(getParam("", 300))
      cacheRollListByRedis(getParam("explain", 100))
      cacheRollListByRedis(getParam("red", 100))
      cacheRollListByRedis(getParam("jpush", 100))
      cacheRollListByRedis(getParam("remind", 100))
      cacheRollListByRedis(getParam("vip_single", 200))
    }
  }()

  // 每10秒钟调度任务
  go func() {
    defer func() {
      if err := recover(); err != nil {
        fmt.Println("每10秒钟调度任务 recover error", err)
      }
    }()
    for range handleRollList10s.C {
      fmt.Println("每10秒钟调度任务")

      // top题材列表
      cacheArticleListByRedis(model.ArticleFilter{CType:[]int{1}, IsFocus:1, Rn:50})

      // pc题材列表
      cacheThemeListByRedis(model.ArticleFilter{CType:[]int{1}, Rn:200})

      // 短题材列表
      cacheArticleListByRedis(model.ArticleFilter{CType:[]int{8}, Rn:200})

      // 深度列表
      cacheArticleListByRedis(model.ArticleFilter{CType: []int{2}, Column_id:-1, Rn: 400})

      // pc版电报刷新列表
      cacheBatchShow()
    }
  }()
}

// 缓存电报页面滚动的专题列表
func cacheSpecialTopicList() {
  var key = KEY_SPECIAL_TOPIC
  var param = model.ArticleFilter{CType:[]int{21}, IsTop:1}
  ress, err := model.GetSpecialTopicList(dbCache, param)
  if err != nil {
    fmt.Println("cacheSpecialTopicList data error", err)
  }
  rt := redisCache.TxPipeline()

  // 删除当前key，如果删除错误终止该事务
  del := rt.Del(key)
  if del.Err() != nil {
    fmt.Println("del error ", del.Err())
    rt.Discard()
    return
  }

  for _, v := range ress {
    // 保存有序分页对象
    abytes, err := json.Marshal(v)
    if err != nil {
      fmt.Println("zadd Marshal--->>> ", err)
      continue
    }
    zadd := rt.ZAdd(key, redis.Z{v.SortScore, abytes})
    if zadd.Err() != nil {
      fmt.Println("zadd error--->>> ", zadd.Err())
      continue
    }
  }
  rt.Exec()
}

// 根据电报类型缓存不同key的电报集合
func cacheRollListByRedis(category string, param model.ArticleFilter) {
  var key = KEY_ROLL_LIST_ALL
  var rollListAllTemp, err = model.FindArticlesTime(dbCache, param)
  if err != nil {
    fmt.Println("cacheRollListByRedis get roll data error", err)
  }
  if len(rollListAllTemp) > 0 {
    if len(category) > 0 {
      key = KEY_ROLL_LIST_ + category
    }
    // 开启事务管道
    rt := redisCache.TxPipeline()

    // 删除当前key，如果删除错误终止该事务
    del := rt.Del(key)
    if del.Err() != nil {
      fmt.Println("del error ", del.Err())
      rt.Discard()
      return
    }
    // 写入新的list
    for _, v := range rollListAllTemp {
      // 保存有序分页对象
      abytes, err := json.Marshal(v)
      if err != nil {
        fmt.Println("zadd Marshal--->>> ", err)
        continue
      }

      zadd := rt.ZAdd(key, redis.Z{v.Sort_score, abytes})
      if zadd.Err() != nil {
        fmt.Println("zadd error--->>> ", zadd.Err())
        continue
      }
    }
    rt.Exec()
  }
}

func getParam(category string, limit uint64) (string, model.ArticleFilter) {
  param := model.ArticleFilter{CType: []int{-1}, Rn:model.PageLimit(limit)}
  switch category {
  case "red":
    param.Recommend = 1
    param.In_roll = 1
    param.CType = []int{-1}
    break
  case "jpush":
    param.Jpush = 1
    param.CType = []int{-1}
    break
  case "remind":
    param.CType = []int{12}
    break
  case "explain":
    param.CType = []int{-1}
    param.Explain = 1
    break
  case "pend_confirm":
    param.CType = []int{15}
    break
  case "special":
    param.CType = []int{16}
    break
  case "vip_single":
    param.CType = []int{15, 16}
    break
  case "vip_all":
    param.CType = []int{-1, 15, 16}
    break
  default:
    param.CType = []int{-1}
    param.In_roll = 1
    break
  }
  return category, param
}

// 根据类型缓存不同key的集合
func cacheArticleListByRedis(param model.ArticleFilter) {
  if len(param.CType) != 1 {
    fmt.Println("cacheArticleListByRedis type error ")
    return
  }
  var key = KEY_ARTICLE_LIST_ + strconv.Itoa(param.CType[0])
  var listTemp, err = model.FindArticlesTime(dbCache, param)
  if err != nil {
    fmt.Println("cacheArticleListByRedis database get roll data error \n", err)
  }
  if len(listTemp) > 0 {
    rt := redisCache.TxPipeline()

    // 删除当前key，如果删除错误终止该事务
    del := rt.Del(key)
    if del.Err() != nil {
      fmt.Println("cacheArticleListByRedis del error ", del.Err())
      rt.Discard()
      return
    }

    for _, v := range listTemp {
      // 保存有序分页对象
      abytes, err := json.Marshal(v)
      if err != nil {
        fmt.Println("cacheArticleListByRedis zadd Marshal--->>> ", err)
        continue
      }
      // 如果is_top 需要把时间设置为新的时间保证top在前面
      if v.Type == 2 && v.Column_id == -1 && v.Is_top > 0 {
        v.Sort_score = float64(time.Now().Unix()) + float64(v.Is_top)
      }
      zadd := rt.ZAdd(key, redis.Z{v.Sort_score, abytes})
      if zadd.Err() != nil {
        fmt.Println("cacheArticleListByRedis zadd error--->>> ", zadd.Err())
        continue
      }
    }
    rt.Exec()
  }
}

// 缓存PC版本题材列表
func cacheThemeListByRedis(param model.ArticleFilter) {
  if len(param.CType) != 1 {
    fmt.Println("cacheThemeListByRedis type error ")
    return
  }
  var key = KEY_THEME_LIST
  var listTemp, err = model.FindArticlesTime(dbCache, param)
  if err != nil {
    fmt.Println("cacheThemeListByRedis database get roll data error \n", err)
  }
  if len(listTemp) > 0 {
    rt := redisCache.TxPipeline()

    // 删除当前key，如果删除错误终止该事务
    del := rt.Del(key)
    if del.Err() != nil {
      fmt.Println("cacheThemeListByRedis del error ", del.Err())
      rt.Discard()
      return
    }

    for _, v := range listTemp {
      // 保存有序分页对象
      abytes, err := json.Marshal(v)
      if err != nil {
        fmt.Println("cacheThemeListByRedis zadd Marshal--->>> ", err)
        continue
      }
      zadd := rt.ZAdd(key, redis.Z{v.Sort_score, abytes})
      if zadd.Err() != nil {
        fmt.Println("cacheThemeListByRedis zadd error--->>> ", zadd.Err())
        continue
      }
    }
    rt.Exec()
  }
}

// 缓存PC版本获取电报列表(包含所有新增，修改，删除的数据)
func cacheBatchShow() {
  var key = KEY_ROLL_SHOW
  articles, err := model.GetArticles(dbCache, model.ArticleFilter{CType:[]int{-1}, In_roll:1, Rn:200})
  if err != nil {
    return
  }
  if len(articles) > 0 {
    rt := redisCache.TxPipeline()

    // 删除当前key，如果删除错误终止该事务
    del := rt.Del(key)
    if del.Err() != nil {
      fmt.Println("cacheBatchShow del error ", del.Err())
      rt.Discard()
      return
    }

    for _, v := range articles {
      // 保存有序分页对象
      abytes, err := json.Marshal(v)
      if err != nil {
        fmt.Println("cacheBatchShow zadd Marshal--->>> ", err)
        continue
      }
      zadd := rt.ZAdd(key, redis.Z{float64(v.Modified_time), abytes})
      if zadd.Err() != nil {
        fmt.Println("cacheBatchShow zadd error--->>> ", zadd.Err())
        continue
      }
    }
    rt.Exec()
  }
}
