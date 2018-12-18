---
title: docker mac 镜像清理
layout: post
category: docker
author: 夏泽民
---
du -sh *
这个命令用来查看根目录下，所有文件的大小分布
604K	Applications
165M	CRClientTools
1.2M	Desktop
802M	Documents
2.0G	Downloads
216M	GitBook
du: Library/Python: Permission denied
 86G	Library
 
du -sh Library
du: Library/Python: Permission denied
 86G	Library
 
cd ~/Library
du -d 1 -h
 16K	./com.lc-tech.licman
  0B	./Compositions
 66G	./Containers
 
$du -d 1 -h
280K	./com.apple.WeatherKitService
 64G	./com.docker.docker

$cd ./com.docker.docker
$du -d 1 -h
 64G	./Data
 64G	.
 
 cd ./Data/
 
 $du -d 1 -h
 64G	./com.docker.driver.amd64-linux
 $cd ./com.docker.driver.amd64-linux
 
 $ls -al
 -rw-r--r--@  1 didi  staff  68667637760 12 17 17:50 Docker.qcow2
 
 docker for mac 有个bug，删除了容器或者镜像后，docker 占用的电脑硬盘空间不会相应的减少（Docker.qcow2文件）。
https://blog.mrtrustor.net/post/clean-docker-for-mac/
https://gist.github.com/MrTrustor/e690ba75cefe844086f5e7da909b35ce#file-clean-docker-for-mac-sh
这位法国老哥写了个脚本自动把docker镜像保存到本地， 删除Docker.qcow2文件，重启docker ，再把保存下来的镜像 load 到docker 中。
<!-- more -->
#!/bin/bash

# Copyright 2017 Théo Chamley
# Permission is hereby granted, free of charge, to any person obtaining a copy of 
# this software and associated documentation files (the "Software"), to deal in the Software
# without restriction, including without limitation the rights to use, copy, modify, merge,
# publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons
# to whom the Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all copies or
# substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
# BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
# NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
# DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,

IMAGES=$@

echo "This will remove all your current containers and images except for:"
echo ${IMAGES}
read -p "Are you sure? [yes/NO] " -n 1 -r
echo    # (optional) move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    exit 1
fi


TMP_DIR=$(mktemp -d)

pushd $TMP_DIR >/dev/null

open -a Docker
echo "=> Saving the specified images"
for image in ${IMAGES}; do
	echo "==> Saving ${image}"
	tar=$(echo -n ${image} | base64)
	docker save -o ${tar}.tar ${image}
	echo "==> Done."
done

echo "=> Cleaning up"
echo -n "==> Quiting Docker"
osascript -e 'quit app "Docker"'
while docker info >/dev/null 2>&1; do
	echo -n "."
	sleep 1
done;
echo ""

echo "==> Removing Docker.qcow2 file"
rm ~/Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux/Docker.qcow2

echo "==> Launching Docker"
open -a Docker
echo -n "==> Waiting for Docker to start"
until docker info >/dev/null 2>&1; do
	echo -n "."
	sleep 1
done;
echo ""

echo "=> Done."

echo "=> Loading saved images"
for image in ${IMAGES}; do
	echo "==> Loading ${image}"
	tar=$(echo -n ${image} | base64)
	docker load -q -i ${tar}.tar || exit 1
	echo "==> Done."
done

popd >/dev/null
rm -r ${TMP_DIR}

$docker images
REPOSITORY                                             TAG                 IMAGE ID            CREATED             SIZE
hub.c.163.com/mrjucn/centos6.5-mysql5.1-php5.7-nginx   latest              726cb1dfd4b7        2 years ago         2.78 GB
hub.c.163.com/public/redis                             2.8.4               4888527e1254        2 years ago         190 MB
hub.c.163.com/longjuxu/microbox/etcd                   latest              6aef84b9ec5a        3 years ago         17.9 MB

$sh file-clean-docker-for-mac.sh 726cb1dfd4b7 4888527e1254 6aef84b9ec5a
This will remove all your current containers and images except for:
726cb1dfd4b7 4888527e1254 6aef84b9ec5a
Are you sure? [yes/NO] y
=> Saving the specified images
==> Saving 726cb1dfd4b7

