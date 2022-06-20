package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Request 请求
type Request struct {
	/*
	   （1）与编码类似，解码的时候golang会将json的数据结构和go的数据结构进行匹配。
	   匹配的原则就是寻找tag的相同的字段，然后查找字段。查询的时候是大小写不敏感的（大小写
	   不敏感只是针对公有字段而言，对于私有的字段，即使tag符合，也不会被解析）。

	   对于无法匹配到的字段，会默认设置其零值。
	*/
	Callback   string `json:"callback"`
	RequestId  string `json:"requestId"`
	privateKey string `json:"privateKey"` //小写字母开头，私有字段，即使tag匹配，也不会解析
	PublicKey  string `json:"publicKey"`  //大写字母开头，没有tag，依然可以匹配到

	/*
	   （2）`-`tag。与编码一样，tag的-也不会被解析，但是会初始化其零值。
	*/
	DropMsg string `json:"-"`

	/*
	   （3）omitempty tag。对于解码来说，omitempty tag可以忽略。对于无法匹配到的字段，会默认设置其零值。
	*/
	Actions []string    `json:"actions,omitempty"`
	Data    []*DataList `json:"data,omitempty"`

	/*
	   （4）string tag。在解码的时候，只有字串类型的数字，才能被正确解析，或者会报错。
	*/
	Str2float float64 `json:"str2float"`

	/*
	   （5）RawMessage is a raw encoded JSON value. 需要透传的字段
	*/
	Context json.RawMessage `json:"context"`

	/*
	   （6）interface{}的使用
	*/
}

// DataList 单个请求对象
type DataList struct {
	DataId   string `json:"dataId"`
	DataType string `json:"dataType"`
}

func main() {
	var jsonString = `{
		"callback":"",
		"requestId":"request_id",
		"publicKey":"public_key",
		"actions":["porn","act"],
		"data":[{"dataId":"data_id","dataType":"data_type"}],
		"str2float":0,
		"context":{"precomputed":true,"test":1}
	}`
	request := Request{}

	err := json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", request)
}

//输出：
//{Callback: RequestId:request_id privateKey: PublicKey:public_key DropMsg: Actions:[porn act] Data:[0xc0000a6280]
//Str2float:0 Context:[123 34 112 114 101 99 111 109 112 117 116 101 100 34 58 116 114 117 101 44 34 116 101 115 116 34 58 49 125]}
