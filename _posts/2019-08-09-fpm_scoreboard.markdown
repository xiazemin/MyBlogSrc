---
title: fpm_scoreboard
layout: post
category: php
author: 夏泽民
---
fpm_scoreboard(以下简称scoreboard模块)是PHP-FPM核心功能之一，源码位于sapi/fpm/fpm_scoreboard.c。从字面意思理解是一个”记分器”，实际上，是FPM内置的一个worker进程统计功能模块。
<!-- more -->
scoreboard模块定义fpm_scoreboard_s和fpm_scoreboard_proc_s两种数据结构。

struct fpm_scoreboard_s {
    union {
        atomic_t lock;
        char dummy[16];
    };//锁状态
    char pool[32];//实例名称 例如：[www]
    int pm; //PM运行模式
    time_t start_epoch; //开始时间
    int idle;//procs的空闲数
    int active;//procs的使用数
    int active_max; //最大procs使用数
    unsigned long int requests;
    unsigned int max_children_reached; //到达最大进程数限制的次数
    int lq; //当前listen queue的请求数(accept操作，可以过tcpi_unacked或getsocketopt获取)
    int lq_max;//listen queue大小
    unsigned int lq_len;
    unsigned int nprocs; //procs总数
    int free_proc; //从procs列表遍历下一个空闲对象的开始下标
    struct fpm_scoreboard_proc_s *procs[]; //列表
};
struct fpm_scoreboard_proc_s {
    union {
        atomic_t lock;
        char dummy[16];
    };//锁状态
    int used; //使用标识 0=未使用 1=正在使用
    time_t start_epoch; //使用开始时间
    pid_t pid; //进程id
    unsigned long requests; //处理请求次数
    enum fpm_request_stage_e request_stage; //处理请求阶段
    struct timeval accepted; //accept请求时间
    struct timeval duration; //脚本总执行时间
    time_t accepted_epoch;//accept请求时间戳(秒)
    struct timeval tv; //活跃时间
    char request_uri[128]; //请求路径
    char query_string[512]; //请求参数
    char request_method[16]; //请求方式
    size_t content_length; //请求内容长度 /* used with POST only */
    char script_filename[256];//脚本名称
    char auth_user[32];
#ifdef HAVE_TIMES
    struct tms cpu_accepted;
    struct timeval cpu_duration;
    struct tms last_request_cpu;
    struct timeval last_request_cpu_duration;
#endif
    size_t memory;//脚本占用的内存大小
};
fpm_scoreboard_s结构记录FPM所有worker进程的统计信息，其中*procs数组结构保存各worker进程统计单元。fpm_scoreboard_proc_s记录worker进程的运行状态信息(后文称为统计单元)。

根据上面代码的注释可以很容易地理解其属性的含义，在这就不再逐一说明。

下面将重点介绍scoreboard的运行流程。

scoreboard模块初始化

FPM内部通过执行fpm_init()->fpm_scoreboard_init_main()来完成scoreboard模块的操作。

int fpm_scoreboard_init_main()
{
    //...省略部分代码...
    wp->scoreboard = fpm_shm_alloc(sizeof(struct fpm_scoreboard_s) + (wp->config->pm_max_children - 1) * sizeof(struct fpm_scoreboard_proc_s *));
    if (!wp->scoreboard) {
        return -1;
    }
    wp->scoreboard->nprocs = wp->config->pm_max_children;
    for (i = 0; i < wp->scoreboard->nprocs; i++) {
        wp->scoreboard->procs[i] = fpm_shm_alloc(sizeof(struct fpm_scoreboard_proc_s));
        if (!wp->scoreboard->procs[i]) {
            return -1;
        }
        memset(wp->scoreboard->procs[i], 0, sizeof(struct fpm_scoreboard_proc_s));
    }
    //...省略部分代码...
}
fpm_scoreboard_init_main()函数调用fpm_shm_alloc()为 wp->scoreboard 分配空间，大小根据 wp->config->pm_max_children 参数计算。wp->config->pm_max_children参数对应php-fpm.conf的pm.max_children配置项，表示FPM允许启动的最大worker进程数。从而保证每个worker进程都能分配到一个可用的统计单元。然后再对wp->scoreboard->procs的每个统计单元进行初始化。

