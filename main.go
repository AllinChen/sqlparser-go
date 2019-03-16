package main

import (
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"reflect"
)

import _ "github.com/pingcap/tidb/types/parser_driver"

func elementExistsInSlice(str string, slice []string) bool {
	for i := range slice {
		if slice[i] == str {
			return true
		}
	}

	return false
}

type visitor struct {
	tableList []string
}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	//fmt.Printf("%T\n", in)
	if reflect.TypeOf(in).String() == "*ast.TableName" {
		tableName := in.(*ast.TableName).Name.L
		if !elementExistsInSlice(tableName, v.tableList) {
			v.tableList = append(v.tableList, tableName)
		}
	}
	return in, false
}

func (v *visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func main() {

	sql := "SELECT /*+ TIDB_SMJ(employees) */ emp_no, first_name, last_name " +
		"FROM employees emp USE INDEX (last_name) " +
		"" +
		"where last_name='Aamodt' and gender='F' and birth_date > '1960-01-01'"

	sql = `select t1.id, t2.name from (select * from tb01 group by name order by id desc limit 10) t1 , tb02 t2 where t1.id=t2.id
    union all
    select t3.id, t4.name from tb03 t3, tb02 t4 where t4.id=t3.id and t3.name = "from table01 to tb02" group by t3.name;`

	sql = `SELECT a.time_updated_server/1000,
content,
nick,
name
FROM table1 a
JOIN table2 b ON a.sender_id = b.user_id
JOIN table3 c ON a.channel_id = c.channel_id
JOIN table4 d ON c.store_id = d.store_id
WHERE sender_id NOT IN
  (SELECT user_id
   FROM table5
   WHERE store_id IN ('agent_store:1',
                                     'ask:1'))
   GROUP BY 1,2,3,4
   HAVING sum(1) > 500
   ORDER BY 1 ASC;`
	sqlParser := parser.New()
	stmtNodes, warn, err := sqlParser.Parse(sql, "", "")
	if warn != nil || err != nil {
		fmt.Printf("parse warn:\n%v\n%s", warn, sql)
		fmt.Printf("parse error:\n%v\n%s", err, sql)
		return
	}
	for _, stmtNode := range stmtNodes {
		v := visitor{}
		//tableList := []string{}
		stmtNode.Accept(&v)
		fmt.Println(v.tableList)
	}
}
