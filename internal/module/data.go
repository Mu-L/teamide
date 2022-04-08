package module

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"teamide/internal/base"
	"teamide/pkg/application/model"
	"teamide/pkg/toolbox"
)

type DataRequest struct {
	Origin   string `json:"origin,omitempty"`
	Pathname string `json:"pathname,omitempty"`
}

type DataResponse struct {
	Url           string                `json:"url,omitempty"`
	Api           string                `json:"api,omitempty"`
	IsStandAlone  bool                  `json:"isStandAlone,omitempty"`
	ColumnTypes   []*model.ColumnType   `json:"columnTypes,omitempty"`
	DataTypes     []*model.DataType     `json:"dataTypes,omitempty"`
	IndexTypes    []*model.IndexType    `json:"indexTypes,omitempty"`
	ModelTypes    []*model.ModelType    `json:"modelTypes,omitempty"`
	DataPlaces    []*model.DataPlace    `json:"dataPlaces,omitempty"`
	DatabaseTypes []*model.DatabaseType `json:"databaseTypes,omitempty"`
	ToolboxTypes  []*toolbox.Worker     `json:"toolboxTypes,omitempty"`
}

func (this_ *Api) apiData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	path := requestBean.Path[0:strings.LastIndex(requestBean.Path, "api/")]
	request := &DataRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	response := &DataResponse{}

	pathname := request.Pathname
	re, _ := regexp.Compile("/+")
	pathname = re.ReplaceAllLiteralString(pathname, "/")

	path = strings.TrimSuffix(path, "/")
	pathname = strings.TrimSuffix(pathname, "/")
	pathname = strings.TrimSuffix(pathname, path)

	if !strings.HasSuffix(pathname, "/") {
		pathname += "/"
	}

	response.Url = request.Origin + pathname
	response.Api = response.Url + "api/"
	response.IsStandAlone = this_.IsStandAlone
	response.ColumnTypes = model.COLUMN_TYPES
	response.DataTypes = model.DATA_TYPES
	response.ModelTypes = model.MODEL_TYPES
	response.IndexTypes = model.INDEX_TYPES
	response.DataPlaces = model.DATA_PLACES
	response.DatabaseTypes = model.DATABASE_TYPES
	response.ToolboxTypes = toolbox.GetWorkers()

	res = response
	return
}