package controllers

import (
	"fmt"

	"io/ioutil"

	"littleform/model"

	"github.com/astaxie/beego"
	"github.com/mitchellh/mapstructure"
	"github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/server"
)

type WeixinController struct {
	beego.Controller
}

type Weixin struct{}

func (c *Weixin) Weixin(accountParam, messageParam, otherParam map[string]interface{}) map[string]interface{} {
	fmt.Println("123\n")

	var reply = model.WeixinReply{}
	reply.Action = "reply"
	reply.Exit = false
	reply.Log = false
	var account model.Account
	err := mapstructure.Decode(accountParam, &account)
	if err != nil {
		fmt.Println(err)
		return model.Struct2Map(reply)
	}
	fmt.Printf("appid: %s, account name: %s", account.Appid, account.Name)

	var message model.WeiXinMessageInfo
	err = mapstructure.Decode(messageParam, &message)
	if err != nil {
		fmt.Println(err)
		return model.Struct2Map(reply)
	}
	fmt.Printf("message type: %s, message content:%s", message.MsgType, message.Content)

	var other model.HostWxOther
	err = mapstructure.Decode(otherParam, &other)
	if err != nil {
		fmt.Println(err)
		return model.Struct2Map(reply)
	}

	if message.MsgType == "text" {
		if message.Content == "a2" {
			reply.Exit = true
			reply.Log = true
			var nMap map[string]interface{}
			nMap = make(map[string]interface{})
			nMap["message"] = model.GetReplyWithSendMsg(message, "a2 + 付雨桐")
			nMap["type"] = "raw"
			reply.Data = nMap
		}
	}

	return model.Struct2Map(reply)
}

func (this *WeixinController) Post() {
	c := this.Ctx
	r := c.Request
	w := c.ResponseWriter

	body, err := ioutil.ReadAll(r.Body)
	fmt.Printf("request medthod:%s\n", r.Method)
	fmt.Printf("request url: %s\n", r.URL.String())

	if err != nil {
		fmt.Printf("body err:%s\n", err)
	}
	fmt.Printf("body:%x\n", body)

	s := server.NewServer(&Weixin{})
	s.Opt.LogLevel = yar.LogLevelDebug | yar.LoglevelNormal | yar.LogLevelError
	s.Register("weixin", "Weixin")

	err1 := s.Handle(c.Input.RequestBody, w)
	if err1 != nil {
		fmt.Printf("handle err:%s\n", err1)
	}

	//	w.Write([]byte("11313"))
}
