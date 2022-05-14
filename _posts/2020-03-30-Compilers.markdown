---
title: go Compilers
layout: post
category: golang
author: 夏泽民
---
https://www.caffeinatedwonders.com/2019/12/26/tour-of-go-compilers/
On a high level, compilers are viewed as a single, solid, single-step-worker entity. A Gopher, aka Go programmer, interacts with the Go compiler via the abstractions of go build or go install commands. However, compilers, in their ideal form, consist of three layers creatively named: frontend, middle-end, and backend. The single go {build,install} command embarks on a journey through the three layers to convert the raw Go source code into machine code.
<!-- more -->
The frontend ingests the source code, perform lexical, syntactical, and semantical analyses to verify all the language level constraints are met and sensical. The frontend then generates what is known as intermediate representation (IR) to handout to the middle-end. The middle-end does not know the syntax of the language nor does it know the machine language. The middle-end is responsible for optimization via symbolic manipulation. The backend converts the IR into machine code suitable for the intended architecture.

In mathematical analogy, translating word problem into a collection of placeholders and a system of equations is akin to the work of the compiler frontend. The middle-end takes the set of matrices representing the system of equations as created by the frontend, and perform linear algebraic manipulation to simplify the system reducing into, for example, a row echelon form. The backend then figures the values of the placeholders to translate them into meaningful representation for the word problem.

Compilers of Go
The Go compiler obtained by downloading the archives available in the Downloads page is called gc, which stands for “Go compiler”, not “Garbage Collection” which is denoted by “GC”. The code for gc resides $GOROOT/src/cmd/compile.

Before discussing the various compilers of Go, we must know on what basis do they compile certain source code. In other words, on what basis should any Go compiler accept certain constructs of submitted, presumably, Go code. The answer to this comes down to The Go Programming Language Specification.

One of the strong points of Go is that it is based on a specification rather than on an implementation. This allows the implementation to be found right or wrong by testing it against edge-cases and improve either the spec or the implementation. There have been many cases reported in Go’s issue tracker where gccgo, gc compiler, and the go/types package disagree on code validity. The fixes vary between fixing go/types, gc, gccgo, or even amending the spec to be more clear and specific about intentions.

Here are some issues caught by the benefit of having multiple compilers:

gccgo: internal compiler error declaring method without body: This is a case where gccgo crashes upon facing certain code but go/types reports an error. The fix was implemented in the gofrontend repo, which is the base for gccgo & gollvm.

spec: should accept method recv base type that is alias to a pointer type: This includes a snippet that is accepted by gc, but rejected by gccgo and go/types. The fix not only was implemented in gofrontend and go/types, but also amended the spec to be more clear and specific.

go/types: failure to reject interface{}(nil) == []int(nil): This resulted in a fix in go/types.

cmd/compile: inconsistent behaviors in judging whether or not two types are identical and cmd/compile: Embedded field names are not tracked correctly with aliased types: Those seem to be bugs in gc.

cmd/compile: incorrect package initialization order for spec example and consequently spec: clarify section on initialization order: In these, inconsistency in behavior between gccgo and gc brought up the matter of solidifying the wording on initialization order in the spec. The fix didn’t only go after gc but also amended the spec to be more specific.

go/types: cannot assign “nil…” to … parameters: The fix was implemented in go/types because the spec is clear about such scenario.

spec: clarify type elision rules for composite literals: Fix implemented in go/types, partly due to gc and gccgo agreeing on interpretation of spec.

gccgo: reject typed non-int len/cap arguments to make: Here go/types has the correct implementation, and the fix was implemented in gccgo and gc compilers.

cmd/compile, go/types: erroneously accepts method expressions with anonymous receiver type and spec: remove unnecessary syntax restrictions for method expressions to match compiler: This is a case where gc and go/types accepted a form that was ambiguous in the spec which led to gccgo to reject it. The fix was to amend the spec to allow code of that form.

Not only is the availability of different compilers good for sanity-check, it is also good to target more platforms than whatever limited list supported had it been a single compiler. The spec defines behavior for all compilers to adhere to, and allows the internals to not matter anymore.

Here we will explore the following compilers:

gc

Gccgo

Gollvm

Gopherjs

TinyGo

TARDIS Go

gc
This is the default Go compiler, as mentioned earlier, and one of two officially supported compilers by the Go team. The compiler was written in C, starting its life as the Plan9 C compiler turned into Go compiler. It was later rewritten in Go in Go 1.4 via mechanical translation from C to Go (Proposal, GopherCon 2014 Slides, GopherCon 2014 Talk). In 2015, Keith Randall published a proposal to convert the Go compiler IR from syntax-tree based to Static Single Assignment-based IR, which allows for better generation of machine code. The new SSA-based IR landed in Go 1.7 then propagated to the other architectures. Keith Randall made a great talk in GopherCon 2017 discussing the effort to overhaul the compiler.

Further Reading:

