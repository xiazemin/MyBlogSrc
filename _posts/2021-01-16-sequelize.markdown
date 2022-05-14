---
title: sequelize
layout: post
category: node
author: 夏泽民
---
https://sequelize.org/

https://github.com/sequelize/sequelize

Installation
$ npm i sequelize # This will install v6

# And one of the following:
$ npm i pg pg-hstore # Postgres
$ npm i mysql2
$ npm i mariadb
$ npm i sqlite3
$ npm i tedious # Microsoft SQL Server
<!-- more -->
https://github.com/demopark/sequelize-docs-Zh-CN

// 方法 2: 分别传递参数 (其它数据库)
const sequelize = new Sequelize('database', 'username', 'password', {
  host: 'localhost',
  dialect: /* 选择 'mysql' | 'mariadb' | 'postgres' | 'mssql' 其一 */
});

模型定义：
sequelize.define('model_name',{filed:value})

创建表：
首先定义模型： model
然后同步：model.sync()
创建表会自动创建主键，默认为 id

增删改查：
新增数据：model.create 相当于 build save两步合并；
批量新增：model.bulkCreate([model,...],{...}) ;
但是默认不会运行验证器，需要手动开启
 User.bulkCreate([
  { username: 'foo' },
  { username: 'bar', admin: true }
], { validate: true,//手动开启验证器
fields: ['username']//限制字段 
});
// 因为限制了字段只存username，foo 和 bar 都不会是管理员.
更新 model.update
相当于 set, save两步合并，通常就直接修改实例属性，然后save()更新；
部分更新:
通过传递一个列名数组,可以定义在调用 save 时应该保存哪些属性
save({fields:[ 'name',... ]}) 只更新数组里面的字段
删除 model.destroy
重载实例：model.reload
查询：
include参数 对应sql的 join连接操作
findAll 查找所有的
findByPk 根据主键查找
findOne 找到第一个实例
findOrCreate 查找到或创建实例
findAndCountAll 分页查找

https://github.com/sequelize/

https://www.sequelize.com.cn/

https://www.liaoxuefeng.com/wiki/1022910821149312/1101571555324224




