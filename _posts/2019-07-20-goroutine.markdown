---
title: Memory Sanitizer
layout: post
category: golang
author: 夏泽民
---
// Race Detector的处理(用于检测线程冲突问题)
if raceenabled {
    racemalloc(x, size)
}

// Memory Sanitizer的处理(用于检测危险指针等内存问题)
if msanenabled {
    msanmalloc(x, size)
}
https://mail.grokbase.com/t/gg/golang-dev/15a7t92e1g/go-vs-memory-sanitizer
<!-- more -->
The memory sanitizer (msan) is a feature of clang and GCC that inserts
automatic detection of memory errors in C, in particular programs that
make choices based on reads of uninitialized memory. This does not
apply to Go, of course, where all memory is initialized. However,
people can use cgo to pass memory back and forth between C and Go, and
people using cgo may want to use msan to detect errors in their C code
(in fact, it's an awfully good idea).

Since msan tracks the current state of all memory, if C creates some
uninitialized memory, passes that memory to Go, Go fills it in and
returns to C, and C then looks at that memory, msan will not know that
Go initialized the memory and will report an error. Here is a test
case, that fails today with an msan error:
http://play.golang.org/p/HoMF9Iq1o9 (I think you need to be running
at least clang 3.6).

I looked into fixing this, and the problem is surprisingly tractable,
because Go's thread sanitizer already instruments all the necessary
code. This test case is fixed by https://golang.org/cl/15494 .

I'd like to see what people think of really implementing this
approach. This would mean adding a -msan option to the go tool, the
compiler, and the linker, that roughly parallels the existing -race
option. There would be some new code in the runtime, and a new very
small package runtime/msan. Compiling with -msan would mean that Go
code would automatically call the appropriate msan library functions
for every memory read and write. This would only be useful when using
cgo (or SWIG) and compiling with -fsanitize=memory.

This is in some ways a fairly specialized use case, but on the other
hand, when debugging a memory crash in a Go program that calls C code,
it would certainly be nice to be able to use msan to make sure that
the C code is moderately well behaved. This would obviously be more
code to maintain in the toolchain, but it's fairly limited since
almost all of it would be right next to existing race code.
