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
	//UnauthorizedCode Login error
	UnauthorizedCode = 20004
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

type SuspendUser struct {
	UserName string `json:"username"`
	OPCode   int64  `json:"opcode"`
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

type CreateAccountRequest struct {
}

type SuspendAccountRequest struct {
	Address string `json:"address"`
}

type FreezeAccountRequest struct {
	Address string `json:"address"`
	OPCode  int64  `json:"opcode"`
}

type TokenInfo struct {
	PhoneNum         string `json:"phonenum"`
	VerificationCode string `json:"verificationcode"`
	TimeStamp        int64  `json:"timestamp"`
}

type User struct {
	ID          int64  `json:"id"`
	UserName    string `json:"username"`
	HashPwd     string `json:"hashpwd"`
	PhoneNum    string `json:"phonenum"`
	IsSuspended int64  `json:"issuspended"` //0代表正常，1代表注销
	Auth        int64  `json:"auth"`        //0代表无审批权，1代表有审批权
	IsApproved  int64  `json:"iapproved"`   //0代表未通过，1代表通过审批
}

type Account struct {
	Address     string `json:"address"`
	PrivateKey  string `json:"privatekey"`
	IsSuspended int64  `json:"issuspended"`
	IsFrozen    int64  `json:"isfrozen"`
}
