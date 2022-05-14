---
title: server格式
layout: post
category: web
author: 夏泽民
---
<!-- more -->
<div class="container">
	<div class="row">
CGI程序不是放在服务器上就能顺利运行，如果要想使其在服务器上顺利的运行并准确的处理用户的请求，则须对所使用的服务器进行必要的设置。
配置：根据所使用的服务器类型以及它的设置把CGI程序放在某一特定的目录中或使其带有特定的扩展名。
⑴CREN格式服务器的配置：
编辑CREN格式服务器的配置文件（通常为/etc/httpd.conf）在文件中加入：Exec cgi-bin/*/home/www/cgi-bin/*.exec。命令中出现的第一个参数cgi-bin/*指出了在URL中出现的目录名字，并表示它出现在系统主机后的第一个目录中，如：http://edgar.stern.nyn.***/cgi-bin/。命令中的第二个参数表示CGI程序目录放在系统中的真实路径。
CGI目录除了可以跟网络文件放在同一目录中，也可以放在系统的其它目录中，但必须保证在你的系统中也具有同样的目录。在对服务器完成设置后，须重新启动服务器（除非HTTP服务器是用inetd启动的）。
⑵NCSA格式服务器的配置
在NCSA格式服务器上有两种方法进行设置：
①在srm.conf文件（通常在conf目录下）中加入：Script Alias/cgi-bin/cgi-bin/。Script Alias命令指出某一目录下的文件是可执行程序，且这个命令是用来执行这些程序的；此命令的两个参数与CERN格式服务器中的Exec命令的参数的含意一样。
②在srm.conf文件加入：Add type application/x-httpd-cgi.cgi。此命令表示在服务器上增加了一种新的文件类型，其后第一个参数为CGI程序的MIME类型，第二个参数是文件的扩展名，表示以这一扩展名为扩展名的文件是CGI程序。
在用上述方法之一设置服务器后，都得重新启动服务器（除非HTTP服务器是用inetd启动的）。
编写语言
CGI可以用任何一种语言编写，只要这种语言具有标准输入、输出和环境变量。对初学者来说，最好选用易于归档和能有效表示大量数据结构的语言，例如UNIX环境中：
· Perl (Practical Extraction and Report Language)
· Bourne Shell或者Tcl (Tool Command Language)
· PHP(Hypertext Preprocessor))
由于C语言有较强的平台无关性，所以也是编写CGI程序的首选。
	</div>
</div>
