---
title: automerge
layout: post
category: node
author: 夏泽民
---
Automerge是一个用于在JavaScript中构建协作应用程序的数据结构库。

构建JavaScript应用程序的常用方法是将应用程序的状态保存在模型对象中，例如JSON文档。例如，假设您正在开发一个任务跟踪应用程序，其中每个任务都由一张卡片表示。在JavaScript中，你可以这样写：

var doc = {cards: []}

// User adds a card
doc.cards.push({title: 'Reticulate splines', done: false})

// User marks a task as done
doc.cards[0].done = true

// Save the document to disk
localStorage.setItem('MyToDoList', JSON.stringify(doc))
特点和设计原则

网络不可知论者 。Automerge是一个纯粹的数据结构库，不关心你使用的是什么类型的网络。

不变的状态 。Automerge对象在某个时间点是应用程序状态的不可变快照。无论何时进行更改，或者合并来自网络的更改，都会返回一个反映该更改的新状态对象。

自动合并 。Automerge是所谓的无冲突复制数据类型（CRDT），它允许在不需要任何中央服务器的情况下自动合并不同设备上的并发更改。

相当便捷 。已经在Node.js，Chrome，Firefox和Electron上测试了Automerge 。

建立
如果您在Node.js中，则可以通过npm安装Automerge：

$ npm install --save automerge
然后你可以require('automerge')导入它。

使用这个存储库，可以使用下面的命令：

npm install - 安装依赖关系。

npm test - 在Node中运行测试套件。

npm run browsertest - 在Web浏览器中运行测试套件。

npm run webpack- dist/automerge.js为Web浏览器创建一个捆绑的JS文件。
<!-- more -->
https://automerge.org/docs/quickstart/
https://github.com/automerge/automerge#usage
