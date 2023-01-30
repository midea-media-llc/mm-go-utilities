package utils

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	err "github.com/go-kratos/kratos/v2/errors"
)

func ErrAuthFail(msg ...interface{}) *err.Error {
	return err.New(401, "AuthFail", fmt.Sprintf("%v", msg))
}

func ErrUserLockUp(msg ...interface{}) *err.Error {
	return err.New(401, "UserLockUp", fmt.Sprintf("%v", msg))
}

func ErrNotAllow(msg ...interface{}) *err.Error {
	return err.New(401, "NotAllow", fmt.Sprintf("%v", msg))
}

func ErrInvalidParams(msg ...interface{}) *err.Error {
	return errors.New(400, "InvalidParams", fmt.Sprintf("%v", msg))
}

func ErrInvalidPass(msg ...interface{}) *err.Error {
	return errors.New(400, "InvalidPass", fmt.Sprintf("%v", msg))
}

func ErrInvalidUID(msg ...interface{}) *err.Error {
	return errors.New(400, "InvalidUID", fmt.Sprintf("%v", msg))
}

func ErrInvalidUsername(msg ...interface{}) *err.Error {
	return errors.New(400, "InvalidUsername", fmt.Sprintf("%v", msg))
}

func ErrMuiltiRecord(msg ...interface{}) *err.Error {
	return errors.New(400, "MuiltiRecord", fmt.Sprintf("%v", msg))
}

func ErrSystemBusy(msg ...interface{}) *err.Error {
	return errors.New(400, "SystemBusy", fmt.Sprintf("%v", msg))
}

func ErrMuiltpleUserName(msg ...interface{}) *err.Error {
	return errors.New(400, "MuiltpleUserName", fmt.Sprintf("%v", msg))
}

func ErrNoNeedCaptcha(msg ...interface{}) *err.Error {
	return errors.New(400, "NoNeedCaptcha", fmt.Sprintf("%v", msg))
}

func ErrInvalidVerifyCode(msg ...interface{}) *err.Error {
	return errors.New(400, "InvalidVerifyCode", fmt.Sprintf("%v", msg))
}

func ErrPermissionDeni(msg ...interface{}) *err.Error {
	return errors.New(400, "PermissionDeni", fmt.Sprintf("%v", msg))
}

func ErrUserNotExit(msg ...interface{}) *err.Error {
	return errors.New(400, "UserNotExit", fmt.Sprintf("%v", msg))
}

func ErrRouteNotExit(msg ...interface{}) *err.Error {
	return errors.New(400, "RouteNotExit", fmt.Sprintf("%v", msg))
}

func ErrPermissionNotExit(msg ...interface{}) *err.Error {
	return errors.New(400, "PermissionNotExit", fmt.Sprintf("%v", msg))
}

func ErrRoleNotExit(msg ...interface{}) *err.Error {
	return errors.New(400, "RoleNotExit", fmt.Sprintf("%v", msg))
}

func ErrInvalidToken(msg ...interface{}) *err.Error {
	return errors.New(400, "InvalidToken", fmt.Sprintf("%v", msg))
}
