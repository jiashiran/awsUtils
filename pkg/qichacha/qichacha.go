package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"math/rand"
	"os"
	"spider/pkg"
	"spider/pkg/qichacha/check"
	sshd "spider/pkg/ssh"
	"spider/pkg/store"
	"strconv"
	"strings"
	"time"
)

var (
	//cookieMap = make(map[string]string)
	//companys = make(map[string]string)
	//"39.107.243.195",
	//,"123.56.17.69"  "39.105.202.8", "47.95.241.44",,,"101.200.52.2"
	//,  "47.93.218.38" "39.96.50.197"
	ips = []string{"47.94.130.119"}
)

func init() {
	/*pkg.ReadLine("拉勾客服外呼北京南京汇总结果.csv", func(s string) {
		defer func() {
			if err := recover();err!=nil{
				log.Println(err,s)
			}
		}()
		ss := strings.Split(s,",",)
		companys[ss[1]] = ss[2]
	})*/

}

func buildCookie() {

	log.Println(store.Redis.Set("cookie_sid", "1572228760854", 0))
	store.Redis.Set("cookie_CNZZDATA1254842228", "1950090108-1566797068-https%3A%2F%2Fwww.baidu.com%2F|1572225867", 0)
	store.Redis.Set("cookie_Hm_lvt_3456bee468c83cc63fb5147f119f1075", "1569745221,1571387179;", 0)

	store.Redis.Set("cookie_info", "1572228760858", 0)
	store.Redis.Set("cookie_cuid", "59fcd2edb004751d734b871932fa908f", 0)
	store.Redis.Set("cookie_UM_distinctid", "16ccca11b11b0-0284524a81e822-38607701-13c680-16ccca11b12b6e", 0)
	store.Redis.Set("cookie_QCCSESSID", "aufoq9fleh2m0nt2njg61ejvq5", 0)
	store.Redis.Set("cookie_acw_tc", "3cdfd94715697452191665001ea7aaa79a944b06c1c3c4c296dcb38642", 0)
	store.Redis.Set("cookie_did", "16ccca116221a2-0f050a5165c1-38607701-13c680-16ccca116235f5", 0)
}

func BatchUpdateImage() {
	sshd.InitServer(ips)
	sshd.ExecuteBatch(`docker pull registry.cn-beijing.aliyuncs.com/-dev/chromedp`)
	m := sshd.ExecuteBatch(`docker images | grep '<none>' | awk -F " " '{print $3}'`)
	for _, ip := range ips {
		if m[ip] != "" {
			sshd.Execute(ip, `docker rmi `+m[ip])
		}
	}
}

func BatchCleanImage() {
	sshd.InitServer(ips)
	m := sshd.ExecuteBatch(`docker images | grep '<none>' | awk -F " " '{print $3}'`)
	for _, ip := range ips {
		if m[ip] != "" {
			sshd.Execute(ip, `docker rmi `+m[ip])
		}
	}
}

func BatchRun() {
	//"39.105.202.8","39.107.243.195",
	//ips := []string{"47.94.130.119","47.93.218.38"}
	sshd.InitServer(ips)
	imageIds := sshd.ExecuteBatch(`docker images | grep chromedp | awk -F ' ' '{print $3}'`)

	cIds := sshd.ExecuteBatch(`docker ps -a | grep chromedp | awk -F ' ' '{print $1}'`)
	log.Println(cIds)
	for _, ip := range ips {
		if cIds[ip] != "" {
			sshd.Execute(ip, `docker start `+cIds[ip])
		} else {
			sshd.Execute(ip, `docker run -d --name chromedp `+imageIds[ip])
		}
	}
	sshd.ExecuteBatch(`docker ps | grep chromedp`)
}

func BatchStop() {
	//"39.105.202.8","39.107.243.195",
	//ips := []string{"47.94.130.119","47.93.218.38"}
	sshd.InitServer(ips)
	cIds := sshd.ExecuteBatch(`docker ps | grep chromedp | awk -F ' ' '{print $1}'`)
	for _, ip := range ips {
		log.Println(ip, cIds[ip])
		sshd.Execute(ip, `docker stop `+cIds[ip])
		sshd.Execute(ip, `docker rm `+cIds[ip])
	}
}

