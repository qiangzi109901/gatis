package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

type ModelItem struct {
	Model  string
	Method string
	attrs  map[string]string
	Tplsql string
}

func (this *ModelItem) String() string {
	return fmt.Sprintf("{\n\tModel : %s \n\tMethod : %s \n\tTplsql : %s \n}\n", this.Model, this.Method, this.Tplsql)
}

func (this *ModelItem) init() {
	this.attrs = make(map[string]string)
}

func (this *ModelItem) initMethod(model, method string) {
	this.attrs = make(map[string]string)
	this.Model = model
	this.Method = method
}

//缓存全部sql模板
var Tpls []ModelItem

func pushMethod(item *ModelItem) {
	Tpls = append(Tpls, *item)
}


func Ps(args ...interface{}) {
	for _, v := range args {
		fmt.Print(v, " ")
	}
	fmt.Println("")
}

func Analysis_sql_file(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		abspath,_ := filepath.Abs(filePath)
		Ps(abspath)
		panic("文件路径不存在")
	}
	defer file.Close()

	xmldoc := xml.NewDecoder(file)
	modelName := ""
	isEnterElement := false
	modelItem := new(ModelItem)
	modelItem.init()
	for {
		t, err := xmldoc.Token()
		if err != nil {
			break
		}
		isClose := false
		switch token := t.(type) {

		case xml.StartElement:
			name := token.Name.Local
			tname := strings.TrimSpace(name)

			//根节点
			if tname == "sql" {
				for _, attr := range token.Attr {
					attrName := attr.Name.Local
					if attrName == "id" {
						modelName = attr.Value
					}
				}
			} else if tname != "" {
				modelItem.initMethod(modelName, name)

				for _, attr := range token.Attr {
					attrName := attr.Name.Local
					attrValue := attr.Value
					modelItem.attrs[attrName] = attrValue
				}
				isEnterElement = true

			}

		case xml.CharData:
			c := string([]byte(token))
			tc := strings.TrimSpace(c)

			if tc != "" && isEnterElement {
				modelItem.Tplsql = c
				pushMethod(modelItem)
				modelItem.init()
			}

		case xml.EndElement:
			name := token.Name.Local
			tname := strings.TrimSpace(name)

			if tname == "sql" {
				isClose = true
			} else if tname != "" {
				isEnterElement = false
			}

		default:

		}

		if isClose {
			break
		}
	}
}

func main() {

	Analysis_sql_file("go/src/gatis/src/mapper/user.xml")

	for _,v := range Tpls {
		Ps(v.String())
	}
}