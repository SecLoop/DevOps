```
package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io"
	"os"
	"regexp"
	"strings"
)
//特殊字符清洗
func filter(str string, blackNames []string) string {
	var result3 string
	for _, name := range blackNames {
		str = strings.Replace(str, name, "", -1)
		//去除端口号
		re := regexp.MustCompile("(?im):(.*?)$")
		result:=re.ReplaceAllString(str, "")
		//去除最后带/和数字不是ip段
		re1 := regexp.MustCompile(`(?im)\/$`)
		result1:=re1.ReplaceAllString(result, "")
		//去除path路径
		re2 := regexp.MustCompile(`(?im)\/\D.*`)
		result2:=re2.ReplaceAllString(result1, "")
		//去除路径为数据和ip段分类起冲突
		match, _ := regexp.MatchString(`(?im)\/([1-9]|[1-2]\d|3[0-2])$`,result2)
		if match==false{
			re3 := regexp.MustCompile(`(?im)\/.*`)
			result3=re3.ReplaceAllString(result2, "")
		}
	}
	return result3
}
//读文件
func fileRead(fileName string)[]string  {
	fi, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	var result []string
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		result=append(result,string(a))
	}
	return result
}
//数组去重
func removeDuplicateElement(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
func main()  {
	var checkSuccess []string
	data:=fileRead("/Users/spirit/Project/golang/text/data.txt")
	for _,name:=range data{
		blackNameList := [] string {`"`,`'`,`https://`,`http://`}
			result:=filter(name,blackNameList)
			ip, _ := regexp.Match(`((25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)`, []byte(result))
			_, icann := publicsuffix.PublicSuffix(result)
			if ip == true || icann ==true{
				checkSuccess=append(checkSuccess,result)
			}else {
				fmt.Println("脏数据:",result)
			}
		}
	SuccessData:=removeDuplicateElement(checkSuccess)
	for _,successIp:=range SuccessData{
		//数据分类
		iCdr1, _ := regexp.MatchString(`^([0,1]?\d{1,2}|2([0-4][0-9]|5[0-5]))(\.([0,1]?\d{1,2}|2([0-4][0-9]|5[0-5]))){3}/\d$`, successIp)
		iCdr2, _ := regexp.MatchString(`^([0,1]?\d{1,2}|2([0-4][0-9]|5[0-5]))(\.([0,1]?\d{1,2}|2([0-4][0-9]|5[0-5]))){3}-([0,1]?\d{1,2}|2([0-4][0-9]|5[0-5]))(\.([0,1]?\d{1,2}|2([0-4][0-9]|5[0-5]))){3}$`, successIp)
		aIpAddress, _ := regexp.MatchString(`^10\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])$`, successIp)
		bIpAddress, _ := regexp.MatchString(`^172\.(1[6789]|2[0-9]|3[01])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])$`, successIp)
		cIpAddress, _ := regexp.MatchString(`^192\.168\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])$`, successIp)
		switch {
		case iCdr1==true||iCdr2==true:
			fmt.Println(fmt.Sprintf("IP段为:%s",successIp))
		case aIpAddress==true||bIpAddress==true||cIpAddress==true:
			fmt.Println(fmt.Sprintf("内网地址为:%s",successIp))
		case successIp=="0.0.0.0" || successIp=="127.0.0.1":
			fmt.Println(fmt.Sprintf("其他地址为:%s",successIp))
		default:
			fmt.Println(fmt.Sprintf("IP地址为:%s",successIp))
		}
	}
}
```
运行main函数
