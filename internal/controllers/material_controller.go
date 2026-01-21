package controllers

import (
	"path/filepath"
	"stock-flow/internal/models"
	"stock-flow/internal/pkg/response"
	"stock-flow/internal/services"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// MaterialController 耗材控制器
// 处理耗材信息的创建与查询
type MaterialController struct {
	materialService services.MaterialService
}

// BatchImport
// @Summary 批量导入耗材
// @Description 上传Excel文件批量导入耗材基础信息(需管理员权限)
// @Tags Material
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel文件"
// @Success 200 {object} response.Response{data=services.BatchImportResult} "导入结果"
// @Router /api/v1/materials/import [post]
func (ctrl *MaterialController) BatchImport(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, response.CodeBadRequest, "请上传文件")
		return
	}

	// 1. Check file size (10MB limit)
	if file.Size > 10*1024*1024 {
		response.Error(c, response.CodeBadRequest, "文件大小不能超过10MB")
		return
	}

	// 2. Check extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		response.Error(c, response.CodeBadRequest, "仅支持 .xlsx 或 .xls 格式")
		return
	}

	// 3. Open file
	f, err := file.Open()
	if err != nil {
		response.Error(c, response.CodeServerError, "文件读取失败")
		return
	}
	defer f.Close()

	// 4. Process import
	result, err := ctrl.materialService.BatchImport(f, ext)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, result)
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
// @Description 删除指定耗材(软删除，需管理员或库管员权限)
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
