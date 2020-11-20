package service

import (
	"context"
	"encoding/json"
	pb "excel2config/api"
	"strings"
)

func (s *Service) LoadExcel(ctx context.Context, req *pb.LoadExcelReq) (resp *pb.LoadExcelResp, err error) {
	sheets, err := s.dao.LoadExcel(ctx, req.GridKey)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(sheets)
	resp = &pb.LoadExcelResp{
		Content: string(bytes),
	}
	return resp, nil
}

func (s *Service) LoadExcelSheet(ctx context.Context, req *pb.LoadExcelSheetReq) (resp *pb.LoadExcelSheetResp, err error) {
	indexs := strings.Split(req.Indexs, ",")
	sheets, err := s.dao.LoadExcelSheet(ctx, req.GridKey, indexs)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(sheets)
	resp = &pb.LoadExcelSheetResp{
		Content: string(bytes),
	}
	return resp, nil
}
