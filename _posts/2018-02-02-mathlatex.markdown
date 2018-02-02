---
title: mathlatex
layout: post
category: web
author: 夏泽民
---
<!-- more -->
1,LaTeX for WordPress
<对于WordPress博客来说，使用MathJax库的一个简单方法，就是直接使用一个叫LaTeX for WordPress插件。安装插件，简单配置，就可以使用MathJax的js库提供的数学公式在网页上的渲染支持。本博没有使用插件，而是直接在博客主题引用MathJax的js库。

然后，在网页上编辑公式，只要把LaTeX语法的公式放入MathJax的界定符号之内即可，默认情况下，$$LaTex语法$$表示换行居中显示数学公式，而\(LaTex语法\)表示在行内显示数学公式，即inline的显示方法。
2,Google Chart
Google Chart接受TeX语言，实时返回数学公式的图片
http://www.ruanyifeng.com/blog/2011/07/formula_online_generator.html
https://developers.google.com/chart/?csw=1
3,MathJax.js
https://www.mathjax.org/
第一步先是引入：你可以通过引入CDN，也可以下载js引入，个人推荐CDN引入。
<script type="text/javascript"
  src="https://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML">
</script>
第二步：有3种使用方法：
①：TeX and LaTeX格式方式
默认运算符是$$...$$和\[... \]为显示数学，\（...\）用于在线数学。请特别注意，在$...$在线分隔符不使用默认值。这是因为，美元符号出现常常在非数学设置，这可能会导致一些文本被视为意外数学。例如，对于单元分隔符，“...的费用是为第一个$2.50和$2.00每增加一...”将导致短语“2.50的第一个，和”被视为数学因为它属于美元符号之间。注意HTML的标签与TeX语法可能有冲突，“小于号/大于号/ampersands&”需要前后空格，比如：$$a < b$$；
②：第二种方法：添加MathML中等标签的形式
③第三种：AsciiMath输入。以``为符号。
4、MathML
mathml 是数学标记语言，是一种基于XML（标准通用标记语言的子集）的标准，用来在互联网上书写数学符号和公式的置标语言。 
5,KateX
https://khan.github.io/KaTeX/ 

