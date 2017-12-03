package model

import (
	sq "github.com/Masterminds/squirrel"
  "github.com/wangyibin/alilog"
  "github.com/Masterminds/structable"
  "fmt"
  "database/sql"
)

type (
   // 用于检查文章是否有更新
  ArticleLatestCtime struct {
    Type     int
    ColumnID int
    Ctime    int
  }

  // 阅读数表
  ArticleReadingNum struct {
    ArticleId int   `stbl:"article_id,PRIMARY_KEY,AUTO_INCREMENT"`
    ListMin   int   `stbl:"list_min"`
    ListMax   int   `stbl:"list_max"`
    ListAll   int   `stbl:"list_all"`
    DetailMin int   `stbl:"detail_min"`
    DetailMax int   `stbl:"detail_max"`
    DetailAll int   `stbl:"detail_all"`
    rec       structable.Recorder
  }

  ArticleV1Filter struct {
    Type         string         //文章类型
    Column_id    int            //内参类型
    Refresh_type int            // 0: 下拉刷新 1: 上拉加载更多
    Rn           PageLimit      // 分页所需字段: Limit
    Last_time    uint64         // 分页所需字段：Start
    IsTop        int            // 题材焦点
  }

  BaiduModel struct {
    Id           int   `stbl:"id"`
    Type         int   `stbl:"type"`
    ColumnId     int   `stbl:"column_id"`
  }

  ArticleFilter struct {
    Id                  int
    CType 		          []int
    IsTop               int
    IsFocus             int
    Column_id           int
    Ctime               int64
    Rn                  PageLimit
    Refresh_Type        int
    Recommend           int
    In_roll             int
    Explain             int
    Jpush               int
    Content             string
    Topics              string
    Uid                 int
  }

  V1_article_type struct {
    Id                  int                   `stbl:"id,PRIMARY_KEY,AUTO_INCREMENT"`
    User_id             int                   `stbl:"user_id"`
    Title               string                `stbl:"title"`
    Short_title         string                `stbl:"short_title"`
    Level               string                `stbl:"level"`
    Type                int                   `stbl:"type"`
    Ctime               int                   `stbl:"ctime"`
    Sort_score          float64               `stbl:"sort_score"`
    Brief               string                `stbl:"brief"`
    Column_info         string                `stbl:"column_info"`
    Content             sql.NullString        `stbl:"content"`
    Reading_num         int                   `stbl:"reading_num"`
    Recommend           int                   `stbl:"recommend"`
    Bold                int                   `stbl:"bold"`
    Is_top              int                   `stbl:"is_top"`
    In_roll             int                   `stbl:"in_roll"`
    Author              string                `stbl:"author"`
    Author_extends      string                `stbl:"author_extends"`
    Img                 string                `stbl:"img"`
    Neican_img          string                `stbl:"neican_img"`
    Jpush               int                   `stbl:"jpush"`
    Jieshuo             sql.NullString        `stbl:"jieshuo"`
    Weibo               string                `stbl:"weibo"`
    Sina_id             string                `stbl:"sina_id"`
    Comment_num         int                   `stbl:"comment_num"`
    Deny_comment        int                   `stbl:"deny_comment"`
    Abstract            string                `stbl:"abstract"`
    Modified_time       int                   `stbl:"modified_time"`
    Status              int                   `stbl:"status"`
    Stocks              sql.NullString        `stbl:"stocks"`
    Depth_extends       sql.NullString        `stbl:"depth_extends"`
    Column_id           int                   `stbl:"column_id"`
    Topics              string                `stbl:"topics"`
    Notes               sql.NullString        `stbl:"notes"`
    Remark              string                `stbl:"remark"`
    Audio_url           string                `stbl:"audio_url"`
    Confirmed           int                   `stbl:"confirmed"`
    Push_content        sql.NullString        `stbl:"push_content"`
    IsExplain           sql.NullInt64         `stbl:"is_explain"`
    rec                 structable.Recorder
  }

  // 专题列表
  SpecialTopicList struct {
    Id	                int	                  `stbl:"id"`	                                                          // 信息ID
    Type	              int	                  `stbl:"type"`	                                                        // 类型
    Img	                string	              `stbl:"img"`	                                                        // 类型
    Title	              string	              `stbl:"title"`	              	                                      // 类型
    Content	            string	              `stbl:"content"`	              	                                    // 帖子内容
    DenyComment	        int	                  `stbl:"deny_comment"`	                                                // 是否禁止评论，1禁止，0允许
    Ctime	              int	                  `stbl:"ctime"`	                                                      // 信息创建时间
    Modified_time       int                   `stbl:"modified_time"`                                                // 更新时间
    SortScore	          float64   	          `stbl:"sort_score"`	                                                  // 信息队列排序时间
    CommentNum	        int	                  `stbl:"comment_num"`	                                                // 评论数
    Status	            int	                  `stbl:"status"`	                                                      // 文章状态，1正常，2隐藏，3正常显示
    IsTop  	            int	                  `stbl:"is_top"`	                                                      // 是否热点专题
    Column_id	          int                   `stbl:"column_id"`	                                                  // column_id
    Collection          int                   `stbl:"collection"`	                                                  // collection
    CollectionNum       int                   `stbl:"collection_num"`	                                              // collection_num
    Reason              string                `stbl:"reason"`
    Stocks              sql.NullString        `stbl:"stockarray"`
    rec                 structable.Recorder
  }
)

