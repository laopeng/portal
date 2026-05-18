# portal

## 关于 gstack

本项目使用 **gstack** —— Garry Tan（Y Combinator CEO）的 AI 辅助开发工具栈。
gstack 将 Claude Code 变成一个虚拟工程团队：CEO、设计、工程经理、代码审查、QA 测试、安全审计、发布管理，全部通过 slash command 调用。

gstack 的核心哲学：
- **Boil the Lake（烧干湖水）**：AI 让完整性的边际成本接近零。90% vs 100% 只差 70 行代码，时间是几秒钟。永远做完整版。
- **Search Before Building（先搜后造）**：1000x 工程师的第一反应是"有人解决过吗？"而不是"我来从零设计"。
- **User Sovereignty（用户主权）**：AI 建议，用户决定。两个模型达成共识不等于必须执行。

---

## sprint 全流程

gstack 按真实 sprint 顺序组织：**Think → Plan → Build → Review → Test → Ship → Reflect**

每个阶段的输出自动传递到下一阶段。`/office-hours` 写出设计文档，`/plan-ceo-review` 读取它，`/plan-eng-review` 写出测试计划，`/qa` 读取它，`/review` 找到的 bug 会被 `/ship` 验证修复。不会遗漏任何东西，因为每一步都知道上一步发生了什么。

### 阶段一：Think（思考产品）

在使用 gstack 构建功能之前，先通过思考技能明确你要构建什么：

| 命令 | 角色 | 功能 |
|------|------|------|
| `/office-hours` | **YC 创业办公室** | 所有事情的起点。6 个强迫性问题重新框定你的产品。在写代码之前挑战你的假设，创建实现方案，生成设计文档。输出会被下游所有技能消费。 |
| `/learn` | **记忆管理** | 查看、搜索、清理 gstack 跨会话学到的经验。项目特定的模式、陷阱和偏好随着时间积累，让 gstack 对你的代码库越来越聪明。 |

### 阶段二：Plan（制定方案）

在写任何代码之前，执行规划评审，确保构建的是正确的东西：

| 命令 | 角色 | 功能 |
|------|------|------|
| `/plan-ceo-review` | **CEO / 创始人** | 重新思考问题。4 种模式：Expand（扩展至 10 星产品）、Selective Expand（选择性扩展）、Hold Scope（保持范围）、Reduce（裁剪）。从战略层面挑战需求。 |
| `/plan-eng-review` | **工程经理** | 锁定架构、数据流、ASCII 图表、边界情况和测试计划。强迫隐藏的假设暴露出来。输出测试矩阵、失败模式和安全隐患。 |
| `/plan-design-review` | **高级设计师** | 以交互方式评审设计的每个维度，0-10 评分，解释 10 分长什么样，然后编辑方案以达成目标。AI Slop 检测。每项设计决策用一次 AskUserQuestion 确认。 |
| `/plan-devex-review` | **开发体验负责人** | 交互式 DX 评审：探索开发者画像，对比竞品 TTHW（Time to Hello World），设计魔法时刻，逐步追踪摩擦点。3 种模式：DX EXPANSION / DX POLISH / DX TRIAGE，20-45 个强迫性问题。 |
| `/autoplan` | **自动评审管线** | 一条命令跑完完整评审：CEO → 设计 → 工程，自动检测哪些评审适用。只用你批准品味决策，其余自动编码。 |

### 阶段三：Build（构建）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/design-consultation` | **设计合伙人** | 从零构建设计系统。研究行业格局，提出创意风险，生成逼真的产品 mockup。 |
| `/design-shotgun` | **设计探索者** | "给我看看选项。"生成 4-6 个 AI mockup 变体，在浏览器中打开比较面板，收集反馈，迭代。品味记忆学习你的喜好。反复直到满意，然后交给 `/design-html`。 |
| `/design-html` | **设计工程师** | 将 mockup 转为生产级 HTML/CSS。Pretext 计算布局：文字自动换行、高度适应内容、布局动态。30KB 零依赖。检测 React/Svelte/Vue 框架，智能 API 路由。输出可交付，不是 demo。 |
| `/pair-agent` | **多 AI 协作** | 与任何 AI agent（OpenClaw、Hermes、Codex、Cursor）共享浏览器。每个 agent 独立 tab，scoped token，速率限制，活动归因。 |

### 阶段四：Review（代码评审）

