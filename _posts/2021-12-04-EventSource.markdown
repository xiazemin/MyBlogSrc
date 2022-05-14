---
title: EventSource
layout: post
category: node
author: 夏泽民
---
An EventSource instance opens a persistent connection to an HTTP server, which sends events in text/event-stream format. The connection remains open until closed by calling

Unlike WebSockets, server-sent events are unidirectional; that is, data messages are delivered in one direction, from the server to the client (such as a user's web browser). That makes them an excellent choice when there's no need to send data from the client to the server in message form. For example, EventSource is a useful approach for handling things like social media status updates, news feeds, or delivering data into a client-side storage mechanism like IndexedDB or web storage.

https://developer.mozilla.org/en-US/docs/Web/API/EventSource
<!-- more -->
https://github.com/wanming/pipelining