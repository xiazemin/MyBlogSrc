---
title: Zookeeper与Paxos
layout: post
category: spark
author: 夏泽民
---
　Zookeeper是一个开源的分布式协调服务，其设计目标是将那些复杂的且容易出错的分布式一致性服务封装起来，构成一个高效可靠的原语集，并以一些列简单的接口提供给用户使用。其是一个典型的分布式数据一致性的解决方案，分布式应用程序可以基于它实现诸如数据发布/发布、负载均衡、命名服务、分布式协调/通知、集群管理、Master选举、分布式锁和分布式队列等功能。其可以保证如下分布式一致性特性。

　　① 顺序一致性，从同一个客户端发起的事务请求，最终将会严格地按照其发起顺序被应用到Zookeeper中去。

　　② 原子性，所有事务请求的处理结果在整个集群中所有机器上的应用情况是一致的，即整个集群要么都成功应用了某个事务，要么都没有应用。

　　③ 单一视图，无论客户端连接的是哪个Zookeeper服务器，其看到的服务端数据模型都是一致的。

　　④ 可靠性，一旦服务端成功地应用了一个事务，并完成对客户端的响应，那么该事务所引起的服务端状态变更将会一直被保留，除非有另一个事务对其进行了变更。

　　⑤ 实时性，Zookeeper保证在一定的时间段内，客户端最终一定能够从服务端上读取到最新的数据状态。

　　2.1 设计目标

　　Zookeeper致力于提供一个高性能、高可用、且具有严格的顺序访问控制能力（主要是写操作的严格顺序性）的分布式协调服务，其具有如下的设计目标。

　　① 简单的数据模型，Zookeeper使得分布式程序能够通过一个共享的树形结构的名字空间来进行相互协调，即Zookeeper服务器内存中的数据模型由一系列被称为ZNode的数据节点组成，Zookeeper将全量的数据存储在内存中，以此来提高服务器吞吐、减少延迟的目的。

　　② 可构建集群，一个Zookeeper集群通常由一组机器构成，组成Zookeeper集群的而每台机器都会在内存中维护当前服务器状态，并且每台机器之间都相互通信。

　　③ 顺序访问，对于来自客户端的每个更新请求，Zookeeper都会分配一个全局唯一的递增编号，这个编号反映了所有事务操作的先后顺序。

　　④ 高性能，Zookeeper将全量数据存储在内存中，并直接服务于客户端的所有非事务请求，因此它尤其适用于以读操作为主的应用场景。

　　2.2 基本概念

　　① 集群角色，最典型的集群就是Master/Slave模式（主备模式），此情况下把所有能够处理写操作的机器称为Master机器，把所有通过异步复制方式获取最新数据，并提供读服务的机器为Slave机器。Zookeeper引入了Leader、Follower、Observer三种角色，Zookeeper集群中的所有机器通过Leaser选举过程来选定一台被称为Leader的机器，Leader服务器为客户端提供写服务，Follower和Observer提供读服务，但是Observer不参与Leader选举过程，不参与写操作的过半写成功策略，Observer可以在不影响写性能的情况下提升集群的性能。

　　② 会话，指客户端会话，一个客户端连接是指客户端和服务端之间的一个TCP长连接，Zookeeper对外的服务端口默认为2181，客户端启动的时候，首先会与服务器建立一个TCP连接，从第一次连接建立开始，客户端会话的生命周期也开始了，通过这个连接，客户端能够心跳检测与服务器保持有效的会话，也能够向Zookeeper服务器发送请求并接受响应，同时还能够通过该连接接受来自服务器的Watch事件通知。

　　③ 数据节点，第一类指构成集群的机器，称为机器节点，第二类是指数据模型中的数据单元，称为数据节点-Znode，Zookeeper将所有数据存储在内存中，数据模型是一棵树，由斜杠/进行分割的路径，就是一个ZNode，如/foo/path1，每个ZNode都会保存自己的数据内存，同时还会保存一些列属性信息。ZNode分为持久节点和临时节点两类，持久节点是指一旦这个ZNode被创建了，除非主动进行ZNode的移除操作，否则这个ZNode将一直保存在Zookeeper上，而临时节点的生命周期和客户端会话绑定，一旦客户端会话失效，那么这个客户端创建的所有临时节点都会被移除。另外，Zookeeper还允许用户为每个节点添加一个特殊的属性：SEQUENTIAL。一旦节点被标记上这个属性，那么在这个节点被创建的时候，Zookeeper会自动在其节点后面追加一个整形数字，其是由父节点维护的自增数字。

　　④ 版本，对于每个ZNode，Zookeeper都会为其维护一个叫作Stat的数据结构，Stat记录了这个ZNode的三个数据版本，分别是version（当前ZNode的版本）、cversion（当前ZNode子节点的版本）、aversion（当前ZNode的ACL版本）。

　　⑤ Watcher，Zookeeper允许用户在指定节点上注册一些Watcher，并且在一些特定事件触发的时候，Zookeeper服务端会将事件通知到感兴趣的客户端。

　　⑥ ACL，Zookeeper采用ACL（Access Control Lists）策略来进行权限控制，其定义了如下五种权限：

　　　　· CREATE：创建子节点的权限。

　　　　· READ：获取节点数据和子节点列表的权限。

　　　　· WRITE：更新节点数据的权限。

　　　　· DELETE：删除子节点的权限。

　　　　· ADMIN：设置节点ACL的权限。

　　2.3 ZAB协议

　　Zookeeper使用了Zookeeper Atomic Broadcast（ZAB，Zookeeper原子消息广播协议）的协议作为其数据一致性的核心算法。ZAB协议是为Zookeeper专门设计的一种支持崩溃恢复的原子广播协议。

　　Zookeeper依赖ZAB协议来实现分布式数据的一致性，基于该协议，Zookeeper实现了一种主备模式的系统架构来保持集群中各副本之间的数据的一致性，即其使用一个单一的诸进程来接收并处理客户端的所有事务请求，并采用ZAB的原子广播协议，将服务器数据的状态变更以事务Proposal的形式广播到所有的副本进程中，ZAB协议的主备模型架构保证了同一时刻集群中只能够有一个主进程来广播服务器的状态变更，因此能够很好地处理客户端大量的并发请求。

