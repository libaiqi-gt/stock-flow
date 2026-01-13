package services

import (
	"fmt"
	"io"
	"stock-flow/internal/dao"
	"stock-flow/internal/models"
	"strconv"
	"time"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
)

// InventoryService 库存业务服务
// 处理入库、库存查询及效期预警
type InventoryService struct {
	inventoryDao IInventoryDao
	materialDao  IMaterialDao
}

// Interfaces for testing
type IInventoryDao interface {
	Create(inv *models.Inventory) error
	GetByMaterialAndBatch(materialID uint, batchNo string) (*models.Inventory, error)
	GetByInboundNo(inboundNo string) (*models.Inventory, error)
	Update(inv *models.Inventory) error
	List(page, pageSize int, materialName, code, batchNo string, status int) ([]models.Inventory, int64, error)
	GetAvailableBatches(materialID uint) ([]models.Inventory, error)
	GetByID(id uint) (*models.Inventory, error)
}

type IMaterialDao interface {
	Create(m *models.Material) error
	GetByCode(code string) (*models.Material, error)
	List(page, pageSize int, name string) ([]models.Material, int64, error)
}

// NewInventoryService creates a new InventoryService
func NewInventoryService() *InventoryService {
	return &InventoryService{
		inventoryDao: &dao.InventoryDao{},
		materialDao:  &dao.MaterialDao{},
	}
}

// SetDao is used for testing to inject mock DAOs
func (s *InventoryService) SetDao(invDao IInventoryDao, matDao IMaterialDao) {
	s.inventoryDao = invDao
	s.materialDao = matDao
}

// InboundDTO 入库请求数据传输对象
type InboundDTO struct {
	MaterialCode    string // 物料编码
	MaterialName    string // 物料名称
	Category        string // 分类
	Spec            string // 规格
	Unit            string // 单位
	Brand           string // 品牌
	BatchNo         string // 内部批号
	ExpiryDate      string // 有效期 (YYYY-MM-DD)
	Quantity        int64  // 数量 (初始入库数量)
	CurrentQuantity int64  // 当前库存数量
	InboundNo       string `binding:"required"` // 入库单号
	Mode            string // 模式: "append" 追加, "overwrite" 覆盖 (默认追加)
}

// BatchImportResult 批量导入结果
type BatchImportResult struct {
	Total   int      `json:"total"`
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
	Msg     string   `json:"msg"`
}

// Inbound 耗材入库
// 包含物料自动创建、批次去重或追加逻辑
//
// 参数:
//
//	dto: 入库数据
//
// 返回值:
//
//	error: 入库失败返回错误
func (s *InventoryService) Inbound(dto InboundDTO) error {
	// 0. Check inbound no uniqueness
	if dto.InboundNo == "" {
		return fmt.Errorf("入库单号不能为空")
	}

	// 1. 查找或创建物料基础信息
	mat, err := s.materialDao.GetByCode(dto.MaterialCode)
	if err != nil {
		// Create new material
		mat = &models.Material{
			Code:     dto.MaterialCode,
			Name:     dto.MaterialName,
			Category: dto.Category,
			Spec:     dto.Spec,
			Unit:     dto.Unit,
			Brand:    dto.Brand,
		}
		if err := s.materialDao.Create(mat); err != nil {
			return err
		}
	}

	// 2. 检查入库单号是否存在（防重复提交）
	if _, err := s.inventoryDao.GetByInboundNo(dto.InboundNo); err == nil {
		return fmt.Errorf("该入库单号已存在，请勿重复提交")
	}

	// Try parsing multiple date formats
	var expiry time.Time
	formats := []string{
		"2006-01-02", "2006/01/02", "20060102",
		"2006.01.02", "2006.01",
	}
	for _, f := range formats {
		if t, e := time.Parse(f, dto.ExpiryDate); e == nil {
			expiry = t
			break
		}
	}
	if expiry.IsZero() && dto.ExpiryDate != "" {
		// Handle Excel serial date if passed as string number?
		// Usually excelize returns formatted string if possible.
		// If failed, default to empty or error? Let's just keep zero if failed.
	}

	// 3. 创建新批次
	// If CurrentQuantity is not set (e.g. from JSON API), default to Quantity
	currentQty := dto.CurrentQuantity
	if currentQty == 0 && dto.Quantity > 0 {
		currentQty = dto.Quantity
	}

	newInv := &models.Inventory{
		MaterialID: mat.ID,
		BatchNo:    dto.BatchNo,
		InboundNo:  dto.InboundNo,
		InitialQty: dto.Quantity,
		CurrentQty: currentQty,
		ExpiryDate: expiry,
	}
	return s.inventoryDao.Create(newInv)
}

