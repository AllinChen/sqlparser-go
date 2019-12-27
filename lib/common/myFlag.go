package common

import (
	"flag"
	"fmt"
	"os"
)

type MyFlag struct {
	Help    bool
	Version bool
	Sql     string
}

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
	flag.BoolVar(&f.Help, "help", false, "show help and exit")
	flag.BoolVar(&f.Version, "version", false, "show version and exit")
	flag.StringVar(&f.Sql, "sql", "", "input sql text")
	flag.Parse()

	flag.Usage = usage

	if f.Help {
		flag.Usage()
		os.Exit(0)
	}

	if f.Version {
		fmt.Println("sqlparser version: v1.0.0")
		os.Exit(0)
	}
}
