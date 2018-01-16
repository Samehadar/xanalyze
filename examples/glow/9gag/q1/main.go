package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chrislusf/glow/flow"
)

// a struct to represent the record
type Record struct {
	Subreddit    string
	User         string
	Time         time.Time
	NumberOfPost int
}

const (
	Delimiter = ":::"
	ResultDir = "results"
)

var (
	inputPath *string
)

func main() {
	// get the input value from command line
	inputPath = flag.String("input", "", "Absolute input file (e.g. /Users/name/input.csv)")

	flag.Parse()

	// make '-input' become required (validation)
	if *inputPath == "" {
		log.Fatal("The flag '-input' is missing. Please specific CSV file location. (e.g. -input /Users/name/input.csv) ")
	}

	// clean output directory to make sure there is no old data
	if _, err := os.Stat(ResultDir); err == nil {
		os.RemoveAll(ResultDir)
	}

	// count the time use
	startTime := time.Now()

	log.Println("Started please wait... It may take some time...")

	// run the program
	run()

	log.Println("Done! Total time in seconds: ", time.Now().Sub(startTime).Seconds())
	log.Println("The result is in the same directory. Pattern: {subreddit}/{date}.csv")
}

func run() {
	// build the map-reduce pipeline
	flow.New().Source(readCSVInput, runtime.NumCPU()).Filter(func(r Record) bool {
		// remove all the data without subreddit or it is empty
		return r.Subreddit != ""
	}).Map(func(r Record) (string, Record) {
		// use delimiter to group by subreddit, date, and user
		return r.Subreddit + Delimiter + r.Time.Format("2006-01-02") + Delimiter + r.User, r
	}).ReduceByKey(func(x, y Record) Record {
		// count the number of the post of each user within a subreddit and a day
		return Record{
			Subreddit:    x.Subreddit,
			Time:         x.Time,
			User:         x.User,
			NumberOfPost: x.NumberOfPost + y.NumberOfPost,
		}
	}).Map(func(key string, r Record) (string, Record) {
		// break the key from {subreddit}:::{date}:::{user}
		// to {subreddit}:::{date} so that we can group by key with this pattern

		arr := strings.Split(key, Delimiter)
		if len(arr) != 3 {
			log.Fatal("The key: " + key + " is not valid")
		}

		return arr[0] + Delimiter + arr[1], r
	}).GroupByKey().Map(writeCSVOutput).Run()
}

// read the csv file and send lines to a channel (pipeline)
func readCSVInput(row chan Record) {
	if inputPath == nil {
		log.Fatal("input path is empty. Please specify the input path.")
	}

	csvfile, err := os.Open(*inputPath)
	if err != nil {
		log.Fatal(err)
	}

	defer csvfile.Close()

	r := csv.NewReader(csvfile)

	if _, err := r.Read(); err != nil { //read header
		log.Fatal(err)
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// parse the 'created_utc' from string to integer
		i, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		// create an object to store the column
		// and send the record to the pipeline
		row <- Record{
			Subreddit:    record[1],
			User:         record[2],
			Time:         time.Unix(i, 0).UTC(),
			NumberOfPost: 1,
		}
	}
}

func writeCSVOutput(key string, rs []Record) {
	// sort the number of submission (DESC order) locally
	sort.SliceStable(rs, func(i, j int) bool {
		return rs[i].NumberOfPost > rs[j].NumberOfPost
	})

	// create a directory from the key format: {subreddit}:::{date}
	// to a directory format: {subreddit}/{date}.csv
	arr := strings.Split(key, Delimiter)
	if len(arr) != 2 {
		log.Fatal("Invalid key :" + key + " found. Cannot split it.")
	}

	subreddit, date := arr[0], arr[1]

	dir := filepath.Join(ResultDir, subreddit)

	// create all the parent directories by given path
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatal(err)
	}

	fileName := date + ".csv"

	file, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	// write the buffer data into file
	// when the function is finished
	defer writer.Flush()

	// write the csv header
	err = writer.Write([]string{
		"subreddit",
		"date",
		"user",
		"num_of_submission",
	})
	if err != nil {
		log.Fatal(err)
	}

	// write the data into the IO buffer
	for _, r := range rs {
		writer.Write([]string{
			r.Subreddit,
			r.Time.Format("2006-01-02"),
			r.User,
			strconv.Itoa(r.NumberOfPost),
		})
	}
}
