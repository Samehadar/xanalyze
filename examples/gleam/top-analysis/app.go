package main

import (
	"flag"
	_ "strings"

	"github.com/chrislusf/gleam/distributed"
	"github.com/chrislusf/gleam/flow"
	"github.com/chrislusf/gleam/gio"
	"github.com/zhangweilun/session/util"
	_"github.com/zhangweilun/session/model"
	"os"


)

var (
	MapperTokenizer = gio.RegisterMapper(tokenize)
	//MapperAddOne    = gio.RegisterMapper(addOne)
	//ReducerSum      = gio.RegisterReducer(sum)

	isDistributed   = flag.Bool("distributed", false, "run in distributed or not")
	isDockerCluster = flag.Bool("onDocker", false, "run in docker cluster")
)

func main() {

	flag.Parse() // optional, since gio.Init() will call this also.
	gio.Init()   // If the command line invokes the mapper or reducer, execute it and exit.

	//f := flow.New().TextFile("/etc/passwd").
	//	Mapper(MapperTokenizer). // invoke the registered "tokenize" mapper function.
	//	Mapper(MapperAddOne).    // invoke the registered "addOne" mapper function.
	//	ReducerBy(ReducerSum).   // invoke the registered "sum" reducer function.
	//	Sort(flow.OrderBy(2, true)).
	//	Fprintf(os.Stdout, "%s\t%d\n")
	f := flow.New()

	dataset := util.Mock(f)

	dataset.Mapper(MapperTokenizer).Fprintf(os.Stdout,"%s\t%d")

	if *isDistributed {
		println("Running in distributed mode.")
		f.Run(distributed.Option())
	} else if *isDockerCluster {
		println("Running in docker cluster.")
		f.Run(distributed.Option().SetMaster("master:45326"))
	} else {
		println("Running in standalone mode.")
		f.Run()
	}

}

func tokenize(row []interface{}) error {

	// map[string]string   key是字段名  value是字段的值
	if visit, ok := row[0].(map[interface{}]interface{}) ;ok{
		for k,v:=range visit{
			if w,ok:=k.(string);ok {
				println("key:",w)
			}
			if t,ok:=v.(string);ok {
				println("v:",t)
			}
		}
	}




	return nil
}

//func addOne(row []interface{}) error {
//	word := string(row[0].([]byte))
//
//	gio.Emit(word, 1)
//
//	return nil
//}
//
//func sum(x, y interface{}) (interface{}, error) {
//	return x.(uint64) + y.(uint64), nil
//}