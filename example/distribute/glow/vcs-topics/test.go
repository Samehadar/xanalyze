package main

/*
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	// jsoniter "github.com/sniperkit/xutil/plugin/format/json"
)

// var json    = jsoniter.ConfigCompatibleWithStandardLibrary

// {"name":"admin-on-rest","owner":"marmelab","path":"src/mui/detail/Tab.js","remote_id":"63226588"}
type File struct {
	name      string `json:"name"`
	owner     string `json:"owner"`
	path      string `json:"path"`
	remote_id int    `json:"remote_id"`
}

type Entries struct {
	Entry []File
}

func main() {
	file, e := ioutil.ReadFile("./files.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	// fmt.Printf("%s\n", string(file))

	// var jsontype []File
	var arrResult []map[string]interface{}
	// var jsontype Entries
	err := json.Unmarshal(file, &arrResult)
	if err != nil {
		panic(err)
	}

	for _, entry := range arrResult {
		// if entry.path != "" {
		fmt.Println("*** entry=", entry)
		// fmt.Println("*** entry.path=", entry.path)
		// out <- entry.path
		//}
	}

	// fmt.Printf("Results: %v\n", jsontype)
}

*/
