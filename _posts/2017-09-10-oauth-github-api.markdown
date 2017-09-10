---
title: oauth-github-api
layout: post
category: jekyll
author: 夏泽民
---
<!-- more -->

According to the documentation: http://developer.github.com/v3/#cross-origin-resource-sharing

Any domain that is registered as an OAuth Application is accepted.
To register you application go to: https://github.com/settings/applications

You need to be posting from the same domain that your application is registered on. If you are trying to test locally you may need to modify your hosts file and run your server on port 80.

$curl -i https://api.github.com -H "Origin: https://xiazemin.github.io/MyBlog/jekyll/2017/09/09/github-api.html"
HTTP/1.1 200 OK
Date: Sun, 10 Sep 2017 14:01:04 GMT
Content-Type: application/json; charset=utf-8
Content-Length: 2165
Server: GitHub.com
Status: 200 OK
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 58
X-RateLimit-Reset: 1505055460
Cache-Control: public, max-age=60, s-maxage=60
Vary: Accept
ETag: "7dc470913f1fe9bb6c7355b50a0737bc"
X-GitHub-Media-Type: github.v3; format=json
Access-Control-Expose-Headers: ETag, Link, X-GitHub-OTP, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, X-OAuth-Scopes, X-Accepted-OAuth-Scopes, X-Poll-Interval
Access-Control-Allow-Origin: *
Content-Security-Policy: default-src 'none'
Strict-Transport-Security: max-age=31536000; includeSubdomains; preload
X-Content-Type-Options: nosniff

https://developer.github.com/v3/#cross-origin-resource-sharing

https://developer.github.com/apps/building-integrations/setting-up-and-registering-oauth-apps/about-authorization-options-for-oauth-apps/

https://developer.github.com/v3/#authentication
https://developer.github.com/v3/oauth_authorizations/#get-a-single-authorization
http://www.ituring.com.cn/book/tupubarticle/11824

原因：
https://github.com/     不允许跨于
https://api.github.com/允许跨域
http://www.membrane-soa.org/service-proxy-doc/4.2/oauth2-github.htm

通过代理解决跨域问题：
$ curl -i https://gh-oauth.imsun.net -d 'client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e&code=212ab8ead2246b853e75'

HTTP/1.1 200 OK
x-powered-by: Express
access-control-allow-origin: *
access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept, Authorization
content-type: application/json; charset=utf-8
content-length: 171
etag: W/"ab-oF8ktyaUXM1tq3Z8pnAn/45npKk"
date: Sun, 10 Sep 2017 15:14:07 GMT
connection: close
Vary: Accept-Encoding

{"error":"bad_verification_code","error_description":"The code passed is incorrect or expired.","error_uri":"https://developer.github.com/v3/oauth/#bad-verification-code"}


XMLHttpRequest cannot load https://gh-oauth.imsun.net/?client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e&code=212ab8ead2246b853e75&code=e06fe81abe307693bd63. The value of the 'Access-Control-Allow-Origin' header in the response must not be the wildcard '*' when the request's credentials mode is 'include'. Origin 'https://xiazemin.github.io' is therefore not allowed access. The credentials mode of requests initiated by the XMLHttpRequest is controlled by the withCredentials attribute.


 xhrFields:{
        withCredentials:true
    },
改为false
