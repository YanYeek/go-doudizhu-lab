package main

import (
	"github.com/dobyte/due/locate/redis/v2"
	"github.com/dobyte/due/network/ws/v2"
	"github.com/dobyte/due/registry/etcd/v2"
	"github.com/dobyte/due/v2"
	"github.com/dobyte/due/v2/cluster/gate"
	"github.com/dobyte/due/v2/cluster/node"

	"github.com/YanYeek/go-doudizhu-lab/server/internal/greet"
)

func main() {
	container := due.NewContainer()

	// Gate：客户端进入服务器集群的网关，负责连接管理与消息转发，不含游戏逻辑。
	gateway := gate.NewGate(
		gate.WithServer(ws.NewServer()),
		gate.WithLocator(redis.NewLocator()),
		gate.WithRegistry(etcd.NewRegistry()),
	)

	// Node：承载游戏业务逻辑的节点。Gate 通过 Registry 发现它，把消息转发过来。
	// 这里注册了项目的第一条路由（问候），用来打通整条通信链路。
	// Node 自带 Locator/Registry，与 Gate 各持一份，避免共享对象被重复关闭。
	gameNode := node.NewNode(
		node.WithLocator(redis.NewLocator()),
		node.WithRegistry(etcd.NewRegistry()),
	)
	greet.Register(gameNode.Proxy().Router())

	// 单进程内同时运行 Gate 与 Node（单体模式），适合学习阶段。
	// 将来需要水平扩容时，可以把 Node 拆成独立进程，代码几乎不用改。
	container.Add(gateway)
	container.Add(gameNode)
	container.Serve()
}