　　ZAB协议的核心是定义了对于那些会改变Zookeeper服务器数据状态的事务请求的处理方式，即：所有事务请求必须由一个全局唯一的服务器来协调处理，这样的服务器被称为Leader服务器，余下的服务器则称为Follower服务器，Leader服务器负责将一个客户端事务请求转化成一个事务Proposal（提议），并将该Proposal分发给集群中所有的Follower服务器，之后Leader服务器需要等待所有Follower服务器的反馈，一旦超过半数的Follower服务器进行了正确的反馈后，那么Leader就会再次向所有的Follower服务器分发Commit消息，要求其将前一个Proposal进行提交。

　　ZAB一些包括两种基本的模式：崩溃恢复和消息广播。

　　当整个服务框架启动过程中或Leader服务器出现网络中断、崩溃退出与重启等异常情况时，ZAB协议就会进入恢复模式并选举产生新的Leader服务器。当选举产生了新的Leader服务器，同时集群中已经有过半的机器与该Leader服务器完成了状态同步之后，ZAB协议就会退出恢复模式，状态同步时指数据同步，用来保证集群在过半的机器能够和Leader服务器的数据状态保持一致。

　　当集群中已经有过半的Follower服务器完成了和Leader服务器的状态同步，那么整个服务框架就可以进入消息广播模式，当一台同样遵守ZAB协议的服务器启动后加入到集群中，如果此时集群中已经存在一个Leader服务器在负责进行消息广播，那么加入的服务器就会自觉地进入数据恢复模式：找到Leader所在的服务器，并与其进行数据同步，然后一起参与到消息广播流程中去。Zookeeper只允许唯一的一个Leader服务器来进行事务请求的处理，Leader服务器在接收到客户端的事务请求后，会生成对应的事务提议并发起一轮广播协议，而如果集群中的其他机器收到客户端的事务请求后，那么这些非Leader服务器会首先将这个事务请求转发给Leader服务器。

　　当Leader服务器出现崩溃或者机器重启、集群中已经不存在过半的服务器与Leader服务器保持正常通信时，那么在重新开始新的一轮的原子广播事务操作之前，所有进程首先会使用崩溃恢复协议来使彼此到达一致状态，于是整个ZAB流程就会从消息广播模式进入到崩溃恢复模式。一个机器要成为新的Leader，必须获得过半机器的支持，同时由于每个机器都有可能会崩溃，因此，ZAB协议运行过程中，前后会出现多个Leader，并且每台机器也有可能会多次成为Leader，进入崩溃恢复模式后，只要集群中存在过半的服务器能够彼此进行正常通信，那么就可以产生一个新的Leader并再次进入消息广播模式。如一个由三台机器组成的ZAB服务，通常由一个Leader、2个Follower服务器组成，某一个时刻，加入其中一个Follower挂了，整个ZAB集群是不会中断服务的。

　　① 消息广播，ZAB协议的消息广播过程使用原子广播协议，类似于一个二阶段提交过程，针对客户端的事务请求，Leader服务器会为其生成对应的事务Proposal，并将其发送给集群中其余所有的机器，然后再分别收集各自的选票，最后进行事务提交。



　　在ZAB的二阶段提交过程中，移除了中断逻辑，所有的Follower服务器要么正常反馈Leader提出的事务Proposal，要么就抛弃Leader服务器，同时，ZAB协议将二阶段提交中的中断逻辑移除意味着我们可以在过半的Follower服务器已经反馈Ack之后就开始提交事务Proposal，而不需要等待集群中所有的Follower服务器都反馈响应，但是，在这种简化的二阶段提交模型下，无法处理Leader服务器崩溃退出而带来的数据不一致问题，因此ZAB采用了崩溃恢复模式来解决此问题，另外，整个消息广播协议是基于具有FIFO特性的TCP协议来进行网络通信的，因此能够很容易保证消息广播过程中消息接受与发送的顺序性。再整个消息广播过程中，Leader服务器会为每个事务请求生成对应的Proposal来进行广播，并且在广播事务Proposal之前，Leader服务器会首先为这个事务Proposal分配一个全局单调递增的唯一ID，称之为事务ID（ZXID），由于ZAB协议需要保证每个消息严格的因果关系，因此必须将每个事务Proposal按照其ZXID的先后顺序来进行排序和处理。

　　② 崩溃恢复，在Leader服务器出现崩溃，或者由于网络原因导致Leader服务器失去了与过半Follower的联系，那么就会进入崩溃恢复模式，在ZAB协议中，为了保证程序的正确运行，整个恢复过程结束后需要选举出一个新的Leader服务器，因此，ZAB协议需要一个高效且可靠的Leader选举算法，从而保证能够快速地选举出新的Leader，同时，Leader选举算法不仅仅需要让Leader自身知道已经被选举为Leader，同时还需要让集群中的所有其他机器也能够快速地感知到选举产生的新的Leader服务器。

　　③ 基本特性，ZAB协议规定了如果一个事务Proposal在一台机器上被处理成功，那么应该在所有的机器上都被处理成功，哪怕机器出现故障崩溃。ZAB协议需要确保那些已经在Leader服务器上提交的事务最终被所有服务器都提交，假设一个事务在Leader服务器上被提交了，并且已经得到了过半Follower服务器的Ack反馈，但是在它Commit消息发送给所有Follower机器之前，Leader服务挂了。如下图所示



　　在集群正常运行过程中的某一个时刻，Server1是Leader服务器，其先后广播了P1、P2、C1、P3、C2（C2是Commit Of Proposal2的缩写），其中，当Leader服务器发出C2后就立即崩溃退出了，针对这种情况，ZAB协议就需要确保事务Proposal2最终能够在所有的服务器上都被提交成功，否则将出现不一致。

　　ZAB协议需要确保丢弃那些只在Leader服务器上被提出的事务。如果在崩溃恢复过程中出现一个需要被丢弃的提议，那么在崩溃恢复结束后需要跳过该事务Proposal，如下图所示

