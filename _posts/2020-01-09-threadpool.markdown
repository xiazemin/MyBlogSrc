---
title: 线程池
layout: post
category: linux
author: 夏泽民
---
0.首先什么是线程池？

线程池就是创建多个线程并且进行管理的容器。（线程池是个容器，可以创建线程和管理线程，并且给线程分配任务）



1.为什么要用线程池呢？

我们都知道，在Java中创建一个线程其实是一个很简单的事情，只要new Thread就可以了，但是这样做并不是一种很好的方式。那么为什么不好呢？



比如在一个项目里，全部都是用的new Thread的方式去启用线程，那么创建好Thread1，而1在运行的时候，创建了Thread2，等等等... 创建了10个线程的时候，

1，2，3都执行完毕了但是没被销毁，就可能导致无限制的新建线程，相互竞争，占用过多的系统资源，导致死锁以及OOM。

而且这些线程缺乏统一的管理的功能，也缺乏定期执行，定时执行，线程中断的功能。

说了这么多，那么线程池的好处是啥啊？

线程池的好处有以下3方面：

（1） 重用已经存在的线程，减少了线程的创建和销毁的开销。

（2）可有效控制最大并发的线程数，提高了系统资源的使用率避免很多竞争，避免了OOM啊 死锁啊等。

（3）可以提供定时和定期的执行方式，单线程，并发数量的控制等功能！

线程池可以使得对线程的管理更加方便，并且对高并发的控制尽在掌握。



2.四种Java线程池功能及分析

线程池都继承了ExecutorService的接口，所以他们都具有ExecutorService的生命周期方法：运行，关闭，终止；

因为继承了ExecutorService接口，所以它在被创建的时候就是处于运行状态，当线程没有任务执行时，就会进入关闭状态，只有调用了shutdown（）的时候才是正式的终止了这个线程池。

java通过Executors工厂类提供我们的线程池一共有4种：

fixedThreadPool() //启动固定线程数的线程池

CachedThreadPool() //按需分配的线程池

ScheduledThreadPoolExecutor()//定时，定期执行任务的线程池

ThreadPoolExecutor()//指定线程数的线程池。



我们先来说参数最多的线程，方便下面对参数的理解，

ThreadPoolExecutor（）

ThreadPoolExecutor是线程池的实现类之一，它主要是启动指定数量线程以及将任务添加到队列，并且将任务发送给空闲的线程。

我们来看一下拥有最多构造函数的线程池的构造函数吧：

public ThreadPoolExecutor(int corePoolSize, 
                              int maximumPoolSize, 
                              long keepAliveTime, 
                              TimeUnit unit, 
                              BlockingQueue<Runnable> workQueue, 
                              RejectedExecutionHandler handler)
corePoolSize ： 线程的核心线程数。

maximumPoolSize：线程允许的最大线程数。

keepAliveTime：当前线程池 线程总数大于核心线程数时，终止多余的空闲线程的时间。

Unit：keepAliveTime参数的时间单位

workQueue：队伍队列，如果线程池达到了核心线程数，并且其他线程都处于活动状态的时候，则将新任务放入此队列。

threadFactory：定制线程的创建过程

Handler：拒绝策略，当workQueue队满时，采取的措施。

ThreadPoolExecutor（）的作用主要是让我们更灵活的使用线程池，这里不详细的举例了。





fixedThreadPool() //启动固定线程数的线程池

在Android中，由于系统资源有限，所以我们最常用的就是fixedThreadPool()。

fixedThreadPool（int size） 就只有一个参数，size，就是线程池中最大可创建多少个线程

我们来看一下该线程池的实现过程



在实现上跟ThreadPoolExecutor类似fixedThreadPool简化了实现过程，把corePoolSize和maximumpoolSize的值都设为传入的size，并且设置keepAliveTime为0ms，然后采用的是LinkedBlockingQueue队列，这个队列是链式结构，所以是无边界的。可以容纳无数个任务。

总结：比如我创建3个线程的fixedThreadPool ，当3个都为活跃的时候，后面的任务会被加入无边界的链式队列（LinkedBlockingQueue），有空闲，就执行任务。

使用过程：

