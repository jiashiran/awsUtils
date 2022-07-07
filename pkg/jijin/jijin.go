package jijin

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
	"log"
	"spider/pkg"
	"strings"
	"time"
)

var (
	codeMap = make(map[string]string)
)

func Run() {
	j_type := "jiegou"
	j_type = "gupiao" //
	j_type = "duichong"
	j_type = "hunhe" //
	j_type = "zhishu"
	j_type = "isUpdate"
	_ = j_type
	DbInit()
	defer Db.Close()
	//createTable()
	//bs, _ := ioutil.ReadFile("pkg/jijin/"+j_type+"_code.txt")

	//getListHTML(j_type) //抓取类型code
	//saveResult(j_type)

	getInfo(j_type)

	//analyses()
	//export()
}

func getListHTML(j_type string) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	//ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()
	//time.Sleep(5 * time.Second)
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			//network.Enable().Do(ctx)
			fmt.Println("step1")
			chromedp.Navigate("https://www.howbuy.com/board/?isUpdate=1").Do(ctx)
			//chromedp.Navigate("https://www.howbuy.com/fund/fundranking/"+j_type+".htm").Do(ctx)
			chromedp.Evaluate(`$("div.dbNew-main-button").click()`, &res).Do(ctx)
			//chromedp.OuterHTML("#tableMain", &res, chromedp.NodeVisible, chromedp.ByID).Do(ctx)
			//dealHtml(res)
			i := 300
			for j := 1; j <= 100; j++ {
				i = i + 2000
				chromedp.Evaluate(`document.scrollingElement.scrollTop=`+fmt.Sprint(i), &res).Do(ctx)
				time.Sleep(4 * time.Second)
				//chromedp.OuterHTML("#tableMain", &res, chromedp.NodeVisible, chromedp.ByID).Do(ctx)
				//
			}
			chromedp.InnerHTML(`div.result_list`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
			pkg.ToTxt(res, "pkg/jijin/"+j_type+"_code.txt")
			dealHtml(res)
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

}

func dealHtml(html string) {
	html = strings.Split(html, `<tbody>`)[1]
	ss := strings.Split(html, `href="https://www.howbuy.com/fund`)

	for _, v := range ss {
		if strings.Index(v, `/`) == 0 {
			code := getSub(v, `/`, `">`)[0]
			if codeMap[code] != "" {
				fmt.Println("已存在：" + code)
			}
			codeMap[code] = code

		}
	}

}

func saveResult(j_type string) {
	i := 0
	for k, _ := range codeMap {
		i++
		//fmt.Println(i,k)
		pkg.ToTxt(k, "pkg/jijin/"+j_type+"_a_code.txt")
	}
}

