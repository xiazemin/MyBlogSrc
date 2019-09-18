---
title: system call table
layout: post
category: linux
author: 夏泽民
---
系统调用号	函数名	入口点	源代码
0	read	sys_read	fs/read_write.c
1	write	sys_write	fs/read_write.c
2	open	sys_open	fs/open.c
3	close	sys_close	fs/open.c
4	stat	sys_newstat	fs/stat.c
5	fstat	sys_newfstat	fs/stat.c
6	lstat	sys_newlstat	fs/stat.c
7	poll	sys_poll	fs/select.c
8	lseek	sys_lseek	fs/read_write.c
9	mmap	sys_mmap	arch/x86/kernel/sys_x86_64.c
10	mprotect	sys_mprotect	mm/mprotect.c
11	munmap	sys_munmap	mm/mmap.c
12	brk	sys_brk	mm/mmap.c
13	rt_sigaction	sys_rt_sigaction	kernel/signal.c
14	rt_sigprocmask	sys_rt_sigprocmask	kernel/signal.c
15	rt_sigreturn	stub_rt_sigreturn	arch/x86/kernel/signal.c
16	ioctl	sys_ioctl	fs/ioctl.c
17	pread64	sys_pread64	fs/read_write.c
18	pwrite64	sys_pwrite64	fs/read_write.c
19	readv	sys_readv	fs/read_write.c
20	writev	sys_writev	fs/read_write.c
21	access	sys_access	fs/open.c
22	pipe	sys_pipe	fs/pipe.c
23	select	sys_select	fs/select.c
24	sched_yield	sys_sched_yield	kernel/sched/core.c
25	mremap	sys_mremap	mm/mmap.c
26	msync	sys_msync	mm/msync.c
27	mincore	sys_mincore	mm/mincore.c
28	madvise	sys_madvise	mm/madvise.c
29	shmget	sys_shmget	ipc/shm.c
30	shmat	sys_shmat	ipc/shm.c
31	shmctl	sys_shmctl	ipc/shm.c
32	dup	sys_dup	fs/file.c
33	dup2	sys_dup2	fs/file.c
34	pause	sys_pause	kernel/signal.c
35	nanosleep	sys_nanosleep	kernel/hrtimer.c
36	getitimer	sys_getitimer	kernel/itimer.c
37	alarm	sys_alarm	kernel/timer.c
38	setitimer	sys_setitimer	kernel/itimer.c
39	getpid	sys_getpid	kernel/sys.c
40	sendfile	sys_sendfile64	fs/read_write.c
41	socket	sys_socket	net/socket.c
42	connect	sys_connect	net/socket.c
43	accept	sys_accept	net/socket.c
44	sendto	sys_sendto	net/socket.c
45	recvfrom	sys_recvfrom	net/socket.c
46	sendmsg	sys_sendmsg	net/socket.c
47	recvmsg	sys_recvmsg	net/socket.c
48	shutdown	sys_shutdown	net/socket.c
49	bind	sys_bind	net/socket.c
50	listen	sys_listen	net/socket.c
51	getsockname	sys_getsockname	net/socket.c
52	getpeername	sys_getpeername	net/socket.c
53	socketpair	sys_socketpair	net/socket.c
54	setsockopt	sys_setsockopt	net/socket.c
55	getsockopt	sys_getsockopt	net/socket.c
56	clone	stub_clone	kernel/fork.c
57	fork	stub_fork	kernel/fork.c
58	vfork	stub_vfork	kernel/fork.c
59	execve	stub_execve	fs/exec.c
60	exit	sys_exit	kernel/exit.c
61	wait4	sys_wait4	kernel/exit.c
62	kill	sys_kill	kernel/signal.c
63	uname	sys_newuname	kernel/sys.c
64	semget	sys_semget	ipc/sem.c
65	semop	sys_semop	ipc/sem.c
66	semctl	sys_semctl	ipc/sem.c
67	shmdt	sys_shmdt	ipc/shm.c
68	msgget	sys_msgget	ipc/msg.c
69	msgsnd	sys_msgsnd	ipc/msg.c
70	msgrcv	sys_msgrcv	ipc/msg.c
71	msgctl	sys_msgctl	ipc/msg.c
72	fcntl	sys_fcntl	fs/fcntl.c
73	flock	sys_flock	fs/locks.c
74	fsync	sys_fsync	fs/sync.c
75	fdatasync	sys_fdatasync	fs/sync.c
76	truncate	sys_truncate	fs/open.c
77	ftruncate	sys_ftruncate	fs/open.c
78	getdents	sys_getdents	fs/readdir.c
79	getcwd	sys_getcwd	fs/dcache.c
80	chdir	sys_chdir	fs/open.c
81	fchdir	sys_fchdir	fs/open.c
82	rename	sys_rename	fs/namei.c
83	mkdir	sys_mkdir	fs/namei.c
84	rmdir	sys_rmdir	fs/namei.c
85	creat	sys_creat	fs/open.c
86	link	sys_link	fs/namei.c
87	unlink	sys_unlink	fs/namei.c
88	symlink	sys_symlink	fs/namei.c
89	readlink	sys_readlink	fs/stat.c
90	chmod	sys_chmod	fs/open.c
91	fchmod	sys_fchmod	fs/open.c
92	chown	sys_chown	fs/open.c
93	fchown	sys_fchown	fs/open.c
94	lchown	sys_lchown	fs/open.c
95	umask	sys_umask	kernel/sys.c
96	gettimeofday	sys_gettimeofday	kernel/time.c
97	getrlimit	sys_getrlimit	kernel/sys.c
98	getrusage	sys_getrusage	kernel/sys.c
99	sysinfo	sys_sysinfo	kernel/sys.c
100	times	sys_times	kernel/sys.c
101	ptrace	sys_ptrace	kernel/ptrace.c
102	getuid	sys_getuid	kernel/sys.c
103	syslog	sys_syslog	kernel/printk/printk.c
104	getgid	sys_getgid	kernel/sys.c
105	setuid	sys_setuid	kernel/sys.c
106	setgid	sys_setgid	kernel/sys.c
107	geteuid	sys_geteuid	kernel/sys.c
108	getegid	sys_getegid	kernel/sys.c
109	setpgid	sys_setpgid	kernel/sys.c
110	getppid	sys_getppid	kernel/sys.c
111	getpgrp	sys_getpgrp	kernel/sys.c
112	setsid	sys_setsid	kernel/sys.c
113	setreuid	sys_setreuid	kernel/sys.c
114	setregid	sys_setregid	kernel/sys.c
115	getgroups	sys_getgroups	kernel/groups.c
116	setgroups	sys_setgroups	kernel/groups.c
117	setresuid	sys_setresuid	kernel/sys.c
118	getresuid	sys_getresuid	kernel/sys.c
119	setresgid	sys_setresgid	kernel/sys.c
120	getresgid	sys_getresgid	kernel/sys.c
121	getpgid	sys_getpgid	kernel/sys.c
122	setfsuid	sys_setfsuid	kernel/sys.c
123	setfsgid	sys_setfsgid	kernel/sys.c
124	getsid	sys_getsid	kernel/sys.c
125	capget	sys_capget	kernel/capability.c
126	capset	sys_capset	kernel/capability.c
127	rt_sigpending	sys_rt_sigpending	kernel/signal.c
128	rt_sigtimedwait	sys_rt_sigtimedwait	kernel/signal.c
129	rt_sigqueueinfo	sys_rt_sigqueueinfo	kernel/signal.c
130	rt_sigsuspend	sys_rt_sigsuspend	kernel/signal.c
131	sigaltstack	sys_sigaltstack	kernel/signal.c
132	utime	sys_utime	fs/utimes.c
133	mknod	sys_mknod	fs/namei.c
134	uselib	 	fs/exec.c
135	personality	sys_personality	kernel/exec_domain.c
136	ustat	sys_ustat	fs/statfs.c
137	statfs	sys_statfs	fs/statfs.c
138	fstatfs	sys_fstatfs	fs/statfs.c
139	sysfs	sys_sysfs	fs/filesystems.c
140	getpriority	sys_getpriority	kernel/sys.c
141	setpriority	sys_setpriority	kernel/sys.c
142	sched_setparam	sys_sched_setparam	kernel/sched/core.c
143	sched_getparam	sys_sched_getparam	kernel/sched/core.c
144	sched_setscheduler	sys_sched_setscheduler	kernel/sched/core.c
145	sched_getscheduler	sys_sched_getscheduler	kernel/sched/core.c
146	sched_get_priority_max	sys_sched_get_priority_max	kernel/sched/core.c
147	sched_get_priority_min	sys_sched_get_priority_min	kernel/sched/core.c
148	sched_rr_get_interval	sys_sched_rr_get_interval	kernel/sched/core.c
149	mlock	sys_mlock	mm/mlock.c
150	munlock	sys_munlock	mm/mlock.c
151	mlockall	sys_mlockall	mm/mlock.c
152	munlockall	sys_munlockall	mm/mlock.c
153	vhangup	sys_vhangup	fs/open.c
154	modify_ldt	sys_modify_ldt	arch/x86/um/ldt.c
155	pivot_root	sys_pivot_root	fs/namespace.c
156	_sysctl	sys_sysctl	kernel/sysctl_binary.c
157	prctl	sys_prctl	kernel/sys.c
158	arch_prctl	sys_arch_prctl	arch/x86/um/syscalls_64.c
159	adjtimex	sys_adjtimex	kernel/time.c
160	setrlimit	sys_setrlimit	kernel/sys.c
161	chroot	sys_chroot	fs/open.c
162	sync	sys_sync	fs/sync.c
163	acct	sys_acct	kernel/acct.c
164	settimeofday	sys_settimeofday	kernel/time.c
165	mount	sys_mount	fs/namespace.c
166	umount2	sys_umount	fs/namespace.c
167	swapon	sys_swapon	mm/swapfile.c
168	swapoff	sys_swapoff	mm/swapfile.c
169	reboot	sys_reboot	kernel/reboot.c
170	sethostname	sys_sethostname	kernel/sys.c
171	setdomainname	sys_setdomainname	kernel/sys.c
172	iopl	stub_iopl	arch/x86/kernel/ioport.c
173	ioperm	sys_ioperm	arch/x86/kernel/ioport.c
174	create_module	 	NOT IMPLEMENTED
175	init_module	sys_init_module	kernel/module.c
176	delete_module	sys_delete_module	kernel/module.c
177	get_kernel_syms	 	NOT IMPLEMENTED
178	query_module	 	NOT IMPLEMENTED
179	quotactl	sys_quotactl	fs/quota/quota.c
180	nfsservctl	 	NOT IMPLEMENTED
181	getpmsg	 	NOT IMPLEMENTED
182	putpmsg	 	NOT IMPLEMENTED
183	afs_syscall	 	NOT IMPLEMENTED
184	tuxcall	 	NOT IMPLEMENTED
185	security	 	NOT IMPLEMENTED
186	gettid	sys_gettid	kernel/sys.c
187	readahead	sys_readahead	mm/readahead.c
188	setxattr	sys_setxattr	fs/xattr.c
189	lsetxattr	sys_lsetxattr	fs/xattr.c
190	fsetxattr	sys_fsetxattr	fs/xattr.c
191	getxattr	sys_getxattr	fs/xattr.c
192	lgetxattr	sys_lgetxattr	fs/xattr.c
193	fgetxattr	sys_fgetxattr	fs/xattr.c
194	listxattr	sys_listxattr	fs/xattr.c
195	llistxattr	sys_llistxattr	fs/xattr.c
196	flistxattr	sys_flistxattr	fs/xattr.c
197	removexattr	sys_removexattr	fs/xattr.c
198	lremovexattr	sys_lremovexattr	fs/xattr.c
199	fremovexattr	sys_fremovexattr	fs/xattr.c
200	tkill	sys_tkill	kernel/signal.c
201	time	sys_time	kernel/time.c
202	futex	sys_futex	kernel/futex.c
203	sched_setaffinity	sys_sched_setaffinity	kernel/sched/core.c
204	sched_getaffinity	sys_sched_getaffinity	kernel/sched/core.c
205	set_thread_area	 	arch/x86/kernel/tls.c
206	io_setup	sys_io_setup	fs/aio.c
207	io_destroy	sys_io_destroy	fs/aio.c
208	io_getevents	sys_io_getevents	fs/aio.c
209	io_submit	sys_io_submit	fs/aio.c
210	io_cancel	sys_io_cancel	fs/aio.c
211	get_thread_area	 	arch/x86/kernel/tls.c
212	lookup_dcookie	sys_lookup_dcookie	fs/dcookies.c
213	epoll_create	sys_epoll_create	fs/eventpoll.c
214	epoll_ctl_old	 	NOT IMPLEMENTED
215	epoll_wait_old	 	NOT IMPLEMENTED
216	remap_file_pages	sys_remap_file_pages	mm/fremap.c
217	getdents64	sys_getdents64	fs/readdir.c
218	set_tid_address	sys_set_tid_address	kernel/fork.c
219	restart_syscall	sys_restart_syscall	kernel/signal.c
220	semtimedop	sys_semtimedop	ipc/sem.c
221	fadvise64	sys_fadvise64	mm/fadvise.c
222	timer_create	sys_timer_create	kernel/posix-timers.c
223	timer_settime	sys_timer_settime	kernel/posix-timers.c
224	timer_gettime	sys_timer_gettime	kernel/posix-timers.c
225	timer_getoverrun	sys_timer_getoverrun	kernel/posix-timers.c
226	timer_delete	sys_timer_delete	kernel/posix-timers.c
227	clock_settime	sys_clock_settime	kernel/posix-timers.c
228	clock_gettime	sys_clock_gettime	kernel/posix-timers.c
229	clock_getres	sys_clock_getres	kernel/posix-timers.c
230	clock_nanosleep	sys_clock_nanosleep	kernel/posix-timers.c
231	exit_group	sys_exit_group	kernel/exit.c
232	epoll_wait	sys_epoll_wait	fs/eventpoll.c
233	epoll_ctl	sys_epoll_ctl	fs/eventpoll.c
234	tgkill	sys_tgkill	kernel/signal.c
235	utimes	sys_utimes	fs/utimes.c
236	vserver	 	NOT IMPLEMENTED
237	mbind	sys_mbind	mm/mempolicy.c
238	set_mempolicy	sys_set_mempolicy	mm/mempolicy.c
239	get_mempolicy	sys_get_mempolicy	mm/mempolicy.c
240	mq_open	sys_mq_open	ipc/mqueue.c
241	mq_unlink	sys_mq_unlink	ipc/mqueue.c
242	mq_timedsend	sys_mq_timedsend	ipc/mqueue.c
243	mq_timedreceive	sys_mq_timedreceive	ipc/mqueue.c
244	mq_notify	sys_mq_notify	ipc/mqueue.c
245	mq_getsetattr	sys_mq_getsetattr	ipc/mqueue.c
246	kexec_load	sys_kexec_load	kernel/kexec.c
247	waitid	sys_waitid	kernel/exit.c
248	add_key	sys_add_key	security/keys/keyctl.c
249	request_key	sys_request_key	security/keys/keyctl.c
250	keyctl	sys_keyctl	security/keys/keyctl.c
251	ioprio_set	sys_ioprio_set	fs/ioprio.c
252	ioprio_get	sys_ioprio_get	fs/ioprio.c
253	inotify_init	sys_inotify_init	fs/notify/inotify/inotify_user.c
254	inotify_add_watch	sys_inotify_add_watch	fs/notify/inotify/inotify_user.c
255	inotify_rm_watch	sys_inotify_rm_watch	fs/notify/inotify/inotify_user.c
256	migrate_pages	sys_migrate_pages	mm/mempolicy.c
257	openat	sys_openat	fs/open.c
258	mkdirat	sys_mkdirat	fs/namei.c
259	mknodat	sys_mknodat	fs/namei.c
260	fchownat	sys_fchownat	fs/open.c
261	futimesat	sys_futimesat	fs/utimes.c
262	newfstatat	sys_newfstatat	fs/stat.c
263	unlinkat	sys_unlinkat	fs/namei.c
264	renameat	sys_renameat	fs/namei.c
265	linkat	sys_linkat	fs/namei.c
266	symlinkat	sys_symlinkat	fs/namei.c
267	readlinkat	sys_readlinkat	fs/stat.c
268	fchmodat	sys_fchmodat	fs/open.c
269	faccessat	sys_faccessat	fs/open.c
270	pselect6	sys_pselect6	fs/select.c
271	ppoll	sys_ppoll	fs/select.c
272	unshare	sys_unshare	kernel/fork.c
273	set_robust_list	sys_set_robust_list	kernel/futex.c
274	get_robust_list	sys_get_robust_list	kernel/futex.c
275	splice	sys_splice	fs/splice.c
276	tee	sys_tee	fs/splice.c
277	sync_file_range	sys_sync_file_range	fs/sync.c
278	vmsplice	sys_vmsplice	fs/splice.c
279	move_pages	sys_move_pages	mm/migrate.c
280	utimensat	sys_utimensat	fs/utimes.c
281	epoll_pwait	sys_epoll_pwait	fs/eventpoll.c
282	signalfd	sys_signalfd	fs/signalfd.c
283	timerfd_create	sys_timerfd_create	fs/timerfd.c
284	eventfd	sys_eventfd	fs/eventfd.c
285	fallocate	sys_fallocate	fs/open.c
286	timerfd_settime	sys_timerfd_settime	fs/timerfd.c
287	timerfd_gettime	sys_timerfd_gettime	fs/timerfd.c
288	accept4	sys_accept4	net/socket.c
289	signalfd4	sys_signalfd4	fs/signalfd.c
290	eventfd2	sys_eventfd2	fs/eventfd.c
291	epoll_create1	sys_epoll_create1	fs/eventpoll.c
292	dup3	sys_dup3	fs/file.c
293	pipe2	sys_pipe2	fs/pipe.c
294	inotify_init1	sys_inotify_init1	fs/notify/inotify/inotify_user.c
295	preadv	sys_preadv	fs/read_write.c
296	pwritev	sys_pwritev	fs/read_write.c
297	rt_tgsigqueueinfo	sys_rt_tgsigqueueinfo	kernel/signal.c
298	perf_event_open	sys_perf_event_open	kernel/events/core.c
299	recvmmsg	sys_recvmmsg	net/socket.c
300	fanotify_init	sys_fanotify_init	fs/notify/fanotify/fanotify_user.c
301	fanotify_mark	sys_fanotify_mark	fs/notify/fanotify/fanotify_user.c
302	prlimit64	sys_prlimit64	kernel/sys.c
303	name_to_handle_at	sys_name_to_handle_at	fs/fhandle.c
304	open_by_handle_at	sys_open_by_handle_at	fs/fhandle.c
305	clock_adjtime	sys_clock_adjtime	kernel/posix-timers.c
306	syncfs	sys_syncfs	fs/sync.c
307	sendmmsg	sys_sendmmsg	net/socket.c
308	setns	sys_setns	kernel/nsproxy.c
309	getcpu	sys_getcpu	kernel/sys.c
310	process_vm_readv	sys_process_vm_readv	mm/process_vm_access.c
311	process_vm_writev	sys_process_vm_writev	mm/process_vm_access.c
312	kcmp	sys_kcmp	kernel/kcmp.c
313	finit_module	sys_finit_module	kernel/module.c
<!-- more -->
在linux 查看32位的系统调用号
cat /usr/include/asm/unistd_32.h 

