---
title: antlr hive
layout: post
category: lang
author: 夏泽民
---
hive是使用antlr来解析的

parser要做的事情，是从无结构的字符串里面，解码产生有结构的数据结构（a parser is a function accepting strings as input and returning some structure as output），参考 Parser_combinator wiki

parser分成两种，一种是parser combinator，一种是parser generator

parser combinator是需要手写parser，a parser combinator is a higher-order function that accepts several parsers as input and returns a new parser as its output

parser generator是需要你用某种指定的描述语言来表示出语法，然后自动把他们转换成parser的代码，比如Antlr里面的g4语法文件，以及calcite的ftl语法文件，缺点是由于代码是生成的，排错比较困难

使用了Antlr的parser有Hive，Presto，Spark SQL

美团点评的文章

1
https://tech.meituan.com/2014/02/12/hive-sql-to-mapreduce.html
以及hive源码的测试用例

1
https://github.com/apache/hive/blob/branch-1.1/ql/src/test/org/apache/hadoop/hive/ql/parse/TestHiveDecimalParse.java
hive的g4文件如下

老版本的hive

1
https://github.com/apache/hive/blob/59d8665cba4fe126df026f334d35e5b9885fc42c/parser/src/java/org/apache/hadoop/hive/ql/parse/HiveParser.g
新版本的hive

1
https://github.com/apache/hive/blob/master/hplsql/src/main/antlr4/org/apache/hive/hplsql/Hplsql.g4
spark的g4文件如下

1
https://github.com/apache/spark/blob/master/sql/catalyst/src/main/antlr4/org/apache/spark/sql/catalyst/parser/SqlBase.g4
Presto的g4文件如下

1
https://github.com/prestodb/presto/blob/master/presto-parser/src/main/antlr4/com/facebook/presto/sql/parser/SqlBase.g4
使用了Apache Calcite的parser有Apache Flink，Mybatis，Apache Storm等

https://tech.meituan.com/2014/02/12/hive-sql-to-mapreduce.html

https://www.cnblogs.com/tonglin0325/p/12212866.html

https://www.yinwang.org/blog-cn/2015/09/19/parser
<!-- more -->
用python+antlr解析hive sql获得数据血缘关系（一）
系列目标
编程获得数据血缘关系的需求对数据仓库来说并不普遍，只有数据规模达到很大的程度，或者存在复杂数据生产关系的报表数量增加到很大的程度，单独的数据血缘关系工作才有必要。
在规模没达到之前，人工的识别和管理是更经济有效的。

本系列想要做到的目标是这个uber的 queryparser的一个子集，在有限知道目标数据表结构的前提下，发现并记录目标字段与来源表和字段的关系。

这种功能queryparser应该是已经具备的，并且它本身是开源的，但queryparser的主体是Haskell写的，为这么一个边缘功能学门新的编程范式，学习代价太大了点。

还是选择python作为开发工具比较靠谱。

可选项比较
自己从头写字符串处理是不可能的，就算是用正则辅助，搞那些语法边角的工作量也难以估计。

于是祭出搜索大法，在各处寻摸一遍后，拿到了这么几个可能的选择项：

queryparser
就是前面说的uber放出的开源项目，因为编程语言的壁垒，最早放弃。

sqlparse
pypi上可以搜索到的模块，github地址https://github.com/andialbrecht/sqlparse
网上也有一些材料，

拿来做了简单试验后，放弃。

放弃主要原因是因为它的功能集合相比要做的hive sql解析，感觉太小了。sqlparse从sql语句解析出来的是 statements tuple，每个statement上会有一个识别出的类型，而在我要解析的sql集合里，大概有三分之一sql语句，识别出的statement类型是UNKNOWN，这个比例太大不能接受。

pyparsing
也是pypi上可以搜索到的模块，github地址https://github.com/pyparsing/pyparsing/ 这是python版本的通用解析工具。

如果有人基于这个pyparsing做过hive sql解析就好了，然而没有。如果要用pyparsing，就要从头写语法文件。python项目用它做表达式解析，或者做新配置语法还好，用来解析hive sql这种量级的，工作量也太大，放弃。

