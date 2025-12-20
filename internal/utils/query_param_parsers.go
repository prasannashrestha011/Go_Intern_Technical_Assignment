package utils

import (
	"main/internal/logger"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func ParsePaginationValues(r *http.Request,defaultPage int ,defaultPageSize int)(page int , pageSize int){
	page=defaultPage
	pageSize=defaultPageSize

	query:=r.URL.Query()
	if pageStr:=query.Get("page");pageStr!=""{
		if p,err:=strconv.Atoi(pageStr);err==nil{
			page=p
		}
	}
	if pageSizeStr:=query.Get("pageSize");pageSizeStr!=""{
		if ps,err:=strconv.Atoi(pageSizeStr);err==nil{
			pageSize=ps
		}
	}

	logger.Log.Info("Page details",zap.Int("Page",page),zap.Int("Page Size",pageSize))
	return page,pageSize
}