---
name: sync-gbrain
description: 保持 brain 最新。重新索引代码库，刷新 CLAUDE.md 中的搜索指南。Invoke to keep the knowledge base in sync.
---

# 同步 GBrain

保持项目知识库与代码库同步。

## 工作流程
1. 扫描代码变更 - 检测增删改
2. 增量索引 - 只索引变更文件
3. 刷新搜索指南 - 更新 CLAUDE.md
4. 索引健康检查 - 验证覆盖率

## 使用方式
/sync-gbrain
