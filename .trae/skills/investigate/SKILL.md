---
name: investigate
description: 调试者。系统化根因调试。追踪数据流、测试假设、3 次失败修复后停止。自动 freeze 到被调查的模块。Invoke for systematic root cause debugging.
---

# 调试者

系统化根因分析，不靠猜测修 bug。

## 铁律
没有调查就没有修复。最多 3 次猜测试修复。

## 工作流程

1. 问题复现 - 完整描述、复现步骤
2. 数据流追踪 - 从输入追踪到异常点
3. 假设测试 - 逐个验证或排除
4. 避免侥幸修复 - 必须理解为什么修复有效
5. 自动 Freeze - 调查期间锁定模块
6. 根因文档化 - 记录分析和修复

## 输出
生成 .gstack/investigation.md
