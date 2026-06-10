# Go 开源游戏服务器框架调研（2026-06）

> 调研日期：2026-06-11。所有 star 数、提交日期、版本号均为当日从 GitHub 实时抓取，非二手资料。
> 调研动机：为本项目（房间式斗地主，WebSocket 通信，学习目的）选择/参考游戏服务器框架，
> 并核实「Zinx 已停更」等流传说法是否属实。

## 一、调研标准

按本项目的诉求，框架按以下维度评估，**前两条是硬性门槛**：

1. **持续更新**：2025–2026 年仍有实质性提交（性能优化、bug 修复、新特性），而不是只改 README
2. **学习价值**：代码可读、文档/教程完善、能学到现代 Go 游戏服务器的架构范式
3. **房间类卡牌游戏适配度**：是否适合 2–5 人小房间、回合制、命令式状态同步的场景
4. **社区与生态**：star 数、issue 响应、示例项目、客户端 SDK

## 二、结论速览（TL;DR）

| 结论 | 说明 |
|------|------|
| **「Zinx 停更」是误传** | v1.2.8 发布于 2026-05-04，2026 年 5–6 月 master 仍有大量实质提交（goroutine 开销优化、C10K 压测、连接 panic 修复）。它没死，只是定位一直是「轻量 TCP 网络层框架」，不含房间/匹配等游戏业务层 |
| **现代范式首选：due** | 2026-06-02 发布 v2.5.8，Kratos 式模块化架构，中文文档，最能代表当前 Go 游戏服务器工程化范式 |
| **房间类游戏最对口：Cherry** | Actor 模型天然匹配「单房间单执行流」，2026-05 仍在更新，体量适中适合读源码 |
| **工业级参照物：Pitaya / Nakama** | Pitaya（Wildlife Studios 出品，2026-06-09 发版）学分布式集群；Nakama（12.7k★）学 BaaS 产品设计，但黑盒多、学习价值偏低 |
| **本项目建议** | 继续按现有规划**手写最小核心**（gorilla/websocket + 单房间单 goroutine + 状态机），把 due 和 Cherry 当「参考读物」对照学习，而不是直接引入重框架 |

## 三、总览对比大表

### 3.1 活跃框架（2025–2026 仍在实质更新）✅

