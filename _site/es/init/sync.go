package init

import (
	"context"
	"encoding/json"
	"es/es"
	"es/format"
	"fmt"
	"io/ioutil"
	"runtime"
)

func Sync(dir string) {
	var blackList []string

	posts, err := ioutil.ReadDir(dir)
	if err != nil {
		funcName, file, line, ok := runtime.Caller(0)
		fmt.Println(funcName, file, line, ok, err)

		return
	}
	for _, post := range posts {
		if post.Name() == ".DS_Store" || post.Name() == "." || post.Name() == ".." || inBlacckList(post.Name()) || post.IsDir() {
			funcName, file, line, ok := runtime.Caller(0)
			fmt.Println(funcName, file, line, ok, "skip:", post.Name())
			continue
		}
		fileName := dir + "/" + post.Name()
		if !syncOne(fileName, post.Name(), blackList) {
			return
		}
	}

	black, _ := json.Marshal(blackList)
	fmt.Println(string(black))
}

func syncOne(path, name string, blackList []string) bool {
	myblog := es.NewMyblog()
	ctx := context.Background()

	content, err := ioutil.ReadFile(path)
	if err != nil {
		funcName, file, line, ok := runtime.Caller(0)
		fmt.Println(funcName, file, line, ok, err)
		return false
	}
	html := format.ToHtml(content)
	//		funcName, file, line, ok := runtime.Caller(0)
	// fmt.Println(funcName, file, line, ok, html)

	header, err := format.ParseHtml(html, name)
	if err != nil {
		blackList = append(blackList, name)
		funcName, file, line, ok := runtime.Caller(0)
		fmt.Println(funcName, file, line, ok, name, err)
		return true
	}
	// return
	//		funcName, file, line, ok := runtime.Caller(0)
	//fmt.Println(funcName, file, line, ok, header)

	if err := myblog.Upsert(ctx, name, map[string]interface{}{
		"title":    header["title"],
		"category": header["category"] + header["categories"],
		"content":  html,
	}); err != nil {
		funcName, file, line, ok := runtime.Caller(0)
		fmt.Println(funcName, file, line, ok, err)
		return false
	}
	return true
}
