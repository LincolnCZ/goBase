package main

import "fmt"

// Go中的Integer有以下几种细分的类型：
//• int8,int16,int32,int64
//• uint8,uint16,uint32,uint64
//• byte
//• rune
//• int,uint
//
//其中8 16 32 64表示该数据类型能存储的bit位数。例如int8表示能存储8位数值，所以这个类型占用1字节，也表示最大能存储的整型数共2^8=256个，
//   所以int8类型允许的最大正数为127，允许的最小负数为-128，共256个数值。
//uint中的u表示unsigned，即无符号整数，只保存0和正数。所以uint8能存储256个数的时候，允许的最小值为0，允许的最大值为255。
//额外的两种Integer是byte和rune，它们分别等价于uint8(即一个字节大小的正数)、int32。
//两种依赖于CPU位数的类型int和uint，它们分别表示一个机器字长。在32位CPU上，一个机器字长为32bit，共4字节，在64位CPU上，一个机器字长
//   为64bit，共8字节。除了int和uint依赖于CPU架构，还有一种uintptr也是依赖于机器字长的。

func integerExample() {

}

func byteExample() {
	//Go中没有专门提供字符类型char，Go内部的所有字符类型(无论是ASCII字符还是其它多字节字符)都使用整数值保存，所以字符可以存放到byte、int等
	//   数据类型变量中。byte类型等价于uint8类型，表示无符号的1字节整数。
	var a byte = 'A' // a=65
	println(a)
	var b uint8 = 'a' // b=97
	println(b)
}

func runeExample() {
	//Go语言同样支持 Unicode（UTF-8），因此字符同样称为 Unicode 代码点或者 runes，并在内存中使用 int 来表示。在文档中，一般使用
	//   格式 U+hhhh 来表示，其中 h 表示一个 16 进制数。
	//在书写 Unicode 字符时，需要在 16 进制数之前加上前缀\u或者\U。因为 Unicode 至少占用 2 个字节，所以我们使用 int16 或者 int 类型来表示。
	//   如果需要使用到 4 字节，则使用\u前缀，如果需要使用到 8 个字节，则使用\U前缀。
	var ch int = '\u0041'
	var ch2 int = '\u03B2'
	var ch3 int = '\U00101234'
	fmt.Printf("%d - %d - %d\n", ch, ch2, ch3) // integer
	fmt.Printf("%c - %c - %c\n", ch, ch2, ch3) // character
	fmt.Printf("%X - %X - %X\n", ch, ch2, ch3) // UTF-8 bytes
	fmt.Printf("%U - %U - %U\n", ch, ch2, ch3) // UTF-8 code point
}

func boolExample() {
}
