---
title: antlr idea
layout: post
category: linux
author: 夏泽民
---
https://www.cntofu.com/book/115/line-between-lexer-and-parser.md
https://blog.csdn.net/qq_36616602/article/details/85858133
 1.安装IDEA.
    2.在File-Settings-Plugins中安装ANTLR v4 grammar plugin插件.
    3.新建一个Maven项目，在pom.xml文件中添加ANTLR4插件和运行库的依赖，注意一定要用最新版的。
    
    项目流程
新建一个g4文件，在里面写入要识别语言的词法规则和语法规则
     .

      2. 写完后，右键prolog.g4，选择Configure ANTLR，配置output路径。
      词法分析器和语法分析器会生成在java文件夹下的com.antlr.out包下，之后如果写的主函数转换程序和生成的这些文件不在同一个文件夹下，就可以通过下面的代码来引用。

import com.antlr.out.prologLexer;
import com.antlr.out.prologParser;
     3.右键prolog.g4，选择Generate ANTLR Recognizer生成所有的监听器Java代码。可以看到生成的结果。

      

     4.最后编写主函数和转换函数。最后的项目目录如下：

      

    我的主函数代码如下，待转化的语言放在t.prolog中，通过文件读取输入，调用转换程序来将带转换语言转换成目标语言。

//主函数
import com.antlr.out.prologLexer;
import com.antlr.out.prologParser;
import org.antlr.v4.runtime.ANTLRInputStream;
import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.tree.ParseTree;
import org.antlr.v4.runtime.tree.ParseTreeWalker;

import java.io.File;
import java.io.FileInputStream;

public class Mytranslation {
    public static void main(String[] args) throws Exception {
        File file = new File("E:\\prolog2NL\\src\\main\\java\\t.prolog");
        FileInputStream is = new FileInputStream(file);
        ANTLRInputStream input = new ANTLRInputStream(is);
        prologLexer lexer = new prologLexer(input);
        CommonTokenStream tokens = new CommonTokenStream(lexer);
        prologParser parser = new prologParser(tokens);
        parser.setBuildParseTree(true);
        ParseTree tree = parser.p_text();
        ParseTreeWalker walker = new ParseTreeWalker();
        prolog2NL.prologEmitter converter = new prolog2NL.prologEmitter();
        walker.walk(converter, tree);
        System.out.println(converter.getprolog(tree));

    }
}
https://blog.csdn.net/qq_36616602/article/details/85858133
https://github.com/mantuoluozk/antlr4_json2xml/blob/master/json2xml/pom.xml
https://github.com/mantuoluozk/antlr4_json2xml
https://github.com/mantuoluozk/antlr4_prolog2NL
<!-- more -->
{% raw %}
Antlr 简介
ANTLR 语言识别的一个工具 (ANother Tool for Language Recognition ) 是一种语言工具，它提供了一个框架，可以通过包含 Java, C++, 或 C# 动作（action）的语法描述来构造语言识别器，编译器和解释器。 计算机语言的解析已经变成了一种非常普遍的工作，在这方面的理论和工具经过近 40 年的发展已经相当成熟，使用 Antlr 等识别工具来识别，解析，构造编译器比手工编程更加容易，同时开发的程序也更易于维护。
语言识别的工具有很多种，比如大名鼎鼎的 Lex 和 YACC，Linux 中有他们的开源版本，分别是 Flex 和 Bison。在 Java 社区里，除了 Antlr 外，语言识别工具还有 JavaCC 和 SableCC 等。
和大多数语言识别工具一样，Antlr 使用上下文无关文法描述语言。最新的 Antlr 是一个基于 LL(*) 的语言识别器。在 Antlr 中通过解析用户自定义的上下文无关文法，自动生成词法分析器 (Lexer)、语法分析器 (Parser) 和树分析器 (Tree Parser)。
Antlr 能做什么
编程语言处理
识别和处理编程语言是 Antlr 的首要任务，编程语言的处理是一项繁重复杂的任务，为了简化处理，一般的编译技术都将语言处理工作分为前端和后端两个部分。其中前端包括词法分析、语法分析、语义分析、中间代码生成等若干步骤，后端包括目标代码生成和代码优化等步骤。

Antlr 致力于解决编译前端的所有工作。使用 Anltr 的语法可以定义目标语言的词法记号和语法规则，Antlr 自动生成目标语言的词法分析器和语法分析器；此外，如果在语法规则中指定抽象语法树的规则，在生成语法分析器的同时，Antlr 还能够生成抽象语法树；最终使用树分析器遍历抽象语法树，完成语义分析和中间代码生成。整个工作在 Anltr 强大的支持下，将变得非常轻松和愉快。
文本处理

