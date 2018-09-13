package fcgi

import (
	"sync"
)

//主要用于解析http
//并且将协议拆解成N份  主要是：头(一般头，自定义头, fast cgi协议头) 主体
//cookies不需要解析 原包不动即可
//通过解析结果获取到请求的配置
//结果不需要拆解 直接发送 提高性能


type CgiProtocol struct {
	mutex sync.Mutex
}

func (cgi *CgiProtocol) writeRecord(recType uint8, reqId uint16, content []byte) (err error) {
	cgi.mutex.Lock()
	defer cgi.mutex.Unlock()
	//cgi.buf.Reset()
	//cgi.h.init(recType, reqId, len(content))
	//
	//if err := binary.Write(&cgi.buf, binary.BigEndian, cgi.h); err != nil {
	//	return err
	//}
	//if _, err := cgi.buf.Write(content); err != nil {
	//	return err
	//}
	//if _, err := cgi.buf.Write(pad[:cgi.h.PaddingLength]); err != nil {
	//	return err
	//}
	//_, err = cgi.rwc.Write(cgi.buf.Bytes())
	return err
}