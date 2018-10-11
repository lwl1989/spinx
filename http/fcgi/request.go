package fcgi

import "bufio"

//用于获取头的位置
const SPLIT_STR  =  "\n\n"
//用户获取特定头
const ENTER_SPACE  = "\n"
//请求头部字符串
const REQUEST_URL  = "Request URL:"

type Request struct {
	Id  uint16
	rwc *bufio.Reader
	Host string
	Port string
	Header map[string]string  //必要设置的Header
	KeepConn bool
	pos position
	content []byte
}

//记录流的位置
type position struct {
	LEN    int64   //请求流的长度
	HStart int64   //头的开始位置
	HEnd   int64   //头的结束位置
	BStart int64   //正文开始位置
}

func (req *Request) setHeader(key, value string)  {

}


func (req *Request) getHeader(key string) string {
	return ""
}


func (req *Request) getHeaderPosition() (start, end int64)  {
	return start,end
}

//执行此方法 获取 Request URL:
//获取解析到HOST PORT 获得CGI转发的配置
//并且获取到双换行的位置
//然后通过配置获取到的参数 修改头 协议需要的 增加头
//document_file  index等等
func (req *Request) Parse()  {

}
