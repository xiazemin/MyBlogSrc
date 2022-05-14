---
title: GOSSAFUNC 查看 Go 的代码优化过程
layout: post
category: golang
author: 夏泽民
---
之前有人在某群里询问 Go 的编译器是怎么识别下面的代码始终为 false，并进行优化的：

package main

func main() {
    var a = 1
    if a != 1 {
        println("oh no")
    }
}
先说是不是，再说为什么。先看看他的结论对不对：

TEXT main.main(SB) /Users/xargin/test/com.go
  com.go:3              0x104ea70               c3                      RET
  .... 后面都是填充物
整个 main 函数的逻辑都被优化掉了，二进制文件中 main 函数什么都没干就直接 RET 了。说明在编译过程中，Go 的编译器确实会对这段无效代码进行优化。

之前有接触过 Go 的静态扫描工具的同学就问了，Go 编译器的这种优化我们能不能进行复用呢。把逻辑从编译器中抽出来，直接做个静态扫描工具来告诉你又写出了垃圾代码。
<!-- more -->
嗯，我们来看看到底行不行，首先需要简单理解 Go 的编译过程。

Go 从代码文本到可执行执行文件的编译过程大致为：

词法分析 ------------> 语法分析 ----------> 中间代码生成 ----------> 目标代码生成
        token stream            ast                    SSA           asm
当前开源社区的静态扫描工具，分析的对象都是 ast，因为 Go 的 compiler 接口是开放的，所以我们可以直接用 go/parser 、 go/ast 库来生成这个 ast。之后再调用 Walk 来遍历语法树，或者我们自己写一个遍历 ast 的流程也不麻烦。在遍历过程中，可以根据单句代码(比如有个东西叫 ineff assign)，或者根据代码的上下文来给出一些建议和警示(比如一些什么 go vet、gosimple 啊之类的东西)。

从词法分析到语法分析一般被称为编译器的前端(frontend)，而中间代码生成和目标代码生成则是编译器后端(backend)。

所以不管怎么说，想做静态扫描，就是在和 ast 打交道，即在编译器前端折腾。这里的问题是，Go 的编译器对前述代码的优化究竟是在编译过程的哪一步进行的呢？

获得代码的 ast 很简单：

