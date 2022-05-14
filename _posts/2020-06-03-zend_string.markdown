---
title: zend_string char
layout: post
category: php
author: 夏泽民
---
zend_string 和char＊不是一个东西，因此转换的时候特别注意
写扩展的时候遇到了一个bug
明明写进去了
    add_assoc_long(MYFILE_G(my_func_set), func_name, timeElipsed);
    但是取出来一直是0
 zval* lastVal = zend_hash_find(Z_ARRVAL_P(MYFILE_G(my_func_set)),key);
 排查发现key的类型不对
 写进去是char＊ 取出来却是zend_string一定要注意
<!-- more -->

https://www.cnblogs.com/pingyeaa/p/9688248.html
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
{% raw %}

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
ZEND_API int ZEND_FASTCALL zend_binary_strcasecmp(const char *s1, size_t len1, const char *s2, size_t len2) /* {{{ */
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
感兴趣的同学可以到源码中查看。

四、参考文献
《PHP7底层设计与源码实现》
《PHP7内核剖析》
http://cs.potsdam.edu/Documentation/php/html/zend-api.add-index-zval.html
https://www.php.net/manual/ru/internals2.variables.tables.php
https://github.com/microsoft/msphpsql/issues/999
https://chrhust.wordpress.com/2015/07/03/php%E6%89%A9%E5%B1%95%E5%BC%80%E5%8F%91%E5%AD%A6%E4%B9%A0%E7%AC%94%E8%AE%B0%EF%BC%88%E4%BA%94%EF%BC%89%EF%BC%9Ahashtable-array%E5%8F%8A%E5%85%B6%E4%BB%96/
https://www.laruence.com/2020/03/23/5605.html
https://www.laruence.com/2020/02/25/3182.html
https://twosee.cn/2018/07/17/custom-zend-object-hack-way/
https://www.laruence.com/2018/04/08/3170.html
https://www.kancloud.cn/lifei6671/php-kernel/674637

https://wiki.jikexueyuan.com/project/extending-embedding-php/8.2.html
Zend把与HashTable有关的API分成了好几类以便于我们寻找，这些API的返回值大多都是常量SUCCESS或者FAILURE。

创建HashTable
下面在介绍函数原型的时候都使用了ht名称，但是我们在编写扩展的时候， 一定不要使用这个名称，因为一些PHP宏展开后会声明这个名称的变量， 进而引发命名冲突。

创建并初始化一个HashTable非常简单，只要使用zend_hash_init函数即可，它的定义如下：

int zend_hash_init(
    HashTable *ht,
    uint nSize,
    hash_func_t pHashFunction,
    dtor_func_t pDestructor,
    zend_bool persistent
);
*ht是指针，指向一个HashTable，我们既可以&一个已存在的HashTable变量， 也可以通过emalloc()、pemalloc()等函数来直接申请一块内存， 不过最常用的方法还是用ALLOC_HASHTABLE(ht)宏来让内核自动的替我们完成这项工作。 ALLOC_HASHTABLE(ht)所做的工作相当于ht = emalloc(sizeof(HashTable));

nSize代表着这个HashTable可以拥有的元素的最大数量(HashTable能够包含任意数量的元素， 这个值只是为了提前申请好内存，提高性能，省的不停的进行rehash操作)。 在我们添加新的元素时，这个值会根据情况决定是否自动增长，有趣的是， 这个值永远都是2的次方，如果你给它的值不是一个2的次方的形式， 那它将自动调整成大于它的最小的2的次方值。 它的计算方法就像这样：nSize = pow(2, ceil(log(nSize, 2)));

pHashFunction是早期的Zend Engine中的一个参数，为了兼容没有去掉它， 但它已经没有用处了，所以我们直接赋成NULL就可以了。在原来， 它其实是一个钩子，用来让用户自己hook一个散列函数，替换php默认的DJBX33A算法实现。

pDestructor也代表着一个回调函数，当我们删除或者修改HashTable中其中一个元素时候便会调用， 它的函数原型必须是这样的：void method_name(void pElement);这里的pElement是一个指针，指向HashTable中那么将要被删除或者修改的那个数据，而数据的类型往往也是个指针。

