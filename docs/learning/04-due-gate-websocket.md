# due 学习 02：Gate 与 WebSocket Server

> 日期：2026-06-13
>
> 本课目标：让 due 第一次真正监听网络端口，并理解 Gate、WebSocket Server、
> Locator、Registry 与 Container 的职责边界。

## 1. 从空 Container 到可以连接的 Gate

上一课的结构是：

```text
Container
└── 暂无组件
```

`container.Serve()` 会让进程持续运行，但没有任何网络组件监听端口，因此客户端
无法连接。

本课加入 Gate 后，结构变为：

```text
Container
└── Gate
    ├── WebSocket Server
    ├── Redis Locator
    └── etcd Registry
```

Container 负责启动和关闭 Gate。Gate 启动时，再启动它持有的 WebSocket Server，
并向 Registry 注册自己。

## 2. 四个角色分别负责什么

### Container：管理组件生命周期

Container 不直接处理网络消息。它负责按统一顺序初始化、启动、关闭和销毁组件：

```text
Init → Start → 等待退出信号 → Close → Destroy
```

### Gate：客户端进入服务器集群的网关

Gate 负责：

- 管理客户端连接与断开；
- 接收客户端发来的消息；
- 将消息转发到对应的 Node；
- 将 Node 的响应发送给客户端。

Gate 通常不编写房间规则和卡牌规则。它更像游戏服务器集群的入口与转发站。

### WebSocket Server：真正监听网络端口

```go
websocketServer := ws.NewServer()
```

WebSocket Server 负责监听 `3553` 端口、执行 WebSocket 握手、收发字节数据。
Gate 通过 `gate.WithServer(websocketServer)` 使用它。

### Locator 与 Registry：支持多个服务协作

due 的 Gate 按分布式架构设计，因此必须提供 Locator 和 Registry。

- Redis Locator：记录某个玩家当前连接在哪个 Gate。
- etcd Registry：记录集群中有哪些 Gate、Node 等服务，以及它们的地址。

当前只有一个 Gate，它们看起来有些超前；以后运行多个 Gate 和 Node 时，这两个
角色会让服务能够找到彼此。

## 3. 当前入口代码的组装顺序

文件：`server/cmd/server/main.go`

```go
container := due.NewContainer()

websocketServer := ws.NewServer()
locator := redis.NewLocator()
registry := etcd.NewRegistry()

gateway := gate.NewGate(
	gate.WithServer(websocketServer),
	gate.WithLocator(locator),
	gate.WithRegistry(registry),
)

container.Add(gateway)
container.Serve()
```

可以按下面的顺序阅读：

1. 创建负责总生命周期的 Container。
2. 创建 Gate 依赖的三个对象。
3. 使用 Functional Options 将依赖交给 Gate。
4. 将 Gate 加入 Container。
5. 启动 Container，Container 再启动 Gate。

`gate.WithServer`、`gate.WithLocator` 和 `gate.WithRegistry` 是 Go 常见的
Functional Options 模式。它允许构造函数接收可选配置，同时避免参数列表过长。

## 4. 为什么 Gate 不能只传 WebSocket Server

阅读 due Gate 的 `Init` 方法可以看到，它会检查三个必要依赖：

```text
Server   必须存在
Locator  必须存在
Registry 必须存在
```

这说明 due 的 Gate 从设计上就是集群网关，而不是简单的 WebSocket 包装器。

我们没有为了省事编写假的 Locator 或 Registry，而是使用 due 官方支持的
Redis 与 etcd 实现。这样当前学习的代码以后可以自然扩展到多个服务实例。

## 5. 配置如何影响组件

文件：`server/etc/etc.toml`

```toml
[network.ws.server]
addr = ":3553"
path = "/ws"

[locate.redis]
addrs = ["127.0.0.1:6379"]

[registry.etcd]
endpoints = ["127.0.0.1:2379"]
```

代码负责表达“使用哪些组件以及如何组装”；配置负责表达“组件运行在哪里”。

因此：

- 修改 WebSocket 端口，不需要改 Go 代码；
- 修改 Redis 或 etcd 地址，不需要重新设计 Gate；
- 不同环境可以使用不同配置运行同一份程序。

## 6. 使用 Docker Compose 启动本地依赖

文件：`server/compose.yaml`

在 `server/` 目录执行：

```powershell
docker compose up -d
docker compose ps
```

其中 `-d` 表示在后台运行。当前 Compose 只启动本课必需的两个服务：

```text
Redis  → 6379
etcd   → 2379
```

停止它们：

```powershell
docker compose down
```

## 7. 启动与验证

先启动依赖，再启动 Gate：

```powershell
docker compose up -d
go run ./cmd/server
```

成功启动后，due 会输出类似信息：

```text
Gate
Name: exploding-minions-gate
Server: [ws] 0.0.0.0:3553
Locator: redis
Registry: etcd
```

本课已经完成以下验证：

```text
2379  etcd             正在监听
6379  Redis            正在监听
3553  WebSocket Gate   正在监听
```

并实际连接 `ws://127.0.0.1:3553/ws`，WebSocket 状态成功从 `Open` 变为
`Closed`。

## 8. 为什么现在仍然不能处理游戏消息

当前 Gate 可以接受连接，但还没有 Node：

```text
客户端 → Gate → 没有可以处理消息的 Node
```

Gate 负责连接与转发，不负责爆炸小黄人的游戏逻辑。下一课需要加入 Node，并注册
第一条消息路由，才能形成完整的请求与响应：

```text
Go 测试客户端 → WebSocket Gate → Node 路由处理器 → 响应
```

## 9. 本课需要记住的四句话

1. Container 管理组件的生命周期。
2. Gate 管理客户端连接，并在客户端与 Node 之间转发消息。
3. WebSocket Server 才是真正监听网络端口的对象。
4. Locator 找玩家连接，Registry 找服务实例。

下一课：加入第一个 Node 与消息路由，并使用 Go 测试客户端完成首次请求响应。
