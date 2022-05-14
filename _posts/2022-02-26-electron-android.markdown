---
title: electron-android cordova
layout: post
category: node
author: 夏泽民
---
改造可行性： electron应用也是在web项目上套了一层应用壳而已，所以移植到混合应用hybird上面，通过webview支持web显示也是可行的，所以！改造只要在浏览器上面跑得通就行！但是有些功能可能会稍有不同，需要做依据需求删除。

https://www.jianshu.com/p/4ed454cf2b95
https://www.jianshu.com/p/a32542277b83

https://stackoverflow.com/questions/62420427/how-do-i-make-my-electron-app-the-default-for-opening-files
https://github.com/electron-userland/electron-prebuilt/issues/161
<!-- more -->
https://www.electronjs.org/apps

What is Apache Cordova?
Apache Cordova is an open-source framework, that allows developers to use standard web technologies — HTML5, CSS3, and JavaScript to build “Hybrid” cross-platform applications.
With this framework, we can deploy to multiple platforms using a single set of source code. It removes the need for learning native languages for simple applications.

Cordova achieved this by using the platform’s native WebView component. Additionally, there is a JavaScript-to-Native bridge that allows an application to leverage native features, for example, geolocation, camera or even system functionality. If a plugin is missing and you need access to a specific native feature, you can always develop a plugin that supports the feature but will require native language understanding.
What is Electron?
Initially, Electron was used as the base of Atom, GitHub’s hackable text editor. Eventually, both Electron and Atom became open sourced. Electron, similar to Apache Cordova, also makes it possible to build cross-platform applications with web standard technologies, except it is targeting the desktop community. Electron was able to accomplish this by combining Chromium and Node.js into a single runtime. Ever since it has been released, it has been growing in popularity and used by many open source developers, startups, and well-established companies.
Apache Cordova Now Supports Electron
On 28th of February 2019, Cordova’s team announced their first official release of the Cordova Electron platform that supports Electron v4.0.1 and electron-builder.
Cordova used to support Ubuntu and continues to support the Windows and macOS platforms, but with the release of Electron platform, it provides an alternative means for building all three major platforms.
Additionally, Electron provides right out of the box its own support implementation for menus, dock integration, finder integration, documents, and etc without the need for additional plugins.
Why Use Cordova Electron?
Takes an existing hybrid mobile or desktop application and expand to both markets simultaneously.
Take advantage of existing browser supported Cordova plugins.

https://medium.com/the-web-tub/electron-on-cordova-29ede5d6d789

https://cordova.apache.org/announcements/2019/02/28/cordova-electron-release-1.0.0.html

https://cordova.apache.org/docs/en/latest/guide/platforms/electron/index.html#requirements-and-support

https://cordova.apache.org/
https://github.com/apache/cordova
https://github.com/apache/cordova-android
https://github.com/apache/cordova-ios
https://cordova.apache.org/docs/en/10.x/guide/platforms/ios/