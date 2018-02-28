---
title: Jekyll目录结构和运行机理
layout: post
category: jekyll
author: 夏泽民
---
https://github.com/xiazemin/jekyll-paginate-plugin
<!-- more -->
Jekyll使用Ruby脚本根据模板生成静态网页，实现了内容与排版的分离。模板以嵌入Liquid脚本的HTML格式存放。内容为markdown或者html。

正常的Jekyll工程包含以下几个目录：

_posts  博客内容
_pages  其他需要生成的网页，如About页
_layouts 网页排版模板
_includes 被模板包含的HTML片段，可在_config.yml中修改位置
assets 辅助资源 css布局 js脚本 图片等
_data 动态数据
_sites  最终生成的静态网页

在项目文件夹(含有_config.yml)中运行jekyll build 指令后，jekyll会依次做如下几件事

加载_layout文件夹内的所有模板，并将其中的%include html% 字段按照_includes文件夹内对应文件填入
遍历_post文件夹及子文件夹，对所有命名符合yyyy-mm--dd-title.md 格式的博客文件放入site.posts 变量(按时间倒序)，并对其进行解析，根据Front Matter 头的内容套入layout生成对应title的博客
遍历整个项目子目录，扫描所有含Front Matter 头的文件，放入site.pages 变量并根据permalink 字段指定的URL目标位置生成index.html
在生成过程中，文件中包含的Liquid脚本{{ content }} 会被解析并替换。Liquid指令 包括Object、Tag、Filter三类，其中object是变量，在解析过程中会被直接文本替换
Tag是控制流，可以做判断和循环
Filter用于对文本进一步处理
.yml文件中，字段的冒号后面必须有空格！

为了添加除博客以外的页面集合(如Projects)，可将含有Front Matter头的文件放入除_posts之外的任意目录，便可被添入site.pages 变量中。为了与posts相区分，一般来说应该在网页头部添加一个变量，如type 并在对应生成循环中逐个判断。

首先说说jeykll , 它是一个静态博客系统，你也可以把jekyll当作是一个工具，它可以将特定格式，如markdown, 或者texttile语法格式编辑的文本文件直接转换成html , 当作网页显示。大家都知道使用markdown等语法来编辑发布博客是一件很轻松愉快的事情，比起直接写html或者jsp，自然轻松许多。 jeykll系统的运行依赖ruby运行时，说明jekyll是使用ruby语言开发的。jekyll 可以在本地环境下安装，使用jekyll --server可以在本地启动一个WebRick的HTTP服务器，浏览器访问localhost:4000便可以预览博客。更多关于jekyll系统生成的文件目录以及每个文件夹里文件的作用，可以查看jeykll的官方文档。其次我们再来说说octopress, 看了octopress的官方文档之后，你会发现他的目录结构和jekyll大同小异。这充分说明了octopress是基于jekyll开发的一套高定制化的静态博客系统，你可以把它理解成属于jekyll的二次开发，类似于android与miui, flyme之间的关系。Github pages在其中充当的角色，仅仅是提供了一个jekyll的运行环境，还有项目托管，让用户不仅仅能够使用jeykll来搭建一个静态博客，而且还能够使用Git的方式来更新和管理博客。最后说说使用octopress在github上搭建博客的基本原理: 关于一些如何安装jekyll，octopress，以及和github项目库的连接这些准备工作我就不啰嗦了，直接看下文。新建一篇文章命令行 执行rake new_title["我的第一篇博客"]这时候会在source/_post目录下自动生成一个[时间][Title].markdown文件(文件名以及后缀可以自己设定)使用octopress可以用多项设置，你可以设置博客的Header, Footer, 已经每篇文章显示样式，字体大小。octopress默认还为用户添加了博客评论，收藏，分享到facebook, google plus , tw等，这也可以通过修改gemfile中的配置信息进行功能添加或者删除。我们还可以自己添加多个模块，例如中国的一些分享，评论插件。还可以自己定义独立页面。编辑完成后，使用rake generate 命令，octopress便会将.markdown文件自动转换成html文件，生成的文件会保存在sass目录下。当用户访问我们的网页时，便会加载各种css样式以及模板文件和配置信息。然后你可以使用rake preview命令在本地启动一个WebRick服务器预览你的博客，ctrl+c可以关闭服务器。最后使用rake deploy命令将本地生成的博客push到github上的远程库里。注意，在使用rake deploy命令后，octopress会首先将generate好的博客文件(包括html,css,js,img等)全部放到_deploy目录下。一般配置好octopress与你github上的repository后，其会自动为你新建一个分支，默认叫做source分支，主分支叫做master。master分支的内容不需要我们手动去pull和push，这些动作octopress会帮助我们完成。我们所有的修改全部是在source分支下完成的。


Bundle介绍：

