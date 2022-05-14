---
title: Go 编程语言规范
layout: post
category: golang
author: 夏泽民
---
https://moego.me/golang_spec.html
https://golang.org/ref/spec
<!-- more -->
本文翻译自 The Go Programming Language Specification (https://golang.org/ref/spec)，原文采用 Creative Commons Attribution 3.0 协议，文档内代码采用 BSD 协议 (https://golang.org/LICENSE)。 本文采用 Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International 许可协议，文档内代码继用相同协议，如果你需要发布本文（或衍生品），也需要注明本文原始链接 (https://moego.me/golang_spec.html) 及译者 Bekcpear。
对应英文原版 为 2019 年 07 月 31 日 版本: golang.org/ref/spec
翻译中针对可能有歧义/不明确/翻译后不易于理解的单词将直接使用原词汇
为了行文工整，代码块内可能使用英文表述
因为学习语言需要，所以翻译
有些翻译可能比较迷糊，我会在进一步学习后完善它们
文中实参仅代表 argument；参数仅代表 parameter，有时候也会译为形参
介绍
这是一篇 Go 编程语言的参考手册。访问 golang.org 以获取更多信息及其它文档。

Go 是一个在设计时便考虑到系统编程的通用语言。它是强类型的、带垃圾回收的并明确支持了并发编程。程序是由包所构建的，包属性支持高效的依赖管理。

语言语法紧凑且常规，以便于如集成开发环境（IDE）这样子的自动化工具所分析。

标记法
标记法语法指定使用扩展巴科斯-瑙尔范式（EBNF）:

Production  = production_name, "=", [ Expression ], "." .
Expression  = Alternative, { "|", Alternative } .
Alternative = Term, { Term } .
Term        = production_name | token, [ "…", token ] | Group | Option | Repetition .
Group       = "(", Expression, ")" .
Option      = "[", Expression, "]" .
Repetition  = "{", Expression, "}" .
产生式Productions 是由 术语terms 和如下操作符所构建的表达式（操作符排列按优先级递增的顺序）:

|   多选一
()  分组
[]  零或一
{}  零或多
小写字母的产生式名是用来标记一个词汇记号（组）的。 非终结符Non-terminals 是以驼峰命名法命名的。词汇记号（ 终结符terminals ）都是使用双引号 "" 或者反引号 `` 包裹起来的。

a … b 这样子的格式表示从 a 连续到 b 的字符集。水平省略号 … 也会用在其它一些地方非正式地表示枚举或者不再进一步说明的代码片段。 字符 … （与三个单独字符 ... 不同）并不是 Go 语言里的 token。

译注： 扩展巴科斯-瑙尔范式extended Backus-Naur form 是一种 元语法metasyntax 符号标记法，可以用于表示 上下文无关文法Context-free grammar 。

针对本文简单说明，其产生式规则由非终结符和终结符所构成，左侧是一个非终结符，右侧则是该非终结符所代表的终结符和非终结符。终结符包括字母数字字符、标点符号和空格字符，其不可再分；非终结符最终指代某种序列组合的多个终结符。

本文用到的上述未说明的范式符号说明： = 定义； , 级联； . 表示表达式终结； " .. " 表示除双引号外的终结符； ` .. ` 表示除反引号外的终结符； ? .. ? 表示特殊序列，用于解释 EBNF 标准以外的文本。

又注：根据维基百科 extended Backus-Naur form 上说明来看，原文的 EBNF 格式并不规范，所以我对原文表达式进行最小程度修改。更详细的 EBNF 说明可以下载 ISO/IEC 14977:1996 PDF 压缩档 查看。

段落名若为中文且在语法标记块中使用英文书写的，均会在段落名上一并附上英文。

源代码表示
源代码是以 UTF-8 编码的 Unicode 文本。该文本并不是规范化的，所以一个单一的带重音符（附加符）的码位和由重音符（附加符）和字母所组成的相同字符不同，该相同字符结构被看成两个码位。为了简便，本文档使用非正规的术语——字符——指代源文本中的 Unicode 码位。

译注： 这里的 规范化 的含义是指，文字处理软件为了对 Unicode 字符串做比较、搜寻和排序操作而不得不考虑其等价性才做的正规化处理，参考维基百科 Unicode 等價性 。
每一个码位都是不同的，比如大写和小写的字母就是不同的字符。

实现限制：为了保证与其它工具的兼容性，编译器可能会不允许源文本中存在 NUL 字符（U+0000）。

实现限制：为了保证与其它工具的兼容性，如果一个 UTF-8 编码的字节顺序标记（U+FEFF）为源文本的第一个 Unicode 码位，编译器可能会忽略它。字节顺序标记也可能会被不允许出现在源中的任何其它位置。

字符
如下术语用于表示特定的 Unicode 字符类:

newline        = ? Unicode 码位 U+000A ? .
unicode_char   = ? newline 以外的任意 Unicode 码位 ? .
unicode_letter = ? 被分类为「字母」的 Unicode 码位 ? .
unicode_digit  = ? 被分类为「数字/十进制数」的 Unicode 码位 ? .
在 The Unicode Standard 8.0 中， 4.5 节 "General Category" 定义了一套字符类别。 Go 语言把类别 Lu, Ll, Lt, Lm 或 Lo 中的字符看作 Unicode 字母，把数字类别 Nd 中的字符看作 Unicode 数字。

译注： Lu 为大写字母， Ll 为小写字母， Lt 为标题字母， Lm 为修饰字母， Lo 为其它字母， Nd 为十进制数字，可以在 Compart 上查到对应分类包含哪些字符。 但是在这里我有一个疑惑，里面明明很多字母和数字是不能用在标识符中的，为什么这里统统包含了进来，并且下文也没有额外的说明？
字母和数字
下划线字符 _ (U+005F) 被认为是一个字母。

letter        = unicode_letter | "_" .
decimal_digit = "0" … "9" .
binary_digit  = "0" | "1" .
octal_digit   = "0" … "7" .
hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
词法元素
注释
注释作为程序的文档，有两种格式：

行内注释从字符序列 // 开始并在一行末尾结束。
通用注释从字符序列 /* 开始并在遇到的第一个字符序列 */ 时结束。
注释不能开始于 rune 或 字符串 字面值或另一个注释的内部。不包含新行的通用注释就像一个空格。任何其它的注释就像一空白行。

Tokens
Tokens 组成了 Go 语言的词汇表。有四个分类： 标识符 、 关键字 、 运算符和标点 以及 字面值 。 空白 是由空格（U+0020）、水平制表（U+0009）、回车（U+000D）和新行（U+000A）所组成的，空白一般会被忽略，除非它分隔了组合在一起会形成单一 token 的 tokens. 并且，新行或者文件结尾可能会触发 分号 的插入。当把输入的内容区分为 tokens 时，每一个 token 都是可组成有效 token 的最长字符序列。

分号
正式的语法使用分号 ; 作为一定数量的产生式的终结符。 Go 程序可以依据如下两条规则来省略大部分这样子的分号：

输入内容被分为 tokens 时，当每一行最后一个 token 为以下 token 时，一个分号会自动插入到其后面：
标识符
整数 、 浮点数 、 虚数 、 rune 或者 字符串 字面值
关键字 break , continue , fallthrough 或 return 之一
运算符和标点 ++ , -- , ) , ] 或 } 之一
为了使复杂的语句可以占据在单一一行上，分号也可以在关闭的 ) 或者 } 前被省略。
为了反应出惯用的使用习惯，本文档中的代码示例将参照这些规则来省略掉分号。

标识符Identifiers
标识符用于命名程序中的实体——比如变量和类型。它是一个或者多个字母和数字的序列组合。标识符的第一个字符必须是一个字母。

identifier = letter, { letter | unicode_digit } .
a
_x9
ThisVariableIsExported
αβ
有一些标识符已经被 预先声明 了。

关键字
如下关键字是保留的，不可以用作标识符。

break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
运算符和标点
如下的字符序列用于代表 运算符 （包括了 赋值运算符 ）和标点:

+    &     +=    &=     &&    ==    !=    (    )
-    |     -=    |=     ||    <     <=    [    ]
*    ^     *=    ^=     <-    >     >=    {    }
/    <<    /=    <<=    ++    =     :=    ,    ;
%    >>    %=    >>=    --    !     ...   .    :
     &^          &^=
整数字面值Integer literals
整数字面值是用来代表整数 常量 的数字序列。可用一个可选前缀来设置非十进制数： 0b 或 0B 代表二进制， 0 , 0o , 0O 代表八进制， 0x 或 0X 代表十六进制。单独的 0 被视作十进制零。在十六进制数字面值中，字母 a 到 f 以及 A 到 F 代表数字值 10 到 15 。

为了可读性，下划线字符 _ 可以出现在基本前缀之后或者连续的数字之间；这样的下划线不改变字面值的值。

int_lit        = decimal_lit | binary_lit | octal_lit | hex_lit .
decimal_lit    = "0" | ( "1" … "9" ), [ [ "_" ], decimal_digits ] .
binary_lit     = "0", ( "b" | "B" ), [ "_" ], binary_digits .
octal_lit      = "0", [ "o" | "O" ], [ "_" ], octal_digits .
hex_lit        = "0", ( "x" | "X" ), [ "_" ], hex_digits .

decimal_digits = decimal_digit, { [ "_" ], decimal_digit } .
binary_digits  = binary_digit, { [ "_" ], binary_digit } .
octal_digits   = octal_digit, { [ "_" ], octal_digit } .
hex_digits     = hex_digit, { [ "_" ], hex_digit } .
42
4_2
0600
0_600
0o600
0O600       // 第二个字符是大写字母 'O'
0xBadFace
0xBad_Face
0x_67_7a_2f_cc_40_c6
170141183460469231731687303715884105727
170_141183_460469_231731_687303_715884_105727

_42         // 这是一个标识符，而不是一个整数字面值
42_         // 无效: _ 必须分隔连续数字
4__2        // 无效: 一次只能有一个 _
0_xBadFace  // 无效: _ 必须分隔连续数字
浮点数字面值Floating-point literals
浮点数字面值是浮点数 常量 的十进制或十六进制表示。

十进制的浮点数字面值由一个整数部分（十进制数字），一个小数点，一个小数部分（十进制数字）和一个指数部分（ e 或 E 后紧跟着带或者不带符号且为十进制的数字）。整数部分和小数部分其中之一可以省略；小数点和指数部分其中之一可以省略。指数值 exp 以 10exp 来缩放 有效数字mantissa （整数和小数部分）。

译注： "An exponent value exp scales the mantissa (integer and fractional part) by 10exp ." 这里的 "mantissa" 存在争议，目前 IEEE 使用的是 "significand" 一词，维基百科 Talk:Significand 整理了相关讨论。
十六进制浮点数字面值由一个 0x 或 0X 前缀，一个整数部分（十六进制数字），一个小数点，一个小数部分（十六进制数字）和一个指数部分（ p 或 P 后紧跟着带或者不带符号且为十六进制的数字）。整数部分和小数部分其中之一可以省略；小数点也可以省略，但是指数部分是必须的。（这个语法匹配 IEEE 754-2008 §5.12.3 章所说的。）指数值 exp 以 2exp 来缩放有效数字（整数和小数部分）。

为了可读性，下划线字符 _ 可以出现在基本前缀后或是连续的数字之间；这样的下划线不会改变字面值的值。

float_lit         = decimal_float_lit | hex_float_lit .

decimal_float_lit = decimal_digits, ".", [ decimal_digits ], [ decimal_exponent ] |
                    decimal_digits, decimal_exponent |
                    ".", decimal_digits, [ decimal_exponent ] .
decimal_exponent  = ( "e" | "E" ), [ "+" | "-" ], decimal_digits .

hex_float_lit     = "0", ( "x" | "X" ), hex_mantissa, hex_exponent .
hex_mantissa      = [ "_" ], hex_digits, ".", [ hex_digits ] |
                    [ "_" ], hex_digits |
                    ".", hex_digits .
hex_exponent      = ( "p" | "P" ), [ "+" | "-" ], decimal_digits .
0.
72.40
072.40       // == 72.40
2.71828
1.e+0
6.67428e-11
1E6
.25
.12345E+5
1_5.         // == 15.0
0.15e+0_2    // == 15.0

0x1p-2       // == 0.25
0x2.p10      // == 2048.0
0x1.Fp+0     // == 1.9375
0X.8p-0      // == 0.5
0X_1FFFP-16  // == 0.1249847412109375
0x15e-2      // == 0x15e - 2 （整数减法）

0x.p1        // 无效的： 有效数字无数字
1p-2         // 无效的： p 指数需要十六进制有效数字
0x1.5e-2     // 无效的： hexadecimal mantissa requires p exponent
1_.5         // 无效的： _ 必须分隔连续的数字
1._5         // 无效的： _ 必须分隔连续的数字
1.5_e1       // 无效的： _ 必须分隔连续的数字
1.5e_1       // 无效的： _ 必须分隔连续的数字
1.5e1_       // 无效的： _ 必须分隔连续的数字
虚数字面值Imaginary literals
虚数字面值表示复数 常量 的虚部。它由 整数 或者 浮点数 字面值紧跟着一个小写的字母 i 组成。这个虚数字面值的值为对应整数或者浮点数字面值的值乘以虚数单位 i 。

imaginary_lit = (decimal_digits | int_lit | float_lit), "i" .
考虑到向后兼容，完全由十进制数字（可能存在下划线）组成的虚数字面值的整数部分被作为十进制整数，即使其以 0 开头也不例外。

0i
0123i         // == 123i 为了向后兼容
0o123i        // == 0o123 * 1i == 83i
0xabci        // == 0xabc * 1i == 2748i
0.i
2.71828i
1.e+0i
6.67428e-11i
1E6i
.25i
.12345E+5i
0x1p-2i       // == 0x1p-2 * 1i == 0.25i
Rune 字面值Rune literals
Rune 字面值代表了一个 rune 常量 ，一个确定了 Unicode 码位的整数值。 Rune 字面值是由一个或者多个字符以单引号包裹来表示的，就像 'x' 或 '\n' 。在引号内，除了新行和未被转义的单引号外的任何字符都可能出现。被单引的字符表示的是该字符的 Unicode 值，不过以反斜杠开头的多字符序列会以不同的格式来编码 Unicode 值。

这是在引号内代表单一字符的最简单的形式；因为 Go 源文件是使用 UTF-8 编码的 Unicode 字符，多个 UTF-8 编码的字节可以表示为一个单一整数值。比如： 'a' 用一个字节代表了字面值 a ， Unicode U+0061，值 0x61 ；但 'ä' 用了两个字节（ 0xc3 0xa4 ）代表了字面值 a 分音符 ， Unicode U+00E4，值 0xe4 。

几个反斜杠转义允许任意值被编码为 ASCII 文本。有四种方法将整数值表达为数值常量： \x 紧跟着两个十六进制数； \u 紧跟着四个十六进制数； \U 紧跟着八个十六进制数；一个单独的反斜杠 \ 紧跟着三个八进制数。每一种情况下的字面值的值都是对应基础上该数所表示的值。

虽然这些表示的最终都是一个整数，但它们有不同的有效范围。八进制转义必须表示 0 到 255 之间的值。十六进制转义满足条件的要求会因为构造不同而不同。 \u 和 \U 代表了 Unicode 码位，所以在这里面有一些值是非法的，尤其是那些超过 0x10FFFF 的和代理了一半的（译注：查阅「 UTF-16 代理对」进行深入阅读）。

在反斜杠后，某些单字符的转义代表了特殊的值:

\a   U+0007 警报或蜂鸣声
\b   U+0008 退格
\f   U+000C 换页
\n   U+000A 换行或新行
\r   U+000D 回车
\t   U+0009 水平制表
\v   U+000b 垂直制表
\\   U+005c 反斜杠
\'   U+0027 单引号（只在 rune 字面值中转义才有效）
\"   U+0022 双引号（只在字符串字面值中转义才有效）
所有其它以反斜杠开头的序列在 rune 字面值中都是非法的。

rune_lit         = "'", ( unicode_value | byte_value ), "'" .
unicode_value    = unicode_char | little_u_value | big_u_value | escaped_char .
byte_value       = octal_byte_value | hex_byte_value .
octal_byte_value = `\`, octal_digit, octal_digit, octal_digit .
hex_byte_value   = `\`, "x", hex_digit, hex_digit .
little_u_value   = `\`, "u", hex_digit, hex_digit, hex_digit, hex_digit .
big_u_value      = `\`, "U", hex_digit, hex_digit, hex_digit, hex_digit,
                             hex_digit, hex_digit, hex_digit, hex_digit, .
escaped_char     = `\`, ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) .
'a'
'ä'
'本'
'\t'
'\000'
'\007'
'\377'
'\x07'
'\xff'
'\u12e4'
'\U00101234'
'\''         // 包含了一个单引号字符的 rune 字面值
'aa'         // 非法：字符太多
'\xa'        // 非法：十六进制数字太少
'\0'         // 非法：八进制数字太少
'\uDFFF'     // 非法：只代理了一半
'\U00110000' // 非法：无效的 Unicode 码位
字符串字面值String literals
字符串字面值代表了通过串联字符序列而获得的字符串 常量 。它有两种形式： 原始raw 字符串字面值和 解释型interpreted 字符串字面值。

原始字符串字面值是在反引号之间的字符序列，就像 `foo` 。除了反引号外的任何字符都可以出现在该引号内。原始字符串字面值的值就是由在引号内未被解释过的（隐式 UTF-8 编码的）字符所组成的字符串；比如，反斜杠在这里没有其它特殊的意义，并且可以包含新行。原始字符串字面值中的回车字符（ '\r' ）是会被从原始字符串值中所丢弃。

译注： 经测试，手动输入的 '\r' 字符是可以正常显示为 '\r' 的，那么理解下来，丢弃的是键盘键入的回车。
解释型字符串字面值是在双引号之间的字符序列，就像 "bar" 。除了新行和未被转义的双引号之外的所有字符都可以出现在该引号内。引号之间的文本组成了字符串字面值的值，反斜杠转义以及限制都和 rune 字面值一样（不同的是，在解释型字符串字面值中， \' 是非法的， \" 是合法的）。三个数字的八进制数（ \nnn ）和两个数字的十六进制数（ \xnn ）的转义代表着所生成字符串的独立的字节；所有其它的转义代表了单独字符的 UTF-8 编码（可能是多字节的）。因此字符串字面值内的 \377 和 \xFF 代表着值为 0xFF=255 的单一字节，而 ÿ , \u00FF , \U000000FF 和 \xc3\xbf 代表着字符 U+00FF 以 UTF-8 编码的双字节 0xc3 0xbf 。

string_lit             = raw_string_lit | interpreted_string_lit .
raw_string_lit         = "`", { unicode_char | newline }, "`" .
interpreted_string_lit = `"`, { unicode_value | byte_value }, `"` .
`abc`                // 同 "abc"
`\n
n`                   // 同 "\\n\n\\n"
"\n"
"\""                 // 同 `"`
"Hello, world!\n"
"日本語"
"\u65e5本\U00008a9e"
"\xff\u00FF"
"\uD800"             // 非法: 代理了一半
"\U00110000"         // 非法: 无效的 Unicode 码位
以下这些例子都代表着相同的字符串：

"日本語"                                 // UTF-8 输入文本
`日本語`                                 // 以原始字面值输入的 UTF-8 文本
"\u65e5\u672c\u8a9e"                    // 明确的 Unicode 码位
"\U000065e5\U0000672c\U00008a9e"        // 明确的 Unicode 码位
"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"  // 明确的 UTF-8 字节
当源代码以两个码位来代表一个字符，比如包含一个重音符和一个字母的组合形式，如果是在 rune 字面值中的话会使得结果出错（因为其并不是一个单一码位），而如果是在字符串字面值中的话则会显示为两个码位。

常量
常量有 布尔值常量 、 rune 常量 、 整数常量 、 浮点数常量 、 复数常量 和 字符串常量 。 Rune、整数、浮点数和复数常量统称为数值常量。

一个常量的值是由如下所表示的： rune ； 整数 ； 浮点数 ； 虚数 ； 字符串 字面值；表示常量的标识符； 常量表达式 ；结果为常量的 变量转换 ；或者一些内置函数所生成的值，这些内置函数比如应用于任意值的 unsafe.Sizeof ，应用于 一些表达式 的 cap 或 len ，应用于复数常量的 real 和 imag 以及应用于数值常量的 complex 。布尔值是由预先声明的常量 true 和 false 所代表的。预先声明的标识符 iota 表示一个整数常量。

