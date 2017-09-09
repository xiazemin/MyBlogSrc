---
title: github-openapi
layout: post
category: jekyll
author: 夏泽民
---

github openapi 学习

参考资料：
1. https://github.com/xiazemin/TinyBlog/blob/master/js/core.js
1. https://developer.github.com/v3/#conditional-requests
1. https://developer.github.com/v3/issues/comments/#create-a-comment
1. https://developer.github.com/v3/issues/#list-issues-for-a-repository
1. https://developer.github.com/v3/auth/#basic-authentication
<!-- more -->

不需要授权的页面：
https://api.github.com/repos/xiazemin/MyBlogComment/issues/33/comments
[
  {
    "url": "https://api.github.com/repos/xiazemin/MyBlogComment/issues/comments/328255639",
    "html_url": "https://github.com/xiazemin/MyBlogComment/issues/33#issuecomment-328255639",
    "issue_url": "https://api.github.com/repos/xiazemin/MyBlogComment/issues/33",
    "id": 328255639,
    "user": {
      "login": "xiazemin",
      "id": 6873571,
      "avatar_url": "https://avatars1.githubusercontent.com/u/6873571?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/xiazemin",
      "html_url": "https://github.com/xiazemin",
      "followers_url": "https://api.github.com/users/xiazemin/followers",
      "following_url": "https://api.github.com/users/xiazemin/following{/other_user}",
      "gists_url": "https://api.github.com/users/xiazemin/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/xiazemin/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/xiazemin/subscriptions",
      "organizations_url": "https://api.github.com/users/xiazemin/orgs",
      "repos_url": "https://api.github.com/users/xiazemin/repos",
      "events_url": "https://api.github.com/users/xiazemin/events{/privacy}",
      "received_events_url": "https://api.github.com/users/xiazemin/received_events",
      "type": "User",
      "site_admin": false
    },
    "created_at": "2017-09-09T05:23:35Z",
    "updated_at": "2017-09-09T05:23:35Z",
    "author_association": "OWNER",
    "body": "101"
  }
]

