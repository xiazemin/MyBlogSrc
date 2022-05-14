---
title: Gossip_zap_raft_paxos 对比
layout: post
category: algorithm
author: 夏泽民
---
1，Raft算法
raft 集群中的每个节点都可以根据集群运行的情况在三种状态间切换：follower, candidate 与 leader。leader 向 follower 同步日志，follower 只从 leader 处获取日志。在节点初始启动时，节点的 raft 状态机将处于 follower 状态并被设定一个 election timeout，如果在这一时间周期内没有收到来自 leader 的 heartbeat，节点将发起选举：节点在将自己的状态切换为 candidate 之后，向集群中其它 follower 节点发送请求，询问其是否选举自己成为 leader。当收到来自集群中过半数节点的接受投票后，节点即成为 leader，开始接收保存 client 的数据并向其它的 follower 节点同步日志。leader 节点依靠定时向 follower 发送 heartbeat 来保持其地位。任何时候如果其它 follower 在 election timeout 期间都没有收到来自 leader 的 heartbeat，同样会将自己的状态切换为 candidate 并发起选举。每成功选举一次，新 leader 的步进数都会比之前 leader 的步进数大1。

Raft一致性算法处理日志复制以保证强一致性。

follower 节点不可用
follower 节点不可用的情况相对容易解决。因为集群中的日志内容始终是从 leader 节点同步的，只要这一节点再次加入集群时重新从 leader 节点处复制日志即可。

leader 不可用
一般情况下，leader 节点定时发送 heartbeat 到 follower 节点。



由于某些异常导致 leader 不再发送 heartbeat ，或 follower 无法收到 heartbeat 。



当某一 follower 发生 election timeout 时，其状态变更为 candidate，并向其他 follower 发起投票。



当超过半数的 follower 接受投票后，这一节点将成为新的 leader，leader 的步进数加1并开始向 follower 同步日志。



当一段时间之后，如果之前的 leader 再次加入集群，则两个 leader 比较彼此的步进数，步进数低的 leader 将切换自己的状态为 follower。



较早前 leader 中不一致的日志将被清除，并与现有 leader 中的日志保持一致。



2，Gossip协议
传统的监控，如ceilometer，由于每个节点都会向server报告状态，随着节点数量的增加server的压力随之增大。分布式健康检查可以解决这类性能瓶颈，降节点数量从数百台扩至数千台，甚至更多。

Agent在每台节点上运行，可以在每个Agent上添加一些健康检查的动作，Agent会周期性的运行这些动作。用户可以添加脚本或者请求一个URL链接。一旦有健康检查报告失败，Agent就把这个事件上报给服务器节点。用户可以在服务器节点上订阅健康检查事件，并处理这些报错消息。

在所有的Agent之间（包括服务器模式和普通模式）运行着Gossip协议。服务器节点和普通Agent都会加入这个Gossip集群，收发Gossip消息。每隔一段时间，每个节点都会随机选择几个节点发送Gossip消息，其他节点会再次随机选择其他几个节点接力发送消息。这样一段时间过后，整个集群都能收到这条消息。

Gossip协议已经是P2P网络中比较成熟的协议了。Gossip协议的最大的好处是，即使集群节点的数量增加，每个节点的负载也不会增加很多，几乎是恒定的。这就允许Consul管理的集群规模能横向扩展到数千个节点。

Consul的每个Agent会利用Gossip协议互相检查在线状态，本质上是节点之间互Ping，分担了服务器节点的心跳压力。如果有节点掉线，不用服务器节点检查，其他普通节点会发现，然后用Gossip广播给整个集群。

Gossip算法又被称为反熵（Anti-Entropy），熵是物理学上的一个概念，代表杂乱无章，而反熵就是在杂乱无章中寻求一致，这充分说明了Gossip的特点：在一个有界网络中，每个节点都随机地与其他节点通信，经过一番杂乱无章的通信，最终所有节点的状态都会达成一致。每个节点可能知道所有其他节点，也可能仅知道几个邻居节点，只要这些节可以通过网络连通，最终他们的状态都是一致的，当然这也是疫情传播的特点。

