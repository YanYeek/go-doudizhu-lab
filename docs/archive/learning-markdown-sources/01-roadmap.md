# 学习区路线图：从 GOPATH 时代到亲手搭出游戏服务器

> 输入：[00-comfort-zone-quiz.md](00-comfort-zone-quiz.md) 的诊断结论
> 对齐：[框架调研](../go-game-server-frameworks-research.md) 5.3 节的分阶段落地路线图
> 原则：**每一步只引入 1~2 个新知识点，学了立刻写进真实代码**——
> 始终踩在学习区，不滑回舒适区（无聊），也不冲进困难区（劝退）。

## 总览

| 阶段 | 名称 | 补的 Go 知识（对应诊断题） | 产出物 |
|------|------|---------------------------|--------|
| 0 | 工具链复健 | go mod、go run/test/vet、internal（Q7） | `server/go.mod` + 能跑的 hello + 第一个测试 |
| 1 | 纯逻辑：牌与牌型 | 切片（Q1）、map（Q2）、接收者（Q3）、表驱动测试、errors（Q8）、泛型读懂（Q9） | `internal/card`、`internal/pattern` 全测试通过 |
| 2 | 并发复健 + 连接层 | channel（Q4）、goroutine/WaitGroup（Q5）、select（Q6）、context | `internal/ws` 回显服务器 + `proto/` 首批消息 |
| 3 | 房间与状态机 | 单 goroutine 串行模型、time.Ticker、锁 vs channel 的取舍 | `internal/room`、`internal/game` |
| 4 | 前后端联调 | （前端为主，Go 侧巩固） | Cocos 接入，跑通一局 |
| 5 | 进阶 | 断线重连、会话、匹配 | 按需展开 |

阶段 3 起与框架调研 5.3 节完全重合，本文重点展开阶段 0~2 的学习设计。

## 阶段 0：工具链复健（预计 1~2 次会话）

目标：把"go mod 之后的世界"建立起来，肌肉记忆级别。

1. 安装/确认 Go ≥1.22，理解为什么 1.22 是分水岭（循环变量语义，Q5）
2. `server/` 下 `go mod init github.com/YanYeek/go-doudizhu-lab/server`
   —— 理解 module 路径 = import 前缀，GOPATH 退役（Q7）
3. 建 `cmd/server/main.go` 跑通 hello；建 `internal/` 理解私有目录含义（Q7）
4. 写第一个 `_test.go` 并 `go test ./...`，体验"测试是一等公民"

每步配一段「当年 vs 现在」对照说明，沉淀到 `docs/learning/02-modern-go-toolchain.md`。

## 阶段 1：纯逻辑——牌与牌型（预计 3~5 次会话）

目标：零网络依赖，把斗地主规则写成可测试的纯函数，同时焊牢第 1 层基础。

| 开发任务 | 顺带焊牢的知识点 |
|----------|------------------|
| 定义 `Card`/`Suit`/`Rank` 类型与 54 张牌生成 | 自定义类型、iota、Stringer 接口 |
| 洗牌（Fisher-Yates）与发牌（17+17+17+3） | 切片底层数组与共享陷阱（Q1）、math/rand/v2 |
| 手牌排序与计数 | sort.Slice、map 零值与 comma-ok（Q2） |
| `Hand` 类型的方法集（加牌/出牌/判空） | 值接收者 vs 指针接收者（Q3）——在真实场景里踩一次 |
| 牌型识别（单/对/三带/顺子/炸弹/王炸…） | 表驱动测试（每种牌型一张用例表） |
| 牌型比较与非法出牌错误 | sentinel error + `%w` + `errors.Is`（Q8） |
| （可选重构）通用计数/分组小工具 | 读懂泛型签名即可，不强求会写（Q9） |

阶段完成标准：`go test ./...` 全绿，且你能对着任意一个测试用例讲出"为什么"。

## 阶段 2：并发复健 + 连接层（预计 3~4 次会话）

目标：先用 3 个微练习恢复并发手感（每个 ≤30 行），再上 WebSocket。

微练习（写在 `docs/learning/` 配套的 playground 里，不进正式代码）：
1. 无缓冲 vs 有缓冲 channel：亲手复现 Q4 的死锁再修好它
2. WaitGroup + 循环变量：复现 Q5，在本机 Go 版本下验证语义
3. select + 超时：复现 Q6，加上 `context.WithTimeout` 版本

然后进入正式开发：
- `proto/` 定义首批消息（心跳、回显、错误码），C2S_/S2C_ 前缀
- `internal/ws`：基于 gorilla/websocket（允许借轮子），实现连接管理、
  读写分离 goroutine、消息编解码，先做到回显
- 对照阅读：Zinx 的 Connection/MsgHandler——这次你能看懂当年没看懂的设计了

## 节奏与复盘约定

- 每次会话结束：当次学到的知识点用 1~3 句话写进对应学习笔记，
  按 Conventional Commits 提交（scope 用 docs / server）
- 每完成一个阶段：与对照读物做一次差异复盘（见框架调研 5.3）
- 感觉无聊 = 滑回舒适区，加快；感觉想放弃 = 误入困难区，拆小步——
  随时说出来，路线图按体感调整
