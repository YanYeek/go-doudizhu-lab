# 项目协作规范

> 本文件是项目协作规范的唯一维护源。
>
> `AGENTS.md` 和 `CLAUDE.md` 只保留入口说明，不再分别维护规则正文。以后修改项目约定、Agent 偏好、文档偏好或开发流程时，只改本文件。

## 仓库定位

前后端分离的斗地主全栈游戏开发**学习项目**（monorepo）。目标是完整走一遍全栈游戏开发流程：客户端、服务器、通信协议、联调部署，而不是单纯做出一个游戏。

写代码和写文档时优先考虑：

- 清晰、可讲解、能复盘；
- 每一步尽量能运行、能观察、能测试；
- 重要设计决策写进 `docs/`；
- 不为了追求简洁或性能牺牲学习价值。

## 技术栈

- 前端：Cocos Creator **3.8.8**（TypeScript），在 `client/`
- 后端：Go，在 `server/`
- 后端框架：due v2
- 通信：WebSocket
- 协议定义：统一放 `proto/`，前后端都以它为准
- 本地依赖：Docker Compose 管理 Redis 与 etcd
- 本地脚本：`scripts/python/dev.py`

## 目录约定

```text
client/              # Cocos Creator 项目
server/
  cmd/server/        # 后端入口（main 包）
  internal/
    card/            # 牌的定义、洗牌发牌
    pattern/         # 牌型识别与比较
    game/            # 对局流程、状态机、规则
    room/            # 房间与匹配
    ws/              # WebSocket 连接与消息收发
proto/               # 协议定义，消息命名 C2S_ / S2C_ 前缀
docs/                # 设计文档、架构决策、学习笔记
scripts/             # 跨端辅助脚本
```

约定：

- 后端业务逻辑放 `server/internal/`，不对外暴露。
- 牌型判定等纯逻辑必须有表驱动测试，`xxx_test.go` 与被测代码同目录。
- 前端逻辑代码放 `client/assets/scripts/`，按场景或模块分子目录。
- 修改通信消息时，`proto/` 的定义和文档必须同步更新。
- `client/library/`、`client/temp/` 等 Cocos 编辑器缓存不要提交。

## 常用命令

```bash
# 一键启动 Redis、etcd，并运行 Gate/Node
python scripts/python/dev.py up

# 查看依赖与 Gate 端口状态
python scripts/python/dev.py status

# 运行测试、静态检查与构建
python scripts/python/dev.py check

# 停止本地 Docker 依赖
python scripts/python/dev.py deps-down
```

后端单独命令（在 `server/` 目录下）：

```bash
go run ./cmd/server
go test ./...
go vet ./...
```

前端：用 Cocos Creator 3.8.8 打开 `client/`，在编辑器内预览或构建。

## 学习笔记规范

`docs/learning/` 使用 HTML 作为活动学习笔记格式，不再直接维护 Markdown 页面。

当前规则：

- 活动入口是 `docs/learning/index.html`。
- 每篇学习笔记都是独立 HTML 页面。
- 页面必须包含更图形化的说明方式，例如 SVG 流程图、结构图、概念关系图。
- 每篇尽量加入精准、恰当的生活类比，帮助把抽象概念落到经验里。
- 左侧导航、首页卡片和文件名前缀必须使用严格一致的数字序号。
- 不能出现两个相同编号；如果新增或插入页面，要顺延后续编号。
- 当前编号范围是 `00` 到 `13`。
- 原 Markdown 源文件归档在 `docs/archive/learning-markdown-sources/`。
- 如需重新生成 HTML，运行：

```bash
python scripts/python/build_learning_html.py
```

编号和文件名必须对应，例如：

```text
03-due-container-lifecycle.html
04-software-engineering-map.html
05-due-gate-websocket.html
```

不要只在页面显示层修序号；文件名、生成器配置、首页、侧边栏必须同步。

## 文档偏好

- 文档优先中文。
- 学习文档要保留“为什么这么做”，不要只写结论。
- 对复杂概念优先使用图示、流程、对照表和生活类比。
- 规则文档、协议文档和 ADR 要区分用途，不要把学习心得混进规则参考。
- 重大方向决策写入 `docs/decisions/`。
- 桌游或商业游戏资料只做学习归档；公开教程避免复用第三方名称、角色、美术、卡面文字和高度近似表达。

## 代码偏好

- 优先匹配仓库已有结构和风格。
- 这是学习项目，写代码时给**非显而易见处**加讲解注释（面向 C#/TS 背景的学习者）：新语法、内置函数（如 `len`/`make`/`copy`）、Go 惯用法、容易踩的坑（如切片共享底层数组）都注；真正显而易见的行不必注，避免噪音。注释讲"这是什么语法 / 为什么这么写"，不复述代码字面。
- 后端规则逻辑先写纯函数和表驱动测试，再接入 due。
- 服务器保持权威：客户端只提交意图，规则判定、手牌、牌堆和随机数都由服务端控制。
- 功能推进采用纵向切片：纯逻辑 -> 单元测试 -> due 路由 -> Go 测试客户端 -> 最小 Cocos 联调。
- 不提前做账号体系、复杂匹配、分布式扩容、Kubernetes 或 AI，除非学习路线走到那里。

## Git 规范

- 提交信息遵循 Conventional Commits。
- 描述用中文，不带 emoji。
- scope 优先使用 `client` / `server` / `proto` / `docs` / `scripts`。
- 一次 commit 只做一件可解释、可验证的事。
- 当前项目是个人独立开发，可以在完整、原子的提交上直接使用 `main`；未完成工作不要作为正式提交推送。
- 远端仓库：https://github.com/YanYeek/go-doudizhu-lab

## Agent 工具分工

本节主要服务于 Claude Code / Codex 同时参与项目时的分工。

| 任务类型 | 默认交给谁 | 说明 |
|---|---|---|
| 文档类工作（`docs/`、README、注释、规则与笔记） | Codex | 文档生成、重构、学习笔记整理优先给 Codex |
| 代码审核 / Code Review | Codex | 可使用 Codex review 能力 |
| Git 管理（暂存、commit、push、分支、合并、提交信息） | Claude Code 或当前明确执行 Git 的 Agent | Git 操作要先检查状态，避免误提交 WIP |
| 高难度设计、疑难排查、架构决策 | Claude Code 或当前主 Agent | 需要完整上下文和判断时不要机械外包 |
| 编写代码、实现功能、写测试 | Claude Code 或当前主 Agent | 以项目结构和测试为准 |

补充规则：

- 用户的临时明确指令优先于默认分工。
- 如果某个工具不可用，当前 Agent 直接完成并说明情况。
- Git 操作必须先看 `git status` 和相关 diff。
- 未完成、未验证、不可独立解释的工作默认不提交，更不要推送到远端。

## 当前重要文档入口

- `docs/learning/index.html`：可视化学习笔记入口
- `docs/rules/doudizhu-rules.md`：斗地主规则与实现规范
- `docs/decisions/0001-select-dou-dizhu.md`：选择斗地主作为教程游戏的决策
- `docs/go-game-server-frameworks-research.md`：Go 游戏服务器框架调研
- `scripts/README.md`：本地脚本说明
