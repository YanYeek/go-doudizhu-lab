# go-doudizhu-lab

用 Go 实现斗地主（Dou Dizhu）的实验项目。目标是从零搭建斗地主的核心玩法逻辑——发牌、牌型判定、出牌规则、对局流程——作为 Go 语言与游戏逻辑的练习场。

> 项目刚初始化，代码尚未开始编写，以下安装与使用方式为预期约定，随开发进度更新。

## 安装

需要 Go 1.22 或更高版本。

```bash
git clone <仓库地址>
cd go-doudizhu-lab
go mod tidy
```

## 使用

```bash
# 运行
go run ./cmd/doudizhu

# 测试
go test ./...
```
