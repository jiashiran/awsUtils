package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"spider/pkg"
	"strings"
	"time"
)

var (
	ctx, cancel = chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
)

func main() {
	defer cancel()
	citys := []string{
		"北京", "上海", "广州", "深圳", "天津", "杭州", "成都", "南京", "西安", "苏州", "石家庄",
		"重庆", "郑州", "济南", "东莞", "青岛", "宁波", "无锡", "佛山", "厦门", "长沙", "大连",
		"合肥", "武汉", "开封", "南宁", "绵阳", "沈阳", "长春", "哈尔滨", "济南", "淄博",
	}
	webcall := []string{"电话催收专员", "电话销售", "电话销售专员", "电话销售代表", "呼叫中心电销", "电销专员"}
	customerervice := []string{"电话客服", "客服代表", "客服专员", "呼叫中心客服专员", "客服总监"}
	for _, city := range citys {
		for _, w := range webcall {
			GetHTML(city, w)
		}
		for _, c := range customerervice {
			GetHTML(city, c)
		}
	}

}

func GetHTML(city, keyWord string) {
	var res string
	tasks := chromedp.Tasks{
		chromedp.Navigate("https://www.lagou.com/jobs/list_" + keyWord + "?px=default&city=" + city + "#filterBox"),

		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.OuterHTML(`#s_position_list`, &res, chromedp.NodeVisible, chromedp.ByID).Do(ctx)
			//log.Println(strings.TrimSpace(res))
			contuine := dealHtml(strings.TrimSpace(res), city, keyWord)
			if !strings.Contains(res, "pager_next") { //没有下一页
				contuine = false
			}
			//contuine := true
			i := 1
			for contuine {
				log.Println("pager_next", i, city, keyWord)
				chromedp.Click(`span.pager_next`, chromedp.ByQuery).Do(ctx)
				log.Println("sleep", i)
				chromedp.Sleep(15 * time.Second).Do(ctx)
				var resTemp string
				log.Println("position_list", i)
				chromedp.OuterHTML(`#s_position_list`, &resTemp, chromedp.NodeVisible, chromedp.ByID).Do(ctx)
				log.Println(i)
				i++
				if res != resTemp {
					res = resTemp
					dealHtml(strings.TrimSpace(res), city, keyWord)
				} else {
					contuine = false
				}
				if !strings.Contains(resTemp, "pager_next") { //没有下一页
					contuine = false
				}

			}

			return nil
		}),
	}

	err := chromedp.Run(ctx,
		//chromedp.Navigate("https://www.lagou.com/jobs/list_"+keyWord+"?px=default&city="+city+"#filterBox"),
		//chromedp.Text(`#s_position_list`, &res, chromedp.NodeVisible, chromedp.ByID),
		/*chromedp.Click(`span.pager_next`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Sleep(3 * time.Second),
		chromedp.Click(`span.pager_next`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Sleep(3 * time.Second),*/
		//chromedp.OuterHTML(`#s_position_list`, &res, chromedp.NodeVisible, chromedp.ByID),
		tasks,
	)
	if err != nil {
		log.Fatal(err)
	}

}

func dealHtml(html string, city, keyWord string) bool {
	log.Println("dealHtml")
	lines := strings.Split(strings.TrimSpace(html), "\n")
	contuine := false
	for _, line := range lines {
		if strings.Contains(line, "data-companyid") {
			companyId := getCompanyId(line)
			pkg.ToTxt(companyId+","+city+","+keyWord, city+"-"+keyWord+"-companyId.csv")
			/*companyInfo := getCompanyInfo(companyId)
			companyInfo.City = city
			companyInfo.WorkName = keyWord
			companyInfo.LagouURL = "https://www.lagou.com/gongsi/"+companyId+".html"
			log.Println(companyInfo)
			toTxt(companyInfo.CompanyId+","+companyInfo.CompanyShortName+","+companyInfo.CompanyFullName+","+companyInfo.IndustryField+","+companyInfo.FinanceStage+","+companyInfo.CompanySize+","+companyInfo.CompanyUrl+","+companyInfo.RegLocation+","+companyInfo.RegCapital+","+companyInfo.LegalPersonName+","+companyInfo.EstablishTime+","+companyInfo.WorkName+","+companyInfo.City+","+companyInfo.LagouURL,
				"深圳-客服类.csv")*/
			contuine = true
		}
	}
	if !contuine {
		log.Println("!contuine:", lines)
	}
	return contuine
}

func getCompanyId(line string) string {
	ss := line
	defer func() {
		if err := recover(); err != nil {
			log.Println(ss)
		}
	}()
	start := strings.Index(line, "data-companyid")
	line = line[start:]
	//log.Println(line)
	start = strings.Index(line, " ")
	line = line[16 : start-1]
	return line
}
