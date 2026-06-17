package pattern

import (
	"testing"

	"github.com/YanYeek/go-doudizhu-lab/server/internal/card"
)

// TestIdentify 用表驱动覆盖当前支持的牌型与几种非法情况。
// 构造牌时只关心点数，花色/ID 留零值即可（Identify 不看它们）。
func TestIdentify(t *testing.T) {
	tests := []struct {
		name  string
		cards []card.Card
		want  Kind
	}{
		// 切片字面量里，内层 {Rank: ...} 省略了类型——Go 会从元素类型 card.Card 推断。
		{name: "单张", cards: []card.Card{{Rank: card.Three}}, want: Single},
		{name: "单张-王也算", cards: []card.Card{{Rank: card.BigJoker}}, want: Single},
		{name: "对子", cards: []card.Card{{Rank: card.Seven}, {Rank: card.Seven}}, want: Pair},
		{name: "两张不同点数不是对子", cards: []card.Card{{Rank: card.Seven}, {Rank: card.Eight}}, want: Invalid},
		{name: "大小王不是对子", cards: []card.Card{{Rank: card.SmallJoker}, {Rank: card.BigJoker}}, want: Invalid},
		{name: "空牌组非法", cards: []card.Card{}, want: Invalid},
		{name: "三张暂不识别", cards: []card.Card{{Rank: card.Five}, {Rank: card.Five}, {Rank: card.Five}}, want: Invalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Identify(tt.cards); got != tt.want {
				t.Errorf("Identify(%v) = %s, 期望 %s", tt.cards, got, tt.want)
			}
		})
	}
}
