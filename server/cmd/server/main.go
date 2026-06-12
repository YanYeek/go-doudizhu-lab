package main

import (
	"github.com/dobyte/due/locate/redis/v2"
	"github.com/dobyte/due/network/ws/v2"
	"github.com/dobyte/due/registry/etcd/v2"
	"github.com/dobyte/due/v2"
	"github.com/dobyte/due/v2/cluster/gate"
)

func main() {
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
}
