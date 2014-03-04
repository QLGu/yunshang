package controllers

import (
	"strconv"
)

// DataTables server-side响应数据结构
type dataTableData struct {
	SEcho                int         `json:"sEcho"`
	ITotalRecords        int64       `json:"iTotalRecords"`
	ITotalDisplayRecords int64       `json:"iTotalDisplayRecords"`
	AaData               interface{} `json:"aaData,omitempty"`
}

// 构建dataTableData
func DataTableData(echo string, total int64, totalDisplay int64, data interface{}) dataTableData {
	ei, err := strconv.Atoi(echo)
	if err != nil {
		ei = 0
	}
	return dataTableData{SEcho: ei, ITotalRecords: total, ITotalDisplayRecords: totalDisplay, AaData: data}
}
