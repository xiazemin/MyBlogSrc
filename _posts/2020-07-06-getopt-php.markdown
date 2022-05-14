---
title: getopt-php
layout: post
category: php
author: 夏泽民
---
https://github.com/getopt-php/getopt-php
https://www.php.net/manual/zh/function.getopt.php
<!-- more -->
php中的getop是用于接收cmd参数的时候用的
例如当你再linxu 中 需要用php调试的时候，往往需要带参数调试
getopt就是可以实现带参数传入的命令

使用方式：<br>
　　　　array getopt ( string $options [, array $longopts ] )
注意：　　$options字符串中的每个字符将被用来作为选项字符和对传递给脚本用一个连字符开始匹配选项（ - ）。例如，一个选项字符“x”对应一个选项-x。只有a - z，A - Z和0-9是允许的　　空格是不能作为选项字符的。
note: 包含当运行于命令行下时传递给当前脚本的参数的数组。

Note: 这个变量仅在 register_argc_argv 打开时可用。


{% raw %}
php script.php -f "value for f" -v -a --required value --optional="optional value" --option will output:
输出：
array(6) {
  ["f"]=>
  string(11) "value for f"
  ["v"]=>
  bool(false)
  ["a"]=>
  bool(false)
  ["required"]=>
  string(5) "value"
  ["optional"]=>
  string(14) "optional value"
  ["option"]=>
  bool(false)
}
{% endraw %}

https://www.gnu.org/savannah-checkouts/gnu/libc/manual/html_node/Getopt.html


使用通常配合php脚本使用
#!/usr/bin/env php
<?php
use GetOpt\Command;
use GetOpt\GetOpt;

$getOpt = new GetOpt();
$getOpt->addOptions([
    \GetOpt\Option::create(null, 'version')->setDescription('Show version'),
]);

$getOpt->addCommands([
    Command::create('build', '\build\\Handler::handle')
        ->setDescription('构建')
        ->addOptions(BuildCliArg::options())
        ->addOperands(BuildCliArg::operands()),
        ....
        ]);
        
 $getOpt->process();
 
 
 <?php
\build\\Handler::handle.php
public static function handle(GetOpt $getOpt) {
   exec();
}

