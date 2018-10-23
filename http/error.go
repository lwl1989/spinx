package http

import "fmt"

type ErrorMsg struct {
	code	uint16
	title 	string
	value   string
}

func GetError(code uint16, msg string) ErrorMsg {
	return ErrorMsg{
		code:code,
		value:msg,
	}
}


func GetRenderHtml(code uint16) []byte {
	return make([]byte, 0)
}

func (msg ErrorMsg) Error() string {
	return fmt.Sprintf("HTTP/1.1  %d \r\n\r\n<h1>%d</h1><p>%s</p>", msg.code, msg.code, msg.value)
}
