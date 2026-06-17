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
		{name: "两张不同点数非法", cards: []card.Card{{Rank: card.Seven}, {Rank: card.Eight}}, want: Invalid},

		{name: "火箭-小大王", cards: []card.Card{{Rank: card.SmallJoker}, {Rank: card.BigJoker}}, want: Rocket},
		{name: "火箭-顺序无关", cards: []card.Card{{Rank: card.BigJoker}, {Rank: card.SmallJoker}}, want: Rocket},

		{name: "三张", cards: []card.Card{{Rank: card.Five}, {Rank: card.Five}, {Rank: card.Five}}, want: Three},
		{name: "三张-含杂牌非法", cards: []card.Card{{Rank: card.Five}, {Rank: card.Five}, {Rank: card.Six}}, want: Invalid},

		{name: "炸弹", cards: []card.Card{{Rank: card.Nine}, {Rank: card.Nine}, {Rank: card.Nine}, {Rank: card.Nine}}, want: Bomb},
		{name: "炸弹-含杂牌非法", cards: []card.Card{{Rank: card.Nine}, {Rank: card.Nine}, {Rank: card.Nine}, {Rank: card.Ten}}, want: Invalid},

		{name: "空牌组非法", cards: []card.Card{}, want: Invalid},

		{name: "顺子-五连", cards: []card.Card{{Rank: card.Three}, {Rank: card.Four}, {Rank: card.Five}, {Rank: card.Six}, {Rank: card.Seven}}, want: Straight},
		{name: "顺子-到A", cards: []card.Card{{Rank: card.Ten}, {Rank: card.Jack}, {Rank: card.Queen}, {Rank: card.King}, {Rank: card.Ace}}, want: Straight},
		{name: "顺子-乱序也认", cards: []card.Card{{Rank: card.Six}, {Rank: card.Three}, {Rank: card.Five}, {Rank: card.Seven}, {Rank: card.Four}}, want: Straight},
		{name: "顺子-含2非法", cards: []card.Card{{Rank: card.Jack}, {Rank: card.Queen}, {Rank: card.King}, {Rank: card.Ace}, {Rank: card.Two}}, want: Invalid},
		{name: "顺子-含王非法", cards: []card.Card{{Rank: card.Ten}, {Rank: card.Jack}, {Rank: card.Queen}, {Rank: card.King}, {Rank: card.SmallJoker}}, want: Invalid},
		{name: "顺子-不连续非法", cards: []card.Card{{Rank: card.Three}, {Rank: card.Four}, {Rank: card.Five}, {Rank: card.Six}, {Rank: card.Eight}}, want: Invalid},
		{name: "顺子-有重复非法", cards: []card.Card{{Rank: card.Three}, {Rank: card.Four}, {Rank: card.Five}, {Rank: card.Six}, {Rank: card.Six}}, want: Invalid},
		{name: "四连不算顺子", cards: []card.Card{{Rank: card.Three}, {Rank: card.Four}, {Rank: card.Five}, {Rank: card.Six}}, want: Invalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Identify(tt.cards); got != tt.want {
				t.Errorf("Identify(%v) = %s, 期望 %s", tt.cards, got, tt.want)
			}
		})
	}
}
