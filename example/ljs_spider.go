```
package ljs

import (
	"Monitor/pkg/public"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"log"
	"regexp"
)

/*
参考:https://juejin.cn/post/7033295246171734030 乐加速破解
*/

type TmpJslClearance struct {
	Code   int    `json:"code"`
	JslUid string `json:"jsluid"`
	JsVale string `json:"js"`
}
type JslClearanceRespData struct {
	Bts   []string `json:"bts"`
	Chars string   `json:"chars"`
	Ct    string   `json:"ct"`
	Ha    string   `json:"ha"`
	Tn    string   `json:"tn"`
	Vt    string   `json:"vt"`
	Wt    string   `json:"wt"`
}

/*
Js执行器
*/
func jsActuator(js string) (jsActuatorData string, err error) {
	vm := otto.New()
	vale, err := vm.Run(js)
	if err != nil {
		return "", err
	}
	jsActuatorData = vale.String()
	return jsActuatorData, nil
}

/*
	data := news.JslClearanceRespData{
		[]string{"1640505971.07|0|x5ai", "WGIiz3pH8Aocl82Pdd8R4%3D"},
		"aavrKPLLUdLygtiBLlsRfs",
		"75d1868de0fdccdcdd9bf833f4ca914804f1fcce55a54bdcb07d7723d1e0877f",
		"sha256",
		"__jsl_clearance_s",
		"3600",
		"1500",
	}

	reslut x5airyWGIiz3pH8Aocl82Pdd8R4%3D
*/

/*
爬虫token选择器
*/

func spiderTokenRule(data JslClearanceRespData) (spiderToken string) {
	var vales, valeData string
	for _, i := range data.Chars {
		for _, j := range data.Chars {
			vales = data.Bts[0] + string(i) + string(j) + data.Bts[1]
			if data.Ha == "md5" {
				md5Hash := md5.Sum([]byte(vales))
				valeData = fmt.Sprintf("%x", md5Hash[:])
			} else if data.Ha == "sha1" {
				sha1Hash := sha1.Sum([]byte(vales))
				valeData = fmt.Sprintf("%x", sha1Hash[:])
			} else if data.Ha == "sha256" {
				sha256Hash := sha256.Sum256([]byte(vales))
				valeData = fmt.Sprintf("%x", sha256Hash[:])
			}
			if valeData == data.Ct {
				spiderToken = vales
			}
		}
	}
	return spiderToken
}

//获取临时的JS

func getJslClearanceTmpToken(url string) TmpJslClearance {
	link := public.Parameter{
		URL:    url,
		Method: "GET",
		Headers: map[string][]string{
			"User-Agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4684.0 Safari/537.36"},
		},
	}
	respData := public.HttpReq(link)
	/*respone header 区分大小写*/
	respSetCookieData := respData.RespHeaders["set-cookie"]
	/*将interface 转换为string */
	respSetCookieDataToString := fmt.Sprintf("%v", respSetCookieData)
	setCookieReg := regexp.MustCompile(`__jsluid_s=(.*?);`) //
	cookieData := setCookieReg.FindString(respSetCookieDataToString)
	//获取响应的JS
	jsReg := regexp.MustCompile(`\<script\>document\.cookie\=(.*?)location\.href\=location\.pathname\+location\.search<\/script\>`)
	jsRegDataS := jsReg.FindAllSubmatch(respData.RespBoydy, -1)
	var jsTmpData string
	jsData := func() string {
		for _, jsRegData := range jsRegDataS {
			jsTmpData = string(jsRegData[1])
		}
		return jsTmpData
	}()
	//jsData := jsRegData[len(jsRegData)-1]
	//执行相应的JS
	jsVale, _ := jsActuator(jsData)
	return TmpJslClearance{
		respData.StatusCode,
		cookieData,
		jsVale,
	}
}

//获取正式的token
func getJslClearanceToken(url string, tmpToken TmpJslClearance) (spiderToken string) {
	tmpCookie := tmpToken.JslUid + tmpToken.JsVale
	link := public.Parameter{
		URL:    url,
		Method: "GET",
		Headers: map[string][]string{
			"User-Agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4684.0 Safari/537.36"},
			"Cookie":     {tmpCookie},
		},
	}
	respData := public.HttpReq(link)
	jsReg := regexp.MustCompile(`;go\((.*?)\)</script>`) // 正则表达式的分组，以括号()表示，每一对括号就是我们匹配到的一个文本，可以把他们提取出来。
	jsRegDatas := jsReg.FindAllSubmatch(respData.RespBoydy, -1)
	var jsTmpData string
	jsData := func() string {
		for _, jsRegData := range jsRegDatas {
			jsTmpData = string(jsRegData[1])
		}
		return jsTmpData
	}()
	var data JslClearanceRespData
	err := json.Unmarshal([]byte(jsData), &data)
	if err != nil {
		log.Println(err.Error())
		return

	}
	token := JslClearanceRespData{
		data.Bts,
		data.Chars,
		data.Ct,
		data.Ha,
		data.Tn,
		data.Vt,
		data.Wt,
	}
	spiderTmpToken := tmpToken.JslUid + `__jsl_clearance_s=` + spiderTokenRule(token)
	spiderToken = spiderTmpToken
	return spiderToken
}
func Token(url string) string {
	tmpToken := getJslClearanceTmpToken(url)
	token := getJslClearanceToken(url, tmpToken)
	return token
}

```
