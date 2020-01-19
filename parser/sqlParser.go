package parser

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/tidb/types/parser_driver"
	"reflect"
	"sqlparser-go/lib/common"
	"strings"
)

type visitor struct {
	SqlType          string
	toParse          bool
	visitSqlList     []string
	funcList         []string
	DbList           []string
	TableList        []string
	TableCommentMap  map[string]string
	ColumnList       []string
	ColumnCommentMap map[string]string
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
	v.SqlType = ""
	v.DbList = []string{}
	v.TableList = []string{}
	v.TableCommentMap = make(map[string]string)
	v.ColumnList = []string{}
	v.ColumnCommentMap = make(map[string]string)
}

func (v *visitor) AddDb(dbName string) {
	if !common.StringInSlice(dbName, v.DbList) {
		v.DbList = append(v.DbList, dbName)
	}
}

func (v *visitor) AddTable(tableName string) {
	if !common.StringInSlice(tableName, v.TableList) {
		v.TableList = append(v.TableList, tableName)
	}
}

func (v *visitor) AddTableComment(tableName string, tableComment string) {
	v.TableCommentMap[tableName] = tableComment
}

func (v *visitor) AddColumn(columnName string) {
	if !common.StringInSlice(columnName, v.ColumnList) {
		v.ColumnList = append(v.ColumnList, columnName)
	}
}

func (v *visitor) AddColumnComment(columnName string, columnComment string) {
	v.ColumnCommentMap[columnName] = columnComment
}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	var funcArgs []ast.ExprNode

	dbName := ""
	tableName := ""
	columnName := ""
	tableComment := ""
	columnComment := ""
	astType := reflect.TypeOf(in).String()

	if common.StringInSlice(astType, v.visitSqlList) {
		v.toParse = true
		// 获取语句类型
		v.SqlType = strings.Split(astType, ".")[1]
	}

	if v.toParse {
		//fmt.Println(astType)

		switch in.(type) {
		case *ast.TableName:
			// 获取表名称
			tableName = in.(*ast.TableName).Name.L
			v.AddTable(tableName)
			// 获取数据库名称
			dbName = in.(*ast.TableName).Schema.L
			if dbName != "" {
				v.AddDb(dbName)
			}
		case *ast.CreateTableStmt:
			// 获取表注释
			tableName = in.(*ast.CreateTableStmt).Table.Name.L

			for _, tableOption := range in.(*ast.CreateTableStmt).Options {
				if tableOption.Tp == ast.TableOptionComment {
					tableComment = tableOption.StrValue
					v.AddTableComment(tableName, tableComment)
					break
				}
			}
		case *ast.AlterTableStmt:
			// 获取表注释
			tableName = in.(*ast.AlterTableStmt).Table.Name.L

			for _, tableSpec := range in.(*ast.AlterTableStmt).Specs {
				for _, tableOption := range tableSpec.Options {
					if tableOption.Tp == ast.TableOptionComment {
						tableComment = tableOption.StrValue
						v.AddTableComment(tableName, tableComment)
						break
					}
				}
			}
		case *ast.SelectField:
			// 获取字段名称
			expr := in.(*ast.SelectField).Expr
			if expr == nil && in.(*ast.SelectField).WildCard != nil {
				columnName = "*"
				v.AddColumn(columnName)
			} else if expr != nil {
				switch expr.(type) {
				case *ast.AggregateFuncExpr:
					funcArgs = expr.(*ast.AggregateFuncExpr).Args
				case *ast.FuncCallExpr:
					funcArgs = expr.(*ast.FuncCallExpr).Args
				case *ast.WindowFuncExpr:
					funcArgs = expr.(*ast.WindowFuncExpr).Args
				case *ast.ColumnNameExpr:
					columnName = expr.(*ast.ColumnNameExpr).Name.Name.L
					v.AddColumn(columnName)
				}

				for _, arg := range funcArgs {
					switch arg.(type) {
					case *ast.ColumnNameExpr:
						columnName = arg.(*ast.ColumnNameExpr).Name.Name.L
						v.AddColumn(columnName)
					}
				}
			}
		case *ast.ColumnDef:
			// 获取字段注释
			columnName := in.(*ast.ColumnDef).Name.Name.L
			v.AddColumn(columnName)

			for _, columnOption := range in.(*ast.ColumnDef).Options {
				if columnOption.Tp == ast.ColumnOptionComment {
					columnComment = columnOption.Expr.(*driver.ValueExpr).Datum.GetString()
				}
			}

			v.AddColumnComment(columnName, columnComment)
		case *ast.ColumnName:
			columnName := in.(*ast.ColumnName).Name.L
			v.AddColumn(columnName)
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