要注意到的一点是，即使有的节点因宕机而重启，有新节点加入，但经过一段时间后，这些节点的状态也会与其他节点达成一致，也就是说，Gossip天然具有分布式容错的优点。

3，Zookeeper ZAB协议
  ZooKeeper为高可用的一致性协调框架，自然的ZooKeeper也有着一致性算法的实现，ZooKeeper使用的是ZAB协议作为数据一致性的算法，ZAB（ZooKeeper Atomic Broadcast ）全称为：原子消息广播协议；ZAB协议设计了支持崩溃恢复，ZooKeeper使用唯一主节点Leader用于处理客户端所有事务请求，采用ZAB协议将服务器数状态以事务形式广播到所有Follower上，还有若干个只提供读，不提供写和投票的observer。由于事务间可能存在着依赖关系，ZAB协议保证Leader广播的变更序列被顺序的处理（维护一个队列）：一个状态被处理那么它所依赖的状态也已经提前被处理；ZAB协议支持的崩溃恢复可以保证在Leader进程崩溃的时候可以重新选出Leader并且保证数据的完整性

      ZAB可以说是在Paxos算法基础上进行了扩展改造而来的，在ZooKeeper中所有的事务请求都由一个主服务器也就是Leader来处理，其他服务器为Follower，Leader将客户端的事务请求转换为事务Proposal，并且将Proposal分发给集群中其他所有的Follower，然后Leader等待Follwer反馈，当有过半数（>=N/2+1）的Follower反馈信息后，Leader将再次向集群内Follower广播Commit信息，Commit为将之前的Proposal提交；

节点状态
      ZAB协议中存在着三种状态，每个节点都属于以下三种中的一种：

Looking：系统刚启动时或者Leader崩溃后正处于选举状态
Following：Follower节点所处的状态，Follower与Leader处于数据同步阶段；
Leading：Leader所处状态，当前集群中有一个Leader为主进程；
      ZooKeeper启动时所有节点初始状态为Looking，这时集群会尝试选举出一个Leader节点，选举出的Leader节点切换为Leading状态；当节点发现集群中已经选举出Leader则该节点会切换到Following状态，然后和Leader节点保持同步；当Follower节点与Leader失去联系时Follower节点则会切换到Looking状态，开始新一轮选举；在ZooKeeper的整个生命周期中每个节点都会在Looking、Following、Leading状态间不断转换；



 

      状态切换图

      选举出Leader节点后ZAB进入原子广播阶段，这时Leader为和自己同步的每个节点Follower创建一个操作序列，一个时期一个Follower只能和一个Leader保持同步，Leader节点与Follower节点使用心跳检测来感知对方的存在；当Leader节点在超时时间内收到来自Follower的心跳检测那Follower节点会一直与该节点保持连接；若超时时间内Leader没有接收到来自过半Follower节点的心跳检测或TCP连接断开，那Leader会结束当前周期的领导，切换到Looking状态，所有Follower节点也会放弃该Leader节点切换到Looking状态，然后开始新一轮选举；