persistent是最后一个参数，它的含义非常简单。 如果它为true，那么这个HashTable将永远存在于内存中，而不会在RSHUTDOWN阶段自动被注销掉。 此时第一个参数ht所指向的地址必须是通过pemalloc()函数申请的。
举个例子，PHP内核在每个Request请求的头部都调用了这个函数来初始化symbol_table。

zend_hash_init(&EG(symbol_table), 50, NULL, ZVAL_PTR_DTOR, 0);

//#define ZVAL_PTR_DTOR (void (*)(void *)) zval_ptr_dtor_wrapper
如你所见，每个元素在从符号表里删除的时候(比如执行"<?php unset($foo);"操作)，都会触发ZVAL_PTR_DTOR宏代表的函数来对其进行与引用计数有关的操作。

因为50不是2的整数幂形式，所以它会在函数执行时被调成成64。

添加&&修改
我们有四个常用的函数来完成这项操作，它们的原型分别如下：

int zend_hash_add(
    HashTable *ht,      //待操作的ht
    char *arKey,            //索引，如"my_key"
    uint nKeyLen,       //字符串索引的长度，如6
    void **pData,       //要插入的数据，注意它是void **类型的。int *p,i=1;p=&i,pData=&p;。
    uint nDataSize,
    void *pDest         //如果操作成功，则pDest=*pData;
);

int zend_hash_update(
    HashTable *ht,
    char *arKey,
    uint nKeyLen,
    void *pData,
    uint nDataSize,
    void **pDest
);

int zend_hash_index_update(
    HashTable *ht,
    ulong h,
    void *pData,
    uint nDataSize,
    void **pDest
);

int zend_hash_next_index_insert(
    HashTable *ht,
    void *pData,
    uint nDataSize,
    void **pDest
);
前两个函数用户添加带字符串索引的数据到HashTable中，就像我们在PHP中使用的那样:$foo['bar'] = 'baz';用C来完成便是：

zend_hash_add(fooHashTbl, "bar", sizeof("bar"), &barZval, sizeof(zval*), NULL);
zend_hash_add()和zend_hash_update()唯一的区别就是如果这个key已经存在了，那么zend_hash_add()将返回FAILURE，而不会修改原有数据。

接下来的两个函数用于像HT中添加数字索引的数据，zend_hash_next_index_insert()函数则不需要索引值参数，而是自己直接计算出下一个数字索引值.

但是如果我们想获取下一个元素的数字索引值，也是有办法的，可以使用zend_hash_next_free_element()函数：

ulong nextid = zend_hash_next_free_element(ht);
zend_hash_index_update(ht, nextid, &data, sizeof(data), NULL);
所有这些函数中，如果pDest不为NULL，内核便会修改其值为被操作的那个元素的地址。在下面的代码中这个参数也有同样的功能。

查找
因为HashTable中有两种类型的索引值，所以需要两个函数来执行find操作。

int zend_hash_find(HashTable *ht, char *arKey, uint nKeyLength,void **pData);
int zend_hash_index_find(HashTable *ht, ulong h, void **pData);
第一种就是我们处理PHP语言中字符串索引数组时使用的，第二种是我们处理PHP语言中数字索引数组使用的。

Recall from Chapter 2 that when data is added to a HashTable, a new memory block is allocated for it and the data passed in is copied; when the data is extracted back out it is the pointer to that data which is returned. The following code fragment adds data1 to the HashTable, and then extracts it back out such that at the end of the routine, data2 contains the same contents as data1 even though the pointers refer to different memory addresses.

void hash_sample(HashTable *ht, sample_data *data1)
{
    sample_data *data2;
    ulong targetID = zend_hash_next_free_element(ht);
    if (zend_hash_index_update(ht, targetID,
            data1, sizeof(sample_data), NULL) == FAILURE) {
            /* Should never happen */
            return;
    }
    if(zend_hash_index_find(ht, targetID, (void **)&data2) == FAILURE) {
        /* Very unlikely since we just added this element */
        return;
    }
    /* data1 != data2, however *data1 == *data2 */
}
除了读取，我们还需要检测某个key是否存在：

int zend_hash_exists(HashTable *ht, char *arKey, uint nKeyLen);
int zend_hash_index_exists(HashTable *ht, ulong h);
这两个函数返回SUCCESS或者FAILURE，分别代表着是否存在：

