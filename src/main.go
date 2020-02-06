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
	stream  *os.File
	EOF     bool
	Items   int64
	Skipped int64

	// timing
	start time.Time
}

func (app *App) Init() {
	app.showVersion()
	app.loadConfig()
	app.openFile()

	app.start = time.Now()

	for {
		// transform the array into a map
		document := app.getNextDocument()
		if app.EOF {
			break
		}
		if document == nil {
			// TODO: error reading document => implement proper logging
			continue
		}

		app.onDocument(document)
	}

	/*
		TODO: send list with all IDs to MDM, if the param "--delete-old" is present
	*/

	app.summary()
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

func (app *App) showVersion() {
	fmt.Printf("AppakaDB importer v%s, by Javier Perez <hallo@javierperez.ch>\n\n", VERSION)
}

func (app *App) openFile() {
	var err error
	app.stream, err = os.Open(app.csvFilePath)
	if err != nil {
		fmt.Printf("ERROR opening file %s!\n", app.csvFilePath)
		// Unable to open file error
		os.Exit(1)
	}

	app.reader = csv.NewReader(bufio.NewReader(app.stream))

	app.headers, err = app.reader.Read()
	if err == io.EOF {
		// No headers error
		os.Exit(1)
	}

	app.EOF = false
	app.Items = 0
	app.Skipped = 0
}

func (app *App) summary() {
	_ = app.stream.Close()

	// final stats
	elapsedSeconds := time.Now().Sub(app.start).Seconds()
	h1("STATS")
	fmt.Printf("%14.4f time elapsed (seconds)\n", elapsedSeconds)
	fmt.Printf("%14.4f documents processed (%d skipped / %.2f%%)\n", float64(app.Items), app.Skipped, float64(app.Skipped/app.Items)*100)
	fmt.Printf("%14.4f documents/second\n", float64(app.Items)/elapsedSeconds)
	fmt.Printf("%14.4f files/second\n", (float64(app.Items)/elapsedSeconds)/float64(app.Items))

	os.Exit(0)
}

func (app *App) getNextDocument() *string {
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

func (app *App) onDocument(document *string) {
	/*
		TODO: execute javascript reader-processor, which will return the final document
		- exec javascript file into our JS VM, and it will return ID and final document (nil if we should skip it)
		- this js file could access to MDM (MongoDB?) to fetch data

		if there is no js file, then calculate the ID with the given params/config
	*/

	/*
		TODO: send document to MDM
		send document to a RabbitMQ queue, which will be consumed by the MDM
	*/
}

func h1(title string) {
	fmt.Printf("%s %s\n", title, strings.Repeat("=", 60-len(title)-1))
}

func main() {
	app := App{}
	app.Init()
}
