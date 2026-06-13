# 软件工程能力全景地图（面向本项目的学习版）

> 研究日期：2026-06-13
> 目的：以本仓库（Go + due 后端、Cocos 客户端、monorepo）为练兵场，全面且系统地
> 提升软件工程能力。本文是「地图」——标出每个能力维度的权威参考与本项目落地点；
> 具体怎么练，挂接到 [01-roadmap.md](01-roadmap.md) 的阶段推进。
> 方法：多源联网检索 + 交叉核对，所有维度均给出一手权威来源；版本信息已核实
> 到 2026-06（SWEBOK v4.0、OWASP Top 10 2025、DORA 2025 报告）。

## 0. 总纲：软件工程的官方知识地图

**SWEBOK v4.0**（IEEE 软件工程知识体系指南，2024-10-15 发布）是行业公认的
学科边界定义，共 **18 个知识域**，v4 新增了三个：软件架构、软件安全、软件工程
运维——这个增量本身就说明了行业重心的变化（架构显学化、安全左移、DevOps 常态化）。

- 来源：[IEEE Computer Society — SWEBOK](https://www.computer.org/education/bodies-of-knowledge/software-engineering)（一手，权威定义）
- 参考：[Wikipedia — SWEBOK](https://en.wikipedia.org/wiki/Software_Engineering_Body_of_Knowledge)（二手，概览用）

本文不照搬 18 个知识域，而是压缩成对单人学习项目**可操作**的 11 个维度，
每个维度标注优先级：

- **P0**＝现在就该练（阶段 0~1 已经在练）
- **P1**＝阶段 2~3 引入
- **P2**＝阶段 4 之后按需引入

## 1. 编码质量与风格 ｜ P0

**要点**：可读性优先级排序——清晰 > 简单 > 简洁 > 可维护 > 一致（Google Go
风格指南的原则顺序）；包组织遵循官方模块布局（`cmd/` 放入口、`internal/` 藏实现）。

**权威参考**：
- [Google Go Style Guide](https://google.github.io/styleguide/go/)（一手，Google 全体 Go 代码的规范）
- [Go 官方：Organizing a Go module](https://go.dev/doc/modules/layout)（一手，cmd/internal 布局的官方出处）
- 书：《代码大全 2》（编码层面的百科全书）、《重构》第 2 版（识别坏味道的词汇表）

**本项目落地**：`server/` 已按官方布局组织；每次写完代码跑 `go vet`、`gofmt`；
阶段 1 写牌组逻辑时对照风格指南的「清晰优先」原则自查命名与函数长度。

**常见误区**：把「风格」理解为个人偏好之争——Go 的世界里风格是工具链强制的
（gofmt 没有配置项），把争论时间省下来。

## 2. 测试策略 ｜ P0

**要点**：测试金字塔——大量单元测试打底、少量集成测试居中、极少端到端测试封顶；
层级越高越慢越脆弱（flaky）。Go 的表驱动测试是单元层的惯用形态。

**权威参考**：
- [The Practical Test Pyramid — Ham Vocke @ martinfowler.com](https://martinfowler.com/articles/practical-test-pyramid.html)（一手，最常被引用的完整论述）

**本项目落地**：金字塔在本项目的具体形状——
- 底座：`internal/card`、`internal/pattern`（未来的牌组/卡牌效果包）的表驱动测试，
  这是规则正确性的生命线，目标是每种卡牌效果至少一张用例表
- 中层：房间状态机的集成测试（不连真实网络，用内存消息驱动状态机）
- 塔尖：极少数 WebSocket 全链路冒烟测试（起真实 Gate，跑一局最短对局）

**常见误区**：单人项目最容易犯「只写跑得通的 happy path」——规则引擎的价值
恰恰在边界用例（爆炸牌插回位置 0、两张 Nope 抵消、Clone 复制 Clone 被拒绝）。

## 3. 版本控制与协作 ｜ P0

**要点**：主干开发（trunk-based development）——分支少而短命、至少每日合回主干。
DORA 多年数据反复验证：主干开发与四大交付指标全部正相关。提交信息结构化
（Conventional Commits），让历史可检索、可自动化。

**权威参考**：
- [DORA Capabilities: Trunk-based Development](https://dora.dev/capabilities/trunk-based-development/)（一手，附研究数据）
- [trunkbaseddevelopment.com](https://trunkbaseddevelopment.com/continuous-integration/)（专题站，实践细节）

**本项目落地**：已在执行——main 主干 + Conventional Commits（scope 区分
client/server/proto/docs）+ 原子提交。下一步可以练：功能较大时开短命分支，
体验「分支活不过三天」的纪律。

**常见误区**：单人项目把 Git 当网盘（一天结束 `git add . && git commit -m "update"`）。
提交粒度 = 未来排查问题时 `git bisect` 的分辨率。

## 4. 代码评审 ｜ P0（单人变体）

**要点**：Google 的评审标准只有一句话——**「CL 让代码库整体健康度变好就该批准，
不存在完美的代码」**。评审者依次看：设计、功能、复杂度、测试、命名、注释、一致性。

**权威参考**：
- [Google Engineering Practices — How to do a code review](https://google.github.io/eng-practices/review/)（一手，业界引用最多的评审手册）

**本项目落地**：单人没有同事，但评审可以这样练——
- 提交前用 Google 清单自查一遍 diff（设计→功能→复杂度→测试→命名）
- 用 AI 评审做「第二双眼睛」：本仓库已装 codex 插件（`/codex:review`），
  Claude Code 也有 `/code-review`——两个模型交叉挑刺
- 隔周回看自己两周前的代码并写下「现在会怎么写」，这是单人最有效的评审替代

**常见误区**：把评审等同于挑语法错——清单里「设计是否合理」排第一，语法排最后。

## 5. CI / CD ｜ P1（阶段 2 引入）

**要点**：持续集成 = 主干开发 + 每次提交自动跑快速测试套件，保证主干永远可用；
持续交付在此之上让发布成为「随时可按的按钮」。

**权威参考**：
- [Continuous Integration — Martin Fowler](https://www.martinfowler.com/articles/continuousIntegration.html)（一手，概念原典，2024 年重写版）
- [DORA Capabilities: Continuous Integration](https://dora.dev/capabilities/continuous-integration/)（一手，研究证据）

**本项目落地**：GitHub Actions 一个 workflow 即可起步：push 时跑
`go build ./... && go test ./... && go vet ./...`（本地 `scripts/python/dev.py check`
的云端版）。阶段 4 联调时再加 Cocos 构建检查。**不需要**部署流水线——
本项目的"CD"先做到「主干永远绿」就够了。

**常见误区**：单人项目搭复杂流水线是典型过度工程；反过来「反正只有我一个人，
不用 CI」也错——CI 防的不是别人，是上周的自己。

## 6. 架构设计与决策记录 ｜ P1

**要点**：架构能力 = 做决策的能力 + 让决策可追溯的能力。ADR（架构决策记录，
Michael Nygard 2011 年定型）用五段式记录每个重要决策：标题、状态、上下文、
决策、后果——几个月后回看「系统为什么长这样」全靠它。

**权威参考**：
- [Documenting Architecture Decisions — Michael Nygard](https://www.cognitect.com/blog/2011/11/15/documenting-architecture-decisions)（一手，ADR 原典）
- [adr.github.io](https://adr.github.io/)（模板与工具集合）
- 书：《Designing Data-Intensive Applications》（理解状态、复制、一致性的底层心智模型，P2 阶段精读）

**本项目落地**：本仓库已经无意识地写过两份"准 ADR"（框架调研、技术选型）。
正式化：建 `docs/adr/` 目录，把「为什么用 due 而不是自研/Zinx」「为什么游戏
从斗地主换成爆炸小黄人」补记为 ADR-0001/0002，之后每个重大决策（协议格式、
状态机设计、Nope 响应窗口裁定）都新增一条。

**常见误区**：把 ADR 写成设计文档——ADR 记录的是**决策和理由**，两页以内，
重点是「当时考虑过哪些选项、为什么没选它们」。

## 7. 文档工程 ｜ P0（已在练，给出框架）

**要点**：Diátaxis 框架把文档分四象限——教程（学习导向）、操作指南（任务导向）、
参考（信息导向）、解释（理解导向）。四类受众和写法完全不同，混写是大多数
文档难读的根因。

**权威参考**：
- [Diátaxis](https://diataxis.fr/)（一手，框架官网）

**本项目落地**：现有文档恰好可以归类——`docs/learning/` 是教程+解释、
`scripts/README.md` 是操作指南、`docs/rules/` 是参考、ADR 是解释。
写新文档前先问「这是四象限里的哪一类」，避免在规则参考里夹学习心得。

**常见误区**：文档求全求长。Diátaxis 的核心洞见是「一篇文档只服务一种需求」。

## 8. 可观测性 ｜ P1（阶段 3 起步，P2 完善）

**要点**：三大信号——日志（带时间戳的事件）、指标（随时间聚合的数值）、
链路追踪（一次请求跨组件的完整路径）；现代实践强调三者**关联**而非孤立。
OpenTelemetry 是厂商中立的事实标准。

**权威参考**：
- [OpenTelemetry — Observability Primer](https://opentelemetry.io/docs/concepts/observability-primer/)（一手，概念入门最佳起点）

**本项目落地**：分三步走，别一上来就上全家桶——
1. 阶段 3：结构化日志先行（Go 标准库 `log/slog`），每条对局事件带 room_id/player_id
2. 阶段 4 联调：学会用日志还原一局完整对局（这就是"可观测"的最朴素定义）
3. P2：给 Gate/Node 加最基础的指标（在线连接数、消息吞吐），体验 OTel SDK

**常见误区**：把可观测性等同于「装个监控面板」。先问自己：服务器半夜出 bug，
我手头的日志能不能还原案发现场？

## 9. 安全 ｜ P2（但有两条 P0 红线）

**要点**：OWASP Top 10（2025 版）是应用安全的标准风险清单，前三位：失效的
访问控制、安全配置错误、软件供应链失效。对游戏服务器，核心心法是
**「永远不信任客户端」**——所有规则判定在服务端，客户端只发意图。

**权威参考**：
- [OWASP Top 10:2025](https://owasp.org/Top10/2025/)（一手，2025 正式版）

**本项目落地**：两条 P0 红线现在就生效——
1. 服务端权威：手牌、牌堆顺序、随机数只存在于服务端（规则文档第 6 节已确立）
2. 机密不入库：`.env` 已在 gitignore；将来有 token/密钥绝不硬编码

其余（输入校验强化、依赖漏洞扫描 `govulncheck`、限流）到阶段 5 再练。

**常见误区**：学习项目跳过安全意识培养——「服务端权威」这一条恰恰是游戏
服务器架构的根，从第一行协议代码就该体现。

## 10. 配置与部署运维 ｜ P2

**要点**：Twelve-Factor App 是云时代应用的体检表，对本项目最相关的三条：
配置存环境变量（III）、依赖显式声明（II）、日志当事件流（XI）。判断配置是否
干净的试金石：「代码库现在开源，会不会泄露任何凭据？」

**权威参考**：
- [The Twelve-Factor App](https://12factor.net/config)（一手，方法论原典）

**本项目落地**：due 框架原生支持配置文件 + 环境变量；Redis/etcd 地址等
环境差异项逐步从代码移到配置。部署（systemd/容器、正式二进制）在
`scripts/README.md` 已有预留笔记，阶段 5 实践。

**常见误区**：本地开发图省事把地址硬编码——养成习惯比改造成本低得多。

## 11. 交付效能度量 ｜ P2（先知道，后实践）

**要点**：DORA 四大指标——部署频率、变更前置时间（吞吐）；变更失败率、
恢复时间（稳定性）。2025 报告新增返工率（Rework Rate），并转向七种团队
画像，把摩擦、倦怠纳入可持续效能。这是「工程实践是否有效」的科学衡量尺。

**权威参考**：
- [DORA 官网与年度报告](https://dora.dev/research/2024/dora-report/)（一手）
- [Google Cloud — State of DevOps](https://cloud.google.com/devops/state-of-devops)（一手，2025 版入口）

**本项目落地**：单人项目不必算指标，但可以用它的**方向感**自检：
提交到推送的间隔是否在变短（前置时间）？推上去的代码多久发现问题（失败率）？
学习节奏是否可持续（倦怠维度）？

## 12. 反过度工程化守则（单人项目专属）

**要点**：YAGNI（You Aren't Gonna Need It，源自极限编程）——不为想象中的
未来需求写代码。配合 KISS 与「先让它工作，再让它正确，最后让它快」。

**权威参考**：
- [YAGNI — Martin Fowler](https://martinfowler.com/bliki/Yagni.html)（一手）

**本项目的三条具体戒律**：
1. 规则引擎按 2~5 人设计但 MVP 只实现 2 人对局（规则文档第 8 节）——这是
   YAGNI 的正确用法：接口留余地，实现不提前
2. 不做：分布式扩容、账号系统、热更新——直到学习路线真的走到那里
3. 每次想加「以后可能用得上」的东西，先在 ADR 里写下来代替写代码

## 13. 优先级总表与路线图挂接

| 维度 | 优先级 | 开始练的阶段 | 第一个动作 |
|------|--------|--------------|------------|
| 编码质量 | P0 | 已在练 | 写牌组逻辑时对照 Go 风格指南自查 |
| 测试策略 | P0 | 已在练 | 卡牌效果全部表驱动测试 |
| 版本控制 | P0 | 已在练 | 保持原子提交，尝试短命分支 |
| 代码评审 | P0 | 阶段 1 | 提交前过 Google 清单 + AI 交叉评审 |
| 文档工程 | P0 | 已在练 | 新文档先定 Diátaxis 象限 |
| CI/CD | P1 | 阶段 2 | 加 GitHub Actions：build+test+vet |
| 架构与 ADR | P1 | 阶段 2 | 建 docs/adr/，补记 due 选型决策 |
| 可观测性 | P1 | 阶段 3 | 用 slog 结构化日志记录对局事件 |
| 安全 | P2（红线 P0） | 阶段 5 | 红线立即生效：服务端权威 + 机密不入库 |
| 配置与部署 | P2 | 阶段 5 | 环境差异项移出代码 |
| 效能度量 | P2 | 随时自检 | 用四指标方向感复盘学习节奏 |

## 14. 来源可信度说明

- **一手权威**（直接采信）：IEEE/SWEBOK、Google eng-practices 与 Go 风格指南、
  Go 官方文档、martinfowler.com、dora.dev、OWASP、OpenTelemetry 官方文档、
  12factor.net、diataxis.fr、Nygard 原文
- **二手参考**（仅佐证概览）：Wikipedia、各厂商博客（Octopus、GitLab 等对
  OWASP/DORA 新版的解读，已与一手来源交叉核对，未发现冲突）
- 所有「2025/2026 最新版」表述（SWEBOK v4.0、OWASP Top 10:2025、DORA 2025
  报告）均经检索核实，非凭记忆。
