package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

var (
	// OKCode Ok
	OKCode = 0
	// ParamCode unmarshal json error
	ParamCode = 20001
	//RequestCode Api request error
	RequestCode = 20002
	//ExecuteCode Api execute error
	ExecuteCode = 20003
)

//APIRespone Api respone data
type APIRespone struct {
	Data    interface{} `json:"data"`
	ErrCode int         `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
}

//MD5 calc hash
func (api *APIRespone) MD5() string {
	bts, _ := json.Marshal(api.Data)
	r := md5.Sum(bts)
	return hex.EncodeToString(r[:])
}

type SendSmsRequest struct {
	PhoneNum string `json:"phonenum"`
}

type SignUpByMobileRequest struct {
	UserName         string `json:"username"`
	Pwd              string `json:"pwd"`
	PhoneNum         string `json:"phonenum"`
	VerificationCode string `json:"verificationcode"`
}

type ResetUserNameRequest struct {
	OldUserName string `json:"oldusername"`
	Pwd         string `json:"pwd"`
	PhoneNum    string `json:"phonenum"`
	NewUserName string `json:"newusername"`
}

type ResetUserPwdRequest struct {
	UserName string `json:"username"`
	OldPwd   string `json:"oldpwd"`
	PhoneNum string `json:"phonenum"`
	NewPwd   string `json:"newpwd"`
}

type ResetUserPhoneRequest struct {
	UserName         string `json:"username"`
	Pwd              string `json:"pwd"`
	OldPhoneNum      string `json:"oldphonenum"`
	NewPhoneNum      string `json:"newphonenum"`
	VerificationCode string `json:"verificationcode"`
}

type TokenInfo struct {
	PhoneNum         string `json:"phonenum"`
	VerificationCode string `json:"verificationcode"`
	TimeStamp        string `json:"timestamp"`
}

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	HashPwd  string `json:"hashpwd"`
	PhoneNum string `json:"phonenum"`
}
