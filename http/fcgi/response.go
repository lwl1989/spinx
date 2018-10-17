package fcgi

type ResponseContent struct {
	received chan bool
	err chan error
	buf []byte
}

func (res *ResponseContent) content() []byte  {
	return res.buf
}