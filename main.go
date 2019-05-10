package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {

	// PARAMS
	fileName := flag.String("file", "", "file path to be imported")
	processScript := flag.String("process", "", "process script file path")
	flag.Parse()

	fmt.Printf("File: %s\n", *fileName)
	fmt.Printf("Script: %s\n", *processScript)

	// LOAD FILE
	csvFile, error := os.Open(*fileName)
	if error != nil {
		fmt.Printf("ERROR!")
		os.Exit(1)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Read()

	/*
		header, error := reader.Read()
		if error == io.EOF {
		} else if error != nil {
			log.Fatal(error)
		}

		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			// TODO: CALL/EXECUTE THE SCRIPT WITH THE ROW DATA

			// TODO: PROCESS THE RETURNED DATA AND SAVE IT INTO REDIS
		}
	*/
}
