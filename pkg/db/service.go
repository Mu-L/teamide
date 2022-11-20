package db

import (
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"go.uber.org/zap"
	"teamide/pkg/util"
	"time"
)

func CreateService(config *DatabaseConfig) (service *Service, err error) {
	service = &Service{
		config: config,
	}
	service.lastUseTime = util.GetNowTime()
	err = service.init()
	return
}

type SqlParam struct {
	Sql    string        `json:"sql,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

type Service struct {
	config         *DatabaseConfig
	lastUseTime    int64
	DatabaseWorker *DatabaseWorker
}

func (this_ *Service) init() (err error) {
	this_.DatabaseWorker, err = NewDatabaseWorker(this_.config)
	if err != nil {
		return
	}
	return
}

func (this_ *Service) GetDatabaseWorker() *DatabaseWorker {
	return this_.DatabaseWorker
}

func (this_ *Service) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *Service) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *Service) SetLastUseTime() {
	this_.lastUseTime = util.GetNowTime()
}

func (this_ *Service) Stop() {
	if this_.DatabaseWorker != nil {
		_ = this_.DatabaseWorker.Close()
	}
}

func (this_ *Service) OwnersSelect(param *Param) (owners []*dialect.OwnerModel, err error) {
	owners, err = worker.OwnersSelect(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel)
	return
}

func (this_ *Service) TablesSelect(param *Param, ownerName string) (tables []*dialect.TableModel, err error) {
	tables, err = worker.TablesSelect(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel, ownerName)
	return
}

func (this_ *Service) TableDetail(param *Param, ownerName string, tableName string) (tableDetail *dialect.TableModel, err error) {
	tableDetail, err = worker.TableDetail(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel, ownerName, tableName, true)
	return
}

func (this_ *Service) OwnerCreate(param *Param, owner *dialect.OwnerModel) (created bool, err error) {
	created, err = worker.OwnerCreate(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel, owner)
	return
}

func (this_ *Service) OwnerDelete(param *Param, ownerName string) (deleted bool, err error) {
	deleted, err = worker.OwnerDelete(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel, ownerName)
	return
}

func (this_ *Service) DDL(param *Param, ownerName string, tableName string) (sqlList []string, err error) {
	var sqlList_ []string
	if param.AppendOwnerCreateSql {
		var owner *dialect.OwnerModel
		owner, err = worker.OwnerSelect(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel, ownerName)
		if err != nil {
			return
		}
		if owner != nil {
			sqlList_, err = this_.GetTargetDialect(param).OwnerCreateSql(param.ParamModel, owner)
			if err != nil {
				return
			}
			sqlList = append(sqlList, sqlList_...)
		}
	}
	var tables []*dialect.TableModel
	if tableName != "" {
		var table *dialect.TableModel
		table, err = worker.TableDetail(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel, ownerName, tableName, false)
		if err != nil {
			return
		}
		if table != nil {
			tables = append(tables, table)
		}
	} else {
		tables, err = worker.TablesDetail(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel, ownerName, false)
		if err != nil {
			return
		}
	}
	for _, table := range tables {
		var appendOwnerName string
		if param.AppendOwnerName {
			appendOwnerName = ownerName
		}
		sqlList_, err = this_.GetTargetDialect(param).TableCreateSql(param.ParamModel, appendOwnerName, table)
		if err != nil {
			return
		}
		sqlList = append(sqlList, sqlList_...)
	}

	return
}

func (this_ *Service) TableCreate(param *Param, ownerName string, table *dialect.TableModel) (err error) {
	err = worker.TableCreate(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel, ownerName, table)
	return
}

func (this_ *Service) TableDelete(param *Param, ownerName string, tableName string) (err error) {
	err = worker.TableDelete(this_.DatabaseWorker.db, this_.DatabaseWorker.Dialect, param.ParamModel, ownerName, tableName)
	return
}

func (this_ *Service) TableTrim(param *Param, ownerName string, tableName string) (err error) {
	sqlInfo := "DELETE FROM " + this_.DatabaseWorker.OwnerNamePack(param.ParamModel, ownerName) + "." + this_.DatabaseWorker.TableNamePack(param.ParamModel, tableName)
	_, err = worker.DoExec(this_.DatabaseWorker.db, sqlInfo, nil)
	return
}

func (this_ *Service) TableCreateSql(param *Param, ownerName string, table *dialect.TableModel) (sqlList []string, err error) {
	var appendOwnerName string
	if param.AppendOwnerName {
		appendOwnerName = ownerName
	}
	sqlList, err = this_.GetTargetDialect(param).TableCreateSql(param.ParamModel, appendOwnerName, table)

	return
}

func (this_ *Service) TableUpdateSql(param *Param, ownerName string, updateTableParam *UpdateTableParam) (sqlList []string, err error) {
	return
}

func (this_ *Service) TableUpdate(param *Param, ownerName string, updateTableParam *UpdateTableParam) (err error) {
	sqlList, err := this_.TableUpdateSql(param, ownerName, updateTableParam)
	if err != nil {
		return
	}
	_, _, _, err = worker.DoExecs(this_.DatabaseWorker.db, sqlList, nil)
	if err != nil {
		util.Logger.Error("TableUpdate error", zap.Error(err))
		return
	}
	return
}

type UpdateTableParam struct {
	TableName string `json:"tableName"`
}

type DataListResult struct {
	Sql      string                   `json:"sql"`
	Total    int                      `json:"total"`
	Params   []interface{}            `json:"params"`
	DataList []map[string]interface{} `json:"dataList"`
}

func (this_ *Service) DataList(param *Param, ownerName string, tableName string, columnList []*dialect.ColumnModel, whereList []*dialect.Where, orderList []*dialect.Order, pageSize int, pageNo int) (dataListResult DataListResult, err error) {

	sql, values, err := this_.DatabaseWorker.Dialect.DataListSelectSql(param.ParamModel, ownerName, tableName, columnList, whereList, orderList)
	if err != nil {
		return
	}

	page := worker.NewPage()
	page.PageSize = pageSize
	page.PageNo = pageNo
	listMap, err := this_.DatabaseWorker.QueryMapPage(sql, values, page)
	if err != nil {
		return
	}
	for _, one := range listMap {
		for k, v := range one {
			if v == nil {
				continue
			}
			switch tV := v.(type) {
			case time.Time:
				if tV.IsZero() {
					one[k] = nil
				} else {
					one[k] = util.GetTimeTime(tV)
				}
			default:
				one[k] = fmt.Sprint(tV)
			}
		}
	}
	dataListResult.Sql = this_.DatabaseWorker.PackPageSql(sql, page.PageSize, page.PageNo)
	dataListResult.Params = values
	dataListResult.Total = page.TotalCount
	dataListResult.DataList = listMap
	return
}

func (this_ *Service) Execs(sqlList []string, paramsList [][]interface{}) (res int64, err error) {
	res, err = this_.DatabaseWorker.Execs(sqlList, paramsList)
	if err != nil {
		return
	}
	return
}

func (this_ *Service) GetTargetDialect(param *Param) (dia dialect.Dialect) {
	if param != nil && param.TargetDatabaseType != "" {
		t := GetDatabaseType(param.TargetDatabaseType)
		if t != nil {
			return t.dia
		}
	}
	return this_.DatabaseWorker.Dialect
}

type Param struct {
	*dialect.ParamModel
	TargetDatabaseType   string `json:"targetDatabaseType"`
	AppendOwnerCreateSql bool   `json:"appendOwnerCreateSql"`
	AppendOwnerName      bool   `json:"appendOwnerName"`
}
