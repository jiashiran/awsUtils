package jijin

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

func DbInit() {
	var err error
	Db, err = gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/spider?charset=utf8&parseTime=True&loc=Local")
	//defer Db.Close()
	if err != nil {
		fmt.Println(err)
	}
	Db.LogMode(true)
}

func Find(info *Info) {
	Db.Table("info").Where("code = ?", info.Code).Find(&info)
}

type Stock struct {
	ID         uint `gorm:"primary_key"`
	Name       string
	Value      float64
	Num        string
	Proportion string
}

func (s Stock) TableName() string {
	return "stock"
}

func SelectAll(pageIndex int32, pageSize int32) []Info {
	c := make([]Info, 0)
	offset := int64((pageIndex - 1)) * int64(pageSize)
	fmt.Println(offset)
	Db.Offset(offset).Limit(pageSize).Order("ID").Find(&c)
	return c
}

func SelectAllStock(pageIndex int32, pageSize int32) []Stock {
	c := make([]Stock, 0)
	offset := int64((pageIndex - 1)) * int64(pageSize)
	fmt.Println(offset)
	Db.Offset(offset).Limit(pageSize).Order("ID").Find(&c)
	return c
}

type Info struct {
	ID      uint `gorm:"primary_key"`
	Code    string
	Amount  string
	Details []Detail `gorm:"ForeignKey:InfoID;save_associations:true"`
}

func (info Info) TableName() string {
	return "info"
}

type Detail struct {
	InfoID     uint
	Name       string
	Proportion string
	Code       string
}

func (detail Detail) TableName() string {
	return "detail"
}

func createTable() {
	Db.DropTable("info")
	Db.DropTable("detail")

	result := Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Table("info").CreateTable(&Info{})
	if result.Error != nil {
		fmt.Println("create info company err:", result.Error)
	}
	result = Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Table("detail").CreateTable(&Detail{})
	if result.Error != nil {
		fmt.Println("create table detail err:", result.Error)
	}

	result = Db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Table("stock").CreateTable(&Stock{})
	if result.Error != nil {
		fmt.Println("create Stock err:", result.Error)
	}
}
