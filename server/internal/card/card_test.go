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

// TestShufflePreservesCards 验证洗牌只改顺序、不增不减不改牌。
// 洗牌是随机的，所以不测"具体洗成什么顺序"，而测一个对任何洗牌都成立的不变量：
// 洗完还是同样的 54 张（ID 集合不变、每张内容不变）。
// 这是测随机函数的正确姿势——测性质，不测输出。
func TestShufflePreservesCards(t *testing.T) {
	deck := NewDeck()
	before := make(map[int]Card, len(deck))
	for _, c := range deck {
		before[c.ID] = c
	}

	Shuffle(deck) // 就地洗牌：deck 这一份直接被打乱

	if len(deck) != 54 {
		t.Fatalf("洗牌后应仍有 54 张，实际 %d", len(deck))
	}
	seen := make(map[int]bool, len(deck))
	for _, c := range deck {
		if seen[c.ID] {
			t.Errorf("洗牌后出现重复 ID: %d", c.ID)
		}
		seen[c.ID] = true
		if before[c.ID] != c {
			t.Errorf("ID %d 的牌内容被改了：%v -> %v", c.ID, before[c.ID], c)
		}
	}
	if len(seen) != 54 {
		t.Errorf("洗牌后应有 54 个唯一 ID，实际 %d", len(seen))
	}
}

// TestDealDistributesAllCards 验证发牌不重不漏：三家各 17、底牌 3，合起来正好 54 张。
func TestDealDistributesAllCards(t *testing.T) {
	deck := NewDeck()
	hands, bottom := Deal(deck)

	for i, h := range hands {
		if len(h) != 17 {
			t.Errorf("第 %d 家应有 17 张，实际 %d", i, len(h))
		}
	}
	if len(bottom) != 3 {
		t.Errorf("底牌应有 3 张，实际 %d", len(bottom))
	}

	seen := make(map[int]bool, 54)
	collect := func(cards []Card) {
		for _, c := range cards {
			if seen[c.ID] {
				t.Errorf("ID %d 出现在多处", c.ID)
			}
			seen[c.ID] = true
		}
	}
	for _, h := range hands {
		collect(h)
	}
	collect(bottom)
	if len(seen) != 54 {
		t.Errorf("发出的牌应共 54 张唯一，实际 %d", len(seen))
	}
}

// TestDealHandsAreIndependent 验证每家手牌是独立副本，而不是 deck 的子切片别名。
// 这是发牌最容易踩的坑：若用 deck[:17] 直接切，手牌会和 deck 共享底层数组，
// 改一处串一片。这里改动手牌后，断言原 deck 不受影响。
func TestDealHandsAreIndependent(t *testing.T) {
	deck := NewDeck()
	original := deck[0] // 发牌前记下 deck[0]
	hands, _ := Deal(deck)

	hands[0][0] = Card{ID: 999, Rank: BigJoker, Suit: Joker} // 篡改第一家第一张

	if deck[0] != original {
		t.Errorf("改手牌竟影响了 deck[0]：说明手牌是别名而非副本（deck[0]=%v）", deck[0])
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