package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
    fset := token.NewFileSet()
    f, _ := parser.ParseFile(fset, "./demo.go", nil, parser.Mode(0))

    for _, d := range f.Decls {
        ast.Print(fset, d)
    }
}
输出 ast：

     0  *ast.FuncDecl {
     1  .  Name: *ast.Ident {
     2  .  .  NamePos: ./com.go:3:6
     3  .  .  Name: "main"
     4  .  .  Obj: *ast.Object {
     5  .  .  .  Kind: func
     6  .  .  .  Name: "main"
     7  .  .  .  Decl: *(obj @ 0)
     8  .  .  }
     9  .  }
    10  .  Type: *ast.FuncType {
    11  .  .  Func: ./com.go:3:1
    12  .  .  Params: *ast.FieldList {
    13  .  .  .  Opening: ./com.go:3:10
    14  .  .  .  Closing: ./com.go:3:11
    15  .  .  }
    16  .  }
    17  .  Body: *ast.BlockStmt {
    18  .  .  Lbrace: ./com.go:3:13
    19  .  .  List: []ast.Stmt (len = 2) {
    20  .  .  .  0: *ast.DeclStmt {
    21  .  .  .  .  Decl: *ast.GenDecl {
    22  .  .  .  .  .  TokPos: ./com.go:4:2
    23  .  .  .  .  .  Tok: var
    24  .  .  .  .  .  Lparen: -
    25  .  .  .  .  .  Specs: []ast.Spec (len = 1) {
    26  .  .  .  .  .  .  0: *ast.ValueSpec {
    27  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    28  .  .  .  .  .  .  .  .  0: *ast.Ident {
    29  .  .  .  .  .  .  .  .  .  NamePos: ./com.go:4:6
    30  .  .  .  .  .  .  .  .  .  Name: "a"
    31  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    32  .  .  .  .  .  .  .  .  .  .  Kind: var
    33  .  .  .  .  .  .  .  .  .  .  Name: "a"
    34  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 26)
    35  .  .  .  .  .  .  .  .  .  .  Data: 0
    36  .  .  .  .  .  .  .  .  .  }
    37  .  .  .  .  .  .  .  .  }
    38  .  .  .  .  .  .  .  }
    39  .  .  .  .  .  .  .  Values: []ast.Expr (len = 1) {
    40  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
    41  .  .  .  .  .  .  .  .  .  ValuePos: ./com.go:4:10
    42  .  .  .  .  .  .  .  .  .  Kind: INT
    43  .  .  .  .  .  .  .  .  .  Value: "1"
    44  .  .  .  .  .  .  .  .  }
    45  .  .  .  .  .  .  .  }
    46  .  .  .  .  .  .  }
    47  .  .  .  .  .  }
    48  .  .  .  .  .  Rparen: -
    49  .  .  .  .  }
    50  .  .  .  }
    51  .  .  .  1: *ast.IfStmt {
    52  .  .  .  .  If: ./com.go:5:2
    53  .  .  .  .  Cond: *ast.BinaryExpr {
    54  .  .  .  .  .  X: *ast.Ident {
    55  .  .  .  .  .  .  NamePos: ./com.go:5:5
    56  .  .  .  .  .  .  Name: "a"
    57  .  .  .  .  .  .  Obj: *(obj @ 31)
    58  .  .  .  .  .  }
    59  .  .  .  .  .  OpPos: ./com.go:5:7
    60  .  .  .  .  .  Op: !=
    61  .  .  .  .  .  Y: *ast.BasicLit {
    62  .  .  .  .  .  .  ValuePos: ./com.go:5:10
    63  .  .  .  .  .  .  Kind: INT
    64  .  .  .  .  .  .  Value: "1"
    65  .  .  .  .  .  }
    66  .  .  .  .  }
    67  .  .  .  .  Body: *ast.BlockStmt {
    68  .  .  .  .  .  Lbrace: ./com.go:5:12
    69  .  .  .  .  .  List: []ast.Stmt (len = 1) {
    70  .  .  .  .  .  .  0: *ast.ExprStmt {
    71  .  .  .  .  .  .  .  X: *ast.CallExpr {
    72  .  .  .  .  .  .  .  .  Fun: *ast.Ident {
    73  .  .  .  .  .  .  .  .  .  NamePos: ./com.go:6:3
    74  .  .  .  .  .  .  .  .  .  Name: "println"
    75  .  .  .  .  .  .  .  .  }
    76  .  .  .  .  .  .  .  .  Lparen: ./com.go:6:10
    77  .  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
    78  .  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
    79  .  .  .  .  .  .  .  .  .  .  ValuePos: ./com.go:6:11
    80  .  .  .  .  .  .  .  .  .  .  Kind: STRING
    81  .  .  .  .  .  .  .  .  .  .  Value: "\"oh no\""
    82  .  .  .  .  .  .  .  .  .  }
    83  .  .  .  .  .  .  .  .  }
    84  .  .  .  .  .  .  .  .  Ellipsis: -
    85  .  .  .  .  .  .  .  .  Rparen: ./com.go:6:18
    86  .  .  .  .  .  .  .  }
    87  .  .  .  .  .  .  }
    88  .  .  .  .  .  }
    89  .  .  .  .  .  Rbrace: ./com.go:7:2
    90  .  .  .  .  }
    91  .  .  .  }
    92  .  .  }
    93  .  .  Rbrace: ./com.go:8:1
    94  .  }
    95  }
显然，到语法分析完毕之后，ast 中的 if 节点还活得好好的。只能看看后端部分了：

GOSSAFUNC=main go build com.go

