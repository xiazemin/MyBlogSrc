---
title: ssa CreateProgram 调用关系生成
layout: post
category: golang
author: 夏泽民
---
https://studygolang.com/articles/19607?fr=sidebar
https://blog.csdn.net/dashuniuniu/article/details/78704741

func AllPackages(initial []*packages.Package, mode ssa.BuilderMode) (*ssa.Program, []*ssa.Package) {
AllPackages creates an SSA program for a set of packages plus all their dependencies.

The packages must have been loaded from source syntax using the golang.org/x/tools/go/packages.Load function in LoadAllSyntax mode.

AllPackages creates an SSA package for each well-typed package in the initial list, plus all their dependencies. The resulting list of packages corresponds to the list of intial packages, and may contain a nil if SSA code could not be constructed for the corresponding initial package due to type errors.

Code for bodies of functions is not built until Build is called on the resulting Program. SSA code is constructed for all packages with well-typed syntax trees.

The mode parameter controls diagnostics and checking during SSA construction.

https://gowalker.org/golang.org/x/tools/go/ssa/ssautil
<!-- more -->
https://s0godoc0org.icopy.site/golang.org/x/tools/go/ssa/ssautil

https://github.com/golang/tools/blob/master/go/ssa/ssautil/load.go

https://sourcegraph.com/github.com/golang/tools/-/commit/9c57c19a58835b00c0e3e283952087842242e49b

https://go.googlesource.com/tools/+/release-branch.go1.6/oracle/pointsto14.go


https://golang.hotexamples.com/de/examples/golang.org.x.tools.go.ssa.ssautil/-/CreateProgram/golang-createprogram-function-examples.html

	ssaprog := ssautil.CreateProgram(prog, ssa.GlobalDebug)
	ssaprog.Build()
	
怎么形成一个项目内部的函数调用关系
使用golang提供的静态编译工具链
我们依赖了如下三个golang工具链：

"golang.org/x/tools/go/loader"
"golang.org/x/tools/go/pointer"
"golang.org/x/tools/go/ssa"

go/loader
Package loader loads a complete Go program from source code, parsing and type-checking the initial packages plus their transitive closure of dependencies. The ASTs and the derived facts are retained for later use.

这个包的官方定义如上，大意是指从源代码加载整个项目，解析代码并作类型校验，分析package之间的依赖关系，返回ASTs和衍生的关系。

go/ssa
Package ssa defines a representation of the elements of Go programs (packages, types, functions, variables and constants) using a static single-assignment (SSA) form intermediate representation (IR) for the bodies of functions.

SSA(Static Single Assignment，静态单赋值），是源代码和机器码中间的表现形式。从AST转换到SSA之后，编译器会进行一系列的优化。这些优化被应用于代码的特定阶段使得处理器能够更简单和快速地执行。

go/pointer
Package pointer implements Andersen's analysis, an inclusion-based pointer analysis algorithm first described in (Andersen, 1994).

指针分析是一类特殊的数据流问题，它是其它静态程序分析的基础。算法最终建立各节点间的指向关系，具体可以参考文章Anderson's pointer analysis。


https://studygolang.com/articles/19607?fr=sidebar
https://blog.csdn.net/dashuniuniu/article/details/78704741

https://github.com/baixiaoustc/go_code_analysis/blob/master/second_post_test.go

http://baixiaoustc.com/2019/01/17/2019-01-17-golang-code-inspector-2-reverse-call-graph/

https://github.com/baixiaoustc/go_code_analysis

http://docs.activestate.com/activego/1.8/pkg/golang.org/x/tools/go/pointer/


