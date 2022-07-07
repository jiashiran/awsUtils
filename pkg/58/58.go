package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/tunny"
	"github.com/buger/jsonparser"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"os"
	"spider/pkg"
	"spider/pkg/store"

	"strings"
	"time"
)

var (
	cityDomin      = make(map[string]string)
	webcall        = []string{"电话催收专员", "电话销售", "电话销售专员", "电话销售代表", "呼叫中心电销", "电销专员"}
	customerervice = []string{"电话客服", "客服代表", "客服专员", "呼叫中心客服专员", "客服总监"}
	companyMap     = make(map[string]string)
	pool           = tunny.NewFunc(1, func(l interface{}) interface{} {
		getCompanyInfoBychromedp(fmt.Sprint(l), "")
		return nil
	})
	ctx       context.Context
	cancel    context.CancelFunc
	optCancel context.CancelFunc
	c         *fasthttp.Client
	cookies   []*network.Cookie = make([]*network.Cookie, 0)
)

type CompanyInfo struct {
	CompanyId        string `json:"company_id"`      //企业ID
	CompanyShortName string `json:"aliasName"`       //公司简称
	CompanyFullName  string `json:"entName"`         //公司全名
	IndustryField    string `json:"industryText"`    //行业领域
	FinanceStage     string `json:""`                //融资
	CompanySize      string `json:"sizeText"`        //公司规模
	CompanyUrl       string `json:"url"`             //公司网址
	RegLocation      string `json:"regLocation"`     //公司地址
	RegCapital       string `json:"regCapital"`      //注册资本
	LegalPersonName  string `json:"legalPersonName"` //注册法人
	EstablishTime    string `json:"createTime"`      //成立日期
	WorkName         string `json:"work_name"`       //岗位名称
	City             string `json:"city"`            //城市
	Contact          string `json:"contact"`         //联系人
	Phone            string
	Email            string `json:"email"`    //
	TypeText         string `json:"typeText"` //公司类型
	URL              string `json:"url"`
}

func init() {
	cityDomin["北京"] = "https://bj.58.com"
	cityDomin["上海"] = "https://sh.58.com"
	cityDomin["广州"] = "https://gz.58.com"
	cityDomin["深圳"] = "https://sz.58.com"
	cityDomin["天津"] = "https://tj.58.com"
	cityDomin["杭州"] = "https://hz.58.com"
	cityDomin["成都"] = "https://cd.58.com"
	cityDomin["南京"] = "https://nj.58.com"
	cityDomin["西安"] = "https://xa.58.com"
	cityDomin["苏州"] = "https://su.58.com"
	cityDomin["石家庄"] = "https://sjz.58.com"
	cityDomin["重庆"] = "https://cq.58.com"
	cityDomin["郑州"] = "https://zz.58.com"
	cityDomin["济南"] = "https://jn.58.com"
	cityDomin["东莞"] = "https://dg.58.com"
	cityDomin["青岛"] = "https://qd.58.com"
	cityDomin["宁波"] = "https://nb.58.com"
	cityDomin["无锡"] = "https://wx.58.com"
	cityDomin["佛山"] = "https://fs.58.com"
	cityDomin["厦门"] = "https://xm.58.com"
	cityDomin["长沙"] = "https://cs.58.com"
	cityDomin["大连"] = "https://dl.58.com"
	cityDomin["合肥"] = "https://hf.58.com"
	cityDomin["武汉"] = "https://wh.58.com"
	cityDomin["开封"] = "https://kaifeng.58.com"
	cityDomin["南宁"] = "https://nn.58.com"
	cityDomin["绵阳"] = "https://mianyang.58.com"
	cityDomin["沈阳"] = "https://sy.58.com"
	cityDomin["长春"] = "https://cc.58.com"
	cityDomin["哈尔滨"] = "https://hrb.58.com"
	cityDomin["济南"] = "https://jn.58.com"
	cityDomin["淄博"] = "https://zb.58.com"

	//f, _ := os.OpenFile("logrus.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	//
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", true),
		chromedp.WindowSize(1280, 1024),
	)
	var allocCtx context.Context
	allocCtx, optCancel = chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

}

func main() {
	//getCompnayId()
	//readcompanys()
	writeCsv()
}

