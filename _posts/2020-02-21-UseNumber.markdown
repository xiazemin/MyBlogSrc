---
title: UseNumber
layout: post
category: golang
author: 夏泽民
---
json的反序列化方式有两种：

Use json.Unmarshal passing the entire response string
// func Unmarshal(data []byte, v interface{}) error
data, err := ioutil.ReadAll(resp.Body)
if err == nil && data != nil {
    err = json.Unmarshal(data, value)
}
using json.NewDecoder.Decode
// func NewDecoder(r io.Reader) *Decoder
// func (dec *Decoder) Decode(v interface{}) error
err = json.NewDecoder(resp.Body).Decode(value)
这两种方法看似差不多，但有不同的应用场景

Use json.Decoder if your data is coming from an io.Reader stream, or you need to decode multiple values from a stream of data.

For the case of reading from an HTTP request, I’d pick json.Decoder since you’re obviously reading from a stream.

Use json.Unmarshal if you already have the JSON data in memory.

例子
从文件中读入一个巨大的json数组用json.Decoder

https://www.cnblogs.com/276815076/p/8583589.html
<!-- more -->
37. 将 JSON 中的数字解码为 interface 类型
在 encode/decode JSON 数据时，Go 默认会将数值当做 float64 处理，比如下边的代码会造成 panic：

func main() {
    var data = []byte(`{"status": 200}`)
    var result map[string]interface{}

    if err := json.Unmarshal(data, &result); err != nil {
        log.Fatalln(err)
    }

    fmt.Printf("%T\n", result["status"])    // float64
    var status = result["status"].(int)    // 类型断言错误
    fmt.Println("Status value: ", status)
}
panic: interface conversion: interface {} is float64, not int
如果你尝试 decode 的 JSON 字段是整型，你可以：

将 int 值转为 float 统一使用
将 decode 后需要的 float 值转为 int 使用
// 将 decode 的值转为 int 使用
func main() {
    var data = []byte(`{"status": 200}`)
    var result map[string]interface{}

    if err := json.Unmarshal(data, &result); err != nil {
        log.Fatalln(err)
    }

    var status = uint64(result["status"].(float64))
    fmt.Println("Status value: ", status)
}
使用 Decoder 类型来 decode JSON 数据，明确表示字段的值类型
// 指定字段类型
func main() {
    var data = []byte(`{"status": 200}`)
    var result map[string]interface{}
    
    var decoder = json.NewDecoder(bytes.NewReader(data))
    decoder.UseNumber()

    if err := decoder.Decode(&result); err != nil {
        log.Fatalln(err)
    }

    var status, _ = result["status"].(json.Number).Int64()
    fmt.Println("Status value: ", status)
}

 // 你可以使用 string 来存储数值数据，在 decode 时再决定按 int 还是 float 使用
 // 将数据转为 decode 为 string
 func main() {
     var data = []byte({"status": 200})
      var result map[string]interface{}
      var decoder = json.NewDecoder(bytes.NewReader(data))
      decoder.UseNumber()
      if err := decoder.Decode(&result); err != nil {
          log.Fatalln(err)
      }
    var status uint64
      err := json.Unmarshal([]byte(result["status"].(json.Number).String()), &status);
    checkError(err)
       fmt.Println("Status value: ", status)
}
​- 使用 struct 类型将你需要的数据映射为数值型

// struct 中指定字段类型
func main() {
      var data = []byte(`{"status": 200}`)
      var result struct {
          Status uint64 `json:"status"`
      }

      err := json.NewDecoder(bytes.NewReader(data)).Decode(&result)
      checkError(err)
    fmt.Printf("Result: %+v", result)
}
可以使用 struct 将数值类型映射为 json.RawMessage 原生数据类型
适用于如果 JSON 数据不着急 decode 或 JSON 某个字段的值类型不固定等情况：

