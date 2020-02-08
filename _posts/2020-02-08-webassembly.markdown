---
title: webassembly syscall js
layout: post
category: golang
author: 夏泽民
---
Golang 标准库中的 syscall/js 包提供了一系列接口。其中 js.Global() 返回一个 js.Value 类型的结构体，它指代 JS 中的全局对象，在浏览器环境中即为 window 对象。可以通过其 Get() 方法获取 window 对象中的字段，也是 js.Value 类型，包括其中的函数对象，并使用其 Invoke() 方法调用 JS 函数。

另一方面，可以使用 js.Value 类型的 Set() 方法向 JS 中注入字段，包括用 js.NewCallback() 封装的 Golang 函数，这样就能在 JS 中调用 Golang 的函数。Golang 函数必须是 func(args []js.Value) 形式的，使用 args 参数接收 JS 调用的参数，且没有返回值。
<!-- more -->
package main

import (
    "sync"
    "syscall/js"
)

type JsFuncTable struct {
    JsFunc func(int, string) (int, string)
}

var jsFuncTable *JsFuncTable

func goFunc(i int, s string) (int, string) {
    i, s = jsFuncTable.JsFunc(i, s)
    return i + 2, s + "c"
}

func main() {
    jsFuncs := js.Global().Get("jsFuncs")
    jsFuncTable = &JsFuncTable{
        JsFunc: func(i int, s string) (int, string) {
            res := jsFuncs.Get("jsFunc").Invoke(i, s)
            return res.Get("i").Int(), res.Get("s").String()
        },
    }

    goFuncs := js.Global().Get("goFuncs")
    goFuncs.Set("goFunc", js.NewCallback(func(args []js.Value) {
        i, s := goFunc(args[0].Int(), args[1].String())
        ret := args[2]
        ret.Set("i", i)
        ret.Set("s", s)
    }))

    wg := &sync.WaitGroup{}
    wg.Add(1)
    wg.Wait()
}
<html>
  <head>
    <meta charset="utf-8">
    <script src="wasm_exec.js"></script>
    <script>
      window.onload = () => {
        document.getElementById("btn").addEventListener("click", event => {
          var ret = {}
          window.goFuncs.goFunc(0, "a", ret)
          console.dir(ret)
        })
      }

      window.goFuncs = {}
      window.jsFuncs = {
        jsFunc: (i, s) => {
          return {i: i + 1, s: s + "b"}
        },
      }

      const go = new Go()
      WebAssembly.instantiateStreaming(fetch("go_main.wasm"), go.importObject).
        then(res => {
          go.run(res.instance)
        })
    </script>
  </head>
  <body>
    <input id="btn" type="button" value="go" />
  </body>
</html>
编译、部署、运行方式与上一节相同。从上述例子也能看出 Golang 和 JS 两端的代码运行的生命周期，JS 中的 go.run() 异步执行对应 Golang 模块中的 main()，main() 作为 Golang 端的 main loop 在整个页面的生命周期中不能返回，因为后续在 JS 中对该模块中 Golang 函数的调用，会在 main loop 的子协程中执行。

内存访问
除了函数调用的交互，还可以通过内存直接共享数据。

Golang 端使用的内存空间，通过 instance.exports.mem 暴露给 JS 端，这里 instance 为 WebAssembly.instantiate* 函数实例化 wasm 模块得到的 instance。可以通过 mem 创建 TypedArray，以此在 JS 直接读写 Golang 使用的内存。

下面的例子会在 JS 端打开一个图片文件，显示在页面上，并将文件内容直接写入 Golang 使用的内存，在 Golang 中将图片的色调改变，再回调 JS 端来读取改变之后的图片，并显示在页面上。

package main

import (
    "bytes"
    "image"
    "reflect"
    "sync"
    "syscall/js"
    "unsafe"

    "github.com/anthonynsimon/bild/adjust"
    "github.com/anthonynsimon/bild/imgio"
)

type Ctx struct {
    SetFileArrCb    js.Value
    SetImageToHueCb js.Value
}

