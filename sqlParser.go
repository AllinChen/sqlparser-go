package main

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/tidb/types/parser_driver"
	"reflect"
)

type visitor struct {
	toParse          bool
	visitSqlList     []string
	funcList         []string
	tableList        []string
	columnList       []string
	columnCommentMap map[string]string
}

func (v *visitor) Init() {
	v.toParse = false
	v.visitSqlList = []string{
		"*ast.CreateTableStmt",
		"*ast.AlterTableStmt",
		"*ast.DropTableStmt",
		"*ast.SelectStmt",
		"*ast.UnionStmt",
		"*ast.InsertStmt",
		"*ast.ReplaceStmt",
		"*ast.InsertStmt",
		"*ast.UpdateStmt",
		"*ast.DeleteStmt",
	}
	v.funcList = []string{
		"*ast.FuncCallExpr",
		"*ast.AggregateFuncExpr",
		"*ast.WindowFuncExpr",
	}
	v.tableList = []string{}
	v.columnList = []string{}
	v.columnCommentMap = make(map[string]string)
}

func (v *visitor) AddTable(tableName string) {
	if !StringInSlice(tableName, v.tableList) {
		v.tableList = append(v.tableList, tableName)
	}
}

func (v *visitor) AddColumn(columnName string) {
	if !StringInSlice(columnName, v.columnList) {
		v.columnList = append(v.columnList, columnName)
	}
}

func (v *visitor) AddComment(columnName string, columnComment string) {
	v.columnCommentMap[columnName] = columnComment
}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	var columnName string
	exprType := ""
	var funcArgs []ast.ExprNode

	astType := reflect.TypeOf(in).String()

	if StringInSlice(astType, v.visitSqlList) {
		v.toParse = true
	}

	if v.toParse {
		if astType == "*ast.ColumnDef" {
			columnName := in.(*ast.ColumnDef).Name.Name.L
			v.columnList = append(v.columnList, columnName)
			v.columnCommentMap[columnName] = ""

			for _, columnOption := range in.(*ast.ColumnDef).Options {
				if columnOption.Tp == ast.ColumnOptionComment {
					columnComment := columnOption.Expr.(*driver.ValueExpr).Datum.GetString()
					v.AddComment(columnName, columnComment)
				}
			}
		}

		if astType == "*ast.TableName" {
			tableName := in.(*ast.TableName).Name.L
			v.AddTable(tableName)
		}

		if astType == "*ast.SelectField" {
			expr := in.(*ast.SelectField).Expr
			if expr == nil && in.(*ast.SelectField).WildCard != nil {
				columnName = "*"
				v.AddColumn(columnName)
			} else if expr != nil {
				exprType = reflect.TypeOf(expr).String()

				if StringInSlice(exprType, v.funcList) {
					funcArgs = []ast.ExprNode{}
					if exprType == "*ast.AggregateFuncExpr" {
						funcArgs = expr.(*ast.AggregateFuncExpr).Args
					} else if exprType == "*ast.FuncCallExpr" {
						funcArgs = expr.(*ast.FuncCallExpr).Args
					} else if exprType == "*ast.WindowFuncExpr" {
						funcArgs = expr.(*ast.WindowFuncExpr).Args
					}

					for _, arg := range funcArgs {
						if reflect.TypeOf(arg).String() == "*ast.ColumnNameExpr" {
							columnName = arg.(*ast.ColumnNameExpr).Name.Name.L
							v.AddColumn(columnName)
						}
					}
				} else if exprType == "*ast.ColumnNameExpr" {
					columnName = expr.(*ast.ColumnNameExpr).Name.Name.L
					v.AddColumn(columnName)
				}
			}
		}
	}

	return in, false
}

func (v *visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type SqlParser struct{}

//parse sql and return *visitor
func (s *SqlParser) ParseSql(sqlText string) (vis *visitor, warns []error, err error) {
	v := visitor{}
	v.Init()

	if stmtNodes, w, e := parser.New().Parse(sqlText, "", ""); w == nil && e == nil {
		for _, stmtNode := range stmtNodes {
			stmtNode.Accept(&v)
			vis = &v
		}
	} else {
		warns = w
		err = e
	}

	return vis, warns, err
}