if( zend_hash_exists(EG(active_symbol_table),"foo", sizeof("foo")) == SUCCESS )
{
    /* $foo is set */
}
else
{
    /* $foo does not exist */
}
提速!
ulong zend_get_hash_value(char *arKey, uint nKeyLen);
当我们需要对同一个字符串的key进行许多操作时候，比如先检测有没，然后插入，然后修改等等，这时我们便可以使用zend_get_hash_value函数来对我们的操作进行加速！这个函数的返回值可以和quick系列函数使用，达到加速的目的(就是不再重复计算这个字符串的散列值，而直接使用已准备好的)！

int zend_hash_quick_add(
    HashTable *ht,
    char *arKey,
    uint nKeyLen,
    ulong hashval,
    void *pData,
    uint nDataSize,
    void **pDest
);

int zend_hash_quick_update(
    HashTable *ht,
    char *arKey,
    uint nKeyLen,
    ulong hashval,
    void *pData,
    uint nDataSize,
    void **pDest
);

int zend_hash_quick_find(
    HashTable *ht,
    char *arKey,
    uint nKeyLen,
    ulong hashval,
    void **pData
);

int zend_hash_quick_exists(
    HashTable *ht,
    char *arKey,
    uint nKeyLen,
    ulong hashval
);
虽然很意外，但你还是要接受没有zend_hash_quick_del()这个函数。quick类函数会在下面这种场合中用到：

void php_sample_hash_copy(HashTable *hta, HashTable *htb,char *arKey, uint nKeyLen TSRMLS_DC)
{
    ulong hashval = zend_get_hash_value(arKey, nKeyLen);
    zval **copyval;

    if (zend_hash_quick_find(hta, arKey, nKeyLen,hashval, (void**)&copyval) == FAILURE)
    {
        //标明不存在这个索引
        return;
    }

    //这个zval已经被其它的Hashtable使用了，这里我们进行引用计数操作。
    (*copyval)->refcount__gc++;
    zend_hash_quick_update(htb, arKey, nKeyLen, hashval,copyval, sizeof(zval*), NULL);
}
复制与合并(Copy And Merge)
在PHP语言中，我们经常需要进行数组间的Copy与Merge操作，所以php语言中的数组在C语言中的实现HashTable也肯定会经常碰到这种情况。为了简化这一类操作，内核中早已准备好了相应的API供我们使用。

void zend_hash_copy(
    HashTable *target,
    HashTable *source,
    copy_ctor_func_t pCopyConstructor,
    void *tmp,
    uint size
);
*source中的所有元素都会通过pCopyConstructor函数Copy到*target中去，我们还是以PHP语言中的数组举例，pCopyConstructor这个hook使得我们可以在copy变量的时候对他们的ref_count进行加一操作。target中原有的与source中索引位置的数据会被替换掉，而其它的元素则会被保留，原封不动。

tmp参数是为了兼容PHP4.0.3以前版本的，现在赋值为NULL即可。
size参数代表每个元素的大小，对于PHP语言中的数组来说，这里的便是sizeof(zval*)了。
void zend_hash_merge(
    HashTable *target,
    HashTable *source,
    copy_ctor_func_t pCopyConstructor,
    void *tmp,
    uint size,
    int overwrite
);
zend_hash_merge()与zend_hash_copy唯一的不同便是多了个int类型的overwrite参数，当其值非0的时候，两个函数的工作是完全一样的；如果overwrite参数为0，则zend_hash_merge函数就不会对target中已有索引的值进行替换了。

typedef zend_bool (*merge_checker_func_t)(HashTable *target_ht,void *source_data, zend_hash_key *hash_key, void *pParam);
void zend_hash_merge_ex(
    HashTable *target,
    HashTable *source,
    copy_ctor_func_t pCopyConstructor, 
    uint size,
    merge_checker_func_t pMergeSource,
    void *pParam
);
这个函数又繁琐了些，与zend_hash_copy相比，其多了两个参数，多出来的pMergeSoure回调函数允许我们选择性的进行merge，而不是全都merge。The final form of this group of functions allows for selective copying using a merge checker function. The following example shows zend_hash_merge_ex() in use to copy only the associatively indexed members of the source HashTable (which happens to be a userspace variable array):

zend_bool associative_only(HashTable *ht, void *pData,zend_hash_key *hash_key, void *pParam)
{
    //如果是字符串索引
    return (hash_key->arKey && hash_key->nKeyLength);
}

