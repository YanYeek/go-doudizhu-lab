# due 学习 03：第一条消息路由，与一个编解码器的坑

> 日期：2026-06-13
>
> 本课目标：让「客户端 → Gate → Node → 响应」整条链路第一次真正跑通，
> 并理解为什么两端的编解码器必须一致。前两课的 Container / Gate 概念不再重复。

## 1. 这一步到底加了什么

上一课 Gate 能接受连接，但没有 Node，消息无处可去。这一课补上 Node 和
第一条路由，链路闭合：

```text
Go 测试客户端 ──ws──▶ Gate ──转发──▶ Node(路由1: Greet) ──Response──▶ 客户端
```

落地为四小块（都很薄，刻意保持最小）：

| 文件 | 作用 | 一句话 |
|------|------|--------|
| `internal/route/route.go` | 协议定义 | 路由号 `Greet=1` + 请求/响应结构体 |
| `internal/greet/greet.go` | 业务处理器 | 解析名字，回一句欢迎语 |
| `cmd/server/main.go` | 组装 | 在原 Gate 旁边加一个 Node，注册 Greet 路由 |
| `cmd/testclient/main.go` | 测试客户端 | 连 Gate、发一条 Greet、打印响应后退出 |

## 2. 一条路由是怎么处理的

Node 的处理器签名是 `func(ctx node.Context)`。`ctx` 装着这条消息的全部上下文
（哪条连接、序列号、路由号）。两个动作就够了：

```go
func handle(ctx node.Context) {
    var req route.GreetReq
    ctx.Parse(&req)                 // 字节 → 请求结构体
    ctx.Response(&route.GreetRes{   // 响应结构体 → 字节，沿原路由号发回
        Message: "你好，" + req.Name + "，欢迎来到斗地主服务器",
    })
}
```

服务端的 `Response` 会用**同一个路由号**把响应发回，所以测试客户端也注册
`AddRouteHandler(route.Greet, ...)` 来接收它。请求和响应共用路由号，这是 due 的
约定。

## 3. 真实踩到的坑：proto vs json 编解码器

第一次跑测试客户端，直接报错：

```text
发送请求失败: can't marshal a value that not implements proto.Buffer interface
```

**线索其实在服务端启动横幅里就给了**：

```text
| Node   | Codec: proto |
| Client | Codec: proto |
```

due 默认用 **proto** 编解码器，而它要求消息体实现 protobuf 接口。我们的
`GreetReq` 只是个带 json tag 的普通结构体，proto 编不了 → 报错。

修复：在配置里把两端都改成 json（结构体本来就是 json tag）：

```toml
[cluster.node]
codec = "json"

[cluster.client]
codec = "json"
```

## 4. 这个坑背后的原理（本课最值得记的一点）

为什么改一个地方不够、必须两端都改？因为：

```text
客户端  ──编码──▶  Gate（只转发字节，从不拆开看）──▶  Node  ──解码
   ↑ 用自己的 codec 编                                    ↑ 用自己的 codec 解
```

**Gate 是编解码器无关的**——它只搬运字节，不关心里面是什么。真正做编码/解码的
是两个端点（客户端和 Node）。所以两端的 codec 必须说同一种"语言"，否则一边写的
另一边读不懂。这和 HTTP 里 `Content-Type` 要前后端一致是同一个道理。

选 json 而不是 proto，是因为本项目第一阶段用「json + 手写结构体」起步，最直观；
等协议稳定、追求性能时再考虑换 proto（那时要写 `.proto` 并生成代码）。

## 5. 怎么自己复跑

```bash
python scripts/python/dev.py up   # 起 Redis/etcd + 前台 Gate+Node
go run ./cmd/testclient           # 另一个终端
# 期望输出：链路打通，收到服务端响应: 你好，测试客户端，欢迎来到斗地主服务器
```

## 6. 本课记住三句话

1. Node 处理器只做两件事：`Parse` 收请求、`Response` 回响应，共用路由号。
2. Gate 只转发字节，不做编解码；编解码发生在客户端和 Node 两个端点。
3. 两端 codec 必须一致——这次的报错就是 proto/json 不一致造成的。

## 7. 下一小步（保持小）

链路通了，但还没有"玩家"概念。下一步可以只做一件事：让客户端先发一条
登录/绑定消息，把连接和一个 uid 绑定起来（`ctx.BindNode`），为之后的房间和
对局做准备——仍然是一条路由的量级，不扩模块。
