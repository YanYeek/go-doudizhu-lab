# scripts — 开发与构建脚本

存放跨前后端的辅助脚本，例如：

- 一键启动本地开发环境（起后端 + 提示打开 Cocos 预览）
- 协议代码生成（如使用 Protobuf）
- 构建与部署脚本

## Python 开发助手

`python/dev.py` 是零第三方依赖的跨平台开发入口，可在 Windows、macOS 和 Linux
上运行。所有命令都可以从仓库任意目录执行。

```bash
# 一键启动 Redis、etcd，并在前台运行 Gate
python scripts/python/dev.py up

# 另开一个终端，运行 Go 测试客户端验证通信链路
python scripts/python/dev.py testclient

# 查看 Redis、etcd 与 Gate 状态
python scripts/python/dev.py status

# 运行测试、静态检查与构建
python scripts/python/dev.py check

# 停止本地 Docker 依赖
python scripts/python/dev.py deps-down
```

完整命令：

| 命令 | 用途 |
|------|------|
| `up` | 启动 Docker 依赖，并在前台运行 Gate |
| `server` | 只运行 Gate，适合依赖已经启动时使用 |
| `testclient` | 运行一次性 Go 测试客户端，向 Gate 发一条消息验证链路（需先 `up`） |
| `deps-up` | 只启动 Redis 与 etcd |
| `deps-down` | 停止并删除 Redis 与 etcd 容器 |
| `status` | 查看依赖与 Gate 端口状态 |
| `test` | 运行单元测试（快，不依赖 Redis/etcd）|
| `test-integration` | 运行集成测试（需先 `deps-up`，带 `integration` build tag）|
| `vet` | 运行 Go 静态检查 |
| `build` | 构建服务器 |
| `check` | 提交前自检：单元测试 + 静态检查 + 构建 |
| `doctor` | 检查 Python、Go、Docker 与端口状态 |

`test` / `vet` / `build` 看着只是一行 go 命令，包成子命令的好处是**从仓库根目录就能跑**，
省去每次 `cd server/`。`test` 是不依赖网络的单元测试；`test-integration` 需要先
`deps-up` 起 Redis/etcd，按 `integration` build tag 隔离端到端测试。

测试输出默认精简（一包一行，便于代码增多后仍看得过来）。想看每个测试用例的细节，
加 `-v`：`python scripts/python/dev.py test -v`。

`up` 会让服务器保持在前台，方便查看日志。按 `Ctrl+C` 会停止 Gate，但保留
Redis 与 etcd，便于继续开发；当天开发结束时再执行 `deps-down`。

执行 `up` 或 `deps-up` 时，如果 Docker 引擎尚未运行，脚本会在 Windows 和
macOS 上自动启动 Docker Desktop，并等待引擎就绪后继续。`status` 和
`deps-down` 只检测当前状态，不会自动启动 Docker。Linux 上仍需自行启动 Docker
服务，因为服务管理方式和权限因发行版而异。

该脚本定位为本地开发助手。未来部署到 Linux 服务器时，不使用 `go run` 长期
运行服务，而是构建正式二进制文件，并交由 systemd 或容器负责启动、重启与日志
管理；届时可以继续在本目录增加独立的部署命令。
