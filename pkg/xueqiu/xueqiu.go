package xueqiu

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"math/rand"
	"regexp"
	"spider/pkg"
	"spider/pkg/jijin"
	"strconv"
	"strings"
	"time"
)

var stocks = make([]Stock, 0)

func Run() {
	//jijin.DbInit()

	//getDetailHTML() // 净资产收益率,如果有值写入 000.csv 按数字顺序遍历

	//getFFC()  //获取自由现金流 写入 stock-FFC.csv

	//https://xueqiu.com/S/SH688677

	getDHTML() //详细数据,名称，股价，简介等信息 写入stock-detail.csv
	//getDetail("",Stock{code: "123"})

	//caculate()
	//v := DCF(3.3,0.15,0.05,10)
	//fmt.Println(v)
	//caculateDCF()

	t("SZ000403")
}

func caculateDCF() {
	pkg.ReadLine("resource/stock-valuation-10.csv", func(s string) {
		ss := strings.Split(s, `,`)
		valuation, _ := strconv.ParseFloat(ss[10], 64)
		price, _ := strconv.ParseFloat(ss[8], 64)
		v := DCF(valuation, 0.15, 0.05, 10)
		bl := fmt.Sprintf("%.6f", price/v)
		value := s + "," + fmt.Sprint(v) + "," + fmt.Sprint(bl)
		pkg.ToTxt(value, "stock-valuation-10----.csv")
	})
}

func caculate() {
	pkg.ReadLine("resource/stock-result-10.csv", func(s string) {
		ss := strings.Split(s, `,`)
		TTM, _ := strconv.ParseFloat(ss[4], 64)
		Sum, _ := strconv.ParseFloat(ss[5], 64)
		Jzc, _ := strconv.ParseFloat(ss[3], 64)
		_ = Jzc
		Price, _ := strconv.ParseFloat(ss[8], 64)
		_ = Price
		flow := ss[9]
		unit := 0.0
		if !strings.Contains(flow, "-") && TTM > 0 {
			if strings.Contains(flow, "千亿") {
				unit = 100000000 * 1000
				flow = strings.ReplaceAll(flow, "千亿", "")

			} else if strings.Contains(flow, "百亿") {
				unit = 100000000 * 100
				flow = strings.ReplaceAll(flow, "百亿", "")

			} else if strings.Contains(flow, "亿") {
				unit = 100000000
				flow = strings.ReplaceAll(flow, "亿", "")

			} else if strings.Contains(flow, "千万") {
				unit = 10000 * 100
				flow = strings.ReplaceAll(flow, "千万", "")

			} else if strings.Contains(flow, "百万") {
				unit = 10000 * 100
				flow = strings.ReplaceAll(flow, "百万", "")

			} else if strings.Contains(flow, "万") {
				unit = 10000
				flow = strings.ReplaceAll(flow, "万", "")

			}

			Flow, err := strconv.ParseFloat(flow, 64)
			if err != nil {
				fmt.Println(err)
			}
			Flow = Flow * unit
			valuation := (Flow / Sum) //每股自由现金流
			if valuation == 0 {
				fmt.Println(ss[0], Flow, Sum, TTM)
			}
			ffc := DCF(valuation, 0.15, 0.05, 10) //ffc 估值
			ttmNew := Price / valuation
			value := s + "," + fmt.Sprint(valuation) + "," + fmt.Sprint(ffc) + "," + fmt.Sprint(ttmNew) + "," + fmt.Sprint(ttmNew/TTM)
			/*if Jzc > 15 {
				pkg.ToTxt(value,"stock-valuation-15.csv")
			}*/
			/*if Jzc >= 20   {
				pkg.ToTxt(value,"stock-valuation-20.csv")
			}*/
			pkg.ToTxt(value, "stock-valuation-10.csv")

		}

	})
}

//自由现金流永续估值
/**
D0  当前现金
g   企业自由现金流增长率
r   折现率为r  假设现在猪肉价格是100块一斤，一年后105块，那么以猪肉为标的，1年后的105块折算到现在就是100块。以猪肉的价格为标准，我们可以得出一个贴现率=（105-100）/100=5%，即一年期的贴现率为5%，一年后的现金如果折算到现在都要除以（1+5%），如果这个值保持恒定，后一年的现金折算到现金需要除以2次（1+5%），以此类推
n   年数
*/
func DCF(D0 float64, g float64, r float64, n uint64) float64 {
	var p0 float64 = 0
	for i := 1; uint64(i) <= n; i++ {
		G := exponent((1.0 + g), uint64(i))
		R := exponent((1.0 + r), uint64(i))
		f, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", G/R), 64)
		D := D0 * f
		fmt.Println(fmt.Sprintf("增长后比例：%.6f，贴现率：%.6f，第：%d年预测值：%.6f", G, R, i, D))
		p0 = p0 + D
	}
	return p0
}

