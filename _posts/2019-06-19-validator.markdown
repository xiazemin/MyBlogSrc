---
title: validator
layout: post
category: golang
author: 夏泽民
---
https://github.com/go-validator/validator
这个包怎么用的
type NewUserRequest struct {
	Username string `validate:"min=3,max=40,regexp=^[a-zA-Z]*$"`
	Name string     `validate:"nonzero"`
	Age int         `validate:"min=21"`
	Password string `validate:"min=8"`
}

nur := NewUserRequest{Username: "something", Age: 20}
if errs := validator.Validate(nur); errs != nil {
	// values not valid, deal with errors here
}
使用validate作为 tag 的名字，然后以逗号分隔验证逻辑，然后对于这样的 struct 值，调用validator.Validate(nur)验证
本包有 6 个内置的验证函数，并且可以手动添加自定义的验证函数
我们在定义 struct 的时候可以使用validate作为 tag 的名字添加一些验证规则，然后使用validator.Validate(nur)验证数据是否满足我们定义的规则，如果不满足的话，就会返回 err
<!-- more -->
显然第一步就是，使用reflect将每个 field 的 tag 取出来，我们规定你必须使用 validate 作为 tag 的名字，这样我们就能拿到规则了

然后我们规定多个规则使用逗号分隔，这样我们就可以对于一个 field，可以拿到一组规则了（可以是 0 个，1 个，任意多个）

接下来就是遍历获取到的规则，执行当前规则对于当前值的校验，很显然，这里有这几个因素：field 的值 + 规则

在本包里面，他把规则解释为一个函数，所以还需要有函数的参数（这里脑洞开大一点：其实不仅仅可以是函数，可以自定义语言？但是成本太大：开发成本和用户学习成本；同一个 field 的不同规则之间可以相互作用？）

ok，所以现在一个规则有这么几个因素：field 值 + 规则函数 + 规则参数

那接下来就很好办了，就是调用这个给定的规则函数，参数是 field 值 + 给定的规则参数，看看合不合法，也就是返回一个 error

甚至我们到这里已经可以猜出规则函数的函数签名了：func(v interface{}, param string) error，第一个参数是 field 值，第二个参数是规则参数

具体代码详解
1 遍历
使用包的入口是 validator.Validate，也就是func (mv *Validator) Validate(v interface{}) error

sv := reflect.ValueOf(v)
st := reflect.TypeOf(v)
这个方法里面使用reflect包遍历reflect.Type.NumField()，然后使用reflect.Type.Field(i).Tag获得了各个 field 的 tag

当然在处理反射的时候，有一些注意事项：

对于 reflect.Value.Kind() 为指针的处理方式，递归 .Elem().Interface()

if sv.Kind() == reflect.Ptr && !sv.IsNil() {
	return mv.Validate(sv.Elem().Interface()) // 递归
}
并且这里的验证只支持结构体：

reflect.Value.Kind() 需要是 reflect.Struct 或者 reflect.Interface

if sv.Kind() != reflect.Struct && sv.Kind() != reflect.Interface {
	return ErrUnsupported
}
ok，然后接下来就是遍历所有的 field 进行验证了，所有的验证规则都没有返回 err，那么就返回 nil

2 每个 field 的处理
首先，只支持 exported 的字段：

if !unicode.IsUpper(rune(st.Field(i).Name[0])) {
	continue
}
处理 field 仍然是指针：

对于 reflect.Value.Field(i).Kind() 为指针的处理方式，一直取 .Elem()

f := sv.Field(i)
for f.Kind() == reflect.Ptr && !f.IsNil() {
	f = f.Elem()
}
然后获取 tag 的值（其实这里我认为应该用 lookup）

tag := st.Field(i).Tag.Get(mv.tagName)
跳过 -

if tag == "-" {
	continue
}
然后用规则去验证，代码： err := mv.Valid(f.Interface(), tag)

然后验证 TODO，代码 mv.deepValidateCollection(f, fname, m)

结束，返回 err 或者 nil

3 处理每个 tag 不为空的 field
调用func (mv *Validator) Valid(val interface{}, tags string)

这个函数是干嘛的呢：这个函数的第一个参数不要求是 struct 了，根据提供的 tag 进行验证

跳过 -

v := reflect.ValueOf(val)
处理指针：对于 reflect.Value.Kind() 为指针的处理方式，递归 .Elem().Interface()

