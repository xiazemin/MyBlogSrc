---
title: WebAssembly
layout: post
category: golang
author: 夏泽民
---
https://github.com/WebAssembly/design
https://juejin.im/post/5e9ee0e7e51d4546f36a5c67
Go 语言源代码的 cmd/compile/internal 中包含了非常多机器码生成相关的包，不同类型的 CPU 分别使用了不同的包进行生成 amd64、arm、arm64、mips、mips64、ppc64、s390x、x86 和 wasm，也就是说 Go 语言能够在上述的 CPU 指令集类型上运行，其中比较有趣的就是 WebAssembly 了。

作为一种在栈虚拟机上使用的二进制指令格式，它的设计的主要目标就是在 Web 浏览器上提供一种具有高可移植性的目标语言。Go 语言的编译器既然能够生成 WASM 格式的指令，那么就能够运行在常见的主流浏览器中。

$ GOARCH=wasm GOOS=js go build -o lib.wasm main.go

<!-- more -->
高效
WebAssembly 有一套完整的语义，实际上 wasm 是体积小且加载快的二进制格式， 其目标就是充分发挥硬件能力以达到原生执行效率

安全
WebAssembly 运行在一个沙箱化的执行环境中，甚至可以在现有的 JavaScript 虚拟机中实现。在web环境中，WebAssembly将会严格遵守同源策略以及浏览器安全策略。

开放
WebAssembly 设计了一个非常规整的文本格式用来、调试、测试、实验、优化、学习、教学或者编写程序。可以以这种文本格式在web页面上查看wasm模块的源码。

标准
WebAssembly 在 web 中被设计成无版本、特性可测试、向后兼容的。WebAssembly 可以被 JavaScript 调用，进入 JavaScript 上下文，也可以像 Web API 一样调用浏览器的功能。当然，WebAssembly 不仅可以运行在浏览器上，也可以运行在非web环境下。

http://webassembly.org.cn/

使用场景
WebAssembly 的整体目标 定义了 WebAssembly 适合做什么。哪些是在 Web 平台可以实现的，哪些是非 Web 平台可以实现的。下面给出了一个不完善的无序列表，包括应用/领域/计算等方向，它们可能将从 WebAssembly 中受益的， WebAssamlby 的设计过程中也会将它们做为用例。

在浏览器中
更好的让一些语言和工具可以编译到 Web 平台运行。
图片/视频编辑。
游戏：
需要快速打开的小游戏
AAA 级，资源量很大的游戏。
游戏门户（代理/原创游戏平台）
P2P 应用（游戏，实时合作编辑）
音乐播放器（流媒体，缓存）
图像识别
视频直播
VR 和虚拟现实
CAD 软件
科学可视化和仿真
互动教育软件和新闻文章。
模拟/仿真平台(ARC, DOSBox, QEMU, MAME, …)。
语言编译器/虚拟机。
POSIX用户空间环境，允许移植现有的POSIX应用程序。
开发者工具（编辑器，编译器，调试器…）
远程桌面。
VPN。
加密工具。
本地 Web 服务器。
使用 NPAPI 分发的插件，但会受限于 Web 安全协议，可以使用 Web APIs。
企业软件功能性客户端（比如：数据库）
脱离浏览器
游戏分发服务（便携、安全）。
服务端执行不可信任的代码。
服务端应用。
移动混合原生应用。
多节点对称计算
如何使用 WebAssembly
整个代码库都用 WebAssembly。
主要使用 WebAssembly 计算，UI 使用 JavaScript/HTML。
在大型 JavaScript/HTML 应用中复用已经存在的 WebAssembly 代码。像使用助手库一样，分担一些计算任务。


目前还存在以下问题：

浏览器兼容性不好，只有最新版本的浏览器支持，并且不同的浏览器对 JS WebAssembly 互调的 API 支持不一致；
生态工具不完善不成熟，目前还不能找到一门体验流畅的编写 WebAssembly 的语言，都还处于起步阶段；
学习资料太少，还需要更多的人去探索去踩坑。；


为什么需要 WebAssembly
自从 JavaScript 诞生起到现在已经变成最流行的编程语言，这背后正是 Web 的发展所推动的。Web 应用变得更多更复杂，但这也渐渐暴露出了 JavaScript 的问题：

语法太灵活导致开发大型 Web 项目困难；
性能不能满足一些场景的需要。
针对以上两点缺陷，近年来出现了一些 JS 的代替语言，例如：