Rails 3中引入Bundle来管理项目中所有gem依赖，该命令只能在一个含有Gemfile的目录下执行，如rails 3项目的根目录。



关于Gemfile和Gemfile.lock

所有Ruby项目的信赖包都在Gemfile中进行配置，不再像以往那样，通过require来查找。Rails 3中如果需要require某个gem包，必须通过修改Gemfile文件来管理。

Gemfile.lock则用来记录本机目前所有依赖的Ruby Gems及其版本。所以强烈建议将该文件放入版本控制器，从而保证大家基于同一环境下工作。



Bundle命令详解：



# 显示所有的依赖包
$ bundle show



# 显示指定gem包的安装位置
$ bundle show [gemname]



# 检查系统中缺少那些项目以来的gem包
# 注：如果系统中存在所有项目以来的包，则会输出：The Gemfile's dependencies are satisfied
$ bundle check



# 安装项目依赖的所有gem包
# 注：此命令会尝试更新系统中已存在的gem包
$ bundle install



# 安装指定的gem包
$ bundle install [gemname]



# 更新系统中存在的项目依赖包，并同时更新项目Gemfile.lock文件
$ bundle update



# 更新系统中指定的gem包信息，并同时更新项目Gemfile.lock中指定的包信息
$ bundle update [gemname]

 

# 向项目中添加新的gem包引用
$ gem [gemname], [ver]



# 你还可以指定包依赖关系
$ gem [gemname], :require => [dependence_gemname]



# 你甚至还可以指定gem包的git源
$ gem [gemname], :git => [git_source_url]



# 锁定当前环境
# 可以使用bundle lock来锁定当前环境，这样便不能通过bundle update来更新依赖包的版本，保证了统一的环境
$ bundle lock



# 解除锁定
$ bundle unlock



# 打包当装环境
# bundle package会把当前所有信赖的包都放到 ./vendor/cache/ 目录下，发布时可用来保证包版本的一致性。
$ bundle package



Make sure all dependencies in your Gemfile are available to your application.
$ bundle install [--binstubs=PATH] [--clean] [--deployment] [--frozen]
                 [--full-index] [--gemfile=FILE] [--local] [--no-cache]
                 [--no-prune] [--path=PATH] [--quiet] [--shebang=STRING]
                 [--standalone=ARRAY] [--system] [--without=GROUP GROUP]
                 [--trust-policy=SECURITYLEVEL]
Options:

--binstubs: Generate bin stubs for bundled gems to ./bin

--clean: Run bundle clean automatically after install

--deployment: Install using defaults tuned for deployment environments

--frozen: Do not allow the Gemfile.lock to be updated after this install

--full-index: Use the rubygems modern index instead of the API endpoint

--gemfile: Use the specified gemfile instead of Gemfile

--jobs: Install gems using parallel workers.

--local: Do not attempt to fetch gems remotely and use the gem cache instead

--no-cache: Don't update the existing gem cache.

--no-prune: Don't remove stale gems from the cache.

--path: Specify a different path than the system default ($BUNDLE_PATH or $GEM_HOME). Bundler will remember this value for future installs on this machine

--quiet: Only output warnings and errors.

--retry: Retry network and git requests that have failed.

--shebang: Specify a different shebang executable name than the default (usually 'ruby')

--standalone: Make a bundle that can work without the Bundler runtime

--system: Install to the system location ($BUNDLE_PATH or $GEM_HOME) even if the bundle was previously installed somewhere else for this application

--trust-policy: Sets level of security when dealing with signed gems. Accepts `LowSecurity`, `MediumSecurity` and `HighSecurity` as values.

--without: Exclude gems that are part of the specified named group.

Gems will be installed to your default system location for gems. If your system gems are stored in a root-owned location (such as in Mac OSX), bundle will ask for your root password to install them there.

While installing gems, Bundler will check vendor/cache and then your system's gems. If a gem isn't cached or installed, Bundler will try to install it from the sources you have declared in your Gemfile.

The --system option is the default. Pass it to switch back after using the --path option as described below.

Install your dependencies, even gems that are already installed to your system gems, to a location other than your system's gem repository. In this case, install them to vendor/bundle.
$ bundle install --path vendor/bundle
Further bundle commands or calls to Bundler.setup or Bundler.require will remember this location.
Learn More: Bundler.setup
Learn More: Bundler.require
Install all dependencies except those in groups that are explicitly excluded.
$ bundle install --without development test
Learn More: Groups
Install all dependencies on to a production server. Do not use this flag on a development machine.
$ bundle install --deployment
The --deployment flag activates a number of deployment-friendly conventions:

Isolate all gems into vendor/bundle
Require an up-to-date Gemfile.lock
If bundle package was run, do not fetch gems from rubygems.org. Instead, only use gems in the checked in vendor/cache
Learn More: Deploying
Install gems parallely by starting the number of workers specificed.
$ bundle install --jobs 4
Retry failed network or git requests.
$ bundle install --retry 3

