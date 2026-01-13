package controllers

import (
	"stock-flow/internal/pkg/response"
	"stock-flow/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// InventoryController 库存控制器
// 处理入库、库存查询和效期预警
type InventoryController struct {
	inventoryService services.InventoryService
}

// Inbound
// @Summary 耗材入库
// @Description 耗材批量入库接口，支持自动创建新物料或追加库存
// @Tags Inventory
// @Accept json
// @Produce json
// @Param request body services.InboundDTO true "入库信息"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/inventory/inbound [post]
func (ctrl *InventoryController) Inbound(c *gin.Context) {
	var dto services.InboundDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	if err := ctrl.inventoryService.Inbound(dto); err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success[any](c, nil)
}

// List
// @Summary 库存总表查询
// @Description 综合查询库存状态，支持效期预警筛选
// @Tags Inventory
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param material_name query string false "物料名称"
// @Param code query string false "物料编码"
// @Param batch_no query string false "批号"
// @Param status query int false "状态: 0全部, 1正常, 2临期, 3过期"
// @Success 200 {object} response.Response "列表数据"
// @Router /api/v1/inventory [get]
func (ctrl *InventoryController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	materialName := c.Query("material_name")
	code := c.Query("code")
	batchNo := c.Query("batch_no")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	list, total, err := ctrl.inventoryService.GetInventoryList(page, pageSize, materialName, code, batchNo, status)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// RecommendedBatches
// @Summary 智能推荐批次 (FEFO)
// @Description 根据 FEFO (先失效先出) 策略推荐领用批次
// @Tags Inventory
// @Param material_id query int true "物料ID"
// @Success 200 {object} response.Response "推荐批次列表"
// @Router /api/v1/inventory/recommend [get]
func (ctrl *InventoryController) RecommendedBatches(c *gin.Context) {
	materialID, _ := strconv.Atoi(c.Query("material_id"))
	if materialID == 0 {
		response.Error(c, response.CodeBadRequest, "material_id is required")
		return
	}

	list, err := ctrl.inventoryService.GetRecommendedBatches(uint(materialID))
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, list)
}
