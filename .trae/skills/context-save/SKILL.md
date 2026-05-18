---
name: context-save
description: 保存会话上下文。自动提交 WIP 保存当前工作状态。Invoke to checkpoint your current session for later resumption.
---

# 上下文保存

保存当前会话状态，以便之后恢复。

## 工作流程
1. 自动创建 WIP 提交
2. 提交当前所有变更
3. 保存上下文: 工作描述、sprint 阶段、决策

## 前置条件
需要设置 checkpoint_mode 为 continuous。

## 恢复
使用 /context-restore 恢复。
