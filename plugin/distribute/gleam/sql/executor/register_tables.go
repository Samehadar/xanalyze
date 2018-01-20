package executor

import (
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/flow"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/sql/model"
)

type TableColumn struct {
	ColumnName string
	ColumnType byte
}

type TableSource struct {
	Dataset   *flow.Dataset
	TableInfo *model.TableInfo
}

var (
	Tables = make(map[string]*TableSource)
)
