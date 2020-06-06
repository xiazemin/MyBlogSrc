---
title: ZEND_HASH_FOREACH_KEY_VAL php扩展hashtable操作
layout: post
category: php
author: 夏泽民
---
php7的遍历数组和php5差很多，7提供了一些专门的宏来遍历元素（或keys）。宏的第一个参数是HashTable，其他的变量被分配到每一步迭代：

ZEND_HASH_FOREACH_VAL(ht, val)
ZEND_HASH_FOREACH_KEY(ht, h, key)
ZEND_HASH_FOREACH_PTR(ht, ptr)
ZEND_HASH_FOREACH_NUM_KEY(ht, h)
ZEND_HASH_FOREACH_STR_KEY(ht, key)
ZEND_HASH_FOREACH_STR_KEY_VAL(ht, key, val)
ZEND_HASH_FOREACH_KEY_VAL(ht, h, key, val)
<!-- more -->
https://segmentfault.com/a/1190000007575322

https://www.cnblogs.com/CoderK/articles/6943291.html

http://nikic.github.io/2014/12/22/PHPs-new-hashtable-implementation.html

1.1     使用数组
曾讲到，PHP数组本质上就是个HashTable，因此访问数组就是对HashTable进行操作，Zend为我们提供的一组数组函数也只是对HashTable操作进行了简单包装而已。

来看创建数组，由于数组也是存在于zval里的，因此要先用MAKE_STD_ZVAL()宏创建一个zval，之后调用如下宏将其转化为一个空数组：

array_init(zval*)
接下来是朝数组中添加元素，这对关联数组元素和非关联数组元素要采用不同操作。

1.1.1 关联数组元素
关联数组采用char*作为key，zval*作为value，可以使用如下宏将已有的zval加入数组或者更新已有元素：

int add_assoc_zval(zval *arr, char *key, zval *value)
需要特别注意的是，Zend不会复制zval，只会简单的储存其指针，并且不关心任何引用计数，因此不能将其他变量的zval或者是栈上的zval传给它，只能用MAKE_STD_ZVAL()宏构建。

Zend为常用的类型定义了相应的API，以简化我们的操作：

add_assoc_long(zval *array, char *key, long n);
add_assoc_bool(zval *array, char *key, int b);
add_assoc_resource(zval *array, char *key, int r);
add_assoc_double(zval *array, char *key, double d);
add_assoc_string(zval *array, char *key, char *str, int duplicate);
add_assoc_stringl(zval *array, char *key, char *str, uint length, int duplicate);
add_assoc_null(zval *array, char *key);
当函数发现目标元素已经存在时，会首先递减其原zval的refcount，然后才插入新zval，这就保证了原zval引用信息的正确性。这种行为是通过HashTable.pDestructor（参见1.2.1）实现的，每次删除一个元素时，HashTable都将对被删元素调用这个函数指针，而数组为其HashTable设置的函数指针就是用来处理被删除zval的引用信息。

另外，查看这些函数的源代码可以发现一个有意思的现象，它们没有直接使用HashTable操作，而是使用变量符号表操作，可见关联数组和变量符号表就是一种东西。

Zend没有提供删除和获取数组元素的函数，此类操作只能使用HashTable函数或者是2.6节的变量符号表操作。

1.1.2非关联数组元素
非关联数组没有key，使用index作为hash，相应函数和上面关联数组的十分类似：

add_index_zval(zval *array, uint idx, zval *value);
add_index_long(zval *array, uint idx, long n);
add_index_bool(zval *array, uint idx, int b);
add_index_resource(zval *array, uint idx, int r);
add_index_double(zval *array, uint idx, double d);
add_index_string(zval *array, uint idx, char *str, int duplicate);
add_index_stringl(zval *array, uint idx, char *str, uint length, int duplicate);
add_index_null(zval *array, uint idx);
如果只是想插入值，而不指定index的话，可以使用如下函数：