antlr
在找到pyparsing时我已经同时在找antlr相关信息了，因为要解析hive sql，最权威的解析器肯定是hive自己用的那个，经过确认，这个工具就是antlr，更具体的说，是antlr 3系列。

antlr自己的历史不是本系列重点，感兴趣的可以自行到https://www.antlr.org/上去查阅

grammar文件
要用hive自身的解析，就要拿到hive的语法文件定义，对于开源的hive来说，这个事还是挺容易的，github上可以很容易按版本访问到历史文件，以hive 1.1.0版本的文件为例，语法文件定义所在的文件夹是
https://github.com/apache/hive/tree/release-1.1.0/ql/src/java/org/apache/hadoop/hive/ql/parse

网上也提到过，hive的语法文件经历过分拆，在1.1.0版本中，一共有5个文件，都是.g后缀名，分别是

FromClauseParser.g
HiveLexer.g
HiveParser.g
IdentifiersParser.g
SelectClauseParser.g
把它们从github上下载回来，或者从页面上复制粘贴到编辑器里，再保存为对应名字的文本文件也可以，主要文件名要严格一样，antlr对文件名和语法文件内容有检查。

antlr版本
antlr有 v2 v3 v4多个版本并存，中文文档多数是v2的， hive 1.1.0版本在注释中提到了antlr 3.4，最新的3.x版本是3.5.2，我选择用3.4版本的。

antlr自己发布的下载地址是在github 上的 ，能下载，但由于美帝的"封锁"，下载速度很难保证，
这是另外一个能下载antlr的网站

java版本
antlr的运行需要jdk版本的java，具体是jdk里面的javac来编译代码，具体支持的最低版本我没确认。

如果环境里没有javac，需要自行安装，在centos下可以这样安装jdk8

# yum install -y java-1.8.0-openjdk-devel
1
pyjnius
antlr这个工具，可以产出多种target语言，这其中也包括了python，不过查阅列表后发现，对python target的有效性验证只持续到antlr 3.1.3 ，到antlr 3.4 版本就很难讲了。继续一番搜索
大法，决定使用pyjnius作为python到java之间的桥梁。

桥接python和java的方案,更具体来说，是在python里调用java代码的方案，其实也有好几个，

pyjnius
Jpype
javabridge
py4j
jcc
其他的方案其实我都还没试，pyjnius的尝试几乎是一次通过，就优先选这个了

安装
通过pip安装，过程很顺利
如果网速慢，推荐使用清华的镜像

pip3 install -i https://pypi.tuna.tsinghua.edu.cn/simple pyjnius
1
如果之前没有装过cython，会连带需要安装这个依赖。
安装cython可能要连带安装gcc

yum install gcc gcc-++
pip3 install -i https://pypi.tuna.tsinghua.edu.cn/simple cython
1
2
编译grammar
前面下载的5个hive grammar文件，最好集中保存到一个目录下，推荐命名为带层次的 grammar/hive110， 原因是grammar文件经过antlr解析后会生成对应的java源代码，然后要再编译为class文件，而java在搜索类时会根据package名和路径名做匹配，一开始就做下目录规划，可以节省后续扩展调整的功夫。

修订 HiveLexer.g
在编译前，要修改一下HiveLexer.g文件，否则会因为找不到hive的相关文件报错

/**
注释掉下面这两段
@lexer::header {
package org.apache.hadoop.hive.ql.parse;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.hive.conf.HiveConf;
}

@lexer::members {
  private Configuration hiveConf;

  public void setHiveConf(Configuration hiveConf) {
    this.hiveConf = hiveConf;
  }

  protected boolean allowQuotedId() {
    String supportedQIds = HiveConf.getVar(hiveConf, HiveConf.ConfVars.HIVE_QUOTEDID_SUPPORT);
    return !"none".equals(supportedQIds);
  }
}
增加下面这段
*/
@lexer::header {
package grammar.hive110;
}
/*
中间部分省略
下面这行要修改
    | {allowQuotedId()}? QuotedIdentifier
*/
    | {true}? QuotedIdentifier

