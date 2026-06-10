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

输出是什么？1,99

### Q2 map 取值

```go
m := map[string]int{"x": 1}
v := m["y"]
v2, ok := m["y"]
fmt.Println(v, v2, ok)
```

输出是什么？不知道

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

输出是什么？不知道

## 第 2 层：并发（zinx 时代应该接触过）

### Q4 channel

```go
func main() {
    ch := make(chan int)
    ch <- 1
    fmt.Println(<-ch)
}
```

这段程序运行后会发生什么？如果有问题，怎么改？

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

输出是什么？（提示：不同 Go 版本下答案可能不一样，知道多少说多少）

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

这段程序的行为是什么？

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

---

## 作答记录与评估

（待作答后补充）
