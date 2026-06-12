// Package greeting 是阶段 0 的工具链复健包：
// 用最小的代码演示 internal 私有目录、包文档注释和表驱动测试。
// 阶段 1 开始写真实的 card 包时，本包将被删除。
package greeting

import "fmt"

// Welcome 返回服务器启动问候语。
func Welcome(name string) string {
	if name == "" {
		name = "无名服务"
	}
	return fmt.Sprintf("%s 已启动", name)
}