协议阶段
    ZAB协议（理论）定义了选举（election）、发现（discovery）、同步（sync）、广播(Broadcast)四个阶段；在实现中，ZAB选举（election）时当Follower存在ZXID（事务ID）时判断所有Follower节点的事务日志，只有lastZXID的节点才有资格成为Leader，这种情况下选举出来的Leader总有最新的事务日志，基于这个原因所以ZooKeeper实现的时候把发现（discovery）与同步（sync）合并为恢复（recovery）阶段；
      1. Election：在Looking状态中选举出准Leader节点，Leader的lastZXID总是最新的(下面详解)；
      2. Discovery：Follower节点向准Leader推送FOllOWERINFO，该信息中包含了上一周期的epoch，接受准Leader的NEWLEADER指令，检查newEpoch有效性，准Leader要确保Follower的epoch与ZXID小于或等于自身的，更新Follower的 acceptedEpoch；
      3. sync：将Follower与Leader的数据进行同步，由Leader发起同步指令，只有当 半数以上同步完成，准 leader 才会成为真正的Leader，最终保持集群数据的一致性；
      4. Broadcast：Leader广播Proposal与Commit，Follower接受Proposal与Commit(下面详解)；
      5. Recovery：在Election阶段选举出Leader后本阶段主要工作就是进行数据的同步，使Leader具有highestZXID，集群保持数据的一致性(下面详解)；

      选举（Election）
      election阶段必须确保选出的Leader具有highestZXID，否则在Recovery阶段没法保证数据的一致性，Recovery阶段Leader要求Follower向自己同步数据而不是Follower要求Leader保持数据同步，所有选举出来的Leader要具有最新的ZXID；
      在选举的过程中会对每个Follower节点的ZXID进行对比只有highestZXID的Follower才可能当选Leader；
      选举流程：
      1. 每个Follower都向其他节点发送选自身为Leader的Vote投票请求，等待回复；
      2. Follower接受到的Vote如果比自身的大（ZXID更新）时则投票，并更新自身的Vote，否则拒绝投票；
      3. 每个Follower中维护着一个投票记录表，当某个节点收到过半的投票时，结束投票并把该Follower选为Leader，投票结束；

      ZAB协议中使用ZXID作为事务编号，ZXID为64位数字，低32位为一个递增的计数器，每一个客户端的一个事务请求时Leader产生新的事务后该计数器都会加1，高32位为Leader周期epoch编号（epoch就像皇帝的年号），当新选举出一个Leader节点时Leader会取出本地日志中最大事务Proposal的ZXID解析出对应的epoch把该值加1作为新的epoch，将低32位从0开始生成新的ZXID；ZAB使用epoch来区分不同的Leader周期；

      恢复（Recovery）
      在election阶段选举出来的Leader已经具有最新的ZXID，所有本阶段的主要工作是根据Leader的事务日志对Follower节点数据进行更新；
      Leader：Leader生成新的ZXID与epoch，接收Follower发送过来的FOllOWERINFO（含有当前节点的LastZXID）然后往Follower发送NEWLEADER；Leader根据Follower发送过来的LastZXID根据数据更新策略向Follower发送更新指令；
      同步策略：
      1. SNAP：如果Follower数据太老，Leader将发送快照SNAP指令给Follower同步数据；
      2. DIFF：Leader发送从Follolwer.lastZXID到Leader.lastZXID议案的DIFF指令给Follower同步数据；
   3. TRUNC：当Follower.lastZXID比Leader.lastZXID大时，Leader发送从Leader.lastZXID到Follower.lastZXID的TRUNC指令让Follower丢弃该段数据；
      Follower：往Leader发送FOLLOWERINFO指令，Leader拒绝就转到Election阶段；接收Leader的NEWLEADER指令，如果该指令中epoch比当前Follower的epoch小那么Follower转到Election阶段；Follower还有主要工作是接收SNAP/DIFF/TRUNC指令同步数据与ZXID，同步成功后回复ACKNETLEADER，然后进入下一阶段；Follower将所有事务都同步完成后Leader会把该节点添加到可用Follower列表中；
      SNAP与DIFF用于保证集群中Follower节点已经Committed的数据的一致性，TRUNC用于抛弃已经被处理但是没有Committed的数据；

      广播(Broadcast)
      客户端提交事务请求时Leader节点为每一个请求生成一个事务Proposal，将其发送给集群中所有的Follower节点，收到过半Follower的反馈后开始对事务进行提交，ZAB协议使用了原子广播协议；在ZAB协议中只需要得到过半的Follower节点反馈Ack就可以对事务进行提交，这也导致了Leader几点崩溃后可能会出现数据不一致的情况，ZAB使用了崩溃恢复来处理数字不一致问题；消息广播使用了TCP协议进行通讯所有保证了接受和发送事务的顺序性。广播消息时Leader节点为每个事务Proposal分配一个全局递增的ZXID（事务ID），每个事务Proposal都按照ZXID顺序来处理；
      Leader节点为每一个Follower节点分配一个队列按事务ZXID顺序放入到队列中，且根据队列的规则FIFO来进行事务的发送。Follower节点收到事务Proposal后会将该事务以事务日志方式写入到本地磁盘中，成功后反馈Ack消息给Leader节点，Leader在接收到过半Follower节点的Ack反馈后就会进行事务的提交，以此同时向所有的Follower节点广播Commit消息，Follower节点收到Commit后开始对事务进行提交；
      
