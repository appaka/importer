package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type CsvReader struct {
	File *os.File
}

func getRow(r *csv.Reader, headers []string) *string {
	record, err := r.Read()
	if err != nil {
		return nil
	}

	row := map[string]string{}
	if headers != nil {
		// transform the array into a map
		for i := 0; i < len(headers); i++ {
			var field string = headers[i]
			var value = record[i]
			row[field] = value
		}
	}

	// transform the map into a JSON object
	document, err := json.Marshal(row)
	if err != nil {
		return nil
	}
	documentString := string(document)

	return &documentString
}

func openFile(FileName string) (Reader *csv.Reader, Headers []string) {
	FileReader, error := os.Open(FileName)
	if error != nil {
		fmt.Printf("ERROR!")
		os.Exit(1)
	}

	Reader = csv.NewReader(bufio.NewReader(FileReader))

	Headers, err := Reader.Read()
	if err == io.EOF {
		os.Exit(1)
	}

	return Reader, Headers
}

func main() {
	// PARAMS
	fileName := flag.String("file", "", "file path to be imported")
	processScript := flag.String("process", "", "process script file path")
	flag.Parse()

	fmt.Printf("File: %s\n", *fileName)
	fmt.Printf("Script: %s\n", *processScript)

	// LOAD FILE
	Reader, Headers := openFile(*fileName)

	start := time.Now()
	counter := float64(0)

	for {
		counter++

		// transform the array into a map
		document := getRow(Reader, Headers)
		if document == nil {
			break
		}

		// print it
		//fmt.Printf("%.0f ", counter)
		//fmt.Println(*document)

		/*
			TODO: execute reader-processor, which will return the final document
			- exec javascript file into our JS VM, and it will return ID and final document (nil if we should skip it)
			- this js file could access to MDM (MongoDB?) to fetch data

			if there is no js file, then calculate the ID with the given params/config
		*/

		/*
			TODO: send document to MDM
			Headers:
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
	t := time.Now()
	elapsedSeconds := t.Sub(start).Seconds()
	rowsPerSecond := counter / elapsedSeconds

	fmt.Printf("time elapsed: %.4f\n", elapsedSeconds)
	fmt.Printf("rows per second: %.4f\n", rowsPerSecond)
}
