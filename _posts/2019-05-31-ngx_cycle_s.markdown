---
title: ngx_cycle_s
layout: post
category: php
author: 夏泽民
---
https://www.kancloud.cn/digest/understandingnginx/202596

Nginx框架是围绕着ngx_cycle_t结构体运行的。ngx_cycle_t结构体中包含的信息主要可以分为以下部分：

所有模块的配置信息
Nginx运行时所需要的一些资源，包括连接池，内存池，打开文件，操作目录等等
<!-- more -->
struct ngx_cycle_s { 
    //保存所有模块的配置结构体
    void                  ****conf_ctx;
    //内存池
    ngx_pool_t               *pool;

    //日志信息
    ngx_log_t                *log;
    ngx_log_t                 new_log;

    ngx_uint_t                log_use_stderr;  /* unsigned  log_use_stderr:1; */

    ngx_connection_t        **files;
    //连接池
    ngx_connection_t         *free_connections;
    ngx_uint_t                free_connection_n;

    ngx_queue_t               reusable_connections_queue;

    //被监听的端口对应的ngx_listen_t数组
    ngx_array_t               listening;
    //操作目录
    ngx_array_t               paths;
    //打开的文件
    ngx_list_t                open_files;
    //共享内存
    ngx_list_t                shared_memory;

    //当前进程中的所有连接对象的总数
    ngx_uint_t                connection_n;
    ngx_uint_t                files_n;

    //指向当前进程中的所有连接对象
    ngx_connection_t         *connections;
    //当前进程中的所有读写事件，每个读写事件对应一个连接，所
    //以读写事件的总数分别都是connection_n
    ngx_event_t              *read_events;
    ngx_event_t              *write_events;

    ngx_cycle_t              *old_cycle;//

    //配置文件信息
    ngx_str_t                 conf_file;
    ngx_str_t                 conf_param;
    ngx_str_t                 conf_prefix;
    ngx_str_t                 prefix;
    //用于进程间同步的文件锁
    ngx_str_t                 lock_file;
    //使用getthehostname得到的主机名
    ngx_str_t                 hostname;
};

关于conf_ctx
这里说明一下里面的conf_ctx数据成员，这是一个多维指针。它主要是保存模块的配置项信息。所谓的模块配置项信息，我们知道nginx是高度模块化的，他的各个功能都由不同的模块构成，这使得系统具有很好的灵活性和可扩展性。在nginx中的配置文件中，会列出很多配置项和配置项对应的配置值，而每个模块都有自己感兴趣的配置项，nginx把一个模块所有感兴趣的配置项放在一个结构体中，称为这个结构体的配置结构体。conf_ctx保存的就是所有模块的配置结构体。，可以看到conf_ctx是void****类型的。因此可以看出，conf_ctx是一个数组，每个数组里面包含一个void***类型的指针，也就是说conf_ctx指向的数组中的每个元素都是一个指针，这个指针再次指向一个指针数组

一般来说，level0中的数组中的每个元素指向一个核心模块构建的配置结构体。核心模块构建的结构体功能主要是为了组织管理归属于这个核心模块的同类型模块。但是按这样说的话，所有归属于这个核心模块的同类型模块只需要一个数组就可以了，为什么这里却需要两层的数组，也就是level2到底是做什么用的呢？

这里就要讲到配置文件的格式了，配置文件具有层级嵌套格式，如下以http类型模块为例

http {

    server {

        location / {

        }

        location {

        }
    }
}

在针对该模块的配置文件中，最外层是http层，也成为main层，main块里面有sever块，server块里面可能嵌套location块，而相同的配置项可能在不同的块中都有被设置，也就是说，不同配置块可能会对同一个模块的配置结构体产生影响，于是nginx为了提高配置灵活性，干脆为每个块都创建一个这个配置块可能会影响到的模块的配置结构体，这就解释了level2的作用


1. ngx_cycle_t
Nginx框架是围绕着ngx_cycle_t结构体来控制进程运行。

//我们来看一看这神奇的ngx_cycle_s结构体吧,看一看庐山真面目.  
struct ngx_cycle_s {  
  
    /* 保存着所有模块配置项的结构体指针p,它首先是一个数组, 
       该数组每个成员又是一个指针,这个指针又指向了存储着指针的数组.*/  
    void                  **** conf_ctx ;  
  
    //内存池  
    ngx_pool_t               * pool ;  
  
    /*日志模块中提供了生成基本ngx_log_t日志对象的功能, 
      这里的log实际上是在在没有执行ngx_init_cycle方法前, 
      也就是还没有解析配置项前,如果有信息需要输出到日志, 
      就会暂时使用log对象,它会输出到屏幕.在执行ngx_init_cycle方法后, 
      会根据nginx.conf配置文件中的配置项,构造出正确的日志文件, 
      此时会对log重新赋值. */  
    ngx_log_t                * log ;  
  