4，Paxos
Paxos算法是莱斯利·兰伯特(Leslie Lamport)1990年提出的一种基于消息传递的一致性算法。Paxos算法解决的问题是一个分布式系统如何就某个值（决议）达成一致。在工程实践意义上来说，就是可以通过Paxos实现多副本一致性，分布式锁，名字管理，序列号分配等。比如，在一个分布式数据库系统中，如果各节点的初始状态一致，每个节点执行相同的操作序列，那么他们最后能得到一个一致的状态。为保证每个节点执行相同的命令序列，需要在每一条指令上执行一个“一致性算法”以保证每个节点看到的指令一致。本文首先会讲原始的Paxos算法(Basic Paxos)，主要描述二阶段提交过程，然后会着重讲Paxos算法的变种(Multi Paxos)，它是对Basic Paxos的优化，而且更适合工程实践，最后我会通过Q&A的方式，给出我在学习Paxos算法中的疑问，以及我对这些疑问的理解。

概念与术语 
Proposer：提议发起者，处理客户端请求，将客户端的请求发送到集群中，以便决定这个值是否可以被批准。
Acceptor：提议批准者，负责处理接收到的提议，他们的回复就是一次投票，会存储一些状态来决定是否接收一个值。
replica：节点或者副本，分布式系统中的一个server，一般是一台单独的物理机或者虚拟机，同时承担paxos中的提议者和接收者角色。
ProposalId：每个提议都有一个编号，编号高的提议优先级高。
Paxos Instance：Paxos中用来在多个节点之间对同一个值达成一致的过程，比如同一个日志序列号：logIndex，不同的logIndex属于不同的Paxos Instance。
acceptedProposal：在一个Paxos Instance内，已经接收过的提议
acceptedValue：在一个Paxos Instance内，已经接收过的提议对应的值。
minProposal：在一个Paxos Instance内，当前接收的最小提议值，会不断更新

Basic-Paxos算法
     基于Paxos协议构建的系统，只需要系统中超过半数的节点在线且相互通信正常即可正常对外提供服务。它的核心实现Paxos Instance主要包括两个阶段:准备阶段(prepare phase)和提议阶段(accept phase)。如下图所示：
	<img src="{{site.url}}{{site.baseurl}}/img/basicPaxos.png"/>
