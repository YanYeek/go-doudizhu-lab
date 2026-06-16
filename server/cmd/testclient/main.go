// 一次性的 Go 测试客户端：连接 Gate，先登录、再发一条问候，打印响应后退出。
//
// 用途：在不启动 Cocos 的情况下，用最快的方式验证后端通信链路
// 「客户端 → Gate → Node → 响应」是否打通（开发学习工作流里的「Go 集成客户端」层）。
//
// 运行前先起依赖与服务端：
//
//	python scripts/python/dev.py up         # 启动 Redis、etcd 并在前台运行 Gate+Node
//	python scripts/python/dev.py testclient # 另开一个终端运行本客户端
package main

import (
	"time"

	"github.com/dobyte/due/network/ws/v2"
	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/cluster/client"
	"github.com/dobyte/due/v2/log"

	"github.com/YanYeek/go-doudizhu-lab/server/internal/route"
)

func main() {
	// 用 due 的 client 组件包一个 WebSocket 客户端，它会复用与服务端一致的
	// 打包格式和 json 编解码器——这正是要验证的协议本身。
	cli := client.NewClient(client.WithClient(ws.NewClient()))

	// 每条响应路由各用一个 channel 把结果交回主流程。
	loginCh := make(chan string, 1)
	greetCh := make(chan string, 1)
	notifyCh := make(chan string, 1)
	cli.Proxy().AddRouteHandler(route.Login, func(ctx *client.Context) {
		var res route.LoginRes
		if err := ctx.Parse(&res); err != nil {
			log.Errorf("解析登录响应失败: %v", err)
			return
		}
		loginCh <- res.Message
	})
	// Notify 没有对应的请求，是服务端登录后主动推过来的。
	cli.Proxy().AddRouteHandler(route.Notify, func(ctx *client.Context) {
		var p route.NotifyPush
		if err := ctx.Parse(&p); err != nil {
			log.Errorf("解析通知失败: %v", err)
			return
		}
		notifyCh <- p.Text
	})
	cli.Proxy().AddRouteHandler(route.Greet, func(ctx *client.Context) {
		var res route.GreetRes
		if err := ctx.Parse(&res); err != nil {
			log.Errorf("解析问候响应失败: %v", err)
			return
		}
		greetCh <- res.Message
	})

	cli.Init()
	cli.Start()
	defer cli.Destroy()

	conn, err := cli.Proxy().Dial()
	if err != nil {
		log.Fatalf("连接 Gate 失败（确认 dev.py up 已启动 Gate）: %v", err)
	}
	defer conn.Close()

	// 1) 先登录：报上 uid，让服务端把这条连接绑定到玩家 1001。
	push(conn, 1, route.Login, &route.LoginReq{UserID: 1001})
	log.Infof("登录响应: %s", await(loginCh))

	// 2) 登录后再发一条普通请求，验证已登录的连接照常工作。
	push(conn, 2, route.Greet, &route.GreetReq{Name: "测试客户端"})
	log.Infof("问候响应: %s", await(greetCh))

	// 3) 收取服务端在登录后主动推来的通知（没有对应请求）。
	log.Infof("服务端主动推送: %s", await(notifyCh))
}

// push 发送一条消息，失败直接终止（测试客户端无需重试）。
func push(conn *client.Conn, seq, r int32, data any) {
	if err := conn.Push(&cluster.Message{Seq: seq, Route: r, Data: data}); err != nil {
		log.Fatalf("发送消息失败 route=%d: %v", r, err)
	}
}

// await 等待一条响应，超时即终止。
func await(ch <-chan string) string {
	select {
	case msg := <-ch:
		return msg
	case <-time.After(5 * time.Second):
		log.Fatal("等待响应超时：检查 Gate 与 Node 是否都已启动、路由号是否一致")
		return ""
	}
}
