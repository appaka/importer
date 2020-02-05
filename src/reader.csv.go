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

func getRow(reader *csv.Reader, headers []string) (*string, error) {
	record, err := reader.Read()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	documentString := string(document)

	return &documentString, nil
}

func openFile(filename string) (reader *csv.Reader, headers []string) {
	stream, error := os.Open(filename)
	if error != nil {
		fmt.Printf("ERROR!")
		os.Exit(1)
	}

	reader = csv.NewReader(bufio.NewReader(stream))

	headers, err := reader.Read()
	if err == io.EOF {
		os.Exit(1)
	}

	return reader, headers
}

func main() {
	// PARAMS
	filename := flag.String("file", "", "file path to be imported")
	preprocessorScript := flag.String("process", "", "process script file path")
	flag.Parse()

	fmt.Printf("File: %s\n", *filename)
	fmt.Printf("Script: %s\n", *preprocessorScript)

	// LOAD FILE
	reader, headers := openFile(*filename)

	start := time.Now()
	counter := float64(0)
	skipped := int16(0)

	for {
		counter++

		// transform the array into a map
		document, err := getRow(reader, headers)
		if err == io.EOF {
			break
		}
		if document == nil {
			skipped++
			continue
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
	fmt.Printf("   processed: %14.0f (%d skipped / %.2f%%)\n", counter, skipped, (float64(skipped)/counter)*100)
	fmt.Printf("time elapsed: %14.4f\n", elapsedSeconds)
	fmt.Printf(" rows/second: %14.4f\n", counter/elapsedSeconds)
	fmt.Printf("files/second: %14.4f\n", (counter/elapsedSeconds)/counter)
}