func setFile(ctx *Ctx, fileJsArr js.Value, length int) {
    bs := make([]byte, length)
    ptr := (*reflect.SliceHeader)(unsafe.Pointer(&bs)).Data
    ctx.SetFileArrCb.Invoke(fileJsArr, ptr)

    img, _, _ := image.Decode(bytes.NewReader(bs))
    buf := &bytes.Buffer{}
    imgio.JPEGEncoder(93)(buf, adjust.Hue(img, -150))

    bs = buf.Bytes()
    ptr = (*reflect.SliceHeader)(unsafe.Pointer(&bs)).Data
    ctx.SetImageToHueCb.Invoke(ptr, len(bs))
}

func main() {
    jsGlobal := js.Global()
    ctx := &Ctx{
        SetFileArrCb:    jsGlobal.Get("setFileArrCb"),
        SetImageToHueCb: jsGlobal.Get("setImageToHueCb"),
    }

    goFuncs := jsGlobal.Get("goFuncs")
    goFuncs.Set("setFile", js.NewCallback(func(args []js.Value) {
        setFile(ctx, args[0], args[1].Int())
    }))

    wg := &sync.WaitGroup{}
    wg.Add(1)
    wg.Wait()
}

<html>
  <head>
    <meta charset="utf-8">
    <script src="wasm_exec.js"></script>
    <script>
      let goMemArr, fileType

      let setImageToElem = (elemId, dateArr) => {
        document.getElementById(elemId).src = URL.createObjectURL(
          new Blob([dateArr], {"type": fileType}))
      }

      window.setFileArrCb = (fileArr, ptr) => {
        goMemArr.set(fileArr, ptr)
      }
      window.setImageToHueCb = (ptr, len) => {
        setImageToElem("img-hue", goMemArr.slice(ptr, ptr + len))
      }
      window.goFuncs = {}

      const go = new Go()
      WebAssembly.instantiateStreaming(fetch("go_main.wasm"), go.importObject).
        then(res => {
          goMemArr = new Uint8Array(res.instance.exports.mem.buffer)
          go.run(res.instance)
        }
      )

      let onFileSelected = event => {
        let reader = new FileReader()
        let file = event.target.files[0]
        fileType = file.type
        reader.onload = event => {
          let fileArr = new Uint8Array(event.target.result)
          setImageToElem("img-ori", fileArr)
          window.goFuncs.setFile(fileArr, fileArr.length)
        }
        reader.readAsArrayBuffer(file)
      }

      window.onload = () => {
        document.getElementById("file-input").addEventListener("change", onFileSelected)
      }
    </script>
  </head>

  <body>
    <input id="file-input" type="file" />
    <br />
    <image id="img-ori" />
    <br />
    <image id="img-hue" />
  </body>
</html>

Golang 在1.11版本中引入了 WebAssembly 支持,意味着以后可以用 go编写可以在浏览器中运行的程序,当然这个肯定也是要受浏览器沙盒环境约束的.

1. 浏览器中运行 Go
1.1 code
package main
func main() {
    println("Hello, WebAssembly!")
}
1.2 编译
必须是 go1.11才行

GOARCH=wasm GOOS=js go build -o test.wasm main.go
1.3 运行
单独的 wasm 文件是无法直接运行的,必须载入浏览器中.

mkdir test
cp test.wasm test
cp $GOROOT/misc/wasm/wasm_exec.{html,js} .
1.3.1 一个测试 http 服务器
chrome 是不支持本地文件中运行 wasm 的,所以必须有一个 http 服务器

//http.go
package main

import (
    "flag"
    "log"
    "net/http"
    "strings"
)

var (
    listen = flag.String("listen", ":8080", "listen address")
    dir    = flag.String("dir", ".", "directory to serve")
)

func main() {
    flag.Parse()
    log.Printf("listening on %q...", *listen)
    log.Fatal(http.ListenAndServe(*listen, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
        if strings.HasSuffix(req.URL.Path, ".wasm") {
            resp.Header().Set("content-type", "application/wasm")
        }

        http.FileServer(http.Dir(*dir)).ServeHTTP(resp, req)
    })))
}
1.3.2 http.go
mv http.go test
cd test
go run http.go 
1.4 效果
在浏览器中打开http://localhost:8080/wasm_exec.html,点击 run 按钮,可以在控制台看到 Hello, WebAssembly!字符串


node中运行 wasm
这个更直接

node wasm_exec.js test.wasm
就可以在控制台看到Hello, WebAssembly!字符串了.


https://github.com/stdiopt/gowasm-experiments 中有许多例子


https://github.com/golang/go/wiki/WebAssembly