微软的 TypeScript 通过为 JS 加入静态类型检查来改进 JS 松散的语法，提升代码健壮性；
谷歌的 Dart 则是为浏览器引入新的虚拟机去直接运行 Dart 程序以提升性能；
火狐的 asm.js 则是取 JS 的子集，JS 引擎针对 asm.js 做性能优化。
以上尝试各有优缺点，其中：

TypeScript 只是解决了 JS 语法松散的问题，最后还是需要编译成 JS 去运行，对性能没有提升；
Dart 只能在 Chrome 预览版中运行，无主流浏览器支持，用 Dart 开发的人不多；
asm.js 语法太简单、有很大限制，开发效率低。
三大浏览器巨头分别提出了自己的解决方案，互不兼容，这违背了 Web 的宗旨； 是技术的规范统一让 Web 走到了今天，因此形成一套新的规范去解决 JS 所面临的问题迫在眉睫。

于是 WebAssembly 诞生了，WebAssembly 是一种新的字节码格式，主流浏览器都已经支持 WebAssembly。 和 JS 需要解释执行不同的是，WebAssembly 字节码和底层机器码很相似可快速装载运行，因此性能相对于 JS 解释执行大大提升。 也就是说 WebAssembly 并不是一门编程语言，而是一份字节码标准，需要用高级编程语言编译出字节码放到 WebAssembly 虚拟机中才能运行， 浏览器厂商需要做的就是根据 WebAssembly 规范实现虚拟机。

WebAssembly 原理
要搞懂 WebAssembly 的原理，需要先搞懂计算机的运行原理。 电子计算机都是由电子元件组成，为了方便处理电子元件只存在开闭两种状态，对应着 0 和 1，也就是说计算机只认识 0 和 1，数据和逻辑都需要由 0 和 1 表示，也就是可以直接装载到计算机中运行的机器码。 机器码可读性极差，因此人们通过高级语言 C、C++、Rust、Go 等编写再编译成机器码。

由于不同的计算机 CPU 架构不同，机器码标准也有所差别，常见的 CPU 架构包括 x86、AMD64、ARM， 因此在由高级编程语言编译成可自行代码时需要指定目标架构。

WebAssembly 字节码是一种抹平了不同 CPU 架构的机器码，WebAssembly 字节码不能直接在任何一种 CPU 架构上运行， 但由于非常接近机器码，可以非常快的被翻译为对应架构的机器码，因此 WebAssembly 运行速度和机器码接近，这听上去非常像 Java 字节码。

相对于 JS，WebAssembly 有如下优点：

体积小：由于浏览器运行时只加载编译成的字节码，一样的逻辑比用字符串描述的 JS 文件体积要小很多；
加载快：由于文件体积小，再加上无需解释执行，WebAssembly 能更快的加载并实例化，减少运行前的等待时间；
兼容性问题少：WebAssembly 是非常底层的字节码规范，制订好后很少变动，就算以后发生变化,也只需在从高级语言编译成字节码过程中做兼容。可能出现兼容性问题的地方在于 JS 和 WebAssembly 桥接的 JS 接口。
每个高级语言都去实现源码到不同平台的机器码的转换工作是重复的，高级语言只需要生成底层虚拟机(LLVM)认识的中间语言(LLVM IR)，LLVM 能实现：

LLVM IR 到不同 CPU 架构机器码的生成；
机器码编译时性能和大小优化。
除此之外 LLVM 还实现了 LLVM IR 到 WebAssembly 字节码的编译功能，也就是说只要高级语言能转换成 LLVM IR，就能被编译成 WebAssembly 字节码，目前能编译成 WebAssembly 字节码的高级语言有：

AssemblyScript:语法和 TypeScript 一致，对前端来说学习成本低，为前端编写 WebAssembly 最佳选择；
c\c++:官方推荐的方式，详细使用见文档;
Rust:语法复杂、学习成本高，对前端来说可能会不适应。详细使用见文档;
Kotlin:语法和 Java、JS 相似，语言学习成本低，详细使用见文档;
Golang:语法简单学习成本低。但对 WebAssembly 的支持还处于未正式发布阶段，详细使用见文档。
通常负责把高级语言翻译到 LLVM IR 的部分叫做编译器前端，把 LLVM IR 编译成各架构 CPU 对应机器码的部分叫做编译器后端； 现在越来越多的高级编程语言选择 LLVM 作为后端，高级语言只需专注于如何提供开发效率更高的语法同时保持翻译到 LLVM IR 的程序执行性能。

