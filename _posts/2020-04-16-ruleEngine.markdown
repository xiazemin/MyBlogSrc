---
title: Java各种规则引擎
layout: post
category: web
author: 夏泽民
---
一. Drools规则引擎
简介：
Drools就是为了解决业务代码和业务规则分离的引擎。
Drools 规则是在 Java 应用程序上运行的，其要执行的步骤顺序由代码确定
，为了实现这一点，Drools 规则引擎将业务规则转换成执行树。
特性：
优点：
　　　1、简化系统架构，优化应用
　　　2、提高系统的可维护性和维护成本
　　　3、方便系统的整合
　　　4、减少编写“硬代码”业务规则的成本和风险
3.原理：
<!-- more -->
	<img src="{{site.url}}{{site.baseurl}}/img/dtools.webp"/>
	使用方式：
(1)Maven 依赖：

<dependencies>
    <dependency>
        <groupId>org.kie</groupId>
        <artifactId>kie-api</artifactId>
        <version>6.5.0.Final</version>
    </dependency>
    <dependency>
        <groupId>org.drools</groupId>
        <artifactId>drools-compiler</artifactId>
        <version>6.5.0.Final</version>
        <scope>runtime</scope>
    </dependency>
    <dependency>
        <groupId>junit</groupId>
        <artifactId>junit</artifactId>
        <version>4.12</version>
    </dependency>
</dependencies>
（2）新建配置文件/src/resources/META-INF/kmodule.xml

<?xml version="1.0" encoding="UTF-8"?>
<kmodule xmlns="http://jboss.org/kie/6.0.0/kmodule">
    <kbase name="rules" packages="rules">
        <ksession name="myAgeSession"/>
    </kbase>
</kmodule>
（3）新建drools规则文件/src/resources/rules/age.drl

import com.lrq.wechatDemo.domain.User               // 导入类

dialect  "mvel"

rule "age"                                      // 规则名，唯一
    when
        $user : User(age<15 || age>60)     //规则的条件部分
    then
        System.out.println("年龄不符合要求！");
end
工程搭建完毕，效果如图：

项目结构.png
测试用例：


/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/26
 */
@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations = {"classpath*:applicationContext.xml"})
public class TestUser {

    private static KieContainer container = null;
    private KieSession statefulKieSession = null;

    @Test
    public void test(){
        KieServices kieServices = KieServices.Factory.get();
        container = kieServices.getKieClasspathContainer();
        statefulKieSession = container.newKieSession("myAgeSession");
        User user = new User("duval yang",12);
        statefulKieSession.insert(user);
        statefulKieSession.fireAllRules();
        statefulKieSession.dispose();

    }



}

二.Aviator表达式求值引擎
简介：
Aviator是一个高性能、轻量级的java语言实现的表达式求值引擎，主要用于各
种表达式的动态求值。现在已经有很多开源可用的java表达式求值引擎，为什
么还需要Avaitor呢？

Aviator的设计目标是轻量级和高性能 ，相比于Groovy、JRuby的笨重，Aviator
非常小，加上依赖包也才450K,不算依赖包的话只有70K；当然，Aviator的语法
是受限的，它不是一门完整的语言，而只是语言的一小部分集合。

其次，Aviator的实现思路与其他轻量级的求值器很不相同，其他求值器一般都
是通过解释的方式运行，而Aviator则是直接将表达式编译成Java字节码，交给
JVM去执行。简单来说，Aviator的定位是介于Groovy这样的重量级脚本语言和
IKExpression这样的轻量级表达式引擎之间。
特性：
（1）支持大部分运算操作符，包括算术操作符、关系运算符、逻辑操作符、
正则匹配操作符(=~)、三元表达式?: ，并且支持操作符的优先级和括号强制优
先级，具体请看后面的操作符列表。
（2）支持函数调用和自定义函数。
（3）支持正则表达式匹配，类似Ruby、Perl的匹配语法，并且支持类Ruby的
$digit指向匹配分组。自动类型转换，当执行操作的时候，会自动判断操作数类
型并做相应转换，无法转换即抛异常。
（4）支持传入变量，支持类似a.b.c的嵌套变量访问。
（5）性能优秀。
（6）Aviator的限制，没有if else、do while等语句，没有赋值语句，仅支持逻
辑表达式、算术表达式、三元表达式和正则匹配。没有位运算符

整体结构：


整体结构.png
maven依赖：

<dependency>
    <groupId>com.googlecode.aviator</groupId>
    <artifactId>aviator</artifactId>
    <version>${aviator.version}</version>
