---
title: FileServer
layout: post
category: golang
author: 夏泽民
---
用Go语言写一个简单的HTTP服务器，及静态文件服务器

需要先httpserver接然后转到静态服务器
<!-- more -->
{% raw %}
package main

import (
	// "fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", SayHello)
	serveMux.HandleFunc("/bye", SayBye)
	// serveMux.HandleFunc("/static", StaticServer)

	server := http.Server{
		Addr:        ":8080",
		Handler:     serveMux,
		ReadTimeout: 5 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func SayHello(w http.ResponseWriter, r *http.Request) {
	if ok, _ := regexp.MatchString("/static/", r.URL.String()); ok {
		StaticServer(w, r)
		return
	}
	io.WriteString(w, "hello world")
}

func SayBye(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Byebye")
}

func StaticServer(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	http.StripPrefix("/static/",
		http.FileServer(http.Dir(wd))).ServeHTTP(w, r)
}




// http.ListenAndServe(addr string, handler Handler)
	// handler		http.Handler		interface(ServeHTTP)

	// handler	   	single				实现了ServeHTTP
	// http.ListenAndServe(":8080", single) 所以可以这样

	// NewServeMux	http.ServeMux		实现了ServeHTTP
	// http.ListenAndServe(":8080", mux) 所以可以这样

	// http.Server	包含了Handler 即	http.DefaultServeMux
	// DefaultServeMux	ServeMux		实现了ServeHTTP

	// http.ListenAndServe(addr string, handler Handler)
	// 底层 Http.Server
	// 调用	server.ListenAndServe()

	// 所以 最后还是调用
	// Http.Server 结构体 下的 ListenAndServe() 方法
	// 参数 addr， handler 主要为了赋值

	// 主要方法为Handler下的ServeHTTP方法
	// 负责总控
{% endraw %}

