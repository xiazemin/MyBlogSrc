---
title: ZEND_HASH_FOREACH_KEY_VAL
layout: post
category: web
author: 夏泽民
---
1、zend_hash_num_elements 获取数组元素个数。宏定义如下：
1 #define zend_hash_num_elements(ht) \
2     (ht)->nNumOfElements
2、ZEND_HASH_FOREACH_KEY_VAL 遍历数组键值。使用方法：
1 ZEND_HASH_FOREACH_KEY_VAL(Z_ARRVAL_P(array), num_key, string_key, entry) {
2             // code
3 } ZEND_HASH_FOREACH_END();
　ZEND_HASH_FOREACH_KEY_VAL是个宏函数：
1 #define ZEND_HASH_FOREACH_KEY_VAL(ht, _h, _key, _val) \
2     ZEND_HASH_FOREACH(ht, 0); \
3     _h = _p->h; \
4     _key = _p->key; \
5     _val = _z;
　  继续展开 ZEND_HASH_FOREACH：

复制代码
1 #define ZEND_HASH_FOREACH(_ht, indirect) do { \
2         Bucket *_p = (_ht)->arData; \
3         Bucket *_end = _p + (_ht)->nNumUsed; \
4         for (; _p != _end; _p++) { \
5             zval *_z = &_p->val; \
6             if (indirect && Z_TYPE_P(_z) == IS_INDIRECT) { \
7                 _z = Z_INDIRECT_P(_z); \
8             } \
9             if (UNEXPECTED(Z_TYPE_P(_z) == IS_UNDEF)) continue;
复制代码
ZEND_HASH_FOREACH_END展开：
1 #define ZEND_HASH_FOREACH_END() \
2         } \
3     } while (0)
  

ZEND_HASH_FOREACH_KEY_VAL(Z_ARRVAL_P(array), num_key, string_key, entry) {
             // code
} ZEND_HASH_FOREACH_END();
完全展开：
复制代码
 1 do { 
 2     Bucket *_p = (_ht)->arData;  // Z_ARRVAL_P(array) ---> ht ---> _ht
 3     Bucket *_end = _p + (_ht)->nNumUsed;  // 起始地址+偏移地址
 4     for (; _p != _end; _p++) { 
 5         zval *_z = &_p->val; 
 6         if (indirect && Z_TYPE_P(_z) == IS_INDIRECT) { 
 7             _z = Z_INDIRECT_P(_z); 
 8         } 
 9         if (UNEXPECTED(Z_TYPE_P(_z) == IS_UNDEF)) continue;
10         _h = _p->h;  // zend_ulong num_key ---> _h
11         _key = _p->key; // zend_string *string_key ---> _key
12         _val = _z; // zval *entry ---> _val
13         {
14            //code
15         } 
16     } 
17 } while (0)
复制代码

<!-- more -->
https://www.cnblogs.com/natian-ws/p/9105338.html
      数组
本节我们讲一下php的数组，在php中，数组使用HashTable实现的。本节中我们先详细的介绍一下HashTable，然后再讲讲如何使用HastTable

1.1     变长结构体
所谓的变长结构体，其实是我们C语言结构体的一种特殊用法，并没有什么新奇之处。我们先来看一下变长结构体的一种通用定义方法。

typedef struct bucket {

    int n;

    char key[30];

    char value[1];

} Bucket;

我们定义了一个结构体Bucket，我们希望用这个结构体存放学生的个人简介。其中key用来存在学生的姓名，value用来存放学生的简介。大家可能很好奇，我们的value声明了长度为1. 1个char能存多少信息呀？

         其实，对于变长结构体，我们在使用的使用不能直接定义变量，例如：Bucket bucket; 您要是这样使用，value肯定存储不了多少信息。对于变长结构体，我们在使用的时候需要先声明一个变长结构体的指针，然后通过malloc函数分配函数空间，我们需要用到的空间长度是多少，我们就可以malloc多少。通用的使用方法如下：

Bucket* pBucket;

pBucket = malloc(sizeof(Bucket)+ n *sizeof(char));

