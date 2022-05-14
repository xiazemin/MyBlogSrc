---
title: composer dump-autoload
layout: post
category: php
author: 夏泽民
---
composer 提供的 autoload 机制使得我们组织代码和引入新类库非常方便，但是也使项目的性能下降了不少 。

composer autoload 慢的主要原因在于来自对 PSR-0 和 PSR-4 的支持，加载器得到一个类名时需要到文件系统里查找对应的类文件位置，这导致了很大的性能损耗，当然这在我们开发时还是有用的，这样我们添加的新的类文件就能即时生效。 但是在生产模式下，我们想要最快的找到这些类文件，并加载他们。

因此 composer 提供了几种优化策略，下面说明下这些优化策略。

第一层级(Level-1)优化： 生成 classmap
如何运行：
执行命令 composer dump-autoload -o （-o 等同于 --optimize）

原理：
这个命令的本质是将 PSR-4/PSR-0 的规则转化为了 classmap 的规则， 因为 classmap 中包含了所有类名与类文件路径的对应关系，所以加载器不再需要到文件系统中查找文件了。可以从 classmap 中直接找到类文件的路径。

注意事项
建议开启 opcache , 这样会极大的加速类的加载。
php5.5 以后的版本中默认自带了 opcache 。

这个命令并没有考虑到当在 classmap 中找不到目标类时的情况，当加载器找不到目标类时，仍旧会根据PSR-4/PSR-0 的规则去文件系统中查找

第二层级(Level-2/A)优化：权威的（Authoritative）classmap
执行命令：
执行命令 composer dump-autoload -a （-a 等同于 --classmap-authoritative）

原理
执行这个命令隐含的也执行了 Level-1 的命令， 即同样也是生成了 classmap，区别在于当加载器在 classmap 中找不到目标类时，不会再去文件系统中查找（即隐含的认为 classmap 中就是所有合法的类，不会有其他的类了，除非法调用）

注意事项
如果你的项目在运行时会生成类，使用这个优化策略会找不到这些新生成的类。

第二层级(Level-2/B)优化：使用 APCu cache
执行命令：
执行命令 composer dump-autoload --apcu

原理：
使用这个策略需要安装 apcu 扩展。
apcu 可以理解为一块内存，并且可以在多进程中共享。
这种策略是为了在 Level-1 中 classmap 中找不到目标类时，将在文件系统中找到的结果存储到共享内存中， 当下次再查找时就可以从内存中直接返回，不用再去文件系统中再次查找。

在生产环境下，这个策略一般也会与 Level-1 一起使用， 执行composer dump-autoload -o --apcu, 这样，即使生产环境下生成了新的类，只需要文件系统中查找一次即可被缓存 ， 弥补了Level-2/A 的缺陷。

如何选择优化策略？
要根据自己项目的实际情况来选择策略，如果你的项目在运行时不会生成类文件并且需要 composer 的 autoload 去加载，那么使用 Level-2/A 即可，否则使用 Level-1 及 Level-2/B是比较好的选择。

几个提示
Level-2的优化基本都是 Level-1 优化的补充，Level-2/A 主要是决定在 classmap 中找不到目标类时是否继续找下去的问题，Level-2/B 主要是在提供了一个缓存机制，将在 classmap 中找不到时，将从文件系统中找到的文件路径缓存起来，加速后续查找的速度。
在执行了 Level-2/A 时，表示在 classmap 中找不到不会继续找，此时 Level-2/B 是不会生效的。
不论那种情况都建议要开启 opcache， 这会极大的提高类的加载速度，我目测有性能提升至少 10倍。
<!-- more -->
https://blog.csdn.net/zhouyuqi1/article/details/81098650

Composer 作为现代 phper 的春天，远离重复造轮子的时代，大部分扩展包遵循 psr-4 规范，使得扩展更加轻松，减轻了工作的部分压力

这篇文章来说一下为什么在生产环境下使用 Composer 加载包后要再使用 dumpautoload 呢？

composer dump-autoload (-o)
composer dumpautoload (-o)
这个就要看一下 vendor/composer 目录下的文件了，先看一下 autoload_real.php

类名为 ComposerAutoloaderInit440563a888dcb3a8c02b3ef8400e84e8，ComposerAutoloaderInit 后为一段 hash 值

这个也是为了避免命名冲突，每次 composer install 都会生成不一样的值

再往下看 getLoader 这个方法