deadcode

SSA 的多轮优化就是编译原理里常说的后端优化，这一步是 deadcode opt，顾名思义。

dump 过程中可能会有权限问题：

# runtime
<unknown line number>: internal compiler error: 'main': open ssa.html: permission denied

Please file a bug report including a short program that triggers the error.
https://golang.org/issue/new
加个 sudo 就好。

既然 Go 是在编译后端进行的死代码消除，那么对于我们来说，想要复用编译器代码，并提前提示就不太方便了。从原理上来讲，我们仍然可以在遍历 ast 的时候存储一些常量、变量的值来完成前文中提出的需求。这就看你有没有兴趣去实现了
https://xargin.com/go-compiler-opt/

首先，编译器的三个阶段：
逐行扫描源代码，将之转换为一系列的 token，交给 parser 解析。
parser，它将一系列 token 转换为 AST（抽象语法树），用于下一步生成代码。
最后一步，代码生成，会利用上一步生成的 AST 并根据目标机器平台的不同，生成目标机器码。
注意：下面使用的代码包（go/scanner，go/parser，go/token，go/ast）主要是让我们可以方便地对 Go 代码进行解析和生成，做出更有趣的事情。但是 Go 本身的编译器并不是用这些代码包实现的。

扫描代码，进行词法分析
任何编译器的第一步都是将源代码文本分解成 token，由扫描程序（也称为词法分析器）完成。token 可以是关键字，字符串，变量名，函数名等等。每一个有效的词都由 token 表示。
在 Go 中，我们写在代码上的 "package"，"main"，"func" 这些都是 token。

token 由代码中的位置，类型和原始文本组成。我们可以使用 go/scanner 和 go/token 包在 Go 程序中自己执行扫描程序。这意味着我们可以像编译器那样扫描检视自己的代码。
下面，我们将通过一个打印 Hello World 的示例来展示 token。

package main

import (
    "fmt"
    "go/scanner"
    "go/token"
)

func main() {
    src := []byte(`
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
`)

    var s scanner.Scanner
    fset := token.NewFileSet()
    file := fset.AddFile("", fset.Base(), len(src))
    s.Init(file, src, nil, 0)

    for {
        pos, tok, lit := s.Scan()
        fmt.Printf("%-6s%-8s%q\n", fset.Position(pos), tok, lit)

        if tok == token.EOF {
            break
        }
    }
}
首先通过源代码字符串创建 token 集合并初始化 scan.Scanner，它将逐行扫描我们的源代码。
接下来循环调用 Scan() 并打印每个 token 的位置，类型和文本字符串，直到遇到文件结束（EOF）标记。

输出：

2:1   package "package"
2:9   IDENT   "main"
2:13  ;       "\n"
4:1   import  "import"
4:8   STRING  "\"fmt\""
4:13  ;       "\n"
6:1   func    "func"
6:6   IDENT   "main"
6:10  (       ""
6:11  )       ""
6:13  {       ""
7:2   IDENT   "fmt"
7:5   .       ""
7:6   IDENT   "Println"
7:13  (       ""
7:14  STRING  "\"Hello, world!\""
7:29  )       ""
7:30  ;       "\n"
8:1   }       ""
8:2   ;       "\n"
8:3   EOF     ""
以第一行为例分析这个输出，第一列 2:1 表示扫描到了源代码第二行第一个字符，第二列 package 表示 token 是 package，第三列 "package" 表示源代码文本。
我们可以看到在 Scanner 执行过程中将 \n 换行符标记成了 ; 分号，像在 C 语言中是用分号表示一行结束的。这就解释了为什么 Go 不需要分号：它们是在词法分析阶段由 Scanner 智能地解释的。

语法分析
源代码扫描完成后，扫描结果将被传递给语法分析器。语法分析是编译的一个阶段，它将 token 转换为 抽象语法树（AST）。
AST 是源代码的结构化表示。在 AST 中，我们将能够看到程序结构，比如函数和常量声明。