add_next_index_zval(zval *array, zval *value);
add_next_index_long(zval *array, long n);
add_next_index_bool(zval *array, int b);
add_next_index_resource(zval *array, int r);
add_next_index_double(zval *array, double d);
add_next_index_string(zval *array, char *str, int duplicate);
add_next_index_stringl(zval *array, char *str, uint length, int duplicate);
add_next_index_null(zval *array);
1.2      使用资源
1.2.1  注册资源类型
1.1.1节曾经提到，所谓资源就是内部数据的handle（但是这句话并不全对），使用资源是比较简单的，首先是注册一个资源类型：

int zend_register_list_destructors_ex(
rsrc_dtor_func_t ld,
rsrc_dtor_func_t pld,
char *type_name,
int module_number);
第一个参数是函数指针，当资源不再被使用或者模块将被卸载时，Zend使用它来销毁资源，稍候再作介绍；第二个参数和第一个类似，只是它被用来销毁持久性资源(*)；type_name是资源名称，用户可以使用var_dump函数来读取；module_number是模块号，在启动函数中可以获取该值。

注册过程其实就是将我们传入的参数放到一个内部数据结构，然后把这个数据结构放入一个没有使用key的HashTable里，该函数返回的值，也就是所谓“资源类型id”，其实就是HashTable的index。

1.2.1  注册资源
注册完资源类型后，就可以注册一个该类型的资源了：

1
ZEND_REGISTER_RESOURCE(
2
rsrc_result,
3
rsrc_pointer,
4
rsrc_type)
src_pointer是个指针类型，就是你的资源的handle, 通常是指向内部数据的指针，当然也可以是index或者其它标志符；rsrc_type是上面获取的资源类型id；rsrc_result是个已有的zval，注册完成后，资源的id就被放入该zval，同时其type也被设为IS_RESOURCE，通常是传入return_value，以将资源返回给用户。

在内部，Zend使用如下数据结构表示一个资源：

1
typedef struct _zend_rsrc_list_entry {
2
    void *ptr;
3
    int type;
4
    int refcount;
5
} zend_rsrc_list_entry;
ptr和type就是我们在上面传入的参数；refcount是引用计数，由Zend维护，当引用减到0时，Zend会销毁该资源。不出所料的是，这个数据结构也被组织在一个HashTable里，并且没有使用key，仅仅使用index——这就是zval里存放的东西。现在资源的整个脉络已经清晰：通过zval可以获得资源id，通过资源id可以获得资源handle和资源类型id，通过资源类型id可以获得资源的销毁函数。
现在讲一下销毁函数：

1
typedef void (*rsrc_dtor_func_t)(
2
zend_rsrc_list_entry *rsrc
3
TSRMLS_DC);
rsrc是需要被销毁的资源，我们在函数的实现中可以通过它获得资源的handle，并且加以处理，比如释放内存块、关闭数据库连接或是关闭文件描述符等。

1.2.3  获取资源
当创建了资源后，用户通常都要调用创建者提供的函数来操作资源，此时我们需要从用户传入的zval中取出资源：

1
ZEND_FETCH_RESOURCE(
2
rsrc,  rsrc_type,
3
passed_id, default_id,
4
resource_type_name, resource_type)
首个参数用于接收handle值，第二个参数是handle值的类型，这个函数会扩展成“rsrc = (rsrc_type) zend_fetch_resource(…)”，因此应该保证rsrc是rsrc_type类型的；passed_id是用户传入的zval，这里使用zval**类型，函数从中取得资源id；default_id用来直接指定资源id，如果该值不是-1，则使用它，并且忽略passed_id，所以通常应该使用-1；resource_type_name是资源名称，当获取资源失败时，函数使用它来输出错误信息；resource_type是资源类型，如果取得的资源不是该类型的，则函数返回NULL，这用于防止用户传入一个其他类型资源的zval。

不过，这个宏确实比较难用，用其底层的宏反倒更加容易些：

1
zend_list_find(id, type)
id是要查找的资源id；type是int*类型，用于接收取出的资源的类型，可以用它来判断这是不是我们想要的资源；函数最后返回资源的handle，失败返回NULL。

