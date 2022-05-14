---
title: interface
layout: post
category: golang
author: 夏泽民
---
package main

import (
    "errors"
    "fmt"
    "log"
)

func main() {
    var e interface{}
    e = func() error {
        return errors.New("err")
    }()
    if e != nil {
        fmt.Printf("%T\n", e)
        log.Println(e)
    }
    fmt.Println(e)
}
输出内容：

*errors.errorString
2019/01/05 18:54:43 err
err

这边很容易将e的类型误认为是error，但是实际运行中却被转换成*errors.errorString。
<!-- more -->
src/builtin/builtin.go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
	Error() string
}


src/errors/errors.go
// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}
