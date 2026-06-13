# CLAUDE.md

## 工具分工（最高优先级，优先于本文件其余所有约定）

本仓库已安装并启用 codex 插件（`codex@openai-codex`，底层为 OpenAI Codex CLI）。
为节省 Claude Code 当前模型（Opus 4.8）的算力，按以下规则分派任务：

| 任务类型 | 交给谁 | 怎么做 |
|----------|--------|--------|
| 文档类工作（写/改 `docs/`、README、注释、规则与笔记） | **codex** | 委派给 codex（`/codex` 相关 skill 或 codex 子代理） |
| 代码审核 / Code Review | **codex** | 用 `/codex:review` |
| Git 管理（暂存、commit、分支、合并、写提交信息） | **codex** | 委派 codex 执行 git 操作 |
| `git push`（推送到远端） | **Claude Code（Opus 4.8）** | codex 沙箱默认禁网，push 由 Claude Code 兜底 |
| 高难度任务（复杂设计、疑难排查、架构决策） | **Claude Code（Opus 4.8）** | 自己处理 |
| 编写代码（实现功能、写测试） | **Claude Code（Opus 4.8）** | 自己处理 |

- 默认优先把上表前三类（文档、评审、Git 本地操作）下派给 codex；只有高难度任务、写代码和 `git push` 才由 Opus 4.8 亲自做。
- **实测约束**：codex 运行在 workspace-write 沙箱里，能读写工作区、能本地 commit，但**不能访问网络**——所以 `git push`、拉取依赖等联网操作一律由 Claude Code 完成。
- codex 的 `/codex:*` skill（如 `/codex:review`）在会话启动时加载；新装或刚启用后需**重启会话**才能用。本会话内若 skill 不可用，可直接调底层 `codex exec` CLI 代替。
- 若 codex 不可用（未登录 / CLI 缺失 / 调用失败），降级由 Claude Code 直接完成，并提示用户。
- 该分工规则可被用户的临时明确指令覆盖（例如用户说"这次你自己写文档"）。

## 仓库定位

前后端分离的斗地主全栈游戏开发**学习项目**（monorepo）。目标是完整走一遍全栈游戏开发流程——客户端、服务器、通信协议、联调部署——而不是单纯做出一个游戏。因此写代码时优先考虑清晰、可讲解，必要处解释设计决策（写进 `docs/`），而不是一味追求简洁或性能。

- 前端：Cocos Creator **3.8.8**（TypeScript），在 `client/`
- 后端：Go，在 `server/`
- 通信：WebSocket，协议定义统一放 `proto/`，前后端都以它为准

## 目录约定

```
client/              # Cocos Creator 项目（assets/ 下是场景与 TS 脚本）
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

- 后端业务逻辑放 `server/internal/`，不对外暴露
- 牌型判定等纯逻辑必须有表驱动测试（`xxx_test.go` 与被测代码同目录）
- 前端逻辑代码放 `client/assets/scripts/`，按场景/模块分子目录
- 修改通信消息时，`proto/` 的定义和文档必须同步更新

## 常用命令

```bash
# 后端（在 server/ 目录下）
go run ./cmd/server     # 启动服务器
go test ./...           # 全量测试
go vet ./...            # 静态检查

# 前端：用 Cocos Creator 3.8.8 打开 client/，编辑器内预览/构建
```

## 注意事项

- 后端模块路径：`github.com/YanYeek/go-doudizhu-lab/server`（在 server/ 下 go mod init）
- `client/library/`、`client/temp/` 等是 Cocos 编辑器缓存，已被 .gitignore 忽略，不要提交
- 提交信息遵循 Conventional Commits，描述用中文，不带 emoji；scope 用 `client` / `server` / `proto` / `docs` 区分改动范围
- 远端仓库：https://github.com/YanYeek/go-doudizhu-lab
