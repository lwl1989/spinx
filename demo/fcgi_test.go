package main

import (
	"net/http/fcgi"
	"net"
	"net/http"
	"fmt"
	"io/ioutil"
)

type handle struct{

}

func (handle *handle) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	env := fcgi.ProcessEnv(r)
	fmt.Println(env)
	reqParams := ""
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		reqParams = string(body)

	}
	w.WriteHeader(200)
	w.Write([]byte("post:"+reqParams))
}

func main()  {
	listen,err := net.Listen("tcp","127.0.0.1:9002")

	if err != nil {
		panic(err)
	}
	handle := &handle{}
	fcgi.Serve(listen,handle)
}