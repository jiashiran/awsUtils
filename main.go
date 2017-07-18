package main

import (
	"awsUtils/dynomadb"
	"awsUtils/s3u"
	_ "awsUtils/s3u"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	region          = ""
	accessKey       = ""
	secretAccessKey = ""
	bucket          = ""
)

func main() {
	/*f,_:= os.Open("result.txt")
	bs,_:=ioutil.ReadAll(f)
	s := string(bs)
	lines := strings.Split(s,"\n")
	m := make(map[string]int)
	for _,l := range lines{
		name := strings.Split(l," ")[0]
		if m[name] > 0{
			fmt.Println(name)
		}
		m[name] = 1
	}*/

	s3Test()
}

func s3Test() {
	s3 := s3u.NewS3Util(region, accessKey, secretAccessKey)
	//s3.List(bucket,"")
	//s3.Presign("",bucket,1 )
	//fmt.Println("over")
	f, err := os.OpenFile("url.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil{
		fmt.Println(err)
	}
	defer f.Close()
	c:=0
	path1:="D:/downloads/筛选录音"
	path2:="D:/downloads/Compressed"
	_,_ = path1,path2
	filepath.Walk(path2, func(path string, info os.FileInfo, err error) error {
		if c >= 1{
			os.Exit(0)
		}
		if info == nil || info.IsDir() {
			return nil
		}
		fmt.Println(path)
		//s3.Upload(path,info.Name(),bucket)
		url := s3.Presign(info.Name(), bucket, 120)
		f.WriteString(info.Name() + ";" + url + "\n")
		f.Sync()
		c++
		return nil
	})
	fmt.Println("Presign count:",c)
}

func Export() {
	filepath.Walk("E:/var", func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		f, _ := os.Open(path)
		bytes, _ := ioutil.ReadAll(f)
		exportDate(string(bytes))
		return nil
	})
}

func exportDate(jsonString string) {
	dy := dynomadb.NewdynoUtil(region, accessKey, secretAccessKey)
	m := make(map[string][]string)
	err := json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		fmt.Println(err)
	}
	title := ""
	heads := strings.Split(title, ",")
	dates := make([][]string, 0)
	dates = append(dates, heads)
	eid := ""
	for enterpriseId, uniqueIds := range m {
		eid = enterpriseId
		for _, uniqueId := range uniqueIds {
			param := make(map[string]string)
			param[""] = enterpriseId
			param[""] = uniqueId
			result := dy.Query("", param)
			value := make([]string, 0)
			for _, k := range heads {
				value = append(value, strings.Replace(aws.StringValue(result[k]), ",", "，", -1))
			}
			dates = append(dates, value)
		}
	}
	tocsv(eid, dates)
}

func QueryOne(eid, uid, tableName string) {
	dy := dynomadb.NewdynoUtil(region, accessKey, secretAccessKey)
	param := make(map[string]string)
	param[""] = eid
	param[""] = uid
	result := dy.Query(tableName, param)
	fmt.Println(result)
}

func tocsv(name string, data [][]string) {
	f, err := os.Create(name + ".csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流

	w.WriteAll(data) //写入数据
	w.Flush()
}
