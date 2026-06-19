package pattern

// 和 pattern.go 同属 pattern 包，只是按职责分文件：
// pattern.go 管「识别牌型」，compare.go 管「比大小」。Go 里同一个包可以拆成多文件。
import "github.com/YanYeek/go-doudizhu-lab/server/internal/card"

// Play 是一手「已识别」的牌：牌型 + 用于比大小的点数。
// 把"它是什么牌型"和"它多大"打包在一起，比较时就不必反复重新识别。
type Play struct {
	Kind Kind      // 牌型
	Rank card.Rank // 比大小用的点数：单点数牌型取共同点数，顺子取最小那张；火箭用不到
	Len  int       // 牌张数：顺子这类"带长度"的牌型靠它决定能否相互比较
}

// Parse 把一组牌识别成 Play。无法识别时返回 ok=false（Go 的 comma-ok 风格）。
func Parse(cards []card.Card) (Play, bool) {
	kind := Identify(cards)
	if kind == Invalid {
		return Play{}, false // Play{} 是零值结构体
	}

	play := Play{Kind: kind, Len: len(cards)}
	switch kind {
	case Rocket:
		// 火箭不看点数，Rank 保持零值。
	case Straight:
		play.Rank = minRank(cards) // 顺子用最小的那张比大小
	case ThreeWithSingle, ThreeWithPair:
		play.Rank = tripletRank(cards) // 三带：看 3 张那个主牌定大小，附带牌不算
	default:
		// 单张/对子/三张/炸弹都是同点数，取第一张就能代表整手的大小。
		play.Rank = cards[0].Rank
	}
	return play, true
}

// Beats 判断 p 能否压过 prev（p 出在 prev 之后，是否合法地更大）。
// 规则：火箭最大；炸弹压一切普通牌型、炸弹之间比点数；普通牌型必须同类型再比点数。
// 注意"压过"是严格更大——一样大压不过（斗地主必须出更大的）。
//
// (p Play) 是结构体接收者：这是挂在 Play 上的方法，调用写成 newPlay.Beats(lastPlay)。
func (p Play) Beats(prev Play) bool {
	// 1) 火箭最大，压一切。
	if p.Kind == Rocket {
		return true
	}
	if prev.Kind == Rocket {
		return false // 没有牌能压火箭
	}

	// 2) 炸弹 vs 普通牌型：炸弹压普通；普通压不过炸弹。
	if p.Kind == Bomb && prev.Kind != Bomb {
		return true
	}
	if p.Kind != Bomb && prev.Kind == Bomb {
		return false
	}

	// 3) 到这里：要么两边都是炸弹，要么两边都是普通牌型。
	//    普通牌型必须「同类型」才能比（单张压单张、对子压对子……），类型不同不能压。
	//    当前每种牌型长度固定，所以同 Kind 即同长度；将来加顺子需要再比长度。
	if p.Kind != prev.Kind {
		return false
	}

	// 4) 带长度的牌型（如顺子）：长度不同不能比（5 连压不了也压不过 6 连）。
	//    单点数牌型同 Kind 必同长度，这条对它们恒成立、无副作用。
	if p.Len != prev.Len {
		return false
	}

	// 5) 同类型、同长度，比点数，严格更大才算压过。
	return p.Rank > prev.Rank
}
