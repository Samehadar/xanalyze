package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	// _ "github.com/sniperkit/xanalyze/model"
	jsoniter "github.com/sniperkit/xutil/plugin/format/json"

	_ "github.com/sniperkit/xanalyze/plugin/distribute/glow/driver"
	"github.com/sniperkit/xanalyze/plugin/distribute/glow/flow"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	f    = flow.New()
)

//eg: {"name":"admin-on-rest","owner":"marmelab","path":"src/mui/detail/Tab.js","remote_id":"63226588"}
type File struct {
	name      string
	owner     string
	path      string
	remote_id int
}

type Manifest struct {
	manager string
	path    string
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func init() {

	// context

	files := f.TextFile(
		"files.json", 4,
	).Filter(func(line string) bool {
		return strings.Contains(line, "\"path\"")
	}).Map(func(line string, ch chan File) {
		var f File
		if strings.HasPrefix(line, ",") {
			line = after(line, ",")
		}
		err := json.Unmarshal([]byte(line), &f)
		if err != nil {
			fmt.Printf("parse source post error %v: %s\n", err, line)
			return
		}
		ch <- f
	}).Filter(func(f File) bool {
		return f.path != ""
	}).Map(func(f File) (m Manifest) {
		m.path = f.path

		// check for kewyords

		return
	})

	/*
		files.Map(func(p Post, out chan flow.KeyValue) {
			if len(p.Tags) > 1 {
				for _, t := range p.Tags {
					out <- flow.KeyValue{t, 1}
				}
			}
		}).ReduceByKey(func(x int, y int) int {
			return x + y
		}).Map(func(tag string, count int) flow.KeyValue {
			return flow.KeyValue{count, tag}
		}).Sort(func(a, b int) bool {
			return a < b
		}).Map(func(count int, tag string) {
			fmt.Printf("%d %s\n", count, tag)
		})

		files.Map(func(p Post) flow.KeyValue {
			return flow.KeyValue{p.CreationDate.Format("2006-01"), 1}
		}).ReduceByKey(func(x int, y int) int {
			return x + y
		}).Sort(nil).Map(func(month string, count int) {
			fmt.Printf("%s %d\n", month, count)
		})
	*/

}

func main() {
	flag.Parse()
	flow.Ready()

	f.Run()

}