func NewArticle(db sq.DBProxyBeginner) *V1_article_type {
  u := new(V1_article_type)
  u.rec = structable.New(db, "mysql").Bind("lian_v1_article", u)
  return u
}

// 自定义类型限制为最大值不能大于20
type PageLimit uint64
func (p PageLimit) MustLt20() uint64 {
  if p == 0 {
    return 10
  } else if p > 20 {
    return 20
  }
  return uint64(p)
}

func (p PageLimit) MustLt300() uint64 {
  if p == 0 {
    return 10
  } else if p > 300 {
    return 300
  }
  return uint64(p)
}

var articleLog = alilog.New("cailian-cron", "model")

// 检查文章是否有更新
func GetNewArticle(db sq.DBProxyBeginner, ctime int64) []ArticleLatestCtime {
  log := articleLog.With("method", "GetNewArticle")
  query := sq.Select("type, column_id, max(ctime)").From("lian_v1_article")
  query = query.Where(sq.Eq{"status":1}).Where("ctime > ?", ctime).GroupBy("type, column_id")
  log.Infof(query.ToSql())
  rows, err := query.RunWith(db).Query()
  defer rows.Close()
  if err != nil {
    log.Debugf("GetNewArticle查询错误")
    return nil
  }

  var result = make([]ArticleLatestCtime, 0)
  for rows.Next() {
    u := ArticleLatestCtime{}
    if err := rows.Scan(&u.Type, &u.ColumnID, &u.Ctime); err != nil {
      log.Debugf(err.Error())
      return nil
    }
    result = append(result, u)
  }
  return result
}

// 获取最新（-1, 1, 2, 3 type 各取 20条）文章id
func GetNewArticleIds(db sq.DBProxyBeginner) ([]string, []int64) {
  log := articleLog.With("method", "GetNewArticleIds")
  sqlstr := "select a.id, a.ctime from (" +
    "(select a.id, a.ctime, a.type from lian_v1_article a where a.`status` = 1 and a.`type` = -1 order by a.id desc limit 20)" +
    " UNION ALL " +
    "(select a.id, a.ctime, a.type from lian_v1_article a where a.`status` = 1 and a.`type` = 12 order by a.id desc limit 20)" +
    " UNION ALL " +
    "(select a.id, a.ctime, a.type from lian_v1_article a where a.`status` = 1 and a.`type` = 15 order by a.id desc limit 20)" +
    " UNION ALL " +
    "(select a.id, a.ctime, a.type from lian_v1_article a where a.`status` = 1 and a.`type` = 16 order by a.id desc limit 20)" +
    " UNION ALL " +
    "(select a.id, a.ctime, a.type from lian_v1_article a where a.`status` = 1 and a.`type` = 1 order by a.id desc limit 20)" +
    " UNION ALL " +
    "(select a.id, a.ctime, a.type from lian_v1_article a where a.`status` = 1 and a.`type` = 2 order by a.id desc limit 20)" +
    " UNION ALL " +
    "(select a.id, a.ctime, a.type from lian_v1_article a where a.`status` = 1 and a.`type` = 8 order by a.id desc limit 20)" +
    " UNION ALL " +
    "(select a.id, a.ctime, a.type from lian_v1_article a where a.`status` = 1 and a.`type` = 3  order by a.id desc limit 20)" +
    ") a GROUP BY a.id"

  var result = make([]string, 0)
  var ctime = make([]int64, 0)
  rows, err := sq.QueryWith(db, sq.Expr(sqlstr))
  defer rows.Close()
  if err != nil {
    log.Debugf("GetNewArticleIds查询错误", err)
    return result, ctime
  }

  for rows.Next() {
    var u string
    var ct int64
    if err := rows.Scan(&u, &ct); err != nil {
      log.Debugf(err.Error())
    }
    result = append(result, u)
    ctime = append(ctime, ct)
  }
  return result, ctime
}