// 状态名称可能是 int 也可能是 string，指定为 json.RawMessage 类型
func main() {
    records := [][]byte{
        []byte(`{"status":200, "tag":"one"}`),
        []byte(`{"status":"ok", "tag":"two"}`),
    }

    for idx, record := range records {
        var result struct {
            StatusCode uint64
            StatusName string
            Status     json.RawMessage `json:"status"`
            Tag        string          `json:"tag"`
        }

        err := json.NewDecoder(bytes.NewReader(record)).Decode(&result)
        checkError(err)

        var name string
        err = json.Unmarshal(result.Status, &name)
        if err == nil {
            result.StatusName = name
        }

        var code uint64
        err = json.Unmarshal(result.Status, &code)
        if err == nil {
            result.StatusCode = code
        }

        fmt.Printf("[%v] result => %+v\n", idx, result)
    }
}



有的时候上游传过来的字段是string类型的，但是我们却想用变成数字来使用。 本来用一个json:”,string” 就可以支持了，如果不知道golang的这些小技巧，就要大费周章了。

1）临时忽略struct字段

type User struct {
     Email    string `json:"email"`
     Password string `json:"password"`
    // many more fields… }

2）临时忽略掉Password字段

json.Marshal(struct {
     *User
     Password bool `json:"password,omitempty"` }{
     User: user, })

3）临时添加额外的字段

type User struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    // many more fields…
}

4）临时忽略掉Password字段，并且添加token字段

json.Marshal(struct {
    *User
    Token    string `json:"token"`
    Password bool `json:"password,omitempty"`
}{
    User: user,
    Token: token,
})

5）临时粘合两个struct

type BlogPost struct {
    URL   string `json:"url"`
    Title string `json:"title"`
}

type Analytics struct {
    Visitors  int `json:"visitors"`
    PageViews int `json:"page_views"`
}

json.Marshal(struct{
    *BlogPost
    *Analytics
}{post, analytics})

6）一个json切分成两个struct

json.Unmarshal([]byte(`{
  "url": "attila@attilaolah.eu",
  "title": "Attila's Blog",
  "visitors": 6,
  "page_views": 14
}`), &struct {
  *BlogPost
  *Analytics
}{&post, &analytics})

7）临时改名struct的字段

type CacheItem struct {
    Key    string `json:"key"`
    MaxAge int    `json:"cacheAge"`
    Value  Value  `json:"cacheValue"`
}

json.Marshal(struct{
    *CacheItem

    // Omit bad keys
    OmitMaxAge omit `json:"cacheAge,omitempty"`
    OmitValue  omit `json:"cacheValue,omitempty"`

    // Add nice keys
    MaxAge int    `json:"max_age"`
    Value  *Value `json:"value"`
}{
    CacheItem: item,

    // Set the int by value:
    MaxAge: item.MaxAge,

    // Set the nested struct by reference, avoid making a copy:
    Value: &item.Value,
})

8）用字符串传递数字

type TestObject struct {
    Field1 int    `json:",string"`
}

这个对应的json是 {“Field1”: “100”}

如果json是 {“Field1”: 100} 则会报错

容忍字符串和数字互转

如果你使用的是jsoniter，可以启动模糊模式来支持 PHP 传递过来的 JSON。

import “github.com/json-iterator/go/extra”

extra.RegisterFuzzyDecoders()
1
这样就可以处理字符串和数字类型不对的问题了。比如

var val string
jsoniter.UnmarshalFromString(`100`, &val)

又比如

var val float32
jsoniter.UnmarshalFromString(`"1.23"`, &val)

9）容忍空数组作为对象

PHP另外一个令人崩溃的地方是，如果 PHP array是空的时候，序列化出来是[]。但是不为空的时候，序列化出来的是{“key”:”value”}。 我们需要把 [] 当成 {} 处理。

如果你使用的是jsoniter，可以启动模糊模式来支持 PHP 传递过来的 JSON。

import "github.com/json-iterator/go/extra"

extra.RegisterFuzzyDecoders()

这样就可以支持了

var val map[string]interface{}
jsoniter.UnmarshalFromString(`[]`, &val)

10）使用 MarshalJSON支持time.Time

golang 默认会把 time.Time 用字符串方式序列化。如果我们想用其他方式表示 time.Time，需要自定义类型并定义 MarshalJSON。

type timeImplementedMarshaler time.Time

