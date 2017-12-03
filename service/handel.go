package service

import (
	"fmt"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "strings"
  "time"
)

var log = AppLog.With("file", "service.handel.go")

type (
  ResData struct {
    Code int                `json:"code"`
    Data map[string]ResDataMap  `json:"data"`
  }

  ResDataMap struct {
    End interface{} `json:"end"`
    Meta struct {
      Ab    string `json:"ab"`
      Board string `json:"board"`
      Name struct {
        ZhCN string `json:"zh_CN"`
      } `json:"name"`
      Py     string `json:"py"`
      Region string `json:"region"`
    } `json:"meta"`
    Mtime      string      `json:"mtime"`
    Start      string      `json:"start"`
    Status     int         `json:"status"`
    Tpl        string      `json:"tpl"`
    Type       string      `json:"type"`
    Underlying interface{} `json:"underlying"`
  }

  StockInfo struct {
    StockCode string      `json:"stockCode"`
    StockName string      `json:"stockName"`
  }

  // 领涨概念板块
  GnRankListData struct {
    Code int `json:"code"`
    Data struct {
      Foo struct {
        Mtime      int `json:"mtime"`
        RankCnGnUp []struct {
          Change float32 `json:"change"`
          Diff   float32 `json:"diff"`
          Last   float32 `json:"last"`
          Lead   []struct {
            Change float32 `json:"change"`
            Diff   float32 `json:"diff"`
            Last   float32 `json:"last"`
            Name   string  `json:"name"`
            Symbol string  `json:"symbol"`
            Tpl    string  `json:"tpl"`
            Tr     float32 `json:"tr"`
          } `json:"lead"`
          Name   string      `json:"name"`
          Symbol string      `json:"symbol"`
          Tpl    string      `json:"tpl"`
          Tr     interface{} `json:"tr"`
        } `json:"rank_cn_gn_up"`
      } `json:"foo"`
    } `json:"data"`
  }

  // 领涨行业板块
  IndustryRankData struct {
    Code int `json:"code"`
    Data struct {
      Foo struct {
        Mtime      int          `json:"mtime"`
        RankCnHyUp RankCnHyUp   `json:"rank_cn_hy_up"`
      } `json:"foo"`
    } `json:"data"`
  }

  RankCnHyUp []struct {
    Change float32 `json:"change"`
    Diff   float32 `json:"diff"`
    Last   float32 `json:"last"`
    Lead   []struct {
      Change float32     `json:"change"`
      Diff   float32     `json:"diff"`
      Last   float32     `json:"last"`
      Name   string      `json:"name"`
      Symbol string      `json:"symbol"`
      Tpl    string      `json:"tpl"`
      Tr     interface{} `json:"tr"`
    } `json:"lead"`
    Name   string      `json:"name"`
    Symbol string      `json:"symbol"`
    Tpl    string      `json:"tpl"`
    Tr     interface{} `json:"tr"`
  }

// 取得看盘宝格式json
  StockData struct {
    Code int        `json:"code"`
    Data Foo        `json:"data"`
  }

  Foo struct {
    Foo struct {
      Data  map[string]Data  `json:"data"`
      Mtime int              `json:"mtime"`
    }

    Bar struct {
      Data  map[string]BarData  `json:"data"`
      Mtime int                 `json:"mtime"`
    }
  }

  // 股票排行
  RankCnAUps struct {
    Code int `json:"code"`
    Data struct {
      Foo struct {
        Mtime int `json:"mtime"`
        RankCnAUp []struct {
          Change float64 `json:"change"`
          Diff   float64 `json:"diff"`
          Last   float64 `json:"last"`
          Name   string  `json:"name"`
          Symbol string  `json:"symbol"`
          Tpl    string  `json:"tpl"`
          Tr     float64 `json:"tr"`
        } `json:"rank_cn_a_up"`
      } `json:"foo"`
    } `json:"data"`
  }

  BarData struct {
    Name     string        `json:"name"`
    Change   float32       `json:"change"`
    Status   string        `json:"status"`
    Vol      int           `json:"vol"`
    Tso      int           `json:"tso"`
    Tr       float32       `json:"tr"`
    Mc       float32       `json:"mc"`
    Tpl      string        `json:"tpl"`
    Cmc      float32       `json:"cmc"`
    Type     string        `json:"type"`
    Time     string        `json:"time"`
    Mtime    int           `json:"mtime"`
    Amt      float32       `json:"amt"`
    Diff     float32       `json:"diff"`
    Open     float32       `json:"open"`
    Preclose float32       `json:"preclose"`
    High     float32       `json:"high"`
    Low      float32       `json:"low"`
    Last     float32       `json:"last"`
  }

  Data struct {
    Large_in   int     `json:"large_in"`
    Large_out  int     `json:"large_out"`
    Medium_in  int     `json:"medium_in"`
    Medium_out int     `json:"medium_out"`
    Little_in  int     `json:"little_in"`
    Little_out int     `json:"little_out"`
    Super_in   int     `json:"super_in"`
    Super_out  int     `json:"super_out"`
  }

  // 保存格式
  TSarray struct {
    RiseRange float32   `json:"rise_range"`
    StockName string    `json:"stock_name"`
    StockID   string    `json:"stock_id"`
    Mtime     int       `json:"mtime"`
    BigIn     int       `json:"big_in"`
    BigOut    int       `json:"big_out"`
    SmallIn   int       `json:"small_in"`
    SmallOut  int       `json:"small_out"`
    Schema    string    `json:"schema"`
    Status    string    `json:"status"`
    Last      float32   `json:"last"`
  }
)

