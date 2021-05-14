package model

import (
	"context"
	"excel2config/internal/def"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	"sort"
	"strings"
)

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

type Article struct {
	ID      int
	Content string
	Author  string
}

type Sheet struct {
	CalcChain                      []CalcChain                      `json:"calcChain,omitempty" bson:"calcChain,omitempty"`
	Celldata                       []Cell                           `json:"celldata" bson:"celldata"`
	ChWidth                        int                              `json:"ch_width,omitempty" bson:"ch_width,omitempty"`
	Column                         int                              `json:"column" bson:"column"`
	Config                         *SheetConfig                     `json:"config,omitempty" bson:"config,omitempty"`
	Index                          string                           `json:"index" bson:"index"`
	Color                          string                           `json:"color" bson:"color"`
	LuckysheetSelectSave           []SelectSave                     `json:"luckysheet_select_save,omitempty" bson:"luckysheet_select_save,omitempty"`
	Name                           string                           `json:"name" bson:"name"`
	Order                          int                              `json:"order" bson:"order"`
	RhHeight                       int                              `json:"rh_height,omitempty" bson:"rh_height,omitempty"`
	Row                            int                              `json:"row,omitempty" bson:"row,omitempty"`
	ScrollLeft                     int                              `json:"scrollLeft,omitempty" bson:"scrollLeft,omitempty"`
	ScrollTop                      int                              `json:"scrollTop,omitempty" bson:"scrollTop,omitempty"`
	Status                         FlexInt                          `json:"status" bson:"status"`
	ShowGridLines                  int                              `json:"showGridLines,omitempty" bson:"showGridLines,omitempty"`
	ZoomRatio                      float64                          `json:"zoomRatio,omitempty" bson:"zoomRatio,omitempty"`
	Filter                         map[string]Filter                `json:"filter,omitempty" bson:"filter,omitempty"`
	FilterSelect                   *FilterSelect                    `json:"filter_select,omitempty" bson:"filter_select,omitempty"`
	Images                         interface{}                      `json:"images,omitempty" bson:"images,omitempty"`
	AlternateFormatSave            []AlternateFormatSave            `json:"luckysheet_alternateformat_save,omitempty" bson:"luckysheet_alternateformat_save,omitempty"`
	AlternateFormatSaveModelCustom []AlternateFormatSaveModelCustom `json:"luckysheet_alternateformat_save_modelCustom,omitempty" bson:"luckysheet_alternateformat_save_modelCustom,omitempty"`
	ConditionFormatSave            []ConditionFormatSave            `json:"luckysheet_conditionformat_save,omitempty" bson:"luckysheet_conditionformat_save,omitempty"`
	Frozen                         *Frozen                          `json:"frozen,omitempty" bson:"frozen,omitempty"`
	Chart                          []Chart                          `json:"chart,omitempty" bson:"chart,omitempty"`
	Image                          []Image                          `json:"image,omitempty" bson:"image,omitempty"`
	IsPivotTable                   bool                             `json:"isPivotTable,omitempty" bson:"isPivotTable,omitempty"`
	PivotTable                     *PivotTable                      `json:"pivotTable,omitempty" bson:"pivotTable,omitempty"`
	DynamicArray                   []DynamicArray                   `json:"dynamicArray,omitempty" bson:"dynamicArray,omitempty"`
}

type Excel struct {
	Id           string   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string   `json:"name" bson:"name"`
	CreateTime   int64    `json:"create_time" bson:"create_time"`
	EditTime     int64    `json:"edit_time" bson:"edit_time"`
	Owner        string   `json:"owner" bson:"owner"`
	Remark       string   `json:"remark,omitempty" bson:"remark,omitempty"`
	GroupId      string   `json:"group_id,omitempty" bson:"group_id,omitempty"`
	Contributers []string `json:"contributers" bson:"contributers"`
}

type ExportRecord struct {
	Id        string `json:"id,omitempty" bson:"_id,omitempty"`
	GridKey   string `json:"gridKey" bson:"gridKey"`
	SheetName string `json:"sheetName" bson:"sheetName"`
	UserName  string `json:"userName" bson:"userName"`
	Time      int64  `json:"time" bson:"time"`
	Remark    string `json:"remark" bson:"remark"`
	Sheet     *Sheet `json:"sheet" bson:"sheet"`
}

