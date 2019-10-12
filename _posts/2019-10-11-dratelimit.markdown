---
title: 分布式限流
layout: post
category: web
author: 夏泽民
---
两个基础限流算法是漏斗算法和令牌算法
分布式限流
如果你的应用是单个进程，那么限流就很简单，请求的计数算法都可以在内存里完成。限流算法几乎没有损耗，都是纯内存的计算。但是互联网世界的应用都是多节点的分布式的，每个节点的请求处理能力还不一定一样。我们需要考虑的是这多个节点的整体请求处理能力。
单个进程的处理能力是1w QPS 并不意味着整体的请求处理能力是 N * 1w QPS，因为整体的处理能力还会有共享资源的能力限制。这个共享资源一般是指数据库，也可以是同一台机器的多个进程共享的 CPU 和 磁盘 等资源，还有网络带宽因素也会制约整体的请求量。
这时候请求的计数算法就需要集中在一个地方（限流中间件）来完成。应用程序在处理请求之前都需要向这种集中管理器申请流量（空气、令牌桶）。 每一个请求都需要一次网络 IO，从应用程序到限流中间件之间。

比如我们可以使用 Redis + Lua 来实现这个限流功能，但是 Lua 的性能要比 C 弱很多，通常这个限流算法能达到 1w 左右的 QPS就到顶了。还可以使用 Redis-Cell 模块，其内部使用 Rust 实现，它能达到 5w 左右的 QPS 也就到极限了。这时候它们都进入了满负荷状态，但是在生产环境中我们不会希望它们一直满负荷工作。
那如何完成 10w QPS 的限流呢？
一个简单的想法就是将限流的 key 分桶，然后使用 Redis 集群来扩容，让限流的申请指令经过客户端的 hash 分桶后打散的集群的多个几点，借此分散压力。
那如何完成 百万 QPS （1M）的限流呢？
如果还使用上面的方法那需要的带宽资源和Redis实例也是惊人的。我们可能需要几十个 Redis 节点，外加上百M（1M * 20 字节 * 8） 的带宽来完成这个工作。
这时我们必须转换思路，不再使用这种集中管控的方式来工作了。
我们将整体的 QPS 按照权重分散多每个子节点，每个字节点在内存中进行单机限流。如果每个节点都是对等的，那么每个子节点就可以分得 1/n 的 QPS。
它的优点在于分散了限流压力，将 IO 操作变成纯内存计算，这样就可以很轻松地应对超高的 QPS 限流。但是这也增加了系统的复杂度，需要有一个集中的配置中心来向每个子节点来分发 QPS 阈值，需要每个应用字节点向这个配置中心进行注册，需要有一个配置管理后台来对系统的 QPS 分配进行管理。
<!-- more -->
1. 接口定义
package ratelimit

import "time"

// 限流器接口
type Limiter interface {
    Acquire() error
    TryAcquire() bool
}

// 限流定义接口
type Limit interface {
    Name() string
    Key() string
    Period() time.Duration
    Count() int32
    LimitType() LimitType
}

// 支持 burst
type BurstLimit interface {
    Limit
    BurstCount() int32
}

// 分布式定义的 burst
type DistLimit interface {
    Limit
    ClusterNum() int32
}

type LimitType int32
const (
    CUSTOM LimitType = iota
    IP
)
Limiter 接口参考了 Google 的 guava 包里的 Limiter 实现。Acquire 接口是阻塞接口，其实还需要加上 context 来保证调用链安全，因为实际项目中并没有用到 Acquire 接口，所以没有实现完善；同理，超时时间的支持也可以通过添加新接口继承自 Limiter 接口来实现。TryAcquire 会立即返回。

