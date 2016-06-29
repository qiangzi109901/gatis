package main

import (

	"fmt"
	"gatis/src/models"
	"reflect"
)

func init() {

}
func main() {
	test_insert()
}


func test_insert(){

	data := make(map[string]interface{})
	data["Id"] = "100"
	data["Name"] = "golang"
	rt := models.Execute("course","insert", data)

	////rt := models.Execute("course", "update", data)
	//
	//rt := models.Execute("course", "delete", data)
	//
	fmt.Println(rt)
	fmt.Println(reflect.TypeOf(rt))



}