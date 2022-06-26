---
title: grpc deadlines 错误处理
layout: post
category: golang
author: 夏泽民
---
当未设置 Deadlines 时，将采用默认的 DEADLINE_EXCEEDED（这个时间非常大）

如果产生了阻塞等待，就会造成大量正在进行的请求都会被保留，并且所有请求都有可能达到最大超时

这会使服务面临资源耗尽的风险，例如内存，这会增加服务的延迟，或者在最坏的情况下可能导致整个进程崩溃
<!-- more -->
context.WithDeadline：会返回最终上下文截止时间。第一个形参为父上下文，第二个形参为调整的截止时间。若父级时间早于子级时间，则以父级时间为准，否则以子级时间为最终截止时间
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
     if cur, ok := parent.Deadline(); ok &amp;&amp; cur.Before(d) {
         // The current deadline is already sooner than the new one.
         return WithCancel(parent)
     }
     
 }
 
 http://www.cppcns.com/jiaoben/golang/494937.html
 
 
 
     e, _ := status.FromError(err)
    s := &spb.Status{
        Code:    int32(e.Code()),
        Message: e.Message(),
        
        
 错误的堆栈也可以跨RPC传输，但是，但是这样你只能拿到当前服务的堆栈，却不能拿到调用方的堆栈，就比如说，A服务调用B服务，当B服务发生错误时，在A服务通过日志打印错误的时候，我们只打印了B服务的调用堆栈，怎样可以把A服务的堆栈打印出来。我们在A服务调用的地方也获取一次堆栈。

func WrapError(err error) error {
    if err == nil {
        return nil
    }

    s := &spb.Status{
        Code:    int32(codes.Unknown),
        Message: err.Error(),
        Details: []*any.Any{
            {
                TypeUrl: TypeUrlStack,
                Value:   util.Str2bytes(stack()),
            },
        },
    }
    return status.FromProto(s).Err()
}
// Stack 获取堆栈信息
func stack() string {
    var pc = make([]uintptr, 20)
    n := runtime.Callers(3, pc)

    var build strings.Builder
    for i := 0; i < n; i++ {
        f := runtime.FuncForPC(pc[i] - 1)
        file, line := f.FileLine(pc[i] - 1)
        n := strings.Index(file, name)
        if n != -1 {
            s := fmt.Sprintf(" %s:%d \n", file[n:], line)
            build.WriteString(s)
        }
    }
    return build.String()
}

https://www.jianshu.com/p/592759734b8e?hmsr=toutiao.io&utm_medium=toutiao.io