其中n就是你要使用value的长度。如果这样使用的话，value指向的字符串不久变长了吗！

 

1.2     Hashtable简介
我们先看一下HashTable的定义

struct _hashtable;

 

typedef struct bucket {

    ulong h;//当元素是数字索引时使用

    uint nKeyLength;//当使用字符串索引时，这个变量表示索引的长度，索引（字符串）保存在最后一个元素aKey

    void *pData;//用来指向保存的数据，如果保存的数据是指针的话，pDataPtr就指向这个数据，pData指向pDataPtr

    void *pDataPtr;

    struct bucket *pListNext;//上一个元素

    struct bucket *pListLast;//下一个元素

    struct bucket *pNext;//指向下一个bucket的指针

    struct bucket *pLast;//指向上一个bucket的指针

    char arKey[1];//必须放在最后，主要是为了实现变长结构体

} Bucket;

 

typedef struct _hashtable {

    uint nTableSize;             //哈希表的大小

    uint nTableMask;             //数值上等于nTableSize- 1

    uint nNumOfElements;         //记录了当前HashTable中保存的记录数

    ulong nNextFreeElement;      //指向下一个空闲的Bucket

    Bucket *pInternalPointer;   //这个变量用于数组反转

    Bucket *pListHead;           //指向Bucket的头

    Bucket *pListTail;           //指向Bucket的尾

    Bucket **arBuckets;

    dtor_func_t pDestructor;     //函数指针，数组增删改查时自动调用，用于某些清理操作

    zend_bool persistent;         //是否持久

    unsigned char nApplyCount;

    zend_bool bApplyProtection;  //和nApplyCount一起起作用，防止数组遍历时无限递归

#if ZEND_DEBUG

    int inconsistent;

#endif

} HashTable;

希望大家能好好看看上面的定义，有些东西我将出来反而会说不明白，不如大家看看代码浅显明了。PHP的数组，其实是一个带有头结点的双向链表，其中HashTable是头，Bucket存储具体的结点信息。

1.3     HashTable内部函数分析
1.3.1    宏HASH_PROTECT_RECURSION
#defineHASH_PROTECT_RECURSION(ht)                                                     \

    if ((ht)->bApplyProtection) {                                                       \

        if ((ht)->nApplyCount++ >= 3){                                                \

            zend_error(E_ERROR, "Nestinglevel too deep - recursive dependency?"); \

        }                                                                               \

    }

这个宏主要用来防止循环引用。

1.3.2    宏ZEND_HASH_IF_FULL_DO_RESIZE
#defineZEND_HASH_IF_FULL_DO_RESIZE(ht)            \

    if ((ht)->nNumOfElements >(ht)->nTableSize) {  \

        zend_hash_do_resize(ht);                    \

    }

         这个宏的作用是检查目前HashTable中的元素个数是否大于了总的HashTable的大小，如果个数大于了HashTable的大小，那么我们就重新分配空间。我们看一下zend_hash_do_resize

static int zend_hash_do_resize(HashTable*ht)

{

    Bucket **t;

    IS_CONSISTENT(ht);

    if ((ht->nTableSize<<1)>0){   /* Let's double the table size */

        t = (Bucket**) perealloc_recoverable(ht->arBuckets,

(ht->nTableSize<<1)* sizeof(Bucket*), ht->persistent);

        if (t){

            HANDLE_BLOCK_INTERRUPTIONS();

            ht->arBuckets = t;

            ht->nTableSize = (ht->nTableSize<<1);

            ht->nTableMask = ht->nTableSize-1;

            zend_hash_rehash(ht);

            HANDLE_UNBLOCK_INTERRUPTIONS();

            return SUCCESS;

        }

        return FAILURE;

    }

    return SUCCESS;

}  

         从上面的代码中我们可以看出，HashTable在分配空间的时候，新分配的空间等于原有空间的2倍。

1.3.3    函数 _zend_hash_init
这个函数是用来初始化HashTable的，我们先看一下代码：

