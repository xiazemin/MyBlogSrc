---
title: swagger php
layout: post
category: php
author: 夏泽民
---
在php文件中写 swagger 格式的 /* 注释 /
用 swagger-php 内的 bin/swagger.phar 命令扫描 php controller 所在目录, 生成 swagger.json 文件
将 swagger.json 文件拷贝到 swagger-ui 中 index.html 指定的目录中
打开 swagger-ui 所在的 url, 就可以看到文档了. 文档中的各个 api 可以在该网址上直接访问得到数据.
实现此需求只需要 swagger 的如下两个项目:
swagger-php: 扫描 php 注释的工具. 内含一个不错的例子.
swagger-ui: 用以将扫描工具生成的 swagger.json 文件内容展示在网页上.
<!-- more -->
$ git clone https://github.com/swagger-api/swagger-ui.git
$ git clone https://github.com/zircote/swagger-php.git

说是部署,主要就是产生 bin/swagger 这个用来生成 swagger.json 文件的命令.
主要工作,就是用 composer 解决下依赖就可以了.
因为国内直接用 composer 比较蛋疼,所以最好设置下国内的那个 composer 源.
这样的话, 整个 文档生成工具的部署 就是下面三行命令:

$ cd swagger-php
$ composer config repo.packagist composer https://packagist.phpcomposer.com
$ composer update
只要中间不报错,就算部署完成了. 完成后可以生成一份文档试一下.
swagger-php 项目下的 Examples 目录下有一个示例php工程,里面已经用 swagger 格式写了各种接口注释, 我们来尝试生成一份文档.
执行下面命令:

$ cd swagger-php
$ mkdir json_docs
$ php ./bin/swagger ./Examples -o json_docs/
上面命令会扫描 Examples 目录中的php文件注释, 然后在 json_docs 目录下生成 swagger.json 文件.这个 swagger.json 文件就是前端 swagger-ui 用来展示的我们的api文档文件.

NOTE: swagger-php 只是个工具,放在哪里都可以.

前端 swagger-ui 部署:
部署方法很简单,就三步:

1. 将 swagger-ui 项目中的 dist 文件夹拷贝到 php_rest_api 根目录下.
NOTE1: 只需要拷贝dist这一个文件夹就可以了.最好重命名下,简单起见,这里不再重命名.
NOTE2: 我们的项目根目录和 nginx 配置的 root 是同一个目录.其实不用放跟目录,只要放到一个不用跨域就跨域访问的目录就可以了. 为啥有跨域问题? 后面会讲.

2. 修改 dist 文件夹下的 index.html 文件,指定 swagger.json 所在目录
只改一行就可以.
简单起见,这里直接将 swagger.json 目录指定在 dist 目录下即可. 我们这里屡一下预设条件:
假设 php_api_project 项目的 host 是 api.my_project.com;
假设 php_api_project 项目在 nginx 中指定的 root 即为其根目录;
假设 swagger-ui 里的 dist 文件夹放在上述根目录中;
假设 swagger.json 文件就打算放在上述 dist 目录下 (php_api_project/dist/swagger.json) ;
那么 index.html 中把下面的片段改成这样:

      var url = window.location.search.match(/url=([^&]+)/);
      if (url && url.length > 1) {
        url = decodeURIComponent(url[1]);
      } else {
        <!-- 就是这行,改成你生成的 swagger.json 可以被访问到的路径即可 -->
        url = "http://api.my_project.com/dist/swagger.json";
      }
3. 拷贝 swagger.json 到上述目录中.
# 把 swagger-php_dir 这个,换成你的 swagger-php 录即可
cp swagger-php_dir/json_docs/swagger.json php_api_project/dist/
上述步骤完成后, 访问 api.my_project.com/dis... 就可以看到 Examples 那个小项目的 api 文档了.

编写 PHP 注释
swagger-php 项目的 Example 中已经有了很多相关例子,照着复制粘贴就可以了.
更具体的相关注释规则的文档,看这里:
bfanger.nl/swagger-exp…

假设我的项目 controller 所在目录为 php_api_project/controller/, 那么我只需要扫描这个目录就可以了,不用扫描整个 php 工程.

为了在 swagger.json 中生成某些统一的配置, 建立 php_api_project/controller/swagger 目录. 目录存放一个没有代码的php文件,里面只写注释.

我给这个文件取名叫 Swagger.php, 大体内容如下:

<?php
 
/**
 * @SWG\Swagger(
 *   schemes={"http"},
 *   host="api.my_project.com",
 *   consumes={"multipart/form-data"},
 *   produces={"application/json"},
 *   @SWG\Info(
 *     version="2.3",
 *     title="my project doc",
 *     description="my project 接口文档, V2-3.<br>
以后大家就在这里愉快的对接口把!<br>
以后大家就在这里愉快的对接口把!<br>
以后大家就在这里愉快的对接口把!<br>
"
 *   ),
 *
 *   @SWG\Tag(
 *     name="User",
 *     description="用户操作",
 *   ),
 *
 *   @SWG\Tag(
 *     name="MainPage",
 *     description="首页模块",
 *   ),
 *
 *   @SWG\Tag(
 *     name="News",
 *     description="新闻资讯",
 *   ),
 *
 *   @SWG\Tag(
 *     name="Misc",
 *     description="其他接口",
 *   ),
 * )
 */
如上所示,我的这个php文件一行php代码也没有,就只有注释,为了定义一些全局的swagger设置:

schemes: 使用协议 (可以填多种协议)
host: 项目地址, 这个地址会作为每个接口的 url base ,拼接起来一期作为访问地址
consumes: 接口默认接收的MIME类型, 我的例子中的 formData 对应post表单类型. 注意这是项目默认值,在单个接口注释里可以复写这个值.
produces: 接口默认的回复MIME类型. api接口用的比较多的就是 application/json 和 application/xml.
@SWG\Info: 这个里面填写的东西,会放在文档的最开头,用作文档说明.
@SWG\Tag: tag是用来给文档分类的,name字段必须唯一.某个接口可以指定多个tag,那它就会出现在多组分类中. tag也可以不用在这里预先定义就可以使用,但那样就没有描述了. 多说无益,稍微用用就啥都明白了.

然后就是给每个接口编写 swagger 格式的注释了

