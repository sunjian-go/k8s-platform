package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	// 定义加密因子
	secret := "sunjiandevops"

	// 创建一个新的Token对象
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置Token的Claim(声明)，这是您自定义的数据
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // 设置Token过期时间（1小时）
	claims["user_id"] = "your_user_id"
	claims["username"] = "your_username"

	// 使用加密因子进行签名，并获取最终的Token字符串
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("生成Token失败:", err)
		return
	}

	fmt.Println("生成的Token:", tokenString)
}