if v.Kind() == reflect.Ptr && !v.IsNil() {
	return mv.Valid(v.Elem().Interface(), tags)
}
然后调用func (mv *Validator) validateVar(v interface{}, tag string) error处理

switch v.Kind() {
case reflect.Invalid:
	err = mv.validateVar(nil, tags)
default:
	err = mv.validateVar(val, tags)
}
在 validateVar 中：首先将 tag 解析为 n 个规则（每个规则包括函数，函数名称，参数），然后遍历调用这些规则

for _, t := range tags {
	if err := t.Fn(v, t.Param); err != nil {
		errs = append(errs, err)
	}
}
4 处理所有的 field
调用func (mv *Validator) deepValidateCollection(f reflect.Value, fname string, m ErrorMap)

刚刚处理了所有 tag 不为空的 field，但是还需要处理所有的 field，比如 []AxxxStruct这个 field 的 tag 就是空，但是他里面的AxxxStruct的 tag 不为空

这个方法的最后一个参数一个 err 的 map，以在递归的过程中拿到所有的 error

在这个函数里面，需要处理三种情形：

4.1 struct，interface，ptr
调用func (mv *Validator) Validate(v interface{}) error，已经在上面讲过了，就是入口函数

相当于：struct 的 struct 的 struct，递归调用 Validate 去处理

4.2 array，slice
对于每个元素递归调用 deepValidateCollection

4.3 map
对于 key 和 map 递归调用 deepValidateCollection

代码详解注释
diff --git validator.go validator.go
index a23f3ee..d1ed91c 100644
--- validator.go
+++ validator.go
@@ -1,369 +1,401 @@
 // Package validator implements value validations
 //
 // Copyright 2014 Roberto Teixeira <robteix@robteix.com>
 //
 // Licensed under the Apache License, Version 2.0 (the "License");
 // you may not use this file except in compliance with the License.
 // You may obtain a copy of the License at
 //
 //    http://www.apache.org/licenses/LICENSE-2.0
 //
 // Unless required by applicable law or agreed to in writing, software
 // distributed under the License is distributed on an "AS IS" BASIS,
 // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 // See the License for the specific language governing permissions and
 // limitations under the License.

 package validator

 import (
        "errors"
        "fmt"
        "reflect"
        "regexp"
        "strings"
        "unicode"
 )

 // TextErr is an error that also implements the TextMarshaller interface for
 // serializing out to various plain text encodings. Packages creating their
 // own custom errors should use TextErr if they're intending to use serializing
 // formats like json, msgpack etc.
 type TextErr struct {
        Err error
 }

 // Error implements the error interface.
 func (t TextErr) Error() string {
        return t.Err.Error()
 }

 // MarshalText implements the TextMarshaller
 func (t TextErr) MarshalText() ([]byte, error) {
        return []byte(t.Err.Error()), nil
 }

 var (
        // ErrZeroValue is the error returned when variable has zero valud
        // and nonzero was specified
        ErrZeroValue = TextErr{errors.New("zero value")}
        // ErrMin is the error returned when variable is less than mininum
        // value specified
        ErrMin = TextErr{errors.New("less than min")}
        // ErrMax is the error returned when variable is more than
        // maximum specified
        ErrMax = TextErr{errors.New("greater than max")}
        // ErrLen is the error returned when length is not equal to
        // param specified
        ErrLen = TextErr{errors.New("invalid length")}
        // ErrRegexp is the error returned when the value does not
        // match the provided regular expression parameter
        ErrRegexp = TextErr{errors.New("regular expression mismatch")}
        // ErrUnsupported is the error error returned when a validation rule
        // is used with an unsupported variable type
        ErrUnsupported = TextErr{errors.New("unsupported type")}
        // ErrBadParameter is the error returned when an invalid parameter
        // is provided to a validation rule (e.g. a string where an int was
        // expected (max=foo,len=bar) or missing a parameter when one is required (len=))
        ErrBadParameter = TextErr{errors.New("bad parameter")}
        // ErrUnknownTag is the error returned when an unknown tag is found
        ErrUnknownTag = TextErr{errors.New("unknown tag")}
        // ErrInvalid is the error returned when variable is invalid
        // (normally a nil pointer)
        ErrInvalid = TextErr{errors.New("invalid value")}
 )

 // ErrorMap is a map which contains all errors from validating a struct.
 type ErrorMap map[string]ErrorArray

 // ErrorMap implements the Error interface so we can check error against nil.
 // The returned error is if existent the first error which was added to the map.
 func (err ErrorMap) Error() string {
        for k, errs := range err {
                if len(errs) > 0 {
                        return fmt.Sprintf("%s: %s", k, errs.Error())
                }
        }

        return ""
 }

 // ErrorArray is a slice of errors returned by the Validate function.
 type ErrorArray []error

 // ErrorArray implements the Error interface and returns the first error as
 // string if existent.
 func (err ErrorArray) Error() string {
        if len(err) > 0 {
                return err[0].Error()
        }
        return ""
 }

