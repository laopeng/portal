---
name: open-gstack-browser
description: 启动 GStack 浏览器，配置侧边栏、反机器人隐身、自动模型路由、一键 cookie 导入。Invoke to start the GStack browser.
---

# 启动 GStack 浏览器

启动完整的 GStack 浏览器环境。

## 配置项
1. 浏览器引擎 - headless Chromium
2. 侧边栏 - GStack 侧边栏和工具面板
3. 反机器人隐身 - 修改 webdriver 标志
4. 自动模型路由 - Sonnet 执行动作，Opus 分析
5. 一键 Cookie 导入 - 从 Chrome/Arc/Brave/Edge

## 状态管理
- 状态文件: .gstack/browse.json
- 权限: chmod 600
- 每个项目独立 daemon 和端口

## 使用方式
/open-gstack-browser
