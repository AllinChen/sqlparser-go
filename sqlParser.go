package main

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/tidb/types/parser_driver"
	"reflect"
	"strings"
)

type visitor struct {
	sqlType          string
	toParse          bool
	visitSqlList     []string
	funcList         []string
	dbList           []string
	tableList        []string
	tableCommentMap  map[string]string
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
	v.sqlType = ""
	v.dbList = []string{}
	v.tableList = []string{}
	v.tableCommentMap = make(map[string]string)
	v.columnList = []string{}
	v.columnCommentMap = make(map[string]string)
}

func (v *visitor) AddDb(dbName string) {
	if !StringInSlice(dbName, v.dbList) {
		v.dbList = append(v.dbList, dbName)
	}
}

func (v *visitor) AddTable(tableName string) {
	if !StringInSlice(tableName, v.tableList) {
		v.tableList = append(v.tableList, tableName)
	}
}

func (v *visitor) AddTableComment(tableName string, tableComment string) {
	v.tableCommentMap[tableName] = tableComment
}

func (v *visitor) AddColumn(columnName string) {
	if !StringInSlice(columnName, v.columnList) {
		v.columnList = append(v.columnList, columnName)
	}
}

func (v *visitor) AddColumnComment(columnName string, columnComment string) {
	v.columnCommentMap[columnName] = columnComment
}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	var funcArgs []ast.ExprNode

	dbName := ""
	tableName := ""
	columnName := ""
	exprType := ""
	tableComment := ""
	columnComment := ""
	astType := reflect.TypeOf(in).String()

	if StringInSlice(astType, v.visitSqlList) {
		v.toParse = true
		// 获取语句类型
		v.sqlType = strings.Split(astType, ".")[1]
	}

	if v.toParse {
		//fmt.Println(astType)

		// 获取表名称
		if astType == "*ast.TableName" {
			tableName = in.(*ast.TableName).Name.L
			v.AddTable(tableName)
			dbName = in.(*ast.TableName).Schema.L
			if dbName != "" {
				v.AddDb(dbName)
			}
		}

		// 获取表注释
		if astType == "*ast.CreateTableStmt" {
			tableName = in.(*ast.CreateTableStmt).Table.Name.L
			for _, tableOption := range in.(*ast.CreateTableStmt).Options {
				if tableOption.Tp == ast.TableOptionComment {
					tableComment = tableOption.StrValue
				}
			}

			v.AddTableComment(tableName, tableComment)
		} else if astType == "*ast.AlterTableStmt" {
			tableName = in.(*ast.AlterTableStmt).Table.Name.L
			for _, tableSpec := range in.(*ast.AlterTableStmt).Specs {
				for _, tableOption := range tableSpec.Options {
					if tableOption.Tp == ast.TableOptionComment {
						tableComment = tableOption.StrValue
					}
				}
			}

			v.AddTableComment(tableName, tableComment)
		}

		// 获取字段名称
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

		// 获取字段注释
		if astType == "*ast.ColumnDef" {
			columnName := in.(*ast.ColumnDef).Name.Name.L
			v.AddColumn(columnName)

			for _, columnOption := range in.(*ast.ColumnDef).Options {
				if columnOption.Tp == ast.ColumnOptionComment {
					columnComment = columnOption.Expr.(*driver.ValueExpr).Datum.GetString()
				}
			}

			v.AddColumnComment(columnName, columnComment)
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