private static void fixedThreadPool(int size) throws ExecutionException, InterruptedException {
        ExecutorService executorService = Executors.newFixedThreadPool(size);
//        ExecutorService executorService2 = Executors.newSingleThreadExecutor();
        for (int i = 0; i < 10; i++) {
            Future<Integer> task = executorService.submit(new Callable<Integer>() {
                @Override
                public Integer call() throws Exception {
                    System.out.println("执行线程" + Thread.currentThread().getName());
                    return fibc(40);
                }
            });
            System.out.println("第"+i+"次计算,结果:"+task.get());
        }
结果：





CachedThreadPool() //按需分配的线程池

对比fixedThreadPool来说 CachedThreadPool就要更快一些，为什么快呢？

我们来看一下CachedThreadPool的实现源码：



可以看到，这里的  maximumPoolSize：线程允许的最大线程数。  为MAX_VALUE(无限大)，

也就是说，只要有任务并且其他线程都在活跃态，就会开启一个新的线程 （因为没有上限）

而当有空闲的线程的时候，就会去调用空闲线程执行任务。



使用过程：

 private static void newCacheThreadPool(){
        ExecutorService executorService = Executors.newCachedThreadPool();
        for (int i = 0; i < 100; i++) {
            Future task = executorService.submit(new Runnable() {
                @Override
                public void run() {
                    System.out.println(Thread.currentThread().getName()+"---->"+fibc(20));
                }
            });
        }
    }
结果：







ScheduledThreadPoolExecutor() 定时定期的线程池



创建ScheduledThreadPoolExecutor很简单，只要传入corePoolSize（）线程的核心线程数。就可以开启这个线程池

然后我们需要使用它的定时任务，就需要实现他的scheduleAtFixedRate方法 这个方法有4个参数：



第一个为要执行的任务。

第二个为每次任务执行的延迟，比如传入1，就会每隔1秒执行一次。

第三个为执行的周期

第四个为第二个参数的时间单位。



*无论实现几个scheduleAtFixedRate方法，他们都互不干扰。



实现如下

private static void newScheduledThreadPool(){
        ScheduledExecutorService executorService = Executors.newScheduledThreadPool(3);
        executorService.scheduleAtFixedRate(new Runnable() {
            @Override
            public void run() {
                System.out.println(Thread.currentThread().getName()+"----->" + fibc(40));
            }
        },5,3,TimeUnit.SECONDS);
 
        executorService.scheduleAtFixedRate(new Runnable() {
            @Override
            public void run() {
                System.out.println(Thread.currentThread().getName()+"----->" + fibc(5));
            }
        },1,1,TimeUnit.SECONDS);
 
 
        executorService.scheduleAtFixedRate(new Runnable() {
            @Override
            public void run() {
                System.out.println(Thread.currentThread().getName()+"----->" + fibc(1));
            }
        },1,2,TimeUnit.SECONDS);
    }
我这里实现了3个scheduleAtFixedRate方法，

第一个为5s一次，3个周期执行一次

第二个为1s一次，1个周期执行一次

第三个为1s一次，2个周期执行一次

结果如下：

总结：线程池的使用，可以更规范的管理和创建销毁线程，也可以更多样化的去使用线程，减小我们的系统开支。
<!-- more -->
Executor框架最核心的类是ThreadPoolExecutor，它是线程池的实现类，主要由下列4个组件构成。

    ·corePool：核心线程池的大小。

    ·maximumPool：最大线程池的大小。

    ·BlockingQueue：用来暂时保存任务的工作队列。

    ·RejectedExecutionHandler：当ThreadPoolExecutor已经关闭或ThreadPoolExecutor已经饱和时（达到了最大线程池大小且工作队列已满），execute()方法将要调用的Handler。

通过Executor框架的工具类Executors，可以创建3种类型的ThreadPoolExecutor。

    ·FixedThreadPool。

    ·SingleThreadExecutor。

    ·CachedThreadPool。 

下面将分别介绍这3种ThreadPoolExecutor。

FixedThreadPool被称为可重用固定线程数的线程池。下面是FixedThreadPool的源代码实现。

public static ExecutorService newFixedThreadPool(int nThreads) { 

    return new ThreadPoolExecutor(nThreads,nThreads,

    0L, TimeUnit.MILLISECONDS,new LinkedBlockingQueue());

}

FixedThreadPool的corePoolSize和maximumPoolSize都被设置为创建FixedThreadPool时指定的参数nThreads。当线程池中的线程数大于corePoolSize时，keepAliveTime为多余的空闲线程等待新任务的最长时间，超过这个时间后多余的线程将被终止。这里把keepAliveTime设置为0L，意味着多余的空闲线程会被立即终止。

FixedThreadPool的execute()方法的运行示意图如图所示。

FixedThreadPool的execute()的运行示意图
上图的说明如下。

    1）如果当前运行的线程数少于corePoolSize，则创建新线程来执行任务。

    2）在线程池完成预热之后（当前运行的线程数等于corePoolSize），将任务加入 LinkedBlockingQueue。

    3）线程执行完1中的任务后，会在循环中反复从LinkedBlockingQueue获取任务来执行。

FixedThreadPool使用无界队列LinkedBlockingQueue作为线程池的工作队列（队列的容量为 Integer.MAX_VALUE）。使用无界队列作为工作队列会对线程池带来如下影响。

    1）当线程池中的线程数达到corePoolSize后，新任务将在无界队列中等待，因此线程池中的线程数不会超过       corePoolSize。

    2）由于1，使用无界队列时maximumPoolSize将是一个无效参数。 

    3）由于1和2，使用无界队列时keepAliveTime将是一个无效参数。

    4）由于使用无界队列，运行中的FixedThreadPool（未执行方法shutdown()或 shutdownNow()）不会拒绝任务   （不会调用RejectedExecutionHandler.rejectedExecution方法）。

SingleThreadExecutor详解

SingleThreadExecutor是使用单个worker线程的Executor。下面是SingleThreadExecutor的源代码实现。

public static ExecutorService newSingleThreadExecutor() { 

    return new FinalizableDelegatedExecutorService(new ThreadPoolExecutor(1, 1,0L, TimeUnit.MILLISECONDS    ,new LinkedBlockingQueue()));

}

SingleThreadExecutor的corePoolSize和maximumPoolSize被设置为1。其他参数与 FixedThreadPool相同。SingleThreadExecutor使用无界队列LinkedBlockingQueue作为线程池的工作队列（队列的容量为Integer.MAX_VALUE）。SingleThreadExecutor使用无界队列作为工作队列对线程池带来的影响与FixedThreadPool相同，这里就不赘述了。

SingleThreadExecutor的运行示意图如图所示。

SingleThreadExecutor的execute()的运行示意
对上图的说明如下。

    1）如果当前运行的线程数少于corePoolSize（即线程池中无运行的线程），则创建一个新线程来执行任务。

    2）在线程池完成预热之后（当前线程池中有一个运行的线程），将任务加入Linked- BlockingQueue。

    3）线程执行完1中的任务后，会在一个无限循环中反复从LinkedBlockingQueue获取任务来执行。

CachedThreadPool详解

CachedThreadPool是一个会根据需要创建新线程的线程池。下面是创建CachedThreadPool的源代码。

public static ExecutorService newCachedThreadPool() {     return new ThreadPoolExecutor(0, Integer.MAX_VALUE,

    60L, TimeUnit.SECONDS,new SynchronousQueue());

}

CachedThreadPool的corePoolSize被设置为0，即corePool为空；maximumPoolSize被设置为 Integer.MAX_VALUE，即maximumPool是无界的。这里把keepAliveTime设置为60L，意味着 CachedThreadPool中的空闲线程等待新任务的最长时间为60秒，空闲线程超过60秒后将会被终止。

FixedThreadPool和SingleThreadExecutor使用无界队列LinkedBlockingQueue作为线程池的工作队列。CachedThreadPool使用没有容量的SynchronousQueue作为线程池的工作队列，但 CachedThreadPool的maximumPool是无界的。这意味着，如果主线程提交任务的速度高于 maximumPool中线程处理任务的速度时，CachedThreadPool会不断创建新线程。极端情况下， CachedThreadPool会因为创建过多线程而耗尽CPU和内存资源。

CachedThreadPool的execute()方法的执行示意图如图所示。


对上图的说明如下。

    1）首先执行SynchronousQueue.offer（Runnable task）。如果当前maximumPool中有空闲线程正在执行       SynchronousQueue.poll（keepAliveTime，TimeUnit.NANOSECONDS），那么主线程执行offer操作与空闲线       程执行的poll操作配对成功，主线程把任务交给空闲线程执行，execute()方法执行完成；否则执行下面的步骤

    2）当初始maximumPool为空，或者maximumPool中当前没有空闲线程时，将没有线程执       行 SynchronousQueue.poll（keepAliveTime，TimeUnit.NANOSECONDS）。这种情况下，步骤1）将失败。       此时CachedThreadPool会创建一个新线程执行任务，execute()方法执行完成。

    3）在步骤2）中新创建的线程将任务执行完后，会执行 SynchronousQueue.poll（keepAliveTime，       TimeUnit.NANOSECONDS）。这个poll操作会让空闲线程最多在SynchronousQueue中等待60秒钟。如果60       秒钟内主线程提交了一个新任务（主线程执行步骤1）），那么这个空闲线程将执行主线程提交的新任务；否       则，这个空闲线程将终止。由于空闲60秒的空闲线程会被终止，因此长时间保持空闲的CachedThreadPool不        会使用任何资源。

前面提到过，SynchronousQueue是一个没有容量的阻塞队列。每个插入操作必须等待另一个线程的对应移除操作，反之亦然。CachedThreadPool使用SynchronousQueue，把主线程提交的任务传递给空闲线程执行

什么是线程池
　　线程池就是以一个或多个线程[循环执行]多个应用逻辑的线程集合.

线程池的作用：
　　线程池作用就是限制系统中执行线程的数量。
　　根据系统的环境情况，可以自动或手动设置线程数量，达到运行的最佳效果；少了浪费了系统资源，多了造成系统拥挤效率不高。用线程池控制线程数量，其他线程排队等候。一个任务执行完毕，再从队列的中取最前面的任务开始执行。若队列中没有等待进程，线程池的这一资源处于等待。当一个新任务需要运行时，如果线程池中有等待的工作线程，就可以开始运行了；否则进入等待队列。

线程池接口:

复制代码
/**
 * 线程池方法定义
 */
public interface ThreadPools<Job extends Runnable>{

    /**
     * 执行一个任务(Job),这个Job必须实现Runnable
     * @param job
     */
    public void execute(Job job);

    /**
     * 关闭线程池
     */
    public void shutdown();

    /**
     * 增加工作者线程，即用来执行任务的线程
     * @param num
     */
    public void addWorkers(int num);

    /**
     * 减少工作者线程
     * @param num
     */
    public void removeWorker(int num);

    /**
     * 获取正在等待执行的任务数量
     */
    public int getJobSize();
}
复制代码
　　客户端可以通过execute(Job)方法将Job提交入线程池来执行，客户端完全不用等待Job的执行完成。除了execute(Job)方法以外，线程池接口提供了增加/减少工作者线程以及关闭线程池的方法。每个客户端提交的Job都会进入到一个工作队列中等待工作者线程的处理。

 

线程池默认实现
复制代码
import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.atomic.AtomicLong;

/**
 * 线程池
 * @param <Job>
 */
public class DefaultThreadPool<Job extends Runnable> implements ThreadPools<Job>{
    /**
     * 线程池维护工作者线程的最大数量
     */
    private static final int MAX_WORKER_NUMBERS=30;

    /**
     * 线程池维护工作者线程的最默认工作数量
     */
    private static final int DEFAULT_WORKER_NUMBERS = 5;

    /**
     * 线程池维护工作者线程的最小数量
     */
    private static final int MIN_WORKER_NUMBERS = 1;

    /**
     * 维护一个工作列表,里面加入客户端发起的工作
     */
    private final LinkedList<Job> jobs = new LinkedList<Job>();

    /**
     * 工作者线程的列表
     */
    private final List<Worker> workers = Collections.synchronizedList(new ArrayList<Worker>());

    /**
     * 工作者线程的数量
     */
    private int workerNum;
    /**
     *每个工作者线程编号生成
     */
    private AtomicLong threadNum = new AtomicLong();

    /**
     * 第一步:构造函数，用于初始化线程池
     * 首先判断初始化线程池的线程个数是否大于最大线程数，如果大于则线程池的默认初始化值为 DEFAULT_WORKER_NUMBERS
     */
    public DefaultThreadPool(int num){
        if (num > MAX_WORKER_NUMBERS) {
            this.workerNum =DEFAULT_WORKER_NUMBERS;
        } else {
            this.workerNum = num;
        }
        initializeWorkers(workerNum);
    }

    /**
     * 初始化每个工作者线程
     */
    private void initializeWorkers(int num) {
        for (int i = 0; i < num; i++) {
            Worker worker = new Worker();
            //添加到工作者线程的列表
            workers.add(worker);
            //启动工作者线程
            Thread thread = new Thread(worker);
            thread.start();
        }
    }

    /**
     * 执行一个任务(Job),这个Job必须实现Runnable
     * @param job
     */
    @Override
    public void execute(Job job) {
        //如果job为null，抛出空指针
        if (job==null){
            throw new NullPointerException();
        }
        //这里进行执行 TODO 当供大于求时候，考虑如何临时添加线程数
        if (job != null) {
            //根据线程的"等待/通知机制"这里必须对jobs加锁
            synchronized (jobs) {
                jobs.addLast(job);
                jobs.notify();
            }
        }

    }

    /**
     * 关闭线程池
     */
    @Override
    public void shutdown() {
        for (Worker worker:workers) {
            worker.shutdown();
        }
    }

    /**
     * 增加工作者线程，即用来执行任务的线程
     * @param num
     */
    @Override
    public void addWorkers(int num) {
        //加锁，防止该线程还没增加完成而下个线程继续增加导致工作者线程超过最大值
        synchronized (jobs) {
            if (num + this.workerNum > MAX_WORKER_NUMBERS) {
                num = MAX_WORKER_NUMBERS - this.workerNum;
            }
            initializeWorkers(num);
            this.workerNum += num;
        }
    }

    /**
     * 减少工作者线程
     * @param num
     */
    @Override
    public void removeWorker(int num) {
        synchronized (jobs) {
            if(num>=this.workerNum){
                throw new IllegalArgumentException("超过了已有的线程数量");
            }
            for (int i = 0; i < num; i++) {
                Worker worker = workers.get(i);
                if (worker != null) {
                    //关闭该线程并从列表中移除
                    worker.shutdown();
                    workers.remove(i);
                }
            }
            this.workerNum -= num;
        }

    }

    /**
     * 获取正在等待执行的任务数量
     */
    @Override
    public int getJobSize() {
        return workers.size();
    }

    /**
     * 消费者
     */
    class Worker implements Runnable {
        // 表示是否运行该worker
        private volatile boolean running = true;

        @Override
        public void run() {
            while (running) { //这个工作线程就一直循环，不停的检测是否还有任务区执行
                Job job = null;
                //线程的等待/通知机制
                synchronized (jobs) {
                    if (jobs.isEmpty()) {//工作队列为空
                        try {
                            jobs.wait();//线程等待唤醒
                        } catch (InterruptedException e) {
                            //感知到外部对该线程的中断操作，返回
                            Thread.currentThread().interrupt();
                            return;
                        }
                    }
                    // 取出一个job
                    job = jobs.removeFirst();
                }
                //执行job
                if (job != null) {
                    job.run();
                }
            }
        }

        /**
         * 终止该线程
         */
        public void shutdown() {
            running = false;
        }
    }

}
复制代码
　　从线程池的实现中可以看出，当客户端调用execute(Job)方法时，会不断地向任务列表jobs中添加Job，而每个工作者线程会不读的从jobs上获取Job来执行，当jobs为空时，工作者线程进入WAITING状态。

　　当添加一个Job后，对工作队列jobs调用其notify()方法来唤醒一个工作者线程。此处我们不调用notifyAll(),避免将等待队列中的线程全部移动到阻塞队列中而造成资源浪费。

　　线程池的本质就是使用了一个线程安全的工作队列连接工作者线程和客户端线程。客户端线程把任务放入工作队列后便返回，而工作者线程则不端的从工作队列中取出工作并执行。当工作队列为空时，工作者线程进入WAITING状态，当有客户端发送任务过来后会通过任意一个工作者线程，随着大量任务的提交，更多的工作者线程被唤醒。

Job实现
复制代码
public class Job implements Runnable{

    @Override
    public void run() {
        try {
            Thread.sleep(2500);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        System.out.println("当前线程名称:"+Thread.currentThread().getName()+";"+"job被指执行了");
    }
}
复制代码
测试程序
复制代码
public class WorkTest {
    public static void main(String[] args) {
        DefaultThreadPool defaultThreadPool = new DefaultThreadPool(10);
        for (int i=0;i<10000;i++){
            if (i==30){
                defaultThreadPool.addWorkers(10);
            }
            Job job = new Job();
            defaultThreadPool.execute(job);
        }
    }
}
复制代码
console结果：

   

 

 

使用Executors工具类创建线程池：

复制代码
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/**
 * 运行结果：总共只会创建5个线程， 开始执行五个线程，
 * 当五个线程都处于活动状态，再次提交的任务都会加入队列等到其他线程运行结束，当线程处于空闲状态时会被下一个任务复用
 *
 */
public class newFixedThreadPoolTest {
    public static void main(String[] args) {
        //Executors工厂类创建一个可重用固定线程数的线程池，以共享的无界队列方式来运行这些线程
        ExecutorService executorService = Executors.newFixedThreadPool(5);
        for(int i = 0; i < 20; i++) {
            Runnable synRunnable = new Runnable() {
                public void run() {
                    System.out.println(Thread.currentThread().getName());
                }
            };
            executorService.execute(synRunnable);
        }
    }
}

ExecutorService是线程池接口。它定义了4中线程池：

1. newCachedThreadPool：
底层：返回ThreadPoolExecutor实例，corePoolSize为0；maximumPoolSize为Integer.MAX_VALUE；keepAliveTime为60L；unit为TimeUnit.SECONDS；workQueue为SynchronousQueue(同步队列)
通俗：当有新任务到来，则插入到SynchronousQueue中，由于SynchronousQueue是同步队列，因此会在池中寻找可用线程来执行，若有可以线程则执行，若没有可用线程则创建一个线程来执行该任务；若池中线程空闲时间超过指定大小，则该线程会被销毁。
适用：执行很多短期异步的小程序或者负载较轻的服务器
2. newFixedThreadPool：
底层：返回ThreadPoolExecutor实例，接收参数为所设定线程数量nThread，corePoolSize为nThread，maximumPoolSize为nThread；keepAliveTime为0L(不限时)；unit为：TimeUnit.MILLISECONDS；WorkQueue为：new LinkedBlockingQueue<Runnable>() 无界阻塞队列
通俗：创建可容纳固定数量线程的池子，每隔线程的存活时间是无限的，当池子满了就不在添加线程了；如果池中的所有线程均在繁忙状态，对于新任务会进入阻塞队列中(无界的阻塞队列)，但是，在线程池空闲时，即线程池中没有可运行任务时，它不会释放工作线程，还会占用一定的系统资源。
适用：执行长期的任务，性能好很多
3. newSingleThreadExecutor:
底层：FinalizableDelegatedExecutorService包装的ThreadPoolExecutor实例，corePoolSize为1；maximumPoolSize为1；keepAliveTime为0L；unit为：TimeUnit.MILLISECONDS；workQueue为：new LinkedBlockingQueue<Runnable>() 无界阻塞队列
通俗：创建只有一个线程的线程池，且线程的存活时间是无限的；当该线程正繁忙时，对于新任务会进入阻塞队列中(无界的阻塞队列)
适用：一个任务一个任务执行的场景
4. NewScheduledThreadPool:
底层：创建ScheduledThreadPoolExecutor实例，corePoolSize为传递来的参数，maximumPoolSize为Integer.MAX_VALUE；keepAliveTime为0；unit为：TimeUnit.NANOSECONDS；workQueue为：new DelayedWorkQueue() 一个按超时时间升序排序的队列
通俗：创建一个固定大小的线程池，线程池内线程存活时间无限制，线程池可以支持定时及周期性任务执行，如果所有线程均处于繁忙状态，对于新任务会进入DelayedWorkQueue队列中，这是一种按照超时时间排序的队列结构
适用：周期性执行任务的场景

Java多线程实现方式主要有四种：继承Thread类、实现Runnable接口、实现Callable接口通过FutureTask包装器来创建Thread线程、使用ExecutorService、Callable、Future实现有返回结果的多线程。

其中前两种方式线程执行完后都没有返回值，后两种是带返回值的。

 

1、继承Thread类创建线程
Thread类本质上是实现了Runnable接口的一个实例，代表一个线程的实例。启动线程的唯一方法就是通过Thread类的start()实例方法。start()方法是一个native方法，它将启动一个新线程，并执行run()方法。这种方式实现多线程很简单，通过自己的类直接extend Thread，并复写run()方法，就可以启动新线程并执行自己定义的run()方法。例如：

复制代码
public class MyThread extends Thread {  
　　public void run() {  
　　 System.out.println("MyThread.run()");  
　　}  
}  
 
MyThread myThread1 = new MyThread();  
MyThread myThread2 = new MyThread();  
myThread1.start();  
myThread2.start();  
复制代码
2、实现Runnable接口创建线程
如果自己的类已经extends另一个类，就无法直接extends Thread，此时，可以实现一个Runnable接口，如下：

public class MyThread extends OtherClass implements Runnable {  
　　public void run() {  
　　 System.out.println("MyThread.run()");  
　　}  
}  
为了启动MyThread，需要首先实例化一个Thread，并传入自己的MyThread实例：

MyThread myThread = new MyThread();  
Thread thread = new Thread(myThread);  
thread.start();  
事实上，当传入一个Runnable target参数给Thread后，Thread的run()方法就会调用target.run()，参考JDK源代码：

public void run() {  
　　if (target != null) {  
　　 target.run();  
　　}  
}  
3、实现Callable接口通过FutureTask包装器来创建Thread线程

Callable接口（也只有一个方法）定义如下：   

public interface Callable<V>   { 
  V call（） throws Exception;   } 
复制代码
public class SomeCallable<V> extends OtherClass implements Callable<V> {

    @Override
    public V call() throws Exception {
        // TODO Auto-generated method stub
        return null;
    }

}
复制代码
复制代码
Callable<V> oneCallable = new SomeCallable<V>();   
//由Callable<Integer>创建一个FutureTask<Integer>对象：   
FutureTask<V> oneTask = new FutureTask<V>(oneCallable);   
//注释：FutureTask<Integer>是一个包装器，它通过接受Callable<Integer>来创建，它同时实现了Future和Runnable接口。 
  //由FutureTask<Integer>创建一个Thread对象：   
Thread oneThread = new Thread(oneTask);   
oneThread.start();   
//至此，一个线程就创建完成了。
复制代码
4、使用ExecutorService、Callable、Future实现有返回结果的线程

ExecutorService、Callable、Future三个接口实际上都是属于Executor框架。返回结果的线程是在JDK1.5中引入的新特征，有了这种特征就不需要再为了得到返回值而大费周折了。而且自己实现了也可能漏洞百出。

可返回值的任务必须实现Callable接口。类似的，无返回值的任务必须实现Runnable接口。

执行Callable任务后，可以获取一个Future的对象，在该对象上调用get就可以获取到Callable任务返回的Object了。

注意：get方法是阻塞的，即：线程无返回结果，get方法会一直等待。

再结合线程池接口ExecutorService就可以实现传说中有返回结果的多线程了。

下面提供了一个完整的有返回结果的多线程测试例子，在JDK1.5下验证过没问题可以直接使用。代码如下：

复制代码
import java.util.concurrent.*;  
import java.util.Date;  
import java.util.List;  
import java.util.ArrayList;  
  
/** 
* 有返回值的线程 
*/  
@SuppressWarnings("unchecked")  
public class Test {  
public static void main(String[] args) throws ExecutionException,  
    InterruptedException {  
   System.out.println("----程序开始运行----");  
   Date date1 = new Date();  
  
   int taskSize = 5;  
   // 创建一个线程池  
   ExecutorService pool = Executors.newFixedThreadPool(taskSize);  
   // 创建多个有返回值的任务  
   List<Future> list = new ArrayList<Future>();  
   for (int i = 0; i < taskSize; i++) {  
    Callable c = new MyCallable(i + " ");  
    // 执行任务并获取Future对象  
    Future f = pool.submit(c);  
    // System.out.println(">>>" + f.get().toString());  
    list.add(f);  
   }  
   // 关闭线程池  
   pool.shutdown();  
  
   // 获取所有并发任务的运行结果  
   for (Future f : list) {  
    // 从Future对象上获取任务的返回值，并输出到控制台  
    System.out.println(">>>" + f.get().toString());  
   }  
  
   Date date2 = new Date();  
   System.out.println("----程序结束运行----，程序运行时间【"  
     + (date2.getTime() - date1.getTime()) + "毫秒】");  
}  
}  
  
class MyCallable implements Callable<Object> {  
private String taskNum;  
  
MyCallable(String taskNum) {  
   this.taskNum = taskNum;  
}  
  
public Object call() throws Exception {  
   System.out.println(">>>" + taskNum + "任务启动");  
   Date dateTmp1 = new Date();  
   Thread.sleep(1000);  
   Date dateTmp2 = new Date();  
   long time = dateTmp2.getTime() - dateTmp1.getTime();  
   System.out.println(">>>" + taskNum + "任务终止");  
   return taskNum + "任务返回运行结果,当前任务时间【" + time + "毫秒】";  
}  
}  
复制代码
 

代码说明：
上述代码中Executors类，提供了一系列工厂方法用于创建线程池，返回的线程池都实现了ExecutorService接口。
public static ExecutorService newFixedThreadPool(int nThreads) 
创建固定数目线程的线程池。
public static ExecutorService newCachedThreadPool() 
创建一个可缓存的线程池，调用execute 将重用以前构造的线程（如果线程可用）。如果现有线程没有可用的，则创建一个新线程并添加到池中。终止并从缓存中移除那些已有 60 秒钟未被使用的线程。
public static ExecutorService newSingleThreadExecutor() 
创建一个单线程化的Executor。
public static ScheduledExecutorService newScheduledThreadPool(int corePoolSize) 
创建一个支持定时及周期性的任务执行的线程池，多数情况下可用来替代Timer类。

ExecutoreService提供了submit()方法，传递一个Callable，或Runnable，返回Future。如果Executor后台线程池还没有完成Callable的计算，这调用返回Future对象的get()方法，会阻塞直到计算完成。

1、什么是线程池：  java.util.concurrent.Executors提供了一个 java.util.concurrent.Executor接口的实现用于创建线程池

    多线程技术主要解决处理器单元内多个线程执行的问题，它可以显著减少处理器单元的闲置时间，增加处理器单元的吞吐能力。    
    假设一个服务器完成一项任务所需时间为：T1 创建线程时间，T2 在线程中执行任务的时间，T3 销毁线程时间。

   如果：T1 + T3 远大于 T2，则可以采用线程池，以提高服务器性能。

一个线程池包括以下四个基本组成部分：
1、线程池管理器（ThreadPool）：用于创建并管理线程池，包括 创建线程池，销毁线程池，添加新任务；
2、工作线程（PoolWorker）：线程池中线程，在没有任务时处于等待状态，可以循环的执行任务；
3、任务接口（Task）：每个任务必须实现的接口，以供工作线程调度任务的执行，它主要规定了任务的入口，任务执行完后的收尾工作，任务的执行状态等；
4、任务队列（taskQueue）：用于存放没有处理的任务。提供一种缓冲机制。

    线程池技术正是关注如何缩短或调整T1,T3时间的技术，从而提高服务器程序性能的。它把T1，T3分别安排在服务器程序的启动和结束的时间段或者一些空闲的时间段，这样在服务器程序处理客户请求时，不会有T1，T3的开销了。
    线程池不仅调整T1,T3产生的时间段，而且它还显著减少了创建线程的数目，看一个例子：
    假设一个服务器一天要处理50000个请求，并且每个请求需要一个单独的线程完成。在线程池中，线程数一般是固定的，所以产生线程总数不会超过线程池中线程的数目，而如果服务器不利用线程池来处理这些请求则线程总数为50000。一般线程池大小是远小于50000。所以利用线程池的服务器程序不会为了创建50000而在处理请求时浪费时间，从而提高效率。

2.常见线程池

①newSingleThreadExecutor

单个线程的线程池，即线程池中每次只有一个线程工作，单线程串行执行任务

②newFixedThreadPool(n)

固定数量的线程池，没提交一个任务就是一个线程，直到达到线程池的最大数量，然后后面进入等待队列直到前面的任务完成才继续执行

③newCacheThreadPool（推荐使用）

可缓存线程池，当线程池大小超过了处理任务所需的线程，那么就会回收部分空闲（一般是60秒无执行）的线程，当有任务来时，又智能的添加新线程来执行。

④newScheduledThreadPool

大小无限制的线程池，支持定时和周期性的执行线程

  java提供的线程池更加强大，相信理解线程池的工作原理，看类库中的线程池就不会感到陌生了。





 

 

Java线程池使用说明

一：简介
线程的使用在java中占有极其重要的地位，在jdk1.4极其之前的jdk版本中，关于线程池的使用是极其简陋的。在jdk1.5之后这一情况有了很大的改观。Jdk1.5之后加入了java.util.concurrent包，这个包中主要介绍java中线程以及线程池的使用。为我们在开发中处理线程的问题提供了非常大的帮助。

二：线程池
线程池的作用：

线程池作用就是限制系统中执行线程的数量。
     根据系统的环境情况，可以自动或手动设置线程数量，达到运行的最佳效果；少了浪费了系统资源，多了造成系统拥挤效率不高。用线程池控制线程数量，其他线程排队等候。一个任务执行完毕，再从队列的中取最前面的任务开始执行。若队列中没有等待进程，线程池的这一资源处于等待。当一个新任务需要运行时，如果线程池中有等待的工作线程，就可以开始运行了；否则进入等待队列。

为什么要用线程池:

1.减少了创建和销毁线程的次数，每个工作线程都可以被重复利用，可执行多个任务。

2.可以根据系统的承受能力，调整线程池中工作线线程的数目，防止因为消耗过多的内存，而把服务器累趴下(每个线程需要大约1MB内存，线程开的越多，消耗的内存也就越大，最后死机)。

Java里面线程池的顶级接口是Executor，但是严格意义上讲Executor并不是一个线程池，而只是一个执行线程的工具。真正的线程池接口是ExecutorService。

比较重要的几个类：

ExecutorService

真正的线程池接口。

ScheduledExecutorService

能和Timer/TimerTask类似，解决那些需要任务重复执行的问题。

ThreadPoolExecutor

ExecutorService的默认实现。

ScheduledThreadPoolExecutor

继承ThreadPoolExecutor的ScheduledExecutorService接口实现，周期性任务调度的类实现。

 

要配置一个线程池是比较复杂的，尤其是对于线程池的原理不是很清楚的情况下，很有可能配置的线程池不是较优的，因此在Executors类里面提供了一些静态工厂，生成一些常用的线程池。

1. newSingleThreadExecutor

创建一个单线程的线程池。这个线程池只有一个线程在工作，也就是相当于单线程串行执行所有任务。如果这个唯一的线程因为异常结束，那么会有一个新的线程来替代它。此线程池保证所有任务的执行顺序按照任务的提交顺序执行。

2.newFixedThreadPool

创建固定大小的线程池。每次提交一个任务就创建一个线程，直到线程达到线程池的最大大小。线程池的大小一旦达到最大值就会保持不变，如果某个线程因为执行异常而结束，那么线程池会补充一个新线程。

3. newCachedThreadPool

创建一个可缓存的线程池。如果线程池的大小超过了处理任务所需要的线程，

那么就会回收部分空闲（60秒不执行任务）的线程，当任务数增加时，此线程池又可以智能的添加新线程来处理任务。此线程池不会对线程池大小做限制，线程池大小完全依赖于操作系统（或者说JVM）能够创建的最大线程大小。

4.newScheduledThreadPool

创建一个大小无限的线程池。此线程池支持定时以及周期性执行任务的需求。

实例

1：newSingleThreadExecutor

package com.thread;
/*
 * 通过实现Runnable接口，实现多线程
 * Runnable类是有run()方法的；
 * 但是没有start方法
 * 参考：
 * 
http://blog.csdn.net/qq_31753145/article/details/50899119

 * */

public class MyThread extends Thread { 

    @Override
    public void run() {
        // TODO Auto-generated method stub
//        super.run();
    System.out.println(Thread.currentThread().getName()+"正在执行....");
    } 
}
package com.thread;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/*
 * 通过实现Runnable接口，实现多线程
 * Runnable类是有run()方法的；
 * 但是没有start方法
 * 参考：
 * 
http://blog.csdn.net/qq_31753145/article/details/50899119

 * */

public class singleThreadExecutorTest{
     
    public static void main(String[] args) {
        ExecutorService pool=Executors.newSingleThreadExecutor();
        
        //创建实现了Runnable接口对象，Thread对象当然也实现了Runnable接口;      
        Thread t1=new MyThread();       
        Thread t2=new MyThread();       
        Thread t3=new MyThread();       
        Thread t4=new MyThread();       
        Thread t5=new MyThread();
        
        //将线程放到池中执行；     
        pool.execute(t1);
        pool.execute(t2);
        pool.execute(t3);
        pool.execute(t4);
        pool.execute(t5);
        
        //关闭线程池       
       pool.shutdown(); 
    }
}
结果：

pool-1-thread-1正在执行....
pool-1-thread-1正在执行....
pool-1-thread-1正在执行....
pool-1-thread-1正在执行....
pool-1-thread-1正在执行....
2：newFixedThreadPool

package com.thread;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/*
 * 通过实现Runnable接口，实现多线程
 * Runnable类是有run()方法的；
 * 但是没有start方法
 * 参考：
 * 
http://blog.csdn.net/qq_31753145/article/details/50899119

 * */

public class fixedThreadExecutorTest{
     

    public static void main(String[] args) {
        ExecutorService pool=Executors.newFixedThreadPool(2);
        
        //创建实现了Runnable接口对象，Thread对象当然也实现了Runnable接口;              
        Thread t1=new MyThread();      
        Thread t2=new MyThread();        
        Thread t3=new MyThread();      
        Thread t4=new MyThread();      
        Thread t5=new MyThread();
       
        //将线程放到池中执行；      
        pool.execute(t1);
        pool.execute(t2);
        pool.execute(t3);
        pool.execute(t4);
        pool.execute(t5);
        
        //关闭线程池       
       pool.shutdown();
    }

}
结果：

pool-1-thread-1正在执行....
pool-1-thread-1正在执行....
pool-1-thread-1正在执行....
pool-1-thread-1正在执行....
pool-1-thread-2正在执行....
3：newCachedThreadPool

package com.thread;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/*
 * 通过实现Runnable接口，实现多线程
 * Runnable类是有run()方法的；
 * 但是没有start方法
 * 参考：
 * 
http://blog.csdn.net/qq_31753145/article/details/50899119

 * */

public class cachedThreadExecutorTest{

    public static void main(String[] args) {
  
        ExecutorService pool=Executors.newCachedThreadPool();
        
        //创建实现了Runnable接口对象，Thread对象当然也实现了Runnable接口;
        Thread t1=new MyThread();    
        Thread t2=new MyThread();  
        Thread t3=new MyThread();    
        Thread t4=new MyThread();   
        Thread t5=new MyThread();
        
        //将线程放到池中执行；      
        pool.execute(t1);
        pool.execute(t2);
        pool.execute(t3);
        pool.execute(t4);
        pool.execute(t5);
        
        //关闭线程池     
       pool.shutdown();
    }

}
结果：

pool-1-thread-2正在执行....
pool-1-thread-1正在执行....
pool-1-thread-3正在执行....
pool-1-thread-4正在执行....
pool-1-thread-5正在执行....
4：newScheduledThreadPool

package com.thread;

import java.util.concurrent.ScheduledThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

/*
 * 通过实现Runnable接口，实现多线程
 * Runnable类是有run()方法的；
 * 但是没有start方法
 * 参考：
 * 
http://blog.csdn.net/qq_31753145/article/details/50899119

 * */

public class scheduledThreadExecutorTest{
    
    public static void main(String[] args) {
        // TODO Auto-generated method stub

       ScheduledThreadPoolExecutor exec =new ScheduledThreadPoolExecutor(1);
       exec.scheduleAtFixedRate(new Runnable(){//每隔一段时间就触发异常

        @Override
        public void run() {
            // TODO Auto-generated method stub
            //throw new RuntimeException();
            System.out.println("===================");
            
        }}, 1000, 5000, TimeUnit.MILLISECONDS);  
       
       exec.scheduleAtFixedRate(new Runnable(){//每隔一段时间打印系统时间，证明两者是互不影响的

        @Override
        public void run() {
            // TODO Auto-generated method stub
            System.out.println(System.nanoTime());
            
        }}, 1000, 2000, TimeUnit.MILLISECONDS);
    }
}
结果：

===================
23119318857491
23121319071841
23123319007891
===================
23125318176937
23127318190359
===================
23129318176148
23131318344312
23133318465896
===================
23135319645812
三：ThreadPoolExecutor详解
ThreadPoolExecutor的完整构造方法的签名是：ThreadPoolExecutor(int corePoolSize, int maximumPoolSize, long keepAliveTime, TimeUnit unit, BlockingQueue<Runnable> workQueue, ThreadFactory threadFactory, RejectedExecutionHandler handler) .

corePoolSize - 池中所保存的线程数，包括空闲线程。

maximumPoolSize-池中允许的最大线程数。

keepAliveTime - 当线程数大于核心时，此为终止前多余的空闲线程等待新任务的最长时间。

unit - keepAliveTime 参数的时间单位。

workQueue - 执行前用于保持任务的队列。此队列仅保持由 execute方法提交的 Runnable任务。

threadFactory - 执行程序创建新线程时使用的工厂。

handler - 由于超出线程范围和队列容量而使执行被阻塞时所使用的处理程序。

ThreadPoolExecutor是Executors类的底层实现。

在JDK帮助文档中，有如此一段话：

“强烈建议程序员使用较为方便的Executors工厂方法Executors.newCachedThreadPool()（无界线程池，可以进行自动线程回收）、Executors.newFixedThreadPool(int)（固定大小线程池）Executors.newSingleThreadExecutor()（单个后台线程）

它们均为大多数使用场景预定义了设置。”

下面介绍一下几个类的源码：

ExecutorService  newFixedThreadPool (int nThreads):固定大小线程池。

可以看到，corePoolSize和maximumPoolSize的大小是一样的（实际上，后面会介绍，如果使用无界queue的话maximumPoolSize参数是没有意义的），keepAliveTime和unit的设值表名什么？-就是该实现不想keep alive！最后的BlockingQueue选择了LinkedBlockingQueue，该queue有一个特点，他是无界的。

public static ExecutorService newFixedThreadPool(int nThreads) {   
           return new ThreadPoolExecutor(nThreads, nThreads,   
                                        0L, TimeUnit.MILLISECONDS,   
                                        new LinkedBlockingQueue<Runnable>());   
      }
ExecutorService  newSingleThreadExecutor()：单线程

public static ExecutorService newSingleThreadExecutor() {   
            return new FinalizableDelegatedExecutorService   
                 (new ThreadPoolExecutor(1, 1,   
                                       0L, TimeUnit.MILLISECONDS,   
                                      new LinkedBlockingQueue<Runnable>()));   
       }
 ExecutorService newCachedThreadPool()：无界线程池，可以进行自动线程回收

这个实现就有意思了。首先是无界的线程池，所以我们可以发现maximumPoolSize为big big。其次BlockingQueue的选择上使用SynchronousQueue。可能对于该BlockingQueue有些陌生，简单说：该QUEUE中，每个插入操作必须等待另一个线程的对应移除操作。

public static ExecutorService newCachedThreadPool() {   
             return new ThreadPoolExecutor(0, Integer.MAX_VALUE,   
                                         60L, TimeUnit.SECONDS,   
                                        new SynchronousQueue<Runnable>());   
    }
先从BlockingQueue<Runnable> workQueue这个入参开始说起。在JDK中，其实已经说得很清楚了，一共有三种类型的queue。

所有BlockingQueue 都可用于传输和保持提交的任务。可以使用此队列与池大小进行交互：

如果运行的线程少于 corePoolSize，则 Executor始终首选添加新的线程，而不进行排队。（如果当前运行的线程小于corePoolSize，则任务根本不会存放，添加到queue中，而是直接抄家伙（thread）开始运行）

如果运行的线程等于或多于 corePoolSize，则 Executor始终首选将请求加入队列，而不添加新的线程。

如果无法将请求加入队列，则创建新的线程，除非创建此线程超出 maximumPoolSize，在这种情况下，任务将被拒绝。

排队有三种通用策略：

直接提交。工作队列的默认选项是 SynchronousQueue，它将任务直接提交给线程而不保持它们。在此，如果不存在可用于立即运行任务的线程，则试图把任务加入队列将失败，因此会构造一个新的线程。此策略可以避免在处理可能具有内部依赖性的请求集时出现锁。直接提交通常要求无界 maximumPoolSizes 以避免拒绝新提交的任务。当命令以超过队列所能处理的平均数连续到达时，此策略允许无界线程具有增长的可能性。

无界队列。使用无界队列（例如，不具有预定义容量的 LinkedBlockingQueue）将导致在所有corePoolSize 线程都忙时新任务在队列中等待。这样，创建的线程就不会超过 corePoolSize。（因此，maximumPoolSize的值也就无效了。）当每个任务完全独立于其他任务，即任务执行互不影响时，适合于使用无界队列；例如，在 Web页服务器中。这种排队可用于处理瞬态突发请求，当命令以超过队列所能处理的平均数连续到达时，此策略允许无界线程具有增长的可能性。

有界队列。当使用有限的 maximumPoolSizes时，有界队列（如 ArrayBlockingQueue）有助于防止资源耗尽，但是可能较难调整和控制。队列大小和最大池大小可能需要相互折衷：使用大型队列和小型池可以最大限度地降低 CPU 使用率、操作系统资源和上下文切换开销，但是可能导致人工降低吞吐量。如果任务频繁阻塞（例如，如果它们是 I/O边界），则系统可能为超过您许可的更多线程安排时间。使用小型队列通常要求较大的池大小，CPU使用率较高，但是可能遇到不可接受的调度开销，这样也会降低吞吐量。  

BlockingQueue的选择。

例子一：使用直接提交策略，也即SynchronousQueue。

首先SynchronousQueue是无界的，也就是说他存数任务的能力是没有限制的，但是由于该Queue本身的特性，在某次添加元素后必须等待其他线程取走后才能继续添加。在这里不是核心线程便是新创建的线程，但是我们试想一样下，下面的场景。

我们使用一下参数构造ThreadPoolExecutor：

new ThreadPoolExecutor(   
               2, 3, 30, TimeUnit.SECONDS,    
               new  SynchronousQueue<Runnable>(),    
               new RecorderThreadFactory("CookieRecorderPool"),    
            new ThreadPoolExecutor.CallerRunsPolicy());
 当核心线程已经有2个正在运行.

此时继续来了一个任务（A），根据前面介绍的“如果运行的线程等于或多于 corePoolSize，则Executor始终首选将请求加入队列，而不添加新的线程。”,所以A被添加到queue中。
又来了一个任务（B），且核心2个线程还没有忙完，OK，接下来首先尝试1中描述，但是由于使用的SynchronousQueue，所以一定无法加入进去。
此时便满足了上面提到的“如果无法将请求加入队列，则创建新的线程，除非创建此线程超出maximumPoolSize，在这种情况下，任务将被拒绝。”，所以必然会新建一个线程来运行这个任务。
暂时还可以，但是如果这三个任务都还没完成，连续来了两个任务，第一个添加入queue中，后一个呢？queue中无法插入，而线程数达到了maximumPoolSize，所以只好执行异常策略了。
所以在使用SynchronousQueue通常要求maximumPoolSize是无界的，这样就可以避免上述情况发生（如果希望限制就直接使用有界队列）。对于使用SynchronousQueue的作用jdk中写的很清楚：此策略可以避免在处理可能具有内部依赖性的请求集时出现锁。

什么意思？如果你的任务A1，A2有内部关联，A1需要先运行，那么先提交A1，再提交A2，当使用SynchronousQueue我们可以保证，A1必定先被执行，在A1么有被执行前，A2不可能添加入queue中。

例子二：使用无界队列策略，即LinkedBlockingQueue

这个就拿newFixedThreadPool来说，根据前文提到的规则：

如果运行的线程少于 corePoolSize，则 Executor 始终首选添加新的线程，而不进行排队。那么当任务继续增加，会发生什么呢？

如果运行的线程等于或多于 corePoolSize，则 Executor 始终首选将请求加入队列，而不添加新的线程。OK，此时任务变加入队列之中了，那什么时候才会添加新线程呢？

如果无法将请求加入队列，则创建新的线程，除非创建此线程超出 maximumPoolSize，在这种情况下，任务将被拒绝。这里就很有意思了，可能会出现无法加入队列吗？不像SynchronousQueue那样有其自身的特点，对于无界队列来说，总是可以加入的（资源耗尽，当然另当别论）。换句说，永远也不会触发产生新的线程！corePoolSize大小的线程数会一直运行，忙完当前的，就从队列中拿任务开始运行。所以要防止任务疯长，比如任务运行的实行比较长，而添加任务的速度远远超过处理任务的时间，而且还不断增加，不一会儿就爆了。

例子三：有界队列，使用ArrayBlockingQueue。

这个是最为复杂的使用，所以JDK不推荐使用也有些道理。与上面的相比，最大的特点便是可以防止资源耗尽的情况发生。

举例来说，请看如下构造方法：

new ThreadPoolExecutor(   
                2, 4, 30, TimeUnit.SECONDS,    
              new ArrayBlockingQueue<Runnable>(2),    
                new RecorderThreadFactory("CookieRecorderPool"),    
               new ThreadPoolExecutor.CallerRunsPolicy());
假设，所有的任务都永远无法执行完。

对于首先来的A,B来说直接运行，接下来，如果来了C,D，他们会被放到queue中，如果接下来再来E,F，则增加线程运行E，F。但是如果再来任务，队列无法再接受了，线程数也到达最大的限制了，所以就会使用拒绝策略来处理。

keepAliveTime

jdk中的解释是：当线程数大于核心时，此为终止前多余的空闲线程等待新任务的最长时间。

有点拗口，其实这个不难理解，在使用了“池”的应用中，大多都有类似的参数需要配置。比如数据库连接池，DBCP中的maxIdle，minIdle参数。

什么意思？接着上面的解释，后来向老板派来的工人始终是“借来的”，俗话说“有借就有还”，但这里的问题就是什么时候还了，如果借来的工人刚完成一个任务就还回去，后来发现任务还有，那岂不是又要去借？这一来一往，老板肯定头也大死了。

合理的策略：既然借了，那就多借一会儿。直到“某一段”时间后，发现再也用不到这些工人时，便可以还回去了。这里的某一段时间便是keepAliveTime的含义，TimeUnit为keepAliveTime值的度量。

RejectedExecutionHandler

另一种情况便是，即使向老板借了工人，但是任务还是继续过来，还是忙不过来，这时整个队伍只好拒绝接受了。

RejectedExecutionHandler接口提供了对于拒绝任务的处理的自定方法的机会。在ThreadPoolExecutor中已经默认包含了4中策略，因为源码非常简单，这里直接贴出来。

CallerRunsPolicy：线程调用运行该任务的 execute 本身。此策略提供简单的反馈控制机制，能够减缓新任务的提交速度。

public void rejectedExecution(Runnable r, ThreadPoolExecutor e) {
           if (!e.isShutdown()) {
               r.run();
           }
       }
这个策略显然不想放弃执行任务。但是由于池中已经没有任何资源了，那么就直接使用调用该execute的线程本身来执行。

AbortPolicy：处理程序遭到拒绝将抛出运行时RejectedExecutionException

public void rejectedExecution(Runnable r, ThreadPoolExecutor e) {
           throw new RejectedExecutionException();
       }
 这种策略直接抛出异常，丢弃任务。

DiscardPolicy：不能执行的任务将被删除

public void rejectedExecution(Runnable r, ThreadPoolExecutor e) {
       }
 这种策略和AbortPolicy几乎一样，也是丢弃任务，只不过他不抛出异常。

DiscardOldestPolicy：如果执行程序尚未关闭，则位于工作队列头部的任务将被删除，然后重试执行程序（如果再次失败，则重复此过程）

public void rejectedExecution(Runnable r, ThreadPoolExecutor e) {
           if (!e.isShutdown()) {
               e.getQueue().poll();
               e.execute(r);
           }
       }
该策略就稍微复杂一些，在pool没有关闭的前提下首先丢掉缓存在队列中的最早的任务，然后重新尝试运行该任务。这个策略需要适当小心。

设想:如果其他线程都还在运行，那么新来任务踢掉旧任务，缓存在queue中，再来一个任务又会踢掉queue中最老任务。

总结：

keepAliveTime和maximumPoolSize及BlockingQueue的类型均有关系。如果BlockingQueue是无界的，那么永远不会触发maximumPoolSize，自然keepAliveTime也就没有了意义。

反之，如果核心数较小，有界BlockingQueue数值又较小，同时keepAliveTime又设的很小，如果任务频繁，那么系统就会频繁的申请回收线程。

public static ExecutorService newFixedThreadPool(int nThreads) {
       return new ThreadPoolExecutor(nThreads, nThreads,
                                     0L, TimeUnit.MILLISECONDS,
                                     new LinkedBlockingQueue<Runnable>());
   }