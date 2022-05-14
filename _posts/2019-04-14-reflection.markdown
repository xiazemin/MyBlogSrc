---
title: reflection
layout: post
category: golang
author: 夏泽民
---
反射指的是程序借助某种手段检查自己结构的一种能力，通常就是借助编程语言中定义的各种类型（types）。
<!-- more -->
类型和接口（Types and interfaces）
因为反射是建立在类型系统（the type system）上的，所以让我们从复习Go中的类型开始讲起。

Go是静态类型化的。每个变量都有一个静态类型，也就是说，在编译的时候变量的类型就被很精确地确定下来了，比如要么是int，或者是float32，或者是MyType类型，或者是[]byte等等。如果我们像下面这样声明：
type MyInt int

var i int
var j MyInt
那么i的类型就是int，而j的类型就是MyInt。这里的变量i和j具有不同的静态类型，虽然它们有相同的底层类型（underlying type），如果不显示的进行强制类型转换它们是不能互相赋值的。

类型（type）中非常重要的一类（category）就是接口类型（interface type），一个接口就表示一组确定的方法（method）集合。一个接口变量能存储任意的具体值（这里的具体concrete就是指非接口的non-interface），只要这个具体值所属的类型实现了这个接口的所有方法。一个大家都很熟悉的例子是io.Reader和io.Writer，类型Reader和类型Writer来自io包:
// Reader is the interface that wraps the basic Read method.
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer is the interface that wraps the basic Write method.
type Writer interface {
    Write(p []byte) (n int, err error)
}
实现了形如上面的Read方法（或者Write方法）的任意类型都可以说它实现了io.Reader接口（或者io.Writer接口）。这就意味着io.Reader接口变量能够保存任意实现了Read方法的类型所定义的值，比如：
var r io.Reader
r = os.Stdin
r = bufio.NewReader(r)
r = new(bytes.Buffer)
// and so on
明确r到底保存了什么样的具体值非常重要，但是这里r的类型却总是io.Reader：注意Go是静态类型化的，而r的静态类型是io.Reader。

一个非常非常重要的接口类型例子就是空接口：

 1
interface{}
空接口表示方法集合为空并且可以保存任意值，因为任意值都有0个或者更多方法。

有些人说Go的接口是动态类型化的，但这是一种误导。Go的接口都是静态类型化的：一个接口类型变量总是保持同一个静态类型，即使在运行时它保存的值的类型发生变化，这些值总是满足这个接口。

我们需要搞清楚上面说的这些，因为反射和接口是紧紧联系在一起的。

接口的表示（The representation of an interface)
Russ Cox曾经写过一篇关于Go中接口值的表示的博客detailed blog post。在这里我们没必要重复他的整篇文章内容了，但是简单概括下还是应该的。

一个接口类型变量存储了一个pair：赋值给这个接口变量的具体值，以及这个值的类型描述符。更进一步的讲，这个”值”是实现了这个接口的底层具体数据项（underlying concrete data item),而这个“类型”是描述了那个项（item）的全类型（full type）。举个例子，执行完下面这些：
var r io.Reader
tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
if err != nil {
  return nil, err
}
r = tty
示意性的讲，此时r就包含了(value, type)对，即(tty, os.File)。注意，除了Read方法以外，类型os.File也实现了其它方法；即使这个接口值仅仅提供了对Read方法的访问，这个接口值内部仍然带有关于这个值的全部类型信息。这就是为什么我们能干下面这些事儿：
var w io.Writer
w = r.(io.Writer)
这个赋值操作中的表达式是一个类型断言（type assertion）；它所断言的是r中存储的项（item）也实现了io.Writer接口，所以我们可以把它赋值给w。赋值操作完毕以后，w将会包含 (tty, *os.File)对。这个pair跟r中的pair是同样的。接口的静态类型决定了能用接口变量调用哪些方法，即使接口里存的具体值内部可能还有一大坨其它方法。（换句话说，接口定义的方法集合是该种接口变量所保存的具体值所含有的方法集合的一个子集，通过这个接口变量只能调用这个接口定义过的方法，没法通过这个接口变量调用其它任何方法。）

继续说，我们可以这样做：

 1
2
var empty interface{}
empty = w
我们的空接口值，empty，也能包含同样的pair即(tty, *os.File)。这样的话就很方便了，一个空接口可以保存任意值和我们所需要的关于所保存值的全部信息。