文本处理
当需要文本处理时，首先想到的是正则表达式，使用 Anltr 的词法分析器生成器，可以很容易的完成正则表达式能够完成的所有工作；除此之外使用 Anltr 还可以完成一些正则表达式难以完成的工作，比如识别左括号和右括号的成对匹配等。

在IDEA中安装使用Antlr
在Settings-Plugins中安装ANTLR v4 grammar plugin
新建一个Maven项目，在pom.xml文件中添加ANTLR4插件和运行库的依赖。注意一定要用最新版的，依赖，不知道最新版本号的可以自己google一下maven antlr4。
<dependencies>

        <dependency>
            <groupId>org.antlr</groupId>
            <artifactId>antlr4-runtime</artifactId>
            <version>4.5.3</version>
        </dependency>
    </dependencies>
    <build>
        <plugins>
            <plugin>
                <groupId>org.antlr</groupId>
                <artifactId>antlr4-maven-plugin</artifactId>
                <version>4.3</version>
                <executions>
                    <execution>
                        <id>antlr</id>
                        <goals>
                            <goal>antlr4</goal>
                        </goals>
                        <phase>none</phase>
                    </execution>
                </executions>
                <configuration>
                    <outputDirectory>src/test/java</outputDirectory>
                    <listener>true</listener>
                    <treatWarningsAsErrors>true</treatWarningsAsErrors>
                </configuration>
            </plugin>
        </plugins>
    </build>

antlr4-maven-plugin用于生产Java代码，antlr4-runtime则是运行时所需的依赖库。把antlr4-maven-plugin的phase设置成none，这样在Maven 的lifecycle种就不会调用ANTLR4。如果你希望每次构建生成文法可以将这个配置去掉。

我们定义一个最简单的领域语言，从一个简单的完成算术运算的例子出发，详细说明 Antlr 的使用。首先我们需要在src\main\java中新建一个 Antlr 的文法文件， 一般以 .g4 为文件名后缀，命名为 Demo.g4 。
表达式定义
文法定义
在这个文法文件 Demo.g4 中根据 Antlr 的语法规则来定义算术表达式的文法，文件的头部是 grammar 关键字，定义文法的名字，必须与文法文件文件的名字相同：

grammar Demo;
1
为了简单起见，假设我们的自定义语言只能输入一个算术表达式。从而整个程序有一个语句构成，语句有表达式或者换行符构成。如清单 1 所示：

清单1.程序和语句

prog: stat 
; 
stat: expr 
  |NEWLINE 
;

在 Anltr 中，算法的优先级需要通过文法规则的嵌套定义来体现，加减法的优先级低于乘除法，表达式 expr 的定义由乘除法表达式 multExpr 和加减法算符 (‘+’|’-‘) 构成；同理，括号的优先级高于乘除法，乘除法表达式 multExpr 通过原子操作数 atom 和乘除法算符 (‘*’|’/’) 构成。整个表达的定义如清单 2 所示：

清单2.表达式

expr : multExpr (('+'|'-') multExpr)* 
; 
multExpr : atom (('*'|'/') atom)* 
; 
atom:  '(' expr ')' 
      | INT  
   | ID  
;

最后需要考虑的词法的定义，在 Antlr 中语法定义和词法定义通过规则的第一个字符来区别， 规定语法定义符号的第一个字母小写，而词法定义符号的第一个字母大写。算术表达式中用到了 4 类记号 ( 在 Antlr 中被称为 Token)，分别是标识符 ID，表示一个变量；常量 INT，表示一个常数；换行符 NEWLINE 和空格 WS，空格字符在语言处理时将被跳过，skip() 是词法分析器类的一个方法。如清单 3 所示：

清单 3. 记号定义

ID:('a'..'z'|'A'..'Z')+;
INT:'0'..'9'+;
NEWLINE:'\r'?'\n';
WS:(' '|'\t'|'\n'|'\r')+{skip();};
1
2
3
4
Antlr 支持多种目标语言，可以把生成的分析器生成为 Java，C#，C，Python，JavaScript 等多种语言，默认目标语言为 Java，通过 options {language=?;} 来改变目标语言。我们的例子中目标语言为 Java。

整个Demo.g4文件内容如下：

grammar Demo;

//parser
prog:stat
;
stat:expr|NEWLINE
;

expr:multExpr(('+'|'-')multExpr)*
;
multExpr:atom(('*'|'/')atom)*
;
atom:'('expr')'
    |INT
    |ID
;

