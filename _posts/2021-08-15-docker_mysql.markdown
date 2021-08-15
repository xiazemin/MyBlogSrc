---
title: Docker容器中Mysql数据的导入/导出
layout: post
category: docker
author: 夏泽民
---

解决办法其实还是用mysqldump命令，但是我们需要进入docker的mysql容器内去执行它，并且通过配置volumes让导出的数据文件可以拷贝到宿主机的磁盘上

所以操作步骤就可以分为：

配置docker的volumes
进入docker的mysql容器，导出数据文件

mysql -h 127.0.0.1 -P 3306 -u root --password=root -c --default-character-set=utf8 容器数据库名< 主机目录下的脚本.sql

https://blog.csdn.net/qq_33326449/article/details/86478766
<!-- more -->
Docker-compose封装mysql并初始化数据以及redis

https://www.cnblogs.com/xiao987334176/p/12669080.html

https://zhuanlan.zhihu.com/p/26129750

导入导出
在容器中运行的 mysql 该怎么导入导出数据或结构呢？照这么做吧：

# Backup
docker exec CONTAINER /usr/bin/mysqldump -u root --password=root DATABASE > backup.sql

# Restore
docker exec -i CONTAINER /usr/bin/mysql -u root --password=root DATABASE < backup.sql

https://cloud.tencent.com/developer/article/1620057

让docker中的mysql启动时自动执行sql文件

```
编写容器启动脚本setup.sh：

#!/bin/bash
set -e

#查看mysql服务的状态，方便调试，这条语句可以删除
echo `service mysql status`

echo '1.启动mysql....'
#启动mysql
service mysql start
sleep 3
echo `service mysql status`

echo '2.开始导入数据....'
#导入数据
mysql < /mysql/schema.sql
echo '3.导入数据完毕....'

sleep 3
echo `service mysql status`

#重新设置mysql密码
echo '4.开始修改密码....'
mysql < /mysql/privileges.sql
echo '5.修改密码完毕....'

#sleep 3
echo `service mysql status`
echo 'mysql容器启动完毕,且数据导入成功'

tail -f /dev/null
```
https://www.imooc.com/article/19894
https://www.jianshu.com/p/e66a1c37bab0
https://www.coder.work/article/41872
https://www.codenong.com/25920029/
