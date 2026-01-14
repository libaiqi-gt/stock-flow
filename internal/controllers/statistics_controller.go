package controllers

import (
	"stock-flow/internal/pkg/response"
	"stock-flow/internal/services"

	"github.com/gin-gonic/gin"
)

type StatisticsController struct {
	statsService services.StatisticsService
}

// GetDashboardStats
// @Summary 获取仪表盘综合统计数据
// @Description 包含库存总批次、临期预警、过期库存及近半年出库趋势
// @Tags Statistics
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=services.DashboardStats} "统计数据"
// @Router /api/v1/statistics/dashboard [get]
func (ctrl *StatisticsController) GetDashboardStats(c *gin.Context) {
	stats, err := ctrl.statsService.GetDashboardStats()
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, stats)
}
