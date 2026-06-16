#!/usr/bin/env python3
"""Build visual HTML learning notes from docs/learning Markdown files."""

from __future__ import annotations

import html
import re
from dataclasses import dataclass
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]
LEARNING_DIR = REPO_ROOT / "docs" / "learning"
ARCHIVE_SOURCE_DIR = REPO_ROOT / "docs" / "archive" / "learning-markdown-sources"
STYLE_PATH = LEARNING_DIR / "learning.css"


@dataclass(frozen=True)
class PageSpec:
    slug: str
    stage: str
    title: str
    subtitle: str
    analogy_title: str
    analogy: str
    visual_title: str
    visual_kind: str
    takeaways: tuple[str, str, str]

    @property
    def number(self) -> str:
        return self.slug.split("-", 1)[0]


PAGES: tuple[PageSpec, ...] = (
    PageSpec(
        "00-comfort-zone-quiz",
        "诊断",
        "Go 舒适区诊断",
        "用代码阅读题判断哪些知识还在舒适区，哪些需要重新焊牢。",
        "像体检，不像考试",
        "这篇不是为了证明会不会，而是像体检报告一样定位血压、视力和心肺。答错不是失败，反而能告诉后面的训练该从哪块肌肉开始练。",
        "舒适区 / 学习区 / 困难区",
        "zones",
        ("诚实标记“不知道”比猜对更有价值", "学习路线从真实薄弱点长出来", "并发、go mod、错误包装是近期重点"),
    ),
    PageSpec(
        "01-roadmap",
        "路线",
        "学习区路线图",
        "从 Go 工具链复健，到纯逻辑、连接层、房间状态机和 Cocos 联调。",
        "像健身房训练计划",
        "每个阶段只加一两块重量。太轻会无聊，太重会受伤；路线图的作用就是让训练一直停在“有点吃力但做得完”的学习区。",
        "从诊断到完整游戏的路径",
        "roadmap",
        ("先纯逻辑，再网络，再真实客户端", "每一步必须可运行、可观察", "Cocos 联调要早，但先保持最小"),
    ),
    PageSpec(
        "02-modern-go-toolchain",
        "Go",
        "现代 Go 工具链",
        "理解 go.mod、internal、测试命令和 Go 版本语义。",
        "像搬家后有了新门牌",
        "GOPATH 时代像必须住在指定小区；go mod 时代像你可以住任意地方，只要门牌号写在 go.mod 上，快递员就知道怎么送包裹。",
        "go mod 如何找到本地代码",
        "gomod",
        ("module 行是项目身份证", "internal 是编译器守门员", "go test ./... 是每天的体温计"),
    ),
    PageSpec(
        "03-due-container-lifecycle",
        "due 01",
        "Container 与生命周期",
        "理解 due 如何托管服务器组件，为什么 Serve 会让进程持续运行。",
        "像剧院舞台监督",
        "演员不是自己随便上下场，舞台监督统一安排开场、换景、谢幕和清场。Container 也是这样管理 Gate、Node 等组件的生命周期。",
        "Container 管理启动与关闭",
        "lifecycle",
        ("Container 是组件总管", "Serve 负责等待退出信号", "进程运行不等于端口可连接"),
    ),
    PageSpec(
        "04-software-engineering-map",
        "地图",
        "软件工程能力全景地图",
        "把编码、测试、Git、评审、CI、安全、运维等能力映射到本项目。",
        "像游戏大地图",
        "主线任务是做斗地主服务器，但地图上还有测试、版本控制、安全、可观测性这些支线。先点亮 P0 区域，别一开始就冲进高等级地图。",
        "P0 / P1 / P2 能力分层",
        "capability",
        ("P0 现在就练，P1 随阶段引入", "文档按 Diataxis 分清用途", "单人项目也要保留工程纪律"),
    ),
    PageSpec(
        "05-due-gate-websocket",
        "due 02",
        "Gate 与 WebSocket",
        "让 due 第一次真正监听端口，并区分 Gate、Server、Locator、Registry。",
        "像酒店大堂",
        "客人先到大堂办理进入，大堂知道房间和服务台在哪里，但它不会亲自做菜。Gate 管连接和转发，真正的游戏逻辑交给 Node。",
        "客户端进入集群的路径",
        "gate",
        ("Gate 是入口，不是业务规则", "WebSocket Server 才监听端口", "Redis/etcd 为多实例协作铺路"),
    ),
    PageSpec(
        "06-development-learning-workflow",
        "工作流",
        "纵向切片开发工作流",
        "先用 Go 客户端快速反馈，再接入最小 Cocos 联调。",
        "像切蛋糕，不像分层烤饼",
        "每次切下一小块都包含业务逻辑、通信和一点界面验证。这样很快能尝到味道，而不是先烤一大张后端饼再发现前端接不上。",
        "从纯逻辑到 Cocos 的纵向切片",
        "slice",
        ("不先写完整服务器", "Go 测试客户端是快速反馈层", "Python 脚本只编排开发命令"),
    ),
    PageSpec(
        "07-first-route-and-codec",
        "due 03",
        "第一条路由与编解码器",
        "跑通客户端 → Gate → Node → 响应，并理解 json/proto 必须一致。",
        "像两个人约定用同一种语言",
        "快递车 Gate 只负责搬箱子，不会拆开看信件。真正写信和读信的是客户端与 Node，所以双方必须约好是中文、英文还是 protobuf。",
        "请求响应链路与 codec 边界",
        "codec",
        ("Route 是请求入口", "Gate 只转发字节", "两端 codec 不一致就读不懂"),
    ),
    PageSpec(
        "08-route-vs-event",
        "due 04",
        "Route 与 Event",
        "解释 connect/disconnect 警告，以及路由和事件的根本区别。",
        "像人说话 vs 门铃响",
        "Route 是玩家主动说“我要做什么”；Event 是门铃、进门、离开这些自动发生的信号。两者都能被前台记录，但来源完全不同。",
        "两条平行通道",
        "route-event",
        ("WARN 无害，只是没人接事件", "Route 来自客户端主动消息", "Event 来自连接生命周期"),
    ),
    PageSpec(
        "09-login-and-binding",
        "due 05",
        "登录与连接绑定",
        "用 BindGate 把匿名连接绑定成可识别的玩家 uid。",
        "像给访客发工牌",
        "刚进楼的人只有临时门禁记录 cid；登录后领到工牌 uid，之后保安、会议室和通知系统都能知道这个人是谁。",
        "cid 到 uid 的身份升级",
        "binding",
        ("登录路由负责绑定身份", "断开日志能验证 uid 是否生效", "BindGate 会触发 Reconnect"),
    ),
    PageSpec(
        "10-server-push",
        "due 06",
        "服务端主动推送",
        "让服务端按 uid 主动给玩家发消息，为房间广播打基础。",
        "像外卖员主动敲门",
        "Response 是你打电话问客服，客服回答你；Push 是外卖到了，骑手不用等你再问，直接按门牌号敲门通知。",
        "从回复请求到主动通知",
        "push",
        ("Response 沿原请求返回", "Push 依赖 uid 绑定反向找人", "Multicast 就是房间广播雏形"),
    ),
)


