# sqlparser-go

#1. how to build

```
GOOS=darwin && GOARCH=amd64 && go build -o ./bin/sql-fingerprint-mac *.go

CGO_ENABLED=0 && GOOS=linux && GOARCH=amd64 && go build -o ./bin/sql-fingerprint-linux *.go

CGO_ENABLED=0 && GOOS=windows && GOARCH=amd64 && go build -o ./bin/sql-fingerprint-windows *.go
```

#2. how to import

```
import https://github.com/romberli/sqlparser-go/parser

func main() {
    sql := `select col1, col2 from t01;`
    
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