（我们不需要在这里做类型断言了，因为可以静态地知道w满足空接口。在上面那个把一个值从一个Reader移到一个Writer的例子里，我们需要显式地用一个类型断言，因为Writer定义的方法集合不是Reader定义的方法集合的子集。）

一个很重要的细节是，一个接口中的pair总有（值，具体类型）这样的格式，而不能有（值，接口类型）这样的格式。接口不能保存接口值（也就是说，你没法把一个接口变量值存储到一个接口变量中，只能把一个具体类型的值存储到一个接口变量中。）

现在，我们终于准备好了可以看看反射是怎么回事儿了。

第一反射定律（The first law of reflection)
1.从接口值到反射对象的反射（Reflection goes from interface value to reflection object）
最最基本的，反射是一种检查存储在接口变量中的（类型，值）对的机制。作为一个开始，我们需要知道reflect包中的两个类型：Type和Value。这两种类型给了我们访问一个接口变量中所包含的内容的途径，另外两个简单的函数reflect.Typeof和reflect.Valueof可以检索一个接口值的reflect.Type和reflect.Value部分。（还有就是，我们可以很容易地从reflect.Value到达reflect.Type，但是现在暂且让我们先把Value和Type的概念分开说。先剧透，从Value到达Type是通过Value中定义的某些方法来实现的，虽然先分开讲，但是后面多注意一下。）

让我们从Typeof开始：
package main

import (
    "fmt"
    "reflect"
)

func main() {
    var x float64 = 3.4
    fmt.Println("type:", reflect.TypeOf(x))
}
上面这段程序将会打印输出：
type: float64
你可能想知道我们所说的接口在上面程序哪个地方，因为这个程序看起来就是把float64类型的变量x，而不是一个接口值，传递给reflect.Typeof函数。但是，它就在那呢！就像godoc reports所描述的，reflect.Typeof 签名里就包含了一个空接口：
// TypeOf returns the reflection Type of the value in the interface{}.
func TypeOf(i interface{}) Type
当我们调用reflect.Typeof(x)的时候，x首先被保存到一个空接口中，这个空接口然后被作为参数传递。reflect.Typeof 会把这个空接口拆包（unpack）恢复出类型信息。

当然，reflect.Valueof可以把值恢复出来（从这里开始，我们将会省略这个样板而是专注与可执行代码）：
var x float64 = 3.4
fmt.Println("value:", reflect.ValueOf(x))//Valueof方法会返回一个Value类型的对象
打印出
value: <float64 Value>
reflect.Type和reflect.Value这两种类型都提供了大量的方法让我们可以检查和操作这两种类型。一个重要的例子是，Value类型有一个Type方法可以返回reflect.Value类型的Type（这个方法返回的是值的静态类型即static type，也就是说如果定义了type MyInt int64，那么这个函数返回的是MyInt类型而不是int64，看后面那个Kind方法就可以理解了）。另外一个重要的例子是，Type和Value都有一个Kind方法可以返回一个常量用于指示一个项到底是以什么形式（也就是底层类型即underlying type，继续前面括号里提到的，Kind返回的是int64而不是MyInt）存储的（what sort of item is stored)，这些常量包括：Unit, Float64, Slice等等。而且，有关Value类型的带有名字诸如Int和Float的方法可让让我们获取存在里面的值（比如int64和float64)：
var x float64 = 3.4
v := reflect.ValueOf(x)
fmt.Println("type:", v.Type())
fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
fmt.Println("value:", v.Float())
打印出
type: float64
kind is float64: true
value: 3.4
还有一些方法像SetInt和SetFloat，但是为了使用它们，我们得理解什么叫settability（在维基百科查不到，google翻译说来自于意大利语，定形的意思，觉得不是太直白，不翻译这个词了，直接用原词），下面会讨论。

