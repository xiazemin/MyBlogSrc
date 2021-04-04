---
title: librdkafka mac m1 安装
layout: post
category: golang
author: 夏泽民
---
pkg-config --cflags -- rdkafka

Package libcrypto was not found in the pkg-config search path. Perhaps you should add the directory containing `libcrypto.pc' to the PKG_CONFIG_PATH environment variable Package 'libcrypto', required by 'rdkafka', not found pkg-config: exit status 1

brew --prefix openssl /usr/local/opt/openssl@1.1
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/lib/pkgconfig

ls /usr/local/opt/openssl@1.1/lib/pkgconfig/ 
libcrypto.pc libssl.pc openssl.pc

 export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/opt/openssl@1.1/lib/pkgconfig/
 
 问题解决
 
 https://stackoverflow.com/questions/57967504/no-package-libcrypto-found-in-mac
 https://github.com/scipr-lab/libsnark/issues/99
 https://github.com/rfjakob/gocryptfs/issues/98
 https://stackoverflow.com/questions/52956290/package-rdkafka-was-not-found-in-the-pkg-config-search-path
<!-- more -->
"_rd_kafka_unsubscribe", referenced from:
  __cgo_13886585fdfe_Cfunc_rd_kafka_unsubscribe in _x007.o
 (maybe you meant: __cgo_13886585fdfe_Cfunc_rd_kafka_unsubscribe)
"_rd_kafka_version", referenced from: __cgo_13886585fdfe_Cfunc_rd_kafka_version in _x014.o (maybe you meant: __cgo_13886585fdfe_Cfunc_rd_kafka_version, __cgo_13886585fdfe_Cfunc_rd_kafka_version_str ) "_rd_kafka_version_str", referenced from: __cgo_13886585fdfe_Cfunc_rd_kafka_version_str in _x009.o (maybe you meant: __cgo_13886585fdfe_Cfunc_rd_kafka_version_str) ld: symbol(s) not found for architecture arm64 clang: error: linker command failed with exit code 1 (use -v to see invocation)

mac m1 是基于arm架构的，原来的lib包无法直接使用，需要源码重新安装
https://github.com/edenhill/librdkafka
https://github.com/edenhill/librdkafka/issues?q=arm64+

./configure --install-deps --source-deps-only

brew reinstall zstd 

Error: Cannot install in Homebrew on ARM processor in Intel default prefix (/usr/local)! Please create a new installation in /opt/homebrew using one of the "Alternative Installs" from: https://docs.brew.sh/Installation You can migrate your previously installed formula list with:

brew bundle dump

 Error: undefined method `bottle_hash' for #Formulary::FormulaNamespace592958f13892655fbf773c98b7dc73a3::PkgConfig:0x000000013d97d040 Please report this bug: https://github.com/Homebrew/homebrew-bundle/issues

 brew uninstall --ignore-dependencies zstd Uninstalling /usr/local/Cellar/zstd/1.4.8... (26 files, 3.4MB)

rosta brew install zstd

./configure --install-deps make

因为以前配置了Rosetta 需要去掉
重新安装brew https://docs.brew.sh/Installation
https://stackoverflow.com/questions/54926712/is-there-a-way-to-list-keys-in-context-context
https://liqiang.io/post/print-all-key-value-in-golang-context-2ac7c19f

homebrew % ./bin/brew install openssl

/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/uninstall)"

https://docs.brew.sh/Installation 

cd /opt mkdir homebrew && curl -L https://github.com/Homebrew/brew/tarball/master | tar xz --strip 1 -C homebrew

/opt % chmod -R 777 homebrew

curl -L https://github.com/Homebrew/brew/tarball/master | tar xz --strip 1 -C homebrew

% brew install librdkafka

https://blog.csdn.net/weixin_30253461/article/details/112518937 

关闭 在访达 -> 应用程序，找到 iTerm2，右键，选择“显示简介”，然后选择“使用 Rosetta 打开”

./configure --install-deps --source-deps-only
  
  
  
 libzstd ()
    module: self
    action: fail
    reason:
Failed to install dependency libzstd

###########################################################
### Installing the following packages might help:       ###
###########################################################
brew install  openssl zstd

 % brew list openssl@1.1
/opt/homebrew/Cellar/openssl@1.1/1.1.1k/bin/c_rehash
/opt/homebrew/Cellar/openssl@1.1/1.1.1k/bin/openssl

 % brew install  openssl@1.1 zstd
Warning: openssl@1.1 1.1.1k is already installed, it's just not linked.
To link this version, run:
  brew link openssl@1.1
  
   % brew link openssl@1.1
Warning: Refusing to link macOS provided/shadowed software: openssl@1.1
If you need to have openssl@1.1 first in your PATH, run:
  echo 'export PATH="/opt/homebrew/opt/openssl@1.1/bin:$PATH"' >> ~/.zshrc

For compilers to find openssl@1.1 you may need to set:
  export LDFLAGS="-L/opt/homebrew/opt/openssl@1.1/lib"
  export CPPFLAGS="-I/opt/homebrew/opt/openssl@1.1/include"

For pkg-config to find openssl@1.1 you may need to set:
  export PKG_CONFIG_PATH="/opt/homebrew/opt/openssl@1.1/lib/pkgconfig"
  
  
xiazemin@xiazemindeMacBook-Pro librdkafka % echo 'export PATH="/opt/homebrew/opt/openssl@1.1/bin:$PATH"' >> ~/.zshrc
xiazemin@xiazemindeMacBook-Pro librdkafka % echo ' export PKG_CONFIG_PATH="/opt/homebrew/opt/openssl@1.1/lib/pkgconfig"' >>  ~/.zshrc
xiazemin@xiazemindeMacBook-Pro librdkafka % echo 'export LDFLAGS="-L/opt/homebrew/opt/openssl@1.1/lib"' >>  ~/.zshrc
xiazemin@xiazemindeMacBook-Pro librdkafka % echo ' export CPPFLAGS="-I/opt/homebrew/opt/openssl@1.1/include"'  >>  ~/.zshrc



% brew install  openssl zstd
Updating Homebrew...
^C
Error: No available formula with the name "openssl".
In formula file: /opt/homebrew/Library/Taps/homebrew/homebrew-core/Aliases/openssl
Expected to find class Openssl, but only found: OpensslAT11.


./configure --arch=arm64
make
 sudo make install
 安装成功
 
 go get -u github.com/confluentinc/confluent-kafka-go/kafka


https://github.com/confluentinc/confluent-kafka-go/issues/439

需要加-tags ，否则使用的是静态包，没法直接使用

 go run -tags dynamic main.go serve
 
 https://github.com/confluentinc/confluent-kafka-go/issues/591
 问题解决
