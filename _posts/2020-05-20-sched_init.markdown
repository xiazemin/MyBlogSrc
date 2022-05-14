---
title: sched_init
layout: post
category: linux
author: 夏泽民
---
这个很熟是进程调度初始化函数，主要做了设置进程的GDT，LDT描述符，设置系统定时器中断，系统调用终端，代码如下：

void sched_init(void)
{
        int i;
        struct desc_struct * p;
/**********************************
描述结构在include/linux/head.h中
typedef struct desc_struct {
        unsigned long a,b;
} 
**************************************************/
        if (sizeof(struct sigaction) != 16)
                panic("Struct sigaction MUST be 16 bytes");

    set_tss_desc(gdt+FIRST_TSS_ENTRY,&(init_task.task.tss));
/***********************************************************
设置GDT，FIRST_TSS_ENTRY，对应任务0的GDT。设置任务状态
**************************************************************/
        set_ldt_desc(gdt+FIRST_LDT_ENTRY,&(init_task.task.ldt));
/**********************************************
设置LDT，FIRST_TSS_ENTRY，对应任务0的LDT。设置任务0的LDT
体调用：设置了对应ldt，gdt的位
#define _set_tssldt_desc(n,addr,type) \
__asm__ ("movw $104,%1\n\t" \
        "movw %%ax,%2\n\t" \
        "rorl $16,%%eax\n\t" \
        "movb %%al,%3\n\t" \
        "movb $" type ",%4\n\t" \
        "movb $0x00,%5\n\t" \
        "movb %%ah,%6\n\t" \
        "rorl $16,%%eax" \
        ::"a" (addr), "m" (*(n)), "m" (*(n+2)), "m" (*(n+4)), \
         "m" (*(n+5)), "m" (*(n+6)), "m" (*(n+7)) \
        )
        
#define set_tss_desc(n,addr) _set_tssldt_desc(((char *) (n)),addr,"0x89")
#define set_ldt_desc(n,addr) _set_tssldt_desc(((char *) (n)),addr,"0x82")
?*************************************************************/

        p = gdt+2+FIRST_TSS_ENTRY;//移动到任务1 gdt
        for(i=1;i                 task[i] = NULL;
                p->a=p->b=0;//状态段清零或者GDT置空
                p++;
                p->a=p->b=0;//ldt置空置空哦
                p++;
        }
//上面这一段，初始货63个进程，全部置空，NR_TASKS为64
/* Clear NT, so that we won't have troubles with that later on */
        __asm__("pushfl ; andl $0xffffbfff,(%esp) ; popfl");
        ltr(0);
/***********************************************************
这个调用#define ltr(n) __asm__("ltr %%ax"::"a" (_TSS(n)))
#define _TSS(n) ((((unsigned long) n)<<4)+(FIRST_TSS_ENTRY<<3))
加载任务寄存器。把任务0的tss段选择述符和段描述符加载到任务寄存器。
**********************************************************/
        lldt(0);
/***************************************************#
这个调用：define lldt(n) __asm__("lldt %%ax"::"a" (_LDT(n)))
_LDT（n）:#define _LDT(n) ((((unsigned long) n)<<4)+(FIRST_LDT_ENTRY<<3))
主要作用加载描述符表寄存器LDTR。
*************************************************************/
        outb_p(0x36,0x43);              /* binary, mode 3, LSB/MSB, ch 0 */
        outb_p(LATCH & 0xff , 0x40);    /* LSB */
        outb(LATCH >> 8 , 0x40);        /* MSB */
        set_intr_gate(0x20,&timer_interrupt);
/************************************
这四句是初始化系统定时器中断，cat   /proc/ioports有   0040-0043 : timer0
说明第一句是设置只用模式3，接下来两句设置计数的高位和低位，#define LATCH (1193180/HZ)
最后一句定义中断服务程序timer_interrupt会执行jmp ret_from_sys_call，检查是否要切换任务
**********************************/
        outb(inb_p(0x21)&~0x01,0x21);
        set_system_gate(0x80,&system_call);//定义系统调用
}

总体来说，这段程序先初始话0进程，包括段选择符，描述符GDT。LDT等。然后将其他63个进程的的 段选择符，描述符GDT。LDT
置空，设置好后好后，将任务0的 段选择符，描述符GDT。LDT等加载进个寄存器中。接着设置系统中断定时器，中断函数判断
是
<!-- more -->
https://www.cnblogs.com/sky-heaven/p/11506521.html

https://www.xuebuyuan.com/3263090.html

https://blog.csdn.net/xiongtiancheng/article/details/78880025
