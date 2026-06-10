# go-doudizhu-lab

前后端分离的斗地主（Dou Dizhu）全栈游戏开发学习项目。重点不是做出一个完整游戏，而是完整走一遍**全栈游戏开发流程**：客户端开发、服务器开发、通信协议设计、联调与部署。

- **前端**：Cocos Creator 3.8.8（TypeScript）
- **后端**：Go
- **通信**：WebSocket（协议定义统一放在 `proto/`）

## 目录结构

```
client/    # 游戏前端（Cocos Creator 3.8.8 项目）
server/    # 游戏后端（Go）
proto/     # 前后端共享的通信协议定义
docs/      # 设计文档与学习笔记
scripts/   # 开发与构建脚本
```

各目录内有独立 README 说明用途和初始化方式。

## 环境要求

- Go 1.22+
- Cocos Creator 3.8.8（通过 Cocos Dashboard 安装）

## 快速开始

```bash
# 后端
cd server
go run ./cmd/server

# 前端：用 Cocos Creator 3.8.8 打开 client/ 目录，点击预览
```

> 项目处于早期阶段，前后端代码均从零搭建中，以上命令随开发进度生效。