1.获取一个ProposalId,为了保证ProposalId递增，可以采用时间戳+serverId方式生成；
2.提议者向所有节点广播prepare(n)请求；
3.接收者比较n和minProposal，如果n>minProposal,表示有更新的提议，minProposal=n；否则将(acceptedProposal,acceptedValue)返回；
4.提议者接收到过半数请求后，如果发现有acceptedValue返回，表示有更新的提议，保存acceptedValue到本地，然后跳转1，生成一个更高的提议；
5.到这里表示在当前paxos instance内，没有优先级更高的提议，可以进入第二阶段，广播accept(n,value)到所有节点；
6.接收者比较n和minProposal，如果n>=minProposal,则acceptedProposal=minProposal=n，acceptedValue=value，本地持久化后，返回；
否则，返回minProposal
7.提议者接收到过半数请求后，如果发现有返回值>n，表示有更新的提议，跳转1；否则value达成一致。
从上述流程可知，并发情况下，可能会出现第4步或者第7步频繁重试的情况，导致性能低下，更严重者可能导致永远都无法达成一致的情况，就是所谓的“活锁”
1.S1作为提议者，发起prepare(3.1),并在S1,S2和S3达成多数派；
2.随后S5作为提议者 ，发起了prepare(3.5)，并在S3,S4和S5达成多数派；
3.S1发起accept(3.1,value1)，由于S3上提议 3.5>3.1,导致accept请求无法达成多数派，S1尝试重新生成提议
4.S1发起prepare(4.1),并在S1，S2和S3达成多数派
5.S5发起accpet(3.5,value5)，由于S3上提议4.1>3.5，导致accept请求无法达成多数派，S5尝试重新生成提议
6.S5发起prepare(5.5),并在S3,S4和S5达成多数派，导致后续的S1发起的accept(4.1,value1)失败

......

prepare阶段的作用
从Basic-Paxos的描述可知，需要通过两阶段来最终确定一个值，由于轮回多，导致性能低下，至少两次网络RTT。那么prepare阶段能否省去？
	<img src="{{site.url}}{{site.baseurl}}/img/paxossimple.png"/>
	1.S1首先发起accept(1,red)，并在S1,S2和S3达成多数派，red在S1，S2，S3上持久化
2.随后S5发起accept(5,blue)，对于S3而言，由于接收到更新的提议，会将acceptedValue值改为blue
3.那么S3，S4和S5达成多数派，blue在S3，S4和S5持久化
4.最后的结果是，S1和S2的值是red，而S3，S4和S5的值是blue，没有达成一致。

所以两阶段必不可少，Prepare阶段的作用是阻塞旧的提议，并且返回已经接收到的acceptedProposal。同时也可以看到的是，假设只有S1提议，则不会出现问题，这就是我们下面要讲的Multi-Paxos。

Multi-paxos算法
     Paxos是对一个值达成一致，Multi-Paxos是连续多个paxos instance来对多个值达成一致，这里最核心的原因是multi-paxos协议中有一个Leader。Leader是系统中唯一的Proposal，在lease租约周期内所有提案都有相同的ProposalId，可以跳过prepare阶段，议案只有accept过程，一个ProposalId可以对应多个Value，所以称为Multi-Paxos。

Multi Paxos一边先运行一次完整的paxos算法选举出leader，有leader处理所有的读写请求，然后省略掉prepare过程.

Multi Paxos要求在各个Proposer中有唯一的Leader，并由这个Leader唯一地提交value给各Acceptor进行表决，在系统中仅有一个Leader进行value提交的情况下，Prepare的过程就可以被跳过：



 

  如上图：

流程图中没有了basic paxos的两阶段，变成了一个一阶段的递交协议：
一阶段a：发起者（leader）直接告诉Acceptor，准备递交协议号为I+1的协议
一阶段b：收到了大部分acceptor的回复后（图中是全部），acceptor就直接回复client协议成功写入
wiki中写的Accept方法，我更愿意把它当做prepare，因为如果没有半数返回，该协议在超时后会返回失败，这种情况下，I+1这个协议号并没有通过，在下个请求是仍是使用I+1这个协议号
选举
     首先我们需要有一个leader，其实选主的实质也是一次Paxos算法的过程，只不过这次Paxos确定的“谁是leader”这个值。由于任何一个节点都可以发起提议，在并发情况下，可能会出现多主的情况，比如A，B先后当选为leader。为了避免频繁选主，当选leader的节点要马上树立自己的leader权威(让其它节点知道它是leader)，写一条特殊日志(start-working日志)确认其身份。根据多数派原则，只有一个leader的startworking日志可以达成多数派。leader确认身份后，可以通过了lease机制(租约)维持自己的leader身份，使得其它proposal不再发起提案，这样就进入了leader任期，由于没有并发冲突，因此可以跳过prepare阶段，直接进入accept阶段。通过分析可知，选出leader后，leader任期内的所有日志都只需要一个网络RTT(Round Trip Time)即可达成一致。

