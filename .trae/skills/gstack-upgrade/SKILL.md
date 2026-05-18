---
name: gstack-upgrade
description: 自升级 gstack。检测全局 vs vendored 安装，同步两者，显示变更内容。Invoke to upgrade gstack.
---

# GStack 自升级

将 gstack 升级到最新版本。

## 工作流程
1. 检测安装类型 - 全局/vendored
2. 获取最新版本 - 检查 npm registry
3. 同步升级 - 全局和项目
4. 验证 - 版本号、冒烟测试

## 使用方式
/gstack-upgrade
