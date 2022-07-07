package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"spider/pkg"
	"strconv"
	"strings"
	"time"
)

var (
	readFile = false
	saveFile = false
	db       *gorm.DB
	start    bool = false
)

func main() {
	var err error
	db, err = gorm.Open("mysql", "root:0l268uPDBPmuMtea8JNfPgZJpMlRQ9w@(bj-cdb-enpzeiqs.sql.tencentcdb.com:61505)/spider?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}
	db.LogMode(true)
	//saveFile = true
	//readFile = true
	//loginQichacha()
	/*
		loginQichacha()*/
	/*readFile = true

	dealHTML(&company,"")*/
	/*var company Company
	var page int8 = 1
	companys,_ := selectAll(&company,page)
	for ;len(companys) > 0;{
		for _,c := range companys{
			mainMembers(&c)
			competitors(&c)
			s := ""
			for _,m := range c.MainMembers{
				s = s + m.Name +":"+m.Position+"\n"
			}
			c.Members = s
			s = ""
			for _,m := range c.Competitors{
				s = s + m.CompatName+"\n"
			}
			c.CompatProducts = s
			c.Competitors = nil
			c.MainMembers = nil
			fmt.Println(c,c.Members,c.CompatProducts,page)
			Update(&c)
		}
		page++
		companys,_ = selectAll(&company,page)
	}*/

	Server()
	//createTable()
	/*var company Company
	dealHTML(&company,"")*/

}

func Server() {
	filepath.Abs(filepath.Dir("."))
	//http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(pkg.HTML))
	})
	http.HandleFunc("/search", func(writer http.ResponseWriter, request *http.Request) {
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println(err)
		}
		url := request.RequestURI
		url = strings.Split(url, "page=")[1]
		var c Company
		json.Unmarshal(bs, &c)
		page, err := strconv.ParseInt(url, 10, 8)
		fmt.Println(page)
		companys, count := selectAll(&c, int8(page))
		writer.Header().Add("count", fmt.Sprint(count))
		bs, _ = json.Marshal(companys)
		writer.Write(bs)
	})
	http.HandleFunc("/mainMembers", func(writer http.ResponseWriter, request *http.Request) {
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(bs))
		var c Company
		json.Unmarshal(bs, &c)
		mainMembers(&c)
		bss, _ := json.Marshal(c.MainMembers)
		writer.Write(bss)
	}) //
	http.HandleFunc("/competitors", func(writer http.ResponseWriter, request *http.Request) {
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(bs))
		var c Company
		json.Unmarshal(bs, &c)
		competitors(&c)
		bss, _ := json.Marshal(c.Competitors)
		writer.Write(bss)
	})
	http.HandleFunc("/startSpider", func(writer http.ResponseWriter, request *http.Request) {
		go loginQichacha()
		fmt.Println("startSpider")
		writer.Write([]byte("ok"))
	})
	http.HandleFunc("/exportAll", func(writer http.ResponseWriter, request *http.Request) {
		exportAll()
		fmt.Println("exportAll")
		writer.Write([]byte("ok"))
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func createTable() {
	//host=rm-2ze81r8ys3tw0s1k8wo.pg.rds.aliyuncs.com port=3432 user=postgres dbname=spider password=postgres123! sslmode=disable
	//root:root@52.83.252.118:3306/dbname?charset=utf8&parseTime=True&loc=Local
	/*db, err := gorm.Open("mysql", "root:root@(52.83.252.118:3306)/spider?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil{
		fmt.Println(err)
	}*/
	db.DropTable("company")
	db.DropTable("main_member")
	db.DropTable("compat_product_info")
	//d := db.CreateTable(Company{})
	//fmt.Println(d.Value,d.Error)
	result := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Table("company").CreateTable(&Company{})
	if result.Error != nil {
		fmt.Println("create table company err:", result.Error)
	}
	result = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Table("main_member").CreateTable(&MainMember{})
	if result.Error != nil {
		fmt.Println("create table main_member err:", result.Error)
	}
	result = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Table("compat_product_info").CreateTable(&CompatProductInfo{})
	if result.Error != nil {
		fmt.Println("create table compat_product_info err:", result.Error)
	}
}

func save(company *Company) {
	if company.CompanyName == "" {
		return
	}
	s := ""
	for _, m := range company.MainMembers {
		s = s + m.Name + ":" + m.Position + "\n"
	}
	company.Members = s
	s = ""
	for _, m := range company.Competitors {
		s = s + m.CompatName + "\n"
	}
	company.CompatProducts = s
	company.Competitors = nil
	company.MainMembers = nil
	result := db.Table("company").Save(company).Commit()
	if result.Error != nil {
		fmt.Println("save err:", result.Error)
	}

}

func Update(company *Company) {
	db.Table("company").
		Where("id = ?", company.ID).
		Update("members", company.Members).
		Update("compat_products", company.CompatProducts)
}

func find(company *Company) {
	db.Where("company_name LIKE ?", "%"+company.CompanyName+"%").Find(&company)
}

func selectAll(company *Company, pageIndex int8) ([]Company, int) {
	return selectAllWithPageSize(company, pageIndex, 30)
}

func selectAllWithPageSize(company *Company, pageIndex int8, pageSize int8) ([]Company, int) {
	c := make([]Company, 0)
	offset := int64((pageIndex - 1)) * int64(pageSize)
	fmt.Println(offset)
	db.Offset(offset).Limit(pageSize).Where("company_name LIKE ?", "%"+company.CompanyName+"%").Order("ID").Find(&c)
	//db.Offset((pageIndex-1) * pageSize).Limit(pageSize).Where("compat_products is  null").Order("ID").Find(&c)
	count := 0
	db.Table("company").Count(&count)
	//fmt.Println(c, count)
	return c, count
}

func exportAll() {
	var c Company
	var pageIndex int8 = 1
	companys, _ := selectAllWithPageSize(&c, pageIndex, 120)
	for len(companys) > 0 {
		offset := int64((pageIndex - 1)) * int64(120)
		writeExcel(companys, offset)
		pageIndex++
		companys, _ = selectAllWithPageSize(&c, pageIndex, 120)

	}
}

func mainMembers(company *Company) {
	db.Where("company_name = ?", company.CompanyName).Find(&company)
	db.Model(&company).Related(&company.MainMembers)
}

func competitors(company *Company) {
	db.Where("company_name = ?", company.CompanyName).Find(&company)
	db.Model(&company).Related(&company.Competitors)
}

func delete() {
	db.Table("company").Delete(&Company{}).Where("id>0")
}

func loginQichacha() {
	if start {
		return
	}
	start = true
	var res string
	_ = res
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
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			network.Enable().Do(ctx)
			chromedp.Navigate("https://www.qcc.com").Do(ctx)
			//btn btn-primary
			Println(time.Now().String())
			//chromedp.Click("div.modal-dialog.nmodal.qcccom-modal button.close", chromedp.ByQuery).Do(ctx)
			//chromedp.OuterHTML("div.container a.navi-btn span", &res, chromedp.ByQuery).Do(ctx)
			//Println(res)
			//chromedp.Click("div.container a.navi-btn span", chromedp.ByQuery).Do(ctx)
			chromedp.Evaluate("showLoginModal();", "").Do(ctx)
			chromedp.Sleep(30 * time.Second).Do(ctx)
			pkg.ReadLine("company.csv", func(s string) {
				var company Company = Company{CompanyName: s}
				find(&company)
				if company.ID <= 0 {
					company = SearchAndGetInfo(s, ctx)
					save(&company)
					fmt.Println(company)
				} else {
					fmt.Println("企业已存在:", s)
				}
			})

			return nil
		}),
	}
	err := chromedp.Run(ctx,
		tasks,
	)
	if err != nil {
		log.Println(err)
	}
}