</dependency>
执行方式
执行表达式的方法有两个：execute()、exec();
execute()，需要传递Map格式参数
exec(),不需要传递Map
示例：

/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/25
 */
public class Test {
    public static void main(String[] args) {
        // exec执行方式，无需传递Map格式
        String age = "18";
        System.out.println(AviatorEvaluator.exec("'His age is '+ age +'!'", age));



        // execute执行方式，需传递Map格式
        Map<String, Object> map = new HashMap<String, Object>();
        map.put("age", "18");
        System.out.println(AviatorEvaluator.execute("'His age is '+ age +'!'", 
map));

    }
}
使用函数
Aviator可以使用两种函数：内置函数、自定义函数
（1）内置函数


Aviator内置函数.png
Aviator内置函数.png

/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/25
 */
public class Test {
    public static void main(String[] args) {
        Map<String,Object> map = new HashMap<>();
        map.put("s1","123qwer");
        map.put("s2","123");

  System.out.println(AviatorEvaluator.execute("string.startsWith(s1,s2)",map));

    }
}


（2）自定义函数

自定义函数要继承AbstractFunction类，重写目标方法。


/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/25
 */
public class Test {
    public static void main(String[] args) {
        // 注册自定义函数
        AviatorEvaluator.addFunction(new MultiplyFunction());
        // 方式1
        System.out.println(AviatorEvaluator.execute("multiply(12.23, -2.3)"));
        // 方式2
        Map<String, Object> params = new HashMap<>();
        params.put("a", 12.23);
        params.put("b", -2.3);
        System.out.println(AviatorEvaluator.execute("multiply(a, b)", params));
    }

}

class MultiplyFunction extends AbstractFunction{
    @Override
    public AviatorObject call(Map<String, Object> env, AviatorObject arg1, AviatorObject arg2) {

        double num1 = FunctionUtils.getNumberValue(arg1, env).doubleValue();
        double num2 = FunctionUtils.getNumberValue(arg2, env).doubleValue();
        return new AviatorDouble(num1 * num2);
    }

    @Override
    public String getName() {
        return "multiply";
    }

}
常用操作符的使用
（1）操作符列表


操作符列表.png
（2）常量和变量


常量和变量.png
（3）编译表达式


/**
* CreateBy: haleyliu
* CreateDate: 2018/12/25
*/
public class Test {
   public static void main(String[] args) {
       String expression = "a+(b-c)>100";
       // 编译表达式
       Expression compiledExp = AviatorEvaluator.compile(expression);

       Map<String, Object> env = new HashMap<>();
       env.put("a", 100.3);
       env.put("b", 45);
       env.put("c", -199.100);

       // 执行表达式
       Boolean result = (Boolean) compiledExp.execute(env);
       System.out.println(result);

   }
}

(4) 访问数组和集合
List和数组用list[0]和array[0]，Map用map.date

/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/25
 */
public class Test {
    public static void main(String[] args) {

        final List<String> list = new ArrayList<>();
        list.add("hello");
        list.add(" world");

        final int[] array = new int[3];
        array[0] = 0;
        array[1] = 1;
        array[2] = 3;

        final Map<String, Date> map = new HashMap<>();
        map.put("date", new Date());

        Map<String, Object> env = new HashMap<>();
        env.put("list", list);
        env.put("array", array);
        env.put("map", map);

        System.out.println(AviatorEvaluator.execute(
                "list[0]+':'+array[0]+':'+'today is '+map.date", env));

    }

}
（5） 三元比较符

/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/25
 */
public class Test {
    public static void main(String[] args) {

        Map<String, Object> env = new HashMap<String, Object>();
        env.put("a", -5);
        String result = (String) AviatorEvaluator.execute("a>0? 'yes':'no'", env);
        System.out.println(result);
    }

}
(6) 正则表达式匹配


/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/25
 */
public class Test {
    public static void main(String[] args) {
        String email = "hello2018@gmail.com";
        Map<String, Object> env = new HashMap<String, Object>();
        env.put("email", email);
        String username = (String) AviatorEvaluator.execute("email=~/([\\w0-8]+)@\\w+[\\.\\w+]+/ ? $1 : 'unknow' ", env);
        System.out.println(username);
    }
}
(7) 变量的语法糖衣


/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/25
 */
public class Test {
    public static void main(String[] args) {
        User user = new User(1,"jack","18");
        Map<String, Object> env = new HashMap<>();
        env.put("user", user);

        String result = (String) AviatorEvaluator.execute(" '[user id='+ user.id + ',name='+user.name + ',age=' +user.age +']' ", env);
        System.out.println(result);
    }
}