代码写完后，在提交 PR 之前运行评审：

| 命令 | 角色 | 功能 |
|------|------|------|
| `/review` | **高级工程师** | 找到通过 CI 但在生产环境中爆炸的 bug。自动修复能确定的，标记完整性缺口。运行检查清单：并发安全、SQL 注入、错误处理、边界测试。 |
| `/design-review` | **会写代码的设计师** | 实时站点的视觉审计 + 修复循环。80 项审计清单，审计结束后修复发现的问题。原子提交，修前/修后截图。 |
| `/devex-review` | **DX 测试者** | 实时开发者体验审计。实际测试你的入门流程：浏览文档、走通 Getting Started、计时 TTHW、截图错误。对比 `/plan-devex-review` 的分数——验证计划是否匹配现实。 |
| `/codex` | **第二意见** | OpenAI Codex CLI 独立评审。3 种模式：review（通过/失败关卡）、adversarial challenge（尝试破坏代码）、open consultation（开放咨询及会话连续性）。当 `/review` 和 `/codex` 都评审过后，生成跨模型分析。 |
| `/cso` | **首席安全官** | OWASP Top 10 + STRIDE 威胁建模。零噪音：17 个误报排除、8/10+ 置信度门槛、独立发现验证。每个发现附带具体攻击场景。 |

**该用哪个评审？**

| 为谁构建 | 方案阶段（写代码前） | 实时审计（发布后） |
|---------|---------------------|--------------------|
| 最终用户（UI、Web App、移动端） | `/plan-design-review` | `/design-review` |
| 开发者（API、CLI、SDK、文档） | `/plan-devex-review` | `/devex-review` |
| 架构（数据流、性能、测试） | `/plan-eng-review` | `/review` |
| 全部 | `/autoplan`（自动检测适用项） | — |

### 阶段五：Test（测试）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/qa` | **QA 负责人** | 测试应用，找 bug，原子提交修复，重新验证。为每个修复自动生成回归测试。 |
| `/qa-only` | **QA 报告者** | 与 `/qa` 相同方法论，但只报告不修复。纯 bug 报告，不改代码。 |
| `/benchmark` | **性能工程师** | 基线页面加载时间、Core Web Vitals、资源大小。每个 PR 对比前后。 |
| `/benchmark-models` | **跨模型基准测试** | 对 Claude、GPT（通过 Codex CLI）、Gemini 运行相同 prompt，比较延迟、token、成本及质量评分。 |

### 阶段六：Ship（发布）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/ship` | **发布工程师** | 同步主分支、运行测试、审计覆盖率、推送、打开 PR。如果没有测试框架则自动搭建。过滤压缩 WIP 提交。 |
| `/land-and-deploy` | **发布工程师** | 合并 PR、等待 CI 和部署、验证生产健康状态。"approved" 到 "verified in production" 一条命令搞定。 |
| `/canary` | **SRE** | 部署后监控循环。观察控制台错误、性能回归、页面故障。 |
| `/setup-deploy` | **部署配置器** | `/land-and-deploy` 的一次性配置。自动检测平台、生产 URL、部署命令。 |
| `/document-release` | **技术文档编写者** | 更新所有项目文档以匹配刚发布的内容。自动检测陈旧的 README。构建 Diataxis 覆盖地图（reference / how-to / tutorial / explanation），让缺口在 PR body 中可见。 |
| `/document-generate` | **文档作者** | 使用 Diataxis 框架从头生成缺失的文档。先研究代码库，然后编写匹配代码的文档。可独立调用或在 `/document-release` 发现覆盖缺口时链式触发。 |

### 阶段七：Reflect（复盘）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/retro` | **工程经理** | 团队感知的周度回顾。每人分解、交付连续天数、测试健康趋势、成长机会。`/retro global` 跨所有项目和 AI 工具运行。 |
| `/investigate` | **调试者** | 系统化根因调试。铁律：没有调查就没有修复。追踪数据流、测试假设、3 次失败修复后停止。自动 freeze 到被调查的模块。 |

---

## 安全与约束工具

| 命令 | 功能 |
|------|------|
| `/careful` | 安全护栏。在执行破坏性命令（rm -rf、DROP TABLE、force-push、git reset --hard）前警告。说 "be careful" 激活。可覆盖任何警告。 |
| `/freeze` | 编辑锁定。限制文件编辑到单一目录，防止调试时意外修改范围外的代码。 |
| `/guard` | 完整安全。`/careful` + `/freeze` 合并。为生产环境提供最大安全。 |
| `/unfreeze` | 解锁。移除 `/freeze` 边界。 |

