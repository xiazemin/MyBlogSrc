---
title: 多个进程绑定相同端口的实现分析[Google Patch]
layout: post
category: linux
author: 夏泽民
---
Google REUSEPORT 新特性，支持多个进程或者线程绑定到相同的 IP 和端口，以提高 server 的性能。
1. 设计思路
该特性实现了 IPv4/IPv6 下 TCP/UDP 协议的支持， 已经集成到 kernel 3.9 中。

核心的实现主要有三点：

扩展 socket option，增加 SO_REUSEPORT 选项，用来设置 reuseport。
修改 bind 系统调用实现，以便支持可以绑定到相同的 IP 和端口
修改处理新建连接的实现，查找 listener 的时候，能够支持在监听相同 IP 和端口的多个 sock 之间均衡选择。
共包含 7 个 patch，其中有两个为 buf fix
<!-- more -->
数据结构调整： 055dc21a1d1d219608cd4baac7d0683fb2cbbe8a
TCP/IPv4: da5e36308d9f7151845018369148201a5d28b46d
UDP/IPv4: ba418fa357a7b3c9d477f4706c6c7c96ddbd1360
TCP/IPv6: 5ba24953e9707387cce87b07f0d5fbdd03c5c11b
UDP/IPv6: 72289b96c943757220ccc681fe2e22b46e21aced
bug fix: 7c0cadc69ca2ac8893aa162ee80d92a805840909 fix: UDP/IPv4
bug fix: 5588d3742da9900323dc3d766845a53bacdfb5ab fix: 数据结构定义
下面根据该特性的实现，简单介绍 IPv4 下多个进程绑定相同 IP 和端口的逻辑分析。 kernel 代码版本：3.11-rc1。





2. 数据结构扩展
通用 sock 结构扩展，增加 skc_reuseport 成员，用于 socket option 配置是记录对应 结果：

--- a/include/net/sock.h
+++ b/include/net/sock.h
@@ -140,6 +140,7 @@ typedef __u64 __bitwise __addrpair;
  *	@skc_family: network address family
  *	@skc_state: Connection state
  *	@skc_reuse: %SO_REUSEADDR setting
+ *	@skc_reuseport: %SO_REUSEPORT setting
  *	@skc_bound_dev_if: bound device index if != 0
  *	@skc_bind_node: bind hash linkage for various protocol lookup tables
  *	@skc_portaddr_node: second hash linkage for UDP/UDP-Lite protocol
