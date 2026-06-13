// Package session 处理连接的生命周期事件（上线、下线）。
//
// 它和 internal/greet 是对称的两半：
//   - greet 处理「客户端主动发来的消息」——路由（Route）；
//   - session 处理「框架自动发来的连接事件」——事件（Event）。
//
// 玩家系统会从这里长出来：之后的登录绑定 uid、掉线清理房间都落在这一层。
package session

import (
	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"
)

// Register 把连接生命周期事件处理器注册到 Node 上。
// 注册后，Gate 转发来的 connect/disconnect 事件就有 Node 接收，
// 不再出现 "not found event" 警告（见 docs/learning/07-route-vs-event.md）。
func Register(proxy *node.Proxy) {
	proxy.AddEventHandler(cluster.Connect, onConnect)
	proxy.AddEventHandler(cluster.Disconnect, onDisconnect)
}

// onConnect 在一条客户端连接建立时触发。
// 现在只记录日志；将来这里会承接登录与玩家上线逻辑。
func onConnect(ctx node.Context) {
	log.Infof("连接建立 cid=%d uid=%d", ctx.CID(), ctx.UID())
}

// onDisconnect 在一条客户端连接断开时触发。
// 现在只记录日志；将来这里会清理玩家所在的房间与对局状态。
func onDisconnect(ctx node.Context) {
	log.Infof("连接断开 cid=%d uid=%d", ctx.CID(), ctx.UID())
}
