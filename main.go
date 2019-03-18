package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"os"
)

func main() {
	var (
		sqlMap  map[string]interface{}
		sqlText string
	)

	f := MyFlag{}
	f.Init()

	parser := SqlParser{}

	if format == "text" {
		sqlText = sql
		if tables, warns, err := parser.GetTables(sqlText); warns == nil && err == nil {
			fmt.Println(tables)
		} else {
			fmt.Printf("parse warn:\n%v\nparse error:\n%v\n", warns, err)
			os.Exit(1)
		}
	} else if format == "json" {
		if err := json.Unmarshal([]byte(sql), &sqlMap); err == nil {
			sqlText = sqlMap["sqlText"].(string)
			if tables, warns, err := parser.GetTables(sqlText); warns == nil && err == nil {
				tableMap := make(map[string]interface{})
				tableMap["tables"] = tables

				if tableMapBytes, err := json.Marshal(tableMap); err == nil {
					fmt.Println(string(tableMapBytes))
				} else {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				fmt.Printf("parse warn:\n%v\nparse error:\n%v\n", warns, err)
				os.Exit(1)
			}
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		usage()
		os.Exit(1)
	}
}
