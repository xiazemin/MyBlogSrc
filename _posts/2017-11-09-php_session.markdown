---
title: php_session
layout: post
category: php
author: 夏泽民
---
<!-- more -->
void session_write_close ( void )

End the current session and store session data.

Session data is usually stored after your script terminated without the need to call session_write_close(), but as session data is locked to prevent concurrent writes only one script may operate on a session at any time. When using framesets together with sessions you will experience the frames loading one by one due to this locking. You can reduce the time needed to load all the frames by ending the session as soon as all changes to session variables are done.
也就是说session是有锁的，为防止并发的写会话数据,php自带的的文件保存会话数据是加了一个互斥锁（在session_start()的时候）。 
程序执行session_start()，此时当前程序就开始持有锁。 
程序结束，此时程序自动释放Session的锁。

如果同一个客户端同时并发发送多个请求（如ajax在页面同时发送多个请求），且脚本执行时间较长，就会导致session文件阻塞，影响性能。因为对于每个请求，PHP执行session_start()，就会取得文件独占锁，只有在该请求处理结束后，才会释放独占锁。这样，同时多个请求就会引起阻塞。解决方案如下： 
修改会话变量后，立即使用session_write_close()来保存会话数据并释放文件锁。
session_start();   
$_SESSION['test'] = 'test';
session_write_close();


(PHP 5 >= 5.3.3, PHP 7)
fastcgi_finish_request — 冲刷(flush)所有响应的数据给客户端
如果有锁的话会使异步作用失效


There are some pitfalls  you should be aware of when using this function.

The script will still occupy a FPM process after fastcgi_finish_request(). So using it excessively for long running tasks may occupy all your FPM threads up to pm.max_children. This will lead to gateway errors on the webserver.

Another important thing is session handling. Sessions are locked as long as they're active (see the documentation for session_write_close()). This means subsequent requests will block until the session is closed.

You should therefore call session_write_close() as soon as possible (even before fastcgi_finish_request()) to allow subsequent requests and a good user experience.

This also applies for all other locking techniques as flock or database locks for example. As long as a lock is active subsequent requests might bock.

