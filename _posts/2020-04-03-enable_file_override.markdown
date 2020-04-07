---
title: enable_file_override
layout: post
category: php
author: 夏泽民
---
opcache.enable_file_override boolean

如果启用，则在调用函数 file_exists()， is_file() 以及 is_readable() 的时候， 都会检查操作码缓存，无论文件是否已经被缓存。 如果应用中包含检查 PHP 脚本存在性和可读性的功能，这样可以提升性能。 但是如果禁用了 opcache.validate_timestamps 选项， 可能存在返回过时数据的风险。
<!-- more -->
PHP执行这段代码会经过如下4个步骤(确切的来说，应该是PHP的语言引擎Zend)

1）Scanning(Lexing) ,将PHP代码转换为语言片段(Tokens)。

2）Parsing, 将Tokens转换成简单而有意义的表达式。

3）Compilation, 将表达式编译成Opocdes。

4）Execution, 顺次执行Opcodes，每次一条，从而实现PHP脚本的功能。
PHP opcache介绍
Optimizer+（Optimizer+于2013年3月中旬改名为Opcache），OPcache通过将PHP脚本预编译的字节码存储到共享内存中来提升PHP的性能，存储预编译字节码的好处就是省去了每次加载和解析PHP脚本的开销。

PHP 5.5.0 及后续版本中已经绑定了 OPcache 扩展。 对于 PHP 5.2，5.3 和 5.4 版本可以使用 » PECL扩展中的OPcache库。

PHP 5.5.0及后续版本

OPcache只能编译为共享扩展。如果你使用–disable-all参数禁用了默认扩展的构建，那么必须使用–enable-opcache选项来开启OPcache。编译之后，就可以使用 zend_extension 指令来将 OPcache 扩展加载到 PHP 中。

推荐的php.ini设置

使用下列推荐设置来获得较好的性能：

opcache.memory_consumption=128
opcache.interned_strings_buffer=8
opcache.max_accelerated_files=4000
opcache.revalidate_freq=60
opcache.fast_shutdown=1
opcache.enable_cli=1
opcache.save_comments=0
你也可以禁用 opcache.save_comments 并且启用 opcache.enable_file_override。 需要提醒的是，在生产环境中使用上述配置之前，必须经过严格测试。 因为上述配置存在一个已知问题，它会引发一些框架和应用的异常， 尤其是在存在文档使用了备注注解的时候。

以下是opcache的配置说明，其中给有值得都是默认配置：

; opcache的开关,关闭时代码不再优化.
opcache.enable=1
 
; Determines if Zend OPCache is enabled for the CLI version of PHP
opcache.enable_cli=1
 
; OPcache的共享内存大小，以兆字节为单位。总共能够存储多少预编译的PHP代码(单位:MB)
; 推荐128
opcache.memory_consumption=64
 
; 用来存储临时字符串的内存大小，以兆字节为单位.
; 推荐8
opcache.interned_strings_buffer=4
 
; 最大缓存的文件数目200到100000之间.
; 推荐4000
opcache.max_accelerated_files=2000
 
; 内存"浪费"达到此值对应的百分比,就会发起一个重启调度.
opcache.max_wasted_percentage=5
 
; 开启这条指令, Zend Optimizer + 会自动将当前工作目录的名字追加到脚本键上,以此消除同名文件间的键值命名冲突.关闭这条指令会提升性能,但是会对已存在的应用造成破坏.
opcache.use_cwd=0
 
; 开启文件时间戳验证
opcache.validate_timestamps=1
 
; 检查脚本时间戳是否有更新的周期，以秒为单位。设置为0会导致针对每个请求，OPcache都会检查脚本更新.
; 推荐60
opcache.revalidate_freq=2
 
; 允许或禁止在include_path中进行文件搜索的优化.
opcache.revalidate_path=0
 
; 如果禁用，脚本文件中的注释内容将不会被包含到操作码缓存文件，这样可以有效减小优化后的文件体积,禁用此配置指令可能会导致一些依赖注释或注解的应用或框架无法正常工作，比如:Doctrine,Zend Framework2等.
; 推荐0
opcache.save_comments=1
 