编写 WebAssembly
AssemblyScript 初体验
接下来详细介绍如何使用 AssemblyScript 来编写 WebAssembly，实现斐波那契序列的计算。 用 TypeScript 实现斐波那契序列计算的模块 f.ts 如下：
export function f(x: i32): i32 {
    if (x === 1 || x === 2) {
        return 1;
    }
    return f(x - 1) + f(x - 2)
}
在按照 AssemblyScript 提供的安装教程成功安装后， 再通过

asc f.ts -o f.wasm
就能把以上代码编译成可运行的 WebAssembly 模块。

为了加载并执行编译出的 f.wasm 模块，需要通过 JS 去加载并调用模块上的 f 函数，为此需要以下 JS 代码：
fetch('f.wasm') // 网络加载 f.wasm 文件
    .then(res => res.arrayBuffer()) // 转成 ArrayBuffer
    .then(WebAssembly.instantiate) // 编译为当前 CPU 架构的机器码 + 实例化
    .then(mod => { // 调用模块实例上的 f 函数计算
    console.log(mod.instance.f(50));
    });
以上代码中出现了一个新的内置类型 i32，这是 AssemblyScript 在 TypeScript 的基础上内置的类型。 AssemblyScript 和 TypeScript 有细微区别，AssemblyScript 是 TypeScript 的子集，为了方便编译成 WebAssembly 在 TypeScript 的基础上加了更严格的类型限制， 区别如下：

比 TypeScript 多了很多更细致的内置类型，以优化性能和内存占用，详情文档;
不能使用 any 和 undefined 类型，以及枚举类型；
可空类型的变量必须是引用类型，而不能是基本数据类型如 string、number、boolean；
函数中的可选参数必须提供默认值，函数必须有返回类型，无返回值的函数返回类型需要是 void；
不能使用 JS 环境中的内置函数，只能使用 AssemblyScript 提供的内置函数。
总体来说 AssemblyScript 比 TypeScript 又多了很多限制，编写起来会觉得局限性很大； 用 AssemblyScript 来写 WebAssembly 经常会出现 tsc 编译通过但运行 WebAssembly 时出错的情况，这很可能就是你没有遵守以上限制导致的；但 AssemblyScript 通过修改 TypeScript 编译器默认配置能在编译阶段找出大多错误。

AssemblyScript 的实现原理其实也借助了 LLVM，它通过 TypeScript 编译器把 TS 源码解析成 AST，再把 AST 翻译成 IR，再通过 LLVM 编译成 WebAssembly 字节码实现； 上面提到的各种限制都是为了方便把 AST 转换成 LLVM IR。

为什么选 AssemblyScript 作为 WebAssembly 开发语言
AssemblyScript 相对于 C、Rust 等其它语言去写 WebAssembly 而言，好处除了对前端来说无额外新语言学习成本外，还有对于不支持 WebAssembly 的浏览器，可以通过 TypeScript 编译器编译成可正常执行的 JS 代码，从而实现从 JS 到 WebAssembly 的平滑迁移。

接入 Webpack 构建
任何新的 Web 开发技术都少不了构建流程，为了提供一套流畅的 WebAssembly 开发流程，接下来介绍接入 Webpack 具体步骤。

1. 安装以下依赖，以便让 TS 源码被 AssemblyScript 编译成 WebAssembly。
{
  "devDependencies": {
    "assemblyscript": "github:AssemblyScript/assemblyscript",
    "assemblyscript-typescript-loader": "^1.3.2",
    "typescript": "^2.8.1",
    "webpack": "^3.10.0",
    "webpack-dev-server": "^2.10.1"
  }
}
2. 修改 webpack.config.js，加入 loader：
module.exports = {
    module: {
        rules: [
            {
                test: /\.ts$/,
                loader: 'assemblyscript-typescript-loader',
                options: {
                    sourceMap: true,
                }
            }
        ]
    },
};
3. 修改 TypeScript 编译器配置 tsconfig.json，以便让 TypeScript 编译器能支持 AssemblyScript 中引入的内置类型和函数。
{
  "extends": "../../node_modules/assemblyscript/std/portable.json",
  "include": [
    "./**/*.ts"
  ]
}
4. 配置直接继承自 assemblyscript 内置的配置文件。

