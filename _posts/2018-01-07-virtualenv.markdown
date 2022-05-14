---
title: virtualenv
layout: post
category: web
author: 夏泽民
---
在开发Python应用程序的时候，系统安装的Python3只有一个版本：3.4。所有第三方的包都会被pip安装到Python3的site-packages目录下。

如果我们要同时开发多个应用程序，那这些应用程序都会共用一个Python，就是安装在系统的Python 3。如果应用A需要jinja 2.7，而应用B需要jinja 2.6怎么办？

这种情况下，每个应用可能需要各自拥有一套“独立”的Python运行环境。virtualenv就是用来为一个应用创建一套“隔离”的Python运行环境。

首先，我们用pip安装virtualenv：

$ pip3 install virtualenv
然后，假定我们要开发一个新的项目，需要一套独立的Python运行环境，可以这么做：

第一步，创建目录：

Mac:~ michael$ mkdir myproject
Mac:~ michael$ cd myproject/
Mac:myproject michael$
第二步，创建一个独立的Python运行环境，命名为venv：

Mac:myproject michael$ virtualenv --no-site-packages venv
Using base prefix '/usr/local/.../Python.framework/Versions/3.4'
New python executable in venv/bin/python3.4
Also creating executable in venv/bin/python
Installing setuptools, pip, wheel...done.
命令virtualenv就可以创建一个独立的Python运行环境，我们还加上了参数--no-site-packages，这样，已经安装到系统Python环境中的所有第三方包都不会复制过来，这样，我们就得到了一个不带任何第三方包的“干净”的Python运行环境。

新建的Python环境被放到当前目录下的venv目录。有了venv这个Python环境，可以用source进入该环境：

Mac:myproject michael$ source venv/bin/activate
(venv)Mac:myproject michael$
注意到命令提示符变了，有个(venv)前缀，表示当前环境是一个名为venv的Python环境。

下面正常安装各种第三方包，并运行python命令：

(venv)Mac:myproject michael$ pip install jinja2
...
Successfully installed jinja2-2.7.3 markupsafe-0.23
(venv)Mac:myproject michael$ python myapp.py
...
在venv环境下，用pip安装的包都被安装到venv这个环境下，系统Python环境不受任何影响。也就是说，venv环境是专门针对myproject这个应用创建的。

退出当前的venv环境，使用deactivate命令：

(venv)Mac:myproject michael$ deactivate 
Mac:myproject michael$
此时就回到了正常的环境，现在pip或python均是在系统Python环境下执行。

完全可以针对每个应用创建独立的Python运行环境，这样就可以对每个应用的Python环境进行隔离。

virtualenv是如何创建“独立”的Python运行环境的呢？原理很简单，就是把系统Python复制一份到virtualenv的环境，用命令source venv/bin/activate进入一个virtualenv环境时，virtualenv会修改相关环境变量，让命令python和pip均指向当前的virtualenv环境。
<!-- more -->

如果在命令行中运行virtualenv --system-site-packages ENV, 会继承/usr/lib/python2.7/site-packages下的所有库, 最新版本virtualenv把把访问全局site-packages作为默认行为
default behavior.

2.1. 激活virtualenv

#ENV目录下使用如下命令
➜  ENV git:(master) ✗ source ./bin/activate  #激活当前virtualenv
(ENV)➜  ENV git:(master) ✗ #注意终端发生了变化
#ENV目录下使用如下命令
➜  ENV git:(master) ✗ source ./bin/activate  #激活当前virtualenv
(ENV)➜  ENV git:(master) ✗ #注意终端发生了变化

#使用pip查看当前库
(ENV)➜  ENV git:(master) ✗ pip list
pip (1.5.6)
setuptools (3.6)
wsgiref (0.1.2) #发现在只有这三个
pip freeze  #显示所有依赖
pip freeze > requirement.txt  #生成requirement.txt文件
pip install -r requirement.txt  #根据requirement.txt生成相同的环境
#使用pip查看当前库
(ENV)➜  ENV git:(master) ✗ pip list
pip (1.5.6)
setuptools (3.6)
wsgiref (0.1.2) #发现在只有这三个
pip freeze  #显示所有依赖
pip freeze > requirement.txt  #生成requirement.txt文件
pip install -r requirement.txt  #根据requirement.txt生成相同的环境
2.2. 关闭virtualenv
使用下面命令
$ deactivate
$ deactivate
2.3. 指定python版本
可以使用-p PYTHON_EXE选项在创建虚拟环境的时候指定python版本


#创建python2.7虚拟环境
➜  Test git:(master) ✗ virtualenv -p /usr/bin/python2.7 ENV2.7
Running virtualenv with interpreter /usr/bin/python2.7
New python executable in ENV2.7/bin/python
Installing setuptools, pip...done.
#创建python2.7虚拟环境
➜  Test git:(master) ✗ virtualenv -p /usr/bin/python2.7 ENV2.7
Running virtualenv with interpreter /usr/bin/python2.7
New python executable in ENV2.7/bin/python
Installing setuptools, pip...done.

