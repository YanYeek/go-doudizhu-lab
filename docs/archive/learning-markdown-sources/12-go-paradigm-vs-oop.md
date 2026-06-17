# Go 范式:从 C#/TS 的"类"到 Go 的"类型 + 方法 + 包"

> 日期:2026-06-17
>
> 从 C#/TypeScript 切到 Go,最大的坎不是语法,是**组织代码的世界观**变了。
> 本篇以 `server/internal/card/card.go` 为主线,逐个对照讲清这层转变。

## 一句话定调

> C#/TS 的世界观是"万物皆类"——数据和行为都塞进 `class {}`,靠继承和 `implements` 组织。
> Go 故意拆开了:**类型只管数据,方法挂在类型外面,接口靠"长得像"自动满足,包(package)取代类做组织单元。**

Go 没有 `class`、`extends`、`implements`、构造函数、`public/private` 关键字。
这些熟悉的东西全被更小的几块拼出来了。

## 转变 1:`type Suit uint8` —— 给基础类型起个有身份的新名字

```go
type Suit uint8   // 底层是 uint8,但它是一个"新类型"
```

C# 写 `enum Suit : byte {...}`。Go 没有 `enum` 关键字,做法是先定义一个新类型。
关键:`Suit` 和 `uint8` 是**两个不同类型**,不能混用(`Suit(5)` 才能转)。这给你类型安全。
而且 Go 的 `type` 更通用——能给任何底层类型(int/string/函数/切片)起新名字并挂方法。

## 转变 2:`const ( ... iota ... )` —— Go 式枚举

```go
const (
    Three Rank = iota + 3 // iota=0 → 3（最小）
    Four                  // iota=1 → 4
    ...
    BigJoker              // iota=14 → 17（最大）
)
```

`iota` 是 const 块里的行计数器,从 0 每行 +1,后续行自动套用首行表达式。
`iota + 3` 让数值刚好等于"牌力强弱":于是比大小直接 `a.Rank > b.Rank`,不用查表。
这是 Go 常见的"用数值编码语义"技巧。

## 转变 3:方法在类型外面,靠"接收者"绑定

视觉上和 C#/TS 最不一样的地方:

```go
func (s Suit) String() string { ... }
//   ^^^^^^^^ 接收者(receiver) = Go 的 this/self,而且要显式写出、可自己命名
```

| | C# / TS | Go |
|---|---|---|
| 方法位置 | 写在 `class { }` **内部** | 写在类型**外部**,用接收者绑定 |
| this | 隐式 | 显式写出(这里叫 `s`) |
| 一个类型的方法 | 必须都在 class 体内 | 可散落在同包任意文件 |

好处:`type Suit uint8` 永远是干净的数据声明,行为单独挂;甚至能给同包里别人定义的类型加方法。

## 转变 4:隐式接口 —— 最大的脑筋急转弯

代码里三个 `String()` 方法,从没写过 `implements Stringer`,却能被 `fmt.Println` 自动识别。因为:

```go
type Stringer interface { String() string }   // 标准库 fmt 里的接口
```

> 规则:任何类型只要"碰巧"有 `String() string` 方法,就自动算实现了 `Stringer`——不声明、不继承。

| | C# / TS | Go |
|---|---|---|
| 声明实现 | `class Card : IStringer`(必须写) | 什么都不写,方法对上就算 |
| 心智模型 | "我**是**一个 IStringer" | "我**能**做 String() 这件事" |

`fmt.Println(card)` 内部逻辑是:"这玩意儿有 `String()` 吗?有就调它。"于是打印 `♠3` 而不是 `{0 3 0}`。
**Go 不问"你是什么类",只问"你会做什么"。** 接口是使用方定义的小契约,不是实现方挂的标签。

## 转变 5:`struct` 是纯数据袋,没有继承

```go
type Card struct {
    ID   int
    Rank Rank
    Suit Suit
}
```

`struct` 就是字段集合,纯数据。和 C# class 的区别:

- **没有继承**:Go 没有 `extends`,复用靠组合(把一个 struct 嵌进另一个),信条是"组合优于继承"。
- **数据与行为分离**:Card 里只有字段,`String()` 在外面。

## 回到 card.go:为什么这样组织

```text
package card        ① 组织单元:一个包 = 一个概念(牌),不是一个类一个文件
type Suit + 常量 + String()   ② 类型紧跟它的常量和方法
type Rank + 常量 + String()
var rankLabels                ③ 小写 = 包私有的辅助数据
type Card + String()          ④ 组合前两者的主类型
func NewDeck() []Card         ⑤ 工厂函数放最后
```

- **包取代类做组织单元**:`package card` 是"牌"这个概念的边界,调用方写 `card.NewDeck()`、`card.Spade`,包名就是命名空间。
- **大小写 = 可见性**:`Suit`/`Card`/`NewDeck` 大写导出(≈ public);`rankLabels` 小写仅包内(≈ private)。没有关键字,看首字母即知。
- **没有构造函数**:`Card{ID: 0, Rank: Three}` 这种结构体字面量就是原生创建方式;复杂创建写个 `NewXxx` 工厂函数(自由函数,不属于任何类型)。
- **返回 `[]Card` 而非自定义集合类**:Go 偏爱内置 slice/map。`make([]Card, 0, 54)` 的 `54` 是预分配容量,append 时不反复扩容。

## C#/TS ↔ Go 总对照表

| 概念 | C# / TypeScript | Go |
|---|---|---|
| 枚举 | `enum Suit {...}` | `type T uint8` + `const(...iota...)` |
| 类 | `class Card {...}` | `type Card struct {...}` |
| 方法 | 写在 class 体内 | 写在类型外,带接收者 |
| this | 隐式 | 显式接收者 `(c Card)` |
| 接口实现 | `implements` 显式声明 | 方法对上就隐式满足 |
| 继承 | `class B : A` | 没有,用组合/嵌入 |
| 可见性 | `public`/`private` | 首字母大小写 |
| 构造函数 | `new Card(...)` | 工厂函数 `NewXxx()` + 字面量 |
| 集合 | `List<T>`/自定义类 | 内置 slice `[]T` / `map` |
| 命名空间 | `namespace` | 包(目录) |

## 接下来会碰到的(本文件还没出现)

- **指针接收者 `func (c *Card)`**:现在的 `func (c Card)` 是值接收者,方法拿到的是拷贝、改不动原物(对只读的 `String()` 正好)。等写"出牌、改手牌"时需要 `*` 才改得到原物。
- **接口真正发力**:写牌型比较时定义小接口,体会隐式满足的威力。
- **组合(嵌入)**:替代继承复用代码。

## 记住三句话

1. Go 没有类:类型管数据、方法挂在外面、包做组织单元。
2. 接口隐式满足——有对应方法就算实现,不写 `implements`;问"会做什么"不问"是什么"。
3. 组合代替继承,工厂函数代替构造函数,大小写代替 `public/private`。