func BatchStatus() {
	//"39.105.202.8","39.107.243.195",

	sshd.InitServer(ips)
	sshd.ExecuteBatch(`docker ps | grep chromedp`)
}

type Zg_de1d1 struct {
	Sid            int64    `json:"sid"`
	Updated        int64    `json:"updated"`
	Info           int64    `json:"info"`
	SuperProperty  struct{} `json:"superProperty"`
	Platform       struct{} `json:"platform"`
	Utm            struct{} `json:"utm"`
	ReferrerDomain string   `json:"referrerDomain"`
	Cuid           string   `json:"cuid"`
}

func main() {
	///buildCookie()
	//BatchStop()
	//BatchUpdateImage()
	//BatchCleanImage()

	///BatchRun()

	//BatchStatus()

	///GrabData()

	//check.BatchDownload()
	//check.BatchClean()
	//CheckCookieStatus()
	//fmt.Println(qichacha("jd"))
	//qichachaRongzi()
}

func CheckCookieStatus() {
	var total int64 = 0
	var count int64 = 0
	for {
		m := check.BatchCheckCount()
		for _, v := range m {
			v = strings.ReplaceAll(v, "\n", "")
			c, _ := strconv.ParseInt(v, 10, 64)
			//log.Println(err,c)
			count = count + c
		}
		log.Println(count)
		if total == count {
			break
		} else {
			total = count
			count = 0
			time.Sleep(15 * time.Second)
		}
	}

}

func GrabData() {
	pkg.ReadLine("/usr/local/data.csv", func(s string) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err, s)
			}
		}()
		ss := strings.Split(s, ",")
		if s == "" {
			return
		}
		key := store.Get("company.Id." + ss[0])
		if key != "" {
			return
		}
		if !store.Setnx("lock.company.Id."+ss[0], 10*time.Second) {
			return
		}
		if ss[4] != "" {
			//log.Println(ss[3])
			//companys[ss[1]] = ss[3]
			info := GetInfo(ss[4])
			if info == "err" {
				return
			}
			store.Redis.Set("company.Id."+ss[0], 1, 0)
			//pkg.ToTxt(s+","+info,"/var/tmp/result.csv")
			store.Push(s + "," + info)
		} else {
			log.Println(s)
			//pkg.ToTxt(s+",","/var/result.csv")
			store.Push(s)
		}
		/*else if ss[2] != ""{
			//log.Println(ss[2])
			info := GetInfo(ss[2])
			//companys[ss[1]] = ss[2]
			store.Redis.Set("company.Id."+ss[1],1,0)
			//pkg.ToTxt(s+","+info,"/var/result.csv")
			store.Push(s+","+info)
		}*/
	})
}

func GetInfo(companyName string) string {
	str, leader := qichacha(companyName)
	defer func() {
		if err := recover(); err != nil {
			log.Println(companyName, err, str)
		}
	}()

	ss := strings.Split(str, "更多号码")
	s := pkg.GetSub(ss[0], "邮箱：", `;">`)[0]
	s = strings.ReplaceAll(
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s,
			" ", ""),
			"	", ""),
			"&quot;", ""),
		`</span><aclass="text-primary"onclick="showHisTel([{t:`, "")
	s = strings.ReplaceAll(s, `<aonclick="showHisEmail([{e:`, "")
	s = strings.ReplaceAll(s, ",", ";")
	if strings.Contains(s, `地址：`) {
		s = strings.Split(s, `地址：`)[0]
	}
	s = strings.ReplaceAll(s, "\n", "")
	if str == "err" {
		return "err"
	}
	return leader + "," + s
}

