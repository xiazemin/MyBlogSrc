---
title: Duck typing
layout: post
category: web
author: 夏泽民
---
<!-- more -->
还是先看定义 duck typing,
    鸭子类型是多态(polymorphism)的一种形式.在这种形式中,不管对象属于哪个,
    也不管声明的具体接口是什么,只要对象实现了相应的方法,函数就可以在对象上执行操作.
    即忽略对象的真正类型，转而关注对象有没有实现所需的方法、签名和语义.
        duck typing
            A form of polymorphism where functions
            operate on any object that implements the
            appropriate methods, regardless of their
            classes or explicit interface declarations.

    Wikipedia 是这样描述 duck typing 的,
        在计算机语言中, duk typing 是一个类型测试的一个具体应用.
        是将对类型的检查推迟到代码运行的时候,由动态类型(dynamic typing)
        或者反省(reflection)实现. duck typing 应用在通过应用规则/协议(protocol)
        建立一个适合的对象 object.
        '如果它走起步来像鸭子,并且叫声像鸭子, 那个它一定是一只鸭子.'
        对于一般类型 normal typing, 假定一个对象的 suitability 只有该对象的类型决定.
        然而,对于 duck typing 来说, 一个对象 object 的 suitability 是通过该对象是否
        实现了特定的方法跟属性来决定 certain methods and properties, 而不是由该对象
        的来类型决定.

        注,
            In computer science, reflection is the ability of a computer program to
            examine,introspect, and modify its own structure and behavior at runtime.

        From Wikipedia,
            In computer programming, duck typing is an application of the duck test
            in type safety.It requires that type checking be deferred to runtime,
            and is implemented by means of dynamic typing or reflection.
            Duck typing is concerned with establishing the suitability of an object
            for some purpose, using the principle, "If it walks like a duck and it
            quacks like a duck, then it must be a duck." With normal typing,
            suitability is assumed to be determined by an object's type only.
            In duck typing, an object's suitability is determined by the presence
            of certain methods and properties (with appropriate meaning),
            rather than the actual type of the object.


    鸭子类型的起源 Origins of duck-typing,
        现在谷歌工程师,Python 社区重要贡献者之一: Alex Martelli 说到,
            我相信是 Ruby 社区推动了 duck typing 这个术语的流行.
            但是这个duck typing 这种表达在 Ruby 和 Python 火之前,
            就是在Python 的讨论中使用过.

        根据 Wikipedia, duck typing 这一术语最早被 Alex Martelli 在 2000 所使用.
        Related Link of Wikipedia - https://en.wikipedia.org/wiki/Duck_typing

    归功于 python 的 数据类型 data model, 你的用户自定义类型的行为可以像 built-in 类型一样自然。
    这并不需要通过继承 inheritance 来获得. 本着 duck typing, 可以在对象中只实现需要的方法, 就能
    保证保证对象的行为符合预期. 对 Python 来说，这基本上是指避免使用 isinstance 检查对象的类,
    更别提 type(foo) is bar 这种更糟的检查方式了，这样做没有任何好处，甚至禁止最简单的继承方式.
    具体使用时,上述建议有一个常见的例外：有些 Python API 接受一个字符串或字符串序列;
    如果只有一个字符串,可以把它放到列表中,从而简化处理. 因为字符串是序列类型,
    所以为了把它和其他不可变序列区分开,最简单的方式是使用 isinstance(x, str) 检查.
    另一方面，如果必须强制执行 API 契约，通常可以使用 isinstance 检查抽象基类。

    在看例子之前, 先看简略一下儿 协议 protocol 相关内容,
        在 Python 中创建功能完善的序列类型无需使用继承, 只需实现符合序列协议的方法.
        在面向对象编程中,协议是非正式的接口,只在文档中定义,在代码中不定义.
        例如,Python 的序列协议只需要 __len__ 和 __getitem__ 两个方法.
        任对象/类型(A)只要使用标准的签名和语义实现了这两个方法,就能用在任何期待序列的地方,
        然而A 是不是哪个类的子类无关紧要,只要提供了所需的方法即可.这就是 python 序列协议.
        协议是非正式的,没有强制力,因此如果你知道类的具体使用场景,通常只需要实现一个协议的部分.
        例如,为了支持迭代,只需实现 __getitem__ 方法，没必要提供 __len__方法.

        经典示例, duck typing 处理一个字符串 string 或 可迭代字符串 iterable of strings
            try:                                                      #1
                field_names = field_names.replace(',', ' ').split()   #2
            except AttributeError:                                    #3
                pass                                                  #4
            field_names = tuple(field_names)                          #5

            #1, 假定 field_names 是一个字符串 string. EAFP, it’s easier to ask forgiveness than permission
            #2, 将 field_names 中的 ',' 替换成空格 ' ' 并 split, 将结果放到 list 中
            #3, sorry, field_names 并不像一个 str, field_names 不能 .replace 或者 .replace 后返回的结果不能 .split()
            #4, 这里我men假设 新的 field_names 是一个可迭代对象
            #5, 确保新的 field_names 是一个可迭代对象, 同事保存一个 copy - create 一个 tuple

            field_names = 'abc'                                       #6
            field_names = 'A,B,C'                                     #7
            try:
                field_names = field_names.replace(',', ' ').split()
            except AttributeError:
                pass
            print(field_names)
            field_names = tuple(field_names)
            print(field_names)
            for item in field_names:
                print(item)

            Output,
                ['abc']            #6
                ('abc',)           #6
                abc                #6
                --------------
                ['A', 'B', 'C']    #7
                ('A', 'B', 'C')    #7
                A                  #7
                B                  #7
                C                  #7

        结论,
            Summarize, Outside of frameworks, duck typing is often sim‐pler and more flexible than type checks.
 
