// model
package model

import (
	"crypto/sha1"
	"fmt"
	//	"regexp" // 正则表达式包
	"sort"
	"strconv"
	"strings"
	//	"database/sql"
	"encoding/json"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mitchellh/mapstructure"
)

type ConfigTypeName string

const (
	CTInput    ConfigTypeName = "input"
	CTPhoto                   = "photo"
	CTLocation                = "location"
	CTSelect                  = "select"
	CTScore                   = "score"
)

type InputType string

const (
	ITText     InputType = "text"
	ITNumber             = "number"
	ITTel                = "tel"
	ITEmail              = "email"
	ITUrl                = "url"
	ITRange              = "range"
	ITTextArea           = "textarea"
	ITDate               = "date"
	ITTime               = "time"
	ITRadio              = "radio"
	ITCheckbox           = "checkbox"
	ITColor              = "color"
)

const (
	// table name
	tbFormList     = "form_list"
	tbConfigList   = "config_list"
	tbSubmitList   = "submit_list"
	tbSubmitIdList = "submit_id_list"
	// serect name
	serForm   = "slkjdfasdf"
	serConfig = "sdafsdf"
	serSubmit = "wedvadsf"
	// normal const
	submitSplitString  = ",|"
	formListCountMax   = 100
	submitListCountMax = 50
)

// mysql model
type FormList struct {
	FormId          string `json:"id" orm:"PK"`
	ComponentAppid  string `json:"omitempty" orm:"index"`
	AccountAppid    string `json:"omitempty" orm:"index"`
	FormTitle       string `json:"title" orm:"size(21)"`
	FormDescription string `json:"description" orm:"size(16)"`
	SubmitCount     int    `json:"submit_count"`
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
	CreateTime      int64  `json:"create_time"`
}

type ConfigList struct {
	ConfigId    string         `orm:"PK" orm:"char(40)" json:"config_id"`
	FormId      string         `orm:"index" orm:"char(40)" json:"form_id"`
	Type        ConfigTypeName `json:"config_type"`
	InputType   InputType      `json:"input_type"`
	Title       string         `json:"title"`
	Order       int            `orm:"" json:"order"`
	InputDetail string         `orm:"null" json:"input_detail"`
	MustType    int            `json:"must_type"`
}

type SubmitIdList struct {
	SubmitId   string `orm:"PK" orm:"char(40)" json:"submit_id"`
	FormId     string `orm:"index" orm:"char(40)" json:"form_id"`
	SubmitTime int64  `json:"submit_time"`
}

type SubmitList struct {
	Id         int
	Order      int            `json:"order"`
	SubmitId   string         `orm:"index" orm:"char(40)"`
	ConfigType ConfigTypeName `json:"config_type"`
	InputType  InputType      `json:"input_type"`
	Content    string         `json:"content"`
	ConfigId   string         `json:"config_id"`
	MustType   int            `json:"must_type"`
}

/* response model */

type FormListResult struct {
	List []FormList `json:"form_list"`
}

type ConfigListResult struct {
	List []ConfigList `json:"config_list"`
}

type SaveConfigResult struct {
	FailedList []string `json:"failed_list"`
}

type ConfigTypeResult struct {
	Type      []string `json:"config_type"`
	InputType []string `json:"input_type"`
}

/* request model */
type ReGetFormList struct {
	Count string `json:"count"`
	Start string `json:"start"`
}

type ReGetFormConfig struct {
	FormId string `json:"form_id"`
}

type ReGetSubmitIdList struct {
	FormId string `json:"form_id"`
	Count  string `json:"count"`
	Start  string `json:"start"`
}

type ReGetSubmitList struct {
	SubmitId string `json:"submit_id"`
}

// 收到的保存表单配置项请求体
type ReSaveFormList struct {
	List       []ConfigList `json:"config_list"`
	FormId     string       `json:"form_id"`
	DeleteList []string     `json:"delete_list"`
}

type ReSubmit struct {
	Order      int            `json:"order"`
	Content    []string       `json:"content"`
	ConfigId   string         `json:"config_id"`
	ConfigType ConfigTypeName `json:"config_type"`
	InputType  InputType      `json:"input_type"`
	MustType   string         `json:"must_type"`
}

// 收到的提交请求体
type ReSubmitList struct {
	FormId string     `json:"form_id"`
	List   []ReSubmit `json:"submit"`
}

// 数据库地址
const (
	mysqlAddr  = "burnlog:1234567burnlog@/burnlogDB?charset=utf8"
	mysqlLocal = "burnlog:1234567burnlog@/burnlogDB?charset=utf8"
	mysqlRoot  = "root:5224733mysql@/burnlogDB?charset=utf8"
)

// 初始化mysql数据库
func init() {
	err := orm.RegisterDataBase("default", "mysql", mysqlRoot, 30)
	if err != nil {
		panic(err)
	}

	orm.RegisterModel(new(FormList), new(ConfigList), new(SubmitIdList), new(SubmitList))
	orm.RunSyncdb("default", false, true)
}