反射库里有俩性质值得单独拿出来说说。第一个性质是，为了保持API简单，Value的”setter”和“getter”类型的方法操作的是可以包含某个值的最大类型：比如，所有的有符号整型，只有针对int64类型的方法，因为它是所有的有符号整型中最大的一个类型。也就是说，Value的Int方法返回的是一个int64，同时SetInt的参数类型采用的是一个int64；所以，必要时要转换成实际类型：
var x uint8 = 'x'
v := reflect.ValueOf(x)
fmt.Println("type:", v.Type())                            // uint8.
fmt.Println("kind is uint8: ", v.Kind() == reflect.Uint8) // true.
x = uint8(v.Uint())// v.Uint returns a uint64.看到啦嘛？这个地方必须进行强制类型转换！
第二个性质是，反射对象（reflection object）的Kind描述的是底层类型（underlying type），而不是静态类型（static type）。如果一个反射对象包含了一个用户定义的整型，比如：（还记得我在上面括号里举例子说明Type方法和Kind方法时说的那一坨嘛？）：
type MyInt int
var x MyInt = 7
v := reflect.ValueOf(x)
v的Kind仍然是reflect.Int，即使x的静态类型是MyInt而不是int。换句话说，Kind不能将一个int从一个MyInt中区别出来，但是Type能做到！

第二反射定律（The second law of reflection）
2.从反射队形到接口值的反射（Reflection goes from reflection object to interface value）
就像物理学上的反射，Go中到反射可以生成它的逆。

给定一个reflect.Value，我们能用Interface方法把它恢复成一个接口值；效果上就是这个Interface方法把类型和值的信息打包成一个接口表示并且返回结果：
// Interface returns v's value as an interface{}.
func (v Value) Interface() interface{}
作为一个结果，我们可以说
y := v.Interface().(float64) // y will have type float64.
fmt.Println(y)
把用反射对象v表示的float64类型的值打印了出来。

我们甚至可以做得更好一些，fmt.Println等方法的参数是一个空接口类型的值，所以我们可以让fmt包自己在内部完成我们在上面代码中做的工作。因此，为了正确打印一个reflect.Value，我们只需把Interface方法的返回值直接传递给这个格式化输出例程：
fmt.Println(v.Interface())
（为什么我们不直接fmt.Println(v)？因为v是一个reflect.Value；我们想要的是v里面保存的具体值。）因为我们的值是float64类型的，所以我们甚至可以用一个floating-point格式来打印：
   fmt.Printf("value is %7.1e\n", v.Interface())
会得到
3.4e+00
还有就是，我们不需要对v.Interface方法的结果调用类型断言（type-assert)为float64；空接口类型值内部包含有具体值的类型信息，并且Printf方法会把它恢复出来。

简要的说，Interface方法是Valueof函数的逆，除了它的返回值的类型总是interface{}静态类型。（不知道会不会有人看到前面这句话既用了方法又用了函数，会觉得奇怪。我推测，Go对于类型里面定义的都叫方法，包级别全局性的不属于任何类型的叫做函数。）

重申一遍：反射就是从接口值到反射对象，然后再反射回来。（Reflection goes from interface value to reflection object and back again.）

第三反射定律（The third law of reflection）
3.为了修改一个反射对象，值必须是settable的（To modify a reflection object, the value must be settable)
这个第三定律是最微妙最让人困惑的了，但是如果我么能从第一定律出发可以很容易的理解它。

下面是一些不能正常运行的代码，但是很值得研究：
var x float64 = 3.4
v := reflect.ValueOf(x)
v.SetFloat(7.1) // Error: will panic.
如果你运行这段代码，将会带着神秘的信息发生panic
panic: reflect.Value.SetFloat using unaddressable value
问题不是出在值7.1不是可以寻址的，而是出在v不是settable的。Settability是Value的一条性质，而且，不是所有的Value都具备这条性质。

Value的CanSet方法用与测试一个Value的settablity；在我们的例子中，
var x float64 = 3.4
v := reflect.ValueOf(x)
fmt.Println("settability of v:", v.CanSet())
输出
settability of v: false
如果对一个non-settable的Value调用Set方法会出现错误。但是，settability到底是什么呢？

settability有点像addressability，但是更加严格。settability是一个性质，描述的是一个反射对象能够修改创造它的那个实际存储的值的能力。settability由反射对象是否保存原始项（original item）而决定。当我们说
var x float64 = 3.4
v := reflect.ValueOf(x)
我们传递了x的一个副本给reflect.Valueof函数，所以作为reflect.Valueof参数被创造出来的接口值只是x的一个副本，而不是x本身。因为，如果下面这条语句
v.SetFloat(7.1)
执行成功（当然不可能执行成功啦，假设而已），它不会更新x，即使v看起来像是从x创造而来，所以它更新的只是存储在反射值内部的x的一个副本，而x本身不受丝毫影响，所以如果真这样的话，将会非常那令人困惑，而且一点用都没有！所以，这么干是非法的，而settability就是用来阻止这种哦给你非法状况出现的。