; 如果禁用，则即使文件中包含注释，也不会加载这些注释内容。本选项可以和opcache.save_comments一起使用，以实现按需加载注释内容.
opcache.load_comments=1

; 打开快速关闭,打开这个在PHP Request Shutdown的时候会收内存的速度会提高.
; 推荐1
opcache.fast_shutdown=1
 
; 允许覆盖文件存在（file_exists等）的优化特性.
opcache.enable_file_override=0 
 
; 定义启动多少个优化过程.
opcache.optimization_level=0xffffffff
 
; 启用此Hack可以暂时性的解决"can’t redeclare class"错误.
opcache.inherited_hack=1
 
; 启用此Hack可以暂时性的解决"can’t redeclare class"错误.
;opcache.dups_fix=0
 
; 通过文件大小屏除大文件的缓存，默认情况下所有的文件都会被缓存.
;opcache.max_file_size=0
 
; 每N次请求检查一次缓存校验.默认值0表示检查被禁用了,由于计算校验值有损性能,这个指令应当紧紧在开发调试的时候开启.
;opcache.consistency_checks=0
 
; 从缓存不被访问后,等待多久后(单位为秒)调度重启.
;opcache.force_restart_timeout=180
 
; 日志记录level，默认只有fatal error和error.
;opcache.error_log=
 
; 将错误信息写入到服务器(Apache等)日志
;opcache.log_verbosity_level=1
 
; 内存共享的首选后台.留空则是让系统选择.
;opcache.preferred_memory_model=
 
; 运行php脚本时保护共享内存防止意外的写入,只对debug时有用.
;opcache.protect_memory=0
最后说一下使用opcache加速php时应该注意的坑：

opcache依靠的是PHP文件的modify time作为文件被修改的检测条件，基于这个会引发两个问题。

第一个问题是做版本回滚时，由于版本回滚后的文件修改时间比现有opcache缓存的文件时间要往前一些，所以可能会导致opcache不会清除缓存，需要手动reload。

第二个问题是做版本发布时，一般都是sync方式，可能会出现文件发布一半时被opcache缓存，用户访问会报程序错误，这个主要是因为文件内容缓存了一半，但是文件的时间戳不会在改变，所以就算opcache检测时也不会去读取新的文件了，需要手动reload。

针对这两个问题，不光reload可以解决，同样调用opcache的接口也可以清除opcache缓存。

你可以使用opcache_reset()或者或者opcache_invalidate()函数来手动重置OPcache。

opcache_reset()：该函数将重置整个字节码缓存，在调用opcache_reset()之后，所有的脚本将会重新载入并且在下次被点击的时候重新解析。

opcache_invalidate()：该函数的作用是使得指定脚本的字节码缓存失效。 如果force没有设置或者传入的是FALSE，那么只有当脚本的修改时间 比对应字节码的时间更新，脚本的缓存才会失效。

但是不推荐使用，个人在生产环境中进行代码发布后调用opcache_reset()清空缓存（测试确实可以清空缓存），出现过奇葩问题（访问量大的应用），后来就果断放弃了，使用了reload的方式。

Opcache优化在著名的《modern php》 中也有重要篇幅。在PHP文档也有详细介绍：http://php.net/manual/zh/opcache.configuration.php#ini.opcache.revalidate-freq

个人觉得这种文章相当有指导意义，所以特地把它的设置方式摘译如下（格式有些修改）。

opcache.revalidate_freq
这个选项用于设置缓存的过期时间（单位是秒），当这个时间达到后，opcache会检查你的代码是否改变，如果改变了PHP会重新编译它，生成新的opcode，并且更新缓存。值为“0”表示每次请求都会检查你的PHP代码是否更新（这意味着会增加很多次stat系统调用，译注：stat系统调用是读取文件的状态，这里主要是获取最近修改时间，这个系统调用会发生磁盘I/O，所以必然会消耗一些CPU时间，当然系统调用本身也会消耗一些CPU时间）。可以在开发环境中把它设置为0，生产环境下不用管，因为下面会介绍另外一个设置选项。

