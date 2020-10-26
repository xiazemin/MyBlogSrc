---
title: php-spider
layout: post
category: php
author: 夏泽民
---
<!-- more -->
$composer install

$php example/example_simple.php

  ENQUEUED:  1
  SKIPPED:   0
  FAILED:    0
  PERSISTED:    1

DOWNLOADED RESOURCES:
 - DMOZ - The Directory of the Web

 参考文档：
 https://doc.phpspider.org/demo-start.html
 
 $git clone https://github.com/xiazemin/phpspider.git
 
 $php qiushibaike_task.php
 
 2017-09-12 10:32:58 [error] Multitasking needs Redis support, The redis extension was not found
 
 $php -r 'echo phpinfo();' |grep ini
 
 $vi /usr/local/etc/php/7.0/php.ini
 
 ;extension_dir = "/usr/local/Cellar/php70/7.0.8/lib/php/extensions/no-debug-non-zts-20151012/"
extension=/usr/local/Cellar/php70/7.0.8/lib/php/extensions/no-debug-non-zts-20151012/redis.so
;extension_dir = "/usr/local/Cellar/php70/7.0.15_8/lib/php/extensions/no-debug-non-zts-20151012/"
extension=/usr/local/Cellar/php70/7.0.15_8/lib/php/extensions/no-debug-non-zts-20151012/test.so
extension=/usr/local/Cellar/php70/7.0.15_8/lib/php/extensions/no-debug-non-zts-20151012/vld.so

php ini   extension_dir  会覆盖

Warning: PHP Startup: Unable to load dynamic library 'ext/php_mysqli.dll' - dlopen(ext/php_mysqli.dll, 9): image not found in Unknown on line 0


vi /usr/local/etc/php/7.0/php.ini
;extension=php_mysqli.dll

$php -m|grep mysql

 
 $which redis-server
/usr/local/bin/redis-server

$redis-server &

$php qiushibaike_task.php


----------------------------- PHPSPIDER -----------------------------
PHPSpider version:2.0.7          PHP version:7.0.15
start time:2017-09-12 12:03:38   run 0 days 0 hours 0 minutes
spider name: 糗事百科测试样例
server id: 1
task number: 3
load average: 2.13, 2.13, 2.09
document: https://doc.phpspider.org
------------------------------- TASKS -------------------------------
taskid    taskpid   mem       collect succ   collect fail   speed
1         21935     2MB       0              1              3.3/s
------------------------------- SERVER ------------------------------
server    tasknum   mem       collect succ   collect fail   speed
1         3         2MB       0              1              3.3/s
--------------------------- COLLECT STATUS --------------------------
find pages      queue         collected      fields         depth
1               0             1              0              0
 ---------------------------------------------------------------------
Press Ctrl-C to quit. Start success.

 