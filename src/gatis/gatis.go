package gatis

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"text/template"
	"path/filepath"
)

type ModelItem struct {
	Model  string
	Tag    string
	Method string
	Attrs  map[string]string
	Tplsql string
}

func (this *ModelItem) String() string {
	return fmt.Sprintf("{\n\tModel : %s \n\tMethod : %s \n\tTplsql : %s \n}\n", this.Model, this.Method, this.Tplsql)
}

func (this *ModelItem) init() {
	this.Attrs = make(map[string]string)
}

func (this *ModelItem) initMethod(model, tag string) {
	this.Attrs = make(map[string]string)
	this.Model = model
	this.Tag = tag
}

//缓存全部sql模板
var Tpls map[string]ModelItem

func FindMethod(model, method string) *ModelItem {
	for _, item := range Tpls {
		if item.Model == model && item.Method == method {
			return &item
		}
	}
	return nil
}

func pushMethod(item *ModelItem) {
	Tpls = append(Tpls, *item)
	key := item.Model + "_" + item.Method
	Tpls[key] = *item
}


var tplFuncs template.FuncMap

func init() {
	Tpls = make(map[string]ModelItem)
	tplFuncs = Tplfuncs
}

func Analysis_sql_file(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		abspath,_ := filepath.Abs(filePath)
		Log.Error("文件路径%s不存在\n", abspath)
		return
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
					modelItem.Attrs[attrName] = attrValue

					if attrName == "id" {
						modelItem.Method = attrValue
					}

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


func Render(tpl string, data interface{}) string {

	t, err := template.New("new").Funcs(tplFuncs).Parse(tpl)
	if err != nil {
		panic(err)
	}
	var bf bytes.Buffer
	err = t.Execute(&bf, data)
	if err != nil {
		panic(err)
	}
	return bf.String()
}
