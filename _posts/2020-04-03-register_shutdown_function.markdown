---
title: __destruct与register_shutdown_function执行的先后顺序
layout: post
category: php
author: 夏泽民
---
根据php手册的解析。

__destruct是

析构函数会在到某个对象的所有引用都被删除或者当对象被显式销毁时执行。

而register_shutdown_function是

Registers a callback to be executed after script execution finishes or exit() is called. 注册一个回调函数，此函数在脚本运行完毕或调用exit()时执行。

从字面上理解，__destruct是对象层面的，而register_shutdown_function是整个脚本层面的，理应register_shutdown_function的级别更高，其所注册的函数也应最后执行。
<!-- more -->
egister_shutdown_function(function(){echo 'global';});
    class A {
        public function __construct(){
        }
        public function __destruct()
        {
            echo __class__,'::',__function__,'<br/>';
        }
    }
    new A;
执行结果：

复制代码代码如下:

A::__destruct
global
完全证实了我们的猜测，它按照对象->脚本的顺序被执行了。

但如果我们在对象中注册了register_shutdown_function呢？它还是一样的顺序吗？！

class A {
        public function __construct(){
            register_shutdown_function(function(){echo 'local', '<br/>';});
        }
        public function __destruct()
        {
            echo __class__,'::',__function__,'<br/>';
        }
    }
    new A;
    
结果：
local
A::__destruct

可以看到register_shutdown_function先被调用了，最后才是执行对象的__destruct。这表明register_shutdown_function注册的函数被当作类中的一个方法？！不得而知，这可能需要查看php源代码才能解析了。

我们可以扩大范围查看情况：

复制代码代码如下:

register_shutdown_function(function(){echo 'global', '<br/>';});
    class A {
        public function __construct(){
            register_shutdown_function(array($this, 'op'));
        }
        public function __destruct()
        {
            echo __class__,'::',__function__,'<br/>';
        }
        public function op()
        {
            echo __class__,'::',__function__,'<br/>';
        }
    }
    class B {
        public function __construct()
        {
            register_shutdown_function(array($this, 'op'));
            $obj = new A;
        }
        public function __destruct()
        {
            echo __class__,'::',__function__,'<br/>';
        }
        public function op()
        {
            echo __class__,'::',__function__,'<br/>';
        }
    }
    $b = new B;
我们在全局注册一个register_shutdown_function函数，在类AB中又各注册了一个，而且类中分别还有析构方法。最后运行结果会怎样呢？

复制代码代码如下:

global
B::op
A::op
A::__destruct
B::__destruct
结果完全颠覆了我们的想像，register_shutdown_function函数无论在类中注册还是在全局注册，它都是先被执行，类中执行的顺序就是它们被注册的先后顺序。如果我们再仔细研究，全局的register_shutdown_function函数无论放在前面还是后面都是这个结果，事情似乎有了结果，那就是register_shutdown_function比__destruct先执行，全局的register_shutdown_function函数又先于类中注册的register_shutdown_function先执行。

且慢，我无法接受这个结果，按照这样的结论，难道说脚本已经结束后还可以再执行__destruct？！因此，我还要继续验证这个结论---去掉类中注册register_shutdown_function，而保留全局register_shutdown_function：

lass A {
        public function __destruct()
        {
            echo __class__,'::',__function__,'<br/>';
        }
    }
    class B {
        public function __construct()
        {
            $obj = new A;
        }
        public function __destruct()
        {
            echo __class__,'::',__function__,'<br/>';
        }
    }
    register_shutdown_function(function(){echo 'global', '<br/>';});
输出：

复制代码代码如下:

A::__destruct
global
B::__destruct
结果令人茫然，A、B两个类的析构函数执行顺序无可质疑，因为B中调用了A，类A肯定比B先销毁，但全局的register_shutdown_function函数又怎么夹在它们中间被执行？！费解。

按照手册的解析，析构函数也可在调用exit时执行。

析构函数即使在使用 exit()终止脚本运行时也会被调用。在析构函数中调用 exit() 将会中止其余关闭操作的运行。

如果在函数中调用exit，它们又如何被调用的呢？

复制代码代码如下:

class A {
        public function __construct(){
            register_shutdown_function(array($this, 'op'));
            exit;
        }
        public function __destruct()
        {
            echo __class__,'::',__function__,'<br/>';
        }
        public function op()
        {
            echo __class__,'::',__function__,'<br/>';
        }
    }
    class B {
        public function __construct()
        {
            register_shutdown_function(array($this, 'op'));
            $obj = new A;
        }
        public function __destruct()
        {
            echo __class__,'::',__function__,'<br/>';
        }
        public function op()
        {
            echo __class__,'::',__function__,'<br/>';
        }
    }
    register_shutdown_function(function(){echo 'global', '<br/>';});
    $b = new B;
输出：

复制代码代码如下:

