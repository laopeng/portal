---
name: context-restore
description: 恢复会话状态。读取 WIP 提交重建之前的会话上下文。Invoke to resume a previous session.
---

# 上下文恢复

从之前的 WIP 提交中恢复会话状态。

## 工作原理
1. 读取最近 WIP 提交
2. 恢复代码变更状态
3. 重建对话上下文摘要
4. 提醒当前 sprint 阶段

## 使用方式
/context-restore
