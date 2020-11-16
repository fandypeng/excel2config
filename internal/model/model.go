package model

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

type Article struct {
	ID      int64
	Content string
	Author  string
}

type MergeCell struct {
	C  int64 `json:"c,omitempty"`
	Cs int64 `json:"cs,omitempty"`
	R  int64 `json:"r,omitempty"`
	Rs int64 `json:"rs,omitempty"`
}

type Cell struct {
	C int64 `json:"c"`
	R int64 `json:"r"`
	V struct {
		Bg string   `json:"bg,omitempty"` // 背景颜色
		Bl int64    `json:"bl,omitempty"` // 0 常规 、 1加粗
		Cl int64    `json:"cl,omitempty"` // 删除线
		Ct struct { // celltype 单元格值格式
			Fa string `json:"fa,omitempty"`
			S  []struct {
				Bl int64       `json:"bl,omitempty"`
				Cl int64       `json:"cl,omitempty"`
				Fc string      `json:"fc,omitempty"`
				Ff string      `json:"ff,omitempty"`
				Fs interface{} `json:"fs,omitempty"`
				It int64       `json:"it,omitempty"`
				Un int64       `json:"un,omitempty"`
				V  string      `json:"v,omitempty"`
			} `json:"s,omitempty"`
			T string `json:"t,omitempty"`
		} `json:"ct,omitempty"`
		F  string      `json:"f,omitempty"`  // 公式
		Fc string      `json:"fc,omitempty"` // 字体颜色
		Ff int64       `json:"ff,omitempty"` // 字体类型
		Fs interface{} `json:"fs,omitempty"` // 字体大小
		It int64       `json:"it,omitempty"` // 斜体
		M  string      `json:"m,omitempty"`  // 显示值
		Mc *MergeCell  `json:"mc,omitempty"` // 合并单元格
		Tb string      `json:"tb,omitempty"` // 文本换行，0 截断、1溢出、2 自动换行
		Tr interface{} `json:"tr,omitempty"` // 竖排文字
		V  interface{} `json:"v,omitempty"`  // 原始值
		Ht interface{} `json:"ht,omitempty"` // 水平对齐，0 居中、1 左、2右
		Vt interface{} `json:"vt,omitempty"` // 垂直对齐，0 中间、1 上、2下
		Rt interface{} `json:"rt,omitempty"` // 文字旋转角度
		Ps struct {
			Height int    `json:"height,omitempty"`
			Width  int    `json:"width,omitempty"`
			Left   int    `json:"left,omitempty"`
			Top    int    `json:"top,omitempty"`
			IsShow bool   `json:"isshow,omitempty"`
			Value  string `json:"value,omitempty"`
		} `json:"ps,omitempty"` //批注
	} `json:"v,omitempty"`
}