我们使用 go/parser 和 go/ast 来打印完整的 AST：

package main

import (
  "go/ast"
  "go/parser"
  "go/token"
  "log"
)

func main() {
  src := []byte(`
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
`)

  fset := token.NewFileSet()

  file, err := parser.ParseFile(fset, "", src, 0)
  if err != nil {
     log.Fatal(err)
  }

  ast.Print(fset, file)
}
输出：

     0  *ast.File {
     1  .  Package: 2:1
     2  .  Name: *ast.Ident {
     3  .  .  NamePos: 2:9
     4  .  .  Name: "main"
     5  .  }
     6  .  Decls: []ast.Decl (len = 2) {
     7  .  .  0: *ast.GenDecl {
     8  .  .  .  TokPos: 4:1
     9  .  .  .  Tok: import
    10  .  .  .  Lparen: -
    11  .  .  .  Specs: []ast.Spec (len = 1) {
    12  .  .  .  .  0: *ast.ImportSpec {
    13  .  .  .  .  .  Path: *ast.BasicLit {
    14  .  .  .  .  .  .  ValuePos: 4:8
    15  .  .  .  .  .  .  Kind: STRING
    16  .  .  .  .  .  .  Value: "\"fmt\""
    17  .  .  .  .  .  }
    18  .  .  .  .  .  EndPos: -
    19  .  .  .  .  }
    20  .  .  .  }
    21  .  .  .  Rparen: -
    22  .  .  }
    23  .  .  1: *ast.FuncDecl {
    24  .  .  .  Name: *ast.Ident {
    25  .  .  .  .  NamePos: 6:6
    26  .  .  .  .  Name: "main"
    27  .  .  .  .  Obj: *ast.Object {
    28  .  .  .  .  .  Kind: func
    29  .  .  .  .  .  Name: "main"
    30  .  .  .  .  .  Decl: *(obj @ 23)
    31  .  .  .  .  }
    32  .  .  .  }
    33  .  .  .  Type: *ast.FuncType {
    34  .  .  .  .  Func: 6:1
    35  .  .  .  .  Params: *ast.FieldList {
    36  .  .  .  .  .  Opening: 6:10
    37  .  .  .  .  .  Closing: 6:11
    38  .  .  .  .  }
    39  .  .  .  }
    40  .  .  .  Body: *ast.BlockStmt {
    41  .  .  .  .  Lbrace: 6:13
    42  .  .  .  .  List: []ast.Stmt (len = 1) {
    43  .  .  .  .  .  0: *ast.ExprStmt {
    44  .  .  .  .  .  .  X: *ast.CallExpr {
    45  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
    46  .  .  .  .  .  .  .  .  X: *ast.Ident {
    47  .  .  .  .  .  .  .  .  .  NamePos: 7:2
    48  .  .  .  .  .  .  .  .  .  Name: "fmt"
    49  .  .  .  .  .  .  .  .  }
    50  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
    51  .  .  .  .  .  .  .  .  .  NamePos: 7:6
    52  .  .  .  .  .  .  .  .  .  Name: "Println"
    53  .  .  .  .  .  .  .  .  }
    54  .  .  .  .  .  .  .  }
    55  .  .  .  .  .  .  .  Lparen: 7:13
    56  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
    57  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
    58  .  .  .  .  .  .  .  .  .  ValuePos: 7:14
    59  .  .  .  .  .  .  .  .  .  Kind: STRING
    60  .  .  .  .  .  .  .  .  .  Value: "\"Hello, world!\""
    61  .  .  .  .  .  .  .  .  }
    62  .  .  .  .  .  .  .  }
    63  .  .  .  .  .  .  .  Ellipsis: -
    64  .  .  .  .  .  .  .  Rparen: 7:29
    65  .  .  .  .  .  .  }
    66  .  .  .  .  .  }
    67  .  .  .  .  }
    68  .  .  .  .  Rbrace: 8:1
    69  .  .  .  }
    70  .  .  }
    71  .  }
    72  .  Scope: *ast.Scope {
    73  .  .  Objects: map[string]*ast.Object (len = 1) {
    74  .  .  .  "main": *(obj @ 27)
    75  .  .  }
    76  .  }
    77  .  Imports: []*ast.ImportSpec (len = 1) {
    78  .  .  0: *(obj @ 12)
    79  .  }
    80  .  Unresolved: []*ast.Ident (len = 1) {
    81  .  .  0: *(obj @ 46)
    82  .  }
    83  }
