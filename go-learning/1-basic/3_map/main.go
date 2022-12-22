package main

import "fmt"

// map 是引用类型，底层本质是指针，所以函数内对 map 的修改对原有值有效
func addOne(m map[string]int, key string, val int) {
	m[key] = val
}

type CheckInfo struct {
	noStreamCount int
	checkCount    int
}

func main() {
	//map --- 引用类型
	//1. map的定义、元素添加、判断 key 是否存在、删除元素
	// len() 函数用于获取 map 中元素的个数，即有多个少 key。delete() 用于删除 map 中的某个 key，即使键不在 map 中，delete 的操作也都是安全的。
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

	//map 元素删除
	delete(m1, "key2")
	fmt.Println(m1) //map[key1:1]

	//map 元素遍历
	for k, v := range m1 {
		fmt.Println(k, v)
	}

	//2. nil map 和 空map
	//nil map 和 empty map 的关系，就像 nil slice 和 empty slice 一样，两者都是空对象，未存储任何数据，但前者不指向底层数据结构，
	//  后者指向底层数据结构，只不过指向的底层对象是空对象。
	//如果向 nil map 中添加元素，会 core dump
	emptyMap := map[string]int{}
	println(emptyMap) // 0xc000110e18
	var nilMap map[string]int
	println(nilMap) //0x0

	//3. 传递 map 给函数
	m4 := make(map[string]int)
	m4["key1"] = 1
	fmt.Println(m4) // map[key1:1]
	addOne(m4, "key2", 2)
	fmt.Println(m4) // map[key1:1 key2:2]

	//4. 常用操作
	m5 := make(map[string]*CheckInfo)
	m5["key2"] = &CheckInfo{
		noStreamCount: 2,
		checkCount:    2,
	}

	var key1 = "key1"
	if info, ok := m5[key1]; ok {
		info.noStreamCount += 1 // 由于info 是指针，所以更改有效
		info.checkCount += 1
	} else {
		info = &CheckInfo{
			noStreamCount: 1,
			checkCount:    1,
		}
		m5[key1] = info // 需要显式添加
	}

	if info, ok := m5["key2"]; ok {
		info.noStreamCount += 1
		info.checkCount += 1
	}

	fmt.Println("m5:", m5)
	for k, v := range m5 {
		fmt.Println(k, v.checkCount, v.noStreamCount)
	}
}
