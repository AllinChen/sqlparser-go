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
	p := parser.NewSQLParser()

	f.SQL = `CREATE TABLE ` + "`t_order_optimize_result`" + `(
	id bigint(20) comment '主键ID',
	planned_id varchar(64) NOT NULL,
	shop_id varchar(64)  NOT NULL,
	task_id varchar(64) NOT NULL,
	planed_orders mediumtext,
	route mediumtext,
	created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
	PRIMARY KEY (id),
	KEY idx_order_optimize_result_planned_id_shop_id_task_id (planned_id,shop_id,task_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;`
	// f.SQL = `alter table posdm_order_item_promotion modify column  p_promo_amt decimal(19,2) DEFAULT '0.00' comment '商品促销行优惠总金额';`
	// f.SQL = `update t01 set id=2 where id=1;`
	// f.SQL = `replace into t01 values(1) ;`
	// f.SQL = `select t.col1 from (select name as name_new from t01 as t_001) t inner join db01.t02 on t.mobile=t02.phone and t.date=t02.datetime order by date limit 0,10;`
	// f.SQL = `select count(name) from t_sms_message where create_time>'2019-02-01 00:00:00' and create_time<'2019-03-01 00:00:00' and sms_type='1' and char_length(message) > 67;`
	// f.SQL = `SELECT t.name, t.* FROM yh_after_sales_order_header t where refund_status="refunding" ;select id from t01;`
	// f.SQL = `select t02.*, id id_1 from db01.t01 inner join db02.t02 on t01.id=t02.id;`
	// f.SQL = `alter table t01 change column col1 col2 int(10) comment 'ddd';`
	// f.SQL = `alter table t01 comment 'ddd';`
	// f.SQL = "CREATE TABLE `yh_cash_gift` (`id` bigint(20) NOT NULL COMMENT 'id，非自增');"
	// f.SQL = `alter table posdm_order_item_promotion change column  p_promo_amt p_promo_amt_new decimal(19,2) DEFAULT '0.00' comment 'ddd';`
	// f.SQL = `insert into t01(col1, col2, col3) values(1, 1, 1, 1);`
	// f.SQL = `select * from t01`

	result, warns, err := p.Parse(f.SQL)
	if warns != nil || err != nil {
		fmt.Printf("parse warn: %v\nparse error: %v\n", warns, err)
		os.Exit(1)
	}

	jsonBytes, err := result.Marshal()
	if err != nil {
		fmt.Printf("marshal error: \n%s", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(jsonBytes))
}
