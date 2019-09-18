---
title: 函数调用帧栈
layout: post
category: linux
author: 夏泽民
---
区域	作用
栈区（stack）	由编译器自动分配和释放，存放函数的参数值，局部变量的值等。操作方式类似与数据结构中的栈
堆区（heap）	一般由程序员分配和释放，若程序员不释放，程序结束时可能由操作系统回收。与数据结构中的堆是两码事，分配方式类似于链表
静态区（static）	全局变量和静态变量存放于此
文字常量区	常量字符串放在此，程序结束后由系统释放
程序代码区	存放函数体的二进制代码
<!-- more -->
栈帧就是一个函数执行的环境。实际上，栈帧可以简单理解为：栈帧就是存储在用户栈上的（当然内核栈同样适用）每一次函数调用涉及的相关信息的记录单元。
栈是从高地址向低地址延伸的。每个函数的每次调用，都有它自己独立的一个栈帧，这个栈帧中维持着所需要的各种信息。寄存器ebp指向当前的栈帧的底部（高地址），寄存器esp指向当前的栈帧的顶部（地址地）。


寄存器名称	作用
eax	累加(Accumulator)寄存器，常用于函数返回值
ebx	基址(Base)寄存器，以它为基址访问内存
ecx	计数器(Counter)寄存器，常用作字符串和循环操作中的计数器
edx	数据(Data)寄存器，常用于乘除法和I/O指针
esi	源变址寄存器
dsi	目的变址寄存器
esp	堆栈(Stack)指针寄存器，指向堆栈顶部
ebp	基址指针寄存器，指向当前堆栈底部
eip	指令寄存器，指向下一条指令的地址

一、数据传送指令
1、传送指令：MOV (move)

格式：mov dst,src
具体用法：

(1) CPU内部寄存器之间的数据传送，如：mov ah,al
 
(2) 立即数送至通用寄存器(非段寄存器)或存储单元，如:mov al,3        mov [bx],1234h
 
(3) 寄存器与存储器间的数据传送，如：mov ax,var        mov ax,[bx]

2、交换指令：XCHG

xchg OPRD1,OPRD2    ;OPRD可以是通用寄存器或存储单元，但不包括段寄存器，不能同时是存储单元，不能有立即数
3、地址传送指令：LEA、LDS、LES

(1) LEA(Load Effective Address)
    格式：   lea REG,OPRD        
    功能：   把操作数OPRD的有效地址传送到操作数REG
     注：    REG必须是16位通用寄存器，OPRD必须是一个存储器操作数
 
     如：    lea ax,buf            ;buf是变量名
            lea ax,[si+2]
(2) LDS(Load pointer into DS)
    格式：   lds REG,OPRD        
    功能：   传送32位地址指针，将OPRD存储的32位数的高16位（段地址）送至DS，低16位（偏移地址）送至REG。（注意OPRD存放的32位数据，不是OPRD本身的地址）
     注：    操作数OPRD必须是一个32位存储器操作数，操作数REG可以时16位通用寄存器，但通常是指令指针寄存器（IP）或变址寄存器（SI，DI，SP，BP）
 
     如：       lds di,[bx]
            lds si,FARPOINTER      ;FARPOINTER是一个32位（双字）变量
(3) LES(Load pointer into ES)
    格式：      les REG,OPRD
    功能：      把操作数OPRD存储的32位数据的高16位（段地址）送至ES，低16位（偏移地址）送至REG
 
    其他同LDS
二、堆栈操作指令
1、进栈指令：push

格式：push src
功能： 把16位数据src压入堆栈。
注：    源操作数src可以是通用寄存器和段寄存器，也可以是字存储单元
 
如：  push si
     push [si]
     push var        ;var是16位（字）变量
2、出栈指令：pop

格式：pop dst
功能：从堆栈弹出16位数据至dst
注： dst可以是通用寄存器和段寄存器，但不能是CS，可以是字存储单元
 
如： pop si
    pop [si]
    pop var            ;var是字变量
三、标志操作指令
1、标志传送指令：LAHF(Load AH with Flags)、SAHF(Store AH into Flags)、PUSHF、POPF

(1) LAHF(Load AH with Flags)
    格式：LAHF
    功能：把标志寄存器的低8位（包括SF(7)、ZF(6)、AF(4)、PF(2)、CF(0)）传送到AH指定位。
 
(2) SAHF(Store AH into Flags)
    格式：SAHF
    功能：把寄存器AH的指定位送至标志寄存器低8位（包括SF(7)、ZF(6)、AF(4)、PF(2)、CF(0)）。
 