//lexer
ID:('a'..'z'|'A'..'Z')+;
INT:'0'..'9'+;
NEWLINE:'\r'?'\n';
WS:(' '|'\t'|'\n'|'\r')+{skip();};

运行ANTLR
右键Demo.g4，选择Configure ANTLR，配置output路径。


右键Demo.g4，选择Generate ANTLR Recognizer。可以看到生成结果结果。
其中Demo.tokens为文法中用到的各种符号做了数字化编号，我们可以不关注这个文件。DemoLexer是Antlr生成的词法分析器，DemoParser是Antlr 生成的语法分析器。


调用分析器。新建一个Main.java。
public static void run(String expr) throws Exception{

        //对每一个输入的字符串，构造一个 ANTLRStringStream 流 in
        ANTLRInputStream in = new ANTLRInputStream(expr);

        //用 in 构造词法分析器 lexer，词法分析的作用是产生记号
        DemoLexer lexer = new DemoLexer(in);

        //用词法分析器 lexer 构造一个记号流 tokens
        CommonTokenStream tokens = new CommonTokenStream(lexer);

        //再使用 tokens 构造语法分析器 parser,至此已经完成词法分析和语法分析的准备工作
        DemoParser parser = new DemoParser(tokens);

        //最终调用语法分析器的规则 prog，完成对表达式的验证
        parser.prog();
    }
完整Main.java代码：

import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.ANTLRInputStream;

public class Main {

    public static void run(String expr) throws Exception{

        //对每一个输入的字符串，构造一个 ANTLRStringStream 流 in
        ANTLRInputStream in = new ANTLRInputStream(expr);

        //用 in 构造词法分析器 lexer，词法分析的作用是产生记号
        DemoLexer lexer = new DemoLexer(in);

        //用词法分析器 lexer 构造一个记号流 tokens
        CommonTokenStream tokens = new CommonTokenStream(lexer);

        //再使用 tokens 构造语法分析器 parser,至此已经完成词法分析和语法分析的准备工作
        DemoParser parser = new DemoParser(tokens);

        //最终调用语法分析器的规则 prog，完成对表达式的验证
        parser.prog();
    }

    public static void main(String[] args) throws Exception{

        String[] testStr={
                "2",
                "a+b+3",
                "(a-b)+3",
                "a+(b*3"
        };

        for (String s:testStr){
            System.out.println("Input expr:"+s);
            run(s);
        }
    }
}

运行Main.java
当输入合法的的表达式时，分析器没有任何输出，表示语言被分析器接受；当输入的表达式违反文法规则时，比如“a + (b * 3”，分析器输出 line 0:-1 mismatched input ‘’ expecting ‘)’；提示期待一个右括号却遇到了结束符号。

文法可视化
打开Antlr Preview。
在Demo.g4中选中一个语法定义符号，如expr。右键选中的符合，选择Text Rule expr。

在ANTLR Preview中选择input,输入表达式，如a+b*c+4/2。则能显示出可视化的文法。
https://blog.csdn.net/sherrywong1220/article/details/53697737?utm_source=blogxgwz4
https://www.antlr.org/tools.html
https://www.cnblogs.com/solvit/p/10097453.html
https://www.jianshu.com/p/628f2a4eb815
https://plugins.jetbrains.com/plugin/7358-antlr-v4-grammar-plugin
https://www.cnblogs.com/wynjauu/articles/9873231.html
源码安装
https://blog.csdn.net/weixin_40038847/article/details/78929254
1. 下载brew，这是一个很好在终端运行的软件包安装器，具体可以去官网查看https://brew.sh。打开终端，输入以下命令 /usr/bin/ruby -e “$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)” 然后回车，过后会有提示，按任何键继续。然后终端就会自行下载brew了。等待一会。
2. Brew下载好了，接着安装wget。继续在终端输入以下命令 brew install wget 意思就是用Brew下载wget。
3. wget 下载好后，就可以安装antlr了。接着输入以下命令 wget http://www.antlr3.org/download/C/libantlr3c-3.4.tar.gz 粘贴 复制就可以了。这时候antlr3.4版本的压缩包就安装好了。
4. tar -xzvf ./libantlr3c-3.4.tar.gz 这个命令是解压之前安装好的antlr的jar包。解压之后，在自己的用户目录下 ls 看是否有解压好的libantlr3c-3.4 文件 如图
5. 然后 cd libantlr3c-3.4 进入到这个文件目录里面。这里lib文件太长，可以先按下lib，之后按下Tab键，神奇的事就发生了。如图

