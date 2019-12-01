---
title: kubernetes
layout: post
category: golang
author: 夏泽民
---
k8s的部署架构
kubernetes中有两类资源，分别是master和nodes，master和nodes上跑的服务如下图，

kube-apiserver                |               kubelet
kube-controller-manager       |          
kube-scheduler                |               kube-proxy
----------------------                   --------------------
     k8s master                              node (non-master)
master：负责管理整个集群，例如，对应用进行调度(扩缩)、维护应用期望的状态、对应用进行发布等。
node：集群中的宿主机（可以是物理机也可以是虚拟机），每个node上都有一个agent，名为kubelet，用于跟master通信。同时一个node需要有管理容器的工具包，用于管理在node上运行的容器(docker或rkt)。一个k8s集群至少要有3个节点。
kubelet通过master暴露的API与master通信，用户也可以直接调用master的API做集群的管理。
<!-- more -->
k8s中的对象Objects
pod
k8s中的最小部署单元，不是一个程序/进程，而是一个环境(包括容器、存储、网络ip:port、容器配置)。其中可以运行1个或多个container（docker或其他容器），在一个pod内部的container共享所有资源，包括共享pod的ip:port和磁盘。
pod是临时性的，用完即丢弃的，当pod中的进程结束、node故障，或者资源短缺时，pod会被干掉。基于此，用户很少直接创建一个独立的pods，而会通过k8s中的controller来对pod进行管理。
controller通过pod templates来创建pod，pod template是一个静态模板，创建出来之后的pod就跟模板没有关系了，模板的修改也不会影响现有的pod。

services
由于pod是临时性的，pod的ip:port也是动态变化的。这种动态变化在k8s集群中就涉及到一个问题：如果一组后端pod作为服务提供方，供一组前端的pod所调用，那服务调用方怎么自动感知服务提供方。这就引入了k8s中的另外一个核心概念，services.
service是通过apiserver创建出来的对象实例，举例，

kind: Service
apiVersion: v1
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9376
这个配置将创建出来一个新的Service对象，名为my-service，后端是所有包含app=MyApp的pod，目标端口是9376，同时这个service也会被分配一个ip，被称为集群ip，对应的端口是80. 如果不指定targetPort, 那么targetPort与port相同。关于targetPort更灵活的设定是，targetPort可以是一个String类型的名字，该名字对应的真实端口值由各个后端pod自己定义，这样同一组pod无需保证同一个port，更加灵活。

在某种场景下，可能你不想用label selector，不想通过选择标签的方式来获取所有的后端pod，如测试环境下同一个service name可能指向自己的pod，而非生产环境中的正式pod集群。或者pod集群并不在k8s中维护。针对这个场景，k8s也有优雅的方案进行适配。就是在创建service的时候不指定selector的部分，而是先创建出来service，之后手动绑定后端pod的ip和端口。
终于知道为什么k8s是目前风头最劲的服务编排技术了，它充分地做了解耦，由于google的业务复杂性，它的组件和组件之间，充分的解耦、灵活，整个系统松散且牢固。
services组件与bns不同的一点，bns的节点是自己指定了name和后端的关联关系，而services是根据pod上的标签(label)自动生成的，更灵活。ali的group就更别提了，group是隶属于app的，扩展性方面更弱一些。

上文说在创建service的时候，系统为service分配了一个集群虚IP和端口，服务使用方通过这个vip:port来访问真实的服务提供方。这里的vip就是kube-proxy提供出来的。

虚ip和service proxies
kube-proxy的模式

userspace: client -> iptables -> kube-proxy -> backend pod(rr), iptables只是把虚ip转换成kube-proxy的ip，通过kube-proxy自己维护的不同端口来轮询转发到后端的pod上。
iptables: client -> iptables -> backend pod(random)，kube-proxy只是监听master上service的创建，之后动态添加/删除本机上的iptables规则
ipvs: client -> ipvs ->backend pod, ipvs是一个内核模块
服务发现
服务使用方如何找到我们定义的Service? 在k8s中用了两个方案，环境变量 && DNS。