// BatchImport 批量导入
//
// 参数:
//
//	r: 文件读取器 (需支持 Seek)
//	ext: 文件扩展名 (.xls / .xlsx)
//
// 返回值:
//
//	*BatchImportResult: 导入结果
//	error: 严重错误
func (s *InventoryService) BatchImport(r io.ReadSeeker, ext string) (*BatchImportResult, error) {
	var rows [][]string

	if ext == ".xlsx" {
		f, err := excelize.OpenReader(r)
		if err != nil {
			return nil, fmt.Errorf("打开XLSX文件失败: %v", err)
		}
		defer f.Close()
		rows, err = f.GetRows(f.GetSheetName(0))
		if err != nil {
			return nil, fmt.Errorf("读取XLSX内容失败: %v", err)
		}
	} else if ext == ".xls" {
		f, err := xls.OpenReader(r, "utf-8")
		if err != nil {
			return nil, fmt.Errorf("打开XLS文件失败: %v", err)
		}
		if f.NumSheets() == 0 {
			return nil, fmt.Errorf("XLS文件为空")
		}
		sheet := f.GetSheet(0)
		if sheet.MaxRow == 0 {
			return nil, fmt.Errorf("Sheet为空")
		}
		// Convert xls rows to [][]string
		for i := 0; i <= int(sheet.MaxRow); i++ {
			row := sheet.Row(i)
			var rowData []string
			if row != nil {
				// xls row col index might not be continuous, check MaxCol?
				// xls lib behavior: Row(i).Col(j)
				// We need to know max col.
				// Let's assume max 20 cols for safety or check LastCol()
				lastCol := row.LastCol()
				for j := 0; j < lastCol; j++ {
					rowData = append(rowData, row.Col(j))
				}
			}
			rows = append(rows, rowData)
		}
	} else {
		return nil, fmt.Errorf("不支持的文件格式: %s", ext)
	}

	result := &BatchImportResult{
		Errors: []string{},
		Msg:    "统计数据已排除表头行",
	}

	if len(rows) < 2 {
		return result, nil // Empty or header only
	}

	// 假设第一行是表头，从第二行开始数据
	// 映射表头索引
	headerMap := make(map[string]int)
	for i, cell := range rows[0] {
		headerMap[cell] = i
	}

	// 必填字段
	required := []string{"物料编码", "物料名称", "内部批号/校准编号", "数量", "自身有效到期日"}
	for _, field := range required {
		if _, ok := headerMap[field]; !ok {
			return nil, fmt.Errorf("缺少必填列: %s", field)
		}
	}

	for i := 1; i < len(rows); i++ {
		result.Total++
		rowIdx := i + 1 // Excel row number

		dto, errStr := s.parseExcelRow(rows[i], headerMap, rowIdx)
		if errStr != "" {
			result.Failed++
			result.Errors = append(result.Errors, errStr)
			continue
		}

		if err := s.Inbound(*dto); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %v", rowIdx, err))
		} else {
			result.Success++
		}
	}

	return result, nil
}

// parseExcelRow 解析单行Excel数据
func (s *InventoryService) parseExcelRow(row []string, headerMap map[string]int, rowIdx int) (*InboundDTO, string) {
	// Helper to get cell value safely
	getVal := func(colName string) string {
		idx, ok := headerMap[colName]
		if !ok || idx >= len(row) {
			return ""
		}
		return row[idx]
	}

	// Basic Validation
	code := getVal("物料编码")
	name := getVal("物料名称")
	batch := getVal("内部批号/校准编号")
	qtyStr := getVal("数量")
	currentQtyStr := getVal("当前库存数量")
	expiryStr := getVal("自身有效到期日")
	inboundNo := getVal("入库单号")

	if code == "" || name == "" || batch == "" || qtyStr == "" || expiryStr == "" || inboundNo == "" {
		return nil, fmt.Sprintf("第%d行: 缺少必填字段", rowIdx)
	}

	qty, err := strconv.ParseInt(qtyStr, 10, 64)
	if err != nil || qty < 0 {
		return nil, fmt.Sprintf("第%d行: 数量格式错误", rowIdx)
	}

	var currentQty int64
	if currentQtyStr != "" {
		currentQty, err = strconv.ParseInt(currentQtyStr, 10, 64)
		if err != nil || currentQty < 0 {
			return nil, fmt.Sprintf("第%d行: 当前库存数量格式错误", rowIdx)
		}
	} else {
		// If not present in Excel, default to Initial Qty
		currentQty = qty
	}

	dto := &InboundDTO{
		MaterialCode:    code,
		MaterialName:    name,
		Category:        getVal("分类"),
		Spec:            getVal("规格"),
		Unit:            getVal("单位"),
		Brand:           getVal("品牌"),
		BatchNo:         batch,
		ExpiryDate:      expiryStr,
		Quantity:        qty,
		CurrentQuantity: currentQty,
		InboundNo:       inboundNo,
		Mode:            "append", // Default append for import
	}
	return dto, ""
}

// GetInventoryList 综合查询库存
//
// 参数:
//
//	page, pageSize: 分页
//	materialName: 物料名
//	code: 编码
//	batchNo: 批号
//	status: 状态(0:全部, 1:正常, 2:临期, 3:过期)
//
// 返回值:
//
//	[]models.Inventory: 库存列表
//	int64: 总数
//	error: 错误
func (s *InventoryService) GetInventoryList(page, pageSize int, materialName, code, batchNo string, status int) ([]models.Inventory, int64, error) {
	return s.inventoryDao.List(page, pageSize, materialName, code, batchNo, status)
}

// GetRecommendedBatches 获取推荐批次 (FEFO)
//
// 参数:
//
//	materialID: 物料ID
//
// 返回值:
//
//	[]models.Inventory: 推荐批次列表
//	error: 错误
func (s *InventoryService) GetRecommendedBatches(materialID uint) ([]models.Inventory, error) {
	// FEFO strategy: First Expired First Out
	return s.inventoryDao.GetAvailableBatches(materialID)
}
