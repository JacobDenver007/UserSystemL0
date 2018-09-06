package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JacobDenver007/UserSystemL0/common"
	"github.com/JacobDenver007/UserSystemL0/utils"
	gin "gopkg.in/gin-gonic/gin.v1"
)

// RegisterAPI 提供的API路由
func RegisterAPI(router *gin.Engine) {
	router.POST(fmt.Sprintf("/sendsms"), SendSmsHandler)

	router.POST(fmt.Sprintf("/signup"), SignUpByMobileHandler)
	router.POST(fmt.Sprintf("/signin"), SignInHandler)

	router.POST(fmt.Sprintf("/suspenduser"), SuspendUserHandler)
	router.POST(fmt.Sprintf("/resetusername"), ResetUserNameHandler)
	router.POST(fmt.Sprintf("/resetuserpwd"), ResetUserPwdHandler)
	router.POST(fmt.Sprintf("/resetuserphone"), ResetUserPhoneHandler)
	router.POST(fmt.Sprintf("/getuserinfo"), GetUserInfoHandler)

	router.POST(fmt.Sprintf("/createaccount"), CreateAccountHandler)
	router.POST(fmt.Sprintf("/suspendaccount"), SuspendAccountHandler)
	router.POST(fmt.Sprintf("/freezeaccount"), FreezeAccountHandler)
	router.POST(fmt.Sprintf("/getuseraccount"), GetUserAccountHandler)

	router.POST(fmt.Sprintf("/sendtransaction"), SendTransactionHandler)
}

func getToken(c *gin.Context) *Token {
	h := c.GetHeader("Authorization")
	t := session.Get(strings.TrimPrefix(h, "Bearer "))

	token := t.(*Token)
	return token
}

func checkUser(userName string) error {
	user, err := DBClient.GetUserInfo(userName)
	if err != nil {
		return err
	}
	if user.IsSuspended == 1 {
		return fmt.Errorf("用户 %s 被注销", userName)
	}
	if user.IsApproved == 0 {
		return fmt.Errorf("用户 %s 未经过审批", userName)
	}
	return nil
}

func checkAccount(userName string, address string) error {
	account, err := DBClient.GetAccountInfo(address)
	if err != nil {
		return err
	}
	if account.IsSuspended == 1 {
		return fmt.Errorf("账户 %s 被注销", address)
	}
	if account.IsFrozen == 1 {
		return fmt.Errorf("账户 %s 被冻结", address)
	}
	if account.User != userName && userName != "" {
		return fmt.Errorf("用户 %s 无权限使用账户 %s", userName, address)
	}
	return nil
}

func SendSmsHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.SendSmsRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		if code, err := sendText(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.Data = fmt.Sprintf("发送成功，验证码为 %s", code)
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func sendText(req *common.SendSmsRequest) (string, error) {
	if err := VailMobile(req.PhoneNum); err != nil {
		return "", err
	}
	code := MakeCode()
	if err := VailCode(code); err != nil {
		return "", err
	}
	if err := SendCode(req.PhoneNum, code); err != nil {
		return "", err
	}
	if err := DBClient.InsertToken(req.PhoneNum, code); err != nil {
		return "", err
	}
	return code, nil
}

func SignUpByMobileHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.SignUpByMobileRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		if err := createUser(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			token := &Token{}
			token.UserName = req.UserName
			token.SignInTime = time.Now().Unix()
			token.Expire = time.Now().Add(600 * time.Second).Unix()

			tokenString := token.token()
			session := getsessions(c)
			session.Set(tokenString, token)

			data := &common.SignUpByMobileResponse{Auth: tokenString}
			respone.Data = data
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func createUser(req *common.SignUpByMobileRequest) error {
	if user, _ := DBClient.GetUserInfo(req.UserName); user != nil {
		return fmt.Errorf("用户 %s 已经存在", req.UserName)
	}
	if err := VerifiedCode(req.PhoneNum, req.VerificationCode); err != nil {
		return err
	}
	hashPwd := utils.HashPassword(req.Pwd)
	user := &common.User{
		UserName: req.UserName,
		HashPwd:  hashPwd,
		PhoneNum: req.PhoneNum,
	}
	if err := DBClient.InsertUser(user); err != nil {
		return err
	}

	return nil
}

func SignInHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.SignInRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		if user, err := signIn(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			token := &Token{}
			token.UserName = user.UserName
			token.SignInTime = time.Now().Unix()
			token.Expire = time.Now().Add(600 * time.Second).Unix()

			tokenString := token.token()
			session := getsessions(c)
			session.Set(tokenString, token)

			data := &common.SignUpByMobileResponse{Auth: tokenString}
			respone.Data = data
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func signIn(req *common.SignInRequest) (*common.User, error) {
	user, err := DBClient.GetUserInfo(req.UserName)
	if err != nil {
		return nil, err
	}
	if !utils.ComparePassword(user.HashPwd, req.Pwd) {
		return nil, fmt.Errorf("密码错误")
	}
	return user, nil
}

func VerifiedCode(phoneNum string, verificationCode string) error {
	token, err := DBClient.GetTokenInfo(phoneNum)
	if err != nil {
		return err
	}
	if token.VerificationCode != verificationCode {
		return fmt.Errorf("验证码错误")
	}
	return nil
}

func SuspendUserHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.SuspendUser{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		token := getToken(c)
		if req.UserName != token.UserName && token.UserName != "admin" {
			respone.ErrCode = common.UnauthorizedCode
			respone.ErrMsg = "无权限"
		} else {
			if err := suspendUser(req); err != nil {
				respone.ErrCode = common.ExecuteCode
				respone.ErrMsg = err.Error()
			} else {
				respone.Data = "操作成功"
				respone.ErrCode = common.OKCode
			}
		}
	}
	c.JSON(http.StatusOK, respone)
}

func suspendUser(req *common.SuspendUser) error {
	if req.OPCode != 0 && req.OPCode != 1 {
		return fmt.Errorf("错误的操作码")
	}
	user, err := DBClient.GetUserInfo(req.UserName)
	if err != nil {
		return err
	}
	user.IsSuspended = req.OPCode
	if err := DBClient.UpdateUserInfo(user); err != nil {
		return err
	}
	return nil
}

func ResetUserNameHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.ResetUserNameRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		token := getToken(c)
		if req.OldUserName != token.UserName {
			respone.ErrCode = common.UnauthorizedCode
			respone.ErrMsg = "无权限"
		} else {
			if err := resetUserName(req); err != nil {
				respone.ErrCode = common.ExecuteCode
				respone.ErrMsg = err.Error()
			} else {
				respone.Data = "重置用户名成功，请重新登录"
				respone.ErrCode = common.OKCode
			}
		}
	}
	c.JSON(http.StatusOK, respone)
}

func resetUserName(req *common.ResetUserNameRequest) error {
	user, err := DBClient.GetUserInfo(req.OldUserName)
	if err != nil {
		return err
	}
	user.UserName = req.NewUserName
	if err := DBClient.UpdateUserInfo(user); err != nil {
		return err
	}
	return nil
}

func ResetUserPwdHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.ResetUserPwdRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		token := getToken(c)
		if req.UserName != token.UserName {
			respone.ErrCode = common.UnauthorizedCode
			respone.ErrMsg = "无权限"
		} else {
			if err := resetUserPwd(req); err != nil {
				respone.ErrCode = common.ExecuteCode
				respone.ErrMsg = err.Error()
			} else {
				respone.Data = "重置用户密码成功"
				respone.ErrCode = common.OKCode
			}
		}
	}
	c.JSON(http.StatusOK, respone)
}

func resetUserPwd(req *common.ResetUserPwdRequest) error {
	user, err := DBClient.GetUserInfo(req.UserName)
	if err != nil {
		return err
	}
	if validPwd := utils.ComparePassword(user.HashPwd, req.OldPwd); !validPwd {
		return fmt.Errorf("旧密码错误")
	}
	user.HashPwd = utils.HashPassword(req.NewPwd)
	if err := DBClient.UpdateUserInfo(user); err != nil {
		return err
	}
	return nil
}

func ResetUserPhoneHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.ResetUserPhoneRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		token := getToken(c)
		if req.UserName != token.UserName {
			respone.ErrCode = common.UnauthorizedCode
			respone.ErrMsg = "无权限"
		} else {
			if err := resetUserPhone(req); err != nil {
				respone.ErrCode = common.ExecuteCode
				respone.ErrMsg = err.Error()
			} else {
				respone.Data = "重置手机号成功"
				respone.ErrCode = common.OKCode
			}
		}
	}
	c.JSON(http.StatusOK, respone)
}

