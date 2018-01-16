package util

import (
	"fmt"

	"github.com/chrislusf/gleam/flow"
	"github.com/sniperkit/xanalyze/model"

	jsoniter "github.com/sniperkit/xutil/plugin/format/json"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func MockTree(context *flow.FlowContext) (ret *flow.Dataset) {

	fileList := make([]model.TreeEntry, 1000)

	var data model.Tree
	err := json.Unmarshal([]byte(gitTree), &data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	for k, v := range data.Entries {
		fileList[k] = v
	}

	input := make(chan interface{})

	go func() {
		for _, data := range fileList {
			input <- data
		}
		close(input)
	}()

	ret = context.Channel(input)
	return
}
