package service

import (
	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model"
	"go_logistics/util"
)

// GetUserByName 获取用户信息
func GetUserByName(c *gin.Context) {
	name := c.Query("name")
	user, err := model.GetUserByName(name)
	if err != nil {
		common.ErrorResponse(c, common.UserNotFound.Code, common.UserNotFound.Message)
		return
	}
	common.SuccessResponseWithData(c, user)
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	password := c.PostForm("password")
	rePassword := c.PostForm("rePassword")
	if name == "" || phone == "" || email == "" || password == "" || rePassword == "" {
		common.ErrorResponse(c, common.UserInfoError.Code, common.UserInfoError.Message)
		return
	}
	user := model.User{}
	user.Name = name
	user.Phone = phone
	user.Email = email
	user.Salt = util.MakeSalt()
	password = util.MakePassword(password, user.Salt)
	user.Password = password
	user.Status = model.Active
	err := model.InsertUser(&user)
	if err != nil {
		common.ErrorResponse(c, common.ServerError.Code, common.ServerError.Message)
		return
	}
	common.SuccessResponse(c)
}
