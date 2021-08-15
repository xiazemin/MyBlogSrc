---
title: singleflight
layout: post
category: golang
author: 夏泽民
---
https://github.com/golang/groupcache/tree/master/singleflight

<!-- more -->
```
// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
```

 防缓存击穿的方式有很多种，比如通过计划任务来跟新缓存使得从前端过来的所有请求都是从缓存读取等等。之前读过 groupCache的源码，发现里面有一个很有意思的库，叫singleFlight, 因为groupCache从节点上获取缓存如果未命中，则会去其他节点寻找，其他节点还没有的话再从数据源获取，所以这个步骤对于防击穿非常有必要。singleFlight使得groupCache在多个并发请求对一个失效的key进行源数据获取时，只让其中一个得到执行，其余阻塞等待到执行的那个请求完成后，将结果传递给阻塞的其他请求达到防止击穿的效果。
 
 https://studygolang.com/articles/18835?fr=sidebar
