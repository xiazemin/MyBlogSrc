---
title: http basic authentication
layout: post
category: golang
author: 夏泽民
---
何为basic authentication
In the context of a HTTP transaction, basic access authentication is a method for a HTTP user agent to provide a user name and password when making a request.

但是这种认证方式有很多的缺点： 
虽然基本认证非常容易实现，但该方案建立在以下的假设的基础上，即：客户端和服务器主机之间的连接是安全可信的。特别是，如果没有使用SSL/TLS这样的传输层安全的协议，那么以明文传输的密钥和口令很容易被拦截。该方案也同样没有对服务器返回的信息提供保护。
SplitN slices s into substrings separated by sep and returns a slice of the substrings between those separators. If sep is empty, SplitN splits after each UTF-8 sequence. 
看到了吗，golang中的strings包为我们提供了强大的字符串处理能力。

package main

import (
    "encoding/base64"
    "net/http"
    "strings"
)

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
    s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
    if len(s) != 2 {
        return false
    }

    b, err := base64.StdEncoding.DecodeString(s[1])
    if err != nil {
        return false
    }

    pair := strings.SplitN(string(b), ":", 2)
    if len(pair) != 2 {
        return false
    }

    return pair[0] == "user" && pair[1] == "pass"
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if checkAuth(w, r) {
            w.Write([]byte("hello world!"))
            return
        }

        w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
        w.WriteHeader(401)
        w.Write([]byte("401 Unauthorized\n"))
    })

    http.ListenAndServe(":8080, nil)
}

客户端：
func basicAuth(username, password string) string {
  auth := username + ":" + password
   return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error{
 req.Header.Add("Authorization","Basic " + basicAuth("username1","password123"))
 return nil
}

func main() {

  client := &http.Client{
    Jar: cookieJar,
    CheckRedirect: redirectPolicyFunc,
}

req, err := http.NewRequest("GET", "http://localhost/", nil)
    req.Header.Add("Authorization","Basic " + basicAuth("username1","password123"))

    resp, err := client.Do(req)
}


goji/httpauth
HTTP Authentication middlewares
获取： 
go get github.com/goji/httpauth
https://github.com/goji/httpauth
package main

import (
    "net/http"

    "github.com/goji/httpauth"
)

func main() {
    finalHandler := http.HandlerFunc(final)
    authHandler := httpauth.SimpleBasicAuth("username", "password")

    http.Handle("/", authHandler(finalHandler))
    http.ListenAndServe(":8080", nil)
}

func final(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}
curl --i username:password@localhost:8080
HTTP/1.1 200 OK
Date: Thu, 01 Feb 2018 05:27:33 GMT
Content-Length: 2
Content-Type: text/plain; charset=utf-8

OK

 或者浏览器
 http://username:password@localhost:8080
<!-- more -->
BasicAuth returns the username and password provided in the request's Authorization header, if the request uses HTTP Basic Authentication. See RFC 2617, Section 2.
https://golang.org/pkg/net/http/#Request.BasicAuth


// BasicAuth wraps a handler requiring HTTP basic auth for it using the given
// username and password and the specified realm, which shouldn't contain quotes.
//
// Most web browser display a dialog with something like:
//
//    The website says: "<realm>"
//
// Which is really stupid so you may want to set the realm to a message rather than
// an actual realm.
func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {

        user, pass, ok := r.BasicAuth()

        if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
            w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
            w.WriteHeader(401)
            w.Write([]byte("Unauthorised.\n"))
            return
        }

        handler(w, r)
    }
}

...

http.HandleFunc("/", BasicAuth(handleIndex, "admin", "123456", "Please enter your username and password for this site"))


func checkAuth(w http.ResponseWriter, r *http.Request) bool {
    s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
    if len(s) != 2 { return false }

    b, err := base64.StdEncoding.DecodeString(s[1])
    if err != nil { return false }

    pair := strings.SplitN(string(b), ":", 2)
    if len(pair) != 2 { return false }

    return pair[0] == "user" && pair[1] == "pass"
}

yourRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if checkAuth(w, r) {
        yourOriginalHandler.ServeHTTP(w, r)
        return
    }

    w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
    w.WriteHeader(401)
    w.Write([]byte("401 Unauthorized\n"))
})


HTTP authentication with PHP
It is possible to use the header() function to send an "Authentication Required" message to the client browser causing it to pop up a Username/Password input window. Once the user has filled in a username and a password, the URL containing the PHP script will be called again with the predefined variables PHP_AUTH_USER, PHP_AUTH_PW, and AUTH_TYPE set to the user name, password and authentication type respectively. These predefined variables are found in the $_SERVER array. Both "Basic" and "Digest" (since PHP 5.1.0) authentication methods are supported. See the header() function for more information.

http://php.net/manual/en/features.http-auth.php
http://www.php.net/manual/en/function.parse-url.php
<?php 

$url = 'http://usr:pss@example.com:81/mypath/myfile.html?a=b&b[]=2&b[]=3#myfragment'; 
if ($url === unparse_url(parse_url($url))) { 
  print "YES, they match!\n"; 
} 

function unparse_url($parsed_url) { 
  $scheme   = isset($parsed_url['scheme']) ? $parsed_url['scheme'] . '://' : ''; 
  $host     = isset($parsed_url['host']) ? $parsed_url['host'] : ''; 
  $port     = isset($parsed_url['port']) ? ':' . $parsed_url['port'] : ''; 
  $user     = isset($parsed_url['user']) ? $parsed_url['user'] : ''; 
  $pass     = isset($parsed_url['pass']) ? ':' . $parsed_url['pass']  : ''; 
  $pass     = ($user || $pass) ? "$pass@" : ''; 
  $path     = isset($parsed_url['path']) ? $parsed_url['path'] : ''; 
  $query    = isset($parsed_url['query']) ? '?' . $parsed_url['query'] : ''; 
  $fragment = isset($parsed_url['fragment']) ? '#' . $parsed_url['fragment'] : ''; 
  return "$scheme$user$pass$host$port$path$query$fragment"; 
} 

?>


