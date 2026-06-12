# Go 舒适区诊断（代码阅读题）

> 目的：按学习区理论（舒适区 / 学习区 / 困难区）定位当前 Go 水平，
> 据此规划游戏服务器框架的分阶段学习路线。
> 作答方式：阅读代码，说出输出或行为 + 一句话原因。
> **不确定就直接写"不知道"或"猜的"**——诚实作答才能测准，这不是考试。
> 日期：2026-06-11

## 第 1 层：语言基础（GOPATH 时代就有）

### Q1 切片

```go
a := []int{1, 2, 3}
b := a[:2]
b = append(b, 99)
fmt.Println(a)
```

输出是什么？答：1,99

### Q2 map 取值

```go
m := map[string]int{"x": 1}
v := m["y"]
v2, ok := m["y"]
fmt.Println(v, v2, ok)
```

输出是什么？答：不知道

### Q3 方法接收者

```go
type Counter struct{ n int }

func (c Counter) IncVal()  { c.n++ }
func (c *Counter) IncPtr() { c.n++ }

func main() {
    c := Counter{}
    c.IncVal()
    c.IncPtr()
    fmt.Println(c.n)
}
```

输出是什么？答：不知道

## 第 2 层：并发（zinx 时代应该接触过）

### Q4 channel

```go
func main() {
    ch := make(chan int)
    ch <- 1
    fmt.Println(<-ch)
}
```

这段程序运行后会发生什么？如果有问题，怎么改？答：细节忘了

### Q5 goroutine 与循环变量

```go
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println(i)
    }()
}
wg.Wait()
```

输出是什么？（提示：不同 Go 版本下答案可能不一样，知道多少说多少）答：不知道

### Q6 select

```go
ch := make(chan int)
select {
case v := <-ch:
    fmt.Println("got", v)
case <-time.After(time.Second):
    fmt.Println("timeout")
}
```

这段程序的行为是什么？答：Channel，信道的变量，接触到的不同变量做出不同的行为

## 第 3 层：现代 Go（go mod 之后的世界）

### Q7 模块与导入

```go
// 文件 server/go.mod 内容：
// module github.com/YanYeek/go-doudizhu-lab/server
//
// 某个 .go 文件里：
import "github.com/YanYeek/go-doudizhu-lab/server/internal/card"
```

这个 import 会去哪里找代码？还需要把项目放进 GOPATH 吗？路径里的 `internal` 有什么特殊含义？

答：不知道

### Q8 错误包装

```go
var ErrNotFound = errors.New("not found")

func find() error {
    return fmt.Errorf("查询用户失败: %w", ErrNotFound)
}

func main() {
    err := find()
    fmt.Println(errors.Is(err, ErrNotFound))
    fmt.Println(err == ErrNotFound)
}
```

两行各输出什么？`%w` 是干什么的？

答：不知道

### Q9 泛型

```go
func Map[T, U any](xs []T, f func(T) U) []U {
    out := make([]U, 0, len(xs))
    for _, x := range xs {
        out = append(out, f(x))
    }
    return out
}

func main() {
    fmt.Println(Map([]int{1, 2, 3}, func(i int) int { return i * 2 }))
}
```

这段代码能编译吗？这是 Go 的什么特性、大约哪个版本引入？输出是什么？

答：编译不知道，泛型和泛型约束 版本忘了 输出不确定

---

## 作答记录与评估

> 评估日期：2026-06-12

### 逐题判定

| 题号 | 考点 | 作答 | 判定 | 正确答案要点 |
|------|------|------|------|--------------|
| Q1 | 切片共享底层数组 | `1,99` | 半对 | 输出 `[1 2 99]`。`b := a[:2]` 与 `a` 共享底层数组且 cap=3，`append` 未扩容，直接把 99 写进了下标 2 的位置。你意识到了 append 会影响 `a`（方向对），但具体写到哪、长度怎么变没说清 |
| Q2 | map 零值与 comma-ok | 不知道 | 不会 | 输出 `0 0 false`。取不存在的 key 返回元素类型的零值，不会 panic；`, ok` 形式用来区分「值就是 0」和「key 不存在」 |
| Q3 | 值接收者 vs 指针接收者 | 不知道 | 不会 | 输出 `1`。值接收者拿到的是副本，`IncVal` 改不到原结构体；指针接收者才能改。这是以后写游戏状态（房间、玩家）必踩的坑 |
| Q4 | 无缓冲 channel 阻塞 | 细节忘了 | 不会 | 死锁 panic。无缓冲 channel 的发送会阻塞到有人接收，main 自己发自己收，永远等不到。改法：`make(chan int, 1)` 加缓冲，或把接收放进另一个 goroutine |
| Q5 | goroutine + 循环变量 | 不知道 | 不会 | Go ≤1.21：很可能打出 `3 3 3`（闭包共享同一个 `i`）；Go ≥1.22 起每轮循环 `i` 是新变量，输出 `0 1 2`（顺序不定）。这是语言近年最大的语义修正之一 |
| Q6 | select 多路复用 | 概念模糊 | 半对 | 没有人往 `ch` 发数据，第一个 case 永远不就绪，1 秒后 `time.After` 触发，打印 `timeout`。select 是「等多个 channel，谁先就绪走谁」——你记得大方向，细节（阻塞行为、超时模式）丢了 |
| Q7 | go mod / internal | 不知道 | 不会 | import 按 go.mod 的 module 路径在**本模块内**解析，直接映射到 `server/internal/card` 目录，完全不需要 GOPATH；`internal` 是编译器强制的私有目录，模块外无法 import |
| Q8 | 错误包装 errors.Is/%w | 不知道 | 不会 | `true` / `false`。`%w` 把原错误包进新错误形成错误链，`errors.Is` 沿链比对；`==` 只比表层所以是 false。Go 1.13 引入，现代 Go 错误处理的基石 |
| Q9 | 泛型 | 知道名词 | 边缘 | 能编译，泛型是 Go 1.18 引入，输出 `[2 4 6]`。你知道「泛型/泛型约束」这两个词但读不出行为 |

### 舒适区地图

```
舒适区（还在）：
  ✔ 能读 Go 语法，不被代码形态吓到
  ✔ 概念名词还在：channel / select / 泛型 / goroutine

学习区（本次路线图的主战场）：
  → 第 1 层基础的"操作性知识"：切片机制、map 零值、值/指针接收者
  → 并发原语的精确行为：channel 阻塞、WaitGroup、select 超时
  → 现代 Go 工具链：go mod、internal、errors 包装、泛型读懂即可

困难区（暂不碰）：
  ✘ 分布式、网关分层、性能调优、反射、unsafe
```

### 结论

GOPATH 时代的知识只剩「语法阅读能力 + 概念索引」，操作性细节基本清零——
这反而是个干净的起点：**不需要纠正坏习惯，只需要把每个知识点在项目里
第一次用到的地方重新焊牢**。路线图原则：不单独刷语法书，每一步开发只引入
1~2 个新知识点，学了立刻写进 `server/` 的真实代码里。

详细路线见 [01-roadmap.md](01-roadmap.md)。
