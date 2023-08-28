package service

import (
	"errors"
	"github.com/wonderivan/logger"
	//"google.golang.org/genproto/googleapis/ads/googleads/v3/errors"
	"k8s-platform/config"
)

var Login login

type login struct {
}
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *login) Login(adminuser *User) error {
	if adminuser.Username != "" && adminuser.Password != "" {
		if adminuser.Username != config.AdminUser || adminuser.Password != config.AdminPasswd {
			logger.Error("username or password is wrong...")
			return errors.New("username or password is wrong")
		}
	} else {
		logger.Error("username or password not is null...")
		return errors.New("username or password not is null")
	}
	return nil
}
