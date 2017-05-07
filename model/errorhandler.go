package model

import (
	"encoding/json"
	"fmt"
)

type ErrorType string

const (
	ErrSuccess ErrorType = "success"
	// data base error
	ErrSysDb          = "System database error"
	ErrOrmGetValues   = "get values from orm date failed"
	ErrMapToStruct    = "map to struct failed"
	ErrStructToMap    = "struct to map failed"
	ErrDeleteForm     = "delete form error"
	ErrAuthorityWrong = "you have no authority to check/update/delete this form"
	// api error
	ErrParamsError = "system params error"
	// form
	ErrNoMoreForm         = "no more form in db"
	ErrFormIdNotExist     = "can not found that form id"
	ErrFormDBToStruct     = "some error occured when read form data from db to golang struct"
	ErrCreateFormDB       = "some error occured when insert new form into db"
	ErrGetFormCountTooBig = "count should not over 100"
	ErrFormStillOnline    = "form still online"
	// config
	ErrDeleteConfig     = "delete config with id failed"
	ErrConfigListFailed = "get config list failed"
	ErrConfigIdNotFound = "can not found that config id"
	ErrCreateConfigDB   = "some error ouccred when create new config into db"
	ErrInputTypeDetail  = "the input detail you subscribe does not match the input type"
	// submit
	ErrCreateSubmitId    = "create submit id failed"
	ErrCreateSubmit      = "create submit detail failed"
	ErrChangeSubmitCount = "can not change submit count"
	ErrSubmitIdNotExist  = "submit id not found"
	ErrDeleteSubmitId    = "delete submit id failed"
	ErrDeleteSubmit      = "delete submit failed"
	ErrMustTypeEmpty     = "some config is must but is empty"
	ErrSubmitCountTooBig = "count should not over 50"
)

var (
	errCode = map[ErrorType]int{
		ErrSuccess:      0,
		ErrSysDb:        -126001,
		ErrMapToStruct:  -126002,
		ErrDeleteForm:   -126003,
		ErrOrmGetValues: -126004,
		ErrStructToMap:  -126005,

		ErrParamsError: 126000,

		// form
		ErrNoMoreForm:         126100,
		ErrFormIdNotExist:     126101,
		ErrFormDBToStruct:     126102,
		ErrCreateFormDB:       126103,
		ErrGetFormCountTooBig: 126104,
		ErrFormStillOnline:    126105,
		ErrAuthorityWrong:     126106,
		// config
		ErrDeleteConfig:     126200,
		ErrConfigListFailed: 126201,
		ErrConfigIdNotFound: 126202,
		ErrCreateConfigDB:   126203,
		ErrInputTypeDetail:  126204,
		// submit
		ErrCreateSubmit:      126300,
		ErrCreateSubmitId:    126301,
		ErrSubmitIdNotExist:  126302,
		ErrChangeSubmitCount: 126303,
		ErrMustTypeEmpty:     126305,
		ErrDeleteSubmitId:    126304,
		ErrDeleteSubmit:      126306,
		ErrSubmitCountTooBig: 126307,
	}
)

type Response struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  ErrorType   `json:"err_msg"`
	Data    interface{} `json:"data,omitempty"`
}

func ParseResult(errMsg ErrorType, data interface{}) string {
	fmt.Println("err type is: ", errMsg)
	code := errCode[errMsg]
	if code == 0 && errMsg != ErrSuccess {
		code = -1
	}
	r := Response{
		ErrCode: code,
		ErrMsg:  errMsg,
		Data:    data,
	}

	result, _ := json.Marshal(r)
	fmt.Println(string(result))
	return string(result)
}
