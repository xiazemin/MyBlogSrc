---
title: current_execute_data
layout: post
category: php
author: 夏泽民
---
1. Zend引擎主要包含两个核心部分：编译、执行：

                           

    执行阶段主要用到的数据结构：

          opcode： php代码编译产生的zend虚拟机可识别的指令，php7有173个opcode，定义在 zend_vm_opcodes.hPHP中的所有语法实现都是由这些opcode组成的。

         

复制代码
struct _zend_op {
    const void *handler; //对应执行的C语言function，即每条opcode都有一个C function处理
    znode_op op1;   //操作数1
    znode_op op2;   //操作数2
    znode_op result; //返回值
    uint32_t extended_value; 
    uint32_t lineno; 
    zend_uchar opcode;  //opcode指令
    zend_uchar op1_type; //操作数1类型
    zend_uchar op2_type; //操作数2类型
    zend_uchar result_type; //返回值类型
};
复制代码
         zend_op_array : zend引擎执行阶段的输入数据结构，整个执行阶段都是操作这个数据结构。

             

                            

 

 

      　　　　　　 zend_op_array有三个核心部分：opcode指令(对应c的指令)

                                                   字面量存储(变量初始值、调用的函数名称、类名称、常量名称等等称之为字面量)

                                                   变量分配的情况 (当前array定义的变量 临时变量的数量 编号，执行初始化一次性分配zval，使用时完全按照标号索引不是根据变量名)

         

           zend_executor_globals     PHP整个生命周期中最主要的一个结构，是一个全局变量，在main执行前分配(非ZTS下)，直到PHP退出，它记录着当前请求全部的信息，经常见到的一个宏EG操作的就是这个结构。

                                定义在zend_globals.h中：

 

                                   

 

                

               zend_execute_data  是执行过程中最核心的一个结构，每次函数的调用、include/require、eval等都会生成一个新的结构，它表示当前的作用域、代码的执行位置以及局部变量的分配等等，等同于机器码执行过程中stack的角色，后面分析具体执行流程的时候会详细分析其作用。 

              zend_execute_data与zend_op_array的关联关系：

                                         

2.执行过程

        Zend的executor与linux二进制程序执行的过程是非常类似的。

        在C程序执行时有两个寄存器ebp、esp分别指向当前作用栈的栈顶、栈底，局部变量全部分配在当前栈，函数调用、返回通过call、ret指令完成，调用时call将当前执行位置压入栈中，返回时ret将之前执行位置出栈，跳回旧的位置继续执行。

        Zend VM中zend_execute_data就扮演了这两个角色，zend_execute_data.prev_execute_data保存的是调用方的信息，实现了call/ret，zend_execute_data后面会分配额外的内存空间用于局部变量的存储，实现了ebp/esp的作用。

                    a. 为当前作用域分配一块内存，充当运行栈，zend_execute_data结构、所有局部变量、中间变量等等都在此内存上分配

                    b.初始化全局变量符号表，然后将全局执行位置指针EG(current_execute_data)指向步骤a新分配的zend_execute_data，然后将zend_execute_data.opline指向op_array的起始位置

                    c.从EX(opline)开始调用各opcode的C处理handler(即_zend_op.handler)，每执行完一条opcode将EX(opline)++继续执行下一条，直到执行完全部opcode

                                if语句将根据条件的成立与否决定EX(opline) + offset所加的偏移量，实现跳转

                                如果是函数调用，则首先从EG(function_table)中根据function_name取出此function对应的编译完成的zend_op_array，然后像步骤a一样新分配一个zend_execute_data结构，将EG(current_execute_data)赋值给新结构的prev_execute_data，再将EG(current_execute_data)指向新的zend_execute_data，最后从新的zend_execute_data.opline开始执行，切换到函数内部，函数执行完以后将EG(current_execute_data)重新指向EX(prev_execute_data)，释放分配的运行栈，销毁局部变量，继续从原来函数调用的位置执行

                                类方法的调用与函数基本相同

                    d.全部opcode执行完成后将步骤a分配的内存释放，这个过程会将所有的局部变量"销毁"，执行阶段结束

                                    

 

                              首先根据zend_execute_data、当前zend_op_array中局部/临时变量数计算需要的内存空间，编译阶段zend_op_array的结果，在编译过程中已经确定当前作用域下有多少个局部变量(func->op_array.last_var)、临时/中间/无用变量(func->op_array.T)，从而在执行之初就将他们全部分配完成。