如果你觉得这个看起来有点怪的话，其实不是的，它实际上是一个披着不寻常外衣的一个你很熟悉的情况。想想下面这个把x传给一个函数：
f(x)
我们不会期待f能够修改x的值，因为我们穿了x值的一个副本，而不是x本身。如果我们想要f直接修改x，我们必须把x的地址传给这个函数（也就是说，给它传x的指针）：
f(x)
这个就很直接了，而且看起来很面熟，其实反射也是按同样的方式来运作。如果我们想通过反射来修改x，我们必须把我们想要修改的值的指针传给一个反射库。

我们来实际操作一下。首先，我们像平常一样初始化x，然后创造一个指向它的反射值，叫做p.
var x float64 = 3.4
p := reflect.ValueOf(&x) // Note: take the address of x.注意这里哦！我们把x地址传进去了！
fmt.Println("type of p:", p.Type())
fmt.Println("settability of p:", p.CanSet())
现在输出就是
type of p: *float64
settability of p: false
反射对象p不是settable的，但是我们想要设置的不是p，而是（效果上来说）*p。为了得到p指向的东西，我们调用Value的Elem方法，这样就能迂回绕过指针，同时把结果保存在叫v的Value中：
v := p.Elem()
fmt.Println("settability of v:", v.CanSet())
现在v就是一个settable的反射对象了，正如输出所描述的，
settability of v: true
并且因为v表示x，我们最终能够通过v.SetFloat方法来修改x的值：
v.SetFloat(7.1)
fmt.Println(v.Interface())
fmt.Println(x)
输出正是我们所期待的，
7.1
7.1
反射理解起来有点困难，但是它确实正在做编程语言要做的，尽管是通过掩盖了所发生的一切的反射Types和Vlues来实现的。这样好了，你就直接记住反射Values为了修改它们所表示的东西必须要有这些东西的地址。

Structs
在我们前面的例子中，v本身不是一个指针，它只是从一个指针派生来的。出现这种情况的一个常见的方法是当使用反射来修改一个structure的各个域的时候。只要我们有这个structure的地址，我们就能修改它的各个域。

下面是分析一个struct值，t，的简单例子。我们用这个struct的地址创建一个反射对象，因为我们想一会改变它的值。然后我们把typeofT变量设置为这个反射对象的类型，接着使用一些直接的方法调用（细节请见reflect包）来迭代各个域。注意，我们从struct类型中提取了各个域的名字，但是这些域本身都是rreflect.Value对象。
type T struct {
    A int
    B string
}
t := T{23, "skidoo"}
s := reflect.ValueOf(&t).Elem()
typeOfT := s.Type()//把s.Type()返回的Type对象复制给typeofT，typeofT也是一个反射。
for i := 0; i < s.NumField(); i++ {
    f := s.Field(i)//迭代s的各个域，注意每个域仍然是反射。
    fmt.Printf("%d: %s %s = %v\n", i,
        typeOfT.Field(i).Name, f.Type(), f.Interface())//提取了每个域的名字
}
这段程序的输出是：
0: A int = 23
1: B string = skidoo
关于settability还有一个要点在这里要介绍一下： 这里T的域的名字都是大写的（被导出的），因为一个struct中只有被导出的域才是settable的。

因为s包含了一个settable的反射对象，所以我们可以修改这个structure的各个域。
s.Field(0).SetInt(77)
s.Field(1).SetString("Sunset Strip")
fmt.Println("t is now", t)
结果在这里：
t is now {77 Sunset Strip}
如果我们修改这个程序，让s从t创建出来而不是&t，那么上面对SetInt和SetString的调用将会统统失败，因为t的各个域不是settable的。

总结（Conclusion）
我们在最后再次列出反射的三大定律：
1.Reflection goes from interface value to reflecton object.
2.Reflection goes from reflection object to interface value.
3.To modify a reflection object, the value must be settable.

一旦你理解了这三条反射定律，Go中的反射用起来就很简单了，虽然它还仍然有点微妙。反射是一个强大的工具，你必须要十分小心的使用它，并且应该尽量避免使用它，除非真的是不用不行了。

关于反射，仍然有大量内容我们没有讲到—-在channel中的发送操作和接收操作，分配内存，使用slices和map，调用方法和函数—但是这篇文章已经够长了。我们将来会在随后的文章中讲到前面提到的这些topics中的一些。