1.2.4  维护引用计数
通常，当用户对资源类型的PHP变量执行赋值或是unset之类操作时，Zend会自动维护资源的引用计数。但有时，我们也需要手动进行，比如我们要复用一个数据库连接或者用户调用我们提供的close操作关闭一个文件，此时可以使用如下宏：

1
zend_list_addref(id)
2
zend_list_delete(id)
id是资源id，这两个宏分别增加和减少目标资源的引用计数，第二个宏还会在引用计数减到0时，调用先前注册的函数销毁资源。

https://blog.csdn.net/tonysz126/article/details/6993665

https://www.cnblogs.com/ling-diary/p/10676109.html

https://www.laruence.com/2020/02/25/3182.html

https://blog.csdn.net/ligupeng7929/article/details/90521059

https://stackoverflow.com/questions/56429192/convert-php-two-dimensional-arrays-to-php-extension/56453407?r=SearchResults#56453407

https://segmentfault.com/q/1010000007890184

https://segmentfault.com/a/1190000007575322

到这已经能声明简单函数，返回静态或者动态值了。定义INI选项，声明内部数值或全局数值。本章节将介绍如何接收从调用脚本(php文件)传入参数的数值，以及 PHP内核 和 Zend引擎 如何操作内部变量。

接收参数

与用户控件的代码不同，内部函数的参数实际上并不是在函数头部声明的，函数声明都形如: PHP_FUNCTION(func_name) 的形式，参数声明不在其中。参数的传入是通过参数列表的地址传入的，并且是传入每一个函数，不论是否存在参数。

通过定义函数hello_str()来看一下，它将接收一个参数然后把它与问候的文本一起输出。

PHP_FUNCTION(hello_greetme){ char*name = NULL; size_tname_len; zend_string *strg; if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "s", &name, &name_len) == FAILURE) { RETURN_NULL(); } strg = strpprintf( 0, "你好: %s", name); RETURN_STR(strg);}

大多数 zend_parse_parameters() 块看起来都差不多。 ZEND_NUM_ARGS() 告诉Zend引擎要取的参数的信息， TSRMLS_CC 用来确保线程安全，返回值检测是SUCCESS还是FAILURE。通常情况下返回是SUCCESS的。除非传入的参数太少或太多或者参数不能被转为适当的类型，Zend会自动输出一条错误信息并将控制权还给调用脚本。

指定 "s" 表明此函数期望只传入一个参数，并且该参数被转化为string数据类型，地址传入char * 变量。

此外，还有一个int变量通过地址传递到 zend_parse_parameters() 。这使Zend引擎提供字符串的字节长度，如此二进制安全的函数不再依赖strlen(name)来确定字符串的长度。因为实际上使用strlen(name)甚至得不到正确的结果，因为name可能在字符串结束之前包含了NULL字符。

在php7中，提供另一种获取参数的方式FAST_ZPP，是为了提高参数解析的性能。

#ifdefFAST_ZPPZEND_PARSE_PARAMETERS_START( 1, 2) Z_PARAM_STR(type) Z_PARAM_OPTIONALZ_PARAM_ZVAL_EX(value, 0, 1) ZEND_PARSE_PARAMETERS_END(); #endif

参数类型表

类型	代码	变量类型
Boolean	b	zend_bool
Long	l	long
Double	d	double
String	s	char*, int
Resource	r	zval *
Array	a	zval *
Object	o	zval *
zval	z	zval *
最后四个类型都是zvals *.这是因为在php的实际使用中，zval数据类型存储所有的用户空间变量。三种“复杂”数据类型：资源、数组、对象。当它们的数据类型代码被用于zend_parse_parameters()时，Zend引擎会进行类型检查，但是因为在C中没有与它们对应的数据类型，所以不会执行类型转换。

Zval

一般而言，zval和php用户空间变量是很伤脑筋的，概念很难懂。到了PHP7，它的结构在Zend/zend_types.h中有定义：

