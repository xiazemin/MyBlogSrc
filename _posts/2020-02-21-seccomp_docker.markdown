---
title: Docker容器内部使用gdb进行debug
layout: post
category: docker
author: 夏泽民
---
安全计算模式（secure computing mode，seccomp）是 Linux 内核功能，可以使用它来限制容器内可用的操作。


Docker 的默认 seccomp 配置文件是一个白名单，它指定了允许的调用。

下表列出了由于不在白名单而被有效阻止的重要（但不是全部）系统调用。该表包含每个系统调用被阻止的原因。
<!-- more -->
Syscall	Description
acct	Accounting syscall which could let containers disable their own resource limits or process accounting. Also gated by CAP_SYS_PACCT.
add_key	Prevent containers from using the kernel keyring, which is not namespaced.
adjtimex	Similar to clock_settime and settimeofday, time/date is not namespaced. Also gated by CAP_SYS_TIME.
bpf	Deny loading potentially persistent bpf programs into kernel, already gated by CAP_SYS_ADMIN.
clock_adjtime	Time/date is not namespaced. Also gated by CAP_SYS_TIME.
clock_settime	Time/date is not namespaced. Also gated by CAP_SYS_TIME.
clone	Deny cloning new namespaces. Also gated by CAP_SYS_ADMIN for CLONE_* flags, except CLONE_USERNS.
create_module	Deny manipulation and functions on kernel modules. Obsolete. Also gated by CAP_SYS_MODULE.
delete_module	Deny manipulation and functions on kernel modules. Also gated by CAP_SYS_MODULE.
finit_module	Deny manipulation and functions on kernel modules. Also gated by CAP_SYS_MODULE.
get_kernel_syms	Deny retrieval of exported kernel and module symbols. Obsolete.
get_mempolicy	Syscall that modifies kernel memory and NUMA settings. Already gated by CAP_SYS_NICE.
init_module	Deny manipulation and functions on kernel modules. Also gated by CAP_SYS_MODULE.
ioperm	Prevent containers from modifying kernel I/O privilege levels. Already gated by CAP_SYS_RAWIO.
iopl	Prevent containers from modifying kernel I/O privilege levels. Already gated by CAP_SYS_RAWIO.
kcmp	Restrict process inspection capabilities, already blocked by dropping CAP_PTRACE.
kexec_file_load	Sister syscall of kexec_load that does the same thing, slightly different arguments. Also gated by CAP_SYS_BOOT.
kexec_load	Deny loading a new kernel for later execution. Also gated by CAP_SYS_BOOT.
keyctl	Prevent containers from using the kernel keyring, which is not namespaced.
lookup_dcookie	Tracing/profiling syscall, which could leak a lot of information on the host. Also gated by CAP_SYS_ADMIN.
mbind	Syscall that modifies kernel memory and NUMA settings. Already gated by CAP_SYS_NICE.
mount	Deny mounting, already gated by CAP_SYS_ADMIN.
move_pages	Syscall that modifies kernel memory and NUMA settings.
name_to_handle_at	Sister syscall to open_by_handle_at. Already gated by CAP_SYS_NICE.
nfsservctl	Deny interaction with the kernel nfs daemon. Obsolete since Linux 3.1.
open_by_handle_at	Cause of an old container breakout. Also gated by CAP_DAC_READ_SEARCH.
perf_event_open	Tracing/profiling syscall, which could leak a lot of information on the host.
personality	Prevent container from enabling BSD emulation. Not inherently dangerous, but poorly tested, potential for a lot of kernel vulns.
pivot_root	Deny pivot_root, should be privileged operation.
process_vm_readv	Restrict process inspection capabilities, already blocked by dropping CAP_PTRACE.
process_vm_writev	Restrict process inspection capabilities, already blocked by dropping CAP_PTRACE.
ptrace	Tracing/profiling syscall, which could leak a lot of information on the host. Already blocked by dropping CAP_PTRACE.
query_module	Deny manipulation and functions on kernel modules. Obsolete.
quotactl	Quota syscall which could let containers disable their own resource limits or process accounting. Also gated by CAP_SYS_ADMIN.
reboot	Don’t let containers reboot the host. Also gated by CAP_SYS_BOOT.
request_key	Prevent containers from using the kernel keyring, which is not namespaced.
set_mempolicy	Syscall that modifies kernel memory and NUMA settings. Already gated by CAP_SYS_NICE.
setns	Deny associating a thread with a namespace. Also gated by CAP_SYS_ADMIN.
settimeofday	Time/date is not namespaced. Also gated by CAP_SYS_TIME.
socket, socketcall	Used to send or receive packets and for other socket operations. All socket and socketcall calls are blocked except communication domains AF_UNIX, AF_INET, AF_INET6, AF_NETLINK, and AF_PACKET.
stime	Time/date is not namespaced. Also gated by CAP_SYS_TIME.
swapon	Deny start/stop swapping to file/device. Also gated by CAP_SYS_ADMIN.
swapoff	Deny start/stop swapping to file/device. Also gated by CAP_SYS_ADMIN.
sysfs	Obsolete syscall.
_sysctl	Obsolete, replaced by /proc/sys.
umount	Should be a privileged operation. Also gated by CAP_SYS_ADMIN.
umount2	Should be a privileged operation. Also gated by CAP_SYS_ADMIN.
unshare	Deny cloning new namespaces for processes. Also gated by CAP_SYS_ADMIN, with the exception of unshare –user.
uselib	Older syscall related to shared libraries, unused for a long time.
userfaultfd	Userspace page fault handling, largely needed for process migration.
ustat	Obsolete syscall.
vm86	In kernel x86 real mode virtual machine. Also gated by CAP_SYS_ADMIN.
vm86old	In kernel x86 real mode virtual machine. Also gated by CAP_SYS_ADMIN.
其中gdb在进行进程debug时，会报错：

