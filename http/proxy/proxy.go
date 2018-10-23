package proxy

type IProxy interface {
	Parse()
	Do()
}

type Proxy struct {

}



func Do(proxy IProxy) {
	//IProxy.Do()
}