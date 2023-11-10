package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/wonderivan/logger"
	"time"

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

func (l *login) Login(adminuser *User) (tok string, err error) {
	if adminuser.Username != "" && adminuser.Password != "" {
		if adminuser.Username != config.AdminUser || adminuser.Password != config.AdminPasswd {
			logger.Error("username or password is wrong...")
			return "", errors.New("username or password is wrong")
		}
	} else {
		logger.Error("username or password not is null...")
		return "", errors.New("username or password not is null")
	}
	//验证账密通过后，生成token
	// 定义加密因子
	secret := "sunjiandevops"
	// 创建一个新的Token对象
	token := jwt.New(jwt.SigningMethodHS256)
	// 设置Token的Claim(声明)，这是您自定义的数据
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1 / 60 / 2).Unix() // 设置Token过期时间（1小时）
	claims["user_id"] = "1234567"
	claims["username"] = adminuser.Username

	// 使用加密因子进行签名，并获取最终的Token字符串
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("生成Token失败:", err)
		return "", errors.New("生成Token失败: " + err.Error())
	}
	fmt.Println("生成的Token:", tokenString)
	return tokenString, nil
}
