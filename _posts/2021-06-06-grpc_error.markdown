---
title: grpc error
layout: post
category: golang
author: 夏泽民
---
https://github.com/grpc-ecosystem/grpc-gateway/blob/v1.14.8/runtime/errors.go

import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    ....
)

return nil, status.Error(codes.PermissionDenied, "PERMISSION_DENIED_TEXT")


<!-- more -->
// client
    assignvar, err := s.MyFunctionCall(ctx, ...)
    if err != nil {
        if e, ok := status.FromError(err); ok {
            switch e.Code() {
            case codes.PermissionDenied:
                fmt.Println(e.Message()) // this will print PERMISSION_DENIED_TEST
            case codes.Internal:
                fmt.Println("Has Internal Error")
            case codes.Aborted:
                fmt.Println("gRPC Aborted the call")
            default:
                fmt.Println(e.Code(), e.Message())
            }
        }
        else {
            fmt.Printf("not able to parse error returned %v", err)
        }
    }

https://stackoverflow.com/questions/52969205/how-to-assert-grpc-error-codes-client-side-in-go/52972944

https://diabloneo.github.io/2018/12/10/golang-grpc-error-code/

https://jbrandhorst.com/post/grpc-errors/
