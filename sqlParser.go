package main

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"reflect"
)

type visitor struct {
	tableList []string
}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	//fmt.Printf("%T\n", in)
	if reflect.TypeOf(in).String() == "*ast.TableName" {
		tableName := in.(*ast.TableName).Name.L
		if !ElementExistsInSlice(tableName, v.tableList) {
			v.tableList = append(v.tableList, tableName)
		}
	}

	return in, false
}

func (v *visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type SqlParser struct{}

func (s *SqlParser) GetTables(sqlText string) (tables []string, warns []error, err error) {
	if stmtNodes, warns, err := parser.New().Parse(sqlText, "", ""); warns == nil && err == nil {
		for _, stmtNode := range stmtNodes {
			v := visitor{}
			stmtNode.Accept(&v)
			tables = v.tableList
		}
	}

	return tables, warns, err
}
