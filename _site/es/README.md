curl -XGET "http://localhost:9200/my_blog/_search" -H 'Content-Type: application/json' -d'
{
  "from": 0,
  "query": {
    "bool": {
      "should": [
        {
          "terms": {
            "category": [
              "怪异模式"
            ]
          }
        },
        {
          "terms": {
            "title": [
              "怪异模式"
            ]
          }
        },
        {
          "terms": {
            "content": [
              "怪异模式"
            ]
          }
        }
      ]
    }
  },
  "size": 1000
}'

curl -XGET "http://localhost:9200/my_blog/_doc/2022-04-17-怪异模式.markdown"


curl -XGET "http://localhost:9200/my_blog/_analyze" -H 'Content-Type: application/json' -d'
{
  "text": ["怪异模式"],
  "analyzer": "ik_max_word"
}'

curl -XGET "http://localhost:9200/my_blog/_search" -H 'Content-Type: application/json' -d'
{
  "profile": true, 
  "from": 0,
  "query": {
    "bool": {
      "should": [
        {
          "match": {
            "content": "怪异模式"
          }
        }
      ]
    }
  },
  "size": 100
}'