struct_zval_struct { zend_value value; /* value */union{ struct{ ZEND_ENDIAN_LOHI_4( zend_uchar type, /* active type */zend_uchar type_flags, zend_uchar const_flags, zend_uchar reserved) /* call info for EX(This) */} v; uint32_ttype_info; } u1; union{ uint32_tnext; /* hash collision chain */uint32_tcache_slot; /* literal cache slot */uint32_tlineno; /* line number (for ast nodes) */uint32_tnum_args; /* arguments number for EX(This) */uint32_tfe_pos; /* foreach position */uint32_tfe_iter_idx; /* foreach iterator index */uint32_taccess_flags; /* class constant access flags */uint32_tproperty_guard; /* single property guard */} u2;};

可以看到，变量是通过_zval_struct结构体存储的，而变量的值是zend_value类型的：

typedefunion_zend_value { zend_long lval; /* long value */doubledval; /* double value */zend_refcounted *counted; zend_string *str; zend_array *arr; zend_object *obj; zend_resource *res; zend_reference *ref; zend_ast_ref *ast; zval *zv; void*ptr; zend_class_entry *ce; zend_function *func; struct{ uint32_tw1; uint32_tw2; } ww;} zend_value;

虽然结构体看起来很大，但细细看，其实都是联合体，value的扩充，u1是type_info，u2是其他各种辅助字段。

zval 类型

变量存储的数据是有数据类型的，php7中总体有以下类型,Zend/zend_types.h中有定义：

/* regular data types */#defineIS_UNDEF 0#defineIS_NULL 1#defineIS_FALSE 2#defineIS_TRUE 3#defineIS_LONG 4#defineIS_DOUBLE 5#defineIS_STRING 6#defineIS_ARRAY 7#defineIS_OBJECT 8#defineIS_RESOURCE 9#defineIS_REFERENCE 10/* constant expressions */#defineIS_CONSTANT 11#defineIS_CONSTANT_AST 12/* fake types */#define_IS_BOOL 13#defineIS_CALLABLE 14#defineIS_ITERABLE 19#defineIS_VOID 18/* internal types */#defineIS_INDIRECT 15#defineIS_PTR 17#define_IS_ERROR 20

测试

书写一个类似gettype()来取得变量的类型的hello_typeof():

PHP_FUNCTION(hello_typeof){ zval *userval = NULL; if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "z", &userval) == FAILURE) { RETURN_NULL(); } switch(Z_TYPE_P(userval)) { caseIS_NULL: RETVAL_STRING( "NULL"); break; caseIS_TRUE: RETVAL_STRING( "true"); break; caseIS_FALSE: RETVAL_STRING( "false"); break; caseIS_LONG: RETVAL_STRING( "integer"); break; caseIS_DOUBLE: RETVAL_STRING( "double"); break; caseIS_STRING: RETVAL_STRING( "string"); break; caseIS_ARRAY: RETVAL_STRING( "array"); break; caseIS_OBJECT: RETVAL_STRING( "object"); break; caseIS_RESOURCE: RETVAL_STRING( "resource"); break; default: RETVAL_STRING( "unknown type"); }}

这里使用RETVAL_STRING()与之前的RETURN_STRING()差别并不大,它们都是宏。只不过RETURN_STRING中包含了RETVAL_STRING的宏代替，详细在 Zend/zend_API.h 中有定义:

#defineRETVAL_STRING(s) ZVAL_STRING(return_value, s)#defineRETVAL_STRINGL(s, l) ZVAL_STRINGL(return_value, s, l)#defineRETURN_STRING(s) { RETVAL_STRING(s); return; }#defineRETURN_STRINGL(s, l) { RETVAL_STRINGL(s, l); return; }

创建zval

前面用到的zval是由Zend引擎分配空间，也通过同样的途径释放。然而有时候需要创建自己的zval，可以参考如下代码：

{ zval temp; ZVAL_LONG(&temp, 1234);}

数组

数组作为运载其他变量的变量。内部实现上使用了众所周知的 HashTable .要创建将被返回PPHP的数组，最简单的方法：