　　假设初始的Leader服务器Server1在提出一个事务Proposal3之后就崩溃退出了，从而导致集群中的其他服务器都没有收到这个事务Proposal，于是，当Server1恢复过来再次加入到集群中的时候，ZAB协议需要确保丢弃Proposal3这个事务。

　　在上述的崩溃恢复过程中需要处理的特殊情况，就决定了ZAB协议必须设计这样的Leader选举算法：能够确保提交已经被Leader提交的事务的Proposal，同时丢弃已经被跳过的事务Proposal。如果让Leader选举算法能够保证新选举出来的Leader服务器拥有集群中所有机器最高编号（ZXID最大）的事务Proposal，那么就可以保证这个新选举出来的Leader一定具有所有已经提交的提议，更为重要的是如果让具有最高编号事务的Proposal机器称为Leader，就可以省去Leader服务器查询Proposal的提交和丢弃工作这一步骤了。

　　④ 数据同步，完成Leader选举后，在正式开始工作前，Leader服务器首先会确认日志中的所有Proposal是否都已经被集群中的过半机器提交了，即是否完成了数据同步。Leader服务器需要确所有的Follower服务器都能够接收到每一条事务Proposal，并且能够正确地将所有已经提交了的事务Proposal应用到内存数据库中。Leader服务器会为每个Follower服务器维护一个队列，并将那些没有被各Follower服务器同步的事务以Proposal消息的形式逐个发送给Follower服务器，并在每一个Proposal消息后面紧接着再发送一个Commit消息，以表示该事务已经被提交，等到Follower服务器将所有其尚未同步的事务Proposal都从Leader服务器上同步过来并成功应用到本地数据库后，Leader服务器就会将该Follower服务器加入到真正的可用Follower列表并开始之后的其他流程。

　　下面分析ZAB协议如何处理需要丢弃的事务Proposal的，ZXID是一个64位的数字，其中32位可以看做是一个简单的单调递增的计数器，针对客户端的每一个事务请求，Leader服务器在产生一个新的事务Proposal时，都会对该计数器进行加1操作，而高32位则代表了Leader周期epoch的编号，每当选举产生一个新的Leader时，就会从这个Leader上取出其本地日志中最大事务Proposal的ZXID，并解析出epoch值，然后加1，之后以该编号作为新的epoch，低32位则置为0来开始生成新的ZXID，ZAB协议通过epoch号来区分Leader周期变化的策略，能够有效地避免不同的Leader服务器错误地使用不同的ZXID编号提出不一样的事务Proposal的异常情况。当一个包含了上一个Leader周期中尚未提交过的事务Proposal的服务器启动时，其肯定无法成为Leader，因为当前集群中一定包含了一个Quorum（过半）集合，该集合中的机器一定包含了更高epoch的事务的Proposal，因此这台机器的事务Proposal并非最高，也就无法成为Leader。

　　2.4 ZAB协议原理

　　ZAB主要包括消息广播和崩溃恢复两个过程，进一步可以分为三个阶段，分别是发现（Discovery）、同步（Synchronization）、广播（Broadcast）阶段。ZAB的每一个分布式进程会循环执行这三个阶段，称为主进程周期。

　　· 发现，选举产生PL(prospective leader)，PL收集Follower epoch(cepoch)，根据Follower的反馈，PL产生newepoch(每次选举产生新Leader的同时产生新epoch)。

　　· 同步，PL补齐相比Follower多数派缺失的状态、之后各Follower再补齐相比PL缺失的状态，PL和Follower完成状态同步后PL变为正式Leader(established leader)。

　　· 广播，Leader处理客户端的写操作，并将状态变更广播至Follower，Follower多数派通过之后Leader发起将状态变更落地(deliver/commit)。

　　在正常运行过程中，ZAB协议会一直运行于阶段三来反复进行消息广播流程，如果出现崩溃或其他原因导致Leader缺失，那么此时ZAB协议会再次进入发现阶段，选举新的Leader。

　　2.4.1 运行分析

　　每个进程都有可能处于如下三种状态之一

　　· LOOKING：Leader选举阶段。

　　· FOLLOWING：Follower服务器和Leader服务器保持同步状态。

　　· LEADING：Leader服务器作为主进程领导状态。

　　所有进程初始状态都是LOOKING状态，此时不存在Leader，此时，进程会试图选举出一个新的Leader，之后，如果进程发现已经选举出新的Leader了，那么它就会切换到FOLLOWING状态，并开始和Leader保持同步，处于FOLLOWING状态的进程称为Follower，LEADING状态的进程称为Leader，当Leader崩溃或放弃领导地位时，其余的Follower进程就会转换到LOOKING状态开始新一轮的Leader选举。

　　一个Follower只能和一个Leader保持同步，Leader进程和所有与所有的Follower进程之间都通过心跳检测机制来感知彼此的情况。若Leader能够在超时时间内正常收到心跳检测，那么Follower就会一直与该Leader保持连接，而如果在指定时间内Leader无法从过半的Follower进程那里接收到心跳检测，或者TCP连接断开，那么Leader会放弃当前周期的领导，比你转换到LOOKING状态，其他的Follower也会选择放弃这个Leader，同时转换到LOOKING状态，之后会进行新一轮的Leader选举，并在选举产生新的Leader之后开始新的一轮主进程周期。

　　2.5 ZAB与Paxos的联系和区别

　　联系：

　　① 都存在一个类似于Leader进程的角色，由其负责协调多个Follower进程的运行。

　　② Leader进程都会等待超过半数的Follower做出正确的反馈后，才会将一个提议进行提交。

　　③ 在ZAB协议中，每个Proposal中都包含了一个epoch值，用来代表当前的Leader周期，在Paxos算法中，同样存在这样的一个标识，名字为Ballot。

　　区别：

　　Paxos算法中，新选举产生的主进程会进行两个阶段的工作，第一阶段称为读阶段，新的主进程和其他进程通信来收集主进程提出的提议，并将它们提交。第二阶段称为写阶段，当前主进程开始提出自己的提议。

　　ZAB协议在Paxos基础上添加了同步阶段，此时，新的Leader会确保存在过半的Follower已经提交了之前的Leader周期中的所有事务Proposal。

