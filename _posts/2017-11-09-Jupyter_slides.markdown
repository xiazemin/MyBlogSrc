---
title: Jupyter_slides
layout: post
category: web
author: 夏泽民
---
<!-- more -->
使用jupyter完成后，需要将后缀为.ipynb文件转换成.html文件才能展示出效果

1. 打开命令行终端，运行jupyter notebook
2. 在制作slides时，首先在view中，将视图切换到Slidesshow
<img src="{{site.url}}{{site.baseurl}}/img/jupyterSlider.png"/>
3. 在要编辑的文本行中，slice type中选择slice或sub-slice。若选择slice，slice之间是左右切换，每个slice和sub-slice相当于一张幻灯片。同一个slice和它的菜单sub-slice之间是上下切换。
<img src="{{site.url}}{{site.baseurl}}/img/jupyterSunSlider.png"/>
4. 完成slice的制作后，我们就可以将.ipynb文件转换生成.html文件，以网页的形式展示幻灯片。
5. 再打开一个命令行终端，进入所要转换的文件目录下，运行一下命令，生成html文件。
jupyter-nbconvert --to slides test.ipynb --reveal-prefix  'https://cdn.bootcss.com/reveal.js/3.5.0' --output test
6.至此，完成了slice的制作，效果图如下
<img src="{{site.url}}{{site.baseurl}}/img/jupyterResult.png"/>

labels
菜单栏选择View—>Toggle Toolbar—>打开

菜单栏选择View—>Cell Toolbar—>Slidesshow—>选择

Slide
单个view，左右滑动切换

Sub-Slide
Cell的sub-cell，上下滑动切换

Fragment
这个是Slide或Sub-Slide的属性，可以按次序展示，单击一次出现一条

Skip
跳过，注释非演示代码用的

Notes
在页面按s就可以跳出来的注释

Reveal
themes
Sky, Beige, Serif, etc.

transitions
Cube, Zoom, None, etc.

gen
jupyter-nbconvert --to slides Python_Share.ipynb --reveal-prefix '//cdn.bootcss.com/reveal.js/3.2.0' --output Python_Share
server
python -m SimpleHTTPServer 8000
