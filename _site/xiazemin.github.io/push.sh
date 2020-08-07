#!/bin/bash
jekyll clean
jekyll build --trace
#git add *
ls |xargs git add -f
t=`date`
git commit -m "new blog  $t"
git push https://github.com/xiazemin/MyBlogSrc.git master
cd ./_site
#git add *
ls |xargs git add -f
t=`date`
git commit -m "new blog $t"
git push https://github.com/xiazemin/MyBlog.git master
cd ..