环境变量
每当有service被创建出来之后，各个node(宿主机)上的kubelet，就会把service name加到自己宿主机的环境变量中，供所有Pod使用。环境变量的命名规则是{SERVICE_NAME}_SERVICE_HOST, ${SERVICE_NAME}SERVICE_PORT，其中SERVICE_NAME是新serviceName的大写形式，serviceName中的横杠-会被替换成下划线.
使用环境变量有一个隐含的创建顺序，即服务使用方在通过环境变量访问一个service的时候，这个service必须已经存在了。
这么简单粗暴的方案...这样做有个好处，就是省的自己搞名字解析服务，相当于本地的agent做了“域名劫持”。serviceName对应到上文提到的，由kube-proxy提供的vip:port上。

DNS
这是官方不推荐的做法，推荐用来跟k8s的外部服务进行交互，ExternalName.

Headless services
在创建service的时候，用户可以给spec.clusterIp指定一个特定的ip。前提是这个ip需要是一个合法的IP，并且要在apiServer定义的service-cluster-ip-range范围内。
当然，如果用户有自己的服务发现服务，也可以不用k8s原生的service服务，这时，需要显式地给spec.clusterIp设置None，这样k8s就不会给service分配clusterIp了。
如果用户将spec.cluster设置为None，但指定了selector，那么endpoints controller还是会创建Endpoints的，会创建一个新的DNS记录直接指向这个service描述的后端pod。否则，不会创建Endpoints记录。 // 这块没看懂

发布服务，服务类型
ClusterIp: 默认配置，给service分配一个集群内部IP，k8s集群外部不识别。
NodePort：宿主机端口，集群外部的服务也访问，使用nodeIp:nodePort
LoadBalancer：通过云负载均衡，将service暴露出去 // 没理解
ExternalName：将serviceName与externalName绑定，让外网可以访问到。要求kube-dns的版本>= 1.7
<img src="{{site.url}}{{site.baseurl}}/img/k82_na.png"/>
POD生命周期
 

需要注意的是pod的生命周期和container的生命周期有一定的联系，但是不能完全混淆一致。pod状态相对来说要简单一些。这里首先列出pod的状态

 

1、pending：pod已经被系统认可了，但是内部的container还没有创建出来。这里包含调度到node上的时间以及下载镜像的时间，会持续一小段时间。
2、Running：pod已经与node绑定了（调度成功），而且pod中所有的container已经创建出来，至少有一个容器在运行中，或者容器的进程正在启动或者重启状态。--这里需要注意pod虽然已经Running了，但是内部的container不一定完全可用。因此需要进一步检测container的状态。
3、Succeeded：这个状态很少出现，表明pod中的所有container已经成功的terminated了，而且不会再被拉起了。
4、Failed：pod中的所有容器都被terminated，至少一个container是非正常终止的。（退出的时候返回了一个非0的值或者是被系统直接终止）
5、unknown：由于某些原因pod的状态获取不到，有可能是由于通信问题。 一般情况下pod最常见的就是前两种状态。而且当Running的时候，需要进一步关注container的状态。下面就来看下container的状态有哪些：
 

Container生命周期
 

1、Waiting：启动到运行中间的一个等待状态。
2、Running：运行状态。
3、Terminated：终止状态。 如果没有任何异常的情况下，container应该会从Waiting状态变为Running状态，这时容器可用。

但如果长时间处于Waiting状态，container会有一个字段reason表明它所处的状态和原因，如果这个原因很容易能标识这个容器再也无法启动起来时，例如ContainerCannotRun，整个服务启动就会迅速返回。
<img src="{{site.url}}{{site.baseurl}}/img/architecture_k8s.png"/>

Master节点上面主要由四个模块组成：APIServer、scheduler、controller manager、etcd。

