---
title: writefile
layout: post
category: golang
author: 夏泽民
---
ioutil.WriteFile(lfile, body, os.ModeAppend),每次执行都会清空原有内容，如何只追加

<!-- more -->
fl, err := os.OpenFile(f.FileName, os.O_APPEND|os.O_CREATE, 0644)
if err != nil {
    return 0, err
}
defer fl.Close()
n, err := fl.Write(data)
if err == nil && n < len(data) {
    err = io.ErrShortWrite
}
return n, err

https://www.golangtc.com/t/530ecaa7320b5261970000a6
