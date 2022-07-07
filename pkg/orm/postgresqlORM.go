package orm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func DBconnect() {
	//host=rm-2ze81r8ys3tw0s1k8wo.pg.rds.aliyuncs.com port=3432 user=postgres dbname=spider password=postgres123! sslmode=disable
	//root:root@52.83.252.118:3306/dbname?charset=utf8&parseTime=True&loc=Local
	db, err := gorm.Open("mysql", "root:root@(52.83.252.118:3306)/spider?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	db.CreateTable()
}