APIServer。APIServer的功能如其名，负责对外提供RESTful的Kubernetes API服务，它是系统管理指令的统一入口，任何对资源进行增删改查的操作都要交给APIServer处理后再提交给etcd。如架构图中所示，kubectl（Kubernetes提供的客户端工具，该工具内部就是对Kubernetes API的调用）是直接和APIServer交互的。
schedule。scheduler的职责很明确，就是负责调度pod（Kubernetes中最小的调度单元，后面马上就会介绍）到合适的Node上。如果把scheduler看成一个黑匣子，那么它的输入是pod和由多个Node组成的列表，输出是Pod和一个Node的绑定（bind），即将这个pod部署到这个Node上。虽然scheduler的职责很简单，但我们知道调度系统的智能程度对于整个集群是非常重要的。Kubernetes目前提供了调度算法，但是同样也保留了接口，用户可以根据自己的需求定义自己的调度算法。
controller manager。如果说APIServer做的是“前台”的工作的话，那controller manager就是负责“后台”的。后面我们马上会介绍到资源，每个资源一般都对应有一个控制器，而controller manager就是负责管理这些控制器的。还是举个例子来说明吧：比如我们通过APIServer创建一个pod，当这个pod创建成功后，APIServer的任务就算完成了。而后面保证Pod的状态始终和我们预期的一样的重任就由controller manager去保证了。
etcd。etcd是一个高可用的键值存储系统，Kubernetes使用它来存储各个资源的状态，从而实现了Restful的API。
至此，Kubernetes Master就简单介绍完了。当然，每个模块内部的实现都很复杂，而且功能也比较复杂，我现在也只是比较浅的了解了一下。如果后续了解的比较清楚了，再做总结分享。

Node
真正干活的来了。每个Node节点主要由三个模块组成：kubelet、kube-proxy、runtime。先从简单的说吧。

runtime。runtime指的是容器运行环境，目前Kubernetes支持docker和rkt两种容器，一般都指的是docker，毕竟docker现在是最主流的容器。
kube-proxy。该模块实现了Kubernetes中的服务发现和反向代理功能。反向代理方面：kube-proxy支持TCP和UDP连接转发，默认基于Round Robin算法将客户端流量转发到与service对应的一组后端pod。服务发现方面，kube-proxy使用etcd的watch机制，监控集群中service和endpoint对象数据的动态变化，并且维护一个service到endpoint的映射关系，从而保证了后端pod的IP变化不会对访问者造成影响。另外kube-proxy还支持session affinity。（这里涉及到了service的概念，可以先跳过，等了解了service之后再来看。）
kubelet。Kubelet是Master在每个Node节点上面的agent，是Node节点上面最重要的模块，它负责维护和管理该Node上面的所有容器，但是如果容器不是通过Kubernetes创建的，它并不会管理。本质上，它负责使Pod得运行状态与期望的状态一致。
至此，Kubernetes的Master和Node就简单介绍完了。下面我们来看Kubernetes中的各种资源/对象。

Kubernetes资源/对象
当然上面已经介绍的Node也算Kubernetes的资源，这里就不再赘述了。

Pod
Pod是Kubernetes里面抽象出来的一个概念，它具有如下特点：

Pod是能够被创建、调度和管理的最小单元；
每个Pod都有一个独立的IP；
一个Pod由一个或多个容器构成；
一个Pod内的容器共享Pod的所有资源，这些资源主要包括：共享存储（以Volumes的形式）、共享网络、共享端口等。
集群内的Pod之间不论是否在同一个Node上都可以任意访问，这一般是通过一个二层网络来实现的。
Kubernetes虽然也是一个容器编排系统，但不同于其他系统，它的最小操作单元不是单个容器，而是Pod。这个特性给Kubernetes带来了很多优势，比如最显而易见的是同一个Pod内的容器可以非常方便的互相访问（通过localhost就可以访问）和共享数据。使用Pod时我们需要注意两点：

