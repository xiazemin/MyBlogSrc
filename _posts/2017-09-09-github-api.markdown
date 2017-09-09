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

<script type="text/javascript">
 $.ajax({
        type: "GET",
        url:"https://github.com/xiazemin/MyBlogComment/issues/",
        dataType: 'json',
        async: false,
        success: function(json) {
           console.log(json);
            for (var i = 0; i < json.length; i++) {
                var title = json[i].title; // Blog title
                var comments_url = json[i].comments_url;
                if (title == blogName) {
                    console.log("该文章存在评论")
                    $('#commentsList').attr("data_comments_url", comments_url);
                    setComment(comments_url);
                    break;
                }
                $("#commentsList").children().remove();
                $("#commentsList").removeAttr('data_comments_url');

            }
        }
    });

 </script>