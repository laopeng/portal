---
name: benchmark-models
description: 跨模型基准测试。对 Claude、GPT、Gemini 运行相同 prompt，比较延迟、token、成本及质量。Invoke to benchmark different AI models.
---

# 跨模型基准测试

对比不同 AI 模型在相同任务上的表现。

## 测试维度
1. 延迟 - 首次响应、完整响应、流式速度
2. Token - 输入/输出/总 token
3. 成本 - 输入/输出/总成本
4. 质量 - 准确性、完整性、代码质量

## 测试模型
Claude、GPT、Gemini

## 输出
生成 .gstack/model-benchmark.md
