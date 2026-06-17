package pattern

import (
	"testing"

	"github.com/YanYeek/go-doudizhu-lab/server/internal/card"
)

// TestParse 验证识别结果被正确打包成 Play（牌型 + 比大小的点数）。
func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		cards    []card.Card
		wantOK   bool
		wantKind Kind
		wantRank card.Rank
		wantLen  int
	}{
		{name: "单张K", cards: []card.Card{{Rank: card.King}}, wantOK: true, wantKind: Single, wantRank: card.King, wantLen: 1},
		{name: "对子7", cards: []card.Card{{Rank: card.Seven}, {Rank: card.Seven}}, wantOK: true, wantKind: Pair, wantRank: card.Seven, wantLen: 2},
		{name: "顺子取最小点", cards: []card.Card{{Rank: card.Seven}, {Rank: card.Three}, {Rank: card.Five}, {Rank: card.Six}, {Rank: card.Four}}, wantOK: true, wantKind: Straight, wantRank: card.Three, wantLen: 5},
		{name: "火箭", cards: []card.Card{{Rank: card.SmallJoker}, {Rank: card.BigJoker}}, wantOK: true, wantKind: Rocket, wantLen: 2},
		{name: "非法", cards: []card.Card{{Rank: card.Seven}, {Rank: card.Eight}}, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			play, ok := Parse(tt.cards)
			if ok != tt.wantOK {
				t.Fatalf("Parse ok = %v, 期望 %v", ok, tt.wantOK)
			}
			if !ok {
				return // 非法牌组不用再比后面的字段
			}
			if play.Kind != tt.wantKind {
				t.Errorf("Kind = %s, 期望 %s", play.Kind, tt.wantKind)
			}
			if play.Len != tt.wantLen {
				t.Errorf("Len = %d, 期望 %d", play.Len, tt.wantLen)
			}
			if tt.wantKind != Rocket && play.Rank != tt.wantRank { // 火箭不看点数
				t.Errorf("Rank = %v, 期望 %v", play.Rank, tt.wantRank)
			}
		})
	}
}

// TestPlayBeats 覆盖比大小的各条规则：同型比点数、类型不同不能压、炸弹压一切、火箭最大。
func TestPlayBeats(t *testing.T) {
	// 局部闭包当「迷你构造器」，让用例表更紧凑。闭包是赋给变量的函数。
	single := func(r card.Rank) Play { return Play{Kind: Single, Rank: r, Len: 1} }
	pair := func(r card.Rank) Play { return Play{Kind: Pair, Rank: r, Len: 2} }
	three := func(r card.Rank) Play { return Play{Kind: Three, Rank: r, Len: 3} }
	bomb := func(r card.Rank) Play { return Play{Kind: Bomb, Rank: r, Len: 4} }
	// 顺子要带长度：minR 是最小那张的点数，n 是张数。
	straight := func(minR card.Rank, n int) Play { return Play{Kind: Straight, Rank: minR, Len: n} }
	rocket := Play{Kind: Rocket}

	tests := []struct {
		name string
		p    Play
		prev Play
		want bool
	}{
		{"大单张压小单张", single(card.King), single(card.Three), true},
		{"小单张压不过大单张", single(card.Three), single(card.King), false},
		{"同点数压不过(要严格更大)", single(card.Three), single(card.Three), false},
		{"大对子压小对子", pair(card.King), pair(card.Three), true},
		{"单张压不过对子(类型不同)", single(card.King), pair(card.Three), false},
		{"对子压不过单张(类型不同)", pair(card.King), single(card.Three), false},
		{"大三张压小三张", three(card.Ace), three(card.Three), true},
		{"炸弹压单张", bomb(card.Three), single(card.Ace), true},
		{"炸弹压对子", bomb(card.Three), pair(card.Ace), true},
		{"大炸弹压小炸弹", bomb(card.King), bomb(card.Three), true},
		{"小炸弹压不过大炸弹", bomb(card.Three), bomb(card.King), false},
		{"同点炸弹压不过", bomb(card.Five), bomb(card.Five), false},
		{"普通牌压不过炸弹", single(card.Two), bomb(card.Three), false},
		{"火箭压单张", rocket, single(card.Two), true},
		{"火箭压炸弹", rocket, bomb(card.King), true},
		{"炸弹压不过火箭", bomb(card.King), rocket, false},

		{"大顺子压小顺子(同长)", straight(card.Four, 5), straight(card.Three, 5), true},
		{"小顺子压不过大顺子", straight(card.Three, 5), straight(card.Four, 5), false},
		{"长度不同的顺子不能比", straight(card.Three, 6), straight(card.Four, 5), false},
		{"炸弹压顺子", bomb(card.Three), straight(card.Ten, 5), true},
		{"顺子压不过炸弹", straight(card.Ten, 5), bomb(card.Three), false},
		{"火箭压顺子", rocket, straight(card.Ten, 5), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// %+v 打印结构体时带上字段名，失败信息更易读。
			if got := tt.p.Beats(tt.prev); got != tt.want {
				t.Errorf("%+v.Beats(%+v) = %v, 期望 %v", tt.p, tt.prev, got, tt.want)
			}
		})
	}
}
