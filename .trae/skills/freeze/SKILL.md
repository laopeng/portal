---
name: freeze
description: 编辑锁定。限制文件编辑到单一目录，防止调试时意外修改范围外的代码。Invoke to lock editing to a specific directory.
---

# 编辑锁定

将文件编辑限制在单一目录内。

## 工作方式
1. 指定锁定的目录
2. 所有编辑操作限制在该目录内
3. 超出范围的修改被拒绝

## 使用场景
- 调试复杂 bug 时防止误改
- 重构时限制变更范围
- 调查问题时防止随手修复

## 激活
/freeze <directory>

## 解除
/unfreeze 移除锁定。