一、Ruby
jekyll提供了很多现成的主题可以使用，里面有很多高大上的款式。



官网上面有专门一节是介绍安装的，不过在实际安装中还是会有一些问题。

需要有下载Ruby环境，选最新的那个版本即可，官网上面安装列中有一个RubyGems，但Ruby1.9.1 以后版本已经自带了，所以无需额外下载。



 

二、切换source源
由于国内网络原因（你懂的），导致 rubygems.org 存放在 Amazon S3 上面的资源文件间歇性连接失败。

有两张方法，一种是切换到淘宝，另外一种是切换到ruby-china，网上大部分的教程都是用淘宝的。

在用淘宝的后老是会出现认证错误，后面上google查问题，发现淘宝的已经不维护了，详见《Ruby China的RubyGems 镜像上线》

我把两种方法都记录了一下，

1） ruby-china

将source改成“https://gems.ruby-china.org/”，在打开的页面中，会告诉你几个指令。

由于我先用了taobao的source，所以这里remove的是淘宝的。



ruby-china中说道：“如果遇到 SSL 证书问题，你又无法解决，请直接用 http://gems.ruby-china.org 避免 SSL 的问题。”

 

2） taobao

将source改成“https://ruby.taobao.org/”，在打开的页面中，会告诉你几个指令。

如果在输入指令出错的话，如下图所示，就是要让你下载下载证书文件。



然后放到某个位置，输入指令set，“D:\Ruby23-x64\cacert.pem”就是文件的具体路径

set SSL_CERT_FILE=D:\Ruby23-x64\cacert.pem


也可以将“SSL_CERT_FILE”设置为环境变量，这样就不用每次都要输入设置的指令。



不知为何，后面我加载包的时候，就是会出现问题，囧，也许是我做了什么操作导致的，额额额。



 

三、安装jekyll
在输入安装指令后，就会看到默认安装了14个包。

gem install jekyll




 

四、启动
从主题列表中选了两个，Minimal Mistakes和Jekyll Clean。前者页面比较全但相对比较复杂，后者页面少但很简洁。

输入指令，

jekyll serve --watch


在显示的文字中有一句让你安装“wdm”，会在下面介绍。

在页面中输入“http://localhost:4000/jekyll-clean/index.html”后就能看到页面了。



 

五、wdm
从 v2.4.0 开始，Jekyll 本地部署时，会相当于以前版本加 --watch 一样，监听其源文件的变化。

而 Windows 似乎有时候并不会奏效，若你碰到，可安装 wdm (Windows Directory Monitor ) 来改善这个问题。

如果要安装“wdm”得要先安装“Devkit”，在打开的网站中下载后，会让你解压到某个文件夹，接下来就是进入到这个文件夹中。

执行指令“ruby dk.rb init”。



再执行指令“ruby dk.rb install”，不过提示我先去修改“config.yml”中的路径。



config.yml文件就在解压出来的文件中。

 

再执行install指令。



 

六、Gemfile文件
Gemfile是一个用于描述gem之间依赖的文件。gem是一堆Ruby代码的集合，它能够为我们提供调用。

Gemfile是可通过Bundler创建。

gem install bundler
bundle init
bundle install
Gemfile文件中设置的内容如下：

source "https://rubygems.org"

gem "jekyll-paginate"
gem "kramdown"
gem "jekyll-watch"
gem "wdm", "~> 0.1.0" if Gem.win_platform?


 

七、自动刷新页面
1）修改Gemfile文件

gem 'guard'
gem 'guard-jekyll-plus'
gem 'guard-livereload'
要添加三个包，执行“bundle install”，如果执行出错，那就一个一个加吧。

 

2）创建guard配置文件

执行指令，将会生成一个Guardfile文件。

guard init
生成的Guardfile文件内有一些代码，在代码的最后添加如下代码：

guard 'jekyll-plus', :serve => true do
  watch /.*/
  ignore /^_site/
end

guard 'livereload' do
  watch /.*/
end
 

3）添加livereload插件

安装Live Reload Extension，如果是chrome，就到Chrome Web Store下载。



安装成功后，在右上角可以看到一个小按钮

如果是运行状态，那么会自动添加一个js文件引用：



 

4）运行

执行运行指令：

bundle exec guard start
这里注意一下，livereload要先关闭。

运行上面指令，当出现下面的内容后，再运行livereload。



然后会出现“connected”连接了，接下来修改内容就会自动刷新页面了。



试用后发现，有时候会刷新不成功，还是原来的样子，看来某些地方还需要改进。

 

demo下载：

http://download.csdn.net/detail/loneleaf1/9508074