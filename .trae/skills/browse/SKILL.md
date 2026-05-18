---
name: browse
description: 给 AI 装上眼睛。真实 Chromium 浏览器，约 70 个命令：导航、截图、快照、点击、填表、JS 执行、网络捕获。Invoke when you need browser automation.
---

# 浏览器交互

给 AI 装上眼睛。使用真实 Chromium 浏览器进行页面交互。

## 核心命令

### 导航
- navigate <url>: 导航到 URL
- back/forward/reload: 后退/前进/刷新

### 页面信息
- snapshot: 获取页面可访问性快照
- screenshot: 截取页面截图
- url/title: 获取 URL/标题

### 交互
- click <uid>: 点击元素
- fill <uid> <value>: 填充输入框
- hover <uid>: 悬停元素
- press <key>: 按键

### 执行
- eval <js>: 执行 JavaScript
- wait <selector>: 等待元素

## 架构
- 首次 ~3s，后续 ~100-200ms
- 独立 daemon 和端口
- Cookies 跨命令保持
- 30min 空闲超时

## 注意
不要使用 mcp__chrome_* 工具，使用  命令替代。