查看64位的系统调用号
cat /usr/include/asm/unistd_64.h 

系统调用执行的流程如下：

应用程序 代码调用系统调用( xyz )，该函数是一个包装系统调用的 库函数 ；
库函数 ( xyz )负责准备向内核传递的参数，并触发 软中断 以切换到内核；
CPU 被 软中断 打断后，执行 中断处理函数 ，即 系统调用处理函数 ( system_call )；
系统调用处理函数 调用 系统调用服务例程 ( sys_xyz )，真正开始处理该系统调用；
执行态切换
应用程序 ( application program )与 库函数 ( libc )之间， 系统调用处理函数 ( system call handler )与 系统调用服务例程 ( system call service routine )之间， 均是普通函数调用，应该不难理解。 而 库函数 与 系统调用处理函数 之间，由于涉及用户态与内核态的切换，要复杂一些。

Linux 通过 软中断 实现从 用户态 到 内核态 的切换。 用户态 与 内核态 是独立的执行流，因此在切换时，需要准备 执行栈 并保存 寄存器 。

内核实现了很多不同的系统调用(提供不同功能)，而 系统调用处理函数 只有一个。 因此，用户进程必须传递一个参数用于区分，这便是 系统调用号 ( system call number )。 在 Linux 中， 系统调用号 一般通过 eax 寄存器 来传递。

