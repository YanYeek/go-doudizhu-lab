// Package card 定义斗地主的牌、点数、花色与一副牌的生成。
//
// 这是项目第一块「纯逻辑」：不依赖任何网络代码，只描述牌本身。
// 按规则文档（docs/rules/doudizhu-rules.md 第 2 节）：
//   - 一张牌由唯一 ID、点数 Rank、花色 Suit 组成；
//   - 花色只用于区分实体牌和前端显示，不参与牌型识别与大小比较；
//   - 大小顺序：3 < 4 < … < K < A < 2 < 小王 < 大王。
package card

import (
	"fmt"
	"math/rand/v2"
)

// Suit 是牌的花色。只用于显示，不参与大小比较。
type Suit uint8

const (
	Spade   Suit = iota // 黑桃 ♠
	Heart               // 红桃 ♥
	Club                // 梅花 ♣
	Diamond             // 方块 ♦
	Joker               // 王（大小王专用，没有真实花色）
)

func (s Suit) String() string {
	switch s {
	case Spade:
		return "♠"
	case Heart:
		return "♥"
	case Club:
		return "♣"
	case Diamond:
		return "♦"
	default:
		return ""
	}
}

// Rank 是牌的点数，也是斗地主里比较大小的唯一依据。
// 常量的数值刻意按强弱递增赋值：数值越大，牌越大。
// 这样以后比大小直接比 Rank 的数值即可，不用额外查表。
type Rank uint8

const (
	Three      Rank = iota + 3 // 3（最小）
	Four                       // 4
	Five                       // 5
	Six                        // 6
	Seven                      // 7
	Eight                      // 8
	Nine                       // 9
	Ten                        // 10
	Jack                       // J
	Queen                      // Q
	King                       // K
	Ace                        // A（比 K 大）
	Two                        // 2（比 A 大，斗地主特有）
	SmallJoker                 // 小王
	BigJoker                   // 大王（最大）
)

// rankLabels 把点数映射成显示文本。用包级变量避免每次调用都重建。
var rankLabels = map[Rank]string{
	Three: "3", Four: "4", Five: "5", Six: "6", Seven: "7",
	Eight: "8", Nine: "9", Ten: "10", Jack: "J", Queen: "Q",
	King: "K", Ace: "A", Two: "2", SmallJoker: "小王", BigJoker: "大王",
}

func (r Rank) String() string {
	if label, ok := rankLabels[r]; ok {
		return label
	}
	return fmt.Sprintf("Rank(%d)", uint8(r))
}

// Card 是一张牌的实例。
// ID 在一副牌内唯一，用来追踪「具体某一张牌」——例如服务端校验
// 「这张牌确实在你手上」时，靠的是 ID 而不是点数（点数会重复）。
type Card struct {
	ID   int  // 一副牌内唯一（0..53）
	Rank Rank // 点数，比大小的唯一依据
	Suit Suit // 花色，仅用于显示
}

func (c Card) String() string {
	switch c.Rank {
	case SmallJoker, BigJoker:
		return c.Rank.String() // 王没有花色，只显示点数
	default:
		return c.Suit.String() + c.Rank.String()
	}
}

// NewDeck 生成一副标准的 54 张斗地主牌：
// 4 种花色 × 13 个点数（3~2）共 52 张，外加小王、大王。
// 每张牌的 ID 在 0..53 内唯一，顺序固定（未洗牌）。
func NewDeck() []Card {
	deck := make([]Card, 0, 54)
	id := 0

	for _, s := range []Suit{Spade, Heart, Club, Diamond} {
		for r := Three; r <= Two; r++ {
			deck = append(deck, Card{ID: id, Rank: r, Suit: s})
			id++
		}
	}

	deck = append(deck, Card{ID: id, Rank: SmallJoker, Suit: Joker})
	id++
	deck = append(deck, Card{ID: id, Rank: BigJoker, Suit: Joker})

	return deck
}

// Shuffle 用 Fisher-Yates 算法就地打乱一副牌。
//
// 关于"就地"——这是 Go 切片最容易踩的点：切片本质是指向底层数组的"句柄"，
// 把它传进函数后交换其中的元素，会直接改到调用方手里那一份，不产生副本。
// 所以这里没有返回值：调用方传进来的 deck 洗完就是乱的了。
//
// 随机用的是 math/rand/v2（Go 1.22+）：它自动随机播种，不再需要老 math/rand
// 那句 rand.Seed(...)；rand.IntN(n) 返回 [0, n) 的随机整数。
func Shuffle(deck []Card) {
	// 从最后一张往前：每一步把当前牌和它之前（含自己）的随机一张交换。
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}
