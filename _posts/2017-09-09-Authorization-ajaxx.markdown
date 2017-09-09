---
title: Authorization-ajax
layout: post
category: jekyll
author: 夏泽民
---

Get Github Authorization Token with proper scope,print to console
<!-- more -->
{% highlight javascript linenos %}
$.ajax({
url: 'https://api.github.com/authorizations',
type: 'POST',
beforeSend: function(xhr) {
 	xhr.setRequestHeader("Authorization",
 	"Basic" + btoa("USERNAME:PASSWORD"));
},
data: '{"scopes":["gist"],"note":"ajax gist test for a user"}'
}).done(function(response) {
 	console.log(response);
});
//Create a Gist with token from above
$.ajax({
url:'https://api.github.com/gists',
type:'POST',
beforeSend: function(xhr) {
 	xhr.setRequestHeader("Authorization",
 	"token TOKEN-FROM-AUTHORIZATION-CALL");
 	},
data: '{"description": "a gist for a user with token api call via ajax","public": true,"files": {"file1.txt": {"content": "String file contents via ajax"}}}'
}).done(function(response) {
 	console.log(response);
});
{% endhighlight %}