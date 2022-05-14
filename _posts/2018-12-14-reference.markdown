---
title: golang 传值传引用
layout: post
category: golang
author: 夏泽民
---
golang中大多数是传值的,有：
基本类型:byte,int,bool,string
复合类型:数组,数组切片,结构体,map,channnel
在函数中参数的传递可以是传值（对象的复制,需要开辟新的空间来存储该新对象）和传引用（指针的复制，和原来的指针指向同一个对象），建议使用指针，原因有两个：能够改变参数的值，避免大对象的复制操作节省内存。struct和数组的用法类似
channel和数组切片，map一样，传参的方式是传值，都可以直接使用，其内部维护着指向真正存储空间的指针。

m = map[value:0]
m1 = map[value:0]
m = map[value:1]
m1 = map[value:1]
我们发现，当修改了m1，m也随着改变了，这看似是传引用，但其实map也是传值的，它的原理和数组切片类似。map内部维护着一个指针，该指针指向真正的map存储空间。我们可以将map描述为如下结构：
type map[key]value struct{
	impl *Map_K_V
}
type Map_K_V struct{
	//......
}
其实，map和slice,channel一样，内部都有一个指向真正存储空间的指针，所以，即使传参时是对值的复制（传值），但都指向同一块存储空间。
<!-- more -->
package main

import "fmt"

func main() {
	fmt.Println("Hello, 世界")
	a:=map[string]string{
	"a":"1",
	"b":"2",
	}
	fmt.Println(a)
	b:=a
	b["a"]="3"
	fmt.Println(b)
	fmt.Println(a)
	changeMap(a)
	fmt.Println(a)
	fmt.Println(b)
	changeMapPoint(&b)
	fmt.Println(a)
	fmt.Println(b)
	c:="123"
	changeString(c)
	fmt.Println(c)
	changeStringPtr(&c)
	fmt.Println(c)
	d:=[]string{"1","2","3"}
	changeSlice(d)
	fmt.Println(d)
	changeSlicePtr(&d)
	fmt.Println(d)
	var array = [3]int{0, 1, 2}
	var array2 = array
	array2[2] = 5
	fmt.Println(array, array2)
	var array3 = [3]int{0, 1, 2}
	var array4 = &array3
	array4[2] = 5
	fmt.Println(array3, *array4)
}

func changeSlicePtr(s* []string){
(*s)[1]="b"
}
func changeSlice(s []string){
s[1]="a"
}
func changeMap(m map[string]string){
m["b"]="4"
}

func changeMapPoint(m *map[string]string){
(*m)["b"]="5"
}

func changeString(s string){
s="abc"
}

func changeStringPtr(s *string){
*s="abcd"
}

Hello, 世界
map[a:1 b:2]
map[b:2 a:3]
map[a:3 b:2]
map[a:3 b:4]
map[a:3 b:4]
map[a:3 b:5]
map[a:3 b:5]
123
abcd
[1 a 3]
[1 b 3]
[0 1 2] [0 1 5]
[0 1 5] [0 1 5]