func writeCsv() {
	pkg.ReadDir("resource/company/", func(path string, info os.FileInfo) {
		if !strings.Contains(info.Name(), "csv") {
			return
		}
		filePath := path + "/" + info.Name()
		ReadLine(filePath, func(s string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("err:", err)
				}
			}()
			ss := strings.Split(s, ",")
			companyId := ss[0]
			url := ss[1]
			_ = url
			city := ss[2]
			_ = city
			keyWord := ss[3]
			_ = keyWord
			r := store.Select(companyId)
			if !strings.Contains(r, "获取信息失败") && r != "" {
				data := []byte(r)
				companyFullName, _ := jsonparser.GetString(data, "data", "entDetail", "aliasName")
				industryField, _ := jsonparser.GetString(data, "data", "entDetail", "industryText")
				companySize, _ := jsonparser.GetString(data, "data", "entDetail", "sizeText")
				companyUrl, _ := jsonparser.GetString(data, "data", "entDetail", "url")
				legalPersonName, _ := jsonparser.GetString(data, "data", "entDetail", "bussiness", "legalPersonName")
				establishTime, _ := jsonparser.GetString(data, "data", "entDetail", "bussiness", "createTime")
				email, _ := jsonparser.GetString(data, "data", "entDetail", "email")
				text := city + "," + keyWord + "," + companyFullName + "," + industryField + "," +
					companySize + "," + companyUrl + "," + legalPersonName +
					"," + establishTime + "," + email
				fmt.Println(text)
				pkg.ToTxt(text, "result/58/"+city+"_"+keyWord+".csv")
			}

		})
	})
}

func readcompanys() {
	pkg.ReadDir("/var/58/", func(path string, info os.FileInfo) {
		//fmt.Printf(info.Name())
		if !strings.Contains(info.Name(), "csv") {
			return
		}
		filePath := path + "/" + info.Name()
		ReadLine(filePath, func(s string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("err:", err)
				}
			}()
			ss := strings.Split(s, ",")
			companyId := ss[0]
			url := ss[1]
			city := ss[2]
			keyWord := ss[3]
			fmt.Println(companyId, url, city, keyWord)
			//getCompanyInfo(companyId,ss,s)
			getCompanyInfoBychromedp(companyId, url)
		})
	})
}

func getCompanyInfoBychromedp(companyId string, url string) (CompanyInfo, string) {
	var c CompanyInfo
	if store.Get(companyId) != "" {
		return c, ""
	}
	time.Sleep(5 * time.Second)
	var res string
	headers := make(map[string]interface{})
	if len(cookies) == 0 {
		headers["cookie"] = `id58=e87rZl1pB4M+Go6LBJFvAg==; 58tj_uuid=d7f0e55d-57f6-455d-80c1-8d6f8002d414; als=0; xxzl_deviceid=6Fo3dL0IQuA87b9dO7uA2cD0MB4BefYYB%2BwFJ2ugwVhqPBWzY4%2FdxeXvSZ%2FVqgE9; wmda_uuid=8cf3ca7c283d19e2eb2c27a64c52100b; wmda_new_uuid=1; wmda_visited_projects=%3B1731916484865; gr_user_id=832c4ac4-e52c-4240-a5b7-f344475e5204; xxzl_smartid=4cb0f68cae46538b6b49695088566898; myfeet_tooltip=end; bangtoptipclose=1; Hm_lvt_b4a22b2e0b326c2da73c447b956d6746=1567753885; city=cs; 58home=cs; __utmz=253535702.1567855106.8.4.utmcsr=passport.58.com|utmccn=(referral)|utmcmd=referral|utmcct=/login/; ppStore_fingerprint=1155522383D8F215ED5ACDB5594FED61A514F51B6DF9072E%EF%BC%BF1567861942268; Hm_lvt_74e7f6ad49774e9bb7a4b0fae4c4f668=1567861914,1567861923,1567861937,1567863131; __utmc=253535702; www58com="UserID=66038324157712&UserName=vknhshwie"; 58cooper="userid=66038324157712&username=vknhshwie"; 58uname=vknhshwie; PPU="UID=66038324157712&UN=vknhshwie&TT=de382cdbb3f95c85e4f49a7bb3fc2058&PBODY=DU3I41RZd-igXGWRZYYPuU5LbKWOpZ05qmmizh9Hp2jH_TFOPZ8vopVU0mNXA2gQHFDUPhqEmYYctfG3CbenCxe4I43f2apPploIdP6-XipRFeLRvMQSq2x7nKXj0id2_V__eoAelJAEUNoMOwnAJURIa-0gqAVxUKdSEvShSC0&VER=1"; __utma=253535702.1528588831.1567666785.1568083987.1568093151.14; new_uv=16; utm_source=; spm=; init_refer=; new_session=0; __utmt_pageTracker=1; JSESSIONID=EEB6B9FAA1D0909F4B65CC9E2B1FE9C6; __utmb=253535702.4.10.1568093151`
	}
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://qy.58.com/ent/detail/"+companyId),
		chromedp.InnerHTML("pre", &res, chromedp.ByQuery),
	)
	if err != nil {
		fmt.Println("err:", err)
	}
	//c := dealRes(res)

	fmt.Println(res)
	/*if (c.CompanyShortName != "" || c.CompanyFullName != "") && !strings.Contains(c.Phone, "登录") {
		bs, _ := json.Marshal(c)
		store.Update("companyInfo.Id."+companyId, string(bs))
	}*/
	store.Update(companyId, res)
	return c, res
}

