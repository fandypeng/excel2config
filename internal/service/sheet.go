package service

import (
	"context"
	"encoding/json"
	pb "excel2config/api"
	"excel2config/internal/def"
	"excel2config/internal/helper"
	"excel2config/internal/model"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

func (s *Service) ExcelList(ctx context.Context, req *pb.ExcelListReq) (resp *pb.ExcelListResp, err error) {
	list, err := s.dao.ExcelList(ctx, req.LastTime, req.Limit, req.GroupId)
	if err != nil {
		return nil, err
	}
	uids := make([]string, 0)
	for _, info := range list {
		if !helper.Contains(uids, info.Owner) {
			uids = append(uids, info.Owner)
		}
	}
	userInfos, err := s.dao.GetUserInfos(ctx, uids)
	if err != nil {
		return
	}
	for _, info := range list {
		if userInfo, ok := userInfos[info.Owner]; ok {
			info.Owner = userInfo.UserName
		}
	}
	resp = &pb.ExcelListResp{
		List: list,
	}
	return resp, nil
}

func (s *Service) CreateExcel(ctx context.Context, req *pb.CreateExcelReq) (resp *pb.CreateExcelResp, err error) {
	groupInfo, err := s.dao.GroupInfo(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	if !groupInfo.IsDev {
		err = ecode.Int(int(def.ErrCannotCreateExcel))
		return
	}
	eid, err := s.dao.CreateExcel(ctx, req.Uid, req.Name, req.Remark, req.GroupId)
	if err == nil {
		_, err = s.dao.CreateExcel(ctx, req.Uid, req.Name, req.Remark, groupInfo.UnionGroupId)
	}
	if err != nil {
		return nil, err
	}
	//bytes, err := json.Marshal(sheets)
	resp = &pb.CreateExcelResp{
		Eid: eid,
	}
	return resp, nil
}

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

func (s *Service) UpdateExcel(ctx context.Context, req *pb.UpdateExcelReq) (resp *pb.UpdateExcelResp, err error) {
	err = s.dao.UpdateExcel(ctx, req.Id, req.Remark, req.Contributers)
	if err != nil {
		return nil, err
	}
	resp = &pb.UpdateExcelResp{}
	return resp, nil
}

func (s *Service) DeleteExcel(ctx context.Context, req *pb.DeleteExcelReq) (resp *pb.DeleteExcelResp, err error) {
	excelInfo, err := s.dao.ExcelInfo(ctx, req.Id)
	if err != nil {
		log.Errorw(ctx, "gridKey not exist, gridKey: ", req.Id, "name: ", req.Name)
		err = ecode.Int(int(def.ErrTableNotExist))
		return
	}
	groupInfo, err := s.dao.GroupInfo(ctx, excelInfo.GroupId)
	if err != nil {
		return nil, err
	}
	err = s.dao.DeleteExcel(ctx, req.Id)
	if err == nil {
		excelInfo, err = s.dao.ExcelInfoByGroupId(ctx, groupInfo.UnionGroupId, excelInfo.Name)
	}
	if err == nil {
		err = s.dao.DeleteExcel(ctx, excelInfo.Id)
	}
	if err != nil {
		return nil, err
	}
	resp = &pb.DeleteExcelResp{}
	return resp, nil
}

func (s *Service) ExportExcel(ctx context.Context, req *pb.ExportExcelReq) (resp *pb.ExportExcelResp, err error) {
	sheet, err := s.dao.LoadSheetByName(ctx, req.GridKey, req.SheetName)
	if err != nil {
		return
	}
	resp = &pb.ExportExcelResp{}
	res, err := sheet.Format()
	if err != nil {
		err = ecode.Int(int(def.ErrTableFormat))
		return
	}
	b, err := json.Marshal(res.Content)
	if err == nil {
		resp.Jsonstr = string(b)
	}
	return
}

// ExportProdExcel 导出正式环境的sheet，根据gridKey查找到测试环境的excel，再根据union_group_id和excelInfo.name找到正式环境的excelInfo
// 最后根据正式环境的excelInfo找到正式环境的gridKey
func (s *Service) ExportProdExcel(ctx context.Context, req *pb.ExportProdExcelReq) (resp *pb.ExportProdExcelResp, err error) {
	prodExcelInfo, err := s.getProdExcelInfo(ctx, req.GridKey, req.Gid)
	if err != nil {
		return
	}
	sheet, err := s.dao.LoadSheetByName(ctx, prodExcelInfo.Id, req.SheetName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			resp = &pb.ExportProdExcelResp{}
			return
		}
		return
	}
	resp = &pb.ExportProdExcelResp{}
	res, err := sheet.Format()
	if err != nil {
		err = ecode.Int(int(def.ErrTableFormat))
		return
	}
	b, err := json.Marshal(res.Content)
	if err == nil {
		resp.Jsonstr = string(b)
	}
	return
}

func (s *Service) SheetList(ctx context.Context, req *pb.SheetListReq) (resp *pb.SheetListResp, err error) {
	sheets, err := s.dao.LoadAllSheet(ctx, req.GridKey)
	if err != nil {
		err = ecode.Int(int(def.ErrTableNotExist))
		return
	}
	sheetNames := make([]string, 0)
	for _, sheet := range sheets {
		sheetNames = append(sheetNames, sheet.Name)
	}
	resp = &pb.SheetListResp{
		SheetName: sheetNames,
	}
	return
}

func (s *Service) getProdExcelInfo(ctx context.Context, gridKey, gid string) (prodExcelInfo *model.Excel, err error) {
	excelInfo, err := s.dao.ExcelInfo(ctx, gridKey)
	if err != nil {
		return
	}
	groupInfo, err := s.dao.GroupInfo(ctx, gid)
	if err != nil {
		return
	}
	prodExcelInfo, err = s.dao.ExcelInfoByGroupId(ctx, groupInfo.UnionGroupId, excelInfo.Name)
	if err != nil {
		return
	}
	return
}
