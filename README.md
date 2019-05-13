# Nginet
An easy, scalable, high performance network communication framework. It's inspired by netty.

Installation
------------

To install this package, you need to install Go and setup your Go workspace on your computer. The simplest way to install the library is to run:

```
$ go get -u github.com/amsalt/nginet
```

Prerequisites
-------------

This sugguest Go 1.11 or later. 

Features
-----------
- Flexible threading model.
- Multiple communication protocol support, such as TCP/UDP/Websocket.
- Handler pipeline, enables the ability to contol inbound&outbound message process.
- Multiple message handlers are built in
- Protocol serialization, integrated with JSON and protobuf.
- AIO, designed for asynchronous call to avoid thread-safe control.
- Optimized bytes buffer, provide read-only and write-only buffer to reduce memory copy and provide more convenient interface
- Integrated with logger, it's just a logger facade, you can choose any logger you like.
- Goroutine tools, such as pool, limited counter, goroutine-safe queue.

The process of message processing
---------------------------------
## receive message
1. network read
1. unpack 
1. parse message body
1. deserialization
1. distribute the message to the processor

## send message
1. write to buffer
1. serialization
1. encode message
1. pack
1. network write