func SearchAndGetInfo(companyName string, ctx context.Context) Company {
	defer func() {
		if err := recover(); err != nil {
			pkg.ToTxt(companyName+":"+fmt.Sprint(err), "error.log")
		}
	}()
	var res string
	var company Company
	chromedp.Navigate("https://www.qcc.com/search?key=" + companyName).Do(ctx)
	chromedp.OuterHTML("div.adsearch-list", &res, chromedp.ByQuery).Do(ctx)
	//pkg.ToTxt(res,"a.txt")
	//res = pkg.GetSub(res,`'内容名称':'`,`class="ma_h1">`)[0]
	ress := strings.Split(res, "</tr>")
	for _, v := range ress {
		//fmt.Println(v)
		if strings.Contains(res, companyName) && strings.Contains(v, `/firm/`) {
			res = strings.ReplaceAll(v, " ", "")
			//fmt.Println(res)
			res = pkg.GetSub(res, `"href="`, `"target="_blank"`)[0]
			url := res
			//chromedp.Navigate("https://www.qcc.com/" + res + "#base").Do(ctx)
			chromedp.Navigate(url + "#base").Do(ctx)
			fmt.Println(url)
			chromedp.OuterHTML("html", &res, chromedp.ByQuery).Do(ctx)
			if saveFile {
				pkg.ToTxt(res, "base.html")
			}
			//chromedp.Sleep(600 * time.Second).Do(ctx)
			//fmt.Println("res:",res)
			dealHTML(&company, res)

			/*chromedp.Evaluate(`$("#firstcaseModal button").click()`, "").Do(ctx)
			chromedp.Evaluate(`$('[data-option="netProfitOption"]').mouseover()`, "").Do(ctx)
			var financialAnalysisTitle string
			fmt.Println("-------1")
			chromedp.InnerHTML("#financialAnalysisTitle", &financialAnalysisTitle, chromedp.ByID).Do(ctx)
			company.NetInterestRate = Println(financialAnalysisTitle)
			fmt.Println("-------2")
			chromedp.OuterHTML("#Mainmember", &res, chromedp.ByID).Do(ctx)
			if saveFile {
				pkg.ToTxt(res, "Mainmember.html")
			}
			Mainmember(&company, res)
			fmt.Println("-------3")
			//营业收入  净利润
			chromedp.Evaluate(`$("#report_title").click()`, "").Do(ctx)
			chromedp.Evaluate("boxScrollNew('#financingInfo');zhugeTrack('企业主页内容点击',{'点击来源':'快捷入口','数据维度':'企业发展-融资信息'});", "").Do(ctx)

			fmt.Println("-------4")
			chromedp.OuterHTML("#financingInfo", &res, chromedp.ByID).Do(ctx)
			if saveFile {
				pkg.ToTxt(res, "FinancingInfo.html")
			}
			financingInfo(&company, res)
			fmt.Println("-------5")
			chromedp.Evaluate("boxScrollNew('#compatProductInfo');zhugeTrack('企业主页内容点击',{'点击来源':'快捷入口','数据维度':'企业发展-竞品信息'});", "").Do(ctx)
			chromedp.OuterHTML("#compatProductInfo", &res, chromedp.ByID).Do(ctx)
			fmt.Println("-------6")
			if saveFile {
				pkg.ToTxt(res, "CompatProductInfo.html")
			}
			compatProductInfo(&company, res)*/

			company.SearchName = companyName
			return company
		}
	}
	return company
}

