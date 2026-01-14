package services

import (
	"stock-flow/internal/dao"
)

type StatisticsService struct {
	statsDao dao.StatisticsDao
}

type DashboardStats struct {
	TotalBatches    int64                 `json:"total_batches"`
	WarningBatches  WarningBatchesStats   `json:"warning_batches"`
	ExpiredBatches  int64                 `json:"expired_batches"`
	OutboundTrend   []dao.MonthlyOutbound `json:"outbound_trend"`
}

type WarningBatchesStats struct {
	Count int64              `json:"count"`
	List  []dao.WarningBatch `json:"list"`
}

// GetDashboardStats 获取仪表盘综合统计数据
func (s *StatisticsService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 1. 当前库存总批次
	total, err := s.statsDao.CountTotalBatches()
	if err != nil {
		return nil, err
	}
	stats.TotalBatches = total

	// 2. 临期预警库存
	warningList, err := s.statsDao.GetWarningBatches()
	if err != nil {
		return nil, err
	}
	stats.WarningBatches = WarningBatchesStats{
		Count: int64(len(warningList)),
		List:  warningList,
	}

	// 3. 已过期库存
	expired, err := s.statsDao.CountExpiredBatches()
	if err != nil {
		return nil, err
	}
	stats.ExpiredBatches = expired

	// 4. 近半年出库趋势
	trend, err := s.statsDao.GetOutboundTrend()
	if err != nil {
		return nil, err
	}
	stats.OutboundTrend = trend

	return stats, nil
}
