package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JacobDenver007/UserSystemL0/common"
	"github.com/JacobDenver007/UserSystemL0/utils"
	gin "gopkg.in/gin-gonic/gin.v1"
)

// RegisterAPI 提供的API路由
func RegisterAPI(router *gin.Engine) {
	router.POST(fmt.Sprintf("/sendsms"), SendSmsHandler)
	router.POST(fmt.Sprintf("/signup"), SignUpByMobileHandler)
	router.POST(fmt.Sprintf("/suspenduser"), SuspendUserHandler)
	router.POST(fmt.Sprintf("/resetusername"), ResetUserNameHandler)
	router.POST(fmt.Sprintf("/resetuserpwd"), ResetUserPwdHandler)
	router.POST(fmt.Sprintf("/resetuserphone"), ResetUserPhoneHandler)

	router.POST(fmt.Sprintf("/createaccount"), CreateAccountHandler)
	router.POST(fmt.Sprintf("/suspendaccount"), SuspendAccountHandler)
	router.POST(fmt.Sprintf("/freezeaccount"), FreezeAccountHandler)
}

func SendSmsHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.SendSmsRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		if err := sendText(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func sendText(req *common.SendSmsRequest) error {
	if err := VailMobile(req.PhoneNum); err != nil {
		return err
	}
	code := MakeCode()
	if err := VailCode(code); err != nil {
		return err
	}
	if err := SendCode(req.PhoneNum, code); err != nil {
		return err
	}
	if err := DBClient.InsertToken(req.PhoneNum, code); err != nil {
		return err
	}
	return nil
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
			token.Phone = req.PhoneNum
			token.SignInCode = req.VerificationCode
			token.SignInTime = time.Now().Unix()
			token.Expire = token.SignInTime + 300

			tokenString := token.token()
			session := getsessions(c)
			session.Set(tokenString, token)

			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func createUser(req *common.SignUpByMobileRequest) error {
	if user, _ := DBClient.GetUserInfo(req.UserName); user != nil {
		return fmt.Errorf("user %s already exist", req.UserName)
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

func VerifiedCode(phoneNum string, verificationCode string) error {
	token, err := DBClient.GetTokenInfo(phoneNum)
	if err != nil {
		return err
	}
	if token.VerificationCode != verificationCode {
		return fmt.Errorf("wrong verification code")
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
		if err := suspendUser(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func suspendUser(req *common.SuspendUser) error {
	user, err := DBClient.GetUserInfo(req.UserName)
	if err != nil {
		return err
	}
	if req.OPCode == user.IsSuspended {
		return nil
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
		if err := resetUserName(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func resetUserName(req *common.ResetUserNameRequest) error {
	user, err := DBClient.GetUserInfo(req.OldUserName)
	if err != nil {
		return err
	}
	if validPwd := utils.ComparePassword(user.HashPwd, req.Pwd); !validPwd {
		return fmt.Errorf("wrong password")
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
		if err := resetUserPwd(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.ErrCode = common.OKCode
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
		return fmt.Errorf("wrong password")
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
		if err := resetUserPhone(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func resetUserPhone(req *common.ResetUserPhoneRequest) error {
	user, err := DBClient.GetUserInfo(req.UserName)
	if err != nil {
		return err
	}
	if validPwd := utils.ComparePassword(user.HashPwd, req.Pwd); !validPwd {
		return fmt.Errorf("wrong password")
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

func CreateAccountHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.ResetUserPhoneRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		address := CreateAccount()
		respone.Data = address
		respone.ErrCode = common.OKCode
	}
	c.JSON(http.StatusOK, respone)
}

func SuspendAccountHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.SuspendAccountRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		if err := suspendAccount(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func suspendAccount(req *common.SuspendAccountRequest) error {
	account, err := DBClient.GetAccountInfo(req.Address)
	if err != nil {
		return err
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
		if err := freezeAccount(req); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func freezeAccount(req *common.FreezeAccountRequest) error {
	account, err := DBClient.GetAccountInfo(req.Address)
	if err != nil {
		return err
	}
	account.IsFrozen = req.OPCode
	if err := DBClient.UpdateAccountInfo(account); err != nil {
		return err
	}
	return nil
}