type Company struct {
	gorm.Model
	CompanyName string
	SearchName  string
	LegalPerson string //法人
	Number      string
	Email       string
	Url         string
	Address     string
	//品牌名称/公司简称
	CompanyAbbreviation string
	//简介
	Introduction string
	//经营状态
	OperatingStatus string
	//成立日期
	DateOfEstablishment string
	//注册资本
	RegisteredCapital string
	//核准日期
	ApprovalDate string
	//所属行业
	Industry string
	//所属地区
	Area string
	//人员规模
	StaffSize string
	//参保人数
	NumberOfParticipants string
	//总市值
	TotalMarketCapitalization string
	//流通市值
	MarketCapitalization string
	//营业收入
	OperatingIncome string
	//净利润
	NetProfit string
	//营业区间
	BusinessArea string
	//净利率
	NetInterestRate string

	FinancingInfo

	Members string

	MainMembers []MainMember `gorm:"ForeignKey:CompanyID;save_associations:true"`

	CompatProducts string

	Competitors []CompatProductInfo `gorm:"ForeignKey:CompanyID;save_associations:true"`
}

func (c Company) TableName() string {
	return "company"
}

func dealHTML(company *Company, html string) {
	if readFile {
		bs, _ := ioutil.ReadFile("base.html")
		html = string(bs)
	}
	companyName := pkg.GetSub(html, `<h1>`, `</h1>`)[0] //工商全称
	company.CompanyName = Println(companyName)
	/*rongzi := pkg.GetSub(html, `class="ntag text-list click"`, `<span`)[0]
	if strings.Contains(rongzi, ">") {
		rongzi += "end"
		rongzi = pkg.GetSub(rongzi, `>`, `end`)[0]
	}
	//融资伦次
	var financingInfo FinancingInfo
	financingInfo.FinancingRounds = Println(rongzi)
	company.FinancingInfo = financingInfo*/
	//电话
	number := pkg.GetSub(html, `style="color: #000;">`, `</span>`)[0]
	//fmt.Println(number)
	company.Number = number
	//邮箱
	email := pkg.GetSub(html, `发送邮件">`, `</a>`)[0]
	company.Email = email
	//官网
	if strings.Contains(html, `进入官网"`) {
		url := pkg.GetSub(html, `"进入官网">`, `</a>`)[0]
		company.Url = url
	}

	productInfo := pkg.GetSub(html, `产品信息：</span>`, `</span>`)[0]
	productInfo += "end"
	productInfo = pkg.GetSub(productInfo, `class="cvlu">`, "end")[0] //品牌名称/公司简称
	company.CompanyAbbreviation = Println(productInfo)
	address := pkg.GetSub(html, `showMapModal('`, `'`)[0] //地址
	company.Address = Println(address)

	/*info := pkg.GetSub(html, `class="m-t-sm m-b-sm">`, `</pre>`)[0]
	info = strings.ReplaceAll(strings.ReplaceAll(info, "<pre>", ""), `"`, "")
	company.Introduction = Println(info)*/

	/*jingyingzhuangtai := pkg.GetSub(html, `经营状态</td>`, `</td>`)[0]
	jingyingzhuangtai += "end"
	jingyingzhuangtai = pkg.GetSub(jingyingzhuangtai, `>`, "end")[0] //经营状态
	_ = jingyingzhuangtai
	company.OperatingStatus = Println(jingyingzhuangtai)*/
	chengliriqi := pkg.GetSub(html, `成立日期</td>`, `</td>`)[0]
	chengliriqi += "end"
	chengliriqi = strings.ReplaceAll(chengliriqi, "\n", "")
	chengliriqi = pkg.GetSub(chengliriqi, ">", "end")[0] //成立日期
	company.DateOfEstablishment = Println(chengliriqi)
	zhuceziben := pkg.GetSub(strings.ReplaceAll(html, " ", ""), `class="tb">注册资本</td><td>`, `</td>`)[0]
	zhuceziben = strings.ReplaceAll(zhuceziben, "<td>", "") //注册资本
	_ = zhuceziben
	company.RegisteredCapital = Println(zhuceziben)
	hezhunriqi := pkg.GetSub(html, `核准日期</td>`, `</td>`)[0]
	hezhunriqi += "end"
	hezhunriqi = pkg.GetSub(hezhunriqi, `>`, `end`)[0] //核准日期
	company.ApprovalDate = Println(hezhunriqi)
	suoshuhangye := pkg.GetSub(html, `所属行业</td>`, `</td>`)[0]
	suoshuhangye += "end"
	suoshuhangye = pkg.GetSub(suoshuhangye, ">", "end")[0] //所属行业
	company.Industry = Println(suoshuhangye)
	renyuanguimuo := pkg.GetSub(strings.ReplaceAll(strings.ReplaceAll(html, "\n", ""), " ", ""), `人员规模</td>`, `</td>`)[0]
	renyuanguimuo += "end"
	renyuanguimuo = pkg.GetSub(renyuanguimuo, ">", "end")[0] //人员规模
	company.StaffSize = Println(renyuanguimuo)
	canbaorenshu := pkg.GetSub(strings.ReplaceAll(strings.ReplaceAll(html, "\n", ""), " ", ""), `参保人数</td>`, `</td>`)[0]
	canbaorenshu += "end"
	canbaorenshu = pkg.GetSub(canbaorenshu, ">", "end")[0] //参保人数
	company.NumberOfParticipants = Println(canbaorenshu)
	if strings.Contains(html, "总市值：") {
		zongshizhi := pkg.GetSub(html, `总市值：</td>`, `</td>`)[0]
		zongshizhi = strings.ReplaceAll(zongshizhi, `<td class="">`, "")
		//总市值
		company.TotalMarketCapitalization = Println(zongshizhi)
		liutongshizhi := pkg.GetSub(html, `流通市值：</td>`, `</td>`)[0]
		liutongshizhi = strings.ReplaceAll(liutongshizhi, `<td class="">`, "")
		//流通市值
		company.MarketCapitalization = Println(liutongshizhi)
	}

	if strings.Contains(html, "营业总收入") {
		yingyeshouru := pkg.GetSub(strings.Split(html, `营业总收入`)[1], `lev2">营业收入</td>`, `</tr>`)[0]
		yingyeshouru = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(yingyeshouru, `<td>`, ""), "</td>", ""), "-", "")
		Println(yingyeshouru) //营业收入
		company.OperatingIncome = yingyeshouru
		/*jinglirun := pkg.GetSub(strings.Split(html, `减:所得税费用`)[1], `lev1">净利润</td>`, `</tr>`)[0]
		jinglirun = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(jinglirun, `<td>`, ""), "</td>", ""), "-", "")
		//净利润
		company.NetProfit = Println(jinglirun)*/
	}
	area := pkg.GetSub(html, `所属地区</td>`, `</td>`)[0]
	area = strings.ReplaceAll(area, `<td class="">`, "")
	company.Area = area
	persion := pkg.GetSub(html, `class="bname">`, `</h2>`)[0]
	persion = strings.ReplaceAll(persion, `<h2 class="seo font-20">`, "")
	company.LegalPerson = persion
	//yingyequjian := pkg.GetSub(html, `营业区间：`, `</div>`)[0]
	//营业区间
	//company.BusinessArea = Println(yingyequjian)
}

