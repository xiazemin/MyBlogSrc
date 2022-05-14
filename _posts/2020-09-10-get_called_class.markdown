---
title: get_called_class
layout: post
category: php
author: 夏泽民
---
在一定的需求场景下，你有一个父类和一些子类，你需要获取这些子类的实例又不想在每个子类中写重复的getInstance()方法。

在各种语言中，一般会用一个工具类去做这个实现。但在php5.3之后，会有另外一种方法。

php5.3之后加入的新的特性：静态延迟绑定。这个特性允许在运行时获取静态继承的上下文。

get_called_class() 可以获取被调用的类。
static 关键字用来访问静态继承的上下文。

https://www.liaohuqiu.net/cn/posts/php-singleton-of-children-class/
<!-- more -->
思考：self关键字适用于类内部代替类的，代替的是当前所在的类本身，随着继承的实现，如果子类子访问父类的方法的时候，self到底代替的是谁呢？

引入：self是一种静态绑定，换言之即使当类进行编译的时候seld已经明确绑定了类名，因此不论多少继承，也不管是子类还是父类自己来进行访问

self代表的都是当前类，如果想要选择性的来支持来访者，就需要使用静态延迟绑定。

 

静态延迟绑定【掌握】
定义：静态延迟绑定，即在类内部用来代表类本身的关键字部分不是在类编译时固定好，而是当方法被访问时动态的选择来访者所属的类，静态

延迟绑定就是利用static 关键字代替静态绑定self，静态延迟绑定需要使用到静态成员的重写。（跟$this比较像  ）

1.静态延迟绑定：使用static关键字代替self进行类成员访问

复制代码
<?php

class posen{
   // 静态属性
   public static $name='posen';

   // 静态方法
   public static function show(){
      echo self::$name.'self::<br>';          //静态绑定
      echo static::$name.'static::<br>';       //静态延迟绑定
   }

}

posen::show();      //两个都能输出 posen 说明两个调用都可以
?>

复制代码
2.静态延迟绑定一定是通过继承后的子类来进行访问才有效果

复制代码
<?php

class posen{
   // 静态属性
   public static $name='posen';

   // 静态方法
   public static function show(){
      echo self::$name.'self::<br>';          //静态绑定
      echo static::$name.'static::<br>';       //静态延迟绑定
   }

}
// 子类继承
class man extends posen{
   // 重写父类中的静态属性name
   public static $name="man";   //有了这个属性，就显示man:static  了

}

man::show();   //还是显示的posen 为什么呢？不是说使用static 延迟绑定就能指向调用的类吗？
               // 因为你man类中没有自己的静态的属性，所以它就向上一层去找 找到了posen
?>
复制代码
 

注意：self关键字 --在你调用这个类加载到内存编译的时候self就绑定了当前的类，而static 在编译的时候则不会绑定，而是调用的时候在绑定调用的类

 

总结：

　　1.静态延迟绑定是指通过static关键字进行类静态成员的访问，是指在被访问时才决定到底使用那个类

　　2.静态延迟绑定对比的是静态绑定self

　　3.静态延迟绑定的意义是用来保证访问的静态成员是根据调用类的不同而选择不同的表现
　　https://segmentfault.com/a/1190000013741642
　　
　　
　　简单理解 PHP 延迟静态绑定
static:: 中的 static 其实是运行时所在类的别名，并不是定义类时所在的那个类名。这个东西可以实现在父类中能够调用子类的方法和属性。

使用 (static) 关键字来表示这个别名，和静态方法，静态类没有半毛钱的关系，static:: 不仅支持静态类，还支持对象（动态类）。

预备概念
转发调用
所谓的 “转发调用”（forwarding call）指的是通过以下几种方式进行的静态调用：self::，parent::，static:: 以及 forward_static_call ()。

非转发调用
那么非转发调用其实就是明确指定类名的静态调用（foo::bar ()）和非静态调用 ($foo->bar ())

后期静态绑定原理
后期静态绑定工作原理是存储了在上一个 “非转发调用”（non-forwarding call）的类名。

例子 1，简单使用 static::
class A {
    public static function who() {
        echo __CLASS__;
    }
    public static function test() {
        static::who(); // 后期静态绑定从这里开始
    }
}
class B extends A {
    public static function who() {
        echo __CLASS__;
    }
}
B::test();
以上例程会输出：

B
例子 2，区分转发调用和非转发调用
class A {
    public static function foo() {
        static::who();
    }

    public static function who() {
        echo __CLASS__."\n";
    }
}

class B extends A {
    public static function test() {
        A::foo();
        parent::foo();
        self::foo();
    }

    public static function who() {
        echo __CLASS__."\n";
    }
}
class C extends B {
    public static function who() {
        echo __CLASS__."\n";
    }
}

C::test();
以上例程会输出：

A
C
C
例子 3，使用场景举例
class Model 
{ 
    public static function find() 
    { 
        echo static::$name; 
    } 
} 

class Product extends Model 
{ 
    protected static $name = 'Product'; 
} 

Product::find();


https://www.php.net/manual/zh/language.oop5.late-static-bindings.php