func exponent(a float64, n uint64) float64 {
	result := float64(1)
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			result *= a
		}
		a *= a
	}
	return result
}

func getFFCDetail(name, res string, ctx context.Context) {
	pkg.ReadLine(name, func(s string) {
		ss := strings.Split(s, `,`)
		code := ss[0]
		/*TTM,_ := strconv.ParseFloat(ss[5],64)
		Sum,_ := strconv.ParseFloat(ss[7],64)
		Price ,_ := strconv.ParseFloat(ss[8],64)
		stock := Stock{code: code,Sum: Sum,Price: Price,Ttm: TTM}*/
		runes := []rune(code)
		if len(runes) < 2 {
			return
		}
		chromedp.Navigate("https://caibaoshuo.com/terms/" + string(runes[2:]) + "/free_cash_flow").Do(ctx)
		chromedp.OuterHTML(`body`, &res, chromedp.ByQuery).Do(ctx)
		fmt.Println(res)
		if strings.Contains(res, `您要访问的页面不存在`) || strings.Contains(res, `服务器内部错误`) {
			return
		}
		chromedp.InnerHTML(`div.scroll-container`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		//fmt.Println(res)
		res = strings.Split(res, `自由现金流(FCF)`)[1]
		ss1 := strings.Split(res, ` </td>`)
		res = strings.ReplaceAll(ss1[len(ss1)-2], `<td>`, ``)
		res = strings.ReplaceAll(res, `
`, ``)
		res = strings.ReplaceAll(res, ` `, ``)
		fmt.Println(res)
		value := code + "," + fmt.Sprint(res)
		pkg.ToTxt(value, "stock-FFC.csv")
		time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
	})
}

//获取自由现金流
func getFFC() {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")
			getFFCDetail("resource/000.csv", res, ctx)
			getFFCDetail("resource/002.csv", res, ctx)
			getFFCDetail("resource/300.csv", res, ctx)
			getFFCDetail("resource/600.csv", res, ctx)
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

//遍历csr,查询详细数据
func getDHTML() {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")
			reCheckByCSR("resource/stock-FFC.csv", res, ctx)
			/*reCheckByCSR("resource/000.csv",res,ctx)
			reCheckByCSR("resource/002.csv",res,ctx)
			reCheckByCSR("resource/300.csv",res,ctx)
			reCheckByCSR("resource/600.csv",res,ctx)*/
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

//雪球主页，获取首页的数据
func reCheckByCSR(name string, res string, ctx context.Context) {
	codeMap := make(map[string]int8)
	pkg.ReadLine("resource/stock-detail.csv", func(s string) {
		if s == "" {
			return
		}
		ss := strings.Split(s, `,`)
		codeMap[ss[0]] = 1
	})
	pkg.ReadLine(name, func(s string) {
		defer func() {
			if err := recover(); err != nil {
				pkg.ToTxt(fmt.Sprint(err), "err.txt")
			}
		}()
		ss := strings.Split(s, `,`)
		code := ss[0]
		_, ok := codeMap[code]
		if ok || strings.Contains(ss[1], `-`) { //自由现金流是负或没有跳过
			return
		}
		/*ffc := 0.0
		if !strings.Contains(ss[1],`<`){
			jzc,_ = strconv.ParseFloat(ss[1],64)
		}
		if jzc < 5 || jzc >= 10{
			return
		}*/
		stock := Stock{code: code, FFC: ss[1]}
		shtml := ""
		chromedp.Navigate("https://xueqiu.com/S/" + code).Do(ctx)
		fmt.Println("code:", code, "has Navigate")
		chromedp.InnerHTML(`div.container-lg`, &shtml, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		chromedp.InnerHTML(`div.quote-container`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		//res = strings.ReplaceAll(strings.ReplaceAll(res,`<span>`,""),`</span>`,"")

		fmt.Println(res)
		stock = getDetail(res, shtml, stock)

		stringV := caculateV(stock.Ttm, stock.Sum, stock.Price, stock.FFC)

		//代码,总市值,市值,静态市盈率,动态市盈率,总股本,当前价格,自由现金流,每股自由现金流,估值,现金流计算市盈率,计算市盈率比值,名称,简介,服务
		value := fmt.Sprint(stock.code) + "," + fmt.Sprint(stock.Zsz) + "," + fmt.Sprint(stock.SZ) + "," +
			//fmt.Sprint(stock.Ltz)+","+fmt.Sprint(stock.Jzcsyl)+","+
			fmt.Sprint(stock.Sylj) + "," + fmt.Sprint(stock.Ttm) + "," +
			//fmt.Sprint(math.Floor(stock.Sylj-stock.Ttm))+","+
			fmt.Sprint(stock.Sum) + "," +
			fmt.Sprint(stock.Price) + "," +
			fmt.Sprint(stock.FFC) + "," + stringV + "," +
			fmt.Sprint(stock.Name) + "," +
			fmt.Sprint(stock.Info) + "," +
			fmt.Sprint(stock.Service)
		pkg.ToTxt(value, "stock-detail.csv")
		time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
	})
}

func getDetail(html, shtml string, s Stock) Stock {
	//fmt.Println(html)
	//html = `<div class="stock-info"><div class="stock-price stock-fall"><div class="stock-current"><strong>¥18.03</strong></div><div class="stock-change">-0.12  -0.66%</div></div><div class="stock-time"><div>&nbsp;53.55 万球友关注</div><div class="quote-market-status"><span>交易中<pan> 09-30 14:12:03 北京时间</span></div></div></div><table class="quote-info"><tbody><tr><td>最高：<span class="stock-rise">18.20</spa<td>今开：<span class="stock-fall">18.09</span></td><td>涨停：<span class="stock-rise">19.97</span></td><td>成交量：<span>66.24万手</sp><tr class="separateTop"><td>最低：<span class="stock-fall">17.71</span></td><td>昨收：<span>18.15</span></td><td>跌停：<span class="st>16.34</span></td><td>成交额：<span>11.85亿</span></td></tr><tr class="separateBottom"><td>量比：<span class="stock-fall">0.66</span></手：<span>0.34%</span></td><td>市盈率(动)：<span>9.95</span></td><td>市盈率(TTM)：<span>10.66</span></td></tr><tr><td>委比：<span class-40.71%</span></td><td>振幅：<span>2.70%</span></td><td>市盈率(静)：<span>12.10</span></td><td>市净率：<span>1.14</span></td></tr><tr><n>1.69</span></td><td>股息(TTM)：<span>0.18</span></td><td>总股本：<span>194.06亿</span></td><td>总市值：<span>3498.89亿</span></td></t资产：<span>15.83</span></td><td>股息率(TTM)：<span>1.00%</span></td><td>流通股：<span>194.06亿</span></td><td>流通值：<span>3498.86亿<tr><td>52周最高：<span>25.16</span></td><td>52周最低：<span>14.64</span></td><td>货币单位：<span>CNY</span></td></tr></tbody></table>`
	zsz := getByTag(html, `总市值：`+`<span>`, "</span>")
	//Zsz,err := strconv.ParseFloat(zsz,64)
	Zsz := stringToMoney(zsz)
	ltz := 0.0 //strings.ReplaceAll(getByTag(html,`流通值：`+`<span>`,""),"亿","")
	ttm := getByTag(html, `市盈率(TTM)：`+`<span>`, "")
	sylj := getByTag(html, `市盈率(静)：`+`<span>`, "")
	sum := getByTag(html, `总股本：<span>`, "")
	price := strings.ReplaceAll(getByTag(html, `<div class="stock-current"><strong>`, "</strong>"), "¥", "")
	re := regexp.MustCompile("[\u4e00-\u9fa5]{1,}")
	service := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<div class="title">业务</div>`, `<!---->`), ",", "，"), -1) //业务
	name := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<div class="stock-name">`, `</div>`), ",", "，"), -1)        //name
	info := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<div class="title">简介</div>`, `<!---->`), ",", "，"), -1)    //简介

	s.Name = strings.Join(name, "，")
	s.Service = strings.Join(service, "，")
	s.Info = strings.Join(info, "，")
	s.SZ = zsz
	s.Zsz = Zsz
	/*Ltz,err := strconv.ParseFloat(ltz,64)
	if err != nil{
		Ltz = 0
	}*/
	s.Ltz = 0.0
	Ttm, err := strconv.ParseFloat(ttm, 64)
	if err != nil {
		Ttm = 10000
	}
	s.Ttm = Ttm
	Sylj, err := strconv.ParseFloat(sylj, 64)
	if err != nil {
		Sylj = 10000
	}
	s.Sylj = Sylj
	c := 1.0
	if strings.Contains(sum, "亿") {
		c = 100000000
		sum = strings.ReplaceAll(sum, "亿", "")
	} else if strings.Contains(sum, "千万") {
		sum = strings.ReplaceAll(sum, "千万", "")
		c = 10000000
	} else if strings.Contains(sum, "百万") {
		sum = strings.ReplaceAll(sum, "百万", "")
		c = 1000000
	} else if strings.Contains(sum, "十万") {
		sum = strings.ReplaceAll(sum, "十万", "")
		c = 100000
	} else if strings.Contains(sum, "万") {
		sum = strings.ReplaceAll(sum, "万", "")
		c = 10000
	}
	Sum, err := strconv.ParseFloat(sum, 64)
	if err != nil {
		Sum = 0.0
	}
	s.Sum = Sum * c
	Price, err := strconv.ParseFloat(price, 64)
	if err != nil {
		Price = 0.0
	}
	s.Price = Price
	//sy := getByTag(html,`市盈`)
	fmt.Println(zsz, ltz, ttm, sylj, s.Sum, s.Price)
	return s
}

type Stock struct {
	Name    string  `名称`
	Service string  `业务`
	Info    string  `简介`
	SZ      string  `市值`
	code    string  `代码`
	Jzcsyl  float64 `净资产收益率`
	Zsz     float64 `总市值`
	Ltz     float64 `流通市值`
	Ttm     float64 `TTM`
	Sylj    float64 `市盈率(静)`
	Sum     float64 `总股本`
	Price   float64 `当前价格`
	FFC     string  `自由现金流`
}

func getByTag(src, tag, end string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	ss := strings.Split(src, tag)[1]
	if end != "" {
		return strings.Split(ss, end)[0]
	}
	return strings.Split(ss, `</span>`)[0]
}

func getListHTML(name string) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")

			var count int32 = 0
			jijin.Db.Table("info").Count(&count)
			var pageIndex int32 = 1
			var pageSize int32 = 100
			for ; pageIndex*pageSize < count; pageIndex++ {
				infos := jijin.SelectAllStock(pageIndex, pageSize)
				for _, stock := range infos {
					if stock.Proportion == "" {
						chromedp.Navigate("https://xueqiu.com/k?q=" + stock.Name).Do(ctx)

						chromedp.InnerHTML(`p.search__stock__bd__code`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)

						res = strings.ReplaceAll(strings.ReplaceAll(res, `<span>`, ""), `</span>`, "")

						fmt.Println(res)

						stock.Num = res

						chromedp.Navigate("https://xueqiu.com/snowman/S/" + res + "/detail#/JJCG").Do(ctx)
						time.Sleep(3 * time.Second)
						chromedp.InnerHTML(`div.container-md.float-left.stock__info__main`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
						if strings.Contains(res, "全部合计") {
							chromedp.InnerHTML(`table.brief-info tbody tr`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
							res = strings.ReplaceAll(res, `\n`, "")
							res = strings.Split(res, `<td>`)[3]
							res = strings.ReplaceAll(res, `</td>`, "")

							fmt.Println(res)

							stock.Proportion = res

							jijin.Db.Save(&stock).Commit()
						} else {
							stock.Proportion = "null"

							jijin.Db.Save(&stock).Commit()
						}
					}
				}
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

}

func getDetailHTML() {
	url := "https://xueqiu.com/snowman/S/%s/detail#/ZYCWZB"
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")
			search300(res, url, ctx)
			search002(res, url, ctx)
			search600(res, url, ctx)
			search000(res, url, ctx)
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

func search300(res, url string, ctx context.Context) {
	for i := 300000; i <= 300999; i++ {
		//SZ 300000  创业板
		//SZ00 2000  中小板
		//SH 600000  601000 603000 沪市A
		//SZ000999 递减  深市A
		cs := []rune(fmt.Sprint(i))
		code := "SZ" + string(cs)
		chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
		chromedp.InnerHTML(`div#app`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Index(res, `净资产收益率`) > 0 {
			s := strings.Split(res, `净资产收益率</td>`)[1]
			s = strings.Split(s, `<span>%</span>`)[0]
			s = strings.ReplaceAll(s, `<td><p>`, "")
			fmt.Println(s)
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f = 0.0
			}
			stock := Stock{code: code, Jzcsyl: f}
			stocks = append(stocks, stock)
			pkg.ToTxt(code+","+s, "300.csv")
		}
		time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)
	}
}