　　ZAB协议主要用于构建一个高可用的分布式数据主备系统，而Paxos算法则用于构建一个分布式的一致性状态机系统。
<!-- more -->
Apache Zookeeper是由Apache Hadoop的子项目发展而来，于2010年11月正式成为Apache顶级项目。Zookeeper为分布式应用提供高效且可靠的分布式协调服务，提供了统一命名服务、配置管理、分布式锁等分布式的基础服务。Zookeeper并没有直接采用Paxos算法，而是采用了一种被称为ZAB(Zookeeper Atomic Broadcast)的一致性协议。

初识Zookeeper
Zookeeper是一个开源的分布式协调服务，由Yahoo创建，是Google Chubby的开源实现。Zookeeper将那些复杂且容易出错的分布式一致性服务封装起来，构成一个高效可靠的原语集，并以一系列简单易用的接口提供给用户使用。

分布式应用程序可以基于Zookeeper实现例如数据发布/订阅、负载均衡、命名服务、协调通知、集群管理、Master选举、分布式锁、分布式队列等功能。Zookeeper可以保证如下分布式一致性特性。

1、顺序一致性：从同一个客户端发起的事务请求，最终将会严格地按照其发起顺序被应用到Zookeeper中；

2、原子性：所有事务的请求结果在整个集群中所有机器上的应用情况是一致的，也就是说，要么在整个集群中所有机器上都成功应用了某一个事务，要么都没有应用，没有中间状态；

3、单一视图：无论客户端连接的是哪个Zookeeper服务器，其看到的服务端数据模型都是一致的。

4、实时性：Zookeeper仅仅保证在一定的时间内，客户端最终一定能够从服务端上读到最终的数据状态。

Zookeeper的设计目标
目标一：简单的数据模型

Zookeeper使得分布式程序能通过一个共享式的、树形结构的名字空间来进行相互协调。

树形结构的名字空间是指Zookeeper服务器内存中的一个数据模型，由一系列被称为ZNode的数据节点组成，类似一个文件系统。

目标二：可以构建集群

组成Zookeeper集群的每台机器都会在内存中维护当前服务器状态，并且每台机器都保持通讯。只要集群中超过一半机器能够正常工作，整个集群就能对外正常服务。

目标三：顺序访问

对于来自客户端的每个更新请求，Zookeeper都会分配一个全局唯一递增编号，这个编号反映了所有事物操作的先后顺序，应用程序可以使用Zookeeper的这个特性来实现更加高层的同步原语。

目标四：高性能

Zookeeper将全量数据存储在内存中直接对外服务(除事务请求)，尤其适合于读操作为主的应用场景。

Zookeeper基本概念


集群角色
最典型的集群角色是Master/Slave模式，Master处理所有写操作，Slave提供读服务。

在Zookeeper中，有Leader、Follower、Observer三种角色。集群中所有机器通过一个Leader选举来决定一台机器作为Leader，Leader为客户端提供读和写服务。Follower和Observer都提供读服务，区别在于Observer机器不参与选举，也不参与写操作的“过半写成功”策略，因此Observe可以在不影响写性能的情况下提升集群的读性能。

会话
客户端与Zookeeper是TCP长连接，默认对外端口是2181，通过这个连接，客户端保持和服务器的心跳以维护连接，也能向Zookeeper发送请求并响应，同时还可以接收到注册通知。

数据节点(ZNode)
Zookeeper所有数据都在内存中，模型类似一颗文件树，ZNode Tree，每个ZNode节点都会保存自己的数据内容和一系列属性。

ZNode分为持久节点和临时节点，后者和客户端会话绑定。

版本
每个ZNode，Zookeeper都会维护一个Stat数据结构记录这个ZNode的三个数据版本：当前ZNode版本version、当前ZNode子节点版本cversion、和当前Node的ACL版本aversion。

Watcher
事件监听是Zookeeper的重要特性，他允许客户端在指定节点注册Watcher，并且在事件被触发后通知客户端。

ACL
Access Control Lists，定义了5中权限：

Create、Read、Write、Delete、Admin

ZAB协议
在ZooKeeper中所有的事务请求都由一个主服务器也就是Leader来处理，其他服务器为Follower，Leader将客户端的事务请求转换为事务Proposal，并且将Proposal分发给集群中其他所有的Follower，然后Leader等待Follwer反馈，当有过半数（>=N/2+1）的Follower反馈信息后，Leader将再次向集群内Follower广播Commit信息，Commit为将之前的Proposal提交。

ZAB协议中存在着三种状态，每个节点都属于以下三种中的一种：

1.Looking：系统刚启动时或者Leader崩溃后正处于选举状态

2.Following：Follower节点所处的状态，Follower与Leader处于数据同步阶段；

3.Leading：Leader所处状态，当前集群中有一个Leader为主进程；

ZooKeeper启动时所有节点初始状态为Looking，这时集群会尝试选举出一个Leader节点，选举出的Leader节点切换为Leading状态；当节点发现集群中已经选举出Leader则该节点会切换到Following状态，然后和Leader节点保持同步；当Follower节点与Leader失去联系时Follower节点则会切换到Looking状态，开始新一轮选举；在ZooKeeper的整个生命周期中每个节点都会在Looking、Following、Leading状态间不断转换；


选举出Leader节点后ZAB进入原子广播阶段，这时Leader为和自己同步的每个节点Follower创建一个操作序列，一个时期一个Follower只能和一个Leader保持同步，Leader节点与Follower节点使用心跳检测来感知对方的存在；当Leader节点在超时时间内收到来自Follower的心跳检测那Follower节点会一直与该节点保持连接；若超时时间内Leader没有接收到来自过半Follower节点的心跳检测或TCP连接断开，那Leader会结束当前周期的领导，切换到Looking状态，所有Follower节点也会放弃该Leader节点切换到Looking状态，然后开始新一轮选举。

ZAB协议定义了选举（election）、发现（discovery）、同步（sync）、广播(Broadcast)四个阶段；ZAB选举（election）时当Follower存在ZXID（事务ID）时判断所有Follower节点的事务日志，只有lastZXID的节点才有资格成为Leader，这种情况下选举出来的Leader总有最新的事务日志，基于这个原因所以ZooKeeper实现的时候把发现（discovery）与同步（sync）合并为恢复（recovery）阶段；

