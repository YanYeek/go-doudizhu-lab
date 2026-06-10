# server — 游戏后端

Go 编写的斗地主游戏服务器。

## 初始化

```bash
cd server
go mod init github.com/YanYeek/go-doudizhu-lab/server
```

## 计划结构

```
cmd/server/      # 程序入口（main 包）
internal/
  card/          # 牌的定义、洗牌发牌
  pattern/       # 牌型识别与比较
  game/          # 对局流程、状态机、规则
  room/          # 房间与匹配
  ws/            # WebSocket 连接与消息收发
```

## 常用命令

```bash
go run ./cmd/server   # 启动服务器
go test ./...         # 全量测试
```