void merge_associative(HashTable *target, HashTable *source)
{
    zend_hash_merge_ex(target, source, zval_add_ref,sizeof(zval*), associative_only, NULL);
}
遍历
在PHP语言中，我们有很多方法来遍历一个数组，对于数组的本质HashTable，我们也有很多办法来对其进行遍历操作。首先最简单的一种办法便是使用一种与PHP语言中forech语句功能类似的函数——zend_hash_apply，它接收一个回调函数，并将HashTable的每一个元素都传递给它。

typedef int (*apply_func_t)(void *pDest TSRMLS_DC);
void zend_hash_apply(HashTable *ht,apply_func_t apply_func TSRMLS_DC);
下面是另外一种遍历函数：

typedef int (*apply_func_arg_t)(void *pDest,void *argument TSRMLS_DC);
void zend_hash_apply_with_argument(HashTable *ht,apply_func_arg_t apply_func, void *data TSRMLS_DC);
通过上面的函数可以在执行遍历时向回调函数传递任意数量的值，这在一些diy操作中非常有用。

上述函数对传给它们的回调函数的返回值有一个共同的约定，详细介绍下下表：

表格 8.1. 回调函数的返回值
Constant                        Meaning

ZEND_HASH_APPLY_KEEP        结束当前请求，进入下一个循环。与PHP语言forech语句中的一次循环执行完毕或者遇到continue关键字的作用一样。
ZEND_HASH_APPLY_STOP        跳出，与PHP语言forech语句中的break关键字的作用一样。
ZEND_HASH_APPLY_REMOVE      删除当前的元素，然后继续处理下一个。相当于在PHP语言中：unset($foo[$key]);continu;
我们来一下PHP语言中的forech循环：

<?php
foreach($arr as $val) {
    echo "The value is: $val\n";
}
?>
那我们的回调函数在C语言中应该这样写：

int php_sample_print_zval(zval **val TSRMLS_DC)
{
    //重新copy一个zval，防止破坏原数据
    zval tmpcopy = **val;
    zval_copy_ctor(&tmpcopy);

    //转换为字符串
    INIT_PZVAL(&tmpcopy);
    convert_to_string(&tmpcopy);

    //开始输出
    php_printf("The value is: ");
    PHPWRITE(Z_STRVAL(tmpcopy), Z_STRLEN(tmpcopy));
    php_printf("\n");

    //毁尸灭迹
    zval_dtor(&tmpcopy);

    //返回，继续遍历下一个～
    return ZEND_HASH_APPLY_KEEP;
}
遍历我们的HashTable：

//生成一个名为arrht、元素为zval*类型的HashTable
zend_hash_apply(arrht, php_sample_print_zval TSRMLS_CC);
再次提醒，保存在HashTable中的元素并不是真正的最终变量，而是指向它的一个指针。我们的上面的遍历函数接收的是一个zval**类型的参数。

typedef int (*apply_func_args_t)(void *pDest,int num_args, va_list args, zend_hash_key *hash_key);
void zend_hash_apply_with_arguments(HashTable *ht,apply_func_args_t apply_func, int numargs, ...);
为了能在遍历时同时接收索引的值，我们必须使用第三种形式的zend_hash_apply！就像PHP语言中这样的功能：

<?php
foreach($arr as $key => $val)
{
    echo "The value of $key is: $val\n";
}
?>
为了配合zend_hash_apply_with_arguments()函数，我们需要对我们的遍历执行函数做一下小小的改动，使其接受索引作为一个参数：

int php_sample_print_zval_and_key(zval **val,int num_args,va_list args,zend_hash_key *hash_key)
{
    //重新copy一个zval，防止破坏原数据
    zval tmpcopy = **val;
    /* tsrm_ls is needed by output functions */
    TSRMLS_FETCH();
    zval_copy_ctor(&tmpcopy);
    INIT_PZVAL(&tmpcopy);

    //转换为字符串
    convert_to_string(&tmpcopy);

    //执行输出
    php_printf("The value of ");
    if (hash_key->nKeyLength)
    {
        //如果是字符串类型的key
        PHPWRITE(hash_key->arKey, hash_key->nKeyLength);
    }
    else
    {
        //如果是数字类型的key
        php_printf("%ld", hash_key->h);
    }

    php_printf(" is: ");
    PHPWRITE(Z_STRVAL(tmpcopy), Z_STRLEN(tmpcopy));
    php_printf("\n");

    //毁尸灭迹
    zval_dtor(&tmpcopy);
    /* continue; */
    return ZEND_HASH_APPLY_KEEP;
}
执行遍历：

