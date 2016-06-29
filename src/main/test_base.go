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
	data["Id"] = "200"
	data["Name"] = "golang"
	rt := gatis.Execute("course","insert", data)

	////rt := models.Execute("course", "update", data)
	//
	//rt := models.Execute("course", "delete", data)
	//
	fmt.Println(rt)
	fmt.Println(reflect.TypeOf(rt))



}