PHP语法	C语法（arr是zval*）	意义
$arr = array();	array_init(arr);	初始化一个新数组
$arr[] = NULL;	add_next_index_null(arr);	向数字索引的数组增加指定类型的值
$arr[] = 42;	add_next_index_long(arr, 42);	
$arr[] = true;	add_next_index_bool(arr, 1);	
$arr[] = 3.14;	add_next_index_double(arr, 3.14);	
$arr[] = 'foo';	add_next_index_string(arr, "foo", 1);	
$arr[] = $myvar;	add_next_index_zval(arr, myvar);	
$arr[0] = NULL;	add_index_null(arr, 0);	向数组中指定的数字索引增加指定类型的值
$arr[1] = 42;	add_index_long(arr, 1, 42);	
$arr[2] = true;	add_index_bool(arr, 2, 1);	
$arr[3] = 3.14;	add_index_double(arr, 3, 3.14);	
$arr[4] = 'foo';	add_index_string(arr, 4, "foo", 1);	
$arr[5] = $myvar;	add_index_zval(arr, 5, myvar);	
$arr['abc'] = NULL;	add_assoc_null(arr, "abc");	
$arr['def'] = 711;	add_assoc_long(arr, "def", 711);	向关联索引的数组增加指定类型的值
$arr['ghi'] = true;	add_assoc_bool(arr, "ghi", 1);	
$arr['jkl'] = 1.44;	add_assoc_double(arr, "jkl", 1.44);	
$arr['mno'] = 'baz';	add_assoc_string(arr, "mno", "baz", 1);	
$arr['pqr'] = $myvar;	add_assoc_zval(arr, "pqr", myvar);	
做一个测试：

PHP_FUNCTION(hello_get_arr){ array_init(return_value); add_next_index_null(return_value); add_next_index_long(return_value, 42); add_next_index_bool(return_value, 1); add_next_index_double(return_value, 3.14); add_next_index_string(return_value, "foo"); add_assoc_string(return_value, "mno", "baz"); add_assoc_bool(return_value, "ghi", 1);}



add_*_string()函数参数从四个改为了三个。

数组遍历

假设我们需要一个取代以下功能的拓展：

<?phpfunctionhello_array_strings($arr){ if(!is_array($arr)) { returnNULL; } printf("The array passed contains %d elementsn", count($arr)); foreach($arr as$data) { if(is_string($data)) echo$data.'n'; }}

php7的遍历数组和php5差很多，7提供了一些专门的宏来遍历元素（或keys）。宏的第一个参数是HashTable，其他的变量被分配到每一步迭代：

ZEND_HASH_FOREACH_VAL(ht, val)

ZEND_HASH_FOREACH_KEY(ht, h, key)

ZEND_HASH_FOREACH_PTR(ht, ptr)

ZEND_HASH_FOREACH_NUM_KEY(ht, h)

ZEND_HASH_FOREACH_STR_KEY(ht, key)

ZEND_HASH_FOREACH_STR_KEY_VAL(ht, key, val)

ZEND_HASH_FOREACH_KEY_VAL(ht, h, key, val)

因此它的对应函数实现如下：

