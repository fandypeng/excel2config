package model

type MergeCell struct {
	C  int `json:"c,omitempty" bson:"c,omitempty"`
	Cs int `json:"cs,omitempty" bson:"cs,omitempty"`
	R  int `json:"r,omitempty" bson:"r,omitempty"`
	Rs int `json:"rs,omitempty" bson:"rs,omitempty"`
}

type PS struct {
	Height int    `json:"height,omitempty" bson:"height,omitempty"`
	Width  int    `json:"width,omitempty" bson:"width,omitempty"`
	Left   int    `json:"left,omitempty" bson:"left,omitempty"`
	Top    int    `json:"top,omitempty" bson:"top,omitempty"`
	IsShow bool   `json:"isshow,omitempty" bson:"isshow,omitempty"`
	Value  string `json:"value,omitempty" bson:"value,omitempty"`
}

type CellType struct { // celltype 单元格值格式
	Fa string `json:"fa,omitempty" bson:"fa,omitempty"`
	S  []struct {
		Bl int         `json:"bl,omitempty" bson:"bl,omitempty"`
		Cl int         `json:"cl,omitempty" bson:"cl,omitempty"`
		Fc string      `json:"fc,omitempty" bson:"fc,omitempty"`
		Ff int      `json:"ff,omitempty" bson:"ff,omitempty"`
		Fs interface{} `json:"fs,omitempty" bson:"fs,omitempty"`
		It int         `json:"it,omitempty" bson:"it,omitempty"`
		Un int         `json:"un,omitempty" bson:"un,omitempty"`
		V  string      `json:"v,omitempty" bson:"v,omitempty"`
	} `json:"s,omitempty" bson:"s,omitempty"`
	T string `json:"t,omitempty" bson:"t,omitempty"`
}

type CellValue struct {
	Bg string      `json:"bg,omitempty" bson:"bg,omitempty"` // 背景颜色
	Bl int         `json:"bl,omitempty" bson:"bl,omitempty"` // 0 常规 、 1加粗
	Cl int         `json:"cl,omitempty" bson:"cl,omitempty"` // 删除线
	Ct *CellType   `json:"ct,omitempty" bson:"ct,omitempty"`
	F  string      `json:"f,omitempty" bson:"f,omitempty"`   // 公式
	Fc string      `json:"fc,omitempty" bson:"fc,omitempty"` // 字体颜色
	Ff int         `json:"ff,omitempty" bson:"ff,omitempty"` // 字体类型
	Fs interface{} `json:"fs,omitempty" bson:"fs,omitempty"` // 字体大小
	It int         `json:"it,omitempty" bson:"it,omitempty"` // 斜体
	M  string      `json:"m,omitempty" bson:"m,omitempty"`   // 显示值
	Mc *MergeCell  `json:"mc,omitempty" bson:"mc,omitempty"` // 合并单元格
	Tb string      `json:"tb,omitempty" bson:"tb,omitempty"` // 文本换行，0 截断、1溢出、2 自动换行
	Tr interface{} `json:"tr,omitempty" bson:"tr,omitempty"` // 竖排文字
	V  interface{} `json:"v,omitempty" bson:"v,omitempty"`   // 原始值
	Ht interface{} `json:"ht,omitempty" bson:"ht,omitempty"` // 水平对齐，0 居中、1 左、2右
	Vt interface{} `json:"vt,omitempty" bson:"vt,omitempty"` // 垂直对齐，0 中间、1 上、2下
	Rt interface{} `json:"rt,omitempty" bson:"rt,omitempty"` // 文字旋转角度
	Ps *PS         `json:"ps,omitempty" bson:"ps,omitempty"` //批注
}

type Cell struct {
	C FlexInt   `json:"c" bson:"c"`
	R FlexInt   `json:"r" bson:"r"`
	V CellValue `json:"v,omitempty" bson:"v,omitempty"`
}

type BorderValueStyle struct {
	Color string `json:"color,omitempty" bson:"color,omitempty"`
	Style int    `json:"style,omitempty" bson:"style,omitempty"`
}

type Border struct {
	BorderType string `json:"borderType,omitempty" bson:"borderType,omitempty"`
	Color      string `json:"color,omitempty" bson:"color,omitempty"`
	Range      []struct {
		Column []int `json:"column,omitempty" bson:"column,omitempty"`
		Row    []int `json:"row,omitempty" bson:"row,omitempty"`
	} `json:"range,omitempty" bson:"range,omitempty"`
	RangeType string `json:"rangeType,omitempty" bson:"rangeType,omitempty"`
	Style     string `json:"style,omitempty" bson:"style,omitempty"`
	Value     struct {
		RowIndex int              `json:"row_index,omitempty" bson:"row_index,omitempty"`
		ColIndex int              `json:"col_index,omitempty" bson:"col_index,omitempty"`
		B        BorderValueStyle `json:"b,omitempty" bson:"b,omitempty"`
		L        BorderValueStyle `json:"l,omitempty" bson:"l,omitempty"`
		R        BorderValueStyle `json:"r,omitempty" bson:"r,omitempty"`
		T        BorderValueStyle `json:"t,omitempty" bson:"t,omitempty"`
	} `json:"value,omitempty" bson:"value,omitempty"`
}

type ConfigMerge struct {
	C  int `json:"c,omitempty" bson:"c,omitempty"`
	Cs int `json:"cs,omitempty" bson:"cs,omitempty"`
	R  int `json:"r,omitempty" bson:"r,omitempty"`
	Rs int `json:"rs,omitempty" bson:"rs,omitempty"`
}

