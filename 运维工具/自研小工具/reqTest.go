```
//批量探测服务
package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var(
	tr =&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//Proxy: proxy,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client = &http.Client{
		Transport: tr,
		Timeout: 5 * time.Second,
	}
	//http代理调试
	//proxy = func(_ *http.Request) (*url.URL, error) {
	//	return url.Parse("http://127.0.0.1:8080")
	//}
	p = flag.String("p", "", "文件路径")
	t = flag.Int("t", 10, "并发数")

)
type Parameter struct {
	URL     string            `json:"featuresurl"`
	Method  string            `json:"method"`
	Data    string            `json:"data"`
	Headers http.Header       `json:"headers"`
}
type Response struct {
	Title     string           			`json:"body"`
	StatusCode 	  int          			`json:"statuscode"`
	RespHeaders map[string][]string 	`json:"respheaders"`
	Err 	error			   			`json:"err"`
}
// Req Request请求
func Req(link Parameter) (reqResp Response) {
	reqBodyReader := strings.NewReader(link.Data)
	request, err := http.NewRequest(link.Method, link.URL, reqBodyReader)
	if err != nil {
		resp:= Response{
			Err: err,
		}
		return resp
	}
	// set headers
	for key, values := range link.Headers {
		for i := range values {
			if i == 0 {
				request.Header.Set(key, values[i])
			} else {
				request.Header.Add(key, values[i])
			}
		}
	}
	response, err := client.Do(request)
	if err!=nil{
		resp:= Response{
			Err: err,
		}
		return resp
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp:= Response{
			Err: err,
		}
		return resp
	}
	r := regexp.MustCompile(`<title>(.*?)</title>`)
	title := r.FindString(string(body))
	resp:= Response{
		Title: title,
		StatusCode: response.StatusCode,
		RespHeaders: map[string][]string{
			"Content-Type": {response.Header.Get("Content-Type")},
		},
	}
	return resp
}

//请求封装
func urlBurst(reqBurstUrl string) (respBody Response)  {
	//reqBurstUrl 组装后的Url
	req:= Parameter{
		URL:    reqBurstUrl,
		Method: "GET",
		Headers: map[string][]string{
			"User-Agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4542.2 Safari/537.36"},
		},
	}
	respBody = Req(req)
	return respBody
}
func readFile(path string)(urls []string) {
	fi, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		url:=string(a)
		urls=append(urls,url)
	}
	return urls
}

/*并发*/
type GliMit struct {
	Num int
	C   chan struct{}
}

func NewG(num int) *GliMit {
	return &GliMit{
		Num: num,
		C : make(chan struct{}, num),
	}
}

func (g *GliMit) Run(f func()){
	g.C <- struct{}{}
	go func() {
		f()
		<-g.C
	}()
}

func concurrent(path string,thread int)  {
	ch := make(chan string)
	//   "/Users/spirit/Project/golang/other/reqTest/test.json"
	urls:=readFile(path)
	for _,url:=range urls{
		url := url
		go func() {
			ch <- url
		}()
	}
	// 限制线程数
	g := NewG(thread)
	wg := sync.WaitGroup{}
	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		goFunc := func() {
			// 做一些业务逻辑处理
			msg := <-ch
			resp:=urlBurst(msg)
			if resp.StatusCode !=0 {
				fmt.Println(msg,resp.Title,resp.StatusCode)
				time.Sleep(1 * time.Second)
				wg.Done()
			}
		}
		g.Run(goFunc)
	}
	wg.Wait()
}
func main()  {
	flag.Parse()
	if *p == "" {
		fmt.Println("请输入文件路径")
		flag.Usage()
		return
	}
	concurrent(*p,*t)
}
```
