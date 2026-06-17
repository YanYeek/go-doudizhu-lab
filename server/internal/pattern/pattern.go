// Package pattern 负责识别一手牌的「牌型」（单张、对子、顺子、炸弹……）。
//
// 它是 card 之上的第二块纯逻辑：只判断「这几张牌算什么牌型」，
// 同样不依赖任何网络代码，可用表驱动测试覆盖。
// 当前只实现最小的两种：单张和对子；后续逐步加顺子、三带、炸弹、火箭等。
package pattern

// 跨包导入：用 card 包定义的牌类型。pattern 依赖 card，card 不反过来依赖 pattern（单向）。
import "github.com/YanYeek/go-doudizhu-lab/server/internal/card"

// Kind 是牌型的种类。和 card 里的 Suit/Rank 一样，用 `type + iota` 当枚举。
type Kind uint8

const (
	Invalid Kind = iota // 0 非法/无法识别的牌组
	Single              // 1 单张：1 张牌
	Pair                // 2 对子：2 张同点数的牌
)

// String 让 Kind 打印成可读文字（满足 fmt.Stringer），方便测试和日志。
func (k Kind) String() string {
	switch k {
	case Single:
		return "单张"
	case Pair:
		return "对子"
	default:
		return "非法"
	}
}

// Identify 判断 cards 是哪种牌型。
// 现在只认得单张和对子，其它一律返回 Invalid（后续再扩展）。
//
// 参数 []card.Card 用了 card 包的类型，写法是「包名.类型名」。
func Identify(cards []card.Card) Kind {
	switch len(cards) {
	case 1:
		// 一张牌永远是单张（哪怕是王，单出也算单张）。
		return Single
	case 2:
		// 两张同点数才算对子。card.Rank 底层是 uint8、可比较，直接用 == 判等。
		// 小王(16)和大王(17)点数不同，两张王 == 不成立，自然不会被当成对子——
		// 正好符合规则「大小王不能组成对子」。
		if cards[0].Rank == cards[1].Rank {
			return Pair
		}
		return Invalid
	default:
		return Invalid
	}
}