(3) PUSHF
    格式：PUSHF
    功能：把标志寄存器的内容（16位）压入堆栈。SP-=2
    注： 这条指令不影响标志位
 
(4) POPF
    格式：POPF
    功能：把当前栈顶的一个字传送到标志寄存器。SP+=2
2、标志位操作指令：CLC、STC、CMC、CLD、STD、CLI、STI

(1) CLC(Clear Carry Flag):                CF置0
(2) STC(Set Carry Flag):                CF置1
 
(3) CMC(Complement Carry Flag):            CF取反
 
(4) CLD(Clear Direction Flag):            DF置0，执行串操作指令时，地址递增
(5) STD(Set Direction Flag):            DF置1，执行串操作指令时，地址递减
 
(6) CLI(Clear Interrupt enable Flag)    IF置0，使CPU不响应来自外部装置的可屏蔽中断，但对不可屏蔽中断和内部中断没有影响
(7) STI(Set Interrupt enable Flag)        IF置1，可以响应可屏蔽中断
四、加减运算指令
加减法运算对无符号数和有符号数的处理一视同仁。即作为无符号数而影响标志位CF和AF，也作为有符号数影响标识OF和SF，且总会影响ZF。加减法运算指令也要影响标志位PF。有些指令稍有例外。
存放运算结果的操作数（两个操作数时即左操作数）只能是通用寄存器或存储单元（变量）。如果参与运算的操作数有两个，则最多只能有一个是存储器操作数。
如果参与运算的操作数有两个，则它们的类型必须一致，如同时为字节或同时为字等。
1、加法指令：add、adc、inc

(1) add(Addtion)
    格式：add OPRD1,OPRD2
    功能：OPRD1 = OPRD1 + OPRD2
    注：  影响FLAG
 
    如：add al,5
           add bl,var        ;var是字节变量
           add var,si        ;var是字变量
 
(2) adc(add with Carry)        ;带进位的加法
    格式：adc OPRD1,OPRD2
    功能：带进位的加法，OPRD1 = OPRD1 + OPRD2 + CF
    注：    影响FLAG，主要用于多字节运算
 
    如：    adc al,[bx]
         adc dx,ax
         adc dx,var        ;var是字变量
 
(3) inc(Increment)
    格式：inc OPRD
    功能：OPRD = OPRD + 1
    注： 不影响CF
 
    如：inc al
        inc var            ;var是字节变量，也可以是字变量
        inc cx
2、减法指令：sub、sbb、dec、neg、cmp

(1) sub(Subtraction)
    格式：sub OPRD1，OPRD2
    功能：OPRD1 = OPRD1 - OPRD2
 
    如：  sub ah,12
         sub bx,bp
         sub al,[bx]
         sub [BP],AX
         sub AX,VAR        ;VAR是字变量
 
(2) sbb(Sub with Borrow)
    格式：sbb OPRD1,OPRD2
    功能：OPRD1 = OPRD1 - OPRD2 - CF
    注：    主要用于多字节数相减的情况
 
(3) dec(decrement)
    格式：dec OPRD
    功能：OPRD = OPRD - 1
    注：    操作数OPRD可以是通用寄存器，也可以是存储单元。相减时把操作数作为一个无符号数对待，这条指令影响ZF、SP、OF、PF、AF，但不影响CF，该指令主要用于调整地址指针和计数器。
 
(4) neg(Negate)
    格式：NEG OPRD
    功能：对操作数取补，即OPRD = 0 - OPRD
    注： 操作数可以是通用寄存器，也可以是存储单元。此指令结果影响CF、ZF、OF、AF、PF，一般会使CF为1，除非OPRD=0
 
(5)    cmp(Compare)
    格式：cmp OPRD1,OPRD2
    功能：执行OPRD1 - OPRD2，但运算结果不运送到OPRD1
    注：    该指令通过OPRD - OPRD2影响标志位CF、ZF、SF、OF、AF、PF来判断OPRD1和OPRD2的大小关系。通过ZF判断是否相等；如果是无符号数，通过CF可判断大小；如果是有符号数，通过SF和OF判断大小
五、乘除运算指令
1、乘法指令：mul、imul

(1) mul(Multiply)            ;无符号数乘法指令
    格式：MUL OPRD
    功能：将OPRD与AX或AL中的操作数相乘，结果保存在DX:AX中或AX中
    注： 无符号数相乘分为16位*16位和8位*8位，结果分别为32位和16位，保存在DX:AX中或AX中，其中结果为32位时，DX为高16位，AX为低16位；结果为16位时，AH为高8位，AL为低8位。
 
