package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os/exec"

	_ "github.com/sniperkit/xanalyze/plugin/distribute/glow/driver"
	"github.com/sniperkit/xanalyze/plugin/distribute/glow/flow"
)

type MapperFunc func(shardFilename string, out chan flow.KeyValue)
type ReducerFunc func(key interface{}, values interface{}) string

func createStreamingMapper(mapperExe string) MapperFunc {
	return func(shardFilename string, out chan flow.KeyValue) {
		cmd := exec.Command(mapperExe, shardFilename)
		var outBuffer, errBuffer bytes.Buffer
		cmd.Stdout = &outBuffer
		cmd.Stderr = &errBuffer
		if err := cmd.Run(); err != nil {
			fmt.Printf(string(errBuffer.Bytes()))
			panic(fmt.Sprintf("%v", err))
		}
		var mapperIface interface{}
		errUnmarshal := json.Unmarshal(outBuffer.Bytes(), &mapperIface)
		if errUnmarshal != nil {
			panic(fmt.Sprintf("%v", errUnmarshal))
		}
		mapperItems := mapperIface.([]interface{})
		for _, itemIface := range mapperItems {
			item := itemIface.(map[string]interface{})
			out <- flow.KeyValue{item["Key"], item["Value"]}
		}
	}
}

type ReducerOutput struct {
	Key   interface{}
	Value interface{}
}

func createStreamingReducer(reducerExe string) ReducerFunc {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("Failed to generate reducer id: %v", err))
	}
	reducerId := fmt.Sprintf("%x%x", b[0:4], b[4:8])
	fmt.Printf("reducer id=%v", reducerId)
	return func(key interface{}, values interface{}) string {
		inputBytes, err := json.Marshal(values)
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}

		cmd := exec.Command(reducerExe, fmt.Sprintf("%v", key), reducerId)
		var outBuffer, errBuffer bytes.Buffer
		cmd.Stdout = &outBuffer
		cmd.Stderr = &errBuffer
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Printf(string(errBuffer.Bytes()))
			panic(fmt.Sprintf("%v", err))
		}

		if err := cmd.Start(); err != nil {
			fmt.Printf(string(errBuffer.Bytes()))
			panic(fmt.Sprintf("%v", err))
		}
		io.Copy(stdin, bytes.NewBuffer(inputBytes))
		stdin.Close()
		if err := cmd.Wait(); err != nil {
			fmt.Printf(string(errBuffer.Bytes()))
			panic(fmt.Sprintf("%v", err))
		}

		if len(outBuffer.Bytes()) == 0 {
			return ""
		}
		var reducerIface interface{}
		if err := json.Unmarshal(outBuffer.Bytes(), &reducerIface); err != nil {
			panic(fmt.Sprintf("%v", err))
		}

		reducerOutput := ReducerOutput{key, reducerIface}
		reducerOutputBytes, err := json.MarshalIndent(reducerOutput, "", "  ")
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}
		return string(reducerOutputBytes)
	}
}

func main() {
	var shardFile = flag.String("shard-file", "", "File-catalog")
	var numShards = flag.Int("num-splits", 1, "How many parts divide input shards")
	var mapper = flag.String("mapper", "", "Executable script with map function")
	var reducer = flag.String("reducer", "", "Executed script with reduce function")
	flag.Parse()

	f := flow.New()
	mapperFunc := createStreamingMapper(*mapper)
	reducerFunc := createStreamingReducer(*reducer)
	f.TextFile(*shardFile, *numShards).Map(mapperFunc).GroupByKey().Map(
		func(key interface{}, values interface{}) {
			reducedValue := reducerFunc(key, values)
			if reducedValue != "" {
				fmt.Printf(reducedValue)
			}
		})
	f.Run()
}