func qichacha(word string) (string, string) {
	random := rand.Int63n(30)
	log.Println(random)
	time.Sleep(time.Duration(random) * time.Second)
	var res string
	var leader string
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, allocCtxCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer allocCtxCancel()
	defer cancel()
	headers := make(map[string]interface{})
	/*microSecond := time.Now().UnixNano() / 1e6
	second := microSecond / 1000
	zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f := `{"sid": ` + store.Get("cookie_sid") + `,"updated": ` + fmt.Sprint(microSecond) + `,"info": ` + store.Get("cookie_info") + `,"superProperty": "{}","platform": "{}","utm": "{}","referrerDomain": "www.qichacha.com","zs": 0,"sc": 0,"cuid": "` + store.Get("cookie_cuid") + `"}`
	v := url.Values{}
	v.Add("zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f", zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f)
	v.Add("UM_distinctid", store.Get("cookie_UM_distinctid"))
	v.Add("QCCSESSID", store.Get("cookie_QCCSESSID"))
	v.Add("acw_tc", store.Get("cookie_acw_tc"))
	v.Add("CNZZDATA1254842228", store.Get("cookie_CNZZDATA1254842228"))
	v.Add("Hm_lvt_3456bee468c83cc63fb5147f119f1075", store.Get("cookie_Hm_lvt_3456bee468c83cc63fb5147f119f1075"))
	v.Add("Hm_lpvt_3456bee468c83cc63fb5147f119f1075", fmt.Sprint(second))
	if store.Get("cookie_did") != "" {
		v.Add("zg_did", `{"did": "`+store.Get("cookie_did")+`"}`)
	}
	body := v.Encode()
	_ = body
	headers["cookie"] = strings.ReplaceAll(body, "&", "; ")*/
	headers["cookie"] = `QCCSESSID=8h0es6o60ch66ndaqnb6lvc732; UM_distinctid=1709f52211f2d0-0b4e77aa81615d-396a7407-13c680-1709f522120a92; acw_tc=6548fe3815832209249725333e6a7c3e59d22e36adead75ececb189374; zg_did=%7B%22did%22%3A%20%221709f52237d223-0280334d3f7084-396a7407-13c680-1709f52237e6ce%22%7D; CNZZDATA1254842228=109524168-1583218884-https%253A%252F%252Fsp0.baidu.com%252F%7C1583995311; hasShow=1; Hm_lvt_3456bee468c83cc63fb5147f119f1075=1583220925,1583749295,1583827520,1583997159; Hm_lpvt_3456bee468c83cc63fb5147f119f1075=1583997985; zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f=%7B%22sid%22%3A%201583997158608%2C%22updated%22%3A%201583997985793%2C%22info%22%3A%201583827519319%2C%22superProperty%22%3A%20%22%7B%7D%22%2C%22platform%22%3A%20%22%7B%7D%22%2C%22utm%22%3A%20%22%7B%7D%22%2C%22referrerDomain%22%3A%20%22www.baidu.com%22%2C%22cuid%22%3A%20%2259fcd2edb004751d734b871932fa908f%22%2C%22zs%22%3A%200%2C%22sc%22%3A%200%7D`
	log.Println(headers["cookie"])
	/*if len(cookieMap) == 0{

	}else {
		cookie := "acw_tc="+cookieMap["acw_tc"]+"; zg_did="+cookieMap["zg_did"]+"; UM_distinctid="+cookieMap["UM_distinctid"] +
			"; QCCSESSID="+cookieMap["QCCSESSID"]+"; hasShow=1; CNZZDATA1254842228="+cookieMap["CNZZDATA1254842228"] +
			"; Hm_lvt_3456bee468c83cc63fb5147f119f1075=" + "1568202279,1568254024,1568614161,1568628815" +
			"; zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f=" + cookieMap["zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f"] + "; Hm_lpvt_3456bee468c83cc63fb5147f119f1075="+ fmt.Sprint(time.Now().Unix())


		headers["cookie"] = cookie
		log.Println(cookie)
	}*/
	tasks := chromedp.Tasks{
		// read network values
		chromedp.ActionFunc(func(ctx context.Context) error {
			network.Enable().Do(ctx)

			network.SetExtraHTTPHeaders(network.Headers(headers)).Do(ctx)

			chromedp.Navigate("https://www.qichacha.com/search?key=" + word).Do(ctx)
			log.Println("Navigate:", word)

			chromedp.OuterHTML("div#V3_SL.container", &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)

			if strings.Contains(res, "小查还没找到数据") {
				log.Println(res)
				return nil
			}

			chromedp.OuterHTML("#search-result", &res, chromedp.NodeVisible, chromedp.ByID).Do(ctx)
			log.Println("OuterHTML")

			if strings.Contains(res, "法定代表人") && strings.Contains(res, `'法定代表人'`) {
				chromedp.InnerHTML("p.m-t-xs a.text-primary", &leader, chromedp.ByQuery).Do(ctx)
				log.Println(leader)
				if strings.Contains(leader, "展示我的企业") {
					return nil
				}
			} else {
				leader = " - "
			}

			//time.Sleep(10*time.Second)
			/*chromedp.InnerHTML("p.m-t-xs a.text-primary", &res , chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
			log.Println(res)
			res = ""*/
			/*cookies, err := network.GetAllCookies().Do(ctx)
			log.Println("GetAllCookies")
			if err != nil {
				return err
			}

			for i, cookie := range cookies {
				log.Printf("chrome cookie %d: %+v", i, cookie)
				cookieMap[cookie.Name] = cookie.Value
			}*/

			return nil
		}),
	}
	err := chromedp.Run(ctx,

		tasks,
	)
	if err != nil {
		log.Println(err)
	}

	if strings.Contains(res, `更多号码','成为VIP会员 即可挖掘企业更多联系方式`) || strings.Contains(res, "登录后可以查看更多数据") {
		log.Println(res)
		os.Exit(0)
	}
	//return "err",leader
	return res, leader
}