def md_source_path(slug: str) -> Path:
    current = LEARNING_DIR / f"{slug}.md"
    if current.exists():
        return current
    return ARCHIVE_SOURCE_DIR / f"{slug}.md"


def inline(text: str) -> str:
    text = html.escape(text)
    text = re.sub(r"`([^`]+)`", r"<code>\1</code>", text)
    text = re.sub(r"\*\*([^*]+)\*\*", r"<strong>\1</strong>", text)

    def repl(match: re.Match[str]) -> str:
        label = match.group(1)
        url = match.group(2)
        if url.endswith(".md"):
            url = url[:-3] + ".html"
        return f'<a href="{html.escape(url)}">{label}</a>'

    return re.sub(r"\[([^\]]+)\]\(([^)]+)\)", repl, text)


def slugify(text: str) -> str:
    value = re.sub(r"[^\w\u4e00-\u9fff]+", "-", text.lower()).strip("-")
    return value or "section"


def parse_table(lines: list[str], start: int) -> tuple[str, int] | None:
    if start + 1 >= len(lines):
        return None
    if not lines[start].lstrip().startswith("|"):
        return None
    if not re.match(r"^\s*\|?[\s:\-|]+\|?\s*$", lines[start + 1]):
        return None

    rows: list[list[str]] = []
    i = start
    while i < len(lines) and lines[i].lstrip().startswith("|"):
        if i != start + 1:
            cells = [cell.strip() for cell in lines[i].strip().strip("|").split("|")]
            rows.append(cells)
        i += 1

    if not rows:
        return None

    header = rows[0]
    body = rows[1:]
    parts = ["<div class=\"table-wrap\"><table>", "<thead><tr>"]
    parts.extend(f"<th>{inline(cell)}</th>" for cell in header)
    parts.append("</tr></thead>")
    if body:
        parts.append("<tbody>")
        for row in body:
            parts.append("<tr>")
            parts.extend(f"<td>{inline(cell)}</td>" for cell in row)
            parts.append("</tr>")
        parts.append("</tbody>")
    parts.append("</table></div>")
    return "\n".join(parts), i


