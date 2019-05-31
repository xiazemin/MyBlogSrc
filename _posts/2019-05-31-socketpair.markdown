---
title: socketpair  popen
layout: post
category: linux
author: 夏泽民
---
先说说我的理解：socketpair创建了一对无名的套接字描述符（只能在AF_UNIX域中使用），描述符存储于一个二元数组,eg. s[2] .这对套接字可以进行双工通信，每一个描述符既可以读也可以写。这个在同一个进程中也可以进行通信，向s[0]中写入，就可以从s[1]中读取（只能从s[1]中读取），也可以在s[1]中写入，然后从s[0]中读取；但是，若没有在0端写入，而从1端读取，则1端的读取操作会阻塞，即使在1端写入，也不能从1读取，仍然阻塞；反之亦然......
<!-- more -->
//
// Created by didi on 2019-05-31.
//


#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
//#include <error.h>
#include <errno.h>
#include <sys/socket.h>
#include <stdlib.h>

#define BUF_SIZE 30

int main(){
    int s[2];
    int w,r;
    char * string = "This is a test string";
    char * buf = (char*)calloc(1 , BUF_SIZE);

    if( socketpair(AF_UNIX,SOCK_STREAM,0,s) == -1 ){
        printf("create unnamed socket pair failed:%s\n",strerror(errno) );
        exit(-1);
    }

    /*******test in a single process ********/
    if( ( w = write(s[0] , string , strlen(string) ) ) == -1 ){
        printf("Write socket error:%s\n",strerror(errno));
        exit(-1);
    }
    /*****read*******/
    if( (r = read(s[1], buf , BUF_SIZE )) == -1){
        printf("Read from socket error:%s\n",strerror(errno) );
        exit(-1);
    }
    printf("Read string in same process : %s \n",buf);
    if( (r = read(s[0], buf , BUF_SIZE )) == -1){
        printf("Read from socket s0 error:%s\n",strerror(errno) );
        exit(-1);
    }
    printf("Read from s0 :%s\n",buf);

    printf("Test successed\n");
    exit(0);
}

若fork子进程，然后在服进程关闭一个描述符eg. s[1] ，在子进程中再关闭另一个 eg. s[0]    ,则可以实现父子进程之间的双工通信，两端都可读可写；当然，仍然遵守和在同一个进程之间工作的原则，一端写，在另一端读取；

     这和pipe有一定的区别，pipe是单工通信，一端要么是读端要么是写端，而socketpair实现了双工套接字，也就没有所谓的读端和写端的区分

//
// Created by didi on 2019-05-31.
//

#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
//#include <error.h>
#include <errno.h>
#include <sys/socket.h>
#include <stdlib.h>

#define BUF_SIZE 30

int main(){
    int s[2];
    int w,r;
    char * string = "This is a test string";
    char * buf = (char*)calloc(1 , BUF_SIZE);
    pid_t pid;

    if( socketpair(AF_UNIX,SOCK_STREAM,0,s) == -1 ){
        printf("create unnamed socket pair failed:%s\n",strerror(errno) );
        exit(-1);
    }

    /***********Test : fork but don't close any fd in neither parent nor child process***********/
    if( ( pid = fork() ) > 0 ){
        printf("Parent process's pid is %d\n",getpid());
        close(s[1]);
        if( ( w = write(s[0] , string , strlen(string) ) ) == -1 ){
            printf("Write socket error:%s\n",strerror(errno));
            exit(-1);
        }
    }else if(pid == 0){
        printf("Fork child process successed\n");
        printf("Child process's pid is :%d\n",getpid());
        close(s[0]);
    }else{
        printf("Fork failed:%s\n",strerror(errno));
        exit(-1);
    }

    /*****read***In parent and child****/
    if( (r = read(s[1], buf , BUF_SIZE )) == -1){
        printf("Pid %d read from socket error:%s\n",getpid() , strerror(errno) );
        exit(-1);
    }
    printf("Pid %d read string in same process : %s \n",getpid(),buf);
    printf("Test successed , %d\n",getpid());
    exit(0);
}