| 框架 | 仓库 | Stars | 最近实质提交 | 最新版本（日期） | License | 一句话定位 |
|------|------|-------|--------------|------------------|---------|------------|
| **Nakama** | [heroiclabs/nakama](https://github.com/heroiclabs/nakama) | 12.7k | 持续高频 | v3.39.0（2026-05-20） | Apache-2.0 | 全功能游戏 BaaS（认证/匹配/排行/聊天开箱即用） |
| **Zinx** | [aceld/zinx](https://github.com/aceld/zinx) | 7.7k | 2026-06-06 | v1.2.8（2026-05-04） | MIT | 轻量 TCP/WS 并发服务器框架（教学向） |
| **Pitaya** | [topfreegames/pitaya](https://github.com/topfreegames/pitaya) | 2.8k | 2026-06 | v2.11.23（2026-06-09） | MIT | 生产验证的分布式集群游戏框架 |
| **Origin** | [duanhf2012/origin](https://github.com/duanhf2012/origin) | 1.7k | 2025 下半年 | v2.2.0（2025-12-15） | Apache-2.0 | Node/Service/Module 三层抽象的分布式引擎 |
| **due** | [dobyte/due](https://github.com/dobyte/due) | 906 | 2026-06-02 | v2.5.8（2026-06-02） | MIT | Kratos 风格模块化分布式游戏框架 |
| **Cherry** | [cherry-game/cherry](https://github.com/cherry-game/cherry) | 783 | 2026-05-11 | v1.5.1（2026-05-11） | MIT | Actor 模型游戏框架，兼容 pomelo 协议 |
| tgf | [thkhxm/tgf](https://github.com/thkhxm/tgf) | 130 | 2025-01-11 | v1.x | Apache-2.0 | 基于 rpcx 的分布式框架（小团队向） |
| gserver | [fish-tennis/gserver](https://github.com/fish-tennis/gserver) | 50 | 中等活跃 | 无正式 release | MIT | 框架 + 完整示例（背包/任务/公会都有实现） |

### 3.2 已停更/半停滞（不建议作为新项目起点）❌

| 框架 | 仓库 | Stars | 最后实质更新 | 状态说明 |
|------|------|-------|--------------|----------|
| Leaf | [name5566/leaf](https://github.com/name5566/leaf) | 5.5k | 2018（最后 release 2016-12） | 经典老框架，2022 年后连 README 都不改了 |
| cellnet | [davyxu/cellnet](https://github.com/davyxu/cellnet) | 4.1k | v4 发布于 2018-05 | 网络库而非游戏框架，多年无版本 |
| Nano | [lonng/nano](https://github.com/lonng/nano) | 3.2k | 2025-03（最后 release v0.5.1 是 2023-05） | 半停滞，偶有社区 PR 合并，作者基本不投入 |
| GoWorld | [xiaonanln/goworld](https://github.com/xiaonanln/goworld) | 2.7k | 2022-08 | MMO 向引擎，已停更，且对卡牌游戏过重 |
| mqant | [liangdas/mqant](https://github.com/liangdas/mqant) | 2.5k | 2023（2024-09 仅版本号变更） | 半停更 |
| Minotaur | [kercylan98/minotaur](https://github.com/kercylan98/minotaur) | 207 | **2026-03-04 已归档** | 作者迁移到 [vivid](https://github.com/kercylan98/vivid)（纯 Actor 库，v0.1.4，12★，太早期） |

> 注：GitHub 上搜索「2023 年后创建、stars>100 的 Go 游戏服务器框架」结果为 0——
> 近三年**没有出现新的高人气框架**，社区力量集中在 due / Cherry / Pitaya / Nakama 这几个存量项目的持续演进上。

## 四、第一梯队详评（按本项目适配度排序）

### 4.1 due —— 现代工程化范式的代表（推荐精读）

- **仓库**：https://github.com/dobyte/due （注意：网上部分文章写的 `devagame/due` 是错误地址）
- **数据**：906★ / 162 fork / 1117 commits / 52 个 release / MIT
- **活跃度**：v2.5.8 发布于 2026-06-02，提交内容是 registry 优化、etcd 配置中心加认证参数等实打实的功能演进
- **架构**：借鉴 Kratos 的模块化设计，三种服务形态——Gate（网关）、Node（有状态逻辑）、Mesh（无状态微服务）
- **能力清单**：
  - 网络：TCP / KCP / WebSocket 一键切换
  - RPC：gRPC、rpcx；序列化：JSON / Protobuf / MessagePack
  - 服务发现：etcd / Consul / Nacos；事件总线：Redis / NATS / Kafka / RabbitMQ
  - 内置 Actor 模型、分布式锁、due-cli 脚手架
- **学习价值**：⭐⭐⭐⭐⭐ 中文文档完善，代码量适中可通读，学到的是「当前国内 Go 游戏服务器的标准工程范式」，与从 Zinx 学到的网络层知识正好衔接
- **对本项目**：作为主要参考对象。它的 Gate/Node 分层、消息路由、房间组织方式都值得对照自己手写的实现

### 4.2 Cherry —— Actor 模型与房间服的天作之合（推荐精读）

- **仓库**：https://github.com/cherry-game/cherry
- **数据**：783★ / 136 fork / 756 commits / MIT
- **活跃度**：v1.5.1 发布于 2026-05-11，近期提交含 NATS 连接池重构、protobuf 导入路径修正等
- **架构**：每个 Actor 一条独立 goroutine 串行处理消息——「一个房间 = 一个 Actor」，天然免锁
- **能力清单**：TCP / WebSocket / HTTP 连接器；NATS 做服务发现与 RPC；兼容 pomelo 协议（Cocos 客户端有现成对接方案）；组件库覆盖 etcd、Gin、GORM、MongoDB
- **学习价值**：⭐⭐⭐⭐⭐ 官方提供单机聊天室和多节点分布式两个完整示例，体量比 due 更小，是理解 Actor 范式最合适的入口
- **对本项目**：斗地主房间逻辑（叫地主→出牌→结算的串行状态流转）与 Actor 模型完全同构，值得对照设计 `internal/room` 和 `internal/game`

### 4.3 Zinx —— 没有停更，但要认清它的定位（适合保留为网络层教材）

- **仓库**：https://github.com/aceld/zinx
- **数据**：7.7k★ / 1.3k fork / MIT
- **活跃度核实**（这是本次调研的关键发现）：
  - v1.2.8 发布于 **2026-05-04**，v1.2.7 发布于 2025-06-30——保持着约一年一个版本的节奏
  - master 分支 2026-05~06 的提交全是硬核内容：goroutine 开销削减、写协程问题修复、连接 panic 修复、C10K 压测示例、批量发送优化
- **学习价值**：⭐⭐⭐⭐ 作者 aceld（刘丹冰）的配套教程（语雀文档、B 站视频、出版书籍）仍是中文圈最完整的「从零手写游戏服务器网络层」教材
- **真正的短板**：它是**网络层框架**（连接管理、消息封包、worker 池），不提供房间、匹配、状态同步、服务发现等游戏业务层能力——这些恰恰是房间类游戏的核心。「跟不上」的不是它的维护，而是它的抽象层次
- **对本项目**：之前学的 Zinx 知识没有过时，连接管理和消息分发的思想直接复用；但房间层以上要看 due/Cherry

### 4.4 Pitaya —— 生产验证最充分的分布式方案

- **仓库**：https://github.com/topfreegames/pitaya
- **数据**：2.8k★ / 533 fork / 80 个 release / MIT
- **活跃度**：v2.11.23 发布于 **2026-06-09**（调研前两天），高频小版本迭代
- **背景**：TFG Co（Wildlife Studios，亿级下载手游厂商）出品并在生产环境使用；思想源自 pomelo/nano
- **能力清单**：etcd 服务发现 + NATS RPC 的集群方案；官方 C SDK 衍生出 Unity / iOS / Android 客户端库；自带 pitaya-cli 调试 REPL
- **学习价值**：⭐⭐⭐⭐ 学「真实商业项目如何做分布式游戏集群」的最佳样本；但文档英文为主，且对 2–5 人小房间场景偏重
- **对本项目**：现阶段不必引入，留作日后学习分布式扩容时的对照

### 4.5 Nakama —— 工业级 BaaS，适合看产品不适合学架构

- **仓库**：https://github.com/heroiclabs/nakama
- **数据**：12.7k★ / 1.4k fork / 103 个 release / Apache-2.0
- **活跃度**：v3.39.0 发布于 2026-05-20，商业公司（Heroic Labs）持续投入，活跃度最有保障
- **能力**：认证、匹配、排行榜、聊天、云存储全部开箱即用；官方提供 Godot / Unity / Cocos（社区）客户端 SDK
- **短板**：业务逻辑跑在受限运行时里（Go 二等公民、Lua/TS 一等）；底层黑盒，部署依赖 PostgreSQL；**学不到服务器架构本身**
- **对本项目**：定位冲突——本项目的目的就是亲手写出 Nakama 替你做掉的那些事。仅当想快速验证玩法原型时值得考虑

### 4.6 Origin —— 被低估的国产分布式引擎（备选）

- **仓库**：https://github.com/duanhf2012/origin （注意作者 ID 是 duanhf2012，网传 duanhf2333 是错的）
- **数据**：1.7k★ / 69 个 release / Apache-2.0
- **活跃度**：v2.2.0 发布于 2025-12-15，仍在维护但节奏比 due/Cherry 慢半拍
- **架构**：Node → Service → Module 三层抽象；内置 HTTP/TCP/WS、MySQL/Redis 模块、排行榜、性能分析器（能检测慢操作和疑似死锁）
- **学习价值**：⭐⭐⭐ 设计朴素直接，但文档和示例不如 due 完善

## 五、针对本项目（斗地主）的最终建议

### 5.1 选型决策

**不直接引入任何框架，继续手写最小核心**，理由：

1. 本项目是学习项目，CLAUDE.md 规划的目录（`internal/ws` / `room` / `game` / `pattern`）本身就是一个微型框架的形状——亲手实现它的学习收益远大于调用现成框架
2. 斗地主固定 3 人一桌、回合制、消息频率低，用不到任何分布式能力，`gorilla/websocket`（或标准库 + `coder/websocket`）加「单房间单 goroutine」就是完整解
3. 框架的价值改为「参考读物」：自己实现一个模块后，去 due / Cherry 里看同一问题的工业解法

### 5.2 对照学习路线

| 自己实现的模块 | 对照阅读 | 重点看什么 |
|----------------|----------|------------|
| `internal/ws` 连接管理、消息封包 | Zinx 的 Connection/MsgHandler | worker 池、读写分离协程、缓冲与批量发送 |
| `internal/room` 房间生命周期 | Cherry 的 Actor 实现 | 邮箱（mailbox）模式、串行处理如何免锁 |
| `internal/game` 状态机 | due 的 Node 事件路由 | 消息如何路由到具体房间/玩家上下文 |
| 协议层（`proto/`） | Pitaya 的 protocol 包 | 路由压缩、心跳、握手设计 |
| 匹配、断线重连（后期） | due 的 Gate/Node 分层、Nakama 的 matchmaker 文档 | 网关与逻辑分离后会话如何迁移 |

### 5.3 关于 Zinx 的最终态度

之前在 Zinx 上投入的学习没有贬值：它仍在维护（2026-06 还有提交），其网络层设计仍是合格的教学样本。真正的版本判断是——**Zinx 教会你的是 2018 年起就不变的网络层基本功，而 due/Cherry 教的是 2023 年后才成型的模块化/Actor 工程范式**，两者是衔接关系，不是替代关系。

## 六、数据来源与时效性说明

- 所有仓库数据（stars、提交日期、release 版本与日期、license、归档状态）于 **2026-06-11** 通过 GitHub 仓库页面、commits 页面与 releases atom feed 逐一抓取核实
- 社区评价参考：[Go 语言中文网框架汇总](https://studygolang.com/articles/29184)、[知乎：Go 游戏服务器框架功能分析对比](https://zhuanlan.zhihu.com/p/693738476)、[GitHub game-server topic (Go)](https://github.com/topics/game-server?l=go)
- star 数会持续变化，活跃度结论建议每 6–12 个月复查一次；复查方法：看 releases atom feed（`https://github.com/<owner>/<repo>/releases.atom`）和 commits 页的最近提交内容是否为实质变更
- 本文档修正了此前归档问答（[go-game-backend-technology-selection.md](go-game-backend-technology-selection.md)）中的两处错误：due 的仓库地址应为 `dobyte/due`（非 `devagame/due`）；「Zinx 已停更、设计过时」的说法与 2026 年的实际维护状态不符
