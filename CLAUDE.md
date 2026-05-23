# 全局规则

## 关于 gstack

本项目使用 **gstack** —— AI 辅助开发工具栈，将 Claude Code 变成虚拟工程团队（CEO、设计、工程经理、代码审查、QA、安全审计、发布管理），全部通过 slash command 调用。

核心哲学：
- **Boil the Lake**：AI 让完整性的边际成本接近零，永远做完整版
- **Search Before Building**：先搜后造，不重复发明轮子
- **User Sovereignty**：AI 建议，用户决定

---

## sprint 全流程

**Think → Plan → Build → Review → Test → Ship → Reflect**

每个阶段输出自动传递到下一阶段。

### 阶段一：Think（思考产品）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/office-hours` | **YC 创业办公室** | 所有事情的起点。6 个强迫性问题重新框定你的产品，创建实现方案，生成设计文档。 |
| `/learn` | **记忆管理** | 查看、搜索、清理 gstack 跨会话学到的经验。 |

### 阶段二：Plan（制定方案）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/plan-ceo-review` | **CEO / 创始人** | 重新思考问题。4 种模式：Expand / Selective Expand / Hold Scope / Reduce。 |
| `/plan-eng-review` | **工程经理** | 锁定架构、数据流、边界情况和测试计划。输出测试矩阵、失败模式和安全隐患。 |
| `/plan-design-review` | **高级设计师** | 交互式设计评审，0-10 评分，AI Slop 检测。 |
| `/plan-devex-review` | **开发体验负责人** | DX 评审：TTHW、魔法时刻、摩擦点追踪。3 种模式：DX EXPANSION / DX POLISH / DX TRIAGE。 |
| `/autoplan` | **自动评审管线** | 一条命令跑完：CEO → 设计 → 工程，自动检测适用项。 |

### 阶段三：Build（构建）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/design-consultation` | **设计合伙人** | 从零构建设计系统，生成逼真的产品 mockup。 |
| `/design-shotgun` | **设计探索者** | 生成 4-6 个 AI mockup 变体，浏览器中比较，迭代直到满意。 |
| `/design-html` | **设计工程师** | 将 mockup 转为生产级 HTML/CSS。检测 React/Svelte/Vue 框架，智能 API 路由。 |
| `/pair-agent` | **多 AI 协作** | 与任何 AI agent 共享浏览器，每个 agent 独立 tab。 |

### 阶段四：Review（代码评审）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/review` | **高级工程师** | 找生产环境 bug，自动修复。检查清单：并发安全、SQL 注入、错误处理、边界测试。 |
| `/design-review` | **会写代码的设计师** | 实时站点视觉审计 + 修复循环，80 项审计清单。 |
| `/devex-review` | **DX 测试者** | 实时开发者体验审计，对比 plan 分数验证计划匹配现实。 |
| `/codex` | **第二意见** | OpenAI Codex CLI 独立评审：review / adversarial / consult 三种模式。 |
| `/cso` | **首席安全官** | OWASP Top 10 + STRIDE 威胁建模，17 个误报排除。 |

**该用哪个评审？**

| 为谁构建 | 方案阶段（写代码前） | 实时审计（发布后） |
|---------|---------------------|--------------------|
| 最终用户（UI、Web App、移动端） | `/plan-design-review` | `/design-review` |
| 开发者（API、CLI、SDK、文档） | `/plan-devex-review` | `/devex-review` |
| 架构（数据流、性能、测试） | `/plan-eng-review` | `/review` |
| 全部 | `/autoplan` | — |

### 阶段五：Test（测试）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/qa` | **QA 负责人** | 测试应用，找 bug，原子提交修复，自动生成回归测试。 |
| `/qa-only` | **QA 报告者** | 只报告不修复，纯 bug 报告。 |
| `/benchmark` | **性能工程师** | 基线页面加载、Core Web Vitals、资源大小，PR 前后对比。 |
| `/benchmark-models` | **跨模型基准测试** | Claude、GPT、Gemini 同 prompt 对比延迟/token/成本/质量。 |

### 阶段六：Ship（发布）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/ship` | **发布工程师** | 同步主分支、运行测试、推送、打开 PR。 |
| `/land-and-deploy` | **发布工程师** | 合并 PR、等待 CI、部署、验证生产。"approved" 到 "verified" 一条命令。 |
| `/canary` | **SRE** | 部署后监控循环：控制台错误、性能回归、页面故障。 |
| `/setup-deploy` | **部署配置器** | `/land-and-deploy` 一次性配置，自动检测平台和部署命令。 |
| `/document-release` | **技术文档编写者** | 更新项目文档匹配发布内容，Diataxis 覆盖地图。 |
| `/document-generate` | **文档作者** | 用 Diataxis 框架从头生成缺失文档。 |

### 阶段七：Reflect（复盘）

| 命令 | 角色 | 功能 |
|------|------|------|
| `/retro` | **工程经理** | 周度回顾：交付连续天数、测试健康趋势、成长机会。 |
| `/investigate` | **调试者** | 系统化根因调试，追踪数据流，3 次失败修复后停止。 |

---

## 安全与约束工具

| 命令 | 功能 |
|------|------|
| `/careful` | 破坏性命令前警告 |
| `/freeze` | 编辑锁定到单一目录 |
| `/guard` | `/careful` + `/freeze` 合并，生产环境最大安全 |
| `/unfreeze` | 移除 `/freeze` 边界 |

---

## 浏览器工具

| 命令 | 功能 |
|------|------|
| `/browse` | 真实 Chromium 浏览器，~70 个命令：导航、截图、点击、填表、JS 执行、网络捕获 |
| `/open-gstack-browser` | 启动 GStack 浏览器，反机器人隐身、自动模型路由 |
| `/setup-browser-cookies` | 从真实浏览器导入 cookies 到 headless 会话 |
| `/scrape <意图>` | 只读页面数据提取，3 条路径（匹配/原型/拒绝） |
| `/skillify` | 将成功的 scrape 原型编码为永久 browser-skill |

---

## 其他工具

| 命令 | 功能 |
|------|------|
| `/setup-gbrain` | GBrain 初始化，5 分钟从零配置项目知识库 |
| `/sync-gbrain` | 重新索引代码库，刷新搜索指南 |
| `/gstack-upgrade` | 自升级 gstack |
| `/context-save` / `/context-restore` | 保存/恢复会话上下文 |
| `/plan-tune` | 微调已有方案 |
| `/make-pdf` | 页面或 Markdown 转生产质量 PDF |
| `/health` | 技能健康面板 |

---

## 建立并行 sprint

1. `/office-hours` — 描述要构建什么
2. `/plan-ceo-review` — 重新思考，找 10 星版本
3. 批准方案，退出 Plan Mode
4. AI 生成代码
5. `/review` — 找 bug、自动修复
6. `/qa https://staging.xxx.com` — 浏览器点击流程，修复 bug
7. `/ship` — 测试、PR、发布

设计管线可并行：
- `/design-consultation` → 设计系统 → DESIGN.md
- `/design-shotgun` → mockup 比较 → 迭代
- `/design-html` → 生产 HTML/CSS