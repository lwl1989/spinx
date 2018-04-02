package main

import (
	"os"
	"fmt"
)

func main(){
	f,e := os.Stat("/www/aaa")
	fmt.Println(e)
	fmt.Println(f)
}