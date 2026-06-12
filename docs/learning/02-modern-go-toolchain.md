# 阶段 0 笔记：现代 Go 工具链（当年 vs 现在）

> 日期：2026-06-12 ｜ 本机 Go 1.25.0
> 对应路线图阶段 0，补诊断题 Q7；Q5 的版本分水岭也在这里交代。

## 1. GOPATH 退役了，但没删除

| | 当年（GOPATH 时代） | 现在（go mod 时代） |
|---|---|---|
| 代码必须放哪 | `$GOPATH/src/github.com/你/项目` | **任意目录**，本仓库在 `d:\work_space\...` |
| 项目身份证 | 目录路径本身 | `go.mod` 里的 `module` 行 |
| 依赖下载到哪 | `$GOPATH/src`（和你的代码混在一起） | `$GOPATH/pkg/mod` 只读缓存，版本号隔离 |
| 依赖版本 | 没有版本概念，`go get` 拉最新 | `go.mod` 锁定版本 + `go.sum` 校验哈希 |

本机验证过：`GOPATH=C:\Users\YanYeek\go` 仍然存在，但它现在只是依赖缓存仓库。

## 2. import 是怎么找到代码的（Q7 答案落地）

```
go.mod:  module github.com/YanYeek/go-doudizhu-lab/server
                     │
import "github.com/YanYeek/go-doudizhu-lab/server/internal/greeting"
         └────────── 前缀匹配 module 行 ──────────┘└── 剩余部分 ──┘
                                                    = 模块根下的相对目录
                                                      server/internal/greeting/
```

- 前缀命中本模块 → 直接在本地目录树里找，**不联网、不看 GOPATH**
- 前缀是第三方（如 `github.com/gorilla/websocket`）→ 去 `go.mod` 查版本，
  从模块缓存读，没有就下载

## 3. internal 目录：编译器强制的私有

`server/internal/...` 下的包，**只有** `server/` 模块自己能 import。
将来如果有人想 `import ".../server/internal/game"` 直接偷用我们的游戏逻辑，
编译器会拒绝。这就是为什么 CLAUDE.md 约定"业务逻辑都放 internal/"——
对外只暴露我们想暴露的东西。当年这个规则刚出现（Go 1.4），现在已是标准实践。

## 4. 测试是一等公民

- `xxx_test.go` 和被测代码同目录同包，`go test ./...` 一条命令全跑
- 第一个表驱动测试见 `server/internal/greeting/greeting_test.go`：
  用例是一张表 + `t.Run` 子测试，失败时能精确报出是哪个用例
- 以后所有牌型判定测试都是这个形状（CLAUDE.md 已约定）

## 5. Go 1.22 分水岭（Q5 伏笔）

Go 1.22（2024 年初）修改了 for 循环语义：**每轮迭代的循环变量是新变量**。
诊断题 Q5 那段代码，老版本大概率打 `3 3 3`，1.22+ 打 `0 1 2`（乱序）。
本机是 1.25，按新语义走。阶段 2 微练习时会亲手验证。

## 6. 本阶段动手记录

```bash
cd server
go mod init github.com/YanYeek/go-doudizhu-lab/server  # 生成 go.mod
go run ./cmd/server    # 斗地主服务器 已启动
go test ./... -v       # 2 个子测试 PASS
go vet ./...           # 无警告
```

产出文件：

- `server/go.mod` —— 模块身份证（3 行）
- `server/cmd/server/main.go` —— 入口，跨包 import 演示
- `server/internal/greeting/` —— 复健包（含表驱动测试），**阶段 1 删除**

## 7. 自检问题（能答上来才算过关）

1. 把整个仓库移动到 E 盘，代码还能编译吗？为什么？
2. `internal/greeting` 改名成 `pkg/greeting`，外部模块能 import 它吗？
3. `go test ./...` 里的 `./...` 是什么意思？