<!-- more -->
{% raw %}
https://www.cnblogs.com/hellohell/p/9101803.html
http://www.lvesu.com/blog/php/migration55.internals.php

先打印一下php调用过程：







在增加一张异常调用的流程图：



今天稍微对php做下总结，首先介绍最重要的两个数据结构，以及两个结构间的数据传递

struct _zend_op_array {
/* Common elements */
zend_uchar type;
const char *function_name;
zend_class_entry *scope;
zend_uint fn_flags;
union _zend_function *prototype;
zend_uint num_args;
zend_uint required_num_args;
zend_arg_info *arg_info;
/* END of common elements */
 
 
zend_uint *refcount;
 
 
zend_op *opcodes;
zend_uint last;
 
 
zend_compiled_variable *vars;
int last_var;
 
 
zend_uint T;
 
zend_literal *literals;
int last_literal;
 
 
...
};



不重点介绍的属性暂时省略


struct _zend_execute_data {
struct _zend_op *opline;
zend_function_state function_state;
zend_op_array *op_array;
zval *object;
HashTable *symbol_table;
struct _zend_execute_data *prev_execute_data;
zval *old_error_reporting;
zend_bool nested;
zval **original_return_value;
zend_class_entry *current_scope;
zend_class_entry *current_called_scope;
zval *current_this;
struct _zend_op *fast_ret; /* used by FAST_CALL/FAST_RET (finally keyword) */
zval *delayed_exception;
call_slot *call_slots;
call_slot *call;
};


zend_execute_data的数据结构大部分是指针，指针指向的内容是这样分配的
5.6很清楚的画出了内存分配图
/*
 * Stack Frame Layout (the whole stack frame is allocated at once)
 * ==================
 *
 *                             +========================================+
 *                             | zend_execute_data                      |<---+
 *                             |     EX(function_state).arguments       |--+ |
 *                             |  ...                                   |  | |
 *                             | ARGUMENT [1]                           |  | |
 *                             | ...                                    |  | |
 *                             | ARGUMENT [ARGS_NUMBER]                 |  | |
 *                             | ARGS_NUMBER                            |<-+ |
 *                             +========================================+    |
 *                                                                           |
 *                             +========================================+    |
 *                             | TMP_VAR[op_arrat->T-1]                 |    |
 *                             | ...                                    |    |
 *     EX_TMP_VAR_NUM(0) ----> | TMP_VAR[0]                             |    |
 *                             +----------------------------------------+    |
 * EG(current_execute_data) -> | zend_execute_data                      |    |
 *                             |     EX(prev_execute_data)              |----+
 *                             +----------------------------------------+
 *     EX_CV_NUM(0) ---------> | CV[0]                                  |--+
 *                             | ...                                    |  |
 *                             | CV[op_array->last_var-1]               |  |
 *                             +----------------------------------------+  |
 *                             | Optional slot for CV[0] zval*          |<-+
 *                             | ...                                    |
 *                             | ...for CV [op_array->last_var-1] zval* |
 *                             +----------------------------------------+
 *           EX(call_slots) -> | CALL_SLOT[0]                           |
 *                             | ...                                    |
 *                             | CALL_SLOT[op_array->nested_calls-1]    |
 *                             +----------------------------------------+
 * zend_vm_stack_frame_base -> | ARGUMENTS STACK [0]                    |
 *                             | ...                                    |
 * zend_vm_stack_top --------> | ...                                    |
 *                             | ...                                    |
 *                             | ARGUMENTS STACK [op_array->used_stack] |
 *                             +----------------------------------------+
 */