func search002(res, url string, ctx context.Context) {
	for i := 2000; i <= 2999; i++ {
		//SZ 300000  创业板
		//SZ00 2000  中小板
		//SH 600000  601000 603000 沪市A
		//SZ000999 递减  深市A
		cs := []rune(fmt.Sprint(i))
		code := "SZ00" + string(cs)
		chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
		chromedp.InnerHTML(`div#app`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Index(res, `净资产收益率`) > 0 {
			s := strings.Split(res, `净资产收益率</td>`)[1]
			s = strings.Split(s, `<span>%</span>`)[0]
			s = strings.ReplaceAll(s, `<td><p>`, "")
			fmt.Println(s)
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f = 0.0
			}
			stock := Stock{code: code, Jzcsyl: f}
			stocks = append(stocks, stock)
			pkg.ToTxt(code+","+s, "002.csv")
		}
		time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)
	}
}

func search600(res, url string, ctx context.Context) {
	for i := 600000; i <= 603999; i++ {
		//SZ 300000  创业板
		//SZ00 2000  中小板
		//SH 600000  601000 603000 沪市A
		if i >= 602000 && i < 603000 {
			continue
		}
		//SZ000999 递减  深市A
		cs := []rune(fmt.Sprint(i))
		code := "SH" + string(cs)
		chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
		chromedp.InnerHTML(`div#app`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Index(res, `净资产收益率`) > 0 {
			s := strings.Split(res, `净资产收益率</td>`)[1]
			s = strings.Split(s, `<span>%</span>`)[0]
			s = strings.ReplaceAll(s, `<td><p>`, "")
			fmt.Println(s)
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f = 0.0
			}
			stock := Stock{code: code, Jzcsyl: f}
			stocks = append(stocks, stock)
			pkg.ToTxt(code+","+s, "600.csv")
		}
		time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)
	}
}

