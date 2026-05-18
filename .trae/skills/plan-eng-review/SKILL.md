---
name: plan-eng-review
description: 工程经理评审。锁定架构、数据流、边界情况和测试计划。Invoke after plan-ceo-review.
---

# 工程经理评审

锁死架构设计，暴露隐藏假设。

## 工作流程

1. 架构设计 - ASCII 图表
2. 接口定义 - API、数据模型
3. 边界情况 - 正常、异常、并发
4. 测试计划 - 单元、集成、E2E
5. 失败模式 - 单点故障、安全漏洞

## 输出
生成 .gstack/engineering-plan.md