global
B::op
A::op
B::__destruct
A::__destruct
这个顺序与上述第三个例子相似，不同的且令人不可思议的是B类的析构函数先于类A执行，难道销毁B后类A的所有引用才被全部销毁？！不得而知。

结论：
1、尽量不要在脚本中将register_shutdown_function与__destruct混搭使用，它们的行为完全不可预测。
1、因为对象在相互引用，因此我们无法测知对象几时被销毁，当需要按顺序输出内容时，不应把内容放在析构函数__destruct里；
2、尽量不要在类中注册register_shutdown_function，因为它的顺序难以预测（只有调用这个对象时才会注册函数），而且__destruct完全可以代替register_shutdown_function；
3、如果需要在脚本退出时执行相关动作，最好在脚本开始时注册register_shutdown_function，并把所有动作放在一个函数里。



在进行开发的过程中，通过register_shutdown_function注册了一个函数进行日志刷新磁盘，但是每次在一个对象的__destruct打印的日志都不能正确的刷新到磁盘，官方文档也没有说明到底是谁先执行，所以看了下源码，跟大家分享一下，避免大家踩坑。

直接看源码：
很清晰的看出来是先调用register_shutdown_function中的方法，而析构函数第二步执行，最后flush buffer。

// php7.0.6线上版本
// 代码路径：php-src/main/main.c
void php_request_shutdown(void *dummy)
{
……
    /* 1. Call all possible shutdown functions registered with register_shutdown_function() */
    if (PG(modules_activated)) zend_try {
        php_call_shutdown_functions();
    } zend_end_try();

    /* 2. Call all possible __destruct() functions */
    zend_try {
        zend_call_destructors();
    } zend_end_try();

    /* 3. Flush all output buffers */
    zend_try {
……
}
验证：
<?php
    // test.php
    register_shutdown_function(function() {
        echo "1.shutdown\n";
    });
    class Test {
        function __destruct() {
            echo "2.destruct\n";
        }
    }
    /* 
     * 这里如果不赋值给$test, 会优先执行__destruct，
     * 因为相当于执行了unset($test)，而不是在请求结束后销毁的对象
     */
    $test = new Test();
未修改前：
php test.php
1.shutdown
2.destruct

修改源码，换了一下顺序：

// php7.0.6线上版本
// 代码路径：php-src/main/main.c
void php_request_shutdown(void *dummy)
{
  ……
    /* 2. Call all possible __destruct() functions */
    zend_try {
        zend_call_destructors();
    } zend_end_try();

    /* 1. Call all possible shutdown functions registered with register_shutdown_function() */
    if (PG(modules_activated)) zend_try {
        php_call_shutdown_functions();
    } zend_end_try();

    /* 3. Flush all output buffers */
    zend_try {
……
}
修改后结果：
php test.php
2.destruct
1.shutdown

结论：
当请求结束后，需要PHP自动释放的对象，PHP优先执行register_shutdown_function注册的函数，后执行对象的__destruct，所以尽量不要混合使用__destruct和register_shutdown_function，除非你清楚他们的执行顺序对你没影响

fastcgi_finish_request,register_shutdown_function和__destruct的理解


针对nginx和php-fpm模式，php定义了一个函数fastcgi_finish_request，可以提高接口返回数据的速度。

nginx的fastcgi模块与php-fpm程序进行交互，获取php-fpm的worker进程执行的结果。一般情况，php的进程完全执行完后，才会吧输出的数据flush到nginx的fastcgi缓存区。在php进程中执行fastcgi_finish_request()，可以主动把进程的输出数据flush到nginx,这时候php-fpm的worker继续执行到程序结束。

__destruct(),析构函数，只有当实例删除时才会执行。

register_shutdown_function()，进程结束时，执行注册的函数。无论php进程是何原因结束，包括err等等，都会执行这里的函数（执行注册的时候，就会把代码加入到内存中）。

class A {
	public function test()
	{
		echo 'aaaaa';
        register_shutdown_function(function () {
            echo 'shutdown';
        });
	}
	public function __destruct()
    {
        echo '__destruct'."<br>";
    }
}
$class = new A();
$class->test();
接口会输出：aaaaashutdown__destruct，三个地方的输出都会有。
class A {
	public function test()
	{
		echo 'aaaaa';
        register_shutdown_function(function () {
            echo 'shutdown';
        });
	}
	public function __destruct()
    {
        echo '__destruct'."<br>";
    }
}
$class = new A();
$class->test();
if (function_exists('fastcgi_finish_request')) {
        fastcgi_finish_request();//主动flush数据给nginx
   }
接口输出：aaaaa。因为__destruct，register_shutdown_function的注册函数都还没有执行，就已经返回数据给nginx了，并且会关闭这个cgi连接。所以在__destruct，register_shutdown_function的注册函数中不要做输出操作，不然nginx接收不到后面的输出了。