# sqlparser-go

# 1. how to build

```
cd /path/to/code
GOOS=darwin && GOARCH=amd64 && go build -o ./bin/sqlparser-mac main.go

CGO_ENABLED=0 && GOOS=linux && GOARCH=amd64 && go build -o ./bin/sqlparser-linux main.go

CGO_ENABLED=0 && GOOS=windows && GOARCH=amd64 && go build -o ./bin/sqlparser-windows main.go

sql="select col1, col2 from db01.t01"
./bin/sqlparser-mac --sql=${sql}
```
output:
```json
{"sql_type":"SelectStmt","db_names":["db01"],"table_names":["t01"],"table_comments":{},"column_names":["col1","col2"],"column_comments":{}}
```

# 2. how to import

```
import github.com/romberli/sqlparser-go/parser

func main() {
    sql := `CREATE TABLE ` + "db01.`t01`" + `(
	id bigint(20) comment '主键ID',
	col1 varchar(64) NOT NULL,
	col2 varchar(64)  NOT NULL,
	col3 varchar(64) NOT NULL,
	col4 mediumtext,
	col5 mediumtext,
	created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
	PRIMARY KEY (id),
	KEY idx_col1_col2_col3 (col1, col2, col3)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;`
    
    result, warns, err := p.Parse(sql)
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
```
output:
```json
{"sql_type":"CreateTableStmt","db_names":["db01"],"table_names":["t01"],"table_comments":{},"column_names":["id","col1","col2","col3","col4","col5","created_at","last_updated_at"],"column_comments":{"col1":"","col2":"","col3":"","col4":"","col5":"","created_at":"","id":"主键ID","last_updated_at":"最后更新时间"}}
```