WebAssembly 相关文件格式
前面提到了 WebAssembly 的二进制文件格式 wasm，这种格式的文件人眼无法阅读，为了阅读 WebAssembly 文件的逻辑，还有一种文本格式叫 wast； 以前面讲到的计算斐波那契序列的模块为例，对应的 wast 文件如下：
func $src/asm/module/f (param f64) (result f64)
(local i32)
  get_local 0
  f64.const 1
  f64.eq
  tee_local 1
  if i32
    get_local 1
  else
    get_local 0
    f64.const 2
    f64.eq
  end
  i32.const 1
  i32.and
  if
    f64.const 1
    return
  end
  get_local 0
  f64.const 1
  f64.sub
  call 0
  get_local 0
  f64.const 2
  f64.sub
  call 0
  f64.add
end
这和汇编语言非常像，里面的 f64 是数据类型，f64.eq f64.sub f64.add 则是 CPU 指令。

为了把二进制文件格式 wasm 转换成人眼可见的 wast 文本，需要安装 WebAssembly 二进制工具箱WABT， 在 Mac 系统下可通过 brew install WABT 安装，安装成功后可以通过命令 wasm2wast f.wasm 获得 wast；除此之外还可以通过 wast2wasm f.wast -o f.wasm 逆向转换回去。

WebAssembly 相关工具
除了前面提到的 WebAssembly 二进制工具箱，WebAssembly 社区还有以下常用工具：

Emscripten: 能把 C、C++代码转换成 wasm、asm.js；
Binaryen: 提供更简洁的 IR，把 IR 转换成 wasm，并且提供 wasm 的编译时优化、wasm 虚拟机，wasm 压缩等功能，前面提到的 AssemblyScript 就是基于它。
WebAssembly JS API
目前 WebAssembly 只能通过 JS 去加载和执行，但未来在浏览器中可以通过像加载 JS 那样 <script src='f.wasm'></script> 去加载和执行 WebAssembly，下面来详细介绍如何用 JS 调 WebAssembly。

JS 调 WebAssembly 分为 3 大步：加载字节码 > 编译字节码 > 实例化，获取到 WebAssembly 实例后就可以通过 JS 去调用了，以上 3 步具体的操作是：

对于浏览器可以通过网络请求去加载字节码，对于 Nodejs 可以通过 fs 模块读取字节码文件；
在获取到字节码后都需要转换成 ArrayBuffer 后才能被编译，通过 WebAssembly 通过的 JS API WebAssembly.compile 编译后会通过 Promise resolve 一个 WebAssembly.Module，这个 module 是不能直接被调用的需要；
在获取到 module 后需要通过 WebAssembly.Instance API 去实例化 module，获取到 Instance 后就可以像使用 JS 模块一个调用了。
其中的第 2、3 步可以合并一步完成，前面提到的 WebAssembly.instantiate 就做了这两个事情。
WebAssembly.instantiate(bytes).then(mod=>{
  mod.instance.f(50);
})
WebAssembly 调 JS
之前的例子都是用 JS 去调用 WebAssembly 模块，但是在有些场景下可能需要在 WebAssembly 模块中调用浏览器 API，接下来介绍如何在 WebAssembly 中调用 JS。

WebAssembly.instantiate 函数支持第二个参数 WebAssembly.instantiate(bytes,importObject)，这个 importObject 参数的作用就是 JS 向 WebAssembly 传入 WebAssembly 中需要调用 JS 的 JS 模块。举个具体的例子，改造前面的计算斐波那契序列在 WebAssembly 中调用 Web 中的 window.alert 函数把计算结果弹出来，为此需要改造加载 WebAssembly 模块的 JS 代码
WebAssembly.instantiate(bytes,{
  window:{
    alert:window.alert
  }
}).then(mod=>{
  mod.instance.f(50);
})
对应的还需要修改 AssemblyScript 编写的源码：
// 声明从外部导入的模块类型
declare namespace window {
    export function alert(v: number): void;
}
 
function _f(x: number): number {
    if (x == 1 || x == 2) {
        return 1;
    }
    return _f(x - 1) + _f(x - 2)
}
 
export function f(x: number): void {
    // 直接调用 JS 模块
    window.alert(_f(x));
}
修改以上 AssemblyScript 源码后重新用 asc 通过命令 asc f.ts 编译后输出的 wast 文件比之前多了几行：
(import "window" "alert" (func $src/asm/module/window.alert (type 0)))
 
(func $src/asm/module/f (type 0) (param f64)
    get_local 0
    call $src/asm/module/_f
    call $src/asm/module/window.alert)
多出的这部分 wast 代码就是在 AssemblyScript 中调用 JS 中传入的模块的逻辑。

