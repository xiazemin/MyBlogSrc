---
title: runtime-gdb.py
layout: post
category: golang
author: 夏泽民
---
(gdb) source /Users/didi/goLang/src/github.com/golang/go/src/runtime/runtime-gdb.py
Loading Go Runtime support.
Traceback (most recent call last):
  File "/Users/didi/goLang/src/github.com/golang/go/src/runtime/runtime-gdb.py", line 27, in <module>
    goobjfile = gdb.current_objfile() or gdb.objfiles()[0]
IndexError: list index out of range
<!-- more -->
原因：
 必须先 gdb program  然后source才可以，

(gdb) source  /Users/didi/goLang/src/github.com/golang/go/src/runtime/runtime-gdb.py
Loading Go Runtime support.


https://github.com/golang/go/issues/4194
https://github.com/golang/go/issues/6963

https://www.cnblogs.com/kaid/p/9698544.html


问题：

(gdb) source /Users/sherlock/documents/go/src/runtime/runtime-gdb.py
Loading Go Runtime support.
Traceback (most recent call last):
  File "/Users/sherlock/documents/go/src/runtime/runtime-gdb.py", line 205, in <module>
    _rctp_type = gdb.lookup_type("struct runtime.rtype").pointer()
gdb.error: No struct type named runtime.rtype.

解决方案:

1. git clone https://github.com/golang/example.git

2.在stringutil目录下 go test -c ./

3.

gdb -q ./stringutil.test
Reading symbols from ./stringutil.test...done.
(gdb) source /usr/src/go/src/runtime/runtime-gdb.py
Loading Go Runtime support.
Traceback (most recent call last):
  File "/usr/src/go/src/runtime/runtime-gdb.py", line 205, in <module>
    _rctp_type = gdb.lookup_type("struct runtime.rtype").pointer()
gdb.error: No struct type named runtime.rtype.
(gdb) set pagination off
(gdb) set logging on
Copying output to gdb.txt.
(gdb) info types
# much output
(gdb) quit
4.命令行


grep runtime gdb.txt   | grep type;

5.得出结论

struct []*runtime._type;
typedef struct runtime.iface error;
typedef struct runtime.iface flag.Value;
typedef struct runtime.iface flag.boolFlag;
typedef struct runtime.iface fmt.Formatter;
typedef struct runtime.iface fmt.GoStringer;
typedef struct runtime.iface fmt.Scanner;
typedef struct runtime.iface fmt.Stringer;
typedef struct runtime.iface fmt.runeUnreader;
typedef void (struct runtime.g *) func(*runtime.g);
typedef void (struct runtime.stkframe *, void *, bool *) func(*runtime.stkframe, unsafe.Pointer) bool;
typedef void (struct []runtime.StackRecord, int *, bool *) func([]runtime.StackRecord) (int, bool);
typedef struct runtime.iface interface { runtime.f() };
typedef struct runtime.eface interface {};
typedef struct runtime.iface io.Reader;
typedef struct runtime.iface io.ReaderFrom;
typedef struct runtime.iface io.RuneReader;
typedef struct runtime.iface io.Writer;
typedef struct runtime.iface io.WriterTo;
typedef struct hash<string,*runtime/pprof.Profile> * map[string]*runtime/pprof.Profile;
typedef struct runtime.iface os.FileInfo;
typedef struct runtime.iface os.Signal;
typedef struct runtime.iface reflect.Type;
typedef struct runtime.iface regexp.input;
struct runtime._type;
typedef runtime.bucketType;
struct runtime.chantype;
typedef struct runtime.iface runtime.fInterface;
struct runtime.functype;
struct runtime.interfacetype;
typedef runtime.intptr;
struct runtime.maptype;
typedef runtime.pageID;
struct runtime.ptrtype;
struct runtime.slicetype;
typedef struct runtime.iface runtime.stringer;
struct runtime.typeAlg;
typedef runtime.uintreg;
struct runtime.uncommontype;
typedef struct runtime.iface runtime/pprof.countProfile;
typedef struct runtime.iface sort.Interface;
typedef struct runtime.iface sync.Locker;

6.修改runtime-gdb.py第205行rtype改成_type
注：我的平台是mac os x 10.10.3 ,go1.4.2

其他平台可能不是

struct []*runtime._type;