总结起来， 执行态切换 过程如下：

应用程序 在 用户态 准备好调用参数，执行 int 指令触发 软中断 ，中断号为 0x80 ；
CPU 被软中断打断后，执行对应的 中断处理函数 ，这时便已进入 内核态 ；
系统调用处理函数 准备 内核执行栈 ，并保存所有 寄存器 (一般用汇编语言实现)；
系统调用处理函数 根据 系统调用号 调用对应的 C 函数—— 系统调用服务例程 ；
系统调用处理函数 准备 返回值 并从 内核栈 中恢复 寄存器 ；
系统调用处理函数 执行 ret 指令切换回 用户态 ；
编程实践
下面，通过一个简单的程序，看看应用程序如何在 用户态 准备参数并通过 int 指令触发 软中断 以陷入 内核态 执行 系统调用 ：
hello_world-int.S
.section .rodata

msg:
    .ascii "Hello, world!\n"

.section .text

.global _start

_start:
    # call SYS_WRITE
    movl $4, %eax
    # push arguments
    movl $1, %ebx
    movl $msg, %ecx
    movl $14, %edx
    int $0x80

    # Call SYS_EXIT
    movl $1, %eax
    # push arguments
    movl $0, %ebx
    # initiate
    int $0x80
这是一个汇编语言程序，程序入口在 _start 标签之后。