func (obj timeImplementedMarshaler) MarshalJSON() ([]byte, error) {
    seconds := time.Time(obj).Unix()
    return []byte(strconv.FormatInt(seconds, 10)), nil
}

11）序列化的时候会调用 MarshalJSON

type TestObject struct {
    Field timeImplementedMarshaler
}
should := require.New(t)
val := timeImplementedMarshaler(time.Unix(123, 0))
obj := TestObject{val}
bytes, err := jsoniter.Marshal(obj)
should.Nil(err)
should.Equal(`{"Field":123}`, string(bytes))

12）使用 RegisterTypeEncoder支持time.Time

jsoniter 能够对不是你定义的type自定义JSON编解码方式。比如对于 time.Time 可以用 epoch int64 来序列化

import "github.com/json-iterator/go/extra"
1
extra.RegisterTimeAsInt64Codec(time.Microsecond)
output, err := jsoniter.Marshal(time.Unix(1, 1002))
should.Equal("1000001", string(output))

如果要自定义的话，参见 RegisterTimeAsInt64Codec 的实现代码

13）使用 MarshalText支持非字符串作为key的map

虽然 JSON 标准里只支持 string 作为 key 的 map。但是 golang 通过 MarshalText() 接口，使得其他类型也可以作为 map 的 key。例如

f, _, _ := big.ParseFloat("1", 10, 64, big.ToZero)
val := map[*big.Float]string{f: "2"}
str, err := MarshalToString(val)
should.Equal(`{"1":"2"}`, str)

其中 big.Float 就实现了 MarshalText()

14）使用 json.RawMessage

如果部分json文档没有标准格式，我们可以把原始的文本信息用string保存下来。

type TestObject struct {
    Field1 string
    Field2 json.RawMessage
}
var data TestObject
json.Unmarshal([]byte(`{"field1": "hello", "field2": [1,2,3]}`), &data)
should.Equal(` [1,2,3]`, string(data.Field2))

15）使用 json.Number

默认情况下，如果是 interface{} 对应数字的情况会是 float64 类型的。如果输入的数字比较大，这个表示会有损精度。所以可以 UseNumber() 启用 json.Number 来用字符串表示数字。

decoder1 := json.NewDecoder(bytes.NewBufferString(`123`))
decoder1.UseNumber()
var obj1 interface{}
decoder1.Decode(&obj1)
should.Equal(json.Number("123"), obj1)

jsoniter 支持标准库的这个用法。同时，扩展了行为使得 Unmarshal 也可以支持 UseNumber 了。

json := Config{UseNumber:true}.Froze()
var obj interface{}
json.UnmarshalFromString("123", &obj)
should.Equal(json.Number("123"), obj)

16）统一更改字段的命名风格

经常 JSON 里的字段名 Go 里的字段名是不一样的。我们可以用 field tag 来修改。

output, err := jsoniter.Marshal(struct {
    UserName      string `json:"user_name"`
    FirstLanguage string `json:"first_language"`
}{
    UserName:      "taowen",
    FirstLanguage: "Chinese",
})
should.Equal(`{"user_name":"taowen","first_language":"Chinese"}`, string(output))

但是一个个字段来设置，太麻烦了。如果使用 jsoniter，我们可以统一设置命名风格。

import "github.com/json-iterator/go/extra"

extra.SetNamingStrategy(LowerCaseWithUnderscores)
output, err := jsoniter.Marshal(struct {
    UserName      string
    FirstLanguage string
}{
    UserName:      "taowen",
    FirstLanguage: "Chinese",
})
should.Nil(err)
should.Equal(`{"user_name":"taowen","first_language":"Chinese"}`, string(output))

17）使用私有的字段

Go 的标准库只支持 public 的 field。jsoniter 额外支持了 private 的 field。需要使用 SupportPrivateFields() 来开启开关。

import "github.com/json-iterator/go/extra"

extra.SupportPrivateFields()
type TestObject struct {
    field1 string
}
obj := TestObject{}
jsoniter.UnmarshalFromString(`{"field1":"Hello"}`, &obj)
should.Equal("Hello", obj.field1)

文中所用第三方库：https://github.com/json-iterator/go