// 获取表单列表
func GetFormList(data ReGetFormList, accountAppId, componentAppId string) (res *[]FormList, errs ErrorType) {
	count := ToInt64(data.Count)
	if count > formListCountMax {
		return nil, ErrGetFormCountTooBig
	}
	start := ToInt64(data.Start)
	o := orm.NewOrm()
	qs := o.QueryTable(tbFormList)
	qs.Filter("component_appid", componentAppId).Filter("account_appid", accountAppId)
	listLength, err := qs.Count()
	// start 容错
	if start < 0 {
		start = 0
	}
	// 判断是否超过已有数据量
	if start > listLength {
		return nil, ErrNoMoreForm
	}
	// 根据已有信息进行数据偏移
	qs.Limit(count, start)
	var result []FormList
	var resultValue FormList
	var maps []orm.Params
	_, err = qs.Values(&maps)
	if err != nil {
		return nil, ErrOrmGetValues
	}
	for _, v := range maps {
		err = mapstructure.Decode(v, &resultValue)
		if err != nil {
			return nil, ErrFormDBToStruct
		}
		result = append(result, resultValue)
	}
	return &result, ErrSuccess
}

// 创建表单
func CreateForm(form *FormList) (formData *FormList, errs ErrorType) {
	now := getNowTime()
	form.CreateTime = now
	// 生成表单id
	form.FormId = getFormatId(form.FormTitle, serForm, now)

	o := orm.NewOrm()
	_, err := o.Insert(form)

	if err != nil {
		fmt.Println(err)
		return nil, ErrCreateFormDB
	}
	return form, ErrSuccess
}

// 删除表单
func DeleteForm(formId string) ErrorType {
	form := FormList{FormId: formId}
	o := orm.NewOrm()
	err := o.Read(&form)
	// err != nil 说明读取表单出错，可能是不存在相应的表单或数据库出错
	if err != nil {
		return ErrFormIdNotExist
	} else {
		if form.EndTime > getNowTime() {
			return ErrFormStillOnline
		}
		_, err := o.Delete(&form)
		if err != nil {
			return ErrDeleteForm
		} else {
			return ErrSuccess
		}
	}
}

// 获取某个表单的配置项
func GetFormConfig(formId string) (res *[]ConfigList, errs ErrorType) {
	o := orm.NewOrm()
	qs := o.QueryTable(tbConfigList)
	// 根据form_id进行query
	qs.Filter("form_id", formId)
	var maps []orm.Params
	_, err := qs.Values(&maps)
	if err != nil {
		return nil, ErrOrmGetValues
	}
	var result []ConfigList
	var resultValue ConfigList
	for _, v := range maps {
		err = mapstructure.Decode(v, &resultValue)
		if err != nil {
			return nil, ErrMapToStruct
		}
		result = append(result, resultValue)
	}
	return &result, ErrSuccess
}

// 保存一个表单的所有配置项
func SaveFormConfig(data *ReSaveFormList) ErrorType {
	o := orm.NewOrm()
	now := getNowTime()
	list := data.List
	formId := data.FormId
	deleteList := data.DeleteList
	// 根据前端上传的list删除config
	for _, v := range deleteList {
		config := ConfigList{ConfigId: v}
		err := o.Read(&config)
		if err != nil {
			return ErrDeleteConfig
		} else {
			_, err = o.Delete(&config)
			if err != nil {
				return ErrDeleteConfig
			}
		}
	}
	// 存储配置项
	for _, v := range list {
		v.FormId = formId
		// 如果config id存在，则更新现有config，如果为空则插入到表中
		if v.ConfigId == "" {
			v.ConfigId = getFormatId(v.Title, serConfig, now)
			_, err := o.Insert(&v)
			if err != nil {
				return ErrCreateConfigDB
			}
		} else {
			_, err := o.Update(&v)
			if err != nil {
				return ErrConfigIdNotFound
			}
		}
	}

	return ErrSuccess
}

// 获取提交信息列表
func GetSubmitList(formId string, count, start int64) (res *[]SubmitIdList, errs ErrorType) {
	if count > submitListCountMax {
		return nil, ErrSubmitCountTooBig
	}
	o := orm.NewOrm()
	qs := o.QueryTable(tbSubmitIdList).Filter("form_id", formId)
	listLength, err := qs.Count()
	// start容错
	if start < 0 {
		start = 0
	}
	// 判断要求的是否超限
	if start > listLength {
		return nil, ErrNoMoreForm
	}
	// 根据要求做数据偏移
	qs.Limit(count, start)
	var maps []orm.Params
	_, err = qs.Values(&maps)
	if err != nil {
		return nil, ErrOrmGetValues
	}
	// 获取提交id列表
	var result []SubmitIdList
	var resultValue SubmitIdList
	for _, v := range maps {
		err = mapstructure.Decode(v, &resultValue)
		if err != nil {
			return nil, ErrMapToStruct
		}
		result = append(result, resultValue)
	}

	return &result, ErrSuccess
}