第 12 行，准备 系统调用号 ：将常数 4 放进 寄存器 eax 。 系统调用号 4 代表 系统调用 SYS_write ， 我们将通过该系统调用向标准输出写入一个字符串。

第 14-16 行， 准备系统调用参数：第一个参数放进 寄存器 ebx ，第二个参数放进 ecx ， 以此类推。

write 系统调用需要 3 个参数：

文件描述符 ，标准输出文件描述符为 1 ；
写入内容(缓冲区)地址；
写入内容长度(字节数)；
第 17 行，执行 int 指令触发软中断 0x80 ，程序将陷入内核态并由内核执行系统调用。 系统调用执行完毕后，内核将负责切换回用户态，应用程序继续执行之后的指令( 从 20 行开始 )。

第 20-24 行，调用 exit 系统调用，以便退出程序。

注解
注意到，这里必须显式调用 exit 系统调用退出程序。 否则，程序将继续往下执行，最终遇到段错误( segmentation fault )！

读者可能很好奇——我在写 C 语言或者其他程序时，这个调用并不是必须的！

这是因为 C 库( libc )已经帮你把脏活累活都干了。

接下来，我们编译并执行这个汇编语言程序：

$ ls
hello_world-int.S
$ as -o hello_world-int.o hello_world-int.S
$ ls
hello_world-int.o  hello_world-int.S
$ ld -o hello_world-int hello_world-int.o
$ ls
hello_world-int  hello_world-int.o  hello_world-int.S
$ ./hello_world-int
Hello, world!
其实，将 系统调用号 和 调用参数 放进正确的 寄存器 并触发正确的 软中断 是个重复的麻烦事。 C 库已经把这脏累活给干了——试试 syscall 函数吧！

hello_world-syscall.c
#include <string.h>
#include <sys/syscall.h>
#include <unistd.h>

int main(int argc, char *argv[])
{
    char *msg = "Hello, world!\n";
    syscall(SYS_write, 1, msg, strlen(msg));

    return 0;
}


系统调用的基本原理
系统调用其实就是函数调用，只不过调用的是内核态的函数，但是我们知道，用户态是不能随意调用内核态的函数的，所以采用软中断的方式从用户态陷入到内核态。在内核中通过软中断0X80，系统会跳转到一个预设好的内核空间地址，它指向了系统调用处理程序（不要和系统调用服务例程混淆），这里指的是在entry.S文件中的system_call函数。就是说，所有的系统调用都会统一跳转到这个地址执行system_call函数，那么system_call函数如何派发它们到各自的服务例程呢？
我们知道每个系统调用都有一个系统调用号。同时，内核中一个有一个system_call_table数组，它是个函数指针数组，每个函数指针都指向了系统调用的服务例程。这个系统调用号是system_call_table的下标，用来指明到底要执行哪个系统调用。当int ox80的软中断执行时，系统调用号会被放进eax寄存器中，system_call函数可以读取eax寄存器获得系统调用号，将其乘以4得到偏移地址，以sys_call_table为基地址，基地址加上偏移地址就是应该执行的系统调用服务例程的地址。

系统调用的传参问题
当一个系统调用的参数个数大于5时（因为5个寄存器（eax, ebx, ecx, edx,esi）已经用完了），执行int 0x80指令时仍需将系统调用功能号保存在寄存器eax中，所不同的只是全部参数应该依次放在一块连续的内存区域里，同时在寄存器ebx中保存指向该内存区域的指针。系统调用完成之后，返回值扔将保存在寄存器eax中。由于只是需要一块连续的内存区域来保存系统调用的参数，因此完全可以像普通函数调用一样使用栈（stack）来传递系统调用所需要的参数。但是要注意一点，Linux采用的是c语言的调用模式，这就意味着所有参数必须以相反的顺序进栈，即最后一个参数先入栈，而第一个参数则最后入栈。如果采用栈来传递系统调用所需要的参数，在执行int 0x80指令时还应该将栈指针的当前值复制到寄存器ebx中。

1.添加系统调用的两种方法
方法一：编译内核法
拿到源码之后

修改内核的系统调用库函数 /usr/include/asm-generic/unistd.h，在这里面可以使用在syscall_table中没有用到的223号
添加系统调用号，让系统根据这个号，去找到syscall_table中的相应表项。在/arch/x86/kernel/syscall_table_32.s文件中添加系统调用号和调用函数的对应关系
接着就是my_syscall的实现了，在这里有两种方法：第一种方法是在kernel下自己新建一个目录添加自己的文件，但是要编写Makefile，而且要修改全局的Makefile。第二种比较简便的方法是，在kernel/sys.c中添加自己的服务函数，这样子不用修改Makefile.
以上准备工作做完之后，然后就要进行编译内核了，以下是我编译内核的一个过程。