(gdb) attach 30721
Attaching to process 30721

ptrace: Operation not permitted.

原因就是因为ptrace被Docker默认禁止的问题。考虑到应用分析的需要，可以有以下几种方法解决：

1、关闭seccomp

docker run --security-opt seccomp=unconfined 

2、采用超级权限模式
docker run --privileged

3、仅开放ptrace限制

docker run --cap-add sys_ptrace 

当然从安全角度考虑，如只是想使用gdb进行debug的话，建议使用第三种。

小提示：

如果使用marathon进行docker的管理，应用JSON的修改在这里：

"parameters": [
        {
          "key": "cap-add",
          "value": "SYS_PTRACE"
        }
      ],


https://www.elastic.co/guide/en/beats/filebeat/6.4/linux-seccomp.html

https://github.com/moby/moby/blob/master/profiles/seccomp/default.json

Secure computing mode (seccomp) is a Linux kernel feature. You can use it to restrict the actions available within the container. The seccomp() system call operates on the seccomp state of the calling process. You can use this feature to restrict your application’s access.

This feature is available only if Docker has been built with seccomp and the kernel is configured with CONFIG_SECCOMP enabled. To check if your kernel supports seccomp:

$ grep CONFIG_SECCOMP= /boot/config-$(uname -r)
CONFIG_SECCOMP=y
Note: seccomp profiles require seccomp 2.2.1 which is not available on Ubuntu 14.04, Debian Wheezy, or Debian Jessie. To use seccomp on these distributions, you must download the latest static Linux binaries (rather than packages).

Pass a profile for a container
The default seccomp profile provides a sane default for running containers with seccomp and disables around 44 system calls out of 300+. It is moderately protective while providing wide application compatibility. The default Docker profile can be found here.

In effect, the profile is a whitelist which denies access to system calls by default, then whitelists specific system calls. The profile works by defining a defaultAction of SCMP_ACT_ERRNO and overriding that action only for specific system calls. The effect of SCMP_ACT_ERRNO is to cause a Permission Denied error. Next, the profile defines a specific list of system calls which are fully allowed, because their action is overridden to be SCMP_ACT_ALLOW. Finally, some specific rules are for individual system calls such as personality, and others, to allow variants of those system calls with specific arguments.

seccomp is instrumental for running Docker containers with least privilege. It is not recommended to change the default seccomp profile.

When you run a container, it uses the default profile unless you override it with the --security-opt option. For example, the following explicitly specifies a policy:

$ docker run --rm \
             -it \
             --security-opt seccomp=/path/to/seccomp/profile.json \
             hello-world
Significant syscalls blocked by the default profile
Docker’s default seccomp profile is a whitelist which specifies the calls that are allowed. The table below lists the significant (but not all) syscalls that are effectively blocked because they are not on the whitelist. The table includes the reason each syscall is blocked rather than white-listed.

