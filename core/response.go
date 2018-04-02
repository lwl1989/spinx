package core

import (
	"net/http"
	"errors"
	"strings"
	"strconv"
)



type Response struct {
	code uint16
	headers map[string]string
	cookies []*http.Cookie
	content string
}

func GetResponseByContent(code int, headers map[string]string, cookies []*http.Cookie,content string) *Response {
	return &Response{code:uint16(code),headers:headers,cookies:cookies,content:content}
}
func GetResponse(content string) (*Response) {
	status, headers, cookies, body, err := parseFastCgiResponse(content)
	if err != nil {
		return &Response{code:500,headers:nil,cookies:nil,content:""}
	}
	return &Response{
		code:status,
		headers:headers,
		cookies:parseCookies(cookies),
		content:body,
	}
}

func (response *Response) getCode() uint16 {
	if response.code > 0 {
		return response.code
	}
	return 200
}

func (response *Response) getContent() string {
	return response.content
}

func (response *Response) getCookies() []*http.Cookie {
	return response.cookies
}

func (response *Response) getHeaders() map[string]string {
	return response.headers
}

func (response *Response) send(write http.ResponseWriter,r *http.Request) {
	code := int(response.getCode())

	if len(response.headers) > 0 {
		for key,value := range response.getHeaders() {
			write.Header().Set(key,value)
		}
	} else {
		write.Header().Set("Content-Type","text/html")
	}

	if len(response.cookies) > 0 {
		for _,cookie := range response.getCookies() {
			http.SetCookie(write,cookie)
		}
	}

	//must set header first
	write.WriteHeader(int(response.getCode()))

	switch code {
	case 302:
		http.Redirect(write,r,response.getHeaders()["Location"],code)
		break
	default:
		write.Write([]byte(response.getContent()))
	}
}


func parseFastCgiResponse(content string) (uint16, map[string]string, []string, string, error) {
	var headers = make(map[string]string)
	var cookies = make([]string,0)
	parts := strings.SplitN(content, "\r\n\r\n", 2)

	if len(parts) < 2 {
		return 502, headers, cookies, "", errors.New("Cannot parse FastCGI Response")
	}

	headerParts := strings.Split(parts[0], "\r\n")
	body := parts[1]
	status := 200

	for _,header := range headerParts {
		lineParts := strings.SplitN(header, ":", 2)

		if len(lineParts) < 2 {
			continue
		}

		lineParts[1] = strings.TrimSpace(lineParts[1])

		if lineParts[0] == "Status" {
			var err error
			code := strings.SplitN(lineParts[1]," ", 2)
			if len(code) < 2 {
				status,err = strconv.Atoi(code[0])
				if err != nil {
					status = 500
				}
			}else{
				status,err = strconv.Atoi(code[0])
				if err != nil {
					status = 500
				}
			}
			continue
		}

		if lineParts[0] == "Set-Cookie" {
			cookies = append(cookies, lineParts[1])
		}

		if lineParts != nil {
			headers[lineParts[0]] = lineParts[1]
		}
	}

	return uint16(status), headers, cookies, body, nil
}



func parseCookies(cookie []string) []*http.Cookie{
	return read(cookie)
}


func respond(w http.ResponseWriter, body string, statusCode int, headers map[string]string) {

	for header, value := range headers {
		w.Header().Set(header, value)
	}
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

