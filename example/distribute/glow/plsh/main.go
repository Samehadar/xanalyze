package main

import (
	"flag"
	"fmt"

	"github.com/sniperkit/xanalyze/plugin/distribute/glow/flow"
	plsh "github.com/sniperkit/xanalyze/plugin/distribute/glow/plugin/lsh"
)

var (
	seed         int64  = 10
	p            uint32 = 65537
	nBands       uint32 = 128
	nRowsPerBand uint32 = 4
	nGram               = 2
	lsh                 = plsh.NewLSH(seed, p, nBands, nRowsPerBand)
)

type Pair struct {
	line string
	id   uint32
}

func main() {
	flag.Parse()
	flow.Ready()

	f1 := flow.New()
	f1.TextFile(
		"lsh_input.txt", 3,
	).Map(func(line string, ch chan Pair) {
		fmt.Println("[DEBUG] mapper 1:", line)
		m := NewMinHashValue(lsh.nMinHashFunc, lsh.minHashParams)
		for j := 0; j < len(line)-nGram; j++ {
			m.Update(line[j : j+nGram])
		}
		bucketIds := lsh.GetBucketIds(m)

		for j := 0; j < len(bucketIds); j++ {
			ch <- Pair{line, bucketIds[j]}
		}
	}).Map(func(pair Pair) (int, string) {
		fmt.Println("[DEBUG] mapper 2:", pair.id, pair.line)
		return int(pair.id), pair.line
	}).GroupByKey().Map(func(id int, lines []string) {
		fmt.Printf("[DEBUG] mapper 3: id: %d, line: %s\n", id, lines)
	}).Run()
}
