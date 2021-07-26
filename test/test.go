package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

/**
api 测试
*/

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容

func httpGet(url string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

// 发送Post请求
// url：         请求地址
// response：    请求返回的内容

func httpPost(urlStr string) string {

	//codeType:= "1"
	//codeType:= "2"
	codeType := "3"
	resp, err := http.PostForm(urlStr,
		url.Values{
			"codeType":   {codeType},
			"drawCount":  {"19"},
			"des":        {"des 这是金币的礼品码"},
			"validTime":  {"16345634354545"},
			"content":    {`{"1":1000,"2":10000}`},
			"createUser": {"qq"},
			"userId":     {"123456"},
		})

	if err != nil {
		// handle error
		log.Println("--resp err")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Println("--ReadAll err")
	}

	fmt.Println(string(body))
	return string(body)
}

// 发送Post请求
// url：         请求地址
// response：    请求返回的内容

func httpPostCheck(urlStr, code, uid string) string {

	resp, err := http.PostForm(urlStr,
		url.Values{
			"code": {code},
			"uid":  {uid},
		})

	if err != nil {
		// handle error
		log.Println("--resp err")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Println("--ReadAll err")
	}

	fmt.Println(string(body))
	return string(body)
}

func main() {
	result := ""
	urlStr := ""
	// type 1-指定用户一次性消耗，2-不指定用户限制兑换次数，3-不限用户不限次数兑换
	// 创建礼品码
	//urlStr = "http://127.0.0.1:8000/createGift"
	//result = httpPost(urlStr)
	//log.Println("--结果", result)

	// 查询礼品码 // type=3 jG7a4lo8, type=2 90KKHauh,
	urlStr = "http://127.0.0.1:8000/getGift?code=jG7a4lo8"
	result = httpGet(urlStr)
	log.Println("--结果", result)
	////
	////// 验证礼品码
	//urlStr = "http://127.0.0.1:8000/checkCode"
	////type=3 jG7a4lo8, type=2 90KKHauh,type=1 C72uloHO uid=123456
	////code:="jG7a4lo8"  //
	////code:="C72uloHO"  //
	////uid:="1234567"
	//
	//code := "90KKHauh" //
	//uid := "777777"
	//result = httpPostCheck(urlStr, code, uid)
	//log.Println("--结果", result)
}
