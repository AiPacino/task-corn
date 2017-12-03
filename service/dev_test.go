package service

import (
  "testing"
  "os"
  "fmt"
  "time"
  "cailian/cron/model"
  "encoding/json"
)

func TestCloseRedis(t *testing.T) {
  // redis service
  os.Setenv("REDIS_ADDR", "101.201.34.157:6333")
  os.Setenv("REDIS_PWD", "Kpe093lss03s")
  OpenRedis()
  defer CloseRedis()

  //GetIndustryRank()
  //GetStockRanking()
  //synMarketDaily()
  //GetIndustryRank()   // 领涨行业板块
  //GetStockRanking()   // 沪深排行榜
  StartLoop()

  time.Sleep(time.Second * 50)
}

func TestSubBaidu(t *testing.T) {
  os.Setenv("WRITER_DB_URL", "cls_dba:!fd@3hAF01P#x7Cm@tcp(rdsuieavryje7vmo.mysql.rds.aliyuncs.com:3306)/cls_1508?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&readTimeout=15s&timeout=5s")
  // mysql writer service
  OpenWriterDB()
  defer CloseWriterDB()

  var i uint64
  for i < 80000 {
    i = i + 2000
    fmt.Println("i", i)
    baiduPushUrl(i)
  }
}

func TestSynMarketDaily(t *testing.T) {
  // redis service
  //os.Setenv("REDIS_ADDR", "101.201.34.251:8779")
  //os.Setenv("REDIS_PWD", "Cailianpress2017")

  os.Setenv("REDIS_ADDR", "101.201.34.157:6333")
  os.Setenv("REDIS_PWD", "Kpe093lss03s")

  OpenRedis()
  defer CloseRedis()

  //synMarketDaily()

  //getBlockList()

  var key = "cls_market_daily_up_lastest"
  var bakKey = "cls_market_daily_up:"
  var url = "http://api.cailianpress.com/v1/Redis_action/redis_get_find?debug=Cp4UFGL8QJ46sXhC26h9cf54fX24HeRD&rediskey=cls_market_daily_up_lastest"
  var res = HttpGet(url)
  var mdr = ResultData{}
  var err = json.Unmarshal([]byte(res), &mdr)
  if err != nil {
    fmt.Println("synMarketDaily Unmarshal error", err)
  }

  var requestData = MarketDataResult{}
  json.Unmarshal([]byte(mdr.Data), &requestData)
  bakKey = bakKey + time.Unix(requestData.Summary.Time, 0).Format("20060102")
  redisCache.Set(bakKey, mdr.Data, 0)
  redisCache.Set(key, mdr.Data, 0)
}

func TestRedDot(t *testing.T) {
  os.Setenv("WRITER_DB_URL", "cailianpress_dba:cailianpress_888@tcp(rm-2ze52963a722633wzo.mysql.rds.aliyuncs.com:3306)/cls_1508?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&readTimeout=15s&timeout=5s")
  // mysql writer service
  OpenWriterDB()
  defer CloseWriterDB()

  os.Setenv("REDIS_ADDR", "47.94.178.76:6379")
  os.Setenv("REDIS_PWD", "")
  OpenRedis()
  defer CloseRedis()

  tu := time.Now().Unix() - 990
  fmt.Println("tu = ", tu)
  res := model.GetNewArticle(dbCacheWriter, tu)
  if len(res) > 0 {
    for _, v := range res {
      AppLog.Infof("handleTime30s \n", v)
      CreateReddot(v.ColumnID, v.Type, v.Ctime)
    }
  }
}

func TestCacheArticleListByRedis(t *testing.T) {
  //os.Setenv("WRITER_DB_URL", "cls_dba:!fd@3hAF01P#x7Cm@tcp(rdsuieavryje7vmo.mysql.rds.aliyuncs.com:3306)/cls_1508?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&readTimeout=15s&timeout=5s")
  //os.Setenv("DB_URL", "cailianpress_dba:cailianpress_888@tcp(rm-2ze52963a722633wzo.mysql.rds.aliyuncs.com:3306)/cls_1508?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&readTimeout=15s&timeout=5s")

  // 删除收藏表中的重复数据
  // mysql writer service
  //OpenDB()
  //defer CloseDB()

  //os.Setenv("REDIS_ADDR", "101.201.34.157:6333")
  //os.Setenv("REDIS_PWD", "Kpe093lss03s")

  os.Setenv("REDIS_ADDR", "101.201.34.251:8779")
  os.Setenv("REDIS_PWD", "Cailianpress2017")
  OpenRedis()
  defer CloseRedis()

  // 缓存电报列表
  //cacheRollListByRedis(getParam(""))
  //cacheRollListByRedis(getParam("explain"))
  //cacheRollListByRedis(getParam("red"))
  //cacheRollListByRedis(getParam("jpush"))
  //cacheRollListByRedis(getParam("remind"))
  //cacheRollListByRedis(getParam("vip_single"))

  //// top题材列表
  //cacheArticleListByRedis(model.ArticleFilter{CType:[]int{1}, IsFocus:1, Rn:50})
  //
  //// 短题材列表
  //cacheArticleListByRedis(model.ArticleFilter{CType:[]int{8}, Rn:200})
  //
  //// 深度列表
  //cacheArticleListByRedis(model.ArticleFilter{CType: []int{2}, Column_id:-1, Rn: 200})
  //
  //// pc版电报刷新列表
  //cacheBatchShow()

  // pc题材列表
  //cacheThemeListByRedis(model.ArticleFilter{CType:[]int{1}, Rn:200})

  StartLoop()

  time.Sleep(time.Second * 50)
}
