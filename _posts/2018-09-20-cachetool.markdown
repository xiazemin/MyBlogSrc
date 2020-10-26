---
title: cachetool 非php环境操作apcu
layout: post
category: php
author: 夏泽民
---
https://packagist.org/packages/gordalina/cachetool
   opcache是个提升php性能的利器，但是在线上服务器真实遇到过偶尔几台服务器代码上线后，一直没有生效，查看opcache的配置也没有问题。后来没有办法，就在上线步骤中增加了重启php-fpm的操作。今天发现了一个小工具cachetool。可以方便的使用命令行清除opcache的缓存。
    当然除了重启php-fpm的进程可以清理opcache缓存外，opcache本身是支持清除缓存的。手动清理缓存涉及到的opcache函数主要为：opcache_reset()和opcache_invalidate() 。    但是opcache_reset()是php中的函数，需要在php脚本中执行，另外当PHP以PHP-FPM的方式运行的时候，opcache的缓存是无法通过php命令进行清除的，只能通过http或cgi到php-fpm进程的方式来清除缓存。
    如果是unix sock方式也可以
    php cachetool.phar opcache:status --fcgi=/var/run/php5-fpm.sock
<!-- more -->
    而opcache_invalidate 废除指定脚本缓存是使指定脚本的字节码缓存失效。可用于明确更新的代码文件列表时使用，但不方便清除整个脚本的缓存。cachetool使用起来也非常方便。如下：