zend_hash_apply_with_arguments(arrht,php_sample_print_zval_and_key, 0);
这个函数通过C语言中的可变参数特性来接收参数。This particular example required no arguments to be passed; for information on extracting variable argument lists from va_list args, see the POSIX documentation pages for va_start(), va_arg(), and va_end().

当我们检查这个hash_key是字符串类型还是数字类型时，是通过nKeyLength属性来检测的,而不是arKey属性。这是因为内核有时候会留在arKey属性里些脏数据，但nKeyLength属性是安全的，可以安全的使用。甚至对于空字符串索引，它也照样能处理。比如：$foo[''] ="Bar";索引的值是NULL字符，但它的长度却是包括最后这个NULL字符的，所以为1。

向前遍历HashTable
有时我们希望不用回调函数也能遍历一个数组的数据，为了实现这个功能，内核特意的为每个HashTable加了个属性：The internal pointer（内部指针）。

我们还是以PHP语言中的数组举例，有以下函数来处理它所对应的那个HashTable的内部指针：reset(), key(), current(), next(), prev(), each(), and end()。

<?php
    $arr = array('a'=>1, 'b'=>2, 'c'=>3);
    reset($arr);
    while (list($key, $val) = each($arr)) {
        /* Do something with $key and $val */
    }
    reset($arr);
    $firstkey = key($arr);
    $firstval = current($arr);
    $bval = next($arr);
    $cval = next($arr);
?>
ZEND内核中有一组操作HashTable的功能与以上函数功能类似的函数：

/* reset() */
void zend_hash_internal_pointer_reset(HashTable *ht);

/* key() */
int zend_hash_get_current_key(HashTable *ht,char **strIdx, unit *strIdxLen,ulong *numIdx, zend_bool duplicate);

/* current() */
int zend_hash_get_current_data(HashTable *ht, void **pData);

/* next()/each() */
int zend_hash_move_forward(HashTable *ht);

/* prev() */
int zend_hash_move_backwards(HashTable *ht);

/* end() */
void zend_hash_internal_pointer_end(HashTable *ht);

/* 其他的...... */
int zend_hash_get_current_key_type(HashTable *ht);
int zend_hash_has_more_elements(HashTable *ht);
PHP语言中的next()、prev()、end()函数在移动完指针之后，都通过调用zend_hash_get_current_data()函数来获取当前所指的元素并返回。而each()虽然和next()很像，却是使用zend_hash_get_current_key()函数的返回值来作为它的返回值。

现在我们用另外一种方法来实现上面的forech：

void php_sample_print_var_hash(HashTable *arrht)
{

    for(
        zend_hash_internal_pointer_reset(arrht);
        zend_hash_has_more_elements(arrht) == SUCCESS;
        zend_hash_move_forward(arrht))
    {
        char *key;
        uint keylen;
        ulong idx;
        int type;
        zval **ppzval, tmpcopy;

        type = zend_hash_get_current_key_ex(arrht, &key, &keylen,&idx, 0, NULL);
        if (zend_hash_get_current_data(arrht, (void**)&ppzval) == FAILURE)
        {
            /* Should never actually fail
             * since the key is known to exist. */
            continue;
        }

        //重新copy一个zval，防止破坏原数据
        tmpcopy = **ppzval;
        zval_copy_ctor(&tmpcopy);
        INIT_PZVAL(&tmpcopy);

        convert_to_string(&tmpcopy);

        /* Output */
        php_printf("The value of ");
        if (type == HASH_KEY_IS_STRING)
        {
            /* String Key / Associative */
            PHPWRITE(key, keylen);
        } else {
            /* Numeric Key */
            php_printf("%ld", idx);
        }
        php_printf(" is: ");
        PHPWRITE(Z_STRVAL(tmpcopy), Z_STRLEN(tmpcopy));
        php_printf("\n");
        /* Toss out old copy */
        zval_dtor(&tmpcopy);
    }
}
上面的代码你应该都能看懂了，唯一还没接触到的可能是zend_hash_get_current_key()函数的返回值，它的返回值见表8.2。