上面已经初始化统计单元列表。那这些统计单元是如何分配给每一个worker进程？

scoreboard统计单元分配

从源码(fpm_children.c)可以看到，FPM每次调用fork()新worker进程之前，系统都会执行fpm_resources_prepare()函数。代码如下：

static struct fpm_child_s *fpm_resources_prepare(struct fpm_worker_pool_s *wp)
{
    struct fpm_child_s *c;

    c = fpm_child_alloc();
    //...省略部分代码...
    //此时c->scoreboard_i=-1
    if (0 > fpm_scoreboard_proc_alloc(wp->scoreboard, &c->scoreboard_i)) {
        fpm_stdio_discard_pipes(c);
        fpm_child_free(c);
        return 0;
    }
    return c;
}
fpm_resources_prepare()函数：首先，初始化child对象；然后调用fpm_scoreboard_proc_alloc()函数。继续追踪fpm_scoreboard_proc_alloc()函数代码。

int fpm_scoreboard_proc_alloc(struct fpm_scoreboard_s *scoreboard, int *child_index)
{
    int i = -1;
    //...省略部分代码...
    if (scoreboard->free_proc >= 0 && scoreboard->free_proc < scoreboard->nprocs) {
        if (scoreboard->procs[scoreboard->free_proc] && !scoreboard->procs[scoreboard->free_proc]->used) {
            i = scoreboard->free_proc;
        }
    }
    if (i < 0) { 
        for (i = 0; i < scoreboard->nprocs; i++) {
            if (scoreboard->procs[i] && !scoreboard->procs[i]->used) {
                break;
            }
        }
    }
    if (i < 0 || i >= scoreboard->nprocs) {
        return -1;
    }
    //打上“使用”标记
    scoreboard->procs[i]->used = 1;
    *child_index = i;
    //重置寻找下一个空闲单元的起始下标
    if (i + 1 >= scoreboard->nprocs) {
        scoreboard->free_proc = 0;
    } else {
        scoreboard->free_proc = i + 1;
    }
    return 0;
}
fpm_scoreboard_proc_alloc()函数的代码比较容易读懂。首先判断scoreboard->free_proc位置对应的元素是否可用；如果不可用，则继续遍历scoreboard->nprocs查找可用的元素；如果找到，则修改该元素的used标识，并重置scoreboard->free_proc属性值。

结合上面两个函数代码，我们可以看出fpm_resources_prepare()的功能是为每个新的worker进程分配一个可用的统计单元。

我们已经知道worker进程统计单元的分配流程，那这些统计单元在worker进程是如何运用的？有哪些功能或者模块已使用？

scoreboard统计单元运用

在介绍之前，我们先来了解scoreboard模块的fpm_scoreboard_update()、fpm_scoreboard_proc_acquire()、fpm_scoreboard_proc_release()三个函数的功能。

void fpm_scoreboard_update(int idle, int active, int lq, int lq_len, int requests, int max_children_reached, int action, struct fpm_scoreboard_s *scoreboard)
{
    //...省略部分代码...
    fpm_spinlock(&scoreboard->lock, 0);
    if (action == FPM_SCOREBOARD_ACTION_SET) {
        if (idle >= 0) {
            scoreboard->idle = idle;
        }
        if (active >= 0) {
            scoreboard->active = active;
        }
        //...省略部分代码...
    } else {
        if (scoreboard->idle + idle > 0) {
            scoreboard->idle += idle;
        } else {
            scoreboard->idle = 0;
        }
        if (scoreboard->active + active > 0) {
            scoreboard->active += active;
        } else {
            scoreboard->active = 0;
        }
        //...省略部分代码...
    }
    if (scoreboard->active > scoreboard->active_max) {
        scoreboard->active_max = scoreboard->active;
    }
    fpm_unlock(scoreboard->lock);
}
fpm_scoreboard_update()：修改wp->scoreboard各属性值，该函数内部引用”锁”机制来保证数据的原子性。action有两个值FPM_SCOREBOARD_ACTION_SET和FPM_SCOREBOARD_ACTION_INC。当action=FPM_SCOREBOARD_ACTION_SET,表示这是一个重置操作。当action=FPM_SCOREBOARD_ACTION_INC时，代表这是一个求和操作。