package main

import (
	"fmt"
)

type People interface {
	Show()
}

type Student struct{}

func (stu *Student) Show() {

}

func live() People {
	var stu *Student
	return stu
}

func main() {
	if live() == nil {
		fmt.Println("AAAAAAA")
	} else {
		fmt.Println("BBBBBBB")
	}
}

考点：interface内部结构
解答：
很经典的题！ 这个考点是很多人忽略的interface内部结构。 go中的接口分为两种一种是空的接口类似这样：

var in interface{}
另一种如题目：

type People interface {
    Show()
}
他们的底层结构如下：

type eface struct {      //空接口
    _type *_type         //类型信息
    data  unsafe.Pointer //指向数据的指针(go语言中特殊的指针类型unsafe.Pointer类似于c语言中的void*)
}
type iface struct {      //带有方法的接口
    tab  *itab           //存储type信息还有结构实现方法的集合
    data unsafe.Pointer  //指向数据的指针(go语言中特殊的指针类型unsafe.Pointer类似于c语言中的void*)
}
type _type struct {
    size       uintptr  //类型大小
    ptrdata    uintptr  //前缀持有所有指针的内存大小
    hash       uint32   //数据hash值
    tflag      tflag
    align      uint8    //对齐
    fieldalign uint8    //嵌入结构体时的对齐
    kind       uint8    //kind 有些枚举值kind等于0是无效的
    alg        *typeAlg //函数指针数组，类型实现的所有方法
    gcdata    *byte
    str       nameOff
    ptrToThis typeOff
}
type itab struct {
    inter  *interfacetype  //接口类型
    _type  *_type          //结构类型
    link   *itab
    bad    int32
    inhash int32
    fun    [1]uintptr      //可变大小 方法集合
}
可以看出iface比eface 中间多了一层itab结构。 itab 存储_type信息和[]fun方法集，从上面的结构我们就可得出，因为data指向了nil 并不代表interface 是nil， 所以返回值并不为空，这里的fun(方法集)定义了接口的接收规则，在编译的过程中需要验证是否实现接口 结果：

BBBBBBB

作为一门编程语言，对方法的处理一般分为两种类型：一是将所有方法组织在一个表格里，静态地调用（C++, java）；二是调用时动态查找方法(python, smalltalk, js)。
而go语言是两者的结合：虽然有table，但是是需要在运行时计算的table。
一个interface值由两个指针组成，第一个指向一个interface table，叫 itable。itable开头是一些描述类型的元字段，后面是一串方法。注意这个方法是interface本身的方法，并非其dynamic value（Binary）的方法
当这个interface无方法时，itable可以省略，直接指向一个type即可。
另一个指针data指向dynamic value的一个拷贝，这里则是b的一份拷贝。也就是，给interface赋值时，会在堆上分配内存，用于存放拷贝的值。
同样，当值本身只有一个字长时，这个指针也可以省略。
var w io.Writer
这时，tab和data都是nil。interface是否为nil取决于itable字段。所以不一定data为nil就是nil，判断时要额外注意。
 switch v := any.(type) {
case int:
    return strconv.Itoa(v)
case float:
    return strconv.Ftoa(v, 'g', -1)
}
 
实际上是any这个interface取了  any. tab->type。
 
而interface的函数调用实际上就变成了：
s.tab->fun[0](s.data)。第一个参数即自身类型指针

itable的生成是理解interface的关键。
如刚开始处提的，为了支持go语言这种接口间仅通过方法来联系的特性，是没有办法像C++一样，在编译时预先生成一个method table的，只能在运行时生成。因此，自然的，所有的实体类型都必须有一个包含此类型所有方法的“类型描述符”(type description structure)；而interface类型也同样有一个类似的描述符，包含了所有的方法。
这样，interface赋值时，计算interface对象的itable时，需要对两种类型的方法列表进行遍历对比。如后面代码所示，这种计算只需要进行一次，而且优化成了O(m+n)。
 
可见，interface与itable之间的关系不是独立的，而是与interface具体的value类型有关。即（interface类型， 具体类型）->itable。
reflection实质上是将interface背后的实现暴露了一部分给应用代码，使应用程序可以使用interface实现的一些内容。只要理解了interface的实现，reflect就好理解了。如reflect.typeof(i)返回interface i的type，Valueof返回value.