+// 上面自定义了三个error： TextErr ErrorMap ErrorArray
+
 // ValidationFunc is a function that receives the value of a
 // field and a parameter used for the respective validation tag.
 type ValidationFunc func(v interface{}, param string) error

 // Validator implements a validator
 type Validator struct {
        // Tag name being used.
+       // 取field的哪个tag去验证，默认是validate
        tagName string
        // validationFuncs is a map of ValidationFuncs indexed
        // by their name.
+       // 验证函数
        validationFuncs map[string]ValidationFunc
 }

 // Helper validator so users can use the
 // functions directly from the package
+// 默认验证器，tag是validate，验证函数有内置的5个：nonzero len min max regexp
 var defaultValidator = NewValidator()

 // NewValidator creates a new Validator
 func NewValidator() *Validator {
        return &Validator{
                tagName: "validate",
                validationFuncs: map[string]ValidationFunc{
                        "nonzero": nonzero,
                        "len":     length,
                        "min":     min,
                        "max":     max,
                        "regexp":  regex,
                },
        }
 }

 // SetTag allows you to change the tag name used in structs
+// 切换验证所使用的tag
 func SetTag(tag string) {
        defaultValidator.SetTag(tag)
 }

 // SetTag allows you to change the tag name used in structs
 func (mv *Validator) SetTag(tag string) {
        mv.tagName = tag
 }

 // WithTag creates a new Validator with the new tag name. It is
 // useful to chain-call with Validate so we don't change the tag
 // name permanently: validator.WithTag("foo").Validate(t)
+// 和SetTag一样，是链式调用，并且不会改变原来的验证器的tag值
 func WithTag(tag string) *Validator {
        return defaultValidator.WithTag(tag)
 }

 // WithTag creates a new Validator with the new tag name. It is
 // useful to chain-call with Validate so we don't change the tag
 // name permanently: validator.WithTag("foo").Validate(t)
 func (mv *Validator) WithTag(tag string) *Validator {
        v := mv.copy()
        v.SetTag(tag)
        return v
 }

 // Copy a validator
+// 克隆验证器
 func (mv *Validator) copy() *Validator {
        newFuncs := map[string]ValidationFunc{}
        for k, f := range mv.validationFuncs {
                newFuncs[k] = f
        }
        return &Validator{
                tagName:         mv.tagName,
                validationFuncs: newFuncs,
        }
 }

 // SetValidationFunc sets the function to be used for a given
 // validation constraint. Calling this function with nil vf
 // is the same as removing the constraint function from the list.
+// 添加自定义的验证函数，或者删除某一个验证函数（函数为nil的情况）
 func SetValidationFunc(name string, vf ValidationFunc) error {
        return defaultValidator.SetValidationFunc(name, vf)
 }

 // SetValidationFunc sets the function to be used for a given
 // validation constraint. Calling this function with nil vf
 // is the same as removing the constraint function from the list.
 func (mv *Validator) SetValidationFunc(name string, vf ValidationFunc) error {
        if name == "" {
                return errors.New("name cannot be empty")
        }
        if vf == nil {
                delete(mv.validationFuncs, name)
                return nil
        }
        mv.validationFuncs[name] = vf
        return nil
 }

 // Validate validates the fields of a struct based
 // on 'validator' tags and returns errors found indexed
 // by the field name.
