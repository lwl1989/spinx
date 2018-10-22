package proxy

type IProxy interface {
	Parse()
	Do()
}

type Proxy struct {
	request
}



func Do(proxy IProxy) {
	IProxy.Do()
}