public static function getLoader()
{
    if (null !== self::$loader) {
        return self::$loader;
    }

    spl_autoload_register(array('ComposerAutoloaderInit440563a888dcb3a8c02b3ef8400e84e8', 'loadClassLoader'), true, true);
    self::$loader = $loader = new \Composer\Autoload\ClassLoader();
    spl_autoload_unregister(array('ComposerAutoloaderInit440563a888dcb3a8c02b3ef8400e84e8', 'loadClassLoader'));

    $useStaticLoader = PHP_VERSION_ID >= 50600 && !defined('HHVM_VERSION') && (!function_exists('zend_loader_file_encoded') || !zend_loader_file_encoded());
    if ($useStaticLoader) {
        require_once __DIR__ . '/autoload_static.php';

        call_user_func(\Composer\Autoload\ComposerStaticInit440563a888dcb3a8c02b3ef8400e84e8::getInitializer($loader));
    } else {
        $map = require __DIR__ . '/autoload_namespaces.php';
        foreach ($map as $namespace => $path) {
            $loader->set($namespace, $path);
        }

        $map = require __DIR__ . '/autoload_psr4.php';
        foreach ($map as $namespace => $path) {
            $loader->setPsr4($namespace, $path);
        }

        $classMap = require __DIR__ . '/autoload_classmap.php';
        if ($classMap) {
            $loader->addClassMap($classMap);
        }
    }

    $loader->register(true);

    return $loader;
}
这个方法是获取 Composer\ClassLoader，如果不存在就是生成一个实例放在 ComposerAutoloaderInit440563a888dcb3a8c02b3ef8400e84e8 中

将 Composer 生成的各种 autoload_psr4、autoload_classmap、autoload_namespaces 全都注册到 Composer\ClassLoader 中

然后 register 注册文件

了解了 autoload.php 是如何工作的，以后那么我们看一下 composer dump-atoload -o 有什么用

autoload_classmap.php 在未执行命名之前 return 了一个空数组

在执行之后会发现所有的扩展包类的 namespace 和 classname 生成成一个 key => value 的数组

这时我们需要分析一下 ClassLoader 这个类的源码

private function findFileWithExtension($class, $ext)
{
    // PSR-4 lookup
    $logicalPathPsr4 = strtr($class, '\\', DIRECTORY_SEPARATOR) . $ext;

    $first = $class[0];
    if (isset($this->prefixLengthsPsr4[$first])) {
        $subPath = $class;
        while (false !== $lastPos = strrpos($subPath, '\\')) {
            $subPath = substr($subPath, 0, $lastPos);
            $search = $subPath . '\\';
            if (isset($this->prefixDirsPsr4[$search])) {
                $pathEnd = DIRECTORY_SEPARATOR . substr($logicalPathPsr4, $lastPos + 1);
                foreach ($this->prefixDirsPsr4[$search] as $dir) {
                    if (file_exists($file = $dir . $pathEnd)) {
                        return $file;
                    }
                }
            }
        }
    }

    // PSR-4 fallback dirs
    foreach ($this->fallbackDirsPsr4 as $dir) {
        if (file_exists($file = $dir . DIRECTORY_SEPARATOR . $logicalPathPsr4)) {
            return $file;
        }
    }

    // PSR-0 lookup
    if (false !== $pos = strrpos($class, '\\')) {
        // namespaced class name
        $logicalPathPsr0 = substr($logicalPathPsr4, 0, $pos + 1)
            . strtr(substr($logicalPathPsr4, $pos + 1), '_', DIRECTORY_SEPARATOR);
    } else {
        // PEAR-like class name
        $logicalPathPsr0 = strtr($class, '_', DIRECTORY_SEPARATOR) . $ext;
    }

    if (isset($this->prefixesPsr0[$first])) {
        foreach ($this->prefixesPsr0[$first] as $prefix => $dirs) {
            if (0 === strpos($class, $prefix)) {
                foreach ($dirs as $dir) {
                    if (file_exists($file = $dir . DIRECTORY_SEPARATOR . $logicalPathPsr0)) {
                        return $file;
                    }
                }
            }
        }
    }

    // PSR-0 fallback dirs
    foreach ($this->fallbackDirsPsr0 as $dir) {
        if (file_exists($file = $dir . DIRECTORY_SEPARATOR . $logicalPathPsr0)) {
            return $file;
        }
    }

    // PSR-0 include paths.
    if ($this->useIncludePath && $file = stream_resolve_include_path($logicalPathPsr0)) {
        return $file;
    }

    return false;
}
我们可以看到会先去查找 autoload_classmap 中所有生成的注册类，如果没有才会加载 psr-4 和 psr-0

所以使用 dumpautoload 后会优先加载需要的类并提前返回，不然的话 compoesr 只能去动态读取 psr-4 和 prs-0 的内容，这样大大减少了 IO 操作和深层次的循环，提升部分性能问题

https://qq52o.me/2688.html
除了psr-4的自动加载规则，其他的需要执行这个命令，composer才能自动加载到

可以看到 psr-4 或者 psr-0 的自动加载都是一件很累人的事儿。基本是个 O(n2) 的复杂度。另外有一大堆 is_file之类的 IO 操作所以性能堪忧

所以给出的解决方案就是空间换时间

Compsoer\ClassLoader 会优先查看 autoload_classmap 中所有生成的注册类。如果在classmap 中没有发现再 fallback 到 psr-4 然后 psr-0

所以当打了 composer dump-autoload -o 之后，composer 就会提前加载需要的类并提前返回。这样大大减少了 IO 和深层次的 loop

（composer依据项目，会生成composer文件夹）

