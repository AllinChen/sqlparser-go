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
	format  string
)

type MyFlag struct{}

func usage() {
	fmt.Fprintf(os.Stderr, "sqlparser can parse input sql and print related tables.\n"+
		"version: v1.0.0\n"+
		"Usage: sqlparser [--help] [--version] [-sql string] [--format string]\n"+
		"Options:\n")
	flag.PrintDefaults()
}

func (f *MyFlag) Init() {
	flag.BoolVar(&help, "help", false, "show version and exit")
	flag.BoolVar(&version, "version", false, "show version and exit")
	flag.StringVar(&format, "format", "text", ""+
		"sql text format: text|json,\n"+
		"\tif use text format, you should input sql as text format,\n"+
		"\t\te.g.: select id from t01;\n"+
		"\tand the output will be a text,\n"+
		"\t\te.g.: [t01]\n\n"+
		"\tif use json format, you should input sql as json format,\n"+
		"\t\te.g.: {\"sql\":\"select id from t01;\"}\n"+
		"\tand the output will be a json string too,\n"+
		"\t\te.g.: {\"tables\":\"[\"t01\"]\"}\n")
	flag.StringVar(&sql, "sql", "", "input sql text")
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}

	if version {
		fmt.Println("sqlparser version: v1.0.0")
		os.Exit(0)
	}
}