通常，复数常量是 常量表达式 的一种形式，会在该节讨论。

数值常量代表任意精度的确切值，而且不会溢出。因此，没有常量表示 IEEE-754 负零，无穷，以及非数字值集。

常量可以是有 类型 的也可以是无类型的。字面值常量， true , false , iota 以及一些仅包含无类型的恒定操作数的 常量表达式 是无类型的。

常量可以通过 常量声明 或 变量转换 被显示地赋予一个类型，也可以在 变量声明 或 赋值 中，或作为一个操作数在 表达式 中使用时隐式地被赋予一个类型。如果常量的值不能按照所对应的类型来表示的话，就会出错。「前一版的内容： 比如， 3.0 可以作为任何整数类型或任何浮点数类型，而 2147483648.0 （相当于 1<<31 ）可以作为 float32 , float64 或 uint32 类型，但不能是 int32 或 string 。」

一个无类型的常量有一个 默认类型 ，当在上下文中需要请求该常量为一个带类型的值时，这个 默认类型 便指向该常量隐式转换后的类型，比如像 i := 0 这样子的 短变量声明 就没有显示的类型。无类型常量的默认类型分别是 bool , rune , int , float64 , complex128 或 string ，取决于它是否是一个布尔值、 rune、整数、浮点数、复数或字符串常量。

实现限制：虽然数值常量在这个语言中可以是任意精度的，但编译器可能会使用精度受限的内部表示法来实现它。也就是说，每一种实现必须：

使用最少 256 位来表示整数。
使用最少 256 位来表示浮点数常量（包括复数常量的对应部分）的小数部分，使用最少 16 位表示其带符号的二进制指数部分。
当无法表示一个整数常量的精度时，需要给出错误。
当因为溢出而无法表示一个浮点数或复数常量时，需要给出错误。
当因为精度限制而无法表示一个浮点数或复数常量时，约到最接近的可表示的常量。
这些要求也适用于字面值常量，以及 常量表达式 的求值结果。

变量
变量是用来放置 值 的存储位置。可允许的值的集是由变量 类型 所确定的。

变量声明 和对于函数参数及其结果而言的 函数声明 或 函数字面值 的签名都为命名的变量保留存储空间。调用内置函数 new 或获取 复合字面值 的地址会在运行时为变量分配存储空间。这样子的一个匿名变量是通过（可能隐式的） 指针间接 引用到的。

结构化的 数组 、 分片 和 结构体 类型变量存在可以独立 寻址 的元素和字段。每一个这样子的元素就像一个变量。

变量的 静态类型 （或者就叫 类型 ）是其声明时确定好的类型，或由 new 调用/复合字面值所提供的类型，或结构化变量的元素类型。接口类型的变量还有一个独特的 动态 类型，该类型是在运行时所分配给变量的值的具体类型（除非那个值是预声明的无类型的标识符 nil ）。动态类型可能会在执行过程中变化，但存储在接口变量中的值始终 可分配 为接口变量的静态类型。

var x interface{}  // x 是 nil，它有一个静态类型 interface{}
var v *T           // v 的值为 nil，静态类型为 *T
x = 42             // x 的值为 42，动态类型为 int
x = v              // x 的值为 (*T)(nil)，动态类型为 *T
变量的值是通过引用 表达式 中的变量来检索的；它总是那个最后 赋 给变量的值。如果一个变量还没有被分配到值，那么它的值是其对应类型的 零值 。

类型
类型确定了一个值集（连同特定于这些值的操作和方法）。类型可能是由 类型名 所表示的（如果它有的话），或者使用 类型字面值 所指定（由已知类型组成的类型）。

Type      = TypeName | TypeLit | "(", Type, ")" .
TypeName  = identifier | QualifiedIdent .
TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
            SliceType | MapType | ChannelType .
语言本身 预先声明 了一些特定的类型名。其它命名的类型则使用 类型声明 引入。 复合类型 ——数组、结构体、指针、函数、接口、分片、映射和信道类型——可以由类型字面值构成。

每个类型 T 都有一个 潜在类型 ：如果 T 是预先声明的布尔值、数值或者字符串类型之一，或一个类型字面值，那对应的潜在类型就是 T 自己。否则，其潜在类型就是在 类型声明 时 T 指定的那个类型的潜在类型。

type (
  A1 = string
  A2 = A1
)

type (
  B1 string
  B2 B1
  B3 []B1
  B4 B3
)
string , A1 , A2 , B1 和 B2 的潜在类型是 string 。 []B1 , B3 和 B4 的潜在类型是 []B1 。

方法集
一个类型可能有一个 方法集method set 与之关联。 接口类型 的方法集就是它的接口。任何其它类型 T 的方法集由以类型 T 为接收者所声明的所有 方法 组成。相应的 指针类型 *T 的方法集是以 *T 或 T 为接收者所声明的所有方法的集合（也就是说，它同样包含了 T 的方法集）。包含嵌入字段的应用于结构体的更多规则，会在 结构体类型 一节描述。任何其它类型会有一个空的方法集。在一个方法集中，每一个方法必须要有一个 唯一的 非 空白 的 方法名 。

类型的方法集确定了这个类型所 实现的接口 和以此类型作为 接收者 所可以 调用 的方法。

布尔类型
布尔类型 代表以预先声明的常量 true 和 false 所表示的布尔真值的集合。预先声明的布尔类型为 bool ，这是一个 定义类型 。

数字类型
数字类型 代表整数或浮点数值的集合。预先声明的结构独立的数字类型有:

uint8       无符号的  8 位整数集合（0 到 255）
uint16      无符号的 16 位整数集合（0 到 65535）
uint32      无符号的 32 位整数集合（0 到 4294967295）
uint64      无符号的 64 位整数集合（0 到 18446744073709551615）

int8        带符号的  8 位整数集合（-128 到 127）
int16       带符号的 16 位整数集合（-32768 到 32767）
int32       带符号的 32 位整数集合（-2147483648 到 2147483647）
int64       带符号的 64 位整数集合（-9223372036854775808 到 9223372036854775807）

float32     所有 IEEE-754 标准的 32 位浮点数数字集合
float64     所有 IEEE-754 标准的 64 位浮点数数字集合

complex64   由 float32 类型的实数和虚数部分所组成的所有复数的集合
complex128  由 float64 类型的实数和虚数部分所组成的所有复数的集合

byte        unit8 的别名
rune        int32 的别名
一个 n 位整数的值是 n 位宽的，是使用 补码 来表示的。

以下是根据实现不同而有特定大小的预先声明的数字类型:

uint     可以是 32 或 64 位
int      和 uint 大小相同
uintptr  一个大到足够用来存储一个指针值的未解释的比特位的无符号整数
为了避免移植性问题，除了 byte （ unit8 的别名）和 rune （ int32 的别名）外的所有数字类型都是截然不同的 定义类型 。当不同的数字类型混合在一个表达式或赋值里时，是需要显示的转换的。比如， int32 和 int 并不是相同的类型，就算在一个特定的架构上它们可能有相同的大小，也是如此。

字符串类型
字符串类型 代表了字符串值的集合。一个字符串值是字节的序列（可能为空）。字节的个数被称为该字符串的长度，并且不能为负。字符串是不可变的：一旦创建好了是不可能去修改其内容的。预先声明的字符串类型是 string ；它是一个 定义类型 。

字符串 s 的长度可以使用内置函数 len 来发现。如果字符串是一个常量，那么长度是一个编译时常量。一个字符串的字节可以通过从 0 索引 到 len(s) - 1 的整数来访问。获取这样子的一个元素的地址是非法的；如果 s[i] 是一个字符串的第 i 个字节，那么 &s[i] 是无效的。

数组类型Array types
数组是单一类型元素的有序序列，该单一类型称为元素类型。元素的个数被称为数组长度，并且不能为负值。

ArrayType   = "[", ArrayLength, "]", ElementType .
ArrayLength = Expression .
ElementType = Type .
长度是数组类型的一部分；它必须为一个可以被 int 类型的值所代表的非负 常量 。数组的长度 a 可以使用内置函数 len 来发现。元素可以被从 0 索引 到 len(a) - 1 的整数所寻址到。数组类型总是一维的，但可以被组合以形成多维类型。

[32]byte
[2*N] struct { x, y int32 }
[1000]*float64
[3][5]int
[2][2][2]float64  // 同 [2]([2]([2]float64))
分片类型Slice types
分片是针对一个底层数组的连续段的描述符，它提供了对该数组内有序序列元素的访问。分片类型表示其元素类型的数组的所有分片的集合。元素的数量被称为分片长度，且不能为负。未初始化的分片的值为 nil 。

SliceType = "[", "]", ElementType .
分片 s 的长度可以被内置函数 len 来发现；和数组不同的是，这个长度可能会在执行过程中改变。元素可以被从 0 索引 到 len(s) - 1 的整数所寻址到。一个给定元素的分片索引可能比其底层数组的相同元素的索引要小。

分片一旦初始化便始终关联到存放其元素的底层数组。因此分片会与其数组和其它相同数组的分片共享存储区；相比之下，不同的数组总是代表不同的存储区域。

分片底层的数组可以延伸超过分片的末端。 容量 便是对这个范围的测量：它是分片长度和数组内除了该分片以外的长度的和；不大于其容量长度的分片可以从原始分片 再分片 新的来创建。分片 a 的容量可以使用内置函数 cap(a) 来找到。

对于给定元素类型 T 的新的初始化好的分片值的创建是使用的内置函数 make ，它需要获取分片类型、指定的长度和可选的容量作为参数。使用 make 创建的分片总是分配一个新的隐藏的数组给返回的分片值去引用。也就是，执行

make([]T, length, capacity)
就像分配个数组然后 再分片 它一样来产生相同的分片，所以如下两个表达式是相等的:

make([]int, 50, 100)
new([100]int)[0:50]
如同数组一样，分片总是一维的但可以通过组合来构造高维的对象。数组间组合时，被构造的内部数组总是拥有相同的长度；但分片与分片（或数组与分片）组合时，内部的长度可能是动态变化的。此外，内部分片必须单独初始化。

结构体类型Struct types
结构体是命名元素的一个序列，这些元素被称为字段，每一个都有一个名字和一个类型。字段名可以被显式指定（IdentifierList）也可以被隐式指定（EmbeddedField）。在结构体中，非 空白 字段名必须是 唯一的 。

StructType    = "struct", "{", { FieldDecl, ";" }, "}" .
FieldDecl     = (IdentifierList Type | EmbeddedField), [ Tag ] .
EmbeddedField = [ "*" ], TypeName .
Tag           = string_lit .
// 一个空的结构体
struct {}

// 一个有六个字段的结构体
struct {
  x, y int
  u float32
  _ float32  // padding
  A *[]int
  F func()
}
一个声明了类型但没有显式的字段名的字段就是 嵌入字段 。嵌入字段必须指定为一个类型名 T 或者为一个到非接口类型的指针名 *T ， 并且 T 不是一个指针类型。这个非限定的类型名就被当作字段名。

// 四个类型分别为 T1, *T2, P.T3, *P.T4 的嵌入字段所组成的结构体
struct {
  T1        // 字段名为 T1
  *T2       // 字段名为 T2
  P.T3      // 字段名为 T3
  *P.T4     // 字段名为 T4
  x, y int  // 字段名为 x 和 y
}
以下声明是非法的，因为在一个结构体类型中，字段名必须是唯一的：

struct {
  T     // 与嵌入字段 *T 和 *P.T 冲突
  *T    // 与嵌入字段  T 和 *P.T 冲突
  *P.T  // 与嵌入字段  T 和   *T 冲突
}
在结构体 x 中，一个嵌入字段的字段或 方法 f 被称为是 promoted ，前提是 x.f 是一个表示那个字段或方法 f 的合法 选择器 。

除了不能在结构体的 复合字面值 中作为字段名外， promoted 字段和结构体的普通字段一样。

给定一个结构体类型 S 和一个 定义类型 T ， promoted 方法包含在这个结构体的方法集中的情况分为：

如果 S 包含一个嵌入字段 T ，那么 S 和 *S 的 方法集 都包括了接收者为 T 的 promoted 方法。 *S 的方法集还包括了接收者为 *T 的 promoted 方法。
如果 S 包含了一个嵌入字段 *T ，那么 S 和 *S 的 方法集 都包括了接收者为 T 或 *T 的 promoted 方法。
字段声明可以紧跟着一个可选的字符串字面值 标签 ，在对应的字段声明中，它将成为针对所有这个字段的属性。空的标签字符串等于没有标签。标签可以通过 反射接口 被可视化，并且可以参与到结构体的 类型一致性 中，但其它情况下都是被忽略的。

struct {
  x, y float64 ""  // 空的标签字面值和没有标签一样
  name string  "any string is permitted as a tag"
  _    [4]byte "ceci n'est pas un champ de structure"
}

// 对应时间戳协议缓冲区的结构体
// 其标签字符串定义了协议缓冲区的字段号
// 它们遵循了由 reflect 包所概述的转换规则
struct {
  microsec  uint64 `protobuf:"1"`
  serverIP6 uint64 `protobuf:"2"`
}
指针类型Pointer types
指针类型表示指向一给定类型的 变量 的所有指针的集合，这个给定类型称为该指针的 基础类型 。未初始化的指针的值为 nil 。

PointerType = "*", BaseType .
BaseType    = Type .
*Point
*[4]int
函数类型Function types
函数类型表示具有相同参数和结果类型的所有函数的集合。函数类型的未初始化的变量的值为 nil 。

FunctionType   = "func", Signature .
Signature      = Parameters, [ Result ] .
Result         = Parameters | Type .
Parameters     = "(", [ ParameterList, [ "," ] ], ")" .
ParameterList  = ParameterDecl, { ",", ParameterDecl } .
ParameterDecl  = [ IdentifierList ], [ "..." ], Type .
在参数或结果的列表中，名字（IdentifierList）要么全部存在，要么全部不存在。如果存在，每个名字代表特定类型的一个条目（参数或者结果），签名中的名字是非 空白 的，且必须是 唯一的 。如果不存在，每个类型代表该类型的一个条目。参数和结果列表总是括起来的，除非只有一个未命名的结果（可以写为不使用括号括起来的类型）。

函数签名中最后的进入参数可以是以 ... 为前缀的类型。带这样一个参数的函数被称为 variadic （可变），它可以携带针对该形参的零或多个实参来调用。

func()
func(x int) int
func(a, _ int, z float32) bool
func(a, b int, z float32) (bool)
func(prefix string, values ...int)
func(a, b int, z float64, opt ...interface{}) (success bool)
func(int, int, float64) (float64, *[]int)
func(n int) func(p *T)
接口类型Interface types
接口类型指定了一个称为 接口 的 方法集 。一个接口变量可以存储任意类型的值，这个类型要带有一个方法集，方法集需要是该接口的任意超集。这样子的类型就被叫做 实现了这个接口 。接口类型的未初始化的变量的值为 nil 。

InterfaceType      = "interface", "{", { MethodSpec, ";" }, "}" .
MethodSpec         = MethodName, Signature | InterfaceTypeName .
MethodName         = identifier .
InterfaceTypeName  = TypeName .
正如所有方法集一样，在接口类型内，每个方法必须有一个 唯一的 非 空白 名称。

// 一个简单的 File 接口。
interface {
  Read([]byte) (int, error)
  Write([]byte) (int, error)
  Close() error
}
interface {
  String() string
  String() string  // 非法: String 不是唯一的
  _(x int)         // 非法: 方法不能是空白名
}
多个类型可以实现一个相同的接口。比如，如果两个类型 S1 和 S2 有方法集

func (p T) Read(p []byte) (n int, err error)   { return … }
func (p T) Write(p []byte) (n int, err error)  { return … }
func (p T) Close() error                       { return … }
（其中 T 代表 S1 或 S2 ）那么 File 接口就被 S1 和 S2 实现了，不考虑 S1 和 S2 是否有其它的（或共享的）方法。

一个类型实现了包括其方法的子集的任意接口，因此可能实现了好几个截然不同的接口。比如，所有的类型都实现了 空 接口：

interface{}
类似的，来看在 类型声明 中用来定义一个叫做 Locker 的接口的规范：

type Locker interface {
  Lock()
  Unlock()
}
如果 S1 和 S2 也实现了

func (p T) Lock() { … }
func (p T) Unlock() { … }
和 File 接口一样，它们也实现了 Locker 接口。

一个接口 T 可以使用（可能是限定的）接口类型名 E 代替方法规范。这叫做在 T 中的 内嵌 接口 E ；它添加所有 E 的（暴露或者非暴露的）方法到接口 T 。

type ReadWriter interface {
  Read(b Buffer) bool
  Write(b Buffer) bool
}

type File interface {
  ReadWriter  // 和添加 ReadWriter 的方法一样
  Locker      // 和添加 Locker 的方法一样
  Close()
}

type LockedFile interface {
  Locker
  File        // 非法: Lock, Unlock 不唯一
  Lock()      // 非法: Lock 不唯一
}
一个接口类型 T 不能递归地嵌入它自己或者其它已经嵌入了 T 的接口类型。

// 非法: Bad 不能嵌入它自己
type Bad interface {
  Bad
}

// 非法: Bad1 不能通过 Bad2 来嵌入它自己
type Bad1 interface {
  Bad2
}
type Bad2 interface {
  Bad1
}
映射类型Map types
映射是单一类型元素所组成的无序组，这个单一类型被称为元素类型。元素由另一个类型的 键 的集合来索引，这个另一个类型被称为键类型。一个未初始化的映射的值为 nil 。

MapType     = "map", "[", KeyType, "]", ElementType .
KeyType     = Type .
比较运算符 == 和 != 对键类型操作而言必须是要完全定义的；因此键类型不能为一个函数、映射或分片。如果键类型是一个接口类型，那么比较运算符必须针对其动态键值做完全定义；失败会导致一个 run-time panic 。

map[string]int
map[*T]struct{ x, y float64 }
map[string]interface{}
映射元素的数目被称为其长度。对于一个映射 m ，长度可以使用内置函数 len 来找到并且可能会在执行过程中改变。元素可以在执行过程中使用 赋值 来进行添加，可以使用 索引表达式 来获取；可以使用内置函数 delete 来移除。

一个新的、空的映射值的创建使用的是内置函数 make ，其获取映射类型和一个可选的容量提示作为实参：

make(map[string]int)
make(map[string]int, 100)
初始化的容量不会限制其大小：映射会增长以适合其存储项目的数量，除了 nil 映射。 nil 映射相当于空映射，但是 nil 映射不能添加元素。

信道类型Channel types
信道针对 并发执行函数 提供了一个 发送 和 接收 特定类型值的机制。未初始化的信道的值为 nil 。

ChannelType = ( "chan" | "chan", "<-" | "<-", "chan" ), ElementType .
可选的 <- 运算符指定了信道的 方向 、 发送 或 接收 。如果没有指定方向，这个信道就是 双向的 。通过 赋值 或显示的 转换 ，信道可以被限制为仅能发送或仅能接收。

chan T          // 可用于发送或接收类型为 T 的值
chan<- float64  // 仅用于发送 float64 类型
<-chan int      // 仅用于接收 int 类型
<- 与最左的 chan 关联的一些可能性：

chan<- chan int    // 和 chan<- (chan int) 一样
chan<- <-chan int  // 和 chan<- (<-chan int) 一样
<-chan <-chan int  // 和 <-chan (<-chan int) 一样
chan (<-chan int)
一个新的，初始化的信道值的创建可以使用内置的函数 make ，它获取信道类型和可选的 容量 作为实参：

make(chan int, 100)
容量（元素的数量）确定了信道中缓冲区的大小。如果容量为零或没有写，那么信道就是无缓冲的，这种情况下，只有在接收端和发送端都准备好的情况下，通信才会成功。不然信道就是有缓冲的，这种情况下只要不阻塞，通信便会成功；阻塞是指缓冲区满了（对于发送端而言）或者缓冲区空了（对于接收端而言）。 一个 nil 的信道是不能用于通信的。

信道可以使用内置函数 close 来关闭。 接收运算符 的多值分配形式报告了在信道关闭前接收到的值是否已经被发送了。

单个信道可以被不需要进一步同步的任意数量的 goroutines 用在 发送语句 ， 接收运算符 和对内置函数 cap 及 len 的调用上。信道是一个先进先出的队列。举例，如果一个 goroutine 在信道上发送了值，第二个 goroutine 接收了这些值，那么这些值是按照发送的顺序被接收的。

