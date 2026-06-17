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
	"math/rand/v2" // Go 1.22+ 的新随机库，自动播种；旧库是 "math/rand"
)

// Suit 是牌的花色。只用于显示，不参与大小比较。
// `type X uint8` 基于 uint8 定义一个新类型——Suit 和 uint8 是不同类型，不能直接混用。
type Suit uint8

// const 块 + iota 是 Go 的「枚举」写法（Go 没有 enum 关键字）。
// iota 是块内的行计数器：第一行 = 0，每往下一行自动 +1。
const (
	Spade   Suit = iota // 0 黑桃 ♠
	Heart               // 1 红桃 ♥
	Club                // 2 梅花 ♣
	Diamond             // 3 方块 ♦
	Joker               // 4 王（大小王专用，没有真实花色）
)

// String() string 让 Suit 满足标准库的 fmt.Stringer 接口（隐式满足，无需写 implements）。
// `(s Suit)` 是「接收者」= 别的语言里的 this/self，这里显式命名为 s。
func (s Suit) String() string {
	switch s { // Go 的 switch 每个 case 自动 break，不会贯穿到下一个
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

// iota + 3 让起点从 3 开始：Three=3、Four=4 …… BigJoker=17，数值正好等于牌力强弱。
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
// `map[Rank]string{...}` 是 map 字面量：键类型 Rank、值类型 string。
var rankLabels = map[Rank]string{
	Three: "3", Four: "4", Five: "5", Six: "6", Seven: "7",
	Eight: "8", Nine: "9", Ten: "10", Jack: "J", Queen: "Q",
	King: "K", Ace: "A", Two: "2", SmallJoker: "小王", BigJoker: "大王",
}

func (r Rank) String() string {
	// `value, ok := map[key]` 是 Go 取 map 的惯用法：ok 为 true 表示 key 存在。
	// 这里还用了 if 的「初始化语句」：先取值再判断，label 的作用域仅限这个 if。
	if label, ok := rankLabels[r]; ok {
		return label
	}
	return fmt.Sprintf("Rank(%d)", uint8(r)) // %d 填整数；uint8(r) 是类型转换
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
	case SmallJoker, BigJoker: // 一个 case 可列多个值，用逗号分隔
		return c.Rank.String() // 王没有花色，只显示点数
	default:
		return c.Suit.String() + c.Rank.String() // 字符串用 + 拼接
	}
}

// NewDeck 生成一副标准的 54 张斗地主牌：
// 4 种花色 × 13 个点数（3~2）共 52 张，外加小王、大王。
// 每张牌的 ID 在 0..53 内唯一，顺序固定（未洗牌）。
func NewDeck() []Card {
	// make([]Card, 0, 54)：长度 0、容量 54。先占好 54 张的位子，
	// 下面 append 时就不必反复扩容、搬动底层数组。
	deck := make([]Card, 0, 54)
	id := 0

	// for ... range 遍历右边的切片字面量 []Suit{...}。
	// range 切片给出 (下标, 值)；这里用 _ 丢弃下标，只要值 s。
	for _, s := range []Suit{Spade, Heart, Club, Diamond} {
		// C 风格 for：r 从 Three 到 Two（含），每轮 r++ 加 1。
		for r := Three; r <= Two; r++ {
			// append 往切片尾部追加，返回「可能换了底层数组的」新切片，必须接住返回值。
			// Card{ID: id, ...} 是结构体字面量，按字段名赋值。
			deck = append(deck, Card{ID: id, Rank: r, Suit: s})
			id++
		}
	}

	deck = append(deck, Card{ID: id, Rank: SmallJoker, Suit: Joker})
	id++
	deck = append(deck, Card{ID: id, Rank: BigJoker, Suit: Joker})

	return deck
}

// Deal 把一副 54 张牌发成斗地主的开局：三家各 17 张，外加 3 张底牌（地主牌）。
// 传入的牌应已洗好；因为牌已随机，这里按顺序切块发即可，和逐张轮流发同样公平。
//
// 关键：每家手牌都是 make+copy 出来的独立副本，而不是 deck 的子切片。
// 若写成 hands[i] = deck[i*17:(i+1)*17]，子切片会和 deck 共享同一底层数组——
// 玩家之后排序、出牌、append 时会串改到 deck 甚至别家手牌。copy 出独立副本才安全。
func Deal(deck []Card) (hands [3][]Card, bottom []Card) {
	// len(x) 是内置函数（不是 x.Length）；!= 是不等于。
	if len(deck) != 54 {
		// panic 立刻中止并打印调用栈，用于表达「不该发生」的程序 bug。
		panic(fmt.Sprintf("Deal 需要一副 54 张的牌，实际 %d 张", len(deck)))
	}

	// range hands：hands 是定长 3 的数组，i 依次是 0、1、2。
	for i := range hands {
		hands[i] = make([]Card, 17) // 给第 i 家造一个长度 17 的独立切片
		// deck[low:high] 是切片表达式，左闭右开 [low, high)；i=0 → [0:17] 即前 17 张。
		// copy(dst, src) 逐元素拷贝值，所以 hands[i] 不与 deck 共享底层数组。
		copy(hands[i], deck[i*17:(i+1)*17])
	}

	bottom = make([]Card, 3)
	copy(bottom, deck[51:54]) // 最后 3 张（下标 51、52、53）

	return hands, bottom // 一次返回多个值
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
		j := rand.IntN(i + 1) // [0, i] 内随机；rand.IntN 取代了旧库的 rand.Intn
		// 多重赋值：Go 一行交换两个值，右边先整体求值，无需临时变量。
		deck[i], deck[j] = deck[j], deck[i]
	}
}
