package model

import (
	"encoding/json"
	//	"fmt"
)

type ErrorType string

const (
	ErrSuccess ErrorType = "success"
	// data base error
	ErrSysDb = "System database error"
	
)

var (
	errCode = map[ErrorType]int{
		ErrSuccess:     0,
		ErrSysDb:       -1,
	}
)

type Response struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  ErrorType   `json:"err_msg"`
	Data    interface{} `json:"data"`
}

func ParseResult(errMsg ErrorType, data interface{}) []byte {
	code := errCode[errMsg]
	if code == 0 && errMsg != ErrSuccess {
		code = -1
	}
	r := Response{
		ErrCode: code,
		ErrMsg:  errMsg,
		Data:    data,
	}
	b, err := json.Marshal(r)
	if err != nil {

	}

	return b
}