package controllers

import (
	"stock-flow/internal/models"
	"stock-flow/internal/pkg/response"
	"stock-flow/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MaterialController 耗材控制器
// 处理耗材信息的创建与查询
type MaterialController struct {
	materialService services.MaterialService
}

// Create
// @Summary 创建耗材基础信息
// @Description 录入新的耗材信息(需管理员或库管员权限)
// @Tags Material
// @Accept json
// @Produce json
// @Param request body models.Material true "耗材信息"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/materials [post]
func (ctrl *MaterialController) Create(c *gin.Context) {
	var m models.Material
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	if err := ctrl.materialService.CreateMaterial(&m); err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success[any](c, nil)
}

// Delete
// @Summary 删除耗材
// @Description 删除指定耗材(需管理员或库管员权限)
// @Tags Material
// @Accept json
// @Produce json
// @Param id path int true "耗材ID"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/materials/{id} [delete]
func (ctrl *MaterialController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "Invalid ID format")
		return
	}

	if err := ctrl.materialService.DeleteMaterial(uint(id)); err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success[any](c, nil)
}

// List
// @Summary 查询耗材列表
// @Description 分页查询耗材基础信息，支持模糊搜索
// @Tags Material
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "物料名称(模糊)"
// @Success 200 {object} response.Response "列表数据"
// @Router /api/v1/materials [get]
func (ctrl *MaterialController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	name := c.Query("name")

	list, total, err := ctrl.materialService.GetMaterialList(page, pageSize, name)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