#不需要安装下载下来即可使用
sudo mkdir /opt/modules/cachetool
cd /opt/modules/cachetool
#下载文件并增加可执行权限
sudo curl -sO http://gordalina.github.io/cachetool/downloads/cachetool.phar
sudo chmod +x cachetool.phar
#查看当前的opcache配置
[onlinedev@BFG-OSER-4471 ~]$ php /opt/modules/cachetool/cachetool.phar opcache:configuration --fcgi=127.0.0.1:9000
Zend OPcache 7.0.10
+---------------------------------------+----------------------+
| Directive                             | Value                |
+---------------------------------------+----------------------+
| opcache.enable                        | true                 |
| opcache.enable_cli                    | false                |
| opcache.use_cwd                       | true                 |
| opcache.validate_timestamps           | true                 |
| opcache.inherited_hack                | true                 |
| opcache.dups_fix                      | false                |
| opcache.revalidate_path               | false                |
| opcache.log_verbosity_level           | 1                    |
| opcache.memory_consumption            | 1073741824           |
| opcache.interned_strings_buffer       | 30                   |
| opcache.max_accelerated_files         | 4000                 |
| opcache.max_wasted_percentage         | 0.050000000000000003 |
| opcache.consistency_checks            | 0                    |
| opcache.force_restart_timeout         | 180                  |
| opcache.revalidate_freq               | 180                  |
| opcache.preferred_memory_model        | ''                   |
| opcache.blacklist_filename            | ''                   |
| opcache.max_file_size                 | 0                    |
| opcache.error_log                     | ''                   |
| opcache.protect_memory                | false                |
| opcache.save_comments                 | false                |
| opcache.fast_shutdown                 | true                 |
| opcache.enable_file_override          | false                |
| opcache.optimization_level            | 2147467263           |
| opcache.lockfile_path                 | '/tmp'               |
| opcache.file_cache                    | ''                   |
| opcache.file_cache_only               | false                |
| opcache.file_cache_consistency_checks | true                 |
+---------------------------------------+----------------------+
#列表查看opcache缓存的脚本
[onlinedev@BFG-OSER-4471 ~]$ php /opt/modules/cachetool/cachetool.phar opcache:status:scripts --fcgi=127.0.0.1:9000
+------+-----------+---------------------------------------------------------------------------+
| Hits | Memory    | Filename                                                                  |
+------+-----------+---------------------------------------------------------------------------+
| 8    | 4.16 KiB  | /opt/data/api/revs/r201710111538_1367/vendor/composer/autoload_static.php |
| 8    | 5.95 KiB  | /opt/data/api/revs/r201710111538_1367/www/index.php                       |
| 8    | 23.85 KiB | /opt/data/api/revs/r201710111538_1367/conf/Config.php                     |
| 8    | 848 b     | /opt/data/api/revs/r201710111538_1367/vendor/autoload.php                 |
| 8    | 84.88 KiB | /opt/data/api/revs/r201710111538_1367/lib/util/Params.php                 |
| 8    | 30.48 KiB | /opt/data/api/revs/r201710111538_1367/lib/Log.php                         |
| 8    | 6.34 KiB  | /opt/data/api/revs/r201710111538_1367/vendor/composer/autoload_real.php   |
| 8    | 51.69 KiB | /opt/data/api/revs/r201710111538_1367/app/controller/DebugController.php  |
| 8    | 6.48 KiB  | /opt/data/api/revs/r201710111538_1367/lib/Controller.php                  |
| 8    | 23.84 KiB | /opt/data/api/revs/r201710111538_1367/vendor/composer/ClassLoader.php     |
+------+-----------+---------------------------------------------------------------------------+
#查看当前的opcache缓存的统计信息
[onlinedev@BFG-OSER-4471 ~]$ php /opt/modules/cachetool/cachetool.phar opcache:status --fcgi=127.0.0.1:9000       
+----------------------+---------------------------------+
| Name                 | Value                           |
+----------------------+---------------------------------+
| Enabled              | Yes                             |
| Cache full           | No                              |
| Restart pending      | No                              |
| Restart in progress  | No                              |
| Memory used          | 66.66 MiB                       |
| Memory free          | 957.34 MiB                      |
| Memory wasted (%)    | 0 b (0%)                        |
| Strings buffer size  | 30 MiB                          |
| Strings memory used  | 221.27 KiB                      |
| Strings memory free  | 29.78 MiB                       |
| Number of strings    | 5501                            |
+----------------------+---------------------------------+
| Cached scripts       | 10                              |
| Cached keys          | 16                              |
| Max cached keys      | 7963                            |
| Start time           | Wed, 11 Oct 2017 03:30:16 +0000 |
| Last restart time    | Wed, 11 Oct 2017 07:38:20 +0000 |
| Oom restarts         | 0                               |
| Hash restarts        | 0                               |
| Manual restarts      | 4                               |
| Hits                 | 100                             |
| Misses               | 22                              |
| Blacklist misses (%) | 0 (0%)                          |
| Opcache hit rate     | 81.967213114754                 |
+----------------------+---------------------------------+
#执行清理opcache缓存
[onlinedev@BFG-OSER-4471 ~]$ php /opt/modules/cachetool/cachetool.phar opcache:reset --fcgi=127.0.0.1:9000   
#再次查看opcache缓存信息，会发现Cached scripts已被清空。                               
[onlinedev@BFG-OSER-4471 ~]$ php /opt/modules/cachetool/cachetool.phar opcache:status --fcgi=127.0.0.1:9000       
+----------------------+---------------------------------+
| Name                 | Value                           |
+----------------------+---------------------------------+
| Enabled              | Yes                             |
| Cache full           | No                              |
| Restart pending      | No                              |
| Restart in progress  | No                              |
| Memory used          | 66.43 MiB                       |
| Memory free          | 957.57 MiB                      |
| Memory wasted (%)    | 0 b (0%)                        |
| Strings buffer size  | 30 MiB                          |
| Strings memory used  | 171.92 KiB                      |
| Strings memory free  | 29.83 MiB                       |
| Number of strings    | 4283                            |
+----------------------+---------------------------------+
| Cached scripts       | 0                               |
| Cached keys          | 0                               |
| Max cached keys      | 7963                            |
| Start time           | Wed, 11 Oct 2017 03:30:16 +0000 |
| Last restart time    | Wed, 11 Oct 2017 08:36:17 +0000 |
| Oom restarts         | 0                               |
| Hash restarts        | 0                               |
| Manual restarts      | 5                               |
| Hits                 | 0                               |
| Misses               | 2                               |
| Blacklist misses (%) | 0 (0%)                          |
| Opcache hit rate     | 0                               |
+----------------------+---------------------------------+ cachetool除了可操作opcache缓存外，还可以操作apc缓存。所有的方法列表如下.
官方文档地址：http://gordalina.github.io/cachetool/
apc
  apc:bin:dump             Get a binary dump of files and user variables
  apc:bin:load             Load a binary dump into the APC file and user variables
  apc:cache:clear          Clears APC cache (user, system or all)
  apc:cache:info           Shows APC user & system cache information
  apc:cache:info:file      Shows APC file cache information
  apc:key:delete           Deletes an APC key
  apc:key:exists           Checks if an APC key exists
  apc:key:fetch            Shows the content of an APC key
  apc:key:store            Store an APC key with given value
  apc:sma:info             Show APC shared memory allocation information
opcache
  opcache:configuration    Get configuration information about the cache
  opcache:reset            Resets the contents of the opcode cache
  opcache:status           Show summary information about the opcode cache
  opcache:status:scripts   Show scripts in the opcode cache
