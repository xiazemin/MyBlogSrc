---
title: gorm Preloading
layout: post
category: golang
author: 夏泽民
---
默认情况下GORM因为性能问题，不会自动加载关联属性的值，gorm通过Preload函数支持预加载（Eager loading）关联数据

// 用户表
type User struct {
  gorm.Model
  Username string
  Orders []Orders // 关联订单，一对多关联关系
}
// 订单表
type Orders struct {
  gorm.Model
  UserID uint // 外键字段 
  Price float64
}

// 预加载Orders字段值，Orders字段是User的关联字段
db.Preload("Orders").Find(&users)
// 下面是自动生成的SQL，自动完成关联查询
//// SELECT * FROM users;
//// SELECT * FROM orders WHERE user_id IN (1,2,3,4);

https://www.tizi365.com/archives/85.html
<!-- more -->
https://gorm.io/zh_CN/docs/preload.html

https://gorm.io/docs/preload.html


