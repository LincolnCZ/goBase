package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Request1 请求
type Request1 struct {
	RequestId json.RawMessage `json:"requestId"`
	Id        json.RawMessage `json:"id"`
}

func jsonRawExample() {
	var jsonString = `{
		"requestId":"request_id",
		"id":1
	}`
	request := Request1{}

	err := json.Unmarshal([]byte(jsonString), &request)
	if err != nil {
		log.Fatalln(err)
	}
	// (1) 使用json.JsonRawMessage 解析出的字段都是 []byte 类型的
	fmt.Printf("%+v\n", request) // 输出 {RequestId:[34 114 101 113 117 101 115 116 95 105 100 34] Id:[49]}

	// (2) 错误使用方式：对于原本是string类型的，如果直接强制转换成string，则输出结果中包含引号
	fmt.Println(string(request.RequestId)) // 输出 "request_id"，注意包含引号。
	fmt.Println(request.Id)                // 输出：[49]

	// (3) 正确使用方式：定义对应的类型，并且使用json.Unmarshal
	var str string
	json.Unmarshal(request.RequestId, &str)
	fmt.Println(str) // 输出 request_id
	var id int
	json.Unmarshal(request.Id, &id)
	fmt.Println(id) // 输出 1
}
