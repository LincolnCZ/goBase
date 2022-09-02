package main

import (
	"fmt"
	"strconv"
	"strings"
)

func stringExample() {
	//string --- 值类型
	// 1. Go 中的 string 类型要使用双引号或反引号包围，它们的区别是：
	//* 双引号是弱引用，其内可以使用反斜线转义符号，如 ab\ncd 表示 ab 后换行加 cd
	//* 反引号是强引用，其内任何符号都被强制解释为字面意义，包括字面的换行。也就是所谓的裸字符串。例如：`ABC\nDEF`
	//注意：
	//* 使用单引号包围的字符实际上是整数数值。例如'a'等价于97。
	println("abc\ndef")
	println(`ABC\nDEF`) //输出：ABC\nDEF
	println('A')        //输出：65

	//2. string 的底层是 byte 数组，每个 string 其实只占用两个机器字长：一个指针和一个长度。只不过这个指针在 Go 中完全不可见，所以对我们来说，
	//    string 是一个底层 byte 数组的值类型而非指针类型。
	//所以，可以将一个 string 使用 copy() 拷贝到一个给定的 byte slice 中，也可以使用 slice 的切片功能截取 string 中的片段。
	var a = "hello world!"
	//使用 slice 的切片功能、len 函数获取长度
	println(a[2:4], len(a)) //ll 12
	s1 := make([]byte, 12)
	//使用 copy 函数，将一个 string 拷贝到一个给定的 byte slice 中
	copy(s1, a)
	println(string(s1)) // hello world!

	//3. 遍历
	//字符串是字符数组，如果字符串中全是 ASCII 字符，直接遍历即可，但如果包含了多字节字符，则可以 []rune(str) 转换后后再遍历。
	s2 := "Hello 你好"
	r := []rune(s2) // 8
	for i := 0; i < len(r); i++ {
		fmt.Printf("%c", r[i]) //Hello 你好
	}

	//4.修改字符串
	//字符串是一个不可变对象，所以对字符串 s 截取后赋值的方式 s[1]="c" 会报错。
	//要想修改字符串中的字符，必须先将字符串拷贝到一个 byte slice 中，然后修改指定索引位置的字符，最后将 byte slice 转换回 string 类型。
	s3 := "hello"
	bs := []byte(s3)
	bs[0] = 'H' // 必须使用单引号
	s3 = string(bs)
	println(s3) //Hello

	//5. 字符串串接
	//使用加号 + 连接两段字符串，字符串连接 + 操作符强制认为它两边的都是 string 类型，所以 "abc" + 2 将报错。需要先将 int 类型的 2 转换为字符串
	//   类型(不能使用 string(2) 的方式转换，因为这种转换方式不能跨大类型转换，只能使用 strconv 包中的函数转换)。
	//另一种更高效的字符串串接方式是使用 strings 包中的 Join() 函数，它可以在缓冲中将字符串串接起来。
	s4 := "first" + "," + "second"
	println(s4)

	s5 := []string{"foo", "bar", "baz"}
	fmt.Println(strings.Join(s5, ", ")) // foo, bar, baz

	//6.字符串和数字的互相转换
	//• 要将整数转换成字符串，一种选择是使用 fmt.Sprintf，另一种选择是 strconv.Itoa 函数。
	//• strconv 包内的 Atoi 函数或 ParseInt 函数用于解释表示整数的字符串。

	//（1）将整型转换成字符串
	x := 123
	y := fmt.Sprintf("%d", x)
	fmt.Println(y, strconv.Itoa(123))

	// 可以按不同的进位制格式化数字
	fmt.Println(strconv.FormatInt(int64(x), 2)) // 1111011
	s := fmt.Sprintf("x=%b", x)
	fmt.Println(s) // x=1111011

	//（2）将字符串转成整型
	z, _ := strconv.Atoi("123")             // z 是整型
	h, _ := strconv.ParseInt("123", 10, 64) // 十进制，最长为64位
	fmt.Println(z, h)

}
