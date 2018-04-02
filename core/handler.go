package core

import (
	"net/http"
	"os"
	"strings"
	"fmt"
	"errors"
)

var serverEnvironment map[string]string
type HttpHandler struct{
	Vhosts Vhosts
	HandlerMap map[string]*http.ServeMux
}

func GetHandler() *HttpHandler {
	return &HttpHandler{}
}

func (httpHandler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		w.WriteHeader(200)
		w.Write([]byte(""))
		return
	}
	//reqParams := ""
	name,port := buildNamePort(r.Host)

	hm,err := httpHandler.Vhosts.GetHostMap(port,name)
	if err != nil {
		respond(w, "<h1>404</h1>", 404, map[string]string{})
		return
	}
	documentRoot := hm.DocumentRoot

	err,env := buildEnv(documentRoot,r)

	var response *Response
	if err != nil {
		response = GetResponseByContent(403, nil,nil, "not allow")
	}else {
		fileCode := httpHandler.buildServerHttp(w,r,env,hm)
		if fileCode == 3 {
			return
		} else if fileCode == 4 {
			tryFiles(r.RequestURI, hm.TryFiles, env)
			fmt.Println(env)
			//response = GetResponseByContent(code, nil,nil, "not allow")
		}
			fcgi, err := New(hm.Net, hm.Addr)

			req := fcgi.GetRequest(w,r,env)
			//fmt.Println(req,res,err)
			if err != nil {
				fmt.Printf("err: %v", err)
			}

			content, _, err := fcgi.DoRequest(req)

			if err != nil {
				fmt.Printf("ERROR: %s - %v", r.URL.Path, err)
			}
			response = GetResponse(fmt.Sprintf("%s", content))

	}

	response.send(w,r)
}


func buildEnv(documentRoot string,r *http.Request) (err error,env map[string]string){
	env = make(map[string]string)
	index := "index.php"
	filename := documentRoot+"/"+index
	if r.URL.Path == "/.env" {
		return errors.New("not allow"),env
	} else if r.URL.Path == "/" || r.URL.Path == "" {
		filename = documentRoot + "/" + index
	} else {
		filename = documentRoot + r.URL.Path
	}

	for name,value := range serverEnvironment {
		env[name] = value
	}

	name,port := buildNamePort(r.Host)

	env["DOCUMENT_ROOT"] 	= documentRoot
	env["SCRIPT_FILENAME"] 	= filename
	env["SCRIPT_NAME"] 		= "/"+index


	env["REQUEST_SCHEME"] 	= "http"
	env["SERVER_NAME"] 		= name
	env["SERVER_PORT"] 		= port
	env["REDIRECT_STATUS"] 	= "200"
	env["HTTP_HOST"] 		= r.Host
	env["SERVER_SOFTWARE"] 	= "spinx/1.0.0"
	env["REMOTE_ADDR"] 		= r.RemoteAddr
	env["SERVER_PROTOCOL"] 	= "HTTP/1.1"
	env["GATEWAY_INTERFACE"] = "CGI/1.1"

	env["PATH_INFO"] = ""
	if r.URL.Path != "/" {
		env["PATH_INFO"] 	= r.URL.Path
	}
	env["REQUEST_METHOD"] 	= r.Method
	env["QUERY_STRING"] 	= r.URL.RawQuery
	if r.URL.RawQuery != "" {
		env["REQUEST_URI"] 	= r.URL.Path + "?" + r.URL.RawQuery
	} else {
		env["REQUEST_URI"] 	= r.URL.Path
	}

	env["DOCUMENT_URI"] 	= env["SCRIPT_NAME"]
	env["PHP_SELF"] 		= env["SCRIPT_NAME"]

	for header, values := range r.Header {
		env["HTTP_" + strings.Replace(strings.ToUpper(header), "-", "_", -1)] = values[0]
	}

	env["CONTENT_LENGTH"] = r.Header.Get("Content-Length")
	env["CONTENT_TYPE"] = r.Header.Get("Content-Type")
	//fmt.Println()
	return errors.New("not allow"),env
}

/**
return file like http code
4 not found
3 static file
success is 0
 */
func  (httpHandler *HttpHandler) buildServerHttp(w http.ResponseWriter, r *http.Request,env map[string]string,hm *HostMap) (fileCode int) {
	name,port := buildNamePort(r.Host)
	filename := env["SCRIPT_FILENAME"]

	staticHandler := httpHandler.HandlerMap[name+port]
	file, err := os.Stat(filename)
	//file exists
	if err == nil {
		if file.IsDir() {
			filename = getScriptFile(filename,hm.Index)
		}
	} else {
		return 4
	}

	if !strings.HasSuffix(filename, ".php") {
		fmt.Println(filename)
		//tryFiles(r.RequestURI, hm.TryFiles, env)

		staticHandler.ServeHTTP(w, r)
		return 3
	}
	env["SCRIPT_FILENAME"] = filename
	return 0
}

func getScriptFile(filename string, index string) string {
	file, err := os.Stat(filename)
	if err != nil {
		return ""
	}
	if !file.IsDir() {
		return filename
	} else {
		return batchIndex(strings.TrimRight(filename,"/"), index)
	}
}

func batchIndex(path string, index string) (real string) {
	if index != "" {
		indexs := strings.Split(index," ")
		for _,v := range indexs {
			info,err := os.Stat(path+"/"+v)
			if err != nil {
				continue
			}
			if info.IsDir() {
				continue
			}
			real = path+"/"+v
			break
		}
	}
	return real
}

func buildNamePort(url string) (name,port string) {
	req := strings.Split(url,":")
	if len(req) == 2 {
		port = req[1]
		name = req[0]
	}else {
		port = "80"
		name = req[0]
	}
	return name,port
}