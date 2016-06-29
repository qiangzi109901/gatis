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

var Tpls map[string]*ModelItem

func init() {
	initTpls()
}

func initTpls(){
	if Tpls == nil {
		Tpls = make(map[string]*ModelItem)
	}
}

func (this *ModelItem) String() string {
	return fmt.Sprintf("{\n\tModel : %s \n\tTag : %s \n\tMethod : %s \n\tTplsql : %s \n}\n", this.Model, this.Tag, this.Method, this.Tplsql)
}

func (this *ModelItem) initItem() {
	this.Attrs = make(map[string]string)
}

func (this *ModelItem) initMethod(model, tag string) {
	this.Attrs = make(map[string]string)
	this.Model = model
	this.Tag = tag
}


func FindMethod(model, method string) *ModelItem {
	if Tpls == nil {
		return nil
	}
	return Tpls[model+"_"+method]
}

func pushMethod(item *ModelItem) {
	key := item.Model + "_" + item.Method
	initTpls()
	Tpls[key] = item
}

func Analysis_sql_file(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		abspath,_ := filepath.Abs(filePath)
		Ps(abspath)
		Log.Error("文件路径%s不存在\n", abspath)
		return
	}
	defer file.Close()

	xmldoc := xml.NewDecoder(file)
	modelName := ""
	isEnterElement := false
	modelItem := new(ModelItem)
	modelItem.initItem()
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
				Ps("tag", name)
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
				modelItem.initItem()
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

	t, err := template.New("new").Funcs(Tplfuncs).Parse(tpl)
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
