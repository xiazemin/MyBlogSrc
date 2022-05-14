---
title: protobuf2json
layout: post
category: node
author: 夏泽民
---
protobuf2json是一个json 和proto互转的命令行工具

https://github.com/revinate/protobuf2json


<!-- more -->
https://www.dllhook.com/post/187.html

$ protobuf2json -h

  Usage: protobuf2json [options]

  Convert protobuf to JSON


  Options:

    -V, --version           output the version number
    -d, --directory <path>  path to base directory containing proto file and all its imports
    -p, --proto <path>      path to proto file relative to base directory
    -t, --type <name>       protobuf message type
    -m, --multi             multiple protobuf messages with prefixed signed 32-bit big endian message length
    -h, --help              output usage information
   


$ protobuf2json -d protos -p my_message.proto -t MyMessage -m < my_messages.bin



$ json2protobuf -h

  Usage: json2protobuf [options]

  Convert JSON to protobuf


  Options:

    -V, --version           output the version number
    -d, --directory <path>  path to base directory containing proto file and all its imports
    -p, --proto <path>      path to proto file relative to base directory
    -t, --type <name>       protobuf message type
    -m, --multi             multiple JSON messages separated by newlines
    -h, --help              output usage information
    
    
$ json2protobuf -d protos -p my_message.proto -t MyMessage -m < my_messages.json


Usage with kafkacat
This utility can be used with kafkacat to publish protobuf-encoded messages to Kafka topics, and to examine the contents of Kafka topics containing protobuf-encoded messages.


$ json2protobuf -d protos -p my_message.proto -t MyMessage < my_message.json | kafkacat -P -b <broker> -t <topic> -D \t


$ kafkacat -C -b <broker> -t <topic> -e -o beginning -f '%R%s' | protobuf2json -d protos -p my_message.proto -t MyMessage -m

