package main

import (
	"bufio"
	"fmt"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	sshd "spider/pkg/ssh"
	"spider/pkg/xueqiu"
	"strings"
)

var (
	serverMap map[string]*ssh.Client

//ips = []string{"39.105.202.8","39.107.243.195","39.96.50.197","123.56.17.69","47.95.241.44","47.94.130.119","47.93.218.38","101.200.52.2"}
)

func main() {
	//jijin.Run()
	xueqiu.Run()

	//areaCode.Get()
	//areaCode.Detail("france")
}

func rep(s string) string {
	if strings.Index(s, ` `) > 0 {
		s = strings.ReplaceAll(strings.Split(s, ` `)[0], " ", "")
	}
	if strings.Index(s, `or`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `or`)[0], " ", "")
	}
	if strings.Index(s, `+`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `+`)[0], " ", "")
	}
	if strings.Index(s, `&lt;`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `&lt;`)[0], " ", "")
	}
	if strings.Index(s, `digits`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `digits`)[0], " ", "")
	}
	if strings.Index(s, `XX`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `XX`)[0], " ", "")
	}
	if strings.Index(s, `)`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `)`)[0], " ", "")
	}
	return s
}

//email verify
func VerifyEmailFormat(email string) string {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	//pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	es := reg.FindAllString(email, -1)
	/*if len(es) > 0{
		log.Println(es)
	} else {
		//log.Println(email)
	}*/

	return fmt.Sprint(es)
}

//mobile verify
func VerifyMobileFormat(mobileNum string) string {
	//regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	regular := `(\(\d{3,4}\)|\d{3,4}-|\s)?\d{7,14}` //\d+`
	reg := regexp.MustCompile(regular)
	ss := reg.FindAllString(mobileNum, -1)
	if len(ss) > 0 {
		//log.Println(ss)
	} else {
		//log.Println(mobileNum)
	}
	return fmt.Sprint(ss)
}

//获得百度token
func getaccess_token() {
	value := make(url.Values)
	resp, _ := http.PostForm("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=loC7VdqRiDBcDtdOxsARjeZc&client_secret=l0Z5yry7YLbGnDIIqQKSKlvIKGOGY6gs&", value)
	bs, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bs))
	// "access_token":"24.cb9afce813b0dc1449529649bfc2a3dc.2592000.1570863615.282335-17235293"
}

func post() {
	//face++
	/*value := make(url.Values)
	value.Add("api_key","-QZz6z1pR_D5_X8KQV3mFwA02PS-Z8AU")
	value.Add("api_secret","MaFmJt0LIon69262LaEsbkD28q6zVjOn")
	value.Add("image_base64","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")
	resp,_ := http.PostForm("https://api-cn.faceplusplus.com/imagepp/v1/recognizetext",value)
	bs,_:=ioutil.ReadAll(resp.Body)
	log.Println(string(bs))*/

	//baidu
	value := make(url.Values)
	value.Add("language_type", "CHN_ENG")
	value.Add("Content-Type", "application/x-www-form-urlencoded")
	value.Add("image", "iVBORw0KGgoAAAANSUhEUgAAAKAAAAARBAMAAACod7rOAAAAG1BMVEX///8AAP9fX/+fn/9/f//f3/8/P/+/v/8fH/8EVqDgAAACKUlEQVQ4jdWTS3PaMBRGv9qS8dIqGLI0kCYsDU2cLIE+l3gy43rp4Dy6NJ1J66UDofCze69kaOAXtPIMMnocnXuvDPwHrQ985H6DEtYYmOpfRDMkIXAKJAVgZzJDMuMXmpNKRUopGv9JD3CruFV/gbyyBpo5FUBk8gUaKALNmeOB/jsBA703dJhLwDniouY09kDBhNmhoW24za5qPlPPwORmjYkZfudFjjYUPh18BOyqjnnZAXdN+tCGGyC2lBoSoQXkesrLHW24CFnS0qec7wxXpntcZeUFAa1Q53RSR26Aew+9t5K+r4FySUPp+NBQNjFUJ+SzWZULAooKsvXKUGqge1/HE/JO6fV0yAtdoeXsACja9xFSUvA/l+6ITnuBVdVRqLffbjs6h8NtyIuTa2PoxclV5RYjCoysxzEHswOm+VQDL/2ktKcEXITDYm94tkSHDZePPdrkfXBrQycK15KsxZqXxkTdA3t9jAgIu6qLYrcoYmEMWxQyc0Ve53AHTKNwdE2L49ApjoCf+pT4tBFYQQ2Uz2ZzOuDcnNq8TtzsgCbkdtsJXZaba8cDIPqyQ8CMLqgBfs1/c8hyLjYMdHQR7SNDOFf3dD2ogvhOwNc5RD8NEIuAPgMGNroriO0XNAYhyrzA3WJ2DBQMdNtFOcbw1/ThLGBDnbMaSF+jc/KUGWB5znOD90uu6o/VZH0BA+Ritxk4VCEB6brKvKHuLiHNjdqae/Dvtz8b53kRh1rTrAAAAABJRU5ErkJggg==")
	resp, _ := http.PostForm("https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token=24.cb9afce813b0dc1449529649bfc2a3dc.2592000.1570863615.282335-17235293", value)
	bs, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bs))
}

//企查查接口
func Qichacha() {

}

