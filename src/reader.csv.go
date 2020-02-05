package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type CSVFile struct {
	reader  *csv.Reader
	headers []string
	EOF     bool
	Items   int64
	Skipped int64
}

func (file *CSVFile) Open(filePath string) {
	stream, error := os.Open(filePath)
	if error != nil {
		fmt.Printf("ERROR opening file %s!\n", filePath)
		// Unable to open file error
		os.Exit(1)
	}

	file.reader = csv.NewReader(bufio.NewReader(stream))

	headers, err := file.reader.Read()
	if err == io.EOF {
		// No headers error
		os.Exit(1)
	}

	file.headers = headers
	file.EOF = false
	file.Items = 0
	file.Skipped = 0
}

func (file *CSVFile) GetDocument() *string {
	record, err := file.reader.Read()
	if err == io.EOF {
		file.EOF = true
		return nil
	}
	if err != nil {
		file.Skipped++
		return nil
	}

	// transform the array into a map
	row := map[string]string{}
	for i := 0; i < len(file.headers); i++ {
		var field string = file.headers[i]
		var value string = record[i]
		row[field] = value
	}

	// transform the map into a JSON object
	document, err := json.Marshal(row)
	if err != nil {
		file.Skipped++
		return nil
	}

	file.Items++

	documentString := string(document)
	return &documentString
}

type Config struct {
	CsvFilePath     string
	ProcessFilePath string
	DeleteOthers    bool
	Debug           bool
}

func (config *Config) Load() {
	// PARAMS
	csvFilePath := flag.String("file", "", "file path to be imported")
	processFilePath := flag.String("process", "", "process script file path")
	deleteOthers := flag.Bool("delete-others", false, "delete items not presents in this file")
	debug := flag.Bool("debug", true, "debug process")
	flag.Parse()

	config.CsvFilePath = *csvFilePath
	config.ProcessFilePath = *processFilePath
	config.DeleteOthers = *deleteOthers
	config.Debug = *debug

	if config.Debug {
		mapBooleanYesNo := map[bool]string{
			false: "No",
			true:  "Yes",
		}

		h1("CONFIGURATION")
		fmt.Printf(" --file (CSV file path)                 = %s\n", config.CsvFilePath)
		fmt.Printf(" --process (Processor script file path) = %s\n", config.ProcessFilePath)
		fmt.Printf(" --delete-others (Delete others items)  = %s\n", mapBooleanYesNo[config.DeleteOthers])
		fmt.Printf(" --debug (Debug process)                = %s\n", mapBooleanYesNo[config.Debug])
		fmt.Println()
	}
}

func h1(title string) {
	fmt.Printf("%s %s\n", title, strings.Repeat("=", 60-len(title)-1))
}

func main() {
	config := Config{}
	config.Load()

	file := CSVFile{}
	file.Open(config.CsvFilePath)

	start := time.Now()

	for {
		// transform the array into a map
		document := file.GetDocument()
		if file.EOF {
			break
		}
		if document == nil {
			continue
		}

		/*
			TODO: execute reader-processor, which will return the final document
			- exec javascript file into our JS VM, and it will return ID and final document (nil if we should skip it)
			- this js file could access to MDM (MongoDB?) to fetch data

			if there is no js file, then calculate the ID with the given params/config
		*/

		/*
			TODO: send document to MDM
			headers:
				Auth-User: manolo@node1
			POST /document/{entity}/{id}
			{
				"id": "123456",
				"name: "manolo",
			}

			--- OR ---

			send document to a RabbitMQ queue, which will be consumed by the MDM
		*/

	}

	/*
		TODO: send list with all IDs to MDM, if the param "--delete-old" is present
	*/

	if config.Debug {
		// final stats
		elapsedSeconds := time.Now().Sub(start).Seconds()
		h1("STATS")
		fmt.Printf("    processed: %14d (%d skipped / %.2f%%)\n", file.Items, file.Skipped, float64(file.Skipped/file.Items)*100)
		fmt.Printf(" time elapsed: %14.4f\n", elapsedSeconds)
		fmt.Printf("  rows/second: %14.4f\n", float64(file.Items)/elapsedSeconds)
		fmt.Printf(" files/second: %14.4f\n", (float64(file.Items)/elapsedSeconds)/float64(file.Items))
	}

	os.Exit(0)
}
