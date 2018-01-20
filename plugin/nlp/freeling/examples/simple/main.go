package main

import (
	"encoding/json"
	"fmt"

	. "github.com/sniperkit/xanalyze/plugin/nlp/freeling/lib"
	. "github.com/sniperkit/xanalyze/plugin/nlp/freeling/models"
)

func main() {
	document := new(DocumentEntity)
	analyzer := NewAnalyzer()
	document.Content = "Hello World"
	output := analyzer.AnalyzeText(document)

	js := output.ToJSON()
	b, err := json.Marshal(js)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
