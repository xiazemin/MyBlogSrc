---
title: redis协议
layout: post
category: web
author: 夏泽民
---
<!-- more -->
redis命令发送格式：
*<参数数量> CRLF 
$<参数 1 的字节数量> CRLF 
<参数 1 的数据> CRLF 
... 
$<参数 N 的字节数量> CRLF 
<参数 N 的数据> CRLF 
其中CRLF表示 rn

举个例子：set name wuzhc

格式化输出：

*3 
$3 
set 
$4 
name 
$5 
wuzhc

说明：

*开头，表示有多少个参数，例如*3表示有3个参数（set, name, wuzhc）
$开头，表示参数的字节长度，例如$3表示set有3个字节，$4表示name有4个字节
每行rn结尾
通信协议为：

*3\r\n$3\r\nset\r\n$4\r\nname\r\n$5\r\nwuzhc\r\n
Redis 回复
状态回复（status reply）的第一个字节是 "+"，例如+OK\r\n
错误回复（error reply）的第一个字节是 "-"，例如-No such key\r\n
整数回复（integer reply）的第一个字节是 ":"，例如:1\r\n
批量回复（bulk reply）的第一个字节是 "$"，例如 $5\r\nwuzhc\r\n
多条批量回复（multi bulk reply）的第一个字节是 "*"，例如*2\r\n$5\r\nwuzhc\r\n$3r\nage\r\n
PHP 实现Redis客户端
{% highlight php linenos %}
<?php
/**
 * Created by PhpStorm.
 * User: wuzhc2016@163.com
 * Date: 2017年09月12日
 * Time: 9:08
 */
class Client
{
    private $_socket = null;
    public function __construct($ip, $port) 
    {
        $this->_socket = stream_socket_client(
            "tcp://{$ip}:{$port}",
            $errno,
            $errstr,
            1,
            STREAM_CLIENT_CONNECT
        );
        if (!$this->_socket) {
            exit($errstr);
        }
    }
    /**
     * 执行redis命令
     * @param $command
     * @return array|bool|string
     */
    public function exec($command)
    {      
        // 拼装发送命令格式
        $command = $this->_execCommand($command);

        // 发送命令到redis
        fwrite($this->_socket, $command);

        // 解析redis响应内容
        return $this->_parseResponse();
    }
    /**
     * 将字符改为redis通讯协议格式
     * 例如mget name age 格式化为 *3\r\n$4\r\nmget\r\n$4\r\nname\r\n$3\r\nage\r\n
     * @param $command
     * @return bool|string
     */
    private function _execCommand($command)
    {
        $line = '';
        $crlf = "\r\n";
        $params = explode(' ', $command);
        if (empty($params)) {
            return $line;
        }

        // 参数个数
        $line .= '*' . count($params) . $crlf;

        // 各个参数拼装
        foreach ((array)$params as $param) {
            $line .= '$' . mb_strlen($param, '8bit') . $crlf;
            $line .= $param . $crlf;
        }

        return $line;
    }
    /**
     * 解析redis回复
     * @return array|bool|string
     */
    private function _parseResponse()
    {
        $line = fgets($this->_socket); 
        $type = $line[0]; 
        $msg = mb_substr($line, 1, -2, '8bit'); 

        switch ($type) {
            // 状态回复
            case '+':
                if ($msg == 'OK' || $msg == 'PONG') {
                    return true;
                } else {
                    return $msg;
                }
            // 错误回复
            case '-':
                exit($msg);
            // 整数回复
            case ':':
                return $msg;
            // 批量回复
            case '$': // $后面跟数据字节数(长度)
                $line = fread($this->_socket, (int)$msg + 2); // 数据字节数 + (\r\n)两个字节
                return mb_substr($line, 0, -2, '8bit'); // 去除最后两个字节
            // 多条批量回复
            case '*': // *表示后面有多少个参数
                $data = [];
                for ($i = 0; $i < $msg; $i++) {
                    $data[] = $this->_parseResponse();
                }
                return $data;
        }
    }
}
// demo
$client = new Client('127.0.0.1', 6379);
$client->exec('set name wuzhc');
$res = $client->exec('get name');
var_dump($res);
{% endhighlight %}
SET
C:
SET HENRY  HENRYFAN
以上命令是设置HENRY 的值为HENRYFAN.在Redis的通讯协议上会以空格把命令拆分成三行;得到最终的命令如下：
*3\r\n
$3\r\n
SET\r\n
$5\r\n
HENRY\r\n
$8\r\n
HENRYFAN\r\n
S:
服务端操作成功
+OK\r\n
如果出现错误服务端会返回
-错误信息\r\n
GET
C:
GET HENRY
产生的通讯指令是：
*2\r\n
$3\r\n
GET\r\n
$5\r\n
HENRY\r\n
S:
 如果存在这个Key则返回
$8\r\n
HENRYFAN\r\n
不存在返回
$-1\r\n
HKEYS
C:
HKEYS HENRY
以上命令是获取对应HENRY有多少个field成员
*2\r\n
$5\r\n
HKEYS\r\n
$5\r\n
HENRY\r\n
S:
如果不存在任何字段信息
*0\r\n
如果存在QQ字段信息
*1\r\n
$2\r\n
QQ\r\n
HMGET
C:
HMGET HENRY QQ
以上命令是获取HENRY的QQ信息。
*3\r\n
$5\r\n
HMGET\r\n
$5\r\n
HENRY\r\n
$2\r\n
QQ\r\n
S:
如果不存在字段值

*1\r\n
$-1\r\n
存在字段值

*1\r\n
$8\r\n
28304340\r\n
