package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

func main() {
	/*companyId := getCompanyId("字节跳动")
	if "" != companyId{
		getCompanyInfo(companyId)
	}*/

	login("17710740409", "Cloud1010")
}

func getCompanyId(companyName string) string {
	var res string
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.Navigate("https://www.itjuzi.com/search?data=" + companyName).Do(ctx)
			chromedp.InnerHTML("div.container.clearfix", &res, chromedp.ByQuery).Do(ctx)
			res = strings.Split(res, "条结果")[1]
			if strings.Contains(res, companyName) {
				res = strings.Split(res, companyName)[0]
				res = strings.Split(res, "</span></div")[0]
				res = strings.Split(res, `href="/company/`)[1]
				res = strings.Split(res, `"`)[0]
			} else {
				res = ""
			}
			//log.Println("Res:", res)
			return nil
		}),
	}
	err := chromedp.Run(ctx,
		tasks,
	)
	if err != nil {
		log.Println(err)
	}
	return res
}

func login(user string, pws string) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			res := ""
			chromedp.Navigate("https://www.itjuzi.com/login").Do(ctx)
			/*chromedp.Focus("input[type='text'].el-input__inner",chromedp.ByQuery).Do(ctx)
			chromedp.SetValue("input[type='text'].el-input__inner", user , chromedp.ByQuery).Do(ctx)
			chromedp.SetAttributeValue("input[type='text'].el-input__inner","value", user , chromedp.ByQuery).Do(ctx)
			chromedp.Blur("input[type='text'].el-input__inner",chromedp.ByQuery).Do(ctx)
			chromedp.Focus("input[type='password'].el-input__inner",chromedp.ByQuery).Do(ctx)
			chromedp.SetValue("input[type='password'].el-input__inner", pws , chromedp.ByQuery).Do(ctx)
			chromedp.SetAttributeValue("input[type='password'].el-input__inner","value", pws , chromedp.ByQuery).Do(ctx)
			chromedp.Blur("input[type='password'].el-input__inner",chromedp.ByQuery).Do(ctx)
			chromedp.Value("input[type='text'].el-input__inner",&res,chromedp.ByQuery).Do(ctx)
			fmt.Println(res)
			time.Sleep(1 * time.Second)
			chromedp.Value("input[type='password'].el-input__inner",&res,chromedp.ByQuery).Do(ctx)
			fmt.Println(res)*/
			time.Sleep(30 * time.Second)
			chromedp.Click("button.el-button.el-button--primary", chromedp.ByQuery).Do(ctx)
			//chromedp.Navigate("https://www.itjuzi.com/company/34048160").Do(ctx)
			chromedp.OuterHTML("body", &res, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
			fmt.Println(res)
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

func getCompanyInfo(companyId string) (string, string) {
	var res string
	var leader string
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", true),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	headers := make(map[string]interface{})
	headers[":method"] = "GET"
	headers[":path"] = "/company/34048160"
	headers[":scheme"] = "https"
	headers["sec-fetch-site"] = "same-origin"
	headers[":authority"] = "www.itjuzi.com"
	headers[":scheme"] = "https"
	headers["user-agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36"
	headers["accept"] = `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`
	cookieStr := `_ga=GA1.2.67751875.1583219519; gr_user_id=8fb5ba2e-46ac-41af-b28c-0d410e400684; MEIQIA_TRACK_ID=1T6FhOb23V8FI0nxuiwkEY3i1VA; MEIQIA_VISIT_ID=1YeXpvG85uPFRAr7AYRcDg5Cqr1; Hm_lvt_1c587ad486cdb6b962e94fc2002edf89=1583219518,1583304030,1583809915; _gid=GA1.2.1762737053.1584003695; _gat_gtag_UA_59006131_1=1; juzi_user=695684; juzi_token=bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczpcL1wvd3d3Lml0anV6aS5jb21cL2FwaVwvYXV0aG9yaXphdGlvbnMiLCJpYXQiOjE1ODQwODg4OTUsImV4cCI6MTU4NDA5MjQ5NSwibmJmIjoxNTg0MDg4ODk1LCJqdGkiOiJ4NWxabnhpdzV0QWpXVmh5Iiwic3ViIjo2OTU2ODQsInBydiI6IjIzYmQ1Yzg5NDlmNjAwYWRiMzllNzAxYzQwMDg3MmRiN2E1OTc2ZjciLCJ1dWlkIjoiMGo2U2NJIn0.QXXmrnllNjlLKj4-X7HFFKM9TVnRqrRf-fO66QTY3HQ; Hm_lpvt_1c587ad486cdb6b962e94fc2002edf89=1584088896`
	headers["cookie"] = cookieStr
	//headers["curlopt_followlocation"] = "true"
	headers["sec-fetch-dest"] = "document"
	headers["sec-fetch-mode"] = "navigate"
	headers["sec-fetch-site"] = "none"
	headers["sec-fetch-user"] = "?1"
	headers["upgrade-insecure-requests"] = 1
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			network.Enable().Do(ctx)
			network.SetExtraHTTPHeaders(network.Headers(headers)).Do(ctx)
			ss := strings.Split(cookieStr, ";")
			for _, kv := range ss {
				sss := strings.Split(kv, "=")
				fmt.Println(sss[0], sss[1])
				c := network.SetCookie(sss[0], sss[1])
				c.Domain = ".itjuzi.com"
				c.Do(ctx)
			}
			chromedp.Navigate("https://www.itjuzi.com/company/" + companyId).Do(ctx)

			/*chromedp.Click("button.driver-next-btn",chromedp.NodeVisible,chromedp.ByQuery).Do(ctx)
			time.Sleep(200 * time.Microsecond)
			chromedp.Click("button.driver-next-btn",chromedp.NodeVisible,chromedp.ByQuery).Do(ctx)
			time.Sleep(200 * time.Microsecond)
			chromedp.Click("button.driver-next-btn",chromedp.NodeVisible,chromedp.ByQuery).Do(ctx)
			time.Sleep(200 * time.Microsecond)
			chromedp.Click("button.driver-next-btn",chromedp.NodeVisible,chromedp.ByQuery).Do(ctx)
			*/
			chromedp.InnerHTML("body", &res, chromedp.ByQuery).Do(ctx)
			log.Println("Res:", res)
			log.Println(strings.Contains(res, "登录后查看"))
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
