package main

import (
	"fmt"
)


var k *int

func init() {

	a := 10
	k = &a
}

func main() {

	fmt.Println(k)
	fmt.Println(*k)


}



