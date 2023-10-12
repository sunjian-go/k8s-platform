package service

import (
	"fmt"
	"k8s-platform/dao"
)

var Styles styles

type styles struct {
}

// 获取颜色信息
func (s *styles) GetColor() (*dao.ColorResp, error) {
	colors, err := dao.Styles.GetColor()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return colors, nil
}

// 更新颜色信息
func (s *styles) UpdateColor(backgroundColor, fontColor string) error {
	err := dao.Styles.UpdateColor(backgroundColor, fontColor)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