Limit 抽象了一个限流定义，Key() 方法返回这个 Limit 的唯一标识，Name() 仅作辅助，Period() 表示周期，单位是秒，Count() 表示周期内的最大次数，LimitType(）表示根据什么来做区分，如 IP，默认是 CUSTOM.

BurstLimit 提供突发的能力，一般是配合令牌桶算法。DistLimit 新增 ClusterNum() 方法，因为 mentor 要求分布式遇到错误的时候，需要退化为单机版本，退化的策略即是：2 节点总共 100QPS，如果出现分区，每个节点需要调整为各 50QPS

https://github.com/GuoZhaoran/rateLimit
https://github.com/GuoZhaoran/Go-RateLimit

令牌桶的状态变量得放在一个 线程安全/一致 的地方，redis 是不二人选。但是令牌桶的算法核心是个延迟计算得到令牌数量，这个是一个很长的临界区，所以要么用分布式锁，要么直接利用 redis 的单线程以原子方式跑。一般业界是后者，即 lua 脚本维护令牌桶的状态变量、计算令牌。

php实现漏桶算法
Redis中设置接口限制1s内访问100次的hash:

 hmset org1/user/list expire 1 limitReq 100
我们使用Predis连接redis进行操作，模拟接口比较简单，我们只获取两个参数，org和pathInfo,RateLimit类中相关方法是：

<?php
/**
 * Description: 漏桶限流
 * User: guozhaoran<guozhaoran@cmcm.com>
 * Date: 2019-06-13
 */

class RateLimit
{
    private $conn = null;       //redis连接
    private $org = '';          //公司标识
    private $pathInfo = '';     //接口路径信息

    /**
     * RateLimit constructor.
     * @param $org
     * @param $pathInfo
     * @param $expire
     * @param $limitReq
     */
    public function __construct($org, $pathInfo)
    {
        $this->conn = $this->getRedisConn();
        $this->org = $org;
        $this->pathInfo = $pathInfo;
    }
    //......此处省略getLuaScript方法
    /**
     * 获取redis连接
     * @return \Predis\Client
     */
    private function getRedisConn()
    {
        require_once('vendor/autoload.php');
        $conn = new Predis\Client(['host' => '127.0.0.1',
            'port' => 6379,]);
        return $conn;
    }
    //......此处省略isActionAllowed方法
}
下边我们看看Lua脚本的设计：

   /**
     * 获取lua脚本
     * @return string
     */
    private function getLuaScript()
    {
        $luaScript = <<<LUA_SCRIPT
-- 限制接口访问频次
local times = redis.call('incr', KEYS[1]);    --将key自增1

if times == 1 then
redis.call('expire', KEYS[1], ARGV[1])    --给key设置过期时间
end

if times > tonumber(ARGV[2]) then
return 0
end

return 1
LUA_SCRIPT;

        return $luaScript;
    }
Lua脚本可以打包到Redis服务端进行执行，因为Redis服务端redis-server在2.6版本默认内置了Lua解析器，php的Redis客户端与Lua脚本交互主要传两个KEYS和ARGV,其中KEYS是对应Redis中操作的key值（示例中的KEYS[1]就是org1/user/list),ARGV是要设置的属性参数。在Lua脚本中Table的索引是从1开始自增的，Lua脚本执行Redis命令可以保证原子性(因为Redis是单线程的)，所以在并发竞态条件下也能保证hash的读写一致。命令首先调用incr设置org/user/list记数,Redis中的list、set、hash、zset这四种数据结构是容器型数据结构，他们共享下面两条通用规则：

1.create if not exists：如果容器不存在，那就创建一个再进行操作。比如incr org/user/list时，如果org/user/list不存在，就相当于设置了org/user/list为1,这就是为什么上边Lua脚本使用expire当times为1时设置org/user/list的过期时间
2.drop if no elements：如果容器里的元素没有了，那么立即删除容器，释放内存。比如lpop操作完一个list之后，list中没有元素内容了，那么这个list也就不存在了
下边的逻辑就很明了了，就是看接口的调用累加次数有没有超限（限制频率通过ARGV[2]）进行判断，超限返回0，否则返回1.

下边我们就可以看看怎样isActionAllowed方法判断是否要进行限流了:

    /**
     * 判断接口是否限制访问
     * @return bool
     */
    public function isActionAllowed()
    {
        $pathInfo = $this->org . $this->pathInfo;
        $config = $this->conn->hgetall($pathInfo);
        //配置中没有对接口进行限制
        if (!$config) return true;

        $pathInfoLimitKey = $this->org . '-' . $this->pathInfo;
        try {
            $ret = $this->conn->evalsha(sha1($this->getLuaScript()), 1, $pathInfoLimitKey, $config['expire'], $config['limitReq']);
        } catch (Exception $e) {
            $ret = $this->conn->eval($this->getLuaScript(), 1, $pathInfoLimitKey, $config['expire'], $config['limitReq']);
        }

        return boolval($ret);
    }
