package main

import (
	"flag"
	"path/filepath"
	"strings"

	"os"

	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/flow"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/gio"
	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/plugins/kafka"
)

var (
	brokers = flag.String("brokers", "127.0.0.1:9092", "a list of comma separated broker:port")
	topic   = flag.String("topic", "", "the topic name")
	group   = flag.String("group", filepath.Base(os.Args[0]), "the consumer group name")
	timeout = flag.Int("timeout", 30, "the number of seconds for timeout connections")
)

func main() {

	gio.Init()
	flag.Parse()

	brokerList := strings.Split(*brokers, ",")

	k := kafka.New(brokerList, *topic, *group)
	k.TimeoutSeconds = *timeout

	f := flow.New("kafka " + *topic).Read(k).Printlnf("%x")

	f.Run()

}