分析这个输出，在 Decls 字段中，包含了代码中所有的声明，例如导入、常量、变量和函数。在本例中，我们只有两个：导入fmt包 和 主函数。
为了进一步理解它，我们可以看看下面这个图，它是上述数据的表示，但只包含类型，红色代表与节点对应的代码：

图片描述

main函数由三个部分组成：Name、Type 和 Body。Name 是值为 main 的标识符。由 Type 字段指定的声明将包含参数列表和返回类型（如果我们指定了的话）。正文由一系列语句组成，里面包含了程序的所有行，在本例中只有一行fmt.Println("Hello, world!")。

我们的一条 fmt.Println 语句由 AST 中很多部分组成。
该语句是一个 ExprStmt表达式语句(expression statement)，例如，它可以像这里一样是一个函数调用，它可以是字面量，可以是一个二元运算（例如加法和减法），当然也可以是一元运算（例如自增++，自减--，否定！等）等等。
同时，在函数调用的参数中可以使用任何表达式。


 
然后，ExprStmt 又包含一个 CallExpr，它是我们实际的函数调用。里面又包括几个部分，其中最重要的部分是 Fun 和 Args。
Fun 包含对函数调用的引用，在这种情况下，它是一个 SelectorExpr，因为我们从 fmt 包中选择 Println 标识符。
但是至此，在 AST 中，编译器还不知道 fmt 是一个包，它也可能是 AST 中的一个变量。

Args 包含一个表达式列表，它是函数的参数。这里，我们将一个文本字符串传递给函数，因而它由一个类型为 STRING 的 BasicLit 表示。

显然，AST 包含了许多信息，我们不仅可以分析出以上结论，还可以进一步检查 AST 并查找文件中的所有函数调用。下面，我们将使用 go/ast 包中的 Inspect 函数来递归地遍历树，并分析所有节点的信息。

package main

import (
    "fmt"
    "go/ast"
    "go/parser"
    "go/printer"
    "go/token"
    "os"
)

func main() {
    src := []byte(`
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
`)

    fset := token.NewFileSet()

    file, err := parser.ParseFile(fset, "", src, 0)
    if err != nil {
        fmt.Println(err)
    }

    ast.Inspect(file, func(n ast.Node) bool {
        call, ok := n.(*ast.CallExpr)
        if !ok {
            return true
        }

        printer.Fprint(os.Stdout, fset, call.Fun)
        
        return false
    })
}
输出：

fmt.Println
上面代码的作用是查找所有节点以及它们是否为 *ast.CallExpr 类型，上面也说过这种类型是函数调用。如果是，则使用 go/printer 包打印 Fun 中存在的函数的名称。

构建出 AST 后，将使用 GOPATH 或者在 Go 1.11 及更高版本中的 modules 解析所有导入。然后，执行类型检查，并做一些让程序运行更快的初级优化。

代码生成
在解析导入并做了类型检查之后，我们可以确认程序是合法的 Go 代码，然后就走到将 AST 转换为（伪）目标机器码的过程。

