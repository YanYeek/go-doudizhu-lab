// 一次性的 Go 测试客户端：连接 Gate，发送一条 Greet 请求，打印服务端响应后退出。
//
// 用途：在不启动 Cocos 的情况下，用最快的方式验证后端通信链路
// 「客户端 → Gate → Node → 响应」是否打通（开发学习工作流里的「Go 集成客户端」层）。
//
// 运行前先起依赖与服务端：
//
//	python scripts/python/dev.py up   # 启动 Redis、etcd 并在前台运行 Gate+Node
//	go run ./cmd/testclient           # 另开一个终端运行本客户端
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

	// 注册响应处理器：服务端的 Response 会沿原路由号发回，这里收下并转给主流程。
	done := make(chan string, 1)
	cli.Proxy().AddRouteHandler(route.Greet, func(ctx *client.Context) {
		var res route.GreetRes
		if err := ctx.Parse(&res); err != nil {
			log.Errorf("解析响应失败: %v", err)
			return
		}
		done <- res.Message
	})

	cli.Init()
	cli.Start()
	defer cli.Destroy()

	conn, err := cli.Proxy().Dial()
	if err != nil {
		log.Fatalf("连接 Gate 失败（确认 dev.py up 已启动 Gate）: %v", err)
	}
	defer conn.Close()

	if err = conn.Push(&cluster.Message{
		Seq:   1,
		Route: route.Greet,
		Data:  &route.GreetReq{Name: "测试客户端"},
	}); err != nil {
		log.Fatalf("发送请求失败: %v", err)
	}
	log.Infof("已发送 Greet 请求，等待响应……")

	select {
	case msg := <-done:
		log.Infof("链路打通，收到服务端响应: %s", msg)
	case <-time.After(5 * time.Second):
		log.Fatal("等待响应超时：检查 Gate 与 Node 是否都已启动、路由号是否一致")
	}
}