type FinancingInfo struct {
	//融资伦次
	FinancingRounds string
	//融资时间
	FinancingTime string
	//融资轮次
	FinancingRound string
	//融资金额
	FinancingAmount string
	//投资机构
	InvestmentAgency string
}

//融资信息
func financingInfo(company *Company, html string) {
	if readFile {
		bs, _ := ioutil.ReadFile("FinancingInfo.html")
		html = string(bs)
	}
	luns := strings.Split(html, `class="tx"`)
	fmt.Println("融资轮数：", len(luns)-1)
	html = luns[1]
	ss := strings.Split(html, `</td>`)
	i := 0
	company.FinancingInfo.FinancingRound = fmt.Sprint(len(luns) - 1)
	for _, s := range ss {
		if strings.Contains(s, "<td") {
			if i == 0 {
				// 融资时间
				s += "end"
				s = pkg.GetSub(s, ">", "end")[0]
				company.FinancingInfo.FinancingTime = Println(s)
			}
			if i == 2 {
				// 融资轮次
				s += "end"
				s = pkg.GetSub(s, ">", "end")[0]
				company.FinancingInfo.FinancingRounds = Println(s)
			}
			if i == 3 {
				// 融资金额
				s += "end"
				s = pkg.GetSub(s, ">", "end")[0]
				company.FinancingInfo.FinancingAmount = Println(s)
			}
			if i == 4 {
				//fmt.Println(s)
				s = strings.ReplaceAll(strings.ReplaceAll(s, "<div>", ""), "</div>", "")
				if strings.Contains(s, `</a>`) {
					s = pkg.GetSub(s, `<a`, `</a>`)[0]
				}
				// 投资机构
				s += "end"
				s = pkg.GetSub(s, ">", "end")[0]
				company.FinancingInfo.InvestmentAgency = Println(s)
			}
			i++
		}
	}
}