Constant                            Meaning

HASH_KEY_IS_STRING              当前元素的索引是字符串类型的。therefore, a pointer to the element's key name will be populated into strIdx, and its length will be populated into stdIdxLen. If the duplicate flag is set to a nonzero value, the key will be estrndup()'d before being populated into strIdx. The calling application is expected to free this duplicated string.

HASH_KEY_IS_LONG                当前元素的索引是数字型的。
HASH_KEY_NON_EXISTANT           HashTable中的内部指针已经移动到尾部，不指向任何元素。
Preserving the Internal Pointer
在我们遍历一个HashTable时，一般是很难陷入死循环的。When iterating through a HashTable, particularly one containing userspace variables, it's not uncommon to encounter circular references, or at least self-overlapping loops. If one iteration context starts looping through a HashTable and the internal pointer reachesfor examplethe halfway mark, a subordinate iterator starts looping through the same HashTable and would obliterate the current internal pointer position, leaving the HashTable at the end when it arrived back at the first loop.

The way this is resolvedboth within the zend_hash_apply implementation and within custom move forward usesis to supply an external pointer in the form of a HashPosition variable.

Each of the zendhash() functions listed previously has a zendhash_ex() counterpart that accepts one additional parameter in the form of a pointer to a HashPostion data type. Because the HashPosition variable is seldom used outside of a short-lived iteration loop, it's sufficient to declare it as an immediate variable. You can then dereference it on usage such as in the following variation on the php_sample_print_var_hash() function you saw earlier:

void php_sample_print_var_hash(HashTable *arrht)
{
    HashPosition pos;
    for(zend_hash_internal_pointer_reset_ex(arrht, &pos);
    zend_hash_has_more_elements_ex(arrht, &pos) == SUCCESS;
    zend_hash_move_forward_ex(arrht, &pos)) {
        char *key;
        uint keylen;
        ulong idx;
        int type;

        zval **ppzval, tmpcopy;

        type = zend_hash_get_current_key_ex(arrht,
                                &key, &keylen,
                                &idx, 0, &pos);
        if (zend_hash_get_current_data_ex(arrht,
                    (void**)&ppzval, &pos) == FAILURE) {
            /* Should never actually fail
             * since the key is known to exist. */
            continue;
        }
        /* Duplicate the zval so that
         * the original's contents are not destroyed */
        tmpcopy = **ppzval;
        zval_copy_ctor(&tmpcopy);
        /* Reset refcount & Convert */
        INIT_PZVAL(&tmpcopy);
        convert_to_string(&tmpcopy);
        /* Output */
        php_printf("The value of ");
        if (type == HASH_KEY_IS_STRING) {
            /* String Key / Associative */
            PHPWRITE(key, keylen);
        } else {
            /* Numeric Key */
            php_printf("%ld", idx);
        }
        php_printf(" is: ");
        PHPWRITE(Z_STRVAL(tmpcopy), Z_STRLEN(tmpcopy));
        php_printf("\n");
        /* Toss out old copy */
        zval_dtor(&tmpcopy);
    }
}
With these very slight additions, the HashTable's true internal pointer is preserved in whatever state it was initially in on entering the function. When it comes to working with internal pointers of userspace variable HashTables (that is, arrays), this extra step will very likely make the difference between whether the scripter's code works as expected.

删除
内核中一共预置了四个删除HashTable元素的函数，头两个是用户删除某个确定索引的数据：

int zendhashdel(HashTable ht, char arKey, uint nKeyLen); int zendhashindex_del(HashTable *ht, ulong h);

它们两个分别用来删除字符串索引和数字索引的数据，操作完成后都返回SUCCESS或者FAILURE表示成功or失败。 回顾一下最上面的叙述，当一个元素被删除时，会激活HashTable的destructor回调函数。

void zend_hash_clean(HashTable *ht);
void zend_hash_destroy(HashTable *ht);
前者用于将HashTable中的元素全部删除，而后者是将这个HashTable自身也毁灭掉。 现在让我们来完整的回顾一下HashTable的创建、添加、删除操作。