type SheetConfig struct {
	BorderInfo      []Border               `json:"borderInfo,omitempty" bson:"borderInfo,omitempty"`
	Columnlen       map[string]float64     `json:"columnlen,omitempty" bson:"columnlen,omitempty"`
	CurentsheetView string                 `json:"curentsheetView,omitempty" bson:"curentsheetView,omitempty"`
	CustomHeight    map[string]float64     `json:"customHeight,omitempty" bson:"customHeight,omitempty"`
	CustomWidth     map[string]float64     `json:"customWidth,omitempty" bson:"customWidth,omitempty"`
	Merge           map[string]ConfigMerge `json:"merge,omitempty" bson:"merge,omitempty"`
	Rowhidden       map[string]float64     `json:"rowhidden,omitempty" bson:"rowhidden,omitempty"`
	Colhidden       map[string]float64     `json:"colhidden,omitempty" bson:"colhidden,omitempty"`
	Rowlen          map[string]float64     `json:"rowlen,omitempty" bson:"rowlen,omitempty"`
	SheetViewZoom   struct {
		ViewLayoutZoomScale float64 `json:"viewLayoutZoomScale,omitempty" bson:"viewLayoutZoomScale,omitempty"`
		ViewNormalZoomScale float64 `json:"viewNormalZoomScale,omitempty" bson:"viewNormalZoomScale,omitempty"`
		ViewPageZoomScale   float64 `json:"viewPageZoomScale,omitempty" bson:"viewPageZoomScale,omitempty"`
	} `json:"sheetViewZoom,omitempty" bson:"sheetViewZoom,omitempty"`
	Authority interface{} `json:"authority,omitempty" bson:"authority,omitempty"`
}

type Filter struct {
	Caljs struct {
		Text  string `json:"text,omitempty" bson:"text,omitempty"`
		Type  string `json:"type,omitempty" bson:"type,omitempty"`
		Value string `json:"value,omitempty" bson:"value,omitempty"`
	} `json:"caljs,omitempty" bson:"caljs,omitempty"`
	Cindex      float64            `json:"cindex,omitempty" bson:"cindex,omitempty"`
	Edc         float64            `json:"edc,omitempty" bson:"edc,omitempty"`
	Edr         float64            `json:"edr,omitempty" bson:"edr,omitempty"`
	Optionstate bool               `json:"optionstate,omitempty" bson:"optionstate,omitempty"`
	Rowhidden   map[string]float64 `json:"rowhidden,omitempty" bson:"rowhidden,omitempty"`
	Stc         float64            `json:"stc,omitempty" bson:"stc,omitempty"`
	Str         float64            `json:"str,omitempty" bson:"str,omitempty"`
}

type FilterSelect struct {
	Column      []float64 `json:"column,omitempty" bson:"column,omitempty"`
	Row         []float64 `json:"row,omitempty" bson:"row,omitempty"`
	ColumnFocus float64   `json:"column_focus,omitempty" bson:"column_focus,omitempty"`
	Height      float64   `json:"height,omitempty" bson:"height,omitempty"`
	HeightMove  float64   `json:"height_move,omitempty" bson:"height_move,omitempty"`
	Left        float64   `json:"left,omitempty" bson:"left,omitempty"`
	LeftMove    float64   `json:"left_move,omitempty" bson:"left_move,omitempty"`
	RowFocus    float64   `json:"row_focus,omitempty" bson:"row_focus,omitempty"`
	Top         float64   `json:"top,omitempty" bson:"top,omitempty"`
	TopMove     float64   `json:"top_move,omitempty" bson:"top_move,omitempty"`
	Width       float64   `json:"width,omitempty" bson:"width,omitempty"`
	WidthMove   float64   `json:"width_move,omitempty" bson:"width_move,omitempty"`
}

type SelectSave struct {
	Column      []float64 `json:"column,omitempty" bson:"column,omitempty"`
	ColumnFocus float64   `json:"column_focus,omitempty" bson:"column_focus,omitempty"`
	Height      float64   `json:"height,omitempty" bson:"height,omitempty"`
	HeightMove  float64   `json:"height_move,omitempty" bson:"height_move,omitempty"`
	Left        float64   `json:"left,omitempty" bson:"left,omitempty"`
	LeftMove    float64   `json:"left_move,omitempty" bson:"left_move,omitempty"`
	Row         []float64 `json:"row,omitempty" bson:"row,omitempty"`
	RowFocus    float64   `json:"row_focus,omitempty" bson:"row_focus,omitempty"`
	Top         float64   `json:"top,omitempty" bson:"top,omitempty"`
	TopMove     float64   `json:"top_move,omitempty" bson:"top_move,omitempty"`
	Width       float64   `json:"width,omitempty" bson:"width,omitempty"`
	WidthMove   float64   `json:"width_move,omitempty" bson:"width_move,omitempty"`
}

type AlternateFormatSave struct {
	Cellrange struct {
		Column []int64 `json:"column,omitempty" bson:"column,omitempty"`
		Row    []int64 `json:"row,omitempty" bson:"row,omitempty"`
	} `json:"cellrange,omitempty" bson:"cellrange,omitempty"`
	Format struct {
		Foot struct {
			Bc string `json:"bc,omitempty" bson:"bc,omitempty"`
			Fc string `json:"fc,omitempty" bson:"fc,omitempty"`
		} `json:"foot,omitempty" bson:"foot,omitempty"`
		Head struct {
			Bc string `json:"bc,omitempty" bson:"bc,omitempty"`
			Fc string `json:"fc,omitempty" bson:"fc,omitempty"`
		} `json:"head,omitempty" bson:"head,omitempty"`
		One struct {
			Bc string `json:"bc,omitempty" bson:"bc,omitempty"`
			Fc string `json:"fc,omitempty" bson:"fc,omitempty"`
		} `json:"one,omitempty" bson:"one,omitempty"`
		Two struct {
			Bc string `json:"bc,omitempty" bson:"bc,omitempty"`
			Fc string `json:"fc,omitempty" bson:"fc,omitempty"`
		} `json:"two,omitempty" bson:"two,omitempty"`
	} `json:"format,omitempty" bson:"format,omitempty"`
	HasRowFooter bool `json:"hasRowFooter,omitempty" bson:"hasRowFooter,omitempty"`
	HasRowHeader bool `json:"hasRowHeader,omitempty" bson:"hasRowHeader,omitempty"`
}

