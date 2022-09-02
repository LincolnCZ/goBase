package main

import "fmt"

// Go 中的 Integer 有以下几种细分的类型：
//• int8,int16,int32,int64
//• uint8,uint16,uint32,uint64
//• byte
//• rune
//• int,uint
//
//其中 8 16 32 64 表示该数据类型能存储的 bit 位数。例如 int8 表示能存储 8 位数值，所以这个类型占用 1 字节，也表示最大能存储的整型数共 2^8=256 个，
//   所以 int8 类型允许的最大正数为 127，允许的最小负数为 -128，共 256 个数值。
//uint 中的 u 表示 unsigned，即无符号整数，只保存 0 和正数。所以 uint8 能存储 256 个数的时候，允许的最小值为 0，允许的最大值为 255。
//额外的两种 Integer 是 byte 和 rune，它们分别等价于 uint8(即一个字节大小的正数)、int32。
//两种依赖于 CPU 位数的类型 int 和 uint，它们分别表示一个机器字长。在 32 位 CPU上，一个机器字长为 32bit，共 4 字节，在 64 位 CPU 上，一个机器字长
//   为 64bit，共 8 字节。除了 int 和 uint 依赖于 CPU 架构，还有一种 uintptr 也是依赖于机器字长的。

func integerExample() {

}

func byteExample() {
	//Go 中没有专门提供字符类型 char，Go 内部的所有字符类型(无论是 ASCII 字符还是其它多字节字符)都使用整数值保存，所以字符可以存放到 byte、int 等
	//   数据类型变量中。byte 类型等价于 uint8 类型，表示无符号的 1 字节整数。
	var a byte = 'A' // a=65
	println(a)
	var b uint8 = 'a' // b=97
	println(b)
}

func runeExample() {
	//Go 语言同样支持 Unicode（UTF-8），因此字符同样称为 Unicode 代码点或者 runes，并在内存中使用 int 来表示。在文档中，一般使用
	//   格式 U+hhhh 来表示，其中 h 表示一个 16 进制数。
	//在书写 Unicode 字符时，需要在 16 进制数之前加上前缀 \u 或者 \U。因为 Unicode 至少占用 2 个字节，所以我们使用 int16 或者 int 类型来表示。
	//   如果需要使用到 4 字节，则使用 \u 前缀，如果需要使用到 8 个字节，则使用 \U 前缀。
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
