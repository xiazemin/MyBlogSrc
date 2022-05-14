---
title: sigaction
layout: post
category: linux
author: 夏泽民
---
sigaction函数的功能是检查或修改与指定信号相关联的处理动作（可同时两种操作）
其中，参数signo是要检测或修改其具体动作的信号编号。若act指针非空，则要修改其动作。如果oact指针非空，则系统经由oact指针返回该信号的上一个动作。此函数使用下列结构

他是POSIX的信号接口，而signal()是标准C的信号接口(如果程序必须在非POSIX系统上运行，那么就应该使用这个接口)

给信号signum设置新的信号处理函数act， 同时保留该信号原有的信号处理函数oldact

int sigaction(int signo,conststruct sigaction*restrict act,struct sigaction*restrict oact);

struct sigaction {
    void (*sa_handler)(int);
    void (*sa_sigaction)(int, siginfo_t *, void *);
    sigset_t sa_mask;
    int sa_flags;
    void (*sa_restorer)(void);
}
sa_handler此参数和signal()的参数handler相同，代表新的信号处理函数
sa_mask 用来设置在处理该信号时暂时将sa_mask 指定的信号集搁置
sa_flags 用来设置信号处理的其他相关操作，下列的数值可用。 
SA_RESETHAND：当调用信号处理函数时，将信号的处理函数重置为缺省值SIG_DFL
SA_RESTART：如果信号中断了进程的某个系统调用，则系统自动启动该系统调用
SA_NODEFER ：一般情况下， 当信号处理函数运行时，内核将阻塞该给定信号。但是如果设置了 SA_NODEFER标记， 那么在该信号处理函数运行时，内核将不会阻塞该信号
<!-- more -->
结构sigaction定义如下：

structsigaction{
  void (*sa_handler)(int);
   sigset_t sa_mask;
  int sa_flag;
  void (*sa_sigaction)(int,siginfo_t*,void*);
};

sa_handler字段包含一个信号捕捉函数的地址

sa_mask字段说明了一个信号集，在调用该信号捕捉函数之前，这一信号集要加进进程的信号屏蔽字中。仅当从信号捕捉函数返回时再将进程的信号屏蔽字复位为原先值。

sa_flag是一个选项，主要理解两个

SA_NODEFER:  当信号处理函数正在进行时，不堵塞对于信号处理函数自身信号功能。
SA_RESETHAND:当用户注册的信号处理函数被执行过一次后，该信号的处理函数被设为系统默认的处理函数。

SA_SIGINFO 提供附加信息，一个指向siginfo结构的指针以及一个指向进程上下文标识符的指针

最后一个参数是一个替代的信号处理程序，当设置SA_SIGINFO时才会用他。


#include <stdio.h>  
#include <signal.h>   
void WrkProcess(int nsig)  
{  
        printf("WrkProcess .I get signal.%d threadid:%d/n",nsig,pthread_self());  
  
  
        int i=0;  
        while(i<5){  
                printf("%d/n",i);  
                sleep(1);  
                i++;  
        }  
}  
  
int main()  
{  
        struct sigaction act,oldact;  
        act.sa_handler  = WrkProcess;  
//      sigaddset(&act.sa_mask,SIGQUIT);  
//      sigaddset(&act.sa_mask,SIGTERM)  
        act.sa_flags = SA_NODEFER | SA_RESETHAND;    
//        act.sa_flags = 0;  
  
        sigaction(SIGINT,&act,&oldact);  
  
        printf("main threadid:%d/n",pthread_self());  
  
        while(1)sleep(5);  
  
        return 0;  
} 

1）执行改程序时，ctrl+c，第一次不会导致程序的结束。而是继续执行，当用户再次执行ctrl+c的时候，程序采用结束。
2）如果对程序稍微进行一下改动，则会出现另外一种情况。
改动为：act.sa_flags = SA_NODEFER；
经过这种改变之后，无论对ctrl+d操作多少次，程序都不会结束。
3）下面如果再对程序进行一次改动，则会出现第三种情况。
For example:  act.sa_flags = 0;

在执行信号处理函数这段期间，多次操作ctrl+c，程序也不会调用信号处理函数，而是在本次信号处理函数完成之后，在执行一次信号处理函数（无论前面产生了多少次ctrl+c信号）。
如果在2）执行信号处理函数的过程中，再次给予ctrl+c信号的时候，会导致再次调用信号处理函数。
4）如果在程序中设置了sigaddset(&act.sa_mask,SIGQUIT);程序在执行信号处理函数的过程中，发送ctrl+/信号，程序也不会已经退出，而是在信号处理函数执行完毕之后才会执行SIGQUIT的信号处理函数，然后程序退出。如果不添加这项设置，则程序将会在接收到ctrl+/信号后马上执行退出，无论是否在ctrl+c的信号处理函数过程中。

原因如下：
1）情况下，第一次产生ctrl+c信号的时候，该信号被自己设定的信号处理函数进行了处理。在处理过程中，由于我们设定了SA_RESETHAND标志位，又将该信号的处理函数设置为默认的信号处理函数（系统默认的处理方式为IGN）,所以在第二次发送ctrl+d信号的时候，是由默认的信号处理函数处理的，导致程序结束；

2）情况下，我们去掉了SA_RESETHAND了标志位，导致程序中所有的ctrl+d信号均是由我们自己的信号处理函数来进行了处理，所以我们发送多少次ctrl+c信号程序都不会退出；

3）情况下，我们去掉了SA_NODEFER标志位。程序在执行信号处理函数过程中，ctrl+c信号将会被阻止，但是在执行信号处理函数期发送的ctrl+c信号将会被阻塞，知道信号处理函数执行完成，才有机会处理信号函数执行期间产生的ctrl+c,但是在信号函数执行产生的多次ctrl+c，最后只会产生ctrl+c。2）情况下，由于设置了SA_NODEF，ctrl+c信号将不会被阻塞。所以能够并行执行下次的信号处理函数。
4）情况下，我们是设置了在执行信号处理函数过程中，我们将屏蔽该信号，当屏蔽该信号的处理函数执行完毕后才会进行处理该信号。

附：
当我们按下ctrl+c的时候，操作为：向系统发送SIGINT信号，SIGINT信号的默认处理，退出程序。
当我们按下ctrl+/的时候，操作为：向系统发送SIGQUIT信号，该信号的默认处理为退出程序。