Syscall	Description
acct	Accounting syscall which could let containers disable their own resource limits or process accounting. Also gated by CAP_SYS_PACCT.
add_key	Prevent containers from using the kernel keyring, which is not namespaced.
bpf	Deny loading potentially persistent bpf programs into kernel, already gated by CAP_SYS_ADMIN.
clock_adjtime	Time/date is not namespaced. Also gated by CAP_SYS_TIME.
clock_settime	Time/date is not namespaced. Also gated by CAP_SYS_TIME.
clone	Deny cloning new namespaces. Also gated by CAP_SYS_ADMIN for CLONE_* flags, except CLONE_USERNS.
create_module	Deny manipulation and functions on kernel modules. Obsolete. Also gated by CAP_SYS_MODULE.
delete_module	Deny manipulation and functions on kernel modules. Also gated by CAP_SYS_MODULE.
finit_module	Deny manipulation and functions on kernel modules. Also gated by CAP_SYS_MODULE.
get_kernel_syms	Deny retrieval of exported kernel and module symbols. Obsolete.
get_mempolicy	Syscall that modifies kernel memory and NUMA settings. Already gated by CAP_SYS_NICE.
init_module	Deny manipulation and functions on kernel modules. Also gated by CAP_SYS_MODULE.
ioperm	Prevent containers from modifying kernel I/O privilege levels. Already gated by CAP_SYS_RAWIO.
iopl	Prevent containers from modifying kernel I/O privilege levels. Already gated by CAP_SYS_RAWIO.
kcmp	Restrict process inspection capabilities, already blocked by dropping CAP_PTRACE.
kexec_file_load	Sister syscall of kexec_load that does the same thing, slightly different arguments. Also gated by CAP_SYS_BOOT.
kexec_load	Deny loading a new kernel for later execution. Also gated by CAP_SYS_BOOT.
keyctl	Prevent containers from using the kernel keyring, which is not namespaced.
lookup_dcookie	Tracing/profiling syscall, which could leak a lot of information on the host. Also gated by CAP_SYS_ADMIN.
mbind	Syscall that modifies kernel memory and NUMA settings. Already gated by CAP_SYS_NICE.
mount	Deny mounting, already gated by CAP_SYS_ADMIN.
move_pages	Syscall that modifies kernel memory and NUMA settings.
name_to_handle_at	Sister syscall to open_by_handle_at. Already gated by CAP_DAC_READ_SEARCH.
nfsservctl	Deny interaction with the kernel nfs daemon. Obsolete since Linux 3.1.
open_by_handle_at	Cause of an old container breakout. Also gated by CAP_DAC_READ_SEARCH.
perf_event_open	Tracing/profiling syscall, which could leak a lot of information on the host.
personality	Prevent container from enabling BSD emulation. Not inherently dangerous, but poorly tested, potential for a lot of kernel vulns.
pivot_root	Deny pivot_root, should be privileged operation.
process_vm_readv	Restrict process inspection capabilities, already blocked by dropping CAP_PTRACE.
process_vm_writev	Restrict process inspection capabilities, already blocked by dropping CAP_PTRACE.
ptrace	Tracing/profiling syscall, which could leak a lot of information on the host. Already blocked by dropping CAP_PTRACE. Blocked in Linux kernel versions before 4.8 to avoid seccomp bypass.
query_module	Deny manipulation and functions on kernel modules. Obsolete.
quotactl	Quota syscall which could let containers disable their own resource limits or process accounting. Also gated by CAP_SYS_ADMIN.
reboot	Don’t let containers reboot the host. Also gated by CAP_SYS_BOOT.
request_key	Prevent containers from using the kernel keyring, which is not namespaced.
set_mempolicy	Syscall that modifies kernel memory and NUMA settings. Already gated by CAP_SYS_NICE.
setns	Deny associating a thread with a namespace. Also gated by CAP_SYS_ADMIN.
settimeofday	Time/date is not namespaced. Also gated by CAP_SYS_TIME.
stime	Time/date is not namespaced. Also gated by CAP_SYS_TIME.
swapon	Deny start/stop swapping to file/device. Also gated by CAP_SYS_ADMIN.
swapoff	Deny start/stop swapping to file/device. Also gated by CAP_SYS_ADMIN.
sysfs	Obsolete syscall.
_sysctl	Obsolete, replaced by /proc/sys.
umount	Should be a privileged operation. Also gated by CAP_SYS_ADMIN.
umount2	Should be a privileged operation. Also gated by CAP_SYS_ADMIN.
unshare	Deny cloning new namespaces for processes. Also gated by CAP_SYS_ADMIN, with the exception of unshare --user.
uselib	Older syscall related to shared libraries, unused for a long time.
userfaultfd	Userspace page fault handling, largely needed for process migration.
ustat	Obsolete syscall.
vm86	In kernel x86 real mode virtual machine. Also gated by CAP_SYS_ADMIN.
vm86old	In kernel x86 real mode virtual machine. Also gated by CAP_SYS_ADMIN.
Run without the default seccomp profile
You can pass unconfined to run a container without the default seccomp profile.

$ docker run --rm -it --security-opt seccomp=unconfined debian:jessie \
    unshare --map-root-user --user sh -c whoami


https://docs.docker.com/engine/security/seccomp/