func search000(res, url string, ctx context.Context) {
	for i := 1000; i <= 1999; i++ {
		//SZ 300000  创业板
		//SZ00 2000  中小板
		//SH 600000  601000 603000 沪市A
		//SZ000999 递减  深市A
		cs := []rune(fmt.Sprint(i))
		code := "SZ000" + string(cs[1:])
		chromedp.Navigate(strings.ReplaceAll(url, "%s", code)).Do(ctx)
		chromedp.InnerHTML(`div#app`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
		if strings.Index(res, `净资产收益率`) > 0 {
			s := strings.Split(res, `净资产收益率</td>`)[1]
			s = strings.Split(s, `<span>%</span>`)[0]
			s = strings.ReplaceAll(s, `<td><p>`, "")
			fmt.Println(s)

			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				f = 0.0
			}
			stock := Stock{code: code, Jzcsyl: f}
			stocks = append(stocks, stock)
			pkg.ToTxt(code+","+s, "000.csv")
		}
		time.Sleep(time.Duration(rand.Int63n(4)) * time.Second)
	}
}

func caculateV(TTM, Sum, Price float64, _Flow string) string {
	//TTM,_ := strconv.ParseFloat(_TTM,64)
	//Sum,_ := strconv.ParseFloat(_SUM,64)
	//Jzc,_ := strconv.ParseFloat(ss[3],64)
	//_ = Jzc
	//Price ,_ := strconv.ParseFloat(_Price,64)
	//_= Price
	flow := _Flow
	unit := 0.0
	if !strings.Contains(flow, "-") {
		if strings.Contains(flow, "千亿") {
			unit = 100000000 * 1000
			flow = strings.ReplaceAll(flow, "千亿", "")

		} else if strings.Contains(flow, "百亿") {
			unit = 100000000 * 100
			flow = strings.ReplaceAll(flow, "百亿", "")

		} else if strings.Contains(flow, "亿") {
			unit = 100000000
			flow = strings.ReplaceAll(flow, "亿", "")

		} else if strings.Contains(flow, "千万") {
			unit = 10000 * 100
			flow = strings.ReplaceAll(flow, "千万", "")

		} else if strings.Contains(flow, "百万") {
			unit = 10000 * 100
			flow = strings.ReplaceAll(flow, "百万", "")

		} else if strings.Contains(flow, "万") {
			unit = 10000
			flow = strings.ReplaceAll(flow, "万", "")

		}

		Flow, err := strconv.ParseFloat(flow, 64)
		if err != nil {
			fmt.Println(err)
		}
		Flow = Flow * unit
		valuation := (Flow / Sum) //每股自由现金流
		/*if valuation ==0 {
			fmt.Println(ss[0],Flow,Sum,TTM)
		}*/
		ffc := DCF(valuation, 0.15, 0.05, 10) //ffc 估值
		ttmNew := Price / valuation
		if TTM <= 0 {
			TTM = 12
		}
		//每股自由现金流,估值,现金流计算市盈率,计算市盈率比值
		value := fmt.Sprint(valuation) + "," + fmt.Sprint(ffc) + "," + fmt.Sprint(ttmNew) + "," + fmt.Sprint(ttmNew/TTM)
		/*if Jzc > 15 {
			pkg.ToTxt(value,"stock-valuation-15.csv")
		}*/
		/*if Jzc >= 20   {
			pkg.ToTxt(value,"stock-valuation-20.csv")
		}*/
		//pkg.ToTxt(value,"stock-valuation-10.csv")
		return value
	}
	return "-,-,-,-"
}

