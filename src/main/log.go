package main


import (
	"github.com/astaxie/beego/logs"
	"fmt"
)

func main() {

	log := logs.NewLogger(2222)

	log.SetLogger("console", "")

	log.Debug("hello")
	log.Error("hello")

	i := 100

	fmt.Println(i)

	for i:=0;i<10;i++ {
		fmt.Println(i)
	}

	fmt.Println(i)

}