新主恢复流程
      由于Paxos中并没有限制，任何节点都可以参与选主并最终成为leader，这就无法保证新选出的leader包含了所有日志，可能存在空洞，因此在真正提供服务前，还存在一个获取所有已提交日志的恢复过程。新主向所有成员查询最大logId的请求，收到多数派响应后，选择最大的logId作为日志恢复结束点，这里多数派的意义在于恢复结束点包含了所有达成一致的日志，当然也可能包含了没有达成多数派的日志。拿到logId后，从头开始对每个logId逐条进行paxos协议，因为在新主获得所有日志之前，系统是无法提供服务的。为了优化，引入了confirm机制，就是将已经达成一致的logId告诉其它acceptor，acceptor写一条confirm日志到日志文件中。那么新主在重启后，扫描本地日志，对于已经拥有confirm日志的log，就不会重新发起paxos了。同样的，在响应客户端请求时，对于没有confirm日志的log，需要重新发起一轮paxos。由于没有严格要求confirm日志的位置，可以批量发送。为了确保重启时，不需要对太多已提价的log进行paxos，需要将confirm日志与最新提交的logId保持一定的距离。

性能优化
      Basic-Paxos一次日志确认，需要至少2次磁盘写操作(prepare,promise)和2次网络RTT(prepare,promise)。Multi-Paxos利用一阶段提交(省去Prepare阶段)，将一次日志确认缩短为一个RTT和一次磁盘写；通过confirm机制，可以缩短新主的恢复时间。为了提高性能，我们还可以实现一批日志作为一个组提交，要么成功一批，要么都不成功，这点类似于group-commit，通过RT换取吞吐量。

安全性(异常处理)
1.Leader异常
Leader在任期内，需要定期给各个节点发送心跳，已告知它还活着(正常工作)，如果一个节点在超时时间内仍然没有收到心跳，它会尝试发起选主流程。Leader异常了，则所有的节点先后都会出现超时，进入选主流程，选出新的主，然后新主进入恢复流程，最后再对外提供服务。我们通常所说的异常包括以下三类：
(1).进程crash(OS crash)
Leader进程crash和Os crash类似，只要重启时间大于心跳超时时间都会导致节点认为leader挂了，触发重新选主流程。
(2).节点网络异常(节点所在网络分区)
Leader网络异常同样会导致其它节点收不到心跳，但有可能leader是活着的，只不过发生了网络抖动，因此心跳超时不能设置的太短，否则容易因为网络抖动造成频繁选主。另外一种情况是，节点所在的IDC发生了分区，则同一个IDC的节点相互还可以通信，如果IDC中节点能构成多数派，则正常对外服务，如果不能，比如总共4个节点，两个IDC，发生分区后会发现任何一个IDC都无法达成多数派，导致无法选出主的问题。因此一般Paxos节点数都是奇数个，而且在部署节点时，IDC节点的分布也要考虑。
(3).磁盘故障
前面两种异常，磁盘都是OK的，即已接收到的日志以及对应confirm日志都在。如果磁盘故障了，节点再加入就类似于一个新节点，上面没有任何日志和Proposal信息。这种情况会导致一个问题就是，这个节点可能会promise一个比已经promise过的最大proposalID更小的proposal，这就违背了Paxos原则。因此重启后，节点不能参与Paxos Instance，它需要先追上Leader，当观察到一次完整的paxos instance时该节点结束不能promise/ack状态。
2.Follower异常(宕机，磁盘损坏等)
对于Follower异常，则处理要简单的多，因为follower本身不对外提供服务(日志可能不全)，对于leader而言，只要能达成多数派，就可以对外提供服务。follower重启后，没有promise能力，直到追上leader为止。