def markdown_to_html(markdown: str) -> tuple[str, list[tuple[int, str, str]]]:
    lines = markdown.splitlines()
    parts: list[str] = []
    headings: list[tuple[int, str, str]] = []
    paragraph: list[str] = []
    list_items: list[str] = []
    ordered_items: list[str] = []
    quote: list[str] = []
    in_code = False
    code_lang = ""
    code_lines: list[str] = []
    i = 0

    def flush_paragraph() -> None:
        nonlocal paragraph
        if paragraph:
            parts.append(f"<p>{inline(' '.join(paragraph))}</p>")
            paragraph = []

    def flush_list() -> None:
        nonlocal list_items, ordered_items
        if list_items:
            parts.append("<ul>" + "".join(f"<li>{item}</li>" for item in list_items) + "</ul>")
            list_items = []
        if ordered_items:
            parts.append("<ol>" + "".join(f"<li>{item}</li>" for item in ordered_items) + "</ol>")
            ordered_items = []

    def flush_quote() -> None:
        nonlocal quote
        if quote:
            parts.append("<blockquote>" + "".join(f"<p>{inline(q)}</p>" for q in quote) + "</blockquote>")
            quote = []

    while i < len(lines):
        line = lines[i]

        if in_code:
            if line.startswith("```"):
                escaped = html.escape("\n".join(code_lines))
                parts.append(f'<pre><code class="language-{html.escape(code_lang)}">{escaped}</code></pre>')
                in_code = False
                code_lang = ""
                code_lines = []
            else:
                code_lines.append(line)
            i += 1
            continue

        if line.startswith("```"):
            flush_paragraph()
            flush_list()
            flush_quote()
            in_code = True
            code_lang = line.strip("`").strip()
            code_lines = []
            i += 1
            continue

        table = parse_table(lines, i)
        if table:
            flush_paragraph()
            flush_list()
            flush_quote()
            parts.append(table[0])
            i = table[1]
            continue

        if not line.strip():
            flush_paragraph()
            flush_list()
            flush_quote()
            i += 1
            continue

        heading = re.match(r"^(#{1,4})\s+(.+)$", line)
        if heading:
            flush_paragraph()
            flush_list()
            flush_quote()
            level = len(heading.group(1))
            text = heading.group(2).strip()
            ident = slugify(text)
            headings.append((level, text, ident))
            parts.append(f'<h{level} id="{ident}">{inline(text)}</h{level}>')
            i += 1
            continue

        if line.startswith(">"):
            flush_paragraph()
            flush_list()
            quote.append(line.lstrip("> ").strip())
            i += 1
            continue

        unordered = re.match(r"^\s*[-*]\s+(.+)$", line)
        if unordered:
            flush_paragraph()
            flush_quote()
            list_items.append(inline(unordered.group(1).strip()))
            i += 1
            continue

        ordered = re.match(r"^\s*\d+\.\s+(.+)$", line)
        if ordered:
            flush_paragraph()
            flush_quote()
            ordered_items.append(inline(ordered.group(1).strip()))
            i += 1
            continue

        flush_list()
        flush_quote()
        paragraph.append(line.strip())
        i += 1

    flush_paragraph()
    flush_list()
    flush_quote()
    return "\n".join(parts), headings


