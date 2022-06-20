## 通过类型嵌入实现“继承”

类型嵌入指的就是在一个类型的定义中嵌入了其他类型。Go 语言支持两种类型嵌入，分别是接口类型的类型嵌入和结构体类型的类型嵌入。
* 接口类型的类型嵌入 
* 结构体类型的类型嵌入 
  * 这种以某个类型名、类型的指针类型名或接口类型名，直接作为结构体字段的方式就叫做结构体的类型嵌入，这些字段也被叫做嵌入字段（Embedded Field）。

## 说明
* interfaceEmbed：interface中嵌入interface
* structEmbed：struct中嵌入struct、interface



