package gatis

import (
	"fmt"
)

func Ps(args ...interface{}) {
	for _, v := range args {
		fmt.Print(v, " ")
	}
	fmt.Println("")
}
