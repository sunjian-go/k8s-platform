package dao

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"k8s-platform/db"
	"k8s-platform/model"
)

var Workflow workflow

type workflow struct {
}

type workflowResp struct {
	Items []*model.WorkFlow `json:"items"`
	Total int               `json:"total"`
}

// 获取workflow列表
func (w *workflow) GetWorkflow(fiterName string, page, limit int) (data *workflowResp, err error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	//定义数据库查询返回的内容
	var workflowList []*model.WorkFlow
	//数据库查询，limit方法用于限制条数，offset方法用于 设置起始位置
	tx := db.GORM.
		Where("name like ?", "%"+fiterName+"%").
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&workflowList)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Workflow列表失败，" + tx.Error.Error())
		return nil, errors.New("获取Workflow列表失败，" + tx.Error.Error())
	}
	fmt.Println("debug4")
	return &workflowResp{
		Items: workflowList,
		Total: len(workflowList),
	}, nil
}

// 获取单挑workflow数据
func (w *workflow) GetById(id int) (workflow *model.WorkFlow, err error) {
	workflow = &model.WorkFlow{}
	tx := db.GORM.Where("id=?", id).Find(workflow)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		logger.Error("获取workflow单条数据失败 " + tx.Error.Error())
		return nil, errors.New("获取workflow单条数据失败 " + tx.Error.Error())
	}
	return
}

// 表数据新增
func (w *workflow) Add(workflow *model.WorkFlow) (err error) {
	tx := db.GORM.Create(&workflow)
	if tx.Error != nil {
		logger.Error("添加workflow失败: " + tx.Error.Error())
		return errors.New("添加workflow失败: " + tx.Error.Error())
	}
	return nil
}

// 表数据删除
func (w *workflow) Delete(id int) (err error) {
	tx := db.GORM.Where("id=?", id).Delete(&model.WorkFlow{})
	if tx.Error != nil {
		logger.Error("删除workflow失败: " + tx.Error.Error())
		return errors.New("删除workflow失败: " + tx.Error.Error())
	}
	return nil
}
