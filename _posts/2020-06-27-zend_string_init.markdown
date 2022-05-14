---
title: zend_string_init
layout: post
category: php
author: 夏泽民
---
char * 转zend_string
zend_string_init

zend_string转char＊
ZSTR_VAL
<!-- more -->
一、字符串的结构
Copy
struct _zend_string {
	zend_refcounted_h gc;       /* 字符串类别及引用计数 */
	zend_ulong        h;        /* 字符串的哈希值 */
	size_t            len;      /* 字符串的长度 */
	char              val[1];   /* 柔性数组，字符串存储位置 */
};
zend_refcounted_h对应的结构体：

Copy
typedef struct _zend_refcounted_h {
	uint32_t         refcount;			/* 引用计数 */
	union {
		struct {
			ZEND_ENDIAN_LOHI_3(
				zend_uchar    type,     
				zend_uchar    flags,    /* 字符串的类型 */
				uint16_t      gc_info   /* 垃圾回收信息 */
			)
		} v;
		uint32_t type_info;
	} u;
} zend_refcounted_h;
image

下面我们来了解一下具体每个成员的作用：

gc：就是_zend_refcounted_h结构体，主要作用是引用计数以及标记变量的类别。
h：字符串的哈希值，在字符串被用来当数组的key时才初始化，这样如果同一个字符串被多次用来做key，就不会重复计算了。
val：这里的char[1]并不意味着只存储1位，char[1]被称为柔性数组，下面来了解一下PHP在字符串内存分配时做了什么。
Copy
static zend_always_inline zend_string *zend_string_alloc(size_t len, int persistent)
{
	zend_string *ret = (zend_string *)pemalloc(ZEND_MM_ALIGNED_SIZE(_ZSTR_STRUCT_SIZE(len)), persistent);
    ......
}
宏替换后：

Copy
static zend_always_inline zend_string *zend_string_alloc(size_t len, int persistent)
{
	zend_string *ret = (zend_string *)pemalloc(ZEND_MM_ALIGNED_SIZE(XtOffsetOf(zend_string, val) + len + 1), persistent);
    ......
}
示例中的代码XtOffsetOf(zend_string, val)表示计算出zend_string结构体的大小，而len就是要分配字符串的长度，最后的+1是留给结束字符\0的。也就是说，分配内存时不仅仅分配结构体大小的内存，还要顾及到长度不可控的val，这样不仅柔性的分配了内存，还使它与其他成员存储在同一块连续的空间中，在分配、释放内存时可以把struct统一处理。

len：字符串的长度，避免重复计算浪费时间，典型的空间换时间做法。
二、字符串的二进制安全
学习过C语言的应该知道，字符串中除了最后一个字符外不允许含有\0，否则会被认为是字符串的结束字符，这就导致了C语言的字符串有很多的限制，比如不存储图片、文件等二进制数据。但是PHP就没有这样的限制，它的字符串可以存储二进制数据，并不会出现任何报错，而PHP的这种能力就叫做字符串的二进制安全。

C语言代码如下：

Copy
main() {
    char a[] = "aaa\0b";    /* 含有\0的字符串 */
    printf("%d\n", strlen(a));  /* 长度为3，\0后的b被忽略 */
}
PHP代码：

Copy
<?php
    $a = "aaa\0b";
    echo strlen($a);    //输出5
?>
但是PHP不是C语言写的吗？为什么PHP不会报错？我们再来回顾一下zend_string结构体，还记得成员变量len吗？它是实现二进制安全的关键，我们不需要像C一样通过\0来判定字符串是否被读取完成，而是通过长度len来判断，这样就保证了字符串的二进制安全。

三、zend_string API
在了解了zend_string结构之后，我们来了解一下用来操作zend_string的函数集合。

函数	作用
zend_interned_strings_init	初始化内部字符串存储哈希表，并把PHP的关键字等字符串信息写进去
zend_new_interned_string	把一个zend_string写入CG(interned_strings)哈希表中
zend_interned_strings_snapshot	将CG(interned_strings)哈希表中的字符串标记为永久字符串，这里标记的只有PHP关键字、内部函数名、内部方法名等
zend_interned_strings_restore	销毁CG(interned_strings)哈希表中类型为非永久字符串的值，在php_request_shutdown阶段释放
zend_interned_strings_dtor	销毁整个CG(interned_strings)哈希表，在php_module_shutdown阶段释放
zend_string_hash_val	得到字符串的哈希值，没有则实时计算
zend_string_forget_hash_val	将字符串的哈希值置为0
zend_string_refcount	读取字符串的引用计数
zend_string_addref	引用计数+1
zend_string_delref	引用计数-1
zend_string_alloc	分配内存及初始化字符串的值
zend_string_init	初始化字符串并在最后追加\0
zend_string_cop	使用引用计数方式复制字符串
zend_string_dup	直接复制一个字符串
zend_string_extend	扩容到len，保留原来的值
zend_string_truncate	截断到len，保留开头到len的值
zend_string_free	释放字符串内存
zend_string_release	GC引用递减，直到为0时释放内存
zend_string_equals	普通判等
zend_string_equals_ci	基于二进制安全，两个zend_string类型字符串判等
zend_string_equals_literal_ci	基于二进制安全，zend_string类型和char*字符串判等
zend_inline_hash_func	计算字符串的哈希值
zend_intern_known_strings	往zend_intern_known_strings全局数组写入str
下面挑几个函数来介绍一下。

