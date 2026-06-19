// Package game 承载斗地主的对局逻辑路由（出牌、校验、对局流程……）。
//
// 它把纯逻辑（card、pattern）接到 due 的网络层：客户端发来意图，
// game 在服务端做权威判定，再把结果回给客户端。
// 目前只有第一条：出牌校验（判牌型是否合法）。后续再加发牌、轮次、比大小、判胜负。
package game

import (
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"

	"github.com/YanYeek/go-doudizhu-lab/server/internal/pattern"
	"github.com/YanYeek/go-doudizhu-lab/server/internal/route"
)

// Register 把对局相关路由注册到 Node 的路由器上。在 main.go 组装 Node 时调用。
func Register(router *node.Router) {
	router.AddRouteHandler(route.PlayCheck, handlePlayCheck)
}

// handlePlayCheck 校验客户端提交的一手牌：用 pattern 引擎识别牌型，回判定结果。
// 牌型判定全在服务端做——这是"服务端权威"的体现，客户端提交的只是"我想出这几张"。
func handlePlayCheck(ctx node.Context) {
	var req route.PlayCheckReq
	if err := ctx.Parse(&req); err != nil {
		log.Warnf("解析出牌校验请求失败: %v", err)
		return
	}

	// pattern.Parse 返回 (Play, ok)；ok=false 表示这几张牌凑不成合法牌型。
	play, ok := pattern.Parse(req.Cards)

	res := route.PlayCheckRes{Valid: ok, Kind: "非法"}
	if ok {
		res.Kind = play.Kind.String() // 牌型枚举的中文名，如"对子""炸弹"
	}
	if err := ctx.Response(&res); err != nil {
		log.Warnf("响应出牌校验失败: %v", err)
	}
}
