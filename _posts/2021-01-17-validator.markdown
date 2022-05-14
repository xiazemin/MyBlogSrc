---
title: validator
layout: post
category: golang
author: 夏泽民
---
https://github.com/go-playground/validator
https://github.com/go-playground/locales
https://github.com/go-playground/universal-translator

原理
将验证规则写在struct对字段tag里，再通过反射（reflect）获取struct的tag，实现数据验证。

安装

go get github.com/go-playground/validator/v10
示例

package main

import (
    "fmt"
    "github.com/go-playground/validator/v10"
)

type Users struct {
    Phone   string `form:"phone" json:"phone" validate:"required"`
    Passwd   string `form:"passwd" json:"passwd" validate:"required,max=20,min=6"`
    Code   string `form:"code" json:"code" validate:"required,len=6"`
}

func main() {

    users := &Users{
        Phone:      "1326654487",
        Passwd:       "123",
        Code:            "123456",
    }
    validate := validator.New()
    err := validate.Struct(users)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            fmt.Println(err)//Key: 'Users.Passwd' Error:Field validation for 'Passwd' failed on the 'min' tag
            return
        }
    }
    return
}
验证规则
required ：必填
email：验证字符串是email格式；例："email"
url：这将验证字符串值包含有效的网址;例："url"
max：字符串最大长度；例："max=20"
min:字符串最小长度；例："min=6"
excludesall:不能包含特殊字符；例："excludesall=0x2C"//注意这里用十六进制表示。
len：字符长度必须等于n，或者数组、切片、map的len值为n，即包含的项目数；例："len=6"
eq：数字等于n，或者或者数组、切片、map的len值为n，即包含的项目数；例："eq=6"
ne：数字不等于n，或者或者数组、切片、map的len值不等于为n，即包含的项目数不为n，其和eq相反；例："ne=6"
gt：数字大于n，或者或者数组、切片、map的len值大于n，即包含的项目数大于n；例："gt=6"
gte：数字大于或等于n，或者或者数组、切片、map的len值大于或等于n，即包含的项目数大于或等于n；例："gte=6"
lt：数字小于n，或者或者数组、切片、map的len值小于n，即包含的项目数小于n；例："lt=6"
lte：数字小于或等于n，或者或者数组、切片、map的len值小于或等于n，即包含的项目数小于或等于n；例："lte=6"
跨字段验证
如想实现比较输入密码和确认密码是否一致等类似场景

eqfield=Field: 必须等于 Field 的值；
nefield=Field: 必须不等于 Field 的值；
gtfield=Field: 必须大于 Field 的值；
gtefield=Field: 必须大于等于 Field 的值；
ltfield=Field: 必须小于 Field 的值；
ltefield=Field: 必须小于等于 Field 的值；
eqcsfield=Other.Field: 必须等于 struct Other 中 Field 的值；
necsfield=Other.Field: 必须不等于 struct Other 中 Field 的值；
gtcsfield=Other.Field: 必须大于 struct Other 中 Field 的值；
gtecsfield=Other.Field: 必须大于等于 struct Other 中 Field 的值；
ltcsfield=Other.Field: 必须小于 struct Other 中 Field 的值；
ltecsfield=Other.Field: 必须小于等于 struct Other 中 Field 的值；
<!-- more -->
https://studygolang.com/articles/28414?fr=sidebar
翻译错误信息为中文
通过以上示例我们看到，validator默认的错误提示信息类似如下

Key: 'Users.Name' Error:Field validation for 'Name' failed on the 'CustomValidationErrors' tag
显然这并不是我们想要，如想翻译成中文，或其他语言怎么办？go-playground上提供了很好的解决方法。

先自行安装需要的两个包

https://github.com/go-playground/locales
https://github.com/go-playground/universal-translator

执行：

go get github.com/go-playground/universal-translator
go get github.com/go-playground/locales


dive
dive用于深入到切片、数组和map中进行数据校验
[][]string with validation tag "gt=0,dive,len=1,dive,required"
// gt=0 will be applied to []
// len=1 will be applied to []string
// required will be applied to string

[][]string with validation tag "gt=0,dive,dive,required"
// gt=0 will be applied to []
// []string will be spared validation
// required will be applied to string

https://www.cnblogs.com/zj420255586/p/13542395.html
https://www.cntofu.com/book/73/ch6-web/ch6-04-validator.md

https://github.com/gookit/validate