type Sheet struct {
	CalcChain []struct {
		C       int64         `json:"c,omitempty"`
		Chidren struct{}      `json:"chidren,omitempty"`
		Color   string        `json:"color,omitempty"`
		Func    []interface{} `json:"func,omitempty"`
		Index   string        `json:"index,omitempty"`
		Parent  interface{}   `json:"parent,omitempty"`
		R       int64         `json:"r,omitempty"`
		Times   int64         `json:"times,omitempty"`
	} `json:"calcChain,omitempty"`
	Celldata []Cell `json:"celldata,omitempty"`
	ChWidth  int64  `json:"ch_width,omitempty"`
	Column   int64  `json:"column,omitempty"`
	Config   struct {
		BorderInfo []struct {
			BorderType string `json:"borderType,omitempty"`
			Color      string `json:"color,omitempty"`
			Range      []struct {
				Column []int64 `json:"column,omitempty"`
				Row    []int64 `json:"row,omitempty"`
			} `json:"range,omitempty"`
			RangeType string `json:"rangeType,omitempty"`
			Style     string `json:"style,omitempty"`
			Value     struct {
				B struct {
					Color string      `json:"color,omitempty"`
					Style interface{} `json:"style,omitempty"`
				} `json:"b,omitempty"`
				ColIndex int64 `json:"col_index,omitempty"`
				L        struct {
					Color string      `json:"color,omitempty"`
					Style interface{} `json:"style,omitempty"`
				} `json:"l,omitempty"`
				R struct {
					Color string      `json:"color,omitempty"`
					Style interface{} `json:"style,omitempty"`
				} `json:"r,omitempty"`
				RowIndex int64 `json:"row_index,omitempty"`
				T        struct {
					Color string      `json:"color,omitempty"`
					Style interface{} `json:"style,omitempty"`
				} `json:"t,omitempty"`
			} `json:"value,omitempty"`
		} `json:"borderInfo,omitempty"`
		Columnlen struct {
			Zero  int64 `json:"0,omitempty"`
			One0  int64 `json:"10,omitempty"`
			Two   int64 `json:"2,omitempty"`
			Three int64 `json:"3,omitempty"`
			Four  int64 `json:"4,omitempty"`
			Five  int64 `json:"5,omitempty"`
			Six   int64 `json:"6,omitempty"`
			Seven int64 `json:"7,omitempty"`
			Eight int64 `json:"8,omitempty"`
			Nine  int64 `json:"9,omitempty"`
		} `json:"columnlen,omitempty"`
		CurentsheetView string `json:"curentsheetView,omitempty"`
		CustomHeight    struct {
			Two9 int64 `json:"29,omitempty"`
		} `json:"customHeight,omitempty"`
		CustomWidth struct {
			Two int64 `json:"2,omitempty"`
		} `json:"customWidth,omitempty"`
		Merge struct {
			One3_5 struct {
				C  int64 `json:"c,omitempty"`
				Cs int64 `json:"cs,omitempty"`
				R  int64 `json:"r,omitempty"`
				Rs int64 `json:"rs,omitempty"`
			} `json:"13_5,omitempty"`
			One3_7 struct {
				C  int64 `json:"c,omitempty"`
				Cs int64 `json:"cs,omitempty"`
				R  int64 `json:"r,omitempty"`
				Rs int64 `json:"rs,omitempty"`
			} `json:"13_7,omitempty"`
			One4_2 struct {
				C  int64 `json:"c,omitempty"`
				Cs int64 `json:"cs,omitempty"`
				R  int64 `json:"r,omitempty"`
				Rs int64 `json:"rs,omitempty"`
			} `json:"14_2,omitempty"`
			One5_10 struct {
				C  int64 `json:"c,omitempty"`
				Cs int64 `json:"cs,omitempty"`
				R  int64 `json:"r,omitempty"`
				Rs int64 `json:"rs,omitempty"`
			} `json:"15_10,omitempty"`
		} `json:"merge,omitempty"`
		Rowhidden struct {
			Three0 int64 `json:"30,omitempty"`
			Three1 int64 `json:"31,omitempty"`
		} `json:"rowhidden,omitempty"`
		Rowlen struct {
			Zero  int64 `json:"0,omitempty"`
			One   int64 `json:"1,omitempty"`
			One0  int64 `json:"10,omitempty"`
			One1  int64 `json:"11,omitempty"`
			One2  int64 `json:"12,omitempty"`
			One3  int64 `json:"13,omitempty"`
			One4  int64 `json:"14,omitempty"`
			One5  int64 `json:"15,omitempty"`
			One6  int64 `json:"16,omitempty"`
			One7  int64 `json:"17,omitempty"`
			One8  int64 `json:"18,omitempty"`
			One9  int64 `json:"19,omitempty"`
			Two   int64 `json:"2,omitempty"`
			Two0  int64 `json:"20,omitempty"`
			Two1  int64 `json:"21,omitempty"`
			Two2  int64 `json:"22,omitempty"`
			Two3  int64 `json:"23,omitempty"`
			Two4  int64 `json:"24,omitempty"`
			Two5  int64 `json:"25,omitempty"`
			Two6  int64 `json:"26,omitempty"`
			Two7  int64 `json:"27,omitempty"`
			Two8  int64 `json:"28,omitempty"`
			Two9  int64 `json:"29,omitempty"`
			Three int64 `json:"3,omitempty"`
			Four  int64 `json:"4,omitempty"`
			Five  int64 `json:"5,omitempty"`
			Six   int64 `json:"6,omitempty"`
			Seven int64 `json:"7,omitempty"`
			Eight int64 `json:"8,omitempty"`
			Nine  int64 `json:"9,omitempty"`
		} `json:"rowlen,omitempty"`
		SheetViewZoom struct {
			ViewLayoutZoomScale int64   `json:"viewLayoutZoomScale,omitempty"`
			ViewNormalZoomScale int64   `json:"viewNormalZoomScale,omitempty"`
			ViewPageZoomScale   float64 `json:"viewPageZoomScale,omitempty"`
		} `json:"sheetViewZoom,omitempty"`
	} `json:"config,omitempty"`
	Index                string `json:"index,omitempty"`
	LuckysheetSelectSave []struct {
		Column      []int64 `json:"column,omitempty"`
		ColumnFocus int64   `json:"column_focus,omitempty"`
		Height      int64   `json:"height,omitempty"`
		HeightMove  int64   `json:"height_move,omitempty"`
		Left        int64   `json:"left,omitempty"`
		LeftMove    int64   `json:"left_move,omitempty"`
		Row         []int64 `json:"row,omitempty"`
		RowFocus    int64   `json:"row_focus,omitempty"`
		Top         int64   `json:"top,omitempty"`
		TopMove     int64   `json:"top_move,omitempty"`
		Width       int64   `json:"width,omitempty"`
		WidthMove   int64   `json:"width_move,omitempty"`
	} `json:"luckysheet_select_save,omitempty"`
	Name       string `json:"name,omitempty"`
	Order      string `json:"order,omitempty"`
	RhHeight   int64  `json:"rh_height,omitempty"`
	Row        int64  `json:"row,omitempty"`
	ScrollLeft int64  `json:"scrollLeft,omitempty"`
	ScrollTop  int64  `json:"scrollTop,omitempty"`
	Status     int64  `json:"status,omitempty"`
	ZoomRatio  int64  `json:"zoomRatio,omitempty"`
}
