package core

import (
	"strings"
)

func tryFiles(uri, rule string,env map[string]string) {
	env["REQUEST_URI"] = strings.Replace(rule, "$uri", uri, -1)
	env["SCRIPT_FILENAME"] = env["DOCUMENT_ROOT"]+"/"+"index.php"
	//env["QUERY_STRING"] = env["REQUEST_URI"]
	//todo get index.php reg
}
