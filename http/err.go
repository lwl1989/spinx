package http

type Error struct {
	code	uint16
	title 	string
	value   string
}

func GetError(code uint16, msg string) Error {
	return Error{
		code:code,
		value:msg,
	}
}


func GetRenderHtml(code uint16) []byte {
	return make([]byte, 0)
}