1.Election：在Looking状态中选举出Leader节点，Leader的lastZXID总是最新的；

2.Discovery：Follower节点向准Leader推送FOllOWERINFO，该信息中包含了上一周期的epoch，接受准Leader的NEWLEADER指令，检查newEpoch有效性，准Leader要确保Follower的epoch与ZXID小于或等于自身的；

3.sync：将Follower与Leader的数据进行同步，由Leader发起同步指令，最总保持集群数据的一致性；

4.Broadcast：Leader广播Proposal与Commit，Follower接受Proposal与Commit；

5.Recovery：在Election阶段选举出Leader后本阶段主要工作就是进行数据的同步，使Leader具有highestZXID，集群保持数据的一致性；

选举（Election）

election阶段必须确保选出的Leader具有highestZXID，否则在Recovery阶段没法保证数据的一致性，Recovery阶段Leader要求Follower向自己同步数据没有Follower要求Leader保持数据同步，所有选举出来的Leader要具有最新的ZXID；

在选举的过程中会对每个Follower节点的ZXID进行对比只有highestZXID的Follower才可能当选Leader；

选举流程：

1. 每个Follower都向其他节点发送选自身为Leader的Vote投票请求，等待回复；

2. Follower接受到的Vote如果比自身的大（ZXID更新）时则投票，并更新自身的Vote，否则拒绝投票；

3. 每个Follower中维护着一个投票记录表，当某个节点收到过半的投票时，结束投票并把该Follower选为Leader，投票结束；

ZAB协议中使用ZXID作为事务编号，ZXID为64位数字，低32位为一个递增的计数器，每一个客户端的一个事务请求时Leader产生新的事务后该计数器都会加1，高32位为Leader周期epoch编号，当新选举出一个Leader节点时Leader会取出本地日志中最大事务Proposal的ZXID解析出对应的epoch把该值加1作为新的epoch，将低32位从0开始生成新的ZXID；ZAB使用epoch来区分不同的Leader周期；

恢复（Recovery）

在election阶段选举出来的Leader已经具有最新的ZXID，所有本阶段的主要工作是根据Leader的事务日志对Follower节点数据进行更新；

Leader：Leader生成新的ZXID与epoch，接收Follower发送过来的FOllOWERINFO（含有当前节点的LastZXID）然后往Follower发送NEWLEADER；Leader根据Follower发送过来的LastZXID根据数据更新策略向Follower发送更新指令；

同步策略：

1.SNAP：如果Follower数据太老，Leader将发送快照SNAP指令给Follower同步数据；

2.DIFF：Leader发送从Follolwer.lastZXID到Leader.lastZXID议案的DIFF指令给Follower同步数据；

3.TRUNC：当Follower.lastZXID比Leader.lastZXID大时，Leader发送从Leader.lastZXID到Follower.lastZXID的TRUNC指令让Follower丢弃该段数据；

Follower：往Leader发送FOLLOERINFO指令，Leader拒绝就转到Election阶段；接收Leader的NEWLEADER指令，如果该指令中epoch比当前Follower的epoch小那么Follower转到Election阶段；Follower还有主要工作是接收SNAP/DIFF/TRUNC指令同步数据与ZXID，同步成功后回复ACKNETLEADER，然后进入下一阶段；Follower将所有事务都同步完成后Leader会把该节点添加到可用Follower列表中；

SNAP与DIFF用于保证集群中Follower节点已经Committed的数据的一致性，TRUNC用于抛弃已经被处理但是没有Committed的数据；

广播(Broadcast)

客户端提交事务请求时Leader节点为每一个请求生成一个事务Proposal，将其发送给集群中所有的Follower节点，收到过半Follower的反馈后开始对事务进行提交，ZAB协议使用了原子广播协议；在ZAB协议中只需要得到过半的Follower节点反馈Ack就可以对事务进行提交，这也导致了Leader几点崩溃后可能会出现数据不一致的情况，ZAB使用了崩溃恢复来处理数字不一致问题；消息广播使用了TCP协议进行通讯所有保证了接受和发送事务的顺序性。广播消息时Leader节点为每个事务Proposal分配一个全局递增的ZXID（事务ID），每个事务Proposal都按照ZXID顺序来处理；

Leader节点为每一个Follower节点分配一个队列按事务ZXID顺序放入到队列中，且根据队列的规则FIFO来进行事务的发送。Follower节点收到事务Proposal后会将该事务以事务日志方式写入到本地磁盘中，成功后反馈Ack消息给Leader节点，Leader在接收到过半Follower节点的Ack反馈后就会进行事务的提交，以此同时向所有的Follower节点广播Commit消息，Follower节点收到Commit后开始对事务进行提交。

总的来说，ZAB主要是用来构建一个高可用的分布式数据主备系统，而Paxos则是构建一个分布式的一致性状态机系统。

<Google的Chubby，Apache的Zookeeper都是基于它的理论来实现的，Paxos还被认为是到目前为止唯一的分布式一致性算法，其它的算法都是Paxos的改进或简化。有个问题要提一下，Paxos有一个前提：没有拜占庭将军问题。就是说Paxos只有在一个可信的计算环境中才能成立，这个环境是不会被入侵所破坏的。

”

 

Paxos 这个算法是Leslie Lamport在1990年提出的一种基于消息传递的一致性算法.Paxos 算法解决的问题是一个分布式系统如何就某个值（决议）达成一致。

part-time parliament Paxos Made Simple里这样描述Paxos算法执行过程：

