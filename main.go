package main

import (
	"fmt"
	"os"

	"github.com/romberli/sqlparser-go/lib/common"
	"github.com/romberli/sqlparser-go/parser"
)

func main() {
	f := common.MyFlag{}
	f.Init()
	p := parser.NewParser()

	// f.SQL = `CREATE TABLE ` + "`t01`" + `(
	// id bigint(20) comment '主键ID',
	// col1 varchar(64) NOT NULL,
	// col2 varchar(64)  NOT NULL,
	// col3 varchar(64) NOT NULL,
	// col4 mediumtext,
	// col5 mediumtext,
	// created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	// last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
	// PRIMARY KEY (id),
	// KEY idx_col1_col2_col3 (col1, col2, col3)
	// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;`
	// f.SQL = `select t02.*, id id_1 from db01.t01 inner join db02.t02 on t01.id=t02.id;`
	// f.SQL = `alter table t01 change column col1 col2 int(10) comment 'ddd';`
	// f.SQL = `alter table t01 comment 'ddd';`
	// f.SQL = `insert into t01(col1, col2, col3) values(1, 1, 1, 1);`
	// f.SQL = `select * from t01`

	result, warns, err := p.Parse(f.SQL)
	if err != nil {
		fmt.Printf("parse error: %s", err.Error())
		os.Exit(1)
	}
	if warns != nil {
		for _, warn := range warns {
			fmt.Printf("parse warn: %s", warn.Error())
		}
		os.Exit(1)
	}

	jsonBytes, err := result.Marshal()
	if err != nil {
		fmt.Printf("marshal error: \n%s", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(jsonBytes))
}
