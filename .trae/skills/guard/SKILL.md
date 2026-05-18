---
name: guard
description: 完整安全。careful + freeze 合并。为生产环境提供最大安全。Invoke for maximum safety during production work.
---

# 完整安全

结合 careful 和 freeze，提供最大级别的安全保护。

## 功能 = Careful + Freeze
- careful: 拦截破坏性命令
- freeze: 限制编辑到指定目录

## 使用场景
- 生产环境操作
- 关键基础设施修改
- 高度敏感的配置变更

## 激活
/guard <directory>

## 解除
/unfreeze 同时解除两者。