我也测试了在父子进程中都不close(s[1])，也就是保持两个读端，则父进程能够读到string串，但子进程读取空串，或者子进程先读了数据，父进程阻塞于read操作！

之所以子进程能读取父进程的string，是因为fork时，子进程继承了父进程的文件描述符的，同时也就得到了一个和父进程指向相同文件表项的指针；若父子进程均不关闭读端，因为指向相同的文件表项，这两个进程就有了竞争关系，争相读取这个字符串．父进程read后将数据转到其应用缓冲区，而子进程就得不到了，只有一份数据拷贝（若将父进程阻塞一段时间，则收到数据的就是子进程了，已经得到验证，让父进程sleep(3)，子进程获得string，而父进程获取不到而是阻塞）


Linux 提供了 popen 和 pclose 函数 (1)，用于创建和关闭管道与另外一个进程进行通信。其接口如下：
FILE *popen(const char *command， const char *mode);
int pclose(FILE *stream);

有一种解决方案是使用 pipe 函数 (2)创建两个单向管道。没有错误检测的代码示意如下：
int pipe_in[2], pipe_out[2];
pid_t pid;
pipe(&pipe_in); // 创建父进程中用于读取数据的管道
pipe(&pipe_out);    // 创建父进程中用于写入数据的管道
if ( (pid = fork()) == 0) { // 子进程
    close(pipe_in[0]);  // 关闭父进程的读管道的子进程读端
    close(pipe_out[1]); // 关闭父进程的写管道的子进程写端
    dup2(pipe_in[1], STDOUT_FILENO);    // 复制父进程的读管道到子进程的标准输出
    dup2(pipe_out[0], STDIN_FILENO);    // 复制父进程的写管道到子进程的标准输入
    close(pipe_in[1]);  // 关闭已复制的读管道
    close(pipe_out[0]); // 关闭已复制的写管道
    /* 使用exec执行命令 */
} else {    // 父进程
    close(pipe_in[1]);  // 关闭读管道的写端
    close(pipe_out[0]); // 关闭写管道的读端
    /* 现在可向pipe_out[1]中写数据，并从pipe_in[0]中读结果 */
    close(pipe_out[1]); // 关闭写管道
    /* 读取pipe_in[0]中的剩余数据 */
    close(pipe_in[0]);  // 关闭读管道
    /* 使用wait系列函数等待子进程退出并取得退出代码 */
}
当然，这样的代码的可读性（特别是加上错误处理代码之后）比较差，也不容易封装成类似于popen/pclose的函数，方便高层代码使用。究其原因，是pipe函数返回的一对文件描述符只能从第一个中读、第二个中写（至少对于Linux是如此）。为了同时读写，就只能采取这么累赘的两个pipe调用、两个文件描述符的形式了。

Linux实现了一个源自BSD的socketpair调用 (3)，可以实现上述在同一个文件描述符中进行读写的功能（该调用目前也是POSIX规范的一部分 (4)）。该系统调用能创建一对已连接的（UNIX族）无名socket。在Linux中，完全可以把这一对socket当成pipe返回的文件描述符一样使用，唯一的区别就是这一对文件描述符中的任何一个都可读和可写。

int fd[2];
pid_t pid;
socketpair(AF_UNIX, SOCKET_STREAM, 0, fd);  // 创建管道
if ( (pid = fork()) == 0) { // 子进程
    close(fd[0]);   // 关闭管道的父进程端
    dup2(fd[1], STDOUT_FILENO); // 复制管道的子进程端到标准输出
    dup2(fd[1], STDIN_FILENO);  // 复制管道的子进程端到标准输入
    close(fd[1]);   // 关闭已复制的读管道
    /* 使用exec执行命令 */
} else {    // 父进程
    close(fd[1]);   // 关闭管道的子进程端
    /* 现在可在fd[0]中读写数据 */
    shutdown(fd[0], SHUT_WR);   // 通知对端数据发送完毕
    /* 读取剩余数据 */
    close(fd[0]);   // 关闭管道
    /* 使用wait系列函数等待子进程退出并取得退出代码 */
}
