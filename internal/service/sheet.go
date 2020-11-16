package service

import (
	"context"
	"encoding/json"
	pb "excel2config/api"
)

func (s *Service) LoadExcel(ctx context.Context, req *pb.LoadExcelReq) (*pb.LoadExcelResp, error) {
	sheets, err := s.dao.LoadExcel(ctx, req.GridKey)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(sheets)
	resp := &pb.LoadExcelResp{
		Jsonstr: string(bytes),
	}
	return resp, nil
}

func (s *Service) LoadExcelSheet(ctx context.Context, req *pb.LoadExcelSheetReq) (*pb.LoadExcelSheetResp, error) {
	sheets, err := s.dao.LoadExcelSheet(ctx, req.GridKey, req.Indexs)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(sheets)
	resp := &pb.LoadExcelSheetResp{
		Jsonstr: string(bytes),
	}
	return resp, nil
}