类型和值的属性
类型一致性
两个类型，要么是 一致的 要么是 不同的 。

定义类型 和其它类型总是不同的。不然的话，如果两个类型所对应的 潜在类型 字面值是结构一致的——也就是说它们拥有相同的字面值结构并且对应的组成部分拥有一致的类型——那么它们便是一致的。详细来说：

如果两个数组类型有一致的元素类型和相同的数组长度，那么它们便是一致的。
如果两个分片类型有一致的元素类型，那么它们便是一致的。
如果两个结构体有相同的字段序列，并且对应的字段有相同的名字、一致的类型和一致的标签，那么它们便是一致的。（不同包的 非暴露的 字段名总是不同的）
如果两个指针类型有一致的基础类型，那么它们便是一致的。
如果两个函数类型有相同的参数数量和结果值，并且对应的参数和结果类型是一致的，并且两者要么都是 variadic 要么都不是，那么它们便是一致的。（参数和结果名不是必须匹配的）
如果两个接口类型有一样的带相同名字和一致的函数类型的方法集，那么它们便是一致的。（不同包的 非暴露的 方法名总是不同的。方法的顺序是无关紧要的）
如果两个映射类型有一致的键类型和值类型，那么它们便是一致的。
如果两个信道类型有一致的值类型和相同的方向，那么它们便是一致的。
给出声明

type (
  A0 = []string
  A1 = A0
  A2 = struct{ a, b int }
  A3 = int
  A4 = func(A3, float64) *A0
  A5 = func(x int, _ float64) *[]string
)

type (
  B0 A0
  B1 []string
  B2 struct{ a, b int }
  B3 struct{ a, c int }
  B4 func(int, float64) *B0
  B5 func(x int, y float64) *A1
)

type  C0 = B0
这些类型是一致的

A0, A1, and []string
A2 and struct{ a, b int }
A3 and int
A4, func(int, float64) *[]string, and A5

B0 and C0
[]int and []int
struct{ a, b *T5 } and struct{ a, b *T5 }
func(x int, y float64) *[]string, func(int, float64) (result *[]string), and A5
B0 和 B1 是不同的，因为它们是被不同的 类型定义 所创建的新类型； func(int, float64) *B0 和 func(x int, y float64) *[]string 是不同的，因为 B0 和 []string 是不同的。

可分配性
在如下这些情况中，值 x 可以分配 给一个类型为 T 的 变量 （「 x 可以分配给 T 」）：

x 的类型和 T 一致。
x 的类型 V 和 T 有一致的 潜在类型 并且二者最少有一个不是 定义类型 。
T 是一个接口类型，而 x 实现了 T 。
x 是一个双向的信道值， T 是一个信道类型， x 的类型 V 和 T 有一致的元素值，并且 V 和 T 中至少有一个不是定义类型。
x 是一个预先声明的标识符 nil 而 T 是一个指针、函数、分片、映射、信道或接口类型。
x 是一个无类型的可以被类型 T 的一个值所代表的 常量 。
可表示性
只要以下条件有一个成立，那么 常量 x 就可以被一个类型为 T 的值所表示：

x 在由 T 所确定的 值集中
T 是一个浮点类型并且 x 可以被不溢出地约到 T 的精度。约数用的是 IEEE 754 round-to-even 规则，但是 IEEE 负零会被进一步简化到一个无符号的零。（注：这种常量值不会出现 IEEE 负零、 NaN 或者无穷。）
T 是一个复合类型并且 x 的 组成 real(x) 和 imag(x) 是可以被 T 的组成类型（ float32 或者 float64 ）所表示的。
x                   T           x 可以被 T 表示的原因是

'a'                 byte        97 在 byte 值集中
97                  rune        rune 是 int32 的别名且 97 在 32 位整数值集中
"foo"               string      "foo" 在 string 值集中
1024                int16       1024 在 16 位整数值集中
42.0                byte        42 在无符号 8 位整数值集中
1e10                uint64      10000000000 在无符号 64 位整数值集中
2.718281828459045   float32     2.718281828459045 约到 2.7182817 后在 float32 值集中
-1e-1000            float64     -1e-1000 约到 IEEE -0.0 后再被进一步简化到 0.0
0i                  int         0 是一个整数值
(42 + 0i)           float32     42.0 （带虚部零）在 float32 值集中
x                   T           x 不能被 T 表示的原因是

0                   bool        0 不在 boolean 值集中
'a'                 string      'a' 是 rune，它不在 string 值集中
1024                byte        1024 不在无符号 8 位整数值集中
-1                  uint16      -1 不在无符号 16 位整数值集中
1.1                 int         1.1 不是一个整数值
42i                 float32     (0 + 42i) 不在 float32 值集中
1e1000              float64     1e1000 约数后溢出了 IEEE +Inf
块Blocks
块 是在一对花括号内的声明和语句序列，这个序列可能是空的。

Block = "{", StatementList, "}" .
StatementList = { Statement, ";" } .
源代码中除了显式的块外，还有隐式的块：

包围所有 Go 原始文本的 宇宙块 。
每个 包 有一个包含针对该包的所有 Go 原始文本的 包块 。
每个文件有一个包含在该文件中所有 Go 原始文本的 文件块 。
每个 "if" , "for" 和 "switch" 语句都被认为是在其自己的隐式块中。
每个在 "switch" 或 "select" 语句中的子句都作为一个隐式的块。
块是嵌套的并影响着 作用域 。

声明和作用域
声明 绑定了非 空白 的标识符到 常量 、 类型 、 变量 、 函数 、 标签 或 包 。程序中的每个标识符都必须要声明。同一个块中不能定义一个标识符两次，并且没有标识符可以同时在文件块和包块中定义。

空白标识符 可以像其它标识符一样在声明中使用，但它不会引出一个绑定，因此不被声明。在包块中，标识符 init 只能用于 init 函数 声明，且和空白标识符一样，它不会引出一个新的绑定。

Declaration   = ConstDecl | TypeDecl | VarDecl .
TopLevelDecl  = Declaration | FunctionDecl | MethodDecl .
声明的标识符的 作用域 是该标识符表示特定常量、类型、变量、函数、标记或包时所处的原始文本的范围。

Go 使用 块 来定作用域：

预先声明的标识符 的作用域为宇宙块。
表示一个常量、类型、变量或函数（但不是方法）的在最上层（在任何函数外）定义的标识符的作用域为包块。
导入的包的包名的作用域为包含导入声明在内的文件的文件块。
表示一个方法接收者、函数参数或结果变量的标识符的作用域为函数主体。
在函数内定义的常量或变量标识符的作用域起始于 ConstSpec 或 VarSpec（对短变量来说为 ShortVarDecl）的尾端，结束于包含着它的最内的块。
在函数内定义的类型标识符的作用域起始于 TypeSpec 的标识符，结束于包含着它的最内的块。
在块中声明的标识符可以在其内的块中重新声明。当内部声明的标识符在作用域内时，它表示内部声明所声明的实体。

包子句 不是一个声明；包名不会在任何作用域中出现。它的目的是确定一个文件属于相同的 包 和针对导入声明指定默认的包名。

标签作用域
标签是由 标签语句 所声明的，它用在 "break" 、 "continue" 和 "goto" 语句中。定义一个不去用的标签是非法的。与其它标识符相对比，标签不按块分作用域，也不和那些不是标签的标识符冲突。标记的作用域是声明时所在的函数的主体，不过要排除所有嵌套函数的主体。

空白标识符
空白标识符 由下划线字符 _ 所代表。它充当一个匿名的占位符替代通常的（非空白的）标识符，并且作为 操作数 在 声明 和 赋值 中有特殊的意义。

预声明的标识符
以下的标识符是在 宇宙块 中被隐式地定义的:

Types:
  bool byte complex64 complex128 error float32 float64
  int int8 int16 int32 int64 rune string
  uint uint8 uint16 uint32 uint64 uintptr

Constants:
  true false iota

Zero value:
  nil

Functions:
  append cap close complex copy delete imag len
  make new panic print println real recover
暴露的标识符
标识符可以被 暴露 用来允许从另一个包访问到它。一个标识符将会被暴露如果同时满足：

标识符的首字母为 Unicode 大写字母（Unicode 类 "Lu"）；以及
标识符是在 包块 中声明的或者它是一个 字段名 或 方法名 。
所有其它的标识符是不暴露的。

标识符的唯一性
给定一个标识符集，如果一个标识符与在该集合中的所有其它都 不同 ，那么其便被称为是 唯一的 。如果两个标识符拼写不同，或它们处于不同的 包 并且没有被暴露，那么它们便是不同的。否则，它们便是相同的。

常量声明
常量声明绑定了一个标识符的列表（常量的名字）到 常量表达式 列表的值。标识符的数量必须等于表达式的数量，并且左侧第 n 个标识符绑定到了右侧第 n 个表达式的值。

ConstDecl      = "const", ( ConstSpec | "(", { ConstSpec, ";" }, ")", ) .
ConstSpec      = IdentifierList, [ [ Type ], "=", ExpressionList ] .

IdentifierList = identifier { ",", identifier } .
ExpressionList = Expression { ",", Expression } .
如果类型提供了，那么所有常量需采用该指定类型，并且表达式必须 可分配 到该类型。如果类型省略了，常量为对应表达式的独立的类型。如果表达式的值为无类型 常量 ，那么声明的常量保持为无类型，常量标识符表示着该常量的值。比如，如果一个表达式为浮点数字面值，那么即使字面值的小数部分为零，常量标识符依旧表示一个浮点数常量。

const Pi float64 = 3.14159265358979323846
const zero = 0.0        // 无类型的浮点数常量
const (
  size int64 = 1024
  eof        = -1       // 无类型的整数常量
)
const a, b, c = 3, 4, "foo"  // a = 3, b = 4, c = "foo", 无类型的整数和字符串常量
const u, v float32 = 0, 3    // u = 0.0, v = 3.0
在括起来的 const 声明列表中，除了第一个常量声明外，其它的表达式列可以省略。这样的一个空列表相当于第一个前面的非空表达式列表及其类型（如果有的话）的文本替换。省略表达式的列表就因此相当于重复之前的列表。标识符的数量必须等于之前列表的表达式的数量。这个机制结合 iota 常量生成器允许了连续值的轻量声明：

const (
  Sunday = iota
  Monday
  Tuesday
  Wednesday
  Thursday
  Friday
  Partyday
  numberOfDays  // 这个常量是不暴露的
)
Iota
在一个 常量声明 中，预先声明的标识符 iota 代表连续的无类型的整数 常量 。它的值从零开始，是在常量声明中各自的 ConstSpec 的索引。其可以用于构造一组相关常量的集合：

const (
  c0 = iota  // c0 == 0
  c1 = iota  // c1 == 1
  c2 = iota  // c2 == 2
)

const (
  a = 1 << iota  // a == 1  (iota == 0)
  b = 1 << iota  // b == 2  (iota == 1)
  c = 3          // c == 3  (iota == 2，没有使用)
  d = 1 << iota  // d == 8  (iota == 3)
)

const (
  u         = iota * 42  // u == 0     （无类型的整数常量）
  v float64 = iota * 42  // v == 42.0  （float64 常量）
  w         = iota * 42  // w == 84    （无类型的整数常量）
)

const x = iota  // x == 0
const y = iota  // y == 0
定义上，在同一个 ConstSpec 中使用的多个 iota 都拥有相同的值：

const (
  bit0, mask0 = 1 << iota, 1<<iota - 1  // bit0 == 1, mask0 == 0  (iota == 0)
  bit1, mask1                           // bit1 == 2, mask1 == 1  (iota == 1)
  _, _                                  //                        (iota == 2，没有使用)
  bit3, mask3                           // bit3 == 8, mask3 == 7  (iota == 3)
)
最后一个例子利用了上一个非空表达式列表的 隐式重复 。

类型声明Type declarations
一个类型声明绑定了一个标识符（也就是 类型名 ）到一个 类型 。类型声明有两种形式：别名声明和类型定义。

TypeDecl     = "type" ( TypeSpec | "(", { TypeSpec, ";" }, ")" ) .
TypeSpec     = AliasDecl | TypeDef .
别名声明Alias declarations
别名声明绑定了一个标识符到一个给定的类型。

AliasDecl = identifier, "=", Type .
在标识符的 作用域 内，它充当了该类型的 别名 。

type (
  nodeList = []*Node  // nodeList 和 []*Node 的类型一致
  Polar    = polar    // Polar 和 polar 表示的类型一致
)
类型定义Type definitions
类型定义创建一个新的，不同的类型，其具有与给定类型相同的 潜在类型 和操作，并将标识符绑定到它。

TypeDef = identifier, Type .
新类型被称为 定义类型 。它和其它任何的类型（包括那个给定类型）都是 不同的 。

type (
  Point struct{ x, y float64 }  // Point 和 struct{x, y float64} 是不同的类型
  polar Point                   // polar 和 Point 表示不同的类型
)

type TreeNode struct {
  left, right *TreeNode
  value *Comparable
}

type Block interface {
  BlockSize() int
  Encrypt(src, dst []byte)
  Decrypt(src, dst []byte)
}
定义类型可能具有与之关联的 方法 。它不会继承任何绑定到给定类型的方法，但接口类型或者复合类型元素的 方法集 是保持不变的：

// Mutex 是带两个方法——Lock 和 Unlock——的数据类型。
type Mutex struct         { /* 互斥对象字段 */ }
func (m *Mutex) Lock()    { /* Lock 实现 */ }
func (m *Mutex) Unlock()  { /* Unlock 实现 */ }

// NewMutex 和 Mutex 有相同的构成，但是其方法集是空的。
type NewMutex Mutex

// PtrMutex 的潜在类型 *Mutex 的方法集是保持不变的，
// 但 PtrMutex 的方法集是空的。
type PtrMutex *Mutex

// *PrintableMutex 的方法集包含了绑定到它的嵌入字段 Mutex 的方法 Lock 和 Unlock 。
type PrintableMutex struct {
  Mutex
}

// MyBlock 是一个和 Block 有着相同方法集的接口类型。
type MyBlock Block
类型声明可以用于定义不同的布尔、数值或字符串类型，并关联方法给它：

type TimeZone int

const (
  EST TimeZone = -(5 + iota)
  CST
  MST
  PST
)

func (tz TimeZone) String() string {
  return fmt.Sprintf("GMT%+dh", tz)
}
变量声明Variable declarations
一个变量声明创建一个或多个变量，给它们绑定对应的标识符，并且给每个分一个类型和一个初始化的值。

VarDecl     = "var", ( VarSpec | "(", { VarSpec, ";" }, ")", ) .
VarSpec     = IdentifierList, ( Type, [ "=", ExpressionList ] | "=", ExpressionList ) .
var i int
var U, V, W float64
var k = 0
var x, y float32 = -1, -2
var (
  i int
  u, v, s = 2.0, 3.0, "bar"
)
var re, im = complexSqrt(-1)
var _, found = entries[name]  // 映射查找；只关心 "found"
如果给出了表达式列表，那么变量会根据 赋值 规则由表达式来初始化。否则，每个变量都被初始化为其 零值 。

如果类型提供了，那么每个变量都会指定为那个类型。否则，每个变量的类型会被给定为赋值中对应的初始化值的类型。如果那个值是无类型的常量，它会先隐式地 转换 为它的 默认类型 ；如果它是一个无类型的布尔值，那么它会先隐式地转换为类型 bool 。预先声明的值 nil 不能用于初始化没有明确类型的变量。

var d = math.Sin(0.5)  // d 是 float64
var i = 42             // i 是 int
var t, ok = x.(T)      // t 是 T, ok 是 bool
var n = nil            // 非法
实现限制：当在 函数实体 中定义的变量没有被使用时，编译器可以认定它为非法的。

短变量声明
短变量声明 使用如下语法:

ShortVarDecl = IdentifierList, ":=", ExpressionList .
这是如下这种带初始化表达式而不带类型的 变量声明 的速记法:

"var", IdentifierList, "=", ExpressionList .
i, j := 0, 10
f := func() int { return 7 }
ch := make(chan int)
r, w, _ := os.Pipe(fd)  // os.Pipe() 返回一个连接着的文件对和一个 error （如果有的话）
_, y, _ := coord(p)  // coord() 返回三个值; 只关心 y 座标
和普通的变量声明不同，短变量声明可以 重复声明 一个变量，这个变量是在同一个块（或者参数列表——如果该块是一个函数实体的话）内之前已经声明过的，且变量类型不能改变，但是重复声明语句最少要存在一个新的非 空白 变量。因此，重复声明仅能出现在多变量短声明中。重复声明不会引进新的变量；它仅赋一个新的值到原变量。


field1, offset := nextField(str, 0)
field2, offset := nextField(str, offset)  // 重复声明了 offset
a, a := 1, 2                              // 非法: a 声明了两次，或者如果 a 已经在其它地方声明的了话那么就没有新的变量了
短变量声明只能出现在函数内。在一些针对诸如 "if" 、 "for" 或 "switch" 这样的初始化器的上下文中，也可以用于声明本地临时变量。

函数声明Function declarations
函数声明绑定一个标识符（也就是 函数名 ）到一个函数。

FunctionDecl = "func", FunctionName, Signature, [ FunctionBody ] .
FunctionName = identifier .
FunctionBody = Block .
如果函数的 签名 声明了结果参数，那么函数体语句列表必须以 终止语句 结尾。

func IndexRune(s string, r rune) int {
  for i, c := range s {
    if c == r {
      return i
    }
  }
  // 无效: 缺少返回语句
}
一个函数声明可以缺少函数体。这样子的声明为 Go 语言外的所实现的函数提供了签名，比如一个汇编程序。

func min(x int, y int) int {
  if x < y {
    return x
  }
  return y
}

func flushICache(begin, end uintptr)  // 由外部实现
方法声明Method declarations
方法是带 接收者 的 函数 。一个方法声明绑定了一个标识符（也就是 方法名 ）为一个方法，并与接收者的 基础类型 关联。

MethodDecl   = "func", Receiver, MethodName, Signature, [ FunctionBody ] .
Receiver     = Parameters .
接收者是使用在方法名之前的额外的参数段来指定的。这个参数段必须声明一个单一非 variadic 参数作为接收者。其类型必须为 定义类型 T 或到定义类型 T 的指针。 T 被称为接收者的 基础类型 。接收者的基本类型不能是一个指针或者接口类型，并且它必须在和方法相同的包中被声明。这个方法就被称为 绑定到了 这个基础类型，方法名只能通过类型 T 或 *T 的 选择器 才可见。

译注：方法的基础类型不能是接口，这边不要混淆，接口是一组方法签名的集合，也就是可以定义一个固定类型为一个接口类型，这个固定类型实现了对应接口类型所声明的方法。
一个非 空白 接收者标识符在方法签名中必须是 唯一的 。如果接收者的值在方法实体内没有被引用，那么其标识符在声明时是可以省略的。一般来说这也同样适用于函数和方法的参数。

对一个基础类型来说，绑定到它的非空白的方法名必须是唯一的。如果基础类型为 结构体类型 。那么非空白的方法和字段名必须是不同的。

给定一个定义类型 Point ，其声明

func (p *Point) Length() float64 {
  return math.Sqrt(p.x * p.x + p.y * p.y)
}

func (p *Point) Scale(factor float64) {
  p.x *= factor
  p.y *= factor
}
绑定了方法 Length 和 Scale ，接收者类型为 *Point ，对应基础类型 Point 。

方法的类型是该函数结合接收者作为的第一个参数的类型。比如，方法 Scale 有类型

func(p *Point, factor float64)
不过，这样子声明的函数并不是一个方法。

表达式
表达式将运算符和函数应用于操作数来规定值的计算。

操作数Operands
操作数表示表达式中基本的值。一个操作数可能是一个字面值；可能是一个（可能为 限定的 ）表示 常量 、 变量 或 函数 的非 空白 标识符或者一个圆括号括起来的表达式。

空白标识符 只有在 赋值 的左侧时才能作为一个操作数。

Operand     = Literal | OperandName | "(", Expression, ")" .
Literal     = BasicLit | CompositeLit | FunctionLit .
BasicLit    = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .
OperandName = identifier | QualifiedIdent.
限定标识符Qualified identifiers
限定标识符是由包名前缀所限定的标识符。包名和标识符都不能为 空白 。