虽然Pod内可以有多个容器，但一般我们只将有亲密关系的容器放在一个Pod内，比如这些容器需要相互访问、共享数据等。举个最典型的例子，比如我们有一个系统，前端是tomcat作为web，后端是存储数据的数据库mysql，那么将这两个容器部署在一个Pod内就非常合理了，因为他们通过localhost就可以访问彼此。
虽然每个Pod都有独立的IP，但是不推荐前台通过IP去访问Pod，因为Pod一旦销毁重建，IP就会变化。如果我们的Pod提供了对外的Web服务，那么我们可以通过Kubernetes提供的service去访问，后面会介绍到。
下面是一个Pod的描述文件nginx-pod.yaml：

apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
    - image: registry.hnaresearch.com/library/nginx:latest
      name:  nginx
      ports:
      - containerPort: 80
apiVersion表示API的版本，kind表示我们要创建的资源的类型。metadata是Pod的一些元数据描述。spec描述了我们期望该Pod内运行的容器。通过kubectl create -f nginx-pod.yaml就可以创建一个Pod，这个Pod里面只有一个nginx容器。

➜  kubectl create -f nginx-pod.yaml
pod "nginx-pod" created
➜  kubectl get pod
NAME        READY     STATUS    RESTARTS   AGE
nginx-pod   1/1       Running   0          1h
这里我们只是为了示例，其实实际应用中我们很少会去直接创建一个Pod，因为这样创建的Pod就像个没人管的孩子，挂了的话也不会有人去重新建立一个新的来顶替它。Kubernetes提供了很多创建Pod的方式，下面我们接着介绍。

Replication Controller
Replication Controller简称RC，一般翻译为副本控制器，这里的副本指的是Pod。如它的名字所言RC的作用就是保证任意时刻集群中都有期望个数的Pod副本在正常运行。我们通过一个简单的RC描述文件（mysql-rc.yaml）来介绍它：

apiVersion: v1
kind: ReplicationController
metadata:
  name: mysql
spec:
  replicas: 1
  selector:
    app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: registry.hnaresearch.com/library/mysql:5.6
        ports:
        - containerPort: 3306
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "123456"
上面这个文件描述了一个RC，名字叫mysql，最上面的spec描述了我们期望有1个副本，这些副本都是按照下面的template去创建的。如果某一时刻副本数比replicas描述的少，就按照template去创建新的，如果多了，就干掉几个。而下面的spec描述了这个Pod内运行的容器。

➜  kubectl create -f mysql-rc.yaml
replicationcontroller "mysql" created
➜  kubectl get rc
NAME      DESIRED   CURRENT   READY     AGE
mysql     1         1         1         7s
➜  kubectl get pod
NAME          READY     STATUS    RESTARTS   AGE
mysql-1l717   1/1       Running   0          27s
nginx-pod     1/1       Running   0          1h
然后我们进行一些删除操作：

➜  kubectl delete pod nginx-pod
pod "nginx-pod" deleted
➜  kubectl get pod
NAME          READY     STATUS    RESTARTS   AGE
mysql-1l717   1/1       Running   0          5m
➜  kubectl delete pod mysql-1l717
pod "mysql-1l717" deleted
➜  kubectl get pod
NAME          READY     STATUS    RESTARTS   AGE
mysql-2vl9k   1/1       Running   0          4s
我们先删掉之前通过Pod描述文件创建的nginx-pod，按照预期它被删除了，并没有重建。然后我们删掉mysql-1l717，发现又出来一个新的mysql-2vl9k。这是因为mysql这个是通过RC创建的，除非删除它的RC，否则任意时刻该RC都会保证有预期个Pod副本在运行。

那么，RC是怎么和Pod产生关联的呢？上面的selector是什么含义？OK，我们来介绍下一个对象。

Labels和Selector
Label是附属在Kubernetes对象上的键值对，它可以在创建的时候指定，也可以随时增删改。一个对象上面可以有任意多个Labels。它往往对于用户是有意义的，对系统是没有特殊含义的。我个人理解你可以简单把他当做Git上面的tag。这里我们只介绍一下它和Selector配合使用时的场景。我们从上面Pod和RC的描述文件中可以看到，每个Pod都有一个Labels，而RC的Selector部分也有一个定义了一个labels。RC会认为凡是和它Selector部分定义的labels相同的Pod都是它预期的副本。比如凡是labels为app=mysql的Pod都是刚才定义的RC的副本。

