---
title: tcp-no-delay
layout: post
category: golang
author: 夏泽民
---
https://github.com/dastergon/gopheracademy-advent2019-tcp-no-delay
https://blog.gopheracademy.com/advent-2019/control-packetflow-tcp-nodelay/

The TCP implementations on most platforms offer algorithms and socket options to dictate the packet flow, connection lifespan and many more things. An algorithm that affects the network performance and is enabled by default on Linux, macOS, and Windows is Nagle’s algorithm. Nagle’s algorithm coalesces small packets and delays their delivery until an ACK is returned from the previously sent packet or an adequate amount of small packets is accumulated after a certain period. This process usually takes milliseconds but, having a latency-sensitive service or tight latency Service Level Objectives (SLOs), shaving off a couple of milliseconds might be worthwhile.

A cross-platform TCP socket option that comes helpful here is TCP_NODELAY. When enabled, it practically disables Nagle’s algorithm. Instead of coalescing small packets, it sends them to the pipe as soon as possible. In general, Nagle’s algorithm’s goal is to reduce the number of packets sent to save bandwidth and increase throughput with the trade-off sometimes to introduce increased latency to services. On the other hand, TCP_NODELAY might decrease throughput for small writes, but there are ways to mitigate this by using buffers on the application side.

In Go, TCP_NODELAY is enabled by default, but the standard library offers the ability to disable the behavior via the net.SetNoDelay method.

<!-- more -->
A small experiment
To observe what’s happening at the packet-level, and see the differences in packet arrival, we will use a tiny TCP client/server written in Go. Usually, we have inter-connected services across different regions, but for the sake of the experiment, we will experiment in our local machine. The full source code is also available on Github.

The server code (server.go):
package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
)

func main() {
    port := ":" + "8000"

    // Create a listening socket.
    l, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatal(err)
    }
    defer l.Close()

    for {
        // Accept new connections.
        c, err := l.Accept()
        if err != nil {
            log.Println(err)
            return
        }

        // Process newly accepted connection.
        go handleConnection(c)
    }
}
func handleConnection(c net.Conn) {
    fmt.Printf("Serving %s\n", c.RemoteAddr().String())

    for {
        // Read what has been sent from the client.
        netData, err := bufio.NewReader(c).ReadString('\n')
        if err != nil {
            log.Println(err)
            return
        }

        cdata := strings.TrimSpace(netData)
        if cdata == "GOPHER" {
            c.Write([]byte("GopherAcademy Advent 2019!"))
        }

        if cdata == "EXIT" {
            break
        }
    }
    c.Close()
}
The client code (client.go):
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    target := "localhost:8000"

    raddr, err := net.ResolveTCPAddr("tcp", target)
    if err != nil {
        log.Fatal(err)
    }

    // Establish a connection with the server.
    conn, err := net.DialTCP("tcp", nil, raddr)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Sending Gophers down the pipe...")

    for i := 0; i < 5; i++ {
        // Send the word "GOPHER" to the open connection.
        _, err = conn.Write([]byte("GOPHER\n"))
        if err != nil {
            log.Fatal(err)
        }
    }
}
To observe the behavior change, first execute tcpdump. You might have to change the network interface to match your own machine:

sudo tcpdump -X  -i lo0 'port 8000'
Then, execute the server (server.go) and the client (client.go).
go run server.go
In another terminal window execute:
go run client.go
Initially, if we look closer at the payload, we’ll notice that each write (Write()) of the word “GOPHER” is transmitted as a separate packet. Five in total. For brevity, I just posted only a couple of packets.

....
14:03:11.057782 IP localhost.58030 > localhost.irdmi: Flags [P.], seq 15:22, ack 1, win 6379, options [nop,nop,TS val 744132314 ecr 744132314], length 7
        0x0000:  4500 003b 0000 4000 4006 0000 7f00 0001  E..;..@.@.......
        0x0010:  7f00 0001 e2ae 1f40 80c5 9759 6171 9822  .......@...Yaq."
        0x0020:  8018 18eb fe2f 0000 0101 080a 2c5a 8eda  ...../......,Z..
        0x0030:  2c5a 8eda 474f 5048 4552 0a              ,Z..GOPHER.
14:03:11.057787 IP localhost.58030 > localhost.irdmi: Flags [P.], seq 22:29, ack 1, win 6379, options [nop,nop,TS val 744132314 ecr 744132314], length 7
        0x0000:  4500 003b 0000 4000 4006 0000 7f00 0001  E..;..@.@.......
        0x0010:  7f00 0001 e2ae 1f40 80c5 9760 6171 9822  .......@...`aq."
        0x0020:  8018 18eb fe2f 0000 0101 080a 2c5a 8eda  ...../......,Z..
        0x0030:  2c5a 8eda 474f 5048 4552 0a              ,Z..GOPHER.

...
If we disable TCP_NODELAY via the SetNoDelay method now, the code of the client looks like the following:
package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
    target := "localhost:8000"

    raddr, err := net.ResolveTCPAddr("tcp", target)
    if err != nil {
        log.Fatal(err)
    }

    // Establish a connection with the server.
    conn, err := net.DialTCP("tcp", nil, raddr)
    if err != nil {
        log.Fatal(err)
    }

    conn.SetNoDelay(false) // Disable TCP_NODELAY; Nagle's Algorithm takes action.

    fmt.Println("Sending Gophers down the pipe...")

    for i := 1; i <= 5; i++ {
        // Send the word "GOPHER" to the open connection.
        _, err = conn.Write([]byte("GOPHER\n"))
        if err != nil {
            log.Fatal(err)
        }
    }
}
Running again the client (go run client.go) with TCP_NODELAY disabled, Nagle’s algorithm is taking action and we get the following results:

14:27:20.120673 IP localhost.64086 > localhost.irdmi: Flags [P.], seq 8:36, ack 1, win 6379, options [nop,nop,TS val 745574362 ecr 745574362], length 28
        0x0000:  4500 0050 0000 4000 4006 0000 7f00 0001  E..P..@.@.......
        0x0010:  7f00 0001 fa56 1f40 07c9 d46f a115 3444  .....V.@...o..4D
        0x0020:  8018 18eb fe44 0000 0101 080a 2c70 8fda  .....D......,p..
        0x0030:  2c70 8fda 474f 5048 4552 0a47 4f50 4845  ,p..GOPHER.GOPHE
        0x0040:  520a 474f 5048 4552 0a47 4f50 4845 520a  R.GOPHER.GOPHER.
If we look closer at the payload, we see there are four coalesced "GOPHER" words that are sent in a single packet instead of separate packets.

Conclusion
To conclude, TCP_NODELAY is no panacea and needs experimentation before deciding to disable it or keep it enabled. However, it’s always good to know whether or not it is enabled by default in our favorite programming language. It might be the case that a service performs better with Nagle’s algorithm enabled (SetNoDelay(false)). The TCP_NODELAY option can be used in both sending and receiving sides. There’s no limitation. In our example, we experimented with it on the client-side. It all depends on the workload and the access we have on both the client and the server.

There are a handful of other socket options such as TCP_QUICKACK and TCP_CORK to experiment. Some of them might be platform-specific. Consequently, Go does not provide a method for controlling these options yet in the same way as TCP_NODELAY. However, we can do this through platform-specific packages. For example, to enable socket options in *nix systems, we can use the golang.org/x/sys/unix package and the SetsockoptInt method.

Example:
err = unix.SetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_QUICKACK, 1)
if err != nil {
  return os.NewSyscallError("setsockopt", err)
}
I highly recommend reading this blog post if you want to learn about Nagle’s algorithm, TCP_NODELAY, and similar algorithms.