prepare 阶段：
proposer 选择一个提案编号 n 并将 prepare 请求发送给 acceptors 中的一个多数派；
acceptor 收到 prepare 消息后，如果提案的编号大于它已经回复的所有 prepare 消息，则 acceptor 将自己上次接受的提案回复给 proposer，并承诺不再回复小于 n 的提案；
批准阶段：
当一个 proposer 收到了多数 acceptors 对 prepare 的回复后，就进入批准阶段。它要向回复 prepare 请求的 acceptors 发送 accept 请求，包括编号 n 和根据 P2c 决定的 value（如果根据 P2c 没有已经接受的 value，那么它可以自由决定 value）。
在不违背自己向其他 proposer 的承诺的前提下，acceptor 收到 accept 请求后即接受这个请求。
wiki上是两个阶段，论文里却是说三阶段，而且默认就有了个proposer相当于leader。查资料有大侠列出了第三个阶段(http://www.wuzesheng.com/?p=2724)：

3. Learn阶段:

当各个Acceptor达到一致之后，需要将达到一致的结果通知给所有的Learner.
 

zookeeper采用org.apache.zookeeper.server.quorum.FastLeaderElection作为其缺省选举算法

这篇文章http://csrd.aliapp.com/?p=162&replytocom=782 直接用paxos实现作为标题，提到  zookeeper在选举leader的时候采用了paxos算法(主要是fast paxos)

偶然看到下边有人反驳：

魏讲文：“

FastLeaderElection根本不是Paxos，也不是Fast Paxos的实现。
FastLeaderElection源码与Paxos的论文相去甚远。

Paxos与 FastPaxos算法中也有一个leader选举的问题。

FastLeaderElection对于zookeeper来讲，只是相当于Paxos中的leader选举。

”

 二、资料证实

好的，查查资料，分析源码开始调研

首先是魏讲文的反驳 ：

There is a very common misunderstanding that the leader election algorithm in zookeeper is paxos or fast paxos. The leader election algorithm is not paxos or fast paxos, please consider the following facts:

There is no the concept of proposal number in the leader election in zookeeper. the proposal number is a key concept to paxos. Some one think the epoch is the proposal number, but different followers may produce proposal with the same epoch which is not allowed in paxos.

Fast paxos requires at least 3t + 1 acceptors, where t is the number of servers which are allowed to fail [3]. This is conflict with the fact that a zookeeper cluster with 3 servers works well even if one server failed.

The leader election algorithm must make sure P1. However paxos does provide such guarantee.

In paxos, a leader is also required to achieve progress. There are some similarities between the leader in paxos and the leader in zookeeper. Even if more than one servers believe they are the leader, the consistency is preserved both in zookeeper and in paxos. this is very clearly discussed in [1] and [2]. 

然后是作者三次对比

1）邮件列表

Our protocol instead, has only two phases, just like a two-phase 
commit protocol. Of course, for Paxos, we can ignore the first phase in runs in 
which we have a single proposer as we can run phase 1 for multiple instances at 
a time, which is what Ben called previously Multi-Paxos, I believe. The trick 
with skipping phase 1 is to deal with leader switching. 
2）出书的访谈

We made a few interesting observations about Paxos when contrasting it to Zab, like problems you could run into if you just implemented Paxos alone. Not that Paxos is broken or anything, just that in our setting, there were some properties it was not giving us. Some people still like to map Zab to Paxos, and they are not completely off, but the way we see it, Zab matches a service like ZooKeeper well.

zk的分布式一致性算法有了个名称叫Zab

3）论文

We use an algorithm that shares some of the character- istics of Paxos, but that combines transaction logging needed for consensus with write-ahead logging needed for data tree recovery to enable an efficient implementa- tion. 

 

 三、leader选举分析

在我理解首先在选举时，并不能用到paxos算法，paxos里选总统也好，zk选leader也好，跟搞个提案让大部分人同意是有区别的。选主才能保证不会出现多proposer的并发提案冲突

谁去作为proposer发提案？是paxos算法进行下去的前提。而提出提案让大部分follower同意则可用到类似paxos的算法实现一致性。zookeeper使用的是Zab算法实现一致性。

zk的选主策略：

there can be at most one leader (proposer) at any time, and we guarantee this by making sure 
that a quorum of replicas recognize the leader as a leader by committing to an 
epoch change. This change in epoch also allows us to get unique zxids since the 
epoch forms part of the zxid. 


每个server有一个id，收到提交的事务时则有一个zxid，随更新数据的变动，事务编号递增，server id各不同。首先选zxid最大的作为leader，如果zxid比较不出来，则选server id最大的为leader

zxid包含一个epoch数字，epoch指示一个server作为leader的时期，随新的leader诞生而递增。

再看代码：

 

四、zookeeper数据更新原理分析

了解完选主的做法后，来了解下同步数据的做法，同步数据则采用Zab协议：Zookeeper Atomic broadcast protocol，是个类似两阶段提交的协议：

The leader sends a PROPOSAL message, p, to all followers.

Upon receiving p, a follower responds to the leader with an ACK, informing the

leader that it has accepted the proposal.

Uponreceivingacknowledgmentsfromaquorum(thequorumincludestheleader

itself), the leader sends a message informing the followers to COMMIT it. 

跟paxos的区别是leaer发送给所有follower，而不是大多数，所有follower都要确认并通知leader，而不是大多数。

保证机制：按顺序广播的两个事务， T 和 Tʹ ，T在前则Tʹ 生效前必须提交T。如果有一个server 提交了T 和 Tʹ ，则所有其他server必须也在Tʹ前提交T。

 

五、leader的探活

为解决leader crash的问题，避免出现多个leader导致事务混乱，Zab算法保证：

1、新事务开启时，leader必须提交上个epoch期间提交的所有事务

2、任何时候都不会有两个leader同时获得足够多的支持者。

一个新leader的起始状态需要大多数server同意

 

六、observer

zk里的第三种角色，观察者和follower的区别就是没有选举权。它主要是1、为系统的读请求扩展性存在 2、满足多机房部署需求，中心机房部署leader、follower，其他机房部署observer，读取配置优先读本地。

 

七、总结

我认为zookeeper只能说是受paxos算法影响，角色划分类似，提案通过方式类似，实现更为简单直观。选主基于voteid(server-id)和zxid做大小优先级排序，信息同步则使用两阶段提交，leader获取follower的全部同意后才提交事务，更新状态。observer角色则是为了增加系统吞吐和满足跨机房部署。

 

参考文献

[1] Reed, B., & Junqueira, F. P. (2008). A simple totally ordered broadcast protocol. Second
Workshop on Large-Scale Distributed Systems and Middleware (LADIS 2008). Yorktown
Heights, NY: ACM. ISBN: 978-1-60558-296-2.
[2] Lamport, L. Paxos made simple. ACM SIGACT News 32, 4 (Dec. 2001), 1825.
[3] F. Junqueira, Y. Mao, and K. Marzullo. Classic paxos vs. fast paxos: caveat emptor. In
Proceedings of the 3rd USENIX/IEEE/IFIP Workshop on Hot Topics in System Dependability
(HotDep.07). Citeseer, 2007.

[4]O'Reilly.ZooKeeper.Distributed process coordination.2013

[5] http://agapple.iteye.com/blog/1184023  zookeeper项目使用几点小结
1.1 基本定义
算法中的参与者主要分为三个角色，同时每个参与者又可兼领多个角色:

⑴proposer 提出提案，提案信息包括提案编号和提议的value;

⑵acceptor 收到提案后可以接受(accept)提案;

⑶learner 只能"学习"被批准的提案;

算法保重一致性的基本语义:

⑴决议(value)只有在被proposers提出后才能被批准(未经批准的决议称为"提案(proposal)");

⑵在一次Paxos算法的执行实例中，只批准(chosen)一个value;

⑶learners只能获得被批准(chosen)的value;

有上面的三个语义可演化为四个约束:

⑴P1:一个acceptor必须接受(accept)第一次收到的提案;

⑵P2a:一旦一个具有value v的提案被批准(chosen)，那么之后任何acceptor 再次接受(accept)的提案必须具有value v;

⑶P2b:一旦一个具有value v的提案被批准(chosen)，那么以后任何 proposer 提出的提案必须具有value v;

⑷P2c:如果一个编号为n的提案具有value v，那么存在一个多数派，要么他们中所有人都没有接受(accept)编号小于n的任何提案，要么他们已经接受(accpet)的所有编号小于n的提案中编号最大的那个提案具有value v;

1.2 基本算法(basic paxos)
算法(决议的提出与批准)主要分为两个阶段:

1. prepare阶段： 

(1). 当Porposer希望提出方案V1，首先发出prepare请求至大多数Acceptor。Prepare请求内容为序列号<SN1>;

(2). 当Acceptor接收到prepare请求<SN1>时，检查自身上次回复过的prepare请求<SN2>

a). 如果SN2>SN1，则忽略此请求，直接结束本次批准过程;

