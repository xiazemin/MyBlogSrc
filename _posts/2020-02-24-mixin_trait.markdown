---
title: mixin trait 多继承
layout: post
category: lang
author: 夏泽民
---
对于Mixin(混合)、Trait（特性）这两个面向对象特性，总是让人觉得说不清道不明的感觉，其实众多设计语言里，这里面的一些概念也是相互参杂的，并不是又那么一个严格的定义或界限说哪种一定是Mixin，或者哪种一定是Trait。这两种语言设施的提出，它的本质实际上都是解决代码复用的问题。

The developers of the Java language were well-versed in C++ and other languages that include multiple inheritance, whereby classes can inherit from an arbitrary number of parents. One of the problems with multiple inheritance is that it's impossible to determine which parent inherited functionality is derived from. This problem is called the diamond problem (see Resources). The diamond problem and other complexities that are inherent in multiple inheritance inspired the Java language designers to opt for single inheritance plus interfaces.

Interfaces define semantics but not behavior. They work well for defining method signatures and data abstractions, and all of the Java.next languages support Java interfaces with no essential changes. However, some cross-cutting concerns don't fit into a single-inheritance-plus-interfaces model.

在Java 中，一个类可以实现任意数量的接口。这个模型在声明一个类实现多个抽象的时候非常有用。不幸的是，它也有一个主要缺点。对于许多接口，大多数功能都可以用对于所有使用这个接口的类都有效的“样板”代码来实现。Java 没有提供一个内置机制来定义和使用这些可重用代码。相反的，Java 程序员必须使用一个特别的转换来重用一个已知接口的实现。在最坏的情况下，程序员必须复制粘贴同样的代码到不同的类中去。

本人认为这两个的涵义根据语言不同，而解释有所不同。但是它们的目的都是作为单继承不足的一种补充，或者是变相地实现多继承。实际上Java的接口也是变相的实现多继承，但是java的接口只是定义signature，没有实现体。在某种意义上Mixin和Trait这两者有点类似于抽象类，或者是有部分或全部实现体的Interface，但是在具体语言中，有表现出不一样的用法。总体上，笔者认为没有特别固定的或者是严格的区别。Mixin和Trait这两者都不能生成实例，否则就跟class没什么区别了。
<!-- more -->
Scala trait
class Person ; //实验用的空类

trait TTeacher extends Person {  

    def teach //虚方法，没有实现  

}  
trait TPianoPlayer extends Person {  

    def playPiano = {println("I’m playing piano. ")} //实方法，已实现  

}  
class PianoplayingTeacher extends Person with TTeacher with TPianoPlayer {  

    def teach = {println("I’m teaching students. ")} //定义虚方法的实现  

} 

PHP traits
  // the template
trait TSingleton {
  private static $_instance = null;

  public static function getInstance() {
    if (null === self::$_instance)
    {
      self::$_instance = new self();
    }

    return self::$_instance;
  }
}
class FrontController {
  use TSingleton;
}
// can also be used in already extended classes
class WebSite extends SomeClass {
  use TSingleton;
}


Ruby mixin
module Foo
  def bar
    puts "foo";
  end
end

然后我们把这个模块混入到对象中去：
class Demo
  include Foo
end 

如上编码后，模块中的实例方法就会被混入到对象中：
d=Demo.new
d.bar

区别：
1）Mixin可能更多的是指动态语言，它是在执行到某个点的时候，将代码插入到其中来达到代码复用的效果。Trait更多的是编译过程中，通过一些静态手段赋值代码到类中使得其拥有Trait中的一些功能以达到代码复用的目的；
2）“Mixins may contain state, (traditional) traits don't.”这个区别比较弱，事实上Scala中Trait已经可以保存状态了（成员变量）；
3）“Mixins use "implicit conflict resolution", traits use "explicit conflict resolution"”。这个区别可能是个明显的区别；但是如果某个语言它可以让Trait implicit resolve，那也没什么大不了。
4）“Mixins depends on linearization, traits are flattened.”这个区别可能有。至少Scala里面貌似Trait是Flattened处理的，跟Java嵌套类差不多

程序设计中，是该用类还是Mixin&Trait

当我们考虑是否一个“概念”应该成为一个Trait 或者一个类的时候，记住作为混入的Trait 对于“附属”行为来说最有意义。如果你发现某一个Trait 经常作为其它类的父类来用，导致子类会有像父Trait 那样的行为，那么考虑把它定义为一个类吧，让这段逻辑关系更加清晰

  对于Mixin(混合)、Trait（特性）这两个面向对象特性，总是让人觉得说不清道不明的感觉，其实众多设计语言里，这里面的一些概念也是相互参杂的，并不是又那么一个严格的定义或界限说哪种一定是Mixin，或者哪种一定是Trait。这两种语言设施的提出，它的本质实际上都是解决代码复用的问题。
  
  The developers of the Java language were well-versed in C++ and other languages that include multiple inheritance, whereby classes can inherit from an arbitrary number of parents. One of the problems with multiple inheritance is that it's impossible to determine which parent inherited functionality is derived from. This problem is called the diamond problem (see Resources). The diamond problem and other complexities that are inherent in multiple inheritance inspired the Java language designers to opt for single inheritance plus interfaces.

Interfaces define semantics but not behavior. They work well for defining method signatures and data abstractions, and all of the Java.next languages support Java interfaces with no essential changes. However, some cross-cutting concerns don't fit into a single-inheritance-plus-interfaces model.

