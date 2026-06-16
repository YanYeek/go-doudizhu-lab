// Package session 处理连接的生命周期事件（上线、下线）。
//
// 它和 internal/greet 是对称的两半：
//   - greet 处理「客户端主动发来的消息」——路由（Route）；
//   - session 处理「框架自动发来的连接事件」——事件（Event）。
//
// 玩家系统会从这里长出来：之后的登录绑定 uid、掉线清理房间都落在这一层。
package session

import (
	"fmt"

	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"
	"github.com/dobyte/due/v2/session"

	"github.com/YanYeek/go-doudizhu-lab/server/internal/route"
)

// Register 把玩家身份相关的处理器注册到 Node 上：
//   - Login 路由：客户端报上 uid，把连接绑定到该玩家；
//   - connect/disconnect 事件：连接的建立与断开。
//
// 注册事件处理器后，Gate 转发来的 connect/disconnect 就有 Node 接收，
// 不再出现 "not found event" 警告（见 docs/learning/08-route-vs-event.html）。
func Register(proxy *node.Proxy) {
	proxy.Router().AddRouteHandler(route.Login, onLogin)
	proxy.AddEventHandler(cluster.Connect, onConnect)
	proxy.AddEventHandler(cluster.Reconnect, onReconnect)
	proxy.AddEventHandler(cluster.Disconnect, onDisconnect)
}

// onLogin 处理登录：把当前连接绑定到客户端报上来的 uid。
// 绑定（BindGate）之后，网关就在「连接」和「玩家」之间建立了对应关系——
// 这条连接后续的断开事件、服务端推送都能落到具体玩家身上。
func onLogin(ctx node.Context) {
	var req route.LoginReq
	if err := ctx.Parse(&req); err != nil {
		log.Warnf("解析登录请求失败: %v", err)
		return
	}

	if err := ctx.BindGate(req.UserID); err != nil {
		log.Warnf("绑定网关失败 uid=%d: %v", req.UserID, err)
		_ = ctx.Response(&route.LoginRes{Message: "登录失败"})
		return
	}

	log.Infof("玩家登录 uid=%d cid=%d", req.UserID, ctx.CID())
	_ = ctx.Response(&route.LoginRes{Message: fmt.Sprintf("玩家 %d 登录成功", req.UserID)})

	// 绑定成功后，服务端主动按 uid 推一条通知。
	// 这和上面的 Response 本质不同：Response 是「回复客户端的请求」，
	// Push 是「服务端主动找到某个玩家发消息」——而能按 uid 找到连接，
	// 正是因为前面 BindGate 建立了 uid ↔ 连接 的对应。房间广播就是它的放大版。
	pushWelcome(ctx, req.UserID)
}

// pushWelcome 按 uid 主动给玩家推一条欢迎通知。
func pushWelcome(ctx node.Context, uid int64) {
	err := ctx.Proxy().Push(ctx.Context(), &cluster.PushArgs{
		Kind:    session.User, // 按「用户(uid)」寻址，而不是「连接(cid)」
		Target:  uid,
		Message: &cluster.Message{
			Route: route.Notify,
			Data:  &route.NotifyPush{Text: "服务端主动推送：欢迎上线，已准备好开始游戏"},
		},
	})
	if err != nil {
		log.Warnf("推送通知失败 uid=%d: %v", uid, err)
	}
}

// onConnect 在一条客户端连接建立时触发。
// 现在只记录日志；将来这里会承接登录与玩家上线逻辑。
func onConnect(ctx node.Context) {
	log.Infof("连接建立 cid=%d uid=%d", ctx.CID(), ctx.UID())
}

// onReconnect 在连接「绑定到一个 uid」时触发。
// 名字叫 reconnect（断线重连），但 due 在 BindGate 成功后也会发这个事件——
// 因为「连接获得身份」和「重连恢复身份」在网关看来是同一类状态变化。
// 所以登录绑定后会走到这里。
func onReconnect(ctx node.Context) {
	log.Infof("连接绑定/重连 cid=%d uid=%d", ctx.CID(), ctx.UID())
}

// onDisconnect 在一条客户端连接断开时触发。
// 现在只记录日志；将来这里会清理玩家所在的房间与对局状态。
func onDisconnect(ctx node.Context) {
	log.Infof("连接断开 cid=%d uid=%d", ctx.CID(), ctx.UID())
}
