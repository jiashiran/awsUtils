package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net"
	"os"
	"spider/pkg/uuid"
	"time"
)

func main() {
	check360()

	/*str := `{"self_imei":"172.16.27.71","header":{"appid":"6512cb299dd5e4e9","nonce":"aed6d2fc-9444-4f3b-ae25-f5ec2f9ed376","timestamp":1568808844017,"v":"1.0","sign":"d9c58a75d5f36459"},"phone_call_ask":[]}`
	strBytes := []byte(str)
	bytes := make([]byte,len(strBytes)+6)
	bytes = append(bytes,0x06)
	bytes = append(bytes,0x51)
	headerShort := 0
	headerShort ++
	bytes = append(bytes,byte(headerShort))
	bytes = append(bytes,byte(headerShort >> 2))
	bytes = append(bytes,0x00)
	bytes = append(bytes,0x00)
	bytes = append(bytes,strBytes...)
	resp,_ := http.Post("http://openapi.shouji.360.cn/OpenapiPhoneQuery","application/json",b.NewBuffer(bytes))
	bs,_:=ioutil.ReadAll(resp.Body)
	log.Println(string(bs))*/
}

/**
测试环境
number.certificate.check.appid=test_trrt_appid
number.certificate.check.secret=test_trrt_secret
number.certificate.check.version=1.0
number.certificate.check.url=http://182.118.28.111/OpenapiPhoneQuery
number.certificate.check.numbers=20000
线上环境
number.certificate.check.appid=6512cb299dd5e4e9
number.certificate.check.secret=49d6e68a09e5627acdee6807724a8889
number.certificate.check.version=1.0
number.certificate.check.url=http://openapi.shouji.360.cn/OpenapiPhoneQuery
number.certificate.check.numbers=20000
*/
func check360() {
	timestamp := time.Now().UnixNano() / 1e6
	fmt.Printf("时间戳（毫秒）：%v;\n", timestamp)
	uuidStr := uuid.GetUUID()
	_ = uuidStr
	//appiD nonce  timestamp v
	var request Request360
	request.Headler = Header{"6512cb299dd5e4e9", "9703f92a61b67e5d", timestamp, "1.0", ""}
	str := fmt.Sprint(request.Headler.Appid) + fmt.Sprint(request.Headler.Nonce) + fmt.Sprint(request.Headler.Timestamp) + fmt.Sprint(request.Headler.V) + "4009608530" + "49d6e68a09e5627acdee6807724a8889"
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	log.Println("sign:", sign)
	request.Headler.Sign = sign[16:]
	request.PhonCallAsk = PhonCallAsk{0, "4009608530"}
	request.SelfImei = GetIntranetIp()
	args := fasthttp.AcquireArgs()
	bs, err := json.Marshal(request)
	log.Println(string(bs))

	headerShort := 0
	headerShort++
	bytes := []byte{0x06, 0x51, byte(headerShort), byte(headerShort >> 2), 0x00, 0x00}

	headler := make([]byte, len(bs)+6)
	for _, b := range bytes {
		headler = append(headler, b)
	} //8659   14:43:59
	for _, b := range bs {
		headler = append(headler, b)
	}
	log.Println(headler)
	args.ParseBytes(headler)

	code, body, err := fasthttp.Post(nil, "http://openapi.shouji.360.cn/OpenapiPhoneQuery", args)
	log.Println(code, string(body), err)
}

func GetIntranetIp() string {
	ip := ""
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("ip:", ipnet.IP.String())
				ip = ipnet.IP.String()
			}

		}
	}
	return ip
}

type Request360 struct {
	Headler     Header      `json:"headler"`
	PhonCallAsk PhonCallAsk `json:"phone_call_ask"`
	SelfImei    string      `json:"self_imei"`
}

type Header struct {
	Appid     string `json:"appid"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	V         string `json:"v"`
	Sign      string `json:"sign"`
}

type PhonCallAsk struct {
	Id       int    `json:"id"`
	PhoneNum string `json:"phone_num"`
}
