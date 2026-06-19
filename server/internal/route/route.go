// Package route 定义前后端共享的消息路由号与消息体（Go 侧的协议单一数据源）。
//
// 约定（与 CLAUDE.md 一致）：
//   - 每条消息占用一个唯一的 int32 路由号；
//   - 请求体（客户端→服务端）以 Req 结尾，响应体（服务端→客户端）以 Res 结尾；
//   - 字段用 json tag，因为 due 默认用 json 编解码器在两端之间传输消息。
//
// proto/ 目录用自然语言记录同一套协议，便于将来 TypeScript 客户端对齐。
package route

// 协议要传「牌」，所以引用 card 包的类型。route 依赖 card（单向）。
import "github.com/YanYeek/go-doudizhu-lab/server/internal/card"

// 消息路由号。新增消息时在这里追加一个常量，避免散落的魔法数字。
const (
	// Greet 问候：客户端发来名字，服务端回一句欢迎语。
	// 这是项目的第一条路由，用于打通「客户端 → Gate → Node → 响应」整条链路。
	Greet int32 = 1

	// Login 登录：客户端报上玩家 id，服务端把这条连接绑定到该 uid。
	// 绑定后网关就知道这条连接属于谁，是玩家系统的起点。
	Login int32 = 2

	// Notify 通知：服务端主动推给客户端的消息（S2C，没有对应的客户端请求）。
	// 用来验证「绑定后服务端能按 uid 反向找到玩家并主动推送」。
	Notify int32 = 3

	// PlayCheck 出牌校验：客户端提交一手牌，服务端判定牌型是否合法、是什么型。
	// 牌型判定只在服务端做（服务端权威），客户端说了不算。
	PlayCheck int32 = 4
)

// GreetReq 是 Greet 路由的请求体（C2S_Greet）。
type GreetReq struct {
	Name string `json:"name"`
}

// GreetRes 是 Greet 路由的响应体（S2C_Greet）。
type GreetRes struct {
	Message string `json:"message"`
}

// LoginReq 是 Login 路由的请求体（C2S_Login）。
// 现在没有真正的鉴权，客户端直接报上自己的 uid——学习阶段够用。
type LoginReq struct {
	UserID int64 `json:"user_id"`
}

// LoginRes 是 Login 路由的响应体（S2C_Login）。
type LoginRes struct {
	Message string `json:"message"`
}

// NotifyPush 是 Notify 的推送体（S2C_Notify）。由服务端主动发起，没有请求体。
type NotifyPush struct {
	Text string `json:"text"`
}

// PlayCheckReq 是 PlayCheck 的请求体（C2S_PlayCheck）：客户端提交的一手牌。
type PlayCheckReq struct {
	Cards []card.Card `json:"cards"`
}

// PlayCheckRes 是 PlayCheck 的响应体（S2C_PlayCheck）：服务端的判定结果。
type PlayCheckRes struct {
	Valid bool   `json:"valid"` // 是否是合法牌型
	Kind  string `json:"kind"`  // 牌型名（如"对子"）；非法时为"非法"
}
