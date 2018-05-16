package core

import (
	"net/http"
	"os"
	"log"
	"strings"
)

//if is static
//need host port to response
//the next version want cache file
type StaticFileHandler struct {
	Host     string
	Port     string
	FilePath string
}


//httpHandler listen host:port
//type HttpHandler struct {
//	Vhosts            Vhosts
//	HandlerMap        map[string]*http.ServeMux
//	Response          chan *Response
//	StaticFile        chan *StaticFileHandler
//	serverEnvironment map[string]string
//	log               *log.Logger
//}
type HttpHandler struct {
	Vhosts            Vhosts
	HandlerMap        map[string]*http.ServeMux
	Response          *Response
	StaticFile        *StaticFileHandler
	serverEnvironment map[string]string
	log               *log.Logger
}

func GetHandler() *HttpHandler {
	return &HttpHandler{
		serverEnvironment: make(map[string]string, 0),
		log:			log.New(os.Stdout,"[spinx]", log.LstdFlags),
	}
}
//set logger to handler
func (httpHandler *HttpHandler) SetLogger(log *log.Logger) {
	httpHandler.log = log
}

//get logger from handler
func (httpHandler *HttpHandler)  GetLogger() *log.Logger {
	return httpHandler.log
}

//func (httpHandler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//
//	if r.RequestURI == "/favicon.ico" {
//		w.WriteHeader(200)
//		w.Write([]byte(""))
//		return
//	}
//	//reqParams := ""
//	name, port := buildNamePort(r.Host)
//
//	hm, err := httpHandler.Vhosts.GetHostMap(port, name)
//	if err != nil {
//		respond(w, "<h1>404</h1>", 404, map[string]string{})
//		return
//	}
//
//	documentRoot := hm.DocumentRoot
//	static := documentRoot + r.URL.Path
//	fi,err := os.Stat(static)
//	if err == nil && !fi.IsDir() {
//		staticHandler := httpHandler.HandlerMap[name+port]
//		staticHandler.ServeHTTP(w, r)
//		return
//	}
//
//	err, env := httpHandler.buildEnv(documentRoot, r)
//
//	var response *Response
//	if err != nil {
//		response = GetResponseByContent(403, nil, nil, "not allow")
//	} else {
//
//		fileCode,_ := httpHandler.buildServerHttp(r, env, hm)
//		if FileCodeTry == fileCode {
//			tryFiles(r.RequestURI, hm.TryFiles, env)
//		}
//		fcgi, err := New(hm.Net, hm.Addr)
//
//		req := fcgi.GetRequest(r, env)
//		//fmt.Println(req,res,err)
//		if err != nil {
//			fmt.Printf("err: %v", err)
//		}
//
//		content, _, err := fcgi.DoRequest(req)
//
//		if err != nil {
//			fmt.Printf("ERROR: %s  %v", r.URL.Path, err)
//		}
//		response = GetResponse(fmt.Sprintf("%s", content))
//	}
//
//	response.send(w, r)
//}
//listen http do some things
func (httpHandler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.RequestURI)
	httpHandler.Run(r)

	if httpHandler.StaticFile != nil {
		staticHandler := httpHandler.HandlerMap[httpHandler.StaticFile.Host+httpHandler.StaticFile.Port]
		staticHandler.ServeHTTP(w, r)
		return
	}

	if httpHandler.Response != nil {
		httpHandler.Response.send(w, r)
		return
	}



	respond(w, "<h1>404</h1>", 404, map[string]string{})
	//go httpHandler.Run(r)
	//
	//select {
	//	case response := <-httpHandler.Response:
	//		response.send(w, r)
	//	case hand := <-httpHandler.StaticFile:
	//		staticHandler := httpHandler.HandlerMap[hand.Host+hand.Port]
	//		staticHandler.ServeHTTP(w, r)
	//
	//	default:
	//		respond(w, "<h1>404</h1>", 404, map[string]string{})
	//}

}

/**
return file like http code
2 try file
1 exit file
success is 0
 */
func (httpHandler *HttpHandler) buildServerHttp(r *http.Request, env map[string]string, hm *HostMap) (code uint8, filename string) {
	//name,port := buildNamePort(r.Host)
	filename = env["SCRIPT_FILENAME"]
	file, err := os.Stat(filename)

	if err == nil {
		if file.IsDir() {
			filename = getScriptFile(filename, hm.Index)
			return FileCodeTry, filename
			//if filename == "" {
			//	return FileCodeNotFound, ""
			//}
		}
	} else {
		return FileCodeTry, filename
	}

	env["SCRIPT_FILENAME"] = filename
	return FileExecute, filename
}

func getScriptFile(filename string, index string) string {
	file, err := os.Stat(filename)
	if err != nil {
		return ""
	}
	if !file.IsDir() {
		return filename
	} else {
		return batchIndex(strings.TrimRight(filename, "/"), index)
	}
}

func batchIndex(path string, index string) (real string) {
	if index != "" {
		indexs := strings.Split(index, " ")
		for _, v := range indexs {
			info, err := os.Stat(path + "/" + v)
			if err != nil {
				continue
			}
			if info.IsDir() {
				continue
			}
			real = path + "/" + v
			break
		}
	}
	return real
}

func buildNamePort(url string) (name, port string) {
	req := strings.Split(url, ":")
	if len(req) == 2 {
		port = req[1]
		name = req[0]
	} else {
		port = "80"
		name = req[0]
	}
	return name, port
}
