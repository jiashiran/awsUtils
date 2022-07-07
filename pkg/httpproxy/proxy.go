package main

import (
	"fmt"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Pxy struct {
	Cfg Cfg
}

// 设置type
type Cfg struct {
	Addr        string // 监听地址
	Port        string // 监听端口
	IsAnonymous bool   // 高匿名模式
	Debug       bool   // 调试模式
}

var (
	remoteAddr = flag.String("remoteAddr", "123.56.62.34:80", "your remote addr")
	remoteHost = flag.String("remoteHost", "data-platform-prometheus.com", "your remote host")
	scheme     = flag.String("scheme", "http", "scheme http or https")
)

func main() {
	addr := flag.String("addr", "0.0.0.0", "your listen addr")
	port := flag.String("port", "8080", "your listen port")
	//faddr := "0.0.0.0"//flag.String("addr","0.0.0.0","监听地址，默认0.0.0.0")
	//fprot := "8080"//flag.String("port","8080","监听端口，默认8080")
	//fanonymous :=  flag.Bool("anonymous",true,"高匿名，默认高匿名")
	//fdebug :=  flag.Bool("debug",false,"调试模式显示更多信息，默认关闭")
	flag.Parse()
	var rootCmd = &cobra.Command{
		Use:   "proxy",
		Short: "proxy is a very fast http proxy",
		Long:  `A Fast`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &Cfg{}
			cfg.Addr = *addr
			cfg.Port = *port
			cfg.IsAnonymous = true
			cfg.Debug = false
			Run(cfg)
		},
	}
	rootCmd.Execute()

}
func a() {
	//定义Flag
	//方式一：通过flag.String(), Bool(), Int() 等flag.Xxx()方法，该种方式返回一个相应的指针
	namePtr := flag.String("name", "Anson", "user's name")
	agePtr := flag.Int("age", 22, "user's age")
	vipPtr := flag.Bool("vip", true, "is a vip user")
	//方式二：通过flag.XxxVar()方法将flag绑定到一个变量，该种方式返回值类型
	var email string
	flag.StringVar(&email, "email", "abc@gmail.com", "user's email")
	//还有第三种方式，通过flag.Var()绑定自定义类型，自定义类型需要实现Value接口(Receives必须为指针)
	//flag.Var(&flagVal, "name", "help message for flagname")

	//解析命令行参数,值保存到定义的flag
	flag.Parse()

	//调用Parse解析后，就可以直接使用flag本身(指针类型)或者绑定的变量了(值类型)
	//还可通过flag.Args(), flag.Arg(i)来获取非flag命令行参数
	others := flag.Args() //保存Flag以外的变量
	fmt.Println("name:", *namePtr)
	fmt.Println("age:", *agePtr)
	fmt.Println("vip:", *vipPtr)
	fmt.Println("email:", email)
	fmt.Println("other:", others)
	fmt.Println("---------")
	for i := 0; i < len(flag.Args()); i++ {
		fmt.Println("Arg", i, "=", flag.Arg(i))
	}
}

func Run(cfg *Cfg) {
	pxy := NewPxy()
	pxy.SetPxyCfg(cfg)
	log.Printf("HttpPxoy is runing on %s:%s \n", cfg.Addr, cfg.Port)
	// http.Handle("/", pxy)
	bindAddr := cfg.Addr + ":" + cfg.Port
	log.Fatalln(http.ListenAndServe(bindAddr, pxy))
}

// 实例化
func NewPxy() *Pxy {
	return &Pxy{
		Cfg: Cfg{
			Addr:        "",
			Port:        "8081",
			IsAnonymous: true,
			Debug:       false,
		},
	}
}

// 配置参数
func (p *Pxy) SetPxyCfg(cfg *Cfg) {
	if cfg.Addr != "" {
		p.Cfg.Addr = cfg.Addr
	}
	if cfg.Port != "" {
		p.Cfg.Port = cfg.Port
	}
	if cfg.IsAnonymous != p.Cfg.IsAnonymous {
		p.Cfg.IsAnonymous = cfg.IsAnonymous
	}
	if cfg.Debug != p.Cfg.Debug {
		p.Cfg.Debug = cfg.Debug
	}

}

// 运行代理服务
func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// debug
	if p.Cfg.Debug {
		log.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)
		// fmt.Println(req)
	}

	// http && https
	if req.Method != "CONNECT" {
		// 处理http
		p.HTTP(rw, req)
	} else {
		// 处理https
		// 直通模式不做任何中间处理
		p.HTTPS(rw, req)
	}

}

// http
func (p *Pxy) HTTP(rw http.ResponseWriter, req *http.Request) {

	transport := http.DefaultTransport

	// 新建一个请求outReq
	outReq := new(http.Request)
	// 复制客户端请求到outReq上
	*outReq = *req // 复制请求
	outReq.URL.Host = *remoteAddr
	outReq.URL.Scheme = *scheme
	//  处理匿名代理
	if p.Cfg.IsAnonymous == false {
		if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
				clientIP = strings.Join(prior, ", ") + ", " + clientIP
			}
			outReq.Header.Set("X-Forwarded-For", clientIP)
		}
	}

	outReq.Host = *remoteHost
	// outReq请求放到传送上
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte(err.Error()))
		return
	}

	// 回写http头
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	// 回写状态码
	rw.WriteHeader(res.StatusCode)
	// 回写body
	io.Copy(rw, res.Body)
	res.Body.Close()
}

// https
func (p *Pxy) HTTPS(rw http.ResponseWriter, req *http.Request) {
	// 拿出host
	host := req.URL.Host
	hij, ok := rw.(http.Hijacker)
	if !ok {
		log.Printf("HTTP Server does not support hijacking")
	}

	client, _, err := hij.Hijack()
	if err != nil {
		return
	}

	// 连接远程
	server, err := net.Dial("tcp", host)
	if err != nil {
		return
	}
	client.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n"))

	// 直通双向复制
	go io.Copy(server, client)
	go io.Copy(client, server)
}
