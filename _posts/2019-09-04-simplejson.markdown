---
title: simplejson  json.Decoder vs json.Unmarshal
layout: post
category: golang
author: 夏泽民
---
1,golang自带的json解析库encoding/json提供了json字符串到json对象的相互转换，在json字符串比较简单的情况下还是挺好用的，但是当json字符串比较复杂或者嵌套比较多的时候，就显得力不从心了，不可能用encoding/json那种为每个嵌套字段定义一个struct类型的方式，这时候使用simplejson库能够很方便的解析。

2，当被解析的json数据不一定完整的时候，使用标准库经常会解析失败，但是解析部分数据也是我们能接受的，这时可以用simplejson
<!-- more -->
可以看到，基本思路是将数据解析进一个interface｛｝，然后进行类型推断。

底层还是用的标准库

func NewJson(body []byte) (*Json, error) {
  j := new(Json)
  err := j.UnmarshalJSON(body)
  }
  
  func (j *Json) UnmarshalJSON(p []byte) error {
  dec := json.NewDecoder(bytes.NewBuffer(p))
  dec.UseNumber()
  return dec.Decode(&j.data)
}
  
  func (j *Json) Map() (map[string]interface{}, error) {
  if m, ok := (j.data).(map[string]interface{}); ok {
    return m, nil
    
    
    
func (j *Json) Array() ([]interface{}, error) {
  if a, ok := (j.data).([]interface{}); ok {
    return a, nil  
  


 json.Decoder vs json.Unmarshal
son的反序列化方式有两种：



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



从文件中读入一个巨大的json数组用json.Decoder

json.Decoder会一个一个元素进行加载，不会把整个json数组读到内存里面

从文件中读入json流用json.Decode
本来就以[]byte存在于内存中的用json.Unmarshal
