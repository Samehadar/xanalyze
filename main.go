package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	// _ "strings"
	"strings"
	"time"

	"github.com/chrislusf/gleam/distributed"
	"github.com/chrislusf/gleam/flow"
	"github.com/chrislusf/gleam/gio"
	"github.com/chrislusf/gleam/plugins/file"
	_ "github.com/sniperkit/xanalyze/model"
	// "github.com/sniperkit/xanalyze/util"

	jsoniter "github.com/sniperkit/xutil/plugin/format/json"
)

var (
	MapperTokenizer      = gio.RegisterMapper(tokenize)
	MapperAddOne         = gio.RegisterMapper(addOne)
	registeredReadConent = gio.RegisterMapper(readContent)
	registeredTfIdf      = gio.RegisterMapper(tfidf)

	isDistributed   = flag.Bool("distributed", false, "run in distributed or not")
	isDockerCluster = flag.Bool("onDocker", false, "run in docker cluster")
)

type SourcePost struct {
	PostTypeId   int    `xml:"PostTypeId,attr"`
	CreationDate string `xml:"CreationDate,attr"` // "2008-07-31T21:42:52.667"
	Tags         string `xml:"Tags,attr"`
}

type Post struct {
	PostTypeId   int
	CreationDate time.Time
	Tags         []string
}

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
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
	f := flow.New() //"word detection in readmes with pipes")

	// f.Datasets

	// dataset := util.MockTree(f)
	// dataset := util.Mock(f)

	dataset := f.Read(file.Txt("Posts100k.xml", 1))

	dataset.
		Map("tokenize", MapperTokenizer).
		Filter(func(line string) bool {
			return strings.Contains(line, "<row")
		}).
		Map(func(line string, ch chan SourcePost) {
			var p SourcePost
			err := xml.Unmarshal([]byte(line), &p)
			if err != nil {
				fmt.Printf("parse source post error %v: %s\n", err, line)
				return
			}
			ch <- p
		}).Filter(func(src SourcePost) bool {
		return src.PostTypeId == 1
	}).
		Map(func(src SourcePost) (p Post) {
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
		}).
		Fprintf(os.Stdout, "%s\t%d")

	// Mapper(MapperTokenizer).Fprintf(os.Stdout, "%s\t%d")

	/*
		}, 4).Filter(func(line string) bool {

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
	*/

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
	if visit, ok := row[0].(map[interface{}]interface{}); ok {
		for k, v := range visit {
			var w, t string
			var ok bool
			if w, ok = k.(string); !ok {
				continue
			}
			if t, ok = v.(string); !ok {
				continue
			}
			println("key:", w, "v:", t)
		}
	}
	return nil
}

func addOne(row []interface{}) error {
	word := string(row[0].([]byte))
	gio.Emit(word, 1)
	return nil
}

func readContent(x []interface{}) error {

	filepath := gio.ToString(x[0])

	f, err := os.Open(filepath)
	if err != nil {
		println("error reading file:", filepath)
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		gio.Emit(scanner.Text(), filepath, 1)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	return nil
}

func tfidf(x []interface{}) error {
	fmt.Fprintf(os.Stderr, "tfidf input: %v\n", x)
	word := gio.ToString(x[0])
	df := uint16(gio.ToInt64(x[1]))
	doc := gio.ToString(x[2])
	tf := uint16(gio.ToInt64(x[3]))

	gio.Emit(word, doc, tf, df, float32(tf)/float32(df))
	return nil
}