int sample_strvec_handler(int argc, char **argv TSRMLS_DC)
{
    HashTable *ht;

    //分配内存
    ALLOC_HASHTABLE(ht);

    //初始化
    if (zend_hash_init(ht, argc, NULL,ZVAL_PTR_DTOR, 0) == FAILURE) {
        FREE_HASHTABLE(ht);
        return FAILURE;
    }

    //填充数据
    while (argc) {
        zval *value;
        MAKE_STD_ZVAL(value);
        ZVAL_STRING(value, argv[argc], 1);
        argv++;
        if (zend_hash_next_index_insert(ht, (void**)&value,
                            sizeof(zval*)) == FAILURE) {
            /* Silently skip failed additions */
            zval_ptr_dtor(&value);
        }
    }

    //完成工作
    process_hashtable(ht);

    //毁尸灭迹
    zend_hash_destroy(ht);

    //释放ht 为什么不在destroy里free呢，求解释！
    FREE_HASHTABLE(ht);
    return SUCCESS;
}
排序、比较and Going to the Extreme(s)
针对HashTable操作的Zend Api中有很多都需要回调函数。首先让我们来处理一下对HashTable中元素大小比较的问题：

typedef int (*compare_func_t)(void *a, void *b TSRMLS_DC);
这很像PHP语言中usort函数需要的参数，它将比较两个值a与b，如果a>b,则返回1，相等则返回0，否则返回-1。下面是zend_hash_minmax函数的声明，它就需要我们上面声明的那个类型的函数作为回调函数： int zend_hash_minmax(HashTable *ht, compare_func_t compar,int flag, void **pData TSRMLS_DC); 这个函数的功能我们从它的名称中便能肯定，它用来比较HashTable中的元素大小。如果flag==0则返回最小值，否则返回最大值！

下面让我们来利用这个函数来对用户端定义的所有函数根据函数名找到最大值与最小值(大小写不敏感～)。

//先定义一个比较函数，作为zend_hash_minmax的回调函数。
int fname_compare(zend_function *a, zend_function *b TSRMLS_DC)
{
    return strcasecmp(a->common.function_name, b->common.function_name);
}

void php_sample_funcname_sort(TSRMLS_D)
{
    zend_function *fe;
    if (zend_hash_minmax(EG(function_table), fname_compare,0, (void **)&fe) == SUCCESS)
    {
        php_printf("Min function: %s\n", fe->common.function_name);
    }
    if (zend_hash_minmax(EG(function_table), fname_compare,1, (void **)&fe) == SUCCESS)
    {
        php_printf("Max function: %s\n", fe->common.function_name);
    }
}
zend_hash_compare()也许要回调函数，它的功能是将HashTable看作一个整体与另一个HashTable做比较，如果前者大于后者返回1，相等返回0，否则返回-1。

int zendhashcompare(HashTable hta, HashTable htb,comparefunct compar, zendbool ordered TSRMLSDC);

默认情况下它往往是先判断各个HashTable元素的个数，个数多的最大！ 如果两者的元素一样多，然后就比较它们各自的第一个元素，If the ordered flag is set, it compares keys/indices with the first element of htb string keys are compared first on length, and then on binary sequence using memcmp(). If the keys are equal, the value of the element is compared with the first element of htb using the comparison callback function.

If the ordered flag is not set, the data portion of the first element of hta is compared against the element with a matching key/index in htb using the comparison callback function. If no matching element can be found for htb, then hta is considered greater than htb and 1 is returned.

If at the end of a given loop, hta and htb are still considered equal, comparison continues with the next element of hta until a difference is found or all elements have been exhausted, in which case 0 is returned.

另外一个重要的需要回调函数的API便是排序函数，它需要的回调函数形式是这样的：

typedef void (*sort_func_t)(void **Buckets, size_t numBuckets,size_t sizBucket, compare_func_t comp TSRMLS_DC);
This callback will be triggered once, and receive a vector of all the Buckets (elements) in the HashTable as a series of pointers. These Buckets may be swapped around within the vector according to the sort function's own logic with or without the use of the comparison callback. In practice, sizBucket will always be sizeof(Bucket*).

Unless you plan on implementing your own alternative bubblesort method, you won't need to implement a sort function yourself. A predefined sort methodzend_qsortalready exists for use as a callback to zend_hash_sort() leaving you to implement the comparison function only.

int zend_hash_sort(HashTable *ht, sort_func_t sort_func,compare_func_t compare_func, int renumber TSRMLS_DC);
最后一个参数如果为TRUE，则会抛弃HashTable中原有的索引-键关系，将对排列好的新值赋予新的数字键值。PHP语言中的sort函数实现如下：