type AlternateFormatSaveModelCustom map[string]struct {
	Bc string `json:"bc,omitempty" bson:"bc,omitempty"`
	Fc string `json:"fc,omitempty" bson:"fc,omitempty"`
}

type ConditionFormatSave struct {
	Cellrange []struct {
		Column []int64 `json:"column,omitempty" bson:"column,omitempty"`
		Row    []int64 `json:"row,omitempty" bson:"row,omitempty"`
	} `json:"cellrange,omitempty" bson:"cellrange,omitempty"`
	ConditionName  string `json:"conditionName,omitempty" bson:"conditionName,omitempty"`
	ConditionRange []struct {
		Column []int64 `json:"column"`
		Row    []int64 `json:"row,omitempty" bson:"row,omitempty"`
	} `json:"conditionRange,omitempty" bson:"conditionRange,omitempty"`
	ConditionValue []int64     `json:"conditionValue,omitempty" bson:"conditionValue,omitempty"`
	Format         interface{} `json:"format,omitempty" bson:"format,omitempty"`
	Type           string      `json:"type,omitempty" bson:"type,omitempty"`
}

type Frozen struct {
	Range struct {
		ColumnFocus int64 `json:"column_focus,omitempty" bson:"column_focus,omitempty"`
		RowFocus    int64 `json:"row_focus,omitempty" bson:"row_focus,omitempty"`
	} `json:"range,omitempty" bson:"range,omitempty"`
	Type string `json:"type,omitempty" bson:"type,omitempty"`
}

