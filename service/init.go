package service

import (
  "time"
  "fmt"
  "strconv"
  "encoding/json"
  "cailian/cron/model"
  "os"
  "strings"
)

const ReadingNumPrefix = "cls_reading_num_" //快讯阅读数
type (
  ReadNumView struct {
    Aid       string    `json:"aid"`
    Num       int64     `json:"num"`       // 放大文章阅读数
    DetailNum int       `json:"detailnum"` // 真实文章阅读总数
    RollNum   int       `json:"rollnum"`   // 真实列表阅读总数
  }

  BaiduModelView struct {
    Id   int   `json:"id"`
    Type int   `json:"type"`
    ColumnId     int   `json:"column_id"`
  }

  ResultData struct {
    Data string `json:"data"`
    Errno int `json:"errno"`
    Time  int `json:"time"`
  }

  MarketDataResult struct {
    Summary Summary `json:"summary"`
  }

  Summary struct {
    D5   int `json:"d5"`
    Ld   int `json:"ld"`
    Lu   int `json:"lu"`
    Time int64 `json:"time"`
    U5   int `json:"u5"`
  }
)

func StartTimingTask() {

  // 每30秒检查一次文章是否有更新。
  handleTime30s := time.NewTicker(time.Second * 30)
  go func() {
    for range handleTime30s.C {
      tu := time.Now().Unix() - 90
      res := model.GetNewArticle(dbCacheWriter, tu)
      if len(res) > 0 {
        for _, v := range res {
          AppLog.Infof("handleTime30s \n", v)
          CreateReddot(v.ColumnID, v.Type, v.Ctime)
        }
      }
    }
  }()

  // 每60秒检查一次 更新手机版本文章阅读数redis 同步到pc版本阅读数redis
  handleTime60s := time.NewTicker(time.Second * 60)
  go func() {
    for range handleTime60s.C {
      AppLog.Infof("handleTime60s \n")
      updateReadNumToPc()
      setReadNumToDatabase()
      synMarketDaily()
    }
  }()

  // 每1小时更新一次百度推送数据
  if os.Getenv("IS_DEBUG") != "open" {
    handleTime6h := time.NewTicker(time.Hour * 1)
    go func() {
      for range handleTime6h.C {
        AppLog.Infof("handleTime6h \n")
        baiduPushUrl(0)
      }
    }()
  }

  // 每12小时检查一次
  handleTime1d := time.NewTicker(time.Hour * 12)
  go func() {
    getBlockList()
    for range handleTime1d.C {
      getBlockList()
    }
  }()
}