QualifiedIdent = PackageName, ".", identifier .
限定标识符可以在不同的包内访问一个标识符，该标识符对应的包必须已经被 导入 。标识符则必须已经在那个包被 暴露 并在 包块 中被声明。

math.Sin      // 表示在包 math 中的 Sin 函数
复合字面值Composite literals
复合字面值为结构体、数组、分片和映射构造值，并在每次被求值时创建一个新的值。复合字面值由字面值类型和紧跟着的花括号绑定的元素列表所组成。每个元素可以选择前缀一个对应的键。

CompositeLit  = LiteralType, LiteralValue .
LiteralType   = StructType | ArrayType | "[", "...", "]", ElementType |
                SliceType | MapType | TypeName .
LiteralValue  = "{", [ ElementList, [ "," ] ], "}" .
ElementList   = KeyedElement, { ",", KeyedElement } .
KeyedElement  = [ Key, ":" ], Element .
Key           = FieldName | Expression | LiteralValue .
FieldName     = identifier .
Element       = Expression | LiteralValue .
LiteralType 的潜在类型必须是结构体、数组、分片或者映射类型（文法强制执行此约束，除非类型是作为 TypeName 给出的）。元素和键的类型必须是 可分配 到字面值类型所对应的字段、元素和键类型的；这里没有额外的转换。该键被解释为结构体字面值的字段名，数组和分片字面值的索引，映射字面值的键。对于映射字面值而言，索引元素必须要有一个键。给多个元素指定相同的字段名或者不变的键值会出错。查阅 求值顺序 一节获取非常量映射键的信息。

对结构体字面值来说，应用如下规则：

键必须为在结构体类型中声明的字段。
不包含任何键的元素列表必须对每个结构体字段（字段声明的顺序）列出一个元素。
只要一个元素有键，那么每个元素都必须要有键。
包含键的元素列表不需要针对每个结构体字段有一个元素。省略的字段会获得一个零值。
字面值可以省略元素列表；这样子的字面值相当于其类型的零值。
针对属于不同包的结构体的非暴露字段来指定一个元素是错误的。
给定一个声明

type Point3D struct { x, y, z float64 }
type Line struct { p, q Point3D }
你可以写

origin := Point3D{}                            // Point3D 为零值
line := Line{origin, Point3D{y: -4, z: 12.3\}\}  // line.q.x 为零值
对数组和分片字面值来说，应用如下规则：

数组中的每个元素有一个关联的标记其位置的整数索引。
带键的元素使用该键作为其索引。这个键必须是可被类型 int 所表示的一个非负常量；而且如果其被赋予了类型的话则必须是整数类型。
不带键的元素使用之前元素的索引加一。如果第一个元素没有键，则其索引为零。
一个复合变量的 寻址 生成了一个到由字面值值初始化的唯一 变量 的指针。

var pointer *Point3D = &Point3D{y: 1000}
注意的是，分片和映射类型的零值不同于同类型的初始化过但为空的值。所以，获取空的分片或映射复合字面值的地址与使用 new 来分配一个新的分片或映射的效果不同。

p1 := &[]int{}    // p1 指向一个初始化过的值为 []int{} 长度为 0 的空分片
p2 := new([]int)  // p2 指向一个值为 nil 长度为 0 的未初始化过的分片
数组字面值的长度是字面值类型所指定的长度。在字面值中，如果少于其长度的元素被提供了，那么缺漏的元素会被设置为数组元素类型的零值。提供其索引值超出了数组索引范围的元素是错误的。符号 ... 指定一个数组长度等于其最大元素索引加一。

buffer := [10]string{}             // len(buffer) == 10
intSet := [6]int{1, 2, 3, 5}       // len(intSet) == 6
days := [...]string{"Sat", "Sun"}  // len(days) == 2
分片字面值描述了整个底层数组字面值。因此一个分片字面值的长度和容量为其最大元素索引加一。分片字面值的格式为

[]T{x1, x2, … xn}
以及针对应用到数组的分片操作的速记为

tmp := [n]T{x1, x2, … xn}
tmp[0 : n]
在数组、分片或者映射类型 T 的复合字面值中，如果元素或映射的键本身为复合字面值，当其字面值类型和 T 的元素或键类型一致时，该字面值类型可以省略。类似的，如果元素或键本身为复合字面值的地址，当元素或键的类型为 *T 时，该元素或键可以省略 &T 。

这边要多看多熟悉
[...]Point\{\{1.5, -3.5}, {0, 0}\}     // 同 [...]Point{Point{1.5, -3.5}, Point{0, 0}\}
[][]int\{\{1, 2, 3}, {4, 5}\}          // 同 [][]int{[]int{1, 2, 3}, []int{4, 5}\}
[][]Point\{\{\{0, 1}, {1, 2}\}\}         // 同 [][]Point{[]Point{Point{0, 1}, Point{1, 2}\}\}
map[string]Point{"orig": {0, 0}\}    // 同 map[string]Point{"orig": Point{0, 0}\}
map[Point]string{\{0, 0}: "orig"}    // 同 map[Point]string{Point{0, 0}: "orig"}

type PPoint *Point
[2]*Point{\{1.5, -3.5}, {\}\}          // 同 [2]*Point{&Point{1.5, -3.5}, &Point{}\}
[2]PPoint{\{1.5, -3.5}, {\}\}          // 同 [2]PPoint{PPoint(&Point{1.5, -3.5}), PPoint(&Point{})}
当一个使用 LiteralType 的 TypeName 形式的复合字面值表现为一个在 关键字 和 "if" 、 "for" 或 "switch" 语句块的左花括号之间的操作数，并且该复合字面值不被圆括号、方括号或花括号所包围时，会出现一个解析歧义。在这样子一个罕见的情况下，复合字面值的左花括号错误地被解析为语句块的引入。为了解决这样子的歧义，这个复合字段必须在圆括号内。

if x == (T{a,b,c}[i]) { … }
if (x == T{a,b,c}[i]) { … }
有效的数组、分片和映射字面值的例子：

// 质数列表
primes := []int{2, 3, 5, 7, 9, 2147483647}

// 当 ch 为元音时 vowels[ch] 为真
vowels := [128]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'y': true}

// 数组 [10]float32{-1, 0, 0, 0, -0.1, -0.1, 0, 0, 0, -1}
filter := [10]float32{-1, 4: -0.1, -0.1, 9: -1}

// 十二平均律以 Hz 为单位的频率（A4 = 440Hz）
noteFrequency := map[string]float32{
  "C0": 16.35, "D0": 18.35, "E0": 20.60, "F0": 21.83,
  "G0": 24.50, "A0": 27.50, "B0": 30.87,
}
函数字面值Function literals
函数字面值代表一个匿名 函数 。

FunctionLit = "func", Signature, FunctionBody .
func(a, b int, z float64) bool { return a*b < int(z) }
函数字面值可以被赋给一个变量或者直接调用。


f := func(x, y int) int { return x + y }
func(ch chan int) { ch <- ACK }(replyChan)

**FLAG: (replayChan) https://stackoverflow.com/questions/16008604/why-add-after-closure-body-in-golang**
函数字面值是 闭包 ：它们可以引用外层函数定义的变量。然后这些变量就在外层函数和函数字面值间共享了，并且只要能被访问就可以一直存活。

主表达式Primary expressions
主表达式是一元表达式和二元表达式的操作数。


PrimaryExpr =
  Operand |
  Conversion |
      MethodExpr |
  PrimaryExpr, Selector |
  PrimaryExpr, Index |
  PrimaryExpr, Slice |
  PrimaryExpr, TypeAssertion |
  PrimaryExpr, Arguments .

Selector       = ".", identifier .
Index          = "[", Expression, "]" .
Slice          = "[", [ Expression ], ":", [ Expression ], "]" |
                 "[", [ Expression ], ":", Expression, ":", Expression, "]" .
TypeAssertion  = ".", "(", Type, ")" .
Arguments      = "(", [ ( ExpressionList | Type, [ ",", ExpressionList ] ), [ "..." ], [ "," ] ], ")" .
x
2
(s + ".txt")
f(3.1415, true)
Point{1, 2}
m["foo"]
s[i : j + 1]
obj.color
f.p[i].x()
选择器
针对一个不为 包名 的 主表达式 x ， 选择器表达式

x.f
表示了值 x （或者有时候为 *x ；见下文）的字段或方法 f 。标识符 f 被称为（字段或方法） 选择器 ，它一定不能为 空白标识符 。选择器表达式的类型为 f 的类型。如果 x 是一个包名，看 限定标识符 一节。

选择器 f 可以表示一个类型 T 的一个字段或方法 f ，或者可以指嵌套在 T 中的 嵌入字段 的字段或方法 f 。遍历以达到 f 所经历的嵌入字段数被称为 f 在 T 中的 深度 。在 T 中声明的字段或者方法 f 的深度为零。在 T 中的嵌入字段 A 中声明的字段或者方法 f 的深度为 A 中 f 的深度加一。

以下规则应用于选择器：

对于为类型 T 或 *T 的值 x （ T 既不是指针类型也不是接口类型）， x.f 表示在 T 中最浅深度的字段或者方法 f 。如果不是恰好 一个 f 在最浅深度的话，那么这个选择器表达式就是非法的。
对于为接口类型 I 的值 x ， x.f 表示动态值 x 的名为 f 的实际的方法。如果在 I 的 方法集 中没有名为 f 的方法，那么这个选择器表达式就是非法的。
作为例外，如果 x 的类型为一个 定义的 指针类型并且 (*x).f 是一个有效的表示一个字段（但不是方法）的选择器表达式，那么 x.f 是 (*x).f 的速记。
在所有其它情况中， x.f 是非法的。
如果 x 是指针类型并且值为 nil 并且 x.f 表示一个结构体字段，那么，给 x.y 赋值或求值会导致一个 run-time panic 。
如果 x 是接口类型并且值为 nil ，那么 调用 或 求值 方法 x.y 会导致一个 run-time panic 。
这边好好熟悉，一头雾水，规则 2 应该要结合方法声明/调用那节一起看
举例，给定声明：

type T0 struct {
  x int
}

func (*T0) M0()

type T1 struct {
  y int
}

func (T1) M1()

type T2 struct {
  z int
  T1
  *T0
}

func (*T2) M2()

type Q *T2

var t T2     // 假定 t.T0 != nil
var p *T2    // 假定 p != nil 并且 (*p).T0 != nil
var q Q = p
你可以写：

t.z          // t.z
t.y          // t.T1.y
t.x          // (*t.T0).x

p.z          // (*p).z
p.y          // (*p).T1.y
p.x          // (*(*p).T0).x

q.x          // (*(*q).T0).x        (*q).x 是一个有效的字段选择器

p.M0()       // ((*p).T0).M0()      M0 期望接收者 *T0
p.M1()       // ((*p).T1).M1()      M1 期望接收者 T1
p.M2()       // p.M2()              M2 期望接收者 *T2
t.M2()       // (&t).M2()           M2 期望接收者 *T2，查看调用一节
但下述是无效的：

q.M0()       // (*q).M0 是有效的，但不是字段选择器
方法表达式
如果 M 在类型 T 的 方法集 中，那么 T.M 是一个函数，该函数可以携带和 M 同样的实参像普通函数一样调用，不过会给其前缀一个额外的实参作为该方法的接收者。

MethodExpr    = ReceiverType, ".", MethodName .
ReceiverType  = Type .
考虑有两个方法的结构体类型 T ，方法一是接收者为类型 T 的 Mv ，其二是接收者为类型 *T 的 Mp 。

type T struct {
  a int
}
func (tv  T) Mv(a int) int         { return 0 }  // 值接收者
func (tp *T) Mp(f float32) float32 { return 1 }  // 指针接收者

var t T
表达式

T.Mv
产生一个等同于 Mv 但带一个明确的接收者作为其第一个实参的函数；它的签名为

func(tv T, a int) int
这个函数可以在带一个明确的接收者情况下被正常地调用，所以以下五种调用是等同的：

t.Mv(7)
T.Mv(t, 7)
(T).Mv(t, 7)
f1 := T.Mv; f1(t, 7)
f2 := (T).Mv; f2(t, 7)
类似的，表达式

(*T).Mp
产生一个签名为如下的代表 Mp 的函数值

func(tp *T, f float32) float32
对于一个带值接收者的方法，可以推导出一个带明确指针接收者的函数，所以

(*T).Mv
产生一个签名为如下的代表 Mv 的函数值

func(tv *T, a int) int
这样的一个函数通过接收者创建一个值间接地将其作为接收者传递给底层函数；这个方法在函数调用中不会覆盖那个地址被传递的值。

最后一种情况——值接收者的函数对指针接收者的方法——是非法的，因为指针接收者的方法不在该值类型的方法集中。

从方法推导出的函数值是用函数调用语法来调用的；接收者作为调用的第一个实参。也就是，给定 f := T.Mv ， f 是作为 f(t, 7) 而非 t.f(7) 被调用的。使用 函数字面值 或 方法值 来构建一个绑定了接收者的函数。

从一个接口类型的方法中得到一个函数值是合法的。所得到的函数使用该接口类型的显式的接收者（原文： The resulting function takes an explicit receiver of that interface type. ）。

方法值
如果表达式 x 有静态类型 T ，并且 M 在类型 T 的 方法集 中，那么 x.M 被称为一个 方法值 。方法值 x.M 是一个可以用与 x.M 的方法调用的相同的实参来调用函数值。表达式 x 在该方法值的求值过程中被求值和保存；保存的副本被用在任意调用中（可能会在后续被执行的）作为接收者。

类型 T 可以为接口或者非接口类型。

就像上面 方法表达式 所讨论的，考虑一个带两个方法的结构体 T ，方法一是接收者为类型 T 的 Mv ，其二是接收者为类型 *T 的 Mp 。

type T struct {
  a int
}
func (tv  T) Mv(a int) int         { return 0 }  // 值接收者
func (tp *T) Mp(f float32) float32 { return 1 }  // 指针接收者

var t T
var pt *T
func makeT() T
表达式

t.Mv
产生了一个类型如下的函数值

func(int) int
这两种调用是等同的：

t.Mv(7)
f := t.Mv; f(7)
类似的，表达式

pt.Mp
产生了一个类型如下的函数值

func(float32) float32
就 选择器 来说，如果以值作为接收者的非接口方法使用了指针来引用，那么会自动解除到该指针的引用： pt.Mv 等同于 (*pt).Mv 。

就 方法调用 来说，如果以指针作为接收者的非接口方法使用了可寻址值来引用，那么会自动获取该值的地址来引用： t.Mp 等同于 (&t).Mp 。

f := t.Mv; f(7)   // 就像 t.Mv(7)
f := pt.Mp; f(7)  // 就像 pt.Mp(7)
f := pt.Mv; f(7)  // 就像 (*pt).Mv(7)
f := t.Mp; f(7)   // 就像 (&t).Mp(7)
f := makeT().Mp   // 无效的: makeT() 的结果是不可寻址的
虽然以上的例子使用了非接口类型，但是从接口类型的值来创建一个方法值同样是合法的。

var i interface { M(int) } = myVal
f := i.M; f(7)  // 就像 i.M(7)
索引表达式
如下形式的主表达式

a[x]
表示了可被 x 索引的数组、到数组的指针、分片、字符串或者被 x 索引的映射 a 的元素。值 x 分别被称为 索引 或 映射键 。以下规则应用于：

如果 a 不是一个映射：

索引 x 必须是整数类型或者无类型常量
常量索引必须为非负且可以被类型 int 所表示的 的一个值
无类型的常量索引会被给定一个类型 int
当 0 <= x < len(a) 时，索引 x 在范围内 ，否则它就 超出了范围
对于为 数组类型 A 的 a ：

常量 索引必须在范围内
如果在运行时 x 超出了范围，那么会发生一个 run-time panic
a[x] 是一个在索引 x 处的数组元素，且 a[x] 的类型是 A 的元素类型
对于到数组类型的 指针 a ：

a[x] 是 (*a)[x] 的速记
对于为 分片类型 S 的 a ：

如果在运行时 x 超出了范围，那么会发生一个 run-time panic
a[x] 是在索引 x 处的分片元素，且 a[x] 的类型是 S 的元素类型
对于 字符串类型 a ：

当字符串 a 是常量时， 常量 索引必须在范围内
如果在运行时 x 超出了范围，那么会发生一个 run-time panic
a[x] 是在索引 x 处的非常量字节，并且 a[x] 的类型为 byte
a[x] 不能被赋值
对于为 映射类型 M 的 a ：

x 的类型必须是 可分配 为 M 的键类型的
如果映射带键为 x 的条目，那么 a[x] 是带键 x 的映射值，并且 a[x] 的类型为 M 的值类型。
如果映射为 nil 或者不存这样这样子的一个条目，那么 a[x] 是针对 M 的值类型的 零值 。
其它情况下 a[x] 是非法的。

在类型为 map[K]v 的映射 a 中的用在 赋值 或特殊格式的初始化中的索引表达式

v, ok = a[x]
v, ok := a[x]
var v, ok = a[x]
产生了一个额外的无类型的布尔值。当键 x 存在于映射中时， ok 的值为 true ，否则为 false 。

给 nil 映射的元素赋值会导致一个 run-time panic 。

分片表达式
分片表达式从一个字符串、数组、到数组的指针或者分片中构建一个子字符串或者一个分片。有两种变体：指定一个低位和高位边界的简单格式，以及同时在容量上有指定的完整格式。

简单的分片表达式
对于一个字符串、数组、到数组的指针或者分片 a ，主表达式

a[low : high]
构造了一个子字符串或者分片。 索引 low 和 high 选择了操作数 a 的哪些元素作为结果被显示。结果有从零开始且长度等于 high - low 的索引。在分片了数组 a 后

a := [5]int{1, 2, 3, 4, 5}
s := a[1:4]
分片 s 有类型 []int ，长度 3，容量 4，以及元素

s[0] == 2
s[1] == 3
s[2] == 4
为了方便，每一个索引都可能被省略。缺少的 low 索引默认为零；缺少的 high 索引默认为被分片的操作数的长度：

a[2:]  // 同 a[2 : len(a)]
a[:3]  // 同 a[0 : 3]
a[:]   // 同 a[0 : len(a)]
如果 a 为到数组的指针，那么 a[low : high] 为 (*a)[low : high] 的速记。

对于数组或者字符串，如果 0 <= low <= high <= len(a) ，那么索引是 在范围内 的，否则就 超出了范围 。对于分片，上索引边界是分片的容量 cap(a) 而不是其长度。 常量 索引必须为非负且是可以被类型 int 所表示的 ；对于数组和常量字符串而言，常量索引也必须在范围内。如果两个索引都是常量，那么它们必须满足 low <= high 。如果在运行时索引超出了范围，那么会发生 run-time panic 。

除了 无类型的字符串 以外，如果被分片的操作数是一个字符串或者分片，那么分片操作的结果为一个和该操作数具有相同类型的非常量值。对于无类型字符串操作数而言，其结果是一个类型为 string 的非常量值。如果被分片的操作数是一个数组，那么它必须是 可被寻址的 ，并且分片操作的结果为和该数组具有相同元素类型的分片。

如果一个有效的分片表达式的被分片的操作数是一个 nil 分片，那么结果是一个 nil 分片。否则，结果会共享该操作数的底层数组。

var a [10]int
s1 := a[3:7]   // s1 的底层数组是数组 a； &s1[2] == &a[5]
s2 := s1[1:4]  // s2 的底层数组是 s1 的底层数组 a； &s2[1] == &a[5]
s2[1] = 42     // s2[1] == s1[2] == a[5] == 42；它们指的都是相同的底层数组元素
完整的分片表达式
对于数组、到数组的指针或者分片 a （不是一个字符串），主表达式

a[low : high : max]
构成了一个有相同类型的分片，并且带有和简单的分片表达式 a[low : high] 一样的长度和元素。此外，它通过设置分片的容量为 max - low 来控制产生的分片的容量。只有第一个索引是可以被省略的；默认为零。在分片了数组 a 后

a := [5]int{1, 2, 3, 4, 5}
t := a[1:3:5]
分片 t 有类型 []int ，长度 2，容量 4，以及元素

t[0] == 2
t[1] == 3
和简单的分片表达式一样，如果 a 是一个到数组的指针，那么 a[low : high : max] 是 (*a)[low : high : max] 的速记。如果被分片的操作数是一个数组，那么它必须是 可被寻址的 。

