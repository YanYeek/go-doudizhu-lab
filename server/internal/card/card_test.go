package card

import "testing"

// TestNewDeckComposition 验证一副牌的构成，对应规则文档第 9 节
// 「牌堆恰好 54 张且每个 ID 唯一」的最低测试要求。
func TestNewDeckComposition(t *testing.T) {
	deck := NewDeck()

	if len(deck) != 54 {
		t.Fatalf("一副牌应有 54 张，实际 %d 张", len(deck))
	}

	seen := make(map[int]bool)
	rankCount := make(map[Rank]int)
	for _, c := range deck {
		if c.ID < 0 || c.ID >= 54 {
			t.Errorf("ID 越界: %d", c.ID)
		}
		if seen[c.ID] {
			t.Errorf("ID 重复: %d", c.ID)
		}
		seen[c.ID] = true
		rankCount[c.Rank]++
	}
	if len(seen) != 54 {
		t.Errorf("应有 54 个唯一 ID，实际 %d 个", len(seen))
	}

	// 13 个普通点数各 4 张。
	for r := Three; r <= Two; r++ {
		if rankCount[r] != 4 {
			t.Errorf("点数 %s 应有 4 张，实际 %d 张", r, rankCount[r])
		}
	}
	// 大小王各 1 张。
	if rankCount[SmallJoker] != 1 {
		t.Errorf("小王应有 1 张，实际 %d 张", rankCount[SmallJoker])
	}
	if rankCount[BigJoker] != 1 {
		t.Errorf("大王应有 1 张，实际 %d 张", rankCount[BigJoker])
	}
}

// TestRankOrdering 验证点数常量按强弱递增赋值——这是「比大小直接比 Rank
// 数值」这一设计能成立的前提。
func TestRankOrdering(t *testing.T) {
	order := []Rank{
		Three, Four, Five, Six, Seven, Eight, Nine, Ten,
		Jack, Queen, King, Ace, Two, SmallJoker, BigJoker,
	}
	for i := 1; i < len(order); i++ {
		if order[i-1] >= order[i] {
			t.Errorf("大小顺序错误：%s 应小于 %s", order[i-1], order[i])
		}
	}
}

// TestCardString 验证显示文本，重点是王没有花色。
func TestCardString(t *testing.T) {
	tests := []struct {
		name string
		card Card
		want string
	}{
		{name: "黑桃3", card: Card{Rank: Three, Suit: Spade}, want: "♠3"},
		{name: "红桃10", card: Card{Rank: Ten, Suit: Heart}, want: "♥10"},
		{name: "方块2", card: Card{Rank: Two, Suit: Diamond}, want: "♦2"},
		{name: "小王无花色", card: Card{Rank: SmallJoker, Suit: Joker}, want: "小王"},
		{name: "大王无花色", card: Card{Rank: BigJoker, Suit: Joker}, want: "大王"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.String(); got != tt.want {
				t.Errorf("String() = %q, 期望 %q", got, tt.want)
			}
		})
	}
}
