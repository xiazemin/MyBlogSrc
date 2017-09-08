---
title: 自动化替换网站引用资源到本地工具
layout: post
category: jekyll
author: 夏泽民
---
{% highlight html linenos %}#!/bin/bash
clear
function replace(){
urlT=${1//\//\\\/}
echo $urlT
}
function replaceDot(){
urlT=${1//\./\\\.}
echo $urlT
}
echo '' > temp.txt
grep  -nrEo  "\<a.*\>|\<script.*\>|\<link.*\>" ./ |grep -E "href=\"http|src=\"http|href=\'http|src=\'http" |grep -v github |grep -v disqus |grep -v 'wb.js' |awk -F ' ' '{for(i=1;i<=NF;i++){split($i,x,"\""); if(x[1]=="src="){print x[2];} }}' >> temp.txt
urls=` sort -u temp.txt |grep js `
for url in $urls
do 
fileCmd=` echo $url |awk -F '/' '{print "curl -o ./js/"$NF " " $0 "\n" }' `
echo $fileCmd
$fileCmd
done
for url in $urls
do
 newUrl=`  echo $url |awk -F '/' '{print "{{site.baseurl}}/js/"$NF }' `
echo $url
echo $newUrl
files=` grep $url .  -rl |grep -v "_site" |grep -v "temp" |sort -u `
for file in $files
do
echo $file
urlT=` replace $url `
urlT=` replaceDot $urlT `
newUrlT=` replace $newUrl `
cmd=` echo " sed -i 'temp.bak' 's/$urlT/$newUrlT/' $file" `
echo $cmd |bash
echo $cmd
$file `
done 
done{% endhighlight %}