如果 0 <= low <= high <= max <= cap(a) ，那么索引是 在范围内 的，否则就 超出了范围 。 常量 索引必须是非负的且可以被类型 int 所代表的值；对于数组，常量索引也必须在范围内。如果多个索引为常量，那么存在的常量必须在相对彼此的范围内。如果在运行时索引超出了范围，那么会出现一个 run-time panic 。

类型断言
对于一个 接口类型 的表达式 x 以及一个类型 T ，主表达式

x.(T)
断言 x 不为 nil 并且存储在 x 中的值具有类型 T 。记法 x.(T) 被称为 类型断言 。

更准确地来说，如果 T 不是一个接口类型，那么 x.(T) 断言 x 的动态类型和类型 T 一致 。在这种情况下， T 必须 实现 x 的（接口）类型；否则类型断言是无效的，因为对于 x 来说存储一个类型为 T 的值是不可能的。如果 T 是一个接口类型，那么 x.(T) 断言 x 的动态类型实现了接口 T 。

如果类型断言成立，那么表达式的值为存储在 x 中的值，并且其类型为 T 。如果类型断言不成立，会发生一个 run-time panic 。换句话来说，即使 x 的动态类型仅在运行时已知， x.(T) 的类型也可以在一个正确的程序中被已知为 T 。

var x interface{} = 7    // x 有动态类型 int 以及值 7
i := x.(int)             // i 有类型 int 以及值 7

type I interface { m() }

func f(y I) {
  s := y.(string)        // 非法: string 没有实现 I (缺少方法 m)
  r := y.(io.Reader)     // r 有类型 io.Reader ，并且 y 的动态类型必须同时实现 I 和 io.Reader
  …
}
用于 赋值 或如下特殊格式的初始化中的类型断言

v, ok = x.(T)
v, ok := x.(T)
var v, ok = x.(T)
var v, ok T1 = x.(T)
产生一个额外的无类型的布尔值。如果断言成功，那么 ok 的值为 true 。否则为 false ，并且 v 的值为类型 T 的 零值 。这种情况下不会发生 run-time panic。

调用
给定一个函数类型为 F 的表达式 f ，

f(a1, a1, … an)
带实参 a1, a2, … an 调用了 f 。除了一种特殊情况以外，实参必须是单一值的 可分配 给 F 的参数类型的表达式，并且它们在函数调用之前就被求值好了。上述函数表达式的类型是 F 的结果类型。方法调用是类似的，但是方法本身是被指定为一个在该方法的接收者的值之上的选择器。

math.Atan2(x, y)  // 函数调用
var pt *Point
pt.Scale(3.5)     // 带接收者 pt 的方法调用
在一个函数调用中，函数值和实参使用 通常的顺序 被求值。在它们求值好后，调用的参数以值传递给函数，然后被调用的函数开始执行。函数的返回参数在函数返回时以值返回给调用的函数。

调用一个 nil 函数会发生 run-time panic 。

作为一个特殊情况，如果一个函数或方法 g 的返回值数量上等于且可以分别被分配给另一个函数或方法 f 的参数，那么调用 f(g(parameters_of_g)) 将会在按序绑定了 g 的返回值到 f 的参数后调用 f 。 f 这个调用必须排除 g 调用以外的参数，并且 g 必须要有最少一个返回值。如果 f 有一个最终的 ... 参数，这个参数会被分配那些在普通参数赋值完之后的剩余的 g 的返回值。

func Split(s string, pos int) (string, string) {
  return s[0:pos], s[pos:]
}

func Join(s, t string) string {
  return s + t
}

if Join(Split(value, len(value)/2)) != value {
  log.Panic("test fails")
}
如果 x 的（类型的）方法集包含了 m ，并且实参列表可以被分配给 m 的形参列表，那么方法调用 x.m() 是有效的。如果 x 是 可被寻址的 并且 &x 的方法集包含了 m ，那么 x.m() 是 (&x).m() 的速记：

var p Point
p.Scale(3.5)
这里没有明确的方法类型，也没有方法字面值。

传递实参给 ... 参数
如果 f 是带最终参数 p （其类型为 ...T ）的 variadic ，那么在 f 内， p 的类型等同于类型 []T 。如果 f 在调用时没有实参给 p ，那么传递给 p 的值为 nil 。否则，传递的值是一个新的类型为 []T 的分片，这个分片带一个底层数组，这个底层数组的连续的元素作为实参，并且必须 可分配 给 T 。因此该分配的长度和容量是绑定到 p 的实参的数量，而且每次调用可能会不同。

给定函数和调用

func Greeting(prefix string, who ...string)
Greeting("nobody")
Greeting("hello:", "Joe", "Anna", "Eileen")
在 Greeting 中， who 的值第一次调用时为 nil ，在第二次调用时为 []string{"Joe", "Anna", "Eileen"} 。

如果最终的实参可分配给一个分片类型 []T ，那么如果这个实参后跟着 ... 的话，它就会在不改变值的情况下传递给一个 ...T 参数。在这种情况下，不会创建新的分片。

给定一个分片 s 和调用

s := []string{"James", "Jasmine"}
Greeting("goodbye", s...)
在 Greeting 内， who 有和 s 有同一个值和同一个底层数组。

运算符
运算符把操作数结合进一个表达式。

Expression = UnaryExpr | Expression, binary_op, Expression .
UnaryExpr  = PrimaryExpr | unary_op, UnaryExpr .

binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
add_op     = "+" | "-" | "|" | "^" .
mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .

unary_op   = "+" | "-" | "!" | "^" | "*" | "&" | "<-" .
比较运算符会在 其它地方 讨论。对于其它二元运算符来说，操作数类型必须是 一致的 ，除非运算涉及位移或者无类型的 常量 。对于只涉及常量的运算，看 常量表达式 一节。

除了位移运算之外，如果一个操作数是无类型 常量 而另一个操作数不是，那么该常量会被隐式地 转换 为另一个操作数的类型。

在位移表达式的右侧的操作数必须为整数类型，或者可以被 uint 类型的值 所表示的 无类型的常量。如果一个非常量位移表达式的左侧的操作数是一个无符号常量，那么它会先被隐式地转换为假如位移表达式被其左侧操作数单独替换后的类型。

译注： 2019 年 7 月版，在上面这句话中，”无符号整数类型“变成了”整数类型“，看后文的描述，应该是正数即可，负数则会恐慌。

var s uint = 33
var i = 1<<s                  // 1 的类型为 int
var j int32 = 1<<s            // 1 的类型为 int32; j == 0
var k = uint64(1<<s)          // 1 的类型为 uint64; k == 1<<33
var m int = 1.0<<s            // 1.0 的类型为 int; 如果此处 int 为 32 比特大小的话， m == 0
var n = 1.0<<s == j           // 1.0 的类型为 int32; n == true
var o = 1<<s == 2<<s          // 1 和 2 的类型为 int; 如果此处 int 为 32 比特大小的话， o == true
var p = 1<<s == 1<<33         // 如果此处 int 为 32 比特大小的话则非法: 1 的类型为 int, 但是 1<<33 溢出了 int
                              // 译注: 1<<33 这种，为无类型的常量，所以在位移操作时不会溢出，但是在赋值时溢出了，所以报错；
                              //  　而 1<<s 会先使 1 带类型 int，所以在位移的时候已经溢出了，而位移溢出并不会报错，也就导致没有报错了。
var u = 1.0<<s                // 非法: 1.0 的类型为 float64, 不能位移
var u1 = 1.0<<s != 0          // 非法: 1.0 的类型为 float64, 不能位移
var u2 = 1<<s != 1.0          // 非法: 1 的类型为 float64, 不能位移
var v float32 = 1<<s          // 非法: 1 的类型为 float32, 不能位移
var w int64 = 1.0<<33         // 1.0<<33 是一个常量位移表达式
var x = a[1.0<<s]             // 1.0 的类型为 int；如果 int 是 32 位的话， x == a[0]
var a = make([]byte, 1.0<<s)  // 1.0 的类型为 int；如果 int 是 32 位的话， len(a) == 0
运算符优先级
一元运算符有最高的优先级。由于 ++ 和 -- 运算符构成了语句（而不是表达式），超出了运算符的结构。因此，语句 *p++ 等同于 (*p)++ 。

对于二元运算符来说有五个优先级。乘法运算符束缚力最强，接下来是加法运算符，比较运算符， && （逻辑与），和最后的 || （逻辑或）:

Precedence     Operator
    5             *  /  %  <<  >>  &  &^
    4             +  -  |  ^
    3             ==  !=  <  <=  >  >=
    2             &&
    1             ||
同一优先级的二元运算符按从左到右的顺序结合。比如， x / y * z 等同于 (x / y) * z 。

+x
23 + 3*x[i]
x <= f()
^a >> b
f() || g()
x == y+1 && <-chanPtr > 0
算数运算符
算数运算符应用于数字值，并产生一个和第一个操作数具有相同类型的结果。四个标准的算数运算符（ + , - , * , / ）应用于整数、浮点数和复数类型， + 还可以应用于字符串。位逻辑运算符和位移运算符仅应用于整数。

+    和                      整数，浮点数，复数值，字符串
-    差                      整数，浮点数，复数值
*    积                      整数，浮点数，复数值
/    商                      整数，浮点数，复数值
%    余                      整数

&    按位与　  (AND)          整数
|    按位或　  (OR)           整数
^    按位异或  (XOR)          整数
&^   按位清除  (AND NOT)      整数

<<   向左位移                 整数 << 无符号整数
>>   向右位移                 整数 >> 无符号整数
整数运算符
对于两个整数值 x 和 y ，其整数商 q = x / y 和余数 r = x % y 满足如下关系:

x = q*y + r 且 |r| < |y|
随着 x / y 截断到零（ 「截断除法」 ）。

 x     y     x / y     x % y
 5     3       1         2
-5     3      -1        -2
 5    -3      -1         2
-5    -3       1        -2
这个规则有一个例外，如果对于 x 的整数类型来说，被除数 x 是该类型中最负的那个值，那么，因为 补码 的 整数溢出 ，商 q = x / -1 等于 x （并且 r = 0 ）。

                         x, q
int8                     -128
int16                  -32768
int32             -2147483648
int64    -9223372036854775808
如果除数是一个 常量 ，那么它一定不能为零。如果在运行时除数为零，那么会发生一个 run-time panic 。如果被除数不为负值并且除数可以表示为以 2 为底数的一个次方常量，那么除法可以被向右位移所替换，计算余数可以被按位与运算符所替换:

 x     x / 4     x % 4     x >> 2     x & 3
 11      2         3         2          3
-11     -2        -3        -3          1
位移运算符通过右侧操作数（必须为正数）所指定的位移数来位移左侧的操作数。如果在运行时位移数为负，那么会发生一个 run-time panic 。如果左侧操作数是一个带符号的整数，那么位移运算符实现算数位移；如果是一个不带符号的整数，那么实现逻辑位移。位移数是没有上限的。对于 n 个位移数来说，位移行为犹如左侧操作数以间隔 1 来位移 n 次。因此， x << 1 等于 x*2 而 x >> 1 等于 x/2 ，不过向右位移会向负无穷截断。

对于整数操作数，一元运算符 + , - 和 ^ 有如下定义:

+x    　　　　               是 0 + x
-x    取其负值               是 0 - x
^x    按位补码               是 m ^ x ，其中对于无符号的 x 来说， m = 「所有位 置 1 」
      　　　　               　       　　　对于带符号的 x 来说， m = -1
      　　　　               　       　　　**-1 的话也是所有位置均为 1 ，但是这里需要考虑符号位**
整数溢出
对于无符号整数值来说， + , - , * 和 << 运算是以 2n 为模来计算的， n 为 无符号整数 类型的位宽。大致来说就是，这些无符号整数运算丢弃了溢出的高位，并且程序可以依赖于 "wrap around" 。

对于带符号整数值来说， + , - , * , / 和 << 运算可以合法地溢出，其产生的值是存在的并且可以被带符号整数表示法、其运算和操作数明确地定义。溢出不会发生 run-time panic 。编译器不会在不发生溢出这个假设情况下来优化代码。比如，它不会假设 x < x + 1 始终是真。

浮点运算符
对于浮点数和复数来说， +x 和 x 是一样的，但是 -x 是负的 x 。除了 IEEE-754 标准外，没有规定浮点数或者复数除以零的值；是否发生 run-time panic 是由具体实现规定的。

某个实现可能会在声明中结合多个浮点运算符为单一的融合运算符，然后产生一个与单独执行指令再取整所不同的值。一个显示的浮点类型 转换 会约到目标类型的精度，避免了融合会丢弃该舍入。

比如，有些架构提供了一个“积和熔加运算”（FMA）指令，该指令在运算 x*y + z 是不会先约取中间结果 x*y 。这些例子展示了什么时候 Go 实现会使用这个指令：

// FMA 允许被用来计算 r, 因为 x*y 不会被明确地约取：
r  = x*y + z
r  = z;   r += x*y
t  = x*y; r = t + z
*p = x*y; r = *p + z
r  = x*y + float64(z)

// FMA 不允许被用来计算 r, 因为它会省略 x*y 的舍入:
r  = float64(x*y) + z
r  = z; r += float64(x*y)
t  = float64(x*y); r = t + z
字符串连接
字符串使用使用 + 运算符或者 += 赋值运算符来连接。

s := "hi" + string(c)
s += " and good bye"
字符串加法通过连接操作数来创建了一个新的字符串。

比较运算符
比较运算符比较两个操作数，然后产生一个无类型的布尔值。

==    等于
!=    不等于
<     小于
<=    小于等于
>     大于
>=    大于等于
在每一个比较中，第一个操作数必须是 可分配 给第二个操作数的类型的，或者反过来。

相等运算符 == 和 != 应用到 可比较的 操作数上。排序运算符 < , <= , > 或 >= 应用到 可排序的 操作数上。术语以及比较的结果定义如下：

布尔值是可比较的。如果两个布尔值都为 true 或者 false ，那么它们相等。
通常情况下，整数值是可比较且可排序的。
浮点数值是可比较且可排序的，就像 IEEE-754 标准定义的。
复数值是可比较的。如果存在两个复数值 u 和 v ，满足 real(u) == real(v) 并且 imag(u) == imag(v) 的话，那么它们相等。
字符串值是可按字节顺序比较且排序的（按照字节的词法）。
指针值是可比较的。如果两个指针指向同一个变量，或者两个都为 nil 的话，那么它们相等。指向不同 零值 变量的指针可能相同也可能不同。
信道值是可比较的。如果两个信道值由同一个 make 调用来创建或者两个的值都为 nil ，那么它们相同。
接口值是可比较的。如果两个接口值有 一致的 动态类型以及相同的动态值，或者两个的值都为 nil ，那么它们相同。
当非接口类型 X 的值是可比较的且 X 实现了接口类型 T ，那么 X 的值 x 和 T 的值 t 是可比较的。如果 t 的动态类型和 X 一致并且 t 的动态值等于 x 的话，那么它们相等。
当结构体的所有字段都是可比较的，那么该结构体是可比较的。如果两个结构体的对应的非 空白 字段相等，那么两个结构体相等。
如果数组的元素值是可比较的，那么该数组是可比较的。如果两个数组对应的元素是相等的，那么两个数组相等。
当两个比较中的接口值的动态类型一致，但是该类型的值是不可比较的时候，会发生一个 run-time panic 。这种情况不仅仅发生在接口值比较上，同样也会发生在比较接口值数组或者带接口值字段的结构体上。

分片、映射和函数值是不可比较的。不过作为一个特例，一个分片、映射或者函数值可以和预先声明的标识符 nil 来比较。指针、信道和接口值与 nil 的比较也是允许的，并遵循上述通用规则。

const c = 3 < 4            // c 是无类型布尔常量"真"

type MyBool bool
var x, y int
var (
  // 比较的结果为一个无类型的布尔值。
  // 应用通用赋值规则。
  b3        = x == y // b3 类型为 bool
  b4 bool   = x == y // b4 类型为 bool
  b5 MyBool = x == y // b5 类型为 MyBool
)
逻辑运算符
逻辑运算符应用于 布尔 值，并产生一个和操作数相同类型的结果。右侧的操作数是按条件来求值的。

&&    条件 与     p && q  是  "如果 p 则 q 否则 false"
||    条件 或     p || q  是  "如果 p 则 true 否则 q"
!     非　 　     !p      是  "非 p"

**这边的条件与(&&)是和按位与(&)区分开来的，其它亦然**
地址运算符
对于类型为 T 的操作数 x 来说，地址运算 &x 生成了一个类型为 *T 的到 x 的指针。该操作数必须是 可被寻址的 ，也就是，一个变量、 指针间接pointer indirection 、或分片索引操作；或一个可寻址的结构体操作数的字段选择器；或一个可寻址的数组的数组索引操作。作为可被寻址要求的一个例外， x 也可以是（可能是括起来的） 复合字面值 。如果 x 的求值会导致一个 run-time panic ，那么 &x 的求值也会。

对于指针类型 *T 的操作数 x 来说，指针间接 *x 表示被 x 指向的类型为 T 的 变量 。如果 x 是 nil ，那么对于 *x 的求值尝试会导致一个 run-time panic 。

&x
&a[f(2)]
&Point{2, 3}
*p
*pf(x)

var x *int = nil
*x   // 导致一个 run-time panic
&*x  // 导致一个 run-time panic
接收运算符
对于 信道类型 的操作数 ch 来说，接收操作 <-ch 的值是从信道 ch 接收到的值。信道方向必须允许接收操作，并且接收操作的类型为信道的元素类型。直到一个值可用前该表达式都会阻塞。从一个 nil 信道接收值会永远阻塞下去。针对一个 closed 信道的接收操作总是会立即进行，并在之前已经发送完成的值被接收完毕后产生一个该元素类型的 零值 。

v1 := <-ch
v2 = <-ch
f(<-ch)
<-strobe  // 等待，直到时钟脉冲一次，并丢弃接收的值
用于 赋值 或特殊格式的初始化中的接收表达式

x, ok = <-ch
x, ok := <-ch
var x, ok = <-ch
var x, ok T = <-ch
产生一个额外的无类型的布尔值用于报告通信是否成功。如果接收的值被到该信道的成功的发送操作传递过来，那么 ok 的值为 true ，如果因为该信道已经关闭且为空，接收到的是零值，那么 ok 为 false 。

转换Conversions
转换会把一个表达式的 类型 改成被该转换所指定的类型。一个转换可能会在字面上出现在源文件中，也可能 隐含在 表达式所在的上下文中。

一个 显示的 转换是 T(x) 这样子形式的表达式，其中 T 是一个类型而 x 是一个可以被转换到类型 T 的一个表达式。

Conversion = Type, "(", Expression, [ "," ], ")" .
如果类型由运算符 * 或者 <- 开头，或者由关键字 func 开头并且没有结果列表，那么当必要时它必须被括起来以避免混淆：

*Point(p)        // 同 *(Point(p))
(*Point)(p)      // p 被转换为 *Point
<-chan int(c)    // 同 <-(chan int(c))
(<-chan int)(c)  // c 被转换为 <-chan int
func()(x)        // 函数签名 func() x
(func())(x)      // x 被转换为 func()
(func() int)(x)  // x 被转换为 func() int
func() int(x)    // x 被转换为 func() int (非歧义表达式)
如果一个 常量 值 x 可以被类型为 T 的值 所表示 ，那么 x 可以被转换为 T 。特殊情况下，整数常量 x 可以使用像非常量 x 一样的规则 被显示地转换为 字符串类型 。

常量转换产生一个带类型的常量来作为结果。

uint(iota)               // unit 类型的 iota 值
float32(2.718281828)     // float32 类型的 2.718281828
complex128(1)            // complex128 类型的 1.0 + 0.0i
float32(0.49999999)      // float32 类型的 0.5
float64(-1e-1000)        // float64 类型的 0.0
string('x')              // string 类型的 "x"
string(0x266c)           // string 类型的 "♬"
MyString("foo" + "bar")  // MyString 类型的 "foobar"
string([]byte{'a'})      // 不是常量: []byte{'a'} 不是常量
(*int)(nil)              // 不是常量: nil 不是常量， *int 不是布尔、数值或字符串类型
int(1.2)                 // 非法: 1.2 不能被 int 表示
string(65.0)             // 非法: 65.0 不是整数常量
非常量值 x 在以下这些情况下可以被转换为类型 T ：

