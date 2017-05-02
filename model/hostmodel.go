// hostmodel
package model

import (
	//	"fmt"
	"encoding/xml"
	"reflect"
	//	"strconv"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data

}

func GetReplyWithSendMsg(oMsg WeiXinMessageInfo, content string) (xmlString string) {
	//	var time string
	//	if reflect.TypeOf(oMsg.CreateTime) == string {
	//		time = oMsg.CreateTime
	//	} else {
	//		time = strconv.FormatInt(oMsg.CreateTime, 10)
	//	}
	//	time := strconv.FormatInt(oMsg.CreateTime, 10)
	time := oMsg.CreateTime
	reply := BaseReply{
		FromUserName: CDATAText{oMsg.ToUserName},
		ToUserName:   CDATAText{oMsg.FromUserName},
		MsgType:      CDATAText{oMsg.MsgType},
		CreateTime:   CDATAText{time},
	}

	text := Text{
		BaseReply: reply,
		Content:   CDATAText{content},
	}
	result, _ := xml.Marshal(text)
	return string(result)
}

type WeixinReply struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
	Exit   bool                   `json:"exit"`
	Log    bool                   `json:"log"`
}

type Text struct {
	XMLName xml.Name `xml:"xml"`
	BaseReply
	Content CDATAText
}

type BaseReply struct {
	FromUserName CDATAText
	ToUserName   CDATAText
	MsgType      CDATAText
	CreateTime   CDATAText
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

type Account struct {
	Id             int                `db:"id" json:"id"`
	IsBan          int                `db:"is_ban" json:"is_ban"`
	Name           string             `db:"name" json:"name"`
	WeixinOriId    string             `db:"weixin_ori_id" json:"weixin_ori_id"`
	WeixinId       string             `db:"weixin_id" json:"weixin_id"`
	Status         int                `db:"status" json:"status"`
	Type           int                `db:"type" json:"type"`
	Verify         int                `db:"verify" json:"verify"`
	ComponentAppid string             `db:"component_appid" json:"component_appid"`
	Appid          string             `db:"appid" json:"appid"`
	Appsecret      string             `db:"appsecret" json:"appsecret"`
	Advance        int                `db:"advance" json:"advance"`
	OwnerId        int                `db:"owner_id" json:"owner_id"`
	Dateline       int                `db:"dateline" json:"dateline"`
	Postip         string             `db:"postip" json:"postip"`
	Actived        int                `db:"actived" json:"actived"`
	WxAvatar       string             `db:"wx_avatar" json:"wx_avatar"`
	WxQrcode       string             `db:"wx_qrcode" json:"wx_qrcode"`
	OwnerOpenid    string             `db:"owner_openid" json:"owner_openid"`
	FromAuth       int                `db:"from_auth" json:"from_auth"`
	AuthOriInfo    *AuthOriInfo       `db:"auth_ori_info" json:"auth_ori_info"`
	Industry       int                `db:"industry" json:"industry"`
	ComponentList  []AccountComponent `db:"-" json:"component_list"`
}

type AuthOriInfo map[string]interface{}

// AccountComponent 公众号第三方对应关系
type AccountComponent struct {
	AccountAppid   string `db:"account_appid" json:"account_appid"`
	ComponentAppid string `db:"component_appid" json:"component_appid"`
	Dateline       int    `db:"dateline" json:"dateline"`
}

type WeiXinMessageInfo struct {
	FromUserName     string  `json:"FromUserName"`
	ToUserName       string  `json:"ToUserName"`
	CreateTime       string  `json:"CreateTime"`
	MsgType          string  `json:"MsgType"`
	MsgID            string  `json:"MsgId"`
	Event            string  `json:"Event,omitempty"`
	EventKey         string  `json:"EventKey,omitempty"`
	Content          string  `json:"Content,omitempty"`
	MediaID          string  `json:"MediaId,omitempty"`
	PicURL           string  `json:"PicUrl,omitempty"`
	Format           string  `json:"Format,omitempty"`
	ThumbMediaID     string  `json:"ThumbMediaId,omitempty"`
	LocationX        float64 `json:"Location_X,omitempty"`
	LocationY        float64 `json:"Location_Y,omitempty"`
	Scale            int     `json:"Scale,omitempty"`
	Label            string  `json:"Label,omitempty"`
	Title            string  `json:"Title,omitempty"`
	Description      string  `json:"Description,omitempty"`
	URL              string  `json:"Url,omitempty"`
	Ticket           string  `json:"Ticket,omitempty"`
	Latitude         float64 `json:"Latitude,omitempty"`
	Longitude        float64 `json:"Longitude,omitempty"`
	Precision        float64 `json:"Precision,omitempty"`
	Recognition      string  `json:"Recognition,omitempty"`
	ResponseThumbURL string  `json:"ResponseThumbUrl,omitempty"` // 侯斯特特有
	ResponseURL      string  `json:"ResponseUrl,omitempty"`      // 侯斯特特有
}

type Component struct {
	ComponentAppid     string `json:"component_appid"`
	ComponentAppsecret string `json:"component_appsecret"`
	Name               string `json:"name"`
	Token              string `json:"token"`
	SymmetricKey       string `json:"symmetric_key"`
	Domain             string `json:"domain"`
	Dateline           int    `json:"dateline"`
	AlipayEmail        string `json:"alipay_email"`
	AlipayPartner      string `json:"alipay_partner"`
	AlipayKey          string `json:"alipay_key"`
}

type HostWxOther struct {
	Component Component `json:"component"`
	SYSId     string    `json:"SYS_last_message_id"`
}