func qichachaRongzi() (string, string) {
	var res string
	var leader string
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, allocCtxCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer allocCtxCancel()
	defer cancel()
	headers := make(map[string]interface{})
	cookieStr := `acw_tc=1bd3c51615832209299923546e1ef8130e59529e3109a1153a0f4e5105; QCCSESSID=g35dnnjlkjibsd9v03ouf0b5j2; zg_did=%7B%22did%22%3A%20%22170d1b84718255-00fe511daa21cd-396b7407-13c680-170d1b8471981c%22%7D; UM_distinctid=170d1b84a7cc97-05d4ee747dea5e-396b7407-13c680-170d1b84a7dd00; hasShow=1; Hm_lvt_78f134d5a9ac3f92524914d0247e70cb=1584066481; _uab_collina=158406648095595909388833; CNZZDATA1254842228=1043190830-1584065522-https%253A%252F%252Fsp0.baidu.com%252F%7C1584092522; zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f=%7B%22sid%22%3A%201584095373875%2C%22updated%22%3A%201584095908798%2C%22info%22%3A%201584066479912%2C%22superProperty%22%3A%20%22%7B%7D%22%2C%22platform%22%3A%20%22%7B%7D%22%2C%22utm%22%3A%20%22%7B%7D%22%2C%22referrerDomain%22%3A%20%22%22%2C%22cuid%22%3A%20%2259fcd2edb004751d734b871932fa908f%22%2C%22zs%22%3A%200%2C%22sc%22%3A%200%7D; Hm_lpvt_78f134d5a9ac3f92524914d0247e70cb=1584095909`
	headers["cookie"] = cookieStr
	log.Println(headers["cookie"])

	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			network.Enable().Do(ctx)
			network.SetExtraHTTPHeaders(network.Headers(headers)).Do(ctx)
			ss := strings.Split(cookieStr, ";")
			for _, kv := range ss {
				sss := strings.Split(kv, "=")
				fmt.Println(sss[0], sss[1])
				c := network.SetCookie(sss[0], sss[1])
				c.Domain = ".qcc.com"
				c.Do(ctx)
			}
			//https://www.qichacha.com/firm_973decc3b2520e594fe74cf15a22ecd8.html#ipo
			//https://www.qcc.com/firm_973decc3b2520e594fe74cf15a22ecd8.html#report
			chromedp.Navigate("https://www.qcc.com/firm_973decc3b2520e594fe74cf15a22ecd8.html#ipo").Do(ctx)
			//btn btn-primary
			chromedp.Click("button.btn.btn-primary", chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
			//ntag text-list click
			//chromedp.InnerHTML("span.ntag.text-list.click", &res, chromedp.ByQuery).Do(ctx)
			//log.Println("span:",res)
			//chromedp.InnerHTML("#realtime", &res, chromedp.ByID).Do(ctx)
			//chromedp.Click("a.btn.btn-sm.btn-default.m-r-sm",chromedp.NodeVisible,chromedp.ByQuery).Do(ctx)
			chromedp.InnerHTML("html", &res, chromedp.ByQuery).Do(ctx)
			//log.Println("OuterHTML:",res)
			pkg.ToTxt(res, "report.html")
			return nil
		}),
	}
	err := chromedp.Run(ctx,
		tasks,
	)
	if err != nil {
		log.Println(err)
	}
	return res, leader
}