可以输入 ls 查看一下libantlr里面的文件。接着输入下面指令
./configure –enable-64bit 这个具体意思我也没有弄懂，大概就是检查是否满足64的配置。
6. 好了，检查好之后，再输入 make 命令，等一会儿。
7. 接着输入 sudo make install 命令这个时候Antlr3.4版本安装好了。可以去你的用户主目录下查看是否有libantlr3c-3.4，使用命令 ls 查看。如下图
8. 这个时候最好可以新建一个第三方的文件目录 mkdir thirdpart，方便管理和使用，这里我们取名为 tirdpart 。
9. 进入新建好的文件目录 cd thirdpart ，再在里面新建两个文件目录 mkdir include （这里面主要是放头文件的）和 mkdir libs （这里面是放一个叫做 libantlr3c.a的静态链接库）。然后在thirdpart目录下输入下面的命令
cp ../libantlr3c-3.4/include/* ./include 这是把头文件拷贝到前面建好的头文件夹中。
接着输入 cp ../libantlr3c-3.4/.libs/libantlr3c.a ./libs 这是拷贝静态链接库到新建的libs中。
解释一下，cp 是复制命令，格式是 cp A B ，三者中间都有一个空格，A代表的是你要拷贝文件的路径 而B表示拷贝到什么文件下的路径。../表示上层目录 ./表示当前目录下 .* 表示隐藏的文件夹 * 可以是文件名，例如.libs 表示这个隐藏文件是.libs 。/*表示目录下所有的文件。
10. 建好上面的文件目录后，在thirdpart 里面用vim编辑器写一个 名字是ExprCppTree.g 的语法文件，注意名字大小写一定要一样，命令是 vim ExprCppTree.g 。回车进入vim模式，按i进入输入模式，然后把下面的代码复制下来。

grammar ExprCppTree;

options {
    language = C;
    output = AST;
    ASTLabelType=pANTLR3_BASE_TREE;
}

@header {
    #include <assert.h>
}

// The suffix '^' means make it a root.
// The suffix '!' means ignore it.

expr: multExpr ((PLUS^ | MINUS^) multExpr)*
    ;

PLUS: '+';
MINUS: '-';

multExpr
    : atom (TIMES^ atom)*
    ;

TIMES: '*';

atom: INT
    | ID
    | '('! expr ')'!
    ;

stmt: expr NEWLINE -> expr  // tree rewrite syntax
    | ID ASSIGN expr NEWLINE -> ^(ASSIGN ID expr) // tree notation
    | NEWLINE ->   // ignore
    ;

ASSIGN: '=';

prog
    : (stmt {pANTLR3_STRING s = $stmt.tree->toStringTree($stmt.tree);
             assert(s->chars);
             printf(" tree \%s\n", s->chars);
            }
        )+
    ;

ID: ('a'..'z'|'A'..'Z')+ ;
INT: '~'? '0'..'9'+ ;
NEWLINE: '\r'? '\n' ;
WS : (' '|'\t')+ {$channel = HIDDEN;};

编辑好后，按esc键退出编辑模式，然后保存退出 命令 :wq 。
11. 这个时候就可以在thirdpart文件里面看到前面写好的ExprCppTree.g文件。接着拷贝下面的命令
java -jar ../antlr-3.4-complete.jar ./ExprCppTree.g 回车
使用java编译前面的那个语法文件。编译之后，在thirdpart下面就可以看到以下的几个文件 ExprCppTree.g ExprCppTreeLexer.c ExprCppTreeLexer.h ExprCppTreeParser.c ExprCppTreeParser.h ExprCppTree.tokens。
12. 接下来编写驱动文件，使用命令 vim main.cpp 建立一个文件，然后复制下面的代码

#include "ExprCppTreeLexer.h"
#include "ExprCppTreeParser.h"
#include <cassert>
#include <map>
#include <string>
#include <iostream>

using std::map;
using std::string;
using std::cout;

class ExprTreeEvaluator {
    map<string,int> memory;
public:
    int run(pANTLR3_BASE_TREE);
};

pANTLR3_BASE_TREE getChild(pANTLR3_BASE_TREE, unsigned);
const char* getText(pANTLR3_BASE_TREE tree);

int main(int argc, char* argv[])
{
  pANTLR3_INPUT_STREAM input;
  pExprCppTreeLexer lex;
  pANTLR3_COMMON_TOKEN_STREAM tokens;
  pExprCppTreeParser parser;

  assert(argc > 1);
  input = antlr3FileStreamNew((pANTLR3_UINT8)argv[1],ANTLR3_ENC_8BIT);
  lex = ExprCppTreeLexerNew(input);

  tokens = antlr3CommonTokenStreamSourceNew(ANTLR3_SIZE_HINT,
                                            TOKENSOURCE(lex));
  parser = ExprCppTreeParserNew(tokens);

  ExprCppTreeParser_prog_return r = parser->prog(parser);

  pANTLR3_BASE_TREE tree = r.tree;

  ExprTreeEvaluator eval;
  int rr = eval.run(tree);
  cout << "Evaluator result: " << rr << '\n';

  parser->free(parser);
  tokens->free(tokens);
  lex->free(lex);
  input->close(input);

  return 0;
}

int ExprTreeEvaluator::run(pANTLR3_BASE_TREE tree)
{
    pANTLR3_COMMON_TOKEN tok = tree->getToken(tree);
    if(tok) {
        switch(tok->type) {
        case INT: {
            const char* s = getText(tree);
            if(s[0] == '~') {
                return -atoi(s+1);
            }
            else {
                return atoi(s);
            }
        }
        case ID: {
            string var(getText(tree));
            return memory[var];
        }
        case PLUS:
            return run(getChild(tree,0)) + run(getChild(tree,1));
        case MINUS:
            return run(getChild(tree,0)) - run(getChild(tree,1));
        case TIMES:
            return run(getChild(tree,0)) * run(getChild(tree,1));
        case ASSIGN: {
            string var(getText(getChild(tree,0)));
            int val = run(getChild(tree,1));
            memory[var] = val;
            return val;
        }
        default:
            cout << "Unhandled token: #" << tok->type << '\n';
            return -1;
        }
    }
    else {
        int k = tree->getChildCount(tree);
        int r = 0;
        for(int i = 0; i < k; i++) {
            r = run(getChild(tree, i));
        }
        return r;
    }
}

pANTLR3_BASE_TREE getChild(pANTLR3_BASE_TREE tree, unsigned i)
{
    assert(i < tree->getChildCount(tree));
    return (pANTLR3_BASE_TREE) tree->getChild(tree, i);
}

const char* getText(pANTLR3_BASE_TREE tree)
{
    return (const char*) tree->getText(tree)->chars;
}

然后保存退出。
13. 接下来我们编译生成可执行的文件test 输入下面的命令
g++ -g -Wall .cpp .c ../lantlr3c/lib/libantlr3c.a -o test -I. -I ../lantlr3c/include/ 这时候就会在thirdpart文件下看到test的文件。
14. 然后我们在thirdpart目录下 新建一个文件测试一下 cat > ./data 接着输入3-4*5-6 然后回车换行，按control+d退出。接着 输入./test ./data 回车，就可以看到下面的内容了。
生成一个tree 表示出这个运算式的运算过程。由内到外，运算符号在最前面，接着是两个参与运算的数。
{% endraw %}
maven 初始化项目
# Maven 使用原型（Archetype）概念为用户提供了大量不同类型的工程模版（614 个）。  
# 创建简单 java 项目  
mvn archetype:generate  
  
mvn archetype:generate \  
-DgroupId=org.darebeat \  
-DartifactId=HelloWorld \  
-DarchetypeArtifactId=maven-archetype-quickstart \  
-DarchetypeCatalog=local \  
-DinteractiveMode=false  
  
# 创建 Web 应用  
mvn archetype:generate \  
-DgroupId=org.darebeat \  
-DartifactId=HelloWorld \  
-DarchetypeArtifactId=maven-archetype-webapp \  
-DarchetypeCatalog=local \  
-DinteractiveMode=false  


打包：
# 1.清理目标目录（clean）  
# 2.打包工程构建的输出为 jar（package）文件  
# 3.测试报告存放在 maven\target\surefire-reports 文件夹中  
# 4.打包好的 jar 文件在 consumerBanking\target 中  
cd maven  
mvn clean package 

下载依赖
 mvn -f pom.xml dependency:copy-dependencies  。在本地仓库就可以看到依赖包下载下来了。
 [INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS

ANTLR Tool version 4.8 used for code generation does not match the current runtime version 4.7.2ANTLR Tool version 4.8 used for code generation does not match the current runtime version 4.7.2324


https://www.jianshu.com/p/628f2a4eb815
https://www.cnblogs.com/niutao/p/11634973.html

Could not resolve dependencies for project antlr:calc:jar:1.0-SNAPSHOT: Could not find artifact org.antlr:antlr4:jar:4.8.1 in central (http://com:80/artifactory/libs-release) -> [Help 1]

https://www.antlr.org/download/

$curl -O  https://www.antlr.org/download/antlr-4.8-complete.jar

https://github.com/antlr/antlr4/blob/master/doc/go-target.md
https://www.antlr.org/download.html