Predis使用evalsha打包Lua脚本发送到服务端执行。evalsha的第一个参数是sha1编码后的Lua脚本。redis-server可以对Lua脚本进行缓存，缓存的方法是key:value的形式，其中key是sha1后的lua脚本内容,这样在Lua脚本比较大时，客户端只需要发送sha1后的值到redis-server就可以了，减小了每次发送命令内容的字节大小。如果evalsha报出错误信息可以改为eval函数，因为redis-server第一次接收到lua脚本，可能还没没有进行缓存。最好是使用try...catch...做一下兼容处理。evalsha的第二个参数是key的个数，这里是一个，$pathInfoLimitKey，下边两个是从Redis中取出的配置值，标示1s内允许$pathInfoLimitKey被操作100次。如果没有对$pathInfoLimitKey做配置限制频率，默认不受限。

以上就是rateLimit类的全部内容了，思路比较简单，下边简单看一下入口文件，也比较简单，就是接收参数，然后将接口是否受限的信息写到stat.log日志文件中去。

<?php
/**
 * Description: 漏斗限流入口文件
 * User: guozhaoran<guozhaoran@cmcm.com>
 * Date: 2019-06-16
 */
require_once('./RateLimit.php');
ini_set('display_errors', true);

$org = $_GET['org'];
$pathInfo = $_GET['path_info'];

$result = (new RateLimit($org, $pathInfo))->isActionAllowed();

$handler = fopen('./stat.log', 'a') or die('can not open file!');
if ($result) {
    fwrite($handler, 'request success!' . PHP_EOL);
} else {
    fwrite($handler, 'request failed!' . PHP_EOL);
}
fclose($handler);
我们通过ab工具压测一下接口信息，程序限制1s内允许100次访问，我们就开10个客户端同时请求110次，理论上应该是前一百次是成功的，后十次是失败的，命令为:

ab -n 110 -c 10 http://localhost/demo/rateLimit/index.php\?org\=org1\&path_info\=/user/list
stat.log中的日志信息和我们预期中的一样，说明我们的接口频次设置达到了预期效果：

...//此处省略96行
request success!
request success!
request success!
request success!
request failed!
request failed!
request failed!
request failed!
request failed!
request failed!
request failed!
request failed!
request failed!
request failed!
但是漏斗限流还是有一些缺点的，它不支持突发流量，我们接口设置1s内限制访问100次，假如说前900毫秒只有80次访问，突然在接下来的100毫秒来了50次访问，那么毫无疑问，后边30次访问是失败的。不过漏斗这种简单粗暴的限流处理方案对于流量集中性访问，比如(1分钟只允许访问1000次)还是非常适合的。

3.2 go语言实现令牌桶算法
我们首先不考虑竞态条件，用go语言实现一个v1版本的令牌桶来体会一下它的算法思想。我们新建一个funnel模块，定义一个结构体，包含了令牌桶需要的属性：

package funnel

import (
    "math"
    "time"
)