def diagram(kind: str) -> str:
    if kind == "zones":
        return """
<svg viewBox="0 0 760 260" role="img" aria-label="舒适区、学习区、困难区">
  <circle cx="250" cy="130" r="70" class="svg-zone zone-safe"/>
  <circle cx="250" cy="130" r="115" class="svg-zone zone-learn"/>
  <circle cx="250" cy="130" r="158" class="svg-zone zone-hard"/>
  <text x="250" y="126" text-anchor="middle" class="svg-title">舒适区</text>
  <text x="250" y="66" text-anchor="middle" class="svg-label">学习区</text>
  <text x="250" y="18" text-anchor="middle" class="svg-label">困难区</text>
  <path d="M435 78 H660" class="svg-line"/>
  <path d="M435 130 H660" class="svg-line"/>
  <path d="M435 182 H660" class="svg-line"/>
  <text x="675" y="84" class="svg-small">会读语法</text>
  <text x="675" y="136" class="svg-small">重新焊牢细节</text>
  <text x="675" y="188" class="svg-small">暂缓分布式和性能调优</text>
</svg>"""
    if kind == "roadmap":
        labels = ["工具链", "牌与牌型", "连接层", "状态机", "Cocos", "进阶"]
        circles = []
        for idx, label in enumerate(labels):
            x = 80 + idx * 120
            circles.append(f'<circle cx="{x}" cy="110" r="34" class="svg-node"/><text x="{x}" y="116" text-anchor="middle" class="svg-small">{label}</text>')
            if idx < len(labels) - 1:
                circles.append(f'<path d="M{x+36} 110 H{x+84}" class="svg-arrow"/>')
        return f'<svg viewBox="0 0 760 220" role="img" aria-label="学习路线">{"".join(circles)}<text x="380" y="178" text-anchor="middle" class="svg-label">每一步只引入 1-2 个新概念</text></svg>'
    if kind == "gomod":
        return """
<svg viewBox="0 0 760 260" role="img" aria-label="go mod 导入解析">
  <rect x="40" y="45" width="250" height="80" rx="18" class="svg-box"/><text x="165" y="78" text-anchor="middle" class="svg-title">go.mod</text><text x="165" y="106" text-anchor="middle" class="svg-small">module github.com/.../server</text>
  <rect x="470" y="45" width="250" height="80" rx="18" class="svg-box"/><text x="595" y="78" text-anchor="middle" class="svg-title">import</text><text x="595" y="106" text-anchor="middle" class="svg-small">.../server/internal/card</text>
  <path d="M292 85 H468" class="svg-arrow"/><text x="380" y="72" text-anchor="middle" class="svg-small">前缀匹配</text>
  <rect x="255" y="165" width="250" height="58" rx="16" class="svg-accent"/><text x="380" y="200" text-anchor="middle" class="svg-title">模块根目录下找代码</text>
</svg>"""
    if kind == "lifecycle":
        labels = ["Init", "Start", "等待信号", "Close", "Destroy"]
        return flow_svg(labels, "Container 统一编排组件生命周期")
    if kind == "capability":
        return """
<svg viewBox="0 0 760 300" role="img" aria-label="软件工程能力分层">
  <rect x="80" y="190" width="600" height="58" rx="14" class="svg-accent"/><text x="380" y="226" text-anchor="middle" class="svg-title">P0：编码 / 测试 / Git / 文档 / 评审</text>
  <rect x="135" y="120" width="490" height="58" rx="14" class="svg-box"/><text x="380" y="156" text-anchor="middle" class="svg-title">P1：CI / ADR / 可观测性</text>
  <rect x="210" y="50" width="340" height="58" rx="14" class="svg-box-muted"/><text x="380" y="86" text-anchor="middle" class="svg-title">P2：安全深化 / 部署 / 效能度量</text>
</svg>"""
    if kind == "gate":
        return layered_svg(("客户端", "Gate", "Node"), ("WebSocket", "Locator / Registry", "游戏逻辑"))
    if kind == "slice":
        return flow_svg(["纯逻辑", "单元测试", "due 路由", "Go 客户端", "Cocos 最小界面"], "纵向切片：每次都能端到端观察")
    if kind == "codec":
        return layered_svg(("客户端", "Gate", "Node"), ("编码 json", "只搬运字节", "解码 json"))
    if kind == "route-event":
        return """
<svg viewBox="0 0 760 260" role="img" aria-label="Route 和 Event">
  <rect x="70" y="55" width="260" height="140" rx="20" class="svg-box"/><text x="200" y="92" text-anchor="middle" class="svg-title">Route</text><text x="200" y="126" text-anchor="middle" class="svg-small">客户端主动发消息</text><text x="200" y="156" text-anchor="middle" class="svg-small">Greet / Login / PlayCard</text>
  <rect x="430" y="55" width="260" height="140" rx="20" class="svg-accent"/><text x="560" y="92" text-anchor="middle" class="svg-title">Event</text><text x="560" y="126" text-anchor="middle" class="svg-small">框架自动触发</text><text x="560" y="156" text-anchor="middle" class="svg-small">connect / disconnect</text>
</svg>"""
    if kind == "binding":
        return flow_svg(["cid=1 匿名", "Login 路由", "BindGate(uid)", "uid=1001 已识别"], "连接从临时号码升级为玩家身份")
    if kind == "push":
        return """
<svg viewBox="0 0 760 260" role="img" aria-label="Response 和 Push">
  <path d="M110 90 H640" class="svg-arrow"/><text x="380" y="72" text-anchor="middle" class="svg-small">Request</text>
  <path d="M640 125 H110" class="svg-arrow"/><text x="380" y="153" text-anchor="middle" class="svg-small">Response：沿原请求返回</text>
  <path d="M640 205 H110" class="svg-arrow accent-arrow"/><text x="380" y="230" text-anchor="middle" class="svg-small">Push：服务端按 uid 主动通知</text>
  <circle cx="100" cy="110" r="42" class="svg-node"/><text x="100" y="116" text-anchor="middle" class="svg-title">客户端</text>
  <circle cx="650" cy="110" r="42" class="svg-node"/><text x="650" y="116" text-anchor="middle" class="svg-title">服务端</text>
</svg>"""
    return flow_svg(["概念", "运行", "观察", "复盘"], "学习闭环")