Iskander Sharipov in a neat blogpost goes through the SSA optimization rules used in the Go compiler.

The $GOROOT/src/cmd/compile directory tree contains the Introduction to the Go compiler and Introduction to the Go compiler’s SSA backend READMEs for those interested in the inner lives of the Go SSA-based compiler.

Go: Overview of the Compiler by Vincent Blanchon.

gccgo
Gccgo is the other officially supported compiler, and covered by the Go1 Compatibilty Promise. You may hear it is being referred to as “the other Go compiler.” The gcc frontend was started by Ian Lance Taylor as he announced it in 2008 mail message to the Go development team saying:

One of my office-mates pointed me at http://…/go_lang.html . It seems like an interesting language, and I threw together a gcc frontend for it.

~ Go: Ten years and climbing, by Rob Pike

Gccgo support was merged into mainline GCC in GCC version 4.7.1. Due to various reasons, mainly different development and release cycles of GCC and packages of OS distros, gccgo, as distributed by package managers, tends to lag behind gc in terms of support of the latest Go stdlib changes. Thus in order to be able to use the latest Go stdlib, it is necessary to compile gofrontend and gcc from source. Here I try to put together in a single snippet all the commands needed to build gccgo from source:

 1
 2
 3
 4
 5
 6
 7
 8
 9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
# Install prerequisites
apt install -y gcc git subversion make g++ flex libgmp-dev libmpfr-dev libmpc-dev

# grab the projects
git clone https://go.googlesource.com/gofrontend
svn checkout svn://gcc.gnu.org/svn/gcc/trunk