此过程的第一步是将 AST 转换为程序的低级表示，特别是转换为 静态单赋值（SSA）表单。这个中间表示不是最终的机器代码，但它确实代表了最终的机器代码。 SSA 具有一组属性，会使应用优化变得更容易，其中最重要的是在使用变量之前总是定义变量，并且每个变量只分配一次。


 
在生成 SSA 的初始版本之后，将执行一些优化。这些优化适用于某些代码，可以使处理器执行起来更简单且更快速。例如，可以做 死码消除。还有比如可以删除某些 nil 检查，因为编译器可以证明这些检查永远不会出错。

现在通过最简单的例子来说明 SSA 和一些优化过程：

package main

import "fmt"

func main() {
    fmt.Println(2)
}
如你所见，此程序只有一个函数和一个导入。它会在运行时打印 2。但是，此例足以让我们了解SSA。

为了显示生成的 SSA，我们需要将 GOSSAFUNC 环境变量设置为我们想要跟踪的函数，在本例中为main 函数。我们还需要将 -S 标识传递给编译器，这样它就会打印代码并创建一个HTML文件。我们还将编译Linux 64位的文件，以确保机器代码与您在这里看到的相同。
在终端执行下面的命令：

GOSSAFUNC=main GOOS=linux GOARCH=amd64 go build -gcflags -S main.go

会在终端打印出所有的 SSA，同时也会生成一个交互式的 ssa.html 文件，我们用浏览器打开它。

图片描述

当你打开 ssa.html 时，将显示阶段，其中大部分都已折叠。start 阶段是从 AST 生成的SSA；lower 阶段将非机器特定的 SSA 转换为机器特定的 SSA，最后的 genssa 就是生成的机器代码。

start 阶段的代码如下：

b1:
    v1  = InitMem <mem>
    v2  = SP <uintptr>
    v3  = SB <uintptr>
    v4  = ConstInterface <interface {}>
    v5  = ArrayMake1 <[1]interface {}> v4
    v6  = VarDef <mem> {.autotmp_0} v1
    v7  = LocalAddr <*[1]interface {}> {.autotmp_0} v2 v6
    v8  = Store <mem> {[1]interface {}} v7 v5 v6
    v9  = LocalAddr <*[1]interface {}> {.autotmp_0} v2 v8
    v10 = Addr <*uint8> {type.int} v3
    v11 = Addr <*int> {"".statictmp_0} v3
    v12 = IMake <interface {}> v10 v11
    v13 = NilCheck <void> v9 v8
    v14 = Const64 <int> [0]
    v15 = Const64 <int> [1]
    v16 = PtrIndex <*interface {}> v9 v14
    v17 = Store <mem> {interface {}} v16 v12 v8
    v18 = NilCheck <void> v9 v17
    v19 = IsSliceInBounds <bool> v14 v15
    v24 = OffPtr <*[]interface {}> [0] v2
    v28 = OffPtr <*int> [24] v2
If v19 → b2 b3 (likely) (line 6)

b2: ← b1
    v22 = Sub64 <int> v15 v14
    v23 = SliceMake <[]interface {}> v9 v22 v22
    v25 = Copy <mem> v17
    v26 = Store <mem> {[]interface {}} v24 v23 v25
    v27 = StaticCall <mem> {fmt.Println} [48] v26
    v29 = VarKill <mem> {.autotmp_0} v27
Ret v29 (line 7)

b3: ← b1
    v20 = Copy <mem> v17
    v21 = StaticCall <mem> {runtime.panicslice} v20
Exit v21 (line 6)
这个简单的程序就已经产生了相当多的 SSA（总共35行）。然而，很多都是引用，可以消除很多（最终的SSA版本有28行，最终的机器代码版本有18行）。


 
每个 v 都是一个新变量，可以点击来查看它被使用的位置。b 是块，这里有三块：b1，b2，b3。b1 始终会执行，b2 和 b3 是条件块，满足条件才执行。
我们来看 b1 结尾处的 If v19 → b2 b3 (likely)。单击该行中的 v19 可以查看它定义的位置。可以看到它定义为 IsSliceInBounds <bool> v14 v15，通过 Go 编译器源代码，我们知道 IsSliceInBounds 的作用是检查 0 <= arg0 <= arg1。然后单击 v14 和 v15 看看在哪定义的，我们会看到 v14 = Const64 <int> [0]，Const64 是一个常量 64 位整数。 v15 定义一样，放在 args1 的位置。所以，实际执行的是 0 <= 0 <= 1，这显然是正确的。

