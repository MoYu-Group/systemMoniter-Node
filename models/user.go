package models

import (
	"errors"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type LoginUser struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// type LoginUserResponse struct {
// 	Error    int    `json:"error"`
// 	ErrorMsg string `json:"error_msg"`
// 	Data     map[string]interface{}
// }

// type UserData struct {
// 	Name   string `json:"name"`
// 	Token  string `json:"token"`
// 	UserID string `json:"userId"`
// }

func SetLoginUser(loginUser *LoginUser) error {
	user := viper.GetString("user")
	password := viper.GetString("password")
	if user == "" || password == "" {
		err := errors.New("Lost LoginUser info")
		zap.L().Error(err.Error())
		return err
	}
	loginUser.User = user
	loginUser.Password = password
	return nil
}
