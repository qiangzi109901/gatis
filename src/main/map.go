package main

import "fmt"

var t map[string]string


type User struct {
	Id int
	Name string
}


func init() {
}

func main() {
	fmt.Println(t == nil)
	t = make(map[string]string)
	fmt.Println(t == nil)
}
