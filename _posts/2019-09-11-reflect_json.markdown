---
title: reflect性能
layout: post
category: golang
author: 夏泽民
---
Go reflect包提供了运行时获取对象的类型和值的能力，它可以帮助我们实现代码的抽象和简化，实现动态的数据获取和方法调用， 提高开发效率和可读性， 也弥补Go在缺乏泛型的情况下对数据的统一处理能力。

通过reflect，我们可以实现获取对象类型、对象字段、对象方法的能力，获取struct的tag信息，动态创建对象，对象是否实现特定的接口，对象的转换、对象值的获取和设置、Select分支动态调用等功能
<!-- more -->
Java中的reflect的使用对性能也有影响， 但是和Java reflect不同， Java中不区分Type和Value类型的， 所以至少Java中我们可以预先将相应的reflect对象缓存起来，减少反射对性能的影响， 但是Go没办法预先缓存reflect, 因为Type类型并不包含对象运行时的值，必须通过ValueOf和运行时实例对象才能获取Value对象。

对象的反射生成和获取都会增加额外的代码指令， 并且也会涉及interface{}装箱/拆箱操作，中间还可能增加临时对象的生成，所以性能下降是肯定的

如果通过反射进行赋值，性能下降是很厉害的，耗时成倍的增长。比较有趣的是，FieldByName方式赋值是Field方式赋值的好几倍， 原因在于FieldByName会有额外的循环进行字段的查找，虽然最终它还是调用Field进行赋值：

 通过反射生成对象和字段赋值都会影响性能，但是通过反射的确确确实实能简化代码，为业务逻辑提供统一的代码， 比如标准库中json的编解码、rpc服务的注册和调用， 一些ORM框架比如gorm等，都是通过反射处理数据的，这是为了能处理通用的类型。
 
在我们追求高性能的场景的时候，我们可能需要尽量避免反射的调用， 比如对json数据的unmarshal， easyjson就通过生成器的方式，避免使用反射。
func (v *Student) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a74e62dDecodeGitABCReflect(&r, v)
	return r.Error()
}
func (v *Student) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a74e62dDecodeGitABCReflect(l, v)
}
func easyjson4a74e62dDecodeGitABCReflect(in *jlexer.Lexer, out *Student) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Name":
			out.Name = string(in.String())
		case "Age":
			out.Age = int(in.Int())
		case "Class":
			out.Class = string(in.String())
		case "Score":
			out.Score = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

将具体对象转换成 interface{}(以及反向操作)确实回带来一点点性能的影响，不过看起来影响倒不是很大。

reflect.Type & reflect.Value
reflect.Type
TypeOf returns the reflection Type that represents the dynamic type of i. If i is a nil interface value, TypeOf returns nil.

func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}

t := reflect.TypeOf(3)   // a reflect.Type
fmt.Println(t.String())  // "int"
fmt.Println(t)           // "int"

var w io.Writer = os.Stdout
fmt.Println(reflect.TypeOf(w)) // "*os.File"
reflect.TypeOf 返回的是i 具体的类型。 所以w的类型是*os.File而不是io.Writer。

为方便进行debug，可以用fmt.Printf(“%T\n”, 3) 来代替 fmt.Println(reflect.TypeOf(3))。

reflect.Value
As with reflect.TypeOf, the results of reflect.ValueOf are always concrete. 和reflect.TypeOf一样，reflect.ValueOf返回也是一个具体的值。

v := reflect.ValueOf(3)   // a reflect.Value
fmt.Println(v)  // "3"
fmt.Printf("%v\n",v)  // "3"
fmt.Println(v.String())  // NOTE: "<int Value>"
Kind
Although there are infinitely many types, there are only a finite number of kinds of type: the basic types Bool, String and all the numbers; the aggregate types Array and Struct, the reference types Chan, Func, Ptr, Slice and Map; interface types; and finally Invalid, meaning no value at all (The zero value of a reflect.Value has kind Invalid).

尽管具体的数据类型可以有无限种，但是它们可以被分为几种类型。这个就是reflect.Kind.

// Any formats any value as a string.
func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
		// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

func main() {
	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	fmt.Println(Any(x))                  // "1"
	fmt.Println(Any(d))                  // "1"
	fmt.Println(Any([]int64{x}))         // "[]int64 0x8202b87b0"
	fmt.Println(Any([]time.Duration{d})) // "[]time.Duration 0x8202b87e0"
}