@lexer::header和@lexer::member都是会被antlr添加到目标文件里的内容，注释掉的部分是对hive里其他部分的引用，因为我只需要lexer和parser，其他部分就不要了。
注释掉的内容里有一个allowQuotedId()的方法，语法文件里有对它的调用，也要一起修改掉。

编译HiveLexer.g成.java文件
需要有下载好的antlr的jar文件，我把它放到和.g文件同一目录下，执行

$ java -jar antlr-3.4-complete.jar HiveLexer.g
1
顺利的话，应该没有报错信息，并且生成HiveLexer.tokens和HiveLexer.java文件。

这个HiveLexer.java，就是antlr自动生成的，可以处理hive sql词法规则的源代码。

词法规则只校验"sql应该怎么写"，处理这部分工作的程序一般叫lexer
还有另外一半"sql应该怎么执行"的工作，一般由叫parser的程序做，也就是前面的HiveParser.g

编译.java文件到.class文件
下一步是把生成的java源代码编译成.class的字节码

$ javac -cp antlr-3.4-complete.jar HiveLexer.java
1
顺利的话，也应该没有报错信息，并且生成HiveLexer.java文件，可能还会同时生成 HiveLexerDFA25.class,HiveLexerDFA25.class, HiveLexerDFA25.class,HiveLexerDFA21.class 这样的文件。

编写测试代码
从pyjnius和antlr的示例代码两个各取一部分，移花接木一番，得到的简单测试代码如下,

保存为antlrtest.py

#antlrtest.py
import jnius_config
jnius_config.set_classpath('./','./grammar/hive110/antlr-3.4-complete.jar')
import jnius
StringStream = jnius.autoclass('org.antlr.runtime.ANTLRStringStream')
Lexer  = jnius.autoclass('grammar.hive110.HiveLexer')
TokenStream  = jnius.autoclass('org.antlr.runtime.CommonTokenStream')

cstream = StringStream("select * from new_table;")
inst = Lexer(cstream)
ts = TokenStream()
ts.setTokenSource(inst)
ts.fill()

jlist = ts.getTokens()
tsize = jlist.size()
for i in range(tsize):
    print(jlist.get(i).getText())

确认目录结构
上面的代码是经过多次调试才得到的，一次运行可能不成功。最有可能的是目录结构不对。能成功运行的目录结构是这样的

antlrtest.py
grammar/
└── hive110
    ├── antlr-3.4-complete.jar
    ├── FromClauseParser.g
    ├── HiveLexer.class
    ├── HiveLexer$DFA21.class
    ├── HiveLexer$DFA25.class
    ├── HiveLexer.g
    ├── HiveLexer.java
    ├── HiveLexer.tokens
    ├── HiveParser.g
    ├── IdentifiersParser.g
    ├── SelectClauseParser.g

运行结果和代码简介
antlrtest.py的输出应该是下面这样的多行文本，每行是被HiveLexer这个类识别出的一个Token

select

*

from

new_table
;
<EOF>

用注释解释代码用途

#antlrtest.py
import jnius_config
#这里是设置java的classpath，必须在import jnius之前做，设置后进程内不能修改了
jnius_config.set_classpath('./','./grammar/hive110/antlr-3.4-complete.jar')
import jnius
#这3个是利用autoclass的自动装载，把java里的类定义反射到python里
#StringStream对应的是是HiveLexer构造函数必须的输入参数类型之一，ANTLRStringStream
StringStream = jnius.autoclass('org.antlr.runtime.ANTLRStringStream')
#注意这里的类名，和前面.g文件里定义的package名要有对应，和HiveLexer.class所在的目录也要有对应
Lexer  = jnius.autoclass('grammar.hive110.HiveLexer')
#TokenStream是要取出token时，保存token的容器类型CommonTokenStream
TokenStream  = jnius.autoclass('org.antlr.runtime.CommonTokenStream')

cstream = StringStream("select * from new_table;")
inst = Lexer(cstream)
ts = TokenStream()
# antlr 3增加的步骤，Lexer和Parser之间用CommonTokenStream为接口
ts.setTokenSource(inst)
#调用fill来消费掉cstream里的所有token
ts.fill()