func stringToMoney(flow string) float64 {
	unit := 0.0
	if strings.Contains(flow, "万亿") {
		unit = 1000000000 * 1000
		flow = strings.ReplaceAll(flow, "万亿", "")

	} else if strings.Contains(flow, "千亿") {
		unit = 100000000 * 1000
		flow = strings.ReplaceAll(flow, "千亿", "")

	} else if strings.Contains(flow, "百亿") {
		unit = 100000000 * 100
		flow = strings.ReplaceAll(flow, "百亿", "")

	} else if strings.Contains(flow, "十亿") {
		unit = 100000000
		flow = strings.ReplaceAll(flow, "十亿", "")

	} else if strings.Contains(flow, "亿") {
		unit = 100000000
		flow = strings.ReplaceAll(flow, "亿", "")

	} else if strings.Contains(flow, "千万") {
		unit = 10000 * 100
		flow = strings.ReplaceAll(flow, "千万", "")

	} else if strings.Contains(flow, "百万") {
		unit = 10000 * 100
		flow = strings.ReplaceAll(flow, "百万", "")

	} else if strings.Contains(flow, "万") {
		unit = 10000
		flow = strings.ReplaceAll(flow, "万", "")
	}

	Flow, err := strconv.ParseFloat(flow, 64)
	if err != nil {
		fmt.Println(err)
	}
	Flow = Flow * unit
	return Flow
}