func dealRes(res string) CompanyInfo {
	var companyInfo CompanyInfo
	if res == "" {
		return companyInfo
	}
	defer func() {
		if err := recover(); err != nil {
			log.Println("dealRes err:", err, res)
		}
	}()

	res = strings.ReplaceAll(res, `class="td_c1">\\n`, "")
	//log.Println(res)
	/*ss := strings.Split(res,"\n")
	for i,v := range ss{
		log.Println(i,v)
	}*/
	c1 := get(res, `<td class="td_c1">`, `</td>`)
	c2 := get(res, `<td class="td_c2">`, `</td>`)
	c3 := get(res, `<td class="td_c3">`, `</td>`)

	companyInfo.CompanyFullName = strings.ReplaceAll(strings.ReplaceAll(c1[0], "\n", ""), "	", "")
	companyInfo.TypeText = strings.ReplaceAll(strings.ReplaceAll(c1[1], "\n", ""), "	", "")
	companyInfo.CompanySize = strings.ReplaceAll(strings.ReplaceAll(c3[0], "\n", ""), "	", "")

	companyInfo.IndustryField = strings.ReplaceAll(strings.ReplaceAll(c2[0], "\n", ""), "	", "")
	companyInfo.Contact = strings.ReplaceAll(strings.ReplaceAll(c2[1], "\n", ""), "	", "")
	companyInfo.Phone = strings.ReplaceAll(strings.ReplaceAll(c3[1], "\n", ""), "	", "")

	if companyInfo.Phone != "" {
		companyInfo.Phone = get(companyInfo.Phone, `<div id="freecall">`, "</div>")[0]
	}

	companyInfo.RegLocation = strings.ReplaceAll(strings.ReplaceAll(c1[2], "\n", ""), "	", "")
	ss := strings.Split(companyInfo.RegLocation, "<span>")
	if len(ss) > 1 {
		ss = strings.Split(ss[1], "</span>")
		if len(ss) > 1 {
			companyInfo.RegLocation = ss[0]
		}
	}
	companyInfo.Email = strings.ReplaceAll(strings.ReplaceAll(c2[2], "\n", ""), "	", "")
	if companyInfo.Email != "" {
		companyInfo.Email = "https://" + strings.ReplaceAll(get(companyInfo.Email, `<img src="//`, `"`)[0], `amp;`, "")
	}
	companyInfo.URL = strings.ReplaceAll(strings.ReplaceAll(c3[2], "\n", ""), "	", "")
	if companyInfo.URL != "" {
		companyInfo.URL = get(companyInfo.URL, `<a href="`, `"`)[0]
	}
	/*cName := c1[0]
	cType := c1[1]
	cPosition := c1[2]

	Contact := c2[0]
	email := c2[1]*/
	/*log.Println(c1)
	log.Println(c2)
	log.Println(c3)*/
	//log.Println(companyInfo)
	return companyInfo
}

func get(res, start, end string) []string {
	result := make([]string, 0)
	ss := strings.Split(res, start)
	for i, v := range ss {
		if i > 0 {
			value := strings.Split(v, end)[0]
			result = append(result, value)
		}
	}
	return result
}

