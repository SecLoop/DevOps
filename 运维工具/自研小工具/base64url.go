```
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	var p = flag.String("e", "", "文件名")
	var p1 = flag.String("d", "", "base64url 内容")
	// 把用户传递的命令行参数解析为对应变量的值
	flag.Parse()
	switch  {
		case *p!="":
			// 如果要用在url中，需要使用URLEncoding
			fmt.Println(*p)
			data, err := ioutil.ReadFile(*p)
			if err != nil {
				fmt.Println(err.Error())
			}
			uEnc:= base64.URLEncoding.EncodeToString(data)
			fmt.Println(uEnc)
		case *p1!="":
			uDec, err := base64.URLEncoding.DecodeString(*p1)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(uDec))
		default:
			fmt.Println("输入异常")
	}
}
```
使用说明

1、go build 编译
base64url -e  文件名   编码
base64url -d  内容     编码
			uEnc:= base64.URLEncoding.EncodeToString(data)
