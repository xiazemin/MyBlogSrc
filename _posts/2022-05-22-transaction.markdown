---
title: Sequelize Transactions 
layout: post
category: node
author: 夏泽民
---
Sequelize 支持两种使用事务的方法：

一个将根据 promise 链的结果自动提交或回滚事务，（如果启用）用回调将该事务传递给所有调用
而另一个 leave committing，回滚并将事务传递给用户。
主要区别在于托管事务使用一个回调，对非托管事务而言期望 promise 返回一个 promise 的结果。
<!-- more -->

托管事务（auto-callback）
托管事务自动处理提交或回滚事务。你可以通过将回调传递给 sequelize.transaction 来启动托管事务。

注意回传传递给 transaction 的回调是否是一个 promise 链，并且没有明确地调用t.commit（）或 t.rollback()。 如果返回链中的所有 promise 都已成功解决，则事务被提交。 如果一个或几个 promise 被拒绝，事务将回滚。


return sequelize.transaction(function (t) {

  // 在这里链接您的所有查询。 确保你返回他们。
  return User.create({
    firstName: 'Abraham',
    lastName: 'Lincoln'
  }, {transaction: t}).then(function (user) {
    return user.setShooter({
      firstName: 'John',
      lastName: 'Boothe'
    }, {transaction: t});
  });

}).then(function (result) {
  // 事务已被提交
  // result 是 promise 链返回到事务回调的结果
}).catch(function (err) {
  // 事务已被回滚
  // err 是拒绝 promise 链返回到事务回调的错误
});



非托管事务（then-callback）
非托管事务强制您手动回滚或提交交易。 如果不这样做，事务将挂起，直到超时。 要启动非托管事务，请调用 sequelize.transaction() 而不用 callback（你仍然可以传递一个选项对象），并在返回的 promise 上调用 then。 请注意，commit() 和 rollback() 返回一个 promise。

return sequelize.transaction().then(function (t) {
  return User.create({
    firstName: 'Bart',
    lastName: 'Simpson'
  }, {transaction: t}).then(function (user) {
    return user.addSibling({
      firstName: 'Lisa',
      lastName: 'Simpson'
    }, {transaction: t});
  }).then(function () {
    return t.commit();
  }).catch(function (err) {
    return t.rollback();
  });
});

https://www.jianshu.com/p/5d63cdc103e4
