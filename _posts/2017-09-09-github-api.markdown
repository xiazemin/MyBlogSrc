---
title: github-api
layout: post
category: jekyll
author: 夏泽民
---
<!-- more -->
curl -d "json_data={
    owner: 'xiazemin',
    repo: 'MyBlogComment',
    oauth: {
        client_id: '981ba8c916c262631ea0',
        client_secret: 'a52260ef92de69011ccd1cf355b973ef11d6da0e',
}"  https://github.com/xiazemin/MyBlogComment/issues/33 
直接curl 授权成功
curl -u xiazemin:xxxxxxxxxxx  https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments
curl -H "Authorization: token xxxxxxxxxxx"  https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments
自己请求不成功

<script type="text/javascript" src="{{site.baseurl}}/js/utils.js">
</script>
<script type="text/javascript">
alert(window.location.search);
var search=window.location.search;
function loadFun(){
     alert(search);
     alert(Query.parse(search));
     code=Query.parse(search).code;
     document.getElementById("code").innerHTML="code:"+code;
    console.log(Query.parse(search));
    if(code){
        url="https://github.com/login/oauth/access_token??scope=public_repo&redirect_uri=https%3a%2f%2fxiazemin.github.io%2fMyBlog%2fjekyll%2f2017%2f09%2f09%2fgithub-api.html&client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e&code="+code;
        alert(url);
        var script = document.createElement('script');
    script.setAttribute('src', url);
    // 把script标签加入head，此时调用开始
    document.getElementsByTagName('head')[0].appendChild(script);
    //     $.ajax({
    // type: "GET",
    // url: "https://github.com/login/oauth/access_token?client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e&code="+code,
    // success: function(json) {
    //         console.log(json);
    //          document.getElementById("token").innerHTML="token:"+json.access_token;
    //         alert(json.access_token);
    //         }
    // });
  }
}
</script>
<body onload="loadFun()">
<span id="code"></span>
<br/>
<hr/>
<span id="token"></span>
<br/>
<hr/>

<a href="https://github.com/login/oauth/authorize?scope=public_repo&redirect_uri=https%3a%2f%2fxiazemin.github.io%2fMyBlog%2fjekyll%2f2017%2f09%2f09%2fgithub-api.html&client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e">登入</a>
</body>
<script src="https://imsun.github.io/gitment/dist/gitment.browser.js"></script>
<!--script type="text/javascript">
//postRequest();
try{
    var flightHandler = function(data){
        alert('你查询的航班结果是：票价 ' + data.price + ' 元，' + '余票 ' + data.tickets + ' 张。');
    };
    // 提供jsonp服务的url地址（不管是什么类型的地址，最终生成的返回值都是一段javascript代码）
    var url = "https://github.com/login/oauth/authorize?client_id=981ba8c916c262631ea0";
    // 创建script标签，设置其属性
    var script = document.createElement('script');
    script.setAttribute('src', url);
    // 把script标签加入head，此时调用开始
    document.getElementsByTagName('head')[0].appendChild(script);
    console.log(Query.parse());
}catch(ex){
    console.log(ex);
    console.log(document);
}
//get access code 
</script-->
<!--script type="text/javascript">
$.ajax({
type: "GET",
url: "https://github.com/login/oauth/authorize?scope=public_repo&redirect_uri=https%3a%2f%2fxiazemin.github.io%2fMyBlog%2fjekyll%2f2017%2f09%2f09%2fgithub-api.html&client_id=981ba8c916c262631ea0&client_secret=a52260ef92de69011ccd1cf355b973ef11d6da0e",
//'"+encodeURIComponent("{{site.url}}{{site.baseurl}}/token.html")+"'",
dataType: 'json',
    async: false,
    xhrFields:{
        withCredentials:true
    },
    crossDomain:true,
    success: function(json) {
        alert(Query.parse());
console.log(Query.parse());//console.log(query);//åconsole.log(json);
    }
});
 </script-->
<script type="text/javascript">
//  $.ajax({
//         type: "GET",
//         url:"https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments",
//         dataType: 'json',
//         async: false,
//         success: function(json) {
//            console.log(json);
//            console.log(json[json.length-1].body);

// 			if(json.length>0){
// 			   json[json.length-1].body+=1;

// 			   $.ajax({
// 					type: "post",
// 					url:"https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments",
// 					dataType: 'json',
// 					async: false,
// 					beforeSend: function(request) {
//             request.setRequestHeader(
//             	"Authorization","token xxxxxxxxxxx");},
// //"Authorization","Basic " + btoa("xiazemin:xxx"));},
// 					//headers: {
//                // "Authorization": "Basic " + btoa("xiazemin :xxxxxxxxxxx")
//            // },
// 					data:{"body": "Me too"},
// 					success: function(json) {
// 					console.log(json);
// 					console.log(json[json.length-1].body);
//                    },
//                    error: function () {
//                 }
// 				});
// 			}


//         }
//     });
 </script>