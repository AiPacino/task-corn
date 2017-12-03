package main

import (
  "cailian/cron/service"
  "time"
  "strconv"
  "net/http"
  "fmt"
)

var log = service.AppLog.With("file", "main.go")

func main() {
  // redis service
  service.OpenRedis()
  defer service.CloseRedis()

  // mysql service
  service.OpenDB()
  defer service.CloseDB()

  // mysql writer service
  service.OpenWriterDB()
  defer service.CloseWriterDB()

  // init调度
  service.StartTimingTask()

  // 文章缓存调度
  service.CacheArticle()

  // 设置每隔1分钟扫描时间在周1-5的9-15点执行任务
  ticker := time.NewTicker(time.Minute * 1)
  go func() {
    for range ticker.C {
      now := time.Now()
      weekday := now.Weekday()
      hour := now.Hour()
      hm := now.Format("1504")      // 20060102150405
      hmi, _ := strconv.Atoi(hm)

      // 循环线程跑
      if weekday <= 5 && weekday >= 1 && hour >= 9 && hour <= 15 {
        // 领涨概念板块只在开盘时间取数
        if (hmi > 900 && hmi < 1131) || (hmi > 1300 && hmi < 1505) {
            log.Infof("start task GetGnRankList..................................................................." + hm)
            service.GetGnRankList()     // 获取领涨概念板块List
        }

        log.Infof("start task GetIndustryRank..................................................................." + hm)
        service.GetIndustryRank()   // 领涨行业板块

        service.GetStockRanking()   // 沪深排行榜

        log.Infof("start task StartLoop...................................................................")
        service.StartLoop()
      }
    }
  }()

  // 卡住服务
  http.HandleFunc("/_status/healthz", func(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(time.Now().Format("2006-01-02 15:04:05.000")))
  })
  log.Infof("server started... \n")
  log.Error(http.ListenAndServe(fmt.Sprintf(":%s", "8080"), nil))
}