b). 否则检查上次批准的accept请求<SNx，Vx>，并且回复<SNx，Vx>；如果之前没有进行过批准，则简单回复<OK>;

2. accept批准阶段： 

(1a). 经过一段时间，收到一些Acceptor回复，回复可分为以下几种:

a). 回复数量满足多数派，并且所有的回复都是<OK>，则Porposer发出accept请求，请求内容为议案<SN1，V1>;

b). 回复数量满足多数派，但有的回复为：<SN2，V2>，<SN3，V3>……则Porposer找到所有回复中超过半数的那个，假设为<SNx，Vx>，则发出accept请求，请求内容为议案<SN1，Vx>;

c). 回复数量不满足多数派，Proposer尝试增加序列号为SN1+，转1继续执行;

         (1b). 经过一段时间，收到一些Acceptor回复，回复可分为以下几种:

a). 回复数量满足多数派，则确认V1被接受;

b). 回复数量不满足多数派，V1未被接受，Proposer增加序列号为SN1+，转1继续执行;

(2). 在不违背自己向其他proposer的承诺的前提下，acceptor收到accept 请求后即接受并回复这个请求。

1.3 算法优化(fast paxos)
Paxos算法在出现竞争的情况下，其收敛速度很慢，甚至可能出现活锁的情况，例如当有三个及三个以上的proposer在发送prepare请求后，很难有一个proposer收到半数以上的回复而不断地执行第一阶段的协议。因此，为了避免竞争，加快收敛的速度，在算法中引入了一个Leader这个角色，在正常情况下同时应该最多只能有一个参与者扮演Leader角色，而其它的参与者则扮演Acceptor的角色，同时所有的人又都扮演Learner的角色。

在这种优化算法中，只有Leader可以提出议案，从而避免了竞争使得算法能够快速地收敛而趋于一致，此时的paxos算法在本质上就退变为两阶段提交协议。但在异常情况下，系统可能会出现多Leader的情况，但这并不会破坏算法对一致性的保证，此时多个Leader都可以提出自己的提案，优化的算法就退化成了原始的paxos算法。

一个Leader的工作流程主要有分为三个阶段：

(1).学习阶段 向其它的参与者学习自己不知道的数据(决议);

(2).同步阶段 让绝大多数参与者保持数据(决议)的一致性;

(3).服务阶段 为客户端服务，提议案;


1.3.1 学习阶段
当一个参与者成为了Leader之后，它应该需要知道绝大多数的paxos实例，因此就会马上启动一个主动学习的过程。假设当前的新Leader早就知道了1-134、138和139的paxos实例，那么它会执行135-137和大于139的paxos实例的第一阶段。如果只检测到135和140的paxos实例有确定的值，那它最后就会知道1-135以及138-140的paxos实例。

1.3.2 同步阶段
此时的Leader已经知道了1-135、138-140的paxos实例，那么它就会重新执行1-135的paxos实例，以保证绝大多数参与者在1-135的paxos实例上是保持一致的。至于139-140的paxos实例，它并不马上执行138-140的paxos实例，而是等到在服务阶段填充了136、137的paxos实例之后再执行。这里之所以要填充间隔，是为了避免以后的Leader总是要学习这些间隔中的paxos实例，而这些paxos实例又没有对应的确定值。

1.3.4 服务阶段
Leader将用户的请求转化为对应的paxos实例，当然，它可以并发的执行多个paxos实例，当这个Leader出现异常之后，就很有可能造成paxos实例出现间断。

1.3.5 问题
(1).Leader的选举原则

(2).Acceptor如何感知当前Leader的失败，客户如何知道当前的Leader

(3).当出现多Leader之后，如何kill掉多余的Leader

(4).如何动态的扩展Acceptor

2. Zookeeper
2.1 整体架构
在Zookeeper集群中，主要分为三者角色，而每一个节点同时只能扮演一种角色，这三种角色分别是：

(1). Leader 接受所有Follower的提案请求并统一协调发起提案的投票，负责与所有的Follower进行内部的数据交换(同步);

