---
name: codex
description: 第二意见。OpenAI Codex CLI 独立评审，3 种模式：review/adversarial/consult。Invoke for independent AI code review.
---

# 第二意见

获取 OpenAI Codex 的独立代码评审。

## 三种模式

1. Review 模式 - 通过/失败判定
2. Adversarial 模式 - 尝试破坏代码
3. Consult 模式 - 开放咨询

## 跨模型分析
当 review 和 codex 都完成后：
- 对比两个 AI 的发现
- 识别共识问题
- 标记分歧点
