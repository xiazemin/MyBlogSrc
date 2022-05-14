#https://www.elastic.co/cn/downloads/elasticsearch
#vi config/elasticsearch.yml
#xpack.security.enabled: false
#./bin/elasticsearch-reset-password -u elastic

#https://github.com/medcl/elasticsearch-analysis-ik/releases
# cd your-es-root/plugins/ && mkdir ik
#unzip plugin to folder your-es-root/plugins/ik

#./bin/elasticsearch

curl -XPUT --user elastic:OngoP+zAoLtfOU-MDPr= --header 'Content-Type: application/json' --header 'Accept: application/json' -d '{"mappings":{"properties":{"category":{"type":"text"},"title":{"type":"text"},"content":{"analyzer":"ik_max_word","search_analyzer":"ik_max_word","type":"text"}}}}' 'http://127.0.0.1:9200/my_blog'

curl -XDELETE --user elastic:OngoP+zAoLtfOU-MDPr= 'http://127.0.0.1:9200/my_blog'

go mod init es
go get -u github.com/olivere/elastic/v7
go get github.com/russross/blackfriday #解析markdown
go get github.com/microcosm-cc/bluemonday #解决跨站脚本攻击，返回安全的html
go get github.com/PuerkitoBio/goquery #查找Document的内容，这个库就实现了类似 jQuery 的功能，让你能方便的使用 Go 语言操作 HTML 文档
go get github.com/sourcegraph/syntaxhighlight #语法高亮syntaxhighlight包提供代码的语法高亮显示。 它目前使用独立于语言的词法分析器， 并在JavaScript，Java，Ruby，Python，Go和C上表现出色。
#主要的AsHTML(src []byte) ([]byte, error)函数，输出就是HTML 与google-code-prettify相同的CSS类，因此任何样式表也应该适用于此包。
go get -u github.com/djimenez/iconv-go #https://github.com/PuerkitoBio/goquery/wiki/Tips-and-tricks

% go run main.go sync
% go run main.go serve

#https://www.elastic.co/cn/downloads/past-releases/kibana-8-1-2
% ./bin/kibana
#http://127.0.0.1:5601/app/home
