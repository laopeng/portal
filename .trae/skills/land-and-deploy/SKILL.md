---
name: land-and-deploy
description: 发布工程师。合并 PR、等待 CI 和部署、验证生产健康状态。一条命令搞定。Invoke after PR approval to deploy.
---

# 发布工程师

从 PR 审核通过到生产验证，一条命令搞定。

## 前置条件
PR 已审核通过。

## 工作流程

1. 合并 PR
2. 等待 CI
3. 部署
4. 生产验证 - 健康检查、核心页面、关键功能
5. 监控 - 5 分钟观察

## 回滚准备
记录部署前版本、准备回滚命令。

## 输出
生成 .gstack/deploy-report.md
