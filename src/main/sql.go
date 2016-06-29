package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func main() {

	db, err := sql.Open("mysql", "root:123456@/test")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}


	insert(db)

}




func insert(db *sql.DB) {
	stmt, err := db.Prepare("insert into t_course(id,name) value(?,?)")
	defer stmt.Close()

	if err != nil {
		panic(err)
	}

	stmt.Exec(1, "java")
	stmt.Exec(2, "golang")


	//queryWithCourse(db)
	queryAndAssignMap(db)
	queryAndAssignCourse(db)
}


func query(db *sql.DB) {

	rows, err := db.Query("select * from t_course")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var id int
	var name string

	for rows.Next() {
		err = rows.Scan(&id, &name)

		if err != nil {
			panic(err)
		}

		fmt.Println(id, name)
	}

}

type Course struct {
	Id int
	Name string
}

func (this Course) String() string {
	return "Id:" + strconv.Itoa(this.Id) + ",Name:" + this.Name
}

func (this Course) Hello(name string) {
	fmt.Println("hello ", name, ", i'm ", this.Name)
}












func queryAndAssignCourse(db *sql.DB) {
	rows,err := db.Query("select * from t_course")
	if err != nil {
		panic(err)
	}
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
	courses := make([]Course, initlen)
	n := 0
	for rows.Next() {
		rows.Scan(scanArgs...)
		course := new(Course)
		fmt.Println(reflect.ValueOf(course).Kind())
		cv := reflect.ValueOf(course).Elem()
		for i,columnValue := range columnValues {
			cf := cv.FieldByName(getFiledName(columns[i]))
			mv := string(columnValue.([]byte))
			fmt.Println(cf.Kind(),cf.Kind() == reflect.Int,"**********")
			switch cf.Kind() {
			case reflect.String:
				cf.SetString(mv)
			case reflect.Int8:
				cf.SetInt(int64(parseInt(mv, 8).(int8)))
			case reflect.Int16:
				cf.SetInt(int64(parseInt(mv, 16).(int16)))
			case reflect.Int32:
				cf.SetInt(int64(parseInt(mv, 32).(int32)))
			case reflect.Int:
				cf.SetInt(int64(parseInt(mv, 0).(int)))
			case reflect.Int64:
				cf.SetInt(int64(parseInt(mv, 64).(int64)))
			case reflect.Bool:
				b,err := strconv.ParseBool(mv)
				if err != nil{
					cf.SetBool(b)
				}
			default:
				fmt.Println("other kind:", cf.Kind())
			}
		}
		if n < initlen {
			courses[n] = *course
		} else {
			courses = append(courses, *course)
			totalcap = cap(courses)
		}
		n += 1
	}
	if n < totalcap {
		courses = courses[:n]
	}
	fmt.Println(courses)
}

func getFiledName(columnName string) string {
	cs := strings.Split(columnName, "_")
	for i,v := range cs {
		cs[i] = strings.ToUpper(v[0:1]) + strings.ToLower(v[1:])
	}
	return strings.Join(cs,"")

}

func parseInt(str string,size int) interface{}{

	switch size {
	case 0:
		d, err := strconv.Atoi(str)
		if err == nil {
			return d
		}
	case 8:
		d, err := strconv.ParseInt(str, 10, 8)
		if err == nil {
			return d
		}
	case 16:
		d, err := strconv.ParseInt(str, 10, 16)
		if err == nil {
			return d
		}
	case 32:
		d, err := strconv.ParseInt(str, 10, 32)
		if err == nil {
			return d
		}
	case 64:
		d, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			return d
		}
	default:
		return -1
	}
	return -1
}


func queryAndAssignMap(db *sql.DB) {
	rows,err := db.Query("select * from t_course")

	if err != nil {
		panic(err)
	}

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
	records := make([]map[string]string, initlen)

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

	fmt.Println(records)

}

func queryWithCourse(db *sql.DB) {

	rows, err := db.Query("select * from t_course")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var courses []Course

	//获取列
	columns,err := rows.Columns()
	lens := len(columns)
	scanArgs := make([]interface{}, lens)
	values := make([]interface{}, lens)

	//这样就将值映射到切片values中
	for i,_ := range values {
		scanArgs[i] = &values[i]
	}

	var initlen ,totalcap = 10,10
	records := make([]map[string]string,initlen)

	n := 0
	for rows.Next() {
		rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i,col := range values {
			if col != nil {

				record[columns[i]] = string(col.([]byte))
			}
		}
		if n < initlen {
			records[n] = record
			//fmt.Println("assign the old address")

		} else {
			records = append(records, record)
			if n < totalcap {
				//fmt.Println("assign the old address")
			} else {
				//fmt.Println("allocate new address")
				totalcap = cap(records)
			}
		}
		n += 1
	}

	if n < totalcap {
		records = records[:n]
	}

	fmt.Println(records)
	fmt.Println(courses)
	fmt.Println(columns)

	testReflect()
}

func testReflect(){

	course := Course{1,"java"}

	v := reflect.ValueOf(&course).Elem()

	mtype := v.Type()
	len := v.NumField()

	fmt.Println(mtype)
	fmt.Println(len)
	fmt.Println(v)


	fmt.Println(v.Kind())
	fmt.Println(v)


	for i:=0;i<len;i++ {
		fmt.Println(mtype.Field(i).Name, v.Field(i).Type(), v.Field(i))
		fmt.Println(v.Field(i).Type())
		if v.Field(i).CanInterface() {
			fmt.Println(v.Field(i).Interface())
		}

		kind := v.Field(i).Kind()
		if kind == reflect.String {
			v.Field(i).SetString("hello")
		}

		if kind == reflect.Int {
			v.Field(i).SetInt(222)
		}


	}

	method := v.MethodByName("String")

	ks := method.Call(nil)

	fmt.Println(ks)


	method2 := v.MethodByName("Hello")

	args := []reflect.Value{reflect.ValueOf("bob")}
	method2.Call(args)



	m := make(map[string]interface{})

	m["Id"] = 100
	m["Name"] = "john"

	c := &m
	map2struct(m, c)
	fmt.Println(c)




}


/**
	构思:  通过反射去调用结构属性的setter方法给属性赋值
 */

func map2struct(m map[string]interface{}, r interface{}) {

	if reflect.TypeOf(r).Kind() != reflect.Ptr {
		//panic("r must be Struct Pointer")
		fmt.Println("r must be Struct Pointer")
		return
	}
	rv := reflect.ValueOf(r).Elem()

	if rv.Type().Kind() != reflect.Struct {
		//panic("r must be Struct Pointer")
		fmt.Println("r must be Struct Pointer")
		return
	}

	//遍历map
	for k,v := range m {
		//k对应了结构体r中的属性
		rf := rv.FieldByName(k)
		fmt.Println(rf.Kind())

		//rf.Set(reflect.ValueOf(v))

		//fmt.Println(rf.CanSet())
		//fmt.Println(rf.CanInterface())
		switch rf.Kind() {
		case reflect.String:
			rf.SetString(v.(string))
		case reflect.Int8:
			rf.SetInt(int64(v.(int8)))
		case reflect.Int16:
			rf.SetInt(int64(v.(int16)))
		case reflect.Int32:
			rf.SetInt(int64(v.(int32)))
		case reflect.Int:
			rf.SetInt(int64(v.(int)))
		case reflect.Int64:
			rf.SetInt(v.(int64))
		case reflect.Bool:
			rf.SetBool(v.(bool))
		default:
			fmt.Println("other kind:", rf.Kind())
		}
		fmt.Println(v)
	}

}