type FormatSheet struct {
	Fields  []string
	Types   []string
	Descs   []string
	Content []map[string]interface{}
}

func (s *Sheet) Format() (formatSheet *FormatSheet, err error) {
	var res []map[string]interface{}
	//先读取表头，前三行分别是字段名，字段类型，字段备注
	var fields []string
	var types []string
	var descs []string
	sort.SliceStable(s.Celldata, func(i, j int) bool {
		return s.Celldata[i].C < s.Celldata[j].C || (s.Celldata[i].C == s.Celldata[j].C && s.Celldata[i].R < s.Celldata[j].R)
	})
	var descMap = make(map[int]interface{})
	for _, cell := range s.Celldata {
		if cell.R < 2 && IsEmptyCell(cell) {
			continue
		}
		if cell.R == 0 { // field
			if field, ok := cell.V.V.(string); ok {
				field = strings.Trim(field, "\r\t\n ")
				fields = append(fields, field)
			}
		}
		if int(cell.R) == 1 { // type
			if t, ok := cell.V.V.(string); ok {
				t = strings.Trim(t, "\r\t\n")
				if t != "string" && t != "int" {
					err = ecode.Int(int(def.ErrTableHead))
					log.Errorw(context.TODO(), "type", t, "cell", cell, "msg", "table type error")
					return
				}
				types = append(types, t)
			}
		}
		if int(cell.R) == 2 { // desc
			if desc, ok := cell.V.V.(string); ok {
				descMap[int(cell.C)] = desc
			}
		}
	}
	for i, _ := range fields {
		if desc, ok := descMap[i]; ok {
			descs = append(descs, desc.(string))
		} else {
			descs = append(descs, "")
		}
	}
	fieldCount := len(fields)
	if fieldCount != len(types) || fieldCount != len(descs) || fieldCount == 0 {
		log.Errorw(context.TODO(), "fieldCount", fieldCount, "type_len", len(types), "desc_len", len(descs), "msg", "head len not valid")
		err = ecode.Int(int(def.ErrTableHead))
		return
	}
	content := make([]map[int]interface{}, 0)
	for _, cell := range s.Celldata {
		if IsEmptyCell(cell) || int(cell.C) >= fieldCount {
			continue
		}
		if len(content) < int(cell.R) {
			continue
		}
		if int(cell.R) >= len(content) {
			content = append(content, make(map[int]interface{}))
		}
		content[int(cell.R)][int(cell.C)] = cell.V.V
	}
	if len(content) < 3 {
		log.Errorw(context.TODO(), "content", content, "msg", "content len not valid")
		err = ecode.Int(int(def.ErrTableHead))
		return
	}
	res = make([]map[string]interface{}, 0)
	for row, columns := range content {
		if row < 3 {
			continue
		}
		rowContent := make(map[string]interface{})
		for k, field := range fields {
			var val interface{}
			if types != nil && k < len(types) {
				if types[k] == "string" {
					val = ""
				} else {
					val = 0
				}
			}
			if v, exist := columns[k]; exist {
				val = v
			}
			rowContent[field] = val
		}
		res = append(res, rowContent)
	}
	formatSheet = new(FormatSheet)
	formatSheet.Content = res
	formatSheet.Fields = fields
	formatSheet.Types = types
	formatSheet.Descs = descs
	return
}

func IsEmptyCell(cell Cell) bool {
	if ok, v := IsInlineCell(cell); ok {
		cell.V.V = v
	}
	if cell.V.V == nil || cell.V.V == 0 || cell.V.V == "" {
		return true
	}
	return false
}

func IsInlineCell(cell Cell) (bool, interface{}) {
	if cell.V.V == nil && cell.V.Ct != nil && cell.V.Ct.T == "inlineStr" && len(cell.V.Ct.S) > 0 {
		return true, cell.V.Ct.S[0].V
	}
	return false, nil
}
