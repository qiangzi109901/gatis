package main
import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)

var db *sql.DB
func init() {

	ks,err := sql.Open("mysql", "root:123456@/test")

	if err != nil {
		panic(err)
	}

	db = ks
}

func main() {
	doadd()
}

func doadd(){
	msql := `insert ignore into
		t_course (
			id,
			name
		)
		value (
			'99',
		'hello'
	)`

	defer db.Close()
	db.Ping()

	rst,err := db.Exec(msql)
	if err != nil {
		panic(err)
	}
	fmt.Println(rst.RowsAffected())
}


func test01(){
	msql := `insert ignore into
		t_course (
			id,
			name
		)
		value (
			'9',
		'hello'
	)`

	db,err := sql.Open("mysql", "root:123456@/test")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.Ping()

	rst,err := db.Exec(msql)
	if err != nil {
		panic(err)
	}
	fmt.Println(rst.RowsAffected())
}



func test02(){
	msql := `insert into
		t_course (
			id,
			name
		)
		value (
			'77',
		'hello'
	)`

	db,err := sql.Open("mysql", "root:123456@/test")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.Ping()

	stmt,err := db.Prepare(msql)
	if err != nil {
		if err != nil {
			panic(err)
		}
	}
	kk,err := stmt.Exec()
	throw(err)
	fmt.Println(kk)
}

func throw(e error) {
	if e != nil {
		panic(e)
	}
}