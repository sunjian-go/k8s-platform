package dao

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"k8s-platform/db"
	"k8s-platform/model"
)

var Styles styles

type styles struct {
}
type ColorResp struct {
	Background string `json:"background" binding:"required"`
	Color      string `json:"color" binding:"required"`
}

// 初始化样式表
func (s *styles) InitStyles() error {
	initStyle := &model.Styles{
		BackgroundColor: "#6c038b",
		FontColor:       "#E4E4E4",
		Uid:             "1111111",
	}
	tx := db.GORM.Create(&initStyle)
	if tx.Error != nil {
		logger.Error("初始化styles失败: " + tx.Error.Error())
		return errors.New("初始化styles失败: " + tx.Error.Error())
	}
	return nil
}

// 获取颜色值
func (s *styles) GetColor() (*ColorResp, error) {
	//定义数据库查询返回的内容
	colors := new(model.Styles)
	//数据库查询，limit方法用于限制条数，offset方法用于 设置起始位置
	tx := db.GORM.
		Where("uid = ?", "1111111").First(&colors)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取颜色失败，" + tx.Error.Error())
		return nil, errors.New("获取颜色失败，" + tx.Error.Error())
	}
	fmt.Println("获取到为：", colors.BackgroundColor, colors.FontColor, colors.Uid)
	//fmt.Println("行数为：", tx.RowsAffected)
	return &ColorResp{
		Background: colors.BackgroundColor,
		Color:      colors.FontColor,
	}, nil

}

// 更新颜色值
func (s *styles) UpdateColor(backgroundColor, fontColor string) error {
	colors := &model.Styles{}
	//db.GORM.Where("uid=?", "1111111").Update(&colors)
	tx := db.GORM.Model(&colors).Where("uid=?", "1111111").Updates(model.Styles{BackgroundColor: backgroundColor, FontColor: fontColor})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("更新颜色失败，" + tx.Error.Error())
		return errors.New("更新颜色失败，" + tx.Error.Error())
	}
	return nil
}
