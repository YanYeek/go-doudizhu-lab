# proto — 前后端共享协议

前后端通信的消息协议定义，单一数据源（single source of truth），两端都从这里生成/引用，避免协议不一致。

可选方案（确定技术选型后二选一）：

- **JSON + TypeScript/Go 手写类型**：简单直观，适合起步
- **Protobuf**：`.proto` 文件定义，`protoc` 分别生成 Go 和 TS 代码，适合后期

约定：消息命名用 `C2S_` / `S2C_` 前缀区分方向（如 `C2S_PlayCards`、`S2C_GameStart`）。