type MainMember struct {
	gorm.Model
	Name string
	//职位
	Position  string
	CompanyID uint
}

func (c MainMember) TableName() string {
	return "main_member"
}

//主要人员
func Mainmember(company *Company, html string) {
	if readFile {
		bs, _ := ioutil.ReadFile("Mainmember.html")
		html = string(bs)
	}
	mainmembers := make([]MainMember, 0)
	html = strings.Split(html, `id="employeeslist"`)[0]
	ss := strings.Split(html, `<h3 class="seo font-14">`)
	for i, s := range ss {
		if i > 0 {
			s = "start" + s
			name := pkg.GetSub(s, "start", "</h3>")[0]
			zhiwei := strings.Split(s, `<td class="text-center">`)[3]
			zhiwei = strings.ReplaceAll(zhiwei, "</td>", "")
			mainmembers = append(mainmembers, MainMember{Name: Println(name), Position: Println(zhiwei)})
		}

	}
	company.MainMembers = mainmembers
}

type CompatProductInfo struct {
	gorm.Model
	CompatName string
	CompanyID  uint
}

func (c CompatProductInfo) TableName() string {
	return "compat_product_info"
}

//竞争企业
func compatProductInfo(company *Company, html string) {
	if readFile {
		bs, _ := ioutil.ReadFile("compatProductInfo.html")
		html = string(bs)
	}
	compats := make([]CompatProductInfo, 0)
	htmls := strings.Split(html, `alt="`)
	for i, s := range htmls {
		if i > 0 {
			s = "start" + s
			s = pkg.GetSub(s, `start`, `">`)[0]

			compats = append(compats, CompatProductInfo{CompatName: Println(s)})
		}

	}
	company.Competitors = compats
}

