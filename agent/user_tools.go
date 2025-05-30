package agent

import (
	"context"
	"go_logistics/model/entity"
	"go_logistics/model/vo"
	"strconv"
	"strings"
)

type GetUserDetailByName struct{}

func (u *GetUserDetailByName) Name() string {
	return "根据用户名字获取用户详情"
}

func (u *GetUserDetailByName) Description() string {
	return "输入用户名字，获取用户详细信息"
}

func (u *GetUserDetailByName) Call(ctx context.Context, input string) (string, error) {
	unquoted, _ := strconv.Unquote(strings.TrimSpace(input))
	if unquoted != "" {
		input = unquoted
	}
	user, err := entity.GetUserByName(input)
	if err != nil {
		return "", err
	}
	// 确保 user 不为 nil 才转换
	if user == nil {
		return "用户信息为空。", nil
	}
	return vo.ToUserVO(user).String(), nil
}