这里分配有个条件判断在5.2是没有的
(op_array->fn_flags & ZEND_ACC_GENERATOR) != 0
当如果用到了yield函数时就会触发该逻辑，从而再分配上面的堆栈结构是会给prev_execute_data单独分配空间，并且指向TMP_VAR变量的上面的内存位置。
这个数据结构重要的是三个属性EX_TMP_VAR_NUM(临时变量的空间)，EX_CV_NUM(缓存变量的空间)，zend_vm_stack_top(函数参数的空间)




注释下php7的结构简化了不少

/*
 * Stack Frame Layout (the whole stack frame is allocated at once)
 * ==================
 *
 *                             +========================================+
 * EG(current_execute_data) -> | zend_execute_data                      |
 *                             +----------------------------------------+
 *     EX_CV_NUM(0) ---------> | VAR[0] = ARG[1]                        |
 *                             | ...                                    |
 *                             | VAR[op_array->num_args-1] = ARG[N]     |
 *                             | ...                                    |
 *                             | VAR[op_array->last_var-1]              |
 *                             | VAR[op_array->last_var] = TMP[0]       |
 *                             | ...                                    |
 *                             | VAR[op_array->last_var+op_array->T-1]  |
 *                             | ARG[N+1] (extra_args)                  |
 *                             | ...                                    |
 *                             +----------------------------------------+
 */



稍微了解了上面两个数据结构后，下面就讲zend_op_array与zend_execute之间的相互关系
php分为几个阶段包括生成opcode阶段和执行opcode阶段，其实分别对应的就是上面两个数据结构，
并且两个数据结构都是在解析到新的函数时分配新的空间，然后层层嵌套，最外层总是有个大的op_array与execute_data,具体点说就是这两个数据结构存储的是当前函数下的变量环境。
然后就是上面两个不同阶段存储该阶段应该存储的数据，然后可供下一层调用。


第一个例子$var = 1
此处省略掉语法解析(| variable '=' expr{ zend_check_writable_variable(&$1); zend_do_assign(&$$, &$1, &$3 TSRMLS_CC); })，

直接到opconde生成阶段


以下代码是var该字符串代表的变量的信息，为什么这么说因为现在它还不能成为变量，opconde解析的$var的返回值才是变量
fetch_simple_variable -> lookup_cv
核心代码
i = op_array->last_var;
op_array->last_var++;

void fetch_simple_variable_ex(znode *result, znode *varname, int bp, zend_uchar op TSRMLS_DC) /* {{{ */
{
	zend_op opline;
	zend_op *opline_ptr;
	zend_llist *fetch_list_ptr;
 
	if (varname->op_type == IS_CONST) {
		ulong hash;
 
		if (Z_TYPE(varname->u.constant) != IS_STRING) {
			convert_to_string(&varname->u.constant);
		}
 
		hash = str_hash(Z_STRVAL(varname->u.constant), Z_STRLEN(varname->u.constant));
		if (!zend_is_auto_global_quick(Z_STRVAL(varname->u.constant), Z_STRLEN(varname->u.constant), hash TSRMLS_CC) &&
		    !(Z_STRLEN(varname->u.constant) == (sizeof("this")-1) &&
		      !memcmp(Z_STRVAL(varname->u.constant), "this", sizeof("this") - 1)) &&
		    (CG(active_op_array)->last == 0 ||
		     CG(active_op_array)->opcodes[CG(active_op_array)->last-1].opcode != ZEND_BEGIN_SILENCE)) {
			result->op_type = IS_CV;
			result->u.op.var = lookup_cv(CG(active_op_array), Z_STRVAL(varname->u.constant), Z_STRLEN(varname->u.constant), hash TSRMLS_CC);
			Z_STRVAL(varname->u.constant) = (char*)CG(active_op_array)->vars[result->u.op.var].name;
			result->EA = 0;
			return;
		}
	}
 
	if (bp) {
		opline_ptr = &opline;
		init_op(opline_ptr TSRMLS_CC);
	} else {
		opline_ptr = get_next_op(CG(active_op_array) TSRMLS_CC);
	}
 
	opline_ptr->opcode = op;
	opline_ptr->result_type = IS_VAR;
	opline_ptr->result.var = get_temporary_variable(CG(active_op_array));
	SET_NODE(opline_ptr->op1, varname);
	GET_NODE(result, opline_ptr->result);
	SET_UNUSED(opline_ptr->op2);
	opline_ptr->extended_value = ZEND_FETCH_LOCAL;
 
	if (varname->op_type == IS_CONST) {
		CALCULATE_LITERAL_HASH(opline_ptr->op1.constant);
		if (zend_is_auto_global_quick(Z_STRVAL(varname->u.constant), Z_STRLEN(varname->u.constant), Z_HASH_P(&CONSTANT(opline_ptr->op1.constant)) TSRMLS_CC)) {
			opline_ptr->extended_value = ZEND_FETCH_GLOBAL;
		}
	}
 
	if (bp) {
		zend_stack_top(&CG(bp_stack), (void **) &fetch_list_ptr);
		zend_llist_add_element(fetch_list_ptr, opline_ptr);
	}
}

