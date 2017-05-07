package addon

import (
	"encoding/json"
	"reflect"
)

type RpcConfig struct {
	//yar
	MagicNum  string `db:"magic_num" json:"magic_num"`
	EncyptKey string `db:"encypt_key" json:"encypt_key"`
	Packager  string `db:"packager" json:"packager"`
}
type AppConfigItem struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type AppConfig []AppConfigItem

func (c *AppConfig) MarshalDB() ([]byte, error) {
	return json.Marshal(c)
}

func (c *AppConfig) UnmarshalDB(data []byte) error {
	return json.Unmarshal(data, c)
}

type AddonConfig struct {
	//addon
	MpUrl        string     `json:"mp_url,omitempty"`
	WxUrl        string     `json:"wx_url,omitempty"`
	ApiUrl       string     `json:"api_url,omitempty"`
	CronUrl      string     `json:"cron_url,omitempty"`
	SystemApiUrl string     `json:"system_api_url,omitempty"`
	AppConfig    *AppConfig `json:"app_config,omitempty" db:"app_config"`
	//wxa
	RequestUrl string `json:"request_url,omitempty"`
	SocketUrl  string `json:"socket_url,omitempty"`
	UploadUrl  string `json:"upload_url,omitempty"`
}

type NetConfig struct {
	AddonKind   string       `json:"addon_kind"`
	RpcKind     string       `json:"rpc_kind"`
	AddonConfig *AddonConfig `json:"addon_config,omitempty"`
	RpcConfig   *RpcConfig   `json:"rpc_config,omitempty"`
}
type AllNetConfig map[string]NetConfig

func (r *AllNetConfig) MarshalDB() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AllNetConfig) UnmarshalDB(data []byte) error {
	return json.Unmarshal(data, r)
}

type Addon struct {
	AppName    string        `db:"app_name" json:"app_name"`
	Kind       string        `db:"kind" json:"kind"`
	Remark     string        `db:"remark" json:"remark"`
	AddonType  int           `db:"addon_type" json:"addon_type"`
	ForSale    int           `db:"for_sale" json:"for_sale"`
	Priority   int           `db:"priority" json:"priority"`
	Appid      string        `db:"appid" json:"appid"`
	Appsecret  string        `db:"appsecret" json:"appsecret"`
	Status     int           `db:"status" json:"status"`
	CreateTime int64         `db:"createtime" json:"createtime"`
	Scheme     string        `db:"scheme" json:"scheme"`
	NetConfig  *AllNetConfig `db:"net_config" json:"net_config"`
}

type AddonAccount struct {
	AccountAppid   string `db:"account_appid" json:"account_appid"`
	ComponentAppid string `db:"component_appid" json:"component_appid"`
	AppName        string `db:"app_name" json:"app_name"`
	ExpireDateline int64  `db:"expire_dateline" json:"expire_dateline"`
	Status         int    `db:"status" json:"status"`
}

const (
	AddonComponentdConditionKindAppointAppid    string = "appoint_appid"      //只允许指定appid使用
	AddonComponentConditionKindAppointAppidList string = "appoint_appid_list" //只允许指定appid 列表使用
	AddonComponentConditionKindAllowSale        string = "allow_sale"         //是否允许上架出售
	AddonComponentConditionKindAllowUse         string = "allow_use"          //是否允许使用
	AddonComponentConditionKindMaxAccountNum    string = "max_account_num"    //最大开通数量
)

type Condition map[string]interface{}

func (c *Condition) MarshalDB() ([]byte, error) {
	return json.Marshal(c)
}
func (c *Condition) UnmarshalDB(data []byte) error {
	return json.Unmarshal(data, c)
}

type AddonComponent struct {
	ComponentAppid string     `json:"component_appid" db:"component_appid"`
	AppName        string     `json:"app_name" db:"app_name"`
	Condition      *Condition `json:"condition" db:"condition"`
	CreateTime     int64      `json:"create_time" db:"create_time"`
	LastUpdateTime int64      `json:"last_update_time" db:"last_update_time"`
}

type Hook struct {
	Id          int64  `db:"id" json:"id"`
	AppName     string `db:"app_name" json:"app_name"`
	Hook        string `db:"hook" json:"hook"`
	Action      string `db:"action" json:"action"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Addon       Addon  `db:"-" json:"addon"`
}

type Button struct {
	Id          int64  `db:"id" json:"id"`
	AppName     string `db:"app_name" json:"app_name"`
	Button      string `db:"button" json:"button"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

type Pattern struct {
	Id          int64  `db:"id" json:"id"`
	AppName     string `db:"app_name" json:"app_name"`
	Pattern     string `db:"pattern" json:"pattern"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

type Keyword struct {
	Id          int64  `db:"id" json:"id"`
	AppName     string `db:"app_name" json:"app_name"`
	Keyword     string `db:"keyword" json:"keyword"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

func (this *Addon) GetUrl(urlType string, functionType string) string {

	var addonUrl string

	if this.NetConfig != nil {

		if _, ok := (*this.NetConfig)[urlType]; ok {

			if (*this.NetConfig)[urlType].AddonConfig != nil {
				addonUrl = getField((*this.NetConfig)[urlType].AddonConfig, functionType)
			}
		}
	}

	return addonUrl
}

func getField(v *AddonConfig, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)

	if f.IsValid() {
		return f.String()
	} else {
		return ""
	}
}