# jlist不能直接在python里迭代，
jlist = ts.getTokens()
tsize = jlist.size()
for i in range(tsize):
    print(jlist.get(i).getText())
    
 https://blog.csdn.net/bigdataolddriver/article/details/103826702
 
 https://eng.uber.com/queryparser/
 
 https://github.com/uber/queryparser
 
 https://github.com/andialbrecht/sqlparse
 
 https://github.com/pyparsing/pyparsing/
 
 SQL转化为MapReduce的过程
了解了MapReduce实现SQL基本操作之后，我们来看看Hive是如何将SQL转化为MapReduce任务的，整个编译过程分为六个阶段：

Antlr定义SQL的语法规则，完成SQL词法，语法解析，将SQL转化为抽象语法树AST Tree
遍历AST Tree，抽象出查询的基本组成单元QueryBlock
遍历QueryBlock，翻译为执行操作树OperatorTree
逻辑层优化器进行OperatorTree变换，合并不必要的ReduceSinkOperator，减少shuffle数据量
遍历OperatorTree，翻译为MapReduce任务
物理层优化器进行MapReduce任务的变换，生成最终的执行计划

https://www.cnblogs.com/yaojingang/p/5446310.html

https://blog.csdn.net/fover717/article/details/69367545
https://www.it610.com/article/4610072.htm
https://www.jianshu.com/p/1a09ead6df21

http://ixiaosi.art/2019/01/28/hive/hive-sql%E6%89%A7%E8%A1%8C%E6%B5%81%E7%A8%8B%E5%88%86%E6%9E%90/

https://code-examples.net/zh-CN/keyword/122735

https://www.xuebuyuan.com/2181078.html
http://www.itkeyword.com/doc/2100804199580563x913
https://www.icode9.com/content-2-615877.html

https://www.cnblogs.com/drawwindows/p/4584326.html
https://blog.csdn.net/bigdataolddriver/article/details/103867682
https://blog.csdn.net/bigdataolddriver/article/details/104000719

ANTLR是一款强大的语法分析器生成工具，可用于读取、处理、执行和翻译结构化的文本或二进制文件。它被广泛应用于学术领域和工业生产实践，是众多语言、工具和框架的基石。Twitter搜索使用ANTLR进行语法分析，每天处理超过20亿次查询；Hadoop生态系统中的Hive、Pig、数据仓库和分析系统所使用的语言都用到了ANTLR；Lex Machina将ANTLR用于分析法律文本；Oracle公司在SQL开发者IDE和迁移工具中使用了ANTLR；NetBeans公司的IDE使用ANTLR来解析C++；Hibernate对象-关系映射框架（ORM）使用ANTLR来处理HQL语言。

除了这些鼎鼎大名的项目之外，还可以利用ANTLR构建各种各样的实用工具，如配置文件读取器、遗留代码转换器、维基文本渲染器，以及JSON解析器。我编写了一些工具，用于创建数据库的对象-关系映射、描述三维可视化以及在Java源代码中插入性能监控代码。我甚至为一次演讲编写了一个简单的DNA模式匹配程序。

一门语言的正式描述称为语法（grammar），ANTLR能够为该语言生成一个语法分析器，并自动建立语法分析树——一种描述语法与输入文本匹配关系的数据结构。ANTLR也能够自动生成树的遍历器，这样你就可以访问树中的节点，执行自定义的业务逻辑代码。

本书既是ANTLR 4的参考手册，也是解决语言识别问题的指南。你会学到如下知识：

识别语言样例和参考手册中的语法模式，从而编写自定义的语法。

循序渐进地为从简单的JSON到复杂的R语言编写语法。同时还能学会解决XML和Python中棘手的识别问题。

基于语法，通过遍历自动生成的语法分析树，实现自己的语言类应用程序。

在特定的应用领域中，自定义识别过程的错误处理机制和错误报告机制。

通过在语法中嵌入Java动作（action），对语法分析过程进行完全的掌控。

本书并非教科书，所有的讨论都是基于实例的，旨在令你巩固所学的知识，并提供语言类应用程序的基本范例。

https://www.jb51.net/books/634176.html