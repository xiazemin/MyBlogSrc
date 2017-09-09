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
或者
https://github.com/login/oauth/authorize?client_id=981ba8c916c262631ea0&redirect_uri=https://xiazemin.github.io/MyBlog/
浏览器返回：
https://xiazemin.github.io/MyBlog/?code=4d4e3333ec04f9bcc5c7

https://github.com/login/oauth/access_token?client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e&code=4d4e3333ec04f9bcc5c7
浏览器返回文件内容：
access_token=7023f1b08df35d30032f1aba02202a0ec1618389&scope=public_repo&token_type=bearer


https://api.github.com/user?access_token=7023f1b08df35d30032f1aba02202a0ec1618389
浏览器返回用户信息json

https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments?access_token=7023f1b08df35d30032f1aba02202a0ec1618389
返回评论信息json




#https://github.com/login/oauth/authorize?client_id=01a8bf26baecfa8db577