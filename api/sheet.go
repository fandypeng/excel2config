package api

import (
	"github.com/gin-gonic/gin/binding"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"net/http"
)

func LoadExcel(c *bm.Context) {
	p := new(LoadExcelReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := SheetSvc.LoadExcel(c, p)
	if err != nil {
		c.String(http.StatusOK, "服务器内部错误")
		return
	}
	c.String(http.StatusOK, resp.Content)
}

func LoadExcelSheet(c *bm.Context) {
	p := new(LoadExcelSheetReq)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := SheetSvc.LoadExcelSheet(c, p)
	if err != nil {
		c.String(http.StatusOK, "服务器内部错误")
		return
	}
	c.String(http.StatusOK, resp.Content)
}