    /* 由nginx.conf配置文件读取到日志路径后,将开始初始化error_log日志文件, 
       由于log对象还在用于输出日志到屏幕,这时候new_log将暂时替代log, 
       待初始化成功后,会用new_log的地址覆盖上面的指针.*/  
    ngx_log_t                 new_log ;  
  
    //下面files数组里元素的总数  
    ngx_uint_t                files_n ;  
    /* 
       对于epoll,rtsig这样的事件模块,会以有效文件句柄树来预先建立 
       这些ngx_connection_t结构体,以加速事件的收集,分发.这时files就会 
       保存所有ngx_connection_t的指针组成的数组,而文件句柄的值用来访问 
       files数组成员. */  
    ngx_connection_t        ** files ;  
  
    //可用连接池,与free_connection_n配合使用  
    ngx_connection_t         * free_connections ;  
    //连接池的个数  
    ngx_uint_t                free_connection_n ;  
  
    //可重用的连接队列  
    ngx_queue_t               reusable_connections_queue ;  
  
    //动态数组,每个成员存储ngx_listening_t成员,表示监听端口以及相关的参数  
    ngx_array_t               listening ;  
  
    /* 动态数组,保存着nginx所有要操作的目录,如果目录不在,将创建, 
       如果创建失败(例如权限不够),nginx将启动失败.*/  
    ngx_array_t               paths ;  
  
    /* 单链表容器,元素类型是ngx_open_file_t结构体,他表示nginx已经打开的 
       所有文件.事实上nginx框架不会向open_files链表中添加文件, 
       而是由对此感兴趣的模块向其中添加文件路径名, 
       nginx会在ngx_init_cycle方法中打开这些文件.*/  
    ngx_list_t                open_files ;  
     
    /* 单链表容器,每一个元素表示一块共享内存*/  
    ngx_list_t                shared_memory ;  
  
  
    //当前进程中所有连接对象的总数,于下面的connections成员配合使用.  
    ngx_uint_t                connection_n ;  
  
    //指向当前进程中多有连接对象  
    ngx_connection_t         * connections ;  
  
    //指向当前进程中所有读事件  
    ngx_event_t              * read_events ;  
  
    //指向当前进程中所有写事件  
    ngx_event_t              * write_events ;  
  
    /* 旧的ngx_cycle_t对象,用于引用上一个ngx_cycle_t对象中的成员. 
       例如 ngx_init_cycle方法在启动初期,需要建立一个临时ngx_cycle_t 
       对象来保存一些变量,再调用ngx_init_cycle方法时,就可以把 
       旧的ngx_cycle_t对象传进去,而这时,这个old_cycle指针 
       就会保存这个前期的ngx_cycle_t对象 */  
    ngx_cycle_t              * old_cycle ;  
  
    //配置文件相对于安装目录的路径  
    ngx_str_t                 conf_file ;  
  
    /* nginx处理配置文件时需要特殊处理的在命令行携带的参数, 
       一般是 -g选项携带的参数 */  
    ngx_str_t                 conf_param ;  
  
    //nginx配置文件路径  
    ngx_str_t                 conf_prefix ;  
  
    //nginx安装路径  
    ngx_str_t                 prefix ;  
  
    //用于进程间同步的文件所名称  
    ngx_str_t                 lock_file ;  
  
    //使用gethostname系统调用得到的主机名  
    ngx_str_t                 hostname ;  
}; 
 

2. ngx_listening_t
ngx_cycle_t对象中有一个动态数组成员叫做listening，它的每个数组元素都是ngx_listening_t结构体，而每个ngx_listening_t结构体又代表着Nginx服务器监听的一个端口。

struct ngx_listening_s {  
    ngx_socket_t        fd;//套接字句柄  
  
    struct sockaddr    *sockaddr;//监听sockaddr地址  
    socklen_t           socklen;    /*sockaddr地址长度 size of sockaddr */  
    size_t              addr_text_max_len;//存储ip地址的字符串addr_text最大长度  
    ngx_str_t           addr_text;//以字符串形式存储ip地址  
    int                 type;   //套接字类型。types是SOCK_STREAM时，表示是tcp 
    int                 backlog;  //TCP实现监听时的backlog队列，它表示允许正在通过三次握手建立tcp连接但还没有任何进程开始处理的连接最大个数
    int                 rcvbuf;//套接字接收缓冲区大小  
    int                 sndbuf;//套接字发送缓冲区大小  
  
    /* handler of accepted connection */  
    ngx_connection_handler_pt   handler;//当新的tcp连接成功建立后的处理方法    
    void               *servers; //实际上框架并不使用servers指针，它更多是作为一个保留指针，目前主要用于HTTP或者mail等模块，用于保存当前监听端口对应着的所有主机名
  