func HttpGet(url string) string {
  defer func(){
    if x := recover(); x != nil {
      log.Errorf("http get time out, url : ", url)
    }
  }()
  timeout := time.Duration(15 * time.Second)
  client := http.Client{ Timeout: timeout }
  req, _ := http.NewRequest("GET", url, nil)
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Set("X-KPB-AUTH-ID", "cailianshe")
  req.Header.Set("X-KPB-AUTH-TOKEN", "cd1b692738c5912ce7b9c9fae4281abf09b68461")
  req.Header.Set("Connection", "close") // 完成后断开连接
  resp, err := client.Do(req)
  defer resp.Body.Close()
  if err != nil {
    log.Errorf("httpget response err \n", err)
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Errorf("httpget body err \n", err)
  }
  return string(body)
}

func HttpKPBGet(url string) string {
  defer func(){
    if x := recover(); x != nil {
      log.Errorf("http get time out, url : ", url)
    }
  }()
  timeout := time.Duration(15 * time.Second)
  client := http.Client{ Timeout: timeout }
  req, _ := http.NewRequest("GET", url, nil)
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Set("X-KPB-AUTH-ID", "cailianshe")
  req.Header.Set("X-KPB-AUTH-TOKEN", "TkFG6Gy7HiLR59ihhwQXprtlp3ki1F8i")
  req.Header.Set("Connection", "close") // 完成后断开连接
  resp, err := client.Do(req)
  defer resp.Body.Close()
  if err != nil {
    log.Errorf("httpget response err \n", err)
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Errorf("httpget body err \n", err)
  }
  return string(body)
}


// 自定义请求
func HttpDo(method string, url1 string, param string) string {
  client := &http.Client{}
  req, _ := http.NewRequest(method, url1, strings.NewReader(param))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  resp, _ := client.Do(req)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  return string(body)
}

// POST text请求
func HttpPost(url1 string, param string) string {
  client := &http.Client{}
  req, _ := http.NewRequest("POST", url1, strings.NewReader(param))
  req.Header.Set("Content-Type", "text/plain")
  resp, _ := client.Do(req)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  return string(body)
}

func StartLoop() bool {
  stocks := make(chan string)

  ui := time.Now().Unix()
  // 获取流入流出信息
  go func() {
    //后期修改为天照的kpb数据接口
    //stocks <- HttpGet("http://kanpanbao.com/private/symbol/list?mkts=sh")
    stocks <- HttpKPBGet("http://kpb.cailianpress.com/private/symbol/list?mkts=sh")
  }()
  go func() {
    //stocks <- HttpGet("http://kanpanbao.com/private/symbol/list?mkts=sz")
    stocks <- HttpKPBGet("http://kpb.cailianpress.com/private/symbol/list?mkts=sz")
  }()

  sis := make([]StockInfo, 0)
  for i := 0; i < 2; i++ {
    str := <-stocks
    handelJsonAndSetRedis(str)

    // 处理股票列表集合
    szData := ResData{}
    json.Unmarshal([]byte(str), &szData)
    for k, v := range szData.Data {
      sis = append(sis, StockInfo{StockCode:k, StockName:v.Meta.Name.ZhCN})
    }
  }

  // 写入reids股票列表
  stockinfo, err := json.Marshal(sis)
  if err == nil {
    sc := redisCache.Set("cls_stockinfo", string(stockinfo), 0)
    log.Infof("股票列表Redis Set：", sc.Val())
  }

  log.Infof("执行一共花了：", time.Now().Unix() - ui)
  return true
}