PHP_FUNCTION(hello_array_strings){ ulongnum_key; zend_string *key; zval *val, *arr; HashTable *arr_hash; intarray_count; if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "a", &arr) == FAILURE) { RETURN_NULL(); } arr_hash = Z_ARRVAL_P(arr); array_count = zend_hash_num_elements(arr_hash); php_printf( "The array passed contains %d elementsn", array_count); ZEND_HASH_FOREACH_KEY_VAL(arr_hash, num_key, key, val) { //if (key) { //HASH_KEY_IS_STRING//}PHPWRITE(Z_STRVAL_P(val), Z_STRLEN_P(val)); php_printf( "n"); }ZEND_HASH_FOREACH_END();}

因为这是新的遍历方法，而我看的还是php5的处理方式，调试出上面的代码花了不少功夫，总的来说，用宏的方式遍历大大减少了编码体积。哈希表是php中很重要的一个内容，有时间再好好研究一下。

遍历数组的其他方式

遍历 HashTable 还有其他方法。Zend引擎针对这个任务展露了三个非常类似的函数：zend_hash_apply(), zend_hash_apply_with_argument(), zend_hash_apply_with_arguments。第一个形式仅仅遍历HashTable，第二种形式允许传入单个void*参数，第三种形式通过var arg列表允许数量不限的参数。hello_array_walk()展示个他们各自的行为。

staticintphp_hello_array_walk(zval *ele TSRMLS_DC){ zval temp = *ele; // 临时zval，避免convert_to_string 污染原元素zval_copy_ctor(&temp); // 分配新 zval 空间并复制 ele 的值convert_to_string(&temp); // 字符串类型转换//简单的打印PHPWRITE(Z_STRVAL(temp), Z_STRLEN(temp)); php_printf( "n"); zval_dtor(&temp); //释放临时的 tempreturnZEND_HASH_APPLY_KEEP;} staticintphp_hello_array_walk_arg(zval *ele, char*greeting TSRMLS_DC){ php_printf( "%s", greeting); php_hello_array_walk(ele TSRMLS_CC); returnZEND_HASH_APPLY_KEEP;} staticintphp_hello_array_walk_args(zval *ele, intnum_args, va_list args, zend_hash_key *hash_key){ char*prefix = va_arg(args, char*); char*suffix = va_arg(args, char*); TSRMLS_FETCH(); php_printf( "%s", prefix); // 打印键值对结果php_printf( "key is : [ "); if(hash_key->key) { PHPWRITE(ZSTR_VAL(hash_key->key), ZSTR_LEN(hash_key->key)); } else{ php_printf( "%ld", hash_key->h); } php_printf( " ]"); php_hello_array_walk(ele TSRMLS_CC); php_printf( "%sn", suffix); returnZEND_HASH_APPLY_KEEP;}

用户调用的函数：

PHP_FUNCTION(hello_array_walk){ zval *arr; HashTable *arr_hash; if(zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "a", &arr) == FAILURE) { RETURN_NULL(); } arr_hash = Z_ARRVAL_P(arr); //第一种遍历 简单遍历各个元素zend_hash_apply(arr_hash, ( apply_func_t)php_hello_array_walk TSRMLS_CC); //第二种遍历 带一个参数的简单遍历各个元素zend_hash_apply_with_argument(arr_hash, ( apply_func_arg_t)php_hello_array_walk_arg, "Hello "TSRMLS_CC); //第三种遍历 带多参数的遍历key->valuezend_hash_apply_with_arguments(arr_hash, ( apply_func_args_t)php_hello_array_walk_args, 2, "Hello ", "Welcome to my extension!"); RETURN_TRUE;}

为了复用，在输出值时调用php_hello_array_walk(ele TSRMLS_CC)。传入hello_array_walk()的数组被遍历了三次,一次不带参数，一次带单个参数，一次带两给参数。三个遍历的函数返回了ZEND_HASH_APPLY_KEEP。这告诉zend_hash_apply()函数离开HashTable中的（当前）元素，继续处理下一个。

这儿也可以返回其他值：ZEND_HASH_APPLY_REMOVE删

除当前元素并继续应用到下一个；ZEND_HASH_APPLY_STOP在当前元素中止数组的遍历并退出zend_hash_apply()函数。

TSRMLS_FETCH() 是一个关于线程安全的动作，用于避免各线程的作用域被其他的侵入。因为zend_hash_apply()的多线程版本用了vararg列表，tsrm_ls标记没有传入walk()函数。

<?php$arr = ["99", "fff", "key1"=>"888", "key2"=>"aaa"];hello_array_walk($arr);

val
zval结构体是Zend内核的非常核心的结构，在PHP5和PHP7之间的差别非常大，我给出2处文章供大家学习，基本上可以代表这块知识点最权威的介绍了。

深入理解PHP7之zval（鸟哥）
变量在 PHP7 内部的实现（Nikita Popov）中文版
PHP7不再使用zval的二级指针，大多数场景下出现的zval*变量都改成zval，相应的使用在这些变量上的宏Z_PP也需要改成Z_P。
在大部分场景下，PHP7是在栈上直接使用zval，不需要去堆上分配内存。这时，zval 就需要改成zval，宏也需要从Z__P改成Z_，创建宏从ZVAL_(var)转换成ZVAL_*(&var)。所以，分配zval内存的宏

ALLOC_ZVAL、ALLOC_INIT_ZVAL、MAKE_STD_ZVAL都被删掉了。
- zval *zv;
- MAKE_STD_ZVAL(zv);
- array_init(zv);
+ zval zv;
+ array_init(&zv);
PHP7中zval的long和double类型是不需要引用计数的，所以相关的宏要做调整。

- Z_ADDREF_P(zv)
+ Z_TRY_ADDREF_P(zv);
PHP7中zval的类型，删除了IS_BOOL，增加了IS_TRUE和IS_FALSE。

- if (Z_TYPE_P(zv) == IS_BOOL) {
- }
+ if (Z_TYPE_P(zv) == IS_TRUE) {
+ } else if (Z_TYPE_P(zv) == IS_FALSE) {
+ }
zend_string
PHP7中增加了一个新的内置字符串类型zend_string，下面是Zend内核中的结构体定义。

struct _zend_string {
zend_refcounted_h gc; /* 垃圾回收结构体 */
zend_ulong h; /* 字符串哈希值 */
size_t len; /* 字符串长度 */
char val[1]; /* 字符串内容 */
};
gc是PHP7中的所有非标量结构都包含的垃圾回收结构体变量；h是字符串哈希值，作为HashTable的key时不需要每次都重新计算哈希值，提高了效率；len是字符串长度，同理每次使用到字符串的长度时不需要再计算，提高了效率；val[1]是C语言的黑科技，此处按照char *理解即可。这里有三个宏帮助我们方便的使用zend_string的变量。

#define ZSTR_VAL(zstr) (zstr)->val
#define ZSTR_LEN(zstr) (zstr)->len
#define ZSTR_H(zstr) (zstr)->h
创建和销毁zend_string使用以下方法。

zend_string *zend_string_init(const char *str, size_t len, int persistent)
void zend_string_release(zend_string *s)
zend_string用来替代PHP5中使用char *和int的场景，尤其是很多API的参数和返回值都做了调整。

- int zend_hash_find(const HashTable *ht, const char *arKey, uint nKeyLength, void **pData)
+ zval* ZEND_FASTCALL zend_hash_find(const HashTable *ht, zend_string *key)
- void zend_mangle_property_name(char **dest, int *dest_length, const char *src1, int src1_length, const char *src2, int src2_length, int internal);
+ zend_string *zend_mangle_property_name(const char *src1, size_t src1_length, const char *src2, size_t src2_length, int internal)
HashTable API
在PHP7中使用HashTable的API方法时，有了非常明显的变化。
查询方法，PHP5使用引用传参的方式，同时返回SUCCESS/FAILURE；PHP7直接返回结果，查询无结果时返回NULL。

- int zend_hash_find(const HashTable *ht, const char *arKey, uint nKeyLength, void **pData)
+ zval* ZEND_FASTCALL zend_hash_find(const HashTable *ht, zend_string *key)
HashTable的API方法中的key，PHP5中使用char 和int代表的字符串；PHP7中使用zend_string代表的字符串，同时提供了对char 和int支持的一组方法，但是需要注意的是这里的字符串长度是不包括结尾的'0'的，在升级扩展时难免会碰到很多地方需要加减一。

- int zend_hash_exists(const HashTable *ht, const char *arKey, uint nKeyLength)
+ zend_bool zend_hash_exists(const HashTable *ht, zend_string *key)
+ zend_bool zend_hash_str_exists(const HashTable *ht, const char *str, size_t len)
PHP7为HashTable的value为指针时设计了一组API，在常规的API方法后添加后缀_ptr即可。

void *zend_hash_find_ptr(const HashTable *ht, zend_string *key)
void *zend_hash_update_ptr(HashTable *ht, zend_string *key, void *pData)
PHP7为HashTable的轮询设计了一组宏，使用起来非常方便。

ZEND_HASH_FOREACH_VAL(ht, val)
ZEND_HASH_FOREACH_KEY(ht, h, key)
ZEND_HASH_FOREACH_PTR(ht, ptr)
ZEND_HASH_FOREACH_NUM_KEY(ht, h)
ZEND_HASH_FOREACH_STR_KEY(ht, key)
ZEND_HASH_FOREACH_STR_KEY_VAL(ht, key, val)
ZEND_HASH_FOREACH_KEY_VAL(ht, h, key, val)
自定义对象
这里有点复杂，我直接附上我的代码，结合代码来做详细说明。

typedef struct{
int max;
int offset;
zend_object zo;
} php_protocolbuffers_unknown_field_set;
static zend_object_handlers php_protocolbuffers_unknown_field_set_object_handlers;
static void php_protocolbuffers_unknown_field_set_free_storage(php_protocolbuffers_unknown_field_set *object TSRMLS_DC)
{
php_protocolbuffers_unknown_field_set *unknown_field_set;
unknown_field_set = (php_protocolbuffers_unknown_field_set*)((char *) object - XtOffsetOf(php_protocolbuffers_unknown_field_set, zo))；
zend_object_std_dtor(&unknown_field_set->zo TSRMLS_CC);
}
zend_object *php_protocolbuffers_unknown_field_set_new(zend_class_entry *ce TSRMLS_DC)
{
php_protocolbuffers_unknown_field_set *intern;
intern = ecalloc(1, sizeof(php_protocolbuffers_unknown_field_set) + zend_object_properties_size(ce));
zend_object_std_init(&intern->zo, ce);
object_properties_init(&intern->zo, ce);
intern->zo.handlers = &php_protocolbuffers_unknown_field_set_object_handlers;
intern->max = 0;
intern->offset = 0;
return &intern->zo;
}
void php_protocolbuffers_unknown_field_set_class(TSRMLS_D)
{
// 此处有省略
php_protocol_buffers_unknown_field_set_class_entry->create_object = php_protocolbuffers_unknown_field_set_new;
memcpy(&php_protocolbuffers_unknown_field_set_object_handlers, zend_get_std_object_handlers(), sizeof(zend_object_handlers));
php_protocolbuffers_unknown_field_set_object_handlers.offset = XtOffsetOf(php_protocolbuffers_unknown_field_set, zo);
php_protocolbuffers_unknown_field_set_object_handlers.free_obj = php_protocolbuffers_unknown_field_set_free_storage;
}
我们想自定义一个php_protocolbuffers_unknown_field_set的对象，在它的结构体里面除了zend_object，还有自定义的max和offset，务必把zend_object放在最后。
实际生成对象的地方基本就是标准写法，先分配内存，包括php_protocolbuffers_unknown_field_set结构体的内存和对象属性的内存；然后对zend_object的handlers赋值；最后再对自己自定义的变量初始化。
实际生成对象handler的地方也是标准写法，先分配内存，offset是必须设置的，可选的设置项有free_obj，dtor_obj，clone_obj。
想取到zend_object，需要(STRUCT_NAME )((char )OBJECT - XtOffsetOf(STRUCT_NAME, zo))

https://www.qdchaoyi.com/hacbv/67623.html

https://blog.csdn.net/caohao0591/article/details/82191001

https://www.sohu.com/a/120741342_505802

https://blog.csdn.net/u013474436/article/details/53485140

https://www.cnblogs.com/CoderK/articles/6943274.html

https://www.sohu.com/a/120741342_505802

https://blog.csdn.net/weixin_33816946/article/details/90590457

https://www.cnblogs.com/natian-ws/p/9105338.html

https://blog.csdn.net/nomius/article/details/94028034

https://www.cnblogs.com/sohuhome/p/9800977.html

https://www.laruence.com/2018/04/08/3170.html