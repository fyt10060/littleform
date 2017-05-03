package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/server"
)

func NewYarServer(class interface{}, rpcName, methodName string, r *http.Request, w http.ResponseWriter) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("body err:%s\n", err)
	}

	s := server.NewServer(class)
	s.Opt.LogLevel = yar.LogLevelDebug | yar.LoglevelNormal | yar.LogLevelError
	s.Register(rpcName, methodName)

	err1 := s.Handle(body, w)
	if err1 != nil {
		fmt.Printf("handle err:%s\n", err1)
	}
}
