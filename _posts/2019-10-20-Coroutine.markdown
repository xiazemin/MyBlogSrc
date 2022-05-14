---
title: php Coroutine
layout: post
category: lang
author: 夏泽民
---
PHP5.5中加入了一个新特性—迭代生成器和协程。

一. 什么是协程（Coroutine）？
在协程出现之前，要实现多任务并发，在无OS（操作系统）时代，可以使用状态机的思想对多任务进行拆解，在单进程环境中运行多任务，但是这种模式下需要开发者对每个任务有清晰的了解，也要开发者自行开发与任务相关功能（如任务间的通讯）。

后来出现了OS（操作系统），咱们就开始使用OS提供的进程和线程功能来轻易实现多任务了。在OS中，进程的上下文切换是OS内核控制。但是后来却出现了一个问题，频繁的进程上下文切换导致了OS性能的降低（主要是短时执行消耗小的任务进程）。

为了解决这个问题，开始提出新的概念，就是在同一进程或线程中运行多个任务，这种问题就相当于回到了早期的无OS时代的多任务实现。而现在解决方案称为协程。其本质是，将将任务切换的部分工作从内核转移到应用层。
<!-- more -->
二. php中协程的基本工具以及基本使用
要实现协程，php给出了两个新东西：生成器和yield关键字。

什么是生成器？
生成器继承了实现了迭代器，在php代码中和函数的定义类似，不过内部使用了yield关键字，如：
<?php
function gen(){
    echo "hello gen".PHP_EOL;//step1
    $ret = (yield "gen1");   //step2
    var_dump($ret);  //step3
    $ret = (yield "gen2");   //step4
    var_dump($ret);  //step5
}
?>

使用时,这样子：

<?php
$my_gen = gen();
var_dump($my_gen->current());
var_dump($my_gen->send("main send"));
?>
1
2
3
4
5
好了，这样使用代表什么意思呢？
（1）首先$my_gen = gen();这句代码只是实例化一个新的生成器，里面的代码并未执行；
（2）\$my_gen->current()；这句代码就执行了生成器里面的step2中的yield “gen1”了，这时代码中断，并且字符串“gen1”被传进了生成器\$my_gen，并且作为current()函数的返回值；
（3）send(“main send”)执行完之后，字符串”main send”被传递进了生成器\$my_gen, 同时生成器作为step2中yield的返回值传递给ret;
（4） 生成器step3执行完后，在step4时，遇到yield就会再次进入中断。

三. 协程的特点
（1）为应用层实现多任务提供了工具;
（2）协程不允许多任务同时执行，要执行其它协程，必须使用关键字yield主动放弃cpu控制权;
（3）协程需要自己写任务管理器，以及任务调度器；
（4）减轻了OS处理零散任务和轻量级任务的负担；


http://www.laruence.com/2015/05/28/3038.html

http://nikic.github.io/2012/12/22/Cooperative-multitasking-using-coroutines-in-PHP.html

https://www.oschina.net/translate/cooperative-multitasking-using-coroutines-in-php

概念理解
到这里，你应该已经大概理解什么是生成器了。下面我们来说下生成器原理。

首先明确一个概念：生成器yield关键字不是返回值，他的专业术语叫产出值，只是生成一个值

那么代码中 foreach 循环的是什么？其实是PHP在使用生成器的时候，会返回一个 Generator 类的对象。 foreach 可以对该对象进行迭代，每一次迭代，PHP会通过 Generator 实例计算出下一次需要迭代的值。这样 foreach 就知道下一次需要迭代的值了。

而且，在运行中 for 循环执行后，会立即停止。等待 foreach 下次循环时候再次和  for  索要下次的值的时候，循环才会再执行一次，然后立即再次停止。直到不满足条件不执行结束。

实际开发应用
很多PHP开发者不了解生成器，其实主要是不了解应用领域。那么，生成器在实际开发中有哪些应用？

读取超大文件
PHP开发很多时候都要读取大文件，比如csv文件、text文件，或者一些日志文件。这些文件如果很大，比如5个G。这时，直接一次性把所有的内容读取到内存中计算不太现实。

这里生成器就可以派上用场啦。

<?php
header("content-type:text/html;charset=utf-8");
function readTxt()
{
    # code...
    $handle = fopen("./test.txt", 'rb');

    while (feof($handle)===false) {
        # code...
        yield fgets($handle);
    }

    fclose($handle);
}

foreach (readTxt() as $key => $value) {
    # code...
    echo $value.'<br />';
}

但是，背后的代码执行规则却一点儿也不一样。使用生成器读取文件，第一次读取了第一行，第二次读取了第二行，以此类推，每次被加载到内存中的文字只有一行，大大的减小了内存的使用。

这样，即使读取上G的文本也不用担心，完全可以像读取很小文件一样编写代码。

 百万级别的访问量

yield生成器是php5.5之后出现的，yield提供了一种更容易的方法来实现简单的迭代对象，相比较定义类实现 Iterator 接口的方式，性能开销和复杂性大大降低。

yield生成器允许你 在 foreach 代码块中写代码来迭代一组数据而不需要在内存中创建一个数组。
https://www.php.net/manual/zh/language.generators.overview.php

