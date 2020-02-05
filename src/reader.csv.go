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

var config struct {
	CsvFilePath     string
	ProcessFilePath string
	DeleteOthers    bool
}

var csvFile struct {
	Reader  *csv.Reader
	Headers []string
	EOF     bool
	Items   int64
	Skipped int64
}

func getDocument() *string {
	record, err := csvFile.Reader.Read()
	if err == io.EOF {
		csvFile.EOF = true
		return nil
	}
	if err != nil {
		csvFile.Skipped++
		return nil
	}

	// transform the array into a map
	row := map[string]string{}
	for i := 0; i < len(csvFile.Headers); i++ {
		var field string = csvFile.Headers[i]
		var value string = record[i]
		row[field] = value
	}

	// transform the map into a JSON object
	document, err := json.Marshal(row)
	if err != nil {
		csvFile.Skipped++
		return nil
	}

	csvFile.Items++

	documentString := string(document)
	return &documentString
}

func openCSVFile() {
	stream, error := os.Open(config.CsvFilePath)
	if error != nil {
		fmt.Printf("ERROR!")
		os.Exit(1)
	}

	csvFile.Reader = csv.NewReader(bufio.NewReader(stream))

	headers, err := csvFile.Reader.Read()
	if err == io.EOF {
		os.Exit(1)
	}

	csvFile.Headers = headers
	csvFile.EOF = false
	csvFile.Items = 0
	csvFile.Skipped = 0
}

func h1(title string) {
	fmt.Printf("%s %s\n", title, strings.Repeat("=", 60-len(title)-1))
}

func loadConfig() {
	// PARAMS
	csvFilePath := flag.String("file", "", "file path to be imported")
	processFilePath := flag.String("process", "", "process script file path")
	deleteOthers := flag.Bool("delete-others", false, "delete items not presents in this file")
	flag.Parse()

	config.CsvFilePath = *csvFilePath
	config.ProcessFilePath = *processFilePath
	config.DeleteOthers = *deleteOthers

	mapBooleanYesNo := map[bool]string{
		false: "No",
		true:  "Yes",
	}

	h1("CONFIGURATION")
	fmt.Printf(" --file (CSV file path)                 = %s\n", config.CsvFilePath)
	fmt.Printf(" --process (Processor script file path) = %s\n", config.ProcessFilePath)
	fmt.Printf(" --delete-others (Delete others items)  = %s\n", mapBooleanYesNo[config.DeleteOthers])
	fmt.Println()
}

func main() {
	loadConfig()
	openCSVFile()

	start := time.Now()

	for {
		// transform the array into a map
		document := getDocument()
		if csvFile.EOF {
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

	// final stats
	elapsedSeconds := time.Now().Sub(start).Seconds()
	h1("STATS")
	fmt.Printf("    processed: %14d (%d skipped / %.2f%%)\n", csvFile.Items, csvFile.Skipped, float64(csvFile.Skipped/csvFile.Items)*100)
	fmt.Printf(" time elapsed: %14.4f\n", elapsedSeconds)
	fmt.Printf("  rows/second: %14.4f\n", float64(csvFile.Items)/elapsedSeconds)
	fmt.Printf(" files/second: %14.4f\n", (float64(csvFile.Items)/elapsedSeconds)/float64(csvFile.Items))
}