func getCompnayId() {
	for k, _ := range cityDomin {
		for _, w := range webcall {
			start := time.Now().Unix()
			fmt.Println("getListByChromedp start:", w, k)
			getListByChromedp(k, w)
			fmt.Println("getListByChromedp end:", w, k, time.Now().Unix()-start)
		}
	}
	for k, _ := range cityDomin {
		for _, w := range customerervice {
			start := time.Now().Unix()
			fmt.Println("getListByChromedp start:", w, k)
			getListByChromedp(k, w)
			fmt.Println("getListByChromedp end:", w, k, time.Now().Unix()-start)
		}
	}
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
	_ = separator
	for _, info := range fileInfo {
		//判断是否是目录
		if info.IsDir() {
			//fmt.Println(path + separator + info.Name())
			//readDir(path + separator + info.Name())
		} else {
			if strings.Contains(info.Name(), ".csv") {
				//fmt.Println("文件：" + info.Name())
				ReadLine(info.Name(), func(s string) {
					defer func() {
						if err := recover(); err != nil {
							if s != "" {
								fmt.Println(s, err)
							}

						}
					}()
					ss := strings.Split(s, ",")
					companyId := ss[0]
					companYUrl := ss[1]
					city := ss[2]
					work := ss[3]
					_ = city
					_ = work
					//fmt.Println(companyId,city,work)
					companyMap[companyId] = companYUrl
					/*for _,w := range webcall {
						if work == w {

						}
					}
					for _,w := range customerervice {
						if work == w {

						}
					}*/

					/*if _,ok := companyIdMap[companyId];ok{
						companyIdMap[companyId] = companyIdMap[companyId] + "," + city+","+work
					}else {
						companyIdMap[companyId] = city+","+work
					}*/
					/*				str,_ := store.Client.Get("companyInfo.Id."+companyId).Result()
									if str == ""{
										return
									}
									fmt.Println(str)
									for _,w := range webcall {
										if work == w {
											var companyInfo CompanyInfo
											json.Unmarshal([]byte(str),&companyInfo)
											util.ToTxt(companyInfo.CompanyId+","+companyInfo.CompanyShortName+"" +
												","+companyInfo.CompanyFullName+","+companyInfo.IndustryField+"" +
												","+companyInfo.FinanceStage+","+companyInfo.CompanySize+","+companyInfo.CompanyUrl+"" +
												","+companyInfo.RegLocation+","+companyInfo.RegCapital+","+companyInfo.LegalPersonName+"" +
												","+companyInfo.EstablishTime+","+companyInfo.WorkName+"" +
												","+companyInfo.City+","+companyInfo.URL+","+companyInfo.Email+","+companyInfo.Contact+","+companyInfo.Phone+","+companyInfo.TypeText+","+work,
												"./bj58/外呼类-"+city+".csv")
										}
									}
									for _,w := range customerervice {
										if work == w {
											var companyInfo CompanyInfo
											json.Unmarshal([]byte(str),&companyInfo)
											util.ToTxt(companyInfo.CompanyId+","+companyInfo.CompanyShortName+"" +
												","+companyInfo.CompanyFullName+","+companyInfo.IndustryField+"" +
												","+companyInfo.FinanceStage+","+companyInfo.CompanySize+","+companyInfo.CompanyUrl+"" +
												","+companyInfo.RegLocation+","+companyInfo.RegCapital+","+companyInfo.LegalPersonName+"" +
												","+companyInfo.EstablishTime+","+companyInfo.WorkName+"" +
												","+companyInfo.City+","+companyInfo.URL+","+companyInfo.Email+","+companyInfo.Contact+","+companyInfo.Phone+","+companyInfo.TypeText+","+work,
												"./bj58/客服类-"+city+".csv")
										}
									}*/

				})
			}
		}
	}
}

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		log.Println("open err:", err)
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

func getListByChromedp(city, keyWord string) string {
	/*opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.DisableGPU,
		//chromedp.Flag("headless", true),
		//chromedp.WindowSize(1280, 1024),
	)*/
	//allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	//defer cancel()
	//ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()
	time.Sleep(5 * time.Second)
	var res string
	var pagesout string
	page := 1
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			//network.Enable().Do(ctx)
			fmt.Println("step1")
			chromedp.Navigate(cityDomin[city] + "/job/pn1/?key=" + keyWord).Do(ctx)
			fmt.Println("step2")
			chromedp.OuterHTML("html", &res).Do(ctx)
			fmt.Println(res)
			chromedp.OuterHTML("#list_con", &res, chromedp.NodeVisible, chromedp.ByID).Do(ctx)
			fmt.Println("step3")
			chromedp.OuterHTML("div.pagesout", &pagesout, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
			fmt.Println("step4")
			fmt.Println("getCompanyId", city, keyWord, page)
			page++
			getCompanyId(res, city, keyWord)
			//fmt.Println(res)
			for strings.Contains(pagesout, `class="next"`) {
				chromedp.Click(`a.next`, chromedp.ByQuery).Do(ctx)
				res = ""
				chromedp.OuterHTML("#list_con", &res, chromedp.NodeVisible, chromedp.ByID).Do(ctx)
				pagesout = ""
				chromedp.OuterHTML("div.pagesout", &pagesout, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)

				fmt.Println("getCompanyId", city, keyWord, page)
				getCompanyId(res, city, keyWord)
				page++
			}
			return nil
		}),
	}
	fmt.Println("start run chromedp tasks!")
	err := chromedp.Run(ctx,

		tasks,
	)
	if err != nil {
		log.Println(err)
	}
	//log.Println(res)
	return res
}

