package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	_ "github.com/sniperkit/xanalyze/model"
	jsoniter "github.com/sniperkit/xutil/plugin/format/json"

	_ "github.com/chrislusf/glow/driver"
	"github.com/chrislusf/glow/flow"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	f    = flow.New()
)

func init() {
	files := f.TextFile(
		"files.json", 4,
	).Filter(func(line string) bool {
		return strings.Contains(line, "<row")
	}).Map(func(line string, ch chan SourcePost) {
		var p SourcePost
		err := xml.Unmarshal([]byte(line), &p)
		if err != nil {
			fmt.Printf("parse source post error %v: %s\n", err, line)
			return
		}
		ch <- p
	}).Filter(func(src SourcePost) bool {
		return src.PostTypeId == 1
	}).Map(func(src SourcePost) (p Post) {
		p.PostTypeId = src.PostTypeId

		t, err := time.Parse("2006-01-02T15:04:05.000", src.CreationDate)
		if err != nil {
			fmt.Printf("error parse creation date %s: %v\n", src.CreationDate, err)
		} else {
			p.CreationDate = t
		}

		if len(src.Tags) > 0 {
			p.Tags = strings.Split(src.Tags[1:len(src.Tags)-1], "><")
		}

		return
	})

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
}

func main() {
	flag.Parse()
	flow.Ready()

	f.Run()

}
