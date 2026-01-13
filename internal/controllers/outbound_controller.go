package controllers

import (
	"stock-flow/internal/pkg/response"
	"stock-flow/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OutboundController 领用控制器
// 处理领用申请、记录查询及状态变更
type OutboundController struct {
	outboundService services.OutboundService
}

// ApplyOutboundReq 领用申请请求参数
type ApplyOutboundReq struct {
	InventoryID uint   `json:"inventory_id" binding:"required"`  // 库存ID
	Quantity    int64  `json:"quantity" binding:"required,gt=0"` // 领用数量(>0)
	Purpose     string `json:"purpose" binding:"required"`       // 领用用途
}

// Apply
// @Summary 领用申请
// @Description 提交领用申请，扣减库存并生成记录
// @Tags Outbound
// @Accept json
// @Produce json
// @Param request body ApplyOutboundReq true "领用信息"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/outbound/apply [post]
func (ctrl *OutboundController) Apply(c *gin.Context) {
	var req ApplyOutboundReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	userID, _ := c.Get("userID")

	dto := services.OutboundApplyDTO{
		InventoryID: req.InventoryID,
		UserID:      userID.(uint),
		Quantity:    req.Quantity,
		Purpose:     req.Purpose,
	}

	if err := ctrl.outboundService.ApplyOutbound(dto); err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success[any](c, nil)
}

// List
// @Summary 我的领用记录
// @Description 查询当前登录用户的领用历史
// @Tags Outbound
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response "列表数据"
// @Router /api/v1/outbound/my [get]
func (ctrl *OutboundController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	userID, _ := c.Get("userID")

	list, total, err := ctrl.outboundService.GetOutboundList(page, pageSize, userID.(uint))
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// UpdateStatus
// @Summary 更新使用状态
// @Description 更新领用记录的状态(如: USING -> FINISHED)
// @Tags Outbound
// @Param id path int true "记录ID"
// @Param status query string true "新状态 (USING/FINISHED)"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/outbound/{id}/status [put]
func (ctrl *OutboundController) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	status := c.Query("status")
	if status == "" {
		response.Error(c, response.CodeBadRequest, "status required")
		return
	}

	// In real world, check if record belongs to user
	if err := ctrl.outboundService.UpdateStatus(uint(id), status); err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success[any](c, nil)
}