除了以上常用的 API 外，WebAssembly 还提供一些 API，你可以通过这个 d.ts 文件去查看所有 WebAssembly JS API 的细节。

不止于浏览器
WebAssembly 作为一种底层字节码，除了能在浏览器中运行外，还能在其它环境运行。

直接执行 wasm 二进制文件
前面提到的 Binaryen 提供了在命令行中直接执行 wasm 二进制文件的工具，在 Mac 系统下通过 brew install binaryen 安装成功后，通过 wasm-shell f.wasm 文件即可直接运行。

在 Node.js 中运行
目前 V8 JS 引擎已经添加了对 WebAssembly 的支持，Chrome 和 Node.js 都采用了 V8 作为引擎，因此 WebAssembly 也可以运行在 Node.js 环境中；

V8 JS 引擎在运行 WebAssembly 时，WebAssembly 和 JS 是在同一个虚拟机中执行，而不是 WebAssembly 在一个单独的虚拟机中运行，这样方便实现 JS 和 WebAssembly 之间的相互调用。

要让上面的例子在 Node.js 中运行，可以使用以下代码：
const fs = require('fs');
 
function toUint8Array(buf) {
    var u = new Uint8Array(buf.length);
    for (var i = 0; i < buf.length; ++i) {
        u[i] = buf[i];
    }
    return u;
}
 
function loadWebAssembly(filename, imports) {
    // 读取 wasm 文件，并转换成 byte 数组
    const buffer = toUint8Array(fs.readFileSync(filename));
    // 编译 wasm 字节码到机器码
    return WebAssembly.compile(buffer)
        .then(module => {
            // 实例化模块
            return new WebAssembly.Instance(module, imports)
        })
}
 
loadWebAssembly('../temp/assembly/module.wasm')
    .then(instance => {
        // 调用 f 函数计算
        console.log(instance.exports.f(10))
    });
在 Nodejs 环境中运行 WebAssembly 的意义其实不大，原因在于 Nodejs 支持运行原生模块，而原生模块的性能比 WebAssembly 要好。 如果你是通过 C、Rust 去编写 WebAssembly，你可以直接编译成 Nodejs 可以调用的原生模块。

WebAssembly 展望
从上面的内容可见 WebAssembly 主要是为了解决 JS 的性能瓶颈，也就是说 WebAssembly 适合用于需要大量计算的场景，例如：

在浏览器中处理音视频，flv.js 用 WebAssembly 重写后性能会有很大提升；
React 的 dom diff 中涉及到大量计算，用 WebAssembly 重写 React 核心模块能提升性能。Safari 浏览器使用的 JS 引擎 JavaScriptCore 也已经支持 WebAssembly，RN 应用性能也能提升；
突破大型 3D 网页游戏性能瓶颈，白鹭引擎已经开始探索用 WebAssembly。
总结
WebAssembly 标准虽然已经定稿并且得到主流浏览器的实现，但目前还存在以下问题：

浏览器兼容性不好，只有最新版本的浏览器支持，并且不同的浏览器对 JS WebAssembly 互调的 API 支持不一致；
生态工具不完善不成熟，目前还不能找到一门体验流畅的编写 WebAssembly 的语言，都还处于起步阶段；
学习资料太少，还需要更多的人去探索去踩坑。；


WebAssembly 完全是围绕提升原始执行速度而构建的。因此，如果我们希望这些代码能够获得快速、可预测的跨浏览器性能，可以考虑使用 WebAssembly。

WebAssembly 实现可预测的性能
通常，JavaScript 和 WebAssembly 可以达到相同的峰值性能。但是，对于 JavaScript 来说，这种性能只能在“快速路径”上实现，而且要保持在“快速路径”上并不容易。WebAssembly 的一个主要优势是可预测的性能，即使是跨浏览器也是如此。严格的类型和低级架构可以让编译器做出更强的保证，只需要对 WebAssembly 代码优化一次，就可以始终使用“快速路径”。

之前我们使用了 C/C++ 库，并将它们编译为 WebAssembly，以便在 Web 上使用它们。但实际上，我们并没有真正触及库的代码，我们只是写了少量的 C/C++ 代码作为浏览器和库之间的桥梁。但这次不一样：我们想要从头开始写一些东西，以便利用 WebAssembly 的优势。

WebAssembly 架构
在开始写代码之前，有必要先了解一下 WebAssembly。

引用 WebAssembly.org 的话：

