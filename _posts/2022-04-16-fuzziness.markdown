---
title: fuzziness
layout: post
category: elasticsearch
author: 夏泽民
---

<!-- more -->
基本的匹配（Query）查询
GET /bookdb_index/book/_search?q=guide　#查询任一字段包含guide的记录。

下面是完整body版的查询：

{
    "query": {
        "multi_match": {
            "query": "guide",
            "fields": [
                "_all"
            ]
        }
    }
}
GET /bookdb_index/book/_search?q=title:in action

查询title字段包含in action的书。

POST /bookdb_index/book/_search

{
    "query": {
        "match": {
            "title": "in action"
        }
    },
    "size": 2,
    "from": 0,
    "_source": [
        "title",
        "summary",
        "publish_date"
    ],
    "highlight": {
        "fields": {
            "title": {}
        }
    }
}
Size指定返回的结果条数，from指定起始位子，_score指定要返回的字段。

多字段（Multi-field）查询
POST /bookdb_index/book/_search

{
    "query": {
        "multi_match": {
            "query": "elasticsearch guide",
            "fields": [
                "title",
                "summary"
            ]
        }
    }
}
Boosting

查询多个字段时，可以把个别字段的分数提高倍数，以此提高此字段的重要程度。

POST /bookdb_index/book/_search

{
    "query": {
        "multi_match": {
            "query": "elasticsearch guide",
            "fields": [
                "title",
                "summary^3"
            ]
        }
    },
    "_source": [
        "title",
        "summary",
        "publish_date"
    ]
}
Bool查询

为了提供更相关或者特定的结果，AND/OR/NOT操作符可以用来调整我们的查询。它是以布尔查询的方式来实现的。布尔查询接受如下参数：

must等同于AND。
must_not等同于NOT。
should等同于OR。
如：查询书名包含Elasticsearch或者Solr，并且它的作者是Clinton Gormley不是Radu Gheorge。

POST /bookdb_index/book/_search

{
    "query": {
        "bool": {
            "must": {
                "bool": {
                    "should": [
                        {
                            "match": {
                                "title": "Elasticsearch"
                            }
                        },
                        {
                            "match": {
                                "title": "Solr"
                            }
                        }
                    ]
                }
            },
            "must": {
                "match": {
                    "authors": "clinton gormely"
                }
            },
            "must_not": {
                "match": {
                    "authors": "radu gheorge"
                }
            }
        }
    }
}
模糊（Fuzzy）查询
在进行匹配和多项匹配时，可以启用模糊匹配来捕捉拼写错误，模糊度是基于原始单词的编辑距离来指定的。

POST /bookdb_index/book/_search

{
    "query": {
        "multi_match": {
            "query": "comprihensiv guide",
            "fields": [
                "title",
                "summary"
            ],
            "fuzziness": "AUTO"
        }
    },
    "_source": [
        "title",
        "summary",
        "publish_date"
    ],
    "size": 1
}
通配符（Wildcard）查询
通配符查询允许你指定匹配的模式，而不是整个术语。

?　匹配任何字符。

*　匹配零个或多个字符。

如：查找名称以字母t开头的所有作者的记录。

POST /bookdb_index/book/_search

{
    "query": {
        "wildcard": {
            "authors": "t*"
        }
    },
    "_source": [
        "title",
        "authors"
    ],
    "highlight": {
        "fields": {
            "authors": {}
        }
    }
}
正则（Regexp）查询
正则查询可以支持比通配符查询更复杂的模式。

POST /bookdb_index/book/_search

{
    "query": {
        "regexp": {
            "authors": "t[a-z]*y"
        }
    },
    "_source": [
        "title",
        "authors"
    ],
    "highlight": {
        "fields": {
            "authors": {}
        }
    }
}
短语匹配查询
短语匹配查询 要求在请求字符串中的所有查询项必须都在文档中存在，文中顺序也得和请求字符串一致，且彼此相连。默认情况下，查询项之间必须紧密相连，但可以设置 slop 值来指定查询项之间可以分隔多远的距离，结果仍将被当作一次成功的匹配。

POST /bookdb_index/book/_search

{
    "query": {
        "multi_match": {
            "query": "search engine",
            "fields": [
                "title",
                "summary"
            ],
            "type": "phrase",
            "slop": 3
        }
    },
    "_source": [
        "title",
        "summary",
        "publish_date"
    ]
}
短语前缀（Match Phrase Prefix）查询
短语前缀式查询 能够进行 即时搜索（search-as-you-type） 类型的匹配，或者说提供一个查询时的初级自动补全功能，无需以任何方式准备你的数据。和 match_phrase 查询类似，它接收slop 参数（用来调整单词顺序和不太严格的相对位置）和 max_expansions 参数（用来限制查询项的数量，降低对资源需求的强度）。

POST /bookdb_index/book/_search

{
    "query": {
        "match_phrase_prefix": {
            "summary": {
                "query": "search en",
                "slop": 3,
                "max_expansions": 10
            }
        }
    },
    "_source": [
        "title",
        "summary",
        "publish_date"
    ]
}
查询字符串（Query String）
查询字符串 类型（query_string）的查询提供了一个方法，用简洁的简写语法来执行 多匹配查询、 布尔查询 、 提权查询、 模糊查询、 通配符查询、 正则查询 和范围查询。下面的例子中，我们在那些作者是 “grant ingersoll” 或 “tom morton” 的某本书当中，使用查询项 “search algorithm” 进行一次模糊查询，搜索全部字段，但给 summary 的权重提升 2 倍。

