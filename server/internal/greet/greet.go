// Package greet 是 Node 上的第一个业务处理器：响应客户端的问候请求。
//
// 它存在的意义不是游戏逻辑本身，而是把 due 的「路由 → 处理器 → 响应」串起来，
// 验证后端通信链路通畅。真正的斗地主逻辑（发牌、牌型、对局）后续在
// internal/card、internal/pattern、internal/game 里实现。
package greet

import (
	"fmt"

	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"

	"github.com/YanYeek/go-doudizhu-lab/server/internal/route"
)

// Register 把问候路由注册到 Node 的路由器上。
// 在 main.go 组装 Node 时调用。
func Register(router *node.Router) {
	router.AddRouteHandler(route.Greet, handle)
}

// handle 处理一次问候请求：解析名字，回一句欢迎语。
// ctx 封装了这条消息的全部上下文（连接、序列号、路由号），
// Parse 把字节解析成请求体，Response 把响应体编码后沿原路发回客户端。
func handle(ctx node.Context) {
	var req route.GreetReq
	if err := ctx.Parse(&req); err != nil {
		log.Warnf("解析问候请求失败: %v", err)
		return
	}

	name := req.Name
	if name == "" {
		name = "玩家"
	}

	res := &route.GreetRes{Message: fmt.Sprintf("你好，%s，欢迎来到斗地主服务器", name)}
	if err := ctx.Response(res); err != nil {
		log.Warnf("响应问候失败: %v", err)
	}
}