opcache.validate_timestamps
当这个选项被启用（设置为1），PHP会在opcache.revalidate_freq设置的时间到达后检测文件的时间戳（timestamp）。

如果这个选项被禁用（设置为0），opcache.revalidate_freq会被忽略，PHP文件永远不会被检查。这意味着如果你修改了你的代码，然后你把它更新到服务器上，再在浏览器上请求更新的代码对应的功能，你会看不到更新的效果，你必须得重新加载你的PHP（使用kill -SIGUSR2强制重新加载）。

这个设定是不是有些蛋疼，但是我强烈建议你在生产环境中使用，why？因为当你在更新服务器代码的时候，如果代码较多，更新操作是有些延迟的，在这个延迟的过程中必然出现老代码和新代码混合的情况，这个时候对用户请求的处理必然存在不确定性。

opcache.max_accelerated_files
这个选项用于控制内存中最多可以缓存多少个PHP文件。这个选项必须得设置得足够大，大于你的项目中的所有PHP文件的总和。我的代码库大概有6000个PHP文件，所以我把这个值设置为一个素数7963。

真实的取值是在质数集合 { 223, 463, 983, 1979, 3907, 7963, 16229, 32531, 65407, 130987 } 中找到的第一个比设置值大的质数。 设置值取值范围最小值是 200，最大值在 PHP 5.5.6 之前是 100000，PHP 5.5.6 及之后是 1000000。
听起来好复杂，但用下面的命令就妥啦

你可以运行 find . -type f -print | grep php | wc -l 这个命令来快速计算你的代码库中的PHP文件数。

opcache.memory_consumption
这个选项的默认值为64MB，我把它设置为192MB，因为我的代码很大。你可以通过调用opcachegetstatus()来获取opcache使用的内存的总量，如果这个值很大，你可以把这个选项设置得更大一些。

opcache.interned_strings_buffer
这是一个很有用的选项，但是似乎完全没有文档说明。PHP使用了一种叫做字符串驻留（string interning）的技术来改善性能。例如，如果你在代码中使用了1000次字符串“foobar”，在PHP内部只会在第一使用这个字符串的时候分配一个不可变的内存区域来存储这个字符串，其他的999次使用都会直接指向这个内存区域。这个选项则会把这个特性提升一个层次——默认情况下这个不可变的内存区域只会存在于单个php-fpm的进程中，如果设置了这个选项，那么它将会在所有的php-fpm进程中共享。在比较大的应用中，这可以非常有效地节约内存，提高应用的性能。

这个选项的值是以兆字节（megabytes）作为单位，如果把它设置为16，则表示16MB，默认是4MB，这是一个比较低的值。

opcache.fast_shutdown
另外一个很有用但也没有文档说明的选项。从字面上理解就是“允许更快速关闭”。它的作用是在单个请求结束时提供一种更快速的机制来调用代码中的析构器，从而加快PHP的响应速度和PHP进程资源的回收速度，这样应用程序可以更快速地响应下一个请求。把它设置为1就可以使用这个机制了。

最终我们对于opcache在php.ini的设置如下：

开发模式下推荐，直接禁用opcache扩展更好

opcache.revalidate_freq=0
opcache.validate_timestamps=1
opcache.max_accelerated_files=3000
opcache.memory_consumption=192
opcache.interned_strings_buffer=16
opcache.fast_shutdown=1
多台机器集群模式或者代码更新频繁时推荐，可以兼顾性能，方便代码更新

opcache.revalidate_freq=300
opcache.validate_timestamps=1
opcache.max_accelerated_files=7963
opcache.memory_consumption=192
opcache.interned_strings_buffer=16
opcache.fast_shutdown=1
稳定项目推荐，性能最好

opcache.revalidate_freq=0
opcache.validate_timestamps=0
opcache.max_accelerated_files=7963
opcache.memory_consumption=192
opcache.interned_strings_buffer=16
opcache.fast_shutdown=1

https://www.php.net/manual/zh/opcache.configuration.php#ini.opcache.revalidate-freq