type Chart struct {
	ChartOptions struct {
		ChartAllType  string `json:"chartAllType,omitempty" bson:"chartAllType,omitempty"`
		ChartID       string `json:"chart_id,omitempty" bson:"chart_id,omitempty"`
		DefaultOption struct {
			Axis struct {
				AxisType  string `json:"axisType,omitempty" bson:"axisType,omitempty"`
				XAxisDown struct {
					Data    []string `json:"data,omitempty" bson:"data,omitempty"`
					Inverse bool     `json:"inverse,omitempty" bson:"inverse,omitempty"`
					Name    string   `json:"name,omitempty" bson:"name,omitempty"`
					NetArea struct {
						ColorOne string `json:"colorOne,omitempty" bson:"colorOne,omitempty"`
						ColorTwo string `json:"colorTwo,omitempty" bson:"colorTwo,omitempty"`
						Interval struct {
							CusNumber int64  `json:"cusNumber,omitempty" bson:"cusNumber,omitempty"`
							Value     string `json:"value,omitempty" bson:"value,omitempty"`
						} `json:"interval,omitempty" bson:"interval,omitempty"`
						Show bool `json:"show,omitempty" bson:"show,omitempty"`
					} `json:"netArea,omitempty" bson:"netArea,omitempty"`
					NetLine struct {
						Color    string `json:"color,omitempty" bson:"color,omitempty"`
						Interval struct {
							CusNumber int64  `json:"cusNumber,omitempty" bson:"cusNumber,omitempty"`
							Value     string `json:"value,omitempty" bson:"value,omitempty"`
						} `json:"interval,omitempty" bson:"interval,omitempty"`
						Show  bool   `json:"show,omitempty" bson:"show,omitempty"`
						Type  string `json:"type,omitempty" bson:"type,omitempty"`
						Width int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"netLine,omitempty" bson:"netLine,omitempty"`
					Show bool `json:"show,omitempty" bson:"show,omitempty"`
					Tick struct {
						Color    string `json:"color,omitempty" bson:"color,omitempty"`
						Length   int64  `json:"length,omitempty" bson:"length,omitempty"`
						Position string `json:"position,omitempty" bson:"position,omitempty"`
						Show     bool   `json:"show,omitempty" bson:"show,omitempty"`
						Width    int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"tick,omitempty" bson:"tick,omitempty"`
					TickLabel struct {
						Digit    string `json:"digit,omitempty" bson:"digit,omitempty"`
						Distance int64  `json:"distance,omitempty" bson:"distance,omitempty"`
						Label    struct {
							Color       string        `json:"color,omitempty" bson:"color,omitempty"`
							CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
							FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
							FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
							FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
						} `json:"label,omitempty" bson:"label,omitempty"`
						Max      interface{} `json:"max,omitempty" bson:"max,omitempty"`
						Min      interface{} `json:"min,omitempty" bson:"min,omitempty"`
						Optimize int64       `json:"optimize,omitempty" bson:"optimize,omitempty"`
						Prefix   string      `json:"prefix,omitempty" bson:"prefix,omitempty"`
						Ratio    int64       `json:"ratio,omitempty" bson:"ratio,omitempty"`
						Rotate   int64       `json:"rotate,omitempty" bson:"rotate,omitempty"`
						Show     bool        `json:"show,omitempty" bson:"show,omitempty"`
						Suffix   string      `json:"suffix,omitempty" bson:"suffix,omitempty"`
					} `json:"tickLabel,omitempty" bson:"tickLabel,omitempty"`
					TickLine struct {
						Color string `json:"color,omitempty" bson:"color,omitempty"`
						Show  bool   `json:"show,omitempty" bson:"show,omitempty"`
						Width int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"tickLine,omitempty" bson:"tickLine,omitempty"`
					Title struct {
						FzPosition string `json:"fzPosition,omitempty" bson:"fzPosition,omitempty"`
						Label      struct {
							Color       string        `json:"color,omitempty" bson:"color,omitempty"`
							CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
							FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
							FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
							FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
						} `json:"label,omitempty" bson:"label,omitempty"`
						NameGap   int64  `json:"nameGap,omitempty" bson:"nameGap,omitempty"`
						Rotate    int64  `json:"rotate,omitempty" bson:"rotate,omitempty"`
						ShowTitle bool   `json:"showTitle,omitempty" bson:"showTitle,omitempty"`
						Text      string `json:"text,omitempty" bson:"text,omitempty"`
					} `json:"title,omitempty" bson:"title,omitempty"`
					Type string `json:"type,omitempty" bson:"type,omitempty"`
				} `json:"xAxisDown,omitempty" bson:"xAxisDown,omitempty"`
				XAxisUp struct {
					AxisLine struct {
						OnZero bool `json:"onZero,omitempty" bson:"onZero,omitempty"`
					} `json:"axisLine,omitempty" bson:"axisLine,omitempty"`
					Inverse bool   `json:"inverse,omitempty" bson:"inverse,omitempty"`
					Name    string `json:"name,omitempty" bson:"name,omitempty"`
					NetArea struct {
						ColorOne string `json:"colorOne,omitempty" bson:"colorOne,omitempty"`
						ColorTwo string `json:"colorTwo,omitempty" bson:"colorTwo,omitempty"`
						Interval struct {
							CusNumber int64  `json:"cusNumber,omitempty" bson:"cusNumber,omitempty"`
							Value     string `json:"value,omitempty" bson:"value,omitempty"`
						} `json:"interval,omitempty" bson:"interval,omitempty"`
						Show bool `json:"show,omitempty" bson:"show,omitempty"`
					} `json:"netArea,omitempty" bson:"netArea,omitempty"`
					NetLine struct {
						Color    string `json:"color,omitempty" bson:"color,omitempty"`
						Interval struct {
							CusNumber int64  `json:"cusNumber,omitempty" bson:"cusNumber,omitempty"`
							Value     string `json:"value,omitempty" bson:"value,omitempty"`
						} `json:"interval,omitempty" bson:"interval,omitempty"`
						Show  bool   `json:"show,omitempty" bson:"show,omitempty"`
						Type  string `json:"type,omitempty" bson:"type,omitempty"`
						Width int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"netLine,omitempty" bson:"netLine,omitempty"`
					Show bool `json:"show,omitempty" bson:"show,omitempty"`
					Tick struct {
						Color    string `json:"color,omitempty" bson:"color,omitempty"`
						Length   int64  `json:"length,omitempty" bson:"length,omitempty"`
						Position string `json:"position,omitempty" bson:"position,omitempty"`
						Show     bool   `json:"show,omitempty" bson:"show,omitempty"`
						Width    int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"tick,omitempty" bson:"tick,omitempty"`
					TickLabel struct {
						Digit    string `json:"digit,omitempty" bson:"digit,omitempty"`
						Distance int64  `json:"distance,omitempty" bson:"distance,omitempty"`
						Label    struct {
							Color       string        `json:"color,omitempty" bson:"color,omitempty"`
							CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
							FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
							FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
							FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
						} `json:"label,omitempty" bson:"label,omitempty"`
						Max      string `json:"max,omitempty" bson:"max,omitempty"`
						Min      string `json:"min,omitempty" bson:"min,omitempty"`
						Optimize int64  `json:"optimize,omitempty" bson:"optimize,omitempty"`
						Prefix   string `json:"prefix,omitempty" bson:"prefix,omitempty"`
						Ratio    int64  `json:"ratio,omitempty" bson:"ratio,omitempty"`
						Rotate   int64  `json:"rotate,omitempty" bson:"rotate,omitempty"`
						Show     bool   `json:"show,omitempty" bson:"show,omitempty"`
						Suffix   string `json:"suffix,omitempty" bson:"suffix,omitempty"`
					} `json:"tickLabel,omitempty" bson:"tickLabel,omitempty"`
					TickLine struct {
						Color string `json:"color,omitempty" bson:"color,omitempty"`
						Show  bool   `json:"show,omitempty" bson:"show,omitempty"`
						Width int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"tickLine,omitempty" bson:"tickLine,omitempty"`
					Title struct {
						FzPosition string `json:"fzPosition,omitempty" bson:"fzPosition,omitempty"`
						Label      struct {
							Color       string        `json:"color,omitempty" bson:"color,omitempty"`
							CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
							FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
							FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
							FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
						} `json:"label,omitempty" bson:"label,omitempty"`
						NameGap   int64  `json:"nameGap,omitempty" bson:"nameGap,omitempty"`
						Rotate    int64  `json:"rotate,omitempty" bson:"rotate,omitempty"`
						ShowTitle bool   `json:"showTitle,omitempty" bson:"showTitle,omitempty"`
						Text      string `json:"text,omitempty" bson:"text,omitempty"`
					} `json:"title,omitempty" bson:"title,omitempty"`
				} `json:"xAxisUp,omitempty" bson:"xAxisUp,omitempty"`
				YAxisLeft struct {
					Inverse bool   `json:"inverse,omitempty" bson:"inverse,omitempty"`
					Name    string `json:"name,omitempty" bson:"name,omitempty"`
					NetArea struct {
						ColorOne string `json:"colorOne,omitempty" bson:"colorOne,omitempty"`
						ColorTwo string `json:"colorTwo,omitempty" bson:"colorTwo,omitempty"`
						Interval struct {
							CusNumber int64  `json:"cusNumber,omitempty" bson:"cusNumber,omitempty"`
							Value     string `json:"value,omitempty" bson:"value,omitempty"`
						} `json:"interval,omitempty" bson:"interval,omitempty"`
						Show bool `json:"show,omitempty" bson:"show,omitempty"`
					} `json:"netArea,omitempty" bson:"netArea,omitempty"`
					NetLine struct {
						Color    string `json:"color,omitempty" bson:"color,omitempty"`
						Interval struct {
							CusNumber int64  `json:"cusNumber,omitempty" bson:"cusNumber,omitempty"`
							Value     string `json:"value,omitempty" bson:"value,omitempty"`
						} `json:"interval,omitempty" bson:"interval,omitempty"`
						Show  bool   `json:"show,omitempty" bson:"show,omitempty"`
						Type  string `json:"type,omitempty" bson:"type,omitempty"`
						Width int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"netLine,omitempty" bson:"netLine,omitempty"`
					Show bool `json:"show,omitempty" bson:"show,omitempty"`
					Tick struct {
						Color    string `json:"color,omitempty" bson:"color,omitempty"`
						Length   int64  `json:"length,omitempty" bson:"length,omitempty"`
						Position string `json:"position,omitempty" bson:"position,omitempty"`
						Show     bool   `json:"show,omitempty" bson:"show,omitempty"`
						Width    int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"tick,omitempty" bson:"tick,omitempty"`
					TickLabel struct {
						Digit     string `json:"digit,omitempty" bson:"digit,omitempty"`
						Distance  int64  `json:"distance,omitempty" bson:"distance,omitempty"`
						Formatter struct {
							Digit  string `json:"digit,omitempty" bson:"digit,omitempty"`
							Prefix string `json:"prefix,omitempty" bson:"prefix,omitempty"`
							Ratio  int64  `json:"ratio,omitempty" bson:"ratio,omitempty"`
							Suffix string `json:"suffix,omitempty" bson:"suffix,omitempty"`
						} `json:"formatter,omitempty" bson:"formatter,omitempty"`
						Label struct {
							Color       string        `json:"color,omitempty" bson:"color,omitempty"`
							CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
							FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
							FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
							FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
						} `json:"label,omitempty" bson:"label,omitempty"`
						Max    interface{} `json:"max,omitempty" bson:"max,omitempty"`
						Min    interface{} `json:"min,omitempty" bson:"min,omitempty"`
						Prefix string      `json:"prefix,omitempty" bson:"prefix,omitempty"`
						Ratio  int64       `json:"ratio,omitempty" bson:"ratio,omitempty"`
						Rotate int64       `json:"rotate,omitempty" bson:"rotate,omitempty"`
						Show   bool        `json:"show,omitempty" bson:"show,omitempty"`
						Split  int64       `json:"split,omitempty" bson:"split,omitempty"`
						Suffix string      `json:"suffix,omitempty" bson:"suffix,omitempty"`
					} `json:"tickLabel,omitempty" bson:"tickLabel,omitempty"`
					TickLine struct {
						Color string `json:"color,omitempty" bson:"color,omitempty"`
						Show  bool   `json:"show,omitempty" bson:"show,omitempty"`
						Width int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"tickLine,omitempty" bson:"tickLine,omitempty"`
					Title struct {
						FzPosition string `json:"fzPosition,omitempty" bson:"fzPosition,omitempty"`
						Label      struct {
							Color       string        `json:"color,omitempty" bson:"color,omitempty"`
							CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
							FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
							FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
							FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
						} `json:"label,omitempty" bson:"label,omitempty"`
						NameGap   int64  `json:"nameGap,omitempty" bson:"nameGap,omitempty"`
						Rotate    int64  `json:"rotate,omitempty" bson:"rotate,omitempty"`
						ShowTitle bool   `json:"showTitle,omitempty" bson:"showTitle,omitempty"`
						Text      string `json:"text,omitempty" bson:"text,omitempty"`
					} `json:"title,omitempty" bson:"title,omitempty"`
					Type string `json:"type,omitempty" bson:"type,omitempty"`
				} `json:"yAxisLeft,omitempty" bson:"yAxisLeft,omitempty"`
				YAxisRight struct {
					Inverse bool   `json:"inverse,omitempty" bson:"inverse,omitempty"`
					Name    string `json:"name,omitempty" bson:"name,omitempty"`
					NetArea struct {
						ColorOne string `json:"colorOne,omitempty" bson:"colorOne,omitempty"`
						ColorTwo string `json:"colorTwo,omitempty" bson:"colorTwo,omitempty"`
						Interval struct {
							CusNumber int64  `json:"cusNumber,omitempty" bson:"cusNumber,omitempty"`
							Value     string `json:"value,omitempty" bson:"value,omitempty"`
						} `json:"interval,omitempty" bson:"interval,omitempty"`
						Show bool `json:"show,omitempty" bson:"show,omitempty"`
					} `json:"netArea,omitempty" bson:"netArea,omitempty"`
					NetLine struct {
						Color    string `json:"color,omitempty" bson:"color,omitempty"`
						Interval struct {
							CusNumber int64  `json:"cusNumber,omitempty" bson:"cusNumber,omitempty"`
							Value     string `json:"value,omitempty" bson:"value,omitempty"`
						} `json:"interval,omitempty" bson:"interval,omitempty"`
						Show  bool   `json:"show,omitempty" bson:"show,omitempty"`
						Type  string `json:"type,omitempty" bson:"type,omitempty"`
						Width int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"netLine,omitempty" bson:"netLine,omitempty"`
					Show bool `json:"show,omitempty" bson:"show,omitempty"`
					Tick struct {
						Color    string `json:"color,omitempty" bson:"color,omitempty"`
						Length   int64  `json:"length,omitempty" bson:"length,omitempty"`
						Position string `json:"position,omitempty" bson:"position,omitempty"`
						Show     bool   `json:"show,omitempty" bson:"show,omitempty"`
						Width    int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"tick,omitempty" bson:"tick,omitempty"`
					TickLabel struct {
						Digit     string `json:"digit,omitempty" bson:"digit,omitempty"`
						Distance  int64  `json:"distance,omitempty" bson:"distance,omitempty"`
						Formatter struct {
							Digit  string `json:"digit,omitempty" bson:"digit,omitempty"`
							Prefix string `json:"prefix,omitempty" bson:"prefix,omitempty"`
							Ratio  int64  `json:"ratio,omitempty" bson:"ratio,omitempty"`
							Suffix string `json:"suffix,omitempty" bson:"suffix,omitempty"`
						} `json:"formatter,omitempty" bson:"formatter,omitempty"`
						Label struct {
							Color       string        `json:"color,omitempty" bson:"color,omitempty"`
							CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
							FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
							FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
							FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
						} `json:"label,omitempty" bson:"label,omitempty"`
						Max    interface{} `json:"max,omitempty" bson:"max,omitempty"`
						Min    interface{} `json:"min,omitempty" bson:"min,omitempty"`
						Prefix string      `json:"prefix,omitempty" bson:"prefix,omitempty"`
						Ratio  int64       `json:"ratio,omitempty" bson:"ratio,omitempty"`
						Rotate int64       `json:"rotate,omitempty" bson:"rotate,omitempty"`
						Show   bool        `json:"show,omitempty" bson:"show,omitempty"`
						Split  int64       `json:"split,omitempty" bson:"split,omitempty"`
						Suffix string      `json:"suffix,omitempty" bson:"suffix,omitempty"`
					} `json:"tickLabel,omitempty" bson:"tickLabel,omitempty"`
					TickLine struct {
						Color string `json:"color,omitempty" bson:"color,omitempty"`
						Show  bool   `json:"show,omitempty" bson:"show,omitempty"`
						Width int64  `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"tickLine,omitempty" bson:"tickLine,omitempty"`
					Title struct {
						FzPosition string `json:"fzPosition,omitempty" bson:"fzPosition,omitempty"`
						Label      struct {
							Color       string        `json:"color,omitempty" bson:"color,omitempty"`
							CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
							FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
							FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
							FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
						} `json:"label,omitempty" bson:"label,omitempty"`
						NameGap   int64  `json:"nameGap,omitempty" bson:"nameGap,omitempty"`
						Rotate    int64  `json:"rotate,omitempty" bson:"rotate,omitempty"`
						ShowTitle bool   `json:"showTitle,omitempty" bson:"showTitle,omitempty"`
						Text      string `json:"text,omitempty" bson:"text,omitempty"`
					} `json:"title,omitempty" bson:"title,omitempty"`
				} `json:"yAxisRight,omitempty" bson:"yAxisRight,omitempty"`
			} `json:"axis,omitempty" bson:"axis,omitempty"`
			Config struct {
				Color      string `json:"color,omitempty" bson:"color,omitempty"`
				FontFamily string `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
				Grid       struct {
					Bottom int64  `json:"bottom,omitempty" bson:"bottom,omitempty"`
					Left   int64  `json:"left,omitempty" bson:"left,omitempty"`
					Right  int64  `json:"right,omitempty" bson:"right,omitempty"`
					Top    int64  `json:"top,omitempty" bson:"top,omitempty"`
					Value  string `json:"value,omitempty" bson:"value,omitempty"`
				} `json:"grid,omitempty" bson:"grid,omitempty"`
			} `json:"config,omitempty" bson:"config,omitempty"`
			Legend struct {
				Data     []string `json:"data,omitempty" bson:"data,omitempty"`
				Distance struct {
					CusGap int64  `json:"cusGap,omitempty" bson:"cusGap,omitempty"`
					Value  string `json:"value,omitempty" bson:"value,omitempty"`
				} `json:"distance,omitempty" bson:"distance,omitempty"`
				Height struct {
					CusSize int64  `json:"cusSize,omitempty" bson:"cusSize,omitempty"`
					Value   string `json:"value,omitempty" bson:"value,omitempty"`
				} `json:"height,omitempty" bson:"height,omitempty"`
				ItemGap int64 `json:"itemGap,omitempty" bson:"itemGap,omitempty"`
				Label   struct {
					Color       string        `json:"color,omitempty" bson:"color,omitempty"`
					CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
					FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
					FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
					FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
				} `json:"label,omitempty" bson:"label,omitempty"`
				Position struct {
					Direction string `json:"direction,omitempty" bson:"direction,omitempty"`
					OffsetX   int64  `json:"offsetX,omitempty" bson:"offsetX,omitempty"`
					OffsetY   int64  `json:"offsetY,omitempty" bson:"offsetY,omitempty"`
					Value     string `json:"value,omitempty" bson:"value,omitempty"`
				} `json:"position,omitempty" bson:"position,omitempty"`
				SelectMode string `json:"selectMode,omitempty" bson:"selectMode,omitempty"`
				Selected   []struct {
					IsShow     bool   `json:"isShow,omitempty" bson:"isShow,omitempty"`
					SeriesName string `json:"seriesName,omitempty" bson:"seriesName,omitempty"`
				} `json:"selected,omitempty" bson:"selected,omitempty"`
				Show  bool `json:"show,omitempty" bson:"show,omitempty"`
				Width struct {
					CusSize int64  `json:"cusSize,omitempty" bson:"cusSize,omitempty"`
					Value   string `json:"value,omitempty" bson:"value,omitempty"`
				} `json:"width,omitempty" bson:"width,omitempty"`
			} `json:"legend,omitempty" bson:"legend,omitempty"`
			Subtitle struct {
				Distance struct {
					CusGap int64  `json:"cusGap,omitempty" bson:"cusGap,omitempty"`
					Value  string `json:"value,omitempty" bson:"value,omitempty"`
				} `json:"distance,omitempty" bson:"distance,omitempty"`
				Label struct {
					Color       string        `json:"color,omitempty" bson:"color,omitempty"`
					CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
					FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
					FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
					FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
				} `json:"label,omitempty" bson:"label,omitempty"`
				Show bool   `json:"show,omitempty" bson:"show,omitempty"`
				Text string `json:"text,omitempty" bson:"text,omitempty"`
			} `json:"subtitle,omitempty" bson:"subtitle,omitempty"`
			Title struct {
				Label struct {
					Color       string        `json:"color,omitempty" bson:"color,omitempty"`
					CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
					FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
					FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
					FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
				} `json:"label,omitempty" bson:"label,omitempty"`
				Position struct {
					OffsetX int64  `json:"offsetX,omitempty" bson:"offsetX,omitempty"`
					OffsetY int64  `json:"offsetY,omitempty" bson:"offsetY,omitempty"`
					Value   string `json:"value,omitempty" bson:"value,omitempty"`
				} `json:"position,omitempty" bson:"position,omitempty"`
				Show bool   `json:"show,omitempty" bson:"show,omitempty"`
				Text string `json:"text,omitempty" bson:"text,omitempty"`
			} `json:"title,omitempty" bson:"title,omitempty"`
			Tooltip struct {
				AxisPointer struct {
					Style struct {
						Color string `json:"color,omitempty" bson:"color,omitempty"`
						Type  string `json:"type,omitempty" bson:"type,omitempty"`
						Width string `json:"width,omitempty" bson:"width,omitempty"`
					} `json:"style,omitempty" bson:"style,omitempty"`
					Type string `json:"type,omitempty" bson:"type,omitempty"`
				} `json:"axisPointer,omitempty" bson:"axisPointer,omitempty"`
				BackgroundColor string `json:"backgroundColor,omitempty" bson:"backgroundColor,omitempty"`
				Format          []struct {
					Digit      string `json:"digit,omitempty" bson:"digit,omitempty"`
					Prefix     string `json:"prefix,omitempty" bson:"prefix,omitempty"`
					Ratio      int64  `json:"ratio,omitempty" bson:"ratio,omitempty"`
					SeriesName string `json:"seriesName,omitempty" bson:"seriesName,omitempty"`
					Suffix     string `json:"suffix,omitempty" bson:"suffix,omitempty"`
				} `json:"format,omitempty" bson:"format,omitempty"`
				Label struct {
					Color       string        `json:"color,omitempty" bson:"color,omitempty"`
					CusFontSize int64         `json:"cusFontSize,omitempty" bson:"cusFontSize,omitempty"`
					FontFamily  string        `json:"fontFamily,omitempty" bson:"fontFamily,omitempty"`
					FontGroup   []interface{} `json:"fontGroup,omitempty" bson:"fontGroup,omitempty"`
					FontSize    int64         `json:"fontSize,omitempty" bson:"fontSize,omitempty"`
				} `json:"label,omitempty" bson:"label,omitempty"`
				Position    string `json:"position,omitempty" bson:"position,omitempty"`
				Show        bool   `json:"show,omitempty" bson:"show,omitempty"`
				TriggerOn   string `json:"triggerOn,omitempty" bson:"triggerOn,omitempty"`
				TriggerType string `json:"triggerType,omitempty" bson:"triggerType,omitempty"`
			} `json:"tooltip,omitempty" bson:"tooltip,omitempty"`
		} `json:"defaultOption,omitempty" bson:"defaultOption,omitempty"`
		RangeArray []struct {
			Column []int64 `json:"column,omitempty" bson:"column,omitempty"`
			Row    []int64 `json:"row,omitempty" bson:"row,omitempty"`
		} `json:"rangeArray,omitempty" bson:"rangeArray,omitempty"`
		RangeColCheck struct {
			Exits bool    `json:"exits,omitempty" bson:"exits,omitempty"`
			Range []int64 `json:"range,omitempty" bson:"range,omitempty"`
		} `json:"rangeColCheck,omitempty" bson:"rangeColCheck,omitempty"`
		RangeConfigCheck bool `json:"rangeConfigCheck,omitempty" bson:"rangeConfigCheck,omitempty"`
		RangeRowCheck    struct {
			Exits bool    `json:"exits,omitempty" bson:"exits,omitempty"`
			Range []int64 `json:"range,omitempty" bson:"range,omitempty"`
		} `json:"rangeRowCheck,omitempty" bson:"rangeRowCheck,omitempty"`
	} `json:"chartOptions,omitempty" bson:"chartOptions,omitempty"`
	ChartID       string `json:"chart_id,omitempty" bson:"chart_id,omitempty"`
	Height        int64  `json:"height,omitempty" bson:"height,omitempty"`
	IsShow        bool   `json:"isShow,omitempty" bson:"isShow,omitempty"`
	Left          int64  `json:"left,omitempty" bson:"left,omitempty"`
	NeedRangeShow bool   `json:"needRangeShow,omitempty" bson:"needRangeShow,omitempty"`
	SheetIndex    string `json:"sheetIndex,omitempty" bson:"sheetIndex,omitempty"`
	Top           int64  `json:"top,omitempty" bson:"top,omitempty"`
	Width         int64  `json:"width,omitempty" bson:"width,omitempty"`
}

type Image struct {
	Border struct {
		Color  string `json:"color,omitempty" bson:"color,omitempty"`
		Radius int64  `json:"radius,omitempty" bson:"radius,omitempty"`
		Style  string `json:"style,omitempty" bson:"style,omitempty"`
		Width  int64  `json:"width,omitempty" bson:"width,omitempty"`
	} `json:"border,omitempty" bson:"border,omitempty"`
	Crop struct {
		Height     int64 `json:"height,omitempty" bson:"height,omitempty"`
		OffsetLeft int64 `json:"offsetLeft,omitempty" bson:"offsetLeft,omitempty"`
		OffsetTop  int64 `json:"offsetTop,omitempty" bson:"offsetTop,omitempty"`
		Width      int64 `json:"width,omitempty" bson:"width,omitempty"`
	} `json:"crop,omitempty" bson:"crop,omitempty"`
	Default struct {
		Height int64 `json:"height,omitempty" bson:"height,omitempty"`
		Left   int64 `json:"left,omitempty" bson:"left,omitempty"`
		Top    int64 `json:"top,omitempty" bson:"top,omitempty"`
		Width  int64 `json:"width,omitempty" bson:"width,omitempty"`
	} `json:"default,omitempty" bson:"default,omitempty"`
	FixedLeft    int64  `json:"fixedLeft,omitempty" bson:"fixedLeft,omitempty"`
	FixedTop     int64  `json:"fixedTop,omitempty" bson:"fixedTop,omitempty"`
	IsFixedPos   bool   `json:"isFixedPos,omitempty" bson:"isFixedPos,omitempty"`
	OriginHeight int64  `json:"originHeight,omitempty" bson:"originHeight,omitempty"`
	OriginWidth  int64  `json:"originWidth,omitempty" bson:"originWidth,omitempty"`
	Src          string `json:"src,omitempty" bson:"src,omitempty"`
	Type         string `json:"type,omitempty" bson:"type,omitempty"`
}

type PivotTable struct {
	Column []struct {
		Fullname string `json:"fullname,omitempty" bson:"fullname,omitempty"`
		Index    int64  `json:"index,omitempty" bson:"index,omitempty"`
		Name     string `json:"name"`
	} `json:"column,omitempty" bson:"column,omitempty"`
	DrawPivotTable      bool          `json:"drawPivotTable,omitempty" bson:"drawPivotTable,omitempty"`
	Filter              []interface{} `json:"filter,omitempty" bson:"filter,omitempty"`
	PivotDataSheetIndex int64         `json:"pivotDataSheetIndex,omitempty" bson:"pivotDataSheetIndex,omitempty"`
	PivotDatas          [][]int64     `json:"pivotDatas,omitempty" bson:"pivotDatas,omitempty"`
	PivotTableBoundary  []int64       `json:"pivotTableBoundary,omitempty" bson:"pivotTableBoundary,omitempty"`
	PivotSelectSave     struct {
		Column []int64 `json:"column,omitempty" bson:"column,omitempty"`
		Row    []int64 `json:"row,omitempty" bson:"row,omitempty"`
	} `json:"pivot_select_save,omitempty" bson:"pivot_select_save,omitempty"`
	Row []struct {
		Fullname string `json:"fullname,omitempty" bson:"fullname,omitempty"`
		Index    int64  `json:"index,omitempty" bson:"index,omitempty"`
		Name     string `json:"name,omitempty" bson:"name,omitempty"`
	} `json:"row,omitempty" bson:"row,omitempty"`
	ShowType string `json:"showType,omitempty" bson:"showType,omitempty"`
	Values   []struct {
		Fullname  string `json:"fullname,omitempty" bson:"fullname,omitempty"`
		Index     int64  `json:"index,omitempty" bson:"index,omitempty"`
		Name      string `json:"name,omitempty" bson:"name,omitempty"`
		Nameindex int64  `json:"nameindex,omitempty" bson:"nameindex,omitempty"`
		Sumtype   string `json:"sumtype,omitempty" bson:"sumtype,omitempty"`
	} `json:"values,omitempty" bson:"values,omitempty"`
}

type CalcChain struct {
	C       int           `json:"c,omitempty" bson:"c,omitempty"`
	Chidren struct{}      `json:"chidren,omitempty" bson:"chidren,omitempty"`
	Color   string        `json:"color,omitempty" bson:"color,omitempty"`
	Func    []interface{} `json:"func,omitempty" bson:"func,omitempty"`
	Index   string        `json:"index,omitempty" bson:"index,omitempty"`
	Parent  interface{}   `json:"parent,omitempty" bson:"parent,omitempty"`
	R       int           `json:"r,omitempty" bson:"r,omitempty"`
	Times   int           `json:"times,omitempty" bson:"times,omitempty"`
}

type DynamicArray struct {
	C    int64     `json:"c,omitempty" bson:"c,omitempty"`
	Data [][]int64 `json:"data,omitempty" bson:"data,omitempty"`
	F    string    `json:"f,omitempty" bson:"f,omitempty"`
	R    int64     `json:"r,omitempty" bson:"r,omitempty"`
}
