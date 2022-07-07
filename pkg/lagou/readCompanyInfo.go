package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"spider/pkg"
	"spider/pkg/store"
	"strings"
)

var (
	companyIdMap   = make(map[string]string)
	webcall        = []string{"电话催收专员", "电话销售", "电话销售专员", "电话销售代表", "呼叫中心电销", "电销专员"}
	customerervice = []string{"电话客服", "客服代表", "客服专员", "呼叫中心客服专员"}
)

func main() {
	path := "."

	readDir(path)

	all := store.SelectAll()
	for k, v := range all {
		fmt.Println(k, v)
	}
	fmt.Println(len(all))
	all = nil
	fmt.Println(len(companyIdMap))

}

func fetchCompanyInfo(companyId string) string {
	companyInfo := GetCompanyInfo(companyId)

	if companyInfo.CompanyShortName != "" {
		bs, _ := json.Marshal(companyInfo)
		log.Println(companyInfo)
		return string(bs)
	}
	return ""
}

func readDir(path string) {
	//以只读的方式打开目录
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeDir)
	if err != nil {
		fmt.Println(err.Error())
	}
	//延迟关闭目录
	defer f.Close()
	fileInfo, _ := f.Readdir(-1)
	//操作系统指定的路径分隔符
	separator := string(os.PathSeparator)

	for _, info := range fileInfo {
		//判断是否是目录
		if info.IsDir() {
			//fmt.Println(path + separator + info.Name())
			readDir(path + separator + info.Name())
		} else {
			if strings.Contains(info.Name(), ".csv") {
				//fmt.Println("文件：" + info.Name())
				ReadLine(info.Name(), func(s string) {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println(s)
						}
					}()
					ss := strings.Split(s, ",")
					companyId := ss[0]
					city := ss[1]
					work := ss[2]
					//fmt.Println(companyId,city,work)
					for _, w := range webcall {
						if work == w {
							var companyInfo CompanyInfo
							json.Unmarshal([]byte(store.Select(companyId)), &companyInfo)
							pkg.ToTxt(companyInfo.CompanyId+","+companyInfo.CompanyShortName+""+
								","+companyInfo.CompanyFullName+","+companyInfo.IndustryField+""+
								","+companyInfo.FinanceStage+","+companyInfo.CompanySize+","+companyInfo.CompanyUrl+""+
								","+companyInfo.RegLocation+","+companyInfo.RegCapital+","+companyInfo.LegalPersonName+""+
								","+companyInfo.EstablishTime+","+companyInfo.WorkName+""+
								","+companyInfo.City+","+companyInfo.LagouURL+","+work,
								"外呼类-"+city+".csv")
						}
					}
					for _, w := range customerervice {
						if work == w {
							var companyInfo CompanyInfo
							json.Unmarshal([]byte(store.Select(companyId)), &companyInfo)
							pkg.ToTxt(companyInfo.CompanyId+","+companyInfo.CompanyShortName+""+
								","+companyInfo.CompanyFullName+","+companyInfo.IndustryField+""+
								","+companyInfo.FinanceStage+","+companyInfo.CompanySize+","+companyInfo.CompanyUrl+""+
								","+companyInfo.RegLocation+","+companyInfo.RegCapital+","+companyInfo.LegalPersonName+""+
								","+companyInfo.EstablishTime+","+companyInfo.WorkName+""+
								","+companyInfo.City+","+companyInfo.LagouURL+","+work,
								"客服类-"+city+".csv")
						}
					}
					/*if _,ok := companyIdMap[companyId];ok{
						companyIdMap[companyId] = companyIdMap[companyId] + "," + city+","+work
					}else {
						companyIdMap[companyId] = city+","+work
					}*/

				})
			}
		}
	}
}

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}
