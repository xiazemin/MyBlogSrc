---
title: oauth-github
layout: post
category: jekyll
author: 夏泽民
---
<!-- more -->

<script src="{{site.baseurl}}/js/oauth.min.js">
OAuth.initialize('981ba8c916c262631ea0');
console.log(OAuth);
OAuth.popup('https://github.com/login/oauth/authorize')
.done(function(result) {
   console.log(result);
    result.post('/message', {
        data: {
            user_id: 93,
            content: 'Hello Mr. 93 !'
        }
    })
    .done(function (response) {
        //this will display the id of the message in the console
        console.log(response.id);
    })
    .fail(function (err) {
        //handle error with err
    });
})
.fail(function (err) {
    //handle error with err
});
</script>

https://github.com/login/oauth/authorize?client_id=981ba8c916c262631ea0
https://github.com/login/oauth/authorize?client_id=981ba8c916c262631ea0&response_type=code
或者
https://github.com/login/oauth/authorize?client_id=981ba8c916c262631ea0&redirect_uri=https://xiazemin.github.io/MyBlog/
浏览器返回：
https://xiazemin.github.io/MyBlog/?code=xxxxx

https://github.com/login/oauth/access_token?client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e&code=212ab8ead2246b853e75
浏览器返回文件内容：
access_token=7023f1b08df35d30032f1aba02202a0ec1618389&scope=public_repo&token_type=bearer


https://api.github.com/user?access_token=7xxxxx
浏览器返回用户信息json

https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments?access_token=xxxxxxx
返回评论信息json

loginlink:
https://github.com/login/oauth/authorize?scope=public_repo&redirect_uri=https%3A%2F%2Fxiazemin.github.io%2FMyBlog%2Fjekyll%2F2017%2F09%2F09%2Fstatics.html&client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e"

返回json
$     curl  -i -X POST -H "Accept: application/json" 'https://github.com/login/oauth/access_token?client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e&callback=parseQueryString&code=39b990e457b27245464e&callback=jQuery112004444319529793137_1505049039697&_=1505049039698'

{"error":"bad_verification_code","error_description":"The code passed is incorrect or expired.","error_uri":"https://developer.github.com/v3/oauth/#bad-verification-code"}





#https://github.com/login/oauth/authorize?client_id=01a8bf26baecfa8db577