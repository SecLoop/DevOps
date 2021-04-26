```
package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/gomail.v2"
	_ "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	info string
	items int
	result string

	httpTransport = &http.Transport{}
	client = &http.Client{
		Transport: httpTransport,
		Timeout:   time.Second*60,
	}
)

func times() (lastWeek string)  {
	nTime := time.Now()
	yesTime := nTime.AddDate(0,0,-1) //修改监控时间
	lastWeek = yesTime.Format("2006-01-02")
	return lastWeek
}
func removeRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
//消息源
func expMonitor()(result []string)  {
	for i:=1;;i++{
		//dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:2333", nil, proxy.Direct)
		//if err != nil {
		//	fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		//	os.Exit(1)
		//}
		//// setup a http client
		//// set our socks5 as the dialer
		//httpTransport.Dial = dialer.Dial
		cveId1:=`cve-`+strconv.Itoa(time.Now().Year())
		cveId2:=`cve_`+strconv.Itoa(time.Now().Year())
		cveId3:=`cve`+strconv.Itoa(time.Now().Year())
		search:=[]string{cveId1,cveId2,cveId3,"exploits","payload","Security","privilege","rce_exp","rce-exp","rce_poc","rce-poc","_rce","-rce","%E6%8F%90%E6%9D%83","%E6%BC%8F%E6%B4%9E"}
		for _,v:=range search{
			url:=`https://api.github.com/search/repositories?q=`+v+`+size:>0+created:>`+ times()+`&sort=updated&order=desc&page=`+strconv.Itoa(i)+`&per_page=100`
			fmt.Println(url)
			req, _ := http.NewRequest("GET",url , nil)
			req.Header.Set("User-Agent", `Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36`)
			req.Header.Set("Content-Type", `application/x-www-form-urlencoded`)
			req.Header.Set("Authorization", `token xxxxxxx`) //设置api token
			resp, _ := client.Do(req)
			gitHubApiRemaining:=resp.Header.Get("x-ratelimit-remaining") //github api 剩余次数
			if gitHubApiRemaining=="0"{
				fmt.Println("Github 超过使用次数")
				return nil
			}
			fmt.Println(fmt.Sprintf("访问URL:%s,API剩余次数%s",url,gitHubApiRemaining))
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			re:= regexp.MustCompile(`"html_url": "([^"]+)",\n\s+"description": ([^,]+)`)
			oneMatch := re.FindAllSubmatch(body,-2)
			for _,n:=range oneMatch{
				htmlUrl:=string(n[1])
				description:=string(n[2])
				switch  {
				case v==cveId1 || v==cveId2 || v==cveId3:
					info="Cve地址："+htmlUrl+"\t\t"+`	`+"Cve描述："+description
				case v=="exploits" || v=="rce_exp" || v=="rce-exp" ||v=="payload" || v=="rce_poc" || v=="rce-poc" || v=="_rce" || v=="-rce" || v=="%E6%BC%8F%E6%B4%9E":
					info="Exp地址："+htmlUrl+"\t\t"+`	`+"Exp描述："+description
				case v=="privilege" || v=="%E6%8F%90%E6%9D%83":
					info="提权地址："+htmlUrl+"\t\t"+`	`+"提权描述："+description
				}
				result=append(result,info)
			}
			//数据是否为空判断
			items = int(gjson.Parse(string(body)).Get("items.#").Int())
			if items == 0{
				continue
			}
			time.Sleep(1e6) //防止请求过快
		}
		return result
	}
}

func arrayToString(arr []string) string {
	for _, i := range arr {  //遍历数组中所有元素追加成string
		result += fmt.Sprintln(fmt.Sprintf("%s\n",i))
	}
	return result
}

func SendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填填授权码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "xxxxxx@163.com",
		"pass": "xxxxxxxx",
		"host": "smtp.163.com",
		"port": "465",
	}
	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int
	m := gomail.NewMessage()
	m.SetHeader("From",  m.FormatAddress(mailConn["user"], "漏洞监控系统")) //这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/plain", body)    //设置邮件正文
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	err := d.DialAndSend(m)
	return err
}
func main()  {
	content := expMonitor()
	removeRepeatContent:=removeRepeatedElement(content)
	monitorContent:=arrayToString(removeRepeatContent)
	//定义收件人
	mailTo := []string{
		"xxxxxx@163.com",
	}
	//邮件主题为"Hello"
	subject := "Bingo 今日份安全监控已到账"+times()
	// 邮件正文
	body := monitorContent
	err := SendMail(mailTo, subject, body)
	if err != nil {
		log.Println(err)
		fmt.Println("send fail")
		return
	}
	fmt.Println("send successfully")
}
```
