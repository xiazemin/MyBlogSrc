---
title: govaluate
layout: post
category: golang
author: 夏泽民
---
https://github.com/Knetic/govaluate
govaluate提供了任意类似C语言的算术/字符串表达式的求值。

为什么你不应该直接在代码中书写表达式
有些时候，你并没有办法提前得知表达式的样子，或者你希望表达式可设置。如果你有一堆运行在你的应用上的数据，或者你想要允许你的用户自定义一些内容，或者你写的是一个监控框架，可以获得很多metrics信息，然后进行一些公式计算，那么这个库就会非常有用。

如何使用
可以创建一个新的EvaluableExpression，然后调用它的”Evaluate”方法。

   expression, err := govaluate.NewEvaluableExpression("10 > 0");
result, err := expression.Evaluate(nil);
// result is now set to "true", the bool value.
那么，如何使用参数？

expression, err := govaluate.NewEvaluableExpression("foo > 0");
parameters := make(map[string]interface{}, 8)
parameters["foo"] = -1;
result, err := expression.Evaluate(parameters);
// result is now set to "false", the bool value.
这很棒，但是这些基本上可以使用代码直接实现。那么如果计算中牵扯到一些数学计算呢？

expression, err := govaluate.NewEvaluableExpression("(requests_made * requests_succeeded / 100) >= 90");
parameters := make(map[string]interface{}, 8)
parameters["requests_made"] = 100;
parameters["requests_succeeded"] = 80;
result, err := expression.Evaluate(parameters);
// result is now set to "false", the bool value.
上述例子返回的都是布尔值，事实上，它是可以返回数字的。

expression, err := govaluate.NewEvaluableExpression("(mem_used / total_mem) * 100");
parameters := make(map[string]interface{}, 8)
parameters["total_mem"] = 1024;
parameters["mem_used"] = 512;
result, err := expression.Evaluate(parameters);
// result is now set to "50.0", the float64 value.
你也可以做一些日期的转化，只要符合RF3339,ISO8061,Unix Date，或者ruby日期格式标准即可。如果你还是不太确定，那么可以看一下支持的日期标准。

   expression, err := govaluate.NewEvaluableExpression("'2014-01-02' > '2014-01-01 23:59:59'");
result, err := expression.Evaluate(nil);
// result is now set to true
表达式只需要进行一次句法分析，就可以多次复用。

   expression, err := govaluate.NewEvaluableExpression("response_time <= 100");
parameters := make(map[string]interface{}, 8)
for {
	parameters["response_time"] = pingSomething();
	result, err := expression.Evaluate(parameters)
}
关于执行顺序，本库支持正常C标准的执行顺序。编写表达式时，请确保您正确地书写操作符，或使用括号来明确表达式的哪些部分应先运行。

govaluate采用\或者[]来完成转义。

支持自定义函数

支持简单的结构体（访问器）

运算符支持
ruleplatform的表达式引擎支持以下运算：
二元计算符 : + - / & | ^ * % >> <<
二元比较符 : > >= < <= == != =~ !~
逻辑操作符 : || &&
括号 : ( )
数组相关 : , IN (例子1 IN (1, 2, ‘foo’)，返回值true)
一元计算符 : ! - ~
三元运算符 : ? :
空值聚合符: ??
<!-- more -->
https://segmentfault.com/a/1190000022235609
为什么你不能在代码中写这些表达式？
有时，你不能知道ahead-of-time是什么样的表达式，或者你希望这些表达式可以配置。 也许你已经通过应用程序运行了一组数据，希望用户在向数据库提交之前指定一些验证。 或者你已经编写了一个监视框架，可以收集一些指标，然后评估一些表达式。

许多人都在写自己的一半烘焙风格的评估语言，但是不完整。 或者他们会把表达式烘焙到实际的可以执行文件中，即使他们知道它有可以能更改。 这些策略可能会工作，但是他们需要时间来实现，用户学习的时间和需求变化。 这个库是用来覆盖所有正常的c 类似的表达式，所以你不必在计算机上重新创建一个最旧的轮子。

我怎么用它？
创建一个新的EvaluableExpression，然后在它上面调用"评估"。

复制
expression, err:= govaluate.NewEvaluableExpression("10> 0");
 result, err:= expression.Evaluate(nil);
 // result is now set to"true", the bool value.
酷，但参数如何？

复制
expression, err:= govaluate.NewEvaluableExpression("foo> 0");
 parameters:=make(map[string]interface{}, 8)
 parameters["foo"] = -1;
 result, err:= expression.Evaluate(parameters);
 // result is now set to"false", the bool value.
很酷但我们几乎可以在代码中完成。 一个复杂的用例涉及一些数学？

复制
expression, err:= govaluate.NewEvaluableExpression("(requests_made * requests_succeeded/100)> = 90");
 parameters:=make(map[string]interface{}, 8)
 parameters["requests_made"] = 100;
 parameters["requests_succeeded"] = 80;
 result, err:= expression.Evaluate(parameters);
 // result is now set to"false", the bool value.
或者你想检查一个活动检查("smoketest") 页的状态，这是一个字符串？

复制
expression, err:= govaluate.NewEvaluableExpression("http_response_body == 'service is ok'");
 parameters:=make(map[string]interface{}, 8)
 parameters["http_response_body"] = "service is ok";
 result, err:= expression.Evaluate(parameters);
 // result is now set to"true", the bool value.
这些示例都返回了布尔值，但同样可能返回数值值。

复制
expression, err:= govaluate.NewEvaluableExpression("(mem_used/total_mem) * 100");
 parameters:=make(map[string]interface{}, 8)
 parameters["total_mem"] = 1024;
 parameters["mem_used"] = 512;
 result, err:= expression.Evaluate(parameters);
 // result is now set to"50.0", the float64 value.
