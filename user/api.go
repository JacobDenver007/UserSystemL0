package user

import (
	"fmt"
	"net/http"

	"github.com/UserSystemL0/common"
	"github.com/UserSystemL0/utils"
	gin "gopkg.in/gin-gonic/gin.v1"
)

// RegisterAPI 提供的API路由
func RegisterAPI(router *gin.Engine) {
	router.POST(fmt.Sprintf("/sendsms"), SendSmsHandler)
	router.POST(fmt.Sprintf("/signupbymobile"), SignUpByMobileHandler)
	router.POST(fmt.Sprintf("/resetusername"), ResetUserNameHandler)
	router.POST(fmt.Sprintf("/resetuserpwd"), ResetUserPwdHandler)
	router.POST(fmt.Sprintf("/resetuserphone"), ResetUserPhoneHandler)
}

func SendSmsHandler(c *gin.Context) {
	var respone common.APIRespone
	req := &common.SendSmsRequest{}
	if err := c.BindJSON(&req); err != nil {
		respone.ErrCode = common.ParamCode
		respone.ErrMsg = err.Error()
	} else {
		if err := VailMobile(req.PhoneNum); err != nil {
			respone.ErrCode = common.ExecuteCode
			respone.ErrMsg = err.Error()
		} else {
			code := MakeCode()
			if err := VailCode(code); err != nil {
				respone.ErrCode = common.ExecuteCode
				respone.ErrMsg = err.Error()
			} else {
				if err := SendCode(req.PhoneNum, code); err != nil {
					respone.ErrCode = common.ExecuteCode
					respone.ErrMsg = err.Error()
				} else {
					respone.ErrCode = common.OKCode
				}
			}
		}
	}
	c.JSON(http.StatusOK, respone)
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
			respone.ErrCode = common.OKCode
		}
	}
	c.JSON(http.StatusOK, respone)
}

func createUser(req *common.SignUpByMobileRequest) error {
	if exist := DBClient.IfExistUserName(req.UserName); exist {
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
	_, err := DBClient.GetTokenInfo(phoneNum, verificationCode)
	if err != nil {
		return err
	}
	return nil
}

//GetHistoryInfoHandler handler
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
	user, err := DBClient.GetUserInfo(req.OldUserName, req.PhoneNum)
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
	user, err := DBClient.GetUserInfo(req.UserName, req.PhoneNum)
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
	user, err := DBClient.GetUserInfo(req.UserName, req.OldPhoneNum)
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
