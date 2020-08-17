package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rolandvarga/batch_proc/job"
)

func parseInput(path string) (job.Data, error) {
	var data job.Data

	// read file into byte array
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return data, fmt.Errorf("unable to read json file: %s", err.Error())
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return data, fmt.Errorf("unable to unmarshal json: %s", err.Error())
	}
	return data, nil
}

func main() {
	fmt.Println("Usage: ./batch_proc <FILE>")

	flag.Parse()

	path := flag.Args()[0]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("input file '%s' cannot be found. Exiting...\n", path)
		os.Exit(1)
	}

	data, err := parseInput(path)
	if err != nil {
		panic(err)
	}

	batchJob := job.NewJob(data)
	batchJob.Run()
}
