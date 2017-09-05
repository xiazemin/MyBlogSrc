#  升级mac自带的rubby
 
 $ curl -L get.rvm.io | bash -s stable
 
 $ rvm -v
 
 $ rvm list known
 
 $ rvm install 2.4.0
 
#  安装jekyll
 
 $ git clone git://github.com/jekyll/jekyll.git
 
 $ cd jekyll
 
 $ script/bootstrap
 
 $ bundle exec rake build
 
 $ ls pkg/*.gem | head -n 1 | xargs gem install -l
 
 参考：https://jekyllrb.com/docs/installation/
 
#  创建jekyll工程，开始使用

$ jekyll new blog

$cd blog 

$jekyll serve

# 修改默认配置

$vi _config.yml

# mac  markdown 编辑器

http://macdown.uranusjr.com/


 
 
 
 