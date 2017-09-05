分页：

vi  _config.yml 

添加
paginate:5
paginatepath: ['topics/study/page/:num','topics/life/page/:num']

问题：
Deprecation: You appear to have pagination turned on, but you haven't included the `jekyll-paginate` gem. Ensure you have `plugins: [jekyll-paginate]` in your configuration file.

解决方案：
1，gem install jekyll-paginate

2，$gem list |grep jekyll-paginate
jekyll-paginate (1.1.0)

3，$vi Gemfile
gem "jekyll-paginate","~> 1.1.0"

4，$vi _config.yml
plugins:
  - jekyll-feed
  - jekyll-paginate
paginate: 1
paginate_path: "page:num"

5，$ bundle install
$ bundle exec jekyll serve

