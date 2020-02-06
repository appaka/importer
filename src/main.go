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

const VERSION = "0.0.1"

type App struct {
	// config values
	csvFilePath     string
	processFilePath string
	deleteOthers    bool
	debug           bool

	// files and documents
	reader  *csv.Reader
	headers []string
	EOF     bool
	Items   int64
	Skipped int64
}

func (app *App) Init() {
	app.ShowVersion()
	app.loadConfig()

}

func (app *App) loadConfig() {
	// PARAMS
	csvFilePath := flag.String("file", "", "file path to be imported")
	processFilePath := flag.String("process", "", "process script file path")
	deleteOthers := flag.Bool("delete-others", false, "delete items not presents in this file")
	debug := flag.Bool("debug", true, "debug process")
	flag.Parse()

	app.csvFilePath = *csvFilePath
	app.processFilePath = *processFilePath
	app.deleteOthers = *deleteOthers
	app.debug = *debug

	mapBooleanYesNo := map[bool]string{
		false: "No",
		true:  "Yes",
	}

	h1("CONFIGURATION")
	fmt.Printf(" --file (CSV file path)                 = %s\n", app.csvFilePath)
	fmt.Printf(" --process (Processor script file path) = %s\n", app.processFilePath)
	fmt.Printf(" --delete-others (Delete others items)  = %s\n", mapBooleanYesNo[app.deleteOthers])
	fmt.Printf(" --debug (Debug process)                = %s\n", mapBooleanYesNo[app.debug])
	fmt.Println()
}

func (app *App) ShowVersion() {
	fmt.Printf("AppakaDB importer v%s, by Javier Perez <hallo@javierperez.ch>\n\n", VERSION)
}

func (app *App) OpenFile() {
	stream, error := os.Open(app.csvFilePath)
	if error != nil {
		fmt.Printf("ERROR opening file %s!\n", app.csvFilePath)
		// Unable to open file error
		os.Exit(1)
	}

	app.reader = csv.NewReader(bufio.NewReader(stream))

	headers, err := app.reader.Read()
	if err == io.EOF {
		// No headers error
		os.Exit(1)
	}

	app.headers = headers
	app.EOF = false
	app.Items = 0
	app.Skipped = 0
}

func (app *App) GetNextDocument() *string {
	record, err := app.reader.Read()
	if err == io.EOF {
		app.EOF = true
		return nil
	}
	if err != nil {
		app.Skipped++
		return nil
	}

	// transform the array into a map
	row := map[string]string{}
	for i := 0; i < len(app.headers); i++ {
		var field string = app.headers[i]
		var value string = record[i]
		row[field] = value
	}

	// transform the map into a JSON object
	document, err := json.Marshal(row)
	if err != nil {
		app.Skipped++
		return nil
	}

	app.Items++

	documentString := string(document)
	return &documentString
}

func h1(title string) {
	fmt.Printf("%s %s\n", title, strings.Repeat("=", 60-len(title)-1))
}

func main() {

	app := App{}
	app.Init()
	app.OpenFile()

	start := time.Now()

	for {
		// transform the array into a map
		document := app.GetNextDocument()
		if app.EOF {
			break
		}
		if document == nil {
			// TODO: error reading document => implement proper logging
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
	fmt.Printf("    processed: %14d (%d skipped / %.2f%%)\n", app.Items, app.Skipped, float64(app.Skipped/app.Items)*100)
	fmt.Printf(" time elapsed: %14.4f\n", elapsedSeconds)
	fmt.Printf("  rows/second: %14.4f\n", float64(app.Items)/elapsedSeconds)
	fmt.Printf(" files/second: %14.4f\n", (float64(app.Items)/elapsedSeconds)/float64(app.Items))

	os.Exit(0)
}