所以就有一个注意点，我们不要单独去创建Pod，更不要单独去创建符合某个RC的Selector的Pod，那样RC会认为是它自己创建的这个Pod而导致与预期Pod数不一致而干掉某些Pod。当然Labels还有很多用途，Selector除了等值外也有一些其他判读规则，这里不细述。

Service
终于轮到Service出场了，之前我们已经多次提到它了。Service是Kubernetes里面抽象出来的一层，它定义了由多个Pods组成的逻辑组（logical set），可以对组内的Pod做一些事情：

对外暴露流量
做负载均衡（load balancing）
服务发现（service-discovery）。
前面我们说了如果我们想将Pod内容器提供的服务暴露出去，就要使用Service。因为Service除了上面的特性外，还有一个集群内唯一的私有IP和对外的端口，用于接收流量。如果我们想将一个Service暴露到集群外，即集群外也可以访问的话，有两种方法：

LoadBalancer - 提供一个公网的IP
NodePort - 使用NAT将Service的端口暴露出去。

为什么不能通过Pod的IP，而要通过Service呢？因为在Kubernetes中，Pod是可能随时死掉被重建的，所以说其IP是不可靠的。但是Service一旦创建，其IP就会一直固定直到这个Service消亡。其实我们能够看到，Kubernetes中一个Service就相当于一个微服务。这里我们就不细述Service的创建方法以及如何使用LB以及NodePort了。

Replica Sets和Deployment
Replica Sets被称为下一代的Replication Controller，它被设计出来的目的是替代RC并提供更多的功能。就目前看，ReplicaSet和RC的唯一区别是对于Labels和Selector的支持。RC只支持等值的方式，而ReplicaSet还支持集合的选择方式（In，Not In）。另外，ReplicaSet很少像RC一样单独使用（当然，它可以单独使用），一般都是配合Deployment一起使用。

Deployment也是Kubernetes新增加的一种资源，从它的名字就可以看出它主要是为部署而设计的，之前的文章中已经有具体的例子了。想像一下我们利用RC创建了一些Pod，但现在我们想要更新Pod内容器使用的镜像或者想更改副本的个数等。这些我们无法通过修改已有的RC去做，只能删除旧的，创建新的。但这样Pod内的容器就会停止，也即业务就会中断，这在生产环境中往往是不可接受的。但有了Deployment以后，这些问题就都可以解决了。通过Deployment我们可以动态的控制副本个数、ReplicaSet和Pod的状态、滚动升级等。Deployment的强大真的需要很长的一篇文章来介绍，后续的博客再介绍吧。

HPA
HPA全称Horizontal Pod Autoscaling，即Pod的水平自动扩展，我觉得这个简直就是Kubernetes的黑科技。因为它可以根据当前系统的负载来自动水平扩容，如果系统负载超过预定值，就开始增加Pod的个数，如果低于某个值，就自动减少Pod的个数。因为被以前的系统扩容缩容深深的折磨过，所以我觉得这个功能是多么的强大。当然，目前Kubernetes的HPA只能根据CPU和内存去度量系统的负载，而且目前还依赖heapster去收集CPU的使用情况，所以功能还不是很强大，但是其完善也只是时间的问题了。让我们期待吧。

Namespace
Kubernetes提供了Namespace来从逻辑上支持多租户的功能，默认有两个Namespace：

default。用户默认的namespace。
kube-system。系统创建的对象都在这个namespace下。
当然我们可以自己创建新的namespace。

Daemon Set
有时我们可能有这样的需求，需要在所有Pod上面（包括将来新创建的）都运行某个容器，比如用于监控、日志收集等。那我们就可以使用DaemonSet，它可以保证所有容器上面都运行一份我们指定的容器的实例。而且，通过Labels和Selector，我们可以实现只在某些Pod上面部署，非常的灵活。