你还可以执行日期解析，尽管格式有些有限。 坚持 RF3339，ISO8061，unix日期或者 ruby 日期格式。 如果你在获取日期字符串时遇到问题，请检查实际使用的格式列表： 解析。转到：248.

复制
expression, err:= govaluate.NewEvaluableExpression("'2014-01-02'> '2014-01-01 23:59:59'");
 result, err:= expression.Evaluate(nil);
 // result is now set to true
表达式被解析一次，并且可以多次使用。 解析是流程的计算密集阶段，因这里如果你打算使用同样的表达式，只需解析一次。 像这样

复制
expression, err:= govaluate.NewEvaluableExpression("response_time <= 100");
 parameters:=make(map[string]interface{}, 8)
 for {
 parameters["response_time"] = pingSomething();
 result, err:= expression.Evaluate(parameters)
 }
遵守正常的c 标准顺序。 在编写表达式时，请确保要正确排序操作符，或者使用括号来阐明表达式的哪些部分。

转义符
有时你会有一些参数，有空格。斜线。优点。符号或者它的他特殊的特征。 例如下面的表达式将不会像预期的那样执行：

复制

"response-time <100"


按照编写的方式，库将把它解析为"。[response] 减去 [time] 小于 100"。 实际上，"响应时间"应该是一个变量，它只有一个折线。

有两种方法可以解决这个问题。 首先，你可以转义整个参数名：

复制

"[response-time] <100"


或者你可以使用反斜线来只转义减号。

复制

"response-time <100"


可以在表达式中的任意位置使用反斜杠来转义下一个字符。 可以在任何时候使用方括号参数名代替普通参数名。

用户定义函数
你可能在执行表达式时希望在参数上调用函数。 也许你希望聚合一些数据集，但不知道在编写表达式本身之前要使用的确切聚合。 或者你有一个数学操作，你需要执行，因为没有操作符，如 log 或者 tan 或者 sqrt。 对于这种情况，你可以提供一个函数映射 NewEvaluableExpressionWithFunctions 在执行过程中，它将能够使用它们。 比如；

复制
functions:=map[string]govaluate.ExpressionFunction {
 "strlen": func(args.. .interface{}) (interface{}, error) {
 length:=len(args[0].(string))
 return (float64)(length), nil },
 }
 expString:="strlen('someReallyLongInputString') <= 16"expression, _:= govaluate.NewEvaluableExpressionWithFunctions(expString, functions)
 result, _:= expression.Evaluate(nil)
 // result is now"false", the boolean value
函数可以接受任意数量的参数，正确处理嵌套函数，并且参数可以是任何类型的( 即使这些库的运算符都不支持该类型的评估)。 例如表达式中函数的每个用法都是有效的( 假设给定的函数和参数是正确的):

复制
"sqrt(x1 ** y1, x2 ** y2)""max(someValue, abs(anotherValue), 10 * lastValue)"
函数不能作为参数传递，它们必须在解析表达式时知道，并且在解析后是不可变的。

访问器
如果参数中有结构，则可以按常规方式访问它们的字段和方法。 例如给定具有方法"回音"的结构，将它的作为 foo 存在于参数中，则以下是有效的：

复制

"foo.Echo('hello world')"


以类似方式访问字段。 假设 foo 有一个名为"长度"的字段：

复制

"foo.Length> 9000"


访问器可以嵌套到任何深度，如下所示

复制

"foo.Bar.Baz.SomeFunction()"


然而，目前它并不支持访问 map的值。 因此，以下将不工作

复制

"foo.SomeMap['key']"


这可能很方便，但是请注意，使用访问器涉及的是反射的 。 这使得表达式的速度比使用参数(。请参考基准测试以更精确地测量你的系统) 慢4 倍。 如果是合理的，作者建议你预先提取到参数映射中的值，或者定义实现 Parameters 接口的结构。 如果有函数需要使用，最好将它们作为表达式函数( 查看上面的部分) 传递。 这些方法不使用反射，并且设计得非常快速和干净。

支持哪些运算符和类型？
修改器：+-/*&|^**%>> TimeoutException
比较器：>>=<<===!==~!~
逻辑运算：||&&
数字常量，如 64位 浮点( 12345.678 )
字符串常量( 单引号：'foobar' )
日期常量( 单引号，使用 RFC3339.ISO8601.ruby 日期或者unix日期) ；日期解析自动使用任何字符串常量进行尝试
布尔常量：truefalse
圆括号以控制评估 ()的顺序
数组( 括号内由 , 分隔的任何内容： (1, 2,'foo') )
前缀： -~!
三元条件：?:
零合并：??
有关每个运算符支持哪些类型的严格细节，请参见 MANUAL.md。

用户定义类型
某些运算符在使用某些类型时没有意义。 例如获得字符串的模是什么意思？ 如果你检查两个数字是逻辑的还是ed的，会发生什么？

每个人对这些问题的答案都有不同的直觉。 为防止混淆，本库将拒绝操作，对操作没有明确意义的类型。 有关运算符对哪些类型有效的详细信息，请参见 MANUAL.md。

基准测试
如果你关心这个库的开销，那么在这个 repo 中会有很好的基准测试范围。 你可以用 go test -bench=. 运行它们。 图书馆是以快速的方式构建的，但并没有被积极地分析和优化。 但是对于大多数应用来说，它完全可以。

这是一个非常粗略的性能概念，这是来自于 3一代 Macbook Pro ( Linux Mint 17.1 )的基准运行的结果。


