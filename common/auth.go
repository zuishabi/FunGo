package common

import (
	"context"
	"errors"
)

func CheckLogin(ctx context.Context) (uint64, error) {
	uidValue := ctx.Value("user_id")
	if uidValue == nil {
		return 0, errors.New("用户没有登录")
	}
	uid := uidValue.(uint64)
	return uid, nil
}
