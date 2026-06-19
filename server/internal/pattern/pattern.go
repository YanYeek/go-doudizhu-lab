// Package pattern 负责识别一手牌的「牌型」（单张、对子、顺子、炸弹……）。
//
// 它是 card 之上的第二块纯逻辑：只判断「这几张牌算什么牌型」，
// 同样不依赖任何网络代码，可用表驱动测试覆盖。
// 当前实现：单张、对子、三张、炸弹、火箭；后续再加顺子、连对、三带、飞机等。
package pattern

import (
	"slices" // 标准库：切片工具，这里用 slices.Equal 比较"次数形状"
	"sort"   // 标准库：排序。判断顺子是否连续要先排序

	// 跨包导入：用 card 包定义的牌类型。pattern 依赖 card，card 不反过来依赖 pattern（单向）。
	"github.com/YanYeek/go-doudizhu-lab/server/internal/card"
)

// Kind 是牌型的种类。和 card 里的 Suit/Rank 一样，用 `type + iota` 当枚举。
type Kind uint8

const (
	Invalid         Kind = iota // 0 非法/无法识别的牌组
	Single                      // 1 单张：1 张牌
	Pair                        // 2 对子：2 张同点数的牌
	Three                       // 3 三张：3 张同点数的牌
	ThreeWithSingle             // 4 三带一：3 张同点 + 1 张单牌
	ThreeWithPair               // 5 三带二：3 张同点 + 1 个对子
	Straight                    // 6 顺子：5 张及以上、点数连续（不含 2 和大小王）
	Bomb                        // 7 炸弹：4 张同点数的牌
	Rocket                      // 8 火箭：小王 + 大王（斗地主最大的牌型）
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
	case ThreeWithSingle:
		return "三带一"
	case ThreeWithPair:
		return "三带二"
	case Straight:
		return "顺子"
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
		// 4 张：炸弹（4 同点）或三带一（3+1）。
		if allSameRank(cards) {
			return Bomb
		}
		if isThreeWithSingle(cards) {
			return ThreeWithSingle
		}
		return Invalid
	case 5:
		// 5 张：顺子（5 连）或三带二（3+2）。
		if isStraight(cards) {
			return Straight
		}
		if isThreeWithPair(cards) {
			return ThreeWithPair
		}
		return Invalid
	default:
		// len 为 0 或 ≥ 6：目前只有更长的顺子落在这里。
		if isStraight(cards) {
			return Straight
		}
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

// isStraight 判断是否是顺子：5 张及以上、点数连续、无重复，且不含 2 和大小王。
func isStraight(cards []card.Card) bool {
	if len(cards) < 5 {
		return false
	}

	// 把点数抽成 []int 以便排序。顺子里不能出现 2 和大小王（点数 > A）。
	ranks := make([]int, len(cards))
	for i, c := range cards {
		if c.Rank > card.Ace { // Ace=14；比它大的是 2、小王、大王
			return false
		}
		ranks[i] = int(c.Rank) // 把 card.Rank 转成 int 才能用 sort.Ints
	}

	sort.Ints(ranks) // 标准库排序，原地升序排列

	// 排好序后，每张必须正好比前一张大 1：既保证连续，也顺带排除了重复。
	for i := 1; i < len(ranks); i++ {
		if ranks[i] != ranks[i-1]+1 {
			return false
		}
	}
	return true
}

// minRank 返回牌组里最小的点数。调用方保证 cards 非空。
func minRank(cards []card.Card) card.Rank {
	m := cards[0].Rank
	for _, c := range cards[1:] {
		if c.Rank < m {
			m = c.Rank
		}
	}
	return m
}

// countByRank 统计每个点数出现的次数。
// map 取不存在的 key 会返回零值 0，所以 counts[r]++ 不必先判断 key 是否存在（呼应诊断题 Q2）。
func countByRank(cards []card.Card) map[card.Rank]int {
	counts := make(map[card.Rank]int) // map 必须 make 才能用
	for _, c := range cards {
		counts[c.Rank]++
	}
	return counts
}

// countShape 返回各点数出现次数的「形状」（升序）。
// 例：三带一 → [1 3]，三带二 → [2 3]，炸弹 → [4]，对子 → [2]。
// 把它配合 slices.Equal 比对，就能很直观地判别这些「按数量组合」的牌型。
func countShape(cards []card.Card) []int {
	shape := make([]int, 0)
	for _, n := range countByRank(cards) { // range map 给出 (键, 值)；这里只要值 n
		shape = append(shape, n)
	}
	sort.Ints(shape) // map 遍历顺序随机，排序后形状才稳定可比
	return shape
}

// isThreeWithSingle 判断三带一：次数形状是 [1 3]（一张单牌 + 三张同点）。
func isThreeWithSingle(cards []card.Card) bool {
	return slices.Equal(countShape(cards), []int{1, 3})
}

// isThreeWithPair 判断三带二：次数形状是 [2 3]（一个对子 + 三张同点）。
func isThreeWithPair(cards []card.Card) bool {
	return slices.Equal(countShape(cards), []int{2, 3})
}

// tripletRank 返回出现 3 次的那个点数（三带类牌型靠主牌比大小）。调用方保证存在。
func tripletRank(cards []card.Card) card.Rank {
	for r, n := range countByRank(cards) {
		if n == 3 {
			return r
		}
	}
	return 0 // 正常不会到这里
}