export GOFRONTEND=$(pwd)/gofrontend
cd trunk
rm -rf gcc/go/gofrontend
ln -s $GOFRONTEND/go gcc/go/gofrontend
rm -rf libgo
mkdir libgo
for f in $GOFRONTEND/libgo/*; do ln -s $f libgo/`basename $f`; done
./contrib/download_prerequisites
mkdir objdir
cd objdir

# depending on your distro, you might need to re-link /bin/sh to /bin/bash instead of /bin/dash.

../configure \
    --enable-languages=c,c++,go \
    --disable-libquadmath \
    --disable-libquadmath-support \
    --disable-werror \
    --disable-multilib  \
    --enable-gold
make
make install

# or /usr/local/lib on 32-bit systems
export LD_LIBRARY_PATH=/usr/local/lib64
To build projects using gccgo, you then invoke go build -compiler gccgo. However, there is a catch: gccgo builds are dynamically linked by default. You can confirm this by running ldd against the resulting executable. For example, running ldd against the binary resulting from building Caddy (v1) from source with gccgo shows:

1
2
3
4
5
6
7
linux-vdso.so.1 (0x00007fff98b21000)
libgo.so.15 => /usr/local/lib64/libgo.so.15 (0x00007f00b12d7000)
libm.so.6 => /lib/x86_64-linux-gnu/libm.so.6 (0x00007f00b0f39000)
libgcc_s.so.1 => /usr/local/lib64/libgcc_s.so.1 (0x00007f00b0d21000)
libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f00b0930000)
/lib64/ld-linux-x86-64.so.2 (0x00007f00b2d9a000)
libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f00b0711000)
We can see in the output how libgo, which is the Go stdlib as defined by the gofrontend, is dynamically linked, among other libs. The first step in our adventure to produce statically linked executable is to link libgo statically by using -gccgoflags '-static-libgo'. Again, building Caddy with go build -compiler gccgo -gccgoflags '-static-libgo' then running ldd shows the list of dynamically linked libraries as before sans the libgo line libgo.so.15 => /usr/local/lib64/libgo.so.15 (0x00007f00b12d7000).

1
2
3
4
5
6
linux-vdso.so.1 (0x00007ffc06dfa000)
libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007f7bfb121000)
libm.so.6 => /lib/x86_64-linux-gnu/libm.so.6 (0x00007f7bfad83000)
libgcc_s.so.1 => /usr/local/lib64/libgcc_s.so.1 (0x00007f7bfab6b000)
libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f7bfa77a000)
/lib64/ld-linux-x86-64.so.2 (0x00007f7bfb340000)
The last step is to link everything statically by using -gccgoflags '-static', thus a full command would be go build -compiler gccgo -gccgoflags '-static'. This brings another catch: some packages have dependency on cgo, and thus fail to compile statically. For example, the golang.org/x/sys package limits gccgo builds to an implementation that relies on C with the build tags here and here. Thus when trying to compile Caddy with CGO_ENABLED=0 go build -compiler gccgo -tags netgo -gccgoflags '-static' you’ll get the error undefined reference to 'gccgoRealSyscall'.

There’s currently an endeavour by Nikhil Benesch to improve gccgo. Nikhil initiated the effort by starting the conversation in gofrontend-dev mailinglist about the biggest problems facing gccgo and gollvm. That was followed by a blogpost about the state of gccgo in 2019. Nikhil also created a gccgo slack channel in the Gopher slack for those who are interested.

Gollvm
The gollvm project uses the gofrontend to generate LLVM IR so LLVM can be used as the backend to generate the machine code. The benefits gained from the gollvm project are very similar to the benfits gained from gccgo, including but not limited to wider support for operating systems and architecture. Building and installing gollvm is simple and can be achieved by following the steps listed out on the project repo. Just as I did with gccgo, I’ve put together all the commands into a single snippet:

 1
 2
 3
 4
 5
 6
 7
 8
 9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
apt install git gcc ninja-build

mkdir workarea
cd workarea

# Sources
git clone https://github.com/llvm/llvm-project.git

cd llvm-project/llvm/tools
git clone https://go.googlesource.com/gollvm

cd gollvm
git clone https://go.googlesource.com/gofrontend

cd libgo
git clone https://github.com/libffi/libffi.git
git clone https://github.com/ianlancetaylor/libbacktrace.git

# back to workdir
cd ../..

mkdir build.rel
cd build.rel

CC=gcc CXX=g++ cmake \
 -DCMAKE_INSTALL_PREFIX=/usr/local \
 -DCMAKE_BUILD_TYPE=Release \
 -DLLVM_USE_LINKER=gold \
 -G Ninja \
 ../llvm-project/llvm

# Build all of gollvm
ninja gollvm

# Install gollvm to /usr/local
ninja install-gollvm

export LD_LIBRARY_PATH=/usr/local/lib64
export PATH=/usr/local/bin:$PATH
This installs a modified go executable that wraps llvm-goc, the compiler driver. It is recommended to use the familiar go command and rest of go tools family, which, in turn, will invoke the llvm-goc to produce the desired output.

Unfortunately, gollvm, just like gccgo, suffers from a few limitations causing it to not be able to build some projects.

Limitations Common to Gollvm & Gccgo
As it stands, gccgo & gollvm have a few limitations and/or concerns:

-X ldflag doesn’t work

x/ packages are not formally tested

Some projects have tightly coupled themselves into the Go runtime and implementation as provided by gc, which makes them not buildable by gccgo and gollvm. For example, starlark-go links into the Go runtime implementation of runtime.stringHash. Etcd, and by extension Kubernetes and CoreDNS, cannot be built by gccgo and gollvm due to one of their dependencies using the reflect package, whose implementation is compiler dependent. This issue is discussed in this golang-nuts thread and tracked in this bug report.

Gopherjs
Before the Go compiler (gc) learning how to target wasm in Go 1.11, the only way to write Go code targeting JavaScript runtimes was via gopherjs. The commandable effort results in a list where the majority of stdlib package are supported. Not only are there bindings for some of the most common JavaScript libraries and frameworks, there’s Vecty for anyone who would like to write their app in Go end-to-end. If you would like to have a peek at something familiar and digestable, you can take a look at the TodoMVC implementation with Gopherjs.

TinyGo
Go programmers who are also into programming microcontrolleres and are tired of C/C++ and assembly have the TinyGo project to the rescue. TinyGo is a Go compiler whose tagline is “Go compiler for small places.” The compiler relies on LLVM to do the heavy lifting, and supports a decent amount of the language features and the standard library. Now you can use Go in your next Arduino Uno project.

TARDIS Go
This is perhaps the most fun compiler. The TARDIS Go compiler generates the SSA form of the Go code, then uses the SSA to generate Haxe code conforming to the behavior of the Go code. Wait there’s more. The Haxe toolchain can then be invoked to translate the Haxe code to C++, C#, Java, or JavaScript, thus allowing Go packages to be used on the respective platforms (e.g. JVM).

The demo page of the project is still running the same JavaScript code generated from Haxe that’s in turn generated from Go via TARDIS Go. Sadly, the project hasn’t been updated since 2016, thus its implementation of the Go runtime is at Go 1.4. It would be nice to see the project resurrected and extended to add Go as Haxe target, thus enabling interoperability across all the supported languages and enriching the entire ecosystem.

In Sum
This list is definitely not exhaustive nor should it. Different compilers serve different goals. New implementations put both the spec and the implementations to the test, more tests means more eyeballs, and “given enough eyeballs, all bugs are shallow.” Even compilers that aim at unconventional targets, namely TARDIS Go, managed to find three bugs in Go toolchain.

Ian Lance Taylor himself has an answer to the question in his Quora response to the prompt Why are there multiple compilers and interpreters for each programming language?

Having multiple compilers/interpreters for a language helps ensure that the language is clearly defined. When the implementations differ, either there is a bug, or there is a failure of definition. Having a language be clearly defined is better for its users, as it means that they can reliably predict how a piece of code is going to behave, rather than have to run it to see what happens. Normally different implementations have a different focus: execution speed, compilation speed, debuggability, analysis, etc. Pick the one that most closely matches your needs.

If you’re interested in learning how compilers work, Computerphile has a brilliant playlist about compilers lead by professor Brailsford.
