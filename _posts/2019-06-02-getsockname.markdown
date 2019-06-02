---
title: getsockname getpeername
layout: post
category: linux
author: 夏泽民
---
　这两个函数或者返回与某个套接字关联的本地协议地址（getsockname），或者返回与某个套接字关联的外地协议地址即得到对方的地址（getpeername）。
#include <sys/socket.h>
 
int getsockname(int sockfd,struct sockaddr* localaddr,socklen_t *addrlen);
int getpeername(int sockfd,struct sockaddr* peeraddr,socklen_t *addrlen);
均返回：若成功则为0，失败则为-1
　　　　　　　　　　　　　　　　　　　　　
getsockname可以获得一个与socket相关的地址，服务器端可以通过它得到相关客户端地址，而客户端也可以得到当前已连接成功的socket的ip和端口。

第二种情况在客户端不进行bind而直接连接服务器时，而且客户端需要知道当前使用哪个ip进行通信时比较有用（如多网卡的情况）
<!-- more -->
getpeername只有在连接建立以后才调用，否则不能正确获得对方地址和端口，所以它的参数描述字一般是已连接描述字而非监听套接口描述字。
    没有连接的UDP不能调用getpeername，但是可以调用getsockname和TCP一样，它的地址和端口不是在调用socket就指定了，而是在第一次调用sendto函数以后。
    已经连接的UDP，在调用connect以后，这2个函数（getsockname，getpeername）都是可以用的。但是这时意义不大，因为已经连接（connect）的UDP已经知道对方的地址。

　　需要这两个函数的理由如下：

　在一个没有调用bind的TCP客户上，connect成功返回后，getsockname用于返回由内核赋予该连接的本地IP地址和本地端口号。

　在以端口号为0调用bind（告知内核去选择本地临时端口号）后，getsockname用于返回由内核赋予的本地端口号。

　在一个以通配IP地址调用bind的TCP服务器上，与某个客户的连接一旦建立（accept成功返回），getsockname就可以用于返回由内核赋予该连接的本地IP地址。在这样的调用中，套接字描述符参数必须是已连接套接字的描述符，而不是监听套接字的描述符。

　当一个服务器的是由调用过accept的某个进程通过调用exec执行程序时，它能够获取客户身份的唯一途径便是调用getpeername。
　
获取该值有两个常用方法：

　　调用exec的进程可以把这个描述符格式化成一个字符串，再把它作为一个命令行参数传递给新程序。

　　约定在调用exec之前，总是把某个特定描述符置为所接受的已连接套接字的描述符。

　　inetd采用的是第二种方法，它总是把描述符0、1、2置为所接受的已连接套接字的描述符（即将已连接套接字描述符dup到描述符0、1、2，然后close原连接套接字）。
　　
#include "unp.h"
int main(int argc, char ** argv)
{ 
　　int listenfd,connfd; 
　　struct sockaddr_in servaddr; 
　　pid_t pid; char temp[20]; 
　　listenfd = Socket(AF_INET, SOCK_STREAM, 0); 

　　bzero(&servaddr, sizeof(servaddr)); 
　　servaddr.sin_family = AF_INET; 
　　servaddr.sin_addr.s_addr = htonl(INADDR_ANY); 
　　servaddr.sin_port = htons(10010); 
　　
　　Bind(listenfd, (SA *)&servaddr, sizeof(servaddr)); 

　　Listen(listenfd, LISTENQ); 
　　
　　for( ; ; ) 
　　{ 
　　　　struct sockaddr_in local; 
　　　　connfd = Accept(listenfd, (SA *)NULL, NULL); 
　　　　if((pid = fork()) == 0) 
　　　　{ 
　　　　　　Close(listenfd);
　　　　　　struct sockaddr_in serv, guest; 
　　　　　　char serv_ip[20]; 
　　　　　　char guest_ip[20]; 
　　　　　　socklen_t serv_len = sizeof(serv); 
　　　　　　socklen_t guest_len = sizeof(guest); 
　　　　　　getsockname(connfd, (struct sockaddr *)&serv, &serv_len); 
　　　　　　getpeername(connfd, (struct sockaddr *)&guest, &guest_len); 
　　　　　　Inet_ntop(AF_INET, &serv.sin_addr, serv_ip, sizeof(serv_ip)); 
　　　　　　Inet_ntop(AF_INET, &guest.sin_addr, guest_ip, sizeof(guest_ip)); 
　　　　　　printf("host %s:%d guest %s:%dn", serv_ip, ntohs(serv.sin_port), guest_ip, ntohs(guest.sin_port)); 
　　　　　　char buf[] = "hello world"; 
　　　　　　Write(connfd, buf, strlen(buf)); 
　　　　　　Close(connfd); exit(0);
 　　　　} 
　　　　Close(connfd); 
　　} 
}


#include <winsock.h>
int PASCAL FAR getsockname( SOCKET s, struct sockaddr FAR* name,
int FAR* namelen);
s：标识一个已捆绑套接口的描述字。
name：接收套接口的地址（名字）。
namelen：名字缓冲区长度。


