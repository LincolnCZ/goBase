package main

import (
	"encoding/json"
	"fmt"
	"log"
)

//Marshal()和MarshalIndent()函数可以将数据封装成json数据：
//• struct、slice、array、map都可以转换成json
//• struct转换成json的时候，只有字段首字母大写的才会被转换
//• map转换的时候，key必须为string
//• 封装的时候，如果是指针，会追踪指针指向的对象进行封装

// Response 请求返回结果
type Response struct {
	/*
		（1）golang在定义结构的时候，只有使用大写字母开头的字段才会被导出。而通常json世界中，
		  更盛行小写字母的方式。golang提供了struct tag的方式可以重命名结构字段的输出形式。

		  下面定义的几个字段都是需要输出的。对于需要输出的值，而又没有手动赋值，则默认输出对应的零值。
	*/
	Code      int    `json:"code"`
	Message   string `json:"message"`
	TraceId   string `json:"traceId"`
	RequestId string `json:"requestId"`
	Timestamp int64  `json:"timestamp"`
	/*
		（2）实际开发中，我们需要某个字段，但是不希望编码到json中，可以使用`-`忽略字段。
	*/
	DropMsg string `json:"-"`
	/*
	   （3）omitempty可选字段。当其有值的时候就输出，而没有值(或者为零值)的时候就不输出。
	   * int零值为0，
	   * string零值为"",
	   * 指针零值为nil。
	   * slice零值：var tmp []int，如果声明为tmp := make([]int,3,3)则为非空的。
	*/
	OmitStr   string          `json:"omitStr,omitempty"`
	OmitInt   int             `json:"omitInt,omitempty"`
	OmitPtr   *int            `json:"omitPtr,omitempty"`
	OmitSlice []int           `json:"omitSlice,omitempty"`
	Data      []*ResultHeader `json:"data,omitempty"`
	/*
	   （4）有时候输出的json希望是数字的字符串，而定义的字段是数字类型，那么就可以使用string选项。
	*/
	Float2Str float32 `json:"float2str,string"`

	/*
	   （5）RawMessage is a raw encoded JSON value.
	   It implements Marshaler and Unmarshaler and can
	   be used to delay JSON decoding or precompute a JSON encoding.

	   有时候需要透传一些别人已经Marshall过的json串。
	*/
	Context json.RawMessage `json:"context,omitempty"`

	/*
	   （6）interface{}的使用
	   golang的数组、切片、map，其value的类型是一样的，如果遇到不同数据类型，则需可以借助interface{}来实现。
	   当interface{}没有初始化其值的时候，零值是 nil。编码成json就是 null。
	*/
	Extra    []interface{}          `json:"extra,omitempty"`
	ExtraMap map[string]interface{} `json:"extraMap,omitempty"`
}

// ResultHeader 返回单个对象信息
type ResultHeader struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Results []json.RawMessage `json:"results,omitempty"`
}

func main() {
	tmp := ResultHeader{
		Code:    0,
		Message: "ok",
	}
	context, _ := json.Marshal(tmp)

	response := Response{
		// 1)
		//Code:      0,   //没有赋值的情况下，最终给出的值。"code":0
		Message:   "ok",
		TraceId:   "dd057796-bc8a-41cf-8cb7-2a9927b2324a",
		RequestId: "51d245c0f4d6fc70418c5a99",
		Timestamp: 0,
		// 2）实际开发中，我们需要某个字段，但是不希望编码到json中，可以使用`-`忽略字段。
		DropMsg: "drop msg",
		// 3) omitempty可选字段说明。Data和OmiStr都是其零值，默认不会输出。
		Data:    nil,
		OmitStr: "",
		// 4) 数字类型转成string
		Float2Str: 1.234, // "float2str": "1.234"
		// 5) RawMessage 可以将任意结构 Marshal之后的结果，赋值给json.RawMessage
		//Context: json.RawMessage(`{"precomputed": true, "test":1}`),
		Context: context,
		// 6) interface{}的使用
		Extra: []interface{}{123, "hello"},
		ExtraMap: map[string]interface{}{
			"int":    1,
			"string": "hello world",
		},
	}

	//rs, err := json.Marshal(response)
	rs, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(rs))
}

//输出结果：
//{
//        "code": 0,
//        "message": "ok",
//        "traceId": "dd057796-bc8a-41cf-8cb7-2a9927b2324a",
//        "requestId": "51d245c0f4d6fc70418c5a99",
//        "timestamp": 0,
//        "float2str": "1.234",
//        "context": {
//                "code": 0,
//                "message": "ok"
//        },
//        "extra": [
//                123,
//                "hello"
//        ],
//        "extraMap": {
//                "int": 1,
//                "string": "hello world"
//        }
//}