x 可分配 给 T 。
忽略结构体标签（见下文）， x 的类型和 T 有 一致的 潜在类型 。
忽略结构体标签（见下文）， x 的类型和 T 都不是 定义的 指针类型，并且它们的基础类型有一致的潜在类型。
x 的类型和 T 都是整数或者浮点数类型。
x 的类型和 T 都是复数类型。
x 是一个整数或者一个字节/ rune 分片，并且 T 是字符串类型。
x 是一个字符串并且 T 是一个字节/ rune 分片。
在为了转换的目的而比较结构体类型是否一致时， 结构体的标签 是被忽略的：

type Person struct {
  Name    string
  Address *struct {
    Street string
    City   string
  }
}

var data *struct {
  Name    string `json:"name"`
  Address *struct {
    Street string `json:"street"`
    City   string `json:"city"`
  } `json:"address"`
}

var person = (*Person)(data)  // 忽略标签，潜在类型是一致的
数字类型之间或者数字类型和字符串类型之间的（非常量）转换有特殊的规则。这些转换可能改变 x 的表现方式并产生运行时成本. 所有其它的转换仅改变其类型而不会改变 x 的表现形式。

没有语言机制可以在指针和整数间做转换。在一些受限制的情况下，包 unsafe 实现了这个功能。

数字类型间的转换
以下的规则应用于非常量数值间的转换：

当在整数类型间做转换时，如果值是一个带符号整数，那么它会用符号位扩展到隐式的无限精度；否则它会用零扩展。然后它会截断以满足结果类型的大小。比如，如果 v := unit16(0x10F0) ，那么 uint32(int(v)) == 0xFFFFFFF0 。这种转换总是会产生一个有效的值；也不会有溢出指示。
当转换浮点数到整数时，小数部分会被丢弃（截断到零）。
当转换整数或者浮点数到浮点数类型，或者复数到其它复数类型时，结果值会约到目标类型所规定的精度。比如， float32 类型变量 x 的值可能使用超过 IEEE-754 32 位数的精度保存着，但是 float32(x) 表示的是把 x 的值约到 32 位精度的结果。类似的， x + 0.1 可能使用了超过 32 位精度，但是 float32(x + 0.1) 则不然。
在所有涉及浮点数或复数的非常量转换中，如果结果类型不能表示转换后的值，转换依旧是成功的，但结果值依赖实现。

