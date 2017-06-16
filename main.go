package main

import (
	"awsUtils/dynomadb"
	_ "awsUtils/s3u"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
	"awsUtils/s3u"
	"path/filepath"
	"io/ioutil"
)

var (

	region                  = ""
	accessKey               = ""
	secretAccessKey         = ""
	bucket                 	= ""
)

func main() {

	QueryOne("","","")
}

func s3Test()  {
	s3 := s3u.NewS3Util(region,accessKey,secretAccessKey)
	s3.List(bucket)
	s3.Presign("",bucket )
	fmt.Println("over")
}

func Export()  {
	filepath.Walk("E:/var", func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir(){
			return nil
		}
		f,_ := os.Open(path)
		bytes,_:=ioutil.ReadAll(f)
		exportDate(string(bytes))
		return nil
	})
}

func exportDate(jsonString string)  {
	dy := dynomadb.NewdynoUtil(region,accessKey,secretAccessKey)
	m := make(map[string][]string)
	err := json.Unmarshal([]byte(jsonString),&m)
	if err!= nil{
		fmt.Println(err)
	}
	title := ""
	heads := strings.Split(title,",")
	dates  := make([][]string,0)
	dates = append(dates,heads)
	eid := ""
	for enterpriseId,uniqueIds := range m{
		eid = enterpriseId
		for _, uniqueId := range uniqueIds{
			param := make(map[string]string)
			param[""] 		= enterpriseId
			param[""] 	= uniqueId
			result := dy.Query("",param)
			value := make([]string,0)
			for _,k := range heads{
				value = append(value,strings.Replace(aws.StringValue(result[k]),",","，",-1))
			}
			dates = append(dates,value)
		}
	}
	tocsv(eid,dates)
}

func QueryOne(eid,uid,tableName string)  {
	dy := dynomadb.NewdynoUtil(region,accessKey,secretAccessKey)
	param := make(map[string]string)
	param[""] 		= eid
	param[""] 	= uid
	result := dy.Query(tableName,param)
	fmt.Println(result)
}

func tocsv(name string, data  [][]string)  {
	f, err := os.Create(name+".csv")//创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f)//创建一个新的写入文件流

	w.WriteAll(data)//写入数据
	w.Flush()
}
