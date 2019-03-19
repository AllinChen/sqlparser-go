package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"os"
)

func main() {
	f := MyFlag{}
	f.Init()
	parser := SqlParser{}
	result := make(map[string]interface{})

	//	sql = `CREATE TABLE ` + "`t_order_optimize_result`" + `(
	//  id bigint(20) comment '主键ID',
	//  planned_id varchar(64) NOT NULL,
	//  shop_id varchar(64)  NOT NULL,
	//  task_id varchar(64) NOT NULL,
	//  planed_orders mediumtext,
	//  route mediumtext,
	//  created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	//  last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
	//  PRIMARY KEY (id),
	//  KEY idx_order_optimize_result_planned_id_shop_id_task_id (planned_id,shop_id,task_id)
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;`
	//	sql = `alter table posdm_order_item_promotion modify column  p_promo_amt decimal(19,2) DEFAULT '0.00' comment '商品促销行优惠总金额';`
	//	sql = `update t01 set id=2 where id=1;`
	//	sql = `replace into t01 values(1) ;`
	//	sql = `select t.col1 from (select name as name_new from t01 as t_001) t inner join t02 on t.mobile=t02.phone and t.date=t02.datetime order by date limit 0,10;`
	//	sql = `select count(name) from t_sms_message where create_time>'2019-02-01 00:00:00' and create_time<'2019-03-01 00:00:00' and sms_type='1' and char_length(message) > 67;`
	//	sql = `SELECT t.name, t.* FROM yh_after_sales_order_header t where refund_status="refunding" ;select id from t01;`
	//	sql = `select *, id from t01;`

	if v, warns, err := parser.ParseSql(sql); warns == nil && err == nil {
		result["tables"] = v.tableList
		result["columns"] = v.columnList
		result["comments"] = v.columnCommentMap

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
