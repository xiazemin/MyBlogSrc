---
title: serverless
layout: post
category: k8s
author: 夏泽民
---
https://github.com/vercel/fun
ocal serverless function λ development runtime.
// example/index.js
exports.handler = function(event, context, callback) {
	callback(null, { hello: 'world' });
};
<!-- more -->
