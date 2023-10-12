package model

import "time"

type Styles struct {
	//gorm:"primaryKey"用于声明主键
	ID              uint       `json:"id" gorm:"primaryKey"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	BackgroundColor string     `json:"background"` //设置前端背景色
	FontColor       string     `json:"font_color"` //设置前端字体颜色
	Uid             string     `json:"uid"`
}

// 定义TableName方法，返回mysql表名，以此来定义mysql中的表名
func (*Styles) TableName() string {
	return "styles"
}