如果是$$var就会走到下面逻辑进行opline赋值，当execute阶段就会执行ZEND_FETCH_W_SPEC_CV_VAR_HANDLER找到$var变量的值再进行变量查询.


op_array->vars[i].name = zend_new_interned_string(name, name_len + 1, 1 TSRMLS_CC);
op_array->vars[i].name_len = name_len;
op_array->vars[i].hash_value = hash_value;
很好理解，为$var这个变量属性分配了空间而不是为变量分配空间，利用了op_array两个属性last_var变量位置和vars数组的对应关系


那$var返回值是什么呢
result->op_type = IS_CV;
result->u.op.var = op_array->last_var;
但是前提大家需要知道语法解析时将字符串或者整形统一解析到znode->u.constant中，

result的数据结构是zonde，

后续opcode阶段如果是变量赋值的是znode_op的var属性，即偏移量，znode_op.zv属性，即常量信息(该值有pass_two赋值)

typedef union _znode_op {
	zend_uint      constant;
	zend_uint      var;
	zend_uint      num;
	zend_ulong     hash;
	zend_uint      opline_num; /*  Needs to be signed */
	zend_op       *jmp_addr;
	zval          *zv;
	zend_literal  *literal;
	void          *ptr;        /* Used for passing pointers from the compile to execution phase, currently used for traits */
} znode_op;
 
typedef struct _znode { /* used only during compilation */
	int op_type;
	union {
		znode_op op;
		zval constant; /* replaced by literal/zv */
		zend_op_array *op_array;
		zend_ast *ast;
	} u;
	zend_uint EA;      /* extended attributes */
} znode;

这里补充下，语法扫描获取znode后，进入语法解析阶段，此刻op_array中有个特殊属性literals，该属性是个数组会提前分配好，最终可进行opcode阶段优化见update_op1_const函数，将变量转为常量。


accel_startup:
accelerator_orig_compile_file = zend_compile_file; // 保存原生handle
zend_compile_file = persistent_compile_file; //赋值新的handle
用persistent_compile_file -> compile_and_cache_file -> cache_script_in_shared_memory -> zend_accel_script_optimize
->zend_accel_optimize->zend_optimize->replace_var_by_const->update_op1_const最后会利用literals数组将变量转换为常量更改opline的op1或者op2

该数组index与value方式进行存储

#define SET_NODE(target, src) do { \
		target ## _type = (src)->op_type; \
		if ((src)->op_type == IS_CONST) { \
			target.constant = zend_add_literal(CG(active_op_array), &(src)->u.constant TSRMLS_CC); \
		} else { \
			target = (src)->u.op; \
		} \
	} while (0)

如果是常量op1->constant = index,  其种value存在将数据存到literals中，接下来的用途见pass_two;

compile->pass_two 此时会生成opcode的回调op->handler，并且会从constant的index中将value赋值给opline中的op1.zv，

这样在真正execute阶段用的就是op1.zv获取常量信息


while (opline < end) {
		if (opline->op1_type == IS_CONST) {
			opline->op1.zv = &op_array->literals[opline->op1.constant].constant;
		}
		if (opline->op2_type == IS_CONST) {
			opline->op2.zv = &op_array->literals[opline->op2.constant].constant;
		}
		。。。。。。
		
		ZEND_VM_SET_OPCODE_HANDLER(opline);
		opline++;
	}
