---
name: careful
description: 安全护栏。在执行破坏性命令前警告：rm -rf、DROP TABLE、force-push、git reset --hard。Invoke before destructive operations.
---

# 安全护栏

在执行破坏性命令前发出警告。

## 拦截的操作
- rm -rf: 强制递归删除
- DROP TABLE/DATABASE: 数据库删除
- git push --force: 强制推送
- git reset --hard: 硬重置
- DELETE FROM: 批量数据删除
- TRUNCATE: 表清空
- git clean -fdx: 清理未跟踪文件

## 工作方式
1. 在执行前暂停
2. 确认操作和影响
3. 用户确认后执行
4. 说 be careful 激活

## 覆盖
用户可覆盖任何警告。
