// Package pattern 负责识别一手牌的「牌型」（单张、对子、顺子、炸弹……）。
//
// 它是 card 之上的第二块纯逻辑：只判断「这几张牌算什么牌型」，
// 同样不依赖任何网络代码，可用表驱动测试覆盖。
// 当前实现：单张、对子、三张、炸弹、火箭；后续再加顺子、连对、三带、飞机等。
package pattern

// 跨包导入：用 card 包定义的牌类型。pattern 依赖 card，card 不反过来依赖 pattern（单向）。
import "github.com/YanYeek/go-doudizhu-lab/server/internal/card"

// Kind 是牌型的种类。和 card 里的 Suit/Rank 一样，用 `type + iota` 当枚举。
type Kind uint8

const (
	Invalid Kind = iota // 0 非法/无法识别的牌组
	Single              // 1 单张：1 张牌
	Pair                // 2 对子：2 张同点数的牌
	Three               // 3 三张：3 张同点数的牌
	Bomb                // 4 炸弹：4 张同点数的牌
	Rocket              // 5 火箭：小王 + 大王（斗地主最大的牌型）
)

// String 让 Kind 打印成可读文字（满足 fmt.Stringer），方便测试和日志。
func (k Kind) String() string {
	switch k {
	case Single:
		return "单张"
	case Pair:
		return "对子"
	case Three:
		return "三张"
	case Bomb:
		return "炸弹"
	case Rocket:
		return "火箭"
	default:
		return "非法"
	}
}

// Identify 判断 cards 是哪种牌型，不认识的一律返回 Invalid。
//
// 参数 []card.Card 用了 card 包的类型，写法是「包名.类型名」。
func Identify(cards []card.Card) Kind {
	switch len(cards) {
	case 1:
		// 一张牌永远是单张（哪怕是王，单出也算单张）。
		return Single
	case 2:
		// 两张牌可能是火箭或对子：先判火箭（小王+大王），再判同点数对子。
		if isRocket(cards) {
			return Rocket
		}
		if allSameRank(cards) {
			return Pair
		}
		return Invalid
	case 3:
		if allSameRank(cards) {
			return Three
		}
		return Invalid
	case 4:
		if allSameRank(cards) {
			return Bomb
		}
		return Invalid
	default:
		return Invalid
	}
}

// allSameRank 判断所有牌点数是否相同。
// 调用方保证 cards 非空——这里从第 2 张起逐一和第 1 张比点数。
func allSameRank(cards []card.Card) bool {
	for _, c := range cards[1:] { // cards[1:]：从下标 1 到末尾的子切片
		if c.Rank != cards[0].Rank {
			return false
		}
	}
	return true
}

// isRocket 判断是否火箭：恰好一张小王 + 一张大王（两张顺序不限）。
func isRocket(cards []card.Card) bool {
	if len(cards) != 2 {
		return false
	}
	a, b := cards[0].Rank, cards[1].Rank // 一行声明并赋值两个变量
	// 行尾的 || 表示表达式还没完，下一行接着写（Go 允许在二元运算符后换行）。
	return (a == card.SmallJoker && b == card.BigJoker) ||
		(a == card.BigJoker && b == card.SmallJoker)
}
