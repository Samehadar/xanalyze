package main

import (
	"flag"

	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/distributed"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/flow"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/gio"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/gio/mapper"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/gio/reducer"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/plugins/file"
)

var (
	isDistributed = flag.Bool("distributed", false, "run in distributed mode")
)

func main() {

	gio.Init()

	join1()

	hashjoin()

}

func join1() {

	f := flow.New("common words count")

	a := f.Read(file.Txt("../../flow/dataset_map.go", 1)).
		Map("tokenize", mapper.Tokenize).
		Map("addOne", mapper.AppendOne).
		ReduceByKey("sum", reducer.SumInt64)

	b := f.Read(file.Txt("../../flow/dataset_reduce.go", 1)).
		Map("tokenize", mapper.Tokenize).
		Map("addOne", mapper.AppendOne).
		ReduceByKey("sum", reducer.SumInt64)

	join := a.JoinByKey("shared words", b).Printlnf("%s\t%d\t%d")

	println("========== joining result=============")

	if *isDistributed {
		join.Run(distributed.Option())
	} else {
		join.Run()
	}

}

func hashjoin() {

	f := flow.New("hash join")

	a := f.Read(file.Txt("../../flow/dataset_map.go", 1)).
		Map("tokenize", mapper.Tokenize).
		Map("addOne", mapper.AppendOne).
		ReduceByKey("sum", reducer.SumInt64)

	b := f.Strings([]string{
		"func",
		"return",
	})

	a.JoinByKey("hash join", b).Printlnf("%s\t%d")

	println("==========hash joining result=============")

	if *isDistributed {
		f.Run(distributed.Option())
	} else {
		f.Run()
	}

}
