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

Structure
---------
```
.
├── README.md
├── aio 异步组件，用于单线程模型，防止阻塞单线程。
│   ├── aio.go
│   ├── channel_queue.go    基于go chan实现的消息队列，容量有上限。
│   ├── event_loop.go       事件循环调度器。
│   ├── limited_counter.go  goroutine计数器，用于限制go chan总数。
│   └── queue.go            基于slice实现的消息队列，容量没有上限，没有阻塞主逻辑线程的风险。
├── bytes 字节流操作接口。
│   ├── bytes.go
│   ├── readonly_buffer.go  只读buffer，用于避免频繁的内存分配，提供便于字节流读的接口。
│   └── writeonly_buffer.go 只写buffer,用于避免频繁的内存拷贝，提供便于字节流写的接口。
├── codes   错误码预定义。
│   └── codes.go    
├── core    核心网络库组件，参考netty拦截器相关设计。
│   ├── acceptor_channel.go 接收器组件，用于服务器端。
│   ├── attr_map.go         线程安全的k-v结构session基类。
│   ├── channel.go          网络连接基类及接口定义。
│   ├── channel_builder.go  连接构建器。
│   ├── channel_context.go  pipeline上下文。
│   ├── channel_pipeline.go 网络事件拦截器。
│   ├── connector_channel.go 连接器组件，用于客户端。
│   ├── core.go             基础接口定义。         
│   ├── handler.go          消息处理器接口。
│   ├── sub_channel.go      代表网络连接。
│   ├── tcp tcp协议实现。
│   │   ├── builder.go
│   │   ├── client.go
│   │   ├── raw_conn.go
│   │   └── server.go
│   ├── udp
│   └── ws  websocket协议实现。
│       ├── builder.go
│       ├── client.go
│       ├── raw_conn.go
│       └── server.go
├── encoding    序列化封装。
│   ├── encoding.go
│   ├── json
│   │   └── json.go
│   └── proto
│       └── proto.go
├── gnetlog
│   └── gnetlog.go
├── handler 网络消息处理器预定义的一些实现。
│   ├── combined_decoder.go
│   ├── combined_encoder.go
│   ├── consts.go
│   ├── encryption_rc4.go
│   ├── idle_state_handler.go
│   ├── message_decoder.go
│   ├── message_deserializer.go
│   ├── message_encoder.go
│   ├── message_processor.go
│   ├── message_serializer.go
│   ├── packet_id_parser.go
│   ├── packet_length_decoder.go
│   ├── packet_length_prepender.go
│   └── string_encoder.go
├── internal
│   └── internal.go
├── message 消息相关逻辑包装。
│   ├── idparser
│   │   ├── idparser.go
│   │   ├── uint16.go
│   │   └── uint32.go
│   ├── message.go
│   ├── packet
│   │   ├── raw_packet.go
│   │   └── var_packet.go
│   ├── processor.go
│   └── register.go
├── nginet.go
├── pool    协程池。
│   └── pool.go
├── rpc
│   └── rpc.go
├── safe
│   └── safe_call.go
├── shortid
│   └── shortid.go
├── test    测试用例。
│   ├── attrmap_test.go
│   ├── buffer_test.go
│   ├── encoding_test.go
│   ├── evtloop_test.go
│   ├── msg_test.go
│   ├── pipeline_test.go
│   ├── pool_test.go
│   ├── queue_test.go
│   ├── tcp_channel_client_test.go
│   ├── tcp_channel_test.go
│   ├── telnet_test.go
│   ├── test.go
│   ├── waitgroup_test.go
│   └── ws_channel_test.go
└── version.go
```