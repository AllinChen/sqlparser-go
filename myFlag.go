package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	help    bool
	version bool
	sql     string
)

type MyFlag struct{}

func usage() {
	fmt.Fprintf(
		os.Stderr, "sqlparser can parse input sql and print json formatted dbNames/tableNames/tableCommments/columnNames/columnComments.\n"+
			"Please notice that ONLY these statements will be parsed:\n"+
			"\tcreate table\n"+
			"\talter table\n"+
			"\tdrop table\n"+
			"\tselect\n"+
			"\tinsert\n"+
			"\treplace\n"+
			"\tupdate\n"+
			"\tdelete\n"+
			"version: v1.0.0\n"+
			"Usage: sqlparser [--help] [--version] [--sql string]\n"+
			"Options:\n")
	flag.PrintDefaults()
}

func (f *MyFlag) Init() {
	flag.BoolVar(&help, "help", false, "show help and exit")
	flag.BoolVar(&version, "version", false, "show version and exit")
	flag.StringVar(&sql, "sql", "", "input sql text")
	flag.Parse()

	flag.Usage = usage

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if version {
		fmt.Println("sqlparser version: v1.0.0")
		os.Exit(0)
	}
}