1.make menuconfig (使用图形化的工具，更新.config文件)
2.make -j3 bzImage  （编译，-j3指的是同时使用3个cpu来编译，bzImage指的是更新grub，以便重新引导）
3.make modules   （对模块进行编译）
4.make modules_install（安装编译好的模块）
5.depmod  （进行依赖关系的处理）
6.reboot  （重启看到自己编译好的内核）
方法二：内核模块法
这种方法是采用系统调用拦截的一种方式，改变某一个系统调用号对应的服务程序为我们自己的编写的程序，从而相当于添加了我们自己的系统调用。具体实现，我们来看下：

2.通过内核模块实现添加系统调用
这种方法其实是系统调用拦截的实现。系统调用服务程序的地址是放在sys_call_table中通过系统调用号定位到具体的系统调用地址，那么我们通过编写内核模块来修改sys_call_table中的系统调用的地址为我们自己定义的函数的地址，就可以实现系统调用的拦截。
想法有了：那就是通过模块加载时，将系统调用表里面的那个系统调用号的那个系统调用号对应的系统调用服务例程改为我们自己实现的系统历程函数地址。但是内核已经不知道从哪个版本就不支持导出sys_call_table了。所以首先要获取sys_call_table的地址。
网上介绍了好多种方法来得到sys_call_table的地址，这里介绍最简单的一种方法

grep sys_call_table /boot/System.map-`uname -r`


这样就得到了sys_call_table的地址，但同时也得到了一个重要的信息，该符号对应的内存区域是只读的。所以我们要修改它，必须对它进行清楚写保护，这里介绍两种方法：
第一种方法：：我们知道控制寄存器cr0的第16位是写保护位。cr0的第16位置为了禁止超级权限，若清零了则允许超级权限往内核中写入数据，这样我们可以再写入之前，将那一位清零，使我们可以写入。然后写完后，又将那一位复原就行了。

unsigned int clear_and_return_cr0(void)
{
 unsigned int cr0 = 0;
 unsigned int ret;
 asm("movl %%cr0, %%eax":"=a"(cr0));
 ret = cr0;
 cr0 &= 0xfffeffff;
 asm("movl %%eax, %%cr0"::"a"(cr0));
 return ret;
}

void setback_cr0(unsigned int val) //读取val的值到eax寄存器，再将eax寄存器的值放入cr0中
{
 asm volatile("movl %%eax, %%cr0"::"a"(val));
}

 
第二种方法：通过设置虚拟地址对应的也表项的读写属性来设置：

int make_rw(unsigned long address)  
{  
        unsigned int level;  
        pte_t *pte = lookup_address(address, &level);//查找虚拟地址所在的页表地址  
        if (pte->pte & ~_PAGE_RW)  //设置页表读写属性
                pte->pte |=  _PAGE_RW;  
          
        return 0;  
}  
  
  
  
int make_ro(unsigned long address)  
{  
        unsigned int level;  
        pte_t *pte = lookup_address(address, &level);  
        pte->pte &= ~_PAGE_RW;  //设置只读属性
  
        return 0;  
} 
 
3.编写系统调用指定自己的系统调用
内核的初始化函数
在这里我使用系统空闲的223号空闲的系统调用号，你也可以换成其他系统调用的调用号，这样你在执行其他函数时，就会调用自己的写的函数的内容。

static int syscall_init_module(void)  
{  
        printk(KERN_ALERT "sys_call_table: 0x%p\n", sys_call_table);//获取系统调用表的地址
        orig_saved = (unsigned long *)(sys_call_table[223]);  //保存原有的223号的系统调用表的地址
        printk(KERN_ALERT "orig_saved : 0x%p\n", orig_saved );  
  
        make_rw((unsigned long)sys_call_table);  //修改页的写属性
        sys_call_table[223] = (unsigned long *)sys_mycall;  //将223号指向自己写的调用函数
        make_ro((unsigned long)sys_call_table);  
  
        return 0;  
}
自己的系统调用服务例程

asmlinkage long sys_mycall(void)
{
    printk(KERN_ALERT "i am hack syscall!\n");
    return 0;
}
移除内核模块时，将原有的系统调用进行还原

static void syscall_cleanup_module(void)  
{  
        printk(KERN_ALERT "Module syscall unloaded.\n");  
  
        make_rw((unsigned long)sys_call_table);  
        sys_call_table[223] = (unsigned long *) orig_saved ;   
        make_ro((unsigned long)sys_call_table);  
}
模块注册相关

module_init(syscall_init_module);  
module_exit(syscall_cleanup_module);  
  
MODULE_LICENSE("GPL");  
MODULE_DESCRIPTION("mysyscall");  
4.编写用户态的测试程序
  1 #include <linux/unistd.h>
  2 #include <syscall.h>
  3 #include <sys/types.h>
  4 #include <stdio.h>
  5 
  6 int main(void)
  7 {
  8     long pid = 0;
  9     pid = syscall(223);
 10     printf("%ld\n",pid);
 11     return 0;
 12 }
当我们使用syscall()这个函数去触发223的系统调用时，dmesg会发现我们自己写的服务函数的输出结果:

start_kernel从内核一启动的时候它会一直存在，这个就是0号进程，idle就是一个while0,一直在循环着，当系统没有进程需要执行的时候就调度到idle进程，我们在windows系统上会经常见到，叫做system idle,这是一个一直会存在的0号进程，然后呢就是0号进程创建了1号进程，这个init_process是我们的1号进程也就是第一个用户态进程，也就是它默认的就是根目录下的程序，也就是常会找默认路径下的程序来作为1号进程，1号进程接下来还创建了kthreadd来管理内核的一些线程，这样整个程序就启动起来了。

二.内核态、用户态、中断等概念的介绍

用户态和内核态的区分：

       现代计算机机中都有几种不同的指令级别，在高执行级别下，代码可以执行特权指令，访问任意的物理地址，这种CPU执行级别就对应着内核态，而在相应的低级别执行状态下，代码的掌控范围会受到限制，只能在对应级别允许的范围内活动。举例：Intrel x86 CPU有四种不同的执行级别0-3，Linux只使用了其中的0级和3级来分别表示内核态和用户态。操作系统让系统本身更为稳定的方式，这样程序员自己写的用户态代码很难把整个系统都给搞崩溃，内核的代码经过仔细的分析有专业的人员写的代码会更加健壮一些，整个程序会更加稳定一些，注意：这里所说的地址空间是逻辑地址而不是物理地址。

     用户态和内核态的很显著的区分就是：CS和EIP， CS寄存器的最低两位表明了当前代码的特权级别；CPU每条指令的读取都是通过CS:EIP这两个寄存器：其中CS是代码段选择寄存器，EIP是偏移量寄存器，上述判断由硬件完成。一般来说在Linux中，地址空间是一个显著的标志：0xc0000000以上的地址空间只能在内核态下访问，0xc00000000-0xbfffffff的地址空间在两种状态下都可以访问。

