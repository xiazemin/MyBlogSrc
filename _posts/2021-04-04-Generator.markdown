---
title: Generator yield
layout: post
category: node
author: 夏泽民
---
注意Generator、yield；async、await要配套使用，且yield 的外层的函数一定要声明称Generator （*）
为了解决异步的嵌套问题，真是操碎了心，先是出了个Promise，然后又是Generator、yield组合，直到ES7的async、await组合。
Generator
生成器对象是由function* 返回的，并且符合可迭代协议和迭代器协议。
这里有几个概念生成器、可迭代协议、迭代器协议。具体的概念可以点击链接查看MDN文档。

function*: 定义一个生成器函数，返回一个Generator对象；
可迭代协议： 允许 JavaScript 对象去定义或定制它们的迭代行为；
迭代器协议： 定义了一种标准的方式来产生一个有限或无限序列的值；当一个对象被认为是一个迭代器时，它实现了一个 next() 的方法，next()返回值如下：

{
 done:true,//false迭代是否结束，
 value:v,//迭代器返回值
}
从这几个基本的概念我们可以了解到，生成器是对象是可以迭代的，那么为什么要可以迭代、可以迭代解决了什么问题。

迭代
下面定义一个简单的迭代生成函数，传入一个数组，则返回一个可以迭代的对象

// 1. 迭代器

let iterator = (items)=>{
  let iter = {
    index:0,
    max:items.length,
    next:function(){ // 返回调用结果
      return this.index === this.max ? {value:undefined,done:true} : {value:items[this.index++],done:false};
    }
  }

  return iter;
}

export default iterator;
调用上面的迭代器，并执行

let iter = iterator([1,2,3,4]);
let result = null;
console.log('``````iterator````````');
do{
  result = iter.next();
  console.log(result);
}while (!result.done)
运行结果如下：

1.JPG
可以看到，迭代器每次调用next()方法，都会返回{value:xx,done:xx}结构的对象，这个就是迭代器协议中next()方法需要遵循的规则，前面说过generator函数也是遵循迭代器协议的，下面用generator实现此功能。

generator的使用
// generator
function *generator(items){
  let index = 0;
  let max = items.length;

  while (index < max){
    yield items[index++];
  }

}

let gene = generator([1,2,3,4]);
result = null;
console.log('``````````generator`````````');
do{
  result = gene.next();
  console.log(result)
}while(!result.done)
此时运行结果如下：


2.JPG
对比两次运行的结果，得出一个结论：生成器(function*)函数，运行时，返回的是一个生成器对象，这个生成器对象是可以迭代(gene.next())的，并且next()的返回值包含value，done两个字段。

进化
生成器是可以迭代的，而且返回值也是符合一定结构的，我们每次再使用生成器的时候，都要用循环去执行，知道返回的done为true，为了简化操作需要把这个循环操作进行封装，下面封装一个简单的run函数，run可以执行迭代器，一直到完成任务

let tick = (duration)=>{
  return new Promise((resolve)=>{
    setTimeout(function () {
      console.log(duration,new Date());
      resolve(duration);
    },duration);
  });
};

function *generator() {
  var result = yield tick(2000);
  console.log('result = ',result);
  result = yield tick(4000);
  console.log('result = ',result);
  result = yield tick(3000);
  console.log('result = ',result);
}


let run = (generator,res)=>{
  var result = generator.next(res);
  if(result.done) return;
  result.value.then((res)=>{
    run(generator,res);
  });
}

run(generator());

以上的运行结果：

3.JPG
看一下run的实现，像极了前面的do...while... 循环，只是做了一个简单的封装，以后就没用每次都手写循环来执行生成器函数了，实际上有一个封装好的库可以使用它叫co

co库执行generator
安装co

npm install --save co
使用

import co from 'co';
co(generator);
运行结果如下

4.JPG
它的作用跟上面实现的run方法的作用是一样的，都是执行generator，并返回结果。这样生成器大概就可以理解了，说白了生成器就是可以返回一个可迭代的对象，这个对象不是通过return返回的，而是通过yield，并且可以实现异步函数的同步调用，我们看上图的时间，虽然tick是异步的，但是打印的结果却是顺序执行的。

async/await
generator可以简化异步的编码，减少嵌套，而async、await组合起来使用，可以更进一步，类似以上的代码，使用async、await改写如下

let tick = (duration)=>{
  return new Promise((resolve)=>{
    setTimeout(function () {
      console.log(new Date());
      resolve(duration);
    },duration);
  });
}


async function asyncFunc(){
  var result = await tick(1000);
  console.log(result);
  result = await tick(2000);
  console.log(result);
  result = await tick(3000);
  console.log(result);
}

asyncFunc();
执行结果

5.JPG
虽然实现的功能是一样的，但是从代码的结构上又简化了一层。

https://www.jianshu.com/p/c94edc0057fe
<!-- more -->
http://www.ruanyifeng.com/blog/2015/04/generator.html

如果依次读取多个文件，就会出现多重嵌套。代码不是纵向发展，而是横向发展，很快就会乱成一团，无法管理。这种情况就称为"回调函数噩梦"（callback hell）。

Promise就是为了解决这个问题而提出的。它不是新的语法功能，而是一种新的写法，允许将回调函数的横向加载，改成纵向加载。采用Promise，连续读取多个文件，写法如下。


var readFile = require('fs-readfile-promise');

readFile(fileA)
.then(function(data){
  console.log(data.toString());
})
.then(function(){
  return readFile(fileB);
})
.then(function(data){
  console.log(data.toString());
})
.catch(function(err) {
  console.log(err);
});


