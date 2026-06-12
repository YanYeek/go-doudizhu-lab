package greeting

import "testing"

// TestWelcome 是本仓库第一个表驱动测试：
// 用例写成一张表，循环逐条跑，加新用例只需加一行——
// 之后牌型判定的所有测试都会是这个形状。
func TestWelcome(t *testing.T) {
	tests := []struct {
		name string // 用例名，失败时显示
		in   string // 输入
		want string // 期望输出
	}{
		{name: "正常名字", in: "斗地主服务器", want: "斗地主服务器 已启动"},
		{name: "空字符串回退默认名", in: "", want: "无名服务 已启动"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Welcome(tt.in)
			if got != tt.want {
				t.Errorf("Welcome(%q) = %q, 期望 %q", tt.in, got, tt.want)
			}
		})
	}
}
