package config

import (
	"flag"
	"os"
)

type Config struct {
	DBConfig *DBConfig

	CsvFilePath     string
	ProcessFilePath string
	FullImport      bool
	Debug           bool
}

type DBConfig struct {
	Dialect  string
	Name     string
	Username string
	Password string
	Host     string
	Port     string
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); len(value) != 0 {
		return value
	}
	return fallback
}

func GetConf() *Config {
	config := Config{
		CsvFilePath:     getenv("APPAKA_CSV_FILE_PATH", ""),
		ProcessFilePath: getenv("APPAKA_PROCESS_FILE_PATH", ""),
		FullImport:      getenv("APPAKA_FULL_IMPOTY", "0"),
		Debug:           getenv("APPAKA_DEBUG", "1"),
	}

	// PARAMS
	csvFilePath := flag.String("file", "", "file path to be imported")
	processFilePath := flag.String("process", "", "process script file path")
	fullImport := flag.Bool("full-import", false, "delete items not presents in this file")
	debug := flag.Bool("debug", true, "debug process")
	flag.Parse()

	config.CsvFilePath = *csvFilePath
	config.ProcessFilePath = *processFilePath
	config.FullImport = *fullImport
	config.Debug = *debug

	config.DBConfig = &DBConfig{
		Dialect:  getenv("TODO_DB_DIALECT", "postgres"),
		Name:     getenv("TODO_DB_NAME", "tododb"),
		Username: getenv("TODO_DB_USERNAME", "tododb"),
		Password: getenv("TODO_DB_PASSWORD", "tododb"),
		Host:     getenv("TODO_DB_HOST", "localhost"),
		Port:     getenv("TODO_DB_PORT", "5432"),
	}

	return &config
}
