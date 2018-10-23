package fcgi


type ResponseContent struct {
	received chan bool
	err chan error
	buf []byte
}

func (res *ResponseContent) content() []byte  {
	return res.buf
}

//Response Interface
func (res *ResponseContent) String() string  {
	return string(res.buf)
}

//Response Interface
func (res *ResponseContent) Bytes() []byte  {
	return res.buf
}