POST /bookdb_index/book/_search

{
    "query": {
        "query_string": {
            "query": "(saerch~1 algorithm~1) AND (grant ingersoll)  OR (tom morton)",
            "fields": [
                "_all",
                "summary^2"
            ]
        }
    },
    "_source": [
        "title",
        "summary",
        "authors"
    ],
    "highlight": {
        "fields": {
            "summary": {}
        }
    }
}
简单查询字符串（Simple Query String）
简单请求字符串 类型（simple_query_string）的查询是请求字符串类型（query_string）查询的一个版本，它更适合那种仅暴露给用户一个简单搜索框的场景；因为它用 +/\|/- 分别替换了 AND/OR/NOT，并且自动丢弃了请求中无效的部分，不会在用户出错时，抛出异常。

POST /bookdb_index/book/_search

{
    "query": {
        "simple_query_string": {
            "query": "(saerch~1 algorithm~1) + (grant ingersoll)  | (tom morton)",
            "fields": [
                "_all",
                "summary^2"
            ]
        }
    },
    "_source": [
        "title",
        "summary",
        "authors"
    ],
    "highlight": {
        "fields": {
            "summary": {}
        }
    }
}
词条（Term）/多词条（Terms）查询
以上例子均为 full-text(全文检索) 的示例。有时我们对结构化查询更感兴趣，希望得到更准确的匹配并返回结果，词条查询 和 多词条查询 可帮我们实现。在下面的例子中，我们要在索引中找到所有由 Manning 出版的图书。

POST /bookdb_index/book/_search

{
    "query": {
        "term": {
            "publisher": "manning"
        }
    },
    "_source": [
        "title",
        "publish_date",
        "publisher"
    ]
}
可使用词条关键字来指定多个词条，将搜索项用数组传入。

{
    "query": {
        "terms": {
            "publisher": [
                "oreilly",
                "packt"
            ]
        }
    }
}
词条（Term）查询-排序（Sorted）
词条查询 的结果（和其他查询结果一样）可以被轻易排序，多级排序也被允许。

POST /bookdb_index/book/_search

{
    "query": {
        "term": {
            "publisher": "manning"
        }
    },
    "_source": [
        "title",
        "publish_date",
        "publisher"
    ],
    "sort": [
        {
            "publish_date": {
                "order": "desc"
            }
        },
        {
            "title": {
                "order": "desc"
            }
        }
    ]
}
范围查询
另一个结构化查询的例子是范围查询。

注：范围查询 用于日期、数字和字符串类型的字段。

如：查找2015年出版的书。

POST /bookdb_index/book/_search

{
    "query": {
        "range": {
            "publish_date": {
                "gte": "2015-01-01",
                "lte": "2015-12-31"
            }
        }
    },
    "_source": [
        "title",
        "publish_date",
        "publisher"
    ]
}
过滤（Filtered）查询
过滤查询允许你可以过滤查询结果。如：要在标题或摘要中检索一些书，查询项为Elasticsearch，但我们又想筛出那些仅有20个以上评论的。

POST /bookdb_index/book/_search

{
    "query": {
        "filtered": {
            "query": {
                "multi_match": {
                    "query": "elasticsearch",
                    "fields": [
                        "title",
                        "summary"
                    ]
                }
            },
            "filter": {
                "range": {
                    "num_reviews": {
                        "gte": 20
                    }
                }
            }
        }
    },
    "_source": [
        "title",
        "summary",
        "publisher",
        "num_reviews"
    ]
}
过滤查询 将在 ElasticSearch 5 中移除，使用 布尔查询 替代。下面有个例子使用 布尔查询 重写上面的例子。

POST /bookdb_index/book/_search

{
    "query": {
        "bool": {
            "must": {
                "multi_match": {
                    "query": "elasticsearch",
                    "fields": [
                        "title",
                        "summary"
                    ]
                }
            },
            "filter": {
                "range": {
                    "num_reviews": {
                        "gte": 20
                    }
                }
            }
        }
    },
    "_source": [
        "title",
        "summary",
        "publisher",
        "num_reviews"
    ]
}
多重过滤（Multiple Filters）
多重过滤 可以结合 布尔查询 使用，下一个例子中，过滤查询决定只返回那些包含至少20条评论，且必须在 2015 年前出版，且由 O’Reilly 出版的结果。

POST /bookdb_index/book/_search

{
    "query": {
        "filtered": {
            "query": {
                "multi_match": {
                    "query": "elasticsearch",
                    "fields": [
                        "title",
                        "summary"
                    ]
                }
            },
            "filter": {
                "bool": {
                    "must": {
                        "range": {
                            "num_reviews": {
                                "gte": 20
                            }
                        }
                    },
                    "must_not": {
                        "range": {
                            "publish_date": {
                                "lte": "2014-12-31"
                            }
                        }
                    },
                    "should": {
                        "term": {
                            "publisher": "oreilly"
                        }
                    }
                }
            }
        }
    },
    "_source": [
        "title",
        "summary",
        "publisher",
        "num_reviews",
        "publish_date"
    ]

https://blog.csdn.net/troubleshooter/article/details/122315749