zend_hash_sort(target_hash, zend_qsort,array_data_compare, 1 TSRMLS_CC);
array_data_compare是一个返回compare_func_t类型数据的函数，它将按照HashTable中zval*值的大小进行排序。

https://www.cnblogs.com/niniwzw/archive/2010/03/09/1681932.html
1.常用的通用功能已经封装好了，在如zen_API.h 头文件中，不用费力查看内部细节，浪费时间。（参考：Extending and Embedding PHP 的附录A）
2.在terminal中运行测试程序，可以看到扩展的内部错误输出，这一点对于解决内存泄漏问题尤其重要。（编译一个debug 的 lib）
3.开发过程中修改Makefile中的“CFLAGS = -g -O2”，去掉优化选项，增加-Wall和-pedantic，便于调试和显示编译警告；
4.某zval*，但其strval非拷贝的，不可用zval_ptr_dtor(zval**)，要用efree(void*)。
5.terminal中的$_SERVER['PWD']有值，但是无法通过zend_getenv()取得，原因应该是该值无意义或不可靠。
6.调用“导出函数”，可利用INTERNAL_FUNCTION_PARAM_PASSTHRU传参；声明的非导出函数可通过INTERNAL_FUNCTION_PARAM使用“导出函数”的参数。
7.注意：RETURN_TYPE用在选择分之和循环等处时，最好置于花括号中，
或者不用分号，因为：#define RETURN_BOOL(b) { RETVAL_BOOL(b); return; }。
8.如果函数的参数是引用的，且非标量，要先析构，以防内存泄露。
9.抛出异常前最好判断EG(exception)中是否已经存在异常，否则会造成内存泄露。
10.当Web服务器API是ISAPI (IIS)的时候，zend_getenv函数是不起作用的。
11.向zend_stack_push()传入数据指针，实际存储(copy)的是该指针指向的数据，换句话说，传入的应该是要存储的数据的指针。
ZEND_API int zend_stack_push(zend_stack *stack, void *element, int size);
ZEND_API int zend_stack_top(zend_stack *stack, void **element);
其中，size == sizeof(*element);
类似地，zend_hash也是如此，比较zend_hash_update和zend_hash_find。
12.使用add_assoc_zval(HashTable*, const char*, zval*)存储的是zval*，而非zval，因此，
存储用户传入的参数时候，要先拷贝一份新的zval，否则会发生不可预料的事情。
13.zval_dtor(zval*)释放变量及其内部的引用内存，zval_ptr_dtor(zval**)先检查refcount
再决定是否调用zval_dtor(zval*)，zval_copy_dtor(zval*)仅执行深层的拷贝，即只拷贝
起内部引用的内存，而不拷贝zval；

14.如使用VC编译win的动态链接库，而且代码中调用了zend函数，如zend_getenv，在zend.h中定义为：



extern "C" {
extern ZEND_API char *(*zend_getenv)(char *name, size_t name_len TSRMLS_DC);
}
需要引入该函数，如要使用ZEND_API，需要事先取消LIBZEND_EXPORTS（包括VC“设置”中的预处理定义），或者使用ZEND_DLIMPORT，
ZEND_DLIMPORT char *(*zend_getenv)(char *name, size_t name_len TSRMLS_DC);
下面取自：zend_config.w32.h

复制代码代码如下:

#ifdef LIBZEND_EXPORTS
# define ZEND_API __declspec(dllexport)
#else
# define ZEND_API __declspec(dllimport)
#endif
#define ZEND_DLEXPORT __declspec(dllexport)
#define ZEND_DLIMPORT __declspec(dllimport)
{% endraw %}
https://zhuanlan.zhihu.com/p/84221387
https://blog.csdn.net/caohao0591/article/details/82191001
https://www.php.net/manual/en/internals2.variables.tables.php
https://learnku.com/articles/9173/2-analysis-of-zval
https://nikic.github.io/2015/05/05/Internal-value-representation-in-PHP-7-part-1.html
https://segmentfault.com/a/1190000007575322
https://www.bo56.com/php7%E6%89%A9%E5%B1%95%E5%BC%80%E5%8F%91%E4%B9%8B%E9%85%8D%E7%BD%AE%E9%A1%B9/

