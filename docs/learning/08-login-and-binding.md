# due 学习 05：登录，与连接到玩家的绑定

> 日期：2026-06-14
>
> 本课目标：加一条登录路由，把「一条连接」绑定到「一个玩家 uid」，
> 让服务端从"匿名连接"升级到"认得玩家"。这是玩家系统的第一块地基。

## 1. 这一步加了什么

| 文件 | 改动 |
|------|------|
| `internal/route` | 新增 `Login=2` 路由 + `LoginReq{UserID}` / `LoginRes` |
| `internal/session` | 新增 `onLogin`（绑定 uid）+ `onReconnect` 事件处理器 |
| `cmd/testclient` | 先登录、再问候，演示真实客户端的先后顺序 |

没有新增模块——登录天然属于 `session`（连接/玩家这一层），和上一课的
connect/disconnect 事件住在一起。

## 2. 核心动作：BindGate

登录处理器只做一件关键的事：

```go
func onLogin(ctx node.Context) {
    var req route.LoginReq
    ctx.Parse(&req)
    ctx.BindGate(req.UserID)   // ← 把当前连接绑定到这个 uid
    ctx.Response(&route.LoginRes{Message: "登录成功"})
}
```

`BindGate` 让**网关**在"连接（cid）"和"玩家（uid）"之间建立对应关系。绑定前，
网关只知道"有一条连接"；绑定后，它知道"这条连接是玩家 1001"。

> 现在没有真正的鉴权——客户端直接报上 uid。学习阶段先打通机制，鉴权是后话。

## 3. 验证：让上一课的事件日志替我们说话

绑定有没有生效，不用猜——上一课注册的 `session` 事件处理器会把整条连接的
一生打出来。跑一次测试客户端，服务端日志是：

```text
连接建立      cid=1 uid=0       ← 刚连上，匿名
玩家登录      uid=1001 cid=1     ← Login 路由，BindGate
连接绑定/重连  cid=1 uid=1001     ← 见第 4 节
连接断开      cid=1 uid=1001     ← 关键：断开时 uid 已是 1001
```

最后一行是铁证：**绑定前断开会是 `uid=0`，绑定后断开变成 `uid=1001`**。
连接终于"带着身份"离开，掉线清理时就能知道是哪个玩家走了。

## 4. 一个意外的小知识：BindGate 会触发 reconnect 事件

第一次跑完，`not found event` 警告没有完全消失，还剩一条：

```text
trigger event failed, uid: 1001, event: reconnect, err: not found event
```

查 due 源码（`cluster/gate/proxy.go`）发现：`BindGate` 成功后，网关会**额外触发
一个 `Reconnect` 事件**。原因是——在网关眼里，"连接获得身份"和"断线重连恢复
身份"是同一类状态变化，都用 `Reconnect` 表示。

我们没注册 `Reconnect` 处理器，于是又是一条"没人接这个事件"的警告（正是上一课
讲过的机制）。补上 `onReconnect` 后，警告归零：

```go
proxy.AddEventHandler(cluster.Reconnect, onReconnect)
```

这也印证了上一课的规律：**每出现一类新事件，就要么处理它、要么容忍那条 WARN。**

## 5. 现在的连接生命周期全景

```text
connect    →  uid=0    （匿名连接建立）
  └─ 客户端发 Login 路由
reconnect  →  uid=1001 （BindGate 绑定身份后触发）
  └─ 客户端正常收发消息（已登录）
disconnect →  uid=1001 （带着身份断开）
```

## 6. 本课记住三句话

1. 登录 = 一条路由里调用 `BindGate(uid)`，把连接和玩家对应起来。
2. 绑定是否生效，看断开事件的 uid 从 0 变成真实玩家号即可确认。
3. `BindGate` 会顺带触发 `Reconnect` 事件——记得也注册它，否则又多一条 WARN。

## 7. 下一小步（保持小）

现在网关知道"这条连接是哪个玩家"，但玩家之间还彼此孤立。下一步可以只做一件事：
**给同一个玩家做一次服务端主动推送**（用 uid 找到连接、Push 一条消息），
体会"绑定之后服务端能反过来找到玩家"——这是将来房间广播、通知出牌的基础。