type Funnel struct {
    Capacity          int64   //令牌桶容量
    LeakingRate       float64 //令牌桶流水速率:每毫秒向令牌桶中添加的令牌数
    RemainingCapacity int64   //令牌桶剩余空间
    LastLeakingTime   int64   //上次流水(放入令牌)时间:毫秒时间戳
Funnel结构体支持导出，分别包含令牌桶的容量、向令牌桶中添加令牌的速率、令牌桶剩余空间
和上次放入令牌时间的四个属性。
我们采用请求进来时实时改变令牌桶状态的思路，改变令牌桶状态的方法如下：

//有请求时更新令牌桶的状态,主要是令牌桶剩余空间和记录取走Token的时间戳
func (rateLimit *Funnel) updateFunnelStatus() {
    nowTs := time.Now().UnixNano() / int64(time.Millisecond)
    //距离上一次取走令牌已经过去了多长时间
    timeDiff := nowTs - rateLimit.LastLeakingTime
    //根据时间差和流水速率计算需要向令牌桶中添加多少令牌
    needAddSpace := int64(math.Floor(rateLimit.LeakingRate * float64(timeDiff)))
    //不需要添加令牌
    if needAddSpace < 1 {
        return
    }
    rateLimit.RemainingCapacity += needAddSpace
    //添加的令牌不能大于令牌桶的剩余空间
    if rateLimit.RemainingCapacity > rateLimit.Capacity {
        rateLimit.RemainingCapacity = rateLimit.Capacity
    }
    //更新上次令牌桶流水(添加令牌)时间戳
    rateLimit.LastLeakingTime = nowTs
}
因为要改变令牌桶的状态，所以我们这里使用指针接收者为结构体Funnel定义方法。主要思路就是根据当前时间和上次放入令牌桶中令牌的时间戳，再结合每毫秒应该放入令牌桶中令牌，计算添加应该放入到令牌桶中的令牌，放入令牌后不能超过令牌桶本身容量的大小。然后取出令牌，更新上次添加令牌时间戳。
判断接口是否限流其实就是看能不能从令牌桶中取出令牌，方法如下：

//判断接口是否被限流
func (rateLimit *Funnel) IsActionAllowed() bool {
    //更新令牌桶状态
    rateLimit.updateFunnelStatus()
    if rateLimit.RemainingCapacity < 1 {
        return false
    }
    rateLimit.RemainingCapacity = rateLimit.RemainingCapacity - 1
    return true
}
到了这里，相信读者已经对令牌桶算法有了一个比较清晰的认识了。我们再来说问题，因为限流最终还是要通过操作Redis来实现的，我们首先来在Redis里初始化好接口限流的配置：

hmset org2/user/list Capacity 100 LeakingRate 0.1 RemainingCapacity 0 LastLeakingTime 1560789716896
我们设置公司二(org2)的接口（/user/list）令牌桶容量100，每隔10ms放入一令牌（计算方法100/1000）。我们将Funnel对象内容的字段存储到一个hash结构中，我们在计算是否限流的时候需要从hash结构中取值，在内存中做运算，再回填到hash结构，尤其对于go语言这种天然并发的程序来讲，我们无法保证整个过程的原子化（这就是为什么要使用Lua脚本的原因，因为如果用程序来实现，就需要加锁，一旦加锁就有加锁失败的可能，失败只能选择重试或放弃，重试会导致性能下降，放弃会影响用户体验，代码复杂度会增加不少）。我们V2版本还是会选择使用Lua脚本来实现：具体调研过程如下：

方案	特点
单服务对操作采用锁机制	文章有提到，这种只能保证单节点下串行且性能差
Redis原子操作incr	这种方案我们在漏斗模型中有使用，它只能应对简单的场景，涉及到复杂场景就比较难处理
Redis分布式事务	虽然Redis的分布式事务能保证原子操作，但是实现复杂，并且网络开销大，需要大量的网络传输
Redis+Lua	这里就不得不夸一下这种方案了，Lua脚本中运行在Redis中，redis又是单线程的，因此能保证操作的串行。另外：减少网络开销，前边我们提到过，Lua代码包装的命令不需要发送多次命令请求，Redis可以对Lua脚本进行缓存，减少了网络传输,另外其他的客户端也可以使用缓存
补充一点：Redis4.0提供了一个限流模块Redis模块，它叫Redis-Cell。该模块也使用了漏斗算法，并提供了原子的限流命令，重试机制也非常简单，有兴趣的可以研究一下。我们这里还是使用Lua + Redis解决方案，废话少说，上V2版本的代码：

const luaScript = `
-- 接口限流
-- last_leaking_time 最后访问时间的毫秒
-- remaining_capacity 当前令牌桶中可用请求令牌数
-- capacity 令牌桶容量
-- leaking_rate    令牌桶添加令牌的速率

-- 把发生数据变更的命令以事务的方式做持久化和主从复制(Redis4.0支持)
redis.replicate_commands()

-- 获取令牌桶的配置信息
local rate_limit_info = redis.call("HGETALL", KEYS[1])

-- 获取当前时间戳
local timestamp = redis.call("TIME")
local now = math.floor((timestamp[1] * 1000000 + timestamp[2]) / 1000)

if rate_limit_info == nil then -- 没有设置限流配置,则默认拿到令牌
    return now * 10 + 1
end

local capacity = tonumber(rate_limit_info[2])
local leaking_rate = tonumber(rate_limit_info[4])
local remaining_capacity = tonumber(rate_limit_info[6])
local last_leaking_time = tonumber(rate_limit_info[8])

-- 计算需要补给的令牌数,更新令牌数和补给时间戳
local supply_token = math.floor((now - last_leaking_time) * leaking_rate)
if (supply_token > 0) then
   last_leaking_time = now
   remaining_capacity = supply_token + remaining_capacity
   if remaining_capacity > capacity then
      remaining_capacity = capacity
   end
end

local result = 0 -- 返回结果是否能够拿到令牌,默认否

-- 计算请求是否能够拿到令牌
if (remaining_capacity > 0) then
    remaining_capacity = remaining_capacity - 1
    result = 1
end

-- 更新令牌桶的配置信息
redis.call("HMSET", KEYS[1], "RemainingCapacity", remaining_capacity, "LastLeakingTime", last_leaking_time)

return now * 10 + result
`
我们这段脚本返回一个int64类型的整数，最后一位0或1表示是否要对接口限流，前边的数字表示毫秒时间戳，将来记录到日志里进行压测统计使用。程序运行时当前时间戳我是调用Redis的time命令计算获得的，原因有二：

Lua命令获得当前时间戳只能精确到秒，而Redis确可以精确到纳秒。
如果时间戳作为脚本调用参数（go程序）传进来会有问题，因为脚本传参到Lua在Redis中执行还有一段时间误差，不能保证最先被接收到的请求先被处理，而Lua中获取时间戳可以保证请求、时间串行
和以前一样，没有设置限流配置，就默认可以请求。
然后根据时间戳补给令牌，计算是否能够取到令牌，然后更新令牌状态，思路和V1版本一样，读者可自行阅读。说明一点，脚本开始处的redis.replicate_commands()命令是因为Redis低版本不支持对Redis既读又写，所以这种方式还是存在版本兼容性，但是解决办法确是最完美的。
接下来我们看go逻辑代码：

func main() {
    http.HandleFunc("/user/list", handleReq)
    http.ListenAndServe(":8082", nil)
}

//初始化redis连接池
func newPool() *redis.Pool {
    return &redis.Pool{
        MaxIdle:   80,
        MaxActive: 12000, // max number of connections
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", ":6379")
            if err != nil {
                panic(err.Error())
            }
            return c, err
        },
    }
}