Q&A
1.Paxos协议数据同步方式相对于基于传统1主N备的同步方式有啥区别？
      一般情况下，传统数据库的高可用都是基于主备来实现，1主1备2个副本，主库crash后，通过HA工具来进行切换，提升备库为主库。在强一致场景下，复制可以开启强同步，Oracle和Mysql都是类似的复制模式。但是如果备库网络抖动，或者crash，都会导致日志同步失败，服务不可用。为此，可以引入1主N备的多副本形式，我们对比都是3副本的情况，一个是基于传统的1主2备，另一种基于paxos的1主2备。传统的1主两备，进行日志同步时，只要有一个副本接收到日志并就持久化成功，就可以返回，在一定程度上解决了网络抖动和备库crash问题。但如果主库出问题后，还是要借助于HA工具来进行切换，那么HA切换工具的可用性如何来保证又成为一个问题。基于Paxos的多副本同步其实是在1主N备的基础上引入了一致性协议，这样整个系统的可用性完全有3个副本控制，不需要额外的HA工具。而实际上，很多系统为了保证多节点HA工具获取主备信息的一致性，采用了zookeeper等第三方接口来实现分布式锁，其实本质也是基于Paxos来实现的。

      我这里以MySQL的主备复制一套体系为例来具体说明传统的主备保持强一致性的一些问题。整个系统中主要包含以下几种角色：Master，Slave，Zookeeper-Service(zk)，HA-Console(HA)，Zookeeper-Agent(Agent)
Master,Slave:分别表示主节点和备节点,主节点提供读写服务，备节点可以提供读服务，或者完全用于容灾。
Zookeeper-Service(zk):分布式一致性服务，负责管理Master/Slave节点的存活信息和切换信息。zk基于zab协议，zab协议是一种一致性协议，与paxos，raft协议类似，它主要有两种模式，包括恢复模式(选主)和广播模式(同步)。一般zk包含N个节点(zk-node)，只要有超过半数的zk-node存活且相互连通，则zk可以对外提供一致性服务。
HA-Console:切换工具，负责具体的切换流程
Zookeeper-Agent(Agent):Master/Slave实例上的监听进程，与监听的实例保持心跳，维护Master/Slave的状态，每个实例有一个对应的Agent。大概工作流程如下：
(1).Master/Slave正常启动并搭建好了复制关系，对应的Agent会调用zk接口去注册alive节点信息，假设分别为A和B。
(2).如果此时Master Crash，则实例对应的Agent发现心跳失败，如果重试几次后仍然失败，则将调用zk接口注销掉A节点信息。
(3).HA工具通过zk接口比较两次的节点信息，发现少了A节点，表示A可能不存在了，需要切换，尝试连接A，如果仍然不通，注册A的dead节点，并开始切换流程。
(4).HA工具读取配置信息，找到对应的Slave节点B，(更改读写比配置信息，设置B提供写)，打开写。
(5).应用程序通过拉取最新的配置信息，得知新主B，新的写入会写入B。
前面几部基本介绍了MySQL借助zk实现高可用的流程，由于zk-node可以多地部署，HA无状态，因此可以很容易实现同城或者是异地的高可用系统，并且动态可扩展，一个HA节点可以同时管理多个Master/Slave的切换。那么能保证一致性吗？前面提到的Agent除了做监听，还有一个作用是尽可能保持主备一致，并且不丢数据。
(6).假设此时A节点重启，Agent检测到，通过zk接口发现A节点在dead目录下，表示被切换过，需要kill上面的所有连接，并回滚crash时A比B多的binlog，为了尽可能的少丢数据，然后再开启binlog后，将这部分数据重做。这里要注意rollback和replay都在old-Master上面进行，rollback时需要关闭binlog，而replay需要开启binlog，保证这部分数据能流向new-Master。
(7).从第6步来看，可以一定程度上保证主备一致性，但是进行rollback和replay时，实际上是往new-Slave上写数据，这一定程度上造成了双写，如果此时new—Master也在写同一条记录，可能会导致写覆盖等问题。
(8).如果开启半同步呢？old-Master crash时，仍然可能比old-Slave多一个group的binlog，所以仍然需要借助于rollback和replay，依然避免不了双写，所以也不能做到严格意义上的强一致。