(2) imul(Signed Multiply)    ;有符号数乘法指令
    格式：IMUL OPRD
    功能：把乘数和被乘数均作为有符号数进行乘法运算。其余与mul类似
    注：    如果乘积结果的高位部分（DX或AH）不是低位的符号扩展，则CF=1，OF=1，否则CF=0，OF=0。即CF=1，OF=1表示AH或DX中含有结果的有效数。
         如果除数为0，或8位数除时商超过8位，16位数除时商超过16位，则认为是除溢出，引起0号中断。除法指令对标志位的影响无定义。
2、除法指令：div、idiv

(1) div(Division)            ;无符号数除法指令
    格式：DIV OPRD
    功能：OPRD为除数，被除数存放在DX:AX或AX中，做除法，结果存放在DX:AX（DX存放余数，AX存放商）或AX（AH余数，AL商）。
    注：    8086中除法有32位除以16位和16位除以8位。前者被除数为32位，高位在DX中，低位在AX中，除数OPRD为16位通用寄存器或16位存储器操作数，结果为16位，其中16位余数存放在DX中，16位商存放在AX中；若为16位除以8位，被除数存放在AX中，OPRD为8位通用寄存器或存储器操作数，结果8位余数存放在AH中，8位商存放在AL中。
 
(2) idiv(Signed Division)    ;有符号数除法指令
    格式：IDIV OPRD
    功能：把除数和被除数看做有符号数做除法，其余与div类似
3、符号扩展指令：cbw、cwd

(1) cbw(Convert Byte to Word)
    格式：CBW
    功能：把寄存器AL中的符号位扩展到寄存器AH
 
(2) cwd(Convert Word to Double Word)
    格式：CWD
    功能：把寄存器AX中的符号扩展到寄存器DX
六、逻辑运算和移位指令
1、逻辑运算指令：not、and、or、xor、test

(1) NOT
    格式：NOT OPRD
    功能：把操作数OPRD取反，然后送回OPRD。
    注：    OPRD可以是通用寄存器，也可以是存储器操作数，此指令对标志没有影响
 
(2) AND
    格式：AND OPRD1,OPRD2
    功能：对两个操作数进行按位逻辑“与”运算，结果送到OPRD1中
    注： 该指令执行后，CF=0，OF=0，标志PF、ZF、SF反映运算结果，AF未定义。
         某个操作数与自身相与，值不变，但可以使CF置0。
 
(3) OR
    格式：OR OPRD1,OPRD2
    功能：对两个操作数进行按位逻辑“或”运算，结果送到OPRD1中
    注： 该指令执行后，CF=0，OF=0，标志PF、ZF、SF反映运算结果，AF未定义。
         某个操作数与自身相或，值不变，但可以使CF置0。
 
(4) XOR
    格式：XOR OPRD1,OPRD2
    功能：对两个操作数进行按位逻辑“异或”运算，结果送到OPRD1中
    注： 该指令执行后，CF=0，OF=0，标志PF、ZF、SF反映运算结果，AF未定义。
 
(5) TEST
    格式：TEST OPRD1,OPRD2
    功能：把OPRD1与OPRD2按位“与”，但结果不送到OPRD1中，仅影响标志位。
    注： 该指令执行后，CF=0，OF=0，标志PF、ZF、SF反映运算结果。常用于检测某些位是否为1
2、一般移位指令：SAL/SHL,SAR/SHR

(1) SAL/SHL(Shift Arithmetic Left / Shift Logic Left)        ;算术左移/逻辑左移
    格式：SAL OPRD,m
         SHL OPRD,m
    功能：把操作数OPRD左移m位，每移动一位，右边用0补足1位，移出的最高位进入标志位CF
    注：    算术左移和逻辑左移进行相同的动作，为了方便提供了两个助记符。
 
(2) SAR(Shift Arithmetic Right)                                ;算数右移指令
    格式：SAR OPRD,m
    功能：操作数右移m位，同时每移1位，左边的符号位保持不变，移出的最低位进入标志位CF
    注：    对有符号数和无符号数，算数右移1位相当于除以2
 
(3) SHR(Shift Logic Right)                                    ;逻辑右移指令
    格式：SHR OPRD,m
    功能：操作数右移m位，同时每移1位，左边用0补足，移出的最低位进入标志位CF
    注：    对无符号数，逻辑右移1位相当于除以2
3、循环移位指令：ROL、ROR、RCL、RCR

格式：

ROL OPRD,m


ROR OPRD,m


RCL OPRD,m


RCR OPRD,m


这些指令只影响CF和OF
七、转移指令
1、无条件转移指令：JMP(Jump)

段内转移：改变IP
段间转移：改变CS:IP

(1) 无条件段内直接转移指令
    格式：JMP 标号
    功能：使控制无条件转移至标号地址处
    原理：把编译时计算出的地址差加到IP上
 
