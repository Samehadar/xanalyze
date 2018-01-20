package executor

import (
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/flow"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/sql/expression"
)

type Executor interface {
	Exec() *flow.Dataset
	Schema() expression.Schema
}
