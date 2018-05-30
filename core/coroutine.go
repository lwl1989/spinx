package core

import (
	"net/http"
	"strings"
	"errors"
	"fmt"
)

const (
	FileCodeStatic uint8 = 1
	FileCodeTry	   = 2
	FileExecute	   = 3
	FileCodeNotFound = 4
)

//is corouine fun
// do  httpHandler.Response <- Response
// or do httpHandler.StaticFile <- StaticFileHandler
func (httpHandler *HttpHandler) Run(r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		httpHandler.Response = &Response{200, map[string]string{}, nil, ""}
		return
	}
	//reqParams := ""
	name, port := buildNamePort(r.Host)

	hm, err := httpHandler.Vhosts.GetHostMap(port, name)
	if err != nil {
		httpHandler.Response = &Response{200, map[string]string{}, nil, "<h1>404</h1>"}
		return
	}

	documentRoot := hm.DocumentRoot

	err, env := httpHandler.buildEnv(documentRoot, r)

	var response *Response
	if err != nil {
		httpHandler.log.Println(err)
		httpHandler.log.Println(env)
		response = GetResponseByContent(403, nil, nil, "not allow")

	} else {

		fileCode,filename := httpHandler.buildServerHttp(r, env, hm)

		if fileCode == FileCodeStatic {


				httpHandler.StaticFile = &StaticFileHandler{
					name,
					port,
					filename,
				}
				return

		}

		if fileCode == FileCodeTry {
			tryFiles(r.RequestURI, hm.TryFiles, env)
		}

		if fileCode == FileCodeNotFound {
			response = &Response{404,nil,nil,"<h1>404</h1>"}
			return
		}


		fcgi, err := New(hm.Net, hm.Addr)

		req := fcgi.GetRequest(r, env)

		if err != nil {
			httpHandler.log.Printf("err: %v", err)
		}

		content, _, err := fcgi.DoRequest(req)

		if err != nil {
			httpHandler.log.Printf("ERROR: %s - %v", r.URL.Path, err)
		}

		response = GetResponse(fmt.Sprintf("%s", content))

	}

	httpHandler.Response = response
}
//func (httpHandler *HttpHandler) Run(r *http.Request) {
//	if r.RequestURI == "/favicon.ico" {
//		httpHandler.Response <- &Response{200, map[string]string{}, nil, ""}
//		return
//	}
//	//reqParams := ""
//	name, port := buildNamePort(r.Host)
//
//	hm, err := httpHandler.Vhosts.GetHostMap(port, name)
//	if err != nil {
//		httpHandler.Response <- &Response{200, map[string]string{}, nil, "<h1>404</h1>"}
//		return
//	}
//
//	documentRoot := hm.DocumentRoot
//
//	err, env := httpHandler.buildEnv(documentRoot, r)
//
//	var response *Response
//	if err != nil {
//		httpHandler.log.Println(err)
//		httpHandler.log.Println(env)
//		response = GetResponseByContent(403, nil, nil, "not allow")
//
//	} else {
//
//		fileCode,filename := httpHandler.buildServerHttp(r, env, hm)
//		switch fileCode {
//		case FileCodeStatic:
//			httpHandler.StaticFile <- &StaticFileHandler{
//				name,
//				port,
//				filename,
//			}
//			return
//
//		case FileCodeNotFound:
//			response = &Response{404,nil,nil,"<h1>404</h1>"}
//			return
//		case FileCodeTry:
//			tryFiles(r.RequestURI, hm.TryFiles, env)
//
//		}
//		fcgi, err := New(hm.Net, hm.Addr)
//
//		req := fcgi.GetRequest(r, env)
//
//		if err != nil {
//			httpHandler.log.Printf("err: %v", err)
//		}
//
//		content, _, err := fcgi.DoRequest(req)
//
//		if err != nil {
//			httpHandler.log.Printf("ERROR: %s - %v", r.URL.Path, err)
//		}
//
//		response = GetResponse(fmt.Sprintf("%s", content))
//
//	}
//
//	httpHandler.Response <- response
//}

//build a header map with fcgi header
func (httpHandler *HttpHandler) buildEnv(documentRoot string, r *http.Request) (err error, env map[string]string) {
	env = make(map[string]string)
	index := "index.php"
	filename := documentRoot + "/" + index
	if r.URL.Path == "/.env" {
		return errors.New("not allow"), env
	} else if r.URL.Path == "/" || r.URL.Path == "" {
		filename = documentRoot + "/" + index
	} else {
		filename = documentRoot + r.URL.Path
	}

	for name, value := range httpHandler.serverEnvironment {
		env[name] = value
	}

	name, port := buildNamePort(r.Host)

	env["DOCUMENT_ROOT"] = documentRoot
	env["SCRIPT_FILENAME"] = filename
	env["SCRIPT_NAME"] = "/" + index

	env["REQUEST_SCHEME"] = "http"
	env["SERVER_NAME"] = name
	env["SERVER_PORT"] = port
	env["REDIRECT_STATUS"] = "200"
	env["HTTP_HOST"] = r.Host
	env["SERVER_SOFTWARE"] = "spinx/1.0.0"
	env["REMOTE_ADDR"] = r.RemoteAddr
	env["SERVER_PROTOCOL"] = "HTTP/1.1"
	env["GATEWAY_INTERFACE"] = "CGI/1.1"

	env["PATH_INFO"] = ""
	if r.URL.Path != "/" {
		env["PATH_INFO"] = r.URL.Path
	}
	env["REQUEST_METHOD"] = r.Method
	env["QUERY_STRING"] = r.URL.RawQuery
	if r.URL.RawQuery != "" {
		env["REQUEST_URI"] = r.URL.Path + "?" + r.URL.RawQuery
	} else {
		env["REQUEST_URI"] = r.URL.Path
	}

	env["DOCUMENT_URI"] = env["SCRIPT_NAME"]
	env["PHP_SELF"] = env["SCRIPT_NAME"]

	for header, values := range r.Header {
		env["HTTP_"+strings.Replace(strings.ToUpper(header), "-", "_", -1)] = values[0]
	}

	env["CONTENT_LENGTH"] = r.Header.Get("Content-Length")
	env["CONTENT_TYPE"] = r.Header.Get("Content-Type")
	//log.Println()
	return nil, env
}