---

## 浏览器工具

| 命令 | 功能 |
|------|------|
| `/browse` | 给 AI 装上眼睛。真实 Chromium 浏览器，真实点击，真实截图。~100ms 每次调用。约 70 个命令：导航、截图、快照、点击、填表、JS 执行、网络捕获等。 |
| `/open-gstack-browser` | 启动 GStack 浏览器，配置侧边栏、反机器人隐身、自动模型路由（Sonnet 执行动作、Opus 分析），一键 cookie 导入。 |
| `/setup-browser-cookies` | 从真实浏览器（Chrome、Arc、Brave、Edge）导入 cookies 到 headless 会话，用于测试需要登录的页面。 |
| `/scrape <意图>` | 只读页面数据提取入口。3 条路径：匹配路径（~200ms 通过已编码的 browser-skill）、原型路径（~30s 通过 $B 命令驱动页面）、写入意图拒绝（/scrape 为只读）。 |
| `/skillify` | 将最近成功的 `/scrape` 原型编码为永久 browser-skill。11 步流程，3 份锁定合约：出处守卫、合成输入切片、原子写入。 |

---

## 浏览器交互说明

**必须使用 `/browse` 或 `$B <command>` 与浏览器交互。永远不要使用 `mcp__claude-in-chrome__*` 工具**——它们慢、不可靠且不是本工具链的标准方案。

### 架构概览

```
Claude Code                     gstack
─────────                      ──────
                               ┌──────────────────────┐
  Tool call: $B snapshot -i    │  CLI (编译后的二进制)   │
  ─────────────────────────→   │  • 读取 state 文件     │
                               │  • POST /command      │
                               │    到 localhost:PORT   │
                               └──────────┬───────────┘
                                          │ HTTP
                               ┌──────────▼───────────┐
                               │  Server (Bun.serve)   │
                               │  • 分发命令            │
                               │  • 与 Chromium 通信    │
                               │  • 返回纯文本          │
                               └──────────┬───────────┘
                                          │ CDP
                               ┌──────────▼───────────┐
                               │  Chromium (headless)  │
                               │  • 持久 tab            │
                               │  • cookies 跨命令保持  │
                               │  • 30min 空闲超时      │
                               └───────────────────────┘
```

- 首次调用启动 daemon + Chromium（~3s），后续调用 ~100-200ms
- 每个项目有独立的 daemon、端口、状态文件和 cookie（通过 `git rev-parse --show-toplevel` 检测项目根目录）
- 状态文件在 `<project>/.gstack/browse.json`（chmod 600）

---

## 其他工具

| 命令 | 功能 |
|------|------|
| `/setup-gbrain` | GBrain 上手。5 分钟内从零到运行。PGLite 本地、Supabase 已有 URL 或自动配置 Supabase 项目。 |
| `/sync-gbrain` | 保持 brain 最新。重新索引代码库，刷新 CLAUDE.md 中的 GBrain 搜索指南。 |
| `/gstack-upgrade` | 自升级。检测全局 vs vendored 安装，同步两者，显示变更内容。 |
| `/skillify` | 为重复性任务生成新的 gstack skill 模板。 |
| `/context-save` / `/context-restore` | 设置连续断点模式后（`gstack-config set checkpoint_mode continuous`），自动提交 WIP 保存上下文。`/context-restore` 读取 WIP 提交重建会话状态。 |
| `/landing-report` | 生成落地页报告。 |
| `/plan-tune` | 微调已有方案。 |
| `/make-pdf` | 将页面或 Markdown 转为生产质量的 PDF。 |
| `/health` | 技能健康面板。 |
| `/codex` | 第二意见：OpenAI Codex CLI 独立代码评审。 |
| `/benchmark-models` | 跨模型基准测试。 |

---

## 建立并行 sprint

gstack 在一个 sprint 中有效，在同时运行 10 个 sprint 时显现威力：

1. **先 `/office-hours`**——描述你要构建什么，让 Claude 挑战你的假设
2. **再 `/plan-ceo-review`**——重新思考问题，找到隐藏在产品需求背后的 10 星版本
3. **你批准方案，退出 Plan Mode**
4. **AI 生成代码**——2,400 行，11 文件，~8 分钟
5. **`/review`**——找 bug、自动修复、标记问题
6. **`/qa https://staging.xxx.com`**——打开真浏览器，点击流程，发现 bug 并修复
7. **`/ship`**——测试、PR、发布

