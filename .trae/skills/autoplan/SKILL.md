---
name: autoplan
description: 自动评审管线。CEO-设计-工程，自动检测适用评审。Invoke to run all applicable plan reviews.
---

# 自动评审管线

一条命令跑完完整评审。

## 工作流程

1. 检测适用性 - 基于 .gstack/design.md
2. 自动执行 - UI/UX 适用设计评审，API 适用 DevEx 评审
3. 品味决策 - AskUserQuestion 确认
4. 自动编码 - 非品味部分自动完成
