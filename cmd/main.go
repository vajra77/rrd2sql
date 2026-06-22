package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	InputPath  string
	OutputPath string
}

func usage() {
	flag.PrintDefaults()
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.InputPath, "input", "", "input path")
	flag.StringVar(&cfg.OutputPath, "output", "", "output path")
	flag.Parse()

	if cfg.InputPath == "" || cfg.OutputPath == "" {
		fmt.Println("ERROR: input and output path must be specified")
		usage()
		os.Exit(1)
	}

	res, err := FetchDailyValues(cfg.InputPath)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(1)
	}

	fmt.Printf("AvgIn: %.2fGbps AvgOut: %.2fGbps\nMaxIn: %.2fGbps MaxOut: %.2fGbps\n", res.AvgIn/1000000000, res.AvgOut/1000000000, res.MaxIn/1000000000, res.MaxOut/1000000000)

	db, err := NewSql3Db(cfg.OutputPath)
	if err != nil {
		fmt.Printf("ERROR: unable to open DB %s, %v", cfg.OutputPath, err)
		os.Exit(1)
	}

	if err := db.WriteResult(res); err != nil {
		fmt.Printf("ERROR: unable to write to DB, %v", err)
	}
	if err := db.Close(); err != nil {
		fmt.Printf("ERROR: unable to close DB: %v", err)
	}
}
