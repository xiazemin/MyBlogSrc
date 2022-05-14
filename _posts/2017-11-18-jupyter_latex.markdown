---
title: jupyter 数学公式
layout: post
category: spark
author: 夏泽民
---
<div class="container">
<div class="row">
	Notebook 文档是由一系列单元（Cell）构成，如何使用Cell？

类型

Code
可执行的代码，Jupyter命令、Unix命令、各种脚本语言代码
Markdown
可书写markdown
Raw NBconvert
应该是默认格式（不确定）
Heading
标题级别，相当于html里面的h1、h2……

主要有两种形式的单元：

代码单元：这里是你编写代码的地方，通过按 Shift + Enter 运行代码，其结果显示在本单元下方。代码单元左边有 In [1]: 这样的序列标记，方便人们查看代码的执行次序。
Markdown 单元：在这里对文本进行编辑，采用 markdown 的语法规范，可以设置文本格式、插入链接、图片甚至数学公式。同样使用 Shift + Enter 运行 markdown 单元来显示格式化的文本。

类似于 Linux 的 Vim 编辑器，在 notebook 中也有两种模式：

编辑模式：编辑文本和代码。选中单元并按 Enter 键进入编辑模式，此时单元左侧显示绿色竖线。
命令模式：用于执行键盘输入的快捷命令。通过 Esc 键进入命令模式，此时单元左侧显示蓝色竖线。
如果要使用快捷键，首先按 Esc 键进入命令模式，然后按相应的键实现对文档的操作。比如切换成代码单元（Y）或 markdown 单元（M），或者在本单元的下方增加一单元（B）。查看所有快捷命令可以按H。
</div>
<div class="row">
数学公式编辑
</div>
<div class="row">
如果你曾做过严肃的学术研究，一定对 LaTeX 并不陌生，这简直是写科研论文的必备工具，不但能实现严格的文档排版，而且能编辑复杂的数学公式。在 Jupyter Notebook 的 markdown 单元中我们也可以使用 LaTeX 的语法来插入数学公式。

在文本行中插入数学公式，使用一对 $符号，比如质能方程 $E = mc^2$。如果要插入一个数学区块，则使用一对美元$符号。比如下面公式表示 z=x/y：
<!-- more -->
</div>
<div class="row">
幻灯片制作
</div>
<div class="row">
既然Jupyter Notebook 擅长展示数据分析的过程，除了通过网页形式分享外，当然也可以将其制作成幻灯片的形式。这里有一个幻灯片示例供参考，其制作风格简洁明晰。

那么如何用 Jupyter Notebook 制作幻灯片呢？首先在 notebook 的菜单栏选择 View > Cell Toolbar > Slideshow，这时在文档的每个单元右上角显示了 Slide Type 的选项。通过设置不同的类型，来控制幻灯片的格式。有如下5中类型：

Slide：主页面，通过按左右方向键进行切换。
Sub-Slide：副页面，通过按上下方向键进行切换。
Fragment：一开始是隐藏的，按空格键或方向键后显示，实现动态效果。
Skip：在幻灯片中不显示的单元。
Notes：作为演讲者的备忘笔记，也不在幻灯片中显示。

当编写好了幻灯片形式的 notebook，如何来演示呢？这时需要使用 nbconvert：
{% highlight bash %}
jupyter nbconvert notebook.ipynb --to slides --post serve
{% endhighlight %}
</div>

<div class="row">
魔术关键字

魔术关键字（magic keywords），正如其名，是用于控制 notebook 的特殊的命令。它们运行在代码单元中，以 % 或者 %% 开头，前者控制一行，后者控制整个单元。

比如，要得到代码运行的时间，则可以使用 %timeit；如果要在文档中显示 matplotlib 包生成的图形，则使用 % matplotlib inline；如果要做代码调试，则使用 %pdb。但注意这些命令大多是在Python kernel 中适用的，其他 kernel 大多不适用。有许许多多的魔术关键字可以使用，更详细的清单请参考 Built-in magic commands 。
</div>
</div>