    ngx_log_t           log;//日志  
    ngx_log_t          *logp;//日志指针  
  
    size_t              pool_size;//如果为新的tcp连接创建内存池，则内存池的初始大小应该是pool_size。  
    /* should be here because of the AcceptEx() preread */  
    size_t              post_accept_buffer_size;  
    /* should be here because of the deferred accept */  
    ngx_msec_t          post_accept_timeout;//连接建立成功后，如果 post_accept_timeout秒后仍然没有收到用户的数据，就丢弃该连接  
     
    ngx_listening_t    *previous;   //前一个ngx_listening_t结构，多个ngx_listening_t结构体之间由previous指针组成单链表
    ngx_connection_t   *connection;//当前监听句柄对应的ngx_connection_t结构体  
  
    unsigned            open:1;//为1表示监听句柄有效，为0表示正常关闭  
    unsigned            remain:1;//为1表示不关闭原先打开的监听端口，为0表示关闭曾经打开的监听端口  
    unsigned            ignore:1;//为1表示跳过设置当前ngx_listening_t结构体中的套接字，为0时正常初始化套接字  
  
    unsigned            bound:1;       /* already bound */  
    unsigned            inherited:1;   /* inherited from previous process */  
    unsigned            nonblocking_accept:1;  
    unsigned            listen:1;//为1表示当前结构体对应的套接字已经监听  
    unsigned            nonblocking:1; //目前该标志位没有意义 
    unsigned            shared:1;    /*目前该标志位没有意义， shared between threads or processes */  
    unsigned            addr_ntop:1;//为1表示将网络地址转变为字符串形式的地址  
  
#if (NGX_HAVE_INET6 && defined IPV6_V6ONLY)  
    unsigned            ipv6only:2;  
#endif  
  
#if (NGX_HAVE_DEFERRED_ACCEPT)  
    unsigned            deferred_accept:1;  
    unsigned            delete_deferred:1;  
    unsigned            add_deferred:1;  
#ifdef SO_ACCEPTFILTER  
    char               *accept_filter;  
#endif  
#endif  
#if (NGX_HAVE_SETFIB)  
    int                 setfib;  
#endif  
  
}; 
3. ngx_connection_t
Nginx定义了基本的数据结构ngx_connection_t来表示连接。由客户端主动发起、Nginx服务器被动接收的TCP连接，这类可以称为被动连接。还有一类连接，在某些请求的处理过程中，Nginx会试图主动向其他上游服务器建立连接，并以此连接与上游服务器通信，Nginx定义ngx_peer_connection_t结构来表示，这类可以称为主动连接。本质上来说，主动连接是以ngx_connection_t结构体为基础实现的。

struct ngx_connection_s {  
    //连接未使用时，data用于充当连接池中空闲链表中的next指针。连接使用时由模块而定，HTTP中，data指向ngx_http_request_t  
    void               *data;  
    ngx_event_t        *read;//连接对应的读事件  
    ngx_event_t        *write;//连接对应的写事件  
  
    ngx_socket_t        fd;//套接字对应的句柄  
  
    ngx_recv_pt         recv;//直接接收网络字符流的方法  
    ngx_send_pt         send;//直接发送网络字符流的方法  
    ngx_recv_chain_pt   recv_chain;//以链表来接收网络字符流的方法  
    ngx_send_chain_pt   send_chain;//以链表来发送网络字符流的方法  
    //这个连接对应的ngx_listening_t监听对象，此连接由listening监听端口的事件建立  
    ngx_listening_t    *listening;  
  
    off_t               sent;//这个连接上已发送的字节数  
  
    ngx_log_t          *log;//日志对象  
       /*内存池。一般在accept一个新的连接时，会创建一个内存池，而在这个连接结束时会销毁内存池。内存池大小是由上面listening成员的pool_size决定的*/  
    ngx_pool_t         *pool;  
  
    struct sockaddr    *sockaddr;//连接客户端的sockaddr  
    socklen_t           socklen;//sockaddr结构体的长度  
    ngx_str_t           addr_text;//连接客户段字符串形式的IP地址  
  
#if (NGX_SSL)  
    ngx_ssl_connection_t  *ssl;  
#endif  
    //本机监听端口对应的sockaddr结构体，实际上就是listening监听对象的sockaddr成员  
    struct sockaddr    *local_sockaddr;  
  
    ngx_buf_t          *buffer;//用户接受、缓存客户端发来的字符流，buffer是由连接内存池分配的，大小自由决定  
    /*用来将当前连接以双向链表元素的形式添加到ngx_cycle_t核心结构体的reuseable_connection_queue双向链表中，表示可以重用的连接*/  
    ngx_queue_t         queue;  
    /*连接使用次数。ngx_connection_t结构体每次建立一条来自客户端的连接，或者主动向后端服务器发起连接时，number都会加1*/  
    ngx_atomic_uint_t   number;  
  
