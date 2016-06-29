package main

import (
	"fmt"
	"reflect"
	"gatis/src/gatis"
)

func init() {

}
func main() {
	test_insert()
}

func test_insert(){

	data := make(map[string]interface{})
	data["Id"] = "1"
	data["Name"] = "java"
	//rt := gatis.Execute("course","insert", data)

	//rt := gatis.Execute("course", "update", data)
	//
	//rt := gatis.Execute("course", "delete", data)

	//rt := gatis.Execute("course", "get", data)

	rt := gatis.Execute("course", "queryAll", data)
	fmt.Println(rt)
	fmt.Println(reflect.TypeOf(rt))



}