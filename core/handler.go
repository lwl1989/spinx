package core

import (
	"net/http"
	"os"
	"strings"
	"fmt"
	"log"
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
type HttpHandler struct {
	Vhosts            Vhosts
	HandlerMap        map[string]*http.ServeMux
	Response          chan *Response
	StaticFile        chan *StaticFileHandler
	serverEnvironment map[string]string
	log               *log.Logger
}

func GetHandler() *HttpHandler {
	return &HttpHandler{
		Response:          make(chan *Response),
		serverEnvironment: make(map[string]string, 0),
		StaticFile:        make(chan *StaticFileHandler),
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

//listen http do some things
func (httpHandler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go httpHandler.Run(r)
	for {
		select {
		case response := <-httpHandler.Response:
			response.send(w, r)
		case hand := <-httpHandler.StaticFile:
			staticHandler := httpHandler.HandlerMap[hand.Host+hand.Port]
			staticHandler.ServeHTTP(w, r)

		default:
			respond(w, "<h1>404</h1>", 404, map[string]string{})
		}
	}
}

/**
return file like http code
2 try file
1 exit file
success is 0
 */
func (httpHandler *HttpHandler) buildServerHttp(r *http.Request, env map[string]string, hm *HostMap) (code int, filename string) {
	//name,port := buildNamePort(r.Host)
	filename = env["SCRIPT_FILENAME"]
	file, err := os.Stat(filename)

	if err == nil {
		if file.IsDir() {
			filename = getScriptFile(filename, hm.Index)
			if filename == "" {
				return FileCodeNotFound, ""
			}
		}
	} else {
		return FileCodeStatic, filename
	}

	if !strings.HasSuffix(filename, ".php") {
		fmt.Println(filename)
		//tryFiles(r.RequestURI, hm.TryFiles, env)

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