中断处理是从用户态进入到内核态的主要的方式：

      也可能是用户态程序执行的过程中调用了一个系统调用陷入了内核态当中，这个叫做trap,系统调用只是一种特殊的中断。
      寄存器上下文：
            ——从用户态切换到内核态的时候
                  必须保存用户态的寄存器上下文
                  要保存哪些？
                  保存在哪里？
      中断/int指令会在堆栈上保存一些寄存器的值
            ——如：用户态栈顶地址、当时的状态字、当时的cs:eip的值
      中断发生的后的第一件事就是保护现场，保护现场就是进入中断的程序保存需要用到的寄存器数据，恢复现场就是退出中断程序，恢复、保存寄存器的数据。
       
1
2
3
4
5
6
7
8
9
10
11
12
13
14
#define SAVE_ALL                                                      RESTORE_ALL
      "cld\n\t"\                                                      popl %ebx;
      "pushl %es\n\t"\                                                popl %ecx;
      "pushl %ds\n\t"\                                                popl %ebx;
      "pushl %eax\n\t"\                                               popl %edx;
      "pushl %ebp\n\t"\                                               popl %esi;
      "pushl %edi\n\t"\                                               popl %edi;
      "pushl %esi\n\t"\                                               popl %ebp; 
      "pushl %edx\n\t"\                                               popl %eax;
      "pushl %ecx\n\t"\                                               popl %ds;             
      "pushl %ebx\n\t"\                                               popl %es;
      "movl $" STR(_KERNEL_DS)",%edx\n\t"\                            addl $4,%esp;
      "movl %edx,%ds\n\t"\                                            iret;
      "movl %edx,%es\n\t"
      iret指令与中断信号(包括int指令),发生时的CPU的动作正好相反。

仔细分析一下中断处理的完整过程：

      interrupt(ex:int0x80)-save//发生系统调用     
       cs:eip/ss:esp/ss:esp/efalgs(current)to kernel stack,then load cs:eip(entry of a specific ISR)and ss:esp(point to kernel stack) //保存了cs:eip的值，保存了堆栈寄存器当前的栈顶，当前的标志寄存器，当前的保存到内核堆栈里面，当前加载了中断信号和系统调用相关联的中断服务程序的入口，把它加载到当前cs:eip的里面，同时也要把当前的esp和堆栈段也就是指向内核的信息也加载到cpu里面,这是由中断向量或者说是int指令完成的。这个时候开始内核态的代码。
SAVE_ALL
     -...//内核代码，完成中断服务，可能会发生进程调度  
    RESTOER_ALL                //完成之后再返回到原来的状态
    iret-pop    
     cs:eip/ss:eip/eflag from kernel stack
三.系统调用概述

系统调用的意义：
      操作系统为用户态进程与硬件设备进行交互提供了一组接口——系统调用:1.把用户从底层的硬件编程中解放了出来;2.极大地提高了系统的安全性使用户程序具有可移植性;用户程序与具体硬件已经被抽象接口所替代。
操作系统提供的API和系统调用的关系:
    API（应用程序编程接口）和系统调用：应用编程接口和系统调用是不同的：1.API只是一个函数定义；2.系统调用通过软中断向内核发出了一个明确的请求。
     Libc库定义的一些API引用了封装例成，唯一目的就是发布系统调用：1.一般每个系统调用对应一个封装例程；2.库函数再用这些封装例程定义出给用户的API（把系统调用封装成很多歌方便程序员使用的函数，不是每个API都对应一个特定的系统调用）
     API可能直接提供用户态的服务 如：一些数学函数 1.一个单独的API可能调用几个系统调用2.不同的API可能调用了同一个系统调用返回：大部分封装例程返回一个整数，其值的含义依赖于相应的系统调用-1在多数情况下表示内核不能满足进程的请求，Libc中定义的errno变量包含特定的出错码；下面一张图可以表示它们的工作过程：
    
      x,y,z就是函数，系统调用应用程序编程接口，这个应用程序编程接口里面封装了一个系统调用，这会触发一个0x80的一个中断，这个中断向量就对应着SYSTEM_CALL这个内核代码的入口的起点，sys_xyz是对应的中断服务程序，在中断服务程序执行完之后，它可能会ret_from_sys_call， 之后就经过这个函数进行处理， 这是一个进程调度的时机，如果没有发生系统调用的时机，如果没有发生系统调用，它就会ireturn可能就会返回到用户态接着执行。
我们要扒开系统调用的三层皮，我们讲这三层皮分别是：xyz、system_call和sys_xyz
     第一个就是API、第二个就是中断向量对应的这些也就是中断服务程序，中断向量对用的系统调用它有很多种不同的服务程序，比如sys_xyz,这就是三层皮。
      我们仔细看一下系统调用的服务历程：中断向量0x80与system_call绑定起来：
      当用户态进程调用一个系统调用时，CPU切换到内核态并开始执行第一个内核函数
      1.在Linux中是通过执行ini $0x80来执行系统调用的，这条汇编指令产生向量为128的编程异常
      2. Intel Pentium ll中引进了sysenter指令（快速系统调用）
系统调用号将xyz和sys_xyz关联起来了：
     传参：
      1.内核实现了很多不同的系统调用
      2.进程必须指明需要哪些系统调用，这需要传递一个系统调用号的参数，使用eax寄存器
     系统调用也需要输入输出参数，例如： 
     1.实际的值 2.用户态进程地址空间的变量的地址 3.甚至是包含指向用户态函数的指针的数据结构的地址
 system_call是linux中所有系统调用的入口点，每个系统调用至少有一个参数，即由eax传递的系统调用号
     2.一个应用程序调用fork(0封装例程，那么在执行int $0x80之前就把eax寄存器的值置为2（即_NR_fork)
     3.这个寄存器的设置是libc库中封装例程进行的，因此用户一般不关心系统调用号
     4.进入sys_call之后，立即将eax的值压入内核堆栈
寄存器传递参数有如下限制：
    1.每个参数的长度不能超过寄存器的长度，即32位
     2.在系统调用号eax之外，参数的个数不能超过6个（ebx,ecx,edx,esi,edi,ebp)
     超过6个怎么办？做一个把某个寄存器作为指针，指向一块内存，这样进入内核态之后可以访问所有内存空间，这就是系统调用的参数传递方式。
 四.库函数API和C代码中嵌入汇编代码两种方式系统调用
 首先选择一个系统调用，我选的是write，然后是用c语言写一段正常熟悉的系统调用代码，如下：