func getCompanyId(res, city, keyWord string) {
	uidStr := strings.Split(res, "uid=")
	for i, us := range uidStr {

		if i > 0 && strings.Contains(us, keyWord) {

			dealUid(us, city, keyWord)
		}

	}
}

func dealUid(us, city, keyWord string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("dealUid err:", err, us)
		}
	}()
	href := strings.Split(us, `<a href="`)[1]
	href = strings.Split(href, `" `)[0]
	//log.Println(href)

	id := strings.Split(us, `" `)[0]
	id = strings.Split(id, `_`)[0][1:]
	//log.Println(id)
	pkg.ToTxt(id+","+href+","+city+","+keyWord, "/var/58/"+city+"-"+keyWord+".csv")
}

//获取企业信息
func getCompanyInfo(companyId string, ss []string, line string) {
	time.Sleep(1 * time.Second)
	resp := fasthttp.AcquireResponse()
	req := fasthttp.AcquireRequest()
	//req.Header.Set("Proxy-Authorization","Basic UDdXVW5zR2t0M3NxZzE1Zjo1Q0MxOG42c2xIZWpDckVB")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("https://qy.58.com/ent/detail/" + companyId)
	//_, data, err := c.Get(nil,"https://qy.58.com/ent/detail/"+companyId)

	err := c.Do(req, resp)
	if err == nil {
		data := resp.Body()
		//data, _ := ioutil.ReadAll(resp.Body)
		res := string(data)
		log.Println(res)
		if strings.Contains(res, `"msg":"获取信息成功"`) {
			//log.Println(companyName)
			companyShortName, _ := jsonparser.GetString(data, "data", "entName", "aliasName")
			companyFullName, _ := jsonparser.GetString(data, "data", "entDetail", "aliasName")
			industryField, _ := jsonparser.GetString(data, "data", "entDetail", "industryText")
			//financeStage,_ := jsonparser.GetString(data, "baseInfo", "financeStage")
			companySize, _ := jsonparser.GetString(data, "data", "entDetail", "sizeText")
			companyUrl, _ := jsonparser.GetString(data, "data", "entDetail", "url")
			regLocation, _ := jsonparser.GetString(data, "data", "entDetail", "address")
			regCapital, _ := jsonparser.GetString(data, "data", "entDetail", "bussiness", "regCapital")
			legalPersonName, _ := jsonparser.GetString(data, "data", "entDetail", "bussiness", "legalPersonName")
			establishTime, _ := jsonparser.GetString(data, "data", "entDetail", "bussiness", "createTime")
			contact, _ := jsonparser.GetString(data, "data", "entDetail", "contact")
			email, _ := jsonparser.GetString(data, "data", "entDetail", "email")
			typeText, _ := jsonparser.GetString(data, "data", "entDetail", "typeText")

			companyInfo := CompanyInfo{
				CompanyId:        companyId,
				CompanyShortName: strings.ReplaceAll(companyShortName, ",", ";"),
				CompanyFullName:  strings.ReplaceAll(companyFullName, ",", ";"),
				IndustryField:    strings.ReplaceAll(industryField, ",", ";"),
				//FinanceStage:strings.ReplaceAll(financeStage,",",";"),
				CompanySize:     strings.ReplaceAll(companySize, ",", ";"),
				CompanyUrl:      strings.ReplaceAll(companyUrl, ",", ";"),
				RegLocation:     strings.ReplaceAll(regLocation, ",", ";"),
				RegCapital:      strings.ReplaceAll(regCapital, ",", ";"),
				LegalPersonName: strings.ReplaceAll(legalPersonName, ",", ";"),
				Contact:         strings.ReplaceAll(contact, ",", ";"),
				EstablishTime:   strings.ReplaceAll(establishTime, ",", ";"),
				Email:           strings.ReplaceAll(email, ",", ";"),
				TypeText:        strings.ReplaceAll(typeText, ",", ";"),
				//City:city,
				//WorkName:word,
			}

			log.Println(companyInfo)

			if companyInfo.CompanyFullName != "" || companyInfo.CompanyShortName != "" {
				bs, _ := json.Marshal(companyInfo)
				//store.Update(companyId,string(bs))
				store.Update("companyInfo.Id."+companyId, string(bs))
			}

		} else {
			pkg.ToTxt("no-data.txt", line)
			/*c := getCompanyInfoBychromedp(companyId)
			log.Println(c)
			if c.CompanyFullName != "" || c.CompanyShortName != ""{
				bs,_ := json.Marshal(c)
				store.Update(companyId,string(bs))
			}*/

		}
	} else {
		log.Println("http client err:", err)
	}
}