func getSub(res, start, end string) []string {
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

//遍历文件逐行查询
func getInfo(j_type string) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	//ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()
	//time.Sleep(5 * time.Second)
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")
			//pkg.ReadLine("pkg/jijin/result.txt", func(s string) {
			pkg.ReadLine("pkg/jijin/"+j_type+"_a_code.txt", func(s string) {
				var info = Info{Code: s}
				Find(&info)
				if info.ID <= 0 || info.Amount == "0亿" {
					chromedp.Navigate("https://www.howbuy.com/fund/" + s).Do(ctx)
					chromedp.InnerHTML(`div.file_t_left`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
					if strings.Contains(res, `重仓持股`) {
						chromedp.InnerHTML(`div.gmfund_num`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
						if !strings.Contains(res, `最新规模<span>--</span>`) && !strings.Contains(res, `最新规模<span>0.00亿</span>`) {
							size := fundSize(res)
							info.Amount = size

							chromedp.InnerHTML(`#content`, &res, chromedp.ByID).Do(ctx)
							details := mainStock(res, s)
							info.Details = details
							Db.Save(&info).Commit()
							time.Sleep(5 * time.Second)
						} else {
							info.Amount = "0亿"
							Db.Save(&info).Commit()
							time.Sleep(5 * time.Second)
						}
					} else {
						info.Amount = "0亿"
						Db.Save(&info).Commit()
						time.Sleep(5 * time.Second)
					}
				} else {

				}

			})
			/*s := "009954"
			chromedp.Navigate("https://www.howbuy.com/fund/"+s).Do(ctx)

			var info = Info{Code: s}
			find(&info)
			if info.ID <= 0 {
				chromedp.InnerHTML(`div.gmfund_num`,&res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
				if !strings.Contains(res,`最新规模<span>--</span>`){
					size := fundSize(res)
					info.Amount = size

					chromedp.InnerHTML(`#content`,&res, chromedp.ByID).Do(ctx)
					details := mainStock(res,s)
					info.Details = details
					Db.Save(&info).Commit()
					time.Sleep(3 * time.Second)
				}
			}
			*/
			//pkg.ToTxt(res,"pkg/jijin/b.txt")
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
}

func mainStock(html string, code string) []Detail {
	ss := strings.Split(html, `<div>`)
	details := make([]Detail, 0)
	for _, v := range ss {
		if strings.Contains(v, `<td class="tdr">`) {
			var d Detail
			d.Code = code
			v = "-start-" + v
			//fmt.Println(v)
			gupiao := pkg.GetSub(v, "-start-", `</div>`)[0]
			fmt.Println(gupiao)
			d.Name = gupiao
			zhanbi := strings.Split(v, `<td class="tdr">`)[2]
			zhanbi = strings.ReplaceAll(zhanbi, `</td>`, "")
			fmt.Println(zhanbi)
			d.Proportion = zhanbi
			details = append(details, d)
		}
	}
	return details
}

func fundSize(html string) string {
	s := strings.Split(html, `<li>最新规模<span>`)[1]
	s = strings.Split(s, `</span></li>`)[0]
	fmt.Println(s)
	return s
}

func analyses() {
	var count int32 = 0
	Db.Table("info").Count(&count)
	m := make(map[string]decimal.Decimal)
	var pageIndex int32 = 1
	var pageSize int32 = 100
	for ; pageIndex*pageSize < count; pageIndex++ {
		infos := SelectAll(pageIndex, pageSize)
		for _, info := range infos {
			ds := make([]Detail, 0)
			Db.Table("detail").Where("code = ?", info.Code).Find(&ds)
			info.Details = ds
			price, err := decimal.NewFromString(strings.ReplaceAll(info.Amount, "亿", ""))
			if err != nil {
				panic(err)
			}
			for _, d := range info.Details {
				p := strings.ReplaceAll(d.Proportion, " ", "")
				priced, _ := decimal.NewFromString(strings.ReplaceAll(p, "%", ""))
				priced = priced.Div(decimal.NewFromFloat(100))
				fmt.Print(info.Amount, price.Mul(priced), ";")
				if m[d.Name].IsZero() {
					m[d.Name] = price.Mul(priced)
				} else {
					m[d.Name] = m[d.Name].Add(price.Mul(priced))
				}
			}
		}
	}

	for k, v := range m {
		var s Stock
		s.Name = k
		f, _ := v.Float64()

		s.Value = f
		Db.Table("stock").Save(&s)
	}

}

func export() {
	var count int32 = 0
	Db.Table("stock").Count(&count)
	var pageIndex int32 = 1
	var pageSize int32 = 100

	for ; (pageIndex-1)*pageSize < count; pageIndex++ {
		c := make([]Stock, 0)
		offset := int64((pageIndex - 1)) * int64(pageSize)
		Db.Offset(offset).Limit(pageSize).Order("value DESC").Find(&c)
		for _, stock := range c {
			pkg.ToTxt(fmt.Sprint(stock.ID)+";"+stock.Name+";"+fmt.Sprint(stock.Value), "pkg/jijin/result.txt")
		}
	}

}
