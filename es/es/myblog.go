package es

import (
	"context"

	"github.com/olivere/elastic/v7"
)

type Blog interface {
	Upsert(ctx context.Context, id string, data map[string]interface{}) error
	Search(ctx context.Context, keyWord string) (*elastic.SearchResult, error)
}

type Myblog struct {
	esClient *elastic.Client
	Index    string
}

func NewMyblog() Blog {
	return &Myblog{
		esClient: GetConn(),
		Index:    "my_blog",
	}
}

func (m *Myblog) Upsert(ctx context.Context, id string, data map[string]interface{}) error {
	client = GetConn()
	if _, err := client.Update().
		Index(m.Index).
		Id(id).
		//Type("doc").
		// data为结构体或map, 需注意的是如果使用结构体零值也会去更新原记录
		//Upsert(data).
		Doc(data).
		// true 无则插入, 有则更新, 设置为false时记录不存在将报错
		DocAsUpsert(true).
		//ScriptedUpsert(false).
		Do(ctx); err != nil {
		return err
	}
	return nil
}

func (m *Myblog) Search(ctx context.Context, keyWord string) (*elastic.SearchResult, error) {
	client = GetConn()
	return client.Search(m.Index).From(0).Size(1000).
		Query(m.BuildQuery(keyWord)).Do(ctx)
}

func (m *Myblog) BuildQuery(keyWord string) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	query.Should(elastic.NewMatchQuery("category", keyWord),
		elastic.NewMatchQuery("title", keyWord),
		elastic.NewMatchQuery("content", keyWord),
	)
	return query
}
