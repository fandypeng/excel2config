package model

type UpdateGridReq struct {
	C int64     `json:"c"`
	I string    `json:"i"`
	R int64     `json:"r"`
	T string    `json:"t"`
	V CellValue `json:"v"`
}