编译器也能够证明这一点，当我们查看 opt 阶段（“机器无关优化”）时，我们可以看到它已经重写了 v19 为 ConstBool <bool> [true]。结果就是，在 opt deadcode 阶段，b3 条件块被删除了，因为永远也不会执行到 b3。

下面来看一下 Go 编译器在把 SSA 转换为 机器特定的SSA 之后所做的另一个更简单的优化，基于amd64体系结构的机器代码。下面，我们将比较 lower 和 lowered deadcode。
lower：

b1:
    BlockInvalid (6)
b2:
    v2 (?) = SP <uintptr>
    v3 (?) = SB <uintptr>
    v10 (?) = LEAQ <*uint8> {type.int} v3
    v11 (?) = LEAQ <*int> {"".statictmp_0} v3
    v15 (?) = MOVQconst <int> [1]
    v20 (?) = MOVQconst <uintptr> [0]
    v25 (?) = MOVQconst <*uint8> [0]
    v1 (?) = InitMem <mem>
    v6 (6) = VarDef <mem> {.autotmp_0} v1
    v7 (6) = LEAQ <*[1]interface {}> {.autotmp_0} v2
    v9 (6) = LEAQ <*[1]interface {}> {.autotmp_0} v2
    v16 (+6) = LEAQ <*interface {}> {.autotmp_0} v2
    v18 (6) = LEAQ <**uint8> {.autotmp_0} [8] v2
    v21 (6) = LEAQ <**uint8> {.autotmp_0} [8] v2
    v30 (6) = LEAQ <*int> [16] v2
    v19 (6) = LEAQ <*int> [8] v2
    v23 (6) = MOVOconst <int128> [0]
    v8 (6) = MOVOstore <mem> {.autotmp_0} v2 v23 v6
    v22 (6) = MOVQstore <mem> {.autotmp_0} v2 v10 v8
    v17 (6) = MOVQstore <mem> {.autotmp_0} [8] v2 v11 v22
    v14 (6) = MOVQstore <mem> v2 v9 v17
    v28 (6) = MOVQstoreconst <mem> [val=1,off=8] v2 v14
    v26 (6) = MOVQstoreconst <mem> [val=1,off=16] v2 v28
    v27 (6) = CALLstatic <mem> {fmt.Println} [48] v26
    v29 (5) = VarKill <mem> {.autotmp_0} v27
Ret v29 (+7)

 
在HTML中，某些行是灰色的，这意味着它们将在下一个阶段中被删除或修改。
例如，v15 (?) = MOVQconst <int> [1] 显示为灰色。点击 v15，我们看到它在其他地方都没有使用，而 MOVQconst 基本上与我们之前看到的 Const64 相同，只针对amd64的特定机器。我们把 v15 设置为1。但是，v15 在其他地方都没有使用，所以它是无用的（死的）代码并且可以消除。

Go 编译器应用了很多这类优化。因此，虽然 AST 生成的初始 SSA 可能不是最快的实现，但编译器将SSA优化为更快的版本。 HTML 文件中的每个阶段都有可能发生优化。

如果你有兴趣了解 Go 编译器中有关 SSA 的更多信息，请查看 Go 编译器的 SSA 源代码。
这里定义了所有的操作以及优化。

结论
Go 是一种非常高效且高性能的语言，由其编译器及其优化支撑。要了解有关 Go 编译器的更多信息，源代码的 README 是不错的选择。

https://studygolang.com/articles/15088?fr=sidebar

https://www.codercto.com/a/77104.html
https://zhuanlan.zhihu.com/p/107665043