3.1、zend_string_init函数#
zend_string_init函数主要负责把一个普通的字符串转化为zend_string结构体。

Copy
static zend_always_inline zend_string *zend_string_init(const char *str, size_t len, int persistent)
{
	zend_string *ret = zend_string_alloc(len, persistent);

	memcpy(ZSTR_VAL(ret), str, len);
	ZSTR_VAL(ret)[len] = '\0';
	return ret;
}
申请一块连续的内存，这个在上文中已经提到，申请的内存大小是zend_string结构体大小+字符串长度+1。
指针偏移到val位置，开始字符串拷贝。
在zend_string.val结尾追加\0。
3.2、zend_string_extend函数#
该函数主要用于对字符串的扩容，注意这里扩容不会改变原来保存的值，只是把长度扩大到len。

Copy
static zend_always_inline zend_string *zend_string_extend(zend_string *s, size_t len, int persistent)
{
	zend_string *ret;

	ZEND_ASSERT(len >= ZSTR_LEN(s));
	if (!ZSTR_IS_INTERNED(s)) {
		if (EXPECTED(GC_REFCOUNT(s) == 1)) {
			ret = (zend_string *)perealloc(s, ZEND_MM_ALIGNED_SIZE(_ZSTR_STRUCT_SIZE(len)), persistent);
			ZSTR_LEN(ret) = len;
			zend_string_forget_hash_val(ret);
			return ret;
		} else {
			GC_REFCOUNT(s)--;
		}
	}
	ret = zend_string_alloc(len, persistent);
	memcpy(ZSTR_VAL(ret), ZSTR_VAL(s), ZSTR_LEN(s) + 1);
	return ret;
}
如果不是内部字符串并且引用计数为1时，直接调用perealloc分配内存。
如果字符串的引用计数大于1或者是内部字符串时，就不能在原来的基础上扩容了，先通过zend_string_alloc申请一块新内存，让后将旧内容拷贝到新内存中。
3.3、zend_string_equals_ci函数#
主要基于二进制安全对两个字符串进行判等，我们来看下PHP是怎么比较两个字符串的。

Copy
#define zend_string_equals_ci(s1, s2) \
	(ZSTR_LEN(s1) == ZSTR_LEN(s2) && !zend_binary_strcasecmp(ZSTR_VAL(s1), ZSTR_LEN(s1), ZSTR_VAL(s2), ZSTR_LEN(s2)))

先比较两个字符串的长度是否相等，注意这里是通过zend_string中的len来比较的。
zend_binary_strcasecmp函数在长度比较完成后，进行逐个字符进行比较。先遍历整个字符串数组，取出每个字符，转换为ASC码进行判等，如果不等则返回差值。循环完了还没发现差异的话就返回两者的长度差，如果长度相等就返回0。感觉这里做的有点多余，参数传进来之前就已经做了长度判等了。
Copy
ZEND_API int ZEND_FASTCALL zend_binary_strcasecmp(const char *s1, size_t len1, const char *s2, size_t len2) /* \{\{\{ */
{
	size_t len;
	int c1, c2;

	if (s1 == s2) {
		return 0;
	}

	len = MIN(len1, len2);
	while (len--) {
		c1 = zend_tolower_ascii(*(unsigned char *)s1++);
		c2 = zend_tolower_ascii(*(unsigned char *)s2++);
		if (c1 != c2) {
			return c1 - c2;
		}
	}

	return (int)(len1 - len2);
}

https://www.cnblogs.com/pingyeaa/p/9688248.html

https://www.cnblogs.com/ITCM/articles/6943600.html

https://www.cnblogs.com/pingyeaa/p/9688248.html

https://blog.csdn.net/yuhezheg/article/details/105028172

https://www.cnblogs.com/yjf512/p/6108444.html
