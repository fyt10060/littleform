package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/server"
)

type RegisterPair struct {
	rpcName    string
	methodName string
}

type RegisterList struct {
	List []RegisterPair
}

// 可复用yar server 声明
func NewYarServer(class interface{}, registers RegisterList, r *http.Request, w http.ResponseWriter) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("body err:%s\n", err)
	}

	s := server.NewServer(class)
	s.Opt.LogLevel = yar.LogLevelDebug | yar.LoglevelNormal | yar.LogLevelError
	// 遍历列表注册所有方法
	for _, v := range registers.List {
		s.Register(v.rpcName, v.methodName)
	}

	err1 := s.Handle(body, w)
	if err1 != nil {
		fmt.Printf("handle err:%s\n", err1)
	}
}