def flow_svg(labels: list[str], caption: str) -> str:
    parts = ['<svg viewBox="0 0 760 230" role="img" aria-label="流程图">']
    step = 620 / max(1, len(labels) - 1)
    for idx, label in enumerate(labels):
        x = 70 + idx * step
        parts.append(f'<circle cx="{x}" cy="105" r="36" class="svg-node"/><text x="{x}" y="111" text-anchor="middle" class="svg-small">{html.escape(label)}</text>')
        if idx < len(labels) - 1:
            parts.append(f'<path d="M{x+38} 105 H{x+step-38}" class="svg-arrow"/>')
    parts.append(f'<text x="380" y="178" text-anchor="middle" class="svg-label">{html.escape(caption)}</text></svg>')
    return "".join(parts)


def layered_svg(top: tuple[str, str, str], bottom: tuple[str, str, str]) -> str:
    xs = (150, 380, 610)
    parts = ['<svg viewBox="0 0 760 250" role="img" aria-label="分层结构图">']
    for idx, x in enumerate(xs):
        klass = "svg-accent" if idx == 1 else "svg-box"
        parts.append(f'<rect x="{x-95}" y="60" width="190" height="100" rx="18" class="{klass}"/>')
        parts.append(f'<text x="{x}" y="98" text-anchor="middle" class="svg-title">{html.escape(top[idx])}</text>')
        parts.append(f'<text x="{x}" y="130" text-anchor="middle" class="svg-small">{html.escape(bottom[idx])}</text>')
        if idx < 2:
            parts.append(f'<path d="M{x+98} 110 H{xs[idx+1]-98}" class="svg-arrow"/>')
    parts.append("</svg>")
    return "".join(parts)


