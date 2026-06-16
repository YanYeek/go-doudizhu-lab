# due 学习 04：路由 vs 事件，与那条 connect 警告

> 日期：2026-06-13
>
> 起因：测试客户端每次连/断，服务端都冒一条 WARN：
>
> ```text
> WARN proxy.go:71 trigger event failed, cid: 1, uid: 0, event: connect, err: not found event
> WARN proxy.go:71 trigger event failed, cid: 1, uid: 0, event: disconnect, err: not found event
> ```
>
> 本课搞清楚它是什么、为什么无害、以及背后的一个核心区分。

## 1. 结论先行：这条 WARN 无害

它是 **WARN 不是 ERROR**。due 源码（`cluster/gate/proxy.go:70`）专门把"没人接这个
事件"的情况降级成警告——因为框架预期你**本来就可能不处理每个事件**。所以它是
提示，不是出错。

## 2. due 的两条平行通道：Route 与 Event

这是本课最该记住的区分：

| | **Route（路由）** | **Event（事件）** |
|---|---|---|
| 谁触发 | 客户端**主动发**一条消息 | 框架在**连接生命周期**自动发 |
| 例子 | `Greet=1` | `connect` / `reconnect` / `disconnect` |
| 在 Node 上怎么接 | `Proxy().Router().AddRouteHandler` | `Proxy().AddEventHandler` |
| 处理器签名 | `func(ctx node.Context)` | `func(ctx node.Context)`（一样） |

两者都最终走到一个 `node.Context`，但来源完全不同：路由是"客户端说了句话"，
事件是"框架替你观察到连接状态变了"。

## 3. 警告的来龙去脉

```text
客户端连接 ──▶ Gate 触发 connect 事件 ──▶ 找一个登记了 connect 处理器的 Node
                                              └─ 没找到 → ErrNotFoundEvent → WARN
```

我们之前只给 Node 注册了 `Greet` 这条**路由**处理器（`internal/greet`），
**没注册任何事件处理器**。所以客户端一连上，Gate 想把 connect 事件交给某个 Node，
dispatcher 在事件表里查不到（`internal/dispatcher/dispatcher.go:79` 返回
`ErrNotFoundEvent`），Gate 就记一条 WARN。断开时同理。

日志里两个字段的含义：

- `cid: 1 / 2`：连接编号（connection id），每个新连接递增。测试客户端跑两次就有两个。
- `uid: 0`：这条连接**还没绑定到任何玩家**（还没登录）。uid 为 0 = 匿名连接。

## 4. 怎么消除它

给 Node 注册连接/断开事件处理器即可——"终于有人接这个事件了"：

```go
gameNode.Proxy().AddEventHandler(cluster.Connect, onConnect)
gameNode.Proxy().AddEventHandler(cluster.Disconnect, onDisconnect)
```

事件常量（`cluster` 包）：`Connect`、`Reconnect`、`Disconnect`。

这一步在 `internal/session` 里实现（见下一课/下一个提交）。注册之后，连接的
建立与断开都有 Node 在处理，WARN 随之消失。

## 5. 为什么这条警告其实是"剧透"

`uid: 0` 和这条 WARN 一起在提醒：**还没有"玩家"概念**。处理 connect/disconnect
事件正是玩家系统的起点——上线时记录"谁来了"，掉线时清理。等之后做登录、用
`ctx.BindNode` 把连接绑定到一个 uid，事件处理器里就能拿到真正的玩家身份，
而不再是 uid 0。

## 6. 本课记住三句话

1. Route 是客户端主动发的消息，Event 是框架自动发的连接生命周期信号。
2. `not found event` 警告 = 有连接事件发生，但没有 Node 登记处理它，无害。
3. 用 `AddEventHandler` 注册 connect/disconnect 处理器即可消除，这也是玩家系统的起点。
