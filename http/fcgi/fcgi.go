package fcgi

//主要用于解析http
//并且将协议拆解成N份  主要是：头(一般头，自定义头, fast cgi协议头) 主体
//cookies不需要解析 原包不动即可
//通过解析结果获取到请求的配置
//结果不需要拆解 直接发送 提高性能

func buildEnv(req *Request) (err error, env map[string]string) {
	env = make(map[string]string)
	index := "index.php"
	//todo:ceshi
	req.cf.DocumentRoot = "/Users/wenglong11/PhpstormProjects/yaf/public"
	filename := req.cf.DocumentRoot + "/" + index


	//for name, value := range req.cf.serverEnvironment {
	//	env[name] = value
	//}


	env["DOCUMENT_ROOT"] = req.cf.DocumentRoot
	env["SCRIPT_FILENAME"] = filename
	env["SCRIPT_NAME"] = "/" + index

	env["REQUEST_SCHEME"] = "http"
	env["SERVER_NAME"] = ""
	env["SERVER_PORT"] = req.Port
	env["REDIRECT_STATUS"] = "200"
	env["HTTP_HOST"] = req.Host
	env["SERVER_SOFTWARE"] = "spinx/1.0.0"
	//env["REMOTE_ADDR"] = "127.0.0.1"
	env["SERVER_PROTOCOL"] = "HTTP/1.1"
	env["GATEWAY_INTERFACE"] = "CGI/1.1"

	env["PATH_INFO"] = ""
	//if r.URL.Path != "/" {
	//	env["PATH_INFO"] = r.URL.Path
	//}
	env["REQUEST_METHOD"] = req.Method
	env["REQUEST_URI"] = "/"
	env["QUERY_STRING"] = ""
	if env["QUERY_STRING"] != "" {
		env["REQUEST_URI"] = env["REQUEST_URI"] + "?" + env["QUERY_STRING"]
	} else {
		env["REQUEST_URI"] = env["REQUEST_URI"]
	}

	env["DOCUMENT_URI"] = env["SCRIPT_NAME"]
	env["PHP_SELF"] = env["SCRIPT_NAME"]


	//只有自定义的需要加上HTTP_xxx_XXX_XXX
	//for header, values := range r.Header {
	//	env["HTTP_"+strings.Replace(strings.ToUpper(header), "-", "_", -1)] = values[0]
	//}

	env["CONTENT_LENGTH"] = "0"
	env["CONTENT_TYPE"] = "text/html"
	//log.Println()
	return nil, env
}