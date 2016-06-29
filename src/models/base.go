package models

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"database/sql"
	"gatis/src/gatis"
	"github.com/patrickmn/go-cache"
	"github.com/astaxie/beego/logs"
	"strings"
	"fmt"
)

//创建基类与方法

type GatisModel struct {
	Id int
	GmtCreate time.Time
	GmtUpdate time.Time
}


type CacheItem struct {
	Tag string
	Tplsql string
	Attrs map[string]string
}

//cache
var presqlCache *cache.Cache
var db *sql.DB

var log *logs.BeeLogger

func init() {

	presqlCache = cache.New(-1, -1)
	var err error
	db, err = sql.Open("mysql", "root:123456@/test")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	//
	if err != nil {
		panic(err)
	}
	log = logs.NewLogger(1000)
	log.SetLogger("console", "")
	gatis.Analysis_sql_file("go/src/gatis/src/mapper/course.xml")


	//利用缓存
	//for _,v := range gatis.Tpls {
	//	item := CacheItem{v.Tag, v.Tplsql, v.Attrs}
	//	presqlCache.Add(v.Model + "_" + v.Method, item, 10 * time.Minute)
	//}

}

func getRealSql(model string, method string, data interface{}) (tag string,sql string,attrs map[string]string){

	//savekey := model + "_" + method
	//先从缓存查找
	//cacheItem,hit := presqlCache.Get(savekey)
	//
	//if hit {
	//	mc := cacheItem.(CacheItem)
	//	tag = mc.Tag
	//	sql = gatis.Render(mc.Tplsql, data)
	//	attrs = mc.Attrs
	//	return
	//}

	item := gatis.FindMethod(model, method)
	if item == nil {
		panic("no such model : " + model + ",method : " + method +" defined,please check your sql xml")
	}
	tplsql := item.Tplsql
	tag = item.Tag
	sql = gatis.Render(tplsql, data)
	attrs = item.Attrs
	//缓存tplsql
	//presqlCache.Set(savekey, CacheItem{tag, tplsql, attrs}, 10 * time.Minute)
	return
}

func Execute(model string, method string, data interface{}) interface{} {
	db.Ping()
	tag, sql, attrs := getRealSql(model, method,data)

	log.Info("执行sql:\n %s", sql)
	switch tag {
	case "insert":
		return insert(db, sql, attrs)
	case "update":
		return update(db, sql, attrs)
	case "delete":
		return del(db, sql, attrs)
	case "select":
		return query(db, method, sql, attrs)
	default:
		panic("非法标签" + tag)
		return nil
	}
}

func insert(db *sql.DB, sql string, attrs map[string]string) int64{

	//fmt.Println(db)
	rst,err := db.Exec(sql)

	if err != nil {

		log.Error("执行错误:%v", err)

		log.Error("执行sql: %s", sql)

		return -1
	}
	//fmt.Println(rst)
	if a,_ := rst.RowsAffected() ; a > 0 {
		if id,_ := rst.LastInsertId(); err == nil{
			return id
		}
	}
	return -1
}

func update(db *sql.DB, sql string, attrs map[string]string) int64{
	rst,err := db.Exec(sql)
	if err != nil {
		panic("execute sql error")
	}

	fmt.Println(rst)

	if a,_ := rst.RowsAffected() ; a > 0 {
		return 1
	}
	return -1
}


func del(db *sql.DB, sql string, attrs map[string]string) int64{
	rst,err := db.Exec(sql)
	if err != nil {
		panic("execute sql error")
	}
	if a,_ := rst.RowsAffected(); a>0 {
		return 1
	}
	return -1
}

func query(db *sql.DB, method string, sql string, attrs map[string]string) interface{}{
	rows,err := db.Query(sql)
	if err != nil {
		panic("execute sql error")
	}
	if method == "get" || attrs["one"] == "true" {
		return queryResultWithOne(rows)
	} else if method == "pageCount" || strings.HasPrefix(method, "count") {
		rt := queryResultWithOne(rows)
		//返回一个int64
		return getMapOne(rt)
	} else if method == "pageQuery" || attrs["more"] == "true" {
		return queryResultWithMore(rows)
	} else {
		rts := queryResultWithMore(rows)
		if len(rts) == 1 {
			return rts[0]
		}
		return rts
	}
}

func getMapOne(m map[string]string) interface{}{
	if len(m) != 1 {
		return m
	}
	for _,v := range m{
		return v
	}
	return m
}


func queryResultWithOne(rows *sql.Rows) (record map[string]string){
	defer rows.Close()
	//获取列
	columns,_ := rows.Columns()
	lens := len(columns)
	//scan参数
	scanArgs := make([]interface{}, lens)
	//列值
	columnValues := make([]interface{}, lens)

	//将scan参数与列值一一对应
	for i,_ := range columnValues {
		scanArgs[i] = &columnValues[i]
	}
	record = make(map[string]string)
	if rows.Next() {
		rows.Scan(scanArgs...)
		for i,columnValue := range columnValues {
			record[columns[i]] = string(columnValue.([]byte))
		}
	}
	return
}

func queryResultWithMore(rows *sql.Rows) (records []map[string]string) {
	defer rows.Close()


	//获取列
	columns,_ := rows.Columns()
	lens := len(columns)

	//scan参数
	scanArgs := make([]interface{}, lens)
	//列值
	columnValues := make([]interface{}, lens)

	//将scan参数与列值一一对应
	for i,_ := range columnValues {
		scanArgs[i] = &columnValues[i]
	}

	var initlen, totalcap = 10,10
	records = make([]map[string]string, initlen)

	n := 0

	for rows.Next() {
		rows.Scan(scanArgs...)

		record := make(map[string]string)

		for i,columnValue := range columnValues {
			record[columns[i]] = string(columnValue.([]byte))
		}

		if n < initlen {
			records[n] = record
		} else {
			records = append(records, record)
			totalcap = cap(records)
		}

		n += 1
	}
	if n < totalcap {
		records = records[:n]
	}

	return
}