#define IS_CONST	(1<<0)
#define IS_TMP_VAR	(1<<1)
#define IS_VAR		(1<<2)
#define IS_UNUSED	(1<<3)	/* Unused variable */
#define IS_CV		(1<<4)	/* Compiled variable */
如：
opline->result_type = IS_TMP_VAR; //
opline->result.var = get_temporary_variable(CG(active_op_array));
opline->result_type = IS_VAR;
opline->result.var = get_temporary_variable(CG(active_op_array));
最后生成的函数就是
ZEND_ASSIGN_SPEC_VAR_TMP_HANDLER
VAR与TMP均来自临时变量，但是两者用的数据结构不同，具体可见以下两个函数
value = _get_zval_ptr_tmp(opline->op2.var, execute_data, &free_op2 TSRMLS_CC);
variable_ptr_ptr = _get_zval_ptr_ptr_var(opline->op1.var, execute_data, &free_op1 TSRMLS_CC);
完！



//又利用了op_array的一个属性op_array>T++,//此刻仍然是个位置变量，要知道一个op_array可以有n多变量，这个属性就像一个全局数组，所以位置var//足可以代表一个变量。


上面分析的是($)(var)的过程
接下来该赋值了通过语法解析知道调用的是zend_do_assign函数
原型zend_do_assign(znode *result, znode *variable, znode *value TSRMLS_DC)
result是返回值，variable就是上面介绍的变量的返回值(result.op_type)，value就是常量1
opline->handle = ZEND_ASSIGN_SPEC_CV_CONST_HANDLER


刚才说过opcode的解析阶段完成，接下来就是要执行opcode了
ZEND_ASSIGN_SPEC_CV_CONST_HANDLER -> _get_zval_ptr_ptr_cv_BP_VAR_W -> _get_zval_cv_lookup_BP_VAR_W -> zend_assign_const_to_variable


variable_ptr_ptr = _get_zval_ptr_ptr_cv_BP_VAR_W(execute_data, opline->op1.var TSRMLS_CC);
该函数原型：
static zend_always_inline zval **_get_zval_ptr_ptr_cv_BP_VAR_W(const zend_execute_data *execute_data, zend_uint var TSRMLS_DC)
{
zval ***ptr = EX_CV_NUM(execute_data, var); //这是execute_data的cv第一部分CV[i]的值其实是个指针


if (UNEXPECTED(*ptr == NULL)) {
return _get_zval_cv_lookup_BP_VAR_W(ptr, var TSRMLS_CC);//这是execute_data的cv第一部分CV[i]的值其实是个指针,这个指针真正分配的空间(该空间可能是execute_data的cv分配，也可能是符号表来分配，两者只能选其一，所以代码里可以看到size_t CVs_size = ZEND_MM_ALIGNED_SIZE(sizeof(zval **) * op_array->last_var * (EG(active_symbol_table) ? 1 : 2));如果没有符号表就需要分配两倍空间)
}
return *ptr;
}
可以看到就是这个函数做了execute_data与op_array之间的关联


总结第一个例子，opcode解析阶段op_array存放的是变量的一些基本信息，opcode的执行阶段execute_data分配空间存放数据,两者联系就是通过opcode的变量位置等




接下来第二个例子函数调用我们关注的是参数的传递，分两部分(仍然忽视掉语法解析直接到达opcode生成)，1 函数解析 2 函数调用：
<?php
function foo($arg1)
{
    print($arg1);
}


$bar = 'hello php';
foo($bar);


函数解析:
zend_do_begin_function_declaration(进行CG(active_op_array)的切换，并且将新的函数注册到CG(function_table)全局变量中) -> zend_do_receive_param -> zend_do_end_function_declaration(退出当前，将之前的op_array重置)


zend_do_receive_param:
(同上对$arg1变量信息进行存储，然后每一个参数就会调用一次receive)
该函数利用到了op_array的两个属性num_args与arg_info


