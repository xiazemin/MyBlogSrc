---
title: golang单例模式
layout: post
category: golang
author: 夏泽民
---
1、定义：单例对象的类必须保证只有一个实例存在，全局有唯一接口访问。

2、分类：　　

懒汉方式：指全局的单例实例在第一次被使用时构建。
饿汉方式：指全局的单例实例在类装载时构建。
<!-- more -->
3、实现：

 （1）懒汉方式　　

1 type singleton struct{}
2 var ins *singleton
3 func GetIns() *singleton{
4     if ins == nil {
5     　　ins = &singleton{}
6     }
7     return ins
8 }
　　缺点：非线程安全。当正在创建时，有线程来访问此时ins = nil就会再创建，单例类就会有多个实例了。

（2）饿汉方式

　　

1 type singleton struct{}
2 var ins *singleton = &singleton{}
3 func GetIns() *singleton{
4     return ins
5 }
　　缺点：如果singleton创建初始化比较复杂耗时时，加载时间会延长。

（3）懒汉加锁

　　

 1 type singleton struct{}
 2 var ins *singleton
 3 var mu sync.Mutex
 4 func GetIns() *singleton{
 5     mu.Lock()
 6     defer mu.Unlock()
 7 
 8     if ins == nil {
 9     　　ins = &singleton{}
10     }
11     return ins
12 }
　　缺点：虽然解决并发的问题，但每次加锁是要付出代价的

（4）双重锁

 1  type singleton struct{}
 2  var ins *singleton
 3  var mu sync.Mutex
 4  func GetIns() *singleton{  
 5  　　if ins == nil {
 6     　　mu.Lock()
 7        defer mu.Unlock()
 8        if ins == nil {
 9       　　ins = &singleton{}
10        }
11     }
12     return ins
13 }
　　避免了每次加锁，提高代码效率

（5）sync.Once实现

1 type singleton struct{}
2 var ins *singleton
3 var once sync.Once
4 func GetIns() *singleton {
5     once.Do(func(){
6         ins = &singleton{}
7     })
8     return ins
9 }