#创建python3.4虚拟环境
➜  Test git:(master) ✗ virtualenv -p /usr/local/bin/python3.4 ENV3.4
Running virtualenv with interpreter /usr/local/bin/python3.4
Using base prefix '/Library/Frameworks/Python.framework/Versions/3.4'
New python executable in ENV3.4/bin/python3.4
Also creating executable in ENV3.4/bin/python
Installing setuptools, pip...done.
#创建python3.4虚拟环境
➜  Test git:(master) ✗ virtualenv -p /usr/local/bin/python3.4 ENV3.4
Running virtualenv with interpreter /usr/local/bin/python3.4
Using base prefix '/Library/Frameworks/Python.framework/Versions/3.4'
New python executable in ENV3.4/bin/python3.4
Also creating executable in ENV3.4/bin/python
Installing setuptools, pip...done.
到此已经可以解决python版本冲突问题和python库不同版本的问题
3. 其他
3.1. 生成可打包环境
某些特殊需求下,可能没有网络, 我们期望直接打包一个ENV, 可以解压后直接使用, 这时候可以使用virtualenv -relocatable指令将ENV修改为可更改位置的ENV

#对当前已经创建的虚拟环境更改为可迁移
➜  ENV3.4 git:(master) ✗ virtualenv --relocatable ./
Making script ./bin/easy_install relative
Making script ./bin/easy_install-3.4 relative
Making script ./bin/pip relative
Making script ./bin/pip3 relative
Making script ./bin/pip3.4 relative
#对当前已经创建的虚拟环境更改为可迁移
➜  ENV3.4 git:(master) ✗ virtualenv --relocatable ./
Making script ./bin/easy_install relative
Making script ./bin/easy_install-3.4 relative
Making script ./bin/pip relative
Making script ./bin/pip3 relative
Making script ./bin/pip3.4 relative
3.2. 获得帮助

$ virtualenv -h
$ virtualenv -h
当前的ENV都被修改为相对路径, 可以打包当前目录, 上传到其他位置使用

python模块以及导入出现ImportError: No module named 'xxx'问题


python中，每个py文件被称之为模块，每个具有__init__.py文件的目录被称为包。只要模
块或者包所在的目录在sys.path中，就可以使用import 模块或import 包来使用
如果你要使用的模块（py文件）和当前模块在同一目录，只要import相应的文件名就好，比
如在a.py中使用b.py： 
import b 

但是如果要import一个不同目录的文件(例如b.py)该怎么做呢？ 
首先需要使用sys.path.append方法将b.py所在目录加入到搜素目录中。然后进行import即
可，例如 
import sys 
sys.path.append('c:\xxxx\b.py') # 这个例子针对 windows 用户来说的 
大多数情况，上面的代码工作的很好。但是如果你没有发现上面代码有什么问题的话，可要

注意了，上面的代码有时会找不到模块或者包（ImportError: No module named 
xxxxxx），这是因为： 
sys模块是使用c语言编写的，因此字符串支持 '\n', '\r', '\t'等来表示特殊字符。所以

上面代码最好写成： 
sys.path.append('c:\\xxx\\b.py') 
或者sys.path.append('c:/xxxx/b.py') 
这样可以避免因为错误的组成转义字符，而造成无效的搜索目录（sys.path）设置。 


sys.path是python的搜索模块的路径集，是一个list
可以在python 环境下使用sys.path.append(path)添加相关的路径，但在退出python环境后
自己添加的路径就会自动消失了！

3、搜索路径和路径搜索

模块的导入需要叫做“路径搜索”的过程。

搜索路径：查找一组目录

路径搜索：查找某个文件的操作

ImportError: No module named myModule
这种错误就是说：模块不在搜索路径里，从而导致路径搜索失败！

导入模块时，不带模块的后缀名，比如.py
Python搜索模块的路径：
1)、程序的主目录
2)、PTYHONPATH目录（如果已经进行了设置）
3)、标准连接库目录（一般在/usr/local/lib/python2.X/）
4)、任何的.pth文件的内容（如果存在的话）.新功能，允许用户把有效果的目录添加到模块搜索路径中去
.pth后缀的文本文件中一行一行的地列出目录。
这四个组建组合起来就变成了sys.path了，

>>> import sys
>>> sys.path
导入时，Python会自动由左到右搜索这个列表中每个目录。




关于 python ImportError: No module named 'xxx'的问题?
解决方法如下：
1. 使用PYTHONPATH环境变量，在这个环境变量中输入相关的路径，不同的路径之间用逗号
（英文的！)分开，如果PYTHONPATH 变量还不存在，可以创建它！
这里的路径会自动加入到sys.path中，永久存在于sys.path中而且可以在不同的python版本
中共享，应该是一样较为方便的方法。
C:\Users\Administrator\Desktop\test\module1.py:
def func1():
    print("func1")

将C:\Users\Administrator\Desktop\test添加到PYTHONPATH即可直接import module1,然后
调用：module1.func1()即可。


2. 将自己做的py文件放到 site_packages 目录下


3. 使用pth文件，在 site-packages 文件中创建 .pth文件，将模块的路径写进去，一行一
个路径，以下是一个示例，pth文件也可以使用注释：

# .pth file for the  my project(这行是注释)，命名为xxx.pth文件
C:\Users\Administrator\Desktop\test
这个不失为一个好的方法，但存在管理上的问题，而且不能在不同的python版本中共享。


4. 在调用文件中添加sys.path.append("模块文件目录")；


5. 直接把模块文件拷贝到$python_dir/Lib目录下。

通过以上5个方法就可以直接使用import module_name了。

