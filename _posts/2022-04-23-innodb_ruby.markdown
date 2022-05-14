---
title: innodb_ruby
layout: post
category: mysql
author: 夏泽民
---
% sudo gem install innodb_ruby
Password:

mysql> SHOW VARIABLES LIKE '%datadir%';
+---------------+--------------------------+
| Variable_name | Value                    |
+---------------+--------------------------+
| datadir       | /opt/homebrew/var/mysql/ |
+---------------+--------------------------+
1 row in set (0.01 sec)


 % innodb_space -s /opt/homebrew/var/mysql/ibdata1 system-spaces
/System/Library/Frameworks/Ruby.framework/Versions/2.6/usr/lib/ruby/2.6.0/universal-darwin21/rbconfig.rb:230: warning: Insecure world writable dir /usr/local/lib/node_modules in PATH, mode 040777
name                            pages       indexes
/Library/Ruby/Gems/2.6.0/gems/innodb_ruby-0.12.0/lib/innodb/index.rb:34:in `page': undefined method `record_describer=' for #<Innodb::Page::FspHdrXdes:0x0000000142944ab0> (NoMethodError)
	from /Library/Ruby/Gems/2.6.0/gems/innodb_ruby-0.12.0/lib/innodb/index.rb:19:in `initialize'
	from /Library/Ruby/Gems/2.6.0/gems/innodb_ruby-0.12.0/lib/innodb/space.rb:312:in `new'


 % brew search mysql
==> Formulae
automysqlbackup            mysql-client@5.7           mysql@5.6
mysql ✔                    mysql-connector-c++        mysql@5.7
mysql++                    mysql-sandbox              mysqltuner
mysql-client               mysql-search-replace       qt-mysql

==> Casks
mysql-connector-python     mysql-utilities            navicat-for-mysql
mysql-shell                mysqlworkbench             sqlpro-for-mysql


 % brew install  mysql@5.7
Running `brew update --preinstall`...
^@


https://github.com/jeremycole/innodb_ruby
https://github.com/jeremycole/innodb_ruby/wiki

echo 'export PATH="/opt/homebrew/opt/mysql@5.7/bin:$PATH"' >> ~/.zshrc
 brew services restart mysql@5.7
 
  % mysql -uroot
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 8
Server version: 8.0.28 Homebrew

 % brew unlink mysql && brew link --overwrite --force mysql@5.7
 
 
 % brew services stop mysql
 
  % brew services start mysql@5.7
==> Successfully started `mysql@5.7` (label: homebrew.mxcl.mysql@5.7)
 
  % mysql -uroot
ERROR 2002 (HY000): Can't connect to local MySQL server through socket '/tmp/mysql.sock' (2)

 %  /opt/homebrew/opt/mysql@5.7/bin/mysql.server start
Starting MySQL
. ERROR! The server quit without updating PID file (/opt/homebrew/var/mysql/xiazemindeMacBook-Pro.local.pid).

% rm -rf /opt/homebrew/opt/mysql@5.7
 % brew reinstall mysql@5.7
 加载8.0的文件失败
 
 % rm -r /opt/homebrew/var/mysql/
 brew services restart mysql@5.7
 
  % brew uninstall mysql
Uninstalling /opt/homebrew/Cellar/mysql/8.0.28_1... (304 files, 294.3MB)


2022-04-21T05:34:26.910874Z 0 [ERROR] unknown variable 'mysqlx-bind-address=127.0.0.1'
2022-04-21T05:34:26.910895Z 0 [ERROR] Aborting

Warning: The post-install step did not complete successfully
You can try again using:
  brew postinstall mysql@5.7
  
  rm -rf /usr/local/var/mysql
  
 https://blog.csdn.net/Wendy_He023/article/details/116508616
 vi /usr/local/etc/my.cnf
 %  brew postinstall mysql@5.7

% vi /opt/homebrew/etc/my.cnf
% brew services start mysql@5.7
 % mysql -uroot
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 2
Server version: 5.7.37 Homebrew
<!-- more -->