//INSERT INTO lian_article_reading_num (article_id,list_min, list_max, detail_min, detail_max, detail_all)
//VALUES (62539,3, 4, 5, 6, 7) ON DUPLICATE KEY UPDATE list_min=10, list_max=20, detail_min=30, detail_max=40, detail_all=50
func UpdateReadNum(db sq.DBProxyBeginner, ns []ArticleReadingNum) error {
  sb := sq.Insert("lian_article_reading_num").Columns("article_id, list_min, list_max, detail_min, detail_max, detail_all")
  if len(ns) == 0 {
    return nil
  }
  for _, n := range ns {
    sb = sb.Values(n.ArticleId, n.ListMin, n.ListMax, n.DetailMin, n.DetailMax, n.DetailAll)
  }
  sb = sb.Suffix("ON DUPLICATE KEY UPDATE list_min=VALUES(list_min), list_max=VALUES(list_max), detail_min=VALUES(detail_min), detail_max=VALUES(detail_max), detail_all=VALUES(detail_all)")
  sql, parm, _ := sb.ToSql()
  fmt.Println("", sql, " \n ", parm)
  _, err := sb.RunWith(db).Exec()
  return err
}

// 百度搜索链接提交
func GetBaiduModel(db sq.DBProxyBeginner, i uint64) []BaiduModel {
  var log = articleLog.With("func", "baiduPushUrl")
  var query = sq.Select("id, type, column_id").From("lian_v1_article").Where("status=1").OrderBy("sort_score desc").Offset(i).Limit(100)
  log.Infof(query.ToSql())
  var bms = make([]BaiduModel, 0)
  var rows, err = query.RunWith(db).Query()
  if err != nil {
    log.Errorf("GetBaiduModel Query error", err)
    return bms
  }

  for rows.Next() {
    var bm = BaiduModel{}
    rows.Scan(&bm.Id, &bm.Type, &bm.ColumnId)
    bms = append(bms, bm)
  }
  return bms
}


// 通用查询写法，动态条件可以添加
func FindArticlesTime(db sq.DBProxyBeginner, filter ArticleFilter) ([]*V1_article_type, error) {
  log := articleLog.With("method", "FindArticlesTime")
  query := sq.Select(NewArticle(db).rec.Columns(true)...).From("lian_v1_article")
  query = queryArticleParams(query, filter) // 通用字段匹配条件

  // 默认0。 0: 下拉刷新（重新查询） 1:上拉加载更多
  if filter.Ctime != 0 && filter.Refresh_Type == 1 {
    query = query.Where(sq.Lt{"sort_score":filter.Ctime})
  } else {
    if len(filter.CType) > 0 && filter.CType[0] == 2 {
      query = query.OrderBy("is_top desc")  // CType=2深度头条设置is_top排序
    }
  }

  query = query.OrderBy("sort_score desc").Limit(filter.Rn.MustLt300())
  rows, err := query.RunWith(db).Query()
  defer rows.Close()
  if err != nil {
    log.Debugf("ArticleFilter查询错误", err)
    return nil, err
  }

  var result = make([]*V1_article_type, 0)
  for rows.Next() {
    u := NewArticle(db)
    if err := rows.Scan(u.rec.FieldReferences(true)...); err != nil {
      log.Debugf(err.Error())
      return nil, err
    }
    result = append(result, u)
  }
  return result, nil
}


