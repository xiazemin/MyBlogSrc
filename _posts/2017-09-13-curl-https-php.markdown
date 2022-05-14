---
title: curl-https-php
layout: post
category: php
author: 夏泽民
---
php使用curl访问https返回无结果的问题
<!-- more -->
用curl发起https请求的时候报错：“SSL certificate problem, verify that the CA cert is OK. Details: error:14090086:SSL routines:SSL3_GET_SERVER_CERTIFICATE:certificate verify failed”
很明显，验证证书的时候出现了问题。

使用curl如果想发起的https请求正常的话有2种做法：

方法一、设定为不验证证书和host。

在执行curl_exec()之前。设置option

$ch = curl_init();

......

curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, FALSE);
curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, FALSE);

方法二、设定一个正确的证书。

本地ssl判别证书太旧，导致链接报错ssl证书不正确。

我们需要下载新的ssl 本地判别文件

http://curl.haxx.se/ca/cacert.pem

放到 程序文件目录

curl 增加下面的配置

curl_setopt($ch,CURLOPT_SSL_VERIFYPEER,true); ;
curl_setopt($ch,CURLOPT_CAINFO,dirname(__FILE__).'/cacert.pem');

$curl http://curl.haxx.se/ca/cacert.pem -o cacert.pem

ErrnoSSL: certificate verification failed

英文mac本地安装了多版本的php，php curl 扩展版本不匹配
旧版本的cur 默认不支持https 需要重新编译php curl 扩展

 certificate verification failed (result: 5)
 vi cacert.pem 发现是空
 curl-config --ca
 也为空
 下载证书
 curl_setopt($ch,CURLOPT_CAINFO,dirname(__FILE__).'/cacert.pem');
 或者
 /.ssh
sudo wget http://curl.haxx.se/ca/cacert.pem
export CURL_CA_BUNDLE=~/.ssh/cacert.pem
发现
$url='https://yq.aliyun.com/ziliao/157066';
成功，但是
$url="https://www.baidu.com";
失败