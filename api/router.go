package api

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

func RegisterUser(e *bm.Engine, auth Auth, server UserBMServer) {
	UserSvc = server
	u := e.Group("user")
	{
		u.POST("login", userLogin)
		u.POST("reg", userRegister)
		u.POST("logout", userLogout)
	}
	uauth := e.Group("user", auth.NeedLogin())
	{
		uauth.GET("info", userInfo)
		uauth.POST("search", userSearch)
	}
}

// RegisterSheetBMServer Register the blademaster route
func RegisterGroup(e *bm.Engine, auth Auth, server GroupBMServer) {
	GroupSvc = server
	s := e.Group("group", auth.NeedLogin())
	{
		s.GET("list", groupGroupList)
		s.POST("add", groupGroupAdd)
		s.POST("update", groupGroupUpdate)
		s.POST("test_connection", groupTestConnection)
		s.POST("get_config_from_db", groupGetConfigFromDB)
		s.POST("export_config_to_db", groupExportConfigToDB)
		s.GET("gen_app_key_secret", groupGenerateAppKeySecret)
		s.POST("sync_to_prod", groupSyncToProd)
		s.POST("export_record", groupExportRecord)
		s.POST("export_record_content", groupExportRecordContent)
		s.POST("export_rollback", groupExportRollback)
	}
}

// RegisterSheetBMServer Register the blademaster route
func RegisterSheet(e *bm.Engine, auth Auth, server SheetBMServer) {
	SheetSvc = server
	s := e.Group("excel", auth.NeedLogin())
	{
		s.POST("", LoadExcel)
		s.POST("sheet", LoadExcelSheet)
		s.GET("list", sheetExcelList)
		s.POST("create", sheetCreateExcel)
		s.POST("update", sheetUpdateExcel)
		s.POST("delete", sheetDeleteExcel)
		s.POST("export", sheetExportExcel)
		s.POST("sheet_list", sheetSheetList)
		s.POST("export_prod", sheetExportProdExcel)
	}
	e.POST("excel/export_all_sheets", sheetExportAllSheets)
}
