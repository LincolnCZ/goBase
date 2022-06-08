package main

import "fmt"

// map 是引用类型，底层本质是指针，所以函数内对map的修改对原有值有效
func addOne(m map[string]int, key string, val int) {
	m[key] = val
}

type CheckInfo struct {
	noStreamCount int
	checkCount    int
}

func main() {
	//map --- 引用类型
	//1. map的定义、元素添加、判断key是否存在、删除元素
	// len()函数用于获取map中元素的个数，即有多个少key。delete()用于删除map中的某个key，即使键不在 map 中，delete 的操作也都是安全的。
	m1 := make(map[string]int)
	m1["key1"] = 1
	m1["key2"] = 2
	fmt.Println(m1) //map[key1:1 key2:2]

	//初始化赋值
	m2 := map[string]int{"test1": 1, "test2": 2}
	fmt.Println(m2)
	m3 := map[string]int{
		"test1": 1,
		"test2": 2, // 注意：结尾的逗号是需要的
	}
	fmt.Println(m3)

	//判断key是否存在
	value, ok := m1["key2"]
	if ok {
		fmt.Println("key2", value)
	}

	//map元素删除
	delete(m1, "key2")
	fmt.Println(m1) //map[key1:1]

	//map元素遍历
	for k, v := range m1 {
		fmt.Println(k, v)
	}

	//2. nil map 和 空map
	//nil map和empty map的关系，就像nil slice和empty slice一样，两者都是空对象，未存储任何数据，但前者不指向底层数据结构，
	//  后者指向底层数据结构，只不过指向的底层对象是空对象。
	//如果向nil map中添加元素，会core dump
	emptyMap := map[string]int{}
	println(emptyMap) // 0xc000110e18
	var nilMap map[string]int
	println(nilMap) //0x0

	//3. 传递map给函数
	m4 := make(map[string]int)
	m4["key1"] = 1
	fmt.Println(m4) // map[key1:1]
	addOne(m4, "key2", 2)
	fmt.Println(m4) // map[key1:1 key2:2]

	//4. 常用操作
	m5 := make(map[string]*CheckInfo)
	info, ok := m5["key1"]
	if ok {
		info.noStreamCount += 1
		info.checkCount += 1
	} else {
		info = &CheckInfo{
			noStreamCount: 1,
			checkCount:    1,
		}
	}
	fmt.Println(m5)
}