ZEND_API int _zend_hash_init(HashTable*ht, uint nSize, hash_func_t pHashFunction, dtor_func_t pDestructor, zend_bool persistent ZEND_FILE_LINE_DC)

{

    uint i = 3; //HashTable的大小默认无2的3次方

    Bucket **tmp;

 

    SET_INCONSISTENT(HT_OK);

 

    if (nSize>=0x80000000){

        ht->nTableSize = 0x80000000;

    } else{

        while ((1U << i)< nSize){

            i++;

        }

        ht->nTableSize = 1 << i;

    }

 

    ht->nTableMask = ht->nTableSize-1;

    ht->pDestructor = pDestructor;

    ht->arBuckets =NULL;

    ht->pListHead =NULL;

    ht->pListTail =NULL;

    ht->nNumOfElements = 0;

    ht->nNextFreeElement = 0;

    ht->pInternalPointer = NULL;

    ht->persistent = persistent;

    ht->nApplyCount =0;

    ht->bApplyProtection = 1;

   

    /* Uses ecalloc() so that Bucket* == NULL */

    if (persistent){

        tmp = (Bucket **) calloc(ht->nTableSize,sizeof(Bucket*));

        if (!tmp){

            return FAILURE;

        }

        ht->arBuckets = tmp;

    } else{

        tmp = (Bucket **) ecalloc_rel(ht->nTableSize,sizeof(Bucket*));

        if (tmp){

            ht->arBuckets = tmp;

        }

    }

   

    return SUCCESS;

}

可以看出，HashTable的大小被初始化为2的n次方，另外我们看到有两种内存方式，一种是calloc，一种是ecalloc_rel，这两中内存分配方式我们细讲了，有兴趣的话大家可以自己查一查。

1.3.4    函数_zend_hash_add_or_update
这个函数向HashTable中添加或者修改元素信息

ZEND_API int _zend_hash_add_or_update(HashTable*ht,constchar *arKey, uint nKeyLength,void*pData, uint nDataSize,void**pDest,int flag ZEND_FILE_LINE_DC)

{

    ulong h;

    uint nIndex;

    Bucket *p;

 

    IS_CONSISTENT(ht);

 

    if (nKeyLength<=0){

#if ZEND_DEBUG

        ZEND_PUTS("zend_hash_update: Can't put inempty key\n");

#endif

        return FAILURE;

    }

 

    h = zend_inline_hash_func(arKey, nKeyLength);

    nIndex = h & ht->nTableMask;

 

    p = ht->arBuckets[nIndex];

    while (p!=NULL){

        if ((p->h== h)&&(p->nKeyLength== nKeyLength)){

            if (!memcmp(p->arKey, arKey, nKeyLength)){

                if (flag & HASH_ADD){

                    return FAILURE;

                }

                HANDLE_BLOCK_INTERRUPTIONS();

#if ZEND_DEBUG

                if (p->pData == pData){

                    ZEND_PUTS("Fatal error in zend_hash_update:p->pData == pData\n");

                    HANDLE_UNBLOCK_INTERRUPTIONS();

                    return FAILURE;

                }

#endif

                if (ht->pDestructor){

                    ht->pDestructor(p->pData);

                }

                UPDATE_DATA(ht, p, pData, nDataSize);

                if (pDest) {

                    *pDest = p->pData;

                }

                HANDLE_UNBLOCK_INTERRUPTIONS();

                return SUCCESS;

            }

        }

        p = p->pNext;

    }

   

    p = (Bucket*) pemalloc(sizeof(Bucket)-1 + nKeyLength, ht->persistent);

    if (!p){

        return FAILURE;

    }

    memcpy(p->arKey, arKey, nKeyLength);

    p->nKeyLength = nKeyLength;

    INIT_DATA(ht, p, pData, nDataSize);

    p->h = h;

    CONNECT_TO_BUCKET_DLLIST(p, ht->arBuckets[nIndex]);

    if (pDest){

        *pDest = p->pData;

    }

 

    HANDLE_BLOCK_INTERRUPTIONS();

    CONNECT_TO_GLOBAL_DLLIST(p, ht);

    ht->arBuckets[nIndex]= p;

    HANDLE_UNBLOCK_INTERRUPTIONS();

 

    ht->nNumOfElements++;

    ZEND_HASH_IF_FULL_DO_RESIZE(ht);       /* If the Hash table is full, resize it */

    return SUCCESS;

}