//写入日志
func writeLog(msg string, logPath string) {
    fd, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    defer fd.Close()
    content := strings.Join([]string{msg, "\r\n"}, "")
    buf := []byte(content)
    fd.Write(buf)
}

//处理请求函数,根据请求将响应结果信息写入日志
func handleReq(w http.ResponseWriter, r *http.Request) {
    //获取url信息
    pathInfo := r.URL.Path
    //获取get传递的公司信息org
    orgInfo, ok := r.URL.Query()["org"]
    if !ok || len(orgInfo) < 1 {
        fmt.Println("Param org is missing!")
    }

    //调用lua脚本原子性进行接口限流统计
    conn := newPool().Get()
    key := orgInfo[0] + pathInfo
    lua := redis.NewScript(1, luaScript)
    reply, err := redis.Int64(lua.Do(conn, key))
    if err != nil {
        fmt.Println(err)
        return
    }
    //接口是否被限制访问
    isLimit := bool(reply % 10 == 1)
    reqTime := int64(math.Floor(float64(reply) / 10))
    //将统计结果写入日志当中
    if !isLimit {
        successLog := strconv.FormatInt(reqTime, 10) + " request failed!"
        writeLog(successLog, "./stat.log")
        return
    }

    failedLog := strconv.FormatInt(reqTime, 10) + " request success!"
    writeLog(failedLog, "./stat.log")
}
脚本监听本地8082端口，使用go的redis框架redigo来操作redis，我们初始化了一个redis连接池，从连接池中取得连接进行操作。我们分析如下代码：

lua := redis.NewScript(1, luaScript)
    reply, err := redis.Int64(lua.Do(conn, key))
NewScript中第一个参数代表要操作Redis的key的个数，这点和Predis的evalsha第二个参数类似。然后采用Do方法执行脚本，返回值使用redis.Int64做处理，然后进行运算判断接口是否允许被访问，然后将访问时间和结果写入到stat.log日志文件中。
逻辑还是非常的简单，我们主要看压测结果，启动代码，使用ab压测命令执行：

 ab -n 110 -c 10 http://127.0.0.1:8082/user/list\?org\=org2
然后我们分析stat.log日志兴许会有些惊讶：

1561263349294 request success!    //第一行日志
...//省略95行
1561263349387 request success!
1561263349388 request success!
1561263349398 request success!
1561263349396 request success!
1561263349404 request success!
1561263349407 request success!
1561263349406 request success!
1561263349406 request success!
1561263349407 request success!
1561263349406 request success!
1561263349406 request success!
1561263349405 request success!
1561263349406 request success!
1561263349406 request success!
1561263349406 request success!
是的，都成功了，为什么呢？我们看统计时间会发现执行这100个请求总共用了110毫秒，在程序执行过程中，每隔10ms会向令牌桶中添加一个令牌，一共添加了11个令牌，所以110次请求都拿到了令牌。可以看出令牌桶适用于大流量下的限流，可以保证流量按照时间均匀分摊，避免出现流量的集中式爆发访问。