CG(active_op_array)->num_args++; 作用显而易见
以及cur_arg_info = &CG(active_op_array)->arg_info[CG(active_op_array)->num_args-1];
将参数信息进行存储，将来用来进行参数类型的对比检查


函数调用：
zend_do_begin_function_call -->  zend_do_pass_param() --> zend_do_end_function_call


zend_do_begin_function_call:
函数解析时曾经注册到函数表，此函数会判断函数名是否存在
获取函数后，会将函数压入堆栈，方便处理下面容易获取函数本身
zend_stack_push(&CG(function_call_stack), (void *) &function, sizeof(zend_function *));


zend_do_pass_param：
会将$bar如同前面介绍的一样进行变量处理，然后生成
op->handle=ZEND_SEND_VAR_NO_REF


在解析opcode时调用的是ZEND_SEND_VAR_NO_REF_SPEC_CV_HANDLER函数：
该函数如同上面介绍一样分配了$bar的存储空间，然后需要注意的是zend_vm_stack_push(varptr TSRMLS_CC);
用到了zend_excute结构中的堆栈，将函数空间的指针压入到了堆栈中。


zend_do_end_function_call：
op->handle=ZEND_DO_FCALL_BY_NAME


在解析opcode时调用的是ZEND_DO_FCALL_BY_NAME_SPEC_HANDLER函数：
会将参数的数量压入到堆栈中，如同c的函数调用压入参数一样，arg1，arg2，argnum
zend_vm_stack_push_args(num_args TSRMLS_CC);//注意先压入参数才切换op_array
然后zend_execute(op_array);
调用完后回到本函数调用zend_vm_stack_clear_multiple进行堆栈释放


从这一步开始进入了函数中，还记得函数的zend_do_receive_param中的op->handle=ZEND_RECV_SPEC_HANDLER:
zval **param = zend_vm_stack_get_arg(arg_num TSRMLS_CC);

static zend_always_inline zval** zend_vm_stack_get_arg(int requested_arg TSRMLS_DC)
{
return zend_vm_stack_get_arg_ex(EG(current_execute_data)->prev_execute_data, requested_arg);
}

为什么要用prev_execute_data因为函数压栈是在当前excute_data之前的excute_data完成，
实际上，在真正执行函数之前，php会将参数个数入栈。



先根据参数个数把堆栈中的参数列表取出来，然后进行参数验证，
顺便var_ptr = _get_zval_ptr_ptr_cv_BP_VAR_W(execute_data, opline->result.var TSRMLS_CC);
将函数的参数变量$arg1进行赋值


整个过程结束



上面是拆分讲解一个函数的调用过程，当将所有程序解析成op_array数组后，就会调用execute_ex来执行所有的opcode数组。


zend_execute_scripts -> zend_execute -> zend_execute_ex -> execute_ex -> i_create_execute_data_from_op_array 


	if (0) {
zend_vm_enter:
		execute_data = i_create_execute_data_from_op_array(EG(active_op_array), 1 TSRMLS_CC);
	}
 
	LOAD_REGS();
	LOAD_OPLINE();
 
	while (1) {
    	int ret;
#ifdef ZEND_WIN32
		if (EG(timed_out)) {
			zend_timeout(0);
		}
#endif
 
		if ((ret = OPLINE->handler(execute_data TSRMLS_CC)) > 0) {
			switch (ret) {
				case 1:
					EG(in_execution) = original_in_execution;
					return;
				case 2:
					goto zend_vm_enter;
					break;
				case 3:
					execute_data = EG(current_execute_data);
					break;
				default:
					break;
			}
		}
 
	}
	zend_error_noreturn(E_ERROR, "Arrived at end of main loop which shouldn't happen");
}
以上是个死循环，解析op_array数组，需要注意的是返回值
#define ZEND_VM_CONTINUE()         return 0
#define ZEND_VM_RETURN()           return 1   返回return是函数终止，
#define ZEND_VM_ENTER()            return 2  函数调用
#define ZEND_VM_LEAVE()            return 3  函数退出
ZEND_VM_RETURN函数返回returen终止，5.2版本很少有返回1，但是5.3增加了yield调用后，yield调用的opcode基本上都会返回1，从而
函数终止，有个疑问，return了后下面如何执行？
这里面有个需要注意的点就是，当我们调用函数的时候大家知道opcode解析的函数是
zend_do_fcall_common_helper_SPEC 
该函数分为两部分，
typedef union _zend_function {
zend_uchar type; /* ...﹚... #define ZEND_USER_FUNCTION 2
MUST be the first element of this struct! */
struct {
zend_uchar type; /* never used */
char *function_name; //ㄧ..
zend_class_entry *scope; //ㄧ.┮..办
zend_uint fn_flags; // ..猭....单ZEND_ACC_STATIC单
union _zend_function *prototype; //ㄧ.
zend_uint num_args; //....
zend_uint required_num_args; //惠璶....
zend_arg_info *arg_info; //..獺.
zend_bool pass_rest_by_reference;
unsigned char return_reference; //
} common;
zend_op_array op_array; //ㄧ.い巨
zend_internal_function internal_function;
} zend_function;