对于一门强类型的静态语言来说，要想通过运行时多态来隔离变化，多个实现类就必须属于同一类型体系。也就是说，它们必须通过继承的方式，与同一抽象类型建立is-a关系。

而Duck Typing则是一种基于特征，而不是基于类型的多态方式。事实上它仍然关心is-a，只不过这种is-a关系是以对方是否具备它所关心的特征来确定的。

James Whitcomb Riley在描述这种is-a的哲学时，使用了所谓的鸭子测试（Duck Test）:

当我看到一只鸟走路像鸭子，游泳像鸭子，叫声像鸭子，那我就把它叫做鸭子。（When I see a bird that walks like a duck and swims like a duck and quacks like a duck, I call that bird a duck.）

鸭子测试
Duck Test基于特征的哲学，给设计提供了强大的灵活性。动态面向对象语言，如Python，Ruby等，都遵从了这种哲学来实现运行时多态。下面给出一个Python的例子：

class Duck:
    def quack(self):
        print("Quaaaaaack!")
    def feathers(self):
        print("The duck has white and gray feathers.")

class Person:
    def quack(self):
        print("The person imitates a duck.")
    def feathers(self):
        print("The person takes a feather from the ground and shows it.")
    def name(self):
        print("John Smith")

def in_the_forest(duck):
    duck.quack()
    duck.feathers()

def game():
    donald = Duck()
    john = Person()
    in_the_forest(donald)
    in_the_forest(john)

game()
但这并不意味着Duck Typing是动态语言的专利。C++作为一门强类型的静态语言，也对此特性有着强有力的支持。只不过，这种支持不是运行时，而是编译时。

其实现的方式为：一个模板类或模版函数，会要求其实例化的类型必须具备某种特征，如某个函数签名，某个类型定义，某个成员变量等等。如果特征不具备，编译器会报错。

比如下面一个模板函数:

template <typename T> 
void f(const T& object) 
{ 
  object.f(0); // 要求类型 T 必须有一个可让此语句编译通过的函数。
} 
对于这样一个函数，下面的四个类均可以用来作为其参数类型。

struct C1 
{
  void f(int); 
};
 
struct C2 
{ 
  int f(char); 
};
 
struct C3 
{ 
  int f(unsigned short, bool isValid = true); 
}; 
 
struct C4
{
  Foo* f(Object*);
};
一旦上述模板函数实现为下面的样子，则只有C2和C3可以和f配合工作。

template <typename T> 
void f(const T& object) 
{ 
  int result = object.f(0); 
  // ... 
} 
通过之前的解释我们不难发现，Duck Typing要表达的多态语义如下图所示：

DuckTyping的语义
适配器：类型萃取
Duck Typing需要实例化的类型具备一致的特征，而模板特化的作用正是为了让不同类型具有统一的特征（统一的操作界面），所以模板特化可以作为Duck Typing与实例化类型之间的适配器。这种模板特化手段称为萃取（Traits)，其中类型萃取最为常见，毕竟类型是模板元编程的核心元素。

所以，类型萃取首先是一种非侵入性的中间层。否则，这些特征就必须被实例化类型提供，而就意味着，当一个实例化类型需要复用多个Duck Typing模板时，就需要迎合多种特征，从而让自己经常被修改，并逐渐变得庞大和难以理解。

Type Traits的语义
另外，一个Duck Typing模板，比如一个通用算法，需要实例化类型提供一些特征时，如果一个类型是类，则是一件很容易的事情，因为你可以在一个类里定义任何需要的特征。但如果一个基本类型也想复用此通用算法，由于基本类型无法靠自己提供算法所需要的特征，就必须借助于类型萃取。

结论
这四篇文章所介绍的，就是C++泛型编程的全部关键知识。

从中可以看出，泛型是一种多态技术。而多态的核心目的是为了消除重复，隔离变化，提高系统的正交性。因而，泛型编程不仅不应该被看做奇技淫巧，而是任何一个追求高效的C++工程师都应该掌握的技术。

同时，我们也可以看出，相关的思想在其它范式和语言中（FP，动态语言）也都存在。因而，对于其它范式和语言的学习，也会有助于更加深刻的理解泛型，从而正确的使用范型。

最后给出关于泛型的缺点：