// 根据id获取提交详情
func GetSubmitDetail(submitId string) (res *[]ReSubmit, errs ErrorType) {
	o := orm.NewOrm()
	// 根据提交id进行query
	qs := o.QueryTable(tbSubmitList).Filter("submit_id", submitId)
	var maps []orm.Params
	_, err := qs.Values(&maps)
	if err != nil {
		return nil, ErrOrmGetValues
	}
	var list []SubmitList
	err = mapstructure.Decode(maps, &list)
	if err != nil {
		return nil, ErrMapToStruct
	}
	// 为了将保存时候的string转换为slice，这一步感觉耗费巨大，以后应考虑改进
	var result []ReSubmit
	for _, v := range list {
		var value = ReSubmit{
			Order:      v.Order,
			Content:    strings.Split(v.Content, submitSplitString),
			ConfigId:   v.ConfigId,
			ConfigType: v.ConfigType,
			InputType:  v.InputType,
			MustType:   strconv.Itoa(v.MustType),
		}
		result = append(result, value)
	}

	return &result, ErrSuccess
}

// 保存新提交
func CreateSubmit(data *ReSubmitList) ErrorType {
	o := orm.NewOrm()
	now := getNowTime()
	list := data.List
	formId := data.FormId
	// 保存submit 详情，先查看是否能够保存各个详情，因为可能出现必填项为空，或者content和对应的config type及input type无法匹配
	var theError = ErrSuccess
	for _, v := range list {
		// 检查必填项是否为空，可能是slice为空或slice不为空但slice元素都为空
		if v.MustType == "1" {
			if len(v.Content) == 0 {
				theError = ErrMustTypeEmpty
				break
			} else {
				var hasContent bool = false
				for _, value := range v.Content {
					if value != "" {
						hasContent = true
						break
					}
				}
				if !hasContent {
					theError = ErrMustTypeEmpty
					break
				}
			}
		}

	}
	if theError != ErrSuccess {
		return theError
	}
	// 将前端上传的数组转为string,同上个API一样，个人认为耗费巨大
	submitId := getFormatId(formId, serSubmit, now)
	var errs error
	for _, v := range list {
		var submit SubmitList
		submit.ConfigId = v.ConfigId
		submit.ConfigType = v.ConfigType
		submit.InputType = v.InputType
		submit.SubmitId = submitId
		submit.Order = v.Order
		submit.MustType = ToInt(v.MustType)
		submit.Content = strings.Join(v.Content, submitSplitString)
		_, errs = o.Insert(&submit)
	}

	if errs != nil {
		return ErrCreateSubmit
	}
	// 为这个表单的填写人数+1
	var form = FormList{FormId: formId}
	err := o.Read(&form)
	if err != nil {
		return ErrFormIdNotExist
	}
	form.SubmitCount += 1
	_, err = o.Update(&form)
	if err != nil {
		return ErrChangeSubmitCount
	}
	// 最终保存submit id
	var newSub = SubmitIdList{
		SubmitId:   submitId,
		FormId:     formId,
		SubmitTime: now,
	}
	_, err = o.Insert(&newSub)
	if err != nil {
		return ErrCreateSubmitId
	}
	return ErrSuccess
}

// 删除提交
func DeleteSubmit(submitId string) ErrorType {
	o := orm.NewOrm()
	var submit = SubmitIdList{SubmitId: submitId}
	// 先查找有无该提交
	err := o.Read(&submit)
	if err != nil {
		return ErrSubmitIdNotExist
	}
	// 删除提交id
	_, err = o.Delete(&submit)
	if err != nil {
		return ErrDeleteSubmitId
	}
	// 删除提交详情
	qs := o.QueryTable(tbSubmitList).Filter("submit_id", submitId)
	_, err = qs.Delete()
	if err != nil {
		return ErrDeleteSubmit
	}
	return ErrSuccess
}

// json转换为对应struct
func FormatJson(jsonData interface{}, class interface{}) error {
	fmt.Println("json params is:", jsonData)
	stringData := jsonData.(string)
	byteData := []byte(stringData)
	err := json.Unmarshal(byteData, class)
	if err != nil {
		fmt.Println("format json err:", err)
	}

	return err
}

// 获取当前时间
func getNowTime() int64 {
	return time.Now().Unix()
}

// 各种id生成器
func getFormatId(value, serect string, time int64) string {
	theTime := fmt.Sprintf("%d", time)
	array := []string{value, theTime, serect}
	sort.Strings(array)

	idString := strings.Join(array, "")

	shId := sha1.New()
	shId.Write([]byte(idString))

	shByte := shId.Sum(nil)
	id := fmt.Sprintf("%x", shByte)

	return id
}

func ToInt(i string) int {
	re, _ := strconv.Atoi(i)
	return re
}

func ToInt64(i string) int64 {
	re, _ := strconv.ParseInt(i, 10, 64)
	return re
}
