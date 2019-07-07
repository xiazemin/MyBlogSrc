---
title: namespace
layout: post
category: docker
author: 夏泽民
---
从原理角度来讲：

docker创建container，说白了就是linux系统中的一次fork的调用，在fork调用的时候，会传入一些flag参数，这些参数可以控制对linux内核的调用使用新的namespace；具体的做法是docker daemon封装好一个Command类，在这个Command类中，有关于namespace的配置；接着docker daemon将这个Command类发给execdriver，execdriver从中提取出关于namespace配置，然后最终将这些配置赋给golang 的exec.Cmd类，这个类就是docker 容器跑起来以后的第一个进程；这个进程在创立的时候就会创建新的namespace；
<!-- more -->
从代码角度来讲：

docker run的启动过程，其中docker daemon 与execdriver之间的交互是通过Command类，这是个重要的类，类的定义在daemon/execdriver/driver.go中；其中主要有：Network、Ipc、Pid、UTS等几种命名空间，uid namespace 现在还不支持创建独立的，也就是docker 容器里面的root和 宿主机的root应该是一样的；而这几个命名空间是怎么放到Command里面来的呢，具体代码在daemon/container_unix.go的populateCommand()函数中，这里面主要的作用就是通过container的config和hostconfig来生成Command中的关于这几个namespace的配置。

 

现在对于docker daemon来说，它已经生成好了发给execdriver的Command类，那么在execdriver中是怎么使用的呢。

在 daemon/execdriver/native/create.go中

func (d *Driver) createContainer(c *execdriver.Command) (*configs.Config, error)

这个函数的作用是接收Command类，根据里面的配置返回一个libcontainer/configs.Config类， 这个类就是execdriver去实际调用去执行关于生成容器的系统调用的时候的配置；

里面有四个这样的函数，就是来从execdriver.Command类中来提取出namespace的相关配置；

if err := d.createIpc(container, c); err != nil {
    return nil, err
}

if err := d.createPid(container, c); err != nil {
    return nil, err
}

if err := d.createUTS(container, c); err != nil {
    return nil, err
}

if err := d.createNetwork(container, c); err != nil {
    return nil, err
}

生成的这个libcontainer/configs.Config类会在vendor/src/github.com/opencontainers/runc/libcontainer/factory_linux.go的函数

func (l *LinuxFactory) Create(id string, config *configs.Config) (Container, error)  中被复制到linuxContainer类里面：

 return &linuxContainer{

        id:            id,

        root:          containerRoot,

        config:        config,   //对，就是这个config

        initPath:      l.InitPath,

        initArgs:      l.InitArgs,

        criuPath:      l.CriuPath,

        cgroupManager: l.NewCgroupsManager(config.Cgroups, nil),

    },nil

接着这个config的配置 在 /vendor/src/github.com/opencontainers/runc/libcontainer/container_linux.go 

func (c *linuxContainer) newInitProcess(p *Process, cmd *exec.Cmd, parentPipe, childPipe *os.File) (*initProcess, error) {     t := "_LIBCONTAINER_INITTYPE=standard"

    cloneFlags := c.config.Namespaces.CloneFlags()

    if cloneFlags&syscall.CLONE_NEWUSER != 0 {

        if err := c.addUidGidMappings(cmd.SysProcAttr); err != nil {

            // user mappings are not supported

            return nil, err

        }

        enableSetgroups(cmd.SysProcAttr)

        // Default to root user when user namespaces are enabled.

        if cmd.SysProcAttr.Credential == nil {

            cmd.SysProcAttr.Credential = &syscall.Credential{}

        }

    }

    cmd.Env = append(cmd.Env, t)

    cmd.SysProcAttr.Cloneflags = cloneFlags

    return &initProcess{

        cmd:        cmd,

        childPipe:  childPipe,

        parentPipe: parentPipe,

        manager:    c.cgroupManager,

        config:     c.newInitConfig(p),

    }, nil

}


