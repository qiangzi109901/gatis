package gatis


var Tplfuncs map[string]interface{}

func init() {

	Tplfuncs = map[string]interface{}{
		"q" : addSingleQuote,
	}
}

func addSingleQuote(arg string) string {
	return "'" + arg + "'"
}