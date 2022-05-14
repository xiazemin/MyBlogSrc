---
title: composer replace
layout: post
category: php
author: 夏泽民
---
原始解释：
“Lists packages that are replaced by this package. This allows you to fork a package, publish it under a different name with its own version numbers, while packages requiring the original package continue to work with your fork because it replaces the original package.”

列出被当前包替换的包。这允许你fork一个包，并以不同的名称和它自己的版本号进行发布，由于替换了原始的包，当包依赖原始包时可以使用fork包继续工作。

实际上，第一句话已经解释清楚了，就是当前包替换哪个（或哪些）包。比如：

1
2
3
4
5
6
{
“name": "b/b",
"replace": {
    "a/a": "x.y.z"
}
}
就是b/b这个包替换了a/a这个对应了x.y.z版本的包。如果b/b包中的replace列出了多个包，那么它将替换多个包。在Composer解决依赖的过程中，如果遇到了replace，那么原始的包会被移除，所以它的应用场景非常清晰，就是原始包停止维护了，但是你的项目或包直接或间接依赖一个原始的包，为了可以修复原始包的bug，你可以fork这个包，并在这个包中声明要替换原始包，然后在你的应用或包中依赖这个你fork的包，这样就可以全局替换原始包（原始包被移除）。

由于是包替换，很明显也可能引入安全问题，比如某个包被替换，替换的包中包含恶意代码。

关于replace的第二段解释：
”This is also useful for packages that contain sub-packages, for example the main symfony/symfony package contains all the Symfony Components which are also available as individual packages. If you require the main package it will automatically fulfill any requirement of one of the individual components, since it replaces them.“

这里描述的就是现代框架的组织方式。一个框架（或一个大组件），是很多子包组成的，每个子包都可以单独使用，通常每个子包都会以只读的方式fork到另一个包名，当单独使用时，可以依赖这个只读的包，当依赖整个框架时，框架中声明了替换，则可以移除单独依赖时下载的包。

比如laravel/framework是由很多包组成的，每个包都以只读的方式fork到另一个包名称：

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
// https://github.com/laravel/framework
"replace": {
        "illuminate/auth": "self.version",
        "illuminate/broadcasting": "self.version",
        "illuminate/bus": "self.version",
        "illuminate/cache": "self.version",
        "illuminate/config": "self.version",
        "illuminate/console": "self.version",
        "illuminate/container": "self.version"
}
这个包包括了子包illuminate/container，而它被只读方式fork了一份到https://github.com/illuminate/container，所以当要仅用这个组件时，可以composer require illuminate/container即可。在依赖了illuminate/container后，如何由依赖laravel/framework，这个时候会引起命名冲突，但是laravel/framework中声明了illuminate/container这个包被laravel/framework替换，所以在依赖laravel/framework后，illuminate/container独立的下载将被移除，从而解决了名称冲突。

第一个用法，提供了一个全局包替换的机会。第二个用法，提供了一个现代化的框架组织方式。
<!-- more -->
http://blog.ifeeline.com/2695.html
https://www.runoob.com/w3cnote/composer-install-and-usage.html

https://www.phpcomposer.com/what-is-composer/

http://json-schema.org/

https://docs.phpcomposer.com/04-schema.html

https://www.jianshu.com/p/4d989f3bd2cd