WebAssembly（缩写为 Wasm）是一种用于栈虚拟机的二进制指令格式。Wasm 被设计为一个可移植的目标，用于编译 C/C++/Rust 等高级语言，支持在 Web 上部署客户端和服务器应用程序。

在将一段 C 语言或 Rust 代码编译为 WebAssembly 后，会得到一个包含模块声明的.wasm 文件。声明中包含了一个导入列表、一个导出列表（函数、常量、内存块）和函数的二进制指令。

有一些需要注意的东西：WebAssembly 虚拟机栈并没有保存在 WebAssembly 模块所使用的内存块中。虚拟机栈完全处在虚拟机内部，Web 开发人员无法访问它（除了通过 DevTools）。因此，我们可以编写完全不需要任何额外内存（只是有虚拟机内部栈）的 WebAssembly 模块。

在我们的例子中，我们需要使用一些额外的内存来访问图像的像素并生成图像的旋转版本。这个时候要用到 WebAssembly.Memory。

内存管理
通常，一旦你使用了额外的内存，就需要以某种方式来管理内存。内存的哪些部分正在使用中？哪些部分是可用的？例如，C 语言提供了 malloc(n) 函数，用来查找连续 n 个字节的内存空间。这种功能也被称为“分配器”。分配器需要被包含在 WebAssembly 模块中，这样会增加文件的大小。根据算法的不同，这些内存管理功能的体积和性能可能会有很大差异，这就是为什么很多语言提供了多种实现（“dmalloc”、“emmalloc”、“wee_alloc”……）。

在我们的例子中，在运行 WebAssembly 模块之前，我们知道输入图像的尺寸（以及输出图像的尺寸）。通常我们会将输入图像的 RGBA 缓冲区作为参数传给 WebAssembly 函数，并将旋转后的图像作为值返回。要生成这个返回值，我们需要使用分配器。但因为我们知道所需的内存总量（输入图像大小的两倍，一次用于输入，一次用于输出），所以可以使用 JavaScript 将输入图像放入 WebAssembly 内存，运行 WebAssembly 模块生成旋转图像，然后使用 JavaScript 回读结果。这样我们就可以不使用内存管理！

https://storage.googleapis.com/webfundamentals-assets/hotpath-with-wasm/animation_2_vp8.webm

如果你看一下原始的 JavaScript 函数，你会发现它其实是一些纯粹的计算代码，没有使用特定的 JavaScript API，所以可以很容易地将这些代码移植到其他语言。我们评估了 3 种可编译为 WebAssembly 的语言：C/C++、Rust 和 AssemblyScript。我们唯一要解决的问题是：如何在不使用内存管理功能的情况下访问原始内存？

C 语言和 Emscripten
Emscripten 是用于将 C 语言编译成 WebAssembly 的编译器。Emscripten 的目标是成为 GCC 或 clang 等知名 C 语言编译器的直接替代品。这是 Emscripten 的核心任务，它旨在尽可能简单地将现有 C 语言和 C++ 代码编译为 WebAssembly。

访问原始内存是 C 语言的本质，指针的存在就是为了这个：

复制代码
uint8_t* ptr = (uint8_t*)0x124;
ptr[0] = 0xFF;
我们将数字 0x124 转换为指向无符号 8 位整数（或字节）的指针，将 ptr 变量变成从内存地址 0x124 开始的数组，并且可以像使用其他数组一样使用它。在我们的例子中，我们想要重新排序图像的 RGBA 缓冲区，以便实现图像旋转。要移动一个像素，我们需要每次移动 4 个连续字节（每个通道一个字节：R、G、B 和 A）。为此，我们创建了一个无符号的 32 位整数数组。按照惯例，我们的输入图像将从地址 4 开始，输出图像从输入图像结束位置开始：

复制代码
int bpp = 4;
int imageSize = inputWidth * inputHeight * bpp;
uint32_t* inBuffer = (uint32_t*) 4;
uint32_t* outBuffer = (uint32_t*) (inBuffer + imageSize);
 
for (int d2 = d2Start; d2 >= 0 && d2 < d2Limit; d2 += d2Advance) {
  for (int d1 = d1Start; d1 >= 0 && d1 < d1Limit; d1 += d1Advance) {
    int in_idx = ((d1 * d1Multiplier) + (d2 * d2Multiplier));
    outBuffer[i] = inBuffer[in_idx];
    i += 1;
  }
}
在将整个 JavaScript 函数移植到 C 语言后，可以使用 emcc 编译 C 文件：