(2). Follower 直接为客户端服务并参与提案的投票，同时与Leader进行数据交换(同步);

(3). Observer 直接为客户端服务但并不参与提案的投票，同时也与Leader进行数据交换(同步);


2.2 QuorumPeer的基本设计

Zookeeper对于每个节点QuorumPeer的设计相当的灵活，QuorumPeer主要包括四个组件：客户端请求接收器(ServerCnxnFactory)、数据引擎(ZKDatabase)、选举器(Election)、核心功能组件(Leader/Follower/Observer)。其中：

(1). ServerCnxnFactory负责维护与客户端的连接(接收客户端的请求并发送相应的响应);

(2). ZKDatabase负责存储/加载/查找数据(基于目录树结构的KV+操作日志+客户端Session);

(3). Election负责选举集群的一个Leader节点;

(4). Leader/Follower/Observer一个QuorumPeer节点应该完成的核心职责;

2.3 QuorumPeer工作流程

2.3.1 Leader职责

Follower确认: 等待所有的Follower连接注册，若在规定的时间内收到合法的Follower注册数量，则确认成功；否则，确认失败。


2.3.2 Follower职责

2.4 选举算法
2.4.1 LeaderElection选举算法

选举线程由当前Server发起选举的线程担任，他主要的功能对投票结果进行统计，并选出推荐的Server。选举线程首先向所有Server发起一次询问(包括自己)，被询问方，根据自己当前的状态作相应的回复，选举线程收到回复后，验证是否是自己发起的询问(验证xid 是否一致)，然后获取对方的id(myid)，并存储到当前询问对象列表中，最后获取对方提议 的

leader 相关信息(id,zxid)，并将这些 信息存储到当次选举的投票记录表中，当向所有Serve r

都询问完以后，对统计结果进行筛选并进行统计，计算出当次询问后获胜的是哪一个Server，并将当前zxid最大的Server 设置为当前Server要推荐的Server(有可能是自己，也有可以是其它的Server，根据投票结果而定，但是每一个Server在第一次投票时都会投自己)，如果此时获胜的Server获得n/2 + 1的Server票数，设置当前推荐的leader为获胜的Server。根据获胜的Server相关信息设置自己的状态。每一个Server都重复以上流程直到选举出Leader。


初始化选票(第一张选票): 每个quorum节点一开始都投给自己;

收集选票: 使用UDP协议尽量收集所有quorum节点当前的选票(单线程/同步方式)，超时设置200ms;

统计选票: 1).每个quorum节点的票数;

         2).为自己产生一张新选票(zxid、myid均最大);

选举成功: 某一个quorum节点的票数超过半数;

更新选票: 在本轮选举失败的情况下，当前quorum节点会从收集的选票中选取合适的选票(zxid、myid均最大)作为自己下一轮选举的投票;

异常问题的处理

1). 选举过程中，Server的加入

当一个Server启动时它都会发起一次选举，此时由选举线程发起相关流程，那么每个 Serve r都会获得当前zxi d最大的哪个Serve r是谁，如果当次最大的Serve r没有获得n/2+1 个票数，那么下一次投票时，他将向zxid最大的Server投票，重复以上流程，最后一定能选举出一个Leader。

2). 选举过程中，Server的退出

只要保证n/2+1个Server存活就没有任何问题，如果少于n/2+1个Server 存活就没办法选出Leader。

3). 选举过程中，Leader死亡

当选举出Leader以后，此时每个Server应该是什么状态(FLLOWING)都已经确定，此时由于Leader已经死亡我们就不管它，其它的Fllower按正常的流程继续下去，当完成这个流程以后，所有的Fllower都会向Leader发送Ping消息，如果无法ping通，就改变自己的状为(FLLOWING ==> LOOKING)，发起新的一轮选举。

4). 选举完成以后，Leader死亡

处理过程同上。

5). 双主问题

Leader的选举是保证只产生一个公认的Leader的，而且Follower重新选举与旧Leader恢复并退出基本上是同时发生的，当Follower无法ping同Leader是就认为Leader已经出问题开始重新选举，Leader收到Follower的ping没有达到半数以上则要退出Leader重新选举。

2.4.2 FastLeaderElection选举算法
FastLeaderElection是标准的fast paxos的实现，它首先向所有Server提议自己要成为leader，当其它Server收到提议以后，解决 epoch 和 zxid 的冲突，并接受对方的提议，然后向对方发送接受提议完成的消息。

FastLeaderElection算法通过异步的通信方式来收集其它节点的选票，同时在分析选票时又根据投票者的当前状态来作不同的处理，以加快Leader的选举进程。

每个Server都一个接收线程池和一个发送线程池, 在没有发起选举时，这两个线程池处于阻塞状态，直到有消息到来时才解除阻塞并处理消息，同时每个Serve r都有一个选举线程(可以发起选举的线程担任)。

1). 主动发起选举端(选举线程)的处理

首先自己的 logicalclock加1，然后生成notification消息，并将消息放入发送队列中， 系统中配置有几个Server就生成几条消息，保证每个Server都能收到此消息，如果当前Server 的状态是LOOKING就一直循环检查接收队列是否有消息，如果有消息，根据消息中对方的状态进行相应的处理。

2).主动发送消息端(发送线程池)的处理

将要发送的消息由Notification消息转换成ToSend消息，然后发送对方，并等待对方的回复。

3). 被动接收消息端(接收线程池)的处理

将收到的消息转换成Notification消息放入接收队列中，如果对方Server的epoch小于logicalclock则向其发送一个消息(让其更新epoch)；如果对方Server处于Looking状态，自己则处于Following或Leading状态，则也发送一个消息(当前Leader已产生，让其尽快收敛)。


2.4.3 AuthFastLeaderElection选举算法
AuthFastLeaderElection算法同FastLeaderElection算法基本一致，只是在消息中加入了认证信息，该算法在最新的Zookeeper中也建议弃用。
2.6 Zookeeper中的请求处理流程
2.6.1 Follower节点处理用户的读写请求

2.6.2 Leader节点处理写请求


    值得注意的是， Follower/Leader上的读操作时并行的，读写操作是串行的，当CommitRequestProcessor处理一个写请求时，会阻塞之后所有的读写请求。