struct fpm_scoreboard_proc_s *fpm_scoreboard_proc_acquire(struct fpm_scoreboard_s *scoreboard, int child_index, int nohang)
{
    //...省略部分代码...
    proc = fpm_scoreboard_proc_get(scoreboard, child_index);
    if (!proc) {
        return NULL;
    }
    //请求锁
    if (!fpm_spinlock(&proc->lock, nohang)) {
        return NULL;
    }
    return proc;
}
struct fpm_scoreboard_proc_s *fpm_scoreboard_proc_get(struct fpm_scoreboard_s *scoreboard, int child_index)
{
    //...省略部分代码...
    return scoreboard->procs[child_index];
}

fpm_scoreboard_proc_acquire()：获取统计单元(wp->scoreboard->procs[i]),并请求对象锁。值得注意的是这里的”锁”与`fpm_scoreboard_update()的”锁”不是同一个。

//释放对象锁
void fpm_scoreboard_proc_release(struct fpm_scoreboard_proc_s *proc) /* {\{\{ */
{
    //...省略部分代码...
    proc->lock = 0;
}
fpm_scoreboard_proc_release()：释放对象锁。

回到刚才那个问题，我们接下来继续分析。worker进程处理客户端请求的完整流程共分为5个阶段。如下。

阶段	备注
FPM_REQUEST_ACCEPTING	空闲状态(等待请求)
FPM_REQUEST_READING_HEADERS	读取头信息
FPM_REQUEST_INFO	获取请求信息
FPM_REQUEST_EXECUTING	执行状态
FPM_REQUEST_END	请求结束状态
1. 接收客户端请求阶段,执行流程：fpm_request_accepting().

void fpm_request_accepting(){
    //...省略部分代码...
    proc = fpm_scoreboard_proc_acquire(NULL, -1, 0);
    //...省略部分代码...
    proc->request_stage = FPM_REQUEST_ACCEPTING;
    proc->tv = now;
    fpm_scoreboard_proc_release(proc);
    /* idle++, active-- */
    fpm_scoreboard_update(1, -1, 0, 0, 0, 0, FPM_SCOREBOARD_ACTION_INC, NULL);
}
fpm_scoreboard_proc_acquire()得到管理该worker进程状态的统计单元，然后修改统计单元(wp->scoreboard->procs[i])的request_stage及tv属性值，最后调用fpm_scoreboard_update修改wp->scoreboard统计信息。 
2. 从FASTCGI读取客户端请求头阶段,执行流程：fpm_request_reading_headers().

void fpm_request_reading_headers()
{
    //...省略部分代码...
    proc = fpm_scoreboard_proc_acquire(NULL, -1, 0);
    //...省略部分代码...
    proc->request_stage = FPM_REQUEST_READING_HEADERS;
    //记录当前时间
    proc->tv = now;
    proc->accepted = now;
    proc->accepted_epoch = now_epoch;
#ifdef HAVE_TIMES
    proc->cpu_accepted = cpu;
#endif
    proc->requests++;
    proc->request_uri[0] = '\0';
    proc->request_method[0] = '\0';
    proc->script_filename[0] = '\0';
    proc->query_string[0] = '\0';
    proc->query_string[0] = '\0';
    proc->auth_user[0] = '\0';
    proc->content_length = 0;
    fpm_scoreboard_proc_release(proc);
    /* idle--, active++, request++ */
    fpm_scoreboard_update(-1, 1, 0, 0, 1, 0, FPM_SCOREBOARD_ACTION_INC, NULL);
}
该阶段重置request_stage为FPM_REQUEST_READING_HEADERS,然后分别修改tv、accepted、accepted_epoch、requests值，并重置统计单元请求头属性(request_uri、request_method、script_filename等)数据。最后调用fpm_scoreboard_update()修改wp->scoreboard的idle、active、request属性。 
3. FPM获取请求信息阶段，执行代码：fpm_request_info().

void fpm_request_info()
{
    //...省略部分代码...
    proc = fpm_scoreboard_proc_acquire(NULL, -1, 0);
    if (proc == NULL) {
        zlog(ZLOG_WARNING, "failed to acquire proc scoreboard");
        return;
    }
    proc->request_stage = FPM_REQUEST_INFO;
    proc->tv = now;
    //请求地址
    if (request_uri) {
        strlcpy(proc->request_uri, request_uri, sizeof(proc->request_uri));
    }
    //请求方法
    if (request_method) {
        strlcpy(proc->request_method, request_method, sizeof(proc->request_method));
    }
    //...省略部分代码...
    fpm_scoreboard_proc_release(proc);
}

从上面的代码可以看出，该阶段只是将请求信息保存在该worker进程的统计单元中。与前两者不同的是，此阶段无需调用fpm_scoreboard_update()修改wp->scoreboard属性值。 
4. zend执行脚本阶段，执行代码：fpm_request_executing().

void fpm_request_executing() 
{
    //...省略部分代码...
    proc = fpm_scoreboard_proc_acquire(NULL, -1, 0);
    if (proc == NULL) {
        zlog(ZLOG_WARNING, "failed to acquire proc scoreboard");
        return;
    }

    proc->request_stage = FPM_REQUEST_EXECUTING;
    proc->tv = now;
    fpm_scoreboard_proc_release(proc);
}
此阶段将request_stage属性值修改为FPM_REQUEST_EXECUTING。 
5. 请求结束阶段，执行代码：fpm_request_end().

void fpm_request_end(TSRMLS_D)
{
    //...省略代码...
    proc = fpm_scoreboard_proc_acquire(NULL, -1, 0);
    if (proc == NULL) {
        zlog(ZLOG_WARNING, "failed to acquire proc scoreboard");
        return;
    }
    proc->request_stage = FPM_REQUEST_FINISHED;
    proc->tv = now;
    timersub(&now, &proc->accepted, &proc->duration);
#ifdef HAVE_TIMES
    timersub(&proc->tv, &proc->accepted, &proc->cpu_duration);
    proc->last_request_cpu.tms_utime = cpu.tms_utime - proc->cpu_accepted.tms_utime;
    proc->last_request_cpu.tms_stime = cpu.tms_stime - proc->cpu_accepted.tms_stime;
    proc->last_request_cpu.tms_cutime = cpu.tms_cutime - proc->cpu_accepted.tms_cutime;
    proc->last_request_cpu.tms_cstime = cpu.tms_cstime - proc->cpu_accepted.tms_cstime;
#endif
    proc->memory = memory;
    fpm_scoreboard_proc_release(proc);
}
此阶段计算程序运行的时间及占用的内存大小，分别保存在统计单元的duration和memory属性。

scoreboard统计定时更新 
FPM内部定义了一个fpm_pctl_perform_idle_server_maintenance_heartbeat定时器,其内部会进行统计worker进程的idle、active等数据，然后调用fpm_scoreboard_update()功能进行更新wp->scoreboard信息。

static void fpm_pctl_perform_idle_server_maintenance(struct timeval *now) /* {\{\{ */
{
    struct fpm_worker_pool_s *wp;
    //...省略部分代码...
    for (wp = fpm_worker_all_pools; wp; wp = wp->next) {
        for (child = wp->children; child; child = child->next) {
            if (fpm_request_is_idle(child)) {
                if (last_idle_child == NULL) {
                    last_idle_child = child;
                } else {
                    if (timercmp(&child->started, &last_idle_child->started, <)) {
                        last_idle_child = child;
                    }
                }
                idle++;
            } else {
                active++;
            }
        }
        fpm_scoreboard_update(idle, active, cur_lq, -1, -1, -1, FPM_SCOREBOARD_ACTION_SET, wp->scoreboard);
        //...省略部分代码...
    }
    //...省略部分代码...
}
以上就是PHP-FPM scoreboard模块介绍的全部内容。通过这个模块，我们可以快速地掌握PHP-FPM各个worker进程运行状态。


