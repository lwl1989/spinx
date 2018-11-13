package http_test

import (
	"testing"
	corehttp "net/http"
	//"github.com/lwl1989/spinx/http"
	"fmt"
	"errors"
	lihttp "github.com/lwl1989/spinx/http"
)

func TestResponse(t *testing.T) {
	res := &corehttp.Response{}
	res1 := &lihttp.Response{}
	maps := make(map[string]interface{})
	maps["http"] = res
	maps["http_li"] = res1
	maps["error"] = errors.New("errors")

	for k,r1 := range maps {
		switch r1.(type) {
		case error:
			fmt.Println(k,r1)
			break
		case *corehttp.Response:
			fmt.Println(k,r1)
			break
		case *lihttp.Response:
			fmt.Println(k,r1)
			break
		default:
			fmt.Println(k,r1)
		}
	}
}
