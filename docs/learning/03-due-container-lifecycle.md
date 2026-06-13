# due 学习 01：Container 与服务器生命周期

> 日期：2026-06-13
>
> 使用版本：Go 1.26.4、due v2.5.8
>
> 本课目标：不接入网络和游戏逻辑，只理解 due 如何托管并运行服务器组件。

## 1. 为什么不先通读全部 due 文档

due 不只是一个 WebSocket 库，它还包含 Gate、Node、Mesh、组件生命周期、
服务注册、RPC 和集群通信等概念。没有运行经验时直接通读，很容易记住名词，
却不知道它们分别解决什么问题。

当前采用渐进式学习：

1. 每次只引入一两个新概念。
2. 先运行并观察现象，再解释框架行为。
3. 学到的概念立即用于斗地主服务器。

## 2. 普通 Go 程序为什么会退出

下面的程序打印完成后，`main` 函数结束，进程也随之退出：

```go
func main() {
	fmt.Println("服务器启动")
}
```

服务器需要持续等待客户端连接和消息，同时还需要正确处理启动和关闭流程。
真实服务器通常需要负责：

- 初始化网络、日志和业务组件；
- 启动组件；
- 等待 `Ctrl+C` 等系统退出信号；
- 停止接收新连接；
- 关闭并清理组件。

## 3. Container 是什么

可以先将 `Container` 理解为：**服务器组件的生命周期总管**。

```go
container := due.NewContainer()
```

这行代码只创建了一个空容器。当前容器里没有 Gate、Node 或 WebSocket Server：

```text
Container
└── 暂无组件
```

后续可以使用 `Add` 将组件交给容器管理：

```go
container.Add(component)
```

未来的服务器可能逐步形成：

```text
Container
├── Gate：管理客户端连接和消息转发
└── Node：处理房间与游戏业务逻辑
```

Container 本身不是游戏服务器的全部业务，它负责组织和托管构成服务器的组件。

## 4. Serve 做了什么

```go
container.Serve()
```

可以先用下面的伪代码理解：

```go
func Serve() {
	初始化所有组件()
	启动所有组件()

	等待系统退出信号()

	关闭所有组件()
	销毁所有组件()
}
```

due Container 的主要生命周期顺序为：

```text
Init → Start → 等待退出信号 → Close → Destroy
```

因此我们不需要在 `main` 中手写无限循环，也不需要为每个组件分别处理
`Ctrl+C` 和关闭顺序。`Serve` 会让进程持续运行，收到退出信号后再执行清理。

## 5. 为什么程序运行了，却不能连接 WebSocket

当前入口只有一个空容器：

```go
func main() {
	container := due.NewContainer()
	container.Serve()
}
```

`Serve` 让进程保持运行，但空容器没有网络组件，也没有任何端口在监听，
所以客户端无法建立 WebSocket 连接。

下一课加入 Gate 和 WebSocket Server 后，结构才会变成：

```text
Container
└── Gate
    └── WebSocket Server：监听端口并接收客户端连接
```

需要区分两个概念：

- **进程正在运行**：程序尚未退出。
- **网络服务可以连接**：某个网络组件正在监听端口。

前者不自动代表后者。

## 6. 当前入口代码

文件：`server/cmd/server/main.go`

```go
package main

import "github.com/dobyte/due/v2"

func main() {
	container := due.NewContainer()
	container.Serve()
}
```

`cmd/server/main.go` 只负责组装与启动。以后真正的游戏规则、房间状态和消息处理，
仍然放在 `server/internal/` 下，不把业务逻辑堆进入口文件。

## 7. 配置目录与本次兼容问题

due 默认读取进程工作目录下的 `./etc`，因此新增了最小配置文件
`server/etc/etc.toml`。从 `server/` 目录运行程序时，框架可以正常找到配置。

due v2.5.8 默认依赖的 `sonic v1.14.2` 无法使用 Go 1.26.4 编译。本项目将
该间接依赖提升到 `sonic v1.15.2` 后通过构建。这说明：

- `go.mod` 中的直接依赖可能继续引入许多间接依赖；
- 最新框架版本不代表它锁定的每个间接依赖都支持最新 Go；
- 遇到编译错误时，应先定位依赖关系，再进行最小范围升级。

另外，due v2.5.8 启动画面显示 v2.5.7，这是框架内部展示版本未同步；
项目实际使用的模块版本以 `go.mod` 和 `go list -m` 的结果为准。

## 8. 本课验证

在 `server/` 目录执行：

```powershell
go run ./cmd/server
```

观察到程序持续运行，按 `Ctrl+C` 后退出。随后执行：

```powershell
go test ./...
go vet ./...
go build ./cmd/server
```

均验证通过。

## 9. 本课需要记住的三句话

1. `due.NewContainer()` 创建负责托管服务器组件的容器。
2. `container.Add(component)` 将组件交给容器管理。
3. `container.Serve()` 启动组件、等待退出信号，并负责关闭与清理。

下一课：向 Container 添加 Gate 与 WebSocket Server，让服务器第一次真正监听端口。
