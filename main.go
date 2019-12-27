package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"os"
	"sqlparser-go/lib/common"
	parser2 "sqlparser-go/parser"
)

func main() {
	f := common.MyFlag{}
	f.Init()
	parser := parser2.SqlParser{}
	result := make(map[string]interface{})

	//f.Sql = `CREATE TABLE ` + "`t_order_optimize_result`" + `(
	//id bigint(20) comment '主键ID',
	//planned_id varchar(64) NOT NULL,
	//shop_id varchar(64)  NOT NULL,
	//task_id varchar(64) NOT NULL,
	//planed_orders mediumtext,
	//route mediumtext,
	//created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	//last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
	//PRIMARY KEY (id),
	//KEY idx_order_optimize_result_planned_id_shop_id_task_id (planned_id,shop_id,task_id)
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;`
	//f.Sql = `alter table posdm_order_item_promotion modify column  p_promo_amt decimal(19,2) DEFAULT '0.00' comment '商品促销行优惠总金额';`
	//f.Sql = `update t01 set id=2 where id=1;`
	//f.Sql = `replace into t01 values(1) ;`
	//f.Sql = `select t.col1 from (select name as name_new from t01 as t_001) t inner join t02 on t.mobile=t02.phone and t.date=t02.datetime order by date limit 0,10;`
	//f.Sql = `select count(name) from t_sms_message where create_time>'2019-02-01 00:00:00' and create_time<'2019-03-01 00:00:00' and sms_type='1' and char_length(message) > 67;`
	//f.Sql = `SELECT t.name, t.* FROM yh_after_sales_order_header t where refund_status="refunding" ;select id from t01;`
	//f.Sql = `select t02.*, id id_1 from db01.t01 inner join db02.t02 on t01.id=t02.id;`
	//f.Sql = `alter table t01 change column col1 col2 int(10) comment 'ddd';`
	//f.Sql = `alter table t01 comment 'ddd';`
	//f.Sql = "CREATE TABLE `yh_cash_gift` (`id` bigint(20) NOT NULL COMMENT 'id，非自增');"
	//f.Sql = `alter table posdm_order_item_promotion change column  p_promo_amt p_promo_amt_new decimal(19,2) DEFAULT '0.00' comment 'ddd';`
	//f.Sql = `insert into t01(col1, col2, col3) values(1, 1, 1, 1);`
	//f.Sql = `select * from t01`
	if v, warns, err := parser.ParseSql(f.Sql); warns == nil && err == nil {
		result["sqlType"] = v.SqlType
		result["dbNames"] = v.DbList
		result["tableNames"] = v.TableList
		result["tableComments"] = v.TableCommentMap
		result["columnNames"] = v.ColumnList
		result["columnComments"] = v.ColumnCommentMap

		if result, err := json.Marshal(result); err == nil {
			fmt.Println(string(result))
		} else {
			fmt.Printf("marshal error: \n%v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("parse warn: %v\nparse error: %v\n", warns, err)
		os.Exit(1)
	}
}
