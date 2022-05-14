---
title: php_json_encode_zval
layout: post
category: web
author: 夏泽民
---
son解析在php，或者说在任何编程语言中都非常常见。但是，你是否知道，json解析在php中是以扩展形式存在。

json处理，最常用的就是json_encode,json_decode。

string json_encode ( mixed $value [, int $options = 0 [, int $depth = 512 ]] )
1
json_encode接受三个参数，但是一般的，我们都是使用一个参数，顶多会使用第二个参数，设置中文不转义，那其他的还有什么呢。

选项	说明
JSON_FORCE_OBJECT	使一个非关联数组输出一个类（Object）而非数组。
JSON_NUMERIC_CHECK	将所有数字字符串编码成数字
JSON_UNESCAPED_UNICODE	以字面编码多字节 Unicode 字符(不使用\u形式编码)
JSON_PRETTY_PRINT	用空白字符格式化返回的数据
json_encode多个选项使用的是多个选项进行或运算得到。json_encode($value,JSON_FORCE_OBJECT|JSON_NUMERIC_CHECK|JSON_UNESCAPED_UNICODE)表示如果空的时候，返回对象。数字返回数字类型，不编码。


$data = ['code'=>0,'data'=>['money'=>'0.12','name'=>'test 你好','id'=>1,'info'=>[] ] ];
$res1 = json_encode($data);
$res2 = json_encode($data,JSON_FORCE_OBJECT);
$res3 = json_encode($data,JSON_NUMERIC_CHECK);
$res4 = json_encode($data,JSON_UNESCAPED_UNICODE);
$res5 = json_encode($data,JSON_FORCE_OBJECT|JSON_NUMERIC_CHECK|JSON_UNESCAPED_UNICODE);
<!-- more -->
上面几个选项对于api中特别重要。像java强类型语言，数据解析如果没做类型判断的就容易导致程序崩溃。虽然也可以强制所有数据都是字符串类型，但是解析过程占用内存就会增加。

7位的整数数字类型数据，如果使用整数的话，占用24bit，3个字节就够了。但是如果是字符串解析则需要7个字节。使用合理的类型对json数据进行编码，既减少了客户端解码后数据内存的占用，也可以减少传输带宽。

但是，有个问题需要注意，JSON_NUMERIC_CHECK是对数字类型数据进行检查。如果数据是['orderid'=>'123456789009876553431234567788909987543886532313344455455']类似这种数据，全部由数字组成，也会转换成数字类型,并以科学计数方式输出 {"orderid":1.2345678900988e+56}，但实际上这种类型在表示成数字类型已经不合适了。

php中json_encode默认对空数组编码后返回的是数组形式。在某些场景下就容易产生问题。例如用户的一些附加属性，只有用户设置了才存在。当用户没有设置的时候，应该是一个对象返回，而不是数组。所以需要对这样的数据进行特殊处理，强制空数组返回对象。但是JSON_FORCE_OBJECT还是很危险的。使用它，会把本来是一个数据列表的空数组转换成对象。所以对于空数组的处理，要根据返回的数据进行特殊处理。如果正常数据是一个对象，则在encode的时候添加JSON_FORCE_OBJECT选项，如果是数组则比添加。但是要注意，JSON_FORCE_OBJECT影响的不仅经是最外层的数据，对于整个json串中所有符合条件的数据都会处理。因此最好的办法是还是单独处理，使用(object)对数据进行强制转换在编码，避免一刀切带来的问题。

json_encode最后一个参数是depth，表示迭代深度。php中json解析是一个递归过程，需要控制最大递归次数。默认限制是512。所以，如果你不设置第三个参数，让php对一个深度为512维的数组进行编码，得到的结果是false,错误提示为:" Maximum stack depth exceeded "

查看php源码中json扩展的内容json_encode.c文件,递归出现在encode的时候。每次进入json_encode_array中层级加1，如果递归次数超过配置次数，直接返回FAILURE。

int php_json_encode_zval(smart_str *buf, zval *val, int options, php_json_encoder *encoder) 
{
	again:
	switch (Z_TYPE_P(val))
	{
		...
		case IS_ARRAY:
			return php_json_encode_array(buf, val, options, encoder);
		...
	}
	return SUCCESS;
}

static int php_json_encode_array(smart_str *buf, zval *val, int options, php_json_encoder *encoder)
{
	...
	++encoder->depth;
	...
	if (php_json_encode_zval(buf, data, options, encoder) == FAILURE &&
					!(options & PHP_JSON_PARTIAL_OUTPUT_ON_ERROR)) {
				PHP_JSON_HASH_APPLY_PROTECTION_DEC(tmp_ht);
				return FAILURE;
	}
	...
	if (encoder->depth > encoder->max_depth) {
		encoder->error_code = PHP_JSON_ERROR_DEPTH;
		if (!(options & PHP_JSON_PARTIAL_OUTPUT_ON_ERROR)) {
			return FAILURE;
		}
	}
	--encoder->depth;
}
string json_decode ( mixed $value [, bool $assoc = false ] [, int $options = 0 ] [, int $depth = 512 ] )

json_decode 的一般使用都是将json转成数组,但是实际上json_encode接受4个参数。除了第二个参数用于标记是否返回数组之外，另外两个参数与json_encode一样。当解析的长度大于depth的时候，json_encode返回false。当json_encode 设置的depth > json_decode 的depth,json_decode返回false,无法正确解析json数据。相反的情况则可以。

整体而言，json_encode提供的option选项和depth选项，在我们明确知道自己在干什么的时候是非常有用的。但是一定要encode,decode使用相同方式。同时注意各种option可能代理的问题才能避免产生bug.

https://segmentfault.com/a/1190000020893651?utm_source=tag-newest

https://blog.csdn.net/agangdi/article/details/48197095
https://www.laruence.com/2009/04/28/719.html