// http://api.cailianpress.com/v1/Redis_action/redis_get_find?debug=Cp4UFGL8QJ46sXhC26h9cf54fX24HeRD&rediskey=cls_market_daily_up_lastest
// 同步涨跌信息
func synMarketDaily() {
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

// 红点
/**
 > type = -1 是标示电报有红点
 > type = 1 是标示资讯模块题材有红点
 > type = 2 是标示资讯模块深度有红点
 > type = 3 是标示资讯模块早报有红点
 > type = 6 是标示内参
 > type = 8 是标示新版本题材
 > type = 9 是标示动向
 > type = -2 是标示监控模块最新有红点
 > type = -3 是标示监控模块深度有红点
 > type = -5 是标示异动
 > topic表示监控里面每个题材的ID
 */
func CreateReddot(columnID, typeVal int, ctime int) {
  if columnID <= 0 {
    t := strconv.Itoa(ctime)
    remindReddot := make(map[string]string, 0)
    remindReddot[strconv.Itoa(typeVal)+"_0"] = t
    if err := redisCache.HMSet("cls_remind_reddot_new", remindReddot).Err(); err != nil {
      AppLog.Errorf("CreateReddot redis HMSet err = %v\n", err)
    }
  }

  // 内参
  if columnID > 0 {
    if err := redisCache.HSet("cls_remind_reddot_new", "6_" + strconv.Itoa(columnID), ctime).Err(); err != nil {
      AppLog.Errorf("CreateReddot redis HSet err = %v\n", err)
    }
  }
}

// 每一分钟执行，更新手机版本文章阅读数redis 同步到pc版本阅读数redis
func updateReadNumToPc() {
  ids, ctimes := model.GetNewArticleIds(dbCacheWriter)
  if len(ids) < 1 {
    AppLog.With("init.func", "updateReadNumToPc").Errorf("GetNewArticleIds error \n")
    return
  }
  rnv := make([]ReadNumView, 0)
  for i, v := range ids {
    // 放大阅读数
    key := ReadingNumPrefix + v
    readNum, _ := redisCache.Get(key).Int64()
    r := ReadNumView{Aid:v, Num:readNum}

    // 真实阅读总数
    readKey := "cls_article_read_" + v
    if redisCache.Exists(readKey).Val() {
      m := redisCache.HGetAll(readKey).Val()
      //r.RollNum, _ = strconv.Atoi(m["lall"])
      r.DetailNum, _ = strconv.Atoi(m["all"])

      // 网页版本和pc版本阅读数 cls_real_detail_reading_  cls_real_roll_reading_
      rn, _ := redisCache.Get("cls_real_roll_reading_" + v).Int64()
      r.RollNum = int(rn)

      if ctimes[i] > time.Now().Unix() - 600 {
        m["lmin"] = strconv.FormatInt(rn, 10)
      }
      if ctimes[i] > time.Now().Unix() - 2400 {
        m["lmax"] = strconv.FormatInt(rn, 10)
      }
      redisCache.HMSet(readKey, m)
    }

    rnv = append(rnv, r)
  }

  param, err := json.Marshal(rnv)
  if err != nil {
    AppLog.With("init.func", "updateReadNumToPc").Error(err)
    return
  }

  // 通过PHP接口更新pc版redis (如果密码为空表示为线上)
  if os.Getenv("IS_DEBUG") != "open" {
    url := "http://api.cailianpress.com/v1/Redis_action/reading_num_redis_edit?debug=Cp4UFGL8QJ46sXhC26h9cf54fX24HeRD"
    params := "data=" + string(param)
    res := HttpDo("POST", url, params)
    fmt.Println("php reading_num_redis_edit res", res)
  }
}

// 定时将redis中的真实阅读数记录到数据库
func setReadNumToDatabase() {
  ids, _ := model.GetNewArticleIds(dbCacheWriter)
  if len(ids) < 1 {
    AppLog.With("init.func", "setReadNumToDatabase").Errorf("setReadNumToDatabase error \n")
    return
  }

  // 批量从redis取出数据
  arns := make([]model.ArticleReadingNum, 0)
  for _, v := range ids {
    readKey := "cls_article_read_" + v
    if redisCache.Exists(readKey).Val() {
      m := redisCache.HGetAll(readKey).Val()
      lmin, _ := strconv.Atoi(m["lmin"])
      lmax, _ := strconv.Atoi(m["lmax"])
      //lall, _ := strconv.Atoi(m["lall"])
      min, _ := strconv.Atoi(m["min"])
      max, _ := strconv.Atoi(m["max"])
      all, _ := strconv.Atoi(m["all"])

      arn := model.ArticleReadingNum{}
      arn.ArticleId, _ = strconv.Atoi(v)
      arn.ListMin = lmin
      arn.ListMax = lmax
      //arn.ListAll = lall
      arn.DetailMin = min
      arn.DetailMax = max
      arn.DetailAll = all
      arns = append(arns, arn)
    }
  }

  // 批量写入数据库
  err := model.UpdateReadNum(dbCacheWriter, arns)
  if err != nil {
    AppLog.With("method", "setReadNumToDatabase").Error(err)
  }
}

// 百度搜索链接提交
func baiduPushUrl(i uint64) {
  var subUrl = "http://data.zz.baidu.com/urls?site=www.cailianpress.com&token=FHB27BJ4PAwlotfj"
  var params = []string{}
  var res = model.GetBaiduModel(dbCacheWriter, i)
  for _, v := range res {
    var vp = "roll"
    switch v.Type {
    case -1,12:
      vp = "roll"
    case 1,8:
      vp = "theme"
    case 3:
      vp = "morning"
    case 2:
      vp = "depth"
      if v.ColumnId == 3 {
        vp = "analyze"
      }
    }
    var qurl = "https://www.cailianpress.com/"+vp+"/"+strconv.Itoa(v.Id)
    params = append(params, qurl)
  }

  var param = strings.Join(params, "\n")

  var b = HttpPost(subUrl, param)
  fmt.Println(b)

}
