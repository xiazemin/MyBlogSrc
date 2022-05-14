---
title: 获取map[interface{}]interface{}的value
layout: post
category: golang
author: 夏泽民
---
定义了一个map[interface{}]interface{} key 也是interface{}类型的，通过反射等动态取得的，这个key一定是存在在map里的。

但是调用map[key]却无法获取到对应的值。

keyValueMap := make(map[interface{}]interface{})

key interface{}

value := keyValueMap[key]

value始终为nil,实际上key是存在Map中的，只是都为interface{}类型所以获取不到。
<!-- more -->
package main

import (
	"fmt"
)

func main() {

	mapInterface := make(map[interface{}]interface{})
	mapString := make(map[string]string)

	mapInterface["k1"] = 1
	mapInterface[3] = "hello"
	mapInterface["world"] = 1.05
	mapInterface["rt"] = true

	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}

	fmt.Printf("%#v", mapString)
}


map[string]string{"3":"hello", "k1":"1", "rt":"true", "world":"1.05"}
我们可以看到，不管int类型，bool类型，都可以传值到map中。

 我们最后只需要把他装换成我们需要的类型就可以了
 
 var key interface{}
value := mapInterface[key]
fmt.Println(value)
key=3
value1 := mapInterface[key]
fmt.Println(value1)
key = "3"
value2 := mapInterface[key]
fmt.Println(value2)

<nil>
hello
<nil>

取的时候类型和值必须和存的时候一致
