#!/bin/bash
bundle exec jekyll clean
bundle exec jekyll build --trace |grep -v 'warning: Using the last argument'
#git add *
ls |xargs git add -f
t=`date`
git commit -m "new blog  $t"
git push https://github.com/xiazemin/MyBlogSrc.git master
cd ./_site
git init
#git add *
ls |xargs git add -f
t=`date`
git commit -m "new blog $t"
git push https://github.com/xiazemin/MyBlog.git master -f
cd ..

#cp -r _site/* ./xiazemin.github.io/
#cd xiazemin.github.io
#git add * -f
#git commit -m 'change to new repository'
#git push https://github.com/xiazemin/xiazemin.github.io master
#cd ..