复制代码
$ emcc -O3 -s ALLOW_MEMORY_GROWTH=1 -o c.js rotate.c
与往常一样，Emscripten 会生成一个叫作 c.js 的胶水代码文件和一个叫作 c.wasm 的 wasm 模块。请注意，wasm 模块被压缩后只有 260 字节左右，而胶水代码在压缩大约是 3.5KB。经过一些调整之后，我们可以去掉胶水代码，并使用普通 API 来实例化 WebAssembly 模块。

Rust
Rust 是一门全新的现代编程语言，它提供了丰富的类型系统，没有运行时和所有权模型，可确保内存安全性和线程安全性。Rust 还将 WebAssembly 视为一等公民，而且 Rust 团队还为 WebAssembly 生态系统贡献了很多优秀的工具。

其中一个工具是由 rustwasm 工作组开发的 wasm-pack（https://rustwasm.github.io/wasm-pack/）。wasm-pack 可以将你的代码转换为一个对 Web 友好的模块，支持 webpack 等捆绑器，但目前仅适用于 Rust。这个工作小组正在考虑增加对其他语言的支持。

Rust 中的切片相当于 C 语言中的数组。就像在 C 语言中一样，我们需要创建切片。这违反了 Rust 的内存安全模型，因此我们必须使用 unsafe 关键字来编写不遵循内存安全模型的代码。

复制代码
let imageSize = (inputWidth * inputHeight) as usize;
let inBuffer: &mut [u32];
let outBuffer: &mut [u32];
unsafe {
  inBuffer = slice::from_raw_parts_mut::<u32>(4 as *mut u32, imageSize);
  outBuffer = slice::from_raw_parts_mut::<u32>((imageSize * 4 + 4) as *mut u32, imageSize);
}
 
for d2 in 0..d2Limit {
  for d1 in 0..d1Limit {
    let in_idx = (d1Start + d1 * d1Advance) * d1Multiplier + (d2Start + d2 * d2Advance) * d2Multiplier;
    outBuffer[i as usize] = inBuffer[in_idx as usize];
    i += 1;
  }
}
编译 Rust 文件：

复制代码
$ wasm-pack build
这个命令将产生一个 7.6KB 的 wasm 模块，以及大约 100 个字节的胶水代码（压缩之后）。

AssemblyScript
AssemblyScript 是一个相当年轻的项目，用于将 TypeScript 编译成 WebAssembly。AssemblyScript 使用与 TypeScript 相同的语法，但使用了自己的标准库。它们的标准库模拟了 WebAssembly 的功能。这意味你无法将任意 TypeScript 代码编译成 WebAssembly，但确实意味着你不必为了编写 WebAssembly 而去学习新的编程语言！

复制代码
for (let d2 = d2Start; d2 >= 0 && d2 < d2Limit; d2 += d2Advance) {
  for (let d1 = d1Start; d1 >= 0 && d1 < d1Limit; d1 += d1Advance) {
    let in_idx = ((d1 * d1Multiplier) + (d2 * d2Multiplier));
    store<u32>(offset + i * 4 + 4, load<u32>(in_idx * 4 + 4));
    i += 1;
  }
}
因为 rotate() 函数具有较小的类型表面，可以很容易将其移植到 AssemblyScript。AssemblyScript 提供了用于访问原始内存的函数load<T>(ptr: usize)和store<T>(ptr: usize, value: T)。要编译我们的AssemblyScript 文件，只需要安装 AssemblyScript/assemblyscript 包，并运行：

复制代码
$ asc rotate.ts -b assemblyscript.wasm --validate -O3
AssemblyScript 将生成约 300 字节的 wasm 模块，并且没有胶水代码。这个模块可以与 WebAssembly API 一起使用。

瘦身
与其他两种语言相比，Rust 的 7.6KB 显得非常大。WebAssembly 生态系统中有一些工具可以用来分析 WebAssembly 文件，可以告诉你发生了什么，并帮你改善这种情况。

twiggy
twiggy是 Rust 团队开发的另一个工具，可以从 WebAssembly 模块中提取大量有用的信息。这个工具并不是特定于 Rust 的，可用来检查模块调用图、找出未使用或多余的部分，以及哪些部分占用了模块的文件大小。可以使用 twiggy 的 top 命令来查看模块文件的组成：

复制代码
$ twiggy top rotate_bg.wasm


我们可以看到，大部分文件大小来自分配器。这个有点让我们感到惊讶，因为我们的代码没有使用动态分配功能。另一个占用较大体积的是“函数名”。

