package main

import (
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"
)

// masscan analysis request

type Nmaprun struct {
	XMLName          xml.Name `xml:"nmaprun"`
	Text             string   `xml:",chardata"`
	Scanner          string   `xml:"scanner,attr"`
	Start            string   `xml:"start,attr"`
	Version          string   `xml:"version,attr"`
	Xmloutputversion string   `xml:"xmloutputversion,attr"`
	Scaninfo         struct {
		Text     string `xml:",chardata"`
		Type     string `xml:"type,attr"`
		Protocol string `xml:"protocol,attr"`
	} `xml:"scaninfo"`
	Host []struct {
		Text    string `xml:",chardata"`
		Endtime string `xml:"endtime,attr"`
		Address struct {
			Text     string `xml:",chardata"`
			Addr     string `xml:"addr,attr"`
			Addrtype string `xml:"addrtype,attr"`
		} `xml:"address"`
		Ports struct {
			Text string `xml:",chardata"`
			Port struct {
				Text     string `xml:",chardata"`
				Protocol string `xml:"protocol,attr"`
				Portid   string `xml:"portid,attr"`
				State    struct {
					Text      string `xml:",chardata"`
					State     string `xml:"state,attr"`
					Reason    string `xml:"reason,attr"`
					ReasonTtl string `xml:"reason_ttl,attr"`
				} `xml:"state"`
			} `xml:"port"`
		} `xml:"ports"`
	} `xml:"host"`
	Runstats struct {
		Text     string `xml:",chardata"`
		Finished struct {
			Text    string `xml:",chardata"`
			Time    string `xml:"time,attr"`
			Timestr string `xml:"timestr,attr"`
			Elapsed string `xml:"elapsed,attr"`
		} `xml:"finished"`
		Hosts struct {
			Text  string `xml:",chardata"`
			Up    string `xml:"up,attr"`
			Down  string `xml:"down,attr"`
			Total string `xml:"total,attr"`
		} `xml:"hosts"`
	} `xml:"runstats"`
}

var (
	//忽略ssl证书，并添socks5代理
	//dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:7890", nil, proxy.Direct)
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
	//	os.Exit(1)
	//}
	//// setup a http client
	//httpTransport := &http.Transport{}
	//client := &http.Client{Transport: httpTransport}
	//// set our socks5 as the dialer
	//httpTransport.Dial = dialer.Dial

	//忽略ssl证书
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
	}
	p = flag.String("p", "", "Masscan Xml 文件路径")
)

func readXml(path string)(urls []string)  {
	file, err := os.Open(path) // For read access.
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return
	}
	v := Nmaprun{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _,v:=range v.Host {
		httpUrl := fmt.Sprintf("http://%s:%s", v.Address.Addr, v.Ports.Port.Portid)
		httpsUrl := fmt.Sprintf("https://%s:%s", v.Address.Addr, v.Ports.Port.Portid)
		urls=append(urls,httpUrl,httpsUrl)
	}
	return urls
}

func getRequest(url string) (respCode,title string ,err error) {
	req,_ := http.NewRequest("GET",url,nil)
	req.Header.Set("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0")
	resp,err := client.Do(req)
	if err!=nil{
		return "","",err
	}
	body, _:= ioutil.ReadAll(resp.Body)
	r := regexp.MustCompile(`<title>(.*?)</title>`)
	title = r.FindString(string(body))
	respCode = resp.Status
	return respCode,title,nil
}

func main()  {
	flag.Parse()
	if *p == "" {
		fmt.Println("请输入Masscan Xml 文件路径")
		flag.Usage()
		return
	}
	urls:=readXml(*p)
	for _,url:=range urls{
		if respCode,title,_:=getRequest(url); respCode!="" {
			fmt.Println(url,title,respCode)
		}

	}
}
