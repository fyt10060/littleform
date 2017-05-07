package controllers

import (
	"fmt"

	//	"io/ioutil"

	"littleform/model"

	"littleform/gopackages/types"

	"github.com/astaxie/beego"
	"github.com/mitchellh/mapstructure"
)

type ApiController struct {
	beego.Controller
}

//type ApiPost struct{} // get server
type ApiPost struct{} // post server

// 注册一个yar post server
func (this *ApiController) Post() {
	c := this.Ctx
	r := c.Request
	w := c.ResponseWriter

	register := RegisterList{
		List: []RegisterPair{
			{"create_form", "CreateForm"},
			{"delete_form", "DeleteForm"},
			{"save_form_config", "SaveFormConfig"},
			{"submit_form", "SubmitForm"},
			// get methods
			{"get_list", "GetList"},
			{"get_form_config", "GetFormConfig"},
			{"config_type", "ConfigType"},
			{"get_submit_list", "GetSubmitList"},
			{"get_submit_detail", "GetSubmitDetail"},
			{"delete_submit", "DeleteSubmit"},
		},
	}

	NewYarServer(&ApiPost{}, register, r, w)
}

// 注册一个yar get server
func (this *ApiController) Get() {
	c := this.Ctx
	r := c.Request
	w := c.ResponseWriter

	register := RegisterList{
		List: []RegisterPair{
			{"get_list", "GetList"},
			{"get_form_config", "GetFormConfig"},
			{"config_type", "ConfigType"},
			{"get_submit_list", "GetSubmitList"},
			{"get_submit_detail", "GetSubmitDetail"},
		},
	}

	NewYarServer(&ApiPost{}, register, r, w)
}

/** get function */
// 获取表单列表， get请求和post请求都会把对应的请求解json化
func (c *ApiPost) GetList(data map[string]interface{}) string {
	params, err := getParamsFromApiRpc(data)
	if err != nil {
		return model.ParseResult(model.ErrParamsError, nil)
	}
	var request model.ReGetFormList
	err = model.FormatJson(params.Get, &request)
	if err != nil {
		return model.ParseResult(model.ErrStructToMap, nil)
	}
	list, errs := model.GetFormList(request, params.Account.Appid, params.Component.ComponentAppid)
	if errs != model.ErrSuccess {
		return model.ParseResult(errs, nil)
	}
	return model.ParseResult(errs, list)
}

// 获取表单配置项
func (c *ApiPost) GetFormConfig(data map[string]interface{}) string {
	params, err := getParamsFromApiRpc(data)
	if err != nil {
		return model.ParseResult(model.ErrParamsError, nil)
	}
	var request model.ReGetFormConfig
	err = model.FormatJson(params.Get, &request)
	if err != nil {
		return model.ParseResult(model.ErrStructToMap, nil)
	}
	list, errs := model.GetFormConfig(request.FormId)
	if errs != model.ErrSuccess {
		return model.ParseResult(errs, nil)
	}
	return model.ParseResult(errs, list)
}

// 获取表单配置项可选类型
func (c *ApiPost) ConfigType(data map[string]interface{}) string {
	var configType = model.ConfigTypeResult{
		Type: []string{
			"input",
			"photo",
			"location",
			"select",
			"score",
		},
		InputType: []string{
			"text",
			"number",
			"tel",
			"email",
			"url",
			"range",
			"textarea",
			"date",
			"time",
			"radio",
			"checkbox",
			"color",
			"range",
		},
	}
	return model.ParseResult(model.ErrSuccess, configType)
}

// 获取提交id列表
func (c *ApiPost) GetSubmitList(data map[string]interface{}) string {
	params, err := getParamsFromApiRpc(data)
	if err != nil {
		return model.ParseResult(model.ErrParamsError, nil)
	}
	var request model.ReGetSubmitIdList
	err = model.FormatJson(params.Get, &request)
	if err != nil {
		return model.ParseResult(model.ErrStructToMap, nil)
	}
	list, errs := model.GetSubmitList(request.FormId, model.ToInt64(request.Count), model.ToInt64(request.Start))
	if errs != model.ErrSuccess {
		return model.ParseResult(errs, nil)
	}
	return model.ParseResult(errs, list)
}

