package agent

import (
	"context"
	"go.uber.org/zap"
	"go_logistics/config"
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

func (u *GetUserDetailByName) Call(ctx context.Context, input string) (response string, err error) {
	config.Log.Debug("agent输入", zap.String("input", input))
	unquoted, _ := strconv.Unquote(strings.TrimSpace(input))
	if unquoted != "" {
		input = unquoted
	}
	user, err := entity.GetUserByName(input)
	if err != nil {
		return
	}
	// 确保 user 不为 nil 才转换
	if user == nil {
		response = "用户不存在"
		return
	}
	response = vo.ToUserVO(user).String()
	config.Log.Debug("agent输出", zap.String("response", response))
	return
}