/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/25
 */
public class User {

    private int id;

    private String name;

    private String age;

    public User() {
    }

    public User(int id, String name, String age) {
        this.id = id;
        this.name = name;
        this.age = age;
    }

    public int getId() {
        return id;
    }

    public void setId(int id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getAge() {
        return age;
    }

    public void setAge(String age) {
        this.age = age;
    }

    @Override
    public String toString() {
        return "User{" +
                "id=" + id +
                ", name='" + name + '\'' +
                ", age='" + age + '\'' +
                '}';
    }

}
(8) nil对象[任何对象都比nil大除了nil本身]

nil是Aviator内置的常量，类似java中的null，表示空的值。nil跟null不同的在
于，在java中null只能使用在==、!=的比较运算符，而nil还可以使用>、>=、
<、<=等比较运算符。Aviator规定，[任何对象都比nil大除了nil本身]。用户传入
的变量如果为null，将自动以nil替代。

        AviatorEvaluator.execute("nil == nil");  //true 
        AviatorEvaluator.execute(" 3> nil");    //true 
        AviatorEvaluator.execute(" true!= nil");    //true 
        AviatorEvaluator.execute(" ' '>nil ");  //true 
        AviatorEvaluator.execute(" a==nil ");   //true,a is null
nil与String相加的时候，跟java一样显示为null

(9) 日期比较


/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/25
 */
public class Test {
    public static void main(String[] args) {
        Map<String, Object> env = new HashMap<String, Object>();
        final Date date = new Date();
        String dateStr = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss:SS").format(date);
        env.put("date", date);
        env.put("dateStr", dateStr);

        Boolean result = (Boolean) AviatorEvaluator.execute("date==dateStr",
 env);
        System.out.println(result);

        result = (Boolean) AviatorEvaluator.execute("date > '2009-12-20 
00:00:00:00' ", env);
        System.out.println(result);

        result = (Boolean) AviatorEvaluator.execute("date < '2200-12-20 
00:00:00:00' ", env);
        System.out.println(result);

        result = (Boolean) AviatorEvaluator.execute("date ==date ", env);
        System.out.println(result);


    }
}
(10) 语法手册

数据类型
Number类型：数字类型，支持两种类型，分别对应Java的Long和Double，也就是说任何整数都将被转换为Long，而任何浮点数都将被转换为Double，包括用户传入的数值也是如此转换。不支持科学计数法，仅支持十进制。如-1、100、2.3等。

String类型： 字符串类型，单引号或者双引号括起来的文本串，如'hello world'，变量如果传入的是String或者Character也将转为String类型。

Bool类型： 常量true和false，表示真值和假值，与java的Boolean.TRUE和Boolean.False对应。

Pattern类型： 类似Ruby、perl的正则表达式，以//括起来的字符串，如//d+/，内部实现为java.util.Pattern。

变量类型： 与Java的变量命名规则相同，变量的值由用户传入，如"a"、"b"等

nil类型: 常量nil,类似java中的null，但是nil比较特殊，nil不仅可以参与==、!=的比较，也可以参与>、>=、<、<=的比较，Aviator规定任何类型都n大于nil除了nil本身，nil==nil返回true。用户传入的变量值如果为null，那么也将作为nil处理，nil打印为null。

算术运算符
Aviator支持常见的算术运算符，包括+ - <tt></tt> / % 五个二元运算符，和一元运算符"-"。其中 - <tt></tt> / %和一元的"-"仅能作用于Number类型。

"+"不仅能用于Number类型，还可以用于String的相加，或者字符串与其他对象的相加。Aviator规定，任何类型与String相加，结果为String。

逻辑运算符
Avaitor的支持的逻辑运算符包括，一元否定运算符"!"，以及逻辑与的"&&"，逻辑或的"||"。逻辑运算符的操作数只能为Boolean。

关系运算符
Aviator支持的关系运算符包括"<" "<=" ">" ">=" 以及"=="和"!=" 。
&&和||都执行短路规则。

关系运算符可以作用于Number之间、String之间、Pattern之间、Boolean之间、变量之间以及其他类型与nil之间的关系比较，不同类型除了nil之外不能相互比较。

Aviator规定任何对象都比nil大除了nil之外。

匹配运算符
匹配运算符"=~"用于String和Pattern的匹配，它的左操作数必须为String，右操作数必须为Pattern。匹配成功后，Pattern的分组将存于变量$num，num为分组索引。

三元运算符
Aviator没有提供if else语句，但是提供了三元运算符 "?:"，形式为 bool ? exp1: exp2。 其中bool必须为结果为Boolean类型的表达式，而exp1和exp2可以为任何合法的Aviator表达式，并且不要求exp1和exp2返回的结果类型一致。

两种模式
默认AviatorEvaluator以编译速度优先：
AviatorEvaluator.setOptimize(AviatorEvaluator.COMPILE);
你可以修改为运行速度优先，这会做更多的编译优化：
AviatorEvaluator.setOptimize(AviatorEvaluator.EVAL);
三.MVEL表达式解析器
1.简介 ：

MVEL在很大程度上受到Java语法的启发，作为一个表达式语言，也有一些根本
的区别，旨在更高的效率，例如：直接支持集合、数组和字符串匹配等操作以
及正则表达式。 MVEL用于执行使用Java语法编写的表达式。
2.特性：

MVEL是一个功能强大的基于Java应用程序的表达式语言。
目前最新的版本是2.0，具有以下特性：
(1). 动态JIT优化器。当负载超过一个确保代码产生的阈值时，选择性地产生字
节代码,这大大减少了内存的使用量。新的静态类型检查和属性支持，允许集成
类型安全表达。
(2). 错误报告的改善。包括行和列的错误信息。
(3). 新的脚本语言特征。MVEL2.0 包含函数定义，如：闭包，lambda定义，
标准循环构造(for, while, do-while, do-until…)，空值安全导航操作，内联with
-context运营 ，易变的（isdef）的测试运营等等。
(4). 改进的集成功能。迎合主流的需求，MVEL2.0支持基础类型的个性化属性处理器，集成到JIT中。
(5). 更快的模板引擎，支持线性模板定义，宏定义和个性化标记定义。
(6). 新的交互式shell（MVELSH）。

(7). 缺少可选类型安全
(8). 集成不良，通常通过映射填入内容。没有字节码不能运作用字节码生成编
译时间慢，还增加了可扩展性问题；不用字节码生成运行时执行非常慢
(9). 内存消耗过大
(10). Jar巨大/依赖规模
3.原理：

与java不同，MVEL是动态类型（带有可选分类），也就是说在源文件中是没有
类型限制的。一条MVEL表达式，简单的可以是单个标识符，复杂的则可能是
一个充满了方法调用和内部集合创建的庞大的布尔表达式。
4.使用方式：
maven引入jar：

<dependency>
            <groupId>org.mvel</groupId>
            <artifactId>mvel2</artifactId>
            <version>2.3.1.Final</version>
        </dependency>

测试：

package com.lrq.wechatdemo.utils;

import com.google.common.collect.Maps;
import org.mvel2.MVEL;

import java.util.Map;

/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/26
 */
public class MvelUtils {

    public static void main(String[] args) {
        String expression = "a == null && b == nil ";
        Map<String,Object> map = Maps.newHashMap();
        map.put("a",null);
        map.put("b",null);

        Object object = MVEL.eval(expression,map);
        System.out.println(object);
    }

}
四.EasyRules规则引擎
1.简介：

easy-rules首先集成了mvel表达式，后续可能集成SpEL的一款轻量
级规则引擎
2.特性：

easy rules是一个简单而强大的java规则引擎，它有以下特性：

轻量级框架，学习成本低
基于POJO
为定义业务引擎提供有用的抽象和简便的应用
从原始的规则组合成复杂的规则
它主要包括几个主要的类或接口：Rule,RulesEngine,RuleListener,Facts 
还有几个主要的注解：@Action,@Condition,@Fact,@Priority,@Rule
3.使用方式：

@Rule可以标注name和description属性，每个rule的name要唯一，
如果没有指定，则RuleProxy则默认取类名
@Condition是条件判断，要求返回boolean值，表示是否满足条件

@Action标注条件成立之后触发的方法

@Priority标注该rule的优先级，默认是Integer.MAX_VALUE - 1，值
越小越优先

@Fact 我们要注意Facts的使用。Facts的用法很像Map，它是客户
端和规则文件之间通信的桥梁。在客户端使用put方法向Facts中添
加数据，在规则文件中通过key来得到相应的数据。
有两种使用方式：

java方式
首先先创建规则并标注属性
package com.lrq.wechatdemo.rules;

import org.jeasy.rules.annotation.Action;
import org.jeasy.rules.annotation.Condition;
import org.jeasy.rules.annotation.Fact;
import org.jeasy.rules.annotation.Rule;
import org.jeasy.rules.support.UnitRuleGroup;

/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/26
 */
public class RuleClass {

    @Rule(priority = 1) //规则设定优先级
    public static class FizzRule {
        @Condition
        public boolean isFizz(@Fact("number") Integer number) {
            return number % 5 == 0;
        }

        @Action
        public void printFizz() {
            System.out.print("fizz\n");
        }
    }

    @Rule(priority = 2)
    public static class BuzzRule {
        @Condition
        public boolean isBuzz(@Fact("number") Integer number) {
            return number % 7 == 0;
        }

        @Action
        public void printBuzz() {
            System.out.print("buzz\n");
        }
    }

    public static class FizzBuzzRule extends UnitRuleGroup {

        public FizzBuzzRule(Object... rules) {
            for (Object rule : rules) {
                addRule(rule);
            }
        }

        @Override
        public int getPriority() {
            return 0;
        }
    }

    @Rule(priority = 3)
    public static class NonFizzBuzzRule {

        @Condition
        public boolean isNotFizzNorBuzz(@Fact("number") Integer number) {
            // can return true, because this is the latest rule to trigger according to
            // assigned priorities
            // and in which case, the number is not fizz nor buzz
            return number % 5 != 0 || number % 7 != 0;
        }

        @Action
        public void printInput(@Fact("number") Integer number) {
            System.out.print(number+"\n");
        }
    }

}
然后客户端调用

package com.lrq.wechatdemo.rules;

import org.jeasy.rules.api.Facts;
import org.jeasy.rules.api.Rules;
import org.jeasy.rules.api.RulesEngine;
import org.jeasy.rules.core.DefaultRulesEngine;
import org.jeasy.rules.core.RulesEngineParameters;

/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/26
 */
public class RuleJavaClient {
    public static void main(String[] args) {
        // 创建规则引擎
        RulesEngineParameters parameters = new RulesEngineParameters().skipOnFirstAppliedRule(true);
        RulesEngine fizzBuzzEngine = new DefaultRulesEngine(parameters);

        // 创建规则集并注册规则
        Rules rules = new Rules();
        rules.register(new RuleClass.FizzRule());
        rules.register(new RuleClass.BuzzRule());
        rules.register(new RuleClass.FizzBuzzRule(new RuleClass.FizzRule(), new RuleClass.BuzzRule()));
        rules.register(new RuleClass.NonFizzBuzzRule());

        // 执行规则
        Facts facts = new Facts();
        for (int i = 1; i <= 100; i++) {
            facts.put("number", i);
            fizzBuzzEngine.fire(rules, facts);
            System.out.println();
        }
    }

}

2.yml方式

resources目录下新建fizzbuzz.yml

---
name: "fizz rule"
description: "print fizz if the number is multiple of 5"
priority: 1
condition: "number % 5 == 0"
actions:
- "System.out.println(\"fizz\")"

---
name: "buzz rule"
description: "print buzz if the number is multiple of 7"
priority: 2
condition: "number % 7 == 0"
actions:
- "System.out.println(\"buzz\")"

---
name: "fizzbuzz rule"
description: "print fizzbuzz if the number is multiple of 5 and 7"
priority: 0
condition: "number % 5 == 0 && number % 7 == 0"
actions:
- "System.out.println(\"fizzbuzz\")"

---
name: "non fizzbuzz rule"
description: "print the number itself otherwise"
priority: 3
condition: "number % 5 != 0 || number % 7 != 0"
actions:
- "System.out.println(number)"
客户端调用：


package com.lrq.wechatdemo.rules;

import org.jeasy.rules.api.Facts;
import org.jeasy.rules.api.Rules;
import org.jeasy.rules.api.RulesEngine;
import org.jeasy.rules.core.DefaultRulesEngine;
import org.jeasy.rules.core.RulesEngineParameters;
import org.jeasy.rules.mvel.MVELRuleFactory;

import java.io.FileNotFoundException;
import java.io.FileReader;

/**
 * CreateBy: haleyliu
 * CreateDate: 2018/12/26
 */
public class RuleYmlClient {

    public static void main(String[] args) throws FileNotFoundException {
        // create a rules engine
        RulesEngineParameters parameters = new RulesEngineParameters().skipOnFirstAppliedRule(true);
        RulesEngine fizzBuzzEngine = new DefaultRulesEngine(parameters);

        // create rules
        Rules rules = MVELRuleFactory.createRulesFrom(new FileReader("fizzbuzz.yml"));

        // fire rules
        Facts facts = new Facts();
        for (int i = 1; i <= 100; i++) {
            facts.put("number", i);
            fizzBuzzEngine.fire(rules, facts);
            System.out.println();
        }
    }
}