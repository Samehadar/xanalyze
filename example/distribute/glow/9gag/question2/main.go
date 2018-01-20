package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/sniperkit/xanalyze/plugin/distribute/glow/flow"
)

type Record struct {
	Subreddit string
	User      string
	Time      time.Time
}

var (
	inputPath *string
)

const (
	OutputPath = "result.output"
)

func main() {

	// get the input path for command line argument
	inputPath = flag.String("input", "", "Absolute input file (e.g. /Users/name/input.csv)")

	flag.Parse()

	// make '-input' become required (validation)
	if *inputPath == "" {
		log.Fatal("The flag '-input' is missing. Please specific CSV file location. (e.g. -input /Users/name/input.csv) ")
	}

	// clean output directory to make sure there is no old data
	if _, err := os.Stat(OutputPath); err == nil {
		os.Remove(OutputPath)
	}

	// timer for benchmark
	startTime := time.Now()

	log.Println("Started please wait... It may take some time...")

	// run the program
	Run(OutputPath)

	log.Println("Done! Total time in seconds: ", time.Now().Sub(startTime).Seconds())
	log.Println("The result is in the same directory with a file name: " + OutputPath)
}

func Run(outputPath string) {
	// build the map-reduce pipeline
	flow.New().Source(readCSVInput, runtime.NumCPU()).Map(func(r Record) (string, time.Time) {
		// use delimiter to group by subreddit, date, and user
		return r.User, r.Time
	}).GroupByKey().Map(func(user string, ts []time.Time) (string, int) {
		// sort the times in a group
		sort.SliceStable(ts, func(i, j int) bool {
			return ts[i].Before(ts[j])
		})

		maxDay := 0
		currentDay := 0

		for i, _ := range ts {
			if i+1 >= len(ts) {
				break
			}

			xYear, xMonth, xDay := ts[i].UTC().Date()
			yYear, yMonth, yDay := ts[i+1].UTC().Date()

			// make sure the year and the month are the same
			// if it is not match, reset the counter
			if xYear == yYear && xMonth == yMonth {
				if yDay-xDay == 1 {
					// the next post is on next day
					// so this is one consecutive day
					currentDay++
				} else if yDay-xDay == 0 {
					// there are two post within the same day
					// so do nothing
				} else {
					// the posts are not continue
					// so reset the day counter
					currentDay = 0
				}
			} else {
				currentDay = 0
			}

			// get the max. number of consecutive days
			if currentDay > maxDay {
				maxDay = currentDay
			}
		}

		return user, maxDay
	}).Filter(func(user string, maxDay int) bool {
		return maxDay > 0
	}).Map(func(user string, maxDay int) (int, string) {
		// change the position so that
		// the library can sort the counter
		return maxDay, user
	}).Sort(func(a, b int) bool {
		return a > b
	}).Map(func(maxDay int, user string) string {
		// prepare output format
		return fmt.Sprint(user, ",", maxDay)
	}).SaveTextToFile(outputPath)
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
			Subreddit: record[1],
			User:      record[2],
			Time:      time.Unix(i, 0).UTC(),
		}
	}
}