+//
+// 基于 validator 的tag对struct进行校验
 func Validate(v interface{}) error {
        return defaultValidator.Validate(v)
 }

 // Validate validates the fields of a struct based
 // on 'validator' tags and returns errors found indexed
 // by the field name.
 func (mv *Validator) Validate(v interface{}) error {
        sv := reflect.ValueOf(v)
        st := reflect.TypeOf(v)
+
+       // 对于 reflect.Value.Kind() 为指针的处理方式，递归 .Elem().Interface()
        if sv.Kind() == reflect.Ptr && !sv.IsNil() {
                return mv.Validate(sv.Elem().Interface())
        }
+
+       // reflect.Value.Kind() 需要是 reflect.Struct 或者 reflect.Interface
        if sv.Kind() != reflect.Struct && sv.Kind() != reflect.Interface {
                return ErrUnsupported
        }

-       nfields := sv.NumField()
        m := make(ErrorMap)
-       for i := 0; i < nfields; i++ {
+       // 遍历 field
+       // reflect.Type.NumField()
+       for i := 0; i < sv.NumField(); i++ {
+               // reflect.Type.Field(i)
+               //                      .Name
                fname := st.Field(i).Name
                if !unicode.IsUpper(rune(fname[0])) {
+                       // 只处理exported的field
                        continue
                }

+               // 对于 reflect.Value.Field(i).Kind() 为指针的处理方式，一直取 .Elem()
                f := sv.Field(i)
-               // deal with pointers
                for f.Kind() == reflect.Ptr && !f.IsNil() {
                        f = f.Elem()
                }
+               // 获取tag： reflect.Type.Field(i).Tag.Get(name)
                tag := st.Field(i).Tag.Get(mv.tagName)
                if tag == "-" {
+                       // 跳过 `-`
                        continue
                }
                var errs ErrorArray

+               // 处理tag不为空的
                if tag != "" {
+                       // reflect.Value.Field(i).Interface() 对应的field的值 ，以及tag的名字，使用.Valid进行校验
                        err := mv.Valid(f.Interface(), tag)
                        if errors, ok := err.(ErrorArray); ok {
                                errs = errors
                        } else {
                                if err != nil {
                                        errs = ErrorArray{err}
                                }
                        }
                }

+               // TODO
                mv.deepValidateCollection(f, fname, m) // no-op if field is not a struct, interface, array, slice or map

                if len(errs) > 0 {
                        m[st.Field(i).Name] = errs
                }
        }

        if len(m) > 0 {
                return m
        }
        return nil
 }

 func (mv *Validator) deepValidateCollection(f reflect.Value, fname string, m ErrorMap) {
        switch f.Kind() {
        case reflect.Struct, reflect.Interface, reflect.Ptr:
+               // struct的struct的struct，递归调用Validate去处理
                e := mv.Validate(f.Interface())
                if e, ok := e.(ErrorMap); ok && len(e) > 0 {
                        for j, k := range e {
                                m[fname+"."+j] = k
                        }
                }
        case reflect.Array, reflect.Slice:
+               // 对于每个元素递归调用deepValidateCollection
                for i := 0; i < f.Len(); i++ {
                        mv.deepValidateCollection(f.Index(i), fmt.Sprintf("%s[%d]", fname, i), m)
                }
        case reflect.Map:
+               // 对于key和map递归调用deepValidateCollection
                for _, key := range f.MapKeys() {
                        mv.deepValidateCollection(key, fmt.Sprintf("%s[%+v](key)", fname, key.Interface()), m) // validate the map key
                        value := f.MapIndex(key)
                        mv.deepValidateCollection(value, fmt.Sprintf("%s[%+v](value)", fname, key.Interface()), m)
                }
        }
 }

 // Valid validates a value based on the provided
 // tags and returns errors found or nil.
+// 这个验证的范围更大，可以在所有类型上验证，并且可以指定tag的值
 func Valid(val interface{}, tags string) error {
        return defaultValidator.Valid(val, tags)
 }

 // Valid validates a value based on the provided
 // tags and returns errors found or nil.
 func (mv *Validator) Valid(val interface{}, tags string) error {
        if tags == "-" {
+               // 跳过 -
                return nil
        }
        v := reflect.ValueOf(val)
+       // 对于 reflect.Value.Kind() 为指针的处理方式，递归 .Elem().Interface()
        if v.Kind() == reflect.Ptr && !v.IsNil() {
                return mv.Valid(v.Elem().Interface(), tags)
        }
        var err error
        switch v.Kind() {
        case reflect.Invalid:
                err = mv.validateVar(nil, tags)
        default:
                err = mv.validateVar(val, tags)
        }
        return err
 }

 // validateVar validates one single variable
 func (mv *Validator) validateVar(v interface{}, tag string) error {
        tags, err := mv.parseTags(tag)
        if err != nil {
                // unknown tag found, give up.
                return err
        }
        errs := make(ErrorArray, 0, len(tags))
        for _, t := range tags {
                if err := t.Fn(v, t.Param); err != nil {
                        errs = append(errs, err)
                }
        }
        if len(errs) > 0 {
                return errs
        }
        return nil
 }

 // tag represents one of the tag items
+// 一个tag计算式，有名字，函数，参数
 type tag struct {
        Name  string         // name of the tag
        Fn    ValidationFunc // validation function to call
        Param string         // parameter to send to the validation function
 }

 // separate by no escaped commas
 var sepPattern *regexp.Regexp = regexp.MustCompile(`((?:^|[^\\])(?:\\\\)*),`)

 func splitUnescapedComma(str string) []string {
        ret := []string{}
        indexes := sepPattern.FindAllStringIndex(str, -1)
        last := 0
        for _, is := range indexes {
                ret = append(ret, str[last:is[1]-1])
                last = is[1]
        }
        ret = append(ret, str[last:])
        return ret
 }

 // parseTags parses all individual tags found within a struct tag.
 func (mv *Validator) parseTags(t string) ([]tag, error) {
        tl := splitUnescapedComma(t)
+       fmt.Printf("tl %v\n", tl)
        tags := make([]tag, 0, len(tl))
        for _, i := range tl {
                i = strings.Replace(i, `\,`, ",", -1)
                tg := tag{}
                v := strings.SplitN(i, "=", 2)
                tg.Name = strings.Trim(v[0], " ")
                if tg.Name == "" {
                        return []tag{}, ErrUnknownTag
                }
                if len(v) > 1 {
                        tg.Param = strings.Trim(v[1], " ")
                }
                var found bool
                if tg.Fn, found = mv.validationFuncs[tg.Name]; !found {
                        return []tag{}, ErrUnknownTag
                }
                tags = append(tags, tg)

        }
        return tags, nil
 }
 
从字段校验的需求来讲，无论我们采用深度优先搜索还是广度优先搜索来对这棵结构体树来进行遍历，都是可以的。


通过golang的structTag来配置验证器
type Class struct {
    Cid       int64  `validate:"required||integer=10000,_"`
    Cname     string `validate:"required||string=1,5||unique"`
    BeginTime string `validate:"required||datetime=H:i"`
}

type Student struct {
    Uid          int64    `validate:"required||integer=10000,_"`
    Name         string   `validate:"required||string=1,5"`
    Age          int64    `validate:"required||integer=10,30"`
    Sex          string   `validate:"required||in=male,female"`
    Email        string   `validate:"email||user||vm"`
    PersonalPage string   `validate:"url"`
    Hobby        []string `validate:"array=_,2||unique||in=swimming,running,drawing"`
    CreateTime   string   `validate:"datetime"`
    Class        []Class  `validate:"array=1,3"`
}
required 判断字段对应的值是否是对应类型的零值
integer 表示字段类型是否是整数类型，如果integer后边不接=?,?，那么表示只判断是否是整数类型，如果后边接=?,?，那么有四种写法
(1). integer=10 表示字段值 = 10
(2). integer=_ ,10 表示字段值 <= 10，字段值最小值为字段对应类型的最小值(比如字段对应类型为int8，那么最小为−128)，最大值为10
(3). integer=10, _ 表示字段值 >= 10，字段值最小值为10，最大值为字段对应类型的最大值(比如字段对应类型为int8，那么最大为127)
(4). integer=1,20 表示字段值 >=1 并且 <= 20
array、string 同 integer，array=?,? 表示元素个数范围，string=?,? 表示字符串长度范围
email 表示字段值是否是合法的email地址
url 表示字段值是否是合法的url地址
in 表示字段值在in指定的值中，比如 Hobby 字段中，in=swimming,running,drawing，表示 Hobby 字段的值，只能是swimming,running,drawing中的一个或多个
datetime 表示字段值符合日期类型，如果datetime后边不接=?，那么默认为Y-m-d H:i:s，否则验证器会按照指定格式判断，比如 datetime=Y-m、datetime=Y/m/d H:i:s等，可以是Y m d H i s 的随意拼接
unique 表示字段值唯一，比如 Hobby 字段的 unique，表示 Hobby 字段值唯一，Class 中，Cname 字段的 unique，表示 Cname 字段值唯一