在Java 中，一个类可以实现任意数量的接口。这个模型在声明一个类实现多个抽象的时候非常有用。不幸的是，它也有一个主要缺点。对于许多接口，大多数功能都可以用对于所有使用这个接口的类都有效的“样板”代码来实现。Java 没有提供一个内置机制来定义和使用这些可重用代码。相反的，Java 程序员必须使用一个特别的转换来重用一个已知接口的实现。在最坏的情况下，程序员必须复制粘贴同样的代码到不同的类中去。

现在排名靠前的面向对象的编程语言中，Java、C#等都是以单继承+接口来实现面向对象，但是这在一定程序了稀释了继承的力量， 因为在业内推荐以组合的方式使用类。这在一些常见的设计模式中有明显的体现，想想在GOF的23个设计模式中有多少个是使用了继承的呢？ 大多数是以接口+组合的方式实现。其实作为一个类来说，它也比较难做，即要能代码复用，又得被实例化，偏向谁呢？ 这个时候Mixin可能就有一些用武之地了。

Mixin最早起源于一个Lisp，Mixin鼓励代码重用，Mixin可以实现运行时的方法绑定，虽然类的属性和实例参数仍然是在编译时定义。 在面向对象编程语言，Mixin是一个提供了一些被用于继承或在子类中重用的功能的类，它类似于一种多继承， 但是实际上它是一种中小粒度的代码复用单元，而不直接用于实例化。 虽然这不是一种专业的方式进行功能复用，这在实现多继承的同时，在一定程序上避免了多继承的明显问题。

PHP和Java类似，也是单继承+接口。 我们知道，一个类可以实现任意数量的接口，这对一个类需要实现多个抽象的时候非常有用。 然而，对于要实现了多个接口的类，每个类都需要实现这些接口，而大多数情况下，这些接口都是可以共用的。 PHP并没有提供内置机制来定义和使用这些可重用代码，虽然我们可以对一地些接口使用一个抽象类来共用代码，但是如果这些类必须继承另一个抽象类呢？ 就算是可以通过抽象类的多次继承实现代码的共用，但是整个继承体系将会变得非常复杂，如果不能实现重用，那么可能我们只得CTRL + C 和 CTRL + V了。 大多数的情况下我们其实只是需要重用一些代码而已。

虽然PHP在之前没有提供完善的解决方案，但在新发布PHP5.4中，出现了一个关键字trait。 通过这个关键字我们可以定义抽象为一个Trait

https://wiki.php.net/rfc/traits

https://stackoverflow.com/questions/925609/mixins-vs-traits

Mixins may contain state, (traditional) traits don't.
Mixins use "implicit conflict resolution", traits use "explicit conflict resolution"
Mixins depends on linearization, traits are flattened.

http://stephane.ducasse.free.fr/Presentations/2009-TraitsAtSC.pdf

ad 1. In mixins you can define instance variables. Traits do not allow this. The state must be provided by composing class (=class using the traits)

ad 2. There may be the name conflict. Two mixins (MA and MB) or traits (TA and TB) define method with the same definition foo():void.

Mixin MA {
    foo():void {
        print 'hello'
    }
}

Mixin MB {
    foo():void {
        print 'bye'
    }
}

Trait TA {
    foo():void {
        print 'hello'
    }
}

Trait TB {
    foo():void {
        print 'bye'
    }
}
In mixins the conflicts in composing class C mixins MA, MB are resolved implicitly.

Class C mixins MA, MB {
    bar():void {
        foo();
    }
}
This will call foo():void from MA

On the other hand while using Traits, composing class has to resolve conflicts.

Class C mixins TA, TB {
    bar():void {
        foo();
    }
}
This code will raise conflict (two definitions of foo():void).

ad 3. The semantics of a method does not depend of whether it is defined in a trait or in a class that uses the trait.

In other words, it does not matter wheter the class consists of the Traits or the Traits code is "copy - pasted" into the class.


作者：刘缙
链接：https://www.zhihu.com/question/49094001/answer/124322283
来源：知乎
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

如果我们有两个类：class A {
    public f();
    public g();
    private int a;
    private int b;
};

class B {
    public h();
    private int x;
};在组合这两个类的所有方法中，一个极端是仅仅把两个类的对象组合起来，并且给两个类的公开方法都作转发：class AplusB {
    public f() { a.f(); }
    public g() { a.g(); }
    public h() { b.h(); }
    private A a;
    private B b;
};特点是：完全不破坏封装（不需要了解A和B的实现）。只需要解决A和B公开方法的名字冲突。另一个极端是把两个类的实现拷贝粘贴在一起：class AplusB {
    public f();
    public g();
    public h();
    private int a;
    private int b;
    private int x;
};
特点是：和继承一样，最大限度地破坏了封装（两个类必须了解对方的实现才能协同工作）需要解决所有可能的名字冲突（方法、属性）多继承、traits、mixin等等，其实都是在这两个极端之间，并且选择不同的名字冲突解决策略，仅此而已。它们没有绝对的优劣之分，分别适用于不同的场景。偏向前者的，需组合的两个类可以独立设计而不考虑配合，各自的内部实现可以随意修改而不影响组合的结果，适用范围广；偏向后者的，需组合的两个类内部实现需要严格配合，结果精巧脆弱，但表达能力强，做好了能发挥出1+1>2的效果，适用于局部、少数人参与的场景。