wasm-strip
wasm-strip 是来自WebAssembly Binary Toolkit，简称为 wabt 的一个工具。它提供了一些工具，可用于检查和操作 WebAssembly 模块。wasm2wat 是一个反汇编程序，可以将二进制 wasm 模块转换为人类可读的格式。wabt 还包含了 wat2wasm，可以将人类可读的格式转换回二进制 wasm 模块。我们确实有使用这两个工具来检查 WebAssembly 文件，不过我们发现 wasm-strip 是最有用的。wasm-strip 从 WebAssembly 模块中移除了不必要的部分和元数据：

复制代码
$ wasm-strip rotate_bg.wasm
这样就可以将 Rust 模块文件大小从 7.5KB 减小到 6.6KB（在压缩之后）。

wasm-opt
wasm-opt 是来自Binaryen的一个工具。它尝试基于字节码对 WebAssembly 模块进行大小和性能方面的优化。Emscripten 已经在使用这个工具，有些编译器则没有。使用这些工具来节省一些额外的字节是个好主意。

复制代码
wasm-opt -O3 -o rotate_bg_opt.wasm rotate_bg.wasm
通过使用 wasm-opt，我们可以减少另外一些字节，在压缩之后只有 6.2KB。

#![no_std]
经过一些咨询和研究，我们使用#![no_std]来重新编写 Rust 代码，这样就可以不使用 Rust 的标准库。这样就可以完全禁用动态内存分配，从而从模块中删除了分配器代码。使用以下命令编译 Rust 文件：

复制代码
$ rustc --target=wasm32-unknown-unknown -C opt-level=3 -o rust.wasm rotate.rs
在使用了 wasm-opt 和 wasm-strip 之后，压缩的 wasm 模块只剩下 1.6KB。虽然它仍然比 C 语言编译器和 AssemblyScript 生成的模块大，但也足以称得上是一个轻量级的模块。

性能
除了文件大小，我们还需要优化性能。那么我们应该如何衡量性能？它们的结果又是怎样的呢？

如何进行基准测试
尽管 WebAssembly 是一种低级的字节码格式，但仍需要通过编译器生成特定于主机的机器码。就像 JavaScript 一样，编译器包含了多个阶段的工作。简单地说：第一阶段编译速度较快，但生成的代码运行速度较慢。在模块开始运行后，浏览器就会观察哪些部分是经常使用的，并通过一个更优化但速度更慢的编译器发送这些部分。

我们的用例很有趣，旋转图像的代码可能会被使用一次，或者两次。因此，在绝大多数情况下，我们无法从优化编译器中获得好处。在进行基准测试时要记住这一点。循环运行 WebAssembly 模块 10,000 次会产生不真实的结果。为了获得更真实的数字，我们应该只运行一次模块，并根据单次运行的结果做出判断。

性能比较



这两张图是相同数据的不同视图。在第一张图中，我们根据浏览器来比较，在第二张图中，我们根据使用的语言来比较。请注意，我使用了对数时间尺度，而且所有基准测试使用了相同的 1600 万像素的测试图像和相同的主机。

从图中可以看出，我们解决了原始性能问题：所有 WebAssembly 模块的运行时间都在大约 500 毫秒或更短的时间内。这证实了我们在开始时的假设：WebAssembly 为我们提供了可预测的性能。无论我们选择哪种语言，浏览器和语言之间的差异都很小。确切地说：JavaScript 的跨浏览器标准偏差约为 400 毫秒，而 WebAssembly 模块的跨浏览器标准偏差约为 80 毫秒。

工作量
另一个度量指标是创建 WebAssembly 模块并将其集成到 squoosh 的工作量。我们很难使用准确的数值来表示工作量，所以我不会创建任何图表，不过我想指出一些东西：

AssemblyScript 不仅让我们可以使用 TypeScript 来编写 WebAssembly，进行代码评审也非常容易，而且还可以生成非常小且具有良好性能的无胶水 WebAssembly 模块。

Rust 与 wasm-pack 结合使用也非常方便，但在大型的 WebAssembly 项目（需要用到绑定和内存管理）中表现更好。我们必须付出额外的工作量才能获得有竞争力的文件大小。

C 语言和 Emscripten 可以生成非常小巧且高性能的 WebAssembly 模块，但是如果没有勇气直接使用胶水代码并将其缩减到最基本的需求，那么总体大小（WebAssembly 模块 + 胶水代码）就会变得非常大。