1.3.5    宏CONNECT_TO_BUCKET_DLLIST
#define CONNECT_TO_BUCKET_DLLIST(element, list_head)        \

    (element)->pNext= (list_head);                         \

    (element)->pLast= NULL;                                \

    if((element)->pNext) {                                 \

        (element)->pNext->pLast =(element);                \

    }

这个宏是将bucket加入到bucket链表中

1.3.6    其他函数或者宏定义
我们只是简单的介绍一下HashTable，如果你想细致的了解HashTable的话，建议您看看php的源代码，关于HashTable的代码在Zend/zend_hash.h 和Zend/zend_hash.c中。

zend_hash_add_empty_element 给函数增加一个空元素

zend_hash_del_key_or_index 根据索引删除元素

zend_hash_reverse_apply  反向遍历HashTable

zend_hash_copy  拷贝

_zend_hash_merge  合并

zend_hash_find  字符串索引方式查找

zend_hash_index_find  数值索引方法查找

zend_hash_quick_find  上面两个函数的封装

zend_hash_exists  是否存在索引

zend_hash_index_exists  是否存在索引

zend_hash_quick_exists  上面两个方法的封装

1.4     C扩展常用HashTable函数
虽然HashTable看起来有点复杂，但是使用却是很方便的，我们可以用下面的函数对HashTable进行初始化和赋值。

2005 年地方院校招生人数

PHP语法

C语法

意义

$arr = array()

array_init(arr);

初始化数组

$arr[] = NULL;

add_next_index_null(arr);

 

$arr[] = 42;

add_next_index_long(arr, 42);

 

$arr[] = true;

add_next_index_bool(arr, 1);

 

$arr[] = 3.14;

add_next_index_double(3.14);

 

$arr[] = ‘foo’;

add_next_index_string(arr, “foo”, 1);

1的意思进行字符串拷贝

$arr[] = $myvar;

add_next_index_zval(arr, myvar);

 

$arr[0] = NULL;

add_index_null(arr, 0);

 

$arr[1] = 42;

add_index_long(arr, 1, 42);

 

$arr[2] = true;

add_index_bool(arr, 2, 1);

 

$arr[3] = 3.14;

add_index_double(arr, 3, 3,14);

 

$arr[4] = ‘foo’;

add_index_string(arr, 4, “foo”, 1);

 

$arr[5] = $myvar;

add_index_zval(arr, 5, myvar);

 

$arr[“abc”] = NULL;

add_assoc_null(arr, “abc”);

 

$arr[“def”] = 711;

add_assoc_long(arr, “def”, 711);

 

$arr[“ghi”] = true;

add_assoc_bool(arr, ghi”, 1);

 

$arr[“jkl”] = 1.44;

add_assoc_double(arr, “jkl”, 1.44);

 

$arr[“mno”] = ‘baz’;

add_assoc_string(arr, “mno”, “baz”, 1);

 

$arr[‘pqr’] = $myvar;

add_assoc_zval(arr, “pqr”, myvar);

 

1.5     任务和实验
说了这么多，我们实验一下。

任务：返回一个数组，数组中的数据如下：

Array

(

   [0] => for test

   [42] => 123

   [for test. for test.] => 1

   [array] => Array

       (

           [0] => 3.34

       )

)

代码实现：

PHP_FUNCTION(test)

{

    zval* t;

 

    array_init(return_value);

    add_next_index_string(return_value,"for test",1);

    add_index_long(return_value,42,123);

    add_assoc_double(return_value,"for test. for test.",1.0);

   

    ALLOC_INIT_ZVAL(t);

    array_init(t);

    add_next_index_double(t,3.34);

 

    add_assoc_zval(return_value,"array", t);

}

很简单吧，还记得return_value吗？

https://blog.csdn.net/niujiaming0819/article/details/8568587