复制代码
#include<stdio.h>
#include<unistd.h>
int main(void)
{
    write(1,"hello world!5124\n",13);
    return 0;
 }
复制代码
   下面是我的命令行内容：



      其中，write有三个参数，第一个是表示写到终端屏幕上，1可以认为是屏幕的代号，第二个参数是写的内容，我是把hello world!写到屏幕上，并换行，第三个参数是写入的字符串长度，长度要大于等于要输出的字符串长度，否则只能输出字符串的一部分。程序执行结果如下：



然后是把这段代码转写为嵌入式汇编，嵌入式汇编的格式如下：
复制代码
_asm_(
     汇编语句模块:
     输出部分:函数调用时候的参数
     输入部分:函数调用时候的参数
     破坏描述部分):
     即格式为asm("statements":output_regs:input_regs:clobbered_regs);
可以看成是一个函数，有时候可以加一个_volatile_来选择让编译器优化或者不让编译器优化。
复制代码
代码如下：

复制代码
#include<stdio.h>
#include<unistd.h>
 int main()
{
   int a;
   char *ch="hello world!\n";
 
    asm volatile(
        "movl $0x4,%%eax\n\t"
        "movl $0x1,%%ebx\n\t"
        "movl $0x1,%%ecx\n\t"
        "movl $0xd,%%edx\n\t"
        "int $0x80\n\t"
        "movl %%eax,%0\n\t"
        :"=m"(a)
        :"s"(ch)
        );
      return 0;
 }
复制代码
 write系统调用有三个参数，分别是：写入的位置，内容和长度，所以转化为汇编对应的寄存器为eax（系统调用号为4）,ebx(参数),ecx(输出位置),edx(参数长度)


执行代码如下：

五.实验感想
      计算机科学中有一句话，任何计算机相关问题都可以通过加一个中间层来解决。操作系统的系统调用也是这样，system_call将api和系统函数连接起来，这样可以保证内核的安全，不会因为用户的失误操作而造成问题。操作系统为了安全，把一些重要的调用放在内核部分，这样只能通过触发系统调用来完成相应功能，这样可以保证内核的安全，但是不可避免的也造成了系统调用的消耗比较大。
   
   
   系统调用大致可分为六大类：进程控制（process control）、文件管理（file manipulation）、设备管理（device manipulation）、信息维护（information maintenance）、通信（communication） 和保护（protection）。
进程控制
执行程序应能正常（end()）或异常（abort()）停止执行。如果一个系统调用异常停止当前执行的程序，或者程序运行遇到问题并引起错误陷阱，那么有时转储内存到磁盘，并生成错误信息。内存信息转储到磁盘后，可用调试器（debugger）来确定问题原因（调试器为系统程序，用以帮助程序员发现和纠正错误（bug））。

无论是正常情况还是异常情况，操作系统都应将控制转到调用命令解释程序。命令解释程序接着读入下个命令。对于交互系统，命令解释程序只是简单读入下个命令，而假定用户会采取合适命令以处理错误。对于 GUI 系统，弹出窗口可用于提醒用户出错，并请求指引。对于批处理系统，命令解释程序通常终止整个作业，并继续下个作业。当出现错误时，有的系统可能允许特殊的恢复操作。

如果程序发现输入有错并且想要异常终止，那么它也可能需要定义错误级别。错误越严重，错误参数的级别也越高。通过将正常终止的错误级别定义为 0，可以把正常和异常终止放在一起处理。命令解释程序或后面的程序可以利用这种错误级别来自动确定下个动作。

执行一个程序的进程或作业可能需要加载（load()）和执行（execute()）另一个程序。这种功能允许命令解释程序来执行一个程序，该命令可以通过用户命令、鼠标点击或批处理命令来给定。一个有趣的问题是：加载程序终止时会将控制返回到哪里？与之相关的问题是：原有程序是否失去或保存了，或者可与新的程序一起并发执行？

如果新程序终止时控制返回到现有程序，那么必须保存现有程序的内存映像。因此，事实上创建了一个机制，以便一个程序调用另一个程序。如果两个程序并发继续，那么也就创建了一个新作业或进程，以便多道执行。通常，有一个系统调用专门用于这一目的（create_process() 或 submit_job()）。

如果创建了一个新的作业或进程或者一组作业或进程，那么我们应能控制执行。这种控制要能判定和重置进程或作业的属性，包括作业的优先级、最大允许执行时间等（get_ process_attributes() 和 set_process_attributes()）。如果发现创建的进程或作业不正确或者不再需要，那么也要能终止它（terminate_process()）。

创建了新的作业或进程后，可能要等待其执行完成，也可能要等待一定时间（wait_time()）。更有可能要等待某个事件的出现（wait_event()）。当事件出现时，作业或进程就会响应（signal_event()）。

通常，两个或多个进程会共享数据。为了确保共享数据的完整性，操作系统通常提供系统调用，以允许一个进程锁定（lock）共享数据。这样，在解锁之前，其他进程不能访问该数据。通常，这样的系统调用包括 acquire_lock() 和 release_lock()。这类系统调用用于协调并发进程，将在后续章节详细讨论。

进程和作业控制差异很大，这里通过两个例子加以说明：一个涉及单任务系统，另一个涉及多任务系统。

MS-DOS 操作系统是个单任务的系统，在计算机启动时它就运行一个命令解释程序（图 1a）。由于 MS-DOS 是单任务的，它采用了一种简单方法来执行程序而且不创建新进程。它加载程序到内存，并对自身进行改写，以便为新程序提供尽可能多的空间（图 1b）。

MS-DOS 执行状态
图 1 MS-DOS 执行状态

接着，它将指令指针设为程序的第一条指令。然后，运行程序，或者错误引起中断，或者程序执行系统调用来终止。无论如何，错误代码会保存在系统内存中以便以后使用。之后，命令解释程序中的尚未改写部分重新开始执行。它首先从磁盘中重新加载命令解释程序的其他部分。然后，命令解释程序会向用户或下个程序提供先前的错误代码。

FreeBSD（源于 Berkeley UNIX）是个多任务系统。在用户登录到系统后，用户所选的外壳就开始运行。这种外壳类似于 MS-DOS 外壳：按用户要求，接受命令并执行程序。不过，由于 FreeBSD 是多任务系统，命令解释程序在另一个程序执行，也可继续执行（图 2）。

运行多个程序的FreeBSD
图 2 运行多个程序的 FreeBSD

