---
title: timekeeper
layout: post
category: linux
author: 夏泽民
---
1.  时间的种类
内核管理着多种时间，它们分别是：

RTC时间
wall time：墙上时间
monotonic time
raw monotonic time
boot time：总启动时间

RTC时间  在PC中，RTC时间又叫CMOS时间，它通常由一个专门的计时硬件来实现，软件可以读取该硬件来获得年月日、时分秒等时间信息，而在嵌入式系统中，有使用专门的RTC芯片，也有直接把RTC集成到Soc芯片中，读取Soc中的某个寄存器即可获取当前时间信息。一般来说，RTC是一种可持续计时的，也就是说，不管系统是否上电，RTC中的时间信息都不会丢失，计时会一直持续进行，硬件上通常使用一个后备电池对RTC硬件进行单独的供电。因为RTC硬件的多样性，开发者需要为每种RTC时钟硬件提供相应的驱动程序，内核和用户空间通过驱动程序访问RTC硬件来获取或设置时间信息。

xtime  xtime和RTC时间一样，都是人们日常所使用的墙上时间，只是RTC时间的精度通常比较低，大多数情况下只能达到毫秒级别的精度，如果是使用外部的RTC芯片，访问速度也比较慢，为此，内核维护了另外一个wall time时间：xtime，取决于用于对xtime计时的clocksource，它的精度甚至可以达到纳秒级别，因为xtime实际上是一个内存中的变量，它的访问速度非常快，内核大部分时间都是使用xtime来获得当前时间信息。xtime记录的是自1970年1月1日24时到当前时刻所经历的纳秒数。

monotonic time  该时间自系统开机后就一直单调地增加，它不像xtime可以因用户的调整时间而产生跳变，不过该时间不计算系统休眠的时间，也就是说，系统休眠时，monotoic时间不会递增。

raw monotonic time  该时间与monotonic时间类似，也是单调递增的时间，唯一的不同是：raw monotonic time“更纯净”，他不会受到NTP时间调整的影响，它代表着系统独立时钟硬件对时间的统计。

boot time  与monotonic时间相同，不过会累加上系统休眠的时间，它代表着系统上电后的总时间。

时间种类	精度（统计单位）	访问速度	累计休眠时间	受NTP调整的影响
RTC	低	慢	Yes	Yes
xtime	高	快	Yes	Yes
monotonic	高	快	No	Yes
raw monotonic	高	快	No	No
boot time	高	快	Yes	Yes

2.  struct timekeeper
内核用timekeeper结构来组织与时间相关的数据，它的定义如下：

struct timekeeper {
	struct clocksource *clock;    /* Current clocksource used for timekeeping. */
	u32	mult;    /* NTP adjusted clock multiplier */
	int	shift;	/* The shift value of the current clocksource. */
	cycle_t cycle_interval;	/* Number of clock cycles in one NTP interval. */
	u64	xtime_interval;	/* Number of clock shifted nano seconds in one NTP interval. */
	s64	xtime_remainder;	/* shifted nano seconds left over when rounding cycle_interval */
	u32	raw_interval;	/* Raw nano seconds accumulated per NTP interval. */
 
	u64	xtime_nsec;	/* Clock shifted nano seconds remainder not stored in xtime.tv_nsec. */
	/* Difference between accumulated time and NTP time in ntp
	 * shifted nano seconds. */
	s64	ntp_error;
	/* Shift conversion between clock shifted nano seconds and
	 * ntp shifted nano seconds. */
	int	ntp_error_shift;
 
	struct timespec xtime;	/* The current time */
 
	struct timespec wall_to_monotonic;
	struct timespec total_sleep_time;	/* time spent in suspend */
	struct timespec raw_time;	/* The raw monotonic time for the CLOCK_MONOTONIC_RAW posix clock. */
 
	ktime_t offs_real;	/* Offset clock monotonic -> clock realtime */
 
	ktime_t offs_boot;	/* Offset clock monotonic -> clock boottime */
 
	seqlock_t lock;	/* Seqlock for all timekeeper values */
};
其中的xtime字段就是上面所说的墙上时间，它是一个timespec结构的变量，它记录了自1970年1月1日以来所经过的时间，因为是timespec结构，所以它的精度可以达到纳秒级，当然那要取决于系统的硬件是否支持这一精度。
内核除了用xtime表示墙上的真实时间外，还维护了另外一个时间：monotonic time，可以把它理解为自系统启动以来所经过的时间，该时间只能单调递增，可以理解为xtime虽然正常情况下也是递增的，但是毕竟用户可以主动向前或向后调整墙上时间，从而修改xtime值。但是monotonic时间不可以往后退，系统启动后只能不断递增。奇怪的是，内核并没有直接定义一个这样的变量来记录monotonic时间，而是定义了一个变量wall_to_monotonic，记录了墙上时间和monotonic时间之间的偏移量，当需要获得monotonic时间时，把xtime和wall_to_monotonic相加即可，因为默认启动时monotonic时间为0，所以实际上wall_to_monotonic的值是一个负数，它和xtime同一时间被初始化，请参考timekeeping_init函数。

计算monotonic时间要去除系统休眠期间花费的时间，内核用total_sleep_time记录休眠的时间，每次休眠醒来后重新累加该时间，并调整wall_to_monotonic的值，使其在系统休眠醒来后，monotonic时间不会发生跳变。因为wall_to_monotonic值被调整。所以如果想获取boot time，需要加入该变量的值
<!-- more -->
