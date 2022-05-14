---
title: clock_event_device
layout: post
category: linux
author: 夏泽民
---
早期的内核版本中，进程的调度基于一个称之为tick的时钟滴答，通常使用时钟中断来定时地产生tick信号，每次tick定时中断都会进行进程的统计和调度，并对tick进行计数，记录在一个jiffies变量中，定时器的设计也是基于jiffies。这时候的内核代码中，几乎所有关于时钟的操作都是在machine级的代码中实现，很多公共的代码要在每个平台上重复实现。随后，随着通用时钟框架的引入，内核需要支持高精度的定时器，为此，通用时间框架为定时器硬件定义了一个标准的接口：clock_event_device，machine级的代码只要按这个标准接口实现相应的硬件控制功能，剩下的与平台无关的特性则统一由通用时间框架层来实现。
<!-- more -->
1.  时钟事件软件架构
本系列文章的第一节中，我们曾经讨论了时钟源设备：clocksource，现在又来一个时钟事件设备：clock_event_device，它们有何区别？看名字，好像都是给系统提供时钟的设备，实际上，clocksource不能被编程，没有产生事件的能力，它主要被用于timekeeper来实现对真实时间进行精确的统计，而clock_event_device则是可编程的，它可以工作在周期触发或单次触发模式，系统可以对它进行编程，以确定下一次事件触发的时间，clock_event_device主要用于实现普通定时器和高精度定时器，同时也用于产生tick事件，供给进程调度子系统使用。时钟事件设备与通用时间框架中的其他模块的关系如下图所示：
<img src="{{site.url}}{{site.baseurl}}/img/clock_event_device.png"/>
与clocksource一样，系统中可以存在多个clock_event_device，系统会根据它们的精度和能力，选择合适的clock_event_device对系统提供时钟事件服务。在smp系统中，为了减少处理器间的通信开销，基本上每个cpu都会具备一个属于自己的本地clock_event_device，独立地为该cpu提供时钟事件服务，smp中的每个cpu基于本地的clock_event_device，建立自己的tick_device，普通定时器和高精度定时器。
在软件架构上看，clock_event_device被分为了两层，与硬件相关的被放在了machine层，而与硬件无关的通用代码则被集中到了通用时间框架层，这符合内核对软件的设计需求，平台的开发者只需实现平台相关的接口即可，无需关注复杂的上层时间框架。
tick_device是基于clock_event_device的进一步封装，用于代替原有的时钟滴答中断，给内核提供tick事件，以完成进程的调度和进程信息统计，负载平衡和时间更新等操作。
struct clock_event_device
时钟事件设备的核心数据结构是clock_event_device结构，它代表着一个时钟硬件设备，该设备就好像是一个具有事件触发能力（通常就是指中断）的clocksource，它不停地计数，当计数值达到预先编程设定的数值那一刻，会引发一个时钟事件中断，继而触发该设备的事件处理回调函数，以完成对时钟事件的处理。clock_event_device结构的定义如下：

struct clock_event_device {
	void			(*event_handler)(struct clock_event_device *);
	int			(*set_next_event)(unsigned long evt,
						  struct clock_event_device *);
	int			(*set_next_ktime)(ktime_t expires,
						  struct clock_event_device *);
	ktime_t			next_event;
	u64			max_delta_ns;
	u64			min_delta_ns;
	u32			mult;
	u32			shift;
	enum clock_event_mode	mode;
	unsigned int		features;
	unsigned long		retries;
 
	void			(*broadcast)(const struct cpumask *mask);
	void			(*set_mode)(enum clock_event_mode mode,
					    struct clock_event_device *);
	unsigned long		min_delta_ticks;
	unsigned long		max_delta_ticks;
 
	const char		*name;
	int			rating;
	int			irq;
	const struct cpumask	*cpumask;
	struct list_head	list;
} ____cacheline_aligned;

event_handler  该字段是一个回调函数指针，通常由通用框架层设置，在时间中断到来时，machine底层的的中断服务程序会调用该回调，框架层利用该回调实现对时钟事件的处理。
set_next_event  设置下一次时间触发的时间，使用类似于clocksource的cycle计数值（离现在的cycle差值）作为参数。

set_next_ktime  设置下一次时间触发的时间，直接使用ktime时间作为参数。

max_delta_ns  可设置的最大时间差，单位是纳秒。

min_delta_ns  可设置的最小时间差，单位是纳秒。

mult shift  与clocksource中的类似，只不过是用于把纳秒转换为cycle。

mode  该时钟事件设备的工作模式，两种主要的工作模式分别是：

CLOCK_EVT_MODE_PERIODIC  周期触发模式，设置后按给定的周期不停地触发事件；
CLOCK_EVT_MODE_ONESHOT  单次触发模式，只在设置好的触发时刻触发一次；
set_mode  函数指针，用于设置时钟事件设备的工作模式。

rating  表示该设备的精度等级。

list  系统中注册的时钟事件设备用该字段挂在全局链表变量clockevent_devices上。