def page_html(spec: PageSpec, content: str, headings: list[tuple[int, str, str]]) -> str:
    content = content.replace("docs/learning/02-modern-go-toolchain.md", "docs/learning/02-modern-go-toolchain.html")
    nav = "".join(
        f'<a class="nav-link {"active" if page.slug == spec.slug else ""}" href="{page.slug}.html"><span>{page.number} · {page.stage}</span>{page.number} {page.title}</a>'
        for page in PAGES
    )
    toc = "".join(
        f'<a href="#{ident}" class="toc-level-{level}">{html.escape(text)}</a>'
        for level, text, ident in headings
        if level in (2, 3)
    )
    cards = "".join(f"<li>{html.escape(item)}</li>" for item in spec.takeaways)
    return f"""<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{html.escape(spec.title)} | go-doudizhu-lab</title>
  <link rel="stylesheet" href="learning.css">
</head>
<body>
  <div class="shell">
    <aside class="sidebar">
      <a class="brand" href="index.html">斗地主学习笔记</a>
      <nav>{nav}</nav>
    </aside>
    <main>
      <section class="hero">
        <p class="eyebrow">{html.escape(spec.stage)}</p>
        <h1>{html.escape(spec.title)}</h1>
        <p>{html.escape(spec.subtitle)}</p>
      </section>
      <section class="visual-grid">
        <article class="panel analogy">
          <span class="panel-label">生活类比</span>
          <h2>{html.escape(spec.analogy_title)}</h2>
          <p>{html.escape(spec.analogy)}</p>
        </article>
        <article class="panel diagram">
          <span class="panel-label">图形化理解</span>
          <h2>{html.escape(spec.visual_title)}</h2>
          {diagram(spec.visual_kind)}
        </article>
        <article class="panel takeaways">
          <span class="panel-label">抓手</span>
          <h2>读完要带走什么</h2>
          <ul>{cards}</ul>
        </article>
      </section>
      <div class="content-layout">
        <aside class="toc"><strong>本页目录</strong>{toc}</aside>
        <article class="content">
          {content}
        </article>
      </div>
    </main>
  </div>
</body>
</html>
"""


def index_html() -> str:
    cards = "".join(
        f"""<a class="index-card" href="{page.slug}.html">
  <span>{html.escape(page.number)} · {html.escape(page.stage)}</span>
  <h2>{html.escape(page.number)} {html.escape(page.title)}</h2>
  <p>{html.escape(page.subtitle)}</p>
</a>"""
        for page in PAGES
    )
    return f"""<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>斗地主学习笔记 | go-doudizhu-lab</title>
  <link rel="stylesheet" href="learning.css">
</head>
<body>
  <main class="index">
    <section class="index-hero">
      <p class="eyebrow">Learning Notes</p>
      <h1>斗地主全栈游戏开发学习地图</h1>
      <p>这些 HTML 页面由原 Markdown 笔记重构而来：保留文字细节，增加图形化关系、阶段导航和生活类比，方便以后边写代码边复盘。</p>
      {diagram("roadmap")}
    </section>
    <section class="index-grid">{cards}</section>
  </main>
</body>
</html>
"""


