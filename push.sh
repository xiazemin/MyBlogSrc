#!/bin/bash
git add *
git commit -m 'new blog'
git push https://github.com/xiazemin/MyBlogSrc.git master
cd ./_site
git add *
git commit -m 'new blog'
git push https://github.com/xiazemin/MyBlog.git master
cd ..