@@ -179,7 +180,8 @@ struct sock_common {
 
 	unsigned short		skc_family;
 	volatile unsigned char	skc_state;
-	unsigned char		skc_reuse;
+	unsigned char		skc_reuse:4;
+	unsigned char		skc_reuseport:4;
 	int			skc_bound_dev_if;
 	union {
 		struct hlist_node	skc_bind_node;
@@ -297,6 +299,7 @@ struct sock {
 #define sk_family		__sk_common.skc_family
 #define sk_state		__sk_common.skc_state
 #define sk_reuse		__sk_common.skc_reuse
+#define sk_reuseport		__sk_common.skc_reuseport
 #define sk_bound_dev_if		__sk_common.skc_bound_dev_if
 #define sk_bind_node		__sk_common.skc_bind_node
 #define sk_prot			__sk_common.skc_prot
 
--- a/net/core/sock.c
+++ b/net/core/sock.c
@@ -665,6 +665,9 @@ int sock_setsockopt(struct socket *sock, int level, int optname,
 	case SO_REUSEADDR:
 		sk->sk_reuse = (valbool ? SK_CAN_REUSE : SK_NO_REUSE);
 		break;
+	case SO_REUSEPORT:
+		sk->sk_reuseport = valbool;
+		break;
 	case SO_TYPE:
 	case SO_PROTOCOL:
 	case SO_DOMAIN:
bind socket 结构扩展，记录 fastreuseport 和 fastuid。这个会在执行 bind 时做相关 的初始化。其中，fastuid 应该是创建 fd 的 uid。

--- a/include/net/inet_hashtables.h
+++ b/include/net/inet_hashtables.h
@@ -81,7 +81,9 @@ struct inet_bind_bucket {
 	struct net		*ib_net;
 #endif
 	unsigned short		port;
-	signed short		fastreuse;
+	signed char		fastreuse;
+	signed char		fastreuseport;
+	kuid_t			fastuid;
 	int			num_owners;
 	struct hlist_node	node;
 	struct hlist_head	owners;
对于 TCP 来讲，owners 记录了使用相同端口号的 sock 列表。这个列表中的 sock 也包含 了监听 IP 不同的情况。而我们要分析的相同 IP 和端口 sock 也在该列表中。

3. bind 系统调用
分析该函数的 callpath，就是为了明确 google patch 中如果是绑定相同 IP 和 端口号的 多个 socket 如何成功的通过 bind 系统调用。如果没有该 patch 的话，应该返回 Address in use 之类的错误。

sys_bind()
-> inet_bind() (TCP)
-> sk->sk_prot->get_port(TCP: inet_csk_get_port)
inet_csk_get_port() 根据 bind 参数中指定的端口，查表 hashinfo->bhash

3.1. 初次绑定某端口
初次绑定某个端口的话，应该查表找不到对应的 struct inet_bind_bucket tb，因此要调用 inet_bind_bucket_create 创建一个表项，并作 resue 方面的初始化：

216 tb_not_found:
217     ret = 1;                                                                                                                      
218     if (!tb && (tb = inet_bind_bucket_create(hashinfo->bind_bucket_cachep,
219                     net, head, snum)) == NULL)
220         goto fail_unlock;
221     if (hlist_empty(&tb->owners)) {
222         if (sk->sk_reuse && sk->sk_state != TCP_LISTEN)
223             tb->fastreuse = 1;
224         else
225             tb->fastreuse = 0;
226         if (sk->sk_reuseport) {
227             tb->fastreuseport = 1;
228             tb->fastuid = uid;
229         } else
230             tb->fastreuseport = 0;
231     } else {
226-228 行： 如果 socket 设置了 reuseport 的话，则新建表项的 fastreuseport 置 1， fastuid 也记录下来，应该就是创建当前 socket fd 的 uid

接着调用 inet_bind_hash() 将当前的 sock 插入到 tb->owners 中，并增加计数

 62 void inet_bind_hash(struct sock *sk, struct inet_bind_bucket *tb,                                                                 
 63             const unsigned short snum)
 64 {       
 65     struct inet_hashinfo *hashinfo = sk->sk_prot->h.hashinfo;
 66             
 67     atomic_inc(&hashinfo->bsockets);
 68             
 69     inet_sk(sk)->inet_num = snum;
 70     sk_add_bind_node(sk, &tb->owners);
 71     tb->num_owners++;
 72     inet_csk(sk)->icsk_bind_hash = tb;
 73 }
并将 sock 对应 inet_connection_sock 的icsk_bind_hash 执行新分配的 tb。

3.2. 再次绑定相同端口
这次应该就可以找到对应的 tb，因此应该进行如下流程：

190 tb_found:
191     if (!hlist_empty(&tb->owners)) {
192         if (sk->sk_reuse == SK_FORCE_REUSE)
193             goto success;
194 
195         if (((tb->fastreuse > 0 &&
196               sk->sk_reuse && sk->sk_state != TCP_LISTEN) ||
197              (tb->fastreuseport > 0 &&
198               sk->sk_reuseport && uid_eq(tb->fastuid, uid))) &&
199             smallest_size == -1) {
200             goto success;
201         } else {
202             ret = 1;
203             if (inet_csk(sk)->icsk_af_ops->bind_conflict(sk, tb, true)) {
204                 if (((sk->sk_reuse && sk->sk_state != TCP_LISTEN) ||
205                      (tb->fastreuseport > 0 &&
206                       sk->sk_reuseport && uid_eq(tb->fastuid, uid))) &&
207                     smallest_size != -1 && --attempts >= 0) {
208                     spin_unlock(&head->lock);
209                     goto again;
210                 }
211 
212                 goto fail_unlock;
213             }
214         }
215     } 
195-196 为 socket reuse 的判断，并且非 LISTEN 的认为可以 bind，如果已经处理 LISTEN 状态的话，这里的条件不成立

197-198 为 Google patch 的检测，tb 配置启用了 reuseport，并且当前 socket 也设置 了reuseport，且 tb 和当前 socket 的 UID 一样，可以认为当前 socket 也可以放到 bind hash 中，随后会调用 inet_bind_hash 将当前 sock 也加入到 tb->owners 链表中。

4. listen 系统调用
sys_listen -> inet_listen -> inet_csk_listen_start

关键的实现就在 inet_csk_listen_start 中。重要的检测主要是再次检查端口是否可用。 因为 bind 和 listen 的执行有时间差，完全有可能被别的进程占去：

769     sk->sk_state = TCP_LISTEN;
770     if (!sk->sk_prot->get_port(sk, inet->inet_num)) {
771         inet->inet_sport = htons(inet->inet_num);                                                                                 
772    
773         sk_dst_reset(sk); 
774         sk->sk_prot->hash(sk);
775    
776         return 0;
777     }    
774 行调用 sk->sk_prot->hash(sk) 将对应的 sock 加入到 listening hash 中。 对于 TCP 而言， hash 指针指向 inet_hash()。这里记录下 listen socket 的 hash 的计算逻辑：

inet_hash
->__inet_hash(sk)
->inet_sk_listen_hashfn
->inet_lhashfn
238 /* These can have wildcards, don't try too hard. */
239 static inline int inet_lhashfn(struct net *net, const unsigned short num)                                                         
240 {
241     return (num + net_hash_mix(net)) & (INET_LHTABLE_SIZE - 1);
242 }
对于 listening socket，可以看出，应该是按照端口做 key 的，最终将 socket 放到了 listening_hash[] 中。

因此，绑定同一个端口的多个 listener sock 最后是放在了同一个 bucket 中。

5. 接受新连接
这里主要就是重点观察 TCP 协议栈将新建连接的请求分发给绑定了相同 IP 和端口的不同 listening socket。

tcp_v4_rcv
-> __inet_lookup_skb
-> __inet_lookup
->  __inet_lookup_listener （新建连接，只能通过 listener hash 查到其所属 listener）
__inet_lookup_listener 函数增加两个参数，saddr 和 sport。没有 Google patch 之前， 查找 listener 的话是不需要这两个参数的：

177 struct sock *__inet_lookup_listener(struct net *net,                                                                              
178                     struct inet_hashinfo *hashinfo,
179                     const __be32 saddr, __be16 sport,
180                     const __be32 daddr, const unsigned short hnum,
181                     const int dif)
182 {
... ...
191 begin:
192     result = NULL;
193     hiscore = 0;
194     sk_nulls_for_each_rcu(sk, node, &ilb->head) {
195         score = compute_score(sk, net, hnum, daddr, dif);
196         if (score > hiscore) {
197             result = sk;
198             hiscore = score;
199             reuseport = sk->sk_reuseport;
200             if (reuseport) {
201                 phash = inet_ehashfn(net, daddr, hnum,
202                              saddr, sport);
203                 matches = 1;
204             }
205         } else if (score == hiscore && reuseport) {
206             matches++;
207             if (((u64)phash * matches) >> 32 == 0)
208                 result = sk;
209             phash = next_pseudo_random32(phash);
210         }
211     }

该函数就是根据 sip+sport+dip+dport+dif 来查找合适的 listener。在没加入 google REUSEPORT patch 之前，是没有 sip 和 sport 的。这两个元素就是用来帮助在多个监 听相同 port 的 listener 之间做选择，并可能尽量保证公平。

这里有个函数调用 compute_score()，用来计算匹配的分数，得分最高的 listener 将作为 result 返回。计算的匹配分数主要是看 listener 的 portnum,rcv_saddr, 目的接口与 listener 的匹配程度。

196-204 行： 查到一个合适的 listener，而且得分比历史记录还高，记下该 sock。同时， 考虑到 reuseport 的问题，根据四元组计算一个 phash，match 置 1.

205 行： 走到这个分支，说明就是出现了 reuseport 的情况，而且是遍历到了第 N 个 （N>1）个监听相同端口的 listener。因此，其得分与历史得分肯定相等。

206-209 行：这几行代码就是实现了是否使用当前 listener 的逻辑。如果不使用的话， 那就继续遍历下一个。最终的结果就会在多个绑定相同端口的 listener 中使用其中一个。 因为 phash 的初次计算中加入了 saddr 和 sport，这个算法在 IP 地址及 port 足够多 的情况下保证了多个 listener 都会被平均分配到请求。

至此，google REUSEPORT 的 patch 简单的分析完毕。
