package gatis

import "text/template"

var Tplfuncs template.FuncMap

func init() {

	Tplfuncs = map[string]interface{}{
		"q" : addSingleQuote,
	}
}

func addSingleQuote(arg string) string {
	return "'" + arg + "'"
}