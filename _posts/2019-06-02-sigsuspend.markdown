---
title: sigsuspend 进程阻塞 与 pause 区别
layout: post
category: linux
author: 夏泽民
---
sigsuspend函数 ：
sigsuspend函数接受一个信号集指针，将信号屏蔽字设置为信号集中的值，在进程接受到一个信号之前，进程会挂起，当捕捉一个信

号，首先执行信号处理程序，然后从sigsuspend返回，最后将信号屏蔽字恢复为调用sigsuspend之前的值。

#include <signal.h>
int sigsuspend(const sigset_t *sigmask);
                                                            // 返回值：－1，并将errno设置为EINTR

pause函数：

pause函数使调用进程挂起直到捕捉到一个信号。只有执行了一个信号处理程序并从其返回时，pause才返回



int sigsuspend(const sigset_t *sigmask);

此函数用于进程的挂起，sigmask指向一个信号集。当此函数被调用时，sigmask所指向的信号集中的信号将赋值给信号掩码。之后进程挂起。直到进程捕捉到信号，并调用处理函数返回时，函数sigsuspend返回。信号掩码恢复为信号调用前的值，同时将errno设为EINTR。进程结束信号可将其立即停止。
<!-- more -->
#include <stdio.h>   
#include <signal.h>   
  
void checkset();  
void func();  
void main()  
{  
     sigset_tblockset,oldblockset,zeroset,pendmask;  
     printf("pid:%ld\n",(long)getpid());  
     signal(SIGINT,func);  
  
     sigemptyset(&blockset);  
     sigemptyset(&zeroset);  
     sigaddset(&blockset,SIGINT);  
  
     sigprocmask(SIG_SETMASK,&blockset,&oldblockset);  
     checkset();  
     sigpending(&pendmask);  
  
     if(sigismember(&pendmask,SIGINT))  
         printf("SIGINTpending\n");  
  
     if(sigsuspend(&zeroset)!= -1)  
     {  
     printf("sigsuspenderror\n");  
     exit(0);  
     }  
  
     printf("afterreturn\n");  
     sigprocmask(SIG_SETMASK,&oldblockset,NULL);  
  
     printf("SIGINTunblocked\n");  
}  
  
void checkset()  
{    sigset_tset;  
     printf("checksetstart:\n");  
     if(sigprocmask(0,NULL,&set)<0)  
     {  
     printf("checksetsigprocmask error!!\n");  
     exit(0);  
     }  
  
     if(sigismember(&set,SIGINT))  
     printf("sigint\n");  
  
     if(sigismember(&set,SIGTSTP))  
     printf("sigtstp\n");  
  
     if(sigismember(&set,SIGTERM))  
     printf("sigterm\n");  
     printf("checksetend\n");  
}  
  
void func()  
{  
     printf("hellofunc\n");  
}  

sigsuspend(sigset_t  sigs);功能： 屏蔽新的信号，原来屏蔽的信号失效。sigsuspend是阻塞函数，对参数信号屏蔽，对参数没有指定的信号不屏蔽，但当没有屏蔽的信号处理函数调用完毕sigsuspend函数返回。

sigsuspend返回条件：

信号发生，并且信号是非屏蔽信号
信号必须要处理，而且处理函数返回后sigsuspend才返回。
sigsuspend设置新的屏蔽信号，保存旧的屏蔽信号，而且当sigsuspend返回的时候，恢复旧的屏蔽信号。

其实可以这样理解：sigsuspend=pause()+指定屏蔽的信号

sigsuspend的整个原子操作过程为：
(1) 设置新的mask阻塞当前进程；
(2) 收到信号，调用该进程设置的信号处理函数；
(3) 待信号处理函数返回后，恢复原先mask；
(4) sigsuspend返回

1. sigpromask(SIG_UNBLOCK,&newmask,&oldmask)和
      sigpromask(SIG_SETMASK,&oldmask,NULL)区别

sigpromask(SIG_UNBLOCK,&newmask,&oldmask)
    它的作用有两个:一是设置新的信号掩码(不阻塞newmask所指的信号集).二是保存原来的信号掩码(放在oldmask所指的信号集中)
sigpromask(SIG_SETMASK,&oldmask,NULL)
    它的作用只有一个:设置新的信号掩码(信号掩码为oldmask所指的信号集)

2. sigsuspend 用实参 sigmask 指定的信号集代替调用进程的信号屏蔽， 然后挂起该进程直到某个不属于 sigmask 成员的信号到

达为止。此信号的动作要么是执行信号句柄，要么是终止该进程。
如果信号终止进程，则 suspend 函数不返回。如果信号的动作是执行信号句柄，则在信号句柄返回后，sigsuspend 函数返回，并使

进程的信号屏蔽恢复到 sigsuspend 调用之前的值。

3. 清晰且可靠的等待信号到达的方法是先阻塞该信号（防止临界区重入，也就是在次期间有另外一个该信号到达），然后使用

sigsuspend 放开此信号并等待句柄设置信号到达标志。如下所示, 等待 SIGUSR1 信号到来：

sigemptyset(&zeromask);
sigaddset(&newmask, SIGUSR1);
......

sigprocmask(SIG_BLOCK, &newmask, NULL);
while(flag)
      sigsuspend(&zeromask);
flag = 0;
......
sigprocmask(SIG_UNBLOCK, &newmask, NULL);

如果在等待信号发生时希望去休眠，则使用sigsuspend函数是非常合适的，但是如果在等待信号期间希望调用其他系统函数，那么将会怎样呢？不幸的是，在单线程环境下对此问题没有妥善的解决方法。如果可以使用多线程，则可专门安排一个线程处理信号。

如果不使用线程，那么我们能尽力做到最好的是，当信号发生时，在信号捕捉程序中对一个全局变量置1.