func resetUserPhone(req *common.ResetUserPhoneRequest) error {
	user, err := DBClient.GetUserInfo(req.UserName)
	if err != nil {
		return err
	}
	if err := VerifiedCode(req.NewPhoneNum, req.VerificationCode); err != nil {
		return err
	}
	user.PhoneNum = req.NewPhoneNum
	if err := DBClient.UpdateUserInfo(user); err != nil {
		return err
	}
	return nil
}

func GetUserInfoHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.GetUserInfoRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		if user, err := DBClient.GetUserInfo(req.UserName); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			data := &common.GetUserInfoResponse{}
			if user.IsSuspended == 1 {
				data.IsSuspended = "用户被注销"
			} else {
				data.IsSuspended = "正常"
			}
			if user.IsApproved == 1 {
				data.IsApproved = "已通过审批"
			} else {
				data.IsApproved = "未通过审批"
			}
			data.ID = user.ID
			data.UserName = user.UserName
			data.PhoneNum = user.PhoneNum

			respone.Data = data
			respone.ErrCode = common.OKCode
		}

	}
	c.JSON(http.StatusOK, respone)
}

func CreateAccountHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.CreateAccountRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		token := getToken(c)
		if response, err := createAccount(token.UserName, req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.Data = response
			respone.ErrCode = common.OKCode
		}

	}
	c.JSON(http.StatusOK, respone)
}

func createAccount(user string, req *common.CreateAccountRequest) (*common.CreateAccountResponse, error) {
	response := &common.CreateAccountResponse{}

	if req.PrivateKey == "" {
		address, privateKey := GernerateAccount()
		response.Address = address
		response.Hex = privateKey
		err := DBClient.InsertAccount(user, address, privateKey)
		return response, err
	} else {
		address, err := GenerateAccountByPrivateKey(req.PrivateKey)
		if err != nil {
			return nil, err
		}
		account, _ := DBClient.GetAccountInfo(address)
		if account == nil {
			response.Address = address
			response.Hex = req.PrivateKey
			err := DBClient.InsertAccount(user, address, req.PrivateKey)
			return response, err
		} else {
			return nil, fmt.Errorf("账户已存在，请勿重复创建")
		}
	}
}

func SuspendAccountHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.SuspendAccountRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		token := getToken(c)
		if err := suspendAccount(token.UserName, req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.Data = fmt.Sprintf("注销账户 %s 成功", req.Address)
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func suspendAccount(user string, req *common.SuspendAccountRequest) error {
	account, err := DBClient.GetAccountInfo(req.Address)
	if err != nil {
		return err
	}
	if user != account.User {
		return fmt.Errorf("此账户非当前用户创建，用户无权操作此账户")
	}
	account.IsSuspended = 1
	if err := DBClient.UpdateAccountInfo(account); err != nil {
		return err
	}
	return nil
}

func FreezeAccountHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.FreezeAccountRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		token := getToken(c)
		if err := freezeAccount(token.UserName, req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.Data = "操作成功"
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func freezeAccount(user string, req *common.FreezeAccountRequest) error {
	account, err := DBClient.GetAccountInfo(req.Address)
	if err != nil {
		return err
	}
	if user != account.User {
		return fmt.Errorf("此账户非当前用户创建，用户无权操作此账户")
	}
	account.IsFrozen = req.OPCode
	if err := DBClient.UpdateAccountInfo(account); err != nil {
		return err
	}
	return nil
}

func GetUserAccountHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.GetUserAccountRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		if accounts, err := DBClient.GetUserAccount(req.UserName); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.Data = accounts
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func SendTransactionHandler(c *gin.Context) {
	var respone common.APIRespone

	req := &common.SendTransactionRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		token := getToken(c)
		if err := checkSendTxOption(token.UserName, req); err != nil {
			respone.ErrCode = common.UnauthorizedCode
			respone.ErrMsg = err.Error()
		} else {
			if err := RPCClient.SendTransaction(req.From, req.To, req.AssetID, req.Value); err != nil {
				respone.ErrCode = common.ExecuteCode
				respone.ErrMsg = err.Error()
			} else {
				respone.Data = "发送成功"
				respone.ErrCode = common.OKCode
			}
		}
	}
	c.JSON(http.StatusOK, respone)
}

func checkSendTxOption(user string, req *common.SendTransactionRequest) error {
	if err := checkUser(user); err != nil {
		return err
	}
	if err := checkAccount(user, req.From); err != nil {
		return err
	}
	if err := checkAccount("", req.To); err != nil {
		return err
	}
	return nil
}