(2) 无条件段内间接转移指令
    格式：JMP OPRD
    功能：使控制指令无条件转移到OPRD的内容给定的目标地址处。操作数OPRD可以是通用寄存器，也可以是字存储单元
    例如：jmp cx                    ;CX寄存器的内容送IP
         jmp word ptr [1234h]    ;字存储单元[1234h]的内容送IP
 
(3) 无条件段间直接转移指令
    格式：jmp far ptr 标号
    功能：使控制指令无条件的转移到标号对应的地址处
    原理：把编译时产生的标号处的段地址和偏移地址分别置入CS和IP
 
(4) 无条件段间间接转移指令
    格式：JMP OPRD
    功能：使控制指令无条件转移到操作数OPRD的内容给定的目标地址处。操作数OPRD必须是双字存储单元
    例如：jmp dword ptr [1234h]    ;双字存储单元的低字内容送IP，高字内容送CS
2、条件转移指令



3、循环指令：LOOP、LOOPE/LOOPZ、LOOPNE/LOOPNZ、JCXZ

(1) Loop            ;计数循环指令
    格式：loop 标号
    功能：使转移标号与Loop指令间的指令循环执行CX次
    原理：指令执行至loop时，cx减1，如果cx不为0，则跳转至标号处，否则继续执行下一条指令
        即：DEC CX
           JNZ 标号
 
(2) LOOPE/LOOPZ        ;等于/全零循环指令
    格式：LOOPE 标号
         LOOPZ 标号
    功能：该指令使CX自减1，若结果不为0，并且ZF=1，则转移至标号，否则顺序执行。注意指令本身实施的CX自减1操作不影响标志
 
(3) LOOPNE/LOOPNZ    ;不等于/非零循环指令
    格式：LOOPNE 标号
         LOOPNZ 标号
    功能：该指令使CX自减1，若结果不为0，并且ZF=0，则转移至标号，否则顺序执行。注意指令本身实施的CX自减1操作不影响标志
 
(4) JCXZ            ;跳转指令
    格式：JCXZ 标号
    功能：当寄存器CX的值为0时跳转到标号，否则顺序执行
 
当发生函数调用的时候,栈空间中存放的数据是这样的:
1、调用者函数把被调函数所需要的参数按照与被调函数的形参顺序相反的顺序压入栈中,即:从右向左依次把被调函数所需要的参数压入栈;
2、调用者函数使用call指令调用被调函数,并把call指令的下一条指令的地址当成返回地址压入栈中(这个压栈操作隐含在call指令中);
3、在被调函数中,被调函数会先保存调用者函数的栈底地址(push ebp),然后再保存调用者函数的栈顶地址,即:当前被调函数的栈底地址(mov ebp,esp);
4、在被调函数中,从ebp的位置处开始存放被调函数中的局部变量和临时变量,并且这些变量的地址按照定义时的顺序依次减小,即:这些变量的地址是按照栈的延伸方向排列的,先定义的变量先入栈,后定义的变量后入栈;

当进程被加载到内存时，会被分成很多段

代码段：保存程序文本，指令指针EIP就是指向代码段，可读可执行不可写，如果发生写操作则会提示segmentation fault
数据段：保存初始化的全局变量和静态变量，可读可写不可执行
BSS：未初始化的全局变量和静态变量
堆(Heap)：动态分配内存，向地址增大的方向增长，可读可写可执行
栈(Stack)：存放局部变量，函数参数，当前状态，函数调用信息等，向地址减小的方向增长，可读可写可执行
环境/参数段（environment/argumentssection）：用来存储系统环境变量的一份复制文件，进程在运行时可能需要。例如，运行中的进程，可以通过环境变量来访问路径、shell 名称、主机名等信息。该节是可写的，因此在缓冲区溢出（buffer overflow）攻击中都可以使用该段
 

寄存器

EAX：累加(Accumulator)寄存器，常用于函数返回值

EBX：基址(Base)寄存器，以它为基址访问内存

ECX：计数器(Counter)寄存器，常用作字符串和循环操作中的计数器

EDX：数据(Data)寄存器，常用于乘除法和I/O指针

ESI：源变址寄存器

DSI：目的变址寄存器

ESP：堆栈(Stack)指针寄存器，指向堆栈顶部

EBP：基址指针寄存器，指向当前堆栈底部

EIP：指令寄存器，指向下一条指令的地址

 

入栈push和出栈pop

push ebp就等于将ebp的值保存到栈中，并且将当前esp下移

pop ebp就等于将ebp的值从栈中取出来，将ebp指向这个值