// 根据提交id获取提交列表
func (c *ApiPost) GetSubmitDetail(data map[string]interface{}) string {
	params, err := getParamsFromApiRpc(data)
	if err != nil {
		return model.ParseResult(model.ErrParamsError, nil)
	}
	var request model.ReGetSubmitList
	err = model.FormatJson(params.Get, &request)
	if err != nil {
		return model.ParseResult(model.ErrStructToMap, nil)
	}
	list, errs := model.GetSubmitDetail(request.SubmitId)
	if errs != model.ErrSuccess {
		return model.ParseResult(errs, nil)
	}
	return model.ParseResult(errs, list)
}

/** post function */
// 创建新表单
func (c *ApiPost) CreateForm(data map[string]interface{}) string {
	params, err := getParamsFromApiRpc(data)
	if err != nil {
		return model.ParseResult(model.ErrParamsError, nil)
	}
	var request model.FormList
	err = model.FormatJson(params.Post, &request)
	fmt.Println("new form:", request)
	if err != nil {
		return model.ParseResult(model.ErrStructToMap, nil)
	}
	request.AccountAppid = params.Account.Appid
	request.ComponentAppid = params.Component.ComponentAppid
	result, errs := model.CreateForm(&request)
	return model.ParseResult(errs, result)
}

// 删除表单
func (c *ApiPost) DeleteForm(data map[string]interface{}) string {
	params, err := getParamsFromApiRpc(data)
	if err != nil {
		return model.ParseResult(model.ErrParamsError, nil)
	}
	var request model.ReGetFormConfig
	err = model.FormatJson(params.Post, &request)
	if err != nil {
		return model.ParseResult(model.ErrStructToMap, nil)
	}
	errs := model.DeleteForm(request.FormId)
	return model.ParseResult(errs, nil)
}

// 保存表单配置项
func (c *ApiPost) SaveFormConfig(data map[string]interface{}) string {
	params, err := getParamsFromApiRpc(data)
	if err != nil {
		return model.ParseResult(model.ErrParamsError, nil)
	}
	var request model.ReSaveFormList
	err = model.FormatJson(params.Post, &request)
	if err != nil {
		return model.ParseResult(model.ErrStructToMap, nil)
	}
	errs := model.SaveFormConfig(&request)
	return model.ParseResult(errs, nil)
}

// 提交表单
func (c *ApiPost) SubmitForm(data map[string]interface{}) string {
	params, err := getParamsFromApiRpc(data)
	if err != nil {
		return model.ParseResult(model.ErrParamsError, nil)
	}
	var request model.ReSubmitList
	err = model.FormatJson(params.Post, &request)
	if err != nil {
		return model.ParseResult(model.ErrStructToMap, nil)
	}
	errs := model.CreateSubmit(&request)
	return model.ParseResult(errs, nil)
}

// 删除表单
func (c *ApiPost) DeleteSubmit(data map[string]interface{}) string {
	params, err := getParamsFromApiRpc(data)
	if err != nil {
		return model.ParseResult(model.ErrParamsError, nil)
	}
	var request model.ReGetSubmitList
	err = model.FormatJson(params.Post, &request)
	if err != nil {
		return model.ParseResult(model.ErrStructToMap, nil)
	}
	errs := model.DeleteSubmit(request.SubmitId)
	return model.ParseResult(errs, nil)
}

func getParamsFromApiRpc(data map[string]interface{}) (result *types.AddonApiNormalParam, err error) {
	fmt.Println("request data is :")
	fmt.Println(data)
	var params types.AddonApiNormalParam
	errs := mapstructure.Decode(data, &params)
	fmt.Println(errs)
	return &params, errs
}
