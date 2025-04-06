package service

import (
	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model"
	"go_logistics/util"
	"strconv"
)

// GetUserByName 获取用户信息
func GetUserByName(c *gin.Context) {
	name := c.Query("name")
	user, err := model.GetUserByName(name)
	if err != nil {
		common.ErrorResponse(c, common.RecordNotFound)
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
		common.ErrorResponse(c, common.ParamError)
		return
	}
	user, _ := model.GetUserByName(name)
	if user.Status != 0 {
		common.ErrorResponse(c, common.RecordExist)
		return
	}
	user = &model.User{}
	user.Name = name
	user.Phone = phone
	user.Email = email
	user.Salt = util.MakeSalt()
	password = util.MakePassword(password, user.Salt)
	user.Password = password
	user.Status = model.Active
	err := model.InsertUser(user)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	var dto model.FindUserListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	users, err := model.GetUserList(dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponseWithData(c, users)
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	status := c.PostForm("status")
	if name == "" || phone == "" || email == "" || status == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	user, err := model.GetUserByName(name)
	if err != nil {
		common.ErrorResponse(c, common.RecordNotFound)
	}
	user.Phone = phone
	user.Email = email
	statusInt, err := strconv.Atoi(status)
	if err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	user.Status = model.UserStatus(statusInt)
	err = model.UpdateUser(user)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	err := model.DeleteUser(name)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	common.SuccessResponse(c)
}

// LoginUser 登录用户
func LoginUser(c *gin.Context) {
	name := c.Query("name")
	password := c.Query("password")
	if name == "" || password == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	user, err := model.GetUserByName(name)
	if err != nil {
		common.ErrorResponse(c, common.RecordNotFound)
		return
	}
	if user.Status == model.Banned {
		common.ErrorResponse(c, common.UserBanned)
		return
	}
	if user.Status == model.Deleted {
		common.ErrorResponse(c, common.UserDeleted)
		return
	}
	if !util.ValidPassword(password, user.Salt, user.Password) {
		common.ErrorResponse(c, common.UserNameOrPasswordError)
	}
	token, err := util.GenerateToken(user.Name)
	if err != nil {
		common.ErrorResponse(c, common.ServerError)
		return
	}
	c.Header("logistics_token", token)
	common.SuccessResponseWithData(c, user)
	return
}