// 获取领涨概念板块List
func GetGnRankList() {
  url := "http://kanpanbao.com/kql?q={\"foo\":[\"kv.mget\",{\"keys\":[\"rank_cn_gn_up\"]}]}"
  res := HttpGet(url)
  sd := GnRankListData{}
  err := json.Unmarshal([]byte(res), &sd)
  if err != nil {
    log.Errorf("Unmarshal Rank res err \n", err)
  }

  str, err := json.Marshal(sd.Data.Foo.RankCnGnUp)
  if err != nil {
    log.Errorf("Marshal Rank sd.Data.Foo.RankCnGnUp err \n", err)
  }

  // 存入redis
  sc := redisCache.GetSet("cls_gn_rank_list", string(str))
  log.Infof("领涨概念板块Redis Set error:=", sc.Err())

  // 单独key存
  rcs := RankCnHyUp{}
  err = json.Unmarshal(str, &rcs)
  for _, v := range rcs {
    ts := TSarray{}
    ts.StockID = v.Symbol
    ts.StockName = v.Name
    ts.RiseRange = v.Change
    str, err := json.Marshal(ts)
    if err != nil {
      log.Errorf("Marshal cls_new_industry_rank_detail_ sd.Data.Foo.RankCnHyUp err \n", err)
    }
    redisCache.Set("cls_new_industry_rank_detail_" + v.Symbol, string(str), 0)
  }
}

// 领涨行业板块
func GetIndustryRank() {
  url := "http://kanpanbao.com/kql?q={\"foo\":[\"kv.mget\",{\"keys\":[\"rank_cn_hy_up\"]}]}"
  res := HttpGet(url)
  sd := IndustryRankData{}
  err := json.Unmarshal([]byte(res), &sd)
  if err != nil {
    log.Errorf("Unmarshal Rank res err \n", err)
  }

  str, err := json.Marshal(sd.Data.Foo.RankCnHyUp)
  if err != nil {
    log.Errorf("Marshal Rank sd.Data.Foo.RankCnHyUp err \n", err)
  }

  // 存入redis
  sc := redisCache.GetSet("cls_new_industry_rank", string(str))
  log.Infof("领涨行业板块Redis Set error:=", sc.Err())

  // 单独key存
  rcs := RankCnHyUp{}
  err = json.Unmarshal(str, &rcs)
  for _, v := range rcs {
    ts := TSarray{}
    ts.StockID = v.Symbol
    ts.StockName = v.Name
    ts.RiseRange = v.Change
    str, err := json.Marshal(ts)
    if err != nil {
      log.Errorf("Marshal cls_new_industry_rank_detail_ sd.Data.Foo.RankCnHyUp err \n", err)
    }
    redisCache.Set("cls_new_industry_rank_detail_" + v.Symbol, string(str), 0)
  }

}

// 获取沪深排行榜 cls_new_rank_cn_a_up
func GetStockRanking() {
  url := `http://kanpanbao.com/kql?q={"foo":["kv.mget",{"keys":["rank_cn_a_up"]}]}`
  res := HttpGet(url)
  sd := RankCnAUps{}
  err := json.Unmarshal([]byte(res), &sd)
  if err != nil {
    log.Errorf("Unmarshal RankCnAUps res err \n", err)
  }

  str, err := json.Marshal(sd.Data.Foo.RankCnAUp)
  if err != nil {
    log.Errorf("Marshal Rank sd.Data.Foo.RankCnAUp err \n", err)
  }

  // 存入redis
  sc := redisCache.GetSet("cls_new_rank_cn_a_up", string(str))
  log.Infof("沪深排行榜 Redis Set error:=", sc.Err())
}

// 获取所有板块信息
func getBlockList() {
  url := `http://kanpanbao.com/private/block/list`
  res := HttpGet(url)
  redisCache.Set("cls_block_list", []byte(res), 0)
}

