package core

import (
	"github.com/jingweno/conf"
	"fmt"
	"log"
	"strings"
	"reflect"
	"errors"
)


type HostMap struct{
	Name string
	Net string
	Addr string
	TryFiles string
	DocumentRoot string
	Index string
}
type Vhosts map[string]map[string]*HostMap


func GetVhosts(path string) Vhosts {

	c, err := conf.NewLoader().Argv().Env().File(path).Load()

	if err != nil {
		fmt.Println(err)
		log.Fatal("lode config err")
		return nil
	}

	config := c.Get("vhosts")
	var hMap = make(Vhosts)
	for _,v := range config.([]interface{}) {
		keyValue := v.(map[string]interface{})
		port := keyValue["port"].(string)
		proxy := getValue(keyValue,"proxy")
		isIp := strings.Contains(getValue(keyValue,"proxy"),".")
		net,addr := "tcp",""
		if !isIp {
			net = "unix"
		}
		addr = proxy
		names  := strings.Split(getValue(keyValue,"name")," ")


		if _,ok := hMap[port]; !ok {
			hMap[port] = make(map[string]*HostMap)
			for _, name := range names {
				hMap[port][name] = &HostMap{
					Name: name,
					Net:  net,
					Addr: addr,
					TryFiles:      getValue(keyValue, "tryFiles"),
					DocumentRoot: getValue(keyValue, "documentRoot"),
					Index:		getValue(keyValue, "index"),
				}
			}
		}else{
			for _, name := range names {
				hMap[port][name] = &HostMap{
					Name: name,
					Net:  net,
					Addr: addr,
					TryFiles:      getValue(keyValue, "tryFiles"),
					DocumentRoot: getValue(keyValue, "documentRoot"),
				}
			}
		}
	}
	return hMap
}

func (vhosts Vhosts) GetPorts() []string {
	var hosts = make([]string,0)
	for k,_ := range vhosts {
		hosts  = append(hosts, k)
	}
	return hosts
}

func (vhosts Vhosts) GetNames(port string) []string {
	var names = make([]string,0)
	for k,_ :=  range vhosts[port] {
		names = append(names, k)
	}
	return names
}

func (vhosts Vhosts) GetHostMap(port,host string) (*HostMap,error) {
	if _,ok := vhosts[port][host]; !ok {
		return &HostMap{},errors.New("config not found")
	}
	return vhosts[port][host],nil
}
func (vhosts Vhosts) Get(port,host,key string) string {
	hMap := vhosts[port][host]
	object := reflect.ValueOf(hMap)
	for i:=0; i<object.NumField(); i++{
		field := object.Field(i)
		if object.Type().Field(i).Name == key {
			return field.String()
		}
	}
	return ""
}
func getValue(keyValue map[string]interface{},key string) string {
	if _,ok := keyValue[key]; !ok {
		return ""
	}
	return  keyValue[key].(string)
}