为了启动新进程，外壳执行系统调用 fork()。接着，所选程序通过系统调用 exec() 加载到内存，程序开始执行。根据命令执行方式，外壳要么等待进程完成，要么后台执行进程。对于后一种情况，外壳可以马上接受下个命令。当进程在后台运行时，它不能直接接受键盘输入，这是因为外壳已在使用键盘。因此 I/O 可通过文件或 GUI 来完成。

同时，用户可以让外壳执行其他程序，监视运行进程状态，改变程序优先级等。当进程完成时，它执行系统调用 exit() 以终止，并将 0 或非 0 的错误代码返回到调用进程。这一状态（或错误）代码可用于外壳或其他程序。后续章节将通过一个使用系统调用 fork() 和 exec() 的程序例子来讨论进程。
文件管理
下面，我们讨论一些有关文件的常用系统调用。

首先要能创建（create()）和删除（delete()）文件。这两个系统调用需要文件名称，还可能需要文件的一些属性。一旦文件创建后，就会打开（open()）并使用它，也会读(read()）、写（write()）或重定位（reposition()）（例如，重新回到文件开头，或直接跳到文件末尾）。最后，需要关闭（close()）文件，表示不再使用它了。

如果采用目录结构来组织文件系统的文件，那么也会需要同样的目录操作。另外，不管是文件还是目录，都要能对各种属性的值加以读取或设置。文件属性包括：文件名、文件类型、保护码、记账信息等。

针对这一功能，至少需要两个系统调用：获取文件属性（get_file_attributes()）和设置文件属性（set_file_attributes()）。有的操作系统还提供其他系统调用，如文件的移动（move()）和复制（copy()）。还有的操作系统通过代码或系统调用来完成这些 API 的功能。其他的操作系统可能通过系统程序来实现这些功能。如果系统程序可被其他程序调用，那么这些系统程序也就相当于 API。
设备管理
进程执行需要一些资源，如内存、磁盘驱动、所需文件等。如果有可用资源，那么系统可以允许请求，并将控制交给用户程序；否则，程序应等待，直到有足够可用的资源为止。

操作系统控制的各种资源可看作设备。有的设备是物理设备（如磁盘驱动），而其他的可当作抽象或虚拟的设备（如文件）。多用户系统要求先请求（request()）设备，以确保设备的专门使用。在设备用完后，要释放（release()）它。这些函数类似于文件的系统调用 open() 和 close()。其他操作系统对设备访问不加管理。这样带来的危害是潜在的设备争用以及可能发生的死锁，这将在后续章节中讨论。

在请求了设备（并得到）后，就能如同对文件一样，对设备进行读（read())、写（write())、重定位（reposition()）。事实上，I/O 设备和文件极为相似，以至于许多操作系统如 UNIX 都将这两者组合成文件-设备结构。这样，一组系统调用不但用于文件而且用于设备。有时，I/O 设备可通过特殊文件名、目录位置或文件属性来辨认。

用户界面可以让文件和设备看起来相似，即便内在系统调用不同。在设计、构建操作系统和用户界面时，这也是要加以考虑的。
信息维护
许多系统调用只不过用于在用户程序与操作系统之间传递信息。例如，大多数操作系统都有一个系统调用，以便返回当前的时间（time()）和日期（date()）。还有的系统调用可以返回系统的其他信息，如当前用户数、操作系统版本、内存或磁盘的可用量等。

还有一组系统调用帮助调试程序。许多系统都提供用于转储内存（dump()）的系统调用。对于调试，这很有用。程序 trace 可以列出程序执行时的所有系统调用。甚至微处理器都有一个 CPU 模式，称为单步（single step），即 CPU 每执行一条指令都会产生一个陷阱。调试器通常可以捕获到这些陷阱。

许多操作系统都提供程序的时间曲线（time profile），用于表示在特定位置或位置组合上的执行时间。时间曲线需要跟踪功能或固定定时中断。当定时中断出现时，就会记录程序计数器的值。如有足够频繁的定时中断，那么就可得到花在程序各个部分的时间统计信息。

再者，操作系统维护所有进程的信息，这些可通过系统调用来访问。通常，也可用系统调用重置进程信息（get_process_attributes() 和 set_process_attributes ()）。
通信
进程间通信的常用模型有两个：消息传递模型和共享内存模型。

对于消息传递模型（message-passing model），通信进程通过相互交换消息来传递信息。进程间的消息交换可以直接进行，也可以通过一个共同邮箱来间接进行。在开始通信前，应先建立连接。应知道另一个通信实体名称，它可能是同一系统的另一个进程，也可能是通过网络相连的另一计算机的进程。

每台网络计算机都有一个主机名（hostname），这是众所周知的。另外，每台主机也都有一个网络标识符，如IP地址。类似地，每个进程有进程名（process name），它通常可转换成标识符，以便操作系统引用。系统调用 get_hostid() 和 get_processid() 可以执行这类转换。这些标识符再传给通用系统调用 open() 和 close()（由文件系统提供），或专用系统调用 open_connection() 和 close_connection()，这取决于系统通信模型。

接受进程应通过系统调用 accept_connection() 来许可通信。大多数可接受连接的进程为专用的守护进程（daemon），即专用系统程序。它们执行系统调用 wait_for_connection()，在有连接时会被唤醒。通信源称为客户机（client），而接受后台程序称为服务器（server），它们通过系统调用 read_message() 和 write_message() 来交换消息。系统调用 close_connection() 终止通信。

对于共享内存模型（shared-memory model），进程通过系统调用 shared_memory_create() 和 shared_memory_attach() 创建共享内存，并访问其他进程拥有的内存区域。

注意，操作系统通常需要阻止一个进程访问另一个进程的内存。共享内存要求两个或多个进程都同意取消这一限制，这样它们就可通过读写共享区域的数据来交换信息。这种数据的类型是由这些进程来决定的，而不受操作系统的控制。进程也负责确保不会同时向同一个地方进行写操作。

上面讨论的两种模型常用于操作系统，而且大多数系统两种都实现了。消息传递对少量数据的交换很有用，因为没有冲突需要避免。与用于计算机间的共享内存相比，它也更容易实现。共享内存在通信方面具有高速和便捷的特点，因为当通信发生在同一计算机内时，它可以按内存传输速度来进行。不过，共享内存的进程在保护和同步方面有问题。
保护
保护提供控制访问计算机的系统资源的机制。过去，只有多用户的多道计算机系统才要考虑保护。随着网络和因特网的出现，所有计算机（从服务器到手持移动设备）都应考虑保护。

通常，提供保护的系统调用包括 set_permission() 和 get_permission()，用于设置资源（如文件和磁盘）权限。系统调用 allow_user() 和 deny_user() 分别用于允许和拒绝特定用户访问某些资源。   
      
