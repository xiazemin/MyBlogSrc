---
title: HTTP_X_FORWARDED_FOR
layout: post
category: web
author: 夏泽民
---
一、如果没有使用代理服务器，

　　REMOTE_ADDR = 客户端IP    HTTP_X_FORWARDED_FOR = 没数值或不显示

　　$ip = $_SERVER['REMOTE_ADDR'];
 

二、使用透明代理，

　　REMOTE_ADDR = 最后一个代理服务器 IP
　　HTTP_X_FORWARDED_FOR = 客户端真实 IP （经过多个代理服务器时，这个值类似：221.5.252.160， 203.98.182.163， 203.129.72.215）

　　这类代理还会将客户真实ip发送到请求对象，无法隐藏真实ip。

　　$ip = $_SERVER['HTTP_X_FORWARDED_FOR'];
<!-- more -->
三、使用普通匿名代理服务器，

　　REMOTE_ADDR = 最后一个代理服务器 IP
　　HTTP_X_FORWARDED_FOR = 代理服务器 IP （经过多个代理服务器时，这个值类似：203.98.182.163， 203.98.182.163， 203.129.72.215）

　　这样就隐藏了客户端的真实ip，但服务器会知道客户端是通过代理服务器去访问的。

四、使用欺骗性代理服务器，

　　REMOTE_ADDR = 代理服务器 IP
　　HTTP_X_FORWARDED_FOR = 随机的 IP（经过多个代理服务器时，这个值类似：220.4.251.159， 203.98.182.163， 203.129.72.215）

　　服务器可以识别到时通过代理服务器访问的，但发送给目标服务器的是虚假ip。

五、使用高匿名代理，

　　REMOTE_ADDR = 代理服务器 IP       HTTP_X_FORWARDED_FOR = 没数值或不显示

　　使用这种代理时，不同浏览器不同设备会返回不同的ip头信息，因此PHP使用$_SERVER["REMOTE_ADDR"] 、$_SERVER["HTTP_X_FORWARDED_FOR"] 获取的值可能是空值也可能是“unknown”值。

PHP获取ip代码如下：
public function ip() {
    //strcasecmp 比较两个字符，不区分大小写。返回0，>0，<0。
    if(getenv('HTTP_CLIENT_IP') && strcasecmp(getenv('HTTP_CLIENT_IP'), 'unknown')) {
        $ip = getenv('HTTP_CLIENT_IP');
    } elseif(getenv('HTTP_X_FORWARDED_FOR') && strcasecmp(getenv('HTTP_X_FORWARDED_FOR'), 'unknown')) {
        $ip = getenv('HTTP_X_FORWARDED_FOR');
    } elseif(getenv('REMOTE_ADDR') && strcasecmp(getenv('REMOTE_ADDR'), 'unknown')) {
        $ip = getenv('REMOTE_ADDR');
    } elseif(isset($_SERVER['REMOTE_ADDR']) && $_SERVER['REMOTE_ADDR'] && strcasecmp($_SERVER['REMOTE_ADDR'], 'unknown')) {
        $ip = $_SERVER['REMOTE_ADDR'];
    }
    $res =  preg_match ( '/[\d\.]{7,15}/', $ip, $matches ) ? $matches [0] : '';
    echo $res;
    //dump(phpinfo());//所有PHP配置信息
}
