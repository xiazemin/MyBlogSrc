---
title: redismock
layout: post
category: golang
author: 夏泽民
---
https://github.com/go-redis/redismock
db, mock := redismock.NewClientMock()
mock.ExpectGet(key).RedisNil()

https://golangrepo.com/repo/go-redis-redismock
<!-- more -->
还有其他实现
https://github.com/elliotchance/redismock
https://github.com/alicebob/miniredis
https://stackoverflow.com/questions/58016501/returning-a-mock-from-a-package-function

https://medium.com/easyread/unit-test-redis-in-golang-c22b5589ea37

// newTestRedis returns a redis.Cmdable.
func newTestRedis() *redismock.ClientMock {
        mr, err := miniredis.Run()
        if err != nil {
                panic(err)
        }
        client := redis.NewClient(&redis.Options{
                Addr: mr.Addr(),
        })
        return redismock.NewNiceMock(client)
}
https://elliotchance.medium.com/mocking-redis-in-unit-tests-in-go-28aff285b98