设计管线可以并行运行：
- `/design-consultation` 构建设计系统 → 写入 DESIGN.md
- `/design-shotgun` 生成 4-6 个 mockup → 在浏览器中比较 → 迭代 → 直到满意
- `/design-html` 将批准的设计转成生产 HTML/CSS

---

## 安全架构（侧边栏 Agent）

GStack 浏览器侧边栏内置了多层 prompt injection 防御：

| 层级 | 模块 | 说明 |
|------|------|------|
| L1–L3 | `content-security.ts` | 数据标记、隐藏元素剥离、ARIA 正则、URL 拦截列表、信封包装 |
| L4 | `security-classifier.ts` (TestSavantAI ONNX, 112MB) | 侧边栏 agent 专用 ML 分类器 |
| L4b | `security-classifier.ts` (Claude Haiku 文本) | 侧边栏 agent 专用文本检查 |
| L5 | `security.ts` (canary) | 在系统 prompt 中注入随机 canary token，捕获跨文本、工具参数、URL 和文件写入的会话窃取 |
| L6 | `security.ts` (集成裁决) | 要求两个分类器一致同意才阻止 |

**环境变量：**
- `GSTACK_SECURITY_OFF=1` — 紧急开关
- `GSTACK_SECURITY_ENSEMBLE=deberta` — 选择启用 DeBERTa-v3 集成（721MB 下载），提供 2-of-3 模型协议

---

## 开发说明

如果你需要修改 gstack 本身：

```bash
# gstack 开发命令
bun install          # 安装依赖
bun test             # 免费测试（browse + snapshot + skill 验证，<2s）
bun run test:evals   # 付费 LLM + E2E 测试（~$4/次），提交前运行
bun run build        # 生成文档 + 编译二进制
bun run dev <cmd>    # 开发模式运行 CLI
bun run gen:skill-docs     # 从 .tmpl 重新生成 SKILL.md
bun run skill:check        # 技能健康面板
bun run slop:diff           # slop-scan（AI 代码质量）当前分支 vs base
```

**SKILL.md 工作流：**
- SKILL.md 文件是**从 `.tmpl` 模板生成的**，不要直接编辑 `.md` 文件
- 编辑 `.tmpl` 文件 → `bun run gen:skill-docs` → 同时提交 `.tmpl` 和生成的 `.md`
- 在 SKILL.md 上发生合并冲突：解决 `.tmpl` 模板的冲突，然后重新生成

**编译后的二进制文件** (`browse/dist/`, `design/dist/`) **永远不要提交**。
这些是 Bun 编译产物（~98MB Mach-O arm64），只在本地使用。

**平台无关设计：** 技能永远不要硬编码框架特定的命令、文件模式或目录结构。
而是读取项目的 CLAUDE.md 获取项目特定配置，缺失时用 AskUserQuestion 询问。

---

## 项目结构

```
gstack/
├── browse/           # Headless 浏览器 CLI (Playwright + Chromium)
│   ├── src/          # CLI + server + ~70 个浏览器命令
│   ├── test/         # 集成测试
│   └── dist/         # 编译后的二进制文件 (browse, find-browse)
├── hosts/            # 各 AI agent 的配置 (claude, codex, factory, kiro, opencode 等)
├── scripts/          # 构建 + 开发工具
│   ├── gen-skill-docs.ts   # 模板 → SKILL.md 生成器
│   ├── resolvers/          # 模板解析模块（preamble, design, review, gbrain 等）
│   └── host-config.ts      # HostConfig 接口 + 验证器
├── extension/        # Chrome 扩展（侧边栏 + 活动源 + CSS inspector + 终端 PTY）
├── design/           # 设计 CLI (GPT Image API 集成)
├── make-pdf/         # PDF 生成 (Playwright + smartypants)
├── cso/              # /cso 技能 (OWASP + STRIDE 安全审计，17 个误报排除)
├── supabase/         # 后端服务（遥测、社区、更新检查）
├── test/             # 技能验证 + eval + E2E 测试
├── docs/             # 设计文档 + 开发指南
│   └── designs/      # 关键技术决策文档
└── package.json      # 构建脚本
```
