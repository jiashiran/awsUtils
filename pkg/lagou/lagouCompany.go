package main

import (
	"github.com/buger/jsonparser"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type CompanyInfo struct {
	CompanyId        string //企业ID
	CompanyShortName string //公司简称
	CompanyFullName  string //公司全名
	IndustryField    string //行业领域
	FinanceStage     string //融资
	CompanySize      string //公司规模
	CompanyUrl       string //公司网址
	RegLocation      string //公司地址
	RegCapital       string //注册资本
	LegalPersonName  string //注册法人
	EstablishTime    string //成立日期
	WorkName         string //岗位名称
	City             string //城市
	LagouURL         string
}

func GetCompanyInfo(companyId string) CompanyInfo {
	var companyInfo CompanyInfo
	data := getCompanyData(companyId)
	if data == nil {
		return companyInfo
	}
	//log.Println(companyName)
	companyShortName, _ := jsonparser.GetString(data, "coreInfo", "companyShortName")
	companyFullName, _ := jsonparser.GetString(data, "companyBusinessInfo", "companyName")
	industryField, _ := jsonparser.GetString(data, "baseInfo", "industryField")
	financeStage, _ := jsonparser.GetString(data, "baseInfo", "financeStage")
	companySize, _ := jsonparser.GetString(data, "baseInfo", "companySize")
	companyUrl, _ := jsonparser.GetString(data, "coreInfo", "companyUrl")
	regLocation, _ := jsonparser.GetString(data, "companyBusinessInfo", "regLocation")
	regCapital, _ := jsonparser.GetString(data, "companyBusinessInfo", "regCapital")
	legalPersonName, _ := jsonparser.GetString(data, "companyBusinessInfo", "legalPersonName")
	establishTime, _ := jsonparser.GetString(data, "companyBusinessInfo", "establishTime")
	companyInfo = CompanyInfo{
		CompanyId:        companyId,
		CompanyShortName: strings.ReplaceAll(companyShortName, ",", ";"),
		CompanyFullName:  strings.ReplaceAll(companyFullName, ",", ";"),
		IndustryField:    strings.ReplaceAll(industryField, ",", ";"),
		FinanceStage:     strings.ReplaceAll(financeStage, ",", ";"),
		CompanySize:      strings.ReplaceAll(companySize, ",", ";"),
		CompanyUrl:       strings.ReplaceAll(companyUrl, ",", ";"),
		RegLocation:      strings.ReplaceAll(regLocation, ",", ";"),
		RegCapital:       strings.ReplaceAll(regCapital, ",", ";"),
		LegalPersonName:  strings.ReplaceAll(legalPersonName, ",", ";"),
		EstablishTime:    strings.ReplaceAll(establishTime, ",", ";"),
	}
	return companyInfo
}

func getCompanyDataByChromedp(companyId string) []byte {
	//ctx, cancel := chromedp.NewContext(context.Background())
	//defer cancel()
	time.Sleep(5 * time.Second)
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.lagou.com/gongsi/"+companyId+".html"),
		chromedp.OuterHTML("body", &res, chromedp.ByQuery),
	)
	if err != nil {
		log.Println("https://www.lagou.com/gongsi/"+companyId+".html err:", err)
	}
	return []byte(res)
}

func getCompanyData(companyId string) []byte {
	url := "https://www.lagou.com/gongsi/" + companyId + ".html"
	resp, err := http.Get(url)
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		if !strings.Contains(string(body), `id="companyInfoData"`) {
			body = getCompanyDataByChromedp(companyId)
		}
		//log.Println(string(body))
		lines := strings.Split(string(body), "\n")
		for _, line := range lines {
			if strings.Contains(line, `id="companyInfoData"`) {
				start := strings.Index(line, "{")
				end := strings.LastIndex(line, "}")
				//log.Println(line[start:end+1])
				return []byte(line[start : end+1])
			}
		}
		//fmt.Println(string(body))
	} else {
		log.Println("get", url, "error:", err)
	}
	return nil
}
