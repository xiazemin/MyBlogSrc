---
title: Assertion
layout: post
category: golang
author: 夏泽民
---
Type assertion(断言)是用于 interface value 的一种操作，语法是 x.(T)，x 是 interface type 的表达式，而 T 是 assertd type，被断言的类型。

断言的使用主要有两种情景:

如果 asserted type 是一个 concrete type，一个实例类 type，断言会检查 x 的 dynamic type 是否和 T 相同，如果相同，断言的结果是 x 的 dynamic value，当然 dynamic value 的 type 就是 T 了。换句话说，对 concrete type 的断言实际上是获取 x 的 dynamic value。

如果 asserted type 是一个 interface type，断言的目的是为了检测 x 的 dynamic type 是否满足 T，如果满足，断言的结果是满足 T 的表达式，但是其 dynamic type 和 dynamic value 与 x 是一样的。换句话说，对 interface type 的断言实际上改变了 x 的 type，通常是一个更大 method set 的 interface type，但是保留原来的 dynamic type 和 dynamic value。
<!-- more -->
Type Switches
Interface 一般被用在这两种场合，一种是像 io.Reader, io.Writer 那样，一个 interface 的 method 真正含义是表达了实现这个 interface 的不同 concrete type 的相似性，意味着这里充分发挥的是 interface method 的表现力。重点在 method，而不是 concrete type。

一种是利用 interface 可以存储不同 concrete type 的能力，在必要的时候根据不同的 concrete type 做不同的处理，这样的用法就是利用 interface 的 assertion 来判断 dynamic type 的类型来做出具体的判断。重点在 concrete type，而不是 method。

Type switch 就是利用 interface 存储不同 concrete type 的能力来实现的 assertion。

switch x.(type) {
    case nil:
case int, uint:
case bool:
case string:
default:
}

类型（type）中非常重要的一类（category）就是接口类型（interface type），一个接口就表示一组确定的方法（method）集合。一个接口变量能存储任意的具体值（这里的具体concrete就是指非接口的non-interface），只要这个具体值所属的类型实现了这个接口的所有方法。

空接口表示方法集合为空并且可以保存任意值，因为任意值都有0个或者更多方法。

有些人说Go的接口是动态类型化的，但这是一种误导。Go的接口都是静态类型化的：一个接口类型变量总是保持同一个静态类型，即使在运行时它保存的值的类型发生变化，这些值总是满足这个接口。
https://research.swtch.com/interfaces

一个接口中的pair总有（值，具体类型）这样的格式，而不能有（值，接口类型）这样的格式。接口不能保存接口值（也就是说，你没法把一个接口变量值存储到一个接口变量中，只能把一个具体类型的值存储到一个接口变量中。）

第一反射定律（The first law of reflection)
1.从接口值到反射对象的反射（Reflection goes from interface value to reflection object）
最最基本的，反射是一种检查存储在接口变量中的（类型，值）对的机制。作为一个开始，我们需要知道reflect包中的两个类型：Type和Value。这两种类型给了我们访问一个接口变量中所包含的内容的途径，另外两个简单的函数reflect.Typeof和reflect.Valueof可以检索一个接口值的reflect.Type和reflect.Value部分。（还有就是，我们可以很容易地从reflect.Value到达reflect.Type，但是现在暂且让我们先把Value和Type的概念分开说。先剧透，从Value到达Type是通过Value中定义的某些方法来实现的，虽然先分开讲，但是后面多注意一下。）

第二反射定律（The second law of reflection）
2.从反射队形到接口值的反射（Reflection goes from reflection object to interface value）
就像物理学上的反射，Go中到反射可以生成它的逆。

给定一个reflect.Value，我们能用Interface方法把它恢复成一个接口值；效果上就是这个Interface方法把类型和值的信息打包成一个接口表示并且返回结果：
// Interface returns v's value as an interface{}.
func (v Value) Interface() interface{}

第三反射定律（The third law of reflection）
3.为了修改一个反射对象，值必须是settable的（To modify a reflection object, the value must be settable)


interface在内存上实际由两个成员组成,tab指向虚表，data则指向实际引用的数据。虚表描绘了实际的类型信息及该接口所需要的方法集

观察itable的结构，首先是描述type信息的一些元数据，然后是满足Stringger接口的函数指针列表（注意，这里不是实际类型Binary的函数指针集哦）。 因此我们如果通过接口进行函数调用，实际的操作其实就是s.tab->fun[0](s.data)。 是不是和C++的虚表很像？接下来我们要看看golang的虚表和C++的虚表区别在哪里。

先看C++，它为每种类型创建了一个方法集，而它的虚表实际上就是这个方法集本身或是它的一部分而已，当面临多继承时（或者叫实现多个接口时，这是很常见的），C++对象结构里就会存在多个虚表指针，每个虚表指针指向该方法集的不同部分，因此，C++方法集里面函数指针有严格的顺序。 许多C++新手在面对多继承时就变得紧张，因为它的这种设计方式，为了保证其虚表能够正常工作，C++引入了很多概念，什么虚继承啊，接口函数同名问题啊，同一个接口在不同的层次上被继承多次的问题啊等等…… 就是老手也很容易因疏忽而写出问题代码出来。

我们再来看golang的实现方式，同C++一样，golang也为每种类型创建了一个方法集，不同的是接口的虚表是在运行时专门生成的。 可能细心的同学能够发现为什么要在运行时生成虚表。 因为太多了，每一种接口类型和所有满足其接口的实体类型的组合就是其可能的虚表数量，实际上其中的大部分是不需要的，因此golang选择在运行时生成它，例如，当例子中当首次遇见s := Stringer(b)这样的语句时，golang会生成Stringer接口对应于Binary类型的虚表，并将其缓存。

理解了golang的内存结构，再来分析诸如类型断言等情况的效率问题就很容易了，当判定一种类型是否满足某个接口时，golang使用类型的方法集和接口所需要的方法集进行匹配，如果类型的方法集完全包含接口的方法集，则可认为该类型满足该接口。 例如某类型有$m$个方法，某接口有$n$个方法，则很容易知道这种判定的时间复杂度为$O(m \times n)$，不过可以使用预先排序的方式进行优化，实际的时间复杂度为$O(m+n)$。

2. 使用interface的注意事项
将对象赋值给接口变量时会复制该对象。
接口使用的是一个名为itab的结构体存储的 type iface struct{ tab *itab // 类型信息 data unsafe.Pointer // 实际对象指针 }
只有接口变量内部的两个指针都为nil的时候，接口才等于nil。
interface实际上是一个引用(只保存了两个值)，因此传递它并不会造成太多的损耗。
3. 其他
Go语言中的interface很灵活，但是也付出了一定的性能代价。
如果是性能关键的代码，可以考虑放弃interface，自己写原生的代码。
