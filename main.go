package main

import (
	"github.com/lwl1989/spinx/http"
	"github.com/lwl1989/spinx/conf"
)

func main(){
	conf.GetVhosts("config/server.json")
	http.Do()
}