func t(code string) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res, shtml string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")
			chromedp.Navigate("https://xueqiu.com/S/" + code).Do(ctx)
			fmt.Println("code:", code, "has Navigate")
			chromedp.InnerHTML(`div.container-lg`, &shtml, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
			chromedp.InnerHTML(`div.quote-container`, &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
			//res = strings.ReplaceAll(strings.ReplaceAll(res,`<span>`,""),`</span>`,"")

			fmt.Println(res)
			stock := getDetail(res, shtml, Stock{})
			value := fmt.Sprint(stock.code) + "," + fmt.Sprint(stock.Zsz) + "," + fmt.Sprint(stock.SZ) + "," +
				//fmt.Sprint(stock.Ltz)+","+fmt.Sprint(stock.Jzcsyl)+","+
				fmt.Sprint(stock.Sylj) + "," + fmt.Sprint(stock.Ttm) + "," +
				//fmt.Sprint(math.Floor(stock.Sylj-stock.Ttm))+","+
				fmt.Sprint(stock.Sum) + "," +
				fmt.Sprint(stock.Price) + "," +
				fmt.Sprint(stock.FFC) + "," +
				fmt.Sprint(stock.Name) + "," +
				fmt.Sprint(stock.Info) + "," +
				fmt.Sprint(stock.Service)
			fmt.Println(value)
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