func get() {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://www.qichacha.com/search?key=小米")
	req.Header.Add("Cookie", "acw_tc=6f7e78a715668012034494032efb97f66ec87cb68fb44b21e12d5f8b0d; zg_did=%7B%22did%22%3A%20%2216ccca116221a2-0f050a5165c1-38607701-13c680-16ccca116235f5%22%7D; UM_distinctid=16ccca11b11b0-0284524a81e822-38607701-13c680-16ccca11b12b6e; QCCSESSID=4ddso719atrpgkbmukneu4spe7; hasShow=1; CNZZDATA1254842228=1950090108-1566797068-https%253A%252F%252Fwww.baidu.com%252F%7C1568626442; Hm_lvt_3456bee468c83cc63fb5147f119f1075=1568202279,1568254024,1568614161,1568628815; Hm_lpvt_3456bee468c83cc63fb5147f119f1075=1568628815; zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f=%7B%22sid%22%3A%201568625694403%2C%22updated%22%3A%201568628815114%2C%22info%22%3A%201568614160398%2C%22superProperty%22%3A%20%22%7B%7D%22%2C%22platform%22%3A%20%22%7B%7D%22%2C%22utm%22%3A%20%22%7B%7D%22%2C%22referrerDomain%22%3A%20%22www.qichacha.com%22%2C%22zs%22%3A%200%2C%22sc%22%3A%200%2C%22cuid%22%3A%20%221d0083f0fc0802d50c303cfd3dedfa79%22%7D")
	req.Header.Add("Host", "www.qichacha.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Cache-Control", "max-age=0")

	var client fasthttp.Client

	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(resp.Body()))
}

func download() {
	imgPath := "."
	imgUrl := "http://data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAKAAAAARBAMAAACod7rOAAAAG1BMVEX///8AAP/f3/+/v/9/f/9fX/8fH/8/P/+fn/83sGquAAACUElEQVQ4jdWTT1PbMBDFH7Js+Wi50HCUBya5OtOOuSadOuTohEI4Ovw/GmggRwHF6cfurhwSQr9Aq8x4HHn3p/d2V8C/v8TXkxTK0psX06NjplrrVNGjRX+V3vG0LoAfwCuHhxXnaL3gKOCAfsCt5lU2wMHts+ytgf45amyl/OUQKlU290pJwHPcJ4ZPKBgYbXUsAoqaIUmXyrwlEBnEbA3E6TvgoyFgWBEwLMNSVug7KXocLZRTKGPI6m8g1IEDRs3GCqgqtiwLsqlsYEkOMHdliubKKexaFum7U46WwJHWxxQvm7O3pxEDjQMSgBRmRkyrrpkA7XGTW4o4dkBxSWHt3geFA6MqZzlkyyKmjL4D9isHPOX36f4ZuDMjy5kiOnaWu65Dl2YTeHWBIwf0t52blcJ6wu3PqYW7cnRAzd9Bp24URknnpQxS6vs95fUS2izegOFDx0A9EyKpV0BXw8ybGaqh3xOxpA7HOMYwWCpUC3sm7gDJwpEQdQUMMuz9Gn2maZiwTgcM0iZMFqRQxH5BmukDzeIbsL2wrzUVL7Eq/QDMM3l9mlFceCbOGTiuoV542gj7PUhzHA4MAWk8GNhYbrWo7Sxu5jRuAL/R2LibcmJAXNG6qyGHXG1ZYpaYHO1dEJAGeL5WSGfe3bgQPBDwfQ2xBNZUEXGdejtkWVwUrl2/Pw2Q71/WY2ztXQg2wEDJwKCVZj08fckfJwUrJIebwPDebd6EBQHRd0X0p8WwuztLsbB01RK+0gR80paAPpmZe/rnM7g3fPqVxX+x/gC4nIDbXy+SzAAAAABJRU5ErkJggg=="

	fileName := "tmp.png"

	res, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("A error occurred:", err)
		return
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	file, err := os.Create(imgPath + fileName)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d", written)
}

func initHttpProxyServer() {
	ips := []string{"39.107.243.195"}
	sshd.InitServer(ips)
	pids := sshd.ExecuteBatch(`ps -ef|grep httpproxy | grep -v grep | awk -F " " '{print $2}' `)

	for ip, _ := range serverMap {
		sshd.Execute(ip, "kill -9 "+pids[ip])
	}

	sshd.UploadBatch("./pkg/httpproxy/httpproxy", "/var/tmp/")
	sshd.ExecuteBatch(`echo "" >> /var/tmp/startHttpProxy.sh`)
	sshd.ExecuteBatch(`echo "/var/tmp/httpproxy 1 >> /dev/null 2>&1 &" > startHttpProxy.sh`)
	sshd.ExecuteBatch(`chmod +x /var/tmp/httpproxy`)
	sshd.ExecuteBatch(`chmod +x startHttpProxy.sh`)
	sshd.ExecuteBatch(`./startHttpProxy.sh`)
	sshd.ExecuteBatch(`ps -ef|grep httpproxy | grep -v grep | awk -F " " '{print $2}' `)
}

func httpd(ip, urlStr string) string {
	proxyUrl, err := url.Parse("http://" + ip + ":8888")
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	res, err := myClient.Get(urlStr)

	if err != nil {
		log.Println(err)
	}

	bs, _ := ioutil.ReadAll(res.Body)
	//log.Println(string(bs))
	return string(bs)
}
