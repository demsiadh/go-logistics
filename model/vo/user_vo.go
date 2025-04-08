package vo

import (
	"go_logistics/model/entity"
	"strconv"
	"time"
)

// UserVO User 用户结构
type UserVO struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Status     string `json:"status"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

// ToUserVO 将User模型转换为UserVO值对象
// 参数:
//
//	user: 要转换的User模型
//
// 返回:
//
//	UserVO: 转换后的值对象
func ToUserVO(user *entity.User) UserVO {
	// 加载中国时区
	loc, _ := time.LoadLocation("Asia/Shanghai")

	return UserVO{
		Name:       user.Name,
		Phone:      user.Phone,
		Email:      user.Email,
		Status:     strconv.Itoa(int(user.Status)),
		CreateTime: user.CreateTime.Time().In(loc).Format("2006-01-02 15:04:05"),
		UpdateTime: user.UpdateTime.Time().In(loc).Format("2006-01-02 15:04:05"),
	}
}

// ToUserVOList 将User模型列表转换为UserVO列表
func ToUserVOList(users []*entity.User) []UserVO {
	var voList []UserVO
	for _, user := range users {
		voList = append(voList, ToUserVO(user))
	}
	return voList
}
