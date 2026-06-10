# AGENTS.md

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
