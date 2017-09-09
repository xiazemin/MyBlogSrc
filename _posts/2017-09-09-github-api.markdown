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
curl -u xiazemin:352e7308a2aa563f6ee6d6a94cc6f58708d2a7e1  https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments
curl -H "Authorization: token 352e7308a2aa563f6ee6d6a94cc6f58708d2a7e1"  https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments
自己请求不成功
<script type="text/javascript">
 $.ajax({
        type: "GET",
        url:"https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments",
        dataType: 'json',
        async: false,
        success: function(json) {
           console.log(json);
           console.log(json[json.length-1].body);

			if(json.length>0){
			   json[json.length-1].body+=1;

			   $.ajax({
					type: "post",
					url:"https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments",
					dataType: 'json',
					async: false,
					beforeSend: function(request) {
            request.setRequestHeader(
            	"Authorization","token 352e7308a2aa563f6ee6d6a94cc6f58708d2a7e1");},
//"Authorization","Basic " + btoa("xiazemin:xxx"));},
					//headers: {
               // "Authorization": "Basic " + btoa("xiazemin :352e7308a2aa563f6ee6d6a94cc6f58708d2a7e1")
           // },
					data:{"body": "Me too"},
					success: function(json) {
					console.log(json);
					console.log(json[json.length-1].body);
                   },
                   error: function () {
                }
				});
			}


        }
    });

 </script>