// 专题列表
func GetSpecialTopicList(dbProxy sq.DBProxyBeginner, f ArticleFilter) ([]SpecialTopicList, error) {
  if len(f.CType) == 0 {
    return nil, nil
  }

  var articles = make([]SpecialTopicList, 0)
  query := sq.Select("a.id, a.type, a.img, a.title, a.content, a.deny_comment, a.ctime, a.modified_time, a.sort_score, a.comment_num, a.collection_num, a.status, a.is_top, a.reason, a.stocks").
    From("lian_special_topic a").
    Where(sq.Eq{"a.status":1, "a.type":f.CType})
  if f.Id > 0 {
    query = query.Where("a.id = ?", f.Id)
  }
  if f.IsTop == 1 {
    query = query.Where("a.is_top >= ?", f.IsTop)
  }
  if f.Refresh_Type > 0 {
    if f.Ctime > 0 {
      query = query.Where("a.sort_score < ?", f.Ctime)
    }
  }
  if f.IsTop == 1 {
    query = query.OrderBy("a.is_top desc, a.sort_score desc ")
  } else {
    query = query.OrderBy("a.sort_score desc ")
  }
  query = query.Limit(f.Rn.MustLt20()) // 按时间分页写法，默认10
  rows, err := query.RunWith(dbProxy).Query()
  defer rows.Close()
  if err != nil {
    return articles, err
  }

  for rows.Next() {
    c := SpecialTopicList{}
    err := rows.Scan(&c.Id, &c.Type, &c.Img, &c.Title, &c.Content, &c.DenyComment, &c.Ctime, &c.Modified_time, &c.SortScore, &c.CommentNum, &c.CollectionNum, &c.Status, &c.IsTop, &c.Reason, &c.Stocks)
    if err != nil {
      return nil, err
    }
    articles = append(articles, c)
  }

  return articles, err
}


// FindArticlesTime通用字段匹配条件
func queryArticleParams(query sq.SelectBuilder, filter ArticleFilter) sq.SelectBuilder {
  query = query.Where(sq.Eq{"status":1})
  if filter.Id > 0 {
    query = query.Where(sq.Eq{"id": filter.Id})
  }
  if len(filter.CType) > 0 {
    query = query.Where(sq.Eq{"type": filter.CType})
  }
  if filter.IsTop == 1 {
    query = query.Where(sq.Eq{"is_top": filter.IsTop})
  }
  if filter.IsFocus == 1 {
    query = query.Where(sq.Eq{"is_focus": filter.IsFocus})
  }
  if filter.Column_id != 0 {
    query = query.Where("column_id = ?", filter.Column_id)
  }
  if filter.Recommend != 0 {
    query = query.Where("recommend = ?", filter.Recommend)
  }
  if filter.Explain == 1 {
    query = query.Where("is_explain >= 0")
  }
  if filter.In_roll != 0 {
    query = query.Where("in_roll = ?", filter.In_roll)
  }
  if filter.Jpush != 0 {
    query = query.Where("jpush = ?", filter.Jpush)
  }
  if filter.Content != "" {
    query = query.Where("content like ?", "%"+filter.Content+"%")
  }
  if filter.Topics != "" {
    query = query.Where("topics like ?", "%"+filter.Topics+"%")
  }
  return query
}

func GetArticles(db sq.DBProxyBeginner, f ArticleFilter) ([]V1_article_type, error) {
  var a = NewArticle(db)
  var articles = make([]V1_article_type, 0)
  var query = sq.Select(a.rec.Columns(true)...).From("lian_v1_article")
  if f.Id > 0 {
    query = query.Where("id = ?", f.Id)
  }
  if len(f.CType) > 0 {
    query = query.Where(sq.Eq{"type":f.CType})
  }
  if f.In_roll == 1 {
    query = query.Where("in_roll = ?", f.In_roll)
  }
  if f.Ctime > 0 {
    query = query.Where("sort_score > ?", f.Ctime)
  }
  query = query.OrderBy("sort_score desc ").Limit(f.Rn.MustLt300())
  rows, err := query.RunWith(db).Query()
  defer rows.Close()
  if err != nil {
    return nil, err
  }

  for rows.Next() {
    if err := rows.Scan(a.rec.FieldReferences(true)...); err != nil {
      return nil, err
    }
    articles = append(articles, *a)
  }
  return articles, nil
}