func handelJsonAndSetRedis(sza string) {
  if sza == "" {
    return
  }

  arr := []string{}
  szData := ResData{}
  json.Unmarshal([]byte(sza), &szData)
  for m := range szData.Data {
    arr = append(arr, m)
  }

  log.Infof("arr length ", len(arr))

  // json格式化后每100次存放在二维数组中进行迭代
  params := make(chan []string)

  //cInsertSize := make(chan int)

  go func() {
    SIZE := 40
    var one = make([]string, 0, SIZE)
    for i:=0; i<len(arr); i++ {
      one = append(one, arr[i])
      if i > 0 && i%SIZE == 0 {
        params <- one
        one = make([]string, 0, SIZE)
      }
    }
    if len(one) > 0 {
      params <- one
    }
    close(params)
  }()

  for i := 1; i <= 10; i++ {
    go func() {
      for v := range params {
        url := "http://kanpanbao.com/kql?q={\"foo\":[\"fundflow.origin\",{\"symbols\":%v}]}"      // 流入流出
        str := fmt.Sprintf("%q", v)
        str = strings.Replace(str, " ", ",", -1)
        url = fmt.Sprintf(url, str)

        gurl := "http://kanpanbao.com/kql?q={\"bar\":[\"quote\",{\"symbols\":%v,\"mtime\":1}]}"   // 获取个股或板块详情
        gurl = fmt.Sprintf(gurl, str)

        res := HttpGet(url)
        sd := StockData{}
        err := json.Unmarshal([]byte(res), &sd)
        if err != nil {
          log.Errorf("Unmarshal res err", err)
          continue
        }

        gres := HttpGet(gurl)
        qus := StockData{}
        err = json.Unmarshal([]byte(gres), &qus)
        if err != nil {
          log.Errorf("Unmarshal gres err", err)
          continue
        }
        // 封装成特定json对象后循环set redis
        for k, v := range sd.Data.Foo.Data {
          b := qus.FindInfo(k)
          ts := TSarray{}
          ts.StockID = k
          ts.StockName = b.Name
          ts.RiseRange = b.Change
          ts.Status = b.Status
          ts.Mtime = sd.Data.Foo.Mtime
          ts.BigIn = v.Large_in + v.Super_in
          ts.BigOut = v.Large_out + v.Super_out
          ts.SmallIn = v.Little_in
          ts.SmallOut = v.Little_out
          ts.Last = b.Last

          ch, _ := json.Marshal(ts)
          redisCache.Set("cls_new_stock_detail_" + k, ch, 0)
          //sc := redisCache.Set("cls_new_stock_detail_" + k, ch, 0)
          //log.Infof(sc.Val(), ts)
        }

      //  res := make(chan string)
      //  gres := make(chan string)
      //  err := make(chan error)
      //  go func() {
      //    res <- HttpGet(url)
      //  }()
      //
      //  go func() {
      //    gres <- HttpGet(gurl)
      //  }()
      //
      //  select {
      //  case resa := <-res:
      //    select {
      //    case gresa := <-gres:
      //      sd := StockData{}
      //      err := json.Unmarshal([]byte(resa), &sd)
      //      if err != nil {
      //        log.Errorf("Unmarshal res err", err)
      //        cInsertSize <- len(v)
      //        continue
      //      }
      //
      //      qus := StockData{}
      //      err = json.Unmarshal([]byte(gresa), &qus)
      //      if err != nil {
      //        log.Errorf("Unmarshal gres err", err)
      //        cInsertSize <- len(v)
      //        continue
      //      }
      //
      //      // 封装成特定json对象后循环set redis
      //      for k, v := range sd.Data.Foo.Data {
      //        b := qus.FindInfo(k)
      //        ts := TSarray{}
      //        ts.StockID = k
      //        ts.StockName = b.Name
      //        ts.RiseRange = b.Change
      //        ts.Status = b.Status
      //        ts.Mtime = sd.Data.Foo.Mtime
      //        ts.BigIn = v.Large_in + v.Super_in
      //        ts.BigOut = v.Large_out + v.Super_out
      //        ts.SmallIn = v.Little_in
      //        ts.SmallOut = v.Little_out
      //
      //        ch, _ := json.Marshal(ts)
      //        redisCache.Set("cls_new_stock_detail_" + k, ch, 0)
      //        //sc := redisCache.Set("cls_new_stock_detail_" + k, ch, 0)
      //        // log.Infof(sc.Val(), ts)
      //      }
      //      log.Infof("共处理数据：", len(sd.Data.Foo.Data))
      //      cInsertSize <- len(v)
      //    case <- err:
      //      log.Errorf("case", err)
      //      cInsertSize <- len(v)
      //      continue
      //    }
      //  case <- err:
      //    log.Errorf("case", err)
      //    cInsertSize <- len(v)
      //    continue
      //  }
      }
    }()
  }

  //sum := 0
  //for true {
  //  sum += <-cInsertSize
  //  log.Infof("==> %d/%d\n", sum, len(arr))
  //  if sum == len(arr) {
  //    return
  //  }
  //}
}

func (sd StockData) FindInfo(stock string) BarData {
  for k, v := range sd.Data.Bar.Data {
    if stock == k {
      return v
    }
  }
  return BarData{}
}