从/到字符串的转换
转换带/不带符号的整数值到字符串类型会产生包含该数 UTF-8 表示形式的字符串。超过有效 Unicode 代码点范围的值会被转换为 \uFFFD 。
string('a')       // "a"
string(-1)        // "\ufffd" == "\xef\xbf\xbd"
string(0xf8)      // "\u00f8" == "ø" == "\xc3\xb8"
type MyString string
MyString(0x65e5)  // "\u65e5" == "日" == "\xe6\x97\xa5"
转换字节分片到字符串类型会产生一个以该分片的元素作为连续字节的字符串。
string([]byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'})   // "hellø"
string([]byte{})                                     // ""
string([]byte(nil))                                  // ""

type MyBytes []byte
string(MyBytes{'h', 'e', 'l', 'l', '\xc3', '\xb8'})  // "hellø"
转换 rune 分片到字符串会产生一个把独立的 rune 值转换为 string 后再级联的字符串。
string([]rune{0x767d, 0x9d6c, 0x7fd4})   // "\u767d\u9d6c\u7fd4" == "白鵬翔"
string([]rune{})                         // ""
string([]rune(nil))                      // ""

type MyRunes []rune
string(MyRunes{0x767d, 0x9d6c, 0x7fd4})  // "\u767d\u9d6c\u7fd4" == "白鵬翔"
转换字符串类型的值到字节类型的分片会产生一个以该字符串的字节作为连续元素的分片。
[]byte("hellø")   // []byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'}
[]byte("")        // []byte{}

MyBytes("hellø")  // []byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'}
转换字符串类型到 rune 类型分片会产生一个包含该字符串独立 Unicode 代码点的分片。
[]rune(MyString("白鵬翔"))  // []rune{0x767d, 0x9d6c, 0x7fd4}
[]rune("")                 // []rune{}

MyRunes("白鵬翔")           // []rune{0x767d, 0x9d6c, 0x7fd4}
常量表达式
常量表达式仅包含 常量 操作数，且是在编译的时候进行计算的。

不管以布尔、数字或者字符串类型作为操作数是否合法，使用对应的无类型的布尔、数字和字符串常量作为操作数是可以的。

常量 比较 总是会产生一个无类型的布尔常量。如果常量 位移表达式 的左侧操作数是一个无类型常量，那么其结果是一个整数常量；否则就是和左侧操作数同一类型的常量（必须是 整数类型 ）。

任何其它在无类型常量上的操作结果是同一个类别的无类型常量；也就是：布尔、整数、浮点数、复数或者字符串常量。如果一个二元运算（非位移）的无类型操作数是不同类的，那么其结果是在如下列表中靠后显示的操作数的类：整数、 rune、浮点数、复数。举例：无类型整数常量除以无类型复数常量会产生一个无类型的复数常量。

const a = 2 + 3.0          // a == 5.0   (无类型浮点数常量)
const b = 15 / 4           // b == 3     (无类型整数常量)
const c = 15 / 4.0         // c == 3.75  (无类型浮点数常量)
const Θ float64 = 3/2      // Θ == 1.0   (类型为 float64, 3/2 是整数除法)
const Π float64 = 3/2.     // Π == 1.5   (类型为 float64, 3/2. 是浮点除法)
const d = 1 << 3.0         // d == 8     (无类型整数常量)
const e = 1.0 << 3         // e == 8     (无类型整数常量)
const f = int32(1) << 33   // 非法的      (常量 8589934592 对于 int32 来说溢出了)
const g = float64(2) >> 1  // 非法的      (float64(2) 是一个带类型的浮点数常量)
const h = "foo" > "bar"    // h == true  (无类型布尔常量)
const j = true             // j == true  (无类型布尔常量)
const k = 'w' + 1          // k == 'x'   (无类型 rune 常量)
const l = "hi"             // l == "hi"  (无类型字符串常量)
const m = string(k)        // m == "x"   (字符串类型)
const Σ = 1 - 0.707i       //            (无类型复数常量)
const Δ = Σ + 2.0e-4       //            (无类型复数常量)
const Φ = iota*1i - 1/1i   //            (无类型复数常量)
把内置函数 complex 应用到无类型整数、 rune 或者浮点数常量会产生一个无类型的复数常量。

const ic = complex(0, c)   // ic == 3.75i  (无类型复数常量)
const iΘ = complex(0, Θ)   // iΘ == 1i     (complex128 类型)
常量表达式总是会被精确地求值；中间值和常量本身可能会需求比任何在语言中预定义的类型所支持的更大的精度。以下都是合法的声明：

const Huge = 1 << 100         // Huge == 1267650600228229401496703205376  (无类型整数常量)
const Four int8 = Huge >> 98  // Four == 4                                (int8 类型)
常量除法或取余操作的除数一定不能是零：

3.14 / 0.0   // 非法的：被零除了
带类型的 的常量的值必须是能被该常量类型所精确得 表示的 。以下常量表达式是非法的：

uint(-1)     // -1 不能作为 uint 来表示
int(3.14)    // 3.14 不能作为 int 来表示
int64(Huge)  // 1267650600228229401496703205376 不能作为 int64 来表示
Four * 300   // 操作数 300 不能作为 int8 (Four 的类型) 来表示
Four * 100   // 乘积 400 不能作为 int8 (Four 的类型) 来表示
用于一元按位补码运算符 ^ 的掩码符合非常量的规则：对于无符号常量来说所有位都是 1，而对于带符号且无类型的常量来说，则是一个 -1。

译注：无符号常量必定是带类型的（上文 常量 一节有写默认类型），所以对于掩码为 -1 的情况来说其实是一种情况，无类型常量也是带符号常量
^1         // 无类型整数常量，等于 -2
uint8(^1)  // 非法的: 相当于 uint8(-2)， -2 不能被 uint8 所表示
^uint8(1)  // 带类型的 uint8 常量， 相当于 0xFF ^ uint8(1) = uint8(0xFE)
int8(^1)   // 相当于 int8(-2)
^int8(1)   // 相当于 -1 ^ int8(1) = -2
实现限制：编译器可能会在计算无类型浮点数或者复数常量表达式时凑整；请参阅 常量 一节中的实现限制。该凑整可能导致在整数上下文内的浮点数常量表达式失效，即使它将在使用无限精度计算时是不可或缺的，反之亦然。

这段翻译应该有问题： This rounding may cause a floating-point constant expression to be invalid in an integer context, even if it would be integral when calculated using infinite precision, and vice versa.
求值顺序
在包的级别上， 初始化依赖 确定了 变量声明 中独立的初始化表达式的求值顺序。其它方面，当对表达式、赋值或者 return 语句 的 操作数 进行求值时，所有的函数调用、方法调用和通讯操作都是以词法的从左至右的顺序被求值的。

比如，在（函数本地）赋值

y[f()], ok = g(h(), i()+x[j()], <-c), k()
函数调用和通讯是按照 f() , h() , i() , <-c , g() 和 k() 的顺序发生的。不过，以上这些事件相比较于 x 的求值和索引，以及 y 的求值的顺序则是没有规定的。

a := 1
f := func() int { a++; return a }
x := []int{a, f()}            // x 可以是 [1, 2] 或是 [2, 2]： a 和 f() 的求值顺序没有被规定
m := map[int]int{a: 1, a: 2}  // m 可以是 {2: 1} 或是 {2: 2}： 两个映射赋值的求值顺序没有被规定
n := map[int]int{a: f()}      // n 可以是 {2: 3} 或是 {3: 3}： 键和值的求值顺序没有被规定
在包的级别上，初始化依赖会覆盖掉针对独立初始化表达式的从左至右的规则，但是不会针对在每个表达式中的操作数：

var a, b, c = f() + v(), g(), sqr(u()) + v()

func f() int        { return c }
func g() int        { return a }
func sqr(x int) int { return x*x }

// 函数 u 和 v 独立于其它所有的变量和函数
函数调用是按照 u() , sqr() , v() , f() , v() 和 g() 的顺序发生的。

在单一表达式中的浮点数操作是根据运算符的结合性来求值的。明确的括号会通过覆盖默认的结合性来影响求值。在表达式 x + (y + z) 中，加法 y + z 在加 x 前被执行。

语句
语句控制着执行。

Statement =
  Declaration | LabeledStmt | SimpleStmt |
  GoStmt | ReturnStmt | BreakStmt | ContinueStmt | GotoStmt |
  FallthroughStmt | Block | IfStmt | SwitchStmt | SelectStmt | ForStmt |
  DeferStmt .

SimpleStmt = EmptyStmt | ExpressionStmt | SendStmt | IncDecStmt | Assignment | ShortVarDecl .
终止语句
终止语句 阻止了同一个 块 中在其后（词法上）出现的语句的执行。以下语句是终结的：

"return" 或者 "goto" 语句。
对内置函数 panic 的调用。
语句列表以终止语句结束的 块 。
满足如下条件的 "if" 语句：
"else" 分支存在，并且
两个分支都是终止语句。
满足如下的 "for" 语句：
没有针对这个 "for" 语句的 "break" 语句，并且
循环条件为空。
满足如下的 "switch" 语句：
没有针对这个 "switch" 语句的 "break" 语句，
有一个 default case，并且
在每个 case 中（包括默认的）的语句列表以终止语句或者一个可能标记为 "fallthrough" 的语句结束。
满足如下的 "select" 语句：
没有针对这个 "select" 语句的 "break" 语句，并且
在每个 case 中（包括默认的）的语句列表是存在的并以终止语句结束。
标记终止语句的 标签语句 。
所有其它语句都不是终止的。

如果语句列表非空且其最后的非空语句是终止的，那么这个 语句列表 以终结语句结束。

空语句Empty statements
空语句什么都不做。

EmptyStmt = .
标签语句Labeled statements
标签语句可以是 goto , break 或 continue 语句的目标。

LabeledStmt = Label, ":", Statement .
Label       = identifier .
Error: log.Panic("error encountered")
表达式语句Expression statements
除了特定的内置函数外，函数/方法 调用 以及 接收操作 可以出现在语句上下文中。这种语句可能会被括起来。

ExpressionStmt = Expression .
下述内置函数不允许出现在语句上下文中：

append cap complex imag len make new real
unsafe.Alignof unsafe.Offsetof unsafe.Sizeof
h(x+y)
f.Close()
<-ch
(<-ch)
len("foo")  // 如果 len 是内置函数，那么是非法的
发送语句Send statements
发送语句在信道上发送一个值。信道表达式必须是 信道类型 ，信道方向必须允许发送操作，并且，发送值的类型必须 可分配 为信道的元素类型。

SendStmt = Channel, "<-", Expression .
Channel  = Expression .
信道和值表达式都会在通讯开始前被求值。直到发送进行前，通讯都是阻塞的。如果接收者准备好了那么在无缓冲的信道上的发送就可以进行了。如果缓冲区还有空间那么在带缓冲的信道上的发送就可以进行。在关闭的信道上进行发送会产生一个 run-time panic 。在值为 nil 的信道上的发送是会永久阻塞的。

ch <- 3 // 发送值 3 到信道 ch
自增/减语句IncDec statements
"++" 和 "--" 语句用无类型 常量 1 来增加或减少其操作数。和赋值一样，这个操作数必须是 可被寻址的 或者是一个映射索引表达式。

IncDecStmt = Expression, ( "++" | "--" ) .
以下 赋值 语句在语义上是等同的：

自增/减语句　         赋值
x++                 x += 1
x--                 x -= 1
赋值Assignments
Assignment = ExpressionList, assign_op, ExpressionList .
assign_op = [ add_op | mul_op ], "=" .
每个左侧的操作数必须是 可被寻址的 、一个映射索引表达式或（只对 = 赋值来说） 空白标识符 。操作数可能会被括起来。

x = 1
*p = f()
a[i] = 23
(k) = <-ch  // 同： k = <-ch
当 op 是一个二元 算数运算符 时， 赋值操作 x op= y 等同于 x = x op (y) ，不过 x 仅求值一次。 op= 构造是一个单独的记号。在赋值操作中，左侧和右侧的表达式列表都必须含有一个确切的单一值表达式，并且左侧的表达式不能为空白标识符。

a[i] <<= 2
i &^= 1<<n
多元赋值分配多值运算的独立的值到一个变量列表。有两种形式。第一种，右侧的操作数是譬如函数调用、 信道 、 映射 运算 、 类型断言 这样的单个多值表达式。左侧的操作数的个数必须和值的个数匹配。比如，如果 f 是一个返回两个值的函数，

x, y = f()
分配第一个值给 x 第二个给 y 。第二种形式，左侧操作数的个数必须等于右侧表达式的个数，每个表达式必须是单一的值，并且右侧第 n 个表达式会分配给左侧第 n 个操作数：

one, two, three = '一', '二', '三'
在赋值中， 空白标识符 提供了一个忽略右侧值的方法：

_ = x       // 对 x 求值，但是会忽略它
x, _ = f()  // 对 f() 求值，但是忽略了它的第二个结果值
赋值会分两个阶段进行。第一阶段，左侧的 索引表达式 和 指针间接 （包括在 选择器 中的隐式的指针间接）以及右侧的表达式都会按照 通常的顺序 来求值。第二阶段，赋值按从左至右的顺序进行。

a, b = b, a  // 交换 a 和 b

x := []int{1, 2, 3}
i := 0
i, x[i] = 1, 2  // 设 i = 1, x[0] = 2

i = 0
x[i], i = 2, 1  // 设 x[0] = 2, i = 1

x[0], x[0] = 1, 2  // 先设 x[0] = 1, 然后 x[0] = 2 （所以最后 x[0] == 2）

x[1], x[3] = 4, 5  // 设 x[1] = 4, 然后设 x[3] = 5

type Point struct { x, y int }
var p *Point
x[2], p.x = 6, 7  // 设 x[2] = 6, 然后设 p.x = 7

i = 2
x = []int{3, 5, 7}
for i, x[i] = range x {  // 设 i, x[2] = 0, x[0]
  break
}
// 循环结束后， i == 0 且 x == []int{3, 5, 3}
在赋值中，每个值都必须是 可分配 给需要分配的操作数的类型的，不过会有以下特殊情况：

任何类型的值都可以被分配给空白标识符。
当无类型常量被分配给一个接口类型变量或是空白标识符时，常量会先被隐式地 转换 为它的 默认类型 。
当无类型布尔值被分配给一个接口类型变量或是空白标识符时，它会先被隐式地转换为布尔类型。
If 语句
"if" 语句根据布尔表达式的值来指定两个分支的条件执行。当表达式求值得真时， "if" 分支被执行，否则执行 "else" 分支（存在的话）。

IfStmt = "if", [ SimpleStmt, ";" ], Expression, Block, [ "else", ( IfStmt | Block ) ] .
if x > max {
  x = max
}
表达式前面可能会有一个简单的语句，这个语句会在表达式求值之前被执行。

if x := f(); x < y {
  return x
} else if x > z {
  return z
} else {
  return y
}
Switch 语句
"switch" 语句提供了多路执行。表达式或者类型指示符会和在 "switch" 内的 "case" 做比较去确定执行哪一个分支。

SwitchStmt = ExprSwitchStmt | TypeSwitchStmt .
有两种形式：表达式开关（switch）和类型开关。在表达式开关中， case 包含了要与 switch 表达式的值比较的表达式。在类型开关中， case 包含了要与特别说明的 switch 表达式的类型比较的类型。 switch 表达式在一个开关语句中仅求值一次。

表达式开关
在表达式开关中， switch 表达式和 case 表达式（不能是常量）是按照从左至右、从上之下的顺序进行求值的；第一个和 switch 表达式相等的 case 中对应的语句会被触发执行；其它 case 会被跳过。如果没有 case 匹配且有一个 "default" case，那么会执行这个 case 的语句。最多有一个默认 case ，它可以出现在 "switch" 语句的任意位置。当 switch 表达式不存在时，相当于是一个布尔值 true 。

ExprSwitchStmt = "switch", [ SimpleStmt, ";" ], [ Expression ], "{", { ExprCaseClause }, "}" .
ExprCaseClause = ExprSwitchCase, ":", StatementList .
ExprSwitchCase = "case", ExpressionList | "default" .
如果 switch 表达式求值为一个无类型常量，它会先被隐式地 转换 为它的 默认类型 ；如果它是一个无类型的布尔值，它会先被隐式地转换为类型 bool 。预定义的无类型值 nil 不能用在 switch 表达式中。

如果 case 表达式是无类型的，那么它会先被隐式地 转换 为 switch 表达式的类型。对于每个（可能是转换过的） case 表达式 x 和 switch 表达式的值 t ， x == t 必定是一个有效的 比较 。

也就是说， switch 表达式就像是被用来声明和初始化一个没有明确类型的临时变量 t ；为了测试相等性，这个临时变量 t 的值会和每一个 case 表达式 x 做判断。

在一个 case 或 default 子句中，最后的非空语句可能是一个（可能是 标签 的） "fallthrough" 语句用来指示控制应该从本子句流出以流入下个子句的第一个语句。不然的话控制会流到 "switch" 语句的末尾。 "fallthrough" 语句可以作为除了表达式开关的最后一个子句外的其它所有子句的最后一条语句出现。

switch 表达式可以前缀一个简单的语句，这个语句会在表达式之前被求值。

switch tag {
default: s3()
case 0, 1, 2, 3: s1()
case 4, 5, 6, 7: s2()
}

switch x := f(); {  // 缺少 switch 表达式就意味着 "true"
case x < 0: return -x
default: return x
}

switch {
case x < y: f1()
case x < z: f2()
case x == 4: f3()
}
实现限制：编译器可能会不允许多个 case 表达式求值结果为相同的常量。例如，现在的编译器不允许重复的整数、浮点数或字符串常量出现在 case 表达式中。

类型开关
类型开关用于比较类型而不是值。其它方面和表达式开关类似。它的标识是一个特殊的 switch 表达式，这个表达式形式是一个使用了保留字 type 而不是一个实际类型的 类型断言 。

switch x.(type) {
// cases
}
然后 case 匹配实际的类型 T 而不是表达式 x 的动态类型。与类型断言一样， x 必须是 接口类型 ，并且在 case 中的每一个非接口类型 T 必须实现 x 的类型。在类型开关的 case 中的类型必须都是 不同的 。


TypeSwitchStmt  = "switch", [ SimpleStmt, ";" ], TypeSwitchGuard, "{", { TypeCaseClause }, "}" .
TypeSwitchGuard = [ identifier, ":=" ], PrimaryExpr, ".", "(", "type", ")" .
TypeCaseClause  = TypeSwitchCase, ":", StatementList .
TypeSwitchCase  = "case", TypeList | "default" .
TypeList        = Type, { ",", Type } .
TypeSwitchGuard 可能会包含一个 短变量声明 。当用了这种形式的话，变量会在每个子句的 TypeSwitchCase 末尾的隐式 块 中被声明。在只列出一个类型的 case 的子句中，变量类型就是这个类型；否则，变量类型为 TypeSwitchGuard 中表达式的类型。

不同于类型， case 可以使用预声明的标识符 nil ；这种会在 TypeSwitchGuard 中的表达式为 nil 接口值时被选择。只能最多一个 nil case。

给定一个 interface{} 类型的表达式 x ，以下类型开关：

switch i := x.(type) {
case nil:
  printString("x is nil")                // i 类型为 x 的类型（interface{}）
case int:
  printInt(i)                            // i 类型为 int
case float64:
  printFloat64(i)                        // i 类型为 float64
case func(int) float64:
  printFunction(i)                       // i 类型为 func(int) float64
case bool, string:
  printString("type is bool or string")  // i 类型为 x 的类型（interface{}）
default:
  printString("don't know the type")     // i 类型为 x 的类型（interface{}）
}
可以被重写为：

v := x  // x 只被求值一次
if v == nil {
  i := v                                 // i 类型为 x 的类型（interface{}）
  printString("x is nil")
} else if i, isInt := v.(int); isInt {
  printInt(i)                            // i 类型为 int
} else if i, isFloat64 := v.(float64); isFloat64 {
  printFloat64(i)                        // i 类型为 float64
} else if i, isFunc := v.(func(int) float64); isFunc {
  printFunction(i)                       // i 类型为 func(int) float64
} else {
  _, isBool := v.(bool)
  _, isString := v.(string)
  if isBool || isString {
    i := v                         // i 类型为 x 的类型（interface{}）
    printString("type is bool or string")
  } else {
    i := v                         // i 类型为 x 的类型（interface{}）
    printString("don't know the type")
  }
}
TypeSwitchGuard 可以前缀一个简单的语句，这个语句在 guard 之前被求值。

"fallthrough" 语句在类型开关中是不被允许的。

For 语句
"for" 语句规定了一个块的重复执行。有三种形式：迭代可以被一个单一条件、一个 "for" 子句或是一个 "range" 子句控制。

ForStmt = "for", [ Condition | ForClause | RangeClause ], Block .
Condition = Expression .
带单一条件的 for 语句
在它最简单的形式中， "for" 语句就像一个求值为真的布尔条件一样来规定一个块的重复执行。这个条件的值会在每次迭代前都被求一下。如果条件为空，那么就相当于布尔值 true 。

for a < b {
  a *= 2
}
带 for 子句的 for 语句
带一个 forClause 的 "for" 子句也是通过其条件来控制的，但是它会额外指定一个 init 和 post 语句，比如一个赋值、增量或减量语句。 Init 语句可以是一个 短变量声明 ，但 post 语句一定不是。通过 init 语句声明的变量会在每次迭代时被重复使用。

ForClause = [ InitStmt ], ";", [ Condition ], ";", [ PostStmt ] .
InitStmt = SimpleStmt .
PostStmt = SimpleStmt .
for i := 0; i < 10; i++ {
  f(i)
}
如果非空， init 语句会在首次迭代的条件求值前被执行一次； post 语句会在每次块执行完后被执行（并且只有在块有执行过后）。 ForClause 每个元素都可以是空的，但是 分号 是必须要有的，除非仅存在一个条件元素。如果条件为空，那么就相当于布尔值 true 。

for cond { S() }    同    for ; cond ; { S() }
for      { S() }    同    for true     { S() }
带 range 子句的 for 语句
带 "range" 子句的 "for" 语句会彻底地迭代数组、分片、字符串或映射的所有条目，或是从信道接收到的值。针对每一个条目，它在分配 迭代值 给对应且存在的 迭代变量 后再执行语句块。

RangeClause = [ ExpressionList, "=" | IdentifierList, ":=" ], "range", Expression .
"range" 子句中右侧的表达式被称为 范围表达式 ，它可以是数组、到数组的指针、分片、字符串、映射或是允许 接收操作 的信道。和赋值一样，如果左侧操作数存在，那么它一定是 可被寻址的 或映射索引表达式；它们表示为迭代变量。如果范围表达式是一个信道，那么最多允许一个迭代变量，其它情况下可以最多到两个。如果最后的迭代变量是 空白标识符 ，那么这个 range 子句和没有这个标识符的子句是相同的。

范围表达式 x 会在开始此循环前被求值一次，但有一个例外：当存在最多一个迭代变量且 len(x) 是 常量 时，范围表达式是不被求值的。

左侧的函数调用在每次迭代时被求值。对于每个迭代，如果迭代变量存在，那么对应的迭代值是按以下说明产生的：

范围表达式　                                第一个值　　         第二个值

array or slice  a  [n]E, *[n]E, or []E    index    i  int    a[i]       E
string          s  string type            index    i  int    看下面的 rune
map             m  map[K]V                key      k  K      m[k]       V
channel         c  chan E, <-chan E       element  e  E
对于数组、到数组的指针或是分片值 a ，其索引迭代值是从索引 0 开始，以递增次序产生的。如果存在最多一个迭代变量， range 循环会创建从 0 到 len(a) - 1 的迭代值，且不会索引进数组或分片内。对于 nil 分片而言，迭代数是 0。
对于字符串值， "range" 子句从字节索引 0 开始迭代字符串中的 Unicode 代码点。在连续的迭代上，索引值是字符串中连续 UTF-8 编码的代码点的第一个字节的索引，而第二个值（类型是 rune ）是对应的代码点的值。如果迭代遇到了无效的 UTF-8 序列，那么第二个值会变成 Unicode 替换字符 0xFFFD ，且下一个迭代将在字符串中前进一个字节。
映射的迭代顺序是未指定的，并且不能保证两次完整的迭代是相同的。如果在迭代中某个未接触到的映射条目被移除了，那么对应的迭代值就不会产生。如果在迭代中新创建了一个映射条目，那这个条目可能会在迭代中被产生也可能被跳过。对于每个条目的创建或是一个迭代到下一个迭代，选择可能很多样。如果映射是 nil ，迭代数为 0。
对于信道，迭代值是在信道上发送的直到信道 关闭 的连续值。如果信道是 nil ，那么范围表达式会永久阻塞。
迭代值会像 赋值语句 一样被赋值给对应的迭代变量。

迭代变量可以被 "range" 子句使用 短变量声明 （:=）的形式声明。这种情况下，它们的类型会被设置为对应迭代值的类型，且它们的 作用域 是 "for" 语句块；这些变量会在每次迭代时复用。

如果迭代变量是在 "for" 语句外被声明的，那么在执行完毕后，它们的值会是最后一次迭代的值。

var testdata *struct {
  a *[7]int
}
for i, _ := range testdata.a {
  // testdata.a 不会被求值; len(testdata.a) 是常量
  // i 范围从 0 到 6
  f(i)
}

var a [10]string
for i, s := range a {
  // i 类型为 int
  // s 类型为 string
  // s == a[i]
  g(i, s)
}

var key string
var val interface {}  // m 的元素类型可赋予 val
m := map[string]int{"mon":0, "tue":1, "wed":2, "thu":3, "fri":4, "sat":5, "sun":6}
for key, val = range m {
  h(key, val)
}
// key == 迭代中遇到的最后一个映射键
// val == map[key]

var ch chan Work = producer()
for w := range ch {
  doWork(w)
}

// 清空信道
for range ch {}
Go 语句
"go" 语句会在同一地址空间执行一个函数调用作为一单独的并发控制流程（ goroutine ）。

GoStmt = "go", Expression .
表达式必须是函数或方法调用；它不能是括起来的。对内置函数的调用会有和 表达式语句 一样的限制。

在调用的 goroutine 中的函数值和参数是按 通常的情况来求值 的，但不同于普通调用的是，程序执行不会等待被调用的函数执行完毕。相反，在新的 goroutine 中的函数是独立执行的。当函数终止，其 goroutine 也会终止。如果函数存在任何返回值，这些值会在函数完成时被丢弃。

go Server()
go func(ch chan<- bool) { for { sleep(10); ch <- true\}\} (c)
Select 语句
"select" 语句会选择一组或是 发送 或是 接收 的操作来进行。它看起来和 "switch" 语句类似，但它所有的 case 只涉及通讯操作。

SelectStmt = "select", "{", { CommClause }, "}" .
CommClause = CommCase, ":", StatementList .
CommCase   = "case", ( SendStmt | RecvStmt ) | "default" .
RecvStmt   = [ ExpressionList, "=" | IdentifierList, ":=" ], RecvExpr .
RecvExpr   = Expression .
带 RecvStmt 的 case 可能会分配 RecvExpr 的结果到一个或两个变量，变量是用 短变量声明 声明的。 RecvExpr 一定是一个（可能是括起来的）接收操作。最多可以有一个默认 case ，它可以出现在 case 列表的任意位置。

"select" 语句的执行按如下几个步骤进行：

对于语句中的所有 case 来说，其接收操作和信道的信道操作数以及发送语句右侧的表达式会在进入 "select" 语句时以源码的顺序被执行仅一次。结果是需要接收或发送的信道集，以及对应的需要发送的值。无论选择哪个（如果有）通讯操作进行，在这个求值中的任何副作用都会发生。 RecvStmt 左侧的带短变量声明或赋值的表达式还不会被求值。
如果可以发生一个或多个通讯，通过统一的伪随机选择确定一个来进行。否则，如果有一个默认的 case，那么这个 case 会被选择。如果没有默认的 case，那么这个 "select" 语句会阻塞，直到至少发生了一个通讯。
除非被选择的 case 是默认的 case，否则各自的通讯操作会被执行。
如果被选择的 case 是一个带短变量声明或赋值的 RecvStmt，那么左侧的表达式会被求值且接收到的值会被分配。
被选择的 case 的语句列表被执行。
由于在 nil 信道上的通讯永不会进行，所以只带 nil 信道且没有默认 case 的 select 会永久阻塞。

var a []int
var c, c1, c2, c3, c4 chan int
var i1, i2 int
select {
case i1 = <-c1:
  print("received ", i1, " from c1\n")
case c2 <- i2:
  print("sent ", i2, " to c2\n")
case i3, ok := (<-c3):  // 同： i3, ok := <-c3
  if ok {
    print("received ", i3, " from c3\n")
  } else {
    print("c3 is closed\n")
  }
case a[f()] = <-c4:
  // 同：
  // case t := <-c4
  //  a[f()] = t
default:
  print("no communication\n")
}

for {  // 发送（伪）随机比特序列到 c
  select {
  case c <- 0:  // 注意：没有语句，没有 fallthrough，没有可折叠的 case
  case c <- 1:
  }
}

select {}  // 永久阻塞
Return 语句
函数 F 中的 "return" 语句会终止 F 的执行，并可选择地提供一个或更多的返回值。任何被 F 推迟 的函数会在 F 返回到它调用者前被执行。

ReturnStmt = "return", [ ExpressionList ] .
在没有结果类型的函数中， "return" 语句一定不指定任何返回值。

func noResult() {
  return
}
有三种从带结果类型的函数内返回值的方法：

返回值会明确地列在 "return" 语句中。每个表达式一定是单一值的且是 可分配 给对应的函数返回类型的元素。
func simpleF() int {
  return 2
}

func complexF1() (re float64, im float64) {
  return -7.0, -4.0
}
在 "return" 语句中的表达式列表可以是对多值函数的单一调用。效果就犹如从这个函数返回的值被分配给带对应值类型的一个临时变量，然后这些变量会跟随在 "return" 语句后，并适用上述情况指明的规则。
func complexF2() (re float64, im float64) {
  return complexF1()
}
如果函数结果值对其 结果参数 规定了名字，那么表达式列表可以为空。结果参数会作为本地变量，函数也可以在需要时给它们赋值。 "return" 语句会返回这些变量的值。
func complexF3() (re float64, im float64) {
  re = 7.0
  im = 4.0
  return
}

func (devnull) Write(p []byte) (n int, _ error) {
  n = len(p)
  return
}
不管它们是如何声明的，在进入函数时，所有结果值都会被初始化为其类型的 零值 。指定结果的 "return" 语句会在任何推迟函数执行前设置结果参数。

实现限制：当一个和结果参数同名的实体（常量、类型或变量）在 return 位置的 作用域 内时，编译器会不允许空的表达式列表出现在 "return" 语句中。

func f(n int) (res int, err error) {
  if _, err := f(n-1); err != nil {
    return  // 无效的返回语句： err 被遮蔽了
  }
  return
}
Break 语句
"break" 语句终止在相同函数内最内层的 "for" , "switch" 或 "select" 语句的执行。

BreakStmt = "break", [ Label ] .
如果这里有一个标签，那它必须是一个封闭的 "for" 、 "switch" 或 "select" 语句，然后这个就是被终止执行的那个。

OuterLoop:
  for i = 0; i < n; i++ {
    for j = 0; j < m; j++ {
      switch a[i][j] {
      case nil:
        state = Error
        break OuterLoop
      case item:
        state = Found
        break OuterLoop
      }
    }
  }
Continue 语句
"continue" 语句在发布位置开始执行最内层 "for" 循环的下一次迭代。 "for" 循环必须在同一个函数内。

ContinueStmt = "continue", [ Label ] .
如果这里有一个标签，那么必须是一个闭合的 "for" 语句，然后这个就是被执行功能的那个。

RowLoop:
  for y, row := range rows {
    for x, data := range row {
      if data == endOfRow {
        continue RowLoop
      }
      row[x] = data + bias(x, y)
    }
  }
Goto 语句
"goto" 语句转移控制到相同函数内对应标签的语句。

GotoStmt = "goto", Label .
goto Error
执行 "goto" 语句一定不会使任何在 goto 点位时还不在 作用域 内的变量进入作用域。例如，这个例子：

  goto L  // 坏的
  v := 3
L:
是错误的，因为跳转到标签 L 越过了 v 创建。

在某个 块 外的 "goto" 语句不能跳转到这个块内。例如，这个例子：

if n%2 == 1 {
  goto L1
}
for n > 0 {
  f()
  n--
L1:
  f()
  n--
}
是错误的，因为标签 L1 在 "for" 语句块内，但是 "goto" 不在。

Fallthrough 语句
"fallthrough" 语句转移控制给 表达式 "switch" 语句 内下一个 case 子句的第一条语句。它仅作为此类子句的最终非空语句使用。

FallthroughStmt = "fallthrough" .
Defer 语句
"defer" 语句会调用一个被推迟到其环绕函数返回瞬间执行的函数，而函数返回的原因要么是执行了一个 返回语句 、到达了 函数体 的底部，要么是对应的 goroutine panicking 了。

DeferStmt = "defer", Expression .
这个表达式一定是一个函数或者方法调用；它不能是括起来的。对内置函数的调用会如 表达式语句 一样被限制。

每次 "defer" 语句执行时，针对调用的函数值和参数是按 通常的情况来求值 并重新保存的，但实际的函数是不调用的。相反，被推迟的函数会在其环绕函数返回前，按照被推迟的反序被瞬间调用。也就是说，如果围绕函数通过一个明确的 return 语句 返回的话，那么被推迟的函数会在所有被 return 语句所设置的结果参数 后 ，在函数返回到其调用者 前 被执行。如果推迟函数求值得 nil ，那么在函数被调用时（而不是在 "defer" 语句被执行时），执行会 恐慌 。

例如，如果被推迟的函数是一个 函数字面值 并且其环绕函数有在该字面值作用域内的 命名的结果参数 ，那么该被推迟的函数可以在那些结果参数被返回前访问并修改它们。如果被推迟的函数有任何返回值，这些值会在函数完成时被丢弃。（也看一下 处理恐慌 一节）

lock(l)
defer unlock(l)  // 解锁发生在环绕函数返回前

// 在环绕函数返回前打印 3 2 1 0
for i := 0; i <= 3; i++ {
  defer fmt.Print(i)
}

// f 会返回 42
func f() (result int) {
  defer func() {
    // 结果会在其被 return 语句设为 6 之后再被访问
    result *= 7
  }()
  return 6
}
内置函数
内置函数是 预先声明 的。它们和其它任何函数一样调用，但是其中有一些能接受类型而不是表达式作为其第一个实参。

内置函数没有标准的 Go 类型，所以它们只能出现在 调用 表达式中；它们不能作为函数值来使用。

Close
对于信道 c ，内置函数 close(c) 标明了将不会再有值被发送到这个信道。如果 c 是一个仅可接收的信道，那么会出错。发送到或者关闭一个已经关闭的信道会发生 run-time panic 。 关闭 nil 信道也会发生 run-time panic 。调用 close 后，以及任何之前被发送的值都被接收后，接收操作不会阻塞而将是会返回对应信道类型的零值。多值 接收操作 会返回一个接收到的值，随同一个信道是否已经被关闭的指示符。

长度和容量
内置函数 len 和 cap 获取各种类型的实参并返回一个 int 类型结果。实现会保证结果总是一个 int 值。

调用　     实参类型　　       结果

len(s)    字符串类型　       按字节表示的字符串长度
          [n]T, *[n]T      数组长度（== n）
          []T              分片长度
          map[K]T          映射长度（定义的键的个数）
          chan T           在信道缓冲区内排队的元素个数

cap(s)    [n]T, *[n]T      数组长度（== n）
          []T              分配容量
          chan T           信道缓冲区容量
分片的容量是为其底层数组所分配的空间所对应的元素个数。任何时间都满足如下关系：

0 <= len(s) <= cap(s)
nil 分片、映射或者信道的长度是 0。 nil 分片或信道的容量是 0。

如果 s 是一个字符串常量，那么 len(s) 是一个 常量 。如果 s 类型是一个数组或到数组的指针且表达式 s 不包含 信道接收 或（非常量的） 函数调用 的话， 那么表达式 len(s) 和 cap(s) 是常量；这种情况下， s 是不求值的。否则的话， len 和 cap 的调用不是常量且 s 会被求值。

const (
  c1 = imag(2i)                    // imag(2i) = 2.0 是一个常量
  c2 = len([10]float64{2})         // [10]float64{2} 不包含函数调用
  c3 = len([10]float64{c1})        // [10]float64{c1} 不包含函数调用
  c4 = len([10]float64{imag(2i)})  // imag(2i) 是一个常量且没有函数调用
  c5 = len([10]float64{imag(z)})   // 无效的: imag(z) 是一个非常量的函数调用
)
var z complex128
分配
内置函数 new 获取一个类型 T ，在运行时为该类型的 变量 分配地址空间，并返回一个 指向 它的类型为 *T 的值。这个变量会按照 初始化值 一节所描述的来初始化。

new(T)
例如：

type S struct { a int; b float64 }
new(S)
为 S 类型变量分配存储空间，初始化它（ a=0, b=0.0 ），然后返回含有位置地址的类型为 *S 的一个值。

制作分片、映射和信道
内置函数 make 获取一个分片、映射或信道类型 T ，可选择性的接一个类型相关的表达式列表。它会返回类型 T 的值（不是 *T ）。内存会按照 初始化值 一节所描述的来初始化。

调用　            类型 T　    结果

make(T, n)       分片　      带 n 长度和容量的类型为 T 的分片
make(T, n, m)    分片　      带 n 长度和 m 容量的类型为 T 的分片

make(T)          映射　      类型为 T 的映射
make(T, n)       映射　      为约 n 个元素分配了初始化空间的类型为 T 的映射

make(T)          信道　      类型为 T 的无缓冲区信道
make(T, n)       信道　      类型为 T 的带缓冲区且缓冲区大小为 n 的信道
每个大小实参 n 和 m ，必须为整数类型或一个无类型的 常量 。常量大小实参必须是非负的且可被 int 类型值 所表示的 ；如果它是个无类型常量，那么会被给定类型 int 。如果 n 和 m 都提供了且为常量，那么 n 一定不能大于 m 。如果在运行时 n 为负值或者大于了 m ，那么会发生 run-time panic 。

s := make([]int, 10, 100)       // len(s) == 10, cap(s) == 100 的分片
s := make([]int, 1e3)           // len(s) == cap(s) == 1000 的分片
s := make([]int, 1<<63)         // 非法的: len(s) 不能被 int 类型的值所表示
s := make([]int, 10, 0)         // 非法的: len(s) > cap(s)
c := make(chan int, 10)         // 带大小为 10 的缓冲区的信道
m := make(map[string]int, 100)  // 带为约 100 个元素初始化空间的映射
带映射类型和大小提示 n 来调用 make 会创建一个带持有 n 个映射元素初始化空间的映射。其精度表现是依赖实现的。

添加到和拷贝分片
内置函数 append 和 copy 会协助常见的切片操作。对于这两个函数，其结果和实参的内存引用是否重叠无关。

variadic 函数 append 附加零个或多个值 x 到必须为分片类型的 S 类型的 s ，并返回结果分片，也是 S 类型。值 x 是传递给类型为 ...T 的一个形参，其中 T 是 S 的 元素类型 并应用对应的 参数传递规则 。作为一个特殊的情况， append 也接受首个为可分配给类型 []byte 的实参，且第二个为后缀 ... 的字符串类型的实参。这种形式会附加字符串的字节。

append(s S, x ...T) S  // T 是 S 的元素类型
如果 s 的容量不足以满足额外的值，那么 append 会分配一个新的足够大的底层数组来同时满足已经存在的分片元素和那些额外的值。否则， append 复用原来的底层数组。

s0 := []int{0, 0}
s1 := append(s0, 2)                // 附加一个单一元素　     s1 == []int{0, 0, 2}
s2 := append(s1, 3, 5, 7)          // 附加多个元素          s2 == []int{0, 0, 2, 3, 5, 7}
s3 := append(s2, s0...)            // 附加一个分片          s3 == []int{0, 0, 2, 3, 5, 7, 0, 0}
s4 := append(s3[3:6], s3[2:]...)   // 附加重叠的分片　　     s4 == []int{3, 5, 7, 2, 3, 5, 7, 0, 0}

var t []interface{}
t = append(t, 42, 3.1415, "foo")   //                     t == []interface{}{42, 3.1415, "foo"}

var b []byte
b = append(b, "bar"...)            // 附加字符串内容　　     b == []byte{'b', 'a', 'r' }
函数 copy 从源 src 拷贝分片元素到目的 dst 并返回拷贝的元素个数。两个实参必须有 一致的 元素类型 T 并且必须是 可分配 给类型为 []T 的分片的。拷贝的元素格式是 len(src) 和 len(dst) 中的最小值。作为一个特殊情况， copy 也接受目标实参可分配为 []byte 类型而源实参为字符串类型。这种形式会从字符串中拷贝字节到字节分片中。

copy(dst, src []T) int
copy(dst []byte, src string) int
例子：

var a = [...]int{0, 1, 2, 3, 4, 5, 6, 7}
var s = make([]int, 6)
var b = make([]byte, 5)
n1 := copy(s, a[0:])            // n1 == 6, s == []int{0, 1, 2, 3, 4, 5}
n2 := copy(s, s[2:])            // n2 == 4, s == []int{2, 3, 4, 5, 4, 5}
n3 := copy(b, "Hello, World!")  // n3 == 5, b == []byte("Hello")
映射元素的删除
内置函数 delete 会根据键 k 从 映射 m 中删除元素。 k 的类型必须是 可分配 给 m 的键类型的。

delete(m, k)  // 从映射 m 中删除元素 m[k]
如果映射 m 是 nil 或元素 m[k] 不存在，那么 delete 是一个空操作。

操纵复数
有三个函数用来聚合和分解复数。内置函数 complex 用浮点的实和虚部来构造一个复值，而 real 和 imag 从一个复值中提取其实部和虚部。

complex(realPart, imaginaryPart floatT) complexT
real(complexT) floatT
imag(complexT) floatT
实参的类型和返回值对应。 对于 complex ，两个实参必须是相同的浮点类型，并且返回值类型是带对应浮点成分的复合类型， complex64 对应 float32 实参， complex128 对应 float64 实参。如果有一个实参求值为一个无类型的常量，那么它会先被隐式地 转换 为另一个实参类型。如果两个实参都求值为无类型常量，那么它们必须是非复合数或者它们的虚部一定为零，然后函数的返回值也是一个无类型复合常量。

对于 real 和 imag ，实参必须是复合类型，返回值是对应的浮点类型： float32 对应一个 complex64 实参， float64 对应一个 complex128 实参。如果实参求值为一个无类型常量，那么它必须是一个数，然后函数的返回类型是一个无类型的浮点常量。

real 和 imag 函数一起组成了 complex 的反相，所以对于一个复合类型 Z 的值 z 来说， z == Z(complex(real(z), imag(z))) 。

如果这些函数的操作数都是常量，那么返回值也是一个常量。

var a = complex(2, -2)             // complex128
const b = complex(1.0, -1.4)       // 无类型复合常量 1 - 1.4i
x := float32(math.Cos(math.Pi/2))  // float32
var c64 = complex(5, -x)           // complex64
var s int = complex(1, 0)          // 无类型复合常量 1 + 0i 可以被转化为 int
_ = complex(1, 2<<s)               // 非法的： 2 被认为是浮点类型，不能位移
var rl = real(c64)                 // float32
var im = imag(a)                   // float64
const c = imag(b)                  // 无类型常量 -1.4
_ = imag(3 << s)                   // 非法的： 3 被认为是复合类型，不能位移
处理恐慌
有两个内置函数， panic 和 recover ，协助报告和处理 run-time panic 和程序定义的错误状态。

func panic(interface{})
func recover() interface{}
当执行函数 F 时，对 panic 的明确调用或 run-time panic 会终止 F 的执行。任何被 F 推迟 的函数会照常执行。然后，任何被 F 的调用者所推迟的函数会运行，以此类推直到被在执行中 goroutine 中的顶层函数所推迟的。在这个阶段，程序会终止并且错误状态会被报告，包括给 panic 的实参的值。这个终止过程被称为 panicking 。

panic(42)
panic("unreachable")
panic(Error("cannot parse"))
recover 函数允许程序管理一个 panicking goroutine 的行为。假设函数 G 推迟了调用 recover 的函数 D ，且恐慌发生在了和 G 执行的同一个 goroutine 的函数中。当运行中的被推迟的函数到达了 D 时， D 对 recover 调用的返回值将是传递给 panic 调用的值。如果 D 没有开始一个新的 panic ，正常返回，那么 panicking 序列会停止。在这种情况中，在 G 和 panic 调用之间的函数状态会被丢弃，然后恢复正常的执行。接着会运行被 G 推迟的在 D 前的函数，然后 G 通过返回到它的调用者来终止执行。

如果以下任何条件成立，那么 recover 的返回值为 nil ：

panic 的实参是 nil ；
goroutine 没有 panicking；
recover 没有被一个延迟函数直接调用。
在以下例子中的 protect 函数调用了函数实参 g 并使调用者免受 g 中发生的 run-time panic 之害。

func protect(g func()) {
  defer func() {
    log.Println("done")  // 即使这里有恐慌， Println 也能正常执行
    if x := recover(); x != nil {
      log.Printf("run time panic: %v", x)
    }
  }()
  log.Println("start")
  g()
}
引导
目前的实现提供了一些在引导时有用的内置函数。这些函数已经被记录完整了但是不能保证会一直存在在语言中。它们不会返回一个结果。

函数　      行为

print      打印所有实参；实参的格式化和实现有关
println    和 print 类似，但是会在每个实参间打印空格，在结尾打印新行
实现限制： print 和 println 不需要支持任意的实参类型，但是布尔、数字和字符串 类型 的打印一定要支持。

包
Go 程序是通过连结 包 来构建的。反过来，包由一个或多个源文件构成，这些源文件一起声明属于包的常量、类型、变量和函数，并且可以在同一包的所有文件中访问。这些元素可能被 暴露 并在其它包中使用。

源文件组织
每个源文件都是由以下组成的：定义其所属包的包子句，一组可能为空的用于声明其想要使用内容的包的导入声明，一组可能为空的函数、类型、变量和常量声明。

SourceFile = PackageClause, ";", { ImportDecl ";" }, { TopLevelDecl, ";" } .
包子句Package clause
包子句开始了每个源文件，并定义了文件所属的包。

PackageClause  = "package", PackageName .
PackageName    = identifier .
PackageName 一定不能是 空白标识符 。

package math
共享同一包名的一组文件形成了一个包的实现。实现可能要求一个包的源文件都在同一文件夹下。

导入声明Import declarations
导入声明陈述了这个包含声明的源文件依赖 被导入的 包的功能（ 程序初始化和执行 ）并启用了对该包被 暴露 的标识符的访问。导入命名了一个标识符（包名）用来被访问，以及一个指定被导入包的导入路径。

ImportDecl       = "import", ( ImportSpec | "(", { ImportSpec, ";" }, ")" ) .
ImportSpec       = [ "." | PackageName ], ImportPath .
ImportPath       = string_lit .
PackageName 是用在 限定标识符 中来访问导入源文件中包的暴露标识符的。它是在文件 块 中被声明的。如果 PackageName 缺失，那它默认为被导入包的 包子句 中指定的标识符。如果明确的句号（ . ）取代名字出现了，所有在包的包 块 中声明的包的暴露标识符将在这个导入包的源文件中被声明，并且必须不带限定符来访问。

导入路径的解释是依赖于实现的，但它通常是已编译包的完整文件名的子字符串，并可能和已安装包的库所相对应。

实现限制：编译器可能会限制导入路径仅使用属于 Unicode 的 L, M, N, P 和 S 主类的字符串（无空格的可见字符）到非空字符串，并也可能去除了字符串 !"#$%&'()*,:;<=>?[\]^`{|} 和 Unicode 替换字符 U+FFFD 。

假定我们已经编译了一个包含包子句 package math 的包，它暴露了函数 Sin ，并将编译好的包安装在由 "lib/math" 标记的文件。此表格说明了 Sin 是如何在在各种导入声明后导入包的文件中被访问的。

导入声明　　                  Sin 的本地名

import   "lib/math"         math.Sin
import m "lib/math"         m.Sin
import . "lib/math"         Sin
导入声明声明了导入者和被导入包的依赖关系。在包中直接/间接导入它自己是非法的，直接导入一个没有引用任何其暴露标识符的包也是非法的。仅仅为了包的副作用（初始化）来导入一个包的话，使用 空白 标识符作为明确的包名：

import _ "lib/math"
一个示例包
这里有一个实现并发质数筛选的完整 Go 包。

package main

import "fmt"

// 发送 2, 3, 4, … 序列到信道 'ch'
func generate(ch chan<- int) {
  for i := 2; ; i++ {
    ch <- i  // 发送 'i' 到信道 'ch'
  }
}

// 从信道 'src' 拷贝值到信道 'dst'
// 移除那些可被 'prime' 整除的
func filter(src <-chan int, dst chan<- int, prime int) {
  for i := range src {  // 遍历从 'src' 接收的值
    if i%prime != 0 {
      dst <- i  // 发送 'i' 到信道 'dst'
    }
  }
}

// 质数筛选: 菊花链过滤器一起处理
func sieve() {
  ch := make(chan int)  // 创建一个新的信道
  go generate(ch)       // 启动 generate() 作为子进程
  for {
    prime := <-ch
    fmt.Print(prime, "\n")
    ch1 := make(chan int)
    go filter(ch, ch1, prime)
    ch = ch1
  }
}

func main() {
  sieve()
}
程序初始化和执行
零值
当存储空间被分配给一个 变量 （无论是通过一个声明、对 new 的调用或是新的值被创建，还是通过一个综合的字面值或对 make 的调用）且没有提供明确的初始化时，这个变量或值会被给定一个默认值。这样一个变量或值的每个元素都会被设定到其类型的 零值 ：布尔是 false ，数字类型是 0 ，字符串类型是 "" ，指针、函数、接口、分片、信道和映射类型是 nil 。初始化会被递归地完成，所以打个比方，如果结构数组的元素未指定值，则都将其每个元素字段置零值。

以下两个简单声明是相等：

var i int
var i int = 0
在

type T struct { i int; f float64; next *T }
t := new(T)
后，会有如下的：

t.i == 0
t.f == 0.0
t.next == nil
在完成如下声明后，也会有相同的结果

var t T
包初始化
在一个包内，包级别变量初始化是逐步进行的，每个步骤以 声明顺序 选择不依赖未初始化变量的最早变量。

更精确地说，如果包级别变量还没被初始化且其没有 初始化表达式 或其初始化表达式没有在未声明变量中有依赖，那么它就被认为是 准备好初始化了 。初始化通过重复初始化下一个最早声明且准备好初始化的包级变量来进行，直到没有变量准备好初始化了。

如果在此过程结束时还有变量没初始化，且这些变量是一个或多个初始化循环的一部分，那么程序是无效的。

由在右侧的单个（多值）表达式所初始化的左侧的多个变量是一起被初始化的：如果任意一个在左侧的变量被初始化了，那么这些变量都在同一个步骤被初始化。

var x = a
var a, b = f() // a 和 b 是在 x 被初始化之前一起被初始化的
为了包初始化的目的， 空白 变量会被像其它被描述的变量一样对待。

在多个文件中声明的变量的声明顺序是由对应文件提交给编译器的顺序来决定的：第一个文件中声明的变量会在任何第二个文件中声明的变量之前，以此类推。

依赖关系分析不依赖实际的变量值，仅依赖于源码内的词汇 引用 ，且按照传递轨迹来分析的。例如，如果一个变量 x 的初始化表达式引用了一个其实体引用了变量 y 的函数，那么 x 依赖 y 。具体来说：

到一个变量或函数的引用是表示这个变量或函数的标识符。
到方法 m 的引用是一个 t.m 形式的 方法值 或 方法表达式 ，其中 t 的（静态）类型不能是接口类型，且方法 m 在 t 的方法集中。结果的函数值 t.m 是否被调用是无关紧要的。
如果一个变量、函数或方法 x 的初始化表达式或实体（对于函数和方法而言）包含一个到变量 y 或到依赖于 y 的函数或方法的引用，那么 x 是依赖 y 的。
比如，给定声明

var (
  a = c + b  // == 9
  b = f()    // == 4
  c = f()    // == 5
  d = 3      // 初始化结束后等于 5
)

func f() int {
  d++
  return d
}
初始化顺序是 d, b, c, a 。注意的是，初始化表达式中的子表达式的顺序是无所谓的：示例中 a = c + b 和 a = b + c 得出的是相同的初始化顺序。

依赖分析是分包执行的；只有涉及到在当前包中声明的变量、函数和（非接口）方法的引用才会被考虑。如果变量间存在其它、隐藏的、数据依赖，那么这些变量间的初始化顺序是不明的。

比如，给定声明

var x = I(T{}).ab()   // x 存在在 a 和 b 上的未被发现的隐藏依赖
var _ = sideEffect()  // 与 x, a, 或 b 无关
var a = b
var b = 42

type I interface      { ab() []int }
type T struct{}
func (T) ab() []int   { return []int{a, b} }
变量 a 会在 b 后被初始化，但是 x 是在 b 之前、在 b 和 a 之间、还是在 a 之后，以及 sideEffect() 会在什么时候被调用（在 x 初始化前还是后）都是不明的。

变量也可以被包块中声明的不带实参和结果类型的名为 init 的函数所初始化。

func init() { … }
单一包中可以定义多个这样的函数，甚至是在单一源文件内也没问题。在包块内， init 标识符仅用于声明 init 函数，但标识符本身是未 声明的 。这样的 init 函数不能在程序中的任何位置被引用。

不带导入声明的包是这样初始化的：分配初始化值到它所有的包级变量（按照出现在源码中的顺序，可能会在多个文件中，那就按照提交到编译器的顺序），接着调用:code:init 函数。如果包有导入声明，那么在初始化包本身之前，被导入的包会先初始化好。如果多个包导入了一个包，那么被导入的包只会初始化一次。通过构造可以保证包的导入不存在循环初始化依赖关系。

包的初始化（变量初始化和对 init 函数的调用）在单一 goroutine 内，循序的，每次一个包地发生。 init 函数可能发起其它的可以与初始化代码并行运行的 goroutine。不过，初始化过程总是会序列化 init 函数：在上一个没有返回前不会调用下一个。

为了确保可重现的初始化行为，建议构建系统以词法文件名顺序将属于同一个包的多个文件呈现给编译器。

程序执行
一个完整的程序是通过按轨迹地连接一个单一的，未导入的被叫做 main package 的包与其它所有其导入的包来创建的。主包的包名一定是 main ，并且声明一个无实参也无返回值的 main 函数。

func main() { … }
程序通过先初始化主包再调用 main 函数来开始执行。当这个函数调用返回时，程序退出。并不会等待其它（非 main ） goroutine 完成。

错误
预先声明的类型 error 定义如下：

type error interface {
  Error() string
}
它是表示错误条件的常见接口， nil 值代表没有错误。例如，从文件读入数据的函数可能被定义为：

func Read(f *File, b []byte) (n int, err error)
Run-time panics
像尝试超出数组边界的索引这样的执行错误会触发一个 run-time panic ，它等同于对内置函数 panic 的调用，该调用使用根据实现定义的接口类型 runtime.Error 的值作为实参。这个类型满足预先声明的接口类型 error 。表示不同运行时错误条件的确切错误值是为指定的。

package runtime

type Error interface {
  error
  // 或许还有其它方法
}
系统注意事项
包 unsafe
编译器已知且可以通过 导入路径 "unsafe" 访问的内置包 unsafe 提供了包括违反类型系统操作在内的低级编程设施。使用 unsafe 的包必须必须手动审查以确保类型安全，且不具备可移植性。该包提供了以下接口：

package unsafe

type ArbitraryType int  // 任意 Go 类型的简写； 它不是一个真实的类型
type Pointer *ArbitraryType

func Alignof(variable ArbitraryType) uintptr
func Offsetof(selector ArbitraryType) uintptr
func Sizeof(variable ArbitraryType) uintptr
Pointer 是一个 指针类型 但是 Pointer 值不能被 解引用 。任何指针或 潜在类型 为 uintptr 的值都可以被转换为潜在类型为 Pointer 的类型，反之亦然。在 Pointer 和 uintptr 间的转换效果是由实现定义的。

var f float64
bits = *(*uint64)(unsafe.Pointer(&f))

type ptr unsafe.Pointer
bits = *(*uint64)(ptr(&f))

var p ptr = nil
函数 Alignof 和 Sizeof 获取任意类型的表达式 x 并分别返回假象变量 v 的定位或大小（ v 就像通过 var v = x 声明的）。

函数 Offsetof 获取一个（可能被括起来的）表示被 s 或 *s 所表示的结构体的字段 f 的 选择器 s.f ，并返回相对于结构体地址的以字节表示的字段的偏移量。如果 f 是一个 嵌入字段 ，那么必须可以在不要指针间接情况下直达结构体字段。对于带字段 f 的结构体 s ：

uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.f) == uintptr(unsafe.Pointer(&s.f))
计算机架构可能会要求内存地址是 对其的 ；也就是说，变量的地址是一个因子的倍数，这个因子是变量类型的 对准值alignment 。函数 Alignof 获取一个表示任意类型变量的表达式，并以字节为单位返回变量（的类型）的对准值。对于一个变量 x ：

uintptr(unsafe.Pointer(&x)) % unsafe.Alignof(x) == 0
对于 Alignof 、 Offsetof 、 Sizeof 的调用是类型 uintptr 的编译时常量表达式。

大小和对准值保证
对于 数字类型 ，以下大小是保证的：

类型　                                以字节为单位的大小

byte, uint8, int8                     1
uint16, int16                         2
uint32, int32, float32                4
uint64, int64, float64, complex64     8
complex128                           16
以下最小对准值属性是保证的：

对于任意类型变量 x ： unsafe.Alignof(x) 最小为 1。
对于结构体类型变量 x ： unsafe.Alignof(x) 是所有 unsafe.Alignof(x.f) （对于 x 的每个字段 f ）中最大的值，但最小为1。
对于数组类型变量 x ： unsafe.Alignof(x) 和数组元素类型变量的对准值相同。
译注：这边我一开始很纠结为什么 complex128 类型的对准值是 8 字节，后来发现 complex64 的对准值是 4 字节，所以大胆猜测它是拆开来算的
如果结构体或数组没有包含大于零大小的字段（或元素，对数组而言），那么它大大小为零。两个不同的零大小的变量在内存中可能拥有同一个地址。