func Println(ss string) string {
	ss = strings.ReplaceAll(ss, " ", "")
	ss = strings.ReplaceAll(ss, "\n", "")
	fmt.Print(ss + ",")
	return ss
}

func Println1(ss ...string) []string {
	for i, s := range ss {
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "\n", "")
		fmt.Print(s + ",")
		ss[i] = s
	}
	fmt.Println()
	return ss
}

func writeExcel(companys []Company, offset int64) {
	f, err := excelize.OpenFile("company.xlsx")
	if err != nil {
		fmt.Println(err)
		if strings.Contains(fmt.Sprint(err), "no such file or directory") {
			f = excelize.NewFile()
			f.SetCellValue("Sheet1", "A1", "公司名称")
			f.SetCellValue("Sheet1", "B1", "企业法人")
			f.SetCellValue("Sheet1", "C1", "联系电话")
			f.SetCellValue("Sheet1", "D1", "邮箱")
			f.SetCellValue("Sheet1", "E1", "官网地址")
			f.SetCellValue("Sheet1", "F1", "地址")
			f.SetCellValue("Sheet1", "G1", "品牌名称")
			f.SetCellValue("Sheet1", "H1", "简介")
			f.SetCellValue("Sheet1", "I1", "经营状态")
			f.SetCellValue("Sheet1", "J1", "成立日期")
			f.SetCellValue("Sheet1", "K1", "注册资本")
			f.SetCellValue("Sheet1", "L1", "核准日期")
			f.SetCellValue("Sheet1", "M1", "所属行业")
			f.SetCellValue("Sheet1", "N1", "人员规模")
			f.SetCellValue("Sheet1", "O1", "参保人数")
			f.SetCellValue("Sheet1", "P1", "所属地区")
			f.SetCellValue("Sheet1", "Q1", "总市值")
			f.SetCellValue("Sheet1", "R1", "流通市值")
			f.SetCellValue("Sheet1", "S1", "营业收入")
			f.SetCellValue("Sheet1", "T1", "净利润")
			f.SetCellValue("Sheet1", "U1", "营业区间")
			f.SetCellValue("Sheet1", "V1", "净利率")
			f.SetCellValue("Sheet1", "W1", "最近融资伦次")
			f.SetCellValue("Sheet1", "X1", "最近融资时间")
			f.SetCellValue("Sheet1", "Y1", "历史融资轮次")
			f.SetCellValue("Sheet1", "Z1", "融资金额")
			f.SetCellValue("Sheet1", "AA1", "投资机构")
			f.SetCellValue("Sheet1", "AB1", "主要人员")
			f.SetCellValue("Sheet1", "AC1", "竞争企业")
			f.SetCellValue("Sheet1", "AD1", "搜索公司名称")
			f.SaveAs("company.xlsx")
			f, _ = excelize.OpenFile("company.xlsx")
		}
	}
	fmt.Println("offset:", offset)
	for j, company := range companys {
		i := int64(j+2) + offset
		f.SetCellValue("Sheet1", "A"+fmt.Sprint(i), company.CompanyName)
		f.SetCellValue("Sheet1", "B"+fmt.Sprint(i), company.LegalPerson)
		f.SetCellValue("Sheet1", "C"+fmt.Sprint(i), company.Number)
		f.SetCellValue("Sheet1", "D"+fmt.Sprint(i), company.Email)
		f.SetCellValue("Sheet1", "E"+fmt.Sprint(i), company.Url)
		f.SetCellValue("Sheet1", "F"+fmt.Sprint(i), company.Address)
		f.SetCellValue("Sheet1", "G"+fmt.Sprint(i), company.CompanyAbbreviation)
		f.SetCellValue("Sheet1", "H"+fmt.Sprint(i), company.Introduction)
		f.SetCellValue("Sheet1", "I"+fmt.Sprint(i), company.OperatingStatus)
		f.SetCellValue("Sheet1", "J"+fmt.Sprint(i), company.DateOfEstablishment)
		f.SetCellValue("Sheet1", "K"+fmt.Sprint(i), company.RegisteredCapital)
		f.SetCellValue("Sheet1", "L"+fmt.Sprint(i), company.ApprovalDate)
		f.SetCellValue("Sheet1", "M"+fmt.Sprint(i), company.Industry)
		f.SetCellValue("Sheet1", "N"+fmt.Sprint(i), company.StaffSize)
		f.SetCellValue("Sheet1", "O"+fmt.Sprint(i), company.NumberOfParticipants)
		f.SetCellValue("Sheet1", "P"+fmt.Sprint(i), company.Area)
		f.SetCellValue("Sheet1", "Q"+fmt.Sprint(i), company.MarketCapitalization)
		f.SetCellValue("Sheet1", "R"+fmt.Sprint(i), company.TotalMarketCapitalization)
		f.SetCellValue("Sheet1", "S"+fmt.Sprint(i), company.OperatingIncome)
		f.SetCellValue("Sheet1", "T"+fmt.Sprint(i), company.NetProfit)
		f.SetCellValue("Sheet1", "U"+fmt.Sprint(i), company.BusinessArea)
		f.SetCellValue("Sheet1", "V"+fmt.Sprint(i), company.NetInterestRate)
		f.SetCellValue("Sheet1", "W"+fmt.Sprint(i), company.FinancingRounds)
		f.SetCellValue("Sheet1", "X"+fmt.Sprint(i), company.FinancingTime)
		f.SetCellValue("Sheet1", "Y"+fmt.Sprint(i), company.FinancingRound)
		f.SetCellValue("Sheet1", "Z"+fmt.Sprint(i), company.FinancingAmount)
		f.SetCellValue("Sheet1", "AA"+fmt.Sprint(i), company.InvestmentAgency)
		f.SetCellValue("Sheet1", "AB"+fmt.Sprint(i), company.Members)
		f.SetCellValue("Sheet1", "AC"+fmt.Sprint(i), company.CompatProducts)
		f.SetCellValue("Sheet1", "AD"+fmt.Sprint(i), company.SearchName)
	}
	f.Save()
}
