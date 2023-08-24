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

type WorkflowResp struct {
	Items []*model.WorkFlow `json:"items"`
	Total int               `json:"total"`
}

// 获取workflow列表
func (w *workflow) GetWorkflow(fiterName, namespace string, page, limit int) (data *WorkflowResp, err error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	//定义数据库查询返回的内容
	var workflowList []*model.WorkFlow
	//数据库查询，limit方法用于限制条数，offset方法用于 设置起始位置
	tx := db.GORM.
		Where("name like ?", "%"+fiterName+"%").Where("namespace = ?", namespace).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&workflowList)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Workflow列表失败，" + tx.Error.Error())
		return nil, errors.New("获取Workflow列表失败，" + tx.Error.Error())
	}
	return &WorkflowResp{
		Items: workflowList,
		Total: len(workflowList),
	}, nil
}

// 获取单条workflow数据
func (w *workflow) GetById(id int) (workflow *model.WorkFlow, err error) { //形参只是声明类型，并没有开辟空间
	work := new(model.WorkFlow) //给结构体开辟空间
	tx := db.GORM.Where("id=?", id).First(work)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取workflow单条数据失败 " + tx.Error.Error())
		return nil, errors.New("获取workflow单条数据失败 " + tx.Error.Error())
	}
	return work, nil
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
// 删除workflow
// 软删除 db.GORM.Delete("id = ?", id)
// 软删除执行的是UPDATE语句，将deleted_at字段设置为时间即可， gorm 默认就是软删。
// 实际执行语句 UPDATE `workflow` SET `deleted_at` = '2021-03-01 08:32:11' WHERE `id` IN ('1'
// 硬删除 db.GORM.Unscoped().Delete("id = ?", id)) 直接从表中删除这条数据
// 实际执行语句 DELETE FROM `workflow` WHERE `id` IN ('1');
func (w *workflow) Delete(id int) (err error) {
	fmt.Println("传过来的id为：", id)
	tx := db.GORM.Where("id=?", id).Delete(&model.WorkFlow{})
	if tx.Error != nil {
		logger.Error("删除workflow失败: " + tx.Error.Error())
		return errors.New("删除workflow失败: " + tx.Error.Error())
	}
	return nil
}