2.分布式事务与Paxos协议的关系？
     在数据库领域，提到分布式系统，就会提到分布式事务。Paxos协议与分布式事务并不是同一层面的东西。分布式事务的作用是保证跨节点事务的原子性，涉及事务的节点要么都提交(执行成功)，要么都不提交(回滚)。分布式事务的一致性通常通过2PC来保证(Two-Phase Commit, 2PC)，这里面涉及到一个协调者和若干个参与者。第一阶段，协调者询问参与者事务是否可以执行，参与者回复同意(本地执行成功)，回复取消(本地执行失败)。第二阶段，协调者根据第一阶段的投票结果进行决策，当且仅当所有的参与者同意提交事务时才能提交，否则回滚。2PC的最大问题是，协调者是单点(需要有一个备用节点)，另外协议是阻塞协议，任何一个参与者故障，都需要等待(可以通过加入超时机制)。Paxos协议用于解决多个副本之间的一致性问题。比如日志同步，保证各个节点的日志一致性，或者选主(主故障情况下)，保证投票达成一致，选主的唯一性。简而言之，2PC用于保证多个数据分片上事务的原子性，Paxos协议用于保证同一个数据分片在多个副本的一致性，所以两者可以是互补的关系，绝不是替代关系。对于2PC协调者单点问题，可以利用Paxos协议解决，当协调者出问题时，选一个新的协调者继续提供服务。工程实践中，Google Spanner，Google Chubby就是利用Paxos来实现多副本日志同步。

3.如何将Paxos应用于传统的数据库复制协议中？
    复制协议相当于是对Paxos的定制应用，通过对一系列日志进行投票确认达成多数派，就相当于日志已经在多数派持久化成功。副本通过应用已经持久化的日志，实现与Master节点同步。由于数据库ACID特性，本质是由一个一致的状态到另外一个一致的状态，每次事务操作都是对应数据库状态的变更，并生成一条日志。由于client操作有先后顺序，因此需要保证日志的先后的顺序，在任何副本中，不仅仅要保证所有日志都持久化了，而且要保证顺序。对于每条日志，通过一个logID标示，logID严格递增(标示顺序)，由leader对每个日志进行投票达成多数派，如果中途发生了leader切换，对于新leader中logID的“空洞”，需要重新投票，确认日志的有效性。

4.Multi-Paxos的非leader节点可以提供服务吗？
     Multi-Paxos协议中只有leader确保包含了所有已经已经持久化的日志，当然本地已经持久化的日志不一定达成了多数派，因此对于没有confirm的日志，需要再进行一次投票，然后将最新的结果返回给client。而非leader节点不一定有所有最新的数据，需要通过leader确认，所以一般工程实现中，所有的读写服务都由leader提供。

5.客户端请求过程中失败了，如何处理？
     client向leader发起一次请求，leader在返回前crash了。对于client而言，这次操作可能成功也可能失败。因此client需要检查操作的结果，确定是否要重新操作。如果leader在本地持久化后，并没有达成多数派时就crash，新leader首先会从各个副本获取最大的logID作为恢复结束点，对于它本地没有confirm的日志进行Paxos确认，如果此时达成多数派，则应用成功，如果没有则不应用。client进行检查时，会知道它的操作是否成功。当然具体工程实践中，这里面涉及到client超时时间，以及选主的时间和日志恢复时间。
<!-- more -->
