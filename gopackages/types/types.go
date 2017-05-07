package types

import "encoding/json"

// Component 第三方信息
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

// User 侯斯特账户
type User struct {
	Id             int    `db:"id" json:"id" structs:"id"`
	Mobile         string `db:"mobile" json:"mobile" structs:"mobile"`
	Email          string `db:"email" json:"email" structs:"email"`
	System         int    `db:"system" json:"system" structs:"system"`
	Developer      int    `db:"developer" json:"developer" structs:"developer"`
	Status         int    `db:"status" json:"status" structs:"status"`
	Dateline       int    `db:"dateline" json:"dateline" structs:"dateline"`
	Postip         string `db:"postip" json:"postip" structs:"postip"`
	ComponentAppid string `db:"component_appid" json:"component_appid" structs:"component_appid"`
}

type AuthOriInfo map[string]interface{}

// AccountComponent 公众号第三方对应关系
type AccountComponent struct {
	AccountAppid   string `db:"account_appid" json:"account_appid"`
	ComponentAppid string `db:"component_appid" json:"component_appid"`
	Dateline       int    `db:"dateline" json:"dateline"`
}

// AccountUser 公众号侯斯特帐号对应关系
type AccountUser struct {
	AccountAppid string `db:"account_appid" json:"account_appid"`
	Uid          string `db:"uid" json:"uid"`
	Dateline     int    `db:"dateline" json:"dateline"`
}

// Account 公众号信息
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

// Addon 插件信息
type Addon struct {
	Id         int    `db:"id" json:"id"`
	Status     int    `db:"status" json:"status"`
	Priority   int    `db:"priority" json:"priority"`
	Name       string `db:"name" json:"name"`
	FolderName string `db:"folder_name" json:"folder_name"`
	Type       int    `db:"type" json:"type"`
	Insite     int    `db:"insite" json:"insite"`
	Url        string `db:"url" json:"url"`
}

type AddonApiParam struct {
	ComponentDomain string
	ComponentAppid  string
	Account         *Account
	Addon           *Addon
	User            *User
	Method          string
	Query           string
	RequestBody     interface{}
	RemoteIP        string
	Token           string
}

type AddonApiNormalParam struct {
	Component   Component   `json:"component"`
	Account     *Account    `json:"account"`
	Addon       *Addon      `json:"addon"`
	User        *User       `json:"user"`
	Method      string      `json:"method"`
	RequestBody interface{} `json:"request"`
	RemoteIP    string      `json:"remote_ip"`
	Token       string      `json:"token"`
	Get         interface{} `json:"get"`
	Post        interface{} `json:"post"`
}

// ApiResponse ApiResponse
type ApiResponse struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data,omitempty"`
}

// Token Token
type Token struct {
	AccessToken  string    `json:"token"`
	ExpiresIn    int64     `json:"expires"`
	RefreshToken string    `json:"refresh_token"`
	Data         TokenData `json:"data"`
}

// TokenData TokenData
type TokenData struct {
	Uid            int         `json:"uid"`
	Email          string      `json:"email"`
	Mobile         string      `json:"mobile"`
	Ip             string      `json:"ip"`
	System         int         `json:"system"`
	AccountList    []Account   `json:"account_list"`
	AccountId      int         `json:"account_id"`
	ExpiryDateline int         `json:"expiry_dateline"`
	Component      Component   `json:"component"`
	Addon          interface{} `json:"addon"`
}

func (this *Token) UnmarshalJSON(b []byte) error {

	type t struct {
		AccessToken  string  `json:"token"`
		ExpiresIn    float64 `json:"expires"`
		RefreshToken string  `json:"refresh_token"`
		Data         string  `json:"data"`
	}

	var token t

	err := json.Unmarshal(b, &token)

	if err != nil {
		return err
	}

	var d TokenData

	err = json.Unmarshal([]byte(token.Data), &d)

	if err != nil {
		return err
	}

	this.AccessToken = token.AccessToken
	this.ExpiresIn = int64(token.ExpiresIn)
	this.RefreshToken = token.RefreshToken
	this.Data = d

	return nil
}

type WeixinWebAuthUserInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        string   `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}