1 内部C函数(ZEND_INTERNAL_FUNCTION):内部函数在zend_register_functions时候就注册到了函数表，其中internal_function.handler指向C函数(函数指针)
通过opcode解析函数名到函数表中查找即可获取到函数指针，进行调用
2 php函数(ZEND_USER_FUNCTION):会继续调用zend_execute，所以刚才说的ZEND_VM_RETURN终止的只是具体某个函数而已，大家如果在一个函数
中写yield，函数就不会继续执行了，就是这个道理。




顺带着介绍下词法解析过程

Zend/zend_language_scanner.l  词法解析规则文件
Zend/zend_language_parser.y   语法分析规则文件

bison -o zend_language_parser.c zend_language_parser.y
在Zend目录下就会生成语法解析器zend_language_parser.c。
yyarse（）是bison生成的分析器的主函数。 调用yyarse()，如果一切顺利，那么上例中的g_root将指向一个完成的语法树。
#define yylex zendlex
#define YYSTYPE znode
yyparse -> zendlex -> lex_scan
lex_scan会解析YYSTYPE字段，并且将znode.u.constant进行词法解析赋值数据

语法扫描(lex_scan)前都会进行该函数调用进行准备，可以参考函数token_get_all的实现

static void yy_scan_buffer(char *str, unsigned int len TSRMLS_DC)
{
	YYCURSOR = (YYCTYPE*)str;
	SCNG(yy_start) = YYCURSOR;
	YYLIMIT  = YYCURSOR + len;
}



顺便说一下php7：

compile_file -> zendparse(yyparse) -> zend_compile_top_stmt -> zend_compile_stmt 生成opcode，因为php7中间用了抽象语法树，
需要根据抽象语法树的节点进行分析后获得最的opcode
在这里补充下一个知识点:
ZEND_ASSIGN_ADD_SPEC_VAR_CONST_HANDLER 
	if (RETURN_VALUE_USED(opline)) {
		PZVAL_LOCK(*var_ptr);
		EX_T(opline->result.var).var.ptr = *var_ptr;
	}
这里放的是临时变量的var属性ptr，为什么是指针，而不是下面的tmp_var,是因为ptr指向的zval是不能立马释放的，是需要assign
赋值给其他变量用，也就是多个*zal 共同指向的结构，这个时候采用的就是存放到临时变量的var中的ptr属性。
再看个例子
ZEND_ADD_SPEC_CV_TMP_HANDLER
static int ZEND_FASTCALL  ZEND_ADD_SPEC_CV_TMP_HANDLER(ZEND_OPCODE_HANDLER_ARGS)
{
	USE_OPLINE
	zend_free_op free_op2;
 
	SAVE_OPLINE();
	fast_add_function(&EX_T(opline->result.var).tmp_var,
		_get_zval_ptr_cv_BP_VAR_R(execute_data, opline->op1.var TSRMLS_CC),
		_get_zval_ptr_tmp(opline->op2.var, execute_data, &free_op2 TSRMLS_CC) TSRMLS_CC);
 
	zval_dtor(free_op2.var);
	CHECK_EXCEPTION();
	ZEND_VM_NEXT_OPCODE();
}
很明显放到了赋值给了临时变量，为什么是临时变量，因为该变量不引用其他指针数据，所以释放比较简单

