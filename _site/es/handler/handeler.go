package handler

import (
	"context"
	"encoding/json"
	"es/es"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, `
	<html>
	<head>
	<!--script type="text/javascript">
	window.onload =function(){
	　　//定位提交按钮
	　　 var inputElement = document.getElementsByTagName("input")[1];
	　　//为提交按钮添加单击事件
	　　 inputElement.onclick = function(){
	  　　//定位<form>标签，forms表示document对象中所有表单的集合，通过下标引用不同的表单，从0开始
		var formElement = document.forms[0];
		//提交表单，提交到action属性指定的地方
		formElement.submit();
	   }
	}
	 </script-->
	</head>
	<body>
	<form action="/search" method="post" enctype="application/x-www-form-urlencoded">
	<input type="text" name="keyword" value=""/>
	<input type="submit" value="提交"/>
	</form>
	</body>
	</html>
	`); err != nil {
		fmt.Println(err)
	}
}

type blog struct {
	Category string
	Content  string
	Title    string
}

const (
	url = "http://127.0.0.1:4000/MyBlog/"
)

func Search(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		http.ServeFile(w, req, "form.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", req.MultipartForm)
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", req.PostForm)
		keyword := req.FormValue("keyword")
		//fmt.Fprintf(w, "keyword = %s\n", keyword)
		myblog := es.NewMyblog()
		ctx := context.TODO()
		result, err := myblog.Search(ctx, keyword)
		if err != nil {
			funcName, file, line, ok := runtime.Caller(0)
			fmt.Println(funcName, file, line, ok, err, result)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		content := ""
		for _, hit := range result.Hits.Hits {
			b := &blog{}
			if err := json.Unmarshal(hit.Source, &b); err != nil {
				fmt.Println(err)
			}
			if b != nil {
				content += "<hr/>" + b.Content
			}
			//http: //127.0.0.1:4000/MyBlog/web/2022/04/17/%E6%80%AA%E5%BC%82%E6%A8%A1%E5%BC%8F.html
			content += "<a href=\"" + url + strings.Trim(b.Category, " ") + "/" + strings.Replace(strings.Replace(hit.Id, "-", "/", -1), ".markdown", ".html", -1) + "\">" + hit.Id + "</a><br/>"
		}

		if _, err := fmt.Fprintf(w, content); err != nil {
			fmt.Println(err)
		}

		//r, _ := json.Marshal(result)
		funcName, file, line, ok := runtime.Caller(0)
		fmt.Println(funcName, file, line, ok, len(result.Hits.Hits))
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