复杂模板的代码非常难以理解;
编译器关于模板的出错信息十分晦涩，尤其当模板存在嵌套时；
模板实例化会进行代码生成，重复信息会被多次生成，这可能会造成目标代码膨胀;
模板的编译可能非常耗时;
编译器对模板的复杂性往往会有自己限制，比如当使用递归时，当递归层次太深,编译器将无法编译;
不同编译器（包括不同版本）之间对于模板的支持程度不一，当存在移植性需求时，可能出现问题;
模板具有传染性，往往一处选择模板，很多地方也必须跟着使用模板，这会恶化之前的提到的所有问题。
我对此的原则是：在使用其它非泛型技术可以同等解决的前提下，就不会选择泛型。

python与鸭子类型
调用不同的子类将会产生不同的行为，而无须明确知道这个子类实际上是什么，这是多态的重要应用场景。而在python中，因为鸭子类型(duck typing)使得其多态不是那么酷。 
鸭子类型是动态类型的一种风格。在这种风格中，一个对象有效的语义，不是由继承自特定的类或实现特定的接口，而是由”当前方法和属性的集合”决定。这个概念的名字来源于由James Whitcomb Riley提出的鸭子测试，“鸭子测试”可以这样表述：“当看到一只鸟走起来像鸭子、游泳起来像鸭子、叫起来也像鸭子，那么这只鸟就可以被称为鸭子。” 
在鸭子类型中，关注的不是对象的类型本身，而是它是如何使用的。例如，在不使用鸭子类型的语言中，我们可以编写一个函数，它接受一个类型为”鸭子”的对象，并调用它的”走”和”叫”方法。在使用鸭子类型的语言中，这样的一个函数可以接受一个任意类型的对象，并调用它的”走”和”叫”方法。如果这些需要被调用的方法不存在，那么将引发一个运行时错误。任何拥有这样的正确的”走”和”叫”方法的对象都可被函数接受的这种行为引出了以上表述，这种决定类型的方式因此得名。 
鸭子类型通常得益于不测试方法和函数中参数的类型，而是依赖文档、清晰的代码和测试来确保正确使用。

静态类型语言和动态类型语言的区别
静态类型语言在编译时便已确定变量的类型，而动态类型语言的变量类型要到程序运行的时候，待变量被赋予某个值之后，才会具有某种类型。 
静态类型语言的优点首先是在编译时就能发现类型不匹配的错误，编辑器可以帮助我们提前避免程序在运行期间有可能发生的一些错误。其次，如果在程序中明确地规定了数据类型，编译器还可以针对这些信息对程序进行一些优化工作，提高程序执行速度。 
静态类型语言的缺点首先是迫使程序员依照强契约来编写程序，为每个变量规定数据类型，归根结底只是辅助我们编写可靠性高程序的一种手段，而不是编写程序的目的，毕竟大部分人编写程序的目的是为了完成需求交付生产。其次，类型的声明也会增加更多的代码，在程序编写过程中，这些细节会让程序员的精力从思考业务逻辑上分散开来。 
动态类型语言的优点是编写的代码数量更少，看起来也更加简洁，程序员可以把精力更多地放在业务逻辑上面。虽然不区分类型在某些情况下会让程序变得难以理解，但整体而言，代码量越少，越专注于逻辑表达，对阅读程序是越有帮助的。 
动态类型语言的缺点是无法保证变量的类型，从而在程序的运行期有可能发生跟类型相关的错误。 
动态类型语言对变量类型的宽容给实际编码带来了很大的灵活性。由于无需进行类型检测，我们可以尝试调用任何对象的任意方法，而无需去考虑它原本是否被设计为拥有该方法。

面向接口编程
动态类型语言的面向对象设计中，鸭子类型的概念至关重要。利用鸭子类型的思想，我们不必借助超类型的帮助，就能轻松地在动态类型语言中实现一个原则：“面向接口编程，而不是面向实现编程”。例如，一个对象若有push和pop方法，并且这些方法提供了正确的实现，它就可以被当作栈来使用。一个对象如果有length属性，也可以依照下标来存取属性（最好还要拥有slice和splice等方法），这个对象就可以被当作数组来使用。

在静态类型语言中，要实现“面向接口编程”并不是一件容易的事情，往往要通过抽象类或者接口等将对象进行向上转型。当对象的真正类型被隐藏在它的超类型身后，这些对象才能在类型检查系统的“监视”之下互相被替换使用。只有当对象能够被互相替换使用，才能体现出对象多态性的价值。

python中的多态
python中的鸭子类型允许我们使用任何提供所需方法的对象，而不需要迫使它成为一个子类。 
由于python属于动态语言，当你定义了一个基类和基类中的方法，并编写几个继承该基类的子类时，由于python在定义变量时不指定变量的类型，而是由解释器根据变量内容推断变量类型的（也就是说变量的类型取决于所关联的对象），这就使得python的多态不像是c++或java中那样，定义一个基类类型变量而隐藏了具体子类的细节


而scala是静态强类型语言,  调用的方法必须在对象类型层次(本类或者超类）中定义。不过scala通过structural types支持所谓的类型安全的鸭子类型:
     类型安全的鸭子类型 - structural types  - structural types as a type-safe approach to duck typing