CSS = r"""
:root {
  color-scheme: light;
  --bg: #f7f2e8;
  --paper: #fffaf0;
  --ink: #231f20;
  --muted: #6f6258;
  --line: #e2d4c1;
  --accent: #f4b942;
  --accent-2: #f26b4f;
  --blue: #4c7df0;
  --green: #5abf90;
  --shadow: 0 24px 70px rgba(75, 55, 30, 0.14);
}

* { box-sizing: border-box; }

body {
  margin: 0;
  background:
    radial-gradient(circle at top left, rgba(244, 185, 66, .22), transparent 34rem),
    linear-gradient(135deg, #fbf7ef 0%, var(--bg) 100%);
  color: var(--ink);
  font-family: "Inter", "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif;
  line-height: 1.72;
}

a { color: #2357c6; text-decoration: none; }
a:hover { text-decoration: underline; }

.shell {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  min-height: 100vh;
}

.sidebar {
  position: sticky;
  top: 0;
  height: 100vh;
  padding: 26px 20px;
  border-right: 1px solid var(--line);
  background: rgba(255, 250, 240, .78);
  backdrop-filter: blur(14px);
  overflow: auto;
}

.brand {
  display: block;
  margin-bottom: 20px;
  color: var(--ink);
  font-weight: 900;
  font-size: 20px;
}

.nav-link {
  display: block;
  padding: 11px 12px;
  margin: 5px 0;
  border-radius: 14px;
  color: var(--muted);
}

.nav-link span {
  display: block;
  color: var(--accent-2);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: .06em;
}

.nav-link.active,
.nav-link:hover {
  background: #fff;
  box-shadow: 0 8px 28px rgba(75, 55, 30, .08);
  color: var(--ink);
  text-decoration: none;
}

main { padding: 36px min(5vw, 68px) 80px; }

.hero, .index-hero {
  max-width: 1100px;
  padding: 42px;
  border: 1px solid var(--line);
  border-radius: 30px;
  background: rgba(255, 250, 240, .88);
  box-shadow: var(--shadow);
}

.hero h1, .index-hero h1 {
  margin: 0 0 14px;
  font-size: clamp(34px, 5vw, 64px);
  line-height: 1.05;
  letter-spacing: -0.04em;
}

.hero p, .index-hero p {
  max-width: 760px;
  margin: 0;
  color: var(--muted);
  font-size: 18px;
}

.eyebrow {
  margin: 0 0 10px !important;
  color: var(--accent-2) !important;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: .12em;
  font-size: 13px !important;
}

.visual-grid {
  max-width: 1100px;
  display: grid;
  grid-template-columns: 1fr 1.4fr;
  gap: 18px;
  margin: 24px 0;
}

.panel {
  border: 1px solid var(--line);
  border-radius: 24px;
  background: rgba(255, 255, 255, .72);
  padding: 24px;
  box-shadow: 0 16px 44px rgba(75, 55, 30, .09);
}

.diagram { grid-row: span 2; }
.panel h2 { margin: 6px 0 12px; font-size: 23px; }
.panel p { color: var(--muted); margin-bottom: 0; }
.panel-label { color: var(--accent-2); font-size: 12px; font-weight: 900; letter-spacing: .1em; }
.takeaways ul { margin-bottom: 0; padding-left: 22px; }

.content-layout {
  max-width: 1100px;
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  gap: 28px;
  align-items: start;
}

.toc {
  position: sticky;
  top: 24px;
  padding: 18px;
  border: 1px solid var(--line);
  border-radius: 20px;
  background: rgba(255, 250, 240, .84);
  font-size: 14px;
}

.toc a {
  display: block;
  margin: 8px 0;
  color: var(--muted);
}
.toc-level-3 { padding-left: 12px; }

.content {
  padding: 34px;
  border: 1px solid var(--line);
  border-radius: 26px;
  background: var(--paper);
  box-shadow: var(--shadow);
}

.content h1 { display: none; }
.content h2 { margin-top: 42px; padding-top: 16px; border-top: 1px solid var(--line); font-size: 30px; }
.content h3 { margin-top: 30px; font-size: 22px; }
.content blockquote {
  margin: 22px 0;
  padding: 16px 20px;
  border-left: 5px solid var(--accent);
  background: #fff4d7;
  border-radius: 14px;
  color: var(--muted);
}

code {
  padding: 2px 6px;
  border-radius: 7px;
  background: #f0e5d6;
  color: #2c2a28;
}

pre {
  overflow: auto;
  padding: 18px;
  border-radius: 18px;
  background: #1f2430;
  color: #f8f2e8;
  line-height: 1.55;
}
pre code { background: transparent; color: inherit; padding: 0; }

.table-wrap { overflow-x: auto; margin: 20px 0; }
table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 16px; overflow: hidden; }
th, td { padding: 12px 14px; border-bottom: 1px solid var(--line); vertical-align: top; }
th { background: #f4ead9; text-align: left; }

svg { width: 100%; height: auto; }
.svg-box { fill: #fff7e8; stroke: #dfc9ac; stroke-width: 2; }
.svg-box-muted { fill: #efe7dc; stroke: #d7c7b1; stroke-width: 2; }
.svg-accent { fill: #ffe4a3; stroke: #e6ac29; stroke-width: 2; }
.svg-node { fill: #fffaf0; stroke: #f26b4f; stroke-width: 3; }
.svg-zone { fill: none; stroke-width: 34; opacity: .55; }
.zone-safe { stroke: var(--green); }
.zone-learn { stroke: var(--accent); }
.zone-hard { stroke: var(--accent-2); }
.svg-arrow { stroke: #38312d; stroke-width: 3; fill: none; }
.accent-arrow { stroke: var(--accent-2); }
.svg-line { stroke: #b9a48d; stroke-width: 2; }
.svg-title { font-size: 18px; font-weight: 900; fill: var(--ink); }
.svg-label { font-size: 16px; font-weight: 800; fill: var(--muted); }
.svg-small { font-size: 14px; fill: var(--muted); font-weight: 700; }

.index {
  max-width: 1180px;
  margin: 0 auto;
  padding: 48px 24px 80px;
}
.index-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 18px;
  margin-top: 24px;
}
.index-card {
  display: block;
  min-height: 210px;
  padding: 24px;
  border: 1px solid var(--line);
  border-radius: 24px;
  background: rgba(255, 250, 240, .86);
  color: var(--ink);
  box-shadow: 0 18px 48px rgba(75, 55, 30, .1);
}
.index-card:hover { transform: translateY(-3px); text-decoration: none; transition: .18s ease; }
.index-card span { color: var(--accent-2); font-weight: 900; font-size: 12px; letter-spacing: .1em; }
.index-card h2 { margin: 12px 0 8px; }
.index-card p { color: var(--muted); margin: 0; }

@media (max-width: 980px) {
  .shell { grid-template-columns: 1fr; }
  .sidebar { position: static; height: auto; }
  .visual-grid, .content-layout { grid-template-columns: 1fr; }
  .toc { position: static; }
  main { padding: 22px; }
}
"""


def write_files() -> None:
    LEARNING_DIR.mkdir(parents=True, exist_ok=True)
    STYLE_PATH.write_text(CSS, encoding="utf-8", newline="\n")
    (LEARNING_DIR / "index.html").write_text(index_html(), encoding="utf-8", newline="\n")
    for spec in PAGES:
        source = md_source_path(spec.slug)
        markdown = source.read_text(encoding="utf-8")
        content, headings = markdown_to_html(markdown)
        (LEARNING_DIR / f"{spec.slug}.html").write_text(
            page_html(spec, content, headings),
            encoding="utf-8",
            newline="\n",
        )


if __name__ == "__main__":
    write_files()