    ngx_uint_t          requests;//处理的请求次数  
    //缓存中的业务类型。  
    unsigned            buffered:8;  
    //本连接的日志级别，占用3位，取值范围为0～7，但实际只定义了5个值，由ngx_connection_log_error_e枚举表示。  
    unsigned            log_error:3;     /* ngx_connection_log_error_e */  
  
    unsigned            single_connection:1;//为1时表示独立的连接，为0表示依靠其他连接行为而建立起来的非独立连接  
    unsigned            unexpected_eof:1;//为1表示不期待字符流结束  
    unsigned            timedout:1;//为1表示连接已经超时  
    unsigned            error:1;//为1表示连接处理过程中出现错误  
    unsigned            destroyed:1;//为1表示连接已经销毁  
  
    unsigned            idle:1;//为1表示连接处于空闲状态，如keepalive两次请求中间的状态  
    unsigned            reusable:1;//为1表示连接可重用，与上面的queue字段对应使用  
    unsigned            close:1;//为1表示连接关闭  
  
    unsigned            sendfile:1;//为1表示正在将文件中的数据发往连接的另一端  
    /*为1表示只有连接套接字对应的发送缓冲区必须满足最低设置的大小阀值时，事件驱动模块才会分发该事件。这与ngx_handle_write_event方法中的lowat参数是对应的*/  
    unsigned            sndlowat:1;  
    unsigned            tcp_nodelay:2;   /* ngx_connection_tcp_nodelay_e */  
    unsigned            tcp_nopush:2;    /* ngx_connection_tcp_nopush_e */  
  
#if (NGX_HAVE_IOCP)  
    unsigned            accept_context_updated:1;  
#endif  
  
#if (NGX_HAVE_AIO_SENDFILE)  
    unsigned            aio_sendfile:1;  
    ngx_buf_t          *busy_sendfile;  
#endif  
  
#if (NGX_THREADS)  
    ngx_atomic_t        lock;  
#endif  
}; 
 

4. ngx_peer_connection_t主动连接
Nginx主动想后端服务器发起的连接称作主动连接，用ngx_peer_connection_t结构表示，其实质是对ngx_connection_t作了一层封装。

/* 通过该方法由连接池中获取一个新连接 */  
typedef ngx_int_t (*ngx_event_get_peer_pt)(ngx_peer_connection_t *pc,  
    void *data);  
      
/* 通过该方法将使用完毕的连接释放给连接池 */  
typedef void (*ngx_event_free_peer_pt)(ngx_peer_connection_t *pc, void *data,  
    ngx_uint_t state);  
  
具体结构如下：  
struct ngx_peer_connection_s {  
    /* 主动连接实际上需要被动连接结构的大部分成员，相当于重用 */  
    ngx_connection_t                *connection;  
  
    /* 远端服务器socket信息 */  
    struct sockaddr                 *sockaddr;  
    socklen_t                        socklen;  
    ngx_str_t                       *name;  
  
    /* 连接失败后可以重试的次数 */  
    ngx_uint_t                       tries;  
  
    /* 从连接池中获取长连接，必须使用该方法 */  
    ngx_event_get_peer_pt            get;  
    ngx_event_free_peer_pt           free;  
      
    /* 该data成员作为上面方法的传递参数 */  
    void                            *data;  
  
    /* 本地地址信息 */  
    ngx_addr_t                      *local;  
  
    /* 接收缓冲区大小 */  
    int                              rcvbuf;  
  
    ngx_log_t                       *log;  
  
    /* 标识该连接已经缓存 */  
    unsigned                         cached:1;  
  
    /* ngx_connection_log_error_e */  
    unsigned                         log_error:2;  
}; 
 

ngx_connection_t连接池
      Nginx在接受客户端的连接所使用的ngx_connection_t结构都是在启动阶段就预分配好的，使用时只需从连接池中获取。
       ngx_cycle_t中的connections和free_connections俩成员构成一个连接池。其中，connections指向整个连接池数组首部，free_connections指向第一个ngx_connection_t空闲连接。所有空闲连接ngx_connection_t都以data成员作为next指针串联成一个单链表。这样，一旦有用户发起连接就从free_connections指向的链表头获取一个空闲连接。在连接释放时，只需把该连接再插入到free_connections链表头即可。

       同时，ngx_cycle_t核心结构还有一套事件池，即Nginx认为每一个连接至少需要一个读事件和一个写事件，有多少个连接就分配多少个读、写事件。这里，读、写、连接池都是由三个大小相同的数组组成，根据数组下标就可将每一个连接、读、写事件对应起来。这三个数组的大小都由nginx.conf的connections配置项决定。