a++的opcode
zend_do_post_incdec 的opcode opline->result_type = IS_TMP_VAR;很明显是个tmp变量

	opline = get_next_op(CG(active_op_array) TSRMLS_CC);
	opline->opcode = op;
	SET_NODE(opline->op1, op1);
	SET_UNUSED(opline->op2);
	opline->result_type = IS_TMP_VAR;
	opline->result.var = get_temporary_variable(CG(active_op_array));
	GET_NODE(result, opline->result);

static int ZEND_FASTCALL  ZEND_POST_INC_SPEC_VAR_HANDLER(ZEND_OPCODE_HANDLER_ARGS)
{
	USE_OPLINE
	zend_free_op free_op1;
	zval **var_ptr, *retval;
 
	SAVE_OPLINE();
	var_ptr = _get_zval_ptr_ptr_var(opline->op1.var, execute_data, &free_op1 TSRMLS_CC);
 
	if (IS_VAR == IS_VAR && UNEXPECTED(var_ptr == NULL)) {
		zend_error_noreturn(E_ERROR, "Cannot increment/decrement overloaded objects nor string offsets");
	}
	if (IS_VAR == IS_VAR && UNEXPECTED(*var_ptr == &EG(error_zval))) {
		ZVAL_NULL(&EX_T(opline->result.var).tmp_var);
		if (free_op1.var) {zval_ptr_dtor_nogc(&free_op1.var);};
		CHECK_EXCEPTION();
		ZEND_VM_NEXT_OPCODE();
	}
 
	retval = &EX_T(opline->result.var).tmp_var;
	ZVAL_COPY_VALUE(retval, *var_ptr);
	zendi_zval_copy_ctor(*retval);
当返回的是tmp_var变量时，说明该变量只是个临时值不会做额外的操作，所以也不需要增加gc_recount++等操作，直接堆栈释放即可。
++a的opcode：opline->result_type = IS_VAR; IS_VAR类型

	opline = get_next_op(CG(active_op_array) TSRMLS_CC);
	opline->opcode = op;
	SET_NODE(opline->op1, op1);
	SET_UNUSED(opline->op2);
	opline->result_type = IS_VAR;
	opline->result.var = get_temporary_variable(CG(active_op_array));
	GET_NODE(result, opline->result);


解析后：
	if (RETURN_VALUE_USED(opline)) {
		PZVAL_LOCK(*var_ptr);
		EX_T(opline->result.var).var.ptr = *var_ptr;
	}

此时返回的var的临时变量增加了gc_recount++，所以该返回值被别人用的时候，就需要有个释放的过程，
static zend_always_inline zval **_get_zval_ptr_ptr_var(zend_uint var, const zend_execute_data *execute_data, zend_free_op *should_free TSRMLS_DC)
{
	zval** ptr_ptr = EX_T(var).var.ptr_ptr;
 
	if (EXPECTED(ptr_ptr != NULL)) {
		PZVAL_UNLOCK(*ptr_ptr, should_free);
	} else {
		/* string offset */
		PZVAL_UNLOCK(EX_T(var).str_offset.str, should_free);
	}
	return ptr_ptr;
}
获取该变量值通过该函数，所以就会进行释放gc_recount--.
var和tmp类型的区别是什么，大家都是放在tmp分配的堆栈中，区别就是，var.ptr_ptr是个指针，为了节省空间大家目前先共用，比如++a，返回的临时变量和变量a返回值一样，所以就用指针指向同一个zval，比如a++，返回值的临时变量和变量a返回值不一样，所以必须重新申请一个zval，所以就干脆扔到了tmp中，可以随时释放。

https://blog.csdn.net/xiaolei1982/article/details/52140544
https://blog.csdn.net/xiaolei1982/article/details/20584291
http://www.phppan.com/2012/02/php-execute-data/
{% endraw %}
https://type.so/c/php-extension-in-action-get-arguments-after-zend-execute-ex.html