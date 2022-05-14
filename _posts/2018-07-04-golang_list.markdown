---
title: golang list jsonMarshal之后一个为null一个为[ ]问题
layout: post
category: golang
author: 夏泽民
---
<!-- more -->
两种不同的初始化方式引起的，函数返回值和var方式都是“声明”，而不是“定义”。fmt的输出和jsonMarshal也不一样

package main

import(
"fmt"
"reflect"
"encoding/json"
)

func main() {
	 var x []string
    fmt.Println(x, reflect.TypeOf(x), len(x), cap(x), x == nil)
    x1 := []string{}
    fmt.Println(x1, reflect.TypeOf(x1), len(x1), cap(x1), x1 == nil)
    var x2 = make([]string, 0)
    fmt.Println(x2, reflect.TypeOf(x2), len(x2), cap(x2), x2 == nil)
    
	b1, _ := json.Marshal(x)
    fmt.Println("x: ", string(b1))
    b2, _ := json.Marshal(x1)
    fmt.Println("x1: ", string(b2))
	  b3, _ := json.Marshal(x2)
    fmt.Println("x1: ", string(b3))
}

[] []string 0 0 true
[] []string 0 0 false
[] []string 0 0 false
x:  null
x1:  []
x1:  []
