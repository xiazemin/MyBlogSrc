package es

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

var esOnce sync.Once
var client *elastic.Client

func GetConn() *elastic.Client {
	if client == nil {
		esOnce.Do(func() {
			options := []elastic.ClientOptionFunc{
				elastic.SetURL("http://127.0.0.1:9200"),
				elastic.SetSniff(true),                                             //是否开启集群嗅探
				elastic.SetHealthcheckInterval(10 * time.Second),                   //设置两次运行状况检查之间的间隔, 默认60s
				elastic.SetGzip(false),                                             //启用或禁用gzip压缩
				elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)), //ERROR日志输出配置
				elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),          //INFO级别日志输出配置
				elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)),
			}
			options = append(options, elastic.SetBasicAuth(
				"elastic",              //账号
				"OngoP+zAoLtfOU-MDPr=", //密码
			))
			var err error
			client, err = elastic.NewClient(options...)
			if err != nil {
				fmt.Println(err)
			}
		})
	}
	return client
}
